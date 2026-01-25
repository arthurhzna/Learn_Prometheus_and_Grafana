package booking

import (
	"context"
	"errors"
	"time"

	"github.com/imrenagicom/demo-app/course/catalog"
	"github.com/imrenagicom/demo-app/internal/db"
	"github.com/imrenagicom/demo-app/internal/metrics"
	v1 "github.com/imrenagicom/demo-app/pkg/apiclient/course/v1"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
)

const (
	maxReservationAttemptRetry = 5
	maxReleaseAttemptRetry     = 5
)

func NewService(db *sqlx.DB,
	bookingStore *Store,
	catalogStore *catalog.Store,
) *Service {
	return &Service{
		db:           db,
		bookingStore: bookingStore,
		catalogStore: catalogStore,
	}
}

type Service struct {
	db           *sqlx.DB
	bookingStore *Store
	catalogStore *catalog.Store
}

// CreateBooking creates a new booking for the given course and batch and emits BookingCreated event.
func (s Service) CreateBooking(ctx context.Context, req *v1.CreateBookingRequest) (*Booking, error) {
	course, err := s.catalogStore.FindCourseByID(ctx, req.Booking.GetCourse())
	if err != nil {
		return nil, err
	}

	batch, err := s.catalogStore.FindCourseBatchByID(ctx, req.Booking.GetBatch())
	if err != nil {
		return nil, err
	}

	if err = batch.Available(ctx); err != nil {
		return nil, err
	}

	builder := For(course, batch)
	if req.Booking.Customer != nil {
		// TODO validate customer data
		c := req.Booking.Customer
		builder.WithCustomer(c.Name, c.Email, c.PhoneNumber)
	}
	b := builder.Build()

	err = s.bookingStore.CreateBooking(ctx, b)
	if err != nil {
		return nil, err
	}

	metrics.RecordBookingCreated(course.ID.String(), course.Name)

	return b, nil
}

func (s Service) ReserveBooking(ctx context.Context, req *v1.ReserveBookingRequest) (*Booking, error) {
	tx, err := s.db.BeginTxx(ctx, nil)
	if err != nil {
		return nil, err
	}

	booking, err := s.bookingStore.FindBookingByID(ctx, req.GetBooking(), WithFindTx(tx))
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	if err = s.reserveWithRetry(ctx, tx, booking, 0); err != nil {
		tx.Rollback()
		return nil, err
	}

	if err = s.bookingStore.UpdateBookingStatus(ctx, booking, WithUpdateTx(tx)); err != nil {
		tx.Rollback()
		return nil, err
	}

	if err = tx.Commit(); err != nil {
		tx.Rollback()
		return nil, err
	}

	log.Info().
		Float64("price", booking.Price).
		Msg("Booking reserved")

	metrics.RecordBookingReserved(
		booking.Course.ID.String(),
		booking.Course.Name,
		booking.Price,
		booking.Currency,
	)

	return booking, nil
}

func (s Service) reserveWithRetry(ctx context.Context, tx *sqlx.Tx, b *Booking, retryCount int) error {
	if retryCount > maxReservationAttemptRetry {
		return ErrReservationMaxRetryExceeded
	}

	tc, err := s.catalogStore.FindCourseBatchByIDAndCourseID(ctx, b.Batch.ID.String(), b.Course.ID.String(), catalog.WithFindTx(tx))
	if err != nil {
		return err
	}

	if err := b.Reserve(ctx, tc); err != nil {
		return err
	}

	// if rand.Intn(5)+1 == 3 {
	// 	<-time.After(300 * time.Millisecond)
	// }

	err = s.catalogStore.UpdateBatchAvailableSeats(ctx, tc, catalog.WithUpdateTx(tx))
	if err != nil && !errors.Is(err, db.ErrNoRowUpdated) {
		return err
	}
	if errors.Is(err, db.ErrNoRowUpdated) {
		return s.reserveWithRetry(ctx, tx, b, retryCount+1)
	}
	return nil
}

func (s Service) GetBooking(ctx context.Context, req *v1.GetBookingRequest) (*Booking, error) {
	return s.bookingStore.FindBookingByID(ctx, req.GetBooking())
}

func (s Service) ExpireBooking(ctx context.Context, req *v1.ExpireBookingRequest) error {
	tx, err := s.db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}

	b, err := s.bookingStore.FindBookingByID(ctx, req.GetBooking(), WithDisableCache(), WithFindTx(tx))
	if err != nil {
		tx.Rollback()
		return err
	}

	if err = b.Expire(ctx); err != nil {
		tx.Rollback()
		return err
	}

	ctx, _ = context.WithTimeout(ctx, 5*time.Millisecond)
	if err = s.bookingStore.UpdateBookingStatus(ctx, b, WithUpdateTx(tx)); err != nil {
		tx.Rollback()
		return err
	}

	if err = s.releaseBooking(ctx, tx, b, 0); err != nil {
		tx.Rollback()
		return err
	}

	if err = tx.Commit(); err != nil {
		return err
	}

	metrics.RecordBookingExpired(
		b.Course.ID.String(),
		b.Course.Name,
		b.Price,
		b.Currency,
	)

	return nil
}

func (s Service) releaseBooking(ctx context.Context, tx *sqlx.Tx, b *Booking, retryCount int) error {
	if retryCount > maxReleaseAttemptRetry {
		return ErrReleaseMaxRetryExceeded
	}

	batch, err := s.catalogStore.FindCourseBatchByIDAndCourseID(ctx, b.Batch.ID.String(), b.Course.ID.String(), catalog.WithFindTx(tx))
	if err != nil {
		return err
	}

	err = batch.Allocate(ctx, 1)
	if err != nil {
		return err
	}

	err = s.catalogStore.UpdateBatchAvailableSeats(ctx, batch, catalog.WithUpdateTx(tx))
	if err != nil && !errors.Is(err, db.ErrNoRowUpdated) {
		return err
	}
	if errors.Is(err, db.ErrNoRowUpdated) {
		return s.releaseBooking(ctx, tx, b, retryCount+1)
	}
	return nil
}

func (s Service) ListBookings(ctx context.Context, req *v1.ListBookingsRequest) ([]Booking, string, error) {
	return s.bookingStore.FindAllBookings(ctx, WithFindAllInvoiceNumber(req.GetInvoice()))
}

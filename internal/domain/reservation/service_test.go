package reservation

import (
	"errors"
	"testing"
	"time"

	"spotsync/internal/apperror"
	"spotsync/internal/domain/parkingzone"
	"spotsync/internal/domain/reservation/dto"
	"spotsync/internal/domain/user"

	"gorm.io/gorm"
)

type fakeReservationRepository struct {
	reserveSpotFunc        func(userID, zoneID uint, licensePlate string) (*Reservation, error)
	findByIDFunc           func(id uint) (*Reservation, error)
	findByUserIDFunc       func(userID uint) ([]Reservation, error)
	findAllFunc            func() ([]Reservation, error)
	updateFunc             func(reservation *Reservation) error
	lastUpdatedReservation *Reservation
}

func (f *fakeReservationRepository) ReserveSpot(userID, zoneID uint, licensePlate string) (*Reservation, error) {
	if f.reserveSpotFunc != nil {
		return f.reserveSpotFunc(userID, zoneID, licensePlate)
	}
	return nil, nil
}

func (f *fakeReservationRepository) FindByID(id uint) (*Reservation, error) {
	if f.findByIDFunc != nil {
		return f.findByIDFunc(id)
	}
	return nil, gorm.ErrRecordNotFound
}

func (f *fakeReservationRepository) FindByIDWithRelations(id uint) (*Reservation, error) {
	return f.FindByID(id)
}

func (f *fakeReservationRepository) FindByUserID(userID uint) ([]Reservation, error) {
	if f.findByUserIDFunc != nil {
		return f.findByUserIDFunc(userID)
	}
	return nil, nil
}

func (f *fakeReservationRepository) FindAll() ([]Reservation, error) {
	if f.findAllFunc != nil {
		return f.findAllFunc()
	}
	return nil, nil
}

func (f *fakeReservationRepository) Update(reservation *Reservation) error {
	f.lastUpdatedReservation = reservation
	if f.updateFunc != nil {
		return f.updateFunc(reservation)
	}
	return nil
}

func TestReserveMapsReservationResponse(t *testing.T) {
	repo := &fakeReservationRepository{
		reserveSpotFunc: func(userID, zoneID uint, licensePlate string) (*Reservation, error) {
			return &Reservation{
				ID:           9,
				UserID:       userID,
				ZoneID:       zoneID,
				LicensePlate: licensePlate,
				Status:       StatusActive,
				CreatedAt:    time.Date(2026, 6, 20, 15, 30, 0, 0, time.UTC),
				UpdatedAt:    time.Date(2026, 6, 20, 15, 30, 0, 0, time.UTC),
			}, nil
		},
	}

	service := NewService(repo)
	response, err := service.Reserve(1, dto.CreateReservationRequest{ZoneID: 5, LicensePlate: "ABC-1234"})
	if err != nil {
		t.Fatalf("expected reserve success, got error: %v", err)
	}

	if response.ID != 9 || response.UserID != 1 || response.ZoneID != 5 {
		t.Fatalf("unexpected response: %+v", response)
	}
	if response.Status != StatusActive {
		t.Fatalf("expected active status, got %q", response.Status)
	}
}

func TestReserveRejectsMissingUser(t *testing.T) {
	service := NewService(&fakeReservationRepository{})
	_, err := service.Reserve(0, dto.CreateReservationRequest{ZoneID: 5, LicensePlate: "ABC-1234"})
	assertAppErrorStatus(t, err, 401)
}

func TestGetMyReservationsMapsZoneInfo(t *testing.T) {
	repo := &fakeReservationRepository{
		findByUserIDFunc: func(userID uint) ([]Reservation, error) {
			return []Reservation{
				{
					ID:           11,
					LicensePlate: "ABC-1234",
					Status:       StatusActive,
					Zone: parkingzone.ParkingZone{
						ID:   5,
						Name: "Terminal 1 EV Charging",
						Type: "ev_charging",
					},
					CreatedAt: time.Date(2026, 6, 20, 15, 30, 0, 0, time.UTC),
				},
			}, nil
		},
	}

	service := NewService(repo)
	responses, err := service.GetMyReservations(1)
	if err != nil {
		t.Fatalf("expected get my reservations success, got error: %v", err)
	}

	if len(responses) != 1 {
		t.Fatalf("expected 1 reservation, got %d", len(responses))
	}
	if responses[0].Zone.Name != "Terminal 1 EV Charging" {
		t.Fatalf("unexpected zone info: %+v", responses[0].Zone)
	}
}

func TestCancelOwnReservationRejectsForeignReservation(t *testing.T) {
	repo := &fakeReservationRepository{
		findByIDFunc: func(id uint) (*Reservation, error) {
			return &Reservation{ID: id, UserID: 99, Status: StatusActive}, nil
		},
	}

	service := NewService(repo)
	err := service.CancelOwnReservation(1, 7, "driver")
	assertAppErrorStatus(t, err, 403)
}

func TestCancelOwnReservationUpdatesStatus(t *testing.T) {
	repo := &fakeReservationRepository{
		findByIDFunc: func(id uint) (*Reservation, error) {
			return &Reservation{ID: id, UserID: 1, Status: StatusActive}, nil
		},
	}

	service := NewService(repo)
	err := service.CancelOwnReservation(1, 7, "driver")
	if err != nil {
		t.Fatalf("expected cancel success, got error: %v", err)
	}

	if repo.lastUpdatedReservation == nil || repo.lastUpdatedReservation.Status != StatusCancelled {
		t.Fatalf("expected cancelled status update, got %+v", repo.lastUpdatedReservation)
	}
}

func TestGetAllReservationsRequiresAdmin(t *testing.T) {
	service := NewService(&fakeReservationRepository{})
	_, err := service.GetAllReservations("driver")
	assertAppErrorStatus(t, err, 403)
}

func TestGetAllReservationsMapsAdminResponse(t *testing.T) {
	repo := &fakeReservationRepository{
		findAllFunc: func() ([]Reservation, error) {
			return []Reservation{
				{
					ID:           3,
					UserID:       1,
					ZoneID:       5,
					LicensePlate: "ABC-1234",
					Status:       StatusActive,
					User: user.User{
						ID:    1,
						Name:  "John Doe",
						Email: "john.doe@spotsync.com",
						Role:  "driver",
					},
					Zone: parkingzone.ParkingZone{
						ID:   5,
						Name: "Terminal 1 EV Charging",
						Type: "ev_charging",
					},
				},
			}, nil
		},
	}

	service := NewService(repo)
	responses, err := service.GetAllReservations("admin")
	if err != nil {
		t.Fatalf("expected admin list success, got error: %v", err)
	}

	if len(responses) != 1 {
		t.Fatalf("expected 1 reservation, got %d", len(responses))
	}
	if responses[0].User.Email != "john.doe@spotsync.com" {
		t.Fatalf("unexpected user mapping: %+v", responses[0].User)
	}
}

func assertAppErrorStatus(t *testing.T, err error, statusCode int) {
	t.Helper()

	if err == nil {
		t.Fatalf("expected error status %d, got nil", statusCode)
	}

	var appErr *apperror.AppError
	if !errors.As(err, &appErr) {
		t.Fatalf("expected AppError, got %T", err)
	}

	if appErr.StatusCode != statusCode {
		t.Fatalf("expected status %d, got %d", statusCode, appErr.StatusCode)
	}
}

package reservation

import (
	"errors"

	"spotsync/internal/apperror"
	"spotsync/internal/domain/reservation/dto"

	"gorm.io/gorm"
)

type Service interface {
	Reserve(userID uint, req dto.CreateReservationRequest) (*dto.ReservationResponse, error)
	GetMyReservations(userID uint) ([]dto.MyReservationResponse, error)
	CancelOwnReservation(userID, reservationID uint, role string) error
	GetAllReservations(role string) ([]dto.AdminReservationResponse, error)
}

type service struct {
	repository Repository
}

func NewService(repository Repository) Service {
	return &service{repository: repository}
}

func (s *service) Reserve(userID uint, req dto.CreateReservationRequest) (*dto.ReservationResponse, error) {
	if userID == 0 {
		return nil, apperror.Unauthorized("User identity is required")
	}

	reservation, err := s.repository.ReserveSpot(userID, req.ZoneID, req.LicensePlate)
	if err != nil {
		return nil, err
	}

	return toReservationResponse(reservation), nil
}

func (s *service) GetMyReservations(userID uint) ([]dto.MyReservationResponse, error) {
	if userID == 0 {
		return nil, apperror.Unauthorized("User identity is required")
	}

	reservations, err := s.repository.FindByUserID(userID)
	if err != nil {
		return nil, apperror.Internal("Failed to retrieve reservations")
	}

	responses := make([]dto.MyReservationResponse, 0, len(reservations))
	for _, reservation := range reservations {
		responses = append(responses, toMyReservationResponse(&reservation))
	}

	return responses, nil
}

func (s *service) CancelOwnReservation(userID, reservationID uint, role string) error {
	if userID == 0 {
		return apperror.Unauthorized("User identity is required")
	}

	reservation, err := s.repository.FindByID(reservationID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return apperror.NotFound("Reservation not found")
		}

		return apperror.Internal("Failed to retrieve reservation")
	}

	if role != "admin" && reservation.UserID != userID {
		return apperror.Forbidden("You can only cancel your own reservation")
	}

	if reservation.Status == StatusCancelled {
		return apperror.Conflict("Reservation is already cancelled")
	}

	reservation.Status = StatusCancelled
	if err := s.repository.Update(reservation); err != nil {
		return apperror.Internal("Failed to cancel reservation")
	}

	return nil
}

func (s *service) GetAllReservations(role string) ([]dto.AdminReservationResponse, error) {
	if role != "admin" {
		return nil, apperror.Forbidden("Admin access required")
	}

	reservations, err := s.repository.FindAll()
	if err != nil {
		return nil, apperror.Internal("Failed to retrieve reservations")
	}

	responses := make([]dto.AdminReservationResponse, 0, len(reservations))
	for _, reservation := range reservations {
		responses = append(responses, toAdminReservationResponse(&reservation))
	}

	return responses, nil
}

func toReservationResponse(reservation *Reservation) *dto.ReservationResponse {
	return &dto.ReservationResponse{
		ID:           reservation.ID,
		UserID:       reservation.UserID,
		ZoneID:       reservation.ZoneID,
		LicensePlate: reservation.LicensePlate,
		Status:       reservation.Status,
		CreatedAt:    reservation.CreatedAt,
		UpdatedAt:    reservation.UpdatedAt,
	}
}

func toMyReservationResponse(reservation *Reservation) dto.MyReservationResponse {
	return dto.MyReservationResponse{
		ID:           reservation.ID,
		LicensePlate: reservation.LicensePlate,
		Status:       reservation.Status,
		Zone: dto.ReservationZoneInfo{
			ID:   reservation.Zone.ID,
			Name: reservation.Zone.Name,
			Type: reservation.Zone.Type,
		},
		CreatedAt: reservation.CreatedAt,
	}
}

func toAdminReservationResponse(reservation *Reservation) dto.AdminReservationResponse {
	return dto.AdminReservationResponse{
		ID:           reservation.ID,
		UserID:       reservation.UserID,
		ZoneID:       reservation.ZoneID,
		LicensePlate: reservation.LicensePlate,
		Status:       reservation.Status,
		User: dto.ReservationUserInfo{
			ID:    reservation.User.ID,
			Name:  reservation.User.Name,
			Email: reservation.User.Email,
			Role:  reservation.User.Role,
		},
		Zone: dto.ReservationZoneInfo{
			ID:   reservation.Zone.ID,
			Name: reservation.Zone.Name,
			Type: reservation.Zone.Type,
		},
		CreatedAt: reservation.CreatedAt,
		UpdatedAt: reservation.UpdatedAt,
	}
}

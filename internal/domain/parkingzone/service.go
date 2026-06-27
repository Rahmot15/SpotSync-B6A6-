package parkingzone

import (
	"errors"

	"spotsync/internal/apperror"
	"spotsync/internal/domain/parkingzone/dto"

	"gorm.io/gorm"
)

type Service interface {
	Create(req dto.CreateZoneRequest) (*dto.ZoneResponse, error)
	GetAll() ([]dto.ZoneResponse, error)
	GetByID(id uint) (*dto.ZoneResponse, error)
	Update(id uint, req dto.UpdateZoneRequest) (*dto.ZoneResponse, error)
	Delete(id uint) error
}

type service struct {
	repository Repository
}

func NewService(repository Repository) Service {
	return &service{repository: repository}
}

func (s *service) Create(req dto.CreateZoneRequest) (*dto.ZoneResponse, error) {
	zone := &ParkingZone{
		Name:          req.Name,
		Type:          req.Type,
		TotalCapacity: req.TotalCapacity,
		PricePerHour:  req.PricePerHour,
	}

	if err := s.repository.Create(zone); err != nil {
		return nil, apperror.Internal("Failed to create parking zone")
	}

	response := toZoneResponse(zone, 0)
	return &response, nil
}

func (s *service) GetAll() ([]dto.ZoneResponse, error) {
	zones, err := s.repository.FindAll()
	if err != nil {
		return nil, apperror.Internal("Failed to retrieve parking zones")
	}

	activeCounts, err := s.repository.CountActiveReservationsByZone()
	if err != nil {
		return nil, apperror.Internal("Failed to calculate available spots")
	}

	responses := make([]dto.ZoneResponse, 0, len(zones))
	for _, zone := range zones {
		responses = append(responses, toZoneResponse(&zone, activeCounts[zone.ID]))
	}

	return responses, nil
}

func (s *service) GetByID(id uint) (*dto.ZoneResponse, error) {
	zone, err := s.repository.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperror.NotFound("Parking zone not found")
		}

		return nil, apperror.Internal("Failed to retrieve parking zone")
	}

	activeCount, err := s.repository.CountActiveReservations(zone.ID)
	if err != nil {
		return nil, apperror.Internal("Failed to calculate available spots")
	}

	response := toZoneResponse(zone, activeCount)
	return &response, nil
}

func (s *service) Update(id uint, req dto.UpdateZoneRequest) (*dto.ZoneResponse, error) {
	zone, err := s.repository.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperror.NotFound("Parking zone not found")
		}

		return nil, apperror.Internal("Failed to retrieve parking zone")
	}

	if req.Name != "" {
		zone.Name = req.Name
	}
	if req.Type != "" {
		zone.Type = req.Type
	}
	if req.TotalCapacity != 0 {
		zone.TotalCapacity = req.TotalCapacity
	}
	if req.PricePerHour != 0 {
		zone.PricePerHour = req.PricePerHour
	}

	if err := s.repository.Update(zone); err != nil {
		return nil, apperror.Internal("Failed to update parking zone")
	}

	activeCount, err := s.repository.CountActiveReservations(zone.ID)
	if err != nil {
		return nil, apperror.Internal("Failed to calculate available spots")
	}

	response := toZoneResponse(zone, activeCount)
	return &response, nil
}

func (s *service) Delete(id uint) error {
	zone, err := s.repository.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return apperror.NotFound("Parking zone not found")
		}

		return apperror.Internal("Failed to retrieve parking zone")
	}

	activeCount, err := s.repository.CountActiveReservations(zone.ID)
	if err != nil {
		return apperror.Internal("Failed to check active reservations")
	}
	if activeCount > 0 {
		return apperror.Conflict("Cannot delete a zone with active reservations")
	}

	if err := s.repository.Delete(zone); err != nil {
		return apperror.Internal("Failed to delete parking zone")
	}

	return nil
}

func toZoneResponse(zone *ParkingZone, activeReservations int64) dto.ZoneResponse {
	availableSpots := zone.TotalCapacity - int(activeReservations)
	if availableSpots < 0 {
		availableSpots = 0
	}

	return dto.ZoneResponse{
		ID:             zone.ID,
		Name:           zone.Name,
		Type:           zone.Type,
		TotalCapacity:  zone.TotalCapacity,
		AvailableSpots: availableSpots,
		PricePerHour:   zone.PricePerHour,
		CreatedAt:      zone.CreatedAt,
		UpdatedAt:      zone.UpdatedAt,
	}
}

package parkingzone

import (
	"net/http"
	"strconv"

	"spotsync/internal/apperror"
	"spotsync/internal/domain/parkingzone/dto"
	"spotsync/internal/httpresponse"

	"github.com/labstack/echo/v4"
)

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) Create(c echo.Context) error {
	var req dto.CreateZoneRequest
	if err := c.Bind(&req); err != nil {
		return httpresponse.Error(c, http.StatusBadRequest, "Invalid request body", err.Error())
	}

	if err := c.Validate(&req); err != nil {
		return err
	}

	zone, err := h.service.Create(req)
	if err != nil {
		return err
	}

	return httpresponse.Created(c, "Parking zone created successfully", zone)
}

func (h *Handler) GetAll(c echo.Context) error {
	zones, err := h.service.GetAll()
	if err != nil {
		return err
	}

	return httpresponse.OK(c, "Parking zones retrieved successfully", zones)
}

func (h *Handler) GetByID(c echo.Context) error {
	id, err := parseIDParam(c, "id")
	if err != nil {
		return err
	}

	zone, err := h.service.GetByID(id)
	if err != nil {
		return err
	}

	return httpresponse.OK(c, "Parking zone retrieved successfully", zone)
}

func (h *Handler) Update(c echo.Context) error {
	id, err := parseIDParam(c, "id")
	if err != nil {
		return err
	}

	var req dto.UpdateZoneRequest
	if err := c.Bind(&req); err != nil {
		return httpresponse.Error(c, http.StatusBadRequest, "Invalid request body", err.Error())
	}

	if err := c.Validate(&req); err != nil {
		return err
	}

	zone, err := h.service.Update(id, req)
	if err != nil {
		return err
	}

	return httpresponse.OK(c, "Parking zone updated successfully", zone)
}

func (h *Handler) Delete(c echo.Context) error {
	id, err := parseIDParam(c, "id")
	if err != nil {
		return err
	}

	if err := h.service.Delete(id); err != nil {
		return err
	}

	return httpresponse.OK(c, "Parking zone deleted successfully", nil)
}

func parseIDParam(c echo.Context, name string) (uint, error) {
	id, err := strconv.ParseUint(c.Param(name), 10, 64)
	if err != nil || id == 0 {
		return 0, apperror.BadRequest("Invalid ID parameter", nil)
	}

	return uint(id), nil
}

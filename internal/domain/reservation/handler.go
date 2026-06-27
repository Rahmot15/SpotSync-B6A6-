package reservation

import (
	"fmt"
	"net/http"
	"strconv"

	"spotsync/internal/apperror"
	"spotsync/internal/domain/reservation/dto"
	"spotsync/internal/httpresponse"
	"spotsync/internal/middlewares"

	"github.com/labstack/echo/v4"
)

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) Reserve(c echo.Context) error {
	userID, ok := middlewares.GetUserID(c)
	if !ok {
		return httpresponse.Error(c, http.StatusUnauthorized, "Unauthorized", nil)
	}

	var req dto.CreateReservationRequest
	if err := c.Bind(&req); err != nil {
		return httpresponse.Error(c, http.StatusBadRequest, "Invalid request body", err.Error())
	}

	if err := c.Validate(&req); err != nil {
		return err
	}

	response, err := h.service.Reserve(userID, req)
	if err != nil {
		return err
	}

	return httpresponse.Created(c, "Reservation confirmed successfully", response)
}

func (h *Handler) GetMyReservations(c echo.Context) error {
	userID, ok := middlewares.GetUserID(c)
	if !ok {
		return httpresponse.Error(c, http.StatusUnauthorized, "Unauthorized", nil)
	}

	responses, err := h.service.GetMyReservations(userID)
	if err != nil {
		return err
	}

	return httpresponse.OK(c, "My reservations retrieved successfully", responses)
}

func (h *Handler) Cancel(c echo.Context) error {
	userID, ok := middlewares.GetUserID(c)
	if !ok {
		return httpresponse.Error(c, http.StatusUnauthorized, "Unauthorized", nil)
	}

	role, _ := middlewares.GetRole(c)

	id, err := parseReservationID(c.Param("id"))
	if err != nil {
		return err
	}

	if err := h.service.CancelOwnReservation(userID, id, role); err != nil {
		return err
	}

	return httpresponse.OK(c, "Reservation cancelled successfully", nil)
}

func (h *Handler) GetAll(c echo.Context) error {
	role, _ := middlewares.GetRole(c)

	responses, err := h.service.GetAllReservations(role)
	if err != nil {
		return err
	}

	return httpresponse.OK(c, "Reservations retrieved successfully", responses)
}

func parseReservationID(idParam string) (uint, error) {
	id, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil || id == 0 {
		return 0, apperror.BadRequest(fmt.Sprintf("Invalid ID parameter: %s", idParam), nil)
	}

	return uint(id), nil
}

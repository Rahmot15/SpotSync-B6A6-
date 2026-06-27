package user

import (
	"net/http"

	"spotsync/internal/domain/user/dto"
	"spotsync/internal/httpresponse"

	"github.com/labstack/echo/v4"
)

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) Register(c echo.Context) error {
	var req dto.RegisterRequest
	if err := c.Bind(&req); err != nil {
		return httpresponse.Error(c, http.StatusBadRequest, "Invalid request body", err.Error())
	}

	if err := c.Validate(&req); err != nil {
		return err
	}

	user, err := h.service.Register(req)
	if err != nil {
		return err
	}

	return httpresponse.Created(c, "User registered successfully", user)
}

func (h *Handler) Login(c echo.Context) error {
	var req dto.LoginRequest
	if err := c.Bind(&req); err != nil {
		return httpresponse.Error(c, http.StatusBadRequest, "Invalid request body", err.Error())
	}

	if err := c.Validate(&req); err != nil {
		return err
	}

	loginResponse, err := h.service.Login(req)
	if err != nil {
		return err
	}

	return httpresponse.Success(c, http.StatusOK, "Login successful", loginResponse)
}

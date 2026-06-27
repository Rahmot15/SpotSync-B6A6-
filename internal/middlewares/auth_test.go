package middlewares

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"spotsync/internal/apperror"
	"spotsync/internal/auth"

	"github.com/labstack/echo/v4"
)

func TestJWTAuthRejectsMissingToken(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	handler := JWTAuth(auth.NewJWTService("test-secret"))(func(c echo.Context) error {
		return c.NoContent(http.StatusOK)
	})

	err := handler(c)
	assertAppErrorStatus(t, err, http.StatusUnauthorized)
}

func TestJWTAuthRejectsInvalidToken(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Authorization", "Bearer invalid-token")
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	handler := JWTAuth(auth.NewJWTService("test-secret"))(func(c echo.Context) error {
		return c.NoContent(http.StatusOK)
	})

	err := handler(c)
	assertAppErrorStatus(t, err, http.StatusUnauthorized)
}

func TestJWTAuthSetsUserContext(t *testing.T) {
	jwtService := auth.NewJWTService("test-secret")
	token, err := jwtService.GenerateToken(7, "driver")
	if err != nil {
		t.Fatalf("failed to generate token: %v", err)
	}

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	handler := JWTAuth(jwtService)(func(c echo.Context) error {
		userID, ok := GetUserID(c)
		if !ok || userID != 7 {
			t.Fatalf("expected user id 7, got %d", userID)
		}

		role, ok := GetRole(c)
		if !ok || role != "driver" {
			t.Fatalf("expected driver role, got %q", role)
		}

		return c.NoContent(http.StatusOK)
	})

	if err := handler(c); err != nil {
		t.Fatalf("expected auth middleware success, got error: %v", err)
	}
}

func TestRequireAdminRejectsDriver(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.Set(ContextRole, "driver")

	handler := RequireAdmin(func(c echo.Context) error {
		return c.NoContent(http.StatusOK)
	})

	err := handler(c)
	assertAppErrorStatus(t, err, http.StatusForbidden)
}

func TestRequireAdminAllowsAdmin(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.Set(ContextRole, RoleAdmin)

	handler := RequireAdmin(func(c echo.Context) error {
		return c.NoContent(http.StatusOK)
	})

	if err := handler(c); err != nil {
		t.Fatalf("expected admin middleware success, got error: %v", err)
	}
}

func assertAppErrorStatus(t *testing.T, err error, statusCode int) {
	t.Helper()

	if err == nil {
		t.Fatalf("expected error status %d, got nil", statusCode)
	}

	appErr, ok := err.(*apperror.AppError)
	if !ok {
		t.Fatalf("expected AppError, got %T", err)
	}

	if appErr.StatusCode != statusCode {
		t.Fatalf("expected status %d, got %d", statusCode, appErr.StatusCode)
	}
}

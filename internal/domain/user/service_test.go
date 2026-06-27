package user

import (
	"errors"
	"net/http"
	"testing"

	"spotsync/internal/apperror"
	"spotsync/internal/auth"
	"spotsync/internal/domain/user/dto"

	"gorm.io/gorm"
)

type fakeRepository struct {
	users  map[string]*User
	nextID uint
}

func newFakeRepository() *fakeRepository {
	return &fakeRepository{
		users:  make(map[string]*User),
		nextID: 1,
	}
}

func (r *fakeRepository) Create(user *User) error {
	user.ID = r.nextID
	r.nextID++
	copiedUser := *user
	r.users[user.Email] = &copiedUser
	return nil
}

func (r *fakeRepository) FindByEmail(email string) (*User, error) {
	user, ok := r.users[email]
	if !ok {
		return nil, gorm.ErrRecordNotFound
	}

	copiedUser := *user
	return &copiedUser, nil
}

func (r *fakeRepository) ExistsByEmail(email string) (bool, error) {
	_, ok := r.users[email]
	return ok, nil
}

func TestRegisterDefaultsToDriverAndHashesPassword(t *testing.T) {
	repo := newFakeRepository()
	service := NewService(repo, auth.NewJWTService("test-secret"))

	response, err := service.Register(dto.RegisterRequest{
		Name:     "Test Driver",
		Email:    "driver@spotsync.com",
		Password: "securePassword123",
	})
	if err != nil {
		t.Fatalf("expected register success, got error: %v", err)
	}

	if response.Role != roleDriver {
		t.Fatalf("expected role %q, got %q", roleDriver, response.Role)
	}

	savedUser := repo.users["driver@spotsync.com"]
	if savedUser.Password == "securePassword123" {
		t.Fatal("expected password to be hashed")
	}
	if !auth.CheckPassword(savedUser.Password, "securePassword123") {
		t.Fatal("expected hashed password to match original password")
	}
}

func TestRegisterRejectsPublicAdminRole(t *testing.T) {
	repo := newFakeRepository()
	service := NewService(repo, auth.NewJWTService("test-secret"))

	_, err := service.Register(dto.RegisterRequest{
		Name:     "Admin Attempt",
		Email:    "admin@spotsync.com",
		Password: "securePassword123",
		Role:     roleAdmin,
	})
	if err == nil {
		t.Fatal("expected public admin registration to fail")
	}

	var appErr *apperror.AppError
	if !errors.As(err, &appErr) {
		t.Fatalf("expected app error, got %T", err)
	}
	if appErr.StatusCode != http.StatusForbidden {
		t.Fatalf("expected status %d, got %d", http.StatusForbidden, appErr.StatusCode)
	}
}

func TestLoginReturnsJWT(t *testing.T) {
	repo := newFakeRepository()
	jwtService := auth.NewJWTService("test-secret")
	service := NewService(repo, jwtService)

	_, err := service.Register(dto.RegisterRequest{
		Name:     "Test Driver",
		Email:    "driver@spotsync.com",
		Password: "securePassword123",
		Role:     roleDriver,
	})
	if err != nil {
		t.Fatalf("expected register success, got error: %v", err)
	}

	response, err := service.Login(dto.LoginRequest{
		Email:    "driver@spotsync.com",
		Password: "securePassword123",
	})
	if err != nil {
		t.Fatalf("expected login success, got error: %v", err)
	}

	if response.Token == "" {
		t.Fatal("expected JWT token")
	}
	if response.User.Email != "driver@spotsync.com" {
		t.Fatalf("expected logged in user email, got %q", response.User.Email)
	}

	claims, err := jwtService.ValidateToken(response.Token)
	if err != nil {
		t.Fatalf("expected token to validate, got error: %v", err)
	}
	if claims.UserID == 0 || claims.Role != roleDriver {
		t.Fatalf("unexpected token claims: user_id=%d role=%s", claims.UserID, claims.Role)
	}
}

func TestLoginRejectsWrongPassword(t *testing.T) {
	repo := newFakeRepository()
	service := NewService(repo, auth.NewJWTService("test-secret"))

	_, err := service.Register(dto.RegisterRequest{
		Name:     "Test Driver",
		Email:    "driver@spotsync.com",
		Password: "securePassword123",
	})
	if err != nil {
		t.Fatalf("expected register success, got error: %v", err)
	}

	_, err = service.Login(dto.LoginRequest{
		Email:    "driver@spotsync.com",
		Password: "wrongPassword",
	})
	if err == nil {
		t.Fatal("expected wrong password login to fail")
	}

	var appErr *apperror.AppError
	if !errors.As(err, &appErr) {
		t.Fatalf("expected app error, got %T", err)
	}
	if appErr.StatusCode != http.StatusUnauthorized {
		t.Fatalf("expected status %d, got %d", http.StatusUnauthorized, appErr.StatusCode)
	}
}

package service

import (
	"alumni-app/app/mongodb/model"
	"alumni-app/app/mongodb/repository/mock"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

func TestRegister(t *testing.T) {
	mockRepo := repository.NewMockUserRepository()
	service := NewUserService(mockRepo)

	app := fiber.New()

	// REGISTER ROUTE
	app.Post("/register", service.Register)

	tests := []struct {
		name       string
		body       string
		wantStatus int
	}{
		{
			name:       "Valid Register",
			body:       `{"username":"shendy","email":"s@example.com","password":"12345"}`,
			wantStatus: 201,
		},
		{
			name:       "Missing Fields",
			body:       `{"username":"","email":"x","password":""}`,
			wantStatus: 400,
		},
		{
			name:       "Invalid JSON",
			body:       `{invalid json}`,
			wantStatus: 400,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			req := httptest.NewRequest("POST", "/register", strings.NewReader(tt.body))
			req.Header.Set("Content-Type", "application/json")

			resp, _ := app.Test(req, -1)

			if resp.StatusCode != tt.wantStatus {
				t.Errorf("got %d, want %d", resp.StatusCode, tt.wantStatus)
			}
		})
	}
}

func TestLogin(t *testing.T) {
	mockRepo := repository.NewMockUserRepository()
	service := NewUserService(mockRepo)

	app := fiber.New()

	// LOGIN ROUTE
	app.Post("/login", service.Login)

	// insert dummy hashed user
	hashed, _ := bcrypt.GenerateFromPassword([]byte("12345"), bcrypt.DefaultCost)
	user := model.User{
		ID:           primitive.NewObjectID(),
		Username:     "shendy",
		Email:        "s@example.com",
		PasswordHash: string(hashed),
		Role:         "user",
	}
	mockRepo.Create(&user)

	tests := []struct {
		name       string
		body       string
		wantStatus int
	}{
		{
			name:       "Login Success",
			body:       `{"username":"shendy","password":"12345"}`,
			wantStatus: 200,
		},
		{
			name:       "Wrong Password",
			body:       `{"username":"shendy","password":"wrong"}`,
			wantStatus: 401,
		},
		{
			name:       "User Not Found",
			body:       `{"username":"unknown","password":"123"}`,
			wantStatus: 401,
		},
		{
			name:       "Invalid JSON",
			body:       `{not json}`,
			wantStatus: 400,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			req := httptest.NewRequest("POST", "/login", strings.NewReader(tt.body))
			req.Header.Set("Content-Type", "application/json")

			resp, _ := app.Test(req, -1)

			if resp.StatusCode != tt.wantStatus {
				t.Errorf("got %d, want %d", resp.StatusCode, tt.wantStatus)
			}
		})
	}
}

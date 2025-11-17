package service

import (
	"alumni-app/app/mongodb/model"
	"alumni-app/app/mongodb/repository/mock"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestGetAll(t *testing.T) {
	mockRepo := repository.NewMockAlumniRepository()
	service := NewAlumniService(mockRepo)

	app := fiber.New()
	app.Get("/alumni", service.GetAll)

	// seed
	id := primitive.NewObjectID()
	mockRepo.Data[id.Hex()] = model.Alumni{
		ID:     id,
		Nama:   "Farid",
		Email:  "farid@example.com",
		Jurusan: "TI",
	}

	req := httptest.NewRequest("GET", "/alumni?page=1&limit=10", nil)
	resp, _ := app.Test(req)

	if resp.StatusCode != 200 {
		t.Errorf("expected 200, got %d", resp.StatusCode)
	}
}

func TestGetByID(t *testing.T) {
	mockRepo := repository.NewMockAlumniRepository()
	service := NewAlumniService(mockRepo)

	app := fiber.New()
	app.Get("/alumni/:id", service.GetByID)

	// seed
	id := primitive.NewObjectID()
	mockRepo.Data[id.Hex()] = model.Alumni{
		ID:    id,
		Nama:  "Farid",
		Email: "farid@example.com",
	}

	t.Run("Valid ID", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/alumni/"+id.Hex(), nil)
		resp, _ := app.Test(req)

		if resp.StatusCode != 200 {
			t.Errorf("expected 200, got %d", resp.StatusCode)
		}
	})

	t.Run("Invalid Hex", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/alumni/abc123", nil)
		resp, _ := app.Test(req)

		if resp.StatusCode != 400 {
			t.Errorf("expected 400, got %d", resp.StatusCode)
		}
	})

	t.Run("Not Found", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/alumni/"+primitive.NewObjectID().Hex(), nil)
		resp, _ := app.Test(req)

		if resp.StatusCode != 404 {
			t.Errorf("expected 404, got %d", resp.StatusCode)
		}
	})
}

func TestDelete(t *testing.T) {
	mockRepo := repository.NewMockAlumniRepository()
	service := NewAlumniService(mockRepo)

	app := fiber.New()
	app.Delete("/alumni/:id", func(c *fiber.Ctx) error {
		// Simulasi JWT
		c.Locals("role", "admin")
		c.Locals("user_id", primitive.NewObjectID())
		return service.Delete(c)
	})

	// seed
	id := primitive.NewObjectID()
	mockRepo.Data[id.Hex()] = model.Alumni{
		ID:     id,
		UserID: primitive.NewObjectID(), // owner random
	}

	t.Run("Valid Delete", func(t *testing.T) {
		req := httptest.NewRequest("DELETE", "/alumni/"+id.Hex(), nil)
		resp, _ := app.Test(req)

		if resp.StatusCode != 200 {
			t.Errorf("expected 200, got %d", resp.StatusCode)
		}
	})

	t.Run("Invalid Hex", func(t *testing.T) {
		req := httptest.NewRequest("DELETE", "/alumni/invalid", nil)
		resp, _ := app.Test(req)

		if resp.StatusCode != 400 {
			t.Errorf("expected 400, got %d", resp.StatusCode)
		}
	})

	t.Run("Not Found", func(t *testing.T) {
		req := httptest.NewRequest("DELETE", "/alumni/"+primitive.NewObjectID().Hex(), nil)
		resp, _ := app.Test(req)

		if resp.StatusCode != 404 {
			t.Errorf("expected 404, got %d", resp.StatusCode)
		}
	})
}


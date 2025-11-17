package service

import (
	"alumni-app/app/mongodb/model"
	"alumni-app/app/mongodb/repository/mock"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestPekerjaanGetByID(t *testing.T) {
	mockRepo := repository.NewMockPekerjaanRepository()
	alumniRepo := repository.NewMockAlumniRepository()
	service := NewPekerjaanService(mockRepo, alumniRepo)

	app := fiber.New()
	app.Get("/pekerjaan/:id", service.GetByID)

	// Seed data
	id := primitive.NewObjectID()
	mockRepo.Data[id.Hex()] = model.PekerjaanAlumni{
		ID:           id,
		NamaPerusahaan: "PT Test",
	}

	t.Run("Valid ID", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/pekerjaan/"+id.Hex(), nil)
		resp, _ := app.Test(req)

		if resp.StatusCode != 200 {
			t.Errorf("expected 200, got %d", resp.StatusCode)
		}
	})

	t.Run("Invalid Hex", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/pekerjaan/xyz", nil)
		resp, _ := app.Test(req)

		if resp.StatusCode != 400 {
			t.Errorf("expected 400, got %d", resp.StatusCode)
		}
	})

	t.Run("Not Found", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/pekerjaan/"+primitive.NewObjectID().Hex(), nil)
		resp, _ := app.Test(req)

		if resp.StatusCode != 404 {
			t.Errorf("expected 404, got %d", resp.StatusCode)
		}
	})
}

func TestPekerjaanDelete(t *testing.T) {
	mockRepo := repository.NewMockPekerjaanRepository()
	alumniRepo := repository.NewMockAlumniRepository()
	service := NewPekerjaanService(mockRepo, alumniRepo)

	app := fiber.New()
	app.Delete("/pekerjaan/:id", func(c *fiber.Ctx) error {
		c.Locals("role", "admin")
		c.Locals("user_id", primitive.NewObjectID())
		return service.Delete(c)
	})

	// seed
	id := primitive.NewObjectID()
	mockRepo.Data[id.Hex()] = model.PekerjaanAlumni{
		ID:       id,
		AlumniID: primitive.NewObjectID(),
	}

	t.Run("Valid Delete", func(t *testing.T) {
		req := httptest.NewRequest("DELETE", "/pekerjaan/"+id.Hex(), nil)
		resp, _ := app.Test(req)

		if resp.StatusCode != 200 {
			t.Errorf("expected 200, got %d", resp.StatusCode)
		}
	})

	t.Run("Invalid Hex", func(t *testing.T) {
		req := httptest.NewRequest("DELETE", "/pekerjaan/xxx", nil)
		resp, _ := app.Test(req)

		if resp.StatusCode != 400 {
			t.Errorf("expected 400, got %d", resp.StatusCode)
		}
	})

	t.Run("Not Found", func(t *testing.T) {
		req := httptest.NewRequest("DELETE", "/pekerjaan/"+primitive.NewObjectID().Hex(), nil)
		resp, _ := app.Test(req)

		if resp.StatusCode != 404 {
			t.Errorf("expected 404, got %d", resp.StatusCode)
		}
	})
}

func TestPekerjaanGetAll(t *testing.T) {
	mockRepo := repository.NewMockPekerjaanRepository()
	alumniRepo := repository.NewMockAlumniRepository()
	service := NewPekerjaanService(mockRepo, alumniRepo)

	app := fiber.New()
	app.Get("/pekerjaan", service.GetAll)

	// Seed
	id := primitive.NewObjectID()
	mockRepo.Data[id.Hex()] = model.PekerjaanAlumni{
		ID:              id,
		NamaPerusahaan: "PT ABC",
	}

	req := httptest.NewRequest("GET", "/pekerjaan?page=1&limit=10", nil)
	resp, _ := app.Test(req)

	if resp.StatusCode != 200 {
		t.Errorf("expected 200, got %d", resp.StatusCode)
	}
}

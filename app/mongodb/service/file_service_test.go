package service

import (
	"alumni-app/app/mongodb/model"
	"alumni-app/app/mongodb/repository/mock"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestGetAllFiles(t *testing.T) {
	mockRepo := repository.NewMockFileRepository()
	service := NewFileService(mockRepo, nil, "uploads")

	app := fiber.New()
	app.Get("/files", service.GetAllFiles)

	mockRepo.Create(&model.File{
		ID:       primitive.NewObjectID(),
		FileName: "test.jpg",
	})

	req := httptest.NewRequest("GET", "/files", nil)
	resp, _ := app.Test(req, -1)

	if resp.StatusCode != 200 {
		t.Errorf("Expected 200, got %d", resp.StatusCode)
	}
}

func TestGetFileByID(t *testing.T) {
	mockRepo := repository.NewMockFileRepository()
	service := NewFileService(mockRepo, nil, "uploads")

	app := fiber.New()
	app.Get("/files/:id", service.GetFileByID)

	file := model.File{
		ID:       primitive.NewObjectID(),
		FileName: "foto.jpg",
	}
	mockRepo.Create(&file)

	t.Run("Valid ID", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/files/"+file.ID.Hex(), nil)
		resp, _ := app.Test(req, -1)

		if resp.StatusCode != 200 {
			t.Errorf("Expected 200, got %d", resp.StatusCode)
		}
	})

	t.Run("Not Found", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/files/"+primitive.NewObjectID().Hex(), nil)
		resp, _ := app.Test(req, -1)

		if resp.StatusCode != 404 {
			t.Errorf("Expected 404, got %d", resp.StatusCode)
		}
	})
}

func TestDeleteFile(t *testing.T) {
	mockRepo := repository.NewMockFileRepository()
	service := NewFileService(mockRepo, nil, "uploads")

	app := fiber.New()
	app.Delete("/files/:id", service.DeleteFile)

	file := model.File{
		ID:       primitive.NewObjectID(),
		FileName: "foto.jpg",
		FilePath: "uploads/foto.jpg",
	}
	mockRepo.Create(&file)

	t.Run("Valid Delete", func(t *testing.T) {
		req := httptest.NewRequest("DELETE", "/files/"+file.ID.Hex(), nil)
		resp, _ := app.Test(req, -1)

		if resp.StatusCode != 200 {
			t.Errorf("Expected 200, got %d", resp.StatusCode)
		}
	})

	t.Run("Not Found", func(t *testing.T) {
		req := httptest.NewRequest("DELETE", "/files/"+primitive.NewObjectID().Hex(), nil)
		resp, _ := app.Test(req, -1)

		if resp.StatusCode != 404 {
			t.Errorf("Expected 404, got %d", resp.StatusCode)
		}
	})
}

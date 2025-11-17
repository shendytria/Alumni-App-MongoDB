package service

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
	"alumni-app/app/mongodb/model"
	"alumni-app/app/mongodb/repository"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type FileService interface {
	UploadFoto(c *fiber.Ctx) error
	UploadSertifikat(c *fiber.Ctx) error
	GetAllFiles(c *fiber.Ctx) error
	GetFileByID(c *fiber.Ctx) error
	DeleteFile(c *fiber.Ctx) error
}

type fileService struct {
	repo       repository.FileRepository
	alumniRepo repository.AlumniRepositoryInterface
	uploadPath string
}

func NewFileService(
	fileRepo repository.FileRepository,
	alumniRepo repository.AlumniRepositoryInterface,
	uploadPath string,
) FileService {
	return &fileService{
		repo:       fileRepo,
		alumniRepo: alumniRepo,
		uploadPath: uploadPath,
	}
}

// UploadFoto godoc
// @Summary Upload foto alumni (JPG/PNG, max 1MB)
// @Description Hanya admin atau pemilik data yang boleh upload foto ke profil alumni
// @Tags Files
// @Accept multipart/form-data
// @Produce json
// @Consumes multipart/form-data
// @Param alumni_id path string true "ID Alumni"
// @Param file formData file true "File foto (jpg/png, max 1MB)"
// @Success 200 {object} model.File
// @Failure 400 {object} model.ErrorResponse
// @Failure 403 {object} model.ErrorResponse
// @Failure 404 {object} model.ErrorResponse
// @Failure 500 {object} model.ErrorResponse
// @Security BearerAuth
// @Router /files/upload-foto/{alumni_id} [post]
func (s *fileService) UploadFoto(c *fiber.Ctx) error {
	alumniID := c.Params("alumni_id")
	role, _ := c.Locals("role").(string)
	userID, _ := c.Locals("user_id").(primitive.ObjectID)

	// Role check: user hanya bisa upload miliknya sendiri
	if role != "admin" {
		alumnis, err := s.alumniRepo.GetAllByUserID(userID)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "Failed to check ownership"})
		}

		isOwner := false
		for _, a := range alumnis {
			if a.ID.Hex() == alumniID {
				isOwner = true
				break
			}
		}

		if !isOwner {
			return c.Status(403).JSON(fiber.Map{"error": "You can only upload your own photo"})
		}
	}

	fileHeader, err := c.FormFile("file")
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "No file uploaded"})
	}

	if fileHeader.Size > 1*1024*1024 {
		return c.Status(400).JSON(fiber.Map{"error": "Max photo size is 1MB"})
	}

	contentType := fileHeader.Header.Get("Content-Type")
	if contentType != "image/jpeg" && contentType != "image/png" && contentType != "image/jpg" {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid file type, only JPG/PNG allowed"})
	}

	newName := fmt.Sprintf("FOTO_%s_%s%s", alumniID, uuid.New().String(), filepath.Ext(fileHeader.Filename))
	folder := filepath.Join(s.uploadPath, "foto")
	os.MkdirAll(folder, os.ModePerm)
	filePath := filepath.Join(folder, newName)

	if err := c.SaveFile(fileHeader, filePath); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	fileModel := &model.File{
		FileName:     newName,
		OriginalName: fileHeader.Filename,
		FilePath:     filePath,
		FileSize:     fileHeader.Size,
		FileType:     contentType,
		UploadedAt:   time.Now(),
	}

	if err := s.repo.Create(fileModel); err != nil {
		os.Remove(filePath)
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Foto uploaded successfully",
		"data":    fileModel,
	})
}

// UploadSertifikat godoc
// @Summary Upload sertifikat alumni (PDF, max 2MB)
// @Description Hanya admin atau pemilik data yang boleh upload sertifikat
// @Tags Files
// @Accept multipart/form-data
// @Produce json
// @Consumes multipart/form-data
// @Param alumni_id path string true "ID Alumni"
// @Param file formData file true "File sertifikat (PDF, max 2MB)"
// @Success 200 {object} model.File
// @Failure 400 {object} model.ErrorResponse
// @Failure 403 {object} model.ErrorResponse
// @Failure 404 {object} model.ErrorResponse
// @Failure 500 {object} model.ErrorResponse
// @Security BearerAuth
// @Router /files/upload-sertifikat/{alumni_id} [post]
func (s *fileService) UploadSertifikat(c *fiber.Ctx) error {
	alumniID := c.Params("alumni_id")
	role, _ := c.Locals("role").(string)
	userID, _ := c.Locals("user_id").(primitive.ObjectID)

	if role != "admin" {
		alumnis, err := s.alumniRepo.GetAllByUserID(userID)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "Failed to check ownership"})
		}

		isOwner := false
		for _, a := range alumnis {
			if a.ID.Hex() == alumniID {
				isOwner = true
				break
			}
		}

		if !isOwner {
			return c.Status(403).JSON(fiber.Map{"error": "You can only upload your own certificate"})
		}
	}

	fileHeader, err := c.FormFile("file")
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "No file uploaded"})
	}

	if fileHeader.Header.Get("Content-Type") != "application/pdf" {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid file type, only PDF allowed"})
	}
	if fileHeader.Size > 2*1024*1024 {
		return c.Status(400).JSON(fiber.Map{"error": "Max certificate size is 2MB"})
	}

	newName := fmt.Sprintf("SERTIF_%s_%s%s", alumniID, uuid.New().String(), filepath.Ext(fileHeader.Filename))
	folder := filepath.Join(s.uploadPath, "sertifikat")
	os.MkdirAll(folder, os.ModePerm)
	filePath := filepath.Join(folder, newName)

	if err := c.SaveFile(fileHeader, filePath); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	fileModel := &model.File{
		FileName:     newName,
		OriginalName: fileHeader.Filename,
		FilePath:     filePath,
		FileSize:     fileHeader.Size,
		FileType:     "application/pdf",
		UploadedAt:   time.Now(),
	}

	if err := s.repo.Create(fileModel); err != nil {
		os.Remove(filePath)
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"success": true, "message": "Certificate uploaded successfully", "data": fileModel})
}

// GetAllFiles godoc
// @Summary Dapatkan semua file yang diunggah
// @Description Mengambil daftar seluruh file dari database
// @Tags Files
// @Accept json
// @Produce json
// @Success 200 {array} model.File
// @Failure 500 {object} model.ErrorResponse
// @Security BearerAuth
// @Router /files [get]
func (s *fileService) GetAllFiles(c *fiber.Ctx) error {
	files, err := s.repo.FindAll()
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"success": true, "data": files})
}

// GetFileByID godoc
// @Summary Dapatkan file berdasarkan ID
// @Description Mengambil metadata file dari database berdasarkan ID
// @Tags Files
// @Accept json
// @Produce json
// @Param id path string true "File ID"
// @Success 200 {object} model.File
// @Failure 404 {object} model.ErrorResponse
// @Failure 500 {object} model.ErrorResponse
// @Security BearerAuth
// @Router /files/{id} [get]
func (s *fileService) GetFileByID(c *fiber.Ctx) error {
	id := c.Params("id")
	file, err := s.repo.FindByID(id)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "File not found"})
	}
	return c.JSON(fiber.Map{"success": true, "data": file})
}

// DeleteFile godoc
// @Summary Hapus file berdasarkan ID
// @Description Menghapus file dari sistem dan database
// @Tags Files
// @Accept json
// @Produce json
// @Param id path string true "File ID"
// @Success 200 {object} model.FileUploadResponse
// @Failure 404 {object} model.ErrorResponse
// @Failure 500 {object} model.ErrorResponse
// @Security BearerAuth
// @Router /files/{id} [delete]
func (s *fileService) DeleteFile(c *fiber.Ctx) error {
	id := c.Params("id")
	file, err := s.repo.FindByID(id)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "File not found"})
	}

	if err := os.Remove(file.FilePath); err != nil {
		fmt.Println("⚠️ Warning: gagal hapus file fisik:", err)
	}

	if err := s.repo.Delete(id); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"success": true, "message": "File deleted successfully"})
}

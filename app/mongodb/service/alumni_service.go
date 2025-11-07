package service

import (
	"alumni-app/app/mongodb/model"
	"alumni-app/app/mongodb/repository"
	// "fmt"
	// "time"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AlumniService struct {
	repo *repository.AlumniRepository
}

func NewAlumniService(r *repository.AlumniRepository) *AlumniService {
	return &AlumniService{repo: r}
}

// GetAll godoc
// @Summary Dapatkan semua data alumni
// @Description Mengambil daftar semua alumni dari database (dengan pagination, pencarian, dan sorting)
// @Tags Alumni
// @Accept json
// @Produce json
// @Param page query int false "Halaman saat ini"
// @Param limit query int false "Jumlah data per halaman"
// @Param search query string false "Kata kunci pencarian"
// @Param sortBy query string false "Kolom untuk sorting"
// @Param order query string false "Urutan sorting (asc/desc)"
// @Success 200 {object} model.AlumniResponse
// @Failure 500 {object} fiber.Map
// @Router /alumni [get]
func (s *AlumniService) GetAll(c *fiber.Ctx) error {
	page := c.QueryInt("page", 1)
	limit := c.QueryInt("limit", 10)
	search := c.Query("search", "")
	sortBy := c.Query("sortBy", "created_at")
	order := c.Query("order", "asc")

	data, err := s.repo.GetAll(search, sortBy, order, page, limit)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	total, _ := s.repo.Count(search)
	return c.JSON(model.AlumniResponse{
		Data: data,
		Meta: model.MetaInfo{
			Page:   page,
			Limit:  limit,
			Total:  int(total),
			Pages:  (int(total) + limit - 1) / limit,
			SortBy: sortBy,
			Order:  order,
			Search: search,
		},
	})
}

// GetByID godoc
// @Summary Dapatkan alumni berdasarkan ID
// @Description Mengambil data alumni berdasarkan ID
// @Tags Alumni
// @Accept json
// @Produce json
// @Param id path string true "ID Alumni"
// @Success 200 {object} model.Alumni
// @Failure 400 {object} fiber.Map
// @Failure 404 {object} fiber.Map
// @Router /alumni/{id} [get]
func (s *AlumniService) GetByID(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "ID tidak valid"})
	}
	data, err := s.repo.GetByID(id)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Data tidak ditemukan"})
	}
	return c.JSON(data)
}

// Delete godoc
// @Summary Hapus data alumni
// @Description Menghapus (soft delete) data alumni berdasarkan ID, hanya admin atau pemilik data yang diizinkan
// @Tags Alumni
// @Accept json
// @Produce json
// @Param id path string true "ID Alumni"
// @Success 200 {object} fiber.Map
// @Failure 400 {object} fiber.Map
// @Failure 403 {object} fiber.Map
// @Failure 404 {object} fiber.Map
// @Failure 500 {object} fiber.Map
// @Router /alumni/{id} [delete]
func (s *AlumniService) Delete(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "ID tidak valid"})
	}
	
	// ambil role & user_id dari JWT
	role, _ := c.Locals("role").(string)
	userID, _ := c.Locals("user_id").(primitive.ObjectID)

	// ambil data alumni yang mau dihapus
	alumni, err := s.repo.GetByID(id)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Data tidak ditemukan"})
	}

	// kalau bukan admin, pastikan pemiliknya sama
	if role != "admin" && alumni.UserID != userID {
		return c.Status(403).JSON(fiber.Map{"error": "Tidak diizinkan menghapus data alumni milik orang lain"})
	}

	// soft delete
	if err := s.repo.SoftDelete(id); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"success": true, "message": "Data berhasil dihapus"})
}

// func (s *AlumniService) UploadFiles(c *fiber.Ctx) error {
// 	alumniID := c.Params("id")

// 	// ambil file foto
// 	foto, err := c.FormFile("foto")
// 	if err == nil && foto != nil {
// 		filename := fmt.Sprintf("%d_%s", time.Now().Unix(), foto.Filename)
// 		path := fmt.Sprintf("uploads/foto/%s", filename)
// 		if err := c.SaveFile(foto, path); err != nil {
// 			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
// 		}
// 		// simpan ke DB pakai repository milik struct
// 		if err := s.repo.UpdateFieldByHex(alumniID, "foto", path); err != nil {
// 			return c.Status(500).JSON(fiber.Map{"error": "Gagal menyimpan path foto ke database"})
// 		}
// 	}

// 	// ambil file sertifikat
// 	sertif, err := c.FormFile("sertifikat")
// 	if err == nil && sertif != nil {
// 		filename := fmt.Sprintf("%d_%s", time.Now().Unix(), sertif.Filename)
// 		path := fmt.Sprintf("uploads/sertifikat/%s", filename)
// 		if err := c.SaveFile(sertif, path); err != nil {
// 			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
// 		}
// 		// simpan ke DB pakai repository milik struct
// 		if err := s.repo.UpdateFieldByHex(alumniID, "sertifikat_path", path); err != nil {
// 			return c.Status(500).JSON(fiber.Map{"error": "Gagal menyimpan path sertifikat ke database"})
// 		}
// 	}

// 	return c.JSON(fiber.Map{
// 		"success": true,
// 		"message": "File berhasil diunggah",
// 	})
// }


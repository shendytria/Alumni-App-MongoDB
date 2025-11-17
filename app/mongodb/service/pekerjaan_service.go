package service

import (
	"alumni-app/app/mongodb/model"
	"alumni-app/app/mongodb/repository"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type PekerjaanService struct {
	repo       repository.PekerjaanRepositoryInterface
	alumniRepo repository.AlumniRepositoryInterface
}

func NewPekerjaanService(
	r repository.PekerjaanRepositoryInterface,
	a repository.AlumniRepositoryInterface,
) *PekerjaanService {
	return &PekerjaanService{repo: r, alumniRepo: a}
}

// GetAll godoc
// @Summary Dapatkan semua data pekerjaan
// @Description Mengambil daftar semua pekerjaan alumni dengan pagination dan filter
// @Tags Pekerjaan
// @Accept json
// @Produce json
// @Param page query int false "Halaman saat ini"
// @Param limit query int false "Jumlah data per halaman"
// @Param search query string false "Kata kunci pencarian"
// @Param sortBy query string false "Kolom untuk sorting"
// @Param order query string false "Urutan sorting (asc/desc)"
// @Success 200 {object} model.PekerjaanResponse
// @Failure 500 {object} fiber.Map
// @Security BearerAuth
// @Router /pekerjaan [get]
func (s *PekerjaanService) GetAll(c *fiber.Ctx) error {
	page := c.QueryInt("page", 1)
	limit := c.QueryInt("limit", 10)
	search := c.Query("search", "")
	sortBy := c.Query("sortBy", "created_at")
	order := c.Query("order", "asc")

	data, err := s.repo.GetAll(search, sortBy, order, limit, (page-1)*limit)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	total, _ := s.repo.Count(search)

	return c.JSON(model.PekerjaanResponse{
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
// @Summary Dapatkan pekerjaan berdasarkan ID
// @Description Mengambil detail pekerjaan berdasarkan ID
// @Tags Pekerjaan
// @Accept json
// @Produce json
// @Param id path string true "ID Pekerjaan"
// @Success 200 {object} model.PekerjaanAlumni
// @Failure 400 {object} fiber.Map
// @Failure 404 {object} fiber.Map
// @Security BearerAuth
// @Router /pekerjaan/{id} [get]
func (s *PekerjaanService) GetByID(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "ID tidak valid"})
	}

	data, err := s.repo.GetByID(id)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Data pekerjaan tidak ditemukan"})
	}

	return c.JSON(fiber.Map{"success": true, "data": data})
}

func (s *PekerjaanService) GetByAlumniID(c *fiber.Ctx) error {
	alumniStr := c.Params("alumni_id")
	alumniID, err := primitive.ObjectIDFromHex(alumniStr)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Alumni ID tidak valid"})
	}

	data, err := s.repo.GetByAlumniID(alumniID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"success": true, "data": data})
}

// Create godoc
// @Summary Tambah pekerjaan baru
// @Description Menambahkan data pekerjaan baru untuk alumni tertentu
// @Tags Pekerjaan
// @Accept json
// @Produce json
// @Param body body model.CreatePekerjaanRequest true "Data pekerjaan baru"
// @Success 201 {object} model.PekerjaanAlumni
// @Failure 400 {object} fiber.Map
// @Failure 403 {object} fiber.Map
// @Failure 404 {object} fiber.Map
// @Failure 500 {object} fiber.Map
// @Security BearerAuth
// @Router /pekerjaan [post]
func (s *PekerjaanService) Create(c *fiber.Ctx) error {
	var req model.CreatePekerjaanRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Request tidak valid"})
	}

	role, _ := c.Locals("role").(string)
	userID, _ := c.Locals("user_id").(primitive.ObjectID)

	if role != "admin" {
		alumni, err := s.alumniRepo.GetByID(req.AlumniID)
		if err != nil {
			return c.Status(404).JSON(fiber.Map{"error": "Alumni tidak ditemukan"})
		}
		if alumni.UserID != userID {
			return c.Status(403).JSON(fiber.Map{"error": "Tidak diizinkan"})
		}
	}

	if req.NamaPerusahaan == "" || req.PosisiJabatan == "" {
		return c.Status(400).JSON(fiber.Map{"error": "Nama perusahaan dan posisi jabatan wajib diisi"})
	}

	mulai, err := time.Parse("2006-01-02", req.TanggalMulaiKerja)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Format tanggal mulai kerja salah (YYYY-MM-DD)"})
	}

	var selesai *time.Time
	if req.TanggalSelesaiKerja != "" {
		t, err := time.Parse("2006-01-02", req.TanggalSelesaiKerja)
		if err != nil {
			return c.Status(400).JSON(fiber.Map{"error": "Format tanggal selesai salah (YYYY-MM-DD)"})
		}
		selesai = &t
	}

	p := model.PekerjaanAlumni{
		AlumniID:            req.AlumniID,
		NamaPerusahaan:      req.NamaPerusahaan,
		PosisiJabatan:       req.PosisiJabatan,
		BidangIndustri:      req.BidangIndustri,
		LokasiKerja:         req.LokasiKerja,
		GajiRange:           req.GajiRange,
		TanggalMulaiKerja:   mulai,
		TanggalSelesaiKerja: selesai,
		StatusPekerjaan:     req.StatusPekerjaan,
		DeskripsiPekerjaan:  req.DeskripsiPekerjaan,
	}

	if err := s.repo.Create(&p); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Gagal menambah pekerjaan: " + err.Error()})
	}

	return c.Status(201).JSON(fiber.Map{"success": true, "data": p})
}

// Update godoc
// @Summary Perbarui data pekerjaan
// @Description Memperbarui data pekerjaan berdasarkan ID
// @Tags Pekerjaan
// @Accept json
// @Produce json
// @Param id path string true "ID Pekerjaan"
// @Param body body model.UpdatePekerjaanRequest true "Data pekerjaan baru"
// @Success 200 {object} fiber.Map
// @Failure 400 {object} fiber.Map
// @Failure 403 {object} fiber.Map
// @Failure 404 {object} fiber.Map
// @Failure 500 {object} fiber.Map
// @Security BearerAuth
// @Router /pekerjaan/{id} [put]
func (s *PekerjaanService) Update(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "ID tidak valid"})
	}

	role, _ := c.Locals("role").(string)
	userID, _ := c.Locals("user_id").(primitive.ObjectID)

	existing, err := s.repo.GetByID(id)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Data pekerjaan tidak ditemukan"})
	}

	if role != "admin" {
		alumni, err := s.alumniRepo.GetByID(existing.AlumniID)
		if err != nil {
			return c.Status(404).JSON(fiber.Map{"error": "Alumni tidak ditemukan"})
		}
		if alumni.UserID != userID {
			return c.Status(403).JSON(fiber.Map{"error": "Tidak diizinkan"})
		}
	}

	var req model.UpdatePekerjaanRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Request tidak valid"})
	}

	mulai, err := time.Parse("2006-01-02", req.TanggalMulaiKerja)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Tanggal mulai kerja salah"})
	}

	var selesai *time.Time
	if req.TanggalSelesaiKerja != "" {
		t, err := time.Parse("2006-01-02", req.TanggalSelesaiKerja)
		if err != nil {
			return c.Status(400).JSON(fiber.Map{"error": "Tanggal selesai kerja salah"})
		}
		selesai = &t
	}

	p := model.PekerjaanAlumni{
		NamaPerusahaan:      req.NamaPerusahaan,
		PosisiJabatan:       req.PosisiJabatan,
		BidangIndustri:      req.BidangIndustri,
		LokasiKerja:         req.LokasiKerja,
		GajiRange:           req.GajiRange,
		TanggalMulaiKerja:   mulai,
		TanggalSelesaiKerja: selesai,
		StatusPekerjaan:     req.StatusPekerjaan,
		DeskripsiPekerjaan:  req.DeskripsiPekerjaan,
	}

	if err := s.repo.Update(id, &p); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"success": true, "message": "Pekerjaan berhasil diperbarui"})
}

// Delete godoc
// @Summary Hapus data pekerjaan (soft delete)
// @Description Menghapus pekerjaan berdasarkan ID (hanya admin atau pemilik data yang diizinkan)
// @Tags Pekerjaan
// @Accept json
// @Produce json
// @Param id path string true "ID Pekerjaan"
// @Success 200 {object} fiber.Map
// @Failure 400 {object} fiber.Map
// @Failure 403 {object} fiber.Map
// @Failure 404 {object} fiber.Map
// @Failure 500 {object} fiber.Map
// @Security BearerAuth
// @Router /pekerjaan/{id} [delete]
func (s *PekerjaanService) Delete(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "ID tidak valid"})
	}

	role, _ := c.Locals("role").(string)
	userID, _ := c.Locals("user_id").(primitive.ObjectID)

	existing, err := s.repo.GetByID(id)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Data tidak ditemukan"})
	}

	if role != "admin" {
		alumni, err := s.alumniRepo.GetByID(existing.AlumniID)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "Validasi gagal"})
		}
		if alumni.UserID != userID {
			return c.Status(403).JSON(fiber.Map{"error": "Tidak diizinkan"})
		}
	}

	if err := s.repo.SoftDelete(id); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"success": true, "message": "Pekerjaan berhasil dihapus (soft delete)"})
}

func (s *PekerjaanService) Restore(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "ID tidak valid"})
	}

	role, _ := c.Locals("role").(string)
	userID, _ := c.Locals("user_id").(primitive.ObjectID)

	existing, err := s.repo.GetByID(id)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Data pekerjaan tidak ditemukan"})
	}

	if role != "admin" {
		alumni, err := s.alumniRepo.GetByID(existing.AlumniID)
		if err != nil {
			return c.Status(404).JSON(fiber.Map{"error": "Alumni tidak ditemukan"})
		}
		if alumni.UserID != userID {
			return c.Status(403).JSON(fiber.Map{"error": "Tidak diizinkan"})
		}
	}

	if err := s.repo.Restore(id); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"success": true, "message": "Pekerjaan berhasil direstore"})
}

func (s *PekerjaanService) HardDelete(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "ID tidak valid"})
	}

	role, _ := c.Locals("role").(string)
	userID, _ := c.Locals("user_id").(primitive.ObjectID)

	existing, err := s.repo.GetByID(id)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Data pekerjaan tidak ditemukan"})
	}

	if role != "admin" {
		alumni, err := s.alumniRepo.GetByID(existing.AlumniID)
		if err != nil {
			return c.Status(404).JSON(fiber.Map{"error": "Alumni tidak ditemukan"})
		}
		if alumni.UserID != userID {
			return c.Status(403).JSON(fiber.Map{"error": "Tidak diizinkan"})
		}
	}

	if err := s.repo.HardDelete(id); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"success": true, "message": "Pekerjaan berhasil dihapus permanen"})
}

func (s *PekerjaanService) GetTrashed(c *fiber.Ctx) error {
	role, _ := c.Locals("role").(string)

	if role == "admin" {
		data, err := s.repo.GetTrashed()
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}
		return c.JSON(fiber.Map{
			"success": true,
			"count":   len(data),
			"data":    data,
		})
	}

	userID, _ := c.Locals("user_id").(primitive.ObjectID)

	// Ambil semua alumni milik user ini
	alumnis, err := s.alumniRepo.GetAllByUserID(userID)
	if err != nil || len(alumnis) == 0 {
		return c.Status(404).JSON(fiber.Map{"error": "Alumni untuk user ini tidak ditemukan"})
	}

	// Kumpulkan semua ID alumni milik user ini
	var alumniIDs []primitive.ObjectID
	for _, a := range alumnis {
		alumniIDs = append(alumniIDs, a.ID)
	}

	// Ambil semua pekerjaan yang dihapus dari alumni tersebut
	data, err := s.repo.GetTrashedByAlumniIDs(alumniIDs)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"count":   len(data),
		"data":    data,
	})
}

package service

import (
	"database/sql"
	"fmt"
	"strconv"
	"strings"
	"time"

	"alumni-app/app/postgresql/model"
	"alumni-app/app/postgresql/repository"

	"github.com/gofiber/fiber/v2"
)

type PekerjaanService struct {
	repo repository.PekerjaanRepository
}

func NewPekerjaanService(repo *repository.PekerjaanRepository) *PekerjaanService {
	return &PekerjaanService{repo: *repo}
}

var allowedSortBy = map[string]bool{
	"id": true, "nama_perusahaan": true, "posisi_jabatan": true,
	"bidang_industri": true, "created_at": true,
}
var allowedOrder = map[string]bool{"asc": true, "desc": true}

func (s *PekerjaanService) GetAll(c *fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "10"))
	sortBy := c.Query("sortBy", "id")
	order := c.Query("order", "asc")
	search := c.Query("search", "")

	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}
	if !allowedSortBy[sortBy] {
		sortBy = "id"
	}
	order = strings.ToLower(order)
	if !allowedOrder[order] {
		order = "asc"
	}

	data, err := s.repo.GetAll(search, sortBy, order, limit, (page-1)*limit)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	total, err := s.repo.Count(search)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	resp := model.PekerjaanResponse{
		Data: data,
		Meta: model.MetaInfo{
			Page:   page,
			Limit:  limit,
			Total:  total,
			Pages:  (total + limit - 1) / limit,
			SortBy: sortBy,
			Order:  order,
			Search: search,
		},
	}
	return c.JSON(resp)
}

func (s *PekerjaanService) GetByID(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "ID tidak valid"})
	}
	data, err := s.repo.GetByID(id)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Pekerjaan tidak ditemukan"})
	}
	return c.JSON(fiber.Map{"success": true, "data": data})
}

func (s *PekerjaanService) GetByAlumniID(c *fiber.Ctx) error {
	alumniID, err := strconv.Atoi(c.Params("alumni_id"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Alumni ID tidak valid"})
	}
	data, err := s.repo.GetByAlumniID(alumniID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Gagal mengambil data pekerjaan"})
	}
	return c.JSON(fiber.Map{"success": true, "data": data})
}

func (s *PekerjaanService) Create(c *fiber.Ctx) error {
	var req model.CreatePekerjaanRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Request body tidak valid"})
	}

	if req.AlumniID == 0 || req.NamaPerusahaan == "" || req.PosisiJabatan == "" {
		return c.Status(400).JSON(fiber.Map{"error": "Field wajib diisi: alumni_id, nama_perusahaan, posisi_jabatan"})
	}

	mulai, err := time.Parse("2006-01-02", req.TanggalMulaiKerja)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Tanggal mulai kerja wajib diisi dengan format YYYY-MM-DD"})
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

	createdData, err := s.repo.GetByID(p.ID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Gagal mengambil data setelah dibuat"})
	}

	return c.Status(201).JSON(fiber.Map{"success": true, "data": createdData})
}

func (s *PekerjaanService) Update(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "ID tidak valid"})
	}
	var req model.UpdatePekerjaanRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Request body tidak valid"})
	}

	if req.NamaPerusahaan == "" || req.PosisiJabatan == "" {
		return c.Status(400).JSON(fiber.Map{"error": "Nama_perusahaan dan posisi_jabatan wajib diisi"})
	}

	mulai, err := time.Parse("2006-01-02", req.TanggalMulaiKerja)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Format tanggal mulai salah (YYYY-MM-DD)"})
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
		return c.Status(500).JSON(fiber.Map{"error": "Gagal mengupdate pekerjaan: " + err.Error()})
	}

	updatedData, err := s.repo.GetByID(id)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Pekerjaan tidak ditemukan setelah update"})
	}

	return c.JSON(fiber.Map{"success": true, "data": updatedData})
}

func (s *PekerjaanService) Delete(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "ID tidak valid"})
	}

	uid := c.Locals("user_id")
	if uid == nil {
		return c.Status(401).JSON(fiber.Map{"error": "Unauthorized"})
	}
	userID := uid.(int)
	isAdmin := c.Locals("role") == "admin"

	err = s.repo.Delete(id, userID, isAdmin)
	if err != nil {
		if err == sql.ErrNoRows {
			return c.Status(403).JSON(fiber.Map{"error": "Tidak punya akses atau data tidak ditemukan"})
		}
		return c.Status(500).JSON(fiber.Map{"error": "Gagal menghapus pekerjaan: " + err.Error()})
	}

	return c.JSON(fiber.Map{"success": true, "message": "Pekerjaan berhasil dihapus (soft delete)"})
}

func (s *PekerjaanService) GetTrashed(c *fiber.Ctx) error {
	uid := c.Locals("user_id")
	if uid == nil {
		return c.Status(401).JSON(fiber.Map{"error": "Unauthorized"})
	}
	userID := uid.(int)
	isAdmin := c.Locals("role") == "admin"

	fmt.Printf("DEBUG: Memanggil GetTrashed dengan userID: %d, isAdmin: %v\n", userID, isAdmin)

	data, err := s.repo.GetTrashed(userID, isAdmin)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Gagal mengambil data trash: " + err.Error()})
	}
	return c.JSON(fiber.Map{"success": true, "count": len(data), "data": data})
}

func (s *PekerjaanService) Restore(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "ID tidak valid"})
	}

	uid := c.Locals("user_id")
	if uid == nil {
		return c.Status(401).JSON(fiber.Map{"error": "Unauthorized"})
	}
	userID := uid.(int)
	isAdmin := c.Locals("role") == "admin"

	if err := s.repo.Restore(id, userID, isAdmin); err != nil {
		if err == sql.ErrNoRows {
			return c.Status(403).JSON(fiber.Map{"error": "Tidak punya akses atau data tidak ditemukan"})
		}
		return c.Status(500).JSON(fiber.Map{"error": "Gagal merestore data: " + err.Error()})
	}
	return c.JSON(fiber.Map{"success": true, "message": "Pekerjaan berhasil direstore"})
}

func (s *PekerjaanService) HardDelete(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "ID tidak valid"})
	}

	uid := c.Locals("user_id")
	if uid == nil {
		return c.Status(401).JSON(fiber.Map{"error": "Unauthorized"})
	}
	userID := uid.(int)
	isAdmin := c.Locals("role") == "admin"

	if err := s.repo.HardDelete(id, userID, isAdmin); err != nil {
		if err == sql.ErrNoRows {
			return c.Status(403).JSON(fiber.Map{"error": "Tidak punya akses atau data tidak ditemukan"})
		}
		return c.Status(500).JSON(fiber.Map{"error": "Gagal menghapus permanen: " + err.Error()})
	}
	return c.JSON(fiber.Map{"success": true, "message": "Pekerjaan berhasil dihapus permanen"})
}

func (s *PekerjaanService) GetByTahunLulusWithGaji(c *fiber.Ctx) error {
	tahun, err := strconv.Atoi(c.Params("tahun"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Parameter tahun tidak valid"})
	}
	gajiMin := int64(4000000)

	raw, err := s.repo.GetByTahunLulusWithGaji(tahun, gajiMin)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Gagal mengambil data alumni dengan filter gaji"})
	}
	var res []model.AlumniPekerjaanResponse
	for _, r := range raw {
		res = append(res, model.AlumniPekerjaanResponse{
			ID:             r["id"].(int),
			Nama:           r["nama"].(string),
			Jurusan:        r["jurusan"].(string),
			TahunLulus:     r["tahun_lulus"].(int),
			BidangIndustri: r["bidang_industri"].(string),
			NamaPerusahaan: r["nama_perusahaan"].(string),
			PosisiJabatan:  r["posisi_jabatan"].(string),
			GajiRange:      r["gaji_range"].(int64),
		})
	}

	return c.JSON(fiber.Map{"success": true, "count": len(res), "data": res})
}

package service

import (
	"alumni-app/app/postgresql/model"
	"alumni-app/app/postgresql/repository"
	"database/sql"
	"strconv"
	"github.com/gofiber/fiber/v2"
)

type AlumniService struct {
	repo *repository.AlumniRepository
}

func NewAlumniService(r *repository.AlumniRepository) *AlumniService {
	return &AlumniService{repo: r}
}

func (s *AlumniService) GetAll(c *fiber.Ctx) error {
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

    data, err := s.repo.GetAll(search, sortBy, order, page, limit)
    if err != nil {
        return c.Status(500).JSON(fiber.Map{"error": "Gagal mengambil data alumni"})
    }
    
    total, err := s.repo.Count(search)
    if err != nil {
        return c.Status(500).JSON(fiber.Map{"error": "Gagal menghitung total data"})
    }

    resp := model.AlumniResponse{
        Data: data,
        Meta: model.MetaInfo{
            Page:  page,
            Limit: limit,
            Total: total,
            Pages: (total + limit - 1) / limit, 
            SortBy: sortBy,
            Order: order,
            Search: search,
        },
    }
    return c.JSON(resp)
}

func (s *AlumniService) GetByID(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "ID tidak valid"})
	}

	data, err := s.repo.GetByID(id)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Alumni tidak ditemukan"})
	}
	return c.JSON(data)
}

func (s *AlumniService) Create(c *fiber.Ctx) error {
	var req model.CreateAlumniRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Request tidak valid"})
	}

	if req.NIM == "" || req.Nama == "" || req.Jurusan == "" || req.Email == "" {
		return c.Status(400).JSON(fiber.Map{"error": "NIM, nama, jurusan, dan email wajib diisi"})
	}

	a := model.Alumni{
		NIM:        req.NIM,
		Nama:       req.Nama,
		Jurusan:    req.Jurusan,
		Angkatan:   req.Angkatan,
		TahunLulus: req.TahunLulus,
		Email:      req.Email,
		NoTelepon:  req.NoTelepon,
		Alamat:     req.Alamat,
		UserID:     req.UserID,
	}

	if err := s.repo.Create(&a); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(201).JSON(fiber.Map{"success": true, "data": a})
}

func (s *AlumniService) Update(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "ID tidak valid"})
	}
	
	var req model.UpdateAlumniRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Request tidak valid"})
	}

	if req.Nama == "" || req.Jurusan == "" || req.Email == "" {
		return c.Status(400).JSON(fiber.Map{"error": "Nama, jurusan, dan email wajib diisi"})
	}
	
	a := model.Alumni{
		Nama:       req.Nama,
		Jurusan:    req.Jurusan,
		Angkatan:   req.Angkatan,
		TahunLulus: req.TahunLulus,
		Email:      req.Email,
		NoTelepon:  req.NoTelepon,
		Alamat:     req.Alamat,
	}

	if err := s.repo.Update(id, &a); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	
	updatedData, err := s.repo.GetByID(id)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Gagal mengambil data terbaru setelah update"})
	}

	return c.JSON(fiber.Map{"success": true, "data": updatedData})
}

func (s *AlumniService) Delete(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "ID tidak valid"})
	}
	
	isAdmin := c.Locals("role") == "admin"

	if !isAdmin {
		return c.Status(403).JSON(fiber.Map{"error": "Hanya admin yang dapat menghapus data alumni"})
	}

	if err := s.repo.DeleteByID(id); err != nil {
		if err == sql.ErrNoRows {
			return c.Status(404).JSON(fiber.Map{"error": "Data alumni tidak ditemukan"})
		}
		return c.Status(500).JSON(fiber.Map{"error": "Gagal menghapus data"})
	}

	return c.JSON(fiber.Map{"success": true, "message": "Alumni berhasil dihapus"})
}
package service

import (
	"alumni-app/app/mongodb/model"
	"alumni-app/app/mongodb/repository"
	"alumni-app/utils/mongodb"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	repo *repository.UserRepository
}

func NewUserService(r *repository.UserRepository) *UserService {
	return &UserService{repo: r}
}

// Register godoc
// @Summary Registrasi user baru
// @Description Membuat akun user baru dengan username, email, dan password terenkripsi
// @Tags Users
// @Accept json
// @Produce json
// @Param body body model.RegisterRequest true "Data user baru"
// @Success 201 {object} model.User
// @Failure 400 {object} fiber.Map
// @Failure 500 {object} fiber.Map
// @Router /users/register [post]
func (s *UserService) Register(c *fiber.Ctx) error {
	var req model.RegisterRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Request tidak valid"})
	}
	if req.Username == "" || req.Email == "" || req.PasswordHash == "" {
		return c.Status(400).JSON(fiber.Map{"error": "Username, email, dan password wajib diisi"})
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(req.PasswordHash), bcrypt.DefaultCost)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Gagal enkripsi password"})
	}

	role := req.Role
	if role == "" {
		role = "user"
	}

	user := model.User{
		ID:           primitive.NewObjectID(),
		Username:     req.Username,
		Email:        req.Email,
		PasswordHash: string(hashed),
		Role:         role,
		CreatedAt:    time.Now(),
	}

	if err := s.repo.Create(&user); err != nil {
	return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(201).JSON(fiber.Map{"success": true, "data": user, "message": "User berhasil didaftarkan"})
}

// Login godoc
// @Summary Login user
// @Description Login dengan username dan password untuk mendapatkan token JWT
// @Tags Users
// @Accept json
// @Produce json
// @Param body body model.LoginRequest true "Data login user"
// @Success 200 {object} fiber.Map
// @Failure 400 {object} fiber.Map
// @Failure 401 {object} fiber.Map
// @Failure 500 {object} fiber.Map
// @Router /users/login [post]
func (s *UserService) Login(c *fiber.Ctx) error {
	var req model.LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Request tidak valid"})
	}

	user, err := s.repo.GetByUsername(req.Username)
	if err != nil {
		return c.Status(401).JSON(fiber.Map{"error": "Username tidak ditemukan"})
	}

	if !utils.CheckPassword(user.PasswordHash, req.PasswordHash) {
		return c.Status(401).JSON(fiber.Map{"error": "Password salah"})
	}

	token, err := utils.GenerateToken(user.ID.Hex(), user.Username, user.Role)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Gagal membuat token"})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Login berhasil",
		"data": fiber.Map{
			"user":  user,
			"token": token,
		},
	})
}

package service

import (
	"alumni-app/app/postgresql/model"
	"alumni-app/app/postgresql/repository"
	"alumni-app/utils/postgresql"

	"github.com/gofiber/fiber/v2"

	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) Login(c *fiber.Ctx) error {
	var req model.LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Request body tidak valid"})
	}

	user, err := s.repo.GetByUsername(req.Username)
	if err != nil {
		return c.Status(401).JSON(fiber.Map{"error": "User tidak ditemukan"})
	}

	if !utils.CheckPassword(user.PasswordHash, req.PasswordHash) {
		return c.Status(401).JSON(fiber.Map{"error": "Password salah"})
	}

	token, err := utils.GenerateToken(user.ID, user.Username, user.Role)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Gagal membuat token"})
	}

	resp := &model.LoginResponse{
		User:  user,
		Token: token,
	}
	return c.JSON(fiber.Map{"success": true, "data": resp, "message": "Login berhasil"})
}

func (s *UserService) Register(c *fiber.Ctx) error {
	var req model.RegisterRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Request body tidak valid"})
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.PasswordHash), bcrypt.DefaultCost)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Gagal meng-hash password"})
	}

	role := req.Role
	if role == "" {
		role = "user"
	}

	user := &model.User{
		Username:     req.Username,
		Email:        req.Email,
		PasswordHash: string(hashedPassword),
		Role:         role,
	}

	if err = s.repo.Create(user); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(201).JSON(fiber.Map{"success": true, "data": user, "message": "User berhasil didaftarkan"})
}

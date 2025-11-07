//go:build postgres
package route

import (
	"alumni-app/app/postgresql/service"
	"alumni-app/middleware/postgresql"

	"github.com/gofiber/fiber/v2"
)

// ====================== ROUTE REGISTER ======================
func RegisterRoutes(app *fiber.App, alumniService *service.AlumniService, pekerjaanService *service.PekerjaanService, userService *service.UserService) {
	api := app.Group("/api")

	// ====================== AUTH ROUTES ======================
	api.Post("/login", userService.Login)
	api.Post("/register", userService.Register)

	// ====================== ALUMNI ROUTES ======================
	pekerjaan := api.Group("/pekerjaan", middleware.AuthRequired())

	// === CRUD Utama ===
	pekerjaan.Get("/", pekerjaanService.GetAll)
	pekerjaan.Post("/", middleware.AdminOnly(), pekerjaanService.Create)

	pekerjaan.Get("/trash", pekerjaanService.GetTrashed)

	pekerjaan.Get("/:id", pekerjaanService.GetByID)
	pekerjaan.Put("/:id", middleware.AdminOnly(), pekerjaanService.Update)
	pekerjaan.Delete("/:id", pekerjaanService.Delete)

	// === Fitur Pencarian & Filter ===
	pekerjaan.Get("/alumni/:alumni_id", pekerjaanService.GetByAlumniID)
	pekerjaan.Get("/tahun-lulus/:tahun", pekerjaanService.GetByTahunLulusWithGaji)

	// === Fitur Soft Delete (Trash) ===
	pekerjaan.Put("/:id/restore", pekerjaanService.Restore)
	pekerjaan.Delete("/:id/permanent", pekerjaanService.HardDelete)

	// ====================== ROOT ROUTE ======================
	app.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"message": "Alumni Management System API is running 🚀",
		})
	})
}

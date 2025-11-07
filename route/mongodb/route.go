package route

import (
	svc "alumni-app/app/mongodb/service"
	middleware "alumni-app/middleware/mongodb"

	"github.com/gofiber/fiber/v2"
)

// ====================== ROUTE REGISTER ======================
func RegisterRoutes(
	app *fiber.App,
	alumniService *svc.AlumniService,
	pekerjaanService *svc.PekerjaanService,
	userService *svc.UserService,
	fileService svc.FileService,
) {
	api := app.Group("/api/v1")

	// ====================== AUTH ROUTES ======================
	api.Post("/login", userService.Login)
	api.Post("/register", userService.Register)

	// ====================== PEKERJAAN ROUTES ======================
	pekerjaan := api.Group("/pekerjaan", middleware.AuthRequired())

	// === CRUD Utama ===
	pekerjaan.Get("/", pekerjaanService.GetAll)
	pekerjaan.Post("/", pekerjaanService.Create)
	pekerjaan.Get("/trash", pekerjaanService.GetTrashed)
	pekerjaan.Get("/:id", pekerjaanService.GetByID)
	pekerjaan.Put("/:id", pekerjaanService.Update)
	pekerjaan.Delete("/:id", pekerjaanService.Delete)

	// === Fitur Trash (Soft Delete / Restore / Permanent) ===
	pekerjaan.Put("/:id/restore", pekerjaanService.Restore)
	pekerjaan.Delete("/:id/permanent", pekerjaanService.HardDelete)

	// === Fitur Tambahan ===
	pekerjaan.Get("/alumni/:alumni_id", pekerjaanService.GetByAlumniID)

	// ====================== ALUMNI ROUTES ======================
	alumni := api.Group("/alumni", middleware.AuthRequired())
	alumni.Get("/", alumniService.GetAll)
	// alumni.Post("/:id/upload", alumniService.UploadFiles)
	alumni.Get("/:id", alumniService.GetByID)
	alumni.Delete("/:id", alumniService.Delete)

	// ====================== FILE UPLOAD ROUTES ======================
	files := api.Group("/files", middleware.AuthRequired())

	// Upload file (foto / sertifikat)
	files.Post("/upload-foto/:alumni_id", fileService.UploadFoto)
	files.Post("/upload-sertifikat/:alumni_id", fileService.UploadSertifikat)
	files.Get("/", fileService.GetAllFiles)
	files.Get("/:id", fileService.GetFileByID)
	files.Delete("/:id", fileService.DeleteFile)

	// ====================== ROOT ROUTE ======================
	app.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"message": "🚀 Alumni Management System API (MongoDB version) is running",
		})
	})
}

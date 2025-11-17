package main

import (
	repo "alumni-app/app/mongodb/repository"
	svc "alumni-app/app/mongodb/service"
	"alumni-app/config"
	dbmongo "alumni-app/database/mongodb"
	routepkg "alumni-app/route/mongodb"
	fiberSwagger "github.com/swaggo/fiber-swagger"
    _ "alumni-app/docs"
)

// @title Alumni-App API
// @version 1.0
// @description API untuk mengelola data user dengan MongoDB menggunakan Clean Architecture
// @host localhost:3000
// @BasePath /api/v1
// @schemes http
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {
	config.LoadEnv()
	dbmongo.ConnectMongo()

	app := config.NewApp()

	// repositories
	userRepo := repo.NewUserRepository(dbmongo.DB)
	alumniRepo := repo.NewAlumniRepository(dbmongo.DB)
	pekerjaanRepo := repo.NewPekerjaanRepository(dbmongo.DB)
	// file repo needs the DB handle; ensure dbmongo.DB is exported: var DB *mongo.Database
	fileRepo := repo.NewFileRepository(dbmongo.DB)

	// services
	userService := svc.NewUserService(userRepo)
	alumniService := svc.NewAlumniService(alumniRepo)
	pekerjaanService := svc.NewPekerjaanService(pekerjaanRepo, alumniRepo)
	fileService := svc.NewFileService(fileRepo, alumniRepo, "./uploads")

	// static files
	app.Static("/uploads", "./uploads")

	// Swagger endpoint (UI Dokumentasi API)
	app.Get("/swagger/*", fiberSwagger.WrapHandler)

	// routes — PASS fileService here
	routepkg.RegisterRoutes(app, alumniService, pekerjaanService, userService, fileService)

	port := config.GetEnv("PORT", "3000")
	config.StartServer(app, port)
}

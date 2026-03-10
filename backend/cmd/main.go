package main

import (
	"net/http"
	"os"
	"path/filepath"

	"auth/service/internal/config"
	"auth/service/internal/db"
	"auth/service/internal/handler"
	"auth/service/internal/logger"
	"auth/service/internal/repository"
	"auth/service/internal/service"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	logger.Init()
	log := logger.L()
	cfg := config.Load()

	webDir := "/mnt/e/sobes/auth/frontend"
	log.Infof("Web dir: %s", webDir)

	db, err := db.Connect(cfg)
	if err != nil {
		log.Fatal("Failed to DB connect")
	}
	defer db.Close()

	userRepo := repository.NewUserRepo(db)
	_ = service.NewUserService(userRepo)
	authService := service.NewAuthService(userRepo, cfg.JWTSecret)
	authHandler := handler.NewAuthHandler(authService)

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Route("/api/auth", func(r chi.Router) {
		r.Post("/register", authHandler.Register)
		r.Post("/login", authHandler.Login)
		r.Get("/me", authHandler.Me)
	})

	fs := http.FileServer(http.Dir(webDir))
	r.Get("/*", func(w http.ResponseWriter, r *http.Request) {
		path := filepath.Join(webDir, r.URL.Path)
		if _, err := os.Stat(path); os.IsNotExist(err) {
			http.NotFound(w, r)
			return
		}
		fs.ServeHTTP(w, r)
	})

	log.Infof("Starting server on http://localhost:%s\n", cfg.ServerPort)
	if err := http.ListenAndServe(":"+cfg.ServerPort, r); err != nil {
		log.Fatal("Failed to start server")
	}
}

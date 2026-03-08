package main

import (
	"net/http"
	"os"
	"path/filepath"

	"auth/service/internal/config"
	"auth/service/internal/db"
	"auth/service/internal/logger"
)

func main() {
	logger.Init()
	log := logger.L()
	cfg := config.Load()
	curDir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	projectDir := curDir
	if filepath.Base(curDir) == "cmd" {
		projectDir = filepath.Dir(curDir)
	}
	webDir := filepath.Join(projectDir, "../frontend")

	db, err := db.Connect(cfg)
	if err != nil {
		log.Fatal("Failed to DB connect")
	}
	defer db.Close()

	http.Handle("/", http.FileServer(http.Dir(webDir)))
	log.Infof("Starting server on http://localhost:%v\n", cfg.ServerPort)
	if err := http.ListenAndServe(":"+cfg.ServerPort, nil); err != nil {
		log.Fatal("Failed to start server")
	}
}

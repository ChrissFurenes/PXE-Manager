package main

import (
	"log"

	"PXE-Manager/internal/config"
	"PXE-Manager/internal/httpboot"
	"PXE-Manager/internal/storage"
	"PXE-Manager/internal/tftp"
)

func main() {
	cfg, err := config.Load("config.yaml")
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	store, err := storage.New(cfg.Database.Path)
	if err != nil {
		log.Fatalf("failed to open database: %v", err)
	}

	go func() {
		if err := tftp.Start(cfg); err != nil {
			log.Fatalf("tftp server failed: %v", err)
		}
	}()

	go func() {
		if err := httpboot.Start(cfg, store); err != nil {
			log.Fatalf("http server failed: %v", err)
		}
	}()

	log.Println("PXE Manager is running")
	select {}
}

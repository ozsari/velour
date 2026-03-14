package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/ozsari/velour/internal/api"
	"github.com/ozsari/velour/internal/auth"
	"github.com/ozsari/velour/internal/config"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// CLI subcommands
	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "reset-password":
			resetPassword(cfg)
			return
		case "version":
			fmt.Printf("Velour v%s\n", cfg.Version)
			return
		case "help":
			printHelp()
			return
		}
	}

	fmt.Printf(`
 __   __   _
 \ \ / /__| | ___  _   _ _ __
  \ V / _ \ |/ _ \| | | | '__|
   | |  __/ | (_) | |_| | |
   |_|\___|_|\___/ \__,_|_|

 🚀 Velour v%s
 📡 Listening on %s:%d

`, cfg.Version, cfg.Host, cfg.Port)

	server := api.NewServer(cfg)
	if err := server.Start(); err != nil {
		log.Fatalf("Server failed: %v", err)
		os.Exit(1)
	}
}

func resetPassword(cfg *config.Config) {
	if len(os.Args) < 4 {
		fmt.Println("Usage: velour reset-password <username> <new-password>")
		fmt.Println()
		fmt.Println("Example: velour reset-password admin mynewpassword123")
		os.Exit(1)
	}

	username := os.Args[2]
	newPassword := os.Args[3]

	if len(newPassword) < 6 {
		fmt.Println("Error: password must be at least 6 characters")
		os.Exit(1)
	}

	db, err := sql.Open("sqlite3", cfg.DBPath)
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}
	defer db.Close()

	authSvc := auth.New(db, cfg.JWTSecret)
	if err := authSvc.ResetPassword(username, newPassword); err != nil {
		log.Fatalf("Failed to reset password: %v", err)
	}

	fmt.Printf("Password for '%s' has been reset successfully.\n", username)
}

func printHelp() {
	fmt.Println(`Velour - Server Management Panel

Usage:
  velour                         Start the server
  velour reset-password <user> <pass>  Reset a user's password
  velour version                 Show version
  velour help                    Show this help`)
}

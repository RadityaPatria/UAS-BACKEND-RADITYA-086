package main

import (
	"context"
	"fmt"
	"log"

	"UAS-backend/config"
	"UAS-backend/database"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func main() {
	fmt.Println("=== Adding Admin User ===")

	cfg := config.LoadConfig()
	database.ConnectPostgres(cfg)

	ctx := context.Background()

	// Hash password
	hashed, _ := bcrypt.GenerateFromPassword([]byte("admin123"), 12)

	adminID := uuid.New().String()

	_, err := database.DB.Exec(ctx,
		`INSERT INTO users 
			(id, username, email, password_hash, full_name, role_id, is_active, created_at, updated_at)
		 VALUES ($1,$2,$3,$4,$5,$6,true,NOW(),NOW())
		 ON CONFLICT (username) DO NOTHING`,
		adminID,
		"admin",                       // username
		"admin@system.local",          // email
		hashed,                        // password hash
		"Administrator Sistem",        // full name
		"857c2ed0-e16d-4978-b86e-6fdad876b929", // ROLE ADMIN
	)

	if err != nil {
		log.Fatal("❌ Gagal menambah admin:", err)
	}

	fmt.Println("✅ Admin berhasil ditambahkan:")
	fmt.Println("- Username: admin")
	fmt.Println("- Password: admin123")
	fmt.Println("- Role: admin")
}

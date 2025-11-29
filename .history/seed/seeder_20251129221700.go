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
	fmt.Println("=== Running Seeder with Unique Passwords ===")

	cfg := config.LoadConfig()
	database.ConnectPostgres(cfg)

	ctx := context.Background()

	// ================================
	// HASH PASSWORDS (UNIQUE)
	// ================================
	hash := func(p string) []byte {
		h, _ := bcrypt.GenerateFromPassword([]byte(p), 10)
		return h
	}

	// ================================
	// ROLES (sudah ada — pakai yang kamu punya)
	// ================================
	roles := map[string]string{
		"admin":     "857c2ed0-e16d-4978-b86e-6fdad876b929",
		"mahasiswa": "c65ba8ef-8302-430a-9d5e-fa1438fa0eff",
		"dosen":     "8e9d2e54-aee2-4013-a21d-d3016e8c8295",
	}

	// ================================
	// USERS (dengan password unik)
	// ================================
	type UserSeed struct {
		Username  string
		FullName  string
		Email     string
		RoleID    string
		Password  string
	}

	users := []UserSeed{
		{"admin", "Administrator Kampus", "admin@kampus.ac.id", roles["admin"], "Admin123"},

		{"Pak Denis", "Denis Wanto", "Denis@kampus.ac.id", roles["dosen"], "Dosen1"},
		{"Bu Mira", "Mira Alya", "Mira@kampus.ac.id", roles["dosen"], "Dosen2"},

		{"mhs001", "Mahasiswa Satu", "mhs001@kampus.ac.id", roles["mahasiswa"], "Mhs001*2025"},
		{"mhs002", "Mahasiswa Dua", "mhs002@kampus.ac.id", roles["mahasiswa"], "Mhs002*2025"},
		{"mhs003", "Mahasiswa Tiga", "mhs003@kampus.ac.id", roles["mahasiswa"], "Mhs003*2025"},
		{"mhs004", "Mahasiswa Empat", "mhs004@kampus.ac.id", roles["mahasiswa"], "Mhs004*2025"},
		{"mhs005", "Mahasiswa Lima", "mhs005@kampus.ac.id", roles["mahasiswa"], "Mhs005*2025"},
	}

	userIDs := map[string]string{}

	fmt.Println("→ Seeding users...")
	for _, u := range users {
		id := uuid.New().String()
        userIDs[u.Username] = id

		_, err := database.DB.Exec(ctx,
			`INSERT INTO users (id, username, email, password_hash, full_name, role_id, is_active, created_at, updated_at)
			 VALUES ($1,$2,$3,$4,$5,$6,true,NOW(),NOW())
			 ON CONFLICT (username) DO NOTHING`,
			id, u.Username, u.Email, hash(u.Password), u.FullName, u.RoleID,
		)
		if err != nil {
			log.Fatal("User seeding error:", err)
		}
	}

	// ================================
	// LECTURERS
	// ================================
	type LecturerSeed struct {
		Username   string
		LecturerID string
		Department string
	}

	lecturers := []LecturerSeed{
		{"dosenjoko", "DSN001", "Teknik Informatika"},
		{"dosenwati", "DSN002", "Sistem Informasi"},
	}

	lecturerIDs := map[string]string{}

	fmt.Println("→ Seeding lecturers...")
	for _, l := range lecturers {
		id := uuid.New().String()
		lecturerIDs[l.LecturerID] = id

		_, err := database.DB.Exec(ctx,
			`INSERT INTO lecturers (id, user_id, lecturer_id, department, created_at)
			 VALUES ($1,$2,$3,$4,NOW())
			 ON CONFLICT (lecturer_id) DO NOTHING`,
			id, userIDs[l.Username], l.LecturerID, l.Department,
		)
		if err != nil {
			log.Fatal("Lecturer seeding error:", err)
		}
	}

	// ================================
	// STUDENTS
	// ================================
	type StudentSeed struct {
		Username     string
		StudentID    string
		Program      string
		Year         string
		Advisor      string
	}

	students := []StudentSeed{
		{"mhs001", "MHS001", "Teknik Informatika", "2022", "DSN001"},
		{"mhs002", "MHS002", "Teknik Informatika", "2022", "DSN001"},
		{"mhs003", "MHS003", "Sistem Informasi", "2022", "DSN002"},
		{"mhs004", "MHS004", "Sistem Informasi", "2022", "DSN002"},
		{"mhs005", "MHS005", "Sistem Informasi", "2022", "DSN002"},
	}

	fmt.Println("→ Seeding students...")
	for _, s := range students {
		_, err := database.DB.Exec(ctx,
			`INSERT INTO students (id, user_id, student_id, program_study, academic_year, advisor_id, created_at)
			 VALUES ($1,$2,$3,$4,$5,$6,NOW())
			 ON CONFLICT (student_id) DO NOTHING`,
			uuid.New().String(),
			userIDs[s.Username],
			s.StudentID,
			s.Program,
			s.Year,
			lecturerIDs[s.Advisor],
		)
		if err != nil {
			log.Fatal("Student seeding error:", err)
		}
	}

	fmt.Println("🎉 Seeder selesai! Semua data berhasil dibuat.")
}

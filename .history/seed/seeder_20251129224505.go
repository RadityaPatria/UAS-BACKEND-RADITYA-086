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
	fmt.Println("=== Running Custom Seeder ===")

	cfg := config.LoadConfig()
	database.ConnectPostgres(cfg)

	ctx := context.Background()

	// Helper hash
	hash := func(p string) []byte {
		h, _ := bcrypt.GenerateFromPassword([]byte(p), 12)
		return h
	}

	// ROLE IDs kamu (sudah ada di DB)
	roles := map[string]string{
		"admin":     "857c2ed0-e16d-4978-b86e-6fdad876b929",
		"mahasiswa": "c65ba8ef-8302-430a-9d5e-fa1438fa0eff",
		"dosen":     "8e9d2e54-aee2-4013-a21d-d3016e8c8295",
	}

	// ================================
	// USERS (DOSEN + MAHASISWA)
	// ================================
	type UserSeed struct {
		Username string
		FullName string
		Email    string
		Password string
		RoleID   string
	}

	users := []UserSeed{
		// DOSEN
		{"deniswanto", "Denis Wanto", "denis@kampus.ac.id", "DosenDenis#2025", roles["dosen"]},
		{"miralya", "Mira Alya", "mira@kampus.ac.id", "DosenMira#2025", roles["dosen"]},

		// MAHASISWA
		{"naeiman", "Naeiman", "naeiman@kampus.ac.id", "MhsNaeiman*01", roles["mahasiswa"]},
		{"sarah", "Sarah", "sarah@kampus.ac.id", "MhsSarah*02", roles["mahasiswa"]},
		{"asrul", "Asrul", "asrul@kampus.ac.id", "MhsAsrul*03", roles["mahasiswa"]},
		{"fariz", "Fariz", "fariz@kampus.ac.id", "MhsFariz*04", roles["mahasiswa"]},
		{"yusuf", "Yusuf", "yusuf@kampus.ac.id", "MhsYusuf*05", roles["mahasiswa"]},
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
			log.Fatal("error seeding users:", err)
		}
	}

	// ================================
	// DOSEN (lecturers table)
	// ================================
	type LecturerSeed struct {
		Username   string
		LecturerID string
		Dept       string
	}

	lecturers := []LecturerSeed{
		{"deniswanto", "DSN101", "Teknik Informatika"},
		{"miralya", "DSN102", "Sistem Informasi"},
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
			id, userIDs[l.Username], l.LecturerID, l.Dept,
		)
		if err != nil {
			log.Fatal("error seeding lecturers:", err)
		}
	}

	// ================================
	// MAHASISWA (students table)
	// ================================
	type StudentSeed struct {
		Username  string
		StudentID string
		Program   string
		Year      string
		Advisor   string
	}

	students := []StudentSeed{
		{"naeiman", "MHS101", "Teknik Informatika", "2022", "DSN101"},
		{"sarah", "MHS102", "Teknik Informatika", "2022", "DSN101"},
		{"asrul", "MHS103", "Sistem Informasi", "2022", "DSN102"},
		{"fariz", "MHS104", "Sistem Informasi", "2022", "DSN102"},
		{"yusuf", "MHS105", "Sistem Informasi", "2022", "DSN102"},
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
			log.Fatal("error seeding students:", err)
		}
	}

	fmt.Println("🎉 Seeder selesai — semua data berhasil di-generate!")
}

package repositories

import (
	"context"

	"UAS-backend/app/models"
	"UAS-backend/database"

	"github.com/google/uuid"
)

// GetStudentIDByUserID -> mapping user ke student | FR-003
func GetStudentIDByUserID(ctx context.Context, userID string) (string, error) {
	row := database.DB.QueryRow(ctx,
		`SELECT id FROM students WHERE user_id=$1 LIMIT 1`, userID)

	var studentID string
	if err := row.Scan(&studentID); err != nil {
		return "", err
	}
	return studentID, nil
}

// GetStudentByUserID -> ambil data mahasiswa dari user_id | FR-003
func GetStudentByUserID(ctx context.Context, userID uuid.UUID) (*models.Student, error) {
	row := database.DB.QueryRow(ctx, `
		SELECT id, user_id, student_id, program_study,
		       academic_year, advisor_id, created_at
		FROM students WHERE user_id=$1`, userID)

	var s models.Student
	if err := row.Scan(
		&s.ID, &s.UserID, &s.StudentID,
		&s.ProgramStudy, &s.AcademicYear,
		&s.AdvisorID, &s.CreatedAt,
	); err != nil {
		return nil, err
	}
	return &s, nil
}

// GetStudentByID -> ambil mahasiswa berdasarkan student.id | FR-006
func GetStudentByID(ctx context.Context, id uuid.UUID) (*models.Student, error) {
	row := database.DB.QueryRow(ctx, `
		SELECT id, user_id, student_id, program_study,
		       academic_year, advisor_id, created_at
		FROM students WHERE id=$1`, id)

	var s models.Student
	if err := row.Scan(
		&s.ID, &s.UserID, &s.StudentID,
		&s.ProgramStudy, &s.AcademicYear,
		&s.AdvisorID, &s.CreatedAt,
	); err != nil {
		return nil, err
	}
	return &s, nil
}

// GetAdviseesByLecturerID -> list mahasiswa bimbingan | FR-006
func GetAdviseesByLecturerID(ctx context.Context, lecturerID uuid.UUID) ([]string, error) {
	rows, err := database.DB.Query(ctx,
		`SELECT id FROM students WHERE advisor_id=$1`, lecturerID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var ids []string
	for rows.Next() {
		var id uuid.UUID
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		ids = append(ids, id.String())
	}
	return ids, nil
}

// IsStudentAdvisee -> validasi mahasiswa milik dosen wali | FR-006
func IsStudentAdvisee(ctx context.Context, studentID string, lecturerID string) (bool, error) {
	row := database.DB.QueryRow(ctx,
		`SELECT COUNT(*) FROM students WHERE id=$1 AND advisor_id=$2`,
		studentID, lecturerID)

	var count int
	if err := row.Scan(&count); err != nil {
		return false, err
	}
	return count > 0, nil
}

// GetAllStudents -> list semua mahasiswa | FR-009
func GetAllStudents(ctx context.Context) ([]models.Student, error) {
	rows, err := database.DB.Query(ctx, `
		SELECT id, user_id, student_id, program_study,
		       academic_year, advisor_id, created_at
		FROM students`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []models.Student
	for rows.Next() {
		var s models.Student
		if err := rows.Scan(
			&s.ID, &s.UserID, &s.StudentID,
			&s.ProgramStudy, &s.AcademicYear,
			&s.AdvisorID, &s.CreatedAt,
		); err != nil {
			return nil, err
		}
		list = append(list, s)
	}
	return list, nil
}

// CreateStudent -> tambah mahasiswa baru | FR-009
func CreateStudent(ctx context.Context, s *models.Student) error {
	_, err := database.DB.Exec(ctx, `
		INSERT INTO students (id, user_id, student_id, program_study,
			academic_year, advisor_id, created_at)
		VALUES ($1,$2,$3,$4,$5,$6,NOW())`,
		s.ID, s.UserID, s.StudentID,
		s.ProgramStudy, s.AcademicYear, s.AdvisorID,
	)
	return err
}

// UpdateStudent -> update data mahasiswa | FR-009
func UpdateStudent(ctx context.Context, s *models.Student) error {
	_, err := database.DB.Exec(ctx, `
		UPDATE students
		SET student_id=$2, program_study=$3,
		    academic_year=$4, advisor_id=$5
		WHERE id=$1`,
		s.ID, s.StudentID, s.ProgramStudy,
		s.AcademicYear, s.AdvisorID,
	)
	return err
}

// UpdateAdvisor -> ganti dosen wali mahasiswa | FR-009
func UpdateAdvisor(ctx context.Context, studentID uuid.UUID, lecturerID uuid.UUID) error {
	_, err := database.DB.Exec(ctx,
		`UPDATE students SET advisor_id=$2 WHERE id=$1`,
		studentID, lecturerID,
	)
	return err
}

// DeleteStudent -> hapus mahasiswa | FR-009
func DeleteStudent(ctx context.Context, id uuid.UUID) error {
	_, err := database.DB.Exec(ctx,
		`DELETE FROM students WHERE id=$1`, id)
	return err
}

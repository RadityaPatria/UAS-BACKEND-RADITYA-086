package repositories

import (
	"context"

	"UAS-backend/app/models"
	"UAS-backend/database"

	"github.com/google/uuid"
)


// ==========================================================
// Get Student ID (students.id) by USER ID (users.id)
// Required by Achievement Service (mapping user ↔ student)
// ==========================================================
func GetStudentIDByUserID(ctx context.Context, userID string) (string, error) {
	row := database.DB.QueryRow(ctx,
		`SELECT id FROM students WHERE user_id = $1 LIMIT 1`, userID)

	var studentID string
	if err := row.Scan(&studentID); err != nil {
		return "", err
	}
	return studentID, nil
}



// ==========================================================
// Get full student row by USER ID
// ==========================================================
func GetStudentByUserID(ctx context.Context, userID uuid.UUID) (*models.Student, error) {
	row := database.DB.QueryRow(ctx,
		`SELECT id, user_id, student_id, program_study, 
		 academic_year, advisor_id, created_at
		 FROM students WHERE user_id=$1`,
		userID,
	)

	var s models.Student
	err := row.Scan(
		&s.ID, &s.UserID, &s.StudentID,
		&s.ProgramStudy, &s.AcademicYear,
		&s.AdvisorID, &s.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &s, nil
}



// ==========================================================
// Get student by primary ID (students.id)
// ==========================================================
func GetStudentByID(ctx context.Context, id uuid.UUID) (*models.Student, error) {
	row := database.DB.QueryRow(ctx,
		`SELECT id, user_id, student_id, program_study,
		 academic_year, advisor_id, created_at
		 FROM students WHERE id=$1`, id)

	var s models.Student
	err := row.Scan(
		&s.ID, &s.UserID, &s.StudentID,
		&s.ProgramStudy, &s.AcademicYear,
		&s.AdvisorID, &s.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &s, nil
}



// ==========================================================
// Get Advisees (students.id) by Lecturer ID (advisor_id)
// Dipakai di FR-006: List Achievement for Dosen Wali
// ==========================================================
func GetAdviseesByLecturerID(ctx context.Context, lecturerID uuid.UUID) ([]string, error) {
	rows, err := database.DB.Query(ctx,
		`SELECT id FROM students WHERE advisor_id=$1`,
		lecturerID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var ids []string
	for rows.Next() {
		var stID uuid.UUID
		if err := rows.Scan(&stID); err != nil {
			return nil, err
		}
		ids = append(ids, stID.String())
	}

	return ids, nil
}



// ==========================================================
// Create Student (Admin)
// ==========================================================
func CreateStudent(ctx context.Context, s *models.Student) error {
	_, err := database.DB.Exec(ctx,
		`INSERT INTO students (id, user_id, student_id, program_study,
		 academic_year, advisor_id, created_at)
		 VALUES ($1,$2,$3,$4,$5,$6,NOW())`,
		s.ID, s.UserID, s.StudentID, s.ProgramStudy,
		s.AcademicYear, s.AdvisorID,
	)
	return err
}



// ==========================================================
// Update Student (Admin)
// ==========================================================
func UpdateStudent(ctx context.Context, s *models.Student) error {
	_, err := database.DB.Exec(ctx,
		`UPDATE students SET student_id=$2, program_study=$3,
		 academic_year=$4, advisor_id=$5 WHERE id=$1`,
		s.ID, s.StudentID, s.ProgramStudy,
		s.AcademicYear, s.AdvisorID,
	)
	return err
}



// ==========================================================
// Delete Student (Admin)
// ==========================================================
func DeleteStudent(ctx context.Context, id uuid.UUID) error {
	_, err := database.DB.Exec(ctx,
		`DELETE FROM students WHERE id=$1`,
		id,
	)
	return err
}
func IsStudentAdvisee(ctx context.Context, studentID string, lecturerID string) (bool, error) {
	row := database.DB.QueryRow(ctx,
		`SELECT COUNT(*) FROM students 
         WHERE id = $1 AND advisor_id = $2`,
		studentID, lecturerID,
	)

	var count int
	if err := row.Scan(&count); err != nil {
		return false, err
	}

	return count > 0, nil
}

package repositories

import (
	"context"

	"UAS-backend/app/models"
	"UAS-backend/database"

	"github.com/google/uuid"
)

type StudentRepository struct{}

var StudentRepo = &StudentRepository{}

//
// =======================================================
// AMBIL DATA STUDENT BERDASARKAN USER_ID (PENTING UNTUK LOGIN)
// =======================================================
//
func (r *StudentRepository) GetByUserID(ctx context.Context, userID uuid.UUID) (*models.Student, error) {
	row := database.DB.QueryRow(ctx,
		`SELECT id, user_id, student_id, program_study, academic_year, advisor_id, created_at
		 FROM students WHERE user_id = $1`,
		userID,
	)

	var s models.Student
	err := row.Scan(
		&s.ID,
		&s.UserID,
		&s.StudentID,
		&s.ProgramStudy,
		&s.AcademicYear,
		&s.AdvisorID,
		&s.CreatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &s, nil
}

//
// =======================================================
// AMBIL STUDENT BERDASARKAN ID
// =======================================================
//
func (r *StudentRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.Student, error) {
	row := database.DB.QueryRow(ctx,
		`SELECT id, user_id, student_id, program_study, academic_year, advisor_id, created_at
		 FROM students WHERE id = $1`,
		id,
	)

	var s models.Student
	err := row.Scan(
		&s.ID,
		&s.UserID,
		&s.StudentID,
		&s.ProgramStudy,
		&s.AcademicYear,
		&s.AdvisorID,
		&s.CreatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &s, nil
}

//
// =======================================================
// ADMIN: GET SEMUA STUDENT
// =======================================================
//
func (r *StudentRepository) GetAll(ctx context.Context) ([]models.Student, error) {
	rows, err := database.DB.Query(ctx,
		`SELECT id, user_id, student_id, program_study, academic_year, advisor_id, created_at
		 FROM students`,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []models.Student

	for rows.Next() {
		var s models.Student
		err := rows.Scan(
			&s.ID,
			&s.UserID,
			&s.StudentID,
			&s.ProgramStudy,
			&s.AcademicYear,
			&s.AdvisorID,
			&s.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		list = append(list, s)
	}

	return list, nil
}

//
// =======================================================
// DOSEN WALI: AMBIL LIST MAHASISWA BIMBINGAN
// =======================================================
//
func (r *StudentRepository) GetByAdvisor(ctx context.Context, lecturerID uuid.UUID) ([]models.Student, error) {
	rows, err := database.DB.Query(ctx,
		`SELECT id, user_id, student_id, program_study, academic_year, advisor_id, created_at
		 FROM students WHERE advisor_id = $1`,
		lecturerID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []models.Student

	for rows.Next() {
		var s models.Student
		err := rows.Scan(
			&s.ID,
			&s.UserID,
			&s.StudentID,
			&s.ProgramStudy,
			&s.AcademicYear,
			&s.AdvisorID,
			&s.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		result = append(result, s)
	}

	return result, nil
}

//
// =======================================================
// ADMIN: CREATE STUDENT
// =======================================================
//
func (r *StudentRepository) Create(ctx context.Context, s *models.Student) error {
	_, err := database.DB.Exec(ctx,
		`INSERT INTO students (id, user_id, student_id, program_study, academic_year, advisor_id)
		 VALUES ($1, $2, $3, $4, $5, $6)`,
		s.ID,
		s.UserID,
		s.StudentID,
		s.ProgramStudy,
		s.AcademicYear,
		s.AdvisorID,
	)
	return err
}

//
// =======================================================
// ADMIN: UPDATE STUDENT
// =======================================================
//
func (r *StudentRepository) Update(ctx context.Context, s *models.Student) error {
	_, err := database.DB.Exec(ctx,
		`UPDATE students
		 SET student_id=$2, program_study=$3, academic_year=$4, advisor_id=$5
		 WHERE id=$1`,
		s.ID,
		s.StudentID,
		s.ProgramStudy,
		s.AcademicYear,
		s.AdvisorID,
	)
	return err
}

//
// =======================================================
// ADMIN: DELETE STUDENT
// =======================================================
//
func (r *StudentRepository) Delete(ctx context.Context, id uuid.UUID) error {
	_, err := database.DB.Exec(ctx,
		`DELETE FROM students WHERE id=$1`,
		id,
	)
	return err
}

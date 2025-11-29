package repositories

import (
	"context"
	"UAS-backend/app/models"
	"UAS-backend/database"
)


// GetStudentByUserID → ambil profil mahasiswa via user_id
func GetStudentByUserID(ctx context.Context, userID string) (*models.Student, error) {
	row := database.DB.QueryRow(ctx,
		`SELECT id, user_id, student_id, program_study, 
		 academic_year, advisor_id, created_at 
		 FROM students WHERE user_id=$1`, userID)

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

// GetStudentByID → ambil mahasiswa via ID
func GetStudentByID(ctx context.Context, id string) (*models.Student, error) {
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
	return &s, err
}


// GetStudentsByAdvisor → dosen wali ambil mahasiswa bimbingan
func GetStudentsByAdvisor(ctx context.Context, lecturerID string) ([]models.Student, error) {
	rows, err := database.DB.Query(ctx,
		`SELECT id, user_id, student_id, program_study, 
		 academic_year, advisor_id, created_at 
		 FROM students WHERE advisor_id=$1`, lecturerID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	students := []models.Student{}
	for rows.Next() {
		var s models.Student
		rows.Scan(
			&s.ID, &s.UserID, &s.StudentID,
			&s.ProgramStudy, &s.AcademicYear,
			&s.AdvisorID, &s.CreatedAt,
		)
		students = append(students, s)
	}
	return students, nil
}


// CreateStudent → admin menambah mahasiswa
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


// UpdateStudent → admin update mahasiswa
func UpdateStudent(ctx context.Context, s *models.Student) error {
	_, err := database.DB.Exec(ctx,
		`UPDATE students SET student_id=$2, program_study=$3, 
		 academic_year=$4, advisor_id=$5 WHERE id=$1`,
		s.ID, s.StudentID, s.ProgramStudy,
		s.AcademicYear, s.AdvisorID,
	)
	return err
}


// DeleteStudent → admin hapus mahasiswa
func DeleteStudent(ctx context.Context, id string) error {
	_, err := database.DB.Exec(ctx, `DELETE FROM students WHERE id=$1`, id)
	return err
}

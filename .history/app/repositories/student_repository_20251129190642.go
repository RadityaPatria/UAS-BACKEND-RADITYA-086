package repository

import (
	"context"
	"UAS-backend/app/models"
	"UAS-backend/database"
)

func GetStudentByUserID(ctx context.Context, userID string) (*models.Student, error) {
	row := database.DB.QueryRow(ctx,
		`SELECT id, user_id, student_id, program_study, academic_year, advisor_id, created_at
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

func GetStudentByID(ctx context.Context, id string) (*models.Student, error) {
	row := database.DB.QueryRow(ctx,
		`SELECT id, user_id, student_id, program_study, academic_year, advisor_id, created_at
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

func GetStudentsByAdvisor(ctx context.Context, lecturerID string) ([]models.Student, error) {
	rows, err := database.DB.Query(ctx,
		`SELECT id, user_id, student_id, program_study, academic_year, advisor_id, created_at
		 FROM students WHERE advisor_id=$1`, lecturerID)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []models.Student

	for rows.Next() {
		var s models.Student
		rows.Scan(
			&s.ID, &s.UserID, &s.StudentID,
			&s.ProgramStudy, &s.AcademicYear,
			&s.AdvisorID, &s.CreatedAt,
		)
		list = append(list, s)
	}

	return list, nil
}

func CreateStudent(ctx context.Context, s *models.Student) error {
	_, err := database.DB.Exec(ctx,
		`INSERT INTO students (id, user_id, student_id, program_study, academic_year, advisor_id, created_at)
		 VALUES ($1,$2,$3,$4,$5,$6,NOW())`,
		s.ID, s.UserID, s.StudentID, s.ProgramStudy, s.AcademicYear, s.AdvisorID)
	return err
}

func UpdateStudent(ctx context.Context, s *models.Student) error {
	_, err := database.DB.Exec(ctx,
		`UPDATE students SET student_id=$2, program_study=$3, academic_year=$4, advisor_id=$5
		 WHERE id=$1`,
		s.ID, s.StudentID, s.ProgramStudy, s.AcademicYear, s.AdvisorID)
	return err
}

func DeleteStudent(ctx context.Context, id string) error {
	_, err := database.DB.Exec(ctx, `DELETE FROM students WHERE id=$1`, id)
	return err
}

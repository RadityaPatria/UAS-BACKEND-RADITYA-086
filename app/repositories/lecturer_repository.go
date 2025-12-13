package repositories

import (
	"context"

	"UAS-backend/app/models"
	"UAS-backend/database"
)

// GetLecturerByUserID -> ambil profil dosen berdasarkan user_id | FR-004
func GetLecturerByUserID(ctx context.Context, userID string) (*models.Lecturer, error) {
	row := database.DB.QueryRow(ctx,
		`SELECT id, user_id, lecturer_id, department, created_at
		 FROM lecturers WHERE user_id = $1`, userID)

	var l models.Lecturer
	if err := row.Scan(&l.ID, &l.UserID, &l.LecturerID, &l.Department, &l.CreatedAt); err != nil {
		return nil, err
	}
	return &l, nil
}

// GetLecturerByID -> ambil dosen berdasarkan lecturers.id | FR-006
func GetLecturerByID(ctx context.Context, id string) (*models.Lecturer, error) {
	row := database.DB.QueryRow(ctx,
		`SELECT id, user_id, lecturer_id, department, created_at
		 FROM lecturers WHERE id = $1`, id)

	var l models.Lecturer
	if err := row.Scan(&l.ID, &l.UserID, &l.LecturerID, &l.Department, &l.CreatedAt); err != nil {
		return nil, err
	}
	return &l, nil
}

// GetAllLecturers -> list semua dosen | FR-009
func GetAllLecturers(ctx context.Context) ([]models.Lecturer, error) {
	rows, err := database.DB.Query(ctx,
		`SELECT id, user_id, lecturer_id, department, created_at FROM lecturers`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []models.Lecturer
	for rows.Next() {
		var l models.Lecturer
		if err := rows.Scan(&l.ID, &l.UserID, &l.LecturerID, &l.Department, &l.CreatedAt); err != nil {
			return nil, err
		}
		list = append(list, l)
	}
	return list, nil
}

// CreateLecturer -> tambah dosen baru | FR-009
func CreateLecturer(ctx context.Context, l *models.Lecturer) error {
	_, err := database.DB.Exec(ctx,
		`INSERT INTO lecturers (id, user_id, lecturer_id, department, created_at)
		 VALUES ($1,$2,$3,$4,NOW())`,
		l.ID, l.UserID, l.LecturerID, l.Department)
	return err
}

// UpdateLecturer -> update data dosen | FR-009
func UpdateLecturer(ctx context.Context, l *models.Lecturer) error {
	_, err := database.DB.Exec(ctx,
		`UPDATE lecturers SET lecturer_id=$2, department=$3 WHERE id=$1`,
		l.ID, l.LecturerID, l.Department)
	return err
}

// DeleteLecturer -> hapus dosen | FR-009
func DeleteLecturer(ctx context.Context, id string) error {
	_, err := database.DB.Exec(ctx, `DELETE FROM lecturers WHERE id=$1`, id)
	return err
}

// GetFirstLecturer -> ambil dosen pertama (default advisor) | FR-009
func GetFirstLecturer(ctx context.Context) (*models.Lecturer, error) {
	row := database.DB.QueryRow(ctx,
		`SELECT id, user_id, created_at
		 FROM lecturers
		 ORDER BY created_at ASC
		 LIMIT 1`)

	var l models.Lecturer
	if err := row.Scan(&l.ID, &l.UserID, &l.CreatedAt); err != nil {
		return nil, err
	}
	return &l, nil
}

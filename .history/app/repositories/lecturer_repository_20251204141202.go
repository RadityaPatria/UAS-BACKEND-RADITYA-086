package repositories

import (
	"context"
	"UAS-backend/app/models"
	"UAS-backend/database"
)

// =========================
// GET lecturer.id BY user_id
// =========================
func GetLecturerIDByUserID(ctx context.Context, userID string) (string, error) {
	var id string
	err := database.DB.QueryRow(ctx,
		`SELECT id FROM lecturers WHERE user_id = $1`,
		userID,
	).Scan(&id)

	return id, err
}

// =========================
// GET lecturer profile by lecturer.id
// =========================
func GetLecturerByID(ctx context.Context, id string) (*models.Lecturer, error) {
	row := database.DB.QueryRow(ctx,
		`SELECT id, user_id, lecturer_id, department, created_at
		 FROM lecturers WHERE id = $1`, id)

	var l models.Lecturer
	err := row.Scan(&l.ID, &l.UserID, &l.LecturerID, &l.Department, &l.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &l, nil
}

// =========================
// GET all lecturers (Admin)
// =========================
func GetAllLecturers(ctx context.Context) ([]models.Lecturer, error) {
	rows, err := database.DB.Query(ctx,
		`SELECT id, user_id, lecturer_id, department, created_at FROM lecturers`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	list := []models.Lecturer{}
	for rows.Next() {
		var l models.Lecturer
		rows.Scan(&l.ID, &l.UserID, &l.LecturerID, &l.Department, &l.CreatedAt)
		list = append(list, l)
	}

	return list, nil
}

// =========================
// CREATE lecturer
// =========================
func CreateLecturer(ctx context.Context, l *models.Lecturer) error {
	_, err := database.DB.Exec(ctx,
		`INSERT INTO lecturers (id, user_id, lecturer_id, department, created_at)
		 VALUES ($1,$2,$3,$4,NOW())`,
		l.ID, l.UserID, l.LecturerID, l.Department,
	)
	return err
}

// =========================
// UPDATE lecturer
// =========================
func UpdateLecturer(ctx context.Context, l *models.Lecturer) error {
	_, err := database.DB.Exec(ctx,
		`UPDATE lecturers SET lecturer_id=$2, department=$3 WHERE id=$1`,
		l.ID, l.LecturerID, l.Department,
	)
	return err
}

// =========================
// DELETE lecturer
// =========================
func DeleteLecturer(ctx context.Context, id string) error {
	_, err := database.DB.Exec(ctx,
		`DELETE FROM lecturers WHERE id=$1`, id)
	return err
}

package repositories

import (
	"context"

	"UAS-backend/app/models"
	"UAS-backend/database"

	"github.com/google/uuid"
)

type LecturerRepository struct{}

var LecturerRepo = &LecturerRepository{}

//
// =======================================================
// AMBIL PROFIL DOSEN BERDASARKAN USER_ID (Dipakai LOGIN)
// =======================================================
//
func (r *LecturerRepository) GetByUserID(ctx context.Context, userID uuid.UUID) (*models.Lecturer, error) {
	row := database.DB.QueryRow(ctx,
		`SELECT id, user_id, lecturer_id, department, created_at
		 FROM lecturers WHERE user_id=$1`,
		userID,
	)

	var l models.Lecturer

	err := row.Scan(
		&l.ID,
		&l.UserID,
		&l.LecturerID,
		&l.Department,
		&l.CreatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &l, nil
}

//
// =======================================================
// AMBIL DATA DOSEN BERDASARKAN LECTURER.ID
// =======================================================
//
func (r *LecturerRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.Lecturer, error) {
	row := database.DB.QueryRow(ctx,
		`SELECT id, user_id, lecturer_id, department, created_at
		 FROM lecturers WHERE id=$1`,
		id,
	)

	var l models.Lecturer

	err := row.Scan(
		&l.ID,
		&l.UserID,
		&l.LecturerID,
		&l.Department,
		&l.CreatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &l, nil
}

//
// =======================================================
// ADMIN: LIST SEMUA DOSEN
// =======================================================
//
func (r *LecturerRepository) GetAll(ctx context.Context) ([]models.Lecturer, error) {
	rows, err := database.DB.Query(ctx,
		`SELECT id, user_id, lecturer_id, department, created_at
		 FROM lecturers`,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []models.Lecturer

	for rows.Next() {
		var l models.Lecturer
		err := rows.Scan(
			&l.ID,
			&l.UserID,
			&l.LecturerID,
			&l.Department,
			&l.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		list = append(list, l)
	}

	return list, nil
}

//
// =======================================================
// ADMIN: CREATE DOSEN
// =======================================================
//
func (r *LecturerRepository) Create(ctx context.Context, l *models.Lecturer) error {
	_, err := database.DB.Exec(ctx,
		`INSERT INTO lecturers (id, user_id, lecturer_id, department)
		 VALUES ($1, $2, $3, $4)`,
		l.ID,
		l.UserID,
		l.LecturerID,
		l.Department,
	)
	return err
}

//
// =======================================================
// ADMIN: UPDATE DOSEN
// =======================================================
//
func (r *LecturerRepository) Update(ctx context.Context, l *models.Lecturer) error {
	_, err := database.DB.Exec(ctx,
		`UPDATE lecturers 
		 SET lecturer_id=$2, department=$3
		 WHERE id=$1`,
		l.ID,
		l.LecturerID,
		l.Department,
	)
	return err
}

//
// =======================================================
// ADMIN: DELETE DOSEN
// =======================================================
//
func (r *LecturerRepository) Delete(ctx context.Context, id uuid.UUID) error {
	_, err := database.DB.Exec(ctx,
		`DELETE FROM lecturers WHERE id=$1`,
		id,
	)
	return err
}

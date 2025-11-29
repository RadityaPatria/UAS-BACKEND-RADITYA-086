package repositories

import (
	"context"


	"UAS-backend/app/models"
	"UAS-backend/database"

	"github.com/google/uuid"
)

type AchievementReferenceRepository struct{}

var AchievementReferenceRepo = &AchievementReferenceRepository{}

// ========================================================
// INSERT REFERENCE PRESTASI (status: draft)
// Dipanggil ketika mahasiswa membuat prestasi baru
// ========================================================
func (r *AchievementReferenceRepository) Create(ctx context.Context, ref *models.AchievementReference) error {
	query := `
		INSERT INTO achievement_references 
		(id, student_id, mongo_achievement_id, status, created_at, updated_at)
		VALUES ($1, $2, $3, $4, NOW(), NOW())
	`

	_, err := database.DB.Exec(ctx, query,
		ref.ID,
		ref.StudentID,
		ref.MongoAchievementID,
		models.StatusDraft,
	)

	return err
}

// ========================================================
// GET REFERENCE BY ID
// ========================================================
func (r *AchievementReferenceRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.AchievementReference, error) {
	query := `
		SELECT id, student_id, mongo_achievement_id, status,
		       submitted_at, verified_at, verified_by, rejection_note,
		       created_at, updated_at
		FROM achievement_references
		WHERE id=$1
	`

	row := database.DB.QueryRow(ctx, query, id)

	var ar models.AchievementReference
	err := row.Scan(
		&ar.ID,
		&ar.StudentID,
		&ar.MongoAchievementID,
		&ar.Status,
		&ar.SubmittedAt,
		&ar.VerifiedAt,
		&ar.VerifiedBy,
		&ar.RejectionNote,
		&ar.CreatedAt,
		&ar.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &ar, nil
}

// ========================================================
// UPDATE STATUS → submitted / verified / rejected
// ========================================================
func (r *AchievementReferenceRepository) UpdateStatus(ctx context.Context, id uuid.UUID, status string) error {
	var query string

	switch status {
	case models.StatusSubmitted:
		query = `
			UPDATE achievement_references
			SET status=$2, submitted_at=NOW(), updated_at=NOW()
			WHERE id=$1
		`

	case models.StatusVerified:
		query = `
			UPDATE achievement_references
			SET status=$2, verified_at=NOW(), updated_at=NOW()
			WHERE id=$1
		`

	case models.StatusRejected:
		query = `
			UPDATE achievement_references
			SET status=$2, updated_at=NOW()
			WHERE id=$1
		`
	}

	_, err := database.DB.Exec(ctx, query, id, status)
	return err
}

// ========================================================
// SET REJECTION NOTE (Dosen Wali)
// ========================================================
func (r *AchievementReferenceRepository) SetRejectionNote(ctx context.Context, id uuid.UUID, note string) error {
	query := `
		UPDATE achievement_references
		SET rejection_note=$2, status='rejected', updated_at=NOW()
		WHERE id=$1
	`
	_, err := database.DB.Exec(ctx, query, id, note)
	return err
}

// ========================================================
// ADMIN & DOSEN WALI → FILTER BY STUDENT IDs
// ========================================================
func (r *AchievementReferenceRepository) GetByStudents(ctx context.Context, studentIDs []uuid.UUID) ([]models.AchievementReference, error) {
	query := `
		SELECT id, student_id, mongo_achievement_id, status,
		       submitted_at, verified_at, verified_by, rejection_note,
		       created_at, updated_at
		FROM achievement_references
		WHERE student_id = ANY($1)
	`

	rows, err := database.DB.Query(ctx, query, studentIDs)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []models.AchievementReference

	for rows.Next() {
		var ar models.AchievementReference
		err := rows.Scan(
			&ar.ID,
			&ar.StudentID,
			&ar.MongoAchievementID,
			&ar.Status,
			&ar.SubmittedAt,
			&ar.VerifiedAt,
			&ar.VerifiedBy,
			&ar.RejectionNote,
			&ar.CreatedAt,
			&ar.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		list = append(list, ar)
	}

	return list, nil
}

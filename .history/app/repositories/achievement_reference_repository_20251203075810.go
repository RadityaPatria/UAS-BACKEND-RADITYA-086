package repositories

import (
	"context"

	"UAS-backend/app/models"
	"UAS-backend/database"
)

// -------------------------------------------------------
// CreateAchievementReference -> insert reference
// -------------------------------------------------------
func CreateAchievementReference(ctx context.Context, ref *models.AchievementReference) error {
	_, err := database.DB.Exec(ctx,
		`INSERT INTO achievement_references 
			(id, student_id, mongo_achievement_id, status, created_at, updated_at)
		 VALUES ($1, $2, $3, $4, NOW(), NOW())`,
		ref.ID, ref.StudentID, ref.MongoAchievementID, ref.Status,
	)
	return err
}

// -------------------------------------------------------
// GetAchievementReferenceByID
// -------------------------------------------------------
func GetAchievementReferenceByID(ctx context.Context, id string) (*models.AchievementReference, error) {

	row := database.DB.QueryRow(ctx,
		`SELECT id, student_id, mongo_achievement_id, status, 
				submitted_at, verified_at, verified_by, rejection_note,
				created_at, updated_at
		 FROM achievement_references WHERE id=$1`,
		id,
	)

	var ar models.AchievementReference
	err := row.Scan(
		&ar.ID, &ar.StudentID, &ar.MongoAchievementID, &ar.Status,
		&ar.SubmittedAt, &ar.VerifiedAt, &ar.VerifiedBy, &ar.RejectionNote,
		&ar.CreatedAt, &ar.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &ar, nil
}

// -------------------------------------------------------
// UpdateAchievementReferenceStatus
// -------------------------------------------------------
func UpdateAchievementReferenceStatus(ctx context.Context, id string, status string) error {

	_, err := database.DB.Exec(ctx,
		`UPDATE achievement_references 
		 SET status=$2, updated_at=NOW()
		 WHERE id=$1`,
		id, status,
	)
	return err
}

// -------------------------------------------------------
// SetAchievementRejectionNote
// -------------------------------------------------------
func SetAchievementRejectionNote(ctx context.Context, id string, note string) error {
	_, err := database.DB.Exec(ctx,
		`UPDATE achievement_references 
		 SET rejection_note=$2, status='rejected', updated_at=NOW() 
		 WHERE id=$1`,
		id, note,
	)
	return err
}

// -------------------------------------------------------
// SetVerifiedBy -> DOSEN VERIFIKASI
// -------------------------------------------------------
func SetVerifiedBy(ctx context.Context, id string, lecturerID string) error {
	_, err := database.DB.Exec(ctx,
		`UPDATE achievement_references
		 SET verified_by=$2, verified_at=NOW(), updated_at=NOW()
		 WHERE id=$1`,
		id, lecturerID,
	)
	return err
}

// -------------------------------------------------------
// GetAchievementReferencesByStudentIDs
// -------------------------------------------------------
func GetAchievementReferencesByStudentIDs(ctx context.Context, studentIDs []string) ([]models.AchievementReference, error) {

	rows, err := database.DB.Query(ctx,
		`SELECT id, student_id, mongo_achievement_id, status, 
				submitted_at, verified_at, verified_by, rejection_note,
				created_at, updated_at
		 FROM achievement_references 
		 WHERE student_id = ANY($1)`,
		studentIDs,
	)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []models.AchievementReference

	for rows.Next() {
		var ar models.AchievementReference

		err := rows.Scan(
			&ar.ID, &ar.StudentID, &ar.MongoAchievementID, &ar.Status,
			&ar.SubmittedAt, &ar.VerifiedAt, &ar.VerifiedBy, &ar.RejectionNote,
			&ar.CreatedAt, &ar.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		list = append(list, ar)
	}

	return list, nil
}

// -------------------------------------------------------
// SoftDeleteAchievementReference
// -------------------------------------------------------
func SoftDeleteAchievementReference(ctx context.Context, id string) error {
	_, err := database.DB.Exec(ctx,
		`UPDATE achievement_references 
		 SET deleted=true, updated_at=NOW() 
		 WHERE id=$1`,
		id,
	)
	return err
}

// -------------------------------------------------------
// ListAllAchievementReferences
// -------------------------------------------------------
func ListAllAchievementReferences(ctx context.Context) ([]models.AchievementReference, error) {

	rows, err := database.DB.Query(ctx,
		`SELECT id, student_id, mongo_achievement_id, status, 
				submitted_at, verified_at, verified_by, rejection_note,
				created_at, updated_at
		 FROM achievement_references
		 ORDER BY created_at DESC`,
	)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []models.AchievementReference

	for rows.Next() {
		var ar models.AchievementReference

		err := rows.Scan(
			&ar.ID, &ar.StudentID, &ar.MongoAchievementID, &ar.Status,
			&ar.SubmittedAt, &ar.VerifiedAt, &ar.VerifiedBy, &ar.RejectionNote,
			&ar.CreatedAt, &ar.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		list = append(list, ar)
	}

	return list, nil
}

package repositories

import (
	"context"

	"UAS-backend/app/models"
	"UAS-backend/database"
)

// -------------------------------------------------------
// CreateAchievementReference -> insert reference (status draft)
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
// GetAchievementReferenceByID -> ambil reference by primary id
// -------------------------------------------------------
func GetAchievementReferenceByID(ctx context.Context, id string) (*models.AchievementReference, error) {

	row := database.DB.QueryRow(ctx,
		`SELECT id, student_id, mongo_achievement_id, status, 
				submitted_at, verified_at, verified_by, rejection_note,
				created_at, updated_at
		 FROM achievement_references 
		 WHERE id = $1`,
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
// UpdateAchievementReferenceStatus -> ubah status submitted/verified/rejected
// -------------------------------------------------------
func UpdateAchievementReferenceStatus(ctx context.Context, id string, status string) error {

	var query string

	switch status {

	case models.StatusSubmitted:
		query = `UPDATE achievement_references 
				 SET status=$2, submitted_at=NOW(), updated_at=NOW() 
				 WHERE id=$1`

	case models.StatusVerified:
		query = `UPDATE achievement_references 
				 SET status=$2, verified_at=NOW(), updated_at=NOW() 
				 WHERE id=$1`

	case models.StatusRejected:
		query = `UPDATE achievement_references 
				 SET status=$2, updated_at=NOW() 
				 WHERE id=$1`

	default:
		query = `UPDATE achievement_references 
				 SET status=$2, updated_at=NOW() 
				 WHERE id=$1`
	}

	_, err := database.DB.Exec(ctx, query, id, status)
	return err
}



// -------------------------------------------------------
// SetAchievementRejectionNote -> dosen input catatan penolakan
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
// GetAchievementReferencesByStudentIDs -> untuk admin/dosen
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

		if err := rows.Scan(
			&ar.ID, &ar.StudentID, &ar.MongoAchievementID, &ar.Status,
			&ar.SubmittedAt, &ar.VerifiedAt, &ar.VerifiedBy, &ar.RejectionNote,
			&ar.CreatedAt, &ar.UpdatedAt,
		); err != nil {
			return nil, err
		}

		list = append(list, ar)
	}

	return list, nil
}

// -------------------------------------------------------
// SoftDeleteAchievementReference -> tandai deleted=true (opsional)
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
// ListAllAchievementReferences -> untuk admin
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

		if err := rows.Scan(
			&ar.ID, &ar.StudentID, &ar.MongoAchievementID, &ar.Status,
			&ar.SubmittedAt, &ar.VerifiedAt, &ar.VerifiedBy, &ar.RejectionNote,
			&ar.CreatedAt, &ar.UpdatedAt,
		); err != nil {
			return nil, err
		}

		list = append(list, ar)
	}

	return list, nil
}

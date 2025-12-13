package repositories

import (
	"context"

	"UAS-backend/app/models"
	"UAS-backend/database"

	"github.com/google/uuid"
)

// CreateAchievementReference -> simpan referensi prestasi mahasiswa | FR-007
func CreateAchievementReference(ctx context.Context, ref *models.AchievementReference) error {
	_, err := database.DB.Exec(ctx,
		`INSERT INTO achievement_references
		 (id, student_id, mongo_achievement_id, status, created_at, updated_at)
		 VALUES ($1,$2,$3,$4,NOW(),NOW())`,
		ref.ID, ref.StudentID, ref.MongoAchievementID, ref.Status,
	)
	return err
}

// GetAchievementReferenceByID -> detail prestasi berdasarkan ID | FR-007
func GetAchievementReferenceByID(ctx context.Context, id string) (*models.AchievementReference, error) {
	row := database.DB.QueryRow(ctx,
		`SELECT id, student_id, mongo_achievement_id, status,
		        submitted_at, verified_at, verified_by, rejection_note,
		        created_at, updated_at
		 FROM achievement_references WHERE id=$1`,
		id,
	)

	var ar models.AchievementReference
	if err := row.Scan(
		&ar.ID, &ar.StudentID, &ar.MongoAchievementID, &ar.Status,
		&ar.SubmittedAt, &ar.VerifiedAt, &ar.VerifiedBy, &ar.RejectionNote,
		&ar.CreatedAt, &ar.UpdatedAt,
	); err != nil {
		return nil, err
	}
	return &ar, nil
}

// UpdateAchievementReferenceStatus -> ubah status prestasi | FR-007
func UpdateAchievementReferenceStatus(ctx context.Context, id string, status string) error {
	_, err := database.DB.Exec(ctx,
		`UPDATE achievement_references
		 SET status=$2, updated_at=NOW()
		 WHERE id=$1`,
		id, status,
	)
	return err
}

// RejectAchievementPostgres -> tolak prestasi + catatan | FR-008
func RejectAchievementPostgres(ctx context.Context, refID string, note string) error {
	_, err := database.DB.Exec(ctx,
		`UPDATE achievement_references
		 SET rejection_note=$1,
		     status='rejected',
		     updated_at=NOW()
		 WHERE id=$2`,
		note, refID,
	)
	return err
}

// SetVerifiedBy -> verifikasi prestasi oleh dosen wali | FR-008
func SetVerifiedBy(ctx context.Context, refID string, lecturerUUID uuid.UUID) error {
	_, err := database.DB.Exec(ctx,
		`UPDATE achievement_references
		 SET verified_by=$1,
		     verified_at=NOW(),
		     updated_at=NOW()
		 WHERE id=$2`,
		lecturerUUID, refID,
	)
	return err
}

// SetSubmitted -> submit prestasi oleh mahasiswa | FR-007
func SetSubmitted(ctx context.Context, refID string) error {
	_, err := database.DB.Exec(ctx,
		`UPDATE achievement_references
		 SET status='submitted',
		     submitted_at=NOW(),
		     updated_at=NOW()
		 WHERE id=$1`,
		refID,
	)
	return err
}

// GetAchievementReferencesByStudentIDs -> list prestasi mahasiswa | FR-006
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

// SoftDeleteAchievementReference -> soft delete prestasi | FR-009
func SoftDeleteAchievementReference(ctx context.Context, id string) error {
	_, err := database.DB.Exec(ctx,
		`UPDATE achievement_references
		 SET is_deleted=TRUE, deleted_at=NOW()
		 WHERE id=$1`,
		id,
	)
	return err
}

// ListAllAchievementReferences -> list semua prestasi | FR-009
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

// TouchAchievementReference -> update timestamp | FR-007
func TouchAchievementReference(ctx context.Context, refID string) error {
	_, err := database.DB.Exec(ctx,
		`UPDATE achievement_references
		 SET updated_at=NOW()
		 WHERE id=$1`,
		refID,
	)
	return err
}

// SetDeleted -> tandai prestasi dihapus | FR-009
func SetDeleted(ctx context.Context, id string) error {
	_, err := database.DB.Exec(ctx,
		`UPDATE achievement_references
		 SET status='deleted',
		     updated_at=NOW(),
		     deleted_at=NOW()
		 WHERE id=$1`,
		id,
	)
	return err
}

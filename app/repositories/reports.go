package repositories

import (
	"context"

	"UAS-backend/database"
)

// GetAchievementStatistics -> statistik jumlah achievement per status | FR-008
func GetAchievementStatistics(ctx context.Context) (map[string]int, error) {
	rows, err := database.DB.Query(ctx,
		`SELECT status, COUNT(*)
		 FROM achievement_references
		 GROUP BY status`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	stats := make(map[string]int)
	for rows.Next() {
		var status string
		var count int
		if err := rows.Scan(&status, &count); err != nil {
			return nil, err
		}
		stats[status] = count
	}

	return stats, nil
}

// GetStudentReport -> laporan jumlah achievement mahasiswa | FR-008
func GetStudentReport(ctx context.Context, studentID string) (map[string]any, error) {
	row := database.DB.QueryRow(ctx,
		`SELECT student_id, COUNT(*)
		 FROM achievement_references
		 WHERE student_id=$1
		 GROUP BY student_id`, studentID)

	var sid string
	var total int
	if err := row.Scan(&sid, &total); err != nil {
		return nil, err
	}

	return map[string]any{
		"student_id": sid,
		"total":      total,
	}, nil
}

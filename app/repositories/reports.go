package repositories

import (
	"context"

	"UAS-backend/database"
)

func GetAchievementStatistics(ctx context.Context) (map[string]int, error) {
	rows, err := database.DB.Query(ctx,
		`SELECT status, COUNT(*) FROM achievement_references GROUP BY status`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	stats := map[string]int{}
	var status string
	var count int

	for rows.Next() {
		rows.Scan(&status, &count)
		stats[status] = count
	}

	return stats, nil
}

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

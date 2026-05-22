package repository

import (
	"context"
	"database/sql"
	"time"

	"github.com/erdinhrmwn/hacktiv8-library-cli/internal/model"
)

type ActivityRepository struct {
	db *sql.DB
}

func NewActivityRepository(db *sql.DB) *ActivityRepository {
	return &ActivityRepository{db: db}
}

func (r *ActivityRepository) GetAll(ctx context.Context) ([]model.ActivityLog, error) {
	rows, err := r.db.QueryContext(ctx, "SELECT * FROM activity_logs ORDER BY date DESC")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var logs []model.ActivityLog
	for rows.Next() {
		var l model.ActivityLog
		if err := rows.Scan(&l.ID, &l.Key, &l.Description, &l.Date); err != nil {
			return nil, err
		}
		logs = append(logs, l)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return logs, nil
}

func (r *ActivityRepository) Log(ctx context.Context, key, description string) error {
	_, err := r.db.ExecContext(ctx,
		`INSERT INTO activity_logs (key, description, date) VALUES (?, ?, ?)`,
		key, description, time.Now())
	return err
}

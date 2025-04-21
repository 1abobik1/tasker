package repository

import (
	"context"
	"errors"
	"time"

	"github.com/1abobik1/tasker/internal/errs"
	"github.com/1abobik1/tasker/internal/models"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresRepo struct {
	pool *pgxpool.Pool
}

func NewPostgresRepo(pool *pgxpool.Pool) *PostgresRepo {
	return &PostgresRepo{pool: pool}
}

func (r *PostgresRepo) Create(ctx context.Context, t models.Task) error {
	_, err := r.pool.Exec(ctx, `
      INSERT INTO tasks (id, task_type, payload, status, created_at, updated_at)
      VALUES ($1, $2, $3, $4, $5, $6)`,
		t.ID, t.Type, t.Payload, t.Status, t.CreatedAt, t.UpdatedAt,
	)
	return err
}

func (r *PostgresRepo) GetByID(ctx context.Context, id string) (models.Task, error) {
	row := r.pool.QueryRow(ctx, `
        SELECT id, task_type, payload, status, result, error, created_at, updated_at
          FROM tasks WHERE id = $1`, id)

	var t models.Task
	err := row.Scan(
		&t.ID, &t.Type,
		&t.Payload, &t.Status, &t.Result,
		&t.Error, &t.CreatedAt, &t.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return models.Task{}, errs.ErrIDNotFound
		}
		return models.Task{}, err
	}
	return t, nil
}

func (r *PostgresRepo) UpdateStatus(ctx context.Context, id string, status models.TaskStatus) error {
	_, err := r.pool.Exec(ctx, `
      UPDATE tasks SET status = $2, updated_at = $3 WHERE id = $1`,
		id, status, time.Now().UTC(),
	)
	return err
}

func (r *PostgresRepo) SaveResult(ctx context.Context, id string, result []byte) error {
	_, err := r.pool.Exec(ctx, `
      UPDATE tasks SET result = $2, status = $3, updated_at = $4 WHERE id = $1`,
		id, result, models.StatusCompleted, time.Now().UTC(),
	)
	return err
}

func (r *PostgresRepo) SaveError(ctx context.Context, id, errMsg string) error {
	_, err := r.pool.Exec(ctx, `
      UPDATE tasks SET error = $2, status = $3, updated_at = $4 WHERE id = $1`,
		id, errMsg, models.StatusFailed, time.Now().UTC(),
	)
	return err
}

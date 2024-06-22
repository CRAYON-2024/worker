package repository

import (
	"context"

	"github.com/CRAYON-2024/worker/internal/entity"
	"github.com/jackc/pgx/v5"
	"go.opentelemetry.io/otel/trace"
)

type UserRepository struct {
	DB    *pgx.Conn
	Trace trace.Tracer
}

func NewUserRepository(db *pgx.Conn, trace trace.Tracer) *UserRepository {
	return &UserRepository{
		DB:    db,
		Trace: trace,
	}
}

func (ur *UserRepository) GetUsers(ctx context.Context) ([]entity.CustomUser, error) {
	ctx, span := ur.Trace.Start(ctx, "repo get users")
	defer span.End()

	rows, err := ur.DB.Query(ctx, "select id, name from users")
	if err != nil {
		return nil, err
	}

	var users []entity.CustomUser
	for rows.Next() {
		var tmp entity.CustomUser
		if err := rows.Scan(
			&tmp.ID,
			&tmp.Name,
		); err != nil {
			return nil, err
		}

		users = append(users, tmp)
	}
	return users, nil
}

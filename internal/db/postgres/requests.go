package postgres

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/dandolya/corp-app/internal/approval"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type RequestRepository struct {
	pool *pgxpool.Pool
}

func NewRequestRepository(pool *pgxpool.Pool) *RequestRepository {
	return &RequestRepository{pool: pool}
}

func (r *RequestRepository) Create(ctx context.Context, request approval.Request) (int64, error) {

	if len(request.Details) == 0 {
		request.Details = json.RawMessage("{}")
	}

	repo := NewEmployeeRepository(r.pool)

	chain, err := repo.GetManagerChain(ctx, request.InitiatorID)
	if err != nil {
		return 0, err
	}

	tx, err := r.pool.BeginTx(ctx, pgx.TxOptions{IsoLevel: pgx.RepeatableRead})
	if err != nil {
		return 0, err
	}
	defer tx.Rollback(ctx)

	var id int64
	err = tx.QueryRow(ctx,
		`insert into requests (type, initiator_id, target_id, details, expires_at)
			values ($1, $2, $3, $4, $5) returning id`,
		request.Type, request.InitiatorID, request.TargetID, request.Details, request.ExpiresAt,
	).Scan(&id)

	if err != nil {
		return 0, fmt.Errorf("insert request from %d: %w", request.InitiatorID, err)
	}

	for i, m := range chain {
		_, err = tx.Exec(ctx,
			`insert into notifications (request_id, recipient_id, is_direct, state)
		 values ($1, $2, $3, 'actionable')`,
			id, m.ID, i == 0)
		if err != nil {
			return 0, fmt.Errorf("insert notification for %d: %w", m.ID, err)
		}
	}

	return id, tx.Commit(ctx)
}

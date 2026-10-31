package postgres

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

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

func (r *RequestRepository) CreateRequest(ctx context.Context, request approval.Request) (int64, error) {

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

func (r *RequestRepository) ResolveRequest(ctx context.Context, id, requestID int64, status approval.Status, comment string) error {
	if status != approval.StatusApproved && status != approval.StatusRejected {
		return approval.ErrInvalidDecision
	}

	chain, err := NewEmployeeRepository(r.pool).GetManagerChain(ctx, id)
	if err != nil {
		return err
	}
	aboveIDs := make([]int64, 0, len(chain))
	for _, m := range chain {
		aboveIDs = append(aboveIDs, m.ID)
	}

	tx, err := r.pool.BeginTx(ctx, pgx.TxOptions{IsoLevel: pgx.RepeatableRead})
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	var exists bool
	err = tx.QueryRow(ctx,
		`select exists (select 1 from notifications where request_id = $1 and recipient_id = $2)`,
		requestID, id,
	).Scan(&exists)
	if err != nil {
		return fmt.Errorf("check access to request %d: %w", requestID, err)
	}
	if !exists {
		return approval.ErrNotAllowed
	}

	tag, err := tx.Exec(ctx,
		`update requests
			set status = $1, resolved_by = $2, resolved_at = now(), resolution_comment = nullif($3, '')
			where id = $4 and status = 'pending' and expires_at > now()`,
		status, id, comment, requestID)
	if err != nil {
		return fmt.Errorf("resolve request %d by %d: %w", requestID, id, err)
	}

	if tag.RowsAffected() == 0 {
		return r.whyNotResolvable(ctx, tx, requestID)
	}

	_, err = tx.Exec(ctx,
		`update notifications
		 set state = case
		     when recipient_id = $2 then 'resolved_by_me'
		     when recipient_id = any($3) then 'resolved_below'
		     else 'resolved_above'
		 end
		 where request_id = $1`,
		requestID, id, aboveIDs)
	if err != nil {
		return fmt.Errorf("update notifications for request %d: %w", requestID, err)
	}

	return tx.Commit(ctx)
}

func (r *RequestRepository) whyNotResolvable(ctx context.Context, tx pgx.Tx, requestID int64) error {
	var status approval.Status
	var expiresAt time.Time

	err := tx.QueryRow(ctx,
		`select status, expires_at from requests where id = $1`, requestID,
	).Scan(&status, &expiresAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return approval.ErrRequestNotFound
	}
	if err != nil {
		return fmt.Errorf("check request %d: %w", requestID, err)
	}

	switch status {
	case approval.StatusApproved, approval.StatusRejected:
		return approval.ErrAlreadyResolved
	case approval.StatusCancelled:
		return approval.ErrCancelled
	default:
		return approval.ErrExpired
	}
}

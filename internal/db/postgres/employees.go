package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/dandolya/corp-app/internal/org"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type EmployeeRepo struct {
	pool *pgxpool.Pool
}

func NewEmployeeRepo(pool *pgxpool.Pool) *EmployeeRepo {
	return &EmployeeRepo{pool: pool}
}

func (r *EmployeeRepo) GetByID(ctx context.Context, id int64) (org.Employee, error) {
	row := r.pool.QueryRow(ctx, "select id, name, role, position, manager_id, salary, status, changed_at, hired_at from employees where id = $1;", id)

	e := org.Employee{}
	var position *string
	err := row.Scan(&e.ID, &e.Name, &e.Role, &position, &e.ManagerID, &e.Salary, &e.Status, &e.ChangedAt, &e.HiredAt)

	if errors.Is(err, pgx.ErrNoRows) {
		return org.Employee{}, org.ErrNotFound
	}

	if err != nil {
		return org.Employee{}, fmt.Errorf("get employee %d: %w", id, err)
	}

	if position != nil {
		e.Position = org.Position(*position)
	}

	return e, nil
}

func (r *EmployeeRepo) GetManagerChain(ctx context.Context, id int64) ([]org.Employee, error) {
	res := make([]org.Employee, 0, 4)
	e, err := r.GetByID(ctx, id)

	if err != nil {
		return nil, err
	}

	next := e.ManagerID

	for next != nil {
		e, err := r.GetByID(ctx, *next)
		if err != nil {
			return nil, err
		}
		res = append(res, e)
		next = e.ManagerID
	}

	return res, nil
}

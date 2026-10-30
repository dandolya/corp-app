package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/dandolya/corp-app/internal/org"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type EmployeeRepository struct {
	pool *pgxpool.Pool
}

func NewEmployeeRepository(pool *pgxpool.Pool) *EmployeeRepository {
	return &EmployeeRepository{pool: pool}
}

func scanEmployee(row pgx.Row) (org.Employee, error) {
	var e org.Employee
	var position *string

	err := row.Scan(&e.ID, &e.Name, &e.Role, &position, &e.ManagerID, &e.Salary, &e.Status, &e.ChangedAt, &e.HiredAt)
	if err != nil {
		return org.Employee{}, err
	}
	if position != nil {
		e.Position = org.Position(*position)

	}
	return e, nil
}

func (r *EmployeeRepository) GetByID(ctx context.Context, id int64) (org.Employee, error) {
	row := r.pool.QueryRow(ctx,
		`select id, name, role, position, manager_id, salary, status, changed_at, hired_at
		 from employees where id = $1`, id)

	e, err := scanEmployee(row)

	if errors.Is(err, pgx.ErrNoRows) {
		return org.Employee{}, org.ErrNotFound
	}
	if err != nil {
		return org.Employee{}, fmt.Errorf("get employee %d: %w", id, err)
	}

	return e, nil
}

const managerChainQuery = `
with recursive chain as (
    select id, manager_id, 0 as depth
    from employees where id = $1

    union all

    select e.id, e.manager_id, c.depth + 1
    from employees e
    join chain c on e.id = c.manager_id
)
select e.id, e.name, e.role, e.position, e.manager_id, e.salary, e.status, e.changed_at, e.hired_at
from chain c
join employees e on e.id = c.id
where c.depth > 0
order by c.depth`

func (r *EmployeeRepository) GetManagerChain(ctx context.Context, id int64) ([]org.Employee, error) {
	rows, err := r.pool.Query(ctx, managerChainQuery, id)
	if err != nil {
		return nil, fmt.Errorf("get manager chain of %d: %w", id, err)
	}
	defer rows.Close()

	res := make([]org.Employee, 0, 4)
	for rows.Next() {
		e, err := scanEmployee(rows)
		if err != nil {
			return nil, fmt.Errorf("scan manager of %d: %w", id, err)
		}
		res = append(res, e)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("get manager chain of %d: %w", id, err)
	}

	return res, nil
}

func (r *EmployeeRepository) IsManagerOf(ctx context.Context, managerID, employeeID int64) (bool, error) {
	chain, err := r.GetManagerChain(ctx, employeeID)
	if err != nil {
		return false, fmt.Errorf("is %d manager of %d: %w", managerID, employeeID, err)
	}

	for _, e := range chain {
		if e.ID == managerID {
			return true, nil
		}
	}

	return false, nil
}

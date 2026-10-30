-- +goose Up
create table employees (
    id bigint generated always as identity primary key,
    name text not null,
    role text not null check (role in ('ceo', 'director', 'admin', 'employee')),
    position text check (position in ('seller', 'security', 'helper')),
    manager_id bigint references employees(id),
    salary bigint not null check (salary > 0),
    status text not null check (status in ('active', 'fired', 'probation')),
    changed_at timestamptz not null default now(),
    hired_at timestamptz not null default now(),

    constraint manager_existing check ((role = 'ceo') = (manager_id is null)),
    constraint employee_position check ((role = 'employee') = (position is not null)),
    constraint self_manager check (id <> manager_id)
);

create index employees_manager_id_idx on employees (manager_id);


create table requests (
    id bigint generated always as identity primary key,
    type text not null check (type in ('salary_change', 'penalty', 'resignation', 'termination', 
    'vacation', 'sick_leave', 'day_off', 'position_change','promotion', 'transfer', 'hire')),
    status text not null check (status in ('pending', 'approved', 'rejected', 'cancelled', 'expired')) default 'pending',
    initiator_id bigint not null references employees(id),
    target_id bigint references employees(id),
    details jsonb not null default '{}',
    created_at timestamptz not null default now(),
    expires_at timestamptz not null,
    resolved_by bigint references employees(id),
    resolved_at timestamptz,

    constraint time_for_resolving check (expires_at > created_at),
    constraint resolving_person check ((resolved_by is null) = (status in ('pending', 'expired'))),
    constraint resolving_time check ((resolved_by is null) = (resolved_at is null))
);

create index requests_pending_expiring_idx on requests (expires_at) where status = 'pending';


create table notifications (
    request_id bigint not null references requests(id) on delete cascade,
    recipient_id bigint not null references employees(id),
    is_direct boolean not null,
    state text not null check (state in ('actionable', 'resolved_above', 'resolved_below', 'resolved_by_me', 'expired')),

    primary key (request_id, recipient_id)
);

create index notifications_recipient_index on notifications (recipient_id, state);

-- +goose Down
drop table notifications;

drop table requests;

drop table employees;
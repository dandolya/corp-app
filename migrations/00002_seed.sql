-- +goose Up
insert into employees (name, role, position, manager_id, salary, status)
values
    ('Иван Иванов', 'ceo', null, null, 1000000, 'active'),
    ('Петр Петров', 'director', null, 1, 300000, 'active'),
    ('Денис Денисов', 'director', null, 1, 350000, 'active'),
    ('Роман Романов', 'admin', null, 2, 100000, 'active'),
    ('Пример Примеров', 'admin', null, 2, 120000, 'active'),
    ('Тест Тестов', 'admin', null, 3, 90000, 'active'),
    ('Анна Аннова', 'employee', 'seller', 4, 60000, 'active'),
    ('Иван Петров', 'employee', 'security', 4, 56000, 'active'),
    ('Денис Романов', 'employee', 'seller', 5, 47000, 'active'),
    ('Роман Примеров', 'employee', 'helper', 5, 62000, 'active'),
    ('Анна Иванова', 'employee', 'seller', 6, 45000, 'probation'),
    ('Петр Иванов', 'employee', 'seller', 5, 19000, 'fired');

-- +goose Down
truncate employees restart identity cascade;
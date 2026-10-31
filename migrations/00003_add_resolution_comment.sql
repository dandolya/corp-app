-- +goose Up
alter table requests add column resolution_comment text;

-- +goose Down
alter table requests drop column resolution_comment;
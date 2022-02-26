-- +goose Up
-- +goose StatementBegin
create table matches (
    etf2l_match_id integer primary key,
    log_id integer default null,
    is_default_win boolean default true
);
create table logs (
    id integer references matches(log_id),
    map text,
    played_at timestamp not null
);
create table users (
    name text primary key,
    token text
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd

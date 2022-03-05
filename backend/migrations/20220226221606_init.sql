-- +goose Up
-- +goose StatementBegin
create table matches (
    competition_id integer references competitions(id),
    etf2l_match_id integer primary key,
    log_id integer default null,
    is_default_win boolean default true
);
create table logs (
    id integer,
    map text,
    played_at timestamp not null
);
create table competitions (
    id integer primary key,
    is_completed bool default false
);
create table users (
    name text primary key,
    token text
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table competitions cascade;
drop table matches cascade;
drop table logs cascade;
drop table users cascade;
-- +goose StatementEnd

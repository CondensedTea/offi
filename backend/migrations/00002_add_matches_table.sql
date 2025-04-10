-- +goose Up
-- +goose StatementBegin
create table matches (
	match_id bigint primary key,
	competition text not null,
	stage text not null,
	tier text not null,
	completed_at timestamptz not null,
	demo_id bigint
);

create table logs (
	log_id bigint primary key,
	title text not null,
	map text not null,
	played_at timestamptz not null,
	match_id bigint not null,
	is_secondary boolean default false
);

create index on logs(match_id);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table matches;
drop index logs_match_id_idx;
drop table logs;
-- +goose StatementEnd

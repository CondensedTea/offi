-- +goose Up
-- +goose StatementBegin
create table recruitments (
	recruitment_id bigint primary key,
	post_type text not null,
	author_id bigint not null,
	team_type text not null,
	skill_level text not null,
	classes text[] not null,
	created_at timestamp with time zone default now()
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table recruitments;
-- +goose StatementEnd

-- +goose Up
-- +goose StatementBegin
alter table matches drop column demo_id;
alter table logs add column demo_id bigint default null;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
alter table logs drop column demo_id;
alter table matches add column demo_id bigint default null;
-- +goose StatementEnd

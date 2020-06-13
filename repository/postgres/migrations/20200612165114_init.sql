-- +goose Up
create table cache_snapshot (
    id text primary key,
    json_value jsonb
);

-- +goose Down
drop table cache_snapshot;

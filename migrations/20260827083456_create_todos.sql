-- +goose Up
CREATE TABLE todos (
id SERIAL PRIMARY KEY,
title TEXT NOT NULL,
success BOOL NOT NULL DEFAULT false
);

-- +goose Down
DROP TABLE todos;

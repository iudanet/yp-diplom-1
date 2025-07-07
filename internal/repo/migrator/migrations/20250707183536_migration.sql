-- +goose Up
-- +goose StatementBegin
CREATE TABLE users (
    id int NOT NULL PRIMARY KEY,
    username text,
    name text,
    surname text
);

INSERT INTO users VALUES
(0, 'root', '', ''),
(1, 'vojtechvitek', 'Vojtech', 'Vitek');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE users;
-- +goose StatementEnd

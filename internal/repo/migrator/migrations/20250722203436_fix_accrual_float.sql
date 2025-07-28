-- +goose Up
-- +goose StatementBegin
ALTER TABLE orders ALTER COLUMN accrual TYPE numeric(10,2);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE orders ALTER COLUMN accrual TYPE integer USING (accrual::integer);
-- +goose StatementEnd

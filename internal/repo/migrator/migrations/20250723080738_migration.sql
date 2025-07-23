-- +goose Up
-- +goose StatementBegin
-- Изменяем тип accrual в orders на integer (копейки)
ALTER TABLE orders ALTER COLUMN accrual TYPE integer USING (ROUND(accrual * 100)::integer);

-- Изменяем тип sum в withdrawals на integer (копейки)
ALTER TABLE withdrawals ALTER COLUMN sum TYPE integer USING (ROUND(sum * 100)::integer);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE orders ALTER COLUMN accrual TYPE numeric(10,2) USING (accrual::numeric / 100);
ALTER TABLE withdrawals ALTER COLUMN sum TYPE numeric(10,2) USING (sum::numeric / 100);
-- +goose StatementEnd

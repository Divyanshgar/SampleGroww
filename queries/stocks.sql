-- name: GetStock :one
SELECT * FROM stocks
WHERE id = $1 LIMIT 1;

-- name: ListStocks :many
SELECT * FROM stocks
ORDER BY id;

-- name: ListStocksByUser :many
SELECT * FROM stocks
WHERE user_id = $1
ORDER BY date DESC;

-- name: CreateStock :one
INSERT INTO stocks (
  user_id, stock_name, action, quantity, price, original_value, date
) VALUES (
  $1, $2, $3, $4, $5, $6, $7
)
RETURNING *;

-- name: UpdateStock :one
UPDATE stocks
SET stock_name = $2, action = $3, quantity = $4, price = $5, original_value = $6, date = $7, updated_at = CURRENT_TIMESTAMP
WHERE id = $1
RETURNING *;

-- name: DeleteStock :exec
DELETE FROM stocks
WHERE id = $1;

-- name: GetUser :one
SELECT * FROM users
WHERE id = $1 LIMIT 1;

-- name: GetUserByEmail :one
SELECT * FROM users
WHERE email = $1 LIMIT 1;

-- name: ListUsers :many
SELECT * FROM users
ORDER BY id;

-- name: CreateUser :one
INSERT INTO users (
  first_name, last_name, email, phone, address, city, country
) VALUES (
  $1, $2, $3, $4, $5, $6, $7
)
RETURNING *;

-- name: UpdateUser :one
UPDATE users
SET first_name = $2, last_name = $3, email = $4, phone = $5, address = $6, city = $7, country = $8, updated_at = CURRENT_TIMESTAMP
WHERE id = $1
RETURNING *;

-- name: DeleteUser :exec
DELETE FROM users
WHERE id = $1;

-- name: UpdateUserOTP :one
UPDATE users
SET otp_code = $2, otp_expires = $3, otp_verified = $4, updated_at = CURRENT_TIMESTAMP
WHERE id = $1
RETURNING *;

-- name: VerifyUserOTP :one
UPDATE users
SET otp_verified = true, updated_at = CURRENT_TIMESTAMP
WHERE id = $1 AND otp_code = $2 AND otp_expires > CURRENT_TIMESTAMP
RETURNING *;

-- name: GetUserOTP :one
SELECT id, email, phone, otp_code, otp_expires, otp_verified FROM users
WHERE id = $1 LIMIT 1;

-- name: GetUserByEmailAndOTP :one
SELECT * FROM users
WHERE email = $1 AND otp_code = $2 LIMIT 1;

INSERT INTO accounts (
    owner,
    balance,
    currency
) VALUES (
    $1, $2, $3
) RETURNING *;

SELECT * FROM accounts
WHERE id = $1 LIMIT 1;

SELECT * FROM accounts
WHERE id = $1 LIMIT 1
FOR NO KEY UPDATE;

SELECT * FROM accounts
ORDER BY id
LIMIT $1
OFFSET $2;

UPDATE accounts
SET balance = balance + sqlc.arg(amount)
WHERE id = sqlc.arg(id)
RETURNING *;

DELETE FROM accounts
WHERE id = $1;
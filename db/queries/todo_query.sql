-- name: GetTodo :one
SELECT *
FROM todos
WHERE id = ?
LIMIT 1;

-- name: ListTodos :many
SELECT *
FROM todos;

-- name: CreateTodo :one
INSERT INTO todos (title, description)
VALUES (?, ?)
RETURNING *;

-- name: UpdateTodo :one
UPDATE todos
SET title       = ?,
    description = ?,
    completed   = ?
WHERE id = ?
RETURNING *;

-- name: DeleteTodo :exec
DELETE
FROM todos
WHERE id = ?;

-- name: CreateUser :one
insert into Users (id, created_at, updated_at, email, password , username)
values (
    gen_random_uuid(),
    now(),
    now(),
    ?,
    ?,
    ?
)
returning *;

-- name: GetUserById :one
SELECT * FROM Users 
WHERE id = ?;

-- name: GetUserByUsername :one
SELECT * from Users where username= ?;



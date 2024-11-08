-- name: CreateUser :one
insert into Users ( created_at,   username , password)
values (
    CURRENT_TIMESTAMP,
    ?,
    ?
)
returning *;

-- name: GetUserById :one
SELECT * FROM Users 
WHERE id = ?;

-- name: GetUserByUsername :one
SELECT * from Users where username= ?;



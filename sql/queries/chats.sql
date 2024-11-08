

-- name: NewChat :one
INSERT INTO Chats(id,creation_date)
VALUES (
    gen_random_uuid(),
    NOW()
)
RETURNING *;



-- name: GetChats :many
select * from Chats;

-- name: GetChatCreationDate :one
select creation_date from Chats where id = ? ;


-- name: DeleteChat :exec
DELETE from  Chats where id = ?  ;

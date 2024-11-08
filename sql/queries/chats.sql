

-- name: NewChat :one
INSERT INTO Chats(creation_date)
VALUES (
    NOW()
)
RETURNING *;



-- name: GetChats :many
select * from Chats;

-- name: GetChatCreationDate :one
select creation_date from Chats where id = ? ;


-- name: DeleteChat :exec
DELETE from  Chats where id = ?  ;

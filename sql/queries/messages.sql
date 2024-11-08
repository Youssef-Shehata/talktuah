
-- name: NewMessage :one
INSERT INTO Messages(sender_id,chat_id , content, sent_at)

VALUES (
    ?,
    ?,
    ?,
    NOW()

)
RETURNING *;



-- name: GetMessageId :one
select * from Messages where id = ? order by sent_at desc ;

-- name: GetMessagesByChatId :many
select * from Messages where sender_id= ? order by sent_at desc ;


-- name: DeleteMessage :exec
DELETE from  Messages where chat_id = ? and id = ?;




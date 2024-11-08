-- name: NewMember :one
INSERT INTO ChatMembers (chat_id, user_id ,join_date)

VALUES (
    ?,
    ?,
    NOW()
)
RETURNING *;




-- name: GetChatMembers :many
select * from ChatMembers where chat_id = ? order by join_date desc ;



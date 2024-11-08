// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: messages.sql

package database

import (
	"context"

	"github.com/google/uuid"
)

const deleteMessage = `-- name: DeleteMessage :exec
;


DELETE from  Messages where chat_id = ? and id = ?
`

type DeleteMessageParams struct {
	ChatID interface{}
	ID     uuid.UUID
}

func (q *Queries) DeleteMessage(ctx context.Context, arg DeleteMessageParams) error {
	_, err := q.db.ExecContext(ctx, deleteMessage, arg.ChatID, arg.ID)
	return err
}

const getMessageId = `-- name: GetMessageId :one
select id, sender_id, chat_id, content, sent_at from Messages where id = ? order by sent_at desc
`

func (q *Queries) GetMessageId(ctx context.Context, id uuid.UUID) (Message, error) {
	row := q.db.QueryRowContext(ctx, getMessageId, id)
	var i Message
	err := row.Scan(
		&i.ID,
		&i.SenderID,
		&i.ChatID,
		&i.Content,
		&i.SentAt,
	)
	return i, err
}

const getMessagesByChatId = `-- name: GetMessagesByChatId :many
;

select id, sender_id, chat_id, content, sent_at from Messages where sender_id= ? order by sent_at desc
`

func (q *Queries) GetMessagesByChatId(ctx context.Context, senderID interface{}) ([]Message, error) {
	rows, err := q.db.QueryContext(ctx, getMessagesByChatId, senderID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Message
	for rows.Next() {
		var i Message
		if err := rows.Scan(
			&i.ID,
			&i.SenderID,
			&i.ChatID,
			&i.Content,
			&i.SentAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const newMessage = `-- name: NewMessage :one
INSERT INTO Messages(id,sender_id,chat_id , content, sent_at)

VALUES (
    gen_random_uuid(),
    ?,
    ?,
    ?,
    NOW()

)
RETURNING id, sender_id, chat_id, content, sent_at
`

type NewMessageParams struct {
	SenderID interface{}
	ChatID   interface{}
	Content  string
}

func (q *Queries) NewMessage(ctx context.Context, arg NewMessageParams) (Message, error) {
	row := q.db.QueryRowContext(ctx, newMessage, arg.SenderID, arg.ChatID, arg.Content)
	var i Message
	err := row.Scan(
		&i.ID,
		&i.SenderID,
		&i.ChatID,
		&i.Content,
		&i.SentAt,
	)
	return i, err
}
-- +goose Up
CREATE TABLE Users (
id uuid PRIMARY KEY,
created_at TIMESTAMP not null DEFAULT CURRENT_TIMESTAMP,
email Text unique not null,
password Text not Null,
username text not null 
);
CREATE TABLE Messages(
id uuid PRIMARY KEY,
sender_id uuid references Users(id),
chat_id uuid references Chats(id),
content text not null ,
sent_at TIMESTAMP not null DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE Chats(
id uuid PRIMARY KEY,
creation_date TIMESTAMP not null DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE ChatMembers(
chat_id uuid references Chats(id),
user_id uuid references Users(id),
join_date TIMESTAMP not null DEFAULT CURRENT_TIMESTAMP,
primary key (chat_id , user_id)
);
-- +goose Down
DROP TABLE Users;
DROP TABLE Messages;
DROP TABLE Chats;
DROP TABLE ChatMembers;


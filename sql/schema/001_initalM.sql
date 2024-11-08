-- +goose Up
CREATE TABLE Users (
id INTEGER PRIMARY KEY AUTOINCREMENT,
created_at TIMESTAMP not null DEFAULT CURRENT_TIMESTAMP,
password Text not Null,
username text not null 
);
CREATE TABLE Messages(
id INTEGER PRIMARY KEY AUTOINCREMENT,
sender_id INTEGEGR references Users(id),
chat_id INTEGER references Chats(id),
content text not null ,
sent_at TIMESTAMP not null DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE Chats(
id INTEGER PRIMARY KEY AUTOINCREMENT,
creation_date TIMESTAMP not null DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE ChatMembers(
chat_id INTEGER references Chats(id),
user_id INTEGER references Users(id),
join_date TIMESTAMP not null DEFAULT CURRENT_TIMESTAMP,
primary key (chat_id , user_id)
);
-- +goose Down
DROP TABLE Users;
DROP TABLE Messages;
DROP TABLE Chats;
DROP TABLE ChatMembers;


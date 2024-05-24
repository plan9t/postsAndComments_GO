CREATE TABLE users (
user_id SERIAL PRIMARY KEY,
first_name varchar(32) NOT NULL,
last_name varchar(32) NOT NULL
);

CREATE TABLE posts (
post_id SERIAL PRIMARY KEY,
title varchar(64) NOT NULL,
content text NOT NULL,
commentable boolean NOT NULL,
created_time timestamp NOT NULL,
user_id INT,
FOREIGN KEY (user_id) REFERENCES users(user_id)
);

CREATE TABLE comments (
comment_id SERIAL PRIMARY KEY,
content varchar(2000) NOT NULL,
created_time timestamp NOT NULL,
user_id integer NOT NULL,
post_id integer NOT NULL,
parent_comment_id integer,
FOREIGN KEY (user_id) REFERENCES users(user_id),
FOREIGN KEY (post_id) REFERENCES posts(post_id),
FOREIGN KEY (parent_comment_id) REFERENCES comments(comment_id)
)
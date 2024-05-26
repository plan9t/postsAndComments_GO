
# OZON Posts && Comments

Привет! Чтобы запустить сервис следуйте шагам, описаным ниже.

Для работы сервиса необходима внешняя база данных PostgreSQL.

Для запуска PostgreSQL есть команда Docker-compose:
(все параметры для неё уже прописаны)
```bash
docker-compose up -d
```
Также для удобной работы с PostgreSQL в Docker-compose есть [pgadmin](http://localhost:8080/browser/):

```bash
Для авторизации:
user: plan9t
password: plan9t
```

Пока наша база данных пуста, давайте заполним её таблицами и небольшим количеством строк, чтобы в будущем протестировать работоспособность.
Необоходимо выполнить эти SQL-запросы в [pgadmin](http://localhost:8080/browser/):
```bash
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

ALTER TABLE posts ADD COLUMN deleted_at TIMESTAMP WITH TIME ZONE;

INSERT INTO users (firstName, lastName) VALUES ('Artyom', 'Vdovenko');
INSERT INTO users (firstName, lastName) VALUES ('Ivan', 'Alekseev');

INSERT INTO posts (title, content, commentable, created_time, user_id, id) 
VALUES ('Название тестового поста', 'Содержание тестового поста', true, now(), 1, 1);

INSERT INTO comments (content, created_time, user_id, post_id, parent_comment_id) 
VALUES 
('Этот комментарий заряжен на успех! А циклические зависимости полный отстой!!!', now(), 1, 1, NULL),
('Да! Я полностью с тобой согласен! Циклические зависимости это фу(((', now(), 2, 1, NULL),
('Ой, я забыл ответить именно на твой комментарий! Исправляюсь =)', now(), 2, 1, 1),
('А, да ничего страшного! Главное, что мы друг друга поняли! =)', now(), 1, 1, 3);
```

Далее необходимо собрать Docker-образ:
```bash
docker build -t ozon_server . 
```
## Запуск сервиса

Для запуска сервиса используйте следующую команду:

```bash
docker run -p 8090:8090 ozon_server
```

# Проверка работоспособности
Для того чтобы проверить как работают запросы на GraphQL понадобится [Postman](https://www.postman.com/downloads/).
Необходимо выбрать тип POST запрос на URL:
```bash
localhost:8090/graphql
```

## Запросы для проверки:
### Запрос на получение всех постов:
```bash
query {
    posts {
        id
        content
    }
}
```


### Запрос для получения всех комментариев к посту (рекурсивно):
```bash
query {
    post(id: 1) {
        id
        title
        content
        comments {
            id
            content
            createdTime
            children {
                id
                content
                createdTime
            }
        }
    }
}
```


### Запрос постов с пагинацией:
```bash
query {
    posts {
        id
        title
        comments(limit: 2, offset: 2) {
            id
            content
            user
        }
    }
}
```

### Запрос на изменение возможности оставлять комментарии под постом:
```bash
mutation {
  updatePostCommentable(id: 1, commentable: false) {
    id
    title
    commentable
  }
}
```



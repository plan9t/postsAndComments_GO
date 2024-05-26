## Запрос на получение всех постов:
```bash
query {
    posts {
        id
        content
    }
}
```


## Запрос для получения всех комментариев к посту (рекурсивно)
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


## Запрос постов с пагинацией
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

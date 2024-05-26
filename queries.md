Запрос для получения всех комментариев к посту (рекурсивно)

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

Запрос постов с пагинацией
{
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
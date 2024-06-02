# Тестовое задание Ozon Fintech 2024

# Сервис работы с постами и комментариями на Go,PostgreSQL,GraphQL

##### Автор: [Виноградов Данил](https://t.me/japsty) 
##### Ссылка на задание [тут](https://github.com/Japsty/ozon_fintch_testtask_2024/blob/main/task.md)

## Содержание
1. [Запуск](#запуск)
2. [Доступные методы](#доступные-методы)

### Запуск

Запуск из корневой директории проекта, выбор хранилища при помощи изменения параметра STORAGE(postgres/inmemory)
```
STORAGE=postgres PORT=8081 docker-compose up
```

Тесты только на успешное выполнение, корнер кейсы не обработаны.

Билд в два этапа, чтобы не тянуть весь проект в образ.

В ветке inefficent первая версия, в которой используется рекурсивный алгоритм для построения вложенности.

## Доступные методы

### **QUERY**
#### **Posts**
Вывод всех доступных постов.

Чтобы смотреть вложенные комментарии, требуется запрашивать внутри комментария слайс Replies, содержащий дочернии комментарии
Для ускорения поиска в базе использовались hash индексы для id и parent_id комментариев.

*Пример ниже показывает принимаемую на вход и выходящую структуры*

*Входящая структура*
```
query Posts {
    posts {
        id
        title
        content
        userID
        commentsAllowed
        createdAt
        comments {
            id
            content
            authorID
            postID
            parentID
            createdAt
            replies {
                id
                content
                authorID
                postID
                parentID
                createdAt
            }
        }
    }
}
```
*Выходящая структура*
```json
{
    "data": {
        "posts": [
            {
                "id": "87879be2-a8a4-4669-bf1f-4554d44f4a75",
                "title": "blablapost",
                "content": "blablacontent",
                "userID": "5594a70f-ad01-427e-be8a-43bf94fc76fd",
                "commentsAllowed": true,
                "createdAt": "2024-05-31 12:59:32.937707 +0000 +0000",
                "comments": [
                    {
                        "id": "f35dce4e-14d5-402c-9a74-44e425f1fcfc",
                        "content": "5594a70f-ad01-427e-be8a-43bf94fc76fd",
                        "authorID": "blablabla",
                        "postID": "87879be2-a8a4-4669-bf1f-4554d44f4a75",
                        "parentID": null,
                        "createdAt": "2024-05-31 13:15:09.333932 +0000 +0000",
                        "replies": [
                            {
                                "id": "541d7ab0-86c7-44e0-a8ba-bc2df54ad11d",
                                "content": "5594a70f-ad01-427e-be8a-43bf94fc76fd",
                                "authorID": "blabla2!",
                                "postID": "87879be2-a8a4-4669-bf1f-4554d44f4a75",
                                "parentID": "f35dce4e-14d5-402c-9a74-44e425f1fcfc",
                                "createdAt": "2024-05-31 13:15:39.772892 +0000 +0000"
                            }
                        ]
                    }
                ]
            }
        ]
    }
}
```

#### **Post**
Вывод поста по id с пагинацией по комментариям, offset и limit должны быть > 0.

Как и в запросе вывода всех постов степень вложенности комментариев определяется количеством запросов слайса Replies

*Пример ниже показывает принимаемую на вход и выходящую структуры*

*Входящая структура*
```
query Post {
    post(id: "87879be2-a8a4-4669-bf1f-4554d44f4a75", limit: 10, offset: 0) {
        id
        title
        content
        userID
        commentsAllowed
        createdAt
        comments {
            id
            content
            authorID
            postID
            parentID
            createdAt
            replies {
                id
                content
                authorID
                postID
                parentID
                createdAt
            }
        }
    }
}
```
*Выходящая структура*
```json
{
    "data": {
        "post": {
            "id": "87879be2-a8a4-4669-bf1f-4554d44f4a75",
            "title": "blablapost",
            "content": "blablacontent",
            "userID": "5594a70f-ad01-427e-be8a-43bf94fc76fd",
            "commentsAllowed": true,
            "createdAt": "2024-05-31 12:59:32.937707 +0000 +0000",
            "comments": [
                {
                    "id": "f35dce4e-14d5-402c-9a74-44e425f1fcfc",
                    "content": "5594a70f-ad01-427e-be8a-43bf94fc76fd",
                    "authorID": "blablabla",
                    "postID": "87879be2-a8a4-4669-bf1f-4554d44f4a75",
                    "parentID": null,
                    "createdAt": "2024-05-31 13:15:09.333932 +0000 +0000",
                    "replies": [
                        {
                            "id": "541d7ab0-86c7-44e0-a8ba-bc2df54ad11d",
                            "content": "5594a70f-ad01-427e-be8a-43bf94fc76fd",
                            "authorID": "blabla2!",
                            "postID": "87879be2-a8a4-4669-bf1f-4554d44f4a75",
                            "parentID": "f35dce4e-14d5-402c-9a74-44e425f1fcfc",
                            "createdAt": "2024-05-31 13:15:39.772892 +0000 +0000"
                        }
                    ]
                }
            ]
        }
    }
}
```

### **MUTATION**
Для использлования этих методов требуется токен авторизации, передается он в хедере в поле Authorization, имеет формат:

Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwidXNlcl9pZCI6IjU1OTRhNzBmLWFkMDEtNDI3ZS1iZThhLTQzYmY5NGZjNzZmZCIsImlhdCI6MTcwMDAwMDAwMX0.YfWAczXjj24YTu0voVRXxMm5byZlNDFzWzvki4XBnH8

Формат содержимого токена
```
{
  "sub": "1234567890",
  "user_id": "5594a70f-ad01-427e-be8a-43bf94fc76fd",
  "iat": 1700000001
}
```

#### **AddPost**
Создание поста

В параметрах запроса должны передоваться параметры title,content,commentsAllowed.

*Пример ниже показывает принимаемую на вход и выходящую структуры*

*Входящая структура*

```
mutation AddPost {
    addPost(
        post: { title: "blablapost2", content: "blablacontent2", commentsAllowed: true }
    ) {
        id
        title
        content
        userID
        commentsAllowed
        createdAt
        comments {
            id
            content
            authorID
            postID
            parentID
            createdAt
        }
    }
}
```
*Выходящая структура*

```json
{
    "data": {
        "addPost": {
            "id": "7b905a06-d3f6-4d67-b876-61862cd4eacd",
            "title": "blablapost2",
            "content": "blablacontent2",
            "userID": "5594a70f-ad01-427e-be8a-43bf94fc76fd",
            "commentsAllowed": true,
            "createdAt": "2024-05-31 13:24:11.712037 +0000 +0000",
            "comments": []
        }
    }
}
```

#### **AddComment**
Добавление комментария

В параметрах запроса должны передоваться параметры postID,parentID(необязателен),content.

Если комментарий будет реплаем на другой, то требуется указать в parentID id того комментария, реплаем которого будет созданный.

*Пример ниже показывает принимаемую на вход и выходящую структуры*

*Входящая структура*

```
mutation AddComment {
    addComment(
        comment: {
            postID: "7b905a06-d3f6-4d67-b876-61862cd4eacd"
            content: "blablacontentcomm"
        }
    ) {
        id
        title
        content
        userID
        commentsAllowed
        createdAt
        comments {
            id
            content
            authorID
            postID
            parentID
            createdAt
        }
    }
}
```

*Возвращаемая структура*

```
{
    "data": {
        "addComment": {
            "id": "7b905a06-d3f6-4d67-b876-61862cd4eacd",
            "title": "blablapost2",
            "content": "blablacontent2",
            "userID": "5594a70f-ad01-427e-be8a-43bf94fc76fd",
            "commentsAllowed": true,
            "createdAt": "2024-05-31 13:24:11.712037 +0000 +0000",
            "comments": [
                {
                    "id": "830a98a7-4b3f-4eed-a7a8-2c417a24ed1e",
                    "content": "5594a70f-ad01-427e-be8a-43bf94fc76fd",
                    "authorID": "blablacontentcomm",
                    "postID": "7b905a06-d3f6-4d67-b876-61862cd4eacd",
                    "parentID": null,
                    "createdAt": "2024-05-31 13:26:00.291665 +0000 +0000"
                }
            ]
        }
    }
}
```

#### **ToggleComments**
Включение/выключение комментариев на посте

В параметрах запроса должны передоваться параметры postID и allowed. Данное действие доступно только тому пользователю, чей userID == AuthorID поста, т.е. только автору.

Если попытаться добавить коммент на пост, где комментарии выключены, то вернется ошибка "comments not allowed".

*Пример ниже показывает принимаемую на вход и выходящую структуры*

*Входящая структура*

```
mutation ToggleComments {
    toggleComments(postId: "7b905a06-d3f6-4d67-b876-61862cd4eacd", allowed: true) {
        id
        title
        content
        userID
        commentsAllowed
        createdAt
        comments {
            id
            content
            authorID
            postID
            parentID
            createdAt
        }
    }
}

```

*Возвращаемая структура*

```json
{
    "data": {
        "toggleComments": {
            "id": "7b905a06-d3f6-4d67-b876-61862cd4eacd",
            "title": "blablapost2",
            "content": "blablacontent2",
            "userID": "5594a70f-ad01-427e-be8a-43bf94fc76fd",
            "commentsAllowed": true,
            "createdAt": "2024-05-31 13:24:11.712037 +0000 +0000",
            "comments": []
        }
    }
}
```

### **Subscriptions**

#### **commentAdded**

Подписка на комментарии поста

Требуется передать postID поста на который нужно подписаться.

*Входящая структура*

```
subscription CommentAdded {
    commentAdded(postId: "1ed26baf-d822-433a-98f0-96d7f12bc0ea") {
        id
        content
        authorID
        postID
        parentID
        createdAt
        replies {
            id
            content
            authorID
            postID
            parentID
            createdAt
        }
    }
}
```

*Возвращаемая структура*

```json
{
    "data": {
        "commentAdded": {
            "id": "365442bb-3293-43f0-9e0a-4841139a6f75",
            "content": "afsfsaf",
            "authorID": "5592a70f-ad01-427e-be8a-43bf94fc76fd",
            "postID": "1ed26baf-d822-433a-98f0-96d7f12bc0ea",
            "parentID": null,
            "createdAt": "2024-06-02 22:52:01.932939 +0000 +0000",
            "replies": null
        }
    }
}
```

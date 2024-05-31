# Тестовое задание Ozon Fintech 2024

# Сервис работы с постами и комментариями на Go,PostgreSQL,GraphQL

##### Автор: [Виноградов Данил](https://t.me/japsty) 
##### Ссылка на задание [тут](https://github.com/Japsty/ozon_fintch_testtask_2024/blob/main/task.md)

## Содержание
1. [Запуск](#запуск)
2. [Доступные методы](#доступные-методы)
3. [Что не успел](#что-не-успел)

### Запуск

Запуск из корневой директории проекта

```shell
  docker-compose up
```

Выбор хранилища находится в docker-compose файле

Билд в два этапа, чтобы не тянуть весь проект в образ.

## Доступные методы

### **QUERY**
#### **Posts**
Вывод всех доступных постов.

Чтобы смотреть вложенные комментарии, требуется запрашивать внутри комментария слайс Replies, содержащий дочернии комментарии
Комментарии находятся рекурсивно, возможно, это не самый оптимальный вариант. Для ускорения поиска в базе использовались hash индексы
для id и parent_id комментариев

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
Вывод поста по id с пагинацией по комментариям.

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
Токен нигде не генерируется, не успел реализовать сервис польователя(

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

### Что не успел

Для меня работа с GraphQL была в новинку, потому пришлось немного разобраться, поэтому не успел реализовать subscriptions, но будь у меня побольше времени, то обязательно бы сделал.

Также не получилось вовремя покрыть тестами часть resolver'ов и обработать корнер кейсы в теста репозиториев.

Не успел написать тесты на inmem репозиторий. Некоторый функционал, что есть в postgres отсутствует в inmem.

Хотел развернуть сервис на своем сервере через gitlab ci/cd.

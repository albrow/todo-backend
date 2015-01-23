A RESTful JSON Backend for TodoMVC
----------------------------------

### Installation

1. Install go version >= 1.0. ([Latest version](https://golang.org/dl/) is recommended).
2. Make sure you have followed [these instructions](https://golang.org/doc/code.html) for setting up your go workspace.
3. Run `go get -u github.com/albrow/todo-backend`
4. Change into $GOPATH/src/github.com/albrow/todo-backend
5. Run `go run main.go`

### About

The server is written in go and runs on port 3000. It accepts Content-Types of
application/json, application/x-www-form-urlencoded, or multipart/form-data and
responds with JSON. It does validations and returns 422 errors when validations
fail.

This backend simply stores todos in memory so you don't have to worry about
setting up a database. This means there is no persistence, so if you restart
the server all your todos will be gone. It's perfect for testing and building
new things, but not recommended for use in production.

### Endpoints

#### GET /todos

List all existing todos, ordered by time of creation.

**Parameters**: none

**Example Responses**:

Success:

```json
[
  {
    "id": 0,
    "title": "Write a frontend framework in Go",
    "isCompleted": false
  },
  {
    "id": 1,
    "title": "???",
    "isCompleted": false
  },
  {
    "id": 2,
    "title": "Profit!",
    "isCompleted": false
  }
]
```

#### POST /todos

Create a new todo item.

**Parameters**:

| Field    | Type    | Description     |
| ---------| ------- | --------------- |
| title    | string  | The title of the new todo. |


**Example Responses**:

Success:

```json
{
  "id": 3,
  "title": "Take out the trash",
  "isCompleted": false
}
```

Validation error:

```json
{
  "title": [
    "title is required."
  ]
}
```

#### PUT /todos/{id}

Edit an existing todo item.

**Parameters**:

| Field       | Type    | Description     |
| ----------- | ------- | --------------- |
| title       | string  | The title of the todo. |
| isCompleted | bool    | Whether or not the todo has been completed. |

**Example Responses**:

Success:

```json
{
  "id": 3,
  "title": "Handle the garbage",
  "isCompleted": false
}
```

#### DELETE /todos/{id}

Delete an existing todo item.

**Parameters**: none

**Example Responses**:

Success:

```json
{}
```

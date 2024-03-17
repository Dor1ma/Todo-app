# Todo-app

Simple educational project to understand basics of 
REST API

## Technologies:
- Golang
- Gorilla mux (web framework)
- PostgreSQL
- Docker
- Postman

## Structure

The server consists of the following main components:

- `server`: the main server object, which contains the router, repository, and session store.
- `router`: used for configuring and handling HTTP requests.
- `repository`: used for interacting with the database.
- `sessions`: used for managing user sessions.

## Routes

The server provides the following routes:

- `/users`: create a new user (POST).
- `/sessions`: create a new session (POST).
- `/private/whoami`: get information about the current user (GET).
- `/private/todos`: create a new todo list (POST), get all todo lists (GET).
- `/private/todos/{id}`: update a todo list (PUT), delete a todo list (DELETE), get a todo list by ID (GET).
- `/private/todos/{id}/items`: get all items of a todo list (GET), create a new item in a todo list (POST).
- `/private/todos/{id}/items/{id}`: get an item by ID from a todo list (GET), update an item in a todo list (PUT), delete an item from a todo list (DELETE).
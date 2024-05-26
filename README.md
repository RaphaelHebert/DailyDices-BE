# DailyDices-BE

Backend for the Daily Dices App

## Overview

### Stack

- go fiber

### Endpoints

#### HTML

| HTTP Method | URL                              | Description |
| ----------- | -------------------------------- | ----------- |
| `GET`       | http://localhost:8000/roll-dices | Get dices   |
| `POST`      | http://localhost:8000/login      | login user  |
| `GET`       | http://localhost:8000/token      | renew token |
| `GET`       | http://localhost:8000/score      | user scores |

#### User Service

| HTTP Method | URL                                  | Description       |
| ----------- | ------------------------------------ | ----------------- |
| `POST`      | http://localhost:8000/users          | Create new User   |
| `PUT`       | http://localhost:8000/users/{userId} | Update User by ID |
| `GET`       | http://localhost:8000/users/{userId} | Get User by ID    |
| `DELETE`    | http://localhost:8000/users/{userId} | Delete User by ID |

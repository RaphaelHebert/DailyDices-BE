# DailyDices-BE

Backend for the Daily Dices App

## Overview

### Stack

- go fiber

### Endpoints

#### HTML

| HTTP Method | URL                              | Description           |
| ----------- | -------------------------------- | --------------------- |
| `GET`       | http://localhost:8000/roll-dices | Get dices protected   |
| `POST`      | http://localhost:8000/login      | login user            |
| `GET`       | http://localhost:8000/token      | renew token protected |

#### User Service

| HTTP Method | URL                                  | Description                       |
| ----------- | ------------------------------------ | --------------------------------- |
| `POST`      | http://localhost:8000/users          | Create new User                   |
| `PUT`       | http://localhost:8000/users/{userId} | Update User by ID protected       |
| `GET`       | http://localhost:8000/users/ALL      | Get All users protected           |
| `GET`       | http://localhost:8000/users/{userId} | Get User by ID protected          |
| `DELETE`    | http://localhost:8000/users/{userId} | Delete User by ID protected admin |
| `DELETE`    | http://localhost:8000/users/delete   | Delete own account protected      |

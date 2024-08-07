# Mongodb init for go project

## Overview

- Simple golang user management project to demonstrate how to create database, collections and indexes on startup.
- The project is written in Golang using JWT to authenticate and authorize user, Gin for router, Go Kit architecture, and MongoDB as the database.

## How the code run

1. Load environment variables in .env file
   - `PORT`: server port
   - `SECRET`: JWT secret key
   - `MONGODB_URL`: database connect url
2. Connect to mongodb
3. Create database, colletions and indexes by calling initilizer.InitializeDatabase() method. This method do the following:
   - Create database
   - Create collection with schema validator (if not exist)
   - Create index (if not exist)
4. Create router and run server to handle request

## API
### Register user
- cURL:
  ```
  curl --location 'http://localhost:3000/auth/signup' \
  --header 'Content-Type: application/json' \
  --data-raw '{
      "email" : "tri1234@gmail.com",
      "password" : "123456",
      "firstName" : "Tri",
      "lastName" : "Nguyen Tran",
      "phoneNumber": "123456789",
      "age": 20
  }'
  ```
- `email` and `password` are required in request body
### Login
- cURL:
  ```
  curl --location 'http://localhost:3000/auth/login' \
  --header 'Content-Type: application/json' \
  --data-raw '{
      "email":"tri1234@gmail.com",
      "password":"123456"
  }'
  ```
- A jwt token will be returned in response body
### Update user
- cURL:
  ```
  curl --location --request PATCH 'http://localhost:3000/user/update' \
  --header 'Content-Type: application/json' \
  --header 'Authorization: Bearer ••••••' \
  --data '{
      "password" : "123456",
      "firstName" : "Tri",
      "lastName" : "Nguyen Tran",
      "phoneNumber": "123456789",
      "age": 20
  }'
  ```
- JWT token is required in request header
### Delete user
- cURL:
  ```
  curl --location --request DELETE 'http://localhost:3000/user/delete?userId=66a85f7b18a3324a003a12de' \
  --header 'Authorization: ••••••'
- JWT token of an account with admin role is requires in request header

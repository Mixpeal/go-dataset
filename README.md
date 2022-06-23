## Endpoints

- GET /api - _USER API_
    - GET /users - _Create all users_
    - POST /users - _Create a user_
    - GET /users/:id - _Get a user_
    - PATCH /user/:id - _Update a user_
    - DELETE /user/:id - _Delete user_

## Filtering Results
- Filtering users by email (or any other unique identifier)
    - GET /users?filters=["email", "hanna.lance@gmail.com"] - _You can replace email property with any other one e.g. name, company, etc._
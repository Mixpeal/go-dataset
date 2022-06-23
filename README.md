## Endpoints

- GET /api - _USER API_
    - GET /users - _Create all users_
    - POST /users - _Create a user_
    - GET /users/:id - _Get a user_
    - PATCH /user/:id - _Update a user_
    - DELETE /user/:id - _Delete user_

## Filtering Results
- Filtering users by email (or any other unique identifier)
    - GET /users?filters=["email", "hanna.lance@gmail.com"] (_You can replace email property with any other one e.g. name, company, etc._)
    
- Filtering users by a date range
    - GET /users?filters=[["date",">=","2022-08-15"],["AND"],["date","<=","2023-01-15"]] (_You can replace email property with any other one e.g. name, company, etc._)

- Filtering users by by page and setting page limit
    - GET /users?page=1&size=2 (_page sets the current page and size sets the number of users shown in a page._)

- The user can even be sorted
    - GET /users?size=10&sort=-name,id (_adding - sign sorts in descending order while without - sorts in ascending order._)
### API

```
/api
    /users
        GET - current user
        PUT - create user
        POST - update password
        DELETE - delete account
    /documents
        GET - document list
        PUT - create document
        /:id
            GET - get document metadata
            POST - update document
            DELETE - delete document
    /persons
        GET - persons list
        PUT - create person
        /:id
            GET - get person
            POST - update person
            DELETE - delete person
    /auth
        /user
            GET - get user
        /session
            PUT - create session
            DELETE - delete session
```
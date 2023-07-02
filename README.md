# Employee-service

## Packaging and running

Commands below must be executed from the project root

### Test

```bash
  make test
```

### Build and Run

```bash
  make run
```
or
```bash
  docker-compose up --build -d
```

### Stop
```bash
  make stop
```
or
```bash
  docker-compose down
```

## Access services after running
### Employee-service

`http://localhost:8080`

### MailHog ui

`http://localhost:8025`

## Employee-service endpoints

* /login
  * used for employee authentication. The endpoint should accept POST requests with 'username' and 'password' and, if authenticated successfully, return a JWT. 
* /employee
  * used for querying and mutating the 'Employee' table. It should only accept requests with a valid JWT in the Authorisation header. The requests should be in the GraphQL format, and the endpoint should perform the requested query/mutation on the 'Employee' table

## Postman collection

there is a postman collection called `employee-service.postman_collection` 
in the root folder of the project

there is one user in the db that can make login requests:
```
username: highlander
password: whoWantsToLiveForEver
```

## Issues or unfinished parts

* logs - the project lacks meaningful logs
* validations - some validations are missing
  * email
  * string sizes (username, firstname, lastname, etc)
  * valid date of birth
  * password format
  * and others
* email sent for new employee is not appropriately branded and professional in appearance, missing templating
* list of all employees is not filtering all possible fields
* not all returned errors are meaningful

## Possible improvements

* mail sending would be better through a queue and an async function
* jwt token is being held by the service, in a scalable env it should be saved in cache or database
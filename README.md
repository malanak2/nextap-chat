# Run requirements

## Env variables

- DB
  - DB_host
  - DB_port
  - DB_user
  - DB_pass
  - DB_name
- Program
  - port
  - JWT_SECRET

# Endpoints quick commands

Create User
```shell
curl --header "Content-Type: application/json" --request POST --data "{\"Username\":\"username\",\"password\":\"password\"}" http://localhost:8080/create-user
```

Login
```shell
curl --header "Content-Type: application/json" --request POST --data "{\"Username\":\"username\",\"password\":\"password\"}" http://localhost:8080/login
```

Send message
```shell
curl --header "Content-Type: application/json" --header "Authorization: Bearer token" --request POST --data  http://localhost:8080/sendMessage
```



# 
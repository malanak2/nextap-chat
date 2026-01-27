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

## Before building

### Generate db types
Make sure the db is up and up to date!
```shell
go install github.com/go-jet/jet/v2/cmd/jet@latest
jet -dsn=postgresql://user:pass@localhost:5432/jetdb?sslmode=disable -schema=chatdb -path=./gen
```

### Generate docs

```shell
go install github.com/swaggo/swag/cmd/swag@latest
swag init
```

## TODO

- In docker verify gen exists?
- Clean up code + comments

## Dev
- Ports -> Database
- Usecases -> email, notifikace, obecne proste akce
- Handlers -> Akorat handlery, convert na system wide veci, validace etc., volani usecases
- Json
  -  https://go.dev/blog/slog


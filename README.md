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


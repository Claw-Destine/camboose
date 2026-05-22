# Instruction for developers

## Webapp templ generation and dev mode

```sh
go tool templ generate --watch --proxy="http://localhost:3000" --cmd="go run ./webapp"
```

## Running tests

```sh
go test ./webapp/... ./core/...
```

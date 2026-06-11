# Instruction for developers

## To run in development mode

Prepare env variables:

```sh
cp .env.example .env
```

To build frontend:

```sh
cd frontend
npm run watch
```

To run backend:

```sh
 go run ./webapp/ #...or setup debuger with cwd in root project directory
```

The backend will serve the `frontend/dist` directory at localhost:3000

## Running tests

```sh
go test ./webapp/... ./core/...
```

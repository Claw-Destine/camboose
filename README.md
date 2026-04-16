# camboose

Recipe driven software management AI workflow orchestrator

## Project structure

```sh
|- recipies # Example recipies
|- service # GoLang backend
|- web # Vue frontend
```

## Frontend

The frontend is a Vite + Vue application in `web/`.

```sh
cd web
npm install
npm run dev
```

The Vite dev server proxies `/api` requests to the Go service on `http://localhost:8080`.

## Backend

The project is using Go workspace. You can start the backebd directly from the root folder:

```sh
go run ./service
```

## 3rd Pary software

All the services necessary for a local deployment are in `docker-compose.yml`.

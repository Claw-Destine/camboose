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
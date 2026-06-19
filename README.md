# camboose

<img src="static/assets/camboose_logo.png" alt="logo" width="200"/>

Recipe driven software development AI workflow orchestrator

## Project structure

```sh
|- recipies # Example recipies
|- core # GoLang - application logic
|- webapp # HTMX fragment renderer
|- frontend # Webpack project with web components
```

## 3rd Pary software

All the services necessary for a local deployment are in `docker-compose.yml`.

## Formatting

From the repository root:

```sh
./scripts/format.sh
```

To verify formatting without rewriting files:

```sh
./scripts/format-check.sh
```

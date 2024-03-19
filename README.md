## Installation and prerequisites

- Go 1.22 or higher
- install [Node.js](https://nodejs.org/en/download/)
- install [Make](https://www.gnu.org/software/make/) tool
- download [tailwindcss CLI](https://tailwindcss.com/docs/installation), rename it to "tailwindcss" and move it to the root of the project
- install air for hot reloading by running `go install github.com/cosmtrek/air@latest` in your terminal (Go 1.22 or higher is required)

## Run in dev mode

### on 3 separate terminals(Powershell) run these commands in order respectively:

```bash
make tw
```

```bash
make dev
```

```bash
make templ-proxy
```

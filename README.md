# finance-tracker-backend

This is a backend built with Go for the FinanceTracker project. It's a RESTful API that provides endpoints for the frontends to consume.

## Tools

- [**Air**](https://github.com/air-verse/air): live reload tool that reloads the server whenever changes are made to the code
  - Change settings in the `.air.toml` if needed.
  - Run `air` to start the server.
- [**Chi**](https://github.com/go-chi/chi): lightweight, idiomatic and composable router for building Go HTTP services
  - Check [documentation](https://pkg.go.dev/github.com/go-chi/chi/v5) and [examples](https://github.com/go-chi/chi/tree/master/_examples)
- [**Jet**](https://github.com/go-jet/jet): a complete solution for efficient and high performance database access, consisting of type-safe SQL builder with code generation and automatic query result data mapping.
# CI6NDEX Development Guide

## Build Commands
- Build project: `make build`
- Run project: `make run`
- Generate database models: `make generate`
- Clean generated code: `make clean`
- Update dependencies: `make update`
- Sync Discord commands: `make sync`
- Run tests: `go test -v ./...`
- Run specific test: `go test -v ./path/to/package -run TestName`
- Lint code: `golangci-lint run --config .golangci.yml`

## Code Style Guidelines
- **Formatting**: Use `gofmt` for consistent formatting
- **Imports**: Group imports in blocks - standard library first, then external packages
- **Error Handling**: Use `github.com/pkg/errors` for wrapping errors with context
- **Logging**: Use `log/slog` for structured logging
- **Naming Conventions**:
  - Use camelCase for variable/function names
  - Use PascalCase for exported names
  - Use snake_case for database fields
- **Error Names**: Error variables should be prefixed with `Err` or `err`
- **Types**: Prefer explicit types over interface{} or any
- **Dependencies**: Use `go.mod` for dependency management
- **Documentation**: Document all exported functions and types

## Project Structure
- `bot/`: Discord bot implementation
- `ci6ndex/`: Core business logic and database operations
- `cmd/`: Command-line interface
- `sql/`: SQL queries and migrations
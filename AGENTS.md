## Build, Lint, and Test

- **Build:** `make build` (includes format and lint)
- **Lint:** `make lint`
- **Test All:** `make test` (run after all work is finished)
- **Test Kohan:** `ginkgo -r components/kohan/`
- **Run targeted tests:** `ginkgo -r <package>` after a code change in that package.
- **Format code:** `make format`

** Note: After every code change run Build and Targeted Tests **

## Code Style (via `make format` and `make lint`)

- **Imports:** Grouped and formatted by `make format`.
- **Formatting:** Enforced by `make format`.
- **Types:** Use built-in types when possible. Define custom types for complex data structures.
- **Naming Conventions:**
    - Use camelCase for variables and functions.
    - Use PascalCase for public functions and types.
    - Keep names short and descriptive.
- **Error Handling:**
    - Use `error` type for error handling.
    - Check for errors immediately after a function call.
    - Use `defer` for cleanup.
- **Comments:** Add comments to explain complex logic.
- **Complexity:** Keep functions short and focused. Max 40 lines and 20 statements per function.
- **Line Length:** Max 200 characters per line.

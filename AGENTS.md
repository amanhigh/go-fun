## Build, Lint, and Test

- **Build:** `make OUT=/dev/stdout build` (includes format and lint)
- **Lint:** `make lint`
- **Test All:** `make OUT=/dev/stdout test` (run after all work is finished)
- **Test Kohan:** `ginkgo -r components/kohan/`
- **Targeted tests:** `ginkgo -r <package>` after a code change in that package.
- **Format code:** `make format`

* When Build or Test Fails do not `OUT=/dev/stdout` is not cause of that it is Supported in Root Make File.
Without this you won't be able to see details. Check for actual Error Instead.*
** IMPORTANT Note: After every code change run Build and Targeted Tests **

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

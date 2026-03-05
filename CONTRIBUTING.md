# Contributing to PromptVault

Thank you for your interest in contributing to PromptVault!

## Local Development Setup

1. Ensure you have Go 1.20 or later installed.
2. Clone the repository:
   ```sh
   git clone https://github.com/Bharath-code/promptvault.git
   cd promptvault
   ```
3. Install dependencies:
   ```sh
   make deps
   ```
4. Run the TUI locally:
   ```sh
   make run
   ```

## Pull Request Guidelines

1. **Create a branch** for your feature or bug fix.
2. **Write tests** for any new features or bug fixes. Run `make test` to ensure everything passes.
3. **Format code** properly. `golangci-lint` will be run via GitHub Actions.
4. Update the `README.md` if your change impacts how the CLI or TUI is used.
5. Create an informative Pull Request describing your changes.

## Architecture Overview
- `internal/db`: SQLite database interactions with FTS5 search
- `internal/model`: Core structs representing a Prompt
- `internal/tui`: Bubble Tea components for the CLI interface
- `internal/cmd`: Cobra commands 

Once again, thank you for making PromptVault better!

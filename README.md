# ci6ndex

### What is it?

This repository is a tool for managing our weekly civ 6 games. Features include:
* Discord bot to manage drafts and roll civs
* CLI commands for bot management (`serve` and `sync`)

## Libraries and Tools Used

- **Kong**: Command-line interface parser and handler.
- **env**: A library for environment-based configuration management.
- **SQLite**: The database engine used for storing data.
- **goose**: Database migration tool.
- **Disgo**: A Discord library for Go.
- **charmbracelet/log**: Logging library.

## Local development

### Needed locally:
* .env
  * `DISCORD_API_TOKEN`
  * `DISCORD_BOT_APPLICATION_ID`
  * `GUILD_IDS`

### How do I get set up?

`make` will provide a list of commands to run.

### CLI Commands

The application provides a CLI with the following commands:
- `./civ bot serve` - Start the bot server
- `./civ bot sync` - Synchronize bot data

## TODO:
* https://github.com/hazzardr/ci6ndex/issues
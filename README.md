# ci6ndex

### What is it?

This repository is a tool for managing our weekly civ 6 games. Features include:
* Discord bot to manage drafts and roll civs

## Libraries and Tools Used

- **Go**: The main programming language used for the project.
- **viper**: A library for creating powerful modern CLI applications.
- **sqlc**: A SQL compiler for generating type-safe Go code from SQL queries.
- **SQLite**: The database engine used for storing data.
- **Disgo**: A Discord library for Go.

## Local development

### Needed locally:
* .env
  * `DISCORD_API_TOKEN`
  * `DISCORD_BOT_APPLICATION_ID`
  * `GUILD_IDS`

### How do I get set up?

`make` will provide a list of commands to run.

## TODO:
* https://github.com/hazzardr/ci6ndex/issues
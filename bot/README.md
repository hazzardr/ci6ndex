# ci6ndex

### What is it?

This repository is a tool for managing our weekly civ 6 games. Features include:
* Google Sheets integration to import rankings from Nick's sheet
* Discord bot to manage drafts and roll civs

### Local development

#### Needed locally:
* gcloud_oauth2.json file for authenticating with gsheets api
* .env
  * `DISCORD_API_TOKEN`
  * `GCLOUD_CREDS_LOC`
  * `POSTGRES_USER`
  * `POSTGRES_PASSWORD`
  * `POSTGRES_DB`
  * `POSTGRES_SCHEMA`

```
make build bot
```

### TODO:
* https://github.com/hazzardr/ci6ndex/issues
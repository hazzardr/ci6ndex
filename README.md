# ci6ndex

## Tools

### APIs
* HTTP server via Gorilla/mux
* Discord Bot via DiscordGo
* Google Sheets Client via Google Sheets API Library

### Bootstrapping / Config
* CLI via Cobra
* Config via Viper
* TODO: GCPSM so i don't have to keep shuffling around secrets

### Storage
* Postgres DB via pgx
* DB Migrations via Atlas

```
source .env
# start postgres
docker-compose up -d
# apply db migrations
atlas schema apply -u "postgresql://${POSTGRES_USER}:${POSTGRES_PASS}@:5432/${POSTGRES_DB}?sslmode=disable" --to file://./db/schema.sql --dev-url "docker://postgres/15/dev"
```

# Needed locally:
* gcloud_oauth2.json file for authenticating with gsheets api
* .env
  * `DISCORD_API_TOKEN`
  * `GCLOUD_CREDS_LOC`
  * `POSTGRES_USER`
  * `POSTGRES_PASSWORD`
  * `POSTGRES_DB`
  * `POSTGRES_SCHEMA`

# go-htmx-data-dashboard

## Videos

This project is a part of a series of videos on the Web Dev Fuel channel.

[Go & HTMX Data Dashboard - Live WebSocket Notifications](https://www.youtube.com/watch?v=DXCJNfgExWk)
[Go & HTMX Data Dashboard - Data Table](https://www.youtube.com/watch?v=oBjsh0A-S6U)
[HTMX & Go Data Dashboard - Charts](https://www.youtube.com/watch?v=tr_MW-y70T0)
[Go & HTMX Data Dashboard - Search](https://www.youtube.com/watch?v=fMwQpH36688)

## Set Up Project

To set up the project, clone the Go project with git.

```bash
git clone https://github.com/webdevfuel/go-htmx-data-dashboard
```

### Air

If you have the [air](https://github.com/air-verse/air) executable, then you can run the command below, and everything will just work.

```bash
air
```

### Manually

If you want to do it manually, you can check the `.air.toml` file and manually run the following commands.

```bash
templ generate
npm run dev
go run .
```

## Set Up Meilisearch

To set up Meilisearch, install it locally with the help of the [documentation](https://www.meilisearch.com/docs/learn/self_hosted/install_meilisearch_locally).

Then, inside the root of the project, export the enviroment variable.

```bash
export MEILISEARCH_HOST=http://localhost:7700
```

## Set Up PostgreSQL

To set up PostgreSQL, install it locally with the help of the [documentation](https://www.postgresql.org/download/).

Then, inside the root of the project, export the enviroment variable.

```bash
export DATABASE_URL=postgres://$USER:$PASSWORD@localhost:5432/$DATABASE?sslmode=disable
```

Also, ensure the database exists. You can use the shell, or a DB client to create the database.

## Seed

The project contains a way to seed both PostgreSQL and Meilisearch. This helps us populate the web application with data.

To run it, run the Go program from within the `cmd/seed` directory.

```bash
cd cmd/seed
go run .
```

If the environment variables are correctly exported, it should work.

## Set Up Chart.js

To be able to render charts on the browser, we're using the Chart.js library.

To install it, run the command below to copy the minified JS file from the CDN to the `static` directory.

```bash
wget https://cdnjs.cloudflare.com/ajax/libs/Chart.js/4.4.1/chart.umd.js -O ./static/chart.min.js
```

## Set Up Static Files

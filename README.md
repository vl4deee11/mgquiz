# Magnus-Quiz

Api docs [here](./etc/queries.txt).

Migrations [here](./db-changelog).

## Database: 
    >= PostgreSQL 13.1

## Configure use env

    MG_PORT=8080 - app http port
    MG_DB_HOST=127.0.0.1 - database host 
    MG_DB_PORT=5432 - database port
    MG_DB_USER=postgres - database username
    MG_DB_PASS=postgres - database password
    MG_DB_NAME=postgres - database name
    MG_KEY=uuid4 - security key
    MG_LOGLVL=debug || info || trace || warn || error - app log level


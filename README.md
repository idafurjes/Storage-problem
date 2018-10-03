I created an API that is able to insert a larger amount of records from a csv file to the database table on POST request and returns the json record on GET reuest given the id.
I stored the records in a postgres database, the table is calles promotions.  

To get a working development environment in Fedora, run the following:
```bash
docker run --name promo-postgres -p 5432:5432 -d postgres

docker exec -it promo-postgres /bin/bash/ -c "createuser -U postgres && createdb --owner postgres -U postgres promodb &&
psql -v ON_ERROR_STOP=1 --username "postgres" --dbname "promodb" <<-EOSQL
    CREATE TABLE promotions (
        id              TEXT NOT NULL UNIQUE,
        price           REAL NOT NULL,
        exp_date        TEXT NOT NULL
);
EOSQL"

Then fetch and run the code:

To load a csv file into the database, run:

```bash
curl -X POST localhost:1321/promotions --data-binary @ids.csv
```

To get a file from the database, run:

```bash
curl localhost:1321/promotions/d018ef0b-dbd9-48f1-ac1a-eb4d90e57118
```

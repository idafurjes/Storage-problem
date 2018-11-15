<h1>API</h1>

I created an API that is able to insert a larger amount of records from a CSV file to a database table on POST request and returns a JSON record on GET request given the ID.
I stored the records in a PostgreSQL database, the table is called "promotions".  

<h3>Setting up test DB</h3>
To set up a test database in Docker, run the following command to start a container:

<code>docker run --name promo-postgres -p 5432:5432 -d postgres</code>

Subsequently, run this additional command to set up the DB:

<code>
docker exec -it promo-postgres /bin/bash/ -c "createuser -U postgres && createdb --owner postgres -U postgres promodb &&
psql -v ON_ERROR_STOP=1 --username "postgres" --dbname "promodb" <<-EOSQL  
    
    CREATE TABLE promotions (  
        id              TEXT NOT NULL UNIQUE,     
        price           REAL NOT NULL,  
        exp_date        TEXT NOT NULL  
    );
    EOSQL"
</code>

<h3>Using the API</h3>

To load a csv file into the database, run from the command line:

<code>curl -X POST localhost:1321/promotions --data-binary @ids.csv</code>


To get a record from the database, run from the command line:

<code>curl localhost:1321/promotions/d018ef0b-dbd9-48f1-ac1a-eb4d90e57118</code>


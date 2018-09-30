To get a working development environment in Fedora, run the following:

```bash
su -
dnf install postgresql-server
/usr/bin/postgresql-setup --initdb
cat > /var/lib/pgsql/data/pg_hba.conf << EOF
local	all	all	trust
host	all	all	127.0.0.1/32	trust
host	all	all	::1/128	trust
EOF
systemctl enable postgresql
systemctl start postgresql
createuser -U postgres pg
createdb --owner pg -U postgres promodb
```

Then fetch and run the code:

```bash
go install -v -u -t github.com/idafurjes/pn-jnr-se
$GOPATH/bin/pn-jnr-se &
./test.sh
```

To load a csv file into the database, run:

```bash
curl -X POST localhost:8080/promotions --data-binary @ids.csv
```

To get a file from the database, run:

```bash
curl localhost:8080/promotions/d018ef0b-dbd9-48f1-ac1a-eb4d90e57118
```


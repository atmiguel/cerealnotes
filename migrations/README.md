# Locally:
1. install & setup postgres
	* `brew install postgres`
	* ``createdb `whoami` ``
2. Run postgres daemon. 
	* `pg_ctl start -D /usr/local/var/postgres`
3. Create cerealnotes databases
	* `psql < tools/createDatabases.sql`
3. Run all the migrations on both "unittest" database (`cerealnotes_test`) and as well as the "live" database (`cerealnotes`). 
	* `psql [DATABASENAME] < [MIGRATION_NAME]`

# On Heroku:

1. `heroku pg:psql < [MIGRATION_NAME]`
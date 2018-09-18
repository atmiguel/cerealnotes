# Locally:
1. first make sure that postgres is running. 
	* If installed via homebrew on a macOS: `pg_ctl start -D /usr/local/var/postgres`
2. then make sure you have the neccisary databases
	* `psql < tools/createDatabases.sql`
3. then run all the migrations on both unittest database (`cerealnotes_test`) and as well as the "live" (`cerealnotes`)database. 
	* `psql [DATABASENAME] < [MIGRATION_NAME]`

# On Heroku:

1. `heroku pg:psql < [MIGRATION_NAME]`

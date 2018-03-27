# Locally: 
1. first make sure that postgres is running. 
	* If installed via homebrew on a macOS: `pg_ctl -D /usr/local/var/postgres start`
2. then run migration locally
	* `psql < *MIGRATION_NAME*`

# On Heroku: 

1. `heroku pg:psql < *MIGRATION_NAME*`

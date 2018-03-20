Locally: 
1. first make sure that postgres is running. 
	* On mac via homebrew: `pg_ctl -D /usr/local/var/postgres start`
2. then run migration locally
	* `psql < *MIGRATION_NAME*`

On Heroku: 

`heroku pg:psql < *MIGRATION_NAME*`

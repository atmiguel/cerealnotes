# Installation
## Local prod like build
* install docker
* cd this repo
* `docker-compose up`
* Visit `localhost:8080/`

## Local Dev instance
1. make sure all go depenecies are ready
	* run `./beam_me_up_scotty.sh bash`
	* run `dep ensure`
	* exit bash
2. prepare test db
	* run `./beam_me_up_scotty.sh bash`
	* run `psql $DATABASE_URL < /docker-entrypoint-initdb.d/*.sql`
	* run `psql $DATABASE_URL_TEST < /docker-entrypoint-initdb.d/*.sql`
3. run `./beam_me_up_scotty.sh` and everything should work

Please note step 1 and step 2 are only necessary the first time you connect set up your test environment

# Bugs
* Running sql migrations is very cumbersome

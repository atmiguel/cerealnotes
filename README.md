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
	* run `./beam_me_up_scotty.sh db`
	* run `\c test` (connect to test database)
	* copy contents of 0000_createDbs.sql and paste into terminal
	* exit bash
3. run `./beam_me_up_scotty.sh` and everything should work

Please note step 1 and step 2 are only necciary the first time you connect set up your test environment

# Bugs
* dep doesn't work
* tables ahve not been migrated

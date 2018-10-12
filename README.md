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
2. run `./beam_me_up_scotty.sh` and everything should work

Please note step 1 is only necessary the first time you connect set up your test environment, or whenever a new package is added

#! /bin/bash

GREEN='\033[0;32m'
RED='\033[0;31m'
NC='\033[0m' # No Color

exit_script() {
	echo -e "${GREEN}Shutting down beam, sir!${NC}"
}

trap exit_script SIGINT SIGTERM

docker-compose --file docker-compose.dev.yml up -d;

if [ "$1" = "db" ]; then
	echo -e "${GREEN}Beaming you into the system mechanics, sir!${NC}"
	docker exec -it cerealnotes_db_1 psql -U docker -W docker
elif [ "$1" = "bash" ]; then
	echo -e "${GREEN}Beaming you straight into quantum space, sir!${NC}"
	docker exec -it cerealnotes_backend_1 /bin/bash;
else
	echo -e "${GREEN}Validating tele-ions and then the beam will start, sir!${NC}"
	echo "Running tests then staring the service"
	docker exec cerealnotes_backend_1 bash -c 'go test ./... && go run main.go';
fi

echo -e "${GREEN}Quantum dientagiling reactor subfluid, for next beam, sir!${NC}"
docker-compose --file docker-compose.dev.yml down;

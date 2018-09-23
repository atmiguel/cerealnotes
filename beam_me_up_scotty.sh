#! /bin/bash

exit_script() {
    echo "Scotty: Preparing for next beam"
    docker-compose --file docker-compose.dev.yml down;
}

trap exit_script SIGINT SIGTERM

docker-compose --file docker-compose.dev.yml up -d;
if [ "$1" = "bash" ]; then
	docker exec -it cerealnotes_backend_1 /bin/bash;
	docker-compose --file docker-compose.dev.yml down;
else
	docker exec cerealnotes_backend_1 bash -c 'go run main.go';
fi
(docker-compose --file docker-compose_dev.yml up -d) && docker exec -it cerealnotes_backend_1 bash -c 'go run main.go'
# docker exec -it cerealnotes_backend_1 /bin/bash
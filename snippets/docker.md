# Docker Snippets

## List all running containers

docker ps

## Stop all running containers

docker stop $(docker ps -a -q)

## Remove all unused volumes

docker volume prune -f

## Build an image from a Dockerfile in current directory

docker build -t snip-cli .

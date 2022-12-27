VERSION:="testing"

serve-frontend:
	echo "TODO: add script here"

serve-api:
	echo "TODO: add script here"

build-frontend:
	docker build ./frontend -t animalfacts-frontend:${VERSION} -f ./frontend/Dockerfile

build-api:
	docker build ./api -t animalfacts-api:${VERSION} -f ./api/Dockerfile

VERSION:="testing"

dev-serve-frontend:
	echo "TODO: add script here"

dev-serve-api:
	echo "TODO: add script here"

ci-build-frontend:
	docker build ./frontend -t animalfacts-frontend:${VERSION} -f ./frontend/Dockerfile

ci-build-api:
	docker build ./api -t animalfacts-api:${VERSION} -f ./api/Dockerfile

# animal-facts

> **Warning**
> Please note that this project is in early development phase. There are no releases or a stable hosted version until now, this will follow in the future.

This projects is about awesome facts about animals and make more people know these.

The core is an [API](./backend/) written in Golang that provides the animal facts. In the background there is a PostgreSQL database for storing the facts data. Additionally there is a [vue.js frontend](./frontend/), that is used as a showcase for the API and also handling requests for new facts or updates of existing facts. And for the CLI users there is an [animal-facts-cli](./cli/), that uses the API in the background.

Feel free to visit [animalfacts.app](https://animalfacts.app) to see the latest version of the project in action!

## Technologies

This repository uses the following technologies and frameworks:

- Database
    - [Firestore](https://firebase.google.com/docs/firestore/)
- API
    - [Golang](https://go.dev/)
- CLI
    - [Golang](https://go.dev/)
    - [Cobra](https://cobra.dev/)
- Frontend
    - [Vue.js](https://vuejs.org/)
    - [TypeScript](https://www.typescriptlang.org/)
    - [NodeJS](https://nodejs.org/en/)
    - [Yarn](https://yarnpkg.com/)

## Use the API

The API has a generated [swagger UI](https://animalfacts.app/swagger/index.html). There are public endpoints, the ones for getting animal facts and also private ones that require authentication (for adding / updating facts).

### Examples

All endpoints are returning JSON format.

1. Get a random animal fact from the database:
```
curl https://animalfacts.app/api/v1/fact
```

2. Get a specific animal fact by ID:
```
curl https://animalfacts.app/api/v1/fact/42
```

## Run the setup locally for development
### Database
To run the whole setup locally, you first need a running PostgreSQL database for the animal facts. An easy way to get this started would be this docker compose setup:
```
version: '3.8'

services:
  animalfacts-db:
    image: postgres:15.1
    container_name: animalfacts-db
    restart: unless-stopped
    environment:
      POSTGRES_USER: animalfacts
      POSTGRES_PASSWORD: animalfacts
      POSTGRES_DB: animalfacts
      PGDATA: /var/lib/postgresql/data/pgdata
    ports:
      - '5432:5432'
    volumes:
      - /etc/animalfacts/db:/var/lib/postgresql/data
```
Make sure to use a volume at your PostgreSQL database to not lose any fact data.
### API
To start the Golang API locally run the following command in the [API directoy](./backend/):
```
go run ./...
```
### CLI
To start the Golang CLI locally run the following command in the [CLI directoy](./cli/):
```
go run ./...
```
### Frontend
To start the Angular frontend locally run the following commands in the [frontend directoy](./frontend/):
```
# install dependencies, if you don't have yarn installed, do that first (with npm i -g yarn)
yarn

# start the Angular frontend
yarn start
```

# animal-facts

Simple API to spread random animal facts to everyone. Animals are awesome, everyone should know this.

## Usage of public api

The public api is built get facts without authentication for public usage.

Detailed information about the available endpoints can be seen on the [swagger page](https://animal-facts.cafo.dev/swagger/index.html).

```shell
# get random fact
curl https://animal-facts.cafo.dev/api/v1/api/v1/facts
# example response
{"id":"6578bf140e487ecc049c7594","fact":"The Blue Whale is the largest animal that has ever lived.","source":"https://factanimal.com/blue-whale/"}

# get fact by id
curl https://animal-facts.cafo.dev/api/v1/api/v1/facts/6578bf140e487ecc049c7594
# example response
{"id":"6578bf140e487ecc049c7594","fact":"The Blue Whale is the largest animal that has ever lived.","source":"https://factanimal.com/blue-whale/"}
```

## Usage of internal api

The internal api is built to manage the facts database. A management UI using the API is built with appsmith [here](https://github.com/cafo13/animal-facts-manager). To get access to be able to manage the public's api database of https://animal-facts.cafo.dev/, feel free to create an issue at this repository.

## Development with own database

Prerequisites:
- [go](https://go.dev/doc/install):
  - version 1.21
- [mongo database](https://www.mongodb.com/):
  - database named "animal-facts" with collection named "facts" (database name can be overwritten with environment variable MONGODB_DATABASE_NAME)
- copy the [.env.dist](.env.dist) file to [.env](.env) and fill the variables for the mongodb connection to your database

```shell
# run the public api locally
make public-api-run

# run the internal api locally
make internal-api-run

# build the public api locally
make public-api-build

# build the internal api locally
make internal-api-build

# generate swagger docs for the public api locally
make public-api-generate-swagger

# generate swagger docs for the internal api locally
make internal-api-generate-swagger
```

FROM golang:1.19.4

ENV GIN_MODE: "release"

WORKDIR /animalfacts

COPY . .
RUN go mod download && go mod verify
RUN CGO_ENABLED=0 GOOS=linux go build -o /animalfacts ./...

CMD ["/animalfacts/api"]

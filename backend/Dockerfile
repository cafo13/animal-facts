FROM golang:1.20.1

WORKDIR /animalfacts

COPY . .
RUN go mod download && go mod verify
RUN CGO_ENABLED=0 GOOS=linux go build -o /animalfacts ./...

CMD ["/animalfacts/backend"]

FROM golang:latest

WORKDIR /asperitas
COPY . .

COPY go.mod .
COPY go.sum .

RUN go mod download

RUN go build -o main ./06_databases/99_hw/redditclone/cmd/redditclone

CMD ["/asperitas/main"]

EXPOSE 8080
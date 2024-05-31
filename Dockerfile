FROM golang:1.22-alpine AS build_stage
COPY . /app
WORKDIR /app
COPY go.mod .
COPY go.sum .
RUN go mod download
RUN CGO_ENABLED=0 go build -o /app_binary/ozon ./cmd/ozon


FROM alpine AS run_stage
WORKDIR /app_binary
COPY --from=build_stage /app_binary/ozon /app_binary/
RUN chmod +x ./ozon
ENTRYPOINT ./ozon
CMD [ "ozon" ]

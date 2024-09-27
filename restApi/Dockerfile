FROM golang:1.23.1 AS build-stage
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY *.go ./
RUN CGO_ENABLED=0 GOOS=linux go build -o /api
FROM scratch as run-release-stage
WORKDIR /app
COPY --from=build-stage /api /api
CMD ["/api"]

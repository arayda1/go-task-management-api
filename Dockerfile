FROM golang:1.23-alpine AS build
WORKDIR /app
COPY go.mod ./
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o task-api .

FROM alpine:3.20
WORKDIR /app
COPY --from=build /app/task-api .
EXPOSE 8080
CMD ["./task-api"]

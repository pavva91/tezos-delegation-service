# build stage
FROM golang:1.22-alpine AS build

# WORKDIR /usr/src/app
WORKDIR /

# pre-copy/cache go.mod for pre-downloading dependencies and only redownloading them in subsequent builds if they change
COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .
# RUN go build -v -o /usr/local/bin/app ./...
# RUN go build -v -o /usr/local/bin/app main.go
RUN go build -v -o ./app main.go

# runtime stage
FROM alpine

WORKDIR /usr/local/bin

# COPY --from=build /usr/local/bin/app .
COPY --from=build /app .

CMD ["app"]

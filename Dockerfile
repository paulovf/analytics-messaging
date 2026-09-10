FROM golang:1.27-alpine

WORKDIR /app

RUN apk add --no-cache gcc musl-dev

COPY go.mod go.sum ./
RUN apk add --no-cache gcc musl-dev bash

COPY . .

CMD ["go", "test", "-v", "./..."]

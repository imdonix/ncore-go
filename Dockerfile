FROM alpine:latest AS builder

WORKDIR /app

RUN apk add --no-cache go git make

RUN go version

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN make


FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /app

ENV GIN_MODE=release

COPY --from=builder /app/bin/ncore .

EXPOSE 8080

CMD ["./ncore"]

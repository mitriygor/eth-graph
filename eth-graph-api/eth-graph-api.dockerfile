FROM golang:1.19

WORKDIR /app/eth-graph-api

RUN go install github.com/cosmtrek/air@latest


COPY eth-graph-api/go.mod /app/eth-graph-api
COPY eth-graph-api/.air.toml /app/eth-graph-api

RUN go mod download

COPY eth-graph-api /app/eth-graph-api

CMD ["air", "-c", ".air.toml"]
FROM golang:1.22 AS builder

WORKDIR /app

ADD go.mod go.sum ./

RUN go mod download

COPY ./src/ ./src

RUN go build -o kpi ./src

FROM ubuntu

WORKDIR /app

RUN apt update
RUN apt install -y ca-certificates

COPY --from=builder /app/kpi .

EXPOSE 8080

CMD [ "./kpi" ]

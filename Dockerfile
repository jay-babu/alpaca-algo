FROM golang:1-bullseye

RUN apt update
RUN apt upgrade -y
WORKDIR /build

# Fetch dependencies
COPY go.mod go.sum ./
RUN go mod download

# Build
COPY . ./

RUN go build -race alpacaAlgo

EXPOSE $PORT

ENTRYPOINT [ "./alpacaAlgo" ]

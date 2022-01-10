FROM ubuntu:latest

# install docker
RUN apt-get update && \
    apt-get install -y \
    apt-transport-https \
    ca-certificates \
    curl \
    gnupg-agent \
    software-properties-common && \
    apt install docker.io -y

# copy golang binary from docker image
COPY --from=golang:latest /usr/local/go /usr/local/go

ENV PATH=$PATH:/usr/local/go/bin

RUN go version

# copy files to go path
COPY go.mod ./
COPY go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ./bin/main .

CMD ["./bin/main"]
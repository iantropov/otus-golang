FROM golang:1.20.1

RUN mkdir -p /opt/functest
WORKDIR /opt/functest

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .
CMD ["go", "test"]
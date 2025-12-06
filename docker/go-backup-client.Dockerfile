FROM golang:1.24 AS build

WORKDIR /go-backup

COPY . .

RUN apt-get update && apt-get -y install protobuf-compiler

RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@latest

RUN make proto
RUN make tidy
RUN make build

FROM debian:bookworm-slim AS runtime

WORKDIR /go-backup

RUN apt-get update && apt-get -y install procps

COPY --from=build /go-backup/build/go-backup-linux-amd64 .
RUN chmod +x go-backup-linux-amd64

ENTRYPOINT ["./go-backup-linux-amd64", "--client"]

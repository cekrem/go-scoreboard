FROM golang:1.12 as builder

ENV CGO_ENABLED=0
CMD mkdir -p /home/go/app

RUN go get github.com/yudai/gotty
RUN go install github.com/yudai/gotty

WORKDIR /home/go/app

COPY go.* ./
RUN go mod download

COPY main.go ./
RUN go build

FROM gcr.io/distroless/static

COPY --from=builder /home/go/app/go-scoreboard .
COPY --from=builder /go/bin/gotty .

EXPOSE 8080
ENTRYPOINT ["/gotty", "-w", "/go-scoreboard"]


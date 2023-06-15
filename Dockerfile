FROM golang:alpine as builder

MAINTAINER github@strnad.ch

WORKDIR /go/src/app

ARG VERSION
ENV CGO_ENABLED=0

COPY go.mod .
COPY go.sum .
RUN go mod download

# Now we copy the code and build, this way we can cache the dependencies
COPY . .
RUN --mount=type=cache,target=/root/.cache/go-build  go build \
    -a -tags netgo,osusergo \
    -ldflags "-extldflags '-static' -s -w" \
    -ldflags "-X main.version=$VERSION" \
    *.go

## Build the Run Container
FROM alpine

WORKDIR /go/src/app
COPY html ./html
COPY static ./static
COPY --from=builder /go/src/app/main .

ENTRYPOINT ["./main"]

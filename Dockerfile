# syntax=docker/dockerfile:1

FROM golang:1.23.5 AS build

WORKDIR $GOPATH/src/github.com/brotherlogic/github_overseer

COPY go.mod ./
COPY go.sum ./

RUN go mod download

COPY *.go ./

RUN mkdir proto
COPY proto/*.go ./proto/

RUN CGO_ENABLED=0 go build -o /github_overseer

##
## Deploy
##
FROM gcr.io/distroless/base-debian11

WORKDIR /

COPY --from=build /github_overseer /github_overseer

EXPOSE 8080
EXPOSE 8081
EXPOSE 8082

USER nonroot:nonroot

ENTRYPOINT ["/github_overseer"]
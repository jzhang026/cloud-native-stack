FROM golang:1.17 AS build-base
WORKDIR .
COPY . .
ENV GOOS=linux \
    GOARCH=amd64 \
    GOPATH=/
RUN go build -o /httpServer .

FROM busybox
LABEL multi.tech="http" multi.type="exercise" other="toy"
COPY --from=build-base /httpServer /httpServer
ENV PORT=80 \
    VERSION="v1.1"
EXPOSE 80
ENTRYPOINT /httpServer
FROM node:current-alpine AS frontend-builder
LABEL authors="Daniel Wu"

COPY . /tmp/remote-serial-port-server
WORKDIR /tmp/remote-serial-port-server/static
RUN rm -rf dist node_modules
RUN npm install && npm run build

FROM golang:alpine AS builder
LABEL authors="Daniel Wu"

ARG VERSION="0.0.0"
ARG BUILD_DATE="unknown"
ARG SHA=""

COPY . /tmp/remote-serial-port-server
COPY --from=frontend-builder /tmp/remote-serial-port-server/static/dist /tmp/remote-serial-port-server/static/dist
WORKDIR /tmp/remote-serial-port-server
RUN go build -trimpath \
        -ldflags="-extldflags \"-static\" -X 'github.com/iceking2nd/remote-serial-port-server/global.Version=${VERSION}' -X 'github.com/iceking2nd/remote-serial-port-server/global.BuildTime=${BUILD_DATE}' -X github.com/iceking2nd/remote-serial-port-server/global.GitCommit=${SHA}" \
        -o /bin/remote-serial-port-server

FROM scratch
LABEL authors="Daniel Wu"

COPY --from=builder /bin/remote-serial-port-server /bin/remote-serial-port-server
EXPOSE 8192
ENTRYPOINT ["/bin/remote-serial-port-server","-l","0.0.0.0","-p","8192"]

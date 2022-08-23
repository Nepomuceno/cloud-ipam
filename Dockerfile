FROM node:18-alpine AS buildnode
ARG NPM_TOKEN
WORKDIR /usr/src/app
COPY ui/package*.json /usr/src/app/
RUN npm ci
COPY ui/ /usr/src/app/
RUN npm run build

# build stage
FROM golang:1.19-alpine AS buildgo
# Install git + SSL ca certificates.
# Git is required for fetching the dependencies.
# Ca-certificates is required to call HTTPS endpoints.
RUN apk update && apk add --no-cache git ca-certificates tzdata && update-ca-certificates

# Create appuser
ENV USER=appuser
ENV UID=10001

# See https://stackoverflow.com/a/55757473/12429735
RUN adduser \
    --disabled-password \
    --gecos "" \
    --home "/nonexistent" \
    --shell "/sbin/nologin" \
    --no-create-home \
    --uid "${UID}" \
    "${USER}"

WORKDIR /app

COPY go.mod go.mod
COPY go.sum go.sum
RUN go mod download

COPY --from=buildnode /usr/src/app/dist ui/dist
COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags='-w -s -extldflags "-static"' -a \
    -o /go/bin/cloud-ipam .

# final stage
FROM busybox:stable

COPY --from=buildgo /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=buildgo /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=buildgo /etc/passwd /etc/passwd
COPY --from=buildgo /etc/group /etc/group

COPY --from=buildgo /go/bin/cloud-ipam /go/bin/cloud-ipam

USER appuser:appuser

EXPOSE 8080

ENTRYPOINT /go/bin/cloud-ipam -s "$STORAGE_ACCOUNT_NAME" serve
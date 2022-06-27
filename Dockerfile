FROM golang:1.18-alpine AS build-env

WORKDIR /app

# We copy go.mod and go.sum before copying the source code. This allows us to
# speed up the build process if there were no changes in the dependencies.
COPY go.mod .
COPY go.sum .
RUN go mod download -x

# Copy over the source code and build
COPY . .
ARG VCS_REF
RUN CGO_ENABLED=0 go build

FROM alpine:latest

ARG ADDITIONAL_PACKAGES

RUN apk add --no-cache tzdata ca-certificates ${ADDITIONAL_PACKAGES} && \
	ln -fs /usr/share/zoneinfo/Europe/Rome /etc/localtime && \
	ls -la /etc/localtime

WORKDIR /root/
COPY --from=build-env /etc/ssl/certs /etc/ssl/certs
COPY --from=build-env /app .

CMD ["./lineameteo-prometheus"]


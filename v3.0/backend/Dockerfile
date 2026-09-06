FROM golang:1.25-alpine AS build
RUN apk add --no-cache build-base
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
ARG SERVICE
RUN --mount=type=cache,target=/root/.cache/go-build CGO_ENABLED=1 GOOS=linux go build -trimpath -o /out/service ./cmd/${SERVICE}

FROM alpine:3.22
RUN apk add --no-cache ca-certificates wget
RUN addgroup -S polaris && adduser -S -G polaris -h /app polaris
WORKDIR /app
COPY --from=build --chown=polaris:polaris /out/service /app/service
COPY --chown=polaris:polaris data /app/data
USER polaris
EXPOSE 6080 6081
STOPSIGNAL SIGTERM
ENTRYPOINT ["/app/service"]

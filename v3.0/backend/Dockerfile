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
WORKDIR /app
COPY --from=build /out/service /app/service
COPY data /app/data
ENTRYPOINT ["/app/service"]

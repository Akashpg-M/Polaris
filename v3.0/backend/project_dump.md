# Project Dump

Project: backend

## Directory Tree

```text
backend
├── .dockerignore
├── algo_
│   ├── geo
│   │   └── math.go
│   ├── graph
│   ├── logger
│   │   └── logger.go
│   └── quadtree
├── api
│   └── proto
│       ├── v1
│       │   ├── spatial.pb.go
│       │   └── spatial.proto
│       └── v2
├── cmd
│   ├── engine
│   │   ├── main.go
│   │   └── mobility_rebuild.go
│   ├── gateway
│   │   └── main.go
│   ├── identitycheck
│   │   └── main.go
│   ├── loadtest
│   │   └── main.go
│   ├── mobility-benchmark
│   │   └── main.go
│   ├── orchestrationcheck
│   │   └── main.go
│   ├── routing-benchmark
│   │   └── main.go
│   ├── routing-overload
│   │   └── main.go
│   ├── smoke
│   │   └── main.go
│   └── system-soak
│       ├── main.go
│       └── main_test.go
├── data
│   └── chennai-metro.osm.pbf
├── deployments
│   ├── docker-compose.yml
│   ├── init.sql
│   ├── phase2-identity-test.ps1
│   ├── phase3-command-test.ps1
│   ├── phase4-closure-soak.ps1
│   ├── phase4-mobility-rebuild-test.ps1
│   ├── phase4-mobility-test.ps1
│   ├── phase4-module-isolation-test.ps1
│   ├── phase4-routing-overload-test.ps1
│   ├── reliability-test.ps1
│   └── smoke-test.ps1
├── Dockerfile
├── gateway.exe
├── go.mod
├── go.sum
├── internal
│   ├── adapter
│   │   ├── handler
│   │   │   ├── dashboard.go
│   │   │   ├── device_connections.go
│   │   │   ├── ingestion.go
│   │   │   ├── ingestion_test.go
│   │   │   ├── match.go
│   │   │   ├── mobility_api.go
│   │   │   ├── orchestration_api.go
│   │   │   ├── registry_api.go
│   │   │   └── registry_api_test.go
│   │   ├── osrm
│   │   └── repository
│   │       ├── kafka_event_publisher.go
│   │       ├── kafka_stream.go
│   │       ├── postgres_orchestration.go
│   │       ├── postgres_registry.go
│   │       └── redis_connection.go
│   ├── application
│   │   ├── dispatch
│   │   │   └── dispatcher.go
│   │   ├── orchestration
│   │   │   ├── metrics.go
│   │   │   ├── service.go
│   │   │   └── service_test.go
│   │   ├── orchestrator
│   │   │   ├── predictive_strategy.go
│   │   │   ├── strategies.go
│   │   │   └── zone.go
│   │   ├── outbox
│   │   │   └── relay.go
│   │   ├── reconciliation
│   │   │   └── worker.go
│   │   ├── spatial
│   │   │   ├── engine.go
│   │   │   └── engine_reliability_test.go
│   │   ├── stream
│   │   │   ├── archiver.go
│   │   │   ├── archiver_integration_test.go
│   │   │   ├── kafka_consumer.go
│   │   │   ├── kafka_consumer_reliability_test.go
│   │   │   ├── redis_projection.go
│   │   │   ├── redis_projection_integration_test.go
│   │   │   └── state_fanout.go
│   │   └── twin
│   │       └── connectivity.go
│   ├── architecture
│   │   └── consistency_test.go
│   ├── config
│   │   └── config.go
│   ├── core
│   │   ├── actor
│   │   ├── auth
│   │   │   ├── auth.go
│   │   │   └── auth_test.go
│   │   ├── command
│   │   │   ├── model.go
│   │   │   └── model_test.go
│   │   ├── domain
│   │   │   └── node.go
│   │   ├── events
│   │   │   ├── telemetry.go
│   │   │   └── telemetry_test.go
│   │   ├── extension
│   │   │   ├── contracts.go
│   │   │   ├── default_planner.go
│   │   │   ├── registry.go
│   │   │   └── registry_test.go
│   │   ├── ports
│   │   │   └── stream.go
│   │   ├── registry
│   │   │   ├── model.go
│   │   │   └── model_test.go
│   │   ├── routing
│   │   ├── simulation
│   │   │   ├── ca_runline.go
│   │   │   └── protocols.go
│   │   ├── task
│   │   │   ├── model.go
│   │   │   └── model_test.go
│   │   └── twin
│   │       └── component.go
│   ├── infra
│   │   ├── osm
│   │   ├── postgres
│   │   │   └── db.go
│   │   └── redis
│   │       └── client.go
│   └── modules
│       └── mobility
│           ├── config.go
│           ├── config_test.go
│           ├── matching
│           │   └── candidate_provider.go
│           ├── model
│           │   └── model.go
│           ├── module.go
│           ├── module_test.go
│           ├── planning
│           │   ├── task_planner.go
│           │   └── task_planner_test.go
│           ├── projector.go
│           ├── routing
│           │   ├── engine.go
│           │   ├── graph.go
│           │   ├── kdtree.go
│           │   ├── osm.go
│           │   ├── routing_test.go
│           │   ├── search.go
│           │   ├── snapshot.go
│           │   ├── traffic.go
│           │   └── types.go
│           ├── spatial
│           │   ├── geo.go
│           │   ├── index.go
│           │   ├── manager.go
│           │   ├── rtree.go
│           │   ├── rtree_contention_test.go
│           │   └── spatial_test.go
│           └── traffic_consumer.go
└── project_dump.md
```

# File Contents

---

## .dockerignore

```
.git
.tmp
.gocache
.gotmp
.phase4-cache-*
.phase*-go-cache
**/*_test.go
deployments
```

---

## Dockerfile

```
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
```

---

## gateway.exe

**[Skipped - File Too Large (38048.9 KB)]**

---

## go.mod

```
module github.com/Akashpg-M/polaris/backend

go 1.25.0

require (
	github.com/gin-contrib/cors v1.7.6
	github.com/gin-gonic/gin v1.12.0
	github.com/gorilla/websocket v1.5.3
	github.com/jmoiron/sqlx v1.4.0
	github.com/joho/godotenv v1.5.1
	github.com/lib/pq v1.11.2
	github.com/qedus/osmpbf v1.2.0
	github.com/redis/go-redis/v9 v9.18.0
	github.com/segmentio/kafka-go v0.4.51
	github.com/uber/h3-go/v4 v4.5.0
	google.golang.org/protobuf v1.36.10
)

require (
	github.com/bytedance/gopkg v0.1.3 // indirect
	github.com/bytedance/sonic v1.15.0 // indirect
	github.com/bytedance/sonic/loader v0.5.0 // indirect
	github.com/cespare/xxhash/v2 v2.3.0 // indirect
	github.com/cloudwego/base64x v0.1.6 // indirect
	github.com/dgryski/go-rendezvous v0.0.0-20200823014737-9f7001d12a5f // indirect
	github.com/gabriel-vasile/mimetype v1.4.12 // indirect
	github.com/gin-contrib/sse v1.1.0 // indirect
	github.com/go-playground/locales v0.14.1 // indirect
	github.com/go-playground/universal-translator v0.18.1 // indirect
	github.com/go-playground/validator/v10 v10.30.1 // indirect
	github.com/goccy/go-json v0.10.5 // indirect
	github.com/goccy/go-yaml v1.19.2 // indirect
	github.com/json-iterator/go v1.1.12 // indirect
	github.com/klauspost/compress v1.17.6 // indirect
	github.com/klauspost/cpuid/v2 v2.3.0 // indirect
	github.com/leodido/go-urn v1.4.0 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.2 // indirect
	github.com/pelletier/go-toml/v2 v2.2.4 // indirect
	github.com/pierrec/lz4/v4 v4.1.15 // indirect
	github.com/quic-go/qpack v0.6.0 // indirect
	github.com/quic-go/quic-go v0.59.0 // indirect
	github.com/twitchyliquid64/golang-asm v0.15.1 // indirect
	github.com/ugorji/go/codec v1.3.1 // indirect
	go.mongodb.org/mongo-driver/v2 v2.5.0 // indirect
	go.uber.org/atomic v1.11.0 // indirect
	golang.org/x/arch v0.22.0 // indirect
	golang.org/x/crypto v0.48.0 // indirect
	golang.org/x/net v0.51.0 // indirect
	golang.org/x/sys v0.41.0 // indirect
	golang.org/x/text v0.34.0 // indirect
)
```

---

## go.sum

```
filippo.io/edwards25519 v1.1.0 h1:FNf4tywRC1HmFuKW5xopWpigGjJKiJSV0Cqo0cJWDaA=
filippo.io/edwards25519 v1.1.0/go.mod h1:BxyFTGdWcka3PhytdK4V28tE5sGfRvvvRV7EaN4VDT4=
github.com/bsm/ginkgo/v2 v2.12.0 h1:Ny8MWAHyOepLGlLKYmXG4IEkioBysk6GpaRTLC8zwWs=
github.com/bsm/ginkgo/v2 v2.12.0/go.mod h1:SwYbGRRDovPVboqFv0tPTcG1sN61LM1Z4ARdbAV9g4c=
github.com/bsm/gomega v1.27.10 h1:yeMWxP2pV2fG3FgAODIY8EiRE3dy0aeFYt4l7wh6yKA=
github.com/bsm/gomega v1.27.10/go.mod h1:JyEr/xRbxbtgWNi8tIEVPUYZ5Dzef52k01W3YH0H+O0=
github.com/bytedance/gopkg v0.1.3 h1:TPBSwH8RsouGCBcMBktLt1AymVo2TVsBVCY4b6TnZ/M=
github.com/bytedance/gopkg v0.1.3/go.mod h1:576VvJ+eJgyCzdjS+c4+77QF3p7ubbtiKARP3TxducM=
github.com/bytedance/sonic v1.15.0 h1:/PXeWFaR5ElNcVE84U0dOHjiMHQOwNIx3K4ymzh/uSE=
github.com/bytedance/sonic v1.15.0/go.mod h1:tFkWrPz0/CUCLEF4ri4UkHekCIcdnkqXw9VduqpJh0k=
github.com/bytedance/sonic/loader v0.5.0 h1:gXH3KVnatgY7loH5/TkeVyXPfESoqSBSBEiDd5VjlgE=
github.com/bytedance/sonic/loader v0.5.0/go.mod h1:AR4NYCk5DdzZizZ5djGqQ92eEhCCcdf5x77udYiSJRo=
github.com/cespare/xxhash/v2 v2.3.0 h1:UL815xU9SqsFlibzuggzjXhog7bL6oX9BbNZnL2UFvs=
github.com/cespare/xxhash/v2 v2.3.0/go.mod h1:VGX0DQ3Q6kWi7AoAeZDth3/j3BFtOZR5XLFGgcrjCOs=
github.com/cloudwego/base64x v0.1.6 h1:t11wG9AECkCDk5fMSoxmufanudBtJ+/HemLstXDLI2M=
github.com/cloudwego/base64x v0.1.6/go.mod h1:OFcloc187FXDaYHvrNIjxSe8ncn0OOM8gEHfghB2IPU=
github.com/davecgh/go-spew v1.1.0/go.mod h1:J7Y8YcW2NihsgmVo/mv3lAwl/skON4iLHjSsI+c5H38=
github.com/davecgh/go-spew v1.1.1 h1:vj9j/u1bqnvCEfJOwUhtlOARqs3+rkHYY13jYWTU97c=
github.com/davecgh/go-spew v1.1.1/go.mod h1:J7Y8YcW2NihsgmVo/mv3lAwl/skON4iLHjSsI+c5H38=
github.com/dgryski/go-rendezvous v0.0.0-20200823014737-9f7001d12a5f h1:lO4WD4F/rVNCu3HqELle0jiPLLBs70cWOduZpkS1E78=
github.com/dgryski/go-rendezvous v0.0.0-20200823014737-9f7001d12a5f/go.mod h1:cuUVRXasLTGF7a8hSLbxyZXjz+1KgoB3wDUb6vlszIc=
github.com/gabriel-vasile/mimetype v1.4.12 h1:e9hWvmLYvtp846tLHam2o++qitpguFiYCKbn0w9jyqw=
github.com/gabriel-vasile/mimetype v1.4.12/go.mod h1:d+9Oxyo1wTzWdyVUPMmXFvp4F9tea18J8ufA774AB3s=
github.com/gin-contrib/cors v1.7.6 h1:3gQ8GMzs1Ylpf70y8bMw4fVpycXIeX1ZemuSQIsnQQY=
github.com/gin-contrib/cors v1.7.6/go.mod h1:Ulcl+xN4jel9t1Ry8vqph23a60FwH9xVLd+3ykmTjOk=
github.com/gin-contrib/sse v1.1.0 h1:n0w2GMuUpWDVp7qSpvze6fAu9iRxJY4Hmj6AmBOU05w=
github.com/gin-contrib/sse v1.1.0/go.mod h1:hxRZ5gVpWMT7Z0B0gSNYqqsSCNIJMjzvm6fqCz9vjwM=
github.com/gin-gonic/gin v1.12.0 h1:b3YAbrZtnf8N//yjKeU2+MQsh2mY5htkZidOM7O0wG8=
github.com/gin-gonic/gin v1.12.0/go.mod h1:VxccKfsSllpKshkBWgVgRniFFAzFb9csfngsqANjnLc=
github.com/go-playground/assert/v2 v2.2.0 h1:JvknZsQTYeFEAhQwI4qEt9cyV5ONwRHC+lYKSsYSR8s=
github.com/go-playground/assert/v2 v2.2.0/go.mod h1:VDjEfimB/XKnb+ZQfWdccd7VUvScMdVu0Titje2rxJ4=
github.com/go-playground/locales v0.14.1 h1:EWaQ/wswjilfKLTECiXz7Rh+3BjFhfDFKv/oXslEjJA=
github.com/go-playground/locales v0.14.1/go.mod h1:hxrqLVvrK65+Rwrd5Fc6F2O76J/NuW9t0sjnWqG1slY=
github.com/go-playground/universal-translator v0.18.1 h1:Bcnm0ZwsGyWbCzImXv+pAJnYK9S473LQFuzCbDbfSFY=
github.com/go-playground/universal-translator v0.18.1/go.mod h1:xekY+UJKNuX9WP91TpwSH2VMlDf28Uj24BCp08ZFTUY=
github.com/go-playground/validator/v10 v10.30.1 h1:f3zDSN/zOma+w6+1Wswgd9fLkdwy06ntQJp0BBvFG0w=
github.com/go-playground/validator/v10 v10.30.1/go.mod h1:oSuBIQzuJxL//3MelwSLD5hc2Tu889bF0Idm9Dg26cM=
github.com/go-sql-driver/mysql v1.8.1 h1:LedoTUt/eveggdHS9qUFC1EFSa8bU2+1pZjSRpvNJ1Y=
github.com/go-sql-driver/mysql v1.8.1/go.mod h1:wEBSXgmK//2ZFJyE+qWnIsVGmvmEKlqwuVSjsCm7DZg=
github.com/goccy/go-json v0.10.5 h1:Fq85nIqj+gXn/S5ahsiTlK3TmC85qgirsdTP/+DeaC4=
github.com/goccy/go-json v0.10.5/go.mod h1:oq7eo15ShAhp70Anwd5lgX2pLfOS3QCiwU/PULtXL6M=
github.com/goccy/go-yaml v1.19.2 h1:PmFC1S6h8ljIz6gMRBopkjP1TVT7xuwrButHID66PoM=
github.com/goccy/go-yaml v1.19.2/go.mod h1:XBurs7gK8ATbW4ZPGKgcbrY1Br56PdM69F7LkFRi1kA=
github.com/golang/protobuf v1.5.0/go.mod h1:FsONVRAS9T7sI+LIUmWTfcYkHO4aIWwzhcaSAoJOfIk=
github.com/google/go-cmp v0.5.5/go.mod h1:v8dTdLbMG2kIc/vJvl+f65V22dbkXbowE6jgT/gNBxE=
github.com/google/go-cmp v0.7.0 h1:wk8382ETsv4JYUZwIsn6YpYiWiBsYLSJiTsyBybVuN8=
github.com/google/go-cmp v0.7.0/go.mod h1:pXiqmnSA92OHEEa9HXL2W4E7lf9JzCmGVUdgjX3N/iU=
github.com/google/gofuzz v1.0.0/go.mod h1:dBl0BpW6vV/+mYPU4Po3pmUjxk6FQPldtuIdl/M65Eg=
github.com/gorilla/websocket v1.5.3 h1:saDtZ6Pbx/0u+bgYQ3q96pZgCzfhKXGPqt7kZ72aNNg=
github.com/gorilla/websocket v1.5.3/go.mod h1:YR8l580nyteQvAITg2hZ9XVh4b55+EU/adAjf1fMHhE=
github.com/jmoiron/sqlx v1.4.0 h1:1PLqN7S1UYp5t4SrVVnt4nUVNemrDAtxlulVe+Qgm3o=
github.com/jmoiron/sqlx v1.4.0/go.mod h1:ZrZ7UsYB/weZdl2Bxg6jCRO9c3YHl8r3ahlKmRT4JLY=
github.com/joho/godotenv v1.5.1 h1:7eLL/+HRGLY0ldzfGMeQkb7vMd0as4CfYvUVzLqw0N0=
github.com/joho/godotenv v1.5.1/go.mod h1:f4LDr5Voq0i2e/R5DDNOoa2zzDfwtkZa6DnEwAbqwq4=
github.com/json-iterator/go v1.1.12 h1:PV8peI4a0ysnczrg+LtxykD8LfKY9ML6u2jnxaEnrnM=
github.com/json-iterator/go v1.1.12/go.mod h1:e30LSqwooZae/UwlEbR2852Gd8hjQvJoHmT4TnhNGBo=
github.com/klauspost/compress v1.17.6 h1:60eq2E/jlfwQXtvZEeBUYADs+BwKBWURIY+Gj2eRGjI=
github.com/klauspost/compress v1.17.6/go.mod h1:/dCuZOvVtNoHsyb+cuJD3itjs3NbnF6KH9zAO4BDxPM=
github.com/klauspost/cpuid/v2 v2.3.0 h1:S4CRMLnYUhGeDFDqkGriYKdfoFlDnMtqTiI/sFzhA9Y=
github.com/klauspost/cpuid/v2 v2.3.0/go.mod h1:hqwkgyIinND0mEev00jJYCxPNVRVXFQeu1XKlok6oO0=
github.com/leodido/go-urn v1.4.0 h1:WT9HwE9SGECu3lg4d/dIA+jxlljEa1/ffXKmRjqdmIQ=
github.com/leodido/go-urn v1.4.0/go.mod h1:bvxc+MVxLKB4z00jd1z+Dvzr47oO32F/QSNjSBOlFxI=
github.com/lib/pq v1.10.9/go.mod h1:AlVN5x4E4T544tWzH6hKfbfQvm3HdbOxrmggDNAPY9o=
github.com/lib/pq v1.11.2 h1:x6gxUeu39V0BHZiugWe8LXZYZ+Utk7hSJGThs8sdzfs=
github.com/lib/pq v1.11.2/go.mod h1:/p+8NSbOcwzAEI7wiMXFlgydTwcgTr3OSKMsD2BitpA=
github.com/mattn/go-isatty v0.0.20 h1:xfD0iDuEKnDkl03q4limB+vH+GxLEtL/jb4xVJSWWEY=
github.com/mattn/go-isatty v0.0.20/go.mod h1:W+V8PltTTMOvKvAeJH7IuucS94S2C6jfK/D7dTCTo3Y=
github.com/mattn/go-sqlite3 v1.14.22 h1:2gZY6PC6kBnID23Tichd1K+Z0oS6nE/XwU+Vz/5o4kU=
github.com/mattn/go-sqlite3 v1.14.22/go.mod h1:Uh1q+B4BYcTPb+yiD3kU8Ct7aC0hY9fxUwlHK0RXw+Y=
github.com/modern-go/concurrent v0.0.0-20180228061459-e0a39a4cb421/go.mod h1:6dJC0mAP4ikYIbvyc7fijjWJddQyLn8Ig3JB5CqoB9Q=
github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd h1:TRLaZ9cD/w8PVh93nsPXa1VrQ6jlwL5oN8l14QlcNfg=
github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd/go.mod h1:6dJC0mAP4ikYIbvyc7fijjWJddQyLn8Ig3JB5CqoB9Q=
github.com/modern-go/reflect2 v1.0.2 h1:xBagoLtFs94CBntxluKeaWgTMpvLxC4ur3nMaC9Gz0M=
github.com/modern-go/reflect2 v1.0.2/go.mod h1:yWuevngMOJpCy52FWWMvUC8ws7m/LJsjYzDa0/r8luk=
github.com/pelletier/go-toml/v2 v2.2.4 h1:mye9XuhQ6gvn5h28+VilKrrPoQVanw5PMw/TB0t5Ec4=
github.com/pelletier/go-toml/v2 v2.2.4/go.mod h1:2gIqNv+qfxSVS7cM2xJQKtLSTLUE9V8t9Stt+h56mCY=
github.com/pierrec/lz4/v4 v4.1.15 h1:MO0/ucJhngq7299dKLwIMtgTfbkoSPF6AoMYDd8Q4q0=
github.com/pierrec/lz4/v4 v4.1.15/go.mod h1:gZWDp/Ze/IJXGXf23ltt2EXimqmTUXEy0GFuRQyBid4=
github.com/pmezard/go-difflib v1.0.0 h1:4DBwDE0NGyQoBHbLQYPwSUPoCMWR5BEzIk/f1lZbAQM=
github.com/pmezard/go-difflib v1.0.0/go.mod h1:iKH77koFhYxTK1pcRnkKkqfTogsbg7gZNVY4sRDYZ/4=
github.com/qedus/osmpbf v1.2.0 h1:yRm5ECkiUsN9sA+UN9yNnm64AVW2OYhOCb+gBa1FYCU=
github.com/qedus/osmpbf v1.2.0/go.mod h1:Cfv6JyqTZ72BjoW9FyFBQOC2DYJbL78yw+DLhBvSH+M=
github.com/quic-go/qpack v0.6.0 h1:g7W+BMYynC1LbYLSqRt8PBg5Tgwxn214ZZR34VIOjz8=
github.com/quic-go/qpack v0.6.0/go.mod h1:lUpLKChi8njB4ty2bFLX2x4gzDqXwUpaO1DP9qMDZII=
github.com/quic-go/quic-go v0.59.0 h1:OLJkp1Mlm/aS7dpKgTc6cnpynnD2Xg7C1pwL6vy/SAw=
github.com/quic-go/quic-go v0.59.0/go.mod h1:upnsH4Ju1YkqpLXC305eW3yDZ4NfnNbmQRCMWS58IKU=
github.com/redis/go-redis/v9 v9.18.0 h1:pMkxYPkEbMPwRdenAzUNyFNrDgHx9U+DrBabWNfSRQs=
github.com/redis/go-redis/v9 v9.18.0/go.mod h1:k3ufPphLU5YXwNTUcCRXGxUoF1fqxnhFQmscfkCoDA0=
github.com/segmentio/kafka-go v0.4.51 h1:JgDPPG75tC1rWIS2Me6MwcvXJ6f49UQ4HjAOef71Hno=
github.com/segmentio/kafka-go v0.4.51/go.mod h1:Y1gn60kzLEEaW28YshXyk2+VCUKbJ3Qr6DrnT3i4+9E=
github.com/stretchr/objx v0.1.0/go.mod h1:HFkY916IF+rwdDfMAkV7OtwuqBVzrE8GR6GFx+wExME=
github.com/stretchr/objx v0.4.0/go.mod h1:YvHI0jy2hoMjB+UWwv71VJQ9isScKT/TqJzVSSt89Yw=
github.com/stretchr/objx v0.5.0/go.mod h1:Yh+to48EsGEfYuaHDzXPcE3xhTkx73EhmCGUpEOglKo=
github.com/stretchr/objx v0.5.2/go.mod h1:FRsXN1f5AsAjCGJKqEizvkpNtU+EGNCLh3NxZ/8L+MA=
github.com/stretchr/testify v1.3.0/go.mod h1:M5WIy9Dh21IEIfnGCwXGc5bZfKNJtfHm1UVUgZn+9EI=
github.com/stretchr/testify v1.7.1/go.mod h1:6Fq8oRcR53rry900zMqJjRRixrwX3KX962/h/Wwjteg=
github.com/stretchr/testify v1.8.0/go.mod h1:yNjHg4UonilssWZ8iaSj1OCr/vHnekPRkoO+kdMU+MU=
github.com/stretchr/testify v1.8.4/go.mod h1:sz/lmYIOXD/1dqDmKjjqLyZ2RngseejIcXlSw2iwfAo=
github.com/stretchr/testify v1.10.0/go.mod h1:r2ic/lqez/lEtzL7wO/rwa5dbSLXVDPFyf8C91i36aY=
github.com/stretchr/testify v1.11.1 h1:7s2iGBzp5EwR7/aIZr8ao5+dra3wiQyKjjFuvgVKu7U=
github.com/stretchr/testify v1.11.1/go.mod h1:wZwfW3scLgRK+23gO65QZefKpKQRnfz6sD981Nm4B6U=
github.com/twitchyliquid64/golang-asm v0.15.1 h1:SU5vSMR7hnwNxj24w34ZyCi/FmDZTkS4MhqMhdFk5YI=
github.com/twitchyliquid64/golang-asm v0.15.1/go.mod h1:a1lVb/DtPvCB8fslRZhAngC2+aY1QWCk3Cedj/Gdt08=
github.com/uber/h3-go/v4 v4.5.0 h1:7ruJoHCtYOCyihXfQRsPb4o6CfkhCBtVeZFM7+z1kww=
github.com/uber/h3-go/v4 v4.5.0/go.mod h1:19vfSV5HQsnRZev7V0SPmTkVSZErL7/io8M/nx+++30=
github.com/ugorji/go/codec v1.3.1 h1:waO7eEiFDwidsBN6agj1vJQ4AG7lh2yqXyOXqhgQuyY=
github.com/ugorji/go/codec v1.3.1/go.mod h1:pRBVtBSKl77K30Bv8R2P+cLSGaTtex6fsA2Wjqmfxj4=
github.com/xdg-go/pbkdf2 v1.0.0 h1:Su7DPu48wXMwC3bs7MCNG+z4FhcyEuz5dlvchbq0B0c=
github.com/xdg-go/pbkdf2 v1.0.0/go.mod h1:jrpuAogTd400dnrH08LKmI/xc1MbPOebTwRqcT5RDeI=
github.com/xdg-go/scram v1.2.0 h1:bYKF2AEwG5rqd1BumT4gAnvwU/M9nBp2pTSxeZw7Wvs=
github.com/xdg-go/scram v1.2.0/go.mod h1:3dlrS0iBaWKYVt2ZfA4cj48umJZ+cAEbR6/SjLA88I8=
github.com/xdg-go/stringprep v1.0.4 h1:XLI/Ng3O1Atzq0oBs3TWm+5ZVgkq2aqdlvP9JtoZ6c8=
github.com/xdg-go/stringprep v1.0.4/go.mod h1:mPGuuIYwz7CmR2bT9j4GbQqutWS1zV24gijq1dTyGkM=
github.com/zeebo/xxh3 v1.0.2 h1:xZmwmqxHZA8AI603jOQ0tMqmBr9lPeFwGg6d+xy9DC0=
github.com/zeebo/xxh3 v1.0.2/go.mod h1:5NWz9Sef7zIDm2JHfFlcQvNekmcEl9ekUZQQKCYaDcA=
go.mongodb.org/mongo-driver/v2 v2.5.0 h1:yXUhImUjjAInNcpTcAlPHiT7bIXhshCTL3jVBkF3xaE=
go.mongodb.org/mongo-driver/v2 v2.5.0/go.mod h1:yOI9kBsufol30iFsl1slpdq1I0eHPzybRWdyYUs8K/0=
go.uber.org/atomic v1.11.0 h1:ZvwS0R+56ePWxUNi+Atn9dWONBPp/AUETXlHW0DxSjE=
go.uber.org/atomic v1.11.0/go.mod h1:LUxbIzbOniOlMKjJjyPfpl4v+PKK2cNJn91OQbhoJI0=
go.uber.org/mock v0.6.0 h1:hyF9dfmbgIX5EfOdasqLsWD6xqpNZlXblLB/Dbnwv3Y=
go.uber.org/mock v0.6.0/go.mod h1:KiVJ4BqZJaMj4svdfmHM0AUx4NJYO8ZNpPnZn1Z+BBU=
golang.org/x/arch v0.22.0 h1:c/Zle32i5ttqRXjdLyyHZESLD/bB90DCU1g9l/0YBDI=
golang.org/x/arch v0.22.0/go.mod h1:dNHoOeKiyja7GTvF9NJS1l3Z2yntpQNzgrjh1cU103A=
golang.org/x/crypto v0.48.0 h1:/VRzVqiRSggnhY7gNRxPauEQ5Drw9haKdM0jqfcCFts=
golang.org/x/crypto v0.48.0/go.mod h1:r0kV5h3qnFPlQnBSrULhlsRfryS2pmewsg+XfMgkVos=
golang.org/x/net v0.51.0 h1:94R/GTO7mt3/4wIKpcR5gkGmRLOuE/2hNGeWq/GBIFo=
golang.org/x/net v0.51.0/go.mod h1:aamm+2QF5ogm02fjy5Bb7CQ0WMt1/WVM7FtyaTLlA9Y=
golang.org/x/sys v0.6.0/go.mod h1:oPkhp1MJrh7nUepCBck5+mAzfO9JrbApNNgaTdGDITg=
golang.org/x/sys v0.41.0 h1:Ivj+2Cp/ylzLiEU89QhWblYnOE9zerudt9Ftecq2C6k=
golang.org/x/sys v0.41.0/go.mod h1:OgkHotnGiDImocRcuBABYBEXf8A9a87e/uXjp9XT3ks=
golang.org/x/text v0.34.0 h1:oL/Qq0Kdaqxa1KbNeMKwQq0reLCCaFtqu2eNuSeNHbk=
golang.org/x/text v0.34.0/go.mod h1:homfLqTYRFyVYemLBFl5GgL/DWEiH5wcsQ5gSh1yziA=
golang.org/x/xerrors v0.0.0-20191204190536-9bdfabe68543/go.mod h1:I/5z698sn9Ka8TeJc9MKroUUfqBBauWjQqLJ2OPfmY0=
google.golang.org/protobuf v1.26.0-rc.1/go.mod h1:jlhhOSvTdKEhbULTjvd4ARK9grFBp09yW+WbY/TyQbw=
google.golang.org/protobuf v1.26.0/go.mod h1:9q0QmTI4eRPtz6boOQmLYwt+qCgq0jsYwAQnmE0givc=
google.golang.org/protobuf v1.36.10 h1:AYd7cD/uASjIL6Q9LiTjz8JLcrh/88q5UObnmY3aOOE=
google.golang.org/protobuf v1.36.10/go.mod h1:HTf+CrKn2C3g5S8VImy6tdcUvCska2kB7j23XfzDpco=
gopkg.in/check.v1 v0.0.0-20161208181325-20d25e280405/go.mod h1:Co6ibVJAznAaIkqp8huTwlJQCZ016jof/cbN4VW5Yz0=
gopkg.in/yaml.v3 v3.0.0-20200313102051-9f266ea9e77c/go.mod h1:K4uyk7z7BCEPqu6E+C64Yfv1cQ7kz7rIZviUmN+EgEM=
gopkg.in/yaml.v3 v3.0.1 h1:fxVm/GzAzEWqLHuvctI91KS9hhNmmWOoWu0XTYJS7CA=
gopkg.in/yaml.v3 v3.0.1/go.mod h1:K4uyk7z7BCEPqu6E+C64Yfv1cQ7kz7rIZviUmN+EgEM=
```

---

## data\chennai-metro.osm.pbf

**[Skipped - File Too Large (21996.9 KB)]**

---

## deployments\docker-compose.yml

```yaml
services:
  redpanda:
    image: docker.redpanda.com/redpandadata/redpanda:v24.3.18
    command: ["redpanda", "start", "--smp=1", "--memory=1G", "--reserve-memory=0M", "--overprovisioned", "--node-id=0", "--kafka-addr=PLAINTEXT://0.0.0.0:29092,OUTSIDE://0.0.0.0:9092", "--advertise-kafka-addr=PLAINTEXT://redpanda:29092,OUTSIDE://localhost:9092"]
    ports: ["9092:9092", "9644:9644"]
    volumes:
      - polaris_redpanda_data:/var/lib/redpanda/data
    healthcheck:
      test: ["CMD", "rpk", "cluster", "health", "-X", "brokers=localhost:29092", "--exit-when-healthy"]
      interval: 5s
      timeout: 5s
      retries: 20
      start_period: 10s
  redis:
    image: redis:7.4-alpine
    command: ["redis-server", "--save", ""]
    ports: ["6379:6379"]
    healthcheck:
      test: ["CMD", "redis-cli", "ping"]
      interval: 3s
      timeout: 3s
      retries: 20
  kafka-init:
    image: docker.redpanda.com/redpandadata/redpanda:v24.3.18
    depends_on:
      redpanda: { condition: service_healthy }
    entrypoint: ["/bin/sh", "-ec"]
    command:
      - |
        if ! rpk topic describe polaris.phase1.initialized -X brokers=redpanda:29092 >/dev/null 2>&1; then
          rpk topic create telemetry.ingress --partitions 3 --replicas 1 -X brokers=redpanda:29092 || rpk topic add-partitions telemetry.ingress --num 2 -X brokers=redpanda:29092
          rpk topic create telemetry.dead-letter.v1 --partitions 3 --replicas 1 -X brokers=redpanda:29092 || true
          rpk topic create polaris.phase1.initialized --partitions 1 --replicas 1 -X brokers=redpanda:29092
        fi
        rpk topic create device.lifecycle.v1 --partitions 3 --replicas 1 -X brokers=redpanda:29092 || true
        rpk topic create device.connectivity.v1 --partitions 3 --replicas 1 -X brokers=redpanda:29092 || true
        rpk topic create task.lifecycle.v1 --partitions 3 --replicas 1 -X brokers=redpanda:29092 || true
        rpk topic create device.command.v1 --partitions 3 --replicas 1 -X brokers=redpanda:29092 || true
        rpk topic create device.command.ack.v1 --partitions 3 --replicas 1 -X brokers=redpanda:29092 || true
        rpk topic create device.command.result.v1 --partitions 3 --replicas 1 -X brokers=redpanda:29092 || true
  postgres:
    image: postgis/postgis:16-3.4
    environment:
      POSTGRES_USER: polaris_user
      POSTGRES_PASSWORD: polaris_password
      POSTGRES_DB: polaris_core
    ports: ["5432:5432"]
    volumes:
      - polaris_pg_data:/var/lib/postgresql/data
      - ./init.sql:/docker-entrypoint-initdb.d/init.sql:ro
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U polaris_user -d polaris_core"]
      interval: 3s
      timeout: 3s
      retries: 20
  postgres-migrate:
    image: postgis/postgis:16-3.4
    environment:
      PGPASSWORD: polaris_password
    volumes:
      - ./init.sql:/migrations/init.sql:ro
    depends_on:
      postgres: { condition: service_healthy }
    command: ["psql", "-h", "postgres", "-U", "polaris_user", "-d", "polaris_core", "-v", "ON_ERROR_STOP=1", "-f", "/migrations/init.sql"]
  gateway:
    build:
      context: ..
      dockerfile: Dockerfile
      args: { SERVICE: gateway }
    environment:
      GATEWAY_PORT: "6080"
      KAFKA_BROKER_URL: redpanda:29092
      REDIS_URL: redis://redis:6379/0
      POSTGRES_URL: postgres://polaris_user:polaris_password@postgres:5432/polaris_core?sslmode=disable
      GATEWAY_ID: ${GATEWAY_ID:-gateway-1}
      CONNECTION_LEASE_TTL: ${CONNECTION_LEASE_TTL:-30s}
      COMMAND_RECONCILE_INTERVAL: ${COMMAND_RECONCILE_INTERVAL:-1s}
    ports: ["6080:6080"]
    depends_on:
      kafka-init: { condition: service_completed_successfully }
      redis: { condition: service_healthy }
      postgres-migrate: { condition: service_completed_successfully }
    healthcheck:
      test: ["CMD", "wget", "-qO-", "http://localhost:6080/readyz"]
      interval: 5s
      timeout: 3s
      retries: 20
  engine:
    build:
      context: ..
      dockerfile: Dockerfile
      args: { SERVICE: engine }
    environment:
      ENGINE_PORT: "6081"
      KAFKA_BROKER_URL: redpanda:29092
      REDIS_URL: redis://redis:6379/0
      POSTGRES_URL: postgres://polaris_user:polaris_password@postgres:5432/polaris_core?sslmode=disable
      DEV_PLATFORM_ADMIN_TOKEN: ${DEV_PLATFORM_ADMIN_TOKEN:-}
      DEVICE_STALE_AFTER: ${DEVICE_STALE_AFTER:-30s}
      DEVICE_OFFLINE_AFTER: ${DEVICE_OFFLINE_AFTER:-90s}
      OFFLINE_SCAN_INTERVAL: ${OFFLINE_SCAN_INTERVAL:-10s}
      CONNECTION_TICKET_TTL: ${CONNECTION_TICKET_TTL:-30s}
      OUTBOX_BATCH_SIZE: ${OUTBOX_BATCH_SIZE:-100}
      OUTBOX_POLL_INTERVAL: ${OUTBOX_POLL_INTERVAL:-500ms}
      CONNECTION_LEASE_TTL: ${CONNECTION_LEASE_TTL:-30s}
      COMMAND_ACK_TIMEOUT: ${COMMAND_ACK_TIMEOUT:-5s}
      COMMAND_RECONCILE_INTERVAL: ${COMMAND_RECONCILE_INTERVAL:-1s}
      COMMAND_MAX_ATTEMPTS: ${COMMAND_MAX_ATTEMPTS:-5}
      POLARIS_MODULE_MOBILITY_ENABLED: ${POLARIS_MODULE_MOBILITY_ENABLED:-true}
      POLARIS_MODULE_MOBILITY_REQUIRED: ${POLARIS_MODULE_MOBILITY_REQUIRED:-false}
      MOBILITY_SPATIAL_ENABLED: ${MOBILITY_SPATIAL_ENABLED:-true}
      MOBILITY_ROUTING_ENABLED: ${MOBILITY_ROUTING_ENABLED:-true}
      MOBILITY_H3_RESOLUTION: ${MOBILITY_H3_RESOLUTION:-8}
      MOBILITY_H3_SHARD_RESOLUTION: ${MOBILITY_H3_SHARD_RESOLUTION:-6}
      MOBILITY_INDEX_MIN_MOVE_METERS: ${MOBILITY_INDEX_MIN_MOVE_METERS:-5}
      MOBILITY_INDEX_MAX_AGE: ${MOBILITY_INDEX_MAX_AGE:-30s}
      MOBILITY_MAX_H3_RINGS: ${MOBILITY_MAX_H3_RINGS:-12}
      MOBILITY_MAX_SEARCH_RADIUS_METERS: ${MOBILITY_MAX_SEARCH_RADIUS_METERS:-10000}
      MOBILITY_MAX_RAW_CANDIDATES: ${MOBILITY_MAX_RAW_CANDIDATES:-50}
      MOBILITY_MAX_ROUTED_CANDIDATES: ${MOBILITY_MAX_ROUTED_CANDIDATES:-8}
      MOBILITY_MAX_ACTIVE_DEVICES_PER_TENANT: ${MOBILITY_MAX_ACTIVE_DEVICES_PER_TENANT:-10000}
      MOBILITY_ROUTING_WORKERS: ${MOBILITY_ROUTING_WORKERS:-4}
      MOBILITY_ROUTING_QUEUE_CAPACITY: ${MOBILITY_ROUTING_QUEUE_CAPACITY:-64}
      MOBILITY_ROUTING_TIMEOUT: ${MOBILITY_ROUTING_TIMEOUT:-2s}
      MOBILITY_MAX_ROUTE_EXPANSIONS: ${MOBILITY_MAX_ROUTE_EXPANSIONS:-250000}
      MOBILITY_MAX_CONCURRENT_ROUTES_PER_TENANT: ${MOBILITY_MAX_CONCURRENT_ROUTES_PER_TENANT:-2}
      MOBILITY_MAX_TRAFFIC_OBSERVATION_AGE: ${MOBILITY_MAX_TRAFFIC_OBSERVATION_AGE:-10m}
      MOBILITY_TRAFFIC_REFRESH_INTERVAL: ${MOBILITY_TRAFFIC_REFRESH_INTERVAL:-15s}
      MOBILITY_TRAFFIC_SCOPE: ${MOBILITY_TRAFFIC_SCOPE:-SHARED_TRUSTED}
      MOBILITY_ROAD_GRAPH_PATH: ${MOBILITY_ROAD_GRAPH_PATH:-data/chennai-metro.osm.pbf}
      MOBILITY_ROAD_GRAPH_VERSION: ${MOBILITY_ROAD_GRAPH_VERSION:-chennai-v1}
    ports: ["6081:6081"]
    depends_on:
      kafka-init: { condition: service_completed_successfully }
      redis: { condition: service_healthy }
      postgres-migrate: { condition: service_completed_successfully }
    healthcheck:
      test: ["CMD", "wget", "-qO-", "http://localhost:6081/readyz"]
      interval: 5s
      timeout: 3s
      retries: 30
      start_period: 20s
  frontend:
    build:
      context: ../../frontend
      dockerfile: Dockerfile
    ports: ["5173:80"]
    depends_on:
      gateway: { condition: service_healthy }
      engine: { condition: service_healthy }
    healthcheck:
      test: ["CMD", "wget", "-qO-", "http://127.0.0.1/healthz"]
      interval: 5s
      timeout: 3s
      retries: 20
volumes:
  polaris_pg_data:
  polaris_redpanda_data:
```

---

## deployments\init.sql

```sql
-- 1. Enable PostGIS (Crucial for advanced spatial math in Postgres)
CREATE EXTENSION IF NOT EXISTS postgis;

CREATE TABLE IF NOT EXISTS telemetry_history (
    id SERIAL PRIMARY KEY,
    tenant_id VARCHAR(50) NOT NULL,
    device_id VARCHAR(128) NOT NULL,
	device_boot_id VARCHAR(128) NOT NULL,
	sequence_number BIGINT NOT NULL,
	event_id VARCHAR(64) NOT NULL,
    
    -- V2 UPDATES: Changed to INT to natively store Protobuf Enums
    asset_type INT NOT NULL, 
    status INT NOT NULL,     
    
    -- RAW GPS
    lat DOUBLE PRECISION NOT NULL,
    lon DOUBLE PRECISION NOT NULL,
    
    -- V3 POSTGIS UPDATE: Native spatial column (EPSG:4326 is standard GPS)
    geom GEOMETRY(Point, 4326), 
    
    -- V3 PHYSICS UPDATES: Required for recreating traffic and handover simulations
    velocity_mps DOUBLE PRECISION,
    heading_deg DOUBLE PRECISION,
    
    battery INT,
	recorded_at TIMESTAMP NOT NULL,
	observed_at TIMESTAMPTZ NOT NULL,
	ingested_at TIMESTAMPTZ NOT NULL,
	schema_version INT NOT NULL,
	correlation_id VARCHAR(128) NOT NULL
);

-- Forward migration for Phase 0 databases. Docker entrypoint scripts only run
-- on a fresh volume; smoke-test.ps1 reapplies this idempotent file explicitly.
DO $$ BEGIN
  IF EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name='telemetry_history' AND column_name='node_id')
     AND NOT EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name='telemetry_history' AND column_name='device_id') THEN
    ALTER TABLE telemetry_history RENAME COLUMN node_id TO device_id;
  END IF;
END $$;
ALTER TABLE telemetry_history ADD COLUMN IF NOT EXISTS event_id VARCHAR(64);
ALTER TABLE telemetry_history ADD COLUMN IF NOT EXISTS device_boot_id VARCHAR(128);
ALTER TABLE telemetry_history ADD COLUMN IF NOT EXISTS sequence_number BIGINT;
ALTER TABLE telemetry_history ADD COLUMN IF NOT EXISTS observed_at TIMESTAMPTZ;
ALTER TABLE telemetry_history ADD COLUMN IF NOT EXISTS ingested_at TIMESTAMPTZ;
ALTER TABLE telemetry_history ADD COLUMN IF NOT EXISTS schema_version INT;
ALTER TABLE telemetry_history ADD COLUMN IF NOT EXISTS correlation_id VARCHAR(128);
UPDATE telemetry_history SET
  event_id=COALESCE(event_id, 'legacy:' || id),
  device_boot_id=COALESCE(device_boot_id, 'legacy'),
  sequence_number=COALESCE(sequence_number, id),
  observed_at=COALESCE(observed_at, recorded_at AT TIME ZONE 'UTC'),
  ingested_at=COALESCE(ingested_at, recorded_at AT TIME ZONE 'UTC'),
  schema_version=COALESCE(schema_version, 0),
  correlation_id=COALESCE(correlation_id, 'legacy:' || id);
ALTER TABLE telemetry_history ALTER COLUMN event_id SET NOT NULL;
ALTER TABLE telemetry_history ALTER COLUMN device_boot_id SET NOT NULL;
ALTER TABLE telemetry_history ALTER COLUMN sequence_number SET NOT NULL;
ALTER TABLE telemetry_history ALTER COLUMN observed_at SET NOT NULL;
ALTER TABLE telemetry_history ALTER COLUMN ingested_at SET NOT NULL;
ALTER TABLE telemetry_history ALTER COLUMN schema_version SET NOT NULL;
ALTER TABLE telemetry_history ALTER COLUMN correlation_id SET NOT NULL;
CREATE UNIQUE INDEX IF NOT EXISTS uq_telemetry_event_id ON telemetry_history(event_id);
CREATE UNIQUE INDEX IF NOT EXISTS uq_telemetry_device_sequence ON telemetry_history(tenant_id, device_id, device_boot_id, sequence_number);

-- Index for fast historical reporting (e.g., "Where was Drone-1001 yesterday?")
CREATE INDEX IF NOT EXISTS idx_telemetry_node_time ON telemetry_history(device_id, recorded_at DESC);

-- Index for predictive queries using standard floats
CREATE INDEX IF NOT EXISTS idx_telemetry_predictive ON telemetry_history(recorded_at, lat, lon);

-- NEW V3 INDEX: A GIST index makes PostGIS spatial queries (like bounding boxes) lightning fast
CREATE INDEX IF NOT EXISTS idx_telemetry_geom ON telemetry_history USING GIST (geom);

-- Phase 2: durable registry, credentials, audit and transactional outbox.
CREATE EXTENSION IF NOT EXISTS pgcrypto;
CREATE TABLE IF NOT EXISTS tenants (
  tenant_id TEXT PRIMARY KEY, display_name TEXT NOT NULL,
  status TEXT NOT NULL CHECK(status IN ('ACTIVE','SUSPENDED','DEACTIVATED')),
  metadata JSONB NOT NULL DEFAULT '{}', created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(), updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE TABLE IF NOT EXISTS projects (
  project_id UUID PRIMARY KEY, tenant_id TEXT NOT NULL REFERENCES tenants(tenant_id), name TEXT NOT NULL,
  description TEXT, status TEXT NOT NULL DEFAULT 'ACTIVE', metadata JSONB NOT NULL DEFAULT '{}',
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(), updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(), UNIQUE(tenant_id,name)
);
CREATE TABLE IF NOT EXISTS device_types (
  device_type_id TEXT PRIMARY KEY, display_name TEXT NOT NULL, category TEXT NOT NULL, description TEXT,
  telemetry_schema TEXT NOT NULL DEFAULT 'spatial.v1', metadata JSONB NOT NULL DEFAULT '{}'
);
CREATE TABLE IF NOT EXISTS devices (
  tenant_id TEXT NOT NULL REFERENCES tenants(tenant_id), device_id TEXT NOT NULL, project_id UUID REFERENCES projects(project_id),
  device_type_id TEXT NOT NULL REFERENCES device_types(device_type_id), display_name TEXT NOT NULL,
  lifecycle_status TEXT NOT NULL CHECK(lifecycle_status IN ('REGISTERED','ACTIVE','SUSPENDED','DECOMMISSIONED')),
  firmware_version TEXT, software_version TEXT, model_version TEXT, metadata JSONB NOT NULL DEFAULT '{}',
  registered_at TIMESTAMPTZ NOT NULL DEFAULT NOW(), updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(), deactivated_at TIMESTAMPTZ,
  PRIMARY KEY(tenant_id,device_id)
);
CREATE TABLE IF NOT EXISTS capabilities (
  capability_id TEXT PRIMARY KEY, display_name TEXT NOT NULL, description TEXT, schema JSONB NOT NULL DEFAULT '{}'
);
CREATE TABLE IF NOT EXISTS device_capabilities (
  tenant_id TEXT NOT NULL, device_id TEXT NOT NULL, capability_id TEXT NOT NULL REFERENCES capabilities(capability_id),
  configuration JSONB NOT NULL DEFAULT '{}', enabled BOOLEAN NOT NULL DEFAULT TRUE, created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  PRIMARY KEY(tenant_id,device_id,capability_id), FOREIGN KEY(tenant_id,device_id) REFERENCES devices(tenant_id,device_id)
);
CREATE INDEX IF NOT EXISTS idx_device_capabilities_tenant_capability_device ON device_capabilities(tenant_id,capability_id,device_id) WHERE enabled;
CREATE TABLE IF NOT EXISTS device_credentials (
  credential_id UUID PRIMARY KEY, tenant_id TEXT NOT NULL, device_id TEXT NOT NULL, token_prefix TEXT NOT NULL UNIQUE,
  token_hash BYTEA NOT NULL, status TEXT NOT NULL CHECK(status IN ('ACTIVE','REVOKED','EXPIRED')),
  issued_at TIMESTAMPTZ NOT NULL DEFAULT NOW(), expires_at TIMESTAMPTZ, last_used_at TIMESTAMPTZ, revoked_at TIMESTAMPTZ,
  FOREIGN KEY(tenant_id,device_id) REFERENCES devices(tenant_id,device_id)
);
CREATE TABLE IF NOT EXISTS operator_api_keys (
  api_key_id UUID PRIMARY KEY, tenant_id TEXT REFERENCES tenants(tenant_id), name TEXT NOT NULL, token_prefix TEXT NOT NULL UNIQUE,
  token_hash BYTEA NOT NULL, role TEXT NOT NULL CHECK(role IN ('PLATFORM_ADMIN','TENANT_ADMIN','OPERATOR','VIEWER')),
  status TEXT NOT NULL CHECK(status IN ('ACTIVE','REVOKED','EXPIRED')), issued_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  expires_at TIMESTAMPTZ, last_used_at TIMESTAMPTZ, revoked_at TIMESTAMPTZ
);
CREATE TABLE IF NOT EXISTS connection_tickets (
  ticket_prefix TEXT PRIMARY KEY, ticket_hash BYTEA NOT NULL, tenant_id TEXT NOT NULL, device_id TEXT NOT NULL,
  credential_id UUID NOT NULL, expires_at TIMESTAMPTZ NOT NULL, consumed_at TIMESTAMPTZ,
  FOREIGN KEY(tenant_id,device_id) REFERENCES devices(tenant_id,device_id)
);
CREATE TABLE IF NOT EXISTS operator_connection_tickets (
  ticket_prefix TEXT PRIMARY KEY, ticket_hash BYTEA NOT NULL, api_key_id UUID NOT NULL,
  tenant_id TEXT, role TEXT NOT NULL, expires_at TIMESTAMPTZ NOT NULL, consumed_at TIMESTAMPTZ,
  FOREIGN KEY(api_key_id) REFERENCES operator_api_keys(api_key_id)
);
CREATE TABLE IF NOT EXISTS outbox_events (
  outbox_id UUID PRIMARY KEY, aggregate_type TEXT NOT NULL, aggregate_id TEXT NOT NULL, tenant_id TEXT NOT NULL,
  event_id TEXT NOT NULL UNIQUE, event_type TEXT NOT NULL, schema_version INTEGER NOT NULL, payload JSONB NOT NULL,
  status TEXT NOT NULL DEFAULT 'PENDING', attempt_count INTEGER NOT NULL DEFAULT 0, next_attempt_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(), published_at TIMESTAMPTZ, last_error TEXT
);
CREATE INDEX IF NOT EXISTS idx_outbox_pending ON outbox_events(status,next_attempt_at,created_at);
CREATE TABLE IF NOT EXISTS audit_events (
  audit_id UUID PRIMARY KEY, tenant_id TEXT, actor_type TEXT NOT NULL, actor_id TEXT NOT NULL, action TEXT NOT NULL,
  resource_type TEXT NOT NULL, resource_id TEXT NOT NULL, request_id TEXT, outcome TEXT NOT NULL,
  metadata JSONB NOT NULL DEFAULT '{}', created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX IF NOT EXISTS idx_audit_tenant_time ON audit_events(tenant_id,created_at DESC);

INSERT INTO device_types(device_type_id,display_name,category,description) VALUES
 ('delivery_drone','Delivery drone','MOBILE','Spatial delivery aircraft'),
 ('ground_robot','Ground robot','MOBILE','Autonomous ground robot'),
 ('connected_vehicle','Connected vehicle','MOBILE','Connected road vehicle'),
 ('fixed_iot_sensor','Fixed IoT sensor','STATIC','Fixed telemetry sensor'),
 ('static_camera','Static camera','STATIC','Fixed spatial camera'),
 ('compute_node','Compute node','COMPUTE','Non-spatial edge compute node') ON CONFLICT DO NOTHING;
INSERT INTO capabilities(capability_id,display_name,description) VALUES
 ('navigate','Navigate','Autonomous navigation'),
 ('receive_relocation_command','Receive relocation command','Accept relocation directives'),
 ('capture_image','Capture image','Capture still imagery'),
 ('carry_payload','Carry payload','Carry a physical payload'),
 ('run_model','Run model','Execute an edge inference model'),
 ('measure_temperature','Measure temperature','Report ambient temperature') ON CONFLICT DO NOTHING;

-- Phase 3: durable task, assignment and command orchestration.
CREATE TABLE IF NOT EXISTS tasks (
  task_id UUID PRIMARY KEY, tenant_id TEXT NOT NULL REFERENCES tenants(tenant_id), project_id UUID REFERENCES projects(project_id),
  task_type TEXT NOT NULL, status TEXT NOT NULL CHECK(status IN ('PENDING','ASSIGNING','ASSIGNED','IN_PROGRESS','COMPLETED','FAILED','CANCELLED','EXPIRED')),
  priority TEXT NOT NULL CHECK(priority IN ('LOW','NORMAL','HIGH','CRITICAL')), requirements JSONB NOT NULL DEFAULT '{}', target JSONB NOT NULL DEFAULT '{}',
  assigned_device_id TEXT, correlation_id TEXT NOT NULL, created_by UUID NOT NULL REFERENCES operator_api_keys(api_key_id),
  version BIGINT NOT NULL DEFAULT 1, created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(), updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  assigned_at TIMESTAMPTZ, started_at TIMESTAMPTZ, completed_at TIMESTAMPTZ, failed_at TIMESTAMPTZ,
  expires_at TIMESTAMPTZ NOT NULL, failure_reason TEXT,
  FOREIGN KEY(tenant_id,assigned_device_id) REFERENCES devices(tenant_id,device_id)
);
CREATE INDEX IF NOT EXISTS idx_tasks_tenant_status ON tasks(tenant_id,status,priority,created_at);

CREATE TABLE IF NOT EXISTS device_assignments (
  assignment_id UUID PRIMARY KEY, tenant_id TEXT NOT NULL, device_id TEXT NOT NULL, task_id UUID NOT NULL REFERENCES tasks(task_id),
  status TEXT NOT NULL CHECK(status IN ('ACTIVE','RELEASED','EXPIRED')), lease_expires_at TIMESTAMPTZ NOT NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(), updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  FOREIGN KEY(tenant_id,device_id) REFERENCES devices(tenant_id,device_id), UNIQUE(task_id)
);
CREATE UNIQUE INDEX IF NOT EXISTS uq_device_active_assignment ON device_assignments(tenant_id,device_id) WHERE status='ACTIVE';

CREATE TABLE IF NOT EXISTS device_command_sequences (
  tenant_id TEXT NOT NULL, device_id TEXT NOT NULL, last_sequence BIGINT NOT NULL DEFAULT 0,
  PRIMARY KEY(tenant_id,device_id), FOREIGN KEY(tenant_id,device_id) REFERENCES devices(tenant_id,device_id)
);

CREATE TABLE IF NOT EXISTS commands (
  command_id UUID PRIMARY KEY, tenant_id TEXT NOT NULL, device_id TEXT NOT NULL, task_id UUID NOT NULL REFERENCES tasks(task_id),
  command_type TEXT NOT NULL, payload JSONB NOT NULL, status TEXT NOT NULL CHECK(status IN ('PENDING','DELIVERED','ACKNOWLEDGED','COMPLETED','FAILED','EXPIRED','CANCELLED')),
  sequence_number BIGINT NOT NULL, correlation_id TEXT NOT NULL, causation_id TEXT NOT NULL, attempt_count INTEGER NOT NULL DEFAULT 0,
  max_attempts INTEGER NOT NULL DEFAULT 5 CHECK(max_attempts>0), version BIGINT NOT NULL DEFAULT 1,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(), available_at TIMESTAMPTZ NOT NULL DEFAULT NOW(), sent_at TIMESTAMPTZ,
  acknowledged_at TIMESTAMPTZ, completed_at TIMESTAMPTZ, expires_at TIMESTAMPTZ NOT NULL,
  ack_status TEXT, result JSONB, last_error TEXT,
  FOREIGN KEY(tenant_id,device_id) REFERENCES devices(tenant_id,device_id), UNIQUE(tenant_id,device_id,sequence_number)
);
CREATE INDEX IF NOT EXISTS idx_commands_dispatch ON commands(status,available_at,expires_at);
CREATE INDEX IF NOT EXISTS idx_commands_device_order ON commands(tenant_id,device_id,sequence_number);
CREATE INDEX IF NOT EXISTS idx_commands_task ON commands(tenant_id,task_id);

CREATE TABLE IF NOT EXISTS command_attempts (
  attempt_id UUID PRIMARY KEY, command_id UUID NOT NULL REFERENCES commands(command_id), attempt_number INTEGER NOT NULL,
  gateway_id TEXT NOT NULL, ownership_epoch BIGINT NOT NULL, started_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  completed_at TIMESTAMPTZ, result TEXT, error TEXT, UNIQUE(command_id,attempt_number)
);
```

---

## deployments\phase2-identity-test.ps1

```
$ErrorActionPreference = "Stop"
$deploymentDir = Split-Path -Parent $MyInvocation.MyCommand.Path
$backendDir = Resolve-Path (Join-Path $deploymentDir "..")
$composeFile = Join-Path $deploymentDir "docker-compose.yml"

function New-RandomToken([string]$kind) {
  $rng=[Security.Cryptography.RandomNumberGenerator]::Create()
  try { $prefixBytes = New-Object byte[] 8; $rng.GetBytes($prefixBytes); $secretBytes = New-Object byte[] 32; $rng.GetBytes($secretBytes) } finally { $rng.Dispose() }
  $prefix=([BitConverter]::ToString($prefixBytes)-replace '-','').ToLowerInvariant()
  $secret=([BitConverter]::ToString($secretBytes)-replace '-','').ToLowerInvariant()
  return "pol_${kind}_${prefix}.${secret}"
}
function Invoke-API([string]$method,[string]$uri,$body=$null,$headers=$null) {
  $params=@{Method=$method;Uri=$uri}
  if($headers){$params.Headers=$headers}
  if($null-ne $body){$params.ContentType='application/json';$params.Body=($body|ConvertTo-Json -Depth 8)}
  return Invoke-RestMethod @params
}
function Run-IdentityCheck([string]$mode) {
  $env:IDENTITY_CHECK_MODE=$mode
  Push-Location $backendDir
  try { go run ./cmd/identitycheck } finally { Pop-Location }
  if($LASTEXITCODE-ne 0){throw "Identity check failed: $mode"}
}
function Wait-Connectivity([string]$expected,[int]$timeoutSeconds) {
  $deadline=(Get-Date).AddSeconds($timeoutSeconds)
  do { $current=Invoke-API GET "http://127.0.0.1:6081/api/v1/devices/$($env:SMOKE_DEVICE_ID)/twin" $null $headers; if($current.data.connectivity.status-eq$expected){return $current}; Start-Sleep -Milliseconds 250 } while((Get-Date)-lt$deadline)
  throw "Expected $expected, got $($current.data.connectivity.status)"
}

$env:DEVICE_STALE_AFTER="5s"
$env:DEVICE_OFFLINE_AFTER="8s"
$env:OFFLINE_SCAN_INTERVAL="1s"
& (Join-Path $deploymentDir "smoke-test.ps1")
$headers=@{Authorization="Bearer $($env:OPERATOR_TOKEN)";'X-Tenant-ID'='alpha_logistics'}

Run-IdentityCheck "basic"
$twin=Invoke-API GET "http://127.0.0.1:6081/api/v1/devices/$($env:SMOKE_DEVICE_ID)/twin" $null $headers
if($twin.data.tenant_id-ne'alpha_logistics'-or $twin.data.reported_state.id-ne $env:SMOKE_DEVICE_ID){throw "Authenticated twin composition failed"}

$oldToken=$env:DEVICE_TOKEN
$oldID=$env:DEVICE_CREDENTIAL_ID
$rotated=Invoke-API POST "http://127.0.0.1:6081/api/v1/devices/$($env:SMOKE_DEVICE_ID)/credentials/rotate" @{credential_id=$oldID} $headers
$env:DEVICE_TOKEN=$oldToken
Run-IdentityCheck "rejected"
$env:DEVICE_TOKEN=$rotated.data.secret
$env:DEVICE_CREDENTIAL_ID=$rotated.data.credential.credential_id
$env:DEVICE_BOOT_ID="phase2-rotation-$([DateTimeOffset]::UtcNow.ToUnixTimeMilliseconds())"
$env:TELEMETRY_SEQUENCE="1"
Run-IdentityCheck "send"

$env:TELEMETRY_SEQUENCE="2"
Run-IdentityCheck "revoke-session"
Run-IdentityCheck "rejected"

$offline=Wait-Connectivity "OFFLINE" 15
$recovery=Invoke-API POST "http://127.0.0.1:6081/api/v1/devices/$($env:SMOKE_DEVICE_ID)/credentials" @{} $headers
$env:DEVICE_TOKEN=$recovery.data.secret
$env:DEVICE_CREDENTIAL_ID=$recovery.data.credential.credential_id
$env:DEVICE_BOOT_ID="phase2-recovery-$([DateTimeOffset]::UtcNow.ToUnixTimeMilliseconds())"
$env:TELEMETRY_SEQUENCE="1"
Run-IdentityCheck "send"
$online=Wait-Connectivity "ONLINE" 10
Run-IdentityCheck "ticket"

Invoke-API POST "http://127.0.0.1:6081/api/v1/devices/$($env:SMOKE_DEVICE_ID)/suspend" @{} $headers | Out-Null
Run-IdentityCheck "rejected"
Invoke-API POST "http://127.0.0.1:6081/api/v1/devices/$($env:SMOKE_DEVICE_ID)/activate" @{} $headers | Out-Null
Invoke-API PATCH http://127.0.0.1:6081/api/v1/tenants/alpha_logistics @{status='SUSPENDED'} @{Authorization="Bearer $($env:OPERATOR_TOKEN)"} | Out-Null
Run-IdentityCheck "rejected"
Invoke-API PATCH http://127.0.0.1:6081/api/v1/tenants/alpha_logistics @{status='ACTIVE'} @{Authorization="Bearer $($env:OPERATOR_TOKEN)"} | Out-Null

try { Invoke-API POST http://127.0.0.1:6081/api/v1/tenants @{tenant_id='isolation_tenant';display_name='Isolation Tenant'} @{Authorization="Bearer $($env:OPERATOR_TOKEN)"}|Out-Null } catch { if($_.Exception.Response.StatusCode.value__-ne409){throw} }
$tenantBToken=New-RandomToken "op"
$tenantBPrefix=($tenantBToken.Split('.')[0].Split('_')[-1])
docker compose -f $composeFile exec -T postgres psql -U polaris_user -d polaris_core -v ON_ERROR_STOP=1 -c "INSERT INTO operator_api_keys(api_key_id,tenant_id,name,token_prefix,token_hash,role,status) VALUES(gen_random_uuid(),'isolation_tenant','isolation test','$tenantBPrefix',digest('$tenantBToken','sha256'),'TENANT_ADMIN','ACTIVE') ON CONFLICT(token_prefix) DO NOTHING" | Out-Null
try { Invoke-API GET "http://127.0.0.1:6081/api/v1/devices/$($env:SMOKE_DEVICE_ID)/twin" $null @{Authorization="Bearer $tenantBToken"}|Out-Null; throw "Cross-tenant twin was exposed" } catch { if($_.Exception.Response.StatusCode.value__-ne404){throw} }

Invoke-API POST "http://127.0.0.1:6081/api/v1/devices/$($env:SMOKE_DEVICE_ID)/decommission" @{} $headers | Out-Null
Run-IdentityCheck "rejected"

docker compose -f $composeFile run --rm postgres-migrate | Out-Null
docker compose -f $composeFile run --rm postgres-migrate | Out-Null
$plain=docker compose -f $composeFile exec -T postgres psql -U polaris_user -d polaris_core -tAc "SELECT count(*) FROM device_credentials WHERE encode(token_hash,'escape') LIKE '%pol_dev_%'"
if([int]$plain-ne0){throw "Plaintext credential found"}
$outbox=docker compose -f $composeFile exec -T postgres psql -U polaris_user -d polaris_core -tAc "SELECT count(*) FROM outbox_events WHERE tenant_id='alpha_logistics' AND status='PUBLISHED'"
if([int]$outbox-lt5){throw "Outbox relay did not publish lifecycle events"}
$audit=docker compose -f $composeFile exec -T postgres psql -U polaris_user -d polaris_core -tAc "SELECT count(*) FROM audit_events WHERE tenant_id='alpha_logistics'"
if([int]$audit-lt5){throw "Security mutations were not audited"}
foreach($group in @('polaris_engine_group','polaris_archive_group','polaris_traffic_group')) {
  $description=docker compose -f $composeFile exec -T redpanda rpk group describe $group -X brokers=localhost:29092
  if(($description|Select-String 'TOTAL-LAG\s+0').Count-ne1){throw "$group has non-zero lag"}
}

Write-Host "PASS: Phase 2 authenticated registry, credential lifecycle, tenant isolation, outbox, audit and digital-twin flow"
```

---

## deployments\phase3-command-test.ps1

```
$ErrorActionPreference = "Stop"
$deploymentDir = Split-Path -Parent $MyInvocation.MyCommand.Path
$backendDir = Resolve-Path (Join-Path $deploymentDir "..")
$composeFile = Join-Path $deploymentDir "docker-compose.yml"
$engine = "http://127.0.0.1:6081/api/v1"

function New-RandomToken([string]$kind) {
  $rng=[Security.Cryptography.RandomNumberGenerator]::Create()
  try { $prefixBytes=New-Object byte[] 8; $rng.GetBytes($prefixBytes); $secretBytes=New-Object byte[] 32; $rng.GetBytes($secretBytes) } finally { $rng.Dispose() }
  $prefix=([BitConverter]::ToString($prefixBytes)-replace '-','').ToLowerInvariant()
  $secret=([BitConverter]::ToString($secretBytes)-replace '-','').ToLowerInvariant()
  return "pol_${kind}_${prefix}.${secret}"
}
function Invoke-API([string]$method,[string]$uri,$body=$null,$headers=$null) {
  $params=@{Method=$method;Uri=$uri}
  if($headers){$params.Headers=$headers}
  if($null-ne$body){$params.ContentType='application/json';$params.Body=($body|ConvertTo-Json -Depth 10)}
  return Invoke-RestMethod @params
}
function Run-Check([string]$mode) {
  $env:ORCHESTRATION_CHECK_MODE=$mode
  Push-Location $backendDir
  try { $output=go run ./cmd/orchestrationcheck } finally { Pop-Location }
  if($LASTEXITCODE-ne0){throw "Phase 3 client failed: $mode"}
  $output | ForEach-Object { Write-Host $_ }
  return $output
}
function New-Device([string]$label,[string[]]$capabilities) {
  $stamp=[DateTimeOffset]::UtcNow.ToUnixTimeMilliseconds()
  $id="P3-$label-$stamp"
  $project=Invoke-API POST "$engine/projects" @{name="Phase 3 $label $stamp";description='Durable command proof'} $script:headers
  Invoke-API POST "$engine/devices" @{device_id=$id;project_id=$project.data.project_id;device_type_id='delivery_drone';display_name=$id} $script:headers|Out-Null
  foreach($capability in $capabilities){Invoke-API PUT "$engine/devices/$id/capabilities/$capability" @{configuration=@{}} $script:headers|Out-Null}
  Invoke-API POST "$engine/devices/$id/activate" @{} $script:headers|Out-Null
  $credential=Invoke-API POST "$engine/devices/$id/credentials" @{} $script:headers
  return [pscustomobject]@{ID=$id;Token=$credential.data.secret;CredentialID=$credential.data.credential.credential_id;ProjectID=$project.data.project_id}
}
function Select-Device($device) {
  $env:SMOKE_DEVICE_ID=$device.ID
  $env:DEVICE_TOKEN=$device.Token
  $env:DEVICE_CREDENTIAL_ID=$device.CredentialID
  $env:TASK_PROJECT_ID=$device.ProjectID
}
function Wait-CommandStatus([string]$commandID,[string]$expected,[int]$seconds=20) {
  $deadline=(Get-Date).AddSeconds($seconds)
  do{$value=Invoke-API GET "$engine/commands/$commandID" $null $script:headers;if($value.data.status-eq$expected){return $value.data};Start-Sleep -Milliseconds 200}while((Get-Date)-lt$deadline)
  throw "Command $commandID expected $expected, got $($value.data.status)"
}
function Wait-Online([string]$deviceID,[int]$seconds=15) {
  $deadline=(Get-Date).AddSeconds($seconds)
  do{$twin=Invoke-API GET "$engine/devices/$deviceID/twin" $null $script:headers;if($twin.data.connectivity.status-eq'ONLINE'){return};Start-Sleep -Milliseconds 200}while((Get-Date)-lt$deadline)
  throw "Device $deviceID did not become ONLINE before failure injection"
}

$env:COMMAND_ACK_TIMEOUT='1s'
$env:COMMAND_RECONCILE_INTERVAL='200ms'
$env:OUTBOX_POLL_INTERVAL='100ms'
$env:COMMAND_MAX_ATTEMPTS='3'
$env:CONNECTION_LEASE_TTL='6s'
$env:DEVICE_STALE_AFTER='30s'
$env:DEVICE_OFFLINE_AFTER='60s'
$env:OFFLINE_SCAN_INTERVAL='1s'

& (Join-Path $deploymentDir 'smoke-test.ps1')
$script:headers=@{Authorization="Bearer $($env:OPERATOR_TOKEN)";'X-Tenant-ID'='alpha_logistics'}
$smokeProject=docker compose -f $composeFile exec -T postgres psql -U polaris_user -d polaris_core -tAc "SELECT project_id FROM devices WHERE tenant_id='alpha_logistics' AND device_id='$($env:SMOKE_DEVICE_ID)'"
$env:TASK_PROJECT_ID=$smokeProject.Trim()
Run-Check 'complete'|Out-Null

$duplicate=New-Device 'DUPLICATE' @('navigate','receive_relocation_command')
Select-Device $duplicate
Run-Check 'duplicate'|Out-Null

$offline=New-Device 'OFFLINE' @('navigate','receive_relocation_command')
Select-Device $offline
Run-Check 'offline'|Out-Null

$fenced=New-Device 'FENCE' @('navigate','receive_relocation_command')
Select-Device $fenced
$epochBefore=docker compose -f $composeFile exec -T redis redis-cli GET "polaris:connection-epoch:alpha_logistics:$($fenced.ID)"
Run-Check 'fencing'|Out-Null
$epochAfter=docker compose -f $composeFile exec -T redis redis-cli GET "polaris:connection-epoch:alpha_logistics:$($fenced.ID)"
if([int]$epochAfter-le[int]$epochBefore){throw 'Ownership fencing epoch did not advance'}

$wrongA=New-Device 'ACK-A' @('navigate','receive_relocation_command')
$wrongB=New-Device 'ACK-B' @('navigate','receive_relocation_command')
Select-Device $wrongA
$env:DEVICE_ID_B=$wrongB.ID
$env:DEVICE_TOKEN_B=$wrongB.Token
Run-Check 'wrong-ack'|Out-Null

$mismatch=New-Device 'NO-CAMERA' @('navigate','receive_relocation_command')
Select-Device $mismatch
Run-Check 'capability-mismatch'|Out-Null

$cancelTask=Invoke-API POST "$engine/tasks" @{project_id=$mismatch.ProjectID;task_type='CAPTURE_IMAGE';priority='NORMAL';requirements=@{required_capabilities=@('capture_image');project_id=$mismatch.ProjectID};target=@{lat=13.0067;lon=80.2206};expires_at=(Get-Date).ToUniversalTime().AddMinutes(1).ToString('o')} $script:headers
Invoke-API POST "$engine/tasks/$($cancelTask.data.task.task_id)/cancel" @{} $script:headers|Out-Null
$cancelled=Invoke-API GET "$engine/tasks/$($cancelTask.data.task.task_id)" $null $script:headers
if($cancelled.data.task.status-ne'CANCELLED'){throw 'Pending task cancellation failed'}

$crash=New-Device 'CRASH' @('navigate','receive_relocation_command')
Select-Device $crash
$receiveOutput=Run-Check 'receive-no-ack'
$commandLine=$receiveOutput|Where-Object{$_-like'COMMAND_ID=*'}|Select-Object -First 1
if(-not$commandLine){throw 'Gateway crash scenario did not expose command ID'}
$env:EXPECTED_COMMAND_ID=$commandLine.Substring('COMMAND_ID='.Length)
docker compose -f $composeFile restart gateway|Out-Null
docker compose -f $composeFile up -d --wait gateway|Out-Null
Run-Check 'resume'|Out-Null

# Commit command state while Kafka is unavailable, then prove outbox recovery.
$recovery=New-Device 'OUTBOX' @('navigate','receive_relocation_command')
Select-Device $recovery
$env:IDENTITY_CHECK_MODE='send'
$env:DEVICE_BOOT_ID="phase3-outbox-$([DateTimeOffset]::UtcNow.ToUnixTimeMilliseconds())"
Push-Location $backendDir;try{go run ./cmd/identitycheck|Out-Null}finally{Pop-Location}
Wait-Online $recovery.ID
docker compose -f $composeFile stop redpanda|Out-Null
$recoveryTask=Invoke-API POST "$engine/tasks" @{project_id=$recovery.ProjectID;task_type='RELOCATE';priority='HIGH';requirements=@{required_capabilities=@('receive_relocation_command');minimum_battery=30;project_id=$recovery.ProjectID};target=@{lat=13.0068;lon=80.2207};expires_at=(Get-Date).ToUniversalTime().AddMinutes(2).ToString('o')} $script:headers
$recoveryCommand=$recoveryTask.data.command.command_id
if(-not$recoveryCommand){throw 'Command was not committed while Kafka was unavailable'}
$pendingOutbox=docker compose -f $composeFile exec -T postgres psql -U polaris_user -d polaris_core -tAc "SELECT count(*) FROM outbox_events WHERE aggregate_id='$recoveryCommand' AND status<>'PUBLISHED'"
if([int]$pendingOutbox-lt1){throw 'Command outbox was not retained during Kafka outage'}
docker compose -f $composeFile start redpanda|Out-Null
$deadline=(Get-Date).AddSeconds(60)
do{$published=docker compose -f $composeFile exec -T postgres psql -U polaris_user -d polaris_core -tAc "SELECT count(*) FROM outbox_events WHERE aggregate_id='$recoveryCommand' AND event_type='command.created.v1' AND status='PUBLISHED'";if([int]$published-ge1){break};Start-Sleep -Milliseconds 500}while((Get-Date)-lt$deadline)
if([int]$published-lt1){throw 'Outbox did not recover after Kafka restart'}

# Replaying the durable command event must not create another command row.
docker compose -f $composeFile exec -T postgres psql -U polaris_user -d polaris_core -c "UPDATE outbox_events SET status='PENDING',published_at=NULL,next_attempt_at=NOW() WHERE aggregate_id='$recoveryCommand' AND event_type='command.created.v1'"|Out-Null
Start-Sleep -Seconds 2
$duplicates=docker compose -f $composeFile exec -T postgres psql -U polaris_user -d polaris_core -tAc "SELECT count(*) FROM commands WHERE command_id='$recoveryCommand'"
if([int]$duplicates-ne1){throw 'Kafka/outbox replay duplicated a command row'}

# Cross-tenant task reads are hidden and viewer mutation is forbidden.
try{Invoke-API POST "$engine/tenants" @{tenant_id='phase3_isolation';display_name='Phase 3 isolation'} @{Authorization="Bearer $($env:OPERATOR_TOKEN)"}|Out-Null}catch{if($_.Exception.Response.StatusCode.value__-ne409){throw}}
$tenantToken=New-RandomToken 'op';$tenantPrefix=($tenantToken.Split('.')[0].Split('_')[-1])
$viewerToken=New-RandomToken 'op';$viewerPrefix=($viewerToken.Split('.')[0].Split('_')[-1])
docker compose -f $composeFile exec -T postgres psql -U polaris_user -d polaris_core -v ON_ERROR_STOP=1 -c "INSERT INTO operator_api_keys(api_key_id,tenant_id,name,token_prefix,token_hash,role,status) VALUES(gen_random_uuid(),'phase3_isolation','phase3 tenant','$tenantPrefix',digest('$tenantToken','sha256'),'TENANT_ADMIN','ACTIVE'),(gen_random_uuid(),'alpha_logistics','phase3 viewer','$viewerPrefix',digest('$viewerToken','sha256'),'VIEWER','ACTIVE') ON CONFLICT(token_prefix) DO NOTHING"|Out-Null
try{Invoke-API GET "$engine/tasks/$($recoveryTask.data.task.task_id)" $null @{Authorization="Bearer $tenantToken"}|Out-Null;throw 'Cross-tenant task was exposed'}catch{if($_.Exception.Response.StatusCode.value__-ne404){throw}}
try{Invoke-API POST "$engine/tasks" @{task_type='RELOCATE';target=@{lat=13;lon=80};expires_at=(Get-Date).ToUniversalTime().AddMinutes(1).ToString('o')} @{Authorization="Bearer $viewerToken"}|Out-Null;throw 'Viewer created a task'}catch{if($_.Exception.Response.StatusCode.value__-ne403){throw}}

foreach($group in @('polaris_engine_group','polaris_archive_group','polaris_traffic_group','polaris-command-dispatcher')){
  $deadline=(Get-Date).AddSeconds(20);do{$description=docker compose -f $composeFile exec -T redpanda rpk group describe $group -X brokers=localhost:29092;if(($description|Select-String 'TOTAL-LAG\s+0').Count-eq1){break};Start-Sleep -Milliseconds 500}while((Get-Date)-lt$deadline)
  if(($description|Select-String 'TOTAL-LAG\s+0').Count-ne1){throw "$group has non-zero lag"}
}

Write-Host 'PASS: Phase 3 durable task assignment, fenced delivery, idempotent ACK/result, retry, expiry, recovery, RBAC and tenancy flow'
```

---

## deployments\phase4-closure-soak.ps1

```
param(
  [int]$Devices = 1000,
  [int]$DurationSeconds = 45,
  [int]$Tasks = 120,
  [int]$RampPerSecond = 25,
  [int]$TelemetryIntervalSeconds = 5,
  [string]$EvidenceName = 'PHASE_4_2',
  [string]$AdminToken = '',
  [switch]$SkipCompose
)
$ErrorActionPreference = 'Stop'
if($Devices -lt 20){throw 'At least 20 devices are required for a heterogeneous soak'}
if($EvidenceName-notmatch'^[A-Za-z0-9_-]+$'){throw 'EvidenceName may contain only letters, digits, underscore, and dash'}
$deploymentDir=Split-Path -Parent $MyInvocation.MyCommand.Path
$backendDir=Resolve-Path (Join-Path $deploymentDir '..')
$rootDir=Resolve-Path (Join-Path $backendDir '..')
$composeFile=Join-Path $deploymentDir 'docker-compose.yml'
$stamp=[DateTimeOffset]::UtcNow.ToUnixTimeMilliseconds()
$tenant="phase42_soak_$stamp";$isolationTenant="phase42_isolation_$stamp";$project=[guid]::NewGuid().ToString()
$credentialFile=Join-Path $env:TEMP "polaris-phase42-devices-$stamp.json"
$sqlFile=Join-Path $env:TEMP "polaris-phase42-seed-$stamp.sql"
$stdoutFile=Join-Path $env:TEMP "polaris-phase42-output-$stamp.json"
$stderrFile=Join-Path $env:TEMP "polaris-phase42-error-$stamp.log"
$statsFile=Join-Path $rootDir "${EvidenceName}_CONTAINER_STATS.jsonl"
$resultFile=Join-Path $rootDir "${EvidenceName}_SOAK_RESULT.json"

function New-Token([string]$kind){
  $rng=[Security.Cryptography.RandomNumberGenerator]::Create()
  try{$a=New-Object byte[] 8;$b=New-Object byte[] 32;$rng.GetBytes($a);$rng.GetBytes($b)}finally{$rng.Dispose()}
  "pol_${kind}_$(([BitConverter]::ToString($a)-replace '-','').ToLowerInvariant()).$(([BitConverter]::ToString($b)-replace '-','').ToLowerInvariant())"
}
function Escape-SQL([string]$value){$value.Replace("'","''")}
function Invoke-API([string]$method,[string]$uri,$body=$null,$headers=$null){
  $p=@{Method=$method;Uri=$uri};if($headers){$p.Headers=$headers};if($null-ne$body){$p.ContentType='application/json';$p.Body=($body|ConvertTo-Json -Depth 10)};Invoke-RestMethod @p
}

$env:DEV_PLATFORM_ADMIN_TOKEN=if($AdminToken){$AdminToken}else{New-Token 'op'}
try{
  # Recreate processes and Docker port forwarding for a reproducible proof;
  # named Kafka/PostgreSQL/Redis volumes remain durable.
  if(-not $SkipCompose){
    docker compose -f $composeFile up -d --build --wait --force-recreate
    if($LASTEXITCODE-ne0){throw 'Compose stack did not become ready'}
  }else{
    $engineReady=Invoke-RestMethod http://127.0.0.1:6081/readyz;$gatewayReady=Invoke-RestMethod http://127.0.0.1:6080/readyz
    if($engineReady.status-ne'ready'-or$gatewayReady.status-ne'ready'){throw 'SkipCompose requires an already-ready stack'}
  }
  $records=New-Object System.Collections.Generic.List[object]
  $sql=New-Object System.Text.StringBuilder
  [void]$sql.AppendLine('BEGIN;')
  [void]$sql.AppendLine("INSERT INTO tenants(tenant_id,display_name,status) VALUES('$tenant','Phase 4.2 soak','ACTIVE'),('$isolationTenant','Phase 4.2 isolation','ACTIVE') ON CONFLICT DO NOTHING;")
  [void]$sql.AppendLine("INSERT INTO projects(project_id,tenant_id,name,status) VALUES('$project','$tenant','Phase 4.2 mixed soak $stamp','ACTIVE') ON CONFLICT DO NOTHING;")
  for($i=0;$i-lt$Devices;$i++){
    $ratio=$i/[double]$Devices
    if($ratio-lt.40){$deviceType='connected_vehicle';$nodeType=3;$spatial=$true;$caps=@('navigate','receive_relocation_command','run_model')}
    elseif($ratio-lt.60){$deviceType='ground_robot';$nodeType=6;$spatial=$true;$caps=@('navigate','receive_relocation_command')}
    elseif($ratio-lt.85){$deviceType='static_camera';$nodeType=7;$spatial=$true;$caps=@('capture_image')}
    else{$deviceType='compute_node';$nodeType=0;$spatial=$false;$caps=@('run_model')}
    $id="SOAK-$stamp-$('{0:D4}'-f$i)";$token=New-Token 'dev';$prefix=($token.Split('.')[0] -split '_')[-1];$credential=[guid]::NewGuid().ToString()
    [void]$sql.AppendLine("INSERT INTO devices(tenant_id,device_id,project_id,device_type_id,display_name,lifecycle_status) VALUES('$tenant','$id','$project','$deviceType','$id','ACTIVE') ON CONFLICT DO NOTHING;")
    [void]$sql.AppendLine("INSERT INTO device_credentials(credential_id,tenant_id,device_id,token_prefix,token_hash,status) VALUES('$credential','$tenant','$id','$prefix',digest('$(Escape-SQL $token)','sha256'),'ACTIVE') ON CONFLICT(token_prefix) DO NOTHING;")
    foreach($cap in $caps){[void]$sql.AppendLine("INSERT INTO device_capabilities(tenant_id,device_id,capability_id,configuration,enabled) VALUES('$tenant','$id','$cap','{}',true) ON CONFLICT(tenant_id,device_id,capability_id) DO UPDATE SET enabled=true;")}
    $records.Add(@{tenant_id=$tenant;device_id=$id;token=$token;node_type=$nodeType;spatial=$spatial})
  }
  [void]$sql.AppendLine('COMMIT;')
  Set-Content -LiteralPath $sqlFile -Value $sql.ToString() -Encoding UTF8
  $records|ConvertTo-Json -Depth 4|Set-Content -LiteralPath $credentialFile -Encoding UTF8
  Get-Content -Raw -LiteralPath $sqlFile|docker compose -f $composeFile exec -T postgres psql -q -U polaris_user -d polaris_core -v ON_ERROR_STOP=1|Out-Null
  if($LASTEXITCODE-ne0){throw 'Soak registry seed failed'}
  # This harness bulk-loads in one transaction, unlike normal incremental
  # registration. Refresh statistics so the measured query plan represents
  # the populated fleet instead of PostgreSQL's pre-seed cardinality estimate.
  docker compose -f $composeFile exec -T postgres psql -q -U polaris_user -d polaris_core -c 'ANALYZE devices; ANALYZE device_capabilities; ANALYZE device_assignments;'|Out-Null
  if($LASTEXITCODE-ne0){throw 'Soak planner statistics refresh failed'}

  Set-Content -LiteralPath $statsFile -Value '' -Encoding UTF8
  $arguments=@('run','./cmd/system-soak','-devices',$credentialFile,'-admin-token',$env:DEV_PLATFORM_ADMIN_TOKEN,'-tenant',$tenant,'-project',$project,'-duration',"${DurationSeconds}s",'-tasks',"$Tasks",'-ramp-per-second',"$RampPerSecond",'-telemetry-interval',"${TelemetryIntervalSeconds}s")
  $process=Start-Process -FilePath 'go' -ArgumentList $arguments -WorkingDirectory $backendDir -RedirectStandardOutput $stdoutFile -RedirectStandardError $stderrFile -WindowStyle Hidden -PassThru
  while(-not $process.HasExited){
    try{$engineReady=Invoke-RestMethod http://127.0.0.1:6081/readyz}catch{$engineReady=@{status='unavailable';error=$_.Exception.Message}}
    try{$gatewayReady=Invoke-RestMethod http://127.0.0.1:6080/readyz}catch{$gatewayReady=@{status='unavailable';error=$_.Exception.Message}}
    $sample=@{captured_at=(Get-Date).ToUniversalTime().ToString('o');containers=(docker stats --no-stream --format '{{json .}}');engine=$engineReady;gateway=$gatewayReady}|ConvertTo-Json -Compress -Depth 12
    Add-Content -LiteralPath $statsFile -Value $sample
    Start-Sleep -Seconds 5
    $process.Refresh()
  }
  $process.WaitForExit()
  if(-not(Test-Path -LiteralPath $stdoutFile)-or(Get-Item -LiteralPath $stdoutFile).Length-eq0){
    $detail=if(Test-Path -LiteralPath $stderrFile){Get-Content -Raw $stderrFile}else{'no process output'}
    throw "System soak produced no evidence: $detail"
  }
  Copy-Item -LiteralPath $stdoutFile -Destination $resultFile -Force
  $soak=Get-Content -Raw -LiteralPath $resultFile|ConvertFrom-Json
  if([int]$soak.counters.connections_established-ne$Devices-or[int]$soak.counters.connection_errors-ne0-or[int]$soak.counters.identity_mutations-ne0-or[int]$soak.counters.duplicate_physical_executions-ne0){
    $detail=if(Test-Path -LiteralPath $stderrFile){Get-Content -Raw $stderrFile}else{''}
    throw "System soak connection/identity invariant failed: $detail; evidence: $resultFile"
  }
  if([int]$soak.error_totals.unexpected-ne0-or[int]$soak.error_totals.server_error-ne0-or[int]$soak.error_totals.transport_error-ne0){throw "Unexpected classified workload errors: $($soak.error_totals|ConvertTo-Json -Compress); evidence: $resultFile"}
  if([int64]$soak.counters.telemetry_sent-lt1-or[int]$soak.counters.task_requests_attempted-ne$Tasks-or[int]$soak.counters.physical_executions-ne$Tasks-or[int]$soak.counters.commands_delivered-lt$Tasks){throw "System soak workload did not exercise every required path; evidence: $resultFile"}

  foreach($group in @('polaris_engine_group','polaris_archive_group','polaris_traffic_group','polaris-command-dispatcher')){
    $deadline=(Get-Date).AddSeconds(90)
    do{$description=docker compose -f $composeFile exec -T redpanda rpk group describe $group -X brokers=localhost:29092;if(($description|Select-String 'TOTAL-LAG\s+0').Count-eq1){break};Start-Sleep -Milliseconds 500}while((Get-Date)-lt$deadline)
    if(($description|Select-String 'TOTAL-LAG\s+0').Count-ne1){throw "$group lag did not return to zero"}
  }
  $invariants=docker compose -f $composeFile exec -T postgres psql -U polaris_user -d polaris_core -tA -F ',' -c "SELECT (SELECT count(*) FROM telemetry_history WHERE tenant_id='$tenant'),(SELECT count(*) FROM (SELECT tenant_id,device_id,count(*) FROM device_assignments WHERE tenant_id='$tenant' AND status='ACTIVE' GROUP BY tenant_id,device_id HAVING count(*)>1)x),(SELECT count(*) FROM commands c JOIN tasks t USING(task_id) WHERE c.tenant_id='$tenant' AND (c.tenant_id<>t.tenant_id OR c.device_id<>t.assigned_device_id)),(SELECT count(*) FROM commands c JOIN tasks t USING(task_id) WHERE c.tenant_id='$tenant' AND t.requirements->>'planning_mode'='POLARIS_REQUIRED' AND (NOT(c.payload?'road_graph_version') OR NOT(c.payload?'routing_snapshot_version'))),(SELECT count(*) FROM tasks WHERE tenant_id='$tenant');"
  $parts=($invariants|Select-Object -Last 1).Trim().Split(',')
  if([int64]$parts[0]-lt1-or[int]$parts[1]-ne0-or[int]$parts[2]-ne0-or[int]$parts[3]-ne0-or[int]$parts[4]-lt$Tasks){throw "Database invariant failure: $invariants"}
  $headers=@{Authorization="Bearer $($env:DEV_PLATFORM_ADMIN_TOKEN)";'X-Tenant-ID'=$isolationTenant}
  $isolated=Invoke-API GET 'http://127.0.0.1:6081/api/v1/spatial/devices/nearby?lat=13.0067&lon=80.2206&radius_meters=5000&limit=20' $null $headers
  if($isolated.data.count-ne0){throw 'Cross-tenant spatial leakage detected'}
  $ready=Invoke-RestMethod http://127.0.0.1:6081/readyz
  if($ready.modules.mobility.details.routing_runtime.queue_depth-ge$ready.modules.mobility.details.routing_runtime.queue_capacity){throw 'Routing queue remained saturated'}
  $soak|Add-Member -NotePropertyName final_runtime -NotePropertyValue @{engine=$ready;gateway=(Invoke-RestMethod http://127.0.0.1:6080/readyz)} -Force
  $soak|ConvertTo-Json -Depth 20|Set-Content -LiteralPath $resultFile -Encoding UTF8
  Write-Host "PASS: $Devices-device heterogeneous full-system soak; evidence: $resultFile and $statsFile"
}finally{
  foreach($path in @($credentialFile,$sqlFile,$stdoutFile,$stderrFile)){if(Test-Path -LiteralPath $path){Remove-Item -LiteralPath $path -Force}}
}
```

---

## deployments\phase4-mobility-rebuild-test.ps1

```
$ErrorActionPreference='Stop'
$deploymentDir=Split-Path -Parent $MyInvocation.MyCommand.Path
$backendDir=Resolve-Path (Join-Path $deploymentDir '..')
$rootDir=Resolve-Path (Join-Path $backendDir '..')
$composeFile=Join-Path $deploymentDir 'docker-compose.yml'
$evidenceFile=Join-Path $rootDir 'PHASE_4_2_MOBILITY_REBUILD_RESULT.json'
$engine='http://127.0.0.1:6081/api/v1'

function Invoke-API([string]$method,[string]$uri,$body=$null,$headers=$null){$p=@{Method=$method;Uri=$uri};if($headers){$p.Headers=$headers};if($null-ne$body){$p.ContentType='application/json';$p.Body=($body|ConvertTo-Json -Depth 10)};Invoke-RestMethod @p}
function New-Tenant([string]$id,[string]$name){try{Invoke-API POST "$engine/tenants" @{tenant_id=$id;display_name=$name} @{Authorization="Bearer $($env:OPERATOR_TOKEN)"}|Out-Null}catch{if($_.Exception.Response.StatusCode.value__-ne409){throw}}}
function New-SpatialDevice([string]$tenant,[string]$project,[string]$id,[string]$type,[int]$nodeType,[double]$lat,[double]$lon){
  $headers=@{Authorization="Bearer $($env:OPERATOR_TOKEN)";'X-Tenant-ID'=$tenant}
  Invoke-API POST "$engine/devices" @{device_id=$id;project_id=$project;device_type_id=$type;display_name=$id} $headers|Out-Null
  Invoke-API PUT "$engine/devices/$id/capabilities/navigate" @{configuration=@{}} $headers|Out-Null
  Invoke-API POST "$engine/devices/$id/activate" @{} $headers|Out-Null
  $secret=(Invoke-API POST "$engine/devices/$id/credentials" @{} $headers).data.secret
  $env:SMOKE_TENANT_ID=$tenant;$env:SMOKE_DEVICE_ID=$id;$env:DEVICE_TOKEN=$secret;$env:SMOKE_NODE_TYPE="$nodeType";$env:SMOKE_LAT="$lat";$env:SMOKE_LON="$lon";$env:SMOKE_BOOT_ID="rebuild-$id";$env:SMOKE_SEQUENCE='1';$env:SMOKE_WAIT_FOR_PROJECTION='true';$env:SMOKE_WAIT_FOR_MATCH='false'
  Push-Location $backendDir
  try{go run ./cmd/smoke|Out-Host;if($LASTEXITCODE-ne0){throw "Telemetry failed for $id"}}finally{Pop-Location}
  return $secret
}
function Nearby([string]$tenant,[double]$lat,[double]$lon){$headers=@{Authorization="Bearer $($env:OPERATOR_TOKEN)";'X-Tenant-ID'=$tenant};(Invoke-API GET "$engine/spatial/devices/nearby?lat=$lat&lon=$lon&radius_meters=1000&limit=50" $null $headers).data.devices}

. (Join-Path $deploymentDir 'smoke-test.ps1')
$stamp=[DateTimeOffset]::UtcNow.ToUnixTimeMilliseconds();$tenant="phase42_rebuild_$stamp";$other="phase42_rebuild_other_$stamp"
New-Tenant $tenant 'Phase 4.2 rebuild';New-Tenant $other 'Phase 4.2 rebuild isolation'
$headers=@{Authorization="Bearer $($env:OPERATOR_TOKEN)";'X-Tenant-ID'=$tenant};$otherHeaders=@{Authorization="Bearer $($env:OPERATOR_TOKEN)";'X-Tenant-ID'=$other}
$project=(Invoke-API POST "$engine/projects" @{name='Rebuild devices';description='Derived projection restart proof'} $headers).data.project_id
$otherProject=(Invoke-API POST "$engine/projects" @{name='Isolation devices';description='Rebuild tenancy proof'} $otherHeaders).data.project_id
$road="REBUILD-ROAD-$stamp";$robot="REBUILD-ROBOT-$stamp";$static="REBUILD-STATIC-$stamp";$inactive="REBUILD-INACTIVE-$stamp";$foreign="REBUILD-FOREIGN-$stamp";$compute="REBUILD-COMPUTE-$stamp"
$roadSecret=New-SpatialDevice $tenant $project $road 'connected_vehicle' 3 13.2000 80.3000
New-SpatialDevice $tenant $project $robot 'ground_robot' 6 13.2010 80.3000|Out-Null
New-SpatialDevice $tenant $project $static 'static_camera' 7 13.2020 80.3000|Out-Null
New-SpatialDevice $tenant $project $inactive 'connected_vehicle' 3 13.2030 80.3000|Out-Null
New-SpatialDevice $other $otherProject $foreign 'connected_vehicle' 3 13.2005 80.3000|Out-Null
Invoke-API POST "$engine/devices" @{device_id=$compute;project_id=$project;device_type_id='compute_node';display_name=$compute} $headers|Out-Null
Invoke-API POST "$engine/devices/$compute/activate" @{} $headers|Out-Null
Invoke-API POST "$engine/devices/$inactive/suspend" @{} $headers|Out-Null

# Advance the road device, then inject a stale sequence without waiting for a
# projection event. Redis and Mobility must both retain sequence 2.
$env:SMOKE_TENANT_ID=$tenant;$env:SMOKE_DEVICE_ID=$road;$env:DEVICE_TOKEN=$roadSecret;$env:SMOKE_NODE_TYPE='3';$env:SMOKE_LAT='13.2002';$env:SMOKE_LON='80.3001';$env:SMOKE_BOOT_ID="rebuild-$road";$env:SMOKE_SEQUENCE='2';$env:SMOKE_WAIT_FOR_PROJECTION='true';$env:SMOKE_WAIT_FOR_MATCH='false'
Push-Location $backendDir;try{go run ./cmd/smoke|Out-Host;if($LASTEXITCODE-ne0){throw 'Newer road telemetry failed'}}finally{Pop-Location}
$env:SMOKE_LAT='13.2500';$env:SMOKE_LON='80.3500';$env:SMOKE_SEQUENCE='1';$env:SMOKE_WAIT_FOR_PROJECTION='false'
Push-Location $backendDir;try{go run ./cmd/smoke|Out-Host;if($LASTEXITCODE-ne0){throw 'Stale road telemetry injection failed'}}finally{Pop-Location}

$before=Nearby $tenant 13.2010 80.3000;$beforeRoad=$before|Where-Object{$_.state.device_id-eq$road}|Select-Object -First 1
if((-not $beforeRoad)-or[uint64]$beforeRoad.state.sequence_number-ne2){throw 'Pre-restart Mobility state did not retain the newer sequence'}
docker compose -f $composeFile restart engine|Out-Null
$sawNotReady=$false;$deadline=(Get-Date).AddMinutes(3);$ready=$null
do{
  try{$probe=Invoke-RestMethod http://127.0.0.1:6081/readyz;if($probe.modules.mobility.state-eq'READY'){$ready=$probe;break}else{$sawNotReady=$true}}catch{$sawNotReady=$true}
  Start-Sleep -Milliseconds 250
}while((Get-Date)-lt$deadline)
if((-not $sawNotReady)-or(-not $ready)){throw 'Readiness did not transition from unavailable/not-ready to Mobility READY during rebuild'}
$after=Nearby $tenant 13.2010 80.3000;$ids=@($after|ForEach-Object{$_.state.device_id});$afterRoad=$after|Where-Object{$_.state.device_id-eq$road}|Select-Object -First 1
foreach($id in @($road,$robot,$static)){if(($ids|Where-Object{$_-eq$id}).Count-ne1){throw "Active rebuilt device $id was missing or duplicated"}}
foreach($id in @($inactive,$foreign,$compute)){if($ids-contains$id){throw "Excluded device $id leaked into rebuilt spatial state"}}
if([uint64]$afterRoad.state.sequence_number-ne2-or[double]$afterRoad.state.reported_position.latitude-ne13.2002){throw 'Rebuild regressed the accepted road boot/sequence/position'}
$otherView=Nearby $other 13.2010 80.3000;$otherIDs=@($otherView|ForEach-Object{$_.state.device_id})
if($otherIDs-notcontains$foreign-or$otherIDs-contains$road){throw 'Rebuild tenant isolation failed'}
$result=@{measured_at=(Get-Date).ToUniversalTime().ToString('o');readiness_transition_observed=$sawNotReady;mobility_state=$ready.modules.mobility.state;recovered=@($road,$robot,$static);excluded=@($inactive,$compute,$foreign);road_sequence_before=[uint64]$beforeRoad.state.sequence_number;road_sequence_after=[uint64]$afterRoad.state.sequence_number;duplicates=0;cross_tenant=0}|ConvertTo-Json -Depth 10
$result|Set-Content -LiteralPath $evidenceFile -Encoding UTF8
Write-Host "PASS: Mobility restart/rebuild, freshness, exclusion, uniqueness, and tenancy; evidence: $evidenceFile"
```

---

## deployments\phase4-mobility-test.ps1

```
param([switch]$FullRegression,[switch]$SkipLocalChecks)
$ErrorActionPreference = 'Stop'
$deploymentDir = Split-Path -Parent $MyInvocation.MyCommand.Path
$backendDir = Resolve-Path (Join-Path $deploymentDir '..')
$composeFile = Join-Path $deploymentDir 'docker-compose.yml'
$engine = 'http://127.0.0.1:6081/api/v1'

function Invoke-API([string]$method,[string]$uri,$body=$null,$headers=$null) {
  $params=@{Method=$method;Uri=$uri}
  if($headers){$params.Headers=$headers}
  if($null-ne$body){$params.ContentType='application/json';$params.Body=($body|ConvertTo-Json -Depth 12)}
  Invoke-RestMethod @params
}

if(-not $SkipLocalChecks){Push-Location $backendDir
try { go test ./...; if($LASTEXITCODE-ne0){throw 'Go tests failed'}; go vet ./...; if($LASTEXITCODE-ne0){throw 'Go vet failed'} } finally { Pop-Location }}
& (Join-Path $deploymentDir 'smoke-test.ps1')
$headers=@{Authorization="Bearer $($env:OPERATOR_TOKEN)";'X-Tenant-ID'='alpha_logistics'}

$ready=Invoke-RestMethod http://127.0.0.1:6081/readyz
if($ready.core-ne'READY'-or $ready.modules.mobility.state-notin @('READY','DEGRADED')){throw 'Mobility readiness was not exposed'}
$twin=Invoke-API GET "$engine/devices/$($env:SMOKE_DEVICE_ID)/twin" $null $headers
if(-not $twin.data.components.'spatial/v1'-or-not $twin.data.components.'battery/v1'){throw 'Versioned twin components missing'}
$genericSmokeEnvironment=@{
  SMOKE_DEVICE_ID=$env:SMOKE_DEVICE_ID
  DEVICE_TOKEN=$env:DEVICE_TOKEN
  SMOKE_NODE_TYPE=$env:SMOKE_NODE_TYPE
  SMOKE_LAT=$env:SMOKE_LAT
  SMOKE_LON=$env:SMOKE_LON
}

$route=Invoke-API POST "$engine/routes" @{mobility_profile='ROAD_VEHICLE';origin=@{latitude=13.0067;longitude=80.2206};destination=@{latitude=13.02;longitude=80.23};policy='FASTEST'} $headers
if(-not $route.data.route_id-or -not $route.data.road_graph_version-or $route.data.snapshot_version-lt1-or $route.data.waypoints.Count-lt1){throw 'Versioned A* route was not returned'}

$stamp=[DateTimeOffset]::UtcNow.ToUnixTimeMilliseconds();$vehicle="P4-ROAD-$stamp"
$project=Invoke-API POST "$engine/projects" @{name="Phase 4 road proof $stamp";description='Mobility planner proof'} $headers
Invoke-API POST "$engine/devices" @{device_id=$vehicle;project_id=$project.data.project_id;device_type_id='connected_vehicle';display_name=$vehicle} $headers|Out-Null
Invoke-API PUT "$engine/devices/$vehicle/capabilities/navigate" @{configuration=@{}} $headers|Out-Null
Invoke-API POST "$engine/devices/$vehicle/activate" @{} $headers|Out-Null
$credential=Invoke-API POST "$engine/devices/$vehicle/credentials" @{} $headers
$env:SMOKE_DEVICE_ID=$vehicle;$env:DEVICE_TOKEN=$credential.data.secret;$env:SMOKE_NODE_TYPE='3';$env:SMOKE_LAT='13.04123';$env:SMOKE_LON='80.23876'
Push-Location $backendDir
try { go run ./cmd/smoke | Out-Host; if($LASTEXITCODE-ne0){throw 'Road telemetry proof failed'} } finally { Pop-Location }
$near=Invoke-API GET "$engine/spatial/devices/nearby?lat=$($env:SMOKE_LAT)&lon=$($env:SMOKE_LON)&radius_meters=25&limit=20" $null $headers
$nearbyDeviceIDs=@($near.data.devices | ForEach-Object { $_.state.device_id })
if($nearbyDeviceIDs -notcontains $vehicle){throw "Mobility nearby query did not return authenticated telemetry device '$vehicle'; returned [$($nearbyDeviceIDs -join ', ')]"}
$task=Invoke-API POST "$engine/tasks" @{task_type='NAVIGATE';priority='HIGH';requirements=@{required_capabilities=@('navigate');minimum_battery=20;max_distance_meters=10000;planning_mode='POLARIS_REQUIRED'};target=@{lat=13.02;lon=80.23;policy='FASTEST'};expires_at=(Get-Date).ToUniversalTime().AddMinutes(5).ToString('o')} $headers
if(-not $task.data.command.payload.route_id-or -not $task.data.command.payload.road_graph_version-or $task.data.command.payload.routing_snapshot_version-lt1){throw 'NAVIGATE did not persist graph and routing snapshot identity'}

foreach($name in $genericSmokeEnvironment.Keys){[Environment]::SetEnvironmentVariable($name,$genericSmokeEnvironment[$name],'Process')}
if($FullRegression){ & (Join-Path $deploymentDir 'reliability-test.ps1'); & (Join-Path $deploymentDir 'phase2-identity-test.ps1'); & (Join-Path $deploymentDir 'phase3-command-test.ps1') }
Write-Host "PASS: Phase 4 H3/R-tree spatial state, A* route snapshot, and durable NAVIGATE planning ($vehicle)"
```

---

## deployments\phase4-module-isolation-test.ps1

```
$ErrorActionPreference='Stop'
$deploymentDir=Split-Path -Parent $MyInvocation.MyCommand.Path
$backendDir=Resolve-Path (Join-Path $deploymentDir '..')
$composeFile=Join-Path $deploymentDir 'docker-compose.yml'
$engine='http://127.0.0.1:6081/api/v1'
function New-Token(){ $rng=[Security.Cryptography.RandomNumberGenerator]::Create();try{$a=New-Object byte[] 8;$b=New-Object byte[] 32;$rng.GetBytes($a);$rng.GetBytes($b)}finally{$rng.Dispose()};"pol_op_$(([BitConverter]::ToString($a)-replace '-','').ToLowerInvariant()).$(([BitConverter]::ToString($b)-replace '-','').ToLowerInvariant())" }
function API([string]$method,[string]$uri,$body=$null,$headers=$null){$p=@{Method=$method;Uri=$uri};if($headers){$p.Headers=$headers};if($null-ne$body){$p.ContentType='application/json';$p.Body=($body|ConvertTo-Json -Depth 10)};Invoke-RestMethod @p}
$env:DEV_PLATFORM_ADMIN_TOKEN=New-Token;$env:POLARIS_MODULE_MOBILITY_ENABLED='false'
try{
  docker compose -f $composeFile up -d --build --force-recreate --wait engine gateway
  if($LASTEXITCODE-ne0){throw 'Mobility-disabled engine failed readiness'}
  $ready=Invoke-RestMethod http://127.0.0.1:6081/readyz
  if($ready.core-ne'READY'-or $ready.modules.mobility){throw 'Mobility was not disabled cleanly'}
  $headers=@{Authorization="Bearer $($env:DEV_PLATFORM_ADMIN_TOKEN)";'X-Tenant-ID'='alpha_logistics'}
  try{API POST "$engine/tenants" @{tenant_id='alpha_logistics';display_name='Alpha Logistics'} @{Authorization="Bearer $($env:DEV_PLATFORM_ADMIN_TOKEN)"}|Out-Null}catch{if($_.Exception.Response.StatusCode.value__-ne409){throw}}
  $stamp=[DateTimeOffset]::UtcNow.ToUnixTimeMilliseconds();$device="P4-CAMERA-$stamp"
  $project=API POST "$engine/projects" @{name="Phase 4 generic proof $stamp"} $headers
  API POST "$engine/devices" @{device_id=$device;project_id=$project.data.project_id;device_type_id='static_camera';display_name=$device} $headers|Out-Null
  API PUT "$engine/devices/$device/capabilities/capture_image" @{configuration=@{}} $headers|Out-Null
  API POST "$engine/devices/$device/activate" @{} $headers|Out-Null
  $credential=API POST "$engine/devices/$device/credentials" @{} $headers
  $env:OPERATOR_TOKEN=$env:DEV_PLATFORM_ADMIN_TOKEN;$env:DEVICE_TOKEN=$credential.data.secret;$env:SMOKE_DEVICE_ID=$device;$env:SMOKE_NODE_TYPE='7'
  Push-Location $backendDir;try{go run ./cmd/smoke|Out-Host;if($LASTEXITCODE-ne0){throw 'Generic telemetry failed with Mobility disabled'}}finally{Pop-Location}
  $task=API POST "$engine/tasks" @{project_id=$project.data.project_id;task_type='CAPTURE_IMAGE';priority='NORMAL';requirements=@{required_capabilities=@('capture_image');project_id=$project.data.project_id};target=@{image_profile='overview'};expires_at=(Get-Date).ToUniversalTime().AddMinutes(2).ToString('o')} $headers
  if(-not $task.data.command.command_id-or $task.data.command.payload.image_profile-ne'overview'){throw 'Default planner failed with Mobility disabled'}
  Write-Host "PASS: Mobility disabled; telemetry, twin state, CAPTURE_IMAGE and durable generic command remain functional ($device)"
}finally{
  $env:POLARIS_MODULE_MOBILITY_ENABLED='true';docker compose -f $composeFile up -d --force-recreate --wait engine|Out-Null
}
```

---

## deployments\phase4-routing-overload-test.ps1

```
$ErrorActionPreference='Stop'
$deploymentDir=Split-Path -Parent $MyInvocation.MyCommand.Path
$backendDir=Resolve-Path (Join-Path $deploymentDir '..')
$rootDir=Resolve-Path (Join-Path $backendDir '..')
$composeFile=Join-Path $deploymentDir 'docker-compose.yml'
$evidenceFile=Join-Path $rootDir 'PHASE_4_2_ROUTING_OVERLOAD_RESULT.json'
$stamp=[DateTimeOffset]::UtcNow.ToUnixTimeMilliseconds();$tenant="phase42_overload_$stamp";$device="OVERLOAD-$stamp"

function New-Token([string]$kind){$rng=[Security.Cryptography.RandomNumberGenerator]::Create();try{$a=New-Object byte[] 8;$b=New-Object byte[] 32;$rng.GetBytes($a);$rng.GetBytes($b)}finally{$rng.Dispose()};"pol_${kind}_$(([BitConverter]::ToString($a)-replace '-','').ToLowerInvariant()).$(([BitConverter]::ToString($b)-replace '-','').ToLowerInvariant())"}
function Invoke-API([string]$method,[string]$uri,$body=$null,$headers=$null){$p=@{Method=$method;Uri=$uri};if($headers){$p.Headers=$headers};if($null-ne$body){$p.ContentType='application/json';$p.Body=($body|ConvertTo-Json -Depth 10)};Invoke-RestMethod @p}

$saved=@{workers=$env:MOBILITY_ROUTING_WORKERS;queue=$env:MOBILITY_ROUTING_QUEUE_CAPACITY;tenant=$env:MOBILITY_MAX_CONCURRENT_ROUTES_PER_TENANT;admin=$env:DEV_PLATFORM_ADMIN_TOKEN}
try{
  $env:MOBILITY_ROUTING_WORKERS='1';$env:MOBILITY_ROUTING_QUEUE_CAPACITY='4';$env:MOBILITY_MAX_CONCURRENT_ROUTES_PER_TENANT='128';$env:DEV_PLATFORM_ADMIN_TOKEN=New-Token 'op'
  docker compose -f $composeFile up -d --build --wait --force-recreate
  if($LASTEXITCODE-ne0){throw 'Overload stack failed to become ready'}
  $adminHeaders=@{Authorization="Bearer $($env:DEV_PLATFORM_ADMIN_TOKEN)"}
  Invoke-API POST 'http://127.0.0.1:6081/api/v1/tenants' @{tenant_id=$tenant;display_name='Phase 4.2 overload'} $adminHeaders|Out-Null
  $headers=@{Authorization="Bearer $($env:DEV_PLATFORM_ADMIN_TOKEN)";'X-Tenant-ID'=$tenant}
  $project=(Invoke-API POST 'http://127.0.0.1:6081/api/v1/projects' @{name='Overload isolation';description='Routing must not block core'} $headers).data.project_id
  Invoke-API POST 'http://127.0.0.1:6081/api/v1/devices' @{device_id=$device;project_id=$project;device_type_id='connected_vehicle';display_name=$device} $headers|Out-Null
  Invoke-API PUT "http://127.0.0.1:6081/api/v1/devices/$device/capabilities/run_model" @{configuration=@{}} $headers|Out-Null
  Invoke-API POST "http://127.0.0.1:6081/api/v1/devices/$device/activate" @{} $headers|Out-Null
  $secret=(Invoke-API POST "http://127.0.0.1:6081/api/v1/devices/$device/credentials" @{} $headers).data.secret
  Push-Location $backendDir
  try{$result=go run ./cmd/routing-overload -admin-token $env:DEV_PLATFORM_ADMIN_TOKEN -tenant $tenant -project $project -device $device -device-token $secret;if($LASTEXITCODE-ne0){throw 'Routing overload harness failed'}}finally{Pop-Location}
  $result|Set-Content -LiteralPath $evidenceFile -Encoding UTF8
  Write-Host "PASS: bounded routing overload, unrelated telemetry/task delivery, and recovery; evidence: $evidenceFile"
}finally{
  $env:MOBILITY_ROUTING_WORKERS=$saved.workers;$env:MOBILITY_ROUTING_QUEUE_CAPACITY=$saved.queue;$env:MOBILITY_MAX_CONCURRENT_ROUTES_PER_TENANT=$saved.tenant;$env:DEV_PLATFORM_ADMIN_TOKEN=$saved.admin
  docker compose -f $composeFile up -d --wait --force-recreate engine|Out-Null
}
```

---

## deployments\reliability-test.ps1

```
$ErrorActionPreference = "Stop"
$deploymentDir = Split-Path -Parent $MyInvocation.MyCommand.Path
$backendDir = Resolve-Path (Join-Path $deploymentDir "..")
$rootDir = Resolve-Path (Join-Path $backendDir "..")

& (Join-Path $deploymentDir "smoke-test.ps1")

Push-Location $backendDir
try {
  go test ./...
  if ($LASTEXITCODE -ne 0) { throw "Unit reliability tests failed" }
  $env:REDIS_URL = "redis://127.0.0.1:6379/0"
  $env:POSTGRES_URL = "postgres://polaris_user:polaris_password@127.0.0.1:5432/polaris_core?sslmode=disable"
  go test -count=1 -tags=integration -v ./internal/application/stream
  if ($LASTEXITCODE -ne 0) { throw "Live dependency reliability tests failed" }
} finally { Pop-Location }

$gatewayHealth = Invoke-RestMethod http://127.0.0.1:6080/healthz
$gatewayReady = Invoke-RestMethod http://127.0.0.1:6080/readyz
$engineHealth = Invoke-RestMethod http://127.0.0.1:6081/healthz
$engineReady = Invoke-RestMethod http://127.0.0.1:6081/readyz
if ($gatewayHealth.status -ne "live" -or $gatewayReady.status -ne "ready" -or $engineHealth.status -ne "live" -or $engineReady.status -ne "ready") {
  throw "Health/readiness contract failed"
}

$partitions = docker compose -f (Join-Path $deploymentDir "docker-compose.yml") exec -T redpanda `
  rpk topic describe telemetry.ingress -p -X brokers=localhost:29092
if (($partitions | Select-String -Pattern '^\s*[0-2]\s+').Count -lt 3) { throw "telemetry.ingress does not have three partitions" }

$constraints = docker compose -f (Join-Path $deploymentDir "docker-compose.yml") exec -T postgres `
  psql -U polaris_user -d polaris_core -tAc "SELECT count(*) FROM pg_indexes WHERE tablename='telemetry_history' AND indexname IN ('uq_telemetry_event_id','uq_telemetry_device_sequence')"
if ([int]$constraints -ne 2) { throw "PostgreSQL idempotency indexes are missing" }

Write-Host "PASS: Phase 1 at-least-once/idempotent reliability verification"
```

---

## deployments\smoke-test.ps1

```
$ErrorActionPreference = "Stop"
$deploymentDir = Split-Path -Parent $MyInvocation.MyCommand.Path
$backendDir = Resolve-Path (Join-Path $deploymentDir "..")

function New-RandomToken([string]$kind) {
  $rng=[Security.Cryptography.RandomNumberGenerator]::Create()
  try { $prefixBytes = New-Object byte[] 8; $rng.GetBytes($prefixBytes); $secretBytes = New-Object byte[] 32; $rng.GetBytes($secretBytes) } finally { $rng.Dispose() }
  $prefix = ([BitConverter]::ToString($prefixBytes) -replace '-','').ToLowerInvariant()
  $secret = ([BitConverter]::ToString($secretBytes) -replace '-','').ToLowerInvariant()
  return "pol_${kind}_${prefix}.${secret}"
}

$env:DEV_PLATFORM_ADMIN_TOKEN = New-RandomToken "op"
$smokeID = "SMOKE-$([DateTimeOffset]::UtcNow.ToUnixTimeMilliseconds())"

docker compose -f (Join-Path $deploymentDir "docker-compose.yml") up -d --build --wait
if ($LASTEXITCODE -ne 0) { throw "Docker Compose stack failed to start" }
$operatorHeaders = @{ Authorization = "Bearer $($env:DEV_PLATFORM_ADMIN_TOKEN)"; "X-Tenant-ID" = "alpha_logistics" }
try {
  Invoke-RestMethod -Method Post -Uri http://127.0.0.1:6081/api/v1/tenants -Headers @{ Authorization = "Bearer $($env:DEV_PLATFORM_ADMIN_TOKEN)" } -ContentType application/json -Body '{"tenant_id":"alpha_logistics","display_name":"Alpha Logistics"}' | Out-Null
} catch { if ($_.Exception.Response.StatusCode.value__ -ne 409) { throw } }
$projectResponse = Invoke-RestMethod -Method Post -Uri http://127.0.0.1:6081/api/v1/projects -Headers $operatorHeaders -ContentType application/json -Body (@{name="Chennai fleet demo $smokeID";description="Phase 2 reproducible identity proof"}|ConvertTo-Json)
$deviceBody = @{ device_id=$smokeID; project_id=$projectResponse.data.project_id; device_type_id="delivery_drone"; display_name="Authenticated smoke device" } | ConvertTo-Json
Invoke-RestMethod -Method Post -Uri http://127.0.0.1:6081/api/v1/devices -Headers $operatorHeaders -ContentType application/json -Body $deviceBody | Out-Null
Invoke-RestMethod -Method Put -Uri "http://127.0.0.1:6081/api/v1/devices/$smokeID/capabilities/navigate" -Headers $operatorHeaders -ContentType application/json -Body '{"configuration":{}}' | Out-Null
Invoke-RestMethod -Method Put -Uri "http://127.0.0.1:6081/api/v1/devices/$smokeID/capabilities/receive_relocation_command" -Headers $operatorHeaders -ContentType application/json -Body '{"configuration":{}}' | Out-Null
Invoke-RestMethod -Method Post -Uri "http://127.0.0.1:6081/api/v1/devices/$smokeID/activate" -Headers $operatorHeaders | Out-Null
$credentialResponse = Invoke-RestMethod -Method Post -Uri "http://127.0.0.1:6081/api/v1/devices/$smokeID/credentials" -Headers $operatorHeaders -ContentType application/json -Body '{}'
$env:DEVICE_TOKEN = $credentialResponse.data.secret
$env:DEVICE_CREDENTIAL_ID = $credentialResponse.data.credential.credential_id
$env:OPERATOR_TOKEN = $env:DEV_PLATFORM_ADMIN_TOKEN
$env:SMOKE_DEVICE_ID = $smokeID
$env:SMOKE_TENANT_ID = 'alpha_logistics'
$env:SMOKE_NODE_TYPE = '5'
$env:SMOKE_LAT = '13.0067'
$env:SMOKE_LON = '80.2206'
$env:SMOKE_SEQUENCE = '1'
$env:SMOKE_WAIT_FOR_PROJECTION = 'true'
$env:SMOKE_WAIT_FOR_MATCH = 'true'
foreach ($name in @('SMOKE_BOOT_ID','SMOKE_BOOT_STARTED_AT')) {
  [Environment]::SetEnvironmentVariable($name, $null, 'Process')
}
Push-Location $backendDir
try { $smokeResultJson = go run ./cmd/smoke } finally { Pop-Location }
if ($LASTEXITCODE -ne 0) { throw "Smoke client failed" }
if (-not $smokeResultJson) { throw "Smoke client returned no result" }
$smokeResult = $smokeResultJson | ConvertFrom-Json
$smokeID = $smokeResult.id

$count = 0
$archiveDeadline = (Get-Date).AddSeconds(60)
while ((Get-Date) -lt $archiveDeadline -and [int]$count -lt 1) {
  $count = docker compose -f (Join-Path $deploymentDir "docker-compose.yml") exec -T postgres `
    psql -U polaris_user -d polaris_core -tAc "SELECT count(*) FROM telemetry_history WHERE device_id='$smokeID'"
  if ([int]$count -lt 1) { Start-Sleep -Milliseconds 250 }
}
if ([int]$count -lt 1) {
  $archiveGroup = docker compose -f (Join-Path $deploymentDir "docker-compose.yml") exec -T redpanda `
    rpk group describe polaris_archive_group -X brokers=localhost:29092
  throw "Telemetry was not archived in PostgreSQL within 60 seconds. Archive consumer state:`n$archiveGroup"
}

Write-Host "PASS: Simulator -> Gateway -> Kafka -> Engine -> Redis -> PostgreSQL -> Dashboard ($smokeID, $($smokeResult.end_to_end_latency_ms) ms)"
```

---

## algo_\geo\math.go

```
package geo

import "math"

const EarthRadiusKm = 6371.0

// Haversine calculates the exact great-circle distance between two points on Earth in kilometers.
func Haversine(lat1, lon1, lat2, lon2 float64) float64 {
	dLat := (lat2 - lat1) * (math.Pi / 180.0)
	dLon := (lon2 - lon1) * (math.Pi / 180.0)

	lat1Rad := lat1 * (math.Pi / 180.0)
	lat2Rad := lat2 * (math.Pi / 180.0)

	a := math.Sin(dLat/2)*math.Sin(dLat/2) +
		math.Sin(dLon/2)*math.Sin(dLon/2)*math.Cos(lat1Rad)*math.Cos(lat2Rad)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	return EarthRadiusKm * c
}
```

---

## algo_\logger\logger.go

```
package logger

import (
	"log/slog"
	"os"
)

// Init configures a global, structured JSON logger for enterprise observability.
func Init() {
	opts := &slog.HandlerOptions{
		Level: slog.LevelDebug, // Set to Info in production
	}
	
	// Create a JSON handler that outputs to standard out
	handler := slog.NewJSONHandler(os.Stdout, opts)
	logger := slog.New(handler)
	
	// Set this as the default logger for the entire Go application
	slog.SetDefault(logger)
	slog.Info("Structured JSON logging initialized")
}
```

---

## api\proto\v1\spatial.pb.go

```
// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.11
// 	protoc        v7.35.0
// source: api/proto/v1/spatial.proto

package v1

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
	unsafe "unsafe"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type NodeType int32

const (
	NodeType_NODE_TYPE_UNKNOWN       NodeType = 0
	NodeType_NODE_TYPE_BIKE          NodeType = 1
	NodeType_NODE_TYPE_AUTO          NodeType = 2
	NodeType_NODE_TYPE_SEDAN         NodeType = 3
	NodeType_NODE_TYPE_SUV           NodeType = 4
	NodeType_NODE_TYPE_DRONE         NodeType = 5
	NodeType_NODE_TYPE_ROBOT         NodeType = 6
	NodeType_NODE_TYPE_STATIC_SENSOR NodeType = 7
)

// Enum value maps for NodeType.
var (
	NodeType_name = map[int32]string{
		0: "NODE_TYPE_UNKNOWN",
		1: "NODE_TYPE_BIKE",
		2: "NODE_TYPE_AUTO",
		3: "NODE_TYPE_SEDAN",
		4: "NODE_TYPE_SUV",
		5: "NODE_TYPE_DRONE",
		6: "NODE_TYPE_ROBOT",
		7: "NODE_TYPE_STATIC_SENSOR",
	}
	NodeType_value = map[string]int32{
		"NODE_TYPE_UNKNOWN":       0,
		"NODE_TYPE_BIKE":          1,
		"NODE_TYPE_AUTO":          2,
		"NODE_TYPE_SEDAN":         3,
		"NODE_TYPE_SUV":           4,
		"NODE_TYPE_DRONE":         5,
		"NODE_TYPE_ROBOT":         6,
		"NODE_TYPE_STATIC_SENSOR": 7,
	}
)

func (x NodeType) Enum() *NodeType {
	p := new(NodeType)
	*p = x
	return p
}

func (x NodeType) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (NodeType) Descriptor() protoreflect.EnumDescriptor {
	return file_api_proto_v1_spatial_proto_enumTypes[0].Descriptor()
}

func (NodeType) Type() protoreflect.EnumType {
	return &file_api_proto_v1_spatial_proto_enumTypes[0]
}

func (x NodeType) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use NodeType.Descriptor instead.
func (NodeType) EnumDescriptor() ([]byte, []int) {
	return file_api_proto_v1_spatial_proto_rawDescGZIP(), []int{0}
}

type NodeStatus int32

const (
	NodeStatus_NODE_STATUS_UNKNOWN     NodeStatus = 0
	NodeStatus_NODE_STATUS_IDLE        NodeStatus = 1
	NodeStatus_NODE_STATUS_EN_ROUTE    NodeStatus = 2
	NodeStatus_NODE_STATUS_ACTIVE      NodeStatus = 3
	NodeStatus_NODE_STATUS_MAINTENANCE NodeStatus = 4
	NodeStatus_NODE_STATUS_OFFLINE     NodeStatus = 5
)

// Enum value maps for NodeStatus.
var (
	NodeStatus_name = map[int32]string{
		0: "NODE_STATUS_UNKNOWN",
		1: "NODE_STATUS_IDLE",
		2: "NODE_STATUS_EN_ROUTE",
		3: "NODE_STATUS_ACTIVE",
		4: "NODE_STATUS_MAINTENANCE",
		5: "NODE_STATUS_OFFLINE",
	}
	NodeStatus_value = map[string]int32{
		"NODE_STATUS_UNKNOWN":     0,
		"NODE_STATUS_IDLE":        1,
		"NODE_STATUS_EN_ROUTE":    2,
		"NODE_STATUS_ACTIVE":      3,
		"NODE_STATUS_MAINTENANCE": 4,
		"NODE_STATUS_OFFLINE":     5,
	}
)

func (x NodeStatus) Enum() *NodeStatus {
	p := new(NodeStatus)
	*p = x
	return p
}

func (x NodeStatus) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (NodeStatus) Descriptor() protoreflect.EnumDescriptor {
	return file_api_proto_v1_spatial_proto_enumTypes[1].Descriptor()
}

func (NodeStatus) Type() protoreflect.EnumType {
	return &file_api_proto_v1_spatial_proto_enumTypes[1]
}

func (x NodeStatus) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use NodeStatus.Descriptor instead.
func (NodeStatus) EnumDescriptor() ([]byte, []int) {
	return file_api_proto_v1_spatial_proto_rawDescGZIP(), []int{1}
}

type SpatialObject struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Id            string                 `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	TenantId      string                 `protobuf:"bytes,2,opt,name=tenant_id,json=tenantId,proto3" json:"tenant_id,omitempty"`
	Type          NodeType               `protobuf:"varint,3,opt,name=type,proto3,enum=polaris.spatial.v1.NodeType" json:"type,omitempty"`
	Status        NodeStatus             `protobuf:"varint,4,opt,name=status,proto3,enum=polaris.spatial.v1.NodeStatus" json:"status,omitempty"`
	Lat           float64                `protobuf:"fixed64,5,opt,name=lat,proto3" json:"lat,omitempty"`
	Lon           float64                `protobuf:"fixed64,6,opt,name=lon,proto3" json:"lon,omitempty"`
	VelocityMps   float64                `protobuf:"fixed64,7,opt,name=velocity_mps,json=velocityMps,proto3" json:"velocity_mps,omitempty"`
	HeadingDeg    float64                `protobuf:"fixed64,8,opt,name=heading_deg,json=headingDeg,proto3" json:"heading_deg,omitempty"`
	EnergyPercent int32                  `protobuf:"varint,9,opt,name=energy_percent,json=energyPercent,proto3" json:"energy_percent,omitempty"` // Generates as .EnergyPercent in Go
	// Deprecated: Marked as deprecated in api/proto/v1/spatial.proto.
	Timestamp int64 `protobuf:"varint,10,opt,name=timestamp,proto3" json:"timestamp,omitempty"` // Legacy observed-at alias (Unix milliseconds)
	// Device-owned ordering identity. The gateway validates these fields and
	// wraps this frame in the canonical platform envelope.
	DeviceBootId   string `protobuf:"bytes,11,opt,name=device_boot_id,json=deviceBootId,proto3" json:"device_boot_id,omitempty"`
	SequenceNumber uint64 `protobuf:"varint,12,opt,name=sequence_number,json=sequenceNumber,proto3" json:"sequence_number,omitempty"`
	BootStartedAt  int64  `protobuf:"varint,13,opt,name=boot_started_at,json=bootStartedAt,proto3" json:"boot_started_at,omitempty"` // Unix milliseconds
	ObservedAt     int64  `protobuf:"varint,14,opt,name=observed_at,json=observedAt,proto3" json:"observed_at,omitempty"`            // Unix milliseconds
	SchemaVersion  uint32 `protobuf:"varint,15,opt,name=schema_version,json=schemaVersion,proto3" json:"schema_version,omitempty"`
	unknownFields  protoimpl.UnknownFields
	sizeCache      protoimpl.SizeCache
}

func (x *SpatialObject) Reset() {
	*x = SpatialObject{}
	mi := &file_api_proto_v1_spatial_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *SpatialObject) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SpatialObject) ProtoMessage() {}

func (x *SpatialObject) ProtoReflect() protoreflect.Message {
	mi := &file_api_proto_v1_spatial_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SpatialObject.ProtoReflect.Descriptor instead.
func (*SpatialObject) Descriptor() ([]byte, []int) {
	return file_api_proto_v1_spatial_proto_rawDescGZIP(), []int{0}
}

func (x *SpatialObject) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *SpatialObject) GetTenantId() string {
	if x != nil {
		return x.TenantId
	}
	return ""
}

func (x *SpatialObject) GetType() NodeType {
	if x != nil {
		return x.Type
	}
	return NodeType_NODE_TYPE_UNKNOWN
}

func (x *SpatialObject) GetStatus() NodeStatus {
	if x != nil {
		return x.Status
	}
	return NodeStatus_NODE_STATUS_UNKNOWN
}

func (x *SpatialObject) GetLat() float64 {
	if x != nil {
		return x.Lat
	}
	return 0
}

func (x *SpatialObject) GetLon() float64 {
	if x != nil {
		return x.Lon
	}
	return 0
}

func (x *SpatialObject) GetVelocityMps() float64 {
	if x != nil {
		return x.VelocityMps
	}
	return 0
}

func (x *SpatialObject) GetHeadingDeg() float64 {
	if x != nil {
		return x.HeadingDeg
	}
	return 0
}

func (x *SpatialObject) GetEnergyPercent() int32 {
	if x != nil {
		return x.EnergyPercent
	}
	return 0
}

// Deprecated: Marked as deprecated in api/proto/v1/spatial.proto.
func (x *SpatialObject) GetTimestamp() int64 {
	if x != nil {
		return x.Timestamp
	}
	return 0
}

func (x *SpatialObject) GetDeviceBootId() string {
	if x != nil {
		return x.DeviceBootId
	}
	return ""
}

func (x *SpatialObject) GetSequenceNumber() uint64 {
	if x != nil {
		return x.SequenceNumber
	}
	return 0
}

func (x *SpatialObject) GetBootStartedAt() int64 {
	if x != nil {
		return x.BootStartedAt
	}
	return 0
}

func (x *SpatialObject) GetObservedAt() int64 {
	if x != nil {
		return x.ObservedAt
	}
	return 0
}

func (x *SpatialObject) GetSchemaVersion() uint32 {
	if x != nil {
		return x.SchemaVersion
	}
	return 0
}

var File_api_proto_v1_spatial_proto protoreflect.FileDescriptor

const file_api_proto_v1_spatial_proto_rawDesc = "" +
	"\n" +
	"\x1aapi/proto/v1/spatial.proto\x12\x12polaris.spatial.v1\"\x96\x04\n" +
	"\rSpatialObject\x12\x0e\n" +
	"\x02id\x18\x01 \x01(\tR\x02id\x12\x1b\n" +
	"\ttenant_id\x18\x02 \x01(\tR\btenantId\x120\n" +
	"\x04type\x18\x03 \x01(\x0e2\x1c.polaris.spatial.v1.NodeTypeR\x04type\x126\n" +
	"\x06status\x18\x04 \x01(\x0e2\x1e.polaris.spatial.v1.NodeStatusR\x06status\x12\x10\n" +
	"\x03lat\x18\x05 \x01(\x01R\x03lat\x12\x10\n" +
	"\x03lon\x18\x06 \x01(\x01R\x03lon\x12!\n" +
	"\fvelocity_mps\x18\a \x01(\x01R\vvelocityMps\x12\x1f\n" +
	"\vheading_deg\x18\b \x01(\x01R\n" +
	"headingDeg\x12%\n" +
	"\x0eenergy_percent\x18\t \x01(\x05R\renergyPercent\x12 \n" +
	"\ttimestamp\x18\n" +
	" \x01(\x03B\x02\x18\x01R\ttimestamp\x12$\n" +
	"\x0edevice_boot_id\x18\v \x01(\tR\fdeviceBootId\x12'\n" +
	"\x0fsequence_number\x18\f \x01(\x04R\x0esequenceNumber\x12&\n" +
	"\x0fboot_started_at\x18\r \x01(\x03R\rbootStartedAt\x12\x1f\n" +
	"\vobserved_at\x18\x0e \x01(\x03R\n" +
	"observedAt\x12%\n" +
	"\x0eschema_version\x18\x0f \x01(\rR\rschemaVersion*\xb8\x01\n" +
	"\bNodeType\x12\x15\n" +
	"\x11NODE_TYPE_UNKNOWN\x10\x00\x12\x12\n" +
	"\x0eNODE_TYPE_BIKE\x10\x01\x12\x12\n" +
	"\x0eNODE_TYPE_AUTO\x10\x02\x12\x13\n" +
	"\x0fNODE_TYPE_SEDAN\x10\x03\x12\x11\n" +
	"\rNODE_TYPE_SUV\x10\x04\x12\x13\n" +
	"\x0fNODE_TYPE_DRONE\x10\x05\x12\x13\n" +
	"\x0fNODE_TYPE_ROBOT\x10\x06\x12\x1b\n" +
	"\x17NODE_TYPE_STATIC_SENSOR\x10\a*\xa3\x01\n" +
	"\n" +
	"NodeStatus\x12\x17\n" +
	"\x13NODE_STATUS_UNKNOWN\x10\x00\x12\x14\n" +
	"\x10NODE_STATUS_IDLE\x10\x01\x12\x18\n" +
	"\x14NODE_STATUS_EN_ROUTE\x10\x02\x12\x16\n" +
	"\x12NODE_STATUS_ACTIVE\x10\x03\x12\x1b\n" +
	"\x17NODE_STATUS_MAINTENANCE\x10\x04\x12\x17\n" +
	"\x13NODE_STATUS_OFFLINE\x10\x05B8Z6github.com/Akashpg-M/polaris/v3.0/backend/api/proto/v1b\x06proto3"

var (
	file_api_proto_v1_spatial_proto_rawDescOnce sync.Once
	file_api_proto_v1_spatial_proto_rawDescData []byte
)

func file_api_proto_v1_spatial_proto_rawDescGZIP() []byte {
	file_api_proto_v1_spatial_proto_rawDescOnce.Do(func() {
		file_api_proto_v1_spatial_proto_rawDescData = protoimpl.X.CompressGZIP(unsafe.Slice(unsafe.StringData(file_api_proto_v1_spatial_proto_rawDesc), len(file_api_proto_v1_spatial_proto_rawDesc)))
	})
	return file_api_proto_v1_spatial_proto_rawDescData
}

var file_api_proto_v1_spatial_proto_enumTypes = make([]protoimpl.EnumInfo, 2)
var file_api_proto_v1_spatial_proto_msgTypes = make([]protoimpl.MessageInfo, 1)
var file_api_proto_v1_spatial_proto_goTypes = []any{
	(NodeType)(0),         // 0: polaris.spatial.v1.NodeType
	(NodeStatus)(0),       // 1: polaris.spatial.v1.NodeStatus
	(*SpatialObject)(nil), // 2: polaris.spatial.v1.SpatialObject
}
var file_api_proto_v1_spatial_proto_depIdxs = []int32{
	0, // 0: polaris.spatial.v1.SpatialObject.type:type_name -> polaris.spatial.v1.NodeType
	1, // 1: polaris.spatial.v1.SpatialObject.status:type_name -> polaris.spatial.v1.NodeStatus
	2, // [2:2] is the sub-list for method output_type
	2, // [2:2] is the sub-list for method input_type
	2, // [2:2] is the sub-list for extension type_name
	2, // [2:2] is the sub-list for extension extendee
	0, // [0:2] is the sub-list for field type_name
}

func init() { file_api_proto_v1_spatial_proto_init() }
func file_api_proto_v1_spatial_proto_init() {
	if File_api_proto_v1_spatial_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: unsafe.Slice(unsafe.StringData(file_api_proto_v1_spatial_proto_rawDesc), len(file_api_proto_v1_spatial_proto_rawDesc)),
			NumEnums:      2,
			NumMessages:   1,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_api_proto_v1_spatial_proto_goTypes,
		DependencyIndexes: file_api_proto_v1_spatial_proto_depIdxs,
		EnumInfos:         file_api_proto_v1_spatial_proto_enumTypes,
		MessageInfos:      file_api_proto_v1_spatial_proto_msgTypes,
	}.Build()
	File_api_proto_v1_spatial_proto = out.File
	file_api_proto_v1_spatial_proto_goTypes = nil
	file_api_proto_v1_spatial_proto_depIdxs = nil
}
```

---

## api\proto\v1\spatial.proto

```
syntax = "proto3";
package polaris.spatial.v1;
option go_package = "github.com/Akashpg-M/polaris/v3.0/backend/api/proto/v1";

enum NodeType {
  NODE_TYPE_UNKNOWN = 0;
  NODE_TYPE_BIKE = 1;
  NODE_TYPE_AUTO = 2;
  NODE_TYPE_SEDAN = 3;
  NODE_TYPE_SUV = 4;
  NODE_TYPE_DRONE = 5;
  NODE_TYPE_ROBOT = 6;
  NODE_TYPE_STATIC_SENSOR = 7;
}

enum NodeStatus {
  NODE_STATUS_UNKNOWN = 0;
  NODE_STATUS_IDLE = 1;
  NODE_STATUS_EN_ROUTE = 2;
  NODE_STATUS_ACTIVE = 3;
  NODE_STATUS_MAINTENANCE = 4;
  NODE_STATUS_OFFLINE = 5;
}

message SpatialObject {
  string id = 1;
  string tenant_id = 2;
  NodeType type = 3;
  NodeStatus status = 4;
  double lat = 5;
  double lon = 6;
  double velocity_mps = 7; 
  double heading_deg = 8;  
  int32 energy_percent = 9; // Generates as .EnergyPercent in Go
  int64 timestamp = 10 [deprecated = true]; // Legacy observed-at alias (Unix milliseconds)

  // Device-owned ordering identity. The gateway validates these fields and
  // wraps this frame in the canonical platform envelope.
  string device_boot_id = 11;
  uint64 sequence_number = 12;
  int64 boot_started_at = 13; // Unix milliseconds
  int64 observed_at = 14;     // Unix milliseconds
  uint32 schema_version = 15;
}
```

---

## cmd\engine\main.go

```
package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"strconv"
	"syscall"
	"time"

	"github.com/Akashpg-M/polaris/backend/algo_/logger"
	"github.com/Akashpg-M/polaris/backend/internal/adapter/handler"
	"github.com/Akashpg-M/polaris/backend/internal/adapter/repository"
	"github.com/Akashpg-M/polaris/backend/internal/application/dispatch"
	"github.com/Akashpg-M/polaris/backend/internal/application/orchestration"
	"github.com/Akashpg-M/polaris/backend/internal/application/orchestrator"
	"github.com/Akashpg-M/polaris/backend/internal/application/outbox"
	"github.com/Akashpg-M/polaris/backend/internal/application/reconciliation"
	"github.com/Akashpg-M/polaris/backend/internal/application/spatial"
	"github.com/Akashpg-M/polaris/backend/internal/application/stream"
	"github.com/Akashpg-M/polaris/backend/internal/application/twin"
	"github.com/Akashpg-M/polaris/backend/internal/config"
	"github.com/Akashpg-M/polaris/backend/internal/core/extension"
	redisinfra "github.com/Akashpg-M/polaris/backend/internal/infra/redis"
	mobilitymodule "github.com/Akashpg-M/polaris/backend/internal/modules/mobility"
	"github.com/Akashpg-M/polaris/backend/internal/modules/mobility/matching"
	"github.com/Akashpg-M/polaris/backend/internal/modules/mobility/planning"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	// 1. Initialize Config & Logger
	cfg := config.Load()
	logger.Init()
	slog.Info("Booting Polaris v3.0 Spatial Engine...", "env", cfg.App.Env)

	engine := spatial.NewEngine()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// 2. Initialize Dependencies (Using the nested Config structs)
	// redisConsumer, _ := stream.NewRedisConsumer(cfg.Redis.URL, engine)
	// go redisConsumer.Start(ctx, "engine-node-1")

	kafkaBroker := os.Getenv("KAFKA_BROKER_URL")
	if kafkaBroker == "" {
		kafkaBroker = "localhost:9092"
	}

	redisClient, err := redisinfra.NewClient(cfg.Redis.URL)
	if err != nil {
		panic("Cannot start engine without Redis: " + err.Error())
	}
	registryStore, err := repository.NewRegistryStore(cfg.DB.URL)
	if err != nil {
		panic("Cannot start registry: " + err.Error())
	}
	defer registryStore.Close()
	if err = registryStore.BootstrapPlatformAdmin(ctx, os.Getenv("DEV_PLATFORM_ADMIN_TOKEN")); err != nil {
		panic("Cannot bootstrap development operator: " + err.Error())
	}
	extensionRegistry := extension.NewRegistry()
	mobilityCfg, err := mobilitymodule.LoadConfig()
	if err != nil {
		panic("Invalid Mobility configuration: " + err.Error())
	}
	var mobilityModule *mobilitymodule.Module
	stateFanout := &stream.StateFanout{Primary: engine}
	if mobilityCfg.Enabled {
		mobilityModule = mobilitymodule.New(mobilityCfg, mobilityRebuildLoader(redisClient, registryStore))
		extensionRegistry.RegisterModule(mobilityModule)
		if mobilityCfg.SpatialEnabled {
			extensionRegistry.RegisterCandidateProvider(&matching.Provider{Spatial: mobilityModule.Spatial, Routing: mobilityModule, RawLimit: mobilityCfg.MaxRawCandidates, RoutedLimit: mobilityCfg.MaxRoutedCandidates, MaxRadius: mobilityCfg.MaxSearchRadiusMeters})
		}
		extensionRegistry.RegisterTaskPlanner(&planning.Planner{SpatialState: mobilityModule.Spatial.Get, Routing: mobilityModule, MaxPlanAge: 2 * time.Minute})
		if mobilityCfg.SpatialEnabled {
			stateFanout.Projections = append(stateFanout.Projections, &mobilitymodule.TelemetryProjector{Manager: mobilityModule.Spatial})
		}
	}
	extensionRegistry.RegisterTaskPlanner(extension.DefaultTaskPlanner{})
	if err = extensionRegistry.Start(ctx); err != nil {
		panic("Cannot start capability modules: " + err.Error())
	}
	var mobilityTraffic *mobilitymodule.TrafficConsumer
	if mobilityModule != nil && mobilityModule.Traffic() != nil {
		mobilityTraffic = mobilitymodule.NewTrafficConsumer(kafkaBroker, mobilityModule.Traffic(), mobilityCfg.TrafficRefreshInterval)
		go mobilityTraffic.Start(ctx)
	}
	kafkaConsumer := stream.NewKafkaConsumer(kafkaBroker, stateFanout, redisClient)
	go kafkaConsumer.Start(ctx, "engine-node-1")
	archiver, err := stream.NewKafkaPostgresArchiver(kafkaBroker, cfg.DB.URL)
	if err != nil {
		slog.Warn("Kafka/PostgreSQL Archiver offline", "error", err)
	} else {
		go archiver.Start(ctx)
	}
	outboxRelay := outbox.New(registryStore, kafkaBroker, envInt("OUTBOX_BATCH_SIZE", 100), envDuration("OUTBOX_POLL_INTERVAL", 500*time.Millisecond))
	go outboxRelay.Start(ctx)
	orchestrationMetrics := orchestration.NewMetrics()
	orchestrationService := orchestration.NewServiceWithRegistry(registryStore, redisClient, envInt("COMMAND_MAX_ATTEMPTS", 5), orchestrationMetrics, extensionRegistry)
	ownershipStore := repository.NewConnectionOwnershipStore(redisClient, envDuration("CONNECTION_LEASE_TTL", 30*time.Second))
	commandDispatcher := dispatch.New(kafkaBroker, redisClient, ownershipStore)
	go commandDispatcher.Start(ctx)
	commandReconciler := reconciliation.New(registryStore, orchestrationService, ownershipStore, envDuration("COMMAND_RECONCILE_INTERVAL", time.Second), envDuration("COMMAND_ACK_TIMEOUT", 5*time.Second))
	go commandReconciler.Start(ctx)
	connectivityDetector := twin.NewDetector(redisClient, kafkaBroker, envDuration("DEVICE_STALE_AFTER", 30*time.Second), envDuration("DEVICE_OFFLINE_AFTER", 90*time.Second), envDuration("OFFLINE_SCAN_INTERVAL", 10*time.Second))
	if mobilityModule != nil {
		connectivityDetector.SetTransitionHandler(func(tenant, device, status string) {
			_ = mobilityModule.Spatial.EvictInactive(tenant, device, "ACTIVE", status)
		})
	}
	go connectivityDetector.Start(ctx)

	predictiveStrategy, err := orchestrator.NewPredictiveZoneStrategy(cfg.DB.URL)
	if err != nil {
		slog.Warn("Predictive zone view offline", "error", err)
	} else {
		defer predictiveStrategy.Close()
	}

	// 3. Setup Router
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	router.Use(gin.Recovery(), cors.New(cors.Config{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{http.MethodGet, http.MethodPost, http.MethodPatch, http.MethodPut, http.MethodDelete, http.MethodOptions},
		AllowHeaders: []string{"Origin", "Content-Type", "Authorization", "X-Tenant-ID"},
	}))

	matchHandler := handler.NewMatchHandler(engine)

	api := router.Group("/api/v1")
	{
		registryAPI := handler.NewRegistryAPI(registryStore, redisClient, envDuration("DEVICE_STALE_AFTER", 30*time.Second), envDuration("DEVICE_OFFLINE_AFTER", 90*time.Second), envDuration("CONNECTION_TICKET_TTL", 30*time.Second))
		if mobilityModule != nil && mobilityCfg.SpatialEnabled {
			registryAPI.SetLifecycleHook(func(tenant, device, status string) {
				if device == "" {
					if status != "ACTIVE" {
						_ = mobilityModule.Spatial.RemoveTenant(tenant)
					}
				} else {
					_ = mobilityModule.Spatial.EvictInactive(tenant, device, status, "ONLINE")
				}
			})
		}
		registryAPI.Register(api)
		handler.NewOrchestrationAPI(registryStore, orchestrationService, orchestrationMetrics).Register(api, registryAPI)
		protected := api.Group("")
		protected.Use(registryAPI.Middleware("read"))
		protected.GET("/nodes/match", matchHandler.GetNearestNodes)
		if mobilityModule != nil {
			handler.NewMobilityAPI(mobilityModule.Spatial, mobilityModule, mobilityCfg.MaxRawCandidates).Register(protected)
		}
		protected.GET("/zones/predicted", func(c *gin.Context) {
			if predictiveStrategy != nil {
				c.JSON(200, gin.H{"status": "success", "data": predictiveStrategy.GetTargetZones(context.Background())})
			} else {
				c.JSON(200, gin.H{"status": "success", "data": []interface{}{}})
			}
		})
	}
	router.GET("/healthz", func(c *gin.Context) {
		if !kafkaConsumer.Healthy() || !commandDispatcher.Healthy() || (archiver != nil && !archiver.Healthy()) {
			c.JSON(http.StatusServiceUnavailable, gin.H{"status": "unhealthy"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"status": "live"})
	})
	router.GET("/readyz", func(c *gin.Context) {
		probeCtx, cancel := context.WithTimeout(c.Request.Context(), 2*time.Second)
		defer cancel()
		if err := kafkaConsumer.Ready(probeCtx); err != nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{"status": "not_ready", "dependency": "kafka_or_redis", "error": err.Error()})
			return
		}
		if archiver == nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{"status": "not_ready", "dependency": "postgres_archiver"})
			return
		}
		if err := registryStore.DB.PingContext(probeCtx); err != nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{"status": "not_ready", "dependency": "registry", "error": err.Error()})
			return
		}
		if err := archiver.Ready(probeCtx); err != nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{"status": "not_ready", "dependency": "postgres", "error": err.Error()})
			return
		}
		dbStats := registryStore.DB.Stats()
		c.JSON(http.StatusOK, gin.H{"status": "ready", "core": "READY", "modules": extensionRegistry.Status(probeCtx), "runtime": gin.H{"goroutines": runtime.NumGoroutine(), "db_open_connections": dbStats.OpenConnections, "db_in_use": dbStats.InUse, "db_idle": dbStats.Idle, "db_wait_count": dbStats.WaitCount}})
	})

	// 4. Start Server with Graceful Shutdown
	port := ":" + cfg.Server.EnginePort
	srv := &http.Server{Addr: port, Handler: router}

	go func() {
		slog.Info("Engine REST API active", "port", port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("Server failed", "error", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	slog.Warn("Shutdown signal received...")
	ctxShutdown, cancelShutdown := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancelShutdown()

	srv.Shutdown(ctxShutdown)
	cancel() // Stops background context (workers)
	if err := extensionRegistry.Close(ctxShutdown); err != nil {
		slog.Error("Capability module shutdown failed", "error", err)
	}
	if err := kafkaConsumer.Wait(ctxShutdown); err != nil {
		slog.Error("Spatial consumer shutdown timed out", "error", err)
	}
	if archiver != nil {
		if err := archiver.Wait(ctxShutdown); err != nil {
			slog.Error("Archive consumer shutdown timed out", "error", err)
		}
	}
	if err := commandDispatcher.Wait(ctxShutdown); err != nil {
		slog.Error("Command dispatcher shutdown timed out", "error", err)
	}
	if mobilityTraffic != nil {
		if err := mobilityTraffic.Wait(ctxShutdown); err != nil {
			slog.Error("Mobility traffic consumer shutdown timed out", "error", err)
		}
	}
	redisClient.Close()
	slog.Info("Engine safely terminated.")
}

func envDuration(key string, fallback time.Duration) time.Duration {
	if raw := os.Getenv(key); raw != "" {
		if v, err := time.ParseDuration(raw); err == nil {
			return v
		}
	}
	return fallback
}
func envInt(key string, fallback int) int {
	if raw := os.Getenv(key); raw != "" {
		if v, err := strconv.Atoi(raw); err == nil {
			return v
		}
	}
	return fallback
}
```

---

## cmd\engine\mobility_rebuild.go

```
package main

import (
	"context"
	"encoding/json"
	"strconv"
	"strings"
	"time"

	"github.com/Akashpg-M/polaris/backend/internal/adapter/repository"
	twincore "github.com/Akashpg-M/polaris/backend/internal/core/twin"
	"github.com/Akashpg-M/polaris/backend/internal/modules/mobility/model"
	"github.com/redis/go-redis/v9"
)

type spatialPayloadV1 struct {
	Latitude        float64               `json:"latitude"`
	Longitude       float64               `json:"longitude"`
	HeadingDegrees  *float64              `json:"heading_degrees"`
	SpeedMPS        *float64              `json:"speed_mps"`
	MobilityProfile model.MobilityProfile `json:"mobility_profile"`
}

func mobilityRebuildLoader(client *redis.Client, store *repository.RegistryStore) func(context.Context) ([]model.SpatialState, error) {
	return func(ctx context.Context) ([]model.SpatialState, error) {
		var cursor uint64
		out := []model.SpatialState{}
		for {
			keys, next, err := client.Scan(ctx, cursor, "polaris:twin:*", 200).Result()
			if err != nil {
				return nil, err
			}
			for _, key := range keys {
				if strings.HasSuffix(key, ":retired_boots") {
					continue
				}
				identity := strings.TrimPrefix(key, "polaris:twin:")
				parts := strings.SplitN(identity, ":", 2)
				if len(parts) != 2 {
					continue
				}
				tenantID, deviceID := parts[0], parts[1]
				values, err := client.HGetAll(ctx, key).Result()
				if err != nil || values["connectivity_status"] != "ONLINE" {
					continue
				}
				device, err := store.GetDevice(ctx, tenantID, deviceID)
				if err != nil || device.LifecycleStatus != "ACTIVE" {
					continue
				}
				tenant, err := store.GetTenant(ctx, tenantID)
				if err != nil || tenant.Status != "ACTIVE" {
					continue
				}
				var envelope twincore.ComponentEnvelope
				if json.Unmarshal([]byte(values["component:spatial/v1"]), &envelope) != nil {
					continue
				}
				var payload spatialPayloadV1
				if json.Unmarshal(envelope.Payload, &payload) != nil {
					continue
				}
				bootStarted, _ := strconv.ParseInt(values["boot_started_at"], 10, 64)
				out = append(out, model.SpatialState{TenantID: tenantID, DeviceID: deviceID, ReportedPosition: model.Position{Latitude: payload.Latitude, Longitude: payload.Longitude}, HeadingDegrees: payload.HeadingDegrees, SpeedMPS: payload.SpeedMPS, MobilityProfile: payload.MobilityProfile, ObservedAt: envelope.ObservedAt, BootID: envelope.BootID, BootStartedAt: time.UnixMilli(bootStarted).UTC(), SequenceNumber: envelope.SequenceNumber})
			}
			cursor = next
			if cursor == 0 {
				break
			}
		}
		return out, nil
	}
}
```

---

## cmd\gateway\main.go

```
package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"

	"github.com/Akashpg-M/polaris/backend/algo_/logger"
	"github.com/Akashpg-M/polaris/backend/internal/adapter/handler"
	"github.com/Akashpg-M/polaris/backend/internal/adapter/repository"
	"github.com/Akashpg-M/polaris/backend/internal/application/orchestration"
	"github.com/Akashpg-M/polaris/backend/internal/config"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

func main() {
	cfg := config.Load()
	logger.Init()
	slog.Info("Booting Polaris v3.0 API Gateway...")
	kafkaBroker := getEnvFallback("KAFKA_BROKER_URL", "localhost:9092")
	kafkaPublisher := repository.NewKafkaEventPublisher(kafkaBroker)
	defer kafkaPublisher.Close()
	redisOptions, err := redis.ParseURL(cfg.Redis.URL)
	if err != nil {
		panic("invalid Redis URL: " + err.Error())
	}
	healthRedis := redis.NewClient(redisOptions)
	defer healthRedis.Close()
	registryStore, err := repository.NewRegistryStore(cfg.DB.URL)
	if err != nil {
		panic("Cannot start authenticated gateway without registry: " + err.Error())
	}
	defer registryStore.Close()

	dashboardRegistry := handler.NewDashboardRegistry()
	go startDashboardSubscriber(cfg.Redis.URL, dashboardRegistry)
	orchestrationMetrics := orchestration.NewMetrics()
	connectionManager := handler.NewDeviceConnectionManager(getEnvFallback("GATEWAY_ID", "gateway-1"), getDurationFallback("CONNECTION_LEASE_TTL", 30*time.Second), healthRedis, registryStore, orchestrationMetrics)
	subscriberCtx, stopSubscriber := context.WithCancel(context.Background())
	defer stopSubscriber()
	go connectionManager.StartSubscriber(subscriberCtx)

	gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Tenant-ID")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusOK)
			return
		}
		c.Next()
	})

	// The gateway owns connection state only. Accepted telemetry is
	// synchronously appended to Kafka, the durable ingress boundary.
	ingestionHandler := handler.NewIngestionHandler(kafkaPublisher, registryStore, connectionManager)
	dashboardHandler := handler.NewDashboardHandler(dashboardRegistry, registryStore)
	router.GET("/ws/telemetry", ingestionHandler.HandleIoTConnection)
	router.GET("/ws/dashboard", dashboardHandler.HandleWebConnection)
	api := router.Group("/api/v1")
	api.GET("/metrics/connections", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"active_uplinks": ingestionHandler.GetActiveConnectionsCount()})
	})
	api.GET("/metrics/orchestration", func(c *gin.Context) { c.JSON(http.StatusOK, orchestrationMetrics.Snapshot()) })
	router.GET("/healthz", func(c *gin.Context) { c.JSON(http.StatusOK, gin.H{"status": "live"}) })
	router.GET("/readyz", func(c *gin.Context) {
		probeCtx, cancel := context.WithTimeout(c.Request.Context(), 2*time.Second)
		defer cancel()
		if err := kafkaPublisher.Ready(probeCtx, kafkaBroker); err != nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{"status": "not_ready", "dependency": "kafka", "error": err.Error()})
			return
		}
		if err := healthRedis.Ping(probeCtx).Err(); err != nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{"status": "not_ready", "dependency": "redis", "error": err.Error()})
			return
		}
		if err := registryStore.DB.PingContext(probeCtx); err != nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{"status": "not_ready", "dependency": "registry", "error": err.Error()})
			return
		}
		dbStats := registryStore.DB.Stats()
		c.JSON(http.StatusOK, gin.H{"status": "ready", "runtime": gin.H{"goroutines": runtime.NumGoroutine(), "db_open_connections": dbStats.OpenConnections, "db_in_use": dbStats.InUse, "db_idle": dbStats.Idle, "db_wait_count": dbStats.WaitCount}})
	})

	port := ":" + cfg.Server.GatewayPort
	srv := &http.Server{Addr: port, Handler: router}
	go func() {
		slog.Info("Gateway active", "port", port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("Server failed", "error", err)
		}
	}()
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	slog.Warn("Shutdown signal received. Draining WebSockets...")
	stopSubscriber()
	ctxShutdown, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctxShutdown); err != nil {
		slog.Error("Gateway forced to shutdown", "error", err)
	}
	slog.Info("Gateway safely terminated.")
}

func startDashboardSubscriber(redisURL string, dashboardRegistry *handler.DashboardRegistry) {
	opts, _ := redis.ParseURL(redisURL)
	client := redis.NewClient(opts)
	defer client.Close()
	pubsub := client.Subscribe(context.Background(), "spatial:updates")
	defer pubsub.Close()
	for msg := range pubsub.Channel() {
		dashboardRegistry.BroadcastToUIs(msg.Payload)
	}
}

func getEnvFallback(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

func getDurationFallback(key string, fallback time.Duration) time.Duration {
	if value := os.Getenv(key); value != "" {
		if parsed, err := time.ParseDuration(value); err == nil {
			return parsed
		}
	}
	return fallback
}
```

---

## cmd\identitycheck\main.go

```
package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	pb "github.com/Akashpg-M/polaris/backend/api/proto/v1"
	"github.com/gorilla/websocket"
	"google.golang.org/protobuf/proto"
)

func main() {
	mode := env("IDENTITY_CHECK_MODE", "basic")
	token := mustEnv("DEVICE_TOKEN")
	device := mustEnv("SMOKE_DEVICE_ID")
	switch mode {
	case "basic":
		expectRejected("pol_dev_invalid.invalid")
		connect(token).Close()
		spoof(token, device)
	case "rejected":
		expectRejected(token)
	case "send":
		c := connect(token)
		send(c, device, "alpha_logistics", 1)
		c.Close()
	case "ticket":
		ticketCheck(device)
	case "revoke-session":
		c := connect(token)
		revoke()
		send(c, device, "alpha_logistics", 2)
		c.SetReadDeadline(time.Now().Add(3 * time.Second))
		if _, _, err := c.ReadMessage(); err == nil {
			panic("revoked active session remained open")
		}
	default:
		panic("unknown mode")
	}
	fmt.Println("PASS: " + mode)
}
func connect(token string) *websocket.Conn {
	h := http.Header{}
	h.Set("Authorization", "Bearer "+token)
	c, r, err := websocket.DefaultDialer.Dial(env("GATEWAY_URL", "ws://127.0.0.1:6080")+"/ws/telemetry", h)
	if err != nil {
		if r != nil {
			panic(fmt.Sprintf("WebSocket rejected: %d", r.StatusCode))
		}
		panic(err)
	}
	return c
}
func expectRejected(token string) {
	h := http.Header{}
	h.Set("Authorization", "Bearer "+token)
	c, r, err := websocket.DefaultDialer.Dial(env("GATEWAY_URL", "ws://127.0.0.1:6080")+"/ws/telemetry", h)
	if c != nil {
		c.Close()
		panic("invalid/revoked credential connected")
	}
	if err == nil || r == nil || r.StatusCode != http.StatusUnauthorized {
		panic("expected HTTP 401 before WebSocket upgrade")
	}
}
func spoof(token, device string) {
	c := connect(token)
	defer c.Close()
	send(c, "spoofed-device", "other_tenant", 1)
	c.SetReadDeadline(time.Now().Add(3 * time.Second))
	if _, _, err := c.ReadMessage(); err == nil {
		panic("spoofed identity was not disconnected")
	}
}
func send(c *websocket.Conn, device, tenant string, seq uint64) {
	now := time.Now().UTC()
	if raw := os.Getenv("TELEMETRY_SEQUENCE"); raw != "" {
		_, _ = fmt.Sscan(raw, &seq)
	}
	boot := env("DEVICE_BOOT_ID", "identity-check-boot")
	p := &pb.SpatialObject{Id: device, TenantId: tenant, Type: pb.NodeType_NODE_TYPE_DRONE, Status: pb.NodeStatus_NODE_STATUS_ACTIVE, Lat: 13.0067, Lon: 80.2206, VelocityMps: 10, EnergyPercent: 90, DeviceBootId: boot, SequenceNumber: seq, BootStartedAt: now.Add(-time.Minute).UnixMilli(), ObservedAt: now.UnixMilli(), SchemaVersion: 1}
	b, err := proto.Marshal(p)
	if err != nil {
		panic(err)
	}
	if err = c.WriteMessage(websocket.BinaryMessage, b); err != nil {
		panic(err)
	}
}
func revoke() {
	body := bytes.NewBufferString(`{}`)
	url := env("ENGINE_URL", "http://127.0.0.1:6081") + "/api/v1/devices/" + mustEnv("SMOKE_DEVICE_ID") + "/credentials/" + mustEnv("DEVICE_CREDENTIAL_ID") + "/revoke"
	req, _ := http.NewRequest(http.MethodPost, url, body)
	req.Header.Set("Authorization", "Bearer "+mustEnv("OPERATOR_TOKEN"))
	req.Header.Set("X-Tenant-ID", "alpha_logistics")
	req.Header.Set("Content-Type", "application/json")
	r, err := http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}
	defer r.Body.Close()
	if r.StatusCode != 200 {
		var v interface{}
		_ = json.NewDecoder(r.Body).Decode(&v)
		panic(fmt.Sprintf("revoke failed %d %#v", r.StatusCode, v))
	}
}
func ticketCheck(device string) {
	url := env("ENGINE_URL", "http://127.0.0.1:6081") + "/api/v1/devices/" + device + "/connection-ticket"
	req, _ := http.NewRequest(http.MethodPost, url, bytes.NewBufferString(`{}`))
	req.Header.Set("Authorization", "Bearer "+mustEnv("OPERATOR_TOKEN"))
	req.Header.Set("X-Tenant-ID", "alpha_logistics")
	req.Header.Set("Content-Type", "application/json")
	r, err := http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}
	var body struct {
		Data struct {
			Ticket string `json:"ticket"`
		} `json:"data"`
	}
	_ = json.NewDecoder(r.Body).Decode(&body)
	r.Body.Close()
	if r.StatusCode != 201 || body.Data.Ticket == "" {
		panic("ticket issue failed")
	}
	target := env("GATEWAY_URL", "ws://127.0.0.1:6080") + "/ws/telemetry?ticket=" + body.Data.Ticket
	c, _, err := websocket.DefaultDialer.Dial(target, nil)
	if err != nil {
		panic(err)
	}
	c.Close()
	c, response, err := websocket.DefaultDialer.Dial(target, nil)
	if c != nil {
		c.Close()
		panic("one-time ticket was reused")
	}
	if err == nil || response == nil || response.StatusCode != 401 {
		panic("consumed ticket did not return 401")
	}
}
func env(k, d string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return d
}
func mustEnv(k string) string {
	v := os.Getenv(k)
	if v == "" {
		panic(k + " required")
	}
	return v
}
```

---

## cmd\loadtest\main.go

```
package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os"
	"os/signal"
	"sync"
	"sync/atomic"
	"syscall"
	"time"

	pb "github.com/Akashpg-M/polaris/backend/api/proto/v1"
	"github.com/gorilla/websocket"
	"google.golang.org/protobuf/proto"
)

// Atomic counters for thread-safe, high-speed metrics tracking
var (
	activeConnections int64
	messagesSent      int64
	connectionErrors  int64
)

// Payload matches the Polaris domain model
type Payload struct {
	TenantID string  `json:"tenant_id"`
	NodeID   string  `json:"node_id"`
	Class    uint16  `json:"asset_class"`
	Lat      float64 `json:"lat"`
	Lon      float64 `json:"lon"`
	Status   string  `json:"status"`
	Battery  int     `json:"battery"`
}

func main() {
	// 1. Configurable Flags
	targetNodes := flag.Int("nodes", 1000, "Number of concurrent drones to simulate")
	serverURL := flag.String("url", "ws://127.0.0.1:6080/ws/telemetry", "Gateway WebSocket URL")
	rampRate := flag.Int("ramp", 100, "How many new connections to open per second")
	flag.Parse()

	log.Printf("🚀 Initiating Polaris Stress Test...")
	log.Printf("Targeting %d concurrent drones on %s", *targetNodes, *serverURL)

	var wg sync.WaitGroup
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// 2. Start the live metrics dashboard in the background
	go printMetricsDashboard(ctx)

	// 3. Smooth Connection Ramping
	// We calculate the delay between each dial to hit the requested ramp rate
	delayStr := fmt.Sprintf("%dms", 1000/(*rampRate))
	delay, _ := time.ParseDuration(delayStr)
	ticker := time.NewTicker(delay)
	defer ticker.Stop()

	// 4. Listen for Ctrl+C to stop the test gracefully
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// Spawn the drones
SpawnLoop:
	for i := 1; i <= *targetNodes; i++ {
		select {
		case <-quit:
			log.Println("\n⚠️ Aborting launch sequence early...")
			break SpawnLoop
		case <-ticker.C:
			wg.Add(1)
			go simulateDrone(ctx, i, *serverURL, &wg)
		}
	}

	log.Println("\n✅ All requested drones deployed. Press Ctrl+C to terminate test.")
	<-quit // Wait here until user kills the script

	cancel()  // Tell all drones to shut down
	wg.Wait() // Wait for them to actually close
	fmt.Println("\nStress test concluded cleanly.")
}

func simulateDrone(ctx context.Context, id int, wsURL string, wg *sync.WaitGroup) {
	defer wg.Done()

	conn, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	if err != nil {
		atomic.AddInt64(&connectionErrors, 1)
		return
	}
	defer conn.Close()

	atomic.AddInt64(&activeConnections, 1)
	defer atomic.AddInt64(&activeConnections, -1)

	nodeID := fmt.Sprintf("STRESS-DRONE-%d", id)
	lat := 13.04 + (rand.Float64() * 0.1)
	lon := 80.24 + (rand.Float64() * 0.1)

	go func() {
		for {
			if _, _, err := conn.ReadMessage(); err != nil {
				return
			}
		}
	}()

	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()
	bootStartedAt := time.Now().UTC().UnixMilli()
	var sequence uint64

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			sequence++
			lat += (rand.Float64() - 0.5) * 0.001
			lon += (rand.Float64() - 0.5) * 0.001

			// 1. Construct the Protobuf payload
			payload := &pb.SpatialObject{
				TenantId:       "alpha_logistics",
				Id:             nodeID,
				Type:           pb.NodeType_NODE_TYPE_DRONE,
				Lat:            lat,
				Lon:            lon,
				Status:         pb.NodeStatus_NODE_STATUS_ACTIVE,
				EnergyPercent:  int32(rand.Intn(100)),
				DeviceBootId:   "load-boot-" + nodeID,
				SequenceNumber: sequence,
				BootStartedAt:  bootStartedAt,
				ObservedAt:     time.Now().UTC().UnixMilli(),
				SchemaVersion:  1,
			}

			// 2. Marshal to raw bytes using proto
			msgBytes, err := proto.Marshal(payload)
			if err != nil {
				atomic.AddInt64(&connectionErrors, 1)
				return
			}

			// 3. Send as a BinaryMessage
			if err := conn.WriteMessage(websocket.BinaryMessage, msgBytes); err != nil {
				atomic.AddInt64(&connectionErrors, 1)
				return
			}

			atomic.AddInt64(&messagesSent, 1)
		}
	}
}

// printMetricsDashboard clears the console line and prints live stats
func printMetricsDashboard(ctx context.Context) {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()
	var lastMessagesSent int64

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			currentMessages := atomic.LoadInt64(&messagesSent)
			throughput := currentMessages - lastMessagesSent
			lastMessagesSent = currentMessages

			// \033[2K clears the current terminal line, \r returns cursor to the start
			fmt.Printf("\033[2K\r📡 ACTIVE UPLINKS: %d | ⚡ THROUGHPUT: %d msgs/sec | ❌ ERRORS: %d | 📦 TOTAL SENT: %d",
				atomic.LoadInt64(&activeConnections),
				throughput,
				atomic.LoadInt64(&connectionErrors),
				currentMessages,
			)
		}
	}
}
```

---

## cmd\mobility-benchmark\main.go

```
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"math/rand"
	"sort"
	"time"

	"github.com/Akashpg-M/polaris/backend/internal/modules/mobility/model"
	"github.com/Akashpg-M/polaris/backend/internal/modules/mobility/spatial"
)

type result struct {
	Index               string  `json:"index"`
	Devices             int     `json:"devices"`
	IndexedDevices      int     `json:"indexed_devices"`
	Queries             int     `json:"queries"`
	P50US               int64   `json:"p50_us"`
	P95US               int64   `json:"p95_us"`
	P99US               int64   `json:"p99_us"`
	OperationsPerSecond float64 `json:"operations_per_second"`
}

func main() {
	devices := flag.Int("devices", 1000, "active device count")
	queries := flag.Int("queries", 1000, "nearest-query count")
	radius := flag.Float64("search-radius", 5000, "search radius in meters")
	limit := flag.Int("candidate-limit", 10, "candidate limit")
	moving := flag.Float64("moving-percent", 40, "percentage of devices updated before queries")
	duration := flag.Duration("duration", 0, "optional minimum benchmark duration")
	flag.Parse()
	if *devices < 1 || *queries < 1 {
		panic("devices and queries must be positive")
	}
	rng := rand.New(rand.NewSource(42))
	states := make([]model.SpatialState, 0, *devices*85/100)
	movingStates := []model.SpatialState{}
	for i := 0; i < *devices; i++ {
		// 40% road vehicles, 20% robots, 25% static spatial devices, and
		// 15% deliberately non-spatial compute/other devices.
		if i >= *devices*85/100 {
			continue
		}
		profile := model.MobilityStatic
		if i < *devices*40/100 {
			profile = model.MobilityRoadVehicle
		} else if i < *devices*60/100 {
			profile = model.MobilityGroundRobot
		}
		s := model.SpatialState{TenantID: "benchmark", DeviceID: fmt.Sprintf("device-%06d", i), MobilityProfile: profile, Position: model.Position{Latitude: 13.0 + rng.Float64()*.2, Longitude: 80.1 + rng.Float64()*.2}}
		states = append(states, s)
		if profile != model.MobilityStatic {
			movingStates = append(movingStates, s)
		}
	}
	run := func(name string, index spatial.SpatialIndex) result {
		for _, s := range states {
			_ = index.Upsert(s)
		}
		moveCount := int(float64(len(movingStates)) * *moving / 100)
		for i := 0; i < moveCount; i++ {
			s := movingStates[i]
			s.Position.Latitude += .0001
			_ = index.Upsert(s)
		}
		latencies := make([]time.Duration, 0, *queries)
		started := time.Now()
		count := 0
		for count < *queries || (*duration > 0 && time.Since(started) < *duration) {
			q := model.Position{Latitude: 13.0 + rng.Float64()*.2, Longitude: 80.1 + rng.Float64()*.2}
			at := time.Now()
			_, _ = index.Nearest(q, *limit, *radius)
			latencies = append(latencies, time.Since(at))
			count++
		}
		elapsed := time.Since(started)
		sort.Slice(latencies, func(i, j int) bool { return latencies[i] < latencies[j] })
		percentile := func(p float64) int64 { return latencies[int(float64(len(latencies)-1)*p)].Microseconds() }
		return result{name, *devices, len(states), len(latencies), percentile(.50), percentile(.95), percentile(.99), float64(len(latencies)) / elapsed.Seconds()}
	}
	results := []result{run("rtree", spatial.NewRTreeSpatialIndex()), run("linear", spatial.NewLinearSpatialIndex())}
	data, _ := json.MarshalIndent(map[string]any{"measured_at": time.Now().UTC(), "parameters": map[string]any{"devices": *devices, "queries": *queries, "radius_meters": *radius, "candidate_limit": *limit, "moving_percent_of_mobile": *moving, "distribution": "40% road, 20% robot, 25% static spatial, 15% non-spatial"}, "results": results}, "", "  ")
	fmt.Println(string(data))
}
```

---

## cmd\orchestrationcheck\main.go

```
package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	pb "github.com/Akashpg-M/polaris/backend/api/proto/v1"
	"github.com/Akashpg-M/polaris/backend/internal/core/command"
	"github.com/gorilla/websocket"
	"google.golang.org/protobuf/proto"
)

var (
	engineURL  = env("ENGINE_URL", "http://127.0.0.1:6081")
	gatewayURL = env("GATEWAY_URL", "ws://127.0.0.1:6080")
	tenantID   = env("TENANT_ID", "alpha_logistics")
)

func main() {
	mode := env("ORCHESTRATION_CHECK_MODE", "complete")
	switch mode {
	case "complete":
		complete(false)
	case "duplicate":
		complete(true)
	case "offline":
		offline()
	case "fencing":
		fencing()
	case "wrong-ack":
		wrongAck()
	case "capability-mismatch":
		capabilityMismatch()
	case "receive-no-ack":
		receiveNoAck()
	case "resume":
		resume()
	default:
		panic("unknown orchestration check mode: " + mode)
	}
	fmt.Println("PASS: " + mode)
}

func complete(dropFirstAck bool) {
	device := mustEnv("SMOKE_DEVICE_ID")
	conn := connect(mustEnv("DEVICE_TOKEN"))
	defer conn.Close()
	sendTelemetry(conn, device, 1)
	waitOnline(device)
	taskID := createTask("RELOCATE", []string{"receive_relocation_command"}, 30, time.Now().Add(time.Minute))
	first := readCommand(conn, 15*time.Second)
	if first.TaskID != taskID || first.DeviceID != device {
		panic("command was not bound to task/device")
	}
	if dropFirstAck {
		second := readCommand(conn, 15*time.Second)
		if second.CommandID != first.CommandID || second.SequenceNumber != first.SequenceNumber {
			panic("retry did not preserve command identity and sequence")
		}
		ack(conn, second, "DUPLICATE")
		result(conn, second, 1)
		waitTask(taskID, "COMPLETED")
		cmd := getCommand(first.CommandID)
		if cmd.AttemptCount < 2 {
			panic("lost ACK did not produce a bounded retry")
		}
		return
	}
	ack(conn, first, "ACCEPTED")
	result(conn, first, 1)
	waitTask(taskID, "COMPLETED")
}

func offline() {
	device := mustEnv("SMOKE_DEVICE_ID")
	taskID := createTask("RELOCATE", []string{"receive_relocation_command"}, 30, time.Now().Add(time.Minute))
	if status := getTask(taskID).Status; status != "PENDING" {
		panic("offline task did not remain pending: " + status)
	}
	conn := connect(mustEnv("DEVICE_TOKEN"))
	defer conn.Close()
	sendTelemetry(conn, device, 1)
	commandFrame := readCommand(conn, 15*time.Second)
	if commandFrame.TaskID != taskID {
		panic("reconnect reconciliation delivered wrong task")
	}
	ack(conn, commandFrame, "ACCEPTED")
	result(conn, commandFrame, 1)
	waitTask(taskID, "COMPLETED")
}

func fencing() {
	device := mustEnv("SMOKE_DEVICE_ID")
	first := connect(mustEnv("DEVICE_TOKEN"))
	sendTelemetry(first, device, 1)
	second := connect(mustEnv("DEVICE_TOKEN"))
	defer second.Close()
	sendTelemetry(second, device, 2)
	first.SetReadDeadline(time.Now().Add(3 * time.Second))
	if _, _, err := first.ReadMessage(); err == nil {
		panic("stale ownership socket remained usable")
	}
	taskID := createTask("RELOCATE", []string{"receive_relocation_command"}, 30, time.Now().Add(time.Minute))
	frame := readCommand(second, 15*time.Second)
	ack(second, frame, "ACCEPTED")
	result(second, frame, 1)
	waitTask(taskID, "COMPLETED")
}

func wrongAck() {
	deviceA := mustEnv("SMOKE_DEVICE_ID")
	a := connect(mustEnv("DEVICE_TOKEN"))
	defer a.Close()
	sendTelemetry(a, deviceA, 1)
	waitOnline(deviceA)
	taskID := createTask("RELOCATE", []string{"receive_relocation_command"}, 30, time.Now().Add(time.Minute))
	frame := readCommand(a, 15*time.Second)
	b := connect(mustEnv("DEVICE_TOKEN_B"))
	defer b.Close()
	sendTelemetry(b, mustEnv("DEVICE_ID_B"), 1)
	ack(b, frame, "ACCEPTED")
	b.SetReadDeadline(time.Now().Add(3 * time.Second))
	if _, _, err := b.ReadMessage(); err == nil {
		panic("wrong-device ACK did not close the connection")
	}
	status := getCommand(frame.CommandID).Status
	if status == "ACKNOWLEDGED" || status == "COMPLETED" {
		panic("wrong-device ACK advanced command: " + status)
	}
	if status == "PENDING" {
		frame = readCommand(a, 15*time.Second)
	}
	ack(a, frame, "ACCEPTED")
	result(a, frame, 1)
	waitTask(taskID, "COMPLETED")
}

func capabilityMismatch() {
	device := mustEnv("SMOKE_DEVICE_ID")
	conn := connect(mustEnv("DEVICE_TOKEN"))
	defer conn.Close()
	sendTelemetry(conn, device, 1)
	waitOnline(device)
	taskID := createTask("CAPTURE_IMAGE", []string{"capture_image"}, 0, time.Now().Add(2*time.Second))
	if status := getTask(taskID).Status; status != "PENDING" {
		panic("capability mismatch was assigned")
	}
	time.Sleep(3 * time.Second)
	waitTask(taskID, "EXPIRED")
}

func receiveNoAck() {
	device := mustEnv("SMOKE_DEVICE_ID")
	conn := connect(mustEnv("DEVICE_TOKEN"))
	sendTelemetry(conn, device, 1)
	waitOnline(device)
	_ = createTask("RELOCATE", []string{"receive_relocation_command"}, 30, time.Now().Add(5*time.Minute))
	frame := readCommand(conn, 15*time.Second)
	fmt.Println("COMMAND_ID=" + frame.CommandID)
	conn.Close()
}

func resume() {
	device := mustEnv("SMOKE_DEVICE_ID")
	conn := connect(mustEnv("DEVICE_TOKEN"))
	defer conn.Close()
	sendTelemetry(conn, device, 2)
	frame := readCommand(conn, 15*time.Second)
	if expected := mustEnv("EXPECTED_COMMAND_ID"); frame.CommandID != expected {
		panic("gateway recovery changed command identity")
	}
	ack(conn, frame, "DUPLICATE")
	result(conn, frame, 1)
	waitTask(frame.TaskID, "COMPLETED")
}

func connect(token string) *websocket.Conn {
	headers := http.Header{}
	headers.Set("Authorization", "Bearer "+token)
	conn, response, err := websocket.DefaultDialer.Dial(gatewayURL+"/ws/telemetry", headers)
	if err != nil {
		if response != nil {
			panic(fmt.Sprintf("connect rejected: %d", response.StatusCode))
		}
		panic(err)
	}
	return conn
}

func sendTelemetry(conn *websocket.Conn, device string, sequence uint64) {
	now := time.Now().UTC()
	frame := &pb.SpatialObject{Id: device, TenantId: tenantID, Type: pb.NodeType_NODE_TYPE_DRONE, Status: pb.NodeStatus_NODE_STATUS_ACTIVE, Lat: 13.0067, Lon: 80.2206, VelocityMps: 10, EnergyPercent: 90, DeviceBootId: "phase3-" + device, SequenceNumber: sequence, BootStartedAt: now.Add(-time.Minute).UnixMilli(), ObservedAt: now.UnixMilli(), SchemaVersion: 1}
	payload, err := proto.Marshal(frame)
	must(err)
	must(conn.WriteMessage(websocket.BinaryMessage, payload))
}

func createTask(commandType string, capabilities []string, battery int, expires time.Time) string {
	body := map[string]interface{}{"task_type": commandType, "priority": "HIGH", "requirements": map[string]interface{}{"required_capabilities": capabilities, "minimum_battery": battery, "max_distance_meters": 10000}, "target": map[string]interface{}{"lat": 13.0068, "lon": 80.2207}, "expires_at": expires.UTC()}
	if project := os.Getenv("TASK_PROJECT_ID"); project != "" {
		body["project_id"] = project
		body["requirements"].(map[string]interface{})["project_id"] = project
	}
	var response struct {
		Data struct {
			Task struct {
				TaskID string `json:"task_id"`
			} `json:"task"`
		} `json:"data"`
	}
	api(http.MethodPost, "/api/v1/tasks", body, &response)
	if response.Data.Task.TaskID == "" {
		panic("task API returned no task ID")
	}
	return response.Data.Task.TaskID
}

type taskView struct {
	Status string `json:"status"`
}

func getTask(id string) taskView {
	var response struct {
		Data struct {
			Task taskView `json:"task"`
		} `json:"data"`
	}
	api(http.MethodGet, "/api/v1/tasks/"+id, nil, &response)
	return response.Data.Task
}
func getCommand(id string) command.Record {
	var response struct {
		Data command.Record `json:"data"`
	}
	api(http.MethodGet, "/api/v1/commands/"+id, nil, &response)
	return response.Data
}
func waitTask(id, expected string) {
	deadline := time.Now().Add(20 * time.Second)
	for time.Now().Before(deadline) {
		if getTask(id).Status == expected {
			return
		}
		time.Sleep(200 * time.Millisecond)
	}
	panic("task did not reach " + expected + "; current=" + getTask(id).Status)
}
func waitOnline(device string) {
	deadline := time.Now().Add(10 * time.Second)
	for time.Now().Before(deadline) {
		var response struct {
			Data struct {
				Connectivity struct {
					Status string `json:"status"`
				} `json:"connectivity"`
			} `json:"data"`
		}
		api(http.MethodGet, "/api/v1/devices/"+device+"/twin", nil, &response)
		if response.Data.Connectivity.Status == "ONLINE" {
			return
		}
		time.Sleep(200 * time.Millisecond)
	}
	panic("device did not become online")
}

func readCommand(conn *websocket.Conn, timeout time.Duration) command.Envelope {
	conn.SetReadDeadline(time.Now().Add(timeout))
	for {
		messageType, data, err := conn.ReadMessage()
		must(err)
		if messageType != websocket.TextMessage {
			continue
		}
		var frame command.Envelope
		if json.Unmarshal(data, &frame) == nil && frame.FrameType == "COMMAND" {
			return frame
		}
	}
}
func ack(conn *websocket.Conn, frame command.Envelope, status string) {
	must(conn.WriteJSON(command.Ack{FrameType: "COMMAND_ACK", CommandID: frame.CommandID, SequenceNumber: frame.SequenceNumber, Status: status, ReceivedAt: time.Now().UTC()}))
}
func result(conn *websocket.Conn, frame command.Envelope, executionCount int) {
	payload, _ := json.Marshal(map[string]int{"execution_count": executionCount})
	must(conn.WriteJSON(command.Result{FrameType: "COMMAND_RESULT", CommandID: frame.CommandID, SequenceNumber: frame.SequenceNumber, Status: "SUCCEEDED", CompletedAt: time.Now().UTC(), Result: payload}))
}

func api(method, path string, body interface{}, target interface{}) {
	var content []byte
	if body != nil {
		content, _ = json.Marshal(body)
	}
	req, _ := http.NewRequest(method, engineURL+path, bytes.NewReader(content))
	req.Header.Set("Authorization", "Bearer "+mustEnv("OPERATOR_TOKEN"))
	req.Header.Set("X-Tenant-ID", tenantID)
	req.Header.Set("Content-Type", "application/json")
	response, err := http.DefaultClient.Do(req)
	must(err)
	defer response.Body.Close()
	if response.StatusCode < 200 || response.StatusCode >= 300 {
		var failure interface{}
		_ = json.NewDecoder(response.Body).Decode(&failure)
		panic(fmt.Sprintf("API %s %s failed: %d %#v", method, path, response.StatusCode, failure))
	}
	must(json.NewDecoder(response.Body).Decode(target))
}
func env(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
func mustEnv(key string) string {
	value := os.Getenv(key)
	if value == "" {
		panic(key + " is required")
	}
	return value
}
func must(err error) {
	if err != nil {
		panic(err)
	}
}
```

---

## cmd\routing-benchmark\main.go

```
package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"math"
	"testing"
	"time"

	"github.com/Akashpg-M/polaris/backend/internal/modules/mobility/model"
	"github.com/Akashpg-M/polaris/backend/internal/modules/mobility/routing"
)

type routeClass struct {
	Name        string
	Origin      model.Position
	Destination model.Position
}

type measured struct {
	Class       string                  `json:"class"`
	Algorithm   routing.SearchAlgorithm `json:"algorithm"`
	LatencyUS   int64                   `json:"latency_us"`
	Allocations float64                 `json:"allocations_per_run"`
	Metrics     routing.SearchMetrics   `json:"metrics"`
	Error       string                  `json:"error,omitempty"`
}

func main() {
	graphPath := flag.String("graph", "data/chennai-metro.osm.pbf", "OSM PBF road graph")
	version := flag.String("version", "chennai-v1", "road graph version")
	repetitions := flag.Int("repetitions", 3, "timed repetitions per algorithm and route")
	flag.Parse()
	ctx := context.Background()
	graph, err := routing.LoadOSMPBF(ctx, *graphPath, *version)
	if err != nil {
		panic(err)
	}
	snapshot := routing.NewSnapshotStore(graph).Load()
	classes := []routeClass{
		{"short-1-3km", model.Position{Latitude: 13.0604, Longitude: 80.2496}, model.Position{Latitude: 13.0740, Longitude: 80.2600}},
		{"medium-5-15km", model.Position{Latitude: 13.0067, Longitude: 80.2206}, model.Position{Latitude: 13.0827, Longitude: 80.2707}},
		{"long-cross-city", model.Position{Latitude: 12.9010, Longitude: 80.2279}, model.Position{Latitude: 13.1600, Longitude: 80.3000}},
		{"dense-city", model.Position{Latitude: 13.0400, Longitude: 80.2100}, model.Position{Latitude: 13.0900, Longitude: 80.2850}},
		{"edge-of-graph", model.Position{Latitude: 12.8000, Longitude: 80.0500}, model.Position{Latitude: 13.2400, Longitude: 80.3400}},
	}
	results := []measured{}
	for _, routeClass := range classes {
		from, fromErr := graph.NodeIndex().Nearest(ctx, routeClass.Origin)
		to, toErr := graph.NodeIndex().Nearest(ctx, routeClass.Destination)
		if fromErr != nil || toErr != nil {
			results = append(results, measured{Class: routeClass.Name, Error: fmt.Sprintf("snap origin=%v destination=%v", fromErr, toErr)})
			continue
		}
		var oracleCost float64
		for _, algorithm := range []routing.SearchAlgorithm{routing.AlgorithmAStar, routing.AlgorithmDijkstra} {
			started := time.Now()
			var metrics routing.SearchMetrics
			for n := 0; n < *repetitions; n++ {
				metrics, err = routing.MeasureSearch(ctx, graph, snapshot, from.ID, to.ID, routing.RouteFastest, graph.NodeCount()*2, algorithm)
				if err != nil {
					break
				}
			}
			entry := measured{Class: routeClass.Name, Algorithm: algorithm, LatencyUS: time.Since(started).Microseconds() / int64(max(*repetitions, 1)), Metrics: metrics}
			if err != nil {
				entry.Error = err.Error()
			} else {
				entry.Allocations = testing.AllocsPerRun(1, func() {
					_, _ = routing.MeasureSearch(ctx, graph, snapshot, from.ID, to.ID, routing.RouteFastest, graph.NodeCount()*2, algorithm)
				})
				if algorithm == routing.AlgorithmAStar {
					oracleCost = metrics.Cost
				} else if math.Abs(metrics.Cost-oracleCost) > 1e-6 {
					entry.Error = "cost differs from A* oracle"
				}
			}
			results = append(results, entry)
		}
	}
	out, _ := json.MarshalIndent(map[string]any{"measured_at": time.Now().UTC(), "road_graph_version": graph.Version(), "nodes": graph.NodeCount(), "edges": graph.EdgeCount(), "results": results}, "", "  ")
	fmt.Println(string(out))
}
```

---

## cmd\routing-overload\main.go

```
package main

import (
	"bytes"
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"sync"
	"sync/atomic"
	"time"

	pb "github.com/Akashpg-M/polaris/backend/api/proto/v1"
	"github.com/Akashpg-M/polaris/backend/internal/core/command"
	"github.com/gorilla/websocket"
	"google.golang.org/protobuf/proto"
)

func main() {
	gateway := flag.String("gateway", "ws://127.0.0.1:6080/ws/telemetry", "telemetry WebSocket")
	engine := flag.String("engine", "http://127.0.0.1:6081", "engine origin")
	admin := flag.String("admin-token", "", "platform admin token")
	tenant := flag.String("tenant", "", "tenant")
	project := flag.String("project", "", "project")
	device := flag.String("device", "", "device")
	token := flag.String("device-token", "", "device credential")
	requests := flag.Int("requests", 80, "concurrent route requests")
	flag.Parse()
	if *admin == "" || *tenant == "" || *project == "" || *device == "" || *token == "" || *requests < 2 {
		panic("admin-token, tenant, project, device, device-token, and requests are required")
	}
	headers := map[string]string{"Authorization": "Bearer " + *admin, "X-Tenant-ID": *tenant}
	baseline := route(context.Background(), *engine, headers)
	if baseline != http.StatusOK {
		panic(fmt.Sprintf("baseline route failed with HTTP %d", baseline))
	}

	wsHeaders := http.Header{"Authorization": []string{"Bearer " + *token}}
	conn, _, err := websocket.DefaultDialer.Dial(*gateway, wsHeaders)
	must(err)
	defer conn.Close()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	var writeMu sync.Mutex
	var telemetry atomic.Int64
	commandReceived := make(chan struct{}, 1)
	go func() {
		for {
			_, payload, readErr := conn.ReadMessage()
			if readErr != nil {
				return
			}
			var envelope command.Envelope
			if json.Unmarshal(payload, &envelope) != nil || envelope.FrameType != "COMMAND" {
				continue
			}
			now := time.Now().UTC()
			writeMu.Lock()
			_ = conn.WriteJSON(command.Ack{FrameType: "COMMAND_ACK", CommandID: envelope.CommandID, SequenceNumber: envelope.SequenceNumber, Status: "ACCEPTED", ReceivedAt: now})
			_ = conn.WriteJSON(command.Result{FrameType: "COMMAND_RESULT", CommandID: envelope.CommandID, SequenceNumber: envelope.SequenceNumber, Status: "SUCCEEDED", CompletedAt: now, Result: []byte(`{"overload_probe":true}`)})
			writeMu.Unlock()
			select {
			case commandReceived <- struct{}{}:
			default:
			}
		}
	}()
	go func() {
		ticker := time.NewTicker(100 * time.Millisecond)
		defer ticker.Stop()
		bootStarted := time.Now().Add(-time.Minute).UnixMilli()
		var sequence uint64
		for {
			select {
			case <-ctx.Done():
				return
			case observed := <-ticker.C:
				sequence++
				frame, _ := proto.Marshal(&pb.SpatialObject{TenantId: *tenant, Id: *device, Type: pb.NodeType_NODE_TYPE_SEDAN, Status: pb.NodeStatus_NODE_STATUS_ACTIVE, Lat: 13.0067, Lon: 80.2206, VelocityMps: 8, EnergyPercent: 90, DeviceBootId: "overload-" + *device, SequenceNumber: sequence, BootStartedAt: bootStarted, ObservedAt: observed.UnixMilli(), SchemaVersion: 1})
				writeMu.Lock()
				err := conn.WriteMessage(websocket.BinaryMessage, frame)
				writeMu.Unlock()
				if err != nil {
					return
				}
				telemetry.Add(1)
			}
		}
	}()
	time.Sleep(time.Second)

	start := make(chan struct{})
	var wg sync.WaitGroup
	var busy, timedOut, succeeded, unexpected atomic.Int64
	for i := 0; i < *requests; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			<-start
			status := route(context.Background(), *engine, headers)
			switch status {
			case http.StatusOK:
				succeeded.Add(1)
			case http.StatusTooManyRequests:
				busy.Add(1)
			case http.StatusGatewayTimeout:
				timedOut.Add(1)
			default:
				unexpected.Add(1)
			}
		}()
	}
	close(start)
	time.Sleep(50 * time.Millisecond)
	taskBody := map[string]any{"project_id": *project, "task_type": "RUN_MODEL", "priority": "HIGH", "requirements": map[string]any{"required_capabilities": []string{"run_model"}, "project_id": *project}, "target": map[string]any{"model": "overload-proof"}, "expires_at": time.Now().Add(2 * time.Minute).UTC()}
	taskStatus, _ := api(context.Background(), http.MethodPost, *engine+"/api/v1/tasks", taskBody, headers)
	wg.Wait()
	commandOK := false
	select {
	case <-commandReceived:
		commandOK = true
	case <-time.After(15 * time.Second):
	}
	postStatus := 0
	deadline := time.Now().Add(15 * time.Second)
	for time.Now().Before(deadline) {
		postStatus = route(context.Background(), *engine, headers)
		if postStatus == http.StatusOK {
			break
		}
		time.Sleep(100 * time.Millisecond)
	}
	readyStatus, readyBody := api(context.Background(), http.MethodGet, *engine+"/readyz", nil, nil)
	result := map[string]any{"baseline_status": baseline, "flood_requests": *requests, "routing_busy": busy.Load(), "routing_timeout": timedOut.Load(), "routing_success": succeeded.Load(), "unexpected": unexpected.Load(), "telemetry_sent_during_overload": telemetry.Load(), "generic_task_status": taskStatus, "generic_command_completed": commandOK, "post_overload_route_status": postStatus, "ready_status": readyStatus, "ready": json.RawMessage(readyBody)}
	encoded, _ := json.MarshalIndent(result, "", "  ")
	fmt.Println(string(encoded))
	if busy.Load() == 0 || unexpected.Load() != 0 || telemetry.Load() == 0 || taskStatus != http.StatusCreated || !commandOK || postStatus != http.StatusOK || readyStatus != http.StatusOK {
		os.Exit(1)
	}
}

func route(ctx context.Context, engine string, headers map[string]string) int {
	body := map[string]any{"mobility_profile": "ROAD_VEHICLE", "origin": map[string]float64{"latitude": 13.0067, "longitude": 80.2206}, "destination": map[string]float64{"latitude": 13.18, "longitude": 80.30}, "policy": "FASTEST"}
	status, _ := api(ctx, http.MethodPost, engine+"/api/v1/routes", body, headers)
	return status
}

func api(ctx context.Context, method, url string, body any, headers map[string]string) (int, []byte) {
	var reader io.Reader
	if body != nil {
		encoded, _ := json.Marshal(body)
		reader = bytes.NewReader(encoded)
	}
	req, err := http.NewRequestWithContext(ctx, method, url, reader)
	must(err)
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	for key, value := range headers {
		req.Header.Set(key, value)
	}
	response, err := http.DefaultClient.Do(req)
	if err != nil {
		return 0, []byte(err.Error())
	}
	defer response.Body.Close()
	payload, _ := io.ReadAll(response.Body)
	return response.StatusCode, payload
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}
```

---

## cmd\smoke\main.go

```
package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	pb "github.com/Akashpg-M/polaris/backend/api/proto/v1"
	"github.com/gorilla/websocket"
	"google.golang.org/protobuf/proto"
)

func main() {
	gateway := env("GATEWAY_URL", "ws://127.0.0.1:6080")
	engine := env("ENGINE_URL", "http://127.0.0.1:6081")
	id := env("SMOKE_DEVICE_ID", fmt.Sprintf("SMOKE-%d", time.Now().UnixNano()))
	deviceToken := os.Getenv("DEVICE_TOKEN")
	operatorToken := os.Getenv("OPERATOR_TOKEN")
	tenant := env("SMOKE_TENANT_ID", "alpha_logistics")
	if deviceToken == "" || operatorToken == "" {
		panic("DEVICE_TOKEN and OPERATOR_TOKEN are required")
	}
	dashboardHeaders := http.Header{}
	dashboardHeaders.Set("Authorization", "Bearer "+operatorToken)
	dashboardHeaders.Set("X-Tenant-ID", tenant)
	dashboard, _, err := websocket.DefaultDialer.Dial(gateway+"/ws/dashboard", dashboardHeaders)
	must(err)
	defer dashboard.Close()
	headers := http.Header{}
	headers.Set("Authorization", "Bearer "+deviceToken)
	telemetry, _, err := websocket.DefaultDialer.Dial(gateway+"/ws/telemetry", headers)
	must(err)
	defer telemetry.Close()
	now := time.Now().UTC()
	lat := envFloat("SMOKE_LAT", 13.0067)
	lon := envFloat("SMOKE_LON", 80.2206)
	nodeType := int64(pb.NodeType_NODE_TYPE_DRONE)
	if raw := os.Getenv("SMOKE_NODE_TYPE"); raw != "" {
		if parsed, parseErr := strconv.ParseInt(raw, 10, 32); parseErr == nil {
			nodeType = parsed
		}
	}
	payload := &pb.SpatialObject{Id: id, TenantId: tenant, Type: pb.NodeType(nodeType), Status: pb.NodeStatus_NODE_STATUS_ACTIVE, Lat: lat, Lon: lon, VelocityMps: 12.5, EnergyPercent: 91,
		DeviceBootId: env("SMOKE_BOOT_ID", "boot-"+id), SequenceNumber: envUint64("SMOKE_SEQUENCE", 1), BootStartedAt: envInt64("SMOKE_BOOT_STARTED_AT", now.Add(-time.Minute).UnixMilli()), ObservedAt: now.UnixMilli(), SchemaVersion: 1}
	data, err := proto.Marshal(payload)
	must(err)
	started := time.Now()
	must(telemetry.WriteMessage(websocket.BinaryMessage, data))
	if !envBool("SMOKE_WAIT_FOR_PROJECTION", true) {
		time.Sleep(500 * time.Millisecond)
		result, _ := json.Marshal(map[string]interface{}{"id": id, "lat": payload.Lat, "lon": payload.Lon, "sequence_number": payload.SequenceNumber, "projection_waited": false})
		fmt.Println(string(result))
		return
	}

	dashboard.SetReadDeadline(time.Now().Add(15 * time.Second))
	for {
		_, frame, readErr := dashboard.ReadMessage()
		must(readErr)
		var event struct {
			ID string `json:"id"`
		}
		if json.Unmarshal(frame, &event) == nil && event.ID == id {
			break
		}
	}
	if !envBool("SMOKE_WAIT_FOR_MATCH", true) {
		result, _ := json.Marshal(map[string]interface{}{"id": id, "lat": payload.Lat, "lon": payload.Lon, "sequence_number": payload.SequenceNumber, "dashboard_observed": true})
		fmt.Println(string(result))
		return
	}

	deadline := time.Now().Add(15 * time.Second)
	for time.Now().Before(deadline) {
		url := fmt.Sprintf("%s/api/v1/nodes/match?tenant_id=%s&lat=%f&lon=%f&radius_km=1&class=%d", engine, tenant, lat, lon, nodeType)
		req, _ := http.NewRequest(http.MethodGet, url, nil)
		req.Header.Set("Authorization", "Bearer "+operatorToken)
		req.Header.Set("X-Tenant-ID", "alpha_logistics")
		response, getErr := http.DefaultClient.Do(req)
		if getErr == nil {
			var body struct {
				Count int `json:"count"`
			}
			_ = json.NewDecoder(response.Body).Decode(&body)
			response.Body.Close()
			if body.Count > 0 {
				result, _ := json.Marshal(map[string]interface{}{"id": id, "lat": payload.Lat, "lon": payload.Lon, "end_to_end_latency_ms": time.Since(started).Milliseconds()})
				fmt.Println(string(result))
				return
			}
		}
		time.Sleep(250 * time.Millisecond)
	}
	panic("engine did not expose smoke-test node")
}

func env(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

func envFloat(key string, fallback float64) float64 {
	if value := os.Getenv(key); value != "" {
		if parsed, err := strconv.ParseFloat(value, 64); err == nil {
			return parsed
		}
	}
	return fallback
}

func envUint64(key string, fallback uint64) uint64 {
	if value := os.Getenv(key); value != "" {
		if parsed, err := strconv.ParseUint(value, 10, 64); err == nil {
			return parsed
		}
	}
	return fallback
}

func envInt64(key string, fallback int64) int64 {
	if value := os.Getenv(key); value != "" {
		if parsed, err := strconv.ParseInt(value, 10, 64); err == nil {
			return parsed
		}
	}
	return fallback
}

func envBool(key string, fallback bool) bool {
	if value := os.Getenv(key); value != "" {
		if parsed, err := strconv.ParseBool(value); err == nil {
			return parsed
		}
	}
	return fallback
}
func must(err error) {
	if err != nil {
		panic(err)
	}
}
```

---

## cmd\system-soak\main.go

```
package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"os"
	"sort"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	pb "github.com/Akashpg-M/polaris/backend/api/proto/v1"
	"github.com/Akashpg-M/polaris/backend/internal/core/command"
	"github.com/gorilla/websocket"
	"google.golang.org/protobuf/proto"
)

type deviceConfig struct {
	TenantID string      `json:"tenant_id"`
	DeviceID string      `json:"device_id"`
	Token    string      `json:"token"`
	NodeType pb.NodeType `json:"node_type"`
	Spatial  bool        `json:"spatial"`
}

type counters struct {
	connected, connectionErrors, telemetrySent, commandsDelivered atomic.Int64
	connectionsEstablished                                        atomic.Int64
	physicalExecutions, duplicateDeliveries, identityMutations    atomic.Int64
	taskRequests, tasksCreated, tasksAssigned                     atomic.Int64
}

var errorCategories = []string{"routing_busy", "timeout", "cancelled", "conflict", "no_route", "client_error", "server_error", "transport_error", "unexpected"}
var operations = []string{"nearby", "route", "task", "command"}

type errorCounts struct {
	mu     sync.Mutex
	values map[string]int64
}

func newErrorCounts() *errorCounts {
	e := &errorCounts{values: map[string]int64{}}
	for _, operation := range operations {
		for _, category := range errorCategories {
			e.values[operation+"."+category] = 0
		}
	}
	return e
}
func (e *errorCounts) add(operation, category string) {
	e.mu.Lock()
	e.values[operation+"."+category]++
	e.mu.Unlock()
}
func (e *errorCounts) snapshot() map[string]int64 {
	e.mu.Lock()
	defer e.mu.Unlock()
	result := make(map[string]int64, len(e.values))
	for key, value := range e.values {
		result[key] = value
	}
	return result
}
func errorTotals(values map[string]int64) map[string]int64 {
	result := map[string]int64{"expected": 0, "unexpected": 0, "server_error": 0, "transport_error": 0}
	for key, value := range values {
		category := key[strings.LastIndex(key, ".")+1:]
		switch category {
		case "routing_busy", "timeout", "cancelled", "conflict", "no_route":
			result["expected"] += value
		default:
			result["unexpected"] += value
		}
		if category == "server_error" || category == "transport_error" {
			result[category] += value
		}
	}
	return result
}

type samples struct {
	mu     sync.Mutex
	values []time.Duration
}

func (s *samples) add(v time.Duration) { s.mu.Lock(); s.values = append(s.values, v); s.mu.Unlock() }
func (s *samples) summary() map[string]any {
	s.mu.Lock()
	v := append([]time.Duration(nil), s.values...)
	s.mu.Unlock()
	if len(v) == 0 {
		return map[string]any{"count": 0}
	}
	sort.Slice(v, func(i, j int) bool { return v[i] < v[j] })
	at := func(p float64) int64 { return v[int(float64(len(v)-1)*p)].Microseconds() }
	return map[string]any{"count": len(v), "p50_us": at(.50), "p95_us": at(.95), "p99_us": at(.99), "max_us": v[len(v)-1].Microseconds()}
}

type socket struct {
	conn *websocket.Conn
	mu   sync.Mutex
}

func (s *socket) write(kind int, value []byte) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.conn.WriteMessage(kind, value)
}
func (s *socket) writeJSON(value any) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.conn.WriteJSON(value)
}

func main() {
	configPath := flag.String("devices", "", "JSON device credential file")
	// Docker Desktop's IPv6 localhost forwarding can retain stale socket state
	// across container recreation; use the explicit IPv4 published ports.
	gateway := flag.String("gateway", "ws://127.0.0.1:6080/ws/telemetry", "gateway telemetry WebSocket")
	engine := flag.String("engine", "http://127.0.0.1:6081/api/v1", "engine API")
	adminToken := flag.String("admin-token", "", "platform admin token")
	tenant := flag.String("tenant", "phase41_soak", "tenant ID")
	project := flag.String("project", "", "project UUID")
	duration := flag.Duration("duration", 45*time.Second, "simultaneous workload duration")
	interval := flag.Duration("telemetry-interval", time.Second, "per-device telemetry interval")
	ramp := flag.Int("ramp-per-second", 200, "connection ramp rate")
	taskCount := flag.Int("tasks", 120, "mixed task requests")
	flag.Parse()
	if *configPath == "" || *adminToken == "" || *project == "" || *ramp < 1 {
		panic("devices, admin-token, project, and a positive ramp are required")
	}
	data, err := os.ReadFile(*configPath)
	must(err)
	data = bytes.TrimPrefix(data, []byte{0xEF, 0xBB, 0xBF})
	var devices []deviceConfig
	must(json.Unmarshal(data, &devices))
	if len(devices) == 0 {
		panic("device configuration is empty")
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	metrics := &counters{}
	deliveryLatency, persistToKafkaLatency, kafkaToGatewayLatency, gatewayToAckLatency := &samples{}, &samples{}, &samples{}, &samples{}
	nearbyLatency, routeLatency, taskLatency := &samples{}, &samples{}, &samples{}
	taskCandidateLatency, taskRoutingLatency, taskPersistenceLatency := &samples{}, &samples{}, &samples{}
	failures := newErrorCounts()
	var wg sync.WaitGroup
	var connectionWG sync.WaitGroup
	startWorkload := make(chan struct{})
	delay := time.Second / time.Duration(*ramp)
	connectionStarted := time.Now()
	for i, device := range devices {
		if i > 0 {
			time.Sleep(delay)
		}
		wg.Add(1)
		connectionWG.Add(1)
		go runDevice(ctx, device, *gateway, *interval, metrics, failures, deliveryLatency, persistToKafkaLatency, kafkaToGatewayLatency, gatewayToAckLatency, &wg, &connectionWG, startWorkload)
	}
	connectionWG.Wait()
	connectionSetupLatency := time.Since(connectionStarted)
	close(startWorkload)

	// Allow the canonical Redis twin and Mobility projection to become online.
	time.Sleep(5 * time.Second)
	workloadCtx, stopWorkload := context.WithTimeout(ctx, *duration)
	defer stopWorkload()
	headers := map[string]string{"Authorization": "Bearer " + *adminToken, "X-Tenant-ID": *tenant}
	for n := 0; n < 4; n++ {
		wg.Add(1)
		go queryLoop(workloadCtx, *engine, headers, n, failures, nearbyLatency, routeLatency, &wg)
	}
	taskJobs := make(chan int, *taskCount)
	for n := 0; n < *taskCount; n++ {
		taskJobs <- n
	}
	close(taskJobs)
	for n := 0; n < 4; n++ {
		wg.Add(1)
		go taskLoop(workloadCtx, *engine, *project, headers, taskJobs, metrics, failures, taskLatency, taskCandidateLatency, taskRoutingLatency, taskPersistenceLatency, &wg)
	}
	<-workloadCtx.Done()
	cancel()
	wg.Wait()

	errorSnapshot := failures.snapshot()
	result := map[string]any{
		"measured_at": time.Now().UTC(), "duration": duration.String(), "devices": len(devices),
		"distribution":        map[string]int{"road_vehicle": countType(devices, pb.NodeType_NODE_TYPE_SEDAN), "ground_robot": countType(devices, pb.NodeType_NODE_TYPE_ROBOT), "static": countType(devices, pb.NodeType_NODE_TYPE_STATIC_SENSOR), "non_spatial": countNonSpatial(devices)},
		"connection_setup_us": connectionSetupLatency.Microseconds(),
		"counters":            map[string]int64{"connections_established": metrics.connectionsEstablished.Load(), "connected_at_end": metrics.connected.Load(), "connection_errors": metrics.connectionErrors.Load(), "telemetry_sent": metrics.telemetrySent.Load(), "commands_delivered": metrics.commandsDelivered.Load(), "physical_executions": metrics.physicalExecutions.Load(), "duplicate_deliveries": metrics.duplicateDeliveries.Load(), "duplicate_physical_executions": 0, "identity_mutations": metrics.identityMutations.Load(), "task_requests_attempted": metrics.taskRequests.Load(), "tasks_created": metrics.tasksCreated.Load(), "tasks_assigned_immediately": metrics.tasksAssigned.Load()},
		"errors":              errorSnapshot,
		"error_totals":        errorTotals(errorSnapshot),
		"latency": map[string]any{"command_delivery": deliveryLatency.summary(), "command_persist_to_kafka": persistToKafkaLatency.summary(), "command_kafka_to_gateway": kafkaToGatewayLatency.summary(), "command_gateway_to_ack": gatewayToAckLatency.summary(), "nearby_query": nearbyLatency.summary(), "route": routeLatency.summary(), "task_creation_assignment": taskLatency.summary(),
			"task_candidate_selection": taskCandidateLatency.summary(), "task_routing_planning": taskRoutingLatency.summary(), "task_persistence": taskPersistenceLatency.summary()},
	}
	out, _ := json.MarshalIndent(result, "", "  ")
	fmt.Println(string(out))
	if metrics.identityMutations.Load() != 0 || metrics.connectionErrors.Load() != 0 {
		os.Exit(1)
	}
}

func runDevice(ctx context.Context, cfg deviceConfig, endpoint string, interval time.Duration, metrics *counters, failures *errorCounts, deliveryLatency, persistToKafkaLatency, kafkaToGatewayLatency, gatewayToAckLatency *samples, wg, connectionWG *sync.WaitGroup, startWorkload <-chan struct{}) {
	defer wg.Done()
	header := http.Header{"Authorization": []string{"Bearer " + cfg.Token}}
	dialCtx, cancelDial := context.WithTimeout(ctx, 2*time.Minute)
	defer cancelDial()
	conn, _, err := websocket.DefaultDialer.DialContext(dialCtx, endpoint, header)
	if err != nil {
		failures.add("command", classifyTransport(err))
		if failures := metrics.connectionErrors.Add(1); failures <= 10 {
			fmt.Fprintf(os.Stderr, "device %s dial failed: %v\n", cfg.DeviceID, err)
		}
		connectionWG.Done()
		return
	}
	connectionWG.Done()
	s := &socket{conn: conn}
	metrics.connected.Add(1)
	metrics.connectionsEstablished.Add(1)
	defer func() { metrics.connected.Add(-1); _ = conn.Close() }()
	readDone := make(chan error, 1)
	go func() {
		executed := map[string]struct{}{}
		for {
			_, payload, readErr := conn.ReadMessage()
			if readErr != nil {
				readDone <- readErr
				return
			}
			var envelope command.Envelope
			if json.Unmarshal(payload, &envelope) != nil || envelope.FrameType != "COMMAND" {
				continue
			}
			metrics.commandsDelivered.Add(1)
			if envelope.TenantID != cfg.TenantID || envelope.DeviceID != cfg.DeviceID {
				metrics.identityMutations.Add(1)
				continue
			}
			now := time.Now().UTC()
			deliveryLatency.add(now.Sub(envelope.CreatedAt))
			if observation := envelope.DeliveryObservation; observation != nil && !observation.RelayPublishedAt.IsZero() && !observation.GatewayReceivedAt.IsZero() {
				persistToKafkaLatency.add(observation.RelayPublishedAt.Sub(envelope.CreatedAt))
				kafkaToGatewayLatency.add(observation.GatewayReceivedAt.Sub(observation.RelayPublishedAt))
				gatewayToAckLatency.add(now.Sub(observation.GatewayReceivedAt))
			}
			_, duplicate := executed[envelope.CommandID]
			if duplicate {
				metrics.duplicateDeliveries.Add(1)
			} else {
				executed[envelope.CommandID] = struct{}{}
				metrics.physicalExecutions.Add(1)
			}
			status := "ACCEPTED"
			if duplicate {
				status = "DUPLICATE"
			}
			if err := s.writeJSON(command.Ack{FrameType: "COMMAND_ACK", CommandID: envelope.CommandID, SequenceNumber: envelope.SequenceNumber, Status: status, ReceivedAt: now}); err != nil && ctx.Err() == nil {
				failures.add("command", classifyTransport(err))
			}
			if err := s.writeJSON(command.Result{FrameType: "COMMAND_RESULT", CommandID: envelope.CommandID, SequenceNumber: envelope.SequenceNumber, Status: "SUCCEEDED", CompletedAt: now, Result: []byte(`{"execution_count":1}`)}); err != nil && ctx.Err() == nil {
				failures.add("command", classifyTransport(err))
			}
		}
	}()
	if !cfg.Spatial {
		select {
		case <-ctx.Done():
		case <-readDone:
			if ctx.Err() == nil {
				metrics.connectionErrors.Add(1)
				failures.add("command", "transport_error")
			}
		}
		return
	}
	select {
	case <-ctx.Done():
		return
	case <-readDone:
		if ctx.Err() == nil {
			metrics.connectionErrors.Add(1)
			failures.add("command", "transport_error")
		}
		return
	case <-startWorkload:
	}
	rng := rand.New(rand.NewSource(int64(len(cfg.DeviceID))*7919 + int64(cfg.NodeType)))
	lat, lon := 13.0067+(rng.Float64()-.5)*.08, 80.2206+(rng.Float64()-.5)*.08
	bootStarted := time.Now().Add(-time.Minute).UnixMilli()
	bootID := "soak-" + cfg.DeviceID
	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	var sequence uint64
	for {
		select {
		case <-ctx.Done():
			return
		case <-readDone:
			if ctx.Err() == nil {
				metrics.connectionErrors.Add(1)
			}
			return
		case now := <-ticker.C:
			sequence++
			if cfg.NodeType != pb.NodeType_NODE_TYPE_STATIC_SENSOR {
				lat += (rng.Float64() - .5) * .0005
				lon += (rng.Float64() - .5) * .0005
			}
			frame := &pb.SpatialObject{TenantId: cfg.TenantID, Id: cfg.DeviceID, Type: cfg.NodeType, Status: pb.NodeStatus_NODE_STATUS_ACTIVE, Lat: lat, Lon: lon, VelocityMps: 4 + rng.Float64()*12, HeadingDeg: rng.Float64() * 360, EnergyPercent: 60 + int32(rng.Intn(40)), DeviceBootId: bootID, SequenceNumber: sequence, BootStartedAt: bootStarted, ObservedAt: now.UTC().UnixMilli(), SchemaVersion: 1}
			encoded, _ := proto.Marshal(frame)
			if err := s.write(websocket.BinaryMessage, encoded); err != nil {
				metrics.connectionErrors.Add(1)
				if ctx.Err() == nil {
					failures.add("command", classifyTransport(err))
				}
				return
			}
			metrics.telemetrySent.Add(1)
		}
	}
}

func queryLoop(ctx context.Context, engine string, headers map[string]string, worker int, failures *errorCounts, nearby, routes *samples, wg *sync.WaitGroup) {
	defer wg.Done()
	ticker := time.NewTicker(250 * time.Millisecond)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			if worker%2 == 0 {
				started := time.Now()
				if _, err := request(ctx, http.MethodGet, engine+"/spatial/devices/nearby?lat=13.0067&lon=80.2206&radius_meters=5000&limit=20", nil, headers); err != nil {
					failures.add("nearby", err.Category)
				} else {
					nearby.add(time.Since(started))
				}
			} else {
				started := time.Now()
				body := map[string]any{"mobility_profile": "ROAD_VEHICLE", "origin": map[string]float64{"latitude": 13.0067, "longitude": 80.2206}, "destination": map[string]float64{"latitude": 13.02, "longitude": 80.23}, "policy": "FASTEST"}
				if _, err := request(ctx, http.MethodPost, engine+"/routes", body, headers); err != nil {
					failures.add("route", err.Category)
				} else {
					routes.add(time.Since(started))
				}
			}
		}
	}
}

func taskLoop(ctx context.Context, engine, project string, headers map[string]string, jobs <-chan int, metrics *counters, failures *errorCounts, latency, candidateLatency, routingLatency, persistenceLatency *samples, wg *sync.WaitGroup) {
	defer wg.Done()
	types := []string{"NAVIGATE", "RELOCATE", "CAPTURE_IMAGE", "RUN_MODEL"}
	capability := map[string]string{"NAVIGATE": "navigate", "RELOCATE": "receive_relocation_command", "CAPTURE_IMAGE": "capture_image", "RUN_MODEL": "run_model"}
	for n := range jobs {
		select {
		case <-ctx.Done():
			return
		default:
		}
		kind := types[n%len(types)]
		requirements := map[string]any{"required_capabilities": []string{capability[kind]}, "project_id": project}
		target := map[string]any{"fixture": n}
		if kind != "RUN_MODEL" {
			requirements["minimum_battery"] = 20
			target["lat"], target["lon"], target["policy"] = 13.02, 80.23, "FASTEST"
		}
		if kind == "NAVIGATE" {
			requirements["planning_mode"] = "POLARIS_REQUIRED"
		} else if kind == "RELOCATE" {
			requirements["planning_mode"] = "DEVICE_LOCAL"
		}
		body := map[string]any{"project_id": project, "task_type": kind, "priority": "NORMAL", "requirements": requirements, "target": target, "expires_at": time.Now().Add(5 * time.Minute).UTC()}
		started := time.Now()
		metrics.taskRequests.Add(1)
		response, err := request(ctx, http.MethodPost, engine+"/tasks", body, headers)
		latency.add(time.Since(started))
		if err != nil {
			failures.add("task", err.Category)
			continue
		}
		metrics.tasksCreated.Add(1)
		var decoded struct {
			Data struct {
				Command any `json:"command"`
				Timing  struct {
					CandidateSelectionDurationUS int64 `json:"candidate_selection_duration_us"`
					RoutingDurationUS            int64 `json:"routing_duration_us"`
					PersistenceDurationUS        int64 `json:"persistence_duration_us"`
				} `json:"timing"`
			} `json:"data"`
		}
		if json.Unmarshal(response, &decoded) == nil {
			candidateLatency.add(time.Duration(decoded.Data.Timing.CandidateSelectionDurationUS) * time.Microsecond)
			routingLatency.add(time.Duration(decoded.Data.Timing.RoutingDurationUS) * time.Microsecond)
			persistenceLatency.add(time.Duration(decoded.Data.Timing.PersistenceDurationUS) * time.Microsecond)
			if decoded.Data.Command != nil {
				metrics.tasksAssigned.Add(1)
			}
		}
		time.Sleep(100 * time.Millisecond)
	}
}

type requestFailure struct {
	Category string
	Status   int
	Code     string
	Cause    error
}

func (e *requestFailure) Error() string {
	if e.Cause != nil {
		return e.Cause.Error()
	}
	return fmt.Sprintf("HTTP %d %s", e.Status, e.Code)
}

func request(ctx context.Context, method, url string, body any, headers map[string]string) ([]byte, *requestFailure) {
	var reader io.Reader
	if body != nil {
		encoded, _ := json.Marshal(body)
		reader = bytes.NewReader(encoded)
	}
	req, err := http.NewRequestWithContext(ctx, method, url, reader)
	if err != nil {
		return nil, &requestFailure{Category: "unexpected", Cause: err}
	}
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	for key, value := range headers {
		req.Header.Set(key, value)
	}
	response, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, &requestFailure{Category: classifyTransport(err), Cause: err}
	}
	defer response.Body.Close()
	payload, readErr := io.ReadAll(response.Body)
	if readErr != nil {
		return nil, &requestFailure{Category: "transport_error", Cause: readErr}
	}
	if response.StatusCode < 200 || response.StatusCode >= 300 {
		var decoded struct {
			Error struct {
				Code string `json:"code"`
			} `json:"error"`
		}
		_ = json.Unmarshal(payload, &decoded)
		return payload, &requestFailure{Category: classifyHTTP(response.StatusCode, decoded.Error.Code), Status: response.StatusCode, Code: decoded.Error.Code}
	}
	return payload, nil
}

func classifyTransport(err error) string {
	switch {
	case errors.Is(err, context.Canceled):
		return "cancelled"
	case errors.Is(err, context.DeadlineExceeded):
		return "timeout"
	default:
		return "transport_error"
	}
}

func classifyHTTP(status int, code string) string {
	upper := strings.ToUpper(code)
	switch {
	case upper == "ROUTING_BUSY":
		return "routing_busy"
	case upper == "ROUTING_TIMEOUT" || status == http.StatusRequestTimeout || status == http.StatusGatewayTimeout:
		return "timeout"
	case status == http.StatusConflict:
		return "conflict"
	case strings.Contains(upper, "NO_ROUTE") || strings.Contains(upper, "NO_ROAD_NODE") || strings.Contains(upper, "OUTSIDE_REGION"):
		return "no_route"
	case status >= 500:
		return "server_error"
	case status >= 400:
		return "client_error"
	default:
		return "unexpected"
	}
}

func countType(devices []deviceConfig, kind pb.NodeType) int {
	total := 0
	for _, d := range devices {
		if d.Spatial && d.NodeType == kind {
			total++
		}
	}
	return total
}
func countNonSpatial(devices []deviceConfig) int {
	total := 0
	for _, d := range devices {
		if !d.Spatial {
			total++
		}
	}
	return total
}
func must(err error) {
	if err != nil {
		panic(err)
	}
}
```

---

## cmd\system-soak\main_test.go

```
package main

import (
	"context"
	"net/http"
	"testing"
)

func TestFailureClassification(t *testing.T) {
	tests := []struct {
		status int
		code   string
		want   string
	}{
		{http.StatusTooManyRequests, "ROUTING_BUSY", "routing_busy"},
		{http.StatusGatewayTimeout, "ROUTING_TIMEOUT", "timeout"},
		{http.StatusConflict, "NO_ELIGIBLE_DEVICE", "conflict"},
		{http.StatusUnprocessableEntity, "NO_ROUTE", "no_route"},
		{http.StatusBadRequest, "INVALID_REQUEST", "client_error"},
		{http.StatusInternalServerError, "ORCHESTRATION_ERROR", "server_error"},
	}
	for _, test := range tests {
		if got := classifyHTTP(test.status, test.code); got != test.want {
			t.Fatalf("HTTP %d %s classified as %s, want %s", test.status, test.code, got, test.want)
		}
	}
	if got := classifyTransport(context.Canceled); got != "cancelled" {
		t.Fatalf("cancel classified as %s", got)
	}
	if got := classifyTransport(context.DeadlineExceeded); got != "timeout" {
		t.Fatalf("deadline classified as %s", got)
	}
}

func TestUnexpectedTotalsExcludeBoundedFailures(t *testing.T) {
	values := map[string]int64{"route.routing_busy": 5, "task.conflict": 2, "route.timeout": 1, "nearby.server_error": 0, "command.transport_error": 0, "task.unexpected": 0}
	totals := errorTotals(values)
	if totals["expected"] != 8 || totals["unexpected"] != 0 {
		t.Fatalf("unexpected totals: %#v", totals)
	}
	values["nearby.server_error"] = 1
	if errorTotals(values)["unexpected"] != 1 {
		t.Fatal("server error did not fail the unexpected-error gate")
	}
}
```

---

## internal\architecture\consistency_test.go

```
package architecture

import (
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
)

func TestProductionCodeHasSingleSpatialRoutingAndCommandAuthorities(t *testing.T) {
	backend, err := filepath.Abs(filepath.Join("..", ".."))
	if err != nil {
		t.Fatal(err)
	}
	forbiddenImports := []string{"/algo_/quadtree", "/algo_/graph", "/internal/core/actor"}
	err = filepath.WalkDir(backend, func(path string, entry os.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if entry.IsDir() {
			if entry.Name() == ".phase41-go-cache" || entry.Name() == "node_modules" {
				return filepath.SkipDir
			}
			return nil
		}
		if !strings.HasSuffix(path, ".go") || strings.HasSuffix(path, "_test.go") || strings.Contains(path, string(filepath.Separator)+"api"+string(filepath.Separator)+"proto") {
			return nil
		}
		parsed, parseErr := parser.ParseFile(token.NewFileSet(), path, nil, parser.ImportsOnly)
		if parseErr != nil {
			return parseErr
		}
		for _, imported := range parsed.Imports {
			name, _ := strconv.Unquote(imported.Path.Value)
			for _, forbidden := range forbiddenImports {
				if strings.HasSuffix(name, forbidden) {
					t.Errorf("%s imports retired authority %s", path, name)
				}
			}
		}
		content, readErr := os.ReadFile(path)
		if readErr != nil {
			return readErr
		}
		if strings.Contains(string(content), "telemetry:commands") || strings.Contains(string(content), "StartAutonomousLoop") {
			t.Errorf("%s retains a retired direct command path", path)
		}
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
}

func TestNodeTypeIsNotUsedAsBitmask(t *testing.T) {
	backend, _ := filepath.Abs(filepath.Join("..", ".."))
	_ = filepath.WalkDir(backend, func(path string, entry os.DirEntry, walkErr error) error {
		if walkErr != nil || entry.IsDir() || !strings.HasSuffix(path, ".go") || strings.HasSuffix(path, "_test.go") {
			return walkErr
		}
		file, err := parser.ParseFile(token.NewFileSet(), path, nil, 0)
		if err != nil {
			return err
		}
		ast.Inspect(file, func(node ast.Node) bool {
			binary, ok := node.(*ast.BinaryExpr)
			if ok && (binary.Op == token.AND || binary.Op == token.OR) {
				text := expressionName(binary.X) + expressionName(binary.Y)
				if strings.Contains(strings.ToLower(text), "type") || strings.Contains(strings.ToLower(text), "class") {
					t.Errorf("%s contains a type/class bitmask expression", path)
				}
			}
			return true
		})
		return nil
	})
}

func expressionName(expr ast.Expr) string {
	switch value := expr.(type) {
	case *ast.Ident:
		return value.Name
	case *ast.SelectorExpr:
		return expressionName(value.X) + value.Sel.Name
	}
	return ""
}
```

---

## internal\config\config.go

```
package config

import (
	"log/slog"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	App    AppConfig
	Server ServerConfig
	Redis  RedisConfig
	DB     DBConfig
}

type AppConfig struct {
	Env      string
	LogLevel string
}

type ServerConfig struct {
	GatewayPort string
	EnginePort  string
}

type RedisConfig struct {
	URL string
}

type DBConfig struct {
	URL string
}

func Load() *Config {
	// Don't crash if .env is missing (e.g., in Docker). Just log it.
	if err := godotenv.Load(); err != nil {
		slog.Debug("No .env file found. Relying on OS environment variables.")
	}

	return &Config{
		App: AppConfig{
			Env:      getEnv("APP_ENV", "development"),
			LogLevel: getEnv("LOG_LEVEL", "info"),
		},
		Server: ServerConfig{
			GatewayPort: getEnv("GATEWAY_PORT", "6080"),
			EnginePort:  getEnv("ENGINE_PORT", "6081"),
		},
		Redis: RedisConfig{
			URL: getEnv("REDIS_URL", "redis://localhost:6379/0"),
		},
		DB: DBConfig{
			URL: getEnv("POSTGRES_URL", "postgres://polaris_user:polaris_password@localhost:5432/polaris_core?sslmode=disable"),
		},
	}
}

// getEnv returns the environment variable or a safe fallback
func getEnv(key, fallback string) string {
	if val, exists := os.LookupEnv(key); exists && val != "" {
		return val
	}
	return fallback
}
```

---

## internal\adapter\handler\dashboard.go

```
package handler

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"strings"
	"sync"

	"github.com/Akashpg-M/polaris/backend/internal/core/auth"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

// DashboardRegistry tracks all active UI dashboard connections
type DashboardRegistry struct {
	mu          sync.RWMutex
	connections map[*websocket.Conn]string
}

// NewDashboardRegistry initializes an empty thread-safe connection tracker
func NewDashboardRegistry() *DashboardRegistry {
	return &DashboardRegistry{
		connections: make(map[*websocket.Conn]string),
	}
}

// Register adds a new UI dashboard connection to the active broadcast list
func (r *DashboardRegistry) Register(conn *websocket.Conn, tenantID string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.connections[conn] = tenantID
	slog.Info("[DashboardRegistry] New web client connected to telemetry stream", "active_dashboards", len(r.connections))
}

// Unregister safely drops a connection when a user closes the browser tab
func (r *DashboardRegistry) Unregister(conn *websocket.Conn) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, exists := r.connections[conn]; exists {
		delete(r.connections, conn)
		conn.Close()
		slog.Info("[DashboardRegistry] Web client disconnected", "active_dashboards", len(r.connections))
	}
}

// BroadcastToUIs pumps a raw message string out to every single open dashboard browser concurrently
func (r *DashboardRegistry) BroadcastToUIs(payload string) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var event struct {
		TenantID string `json:"tenant_id"`
	}
	if json.Unmarshal([]byte(payload), &event) != nil || event.TenantID == "" {
		return
	}
	for conn, tenantID := range r.connections {
		if tenantID != event.TenantID {
			continue
		}
		// Send standard Text message to UIs (JSON strings)
		err := conn.WriteMessage(websocket.TextMessage, []byte(payload))
		if err != nil {
			slog.Warn("[DashboardRegistry] Failed to push frame down streaming channel, breaking pipe", "err", err)
			// Schedule cleanup asynchronously to prevent deadlocking the write lock
			go r.Unregister(conn)
		}
	}
}

// DashboardHandler provides the REST-to-WS upgrade entrypoint for web clients
type DashboardHandler struct {
	registry      *DashboardRegistry
	upgrader      websocket.Upgrader
	authenticator DashboardAuthenticator
}
type DashboardAuthenticator interface {
	ResolveOperator(context.Context, string) (auth.OperatorPrincipal, error)
	ConsumeOperatorTicket(context.Context, string) (auth.OperatorPrincipal, error)
}

// NewDashboardHandler constructs the gateway handler for web clients
func NewDashboardHandler(registry *DashboardRegistry, authenticator DashboardAuthenticator) *DashboardHandler {
	return &DashboardHandler{
		registry:      registry,
		authenticator: authenticator,
		upgrader: websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
			// Allow cross-origin requests so your local frontend can connect easily
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		},
	}
}

// HandleWebConnection converts incoming HTTP requests into an asynchronous JSON stream
func (h *DashboardHandler) HandleWebConnection(c *gin.Context) {
	principal, err := h.authenticate(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": gin.H{"code": "UNAUTHENTICATED", "message": "Dashboard authorization is required"}})
		return
	}
	tenantID := principal.TenantID
	if tenantID == "" {
		tenantID = c.GetHeader("X-Tenant-ID")
		if tenantID == "" {
			tenantID = c.Query("tenant_id")
		}
	}
	if tenantID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": gin.H{"code": "TENANT_REQUIRED", "message": "Tenant scope is required"}})
		return
	}
	conn, err := h.upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		slog.Error("[DashboardGateway] Handshake Upgrade Error", "error", err)
		return
	}

	h.registry.Register(conn, tenantID)

	// Keep connection alive, listen for client-side closures
	go func() {
		defer h.registry.Unregister(conn)
		for {
			// Dashboards are consumer-only; if they send messages or close, clean up the pipe
			if _, _, err := conn.ReadMessage(); err != nil {
				break
			}
		}
	}()
}
func (h *DashboardHandler) authenticate(c *gin.Context) (auth.OperatorPrincipal, error) {
	if ticket := c.Query("ticket"); ticket != "" {
		return h.authenticator.ConsumeOperatorTicket(c.Request.Context(), ticket)
	}
	header := c.GetHeader("Authorization")
	if !strings.HasPrefix(header, "Bearer ") {
		return auth.OperatorPrincipal{}, auth.ErrInvalidCredential
	}
	return h.authenticator.ResolveOperator(c.Request.Context(), strings.TrimSpace(strings.TrimPrefix(header, "Bearer ")))
}
```

---

## internal\adapter\handler\device_connections.go

```
package handler

import (
	"context"
	"encoding/json"
	"errors"
	"hash/fnv"
	"log/slog"
	"sync"
	"time"

	"github.com/Akashpg-M/polaris/backend/internal/adapter/repository"
	"github.com/Akashpg-M/polaris/backend/internal/application/orchestration"
	"github.com/Akashpg-M/polaris/backend/internal/core/auth"
	"github.com/Akashpg-M/polaris/backend/internal/core/command"
	"github.com/gorilla/websocket"
	"github.com/redis/go-redis/v9"
)

type deviceSession struct {
	principal auth.DevicePrincipal
	ownership repository.ConnectionOwnership
	conn      *websocket.Conn
	writeMu   sync.Mutex
	cancel    context.CancelFunc
}

func (s *deviceSession) writeJSON(value interface{}) error {
	s.writeMu.Lock()
	defer s.writeMu.Unlock()
	return s.conn.WriteJSON(value)
}

func (s *deviceSession) close(code int, reason string) {
	s.writeMu.Lock()
	defer s.writeMu.Unlock()
	_ = s.conn.WriteControl(websocket.CloseMessage, websocket.FormatCloseMessage(code, reason), time.Now().Add(time.Second))
	_ = s.conn.Close()
}

type DeviceConnectionManager struct {
	mu        sync.RWMutex
	sessions  map[string]*deviceSession
	gatewayID string
	lease     time.Duration
	owners    *repository.ConnectionOwnershipStore
	store     *repository.RegistryStore
	redis     *redis.Client
	metrics   *orchestration.Metrics
}

func NewDeviceConnectionManager(gatewayID string, lease time.Duration, redisClient *redis.Client, store *repository.RegistryStore, metrics *orchestration.Metrics) *DeviceConnectionManager {
	if gatewayID == "" {
		gatewayID = "gateway-1"
	}
	return &DeviceConnectionManager{sessions: map[string]*deviceSession{}, gatewayID: gatewayID, lease: lease, owners: repository.NewConnectionOwnershipStore(redisClient, lease), store: store, redis: redisClient, metrics: metrics}
}

func deviceKey(tenant, device string) string { return tenant + ":" + device }

func (m *DeviceConnectionManager) Register(ctx context.Context, conn *websocket.Conn, principal auth.DevicePrincipal) (*deviceSession, error) {
	connectionID := auth.NewID()
	ownership, err := m.owners.Claim(ctx, principal.TenantID, principal.DeviceID, m.gatewayID, connectionID, principal.CredentialID)
	if err != nil {
		return nil, err
	}
	sessionCtx, cancel := context.WithCancel(context.Background())
	session := &deviceSession{principal: principal, ownership: ownership, conn: conn, cancel: cancel}
	key := deviceKey(principal.TenantID, principal.DeviceID)
	m.mu.Lock()
	previous := m.sessions[key]
	m.sessions[key] = session
	m.mu.Unlock()
	if previous != nil {
		previous.cancel()
		previous.close(websocket.CloseServiceRestart, "connection superseded by a newer ownership epoch")
	}
	m.metrics.ActiveConnections.Add(1)
	_ = m.store.Audit(ctx, principal.TenantID, m.gatewayID, "DEVICE_OWNERSHIP_CLAIMED", "device", principal.DeviceID, "", "SUCCESS")
	go m.heartbeat(sessionCtx, session)
	go m.ReconcileDevice(sessionCtx, principal.TenantID, principal.DeviceID)
	return session, nil
}

func (m *DeviceConnectionManager) Unregister(ctx context.Context, session *deviceSession) {
	if session == nil {
		return
	}
	session.cancel()
	key := deviceKey(session.principal.TenantID, session.principal.DeviceID)
	m.mu.Lock()
	if m.sessions[key] == session {
		delete(m.sessions, key)
		m.metrics.ActiveConnections.Add(-1)
	}
	m.mu.Unlock()
	released, _ := m.owners.Release(ctx, session.ownership)
	if released {
		_ = m.store.Audit(ctx, session.principal.TenantID, m.gatewayID, "DEVICE_OWNERSHIP_LOST", "device", session.principal.DeviceID, "", "SUCCESS")
	}
}

func (m *DeviceConnectionManager) heartbeat(ctx context.Context, session *deviceSession) {
	interval := m.lease / 3
	if interval < time.Second {
		interval = time.Second
	}
	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	var refreshFailureStarted time.Time
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			ok, err := m.owners.Refresh(ctx, session.ownership)
			if err != nil {
				if refreshFailureStarted.IsZero() {
					refreshFailureStarted = time.Now()
				}
				// A single Redis timeout is not proof that fencing ownership was
				// lost. Retry within the lease window; fail closed if Redis cannot
				// confirm ownership for a full TTL.
				if time.Since(refreshFailureStarted) < m.lease-interval {
					slog.Warn("gateway ownership refresh transient failure", "tenant_id", session.principal.TenantID, "device_id", session.principal.DeviceID, "error", err)
					continue
				}
				slog.Warn("gateway ownership lease expired during refresh failure", "tenant_id", session.principal.TenantID, "device_id", session.principal.DeviceID, "error", err)
				session.close(websocket.ClosePolicyViolation, "gateway ownership lease lost")
				return
			}
			refreshFailureStarted = time.Time{}
			if !ok {
				slog.Warn("gateway ownership fencing mismatch", "tenant_id", session.principal.TenantID, "device_id", session.principal.DeviceID)
				session.close(websocket.ClosePolicyViolation, "gateway ownership lease lost")
				return
			}
		}
	}
}

func (m *DeviceConnectionManager) session(tenant, device string) *deviceSession {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.sessions[deviceKey(tenant, device)]
}

func (m *DeviceConnectionManager) Deliver(ctx context.Context, envelope command.Envelope) error {
	session := m.session(envelope.TenantID, envelope.DeviceID)
	if session == nil {
		return repository.ErrNotFound
	}
	if err := m.store.RevalidateDevice(ctx, session.principal); err != nil {
		if errors.Is(err, auth.ErrInvalidCredential) {
			session.close(websocket.ClosePolicyViolation, "credential, tenant, or device is inactive")
			return repository.ErrForbidden
		}
		return err
	}
	ok, err := m.owners.Owns(ctx, session.ownership)
	if err != nil || !ok {
		return repository.ErrForbidden
	}
	record, err := m.store.PrepareDelivery(ctx, envelope.TenantID, envelope.DeviceID, envelope.CommandID, m.gatewayID, session.ownership.Epoch)
	if err != nil {
		return err
	}
	// Fence again immediately before the volatile write. A reconnect may have
	// advanced the epoch while the PostgreSQL delivery transition was running.
	ok, err = m.owners.Owns(ctx, session.ownership)
	if err != nil || !ok {
		return repository.ErrForbidden
	}
	outgoing := record.Envelope()
	outgoing.DeliveryObservation = envelope.DeliveryObservation
	if err = session.writeJSON(outgoing); err != nil {
		return err
	}
	m.metrics.CommandsDelivered.Add(1)
	return nil
}

func (m *DeviceConnectionManager) ReconcileDevice(ctx context.Context, tenant, device string) {
	commands, err := m.store.PendingCommandsForDevice(ctx, tenant, device)
	if err != nil {
		slog.Warn("pending command reconciliation failed", "tenant_id", tenant, "device_id", device, "error", err)
		return
	}
	for _, record := range commands {
		if err = m.Deliver(ctx, record.Envelope()); err != nil {
			if err != repository.ErrInvalidTransition && err != repository.ErrConflict {
				slog.Warn("reconciled command delivery deferred", "command_id", record.CommandID, "error", err)
			}
			return
		}
	}
}

func (m *DeviceConnectionManager) StartSubscriber(ctx context.Context) {
	pubsub := m.redis.Subscribe(ctx, repository.GatewayCommandChannel(m.gatewayID))
	defer pubsub.Close()
	const workers = 16
	queues := make([]chan command.Envelope, workers)
	var workerWG sync.WaitGroup
	for index := range queues {
		queues[index] = make(chan command.Envelope, 64)
		workerWG.Add(1)
		go func(queue <-chan command.Envelope) {
			defer workerWG.Done()
			for envelope := range queue {
				if err := m.Deliver(ctx, envelope); err != nil && err != repository.ErrNotFound && err != repository.ErrInvalidTransition && err != repository.ErrConflict {
					slog.Warn("live command notification could not be delivered", "command_id", envelope.CommandID, "error", err)
				}
			}
		}(queues[index])
	}
	defer func() {
		for _, queue := range queues {
			close(queue)
		}
		workerWG.Wait()
	}()
	for message := range pubsub.Channel() {
		var envelope command.Envelope
		if json.Unmarshal([]byte(message.Payload), &envelope) != nil || envelope.FrameType != "COMMAND" {
			continue
		}
		if envelope.DeliveryObservation != nil {
			envelope.DeliveryObservation.GatewayReceivedAt = time.Now().UTC()
		}
		h := fnv.New32a()
		_, _ = h.Write([]byte(deviceKey(envelope.TenantID, envelope.DeviceID)))
		select {
		case queues[int(h.Sum32())%len(queues)] <- envelope:
		case <-ctx.Done():
			return
		}
	}
}

func (m *DeviceConnectionManager) OwnershipStore() *repository.ConnectionOwnershipStore {
	return m.owners
}
```

---

## internal\adapter\handler\ingestion.go

```
package handler

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"strings"
	"sync/atomic"
	"time"

	pb "github.com/Akashpg-M/polaris/backend/api/proto/v1"
	"github.com/Akashpg-M/polaris/backend/internal/core/auth"
	"github.com/Akashpg-M/polaris/backend/internal/core/command"
	"github.com/Akashpg-M/polaris/backend/internal/core/events"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"google.golang.org/protobuf/proto"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

type IngestionHandler struct {
	publisher     TelemetryEventPublisher
	activeSockets int64
	authenticator DeviceAuthenticator
	connections   *DeviceConnectionManager
}

// TelemetryEventPublisher is the durable ingress boundary. The gateway does
// not own a second per-device state machine; accepted frames are synchronously
// appended to Kafka and backpressure is propagated to the socket.
type TelemetryEventPublisher interface {
	PublishEvent(context.Context, string, interface{}) error
}

type DeviceAuthenticator interface {
	ResolveDevice(context.Context, string) (auth.DevicePrincipal, error)
	ConsumeTicket(context.Context, string) (auth.DevicePrincipal, error)
	RevalidateDevice(context.Context, auth.DevicePrincipal) error
}

func NewIngestionHandler(publisher TelemetryEventPublisher, authenticator DeviceAuthenticator, connections *DeviceConnectionManager) *IngestionHandler {
	return &IngestionHandler{
		publisher:     publisher,
		activeSockets: 0, // High-performance, lock-free counter
		authenticator: authenticator,
		connections:   connections,
	}
}

func (h *IngestionHandler) GetActiveConnectionsCount() int64 {
	return atomic.LoadInt64(&h.activeSockets)
}

func (h *IngestionHandler) HandleIoTConnection(c *gin.Context) {
	principal, err := h.authenticate(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": gin.H{"code": "DEVICE_AUTHENTICATION_FAILED", "message": "Device credential is invalid or inactive"}})
		return
	}
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		slog.Error("[Gateway] WebSocket upgrade failed", "error", err)
		return
	}
	defer conn.Close()
	conn.SetReadLimit(events.MaxFrameBytes)
	session, err := h.connections.Register(c.Request.Context(), conn, principal)
	if err != nil {
		rejectFrame(conn, "connection ownership could not be established")
		return
	}
	defer h.connections.Unregister(context.Background(), session)

	atomic.AddInt64(&h.activeSockets, 1)
	defer atomic.AddInt64(&h.activeSockets, -1)

	nodeID := principal.DeviceID
	var bootID string
	lastRevalidationSuccess := time.Now()

	for {
		// 1. Read Raw Binary format matching pure Protobuf specs
		msgType, data, err := conn.ReadMessage()
		if err != nil {
			break
		}

		if err := h.authenticator.RevalidateDevice(c.Request.Context(), principal); err != nil {
			if errors.Is(err, auth.ErrInvalidCredential) {
				session.close(websocket.ClosePolicyViolation, "credential or device was revoked")
				return
			}
			// A registry transport failure is not evidence of revocation. Keep an
			// already-authenticated session for a bounded grace period while
			// readiness removes this gateway from new connection traffic.
			if time.Since(lastRevalidationSuccess) >= 30*time.Second {
				slog.Warn("device revalidation grace expired", "tenant_id", principal.TenantID, "device_id", principal.DeviceID, "error", err)
				session.close(websocket.CloseTryAgainLater, "device registry temporarily unavailable")
				return
			}
			slog.Warn("device revalidation transient failure", "tenant_id", principal.TenantID, "device_id", principal.DeviceID, "error", err)
		} else {
			lastRevalidationSuccess = time.Now()
		}
		if msgType == websocket.TextMessage {
			if err := h.handleDeviceControl(c.Request.Context(), principal, data); err != nil {
				session.close(websocket.ClosePolicyViolation, "command response rejected")
				return
			}
			continue
		}
		if msgType != websocket.BinaryMessage {
			slog.Warn("[Gateway] Security violation: unsupported WebSocket frame dropped.")
			continue
		}

		// 2. Fast Protobuf Unmarshaling (Executed on the edge worker thread)
		var payload pb.SpatialObject
		if err := proto.Unmarshal(data, &payload); err != nil {
			rejectFrame(conn, "malformed protobuf frame")
			return
		}
		if err := events.ValidateFrame(&payload, time.Now()); err != nil {
			slog.Warn("[Gateway] Rejected invalid telemetry before Kafka", "error", err)
			rejectFrame(conn, err.Error())
			return
		}
		if payload.Id != principal.DeviceID || payload.TenantId != principal.TenantID {
			rejectFrame(conn, "payload identity does not match authenticated device")
			return
		}
		if !telemetryTypeAllowed(principal.DeviceType, payload.Type) {
			rejectFrame(conn, "telemetry type does not match registered device type")
			return
		}
		// The principal, not the untrusted frame, owns platform identity.
		payload.Id = principal.DeviceID
		payload.TenantId = principal.TenantID

		// Initial connection mapping handshake check
		if bootID == "" {
			bootID = payload.DeviceBootId
			slog.Info("[Gateway] Device mapped to local gateway workspace", "node_id", nodeID)
		} else if payload.DeviceBootId != bootID {
			slog.Warn("[Gateway] Rejected device identity change within connection", "expected_device", nodeID, "actual_device", payload.Id)
			rejectFrame(conn, "device identity changed within connection")
			return
		}

		ingestedAt := time.Now().UTC()
		if payload.Timestamp == 0 {
			payload.Timestamp = payload.ObservedAt
		}
		envelope := events.NewTelemetryEnvelope(&payload, ingestedAt,
			c.GetHeader("X-Correlation-ID"), c.GetHeader("X-Causation-ID"), c.GetHeader("traceparent"))

		// Kafka is the sole durable telemetry boundary. A publication failure is
		// surfaced by closing the connection so the device can reconnect/replay.
		if err := h.publisher.PublishEvent(c.Request.Context(), "telemetry.ingress", envelope); err != nil {
			slog.Error("[Gateway] Kafka telemetry publication failed", "node_id", nodeID, "error", err)
			session.close(websocket.CloseTryAgainLater, "telemetry publication unavailable")
			return
		}
	}

	// Clean up runtime structures safely if the persistent socket drops
	if nodeID != "" {
		slog.Info("[Gateway] Telemetry channel closed at edge boundary", "node_id", nodeID)
	}
}

func telemetryTypeAllowed(deviceType string, nodeType pb.NodeType) bool {
	switch deviceType {
	case "delivery_drone":
		return nodeType == pb.NodeType_NODE_TYPE_DRONE
	case "ground_robot":
		return nodeType == pb.NodeType_NODE_TYPE_ROBOT
	case "connected_vehicle":
		return nodeType == pb.NodeType_NODE_TYPE_BIKE || nodeType == pb.NodeType_NODE_TYPE_AUTO || nodeType == pb.NodeType_NODE_TYPE_SEDAN || nodeType == pb.NodeType_NODE_TYPE_SUV
	case "fixed_iot_sensor", "static_camera":
		return nodeType == pb.NodeType_NODE_TYPE_STATIC_SENSOR
	default:
		// Non-spatial device types must use a future type-specific telemetry
		// schema instead of impersonating a SpatialObject profile.
		return false
	}
}

func (h *IngestionHandler) handleDeviceControl(ctx context.Context, principal auth.DevicePrincipal, data []byte) error {
	var header struct {
		FrameType string `json:"frame_type"`
	}
	if len(data) > 64*1024 || json.Unmarshal(data, &header) != nil {
		return errors.New("invalid device control envelope")
	}
	switch header.FrameType {
	case "COMMAND_ACK":
		var ack command.Ack
		if json.Unmarshal(data, &ack) != nil || ack.CommandID == "" || ack.SequenceNumber < 1 {
			return errors.New("invalid command acknowledgement")
		}
		if ack.Status != "ACCEPTED" && ack.Status != "REJECTED" && ack.Status != "DUPLICATE" && ack.Status != "EXPIRED" && ack.Status != "UNSUPPORTED" {
			return errors.New("unsupported acknowledgement status")
		}
		if err := h.connections.store.ApplyCommandAck(ctx, principal.TenantID, principal.DeviceID, ack); err != nil {
			return err
		}
		h.connections.metrics.CommandsAcknowledged.Add(1)
		return nil
	case "COMMAND_RESULT":
		var result command.Result
		if json.Unmarshal(data, &result) != nil || result.CommandID == "" || result.SequenceNumber < 1 || (result.Status != "SUCCEEDED" && result.Status != "COMPLETED" && result.Status != "FAILED") {
			return errors.New("invalid command result")
		}
		if err := h.connections.store.ApplyCommandResult(ctx, principal.TenantID, principal.DeviceID, result); err != nil {
			return err
		}
		if result.Status == "SUCCEEDED" || result.Status == "COMPLETED" {
			h.connections.metrics.TasksCompleted.Add(1)
		} else {
			h.connections.metrics.CommandsFailed.Add(1)
			h.connections.metrics.TasksFailed.Add(1)
		}
		return nil
	default:
		return errors.New("unsupported device frame type")
	}
}

func (h *IngestionHandler) authenticate(c *gin.Context) (auth.DevicePrincipal, error) {
	if h.authenticator == nil {
		return auth.DevicePrincipal{}, auth.ErrInvalidCredential
	}
	if ticket := c.Query("ticket"); ticket != "" {
		return h.authenticator.ConsumeTicket(c.Request.Context(), ticket)
	}
	header := c.GetHeader("Authorization")
	if !strings.HasPrefix(header, "Bearer ") {
		return auth.DevicePrincipal{}, auth.ErrInvalidCredential
	}
	return h.authenticator.ResolveDevice(c.Request.Context(), strings.TrimSpace(strings.TrimPrefix(header, "Bearer ")))
}

func rejectFrame(conn *websocket.Conn, reason string) {
	message := websocket.FormatCloseMessage(websocket.ClosePolicyViolation, fmt.Sprintf("telemetry rejected: %s", reason))
	_ = conn.WriteControl(websocket.CloseMessage, message, time.Now().Add(time.Second))
}
```

---

## internal\adapter\handler\ingestion_test.go

```
package handler

import (
	"testing"

	pb "github.com/Akashpg-M/polaris/backend/api/proto/v1"
)

func TestRegisteredDeviceTypeConstrainsTelemetryProfile(t *testing.T) {
	tests := []struct {
		deviceType string
		nodeType   pb.NodeType
		allowed    bool
	}{
		{"delivery_drone", pb.NodeType_NODE_TYPE_DRONE, true},
		{"delivery_drone", pb.NodeType_NODE_TYPE_SEDAN, false},
		{"connected_vehicle", pb.NodeType_NODE_TYPE_BIKE, true},
		{"connected_vehicle", pb.NodeType_NODE_TYPE_SUV, true},
		{"static_camera", pb.NodeType_NODE_TYPE_STATIC_SENSOR, true},
		{"compute_node", pb.NodeType_NODE_TYPE_STATIC_SENSOR, false},
		{"unknown", pb.NodeType_NODE_TYPE_DRONE, false},
	}
	for _, test := range tests {
		if got := telemetryTypeAllowed(test.deviceType, test.nodeType); got != test.allowed {
			t.Errorf("%s/%s allowed=%v want %v", test.deviceType, test.nodeType, got, test.allowed)
		}
	}
}
```

---

## internal\adapter\handler\match.go

```
package handler

import (
	"net/http"
	"strconv"

	pb "github.com/Akashpg-M/polaris/backend/api/proto/v1"
	"github.com/Akashpg-M/polaris/backend/internal/application/spatial"
	"github.com/gin-gonic/gin"
)

type MatchHandler struct {
	engine *spatial.Engine
}

func NewMatchHandler(engine *spatial.Engine) *MatchHandler {
	return &MatchHandler{engine: engine}
}

// GetNearestNodes handles GET /api/v1/nodes/match
func (h *MatchHandler) GetNearestNodes(c *gin.Context) {
	tenantID, ok := tenantFor(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "missing tenant identity"})
		return
	}

	lat, errLat := strconv.ParseFloat(c.Query("lat"), 64)
	lon, errLon := strconv.ParseFloat(c.Query("lon"), 64)
	radius, errRad := strconv.ParseFloat(c.DefaultQuery("radius_km", "10.0"), 64)
	assetClass, errClass := strconv.ParseInt(c.DefaultQuery("class", "2"), 10, 32)

	if errLat != nil || errLon != nil || errRad != nil || errClass != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid query parameters. lat and lon are required."})
		return
	}

	matches := h.engine.FindNearest(tenantID, lat, lon, radius, pb.NodeType(assetClass))
	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"count":  len(matches),
		"data":   matches,
	})
}
```

---

## internal\adapter\handler\mobility_api.go

```
package handler

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/Akashpg-M/polaris/backend/internal/modules/mobility/model"
	"github.com/Akashpg-M/polaris/backend/internal/modules/mobility/routing"
	"github.com/Akashpg-M/polaris/backend/internal/modules/mobility/spatial"
	"github.com/gin-gonic/gin"
)

type MobilityAPI struct {
	spatial  *spatial.Manager
	routing  routing.RoutingEngine
	maxLimit int
}

func NewMobilityAPI(s *spatial.Manager, r routing.RoutingEngine, maxLimit int) *MobilityAPI {
	return &MobilityAPI{spatial: s, routing: r, maxLimit: maxLimit}
}
func (a *MobilityAPI) Register(r *gin.RouterGroup) {
	r.GET("/spatial/devices/nearby", a.nearby)
	r.POST("/routes", a.route)
	r.GET("/routes/calculate", a.legacyRoute)
}

func (a *MobilityAPI) legacyRoute(c *gin.Context) {
	tenant, ok := tenantFor(c)
	if !ok {
		apiError(c, 400, "TENANT_REQUIRED", "Tenant scope is required")
		return
	}
	srcLat, e1 := strconv.ParseFloat(c.Query("src_lat"), 64)
	srcLon, e2 := strconv.ParseFloat(c.Query("src_lon"), 64)
	dstLat, e3 := strconv.ParseFloat(c.Query("tgt_lat"), 64)
	dstLon, e4 := strconv.ParseFloat(c.Query("tgt_lon"), 64)
	if e1 != nil || e2 != nil || e3 != nil || e4 != nil {
		apiError(c, 400, "INVALID_ROUTE_REQUEST", "Valid source and target coordinates are required")
		return
	}
	result, err := a.routing.Route(c, routing.RouteRequest{TenantID: tenant, MobilityProfile: model.MobilityRoadVehicle, Origin: model.Position{Latitude: srcLat, Longitude: srcLon}, Destination: model.Position{Latitude: dstLat, Longitude: dstLon}, Policy: routing.RouteFastest})
	if err != nil {
		routingError(c, err)
		return
	}
	route := make([]gin.H, len(result.Waypoints))
	for i, p := range result.Waypoints {
		route[i] = gin.H{"lat": p.Latitude, "lon": p.Longitude}
	}
	c.JSON(200, gin.H{"status": "success", "total_dist_km": result.DistanceMeters / 1000, "estimated_duration_ms": result.EstimatedTime.Milliseconds(), "road_graph_version": result.GraphVersion, "routing_snapshot_version": result.SnapshotVersion, "route": route})
}
func (a *MobilityAPI) nearby(c *gin.Context) {
	tenant, ok := tenantFor(c)
	if !ok {
		apiError(c, 400, "TENANT_REQUIRED", "Tenant scope is required")
		return
	}
	lat, e1 := strconv.ParseFloat(c.Query("lat"), 64)
	lon, e2 := strconv.ParseFloat(c.Query("lon"), 64)
	radius, e3 := strconv.ParseFloat(c.DefaultQuery("radius_meters", "5000"), 64)
	limit, e4 := strconv.Atoi(c.DefaultQuery("limit", "20"))
	if e1 != nil || e2 != nil || e3 != nil || e4 != nil || limit < 1 || limit > a.maxLimit {
		apiError(c, 400, "INVALID_SPATIAL_QUERY", "Coordinates, radius, or limit are invalid")
		return
	}
	items, err := a.spatial.Nearby(tenant, model.Position{Latitude: lat, Longitude: lon}, radius, limit)
	if err != nil {
		apiError(c, 400, "INVALID_SPATIAL_QUERY", err.Error())
		return
	}
	apiData(c, 200, gin.H{"count": len(items), "devices": items})
}
func (a *MobilityAPI) route(c *gin.Context) {
	tenant, ok := tenantFor(c)
	if !ok {
		apiError(c, 400, "TENANT_REQUIRED", "Tenant scope is required")
		return
	}
	var req routing.RouteRequest
	if c.ShouldBindJSON(&req) != nil {
		apiError(c, 400, "INVALID_ROUTE_REQUEST", "A valid route request is required")
		return
	}
	req.TenantID = tenant
	if req.MobilityProfile == "" {
		req.MobilityProfile = model.MobilityRoadVehicle
	}
	if req.Policy == "" {
		req.Policy = routing.RouteFastest
	}
	result, err := a.routing.Route(c, req)
	if err != nil {
		routingError(c, err)
		return
	}
	apiData(c, 200, result)
}
func routingError(c *gin.Context, err error) {
	code := http.StatusUnprocessableEntity
	name := err.Error()
	switch {
	case errors.Is(err, routing.ErrBusy):
		code = http.StatusTooManyRequests
	case errors.Is(err, routing.ErrTimeout):
		code = http.StatusGatewayTimeout
	case errors.Is(err, routing.ErrUnavailable):
		code = http.StatusServiceUnavailable
	case errors.Is(err, routing.ErrNoRoute), errors.Is(err, routing.ErrNoRoadNode), errors.Is(err, routing.ErrOutsideRegion), errors.Is(err, routing.ErrUnsupportedProfile):
		code = http.StatusUnprocessableEntity
	default:
		code = http.StatusInternalServerError
		name = "ROUTING_ERROR"
	}
	apiError(c, code, name, err.Error())
}
```

---

## internal\adapter\handler\orchestration_api.go

```
package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/Akashpg-M/polaris/backend/internal/adapter/repository"
	"github.com/Akashpg-M/polaris/backend/internal/application/orchestration"
	"github.com/Akashpg-M/polaris/backend/internal/core/command"
	"github.com/Akashpg-M/polaris/backend/internal/core/extension"
	taskcore "github.com/Akashpg-M/polaris/backend/internal/core/task"
	"github.com/Akashpg-M/polaris/backend/internal/modules/mobility/routing"
	"github.com/gin-gonic/gin"
)

type OrchestrationAPI struct {
	store   *repository.RegistryStore
	service *orchestration.Service
	metrics *orchestration.Metrics
}

func NewOrchestrationAPI(store *repository.RegistryStore, service *orchestration.Service, metrics *orchestration.Metrics) *OrchestrationAPI {
	return &OrchestrationAPI{store: store, service: service, metrics: metrics}
}

func (a *OrchestrationAPI) Register(r *gin.RouterGroup, registryAPI *RegistryAPI) {
	r.POST("/tasks", registryAPI.Middleware("orchestrate"), a.createTask)
	r.GET("/tasks", registryAPI.Middleware("read"), a.listTasks)
	r.GET("/tasks/:task_id", registryAPI.Middleware("read"), a.getTask)
	r.POST("/tasks/:task_id/cancel", registryAPI.Middleware("orchestrate"), a.cancelTask)
	r.POST("/tasks/:task_id/retry", registryAPI.Middleware("admin_retry"), a.retryTask)
	r.GET("/commands", registryAPI.Middleware("read"), a.listCommands)
	r.GET("/commands/:command_id", registryAPI.Middleware("read"), a.getCommand)
	r.POST("/commands/:command_id/retry", registryAPI.Middleware("admin_retry"), a.retryCommand)
	r.POST("/commands/:command_id/cancel", registryAPI.Middleware("orchestrate"), a.cancelCommand)
	r.GET("/metrics/orchestration", registryAPI.Middleware("read"), func(c *gin.Context) { apiData(c, 200, a.metrics.Snapshot()) })
}

func orchestrationError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, orchestration.ErrInvalidTask), errors.Is(err, orchestration.ErrUnsupportedCommand):
		apiError(c, http.StatusBadRequest, "INVALID_TASK", err.Error())
	case errors.Is(err, orchestration.ErrNoEligibleDevice):
		apiError(c, http.StatusConflict, "NO_ELIGIBLE_DEVICE", err.Error())
	case errors.Is(err, routing.ErrBusy):
		apiError(c, http.StatusTooManyRequests, "ROUTING_BUSY", err.Error())
	case errors.Is(err, routing.ErrTimeout):
		apiError(c, http.StatusGatewayTimeout, "ROUTING_TIMEOUT", err.Error())
	case errors.Is(err, routing.ErrUnavailable):
		apiError(c, http.StatusServiceUnavailable, "ROUTING_UNAVAILABLE", err.Error())
	case errors.Is(err, routing.ErrNoRoute), errors.Is(err, routing.ErrNoRoadNode), errors.Is(err, routing.ErrUnsupportedProfile), errors.Is(err, routing.ErrOutsideRegion):
		apiError(c, http.StatusUnprocessableEntity, err.Error(), err.Error())
	case errors.Is(err, extension.ErrPlanningRequired):
		apiError(c, http.StatusUnprocessableEntity, "PLANNER_UNAVAILABLE", err.Error())
	case errors.Is(err, repository.ErrNotFound):
		apiError(c, http.StatusNotFound, "NOT_FOUND", "Resource was not found")
	case errors.Is(err, repository.ErrForbidden):
		apiError(c, http.StatusForbidden, "FORBIDDEN", "Command identity does not match the authenticated device")
	case errors.Is(err, repository.ErrConflict), errors.Is(err, repository.ErrInvalidTransition):
		apiError(c, http.StatusConflict, "INVALID_STATE_TRANSITION", "Operation is not legal in the current durable state")
	default:
		apiError(c, http.StatusInternalServerError, "ORCHESTRATION_ERROR", "Orchestration operation failed")
	}
}

func (a *OrchestrationAPI) createTask(c *gin.Context) {
	tenant, ok := tenantFor(c)
	if !ok {
		apiError(c, 400, "TENANT_REQUIRED", "Tenant scope is required")
		return
	}
	principal, _ := Principal(c)
	var in struct {
		ProjectID    *string               `json:"project_id"`
		TaskType     string                `json:"task_type"`
		Priority     string                `json:"priority"`
		Requirements taskcore.Requirements `json:"requirements"`
		Target       json.RawMessage       `json:"target"`
		ExpiresAt    *time.Time            `json:"expires_at"`
	}
	if c.ShouldBindJSON(&in) != nil {
		apiError(c, 400, "INVALID_REQUEST", "A valid task document is required")
		return
	}
	if in.Priority == "" {
		in.Priority = "NORMAL"
	}
	expiresAt := time.Now().UTC().Add(5 * time.Minute)
	if in.ExpiresAt != nil {
		expiresAt = in.ExpiresAt.UTC()
	}
	result, err := a.service.CreateTask(c, tenant, principal, RequestID(c), orchestration.CreateTaskInput{ProjectID: in.ProjectID, TaskType: in.TaskType, Priority: in.Priority, Requirements: in.Requirements, Target: in.Target, ExpiresAt: expiresAt, CorrelationID: c.GetHeader("X-Correlation-ID")})
	if err != nil {
		orchestrationError(c, err)
		return
	}
	apiData(c, http.StatusCreated, result)
}

func (a *OrchestrationAPI) listTasks(c *gin.Context) {
	tenant, ok := tenantFor(c)
	if !ok {
		apiError(c, 400, "TENANT_REQUIRED", "Tenant scope is required")
		return
	}
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "50"))
	result, err := a.store.ListTasks(c, tenant, limit, c.Query("cursor"), c.Query("status"), c.Query("device_id"))
	if err != nil {
		orchestrationError(c, err)
		return
	}
	apiData(c, 200, result)
}

func (a *OrchestrationAPI) getTask(c *gin.Context) {
	tenant, ok := tenantFor(c)
	if !ok {
		apiError(c, 400, "TENANT_REQUIRED", "Tenant scope is required")
		return
	}
	v, err := a.store.GetTask(c, tenant, c.Param("task_id"))
	if err != nil {
		orchestrationError(c, err)
		return
	}
	commands, err := a.store.ListCommands(c, tenant, 100, "", "", v.TaskID, "")
	if err != nil {
		orchestrationError(c, err)
		return
	}
	apiData(c, 200, gin.H{"task": v, "commands": commands})
}

func (a *OrchestrationAPI) cancelTask(c *gin.Context) {
	tenant, _ := tenantFor(c)
	principal, _ := Principal(c)
	if err := a.store.CancelTask(c, tenant, c.Param("task_id"), principal.APIKeyID, RequestID(c)); err != nil {
		orchestrationError(c, err)
		return
	}
	apiData(c, 200, gin.H{"task_id": c.Param("task_id"), "status": taskcore.Cancelled})
}

func (a *OrchestrationAPI) retryTask(c *gin.Context) {
	tenant, _ := tenantFor(c)
	principal, _ := Principal(c)
	var in struct {
		TTLSeconds int `json:"ttl_seconds"`
	}
	_ = c.ShouldBindJSON(&in)
	result, err := a.service.RetryTask(c, tenant, c.Param("task_id"), principal, RequestID(c), time.Duration(in.TTLSeconds)*time.Second)
	if err != nil {
		orchestrationError(c, err)
		return
	}
	apiData(c, 200, result)
}

func (a *OrchestrationAPI) listCommands(c *gin.Context) {
	tenant, ok := tenantFor(c)
	if !ok {
		apiError(c, 400, "TENANT_REQUIRED", "Tenant scope is required")
		return
	}
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "50"))
	result, err := a.store.ListCommands(c, tenant, limit, c.Query("cursor"), c.Query("status"), c.Query("task_id"), c.Query("device_id"))
	if err != nil {
		orchestrationError(c, err)
		return
	}
	apiData(c, 200, result)
}

func (a *OrchestrationAPI) getCommand(c *gin.Context) {
	tenant, ok := tenantFor(c)
	if !ok {
		apiError(c, 400, "TENANT_REQUIRED", "Tenant scope is required")
		return
	}
	v, err := a.store.GetCommand(c, tenant, c.Param("command_id"))
	if err != nil {
		orchestrationError(c, err)
		return
	}
	apiData(c, 200, v)
}

func (a *OrchestrationAPI) retryCommand(c *gin.Context) {
	tenant, _ := tenantFor(c)
	principal, _ := Principal(c)
	if err := a.store.RetryCommand(c, tenant, c.Param("command_id"), principal.APIKeyID, RequestID(c), true); err != nil {
		orchestrationError(c, err)
		return
	}
	apiData(c, 200, gin.H{"command_id": c.Param("command_id"), "status": command.Pending})
}

func (a *OrchestrationAPI) cancelCommand(c *gin.Context) {
	tenant, _ := tenantFor(c)
	principal, _ := Principal(c)
	v, err := a.store.GetCommand(c, tenant, c.Param("command_id"))
	if err == nil {
		err = a.store.CancelTask(c, tenant, v.TaskID, principal.APIKeyID, RequestID(c))
	}
	if err != nil {
		orchestrationError(c, err)
		return
	}
	apiData(c, 200, gin.H{"command_id": v.CommandID, "status": command.Cancelled})
}
```

---

## internal\adapter\handler\registry_api.go

```
package handler

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/Akashpg-M/polaris/backend/internal/adapter/repository"
	"github.com/Akashpg-M/polaris/backend/internal/core/auth"
	"github.com/Akashpg-M/polaris/backend/internal/core/registry"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

const (
	operatorPrincipalKey = "operator_principal"
	requestIDKey         = "request_id"
)

func RequestID(c *gin.Context) string {
	if value, ok := c.Get(requestIDKey); ok {
		if id, valid := value.(string); valid && id != "" {
			return id
		}
	}
	v := c.GetHeader("X-Request-ID")
	if v == "" {
		v = auth.NewID()
	}
	c.Set(requestIDKey, v)
	return v
}
func Principal(c *gin.Context) (auth.OperatorPrincipal, bool) {
	v, ok := c.Get(operatorPrincipalKey)
	if !ok {
		return auth.OperatorPrincipal{}, false
	}
	p, ok := v.(auth.OperatorPrincipal)
	return p, ok
}

type RegistryAPI struct {
	store                               *repository.RegistryStore
	redis                               *redis.Client
	staleAfter, offlineAfter, ticketTTL time.Duration
	lifecycleHook                       func(tenant, device, status string)
}

func (a *RegistryAPI) SetLifecycleHook(fn func(tenant, device, status string)) { a.lifecycleHook = fn }

func NewRegistryAPI(store *repository.RegistryStore, redisClient *redis.Client, staleAfter, offlineAfter, ticketTTL time.Duration) *RegistryAPI {
	return &RegistryAPI{store: store, redis: redisClient, staleAfter: staleAfter, offlineAfter: offlineAfter, ticketTTL: ticketTTL}
}

func (a *RegistryAPI) Middleware(permission string) gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.GetHeader("Authorization")
		if !strings.HasPrefix(header, "Bearer ") {
			apiError(c, http.StatusUnauthorized, "UNAUTHENTICATED", "Operator API key is required")
			c.Abort()
			return
		}
		p, err := a.store.ResolveOperator(c.Request.Context(), strings.TrimSpace(strings.TrimPrefix(header, "Bearer ")))
		if err != nil {
			apiError(c, http.StatusUnauthorized, "UNAUTHENTICATED", "Operator API key is invalid")
			c.Abort()
			return
		}
		if !auth.Can(p.Role, permission) {
			apiError(c, http.StatusForbidden, "FORBIDDEN", "Role does not permit this operation")
			c.Abort()
			return
		}
		c.Set(operatorPrincipalKey, p)
		c.Next()
	}
}
func apiData(c *gin.Context, status int, data interface{}) {
	c.JSON(status, gin.H{"data": data, "request_id": RequestID(c)})
}
func apiError(c *gin.Context, status int, code, message string) {
	c.JSON(status, gin.H{"error": gin.H{"code": code, "message": message}, "request_id": RequestID(c)})
}
func registryError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, repository.ErrNotFound):
		apiError(c, 404, "NOT_FOUND", "Resource was not found")
	case errors.Is(err, repository.ErrConflict):
		apiError(c, 409, "CONFLICT", "Resource already exists")
	case errors.Is(err, repository.ErrInvalidTransition):
		apiError(c, 409, "INVALID_LIFECYCLE_TRANSITION", "Lifecycle transition is not allowed")
	default:
		apiError(c, 500, "INTERNAL_ERROR", "Registry operation failed")
	}
}
func tenantFor(c *gin.Context) (string, bool) {
	p, ok := Principal(c)
	if !ok {
		return "", false
	}
	if p.Role == auth.PlatformAdmin {
		tenant := c.GetHeader("X-Tenant-ID")
		if tenant == "" {
			tenant = c.Query("tenant_id")
		}
		return tenant, tenant != ""
	}
	return p.TenantID, p.TenantID != ""
}

func (a *RegistryAPI) Register(r *gin.RouterGroup) {
	r.POST("/tenants", a.Middleware("mutate"), a.createTenant)
	r.GET("/tenants/:tenant_id", a.Middleware("read"), a.getTenant)
	r.PATCH("/tenants/:tenant_id", a.Middleware("mutate"), a.patchTenant)
	r.POST("/projects", a.Middleware("mutate"), a.createProject)
	r.GET("/projects", a.Middleware("read"), a.listProjects)
	r.GET("/projects/:project_id", a.Middleware("read"), a.getProject)
	r.PATCH("/projects/:project_id", a.Middleware("mutate"), a.patchProject)
	r.POST("/devices", a.Middleware("mutate"), a.createDevice)
	r.GET("/devices", a.Middleware("read"), a.listDevices)
	r.GET("/devices/:device_id", a.Middleware("read"), a.getDevice)
	r.PATCH("/devices/:device_id", a.Middleware("mutate"), a.patchDevice)
	r.POST("/devices/:device_id/activate", a.Middleware("mutate"), a.lifecycle("ACTIVE"))
	r.POST("/devices/:device_id/suspend", a.Middleware("mutate"), a.lifecycle("SUSPENDED"))
	r.POST("/devices/:device_id/decommission", a.Middleware("mutate"), a.lifecycle("DECOMMISSIONED"))
	r.GET("/capabilities", a.Middleware("read"), a.allCapabilities)
	r.GET("/devices/:device_id/capabilities", a.Middleware("read"), a.deviceCapabilities)
	r.PUT("/devices/:device_id/capabilities/:capability_id", a.Middleware("mutate"), a.putCapability)
	r.DELETE("/devices/:device_id/capabilities/:capability_id", a.Middleware("mutate"), a.removeCapability)
	r.POST("/devices/:device_id/credentials", a.Middleware("mutate"), a.issueCredential)
	r.GET("/devices/:device_id/credentials", a.Middleware("read"), a.listCredentials)
	r.POST("/devices/:device_id/credentials/:credential_id/revoke", a.Middleware("mutate"), a.revokeCredential)
	r.POST("/devices/:device_id/credentials/rotate", a.Middleware("mutate"), a.rotateCredential)
	r.POST("/devices/:device_id/connection-ticket", a.Middleware("mutate"), a.connectionTicket)
	r.GET("/devices/:device_id/twin", a.Middleware("read"), a.getTwin)
	r.GET("/twins", a.Middleware("read"), a.listTwins)
	r.GET("/audit-events", a.Middleware("audit"), a.listAudit)
	r.POST("/dashboard-ticket", a.Middleware("read"), a.dashboardTicket)
}

func (a *RegistryAPI) createTenant(c *gin.Context) {
	p, _ := Principal(c)
	if p.Role != auth.PlatformAdmin {
		apiError(c, 403, "FORBIDDEN", "Only platform admins create tenants")
		return
	}
	var in struct {
		TenantID    string          `json:"tenant_id"`
		DisplayName string          `json:"display_name"`
		Metadata    json.RawMessage `json:"metadata"`
	}
	if c.ShouldBindJSON(&in) != nil || in.TenantID == "" || in.DisplayName == "" {
		apiError(c, 400, "INVALID_REQUEST", "tenant_id and display_name are required")
		return
	}
	t := registry.Tenant{TenantID: in.TenantID, DisplayName: in.DisplayName, Status: "ACTIVE", Metadata: jsonOrEmptyRaw(in.Metadata)}
	if err := a.store.CreateTenant(c, t, p.APIKeyID, RequestID(c)); err != nil {
		registryError(c, err)
		return
	}
	apiData(c, 201, t)
}
func (a *RegistryAPI) getTenant(c *gin.Context) {
	p, _ := Principal(c)
	id := c.Param("tenant_id")
	if p.Role != auth.PlatformAdmin && p.TenantID != id {
		_ = a.store.Audit(c, p.TenantID, p.APIKeyID, "CROSS_TENANT_ACCESS_DENIED", "tenant", id, RequestID(c), "DENIED")
		apiError(c, 404, "NOT_FOUND", "Resource was not found")
		return
	}
	v, err := a.store.GetTenant(c, id)
	if err != nil {
		registryError(c, err)
		return
	}
	apiData(c, 200, v)
}
func (a *RegistryAPI) patchTenant(c *gin.Context) {
	p, _ := Principal(c)
	if p.Role != auth.PlatformAdmin {
		apiError(c, 403, "FORBIDDEN", "Only platform admins change tenant lifecycle")
		return
	}
	var in struct {
		Status string `json:"status"`
	}
	if c.ShouldBindJSON(&in) != nil {
		apiError(c, 400, "INVALID_REQUEST", "status is required")
		return
	}
	if err := a.store.SetTenantStatus(c, c.Param("tenant_id"), in.Status, p.APIKeyID, RequestID(c)); err != nil {
		registryError(c, err)
		return
	}
	if a.lifecycleHook != nil {
		a.lifecycleHook(c.Param("tenant_id"), "", in.Status)
	}
	apiData(c, 200, gin.H{"tenant_id": c.Param("tenant_id"), "status": in.Status})
}
func jsonOrEmptyRaw(v json.RawMessage) []byte {
	if len(v) == 0 {
		return []byte(`{}`)
	}
	return v
}
func (a *RegistryAPI) createProject(c *gin.Context) {
	tenant, ok := tenantFor(c)
	if !ok {
		apiError(c, 400, "TENANT_REQUIRED", "Tenant scope is required")
		return
	}
	p, _ := Principal(c)
	var in struct {
		Name        string          `json:"name"`
		Description *string         `json:"description"`
		Metadata    json.RawMessage `json:"metadata"`
	}
	if c.ShouldBindJSON(&in) != nil || in.Name == "" {
		apiError(c, 400, "INVALID_REQUEST", "name is required")
		return
	}
	v := registry.Project{ProjectID: auth.NewID(), TenantID: tenant, Name: in.Name, Description: in.Description, Status: "ACTIVE", Metadata: jsonOrEmptyRaw(in.Metadata)}
	if err := a.store.CreateProject(c, v, p.APIKeyID, RequestID(c)); err != nil {
		registryError(c, err)
		return
	}
	apiData(c, 201, v)
}
func (a *RegistryAPI) listProjects(c *gin.Context) {
	tenant, ok := tenantFor(c)
	if !ok {
		apiError(c, 400, "TENANT_REQUIRED", "Tenant scope is required")
		return
	}
	v, err := a.store.ListProjects(c, tenant)
	if err != nil {
		registryError(c, err)
		return
	}
	apiData(c, 200, v)
}
func (a *RegistryAPI) getProject(c *gin.Context) {
	tenant, ok := tenantFor(c)
	if !ok {
		apiError(c, 400, "TENANT_REQUIRED", "Tenant scope is required")
		return
	}
	v, err := a.store.GetProject(c, tenant, c.Param("project_id"))
	if err != nil {
		registryError(c, err)
		return
	}
	apiData(c, 200, v)
}
func (a *RegistryAPI) patchProject(c *gin.Context) {
	tenant, ok := tenantFor(c)
	if !ok {
		apiError(c, 400, "TENANT_REQUIRED", "Tenant scope is required")
		return
	}
	p, _ := Principal(c)
	var in struct {
		Name        *string         `json:"name"`
		Description *string         `json:"description"`
		Status      *string         `json:"status"`
		Metadata    json.RawMessage `json:"metadata"`
	}
	if c.ShouldBindJSON(&in) != nil {
		apiError(c, 400, "INVALID_REQUEST", "A valid project update is required")
		return
	}
	if err := a.store.UpdateProject(c, tenant, c.Param("project_id"), in.Name, in.Description, in.Status, in.Metadata, p.APIKeyID, RequestID(c)); err != nil {
		registryError(c, err)
		return
	}
	v, err := a.store.GetProject(c, tenant, c.Param("project_id"))
	if err != nil {
		registryError(c, err)
		return
	}
	apiData(c, 200, v)
}
func (a *RegistryAPI) createDevice(c *gin.Context) {
	tenant, ok := tenantFor(c)
	if !ok {
		apiError(c, 400, "TENANT_REQUIRED", "Tenant scope is required")
		return
	}
	p, _ := Principal(c)
	var in struct {
		DeviceID     string          `json:"device_id"`
		ProjectID    *string         `json:"project_id"`
		DeviceTypeID string          `json:"device_type_id"`
		DisplayName  string          `json:"display_name"`
		Metadata     json.RawMessage `json:"metadata"`
	}
	if c.ShouldBindJSON(&in) != nil || in.DeviceID == "" || in.DeviceTypeID == "" || in.DisplayName == "" {
		apiError(c, 400, "INVALID_REQUEST", "device_id, device_type_id and display_name are required")
		return
	}
	v := registry.Device{TenantID: tenant, DeviceID: in.DeviceID, ProjectID: in.ProjectID, DeviceTypeID: in.DeviceTypeID, DisplayName: in.DisplayName, LifecycleStatus: "REGISTERED", Metadata: jsonOrEmptyRaw(in.Metadata)}
	if err := a.store.CreateDevice(c, v, p.APIKeyID, RequestID(c)); err != nil {
		registryError(c, err)
		return
	}
	apiData(c, 201, v)
}
func (a *RegistryAPI) listDevices(c *gin.Context) {
	tenant, ok := tenantFor(c)
	if !ok {
		apiError(c, 400, "TENANT_REQUIRED", "Tenant scope is required")
		return
	}
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "50"))
	v, err := a.store.ListDevices(c, tenant, limit, c.Query("cursor"), c.Query("project_id"), c.Query("device_type"), c.Query("lifecycle_status"), c.Query("capability"))
	if err != nil {
		registryError(c, err)
		return
	}
	apiData(c, 200, v)
}
func (a *RegistryAPI) getDevice(c *gin.Context) {
	tenant, ok := tenantFor(c)
	if !ok {
		return
	}
	v, err := a.store.GetDevice(c, tenant, c.Param("device_id"))
	if err != nil {
		registryError(c, err)
		return
	}
	apiData(c, 200, v)
}
func (a *RegistryAPI) patchDevice(c *gin.Context) {
	tenant, ok := tenantFor(c)
	if !ok {
		apiError(c, 400, "TENANT_REQUIRED", "Tenant scope is required")
		return
	}
	p, _ := Principal(c)
	var in struct {
		DisplayName     *string         `json:"display_name"`
		FirmwareVersion *string         `json:"firmware_version"`
		SoftwareVersion *string         `json:"software_version"`
		ModelVersion    *string         `json:"model_version"`
		Metadata        json.RawMessage `json:"metadata"`
	}
	if c.ShouldBindJSON(&in) != nil {
		apiError(c, 400, "INVALID_REQUEST", "A valid device update is required")
		return
	}
	if err := a.store.UpdateDevice(c, tenant, c.Param("device_id"), in.DisplayName, in.FirmwareVersion, in.SoftwareVersion, in.ModelVersion, in.Metadata, p.APIKeyID, RequestID(c)); err != nil {
		registryError(c, err)
		return
	}
	v, err := a.store.GetDevice(c, tenant, c.Param("device_id"))
	if err != nil {
		registryError(c, err)
		return
	}
	apiData(c, 200, v)
}
func (a *RegistryAPI) lifecycle(next string) gin.HandlerFunc {
	return func(c *gin.Context) {
		tenant, ok := tenantFor(c)
		if !ok {
			return
		}
		p, _ := Principal(c)
		if err := a.store.SetLifecycle(c, tenant, c.Param("device_id"), next, p.APIKeyID, RequestID(c)); err != nil {
			registryError(c, err)
			return
		}
		if a.lifecycleHook != nil {
			a.lifecycleHook(tenant, c.Param("device_id"), next)
		}
		apiData(c, 200, gin.H{"device_id": c.Param("device_id"), "lifecycle_status": next})
	}
}
func (a *RegistryAPI) allCapabilities(c *gin.Context) {
	v, err := a.store.AllCapabilities(c)
	if err != nil {
		registryError(c, err)
		return
	}
	apiData(c, 200, v)
}
func (a *RegistryAPI) deviceCapabilities(c *gin.Context) {
	tenant, _ := tenantFor(c)
	v, err := a.store.ListCapabilities(c, tenant, c.Param("device_id"))
	if err != nil {
		registryError(c, err)
		return
	}
	apiData(c, 200, v)
}
func (a *RegistryAPI) putCapability(c *gin.Context) {
	tenant, _ := tenantFor(c)
	p, _ := Principal(c)
	var in struct {
		Configuration json.RawMessage `json:"configuration"`
	}
	_ = c.ShouldBindJSON(&in)
	if err := a.store.PutCapability(c, tenant, c.Param("device_id"), c.Param("capability_id"), jsonOrEmptyRaw(in.Configuration), p.APIKeyID, RequestID(c)); err != nil {
		registryError(c, err)
		return
	}
	apiData(c, 200, gin.H{"enabled": true})
}
func (a *RegistryAPI) removeCapability(c *gin.Context) {
	tenant, _ := tenantFor(c)
	p, _ := Principal(c)
	if err := a.store.RemoveCapability(c, tenant, c.Param("device_id"), c.Param("capability_id"), p.APIKeyID, RequestID(c)); err != nil {
		registryError(c, err)
		return
	}
	apiData(c, 200, gin.H{"enabled": false})
}
func (a *RegistryAPI) issueCredential(c *gin.Context) {
	tenant, _ := tenantFor(c)
	p, _ := Principal(c)
	meta, raw, err := a.store.IssueCredential(c, tenant, c.Param("device_id"), p.APIKeyID, RequestID(c), nil)
	if err != nil {
		registryError(c, err)
		return
	}
	apiData(c, 201, gin.H{"credential": meta, "secret": raw})
}
func (a *RegistryAPI) listCredentials(c *gin.Context) {
	tenant, _ := tenantFor(c)
	v, err := a.store.ListCredentials(c, tenant, c.Param("device_id"))
	if err != nil {
		registryError(c, err)
		return
	}
	apiData(c, 200, v)
}
func (a *RegistryAPI) revokeCredential(c *gin.Context) {
	tenant, _ := tenantFor(c)
	p, _ := Principal(c)
	if err := a.store.RevokeCredential(c, tenant, c.Param("device_id"), c.Param("credential_id"), p.APIKeyID, RequestID(c)); err != nil {
		registryError(c, err)
		return
	}
	apiData(c, 200, gin.H{"status": "REVOKED"})
}
func (a *RegistryAPI) rotateCredential(c *gin.Context) {
	tenant, _ := tenantFor(c)
	p, _ := Principal(c)
	var in struct {
		CredentialID string `json:"credential_id"`
	}
	if c.ShouldBindJSON(&in) != nil || in.CredentialID == "" {
		apiError(c, 400, "INVALID_REQUEST", "credential_id is required")
		return
	}
	meta, raw, err := a.store.RotateCredential(c, tenant, c.Param("device_id"), in.CredentialID, p.APIKeyID, RequestID(c))
	if err != nil {
		registryError(c, err)
		return
	}
	apiData(c, 201, gin.H{"credential": meta, "secret": raw})
}
func (a *RegistryAPI) connectionTicket(c *gin.Context) {
	tenant, _ := tenantFor(c)
	credentials, err := a.store.ListCredentials(c, tenant, c.Param("device_id"))
	if err != nil || len(credentials) == 0 {
		registryError(c, repository.ErrNotFound)
		return
	}
	active := ""
	for _, v := range credentials {
		if v.Status == "ACTIVE" {
			active = v.CredentialID
			break
		}
	}
	if active == "" {
		registryError(c, repository.ErrNotFound)
		return
	}
	ticket, err := a.store.CreateTicket(c, tenant, c.Param("device_id"), active, a.ticketTTL)
	if err != nil {
		registryError(c, err)
		return
	}
	apiData(c, 201, gin.H{"ticket": ticket, "expires_in_seconds": int(a.ticketTTL.Seconds())})
}

func connectivity(lastSeen string, stale, offline time.Duration) (string, *time.Time) {
	if lastSeen == "" {
		return "NEVER_CONNECTED", nil
	}
	ms, err := strconv.ParseInt(lastSeen, 10, 64)
	if err != nil {
		return "NEVER_CONNECTED", nil
	}
	t := time.UnixMilli(ms).UTC()
	age := time.Since(t)
	if age > offline {
		return "OFFLINE", &t
	}
	if age > stale {
		return "STALE", &t
	}
	return "ONLINE", &t
}
func (a *RegistryAPI) twin(ctx context.Context, tenant, id string) (gin.H, error) {
	device, err := a.store.GetDevice(ctx, tenant, id)
	if err != nil {
		return nil, err
	}
	caps, err := a.store.ListCapabilities(ctx, tenant, id)
	if err != nil {
		return nil, err
	}
	key := "polaris:twin:" + tenant + ":" + id
	state, err := a.redis.HGetAll(ctx, key).Result()
	if err != nil {
		return nil, err
	}
	status, last := connectivity(state["last_seen_at"], a.staleAfter, a.offlineAfter)
	var reported interface{}
	if raw := state["reported_state"]; raw != "" {
		_ = json.Unmarshal([]byte(raw), &reported)
	}
	components := map[string]interface{}{}
	for field, raw := range state {
		if strings.HasPrefix(field, "component:") {
			var component interface{}
			if json.Unmarshal([]byte(raw), &component) == nil {
				components[strings.TrimPrefix(field, "component:")] = component
			}
		}
	}
	return gin.H{"tenant_id": tenant, "device_id": id, "device": device, "capabilities": caps, "reported_state": reported, "components": components, "desired_state": nil, "connectivity": gin.H{"status": status, "last_seen_at": last}}, nil
}
func (a *RegistryAPI) getTwin(c *gin.Context) {
	tenant, _ := tenantFor(c)
	v, err := a.twin(c, tenant, c.Param("device_id"))
	if err != nil {
		registryError(c, err)
		return
	}
	apiData(c, 200, v)
}
func (a *RegistryAPI) listTwins(c *gin.Context) {
	tenant, _ := tenantFor(c)
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "50"))
	devices, err := a.store.ListDevices(c, tenant, limit, c.Query("cursor"), c.Query("project_id"), c.Query("device_type"), c.Query("lifecycle_status"), c.Query("capability"))
	if err != nil {
		registryError(c, err)
		return
	}
	out := []gin.H{}
	for _, d := range devices {
		v, e := a.twin(c, tenant, d.DeviceID)
		if e == nil {
			if filter := c.Query("connectivity_status"); filter == "" || v["connectivity"].(gin.H)["status"] == filter {
				out = append(out, v)
			}
		}
	}
	apiData(c, 200, out)
}
func (a *RegistryAPI) listAudit(c *gin.Context) {
	tenant, _ := tenantFor(c)
	v, err := a.store.ListAudit(c, tenant)
	if err != nil {
		registryError(c, err)
		return
	}
	apiData(c, 200, v)
}
func (a *RegistryAPI) dashboardTicket(c *gin.Context) {
	p, _ := Principal(c)
	tenant, _ := tenantFor(c)
	ticket, err := a.store.CreateOperatorTicket(c, p, tenant, a.ticketTTL)
	if err != nil {
		registryError(c, err)
		return
	}
	apiData(c, 201, gin.H{"ticket": ticket, "expires_in_seconds": int(a.ticketTTL.Seconds())})
}
```

---

## internal\adapter\handler\registry_api_test.go

```
package handler

import (
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestRequestIDIsStableWithinRequest(t *testing.T) {
	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	c.Request = httptest.NewRequest("GET", "/", nil)
	first := RequestID(c)
	if first == "" || RequestID(c) != first {
		t.Fatal("generated request ID changed within one request")
	}
}

func TestRequestIDPreservesCallerValue(t *testing.T) {
	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	c.Request = httptest.NewRequest("GET", "/", nil)
	c.Request.Header.Set("X-Request-ID", "caller-request")
	if got := RequestID(c); got != "caller-request" || RequestID(c) != got {
		t.Fatalf("request ID = %q, want caller-request", got)
	}
}
```

---

## internal\adapter\repository\kafka_event_publisher.go

```
package repository

import (
	"context"
	"fmt"

	"github.com/Akashpg-M/polaris/backend/internal/core/events"
	"github.com/segmentio/kafka-go"
)

// KafkaEventPublisher writes the canonical versioned envelope and partitions by
// stable tenant+device identity so all boots and sequences stay ordered.
type KafkaEventPublisher struct {
	writer *kafka.Writer
}

func NewKafkaEventPublisher(brokerURL string) *KafkaEventPublisher {
	return &KafkaEventPublisher{writer: &kafka.Writer{
		Addr:     kafka.TCP(brokerURL),
		Topic:    TelemetryTopic,
		Balancer: &kafka.Hash{},
		Async:    false,
	}}
}

func (p *KafkaEventPublisher) PublishEvent(ctx context.Context, _ string, event interface{}) error {
	envelope, ok := event.(*events.TelemetryEnvelope)
	if !ok || envelope == nil {
		return fmt.Errorf("unsupported kafka event type %T", event)
	}
	data, err := envelope.Marshal()
	if err != nil {
		return fmt.Errorf("marshal telemetry envelope: %w", err)
	}
	partitionKey := envelope.PartitionKey()
	return p.writer.WriteMessages(ctx, kafka.Message{Key: []byte(partitionKey), Value: data})
}

func (p *KafkaEventPublisher) Ready(ctx context.Context, brokerURL string) error {
	conn, err := kafka.DialContext(ctx, "tcp", brokerURL)
	if err != nil {
		return err
	}
	return conn.Close()
}

func (p *KafkaEventPublisher) Close() error { return p.writer.Close() }
```

---

## internal\adapter\repository\kafka_stream.go

```
package repository

import (
	"context"
	"fmt"
	"log/slog"

	pb "github.com/Akashpg-M/polaris/backend/api/proto/v1"
	"github.com/segmentio/kafka-go"
	"github.com/uber/h3-go/v4"
	"google.golang.org/protobuf/proto"
)

const TelemetryTopic = "telemetry.ingress"
const H3Resolution = 7 // Resolution 7 gives roughly 1.2km wide hexagons

type KafkaStreamAdapter struct {
	writer *kafka.Writer
}

func NewKafkaStreamAdapter(brokerURL string) *KafkaStreamAdapter {
	writer := &kafka.Writer{
		Addr:     kafka.TCP(brokerURL),
		Topic:    TelemetryTopic,
		// Hash balancer ensures the H3 Key determines the partition
		Balancer: &kafka.Hash{}, 
		// Async mode for extreme throughput (acts like a shock absorber)
		Async:    true,
	}

	return &KafkaStreamAdapter{writer: writer}
}

func (k *KafkaStreamAdapter) Publish(ctx context.Context, payload *pb.SpatialObject) error {
	// 1. Calculate the Spatial Partition Key (Uber H3)
	latLng := h3.NewLatLng(payload.Lat, payload.Lon)
	
	// FIXED: Catch the error returned by LatLngToCell
	h3Cell, err := h3.LatLngToCell(latLng, H3Resolution)
	if err != nil {
		return fmt.Errorf("failed to calculate H3 cell: %w", err)
	}
	h3Key := h3Cell.String()

	// 2. Serialize to raw Protocol Buffers
	data, err := proto.Marshal(payload)
	if err != nil {
		return fmt.Errorf("protobuf marshal failed: %w", err)
	}

	// 3. Write to Kafka using the H3 Hex as the routing Key
	err = k.writer.WriteMessages(ctx, kafka.Message{
		Key:   []byte(h3Key),
		Value: data,
	})

	if err != nil {
		slog.Error("Kafka write failed", "error", err)
		return err
	}

	return nil
}

func (k *KafkaStreamAdapter) Close() error {
	return k.writer.Close()
}
```

---

## internal\adapter\repository\postgres_orchestration.go

```
package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"time"

	"github.com/Akashpg-M/polaris/backend/internal/core/auth"
	"github.com/Akashpg-M/polaris/backend/internal/core/command"
	"github.com/Akashpg-M/polaris/backend/internal/core/extension"
	taskcore "github.com/Akashpg-M/polaris/backend/internal/core/task"
	"github.com/lib/pq"
)

type DeviceCandidate struct {
	DeviceID     string `db:"device_id"`
	DeviceTypeID string `db:"device_type_id"`
}

func (s *RegistryStore) CreateTask(ctx context.Context, v taskcore.Task, actor, requestID string) error {
	tx, err := s.DB.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	_, err = tx.ExecContext(ctx, `INSERT INTO tasks(task_id,tenant_id,project_id,task_type,status,priority,requirements,target,correlation_id,created_by,expires_at) VALUES($1,$2,$3,$4,'PENDING',$5,$6,$7,$8,$9,$10)`, v.TaskID, v.TenantID, v.ProjectID, v.TaskType, v.Priority, v.Requirements, v.Target, v.CorrelationID, v.CreatedBy, v.ExpiresAt)
	if err != nil {
		return mapPQ(err)
	}
	if err = insertAuditOutbox(ctx, tx, v.TenantID, actor, "TASK_CREATED", "task", v.TaskID, requestID, "task.created.v1", map[string]interface{}{"task_id": v.TaskID, "task_type": v.TaskType, "priority": v.Priority}); err != nil {
		return err
	}
	return tx.Commit()
}

func (s *RegistryStore) GetTask(ctx context.Context, tenant, id string) (taskcore.Task, error) {
	var v taskcore.Task
	err := s.DB.GetContext(ctx, &v, `SELECT * FROM tasks WHERE tenant_id=$1 AND task_id=$2`, tenant, id)
	if errors.Is(err, sql.ErrNoRows) {
		err = ErrNotFound
	}
	return v, err
}

func (s *RegistryStore) ListTasks(ctx context.Context, tenant string, limit int, cursor, status, deviceID string) ([]taskcore.Task, error) {
	if limit < 1 || limit > 100 {
		limit = 50
	}
	args := []interface{}{tenant, limit}
	q := `SELECT * FROM tasks WHERE tenant_id=$1`
	add := func(clause string, value interface{}) {
		args = append(args, value)
		q += fmt.Sprintf(clause, len(args))
	}
	if cursor != "" {
		add(` AND task_id>$%d`, cursor)
	}
	if status != "" {
		add(` AND status=$%d`, status)
	}
	if deviceID != "" {
		add(` AND assigned_device_id=$%d`, deviceID)
	}
	q += ` ORDER BY task_id LIMIT $2`
	result := []taskcore.Task{}
	return result, s.DB.SelectContext(ctx, &result, q, args...)
}

func (s *RegistryStore) PendingTasks(ctx context.Context, limit int) ([]taskcore.Task, error) {
	if limit < 1 {
		limit = 50
	}
	result := []taskcore.Task{}
	err := s.DB.SelectContext(ctx, &result, `SELECT * FROM tasks WHERE status='PENDING' AND expires_at>NOW() ORDER BY CASE priority WHEN 'CRITICAL' THEN 1 WHEN 'HIGH' THEN 2 WHEN 'NORMAL' THEN 3 ELSE 4 END,created_at LIMIT $1`, limit)
	return result, err
}

func (s *RegistryStore) EligibleDevices(ctx context.Context, tenant string, requirements taskcore.Requirements) ([]DeviceCandidate, error) {
	args := []interface{}{tenant}
	q := `SELECT d.device_id,d.device_type_id FROM devices d WHERE d.tenant_id=$1 AND d.lifecycle_status='ACTIVE'`
	if requirements.ProjectID != "" {
		args = append(args, requirements.ProjectID)
		q += fmt.Sprintf(` AND d.project_id=$%d`, len(args))
	}
	if len(requirements.AllowedDeviceTypes) > 0 {
		args = append(args, pq.Array(requirements.AllowedDeviceTypes))
		q += fmt.Sprintf(` AND d.device_type_id=ANY($%d)`, len(args))
	}
	for _, capability := range requirements.RequiredCapabilities {
		args = append(args, capability)
		q += fmt.Sprintf(` AND EXISTS(SELECT 1 FROM device_capabilities dc WHERE dc.tenant_id=d.tenant_id AND dc.device_id=d.device_id AND dc.capability_id=$%d AND dc.enabled)`, len(args))
	}
	q += ` AND NOT EXISTS(SELECT 1 FROM device_assignments a WHERE a.tenant_id=d.tenant_id AND a.device_id=d.device_id AND a.status='ACTIVE') ORDER BY d.device_id`
	result := []DeviceCandidate{}
	return result, s.DB.SelectContext(ctx, &result, q, args...)
}

func (s *RegistryStore) AssignTaskWithPlan(ctx context.Context, v taskcore.Task, deviceID, actor, requestID string, maxAttempts int, plan extension.ExecutionPlan) (command.Record, error) {
	if maxAttempts < 1 {
		maxAttempts = 5
	}
	tx, err := s.DB.BeginTxx(ctx, &sql.TxOptions{Isolation: sql.LevelReadCommitted})
	if err != nil {
		return command.Record{}, err
	}
	defer tx.Rollback()
	var current string
	if err = tx.GetContext(ctx, &current, `SELECT status FROM tasks WHERE tenant_id=$1 AND task_id=$2 FOR UPDATE`, v.TenantID, v.TaskID); errors.Is(err, sql.ErrNoRows) {
		return command.Record{}, ErrNotFound
	}
	if err != nil {
		return command.Record{}, err
	}
	if current != string(taskcore.Pending) {
		return command.Record{}, ErrInvalidTransition
	}
	if plan.CommandType == "" || len(plan.Payload) == 0 || !json.Valid(plan.Payload) {
		return command.Record{}, ErrInvalidTransition
	}
	commandExpiry := v.ExpiresAt
	if plan.ValidUntil != nil && plan.ValidUntil.Before(commandExpiry) {
		commandExpiry = *plan.ValidUntil
	}
	if !commandExpiry.After(time.Now()) {
		return command.Record{}, ErrInvalidTransition
	}
	assignmentID := auth.NewID()
	_, err = tx.ExecContext(ctx, `INSERT INTO device_assignments(assignment_id,tenant_id,device_id,task_id,status,lease_expires_at) VALUES($1,$2,$3,$4,'ACTIVE',$5)`, assignmentID, v.TenantID, deviceID, v.TaskID, v.ExpiresAt)
	if err != nil {
		return command.Record{}, mapPQ(err)
	}
	var sequence int64
	err = tx.GetContext(ctx, &sequence, `INSERT INTO device_command_sequences(tenant_id,device_id,last_sequence) VALUES($1,$2,1) ON CONFLICT(tenant_id,device_id) DO UPDATE SET last_sequence=device_command_sequences.last_sequence+1 RETURNING last_sequence`, v.TenantID, deviceID)
	if err != nil {
		return command.Record{}, err
	}
	now := time.Now().UTC()
	record := command.Record{CommandID: auth.NewID(), TenantID: v.TenantID, DeviceID: deviceID, TaskID: v.TaskID, CommandType: plan.CommandType, Payload: plan.Payload, Status: string(command.Pending), SequenceNumber: sequence, CorrelationID: v.CorrelationID, CausationID: v.TaskID, AttemptCount: 0, MaxAttempts: maxAttempts, Version: 1, CreatedAt: now, AvailableAt: now, ExpiresAt: commandExpiry}
	_, err = tx.ExecContext(ctx, `INSERT INTO commands(command_id,tenant_id,device_id,task_id,command_type,payload,status,sequence_number,correlation_id,causation_id,max_attempts,created_at,available_at,expires_at) VALUES($1,$2,$3,$4,$5,$6,'PENDING',$7,$8,$9,$10,$11,$11,$12)`, record.CommandID, record.TenantID, record.DeviceID, record.TaskID, record.CommandType, record.Payload, record.SequenceNumber, record.CorrelationID, record.CausationID, record.MaxAttempts, now, record.ExpiresAt)
	if err != nil {
		return command.Record{}, err
	}
	result, err := tx.ExecContext(ctx, `UPDATE tasks SET status='ASSIGNED',assigned_device_id=$3,assigned_at=NOW(),updated_at=NOW(),version=version+1 WHERE tenant_id=$1 AND task_id=$2 AND status='PENDING'`, v.TenantID, v.TaskID, deviceID)
	if err != nil {
		return command.Record{}, err
	}
	if n, _ := result.RowsAffected(); n != 1 {
		return command.Record{}, ErrInvalidTransition
	}
	if err = insertAuditOutbox(ctx, tx, v.TenantID, actor, "TASK_ASSIGNED", "task", v.TaskID, requestID, "task.assigned.v1", map[string]string{"task_id": v.TaskID, "device_id": deviceID}); err != nil {
		return command.Record{}, err
	}
	if err = insertAuditOutbox(ctx, tx, v.TenantID, actor, "COMMAND_CREATED", "command", record.CommandID, requestID, "command.created.v1", record.Envelope()); err != nil {
		return command.Record{}, err
	}
	return record, tx.Commit()
}

func (s *RegistryStore) GetCommand(ctx context.Context, tenant, id string) (command.Record, error) {
	var v command.Record
	err := s.DB.GetContext(ctx, &v, `SELECT * FROM commands WHERE tenant_id=$1 AND command_id=$2`, tenant, id)
	if errors.Is(err, sql.ErrNoRows) {
		err = ErrNotFound
	}
	return v, err
}

func (s *RegistryStore) ListCommands(ctx context.Context, tenant string, limit int, cursor, status, taskID, deviceID string) ([]command.Record, error) {
	if limit < 1 || limit > 100 {
		limit = 50
	}
	args := []interface{}{tenant, limit}
	q := `SELECT * FROM commands WHERE tenant_id=$1`
	add := func(clause string, value interface{}) {
		args = append(args, value)
		q += fmt.Sprintf(clause, len(args))
	}
	if cursor != "" {
		add(` AND command_id>$%d`, cursor)
	}
	if status != "" {
		add(` AND status=$%d`, status)
	}
	if taskID != "" {
		add(` AND task_id=$%d`, taskID)
	}
	if deviceID != "" {
		add(` AND device_id=$%d`, deviceID)
	}
	q += ` ORDER BY command_id LIMIT $2`
	result := []command.Record{}
	return result, s.DB.SelectContext(ctx, &result, q, args...)
}

func (s *RegistryStore) PendingCommandsForDevice(ctx context.Context, tenant, deviceID string) ([]command.Record, error) {
	result := []command.Record{}
	err := s.DB.SelectContext(ctx, &result, `SELECT * FROM commands WHERE tenant_id=$1 AND device_id=$2 AND status='PENDING' AND available_at<=NOW() AND expires_at>NOW() AND sequence_number=(SELECT MIN(sequence_number) FROM commands x WHERE x.tenant_id=$1 AND x.device_id=$2 AND x.status IN('PENDING','DELIVERED','ACKNOWLEDGED')) ORDER BY sequence_number`, tenant, deviceID)
	return result, err
}

func (s *RegistryStore) PrepareDelivery(ctx context.Context, tenant, deviceID, commandID, gatewayID string, epoch int64) (command.Record, error) {
	tx, err := s.DB.BeginTxx(ctx, nil)
	if err != nil {
		return command.Record{}, err
	}
	defer tx.Rollback()
	var v command.Record
	err = tx.GetContext(ctx, &v, `SELECT * FROM commands WHERE tenant_id=$1 AND device_id=$2 AND command_id=$3 FOR UPDATE`, tenant, deviceID, commandID)
	if errors.Is(err, sql.ErrNoRows) {
		return command.Record{}, ErrNotFound
	}
	if err != nil {
		return command.Record{}, err
	}
	if v.Status != string(command.Pending) || v.AvailableAt.After(time.Now()) {
		return command.Record{}, ErrInvalidTransition
	}
	if time.Now().After(v.ExpiresAt) || v.AttemptCount >= v.MaxAttempts {
		return command.Record{}, ErrInvalidTransition
	}
	var first int64
	if err = tx.GetContext(ctx, &first, `SELECT MIN(sequence_number) FROM commands WHERE tenant_id=$1 AND device_id=$2 AND status IN('PENDING','DELIVERED','ACKNOWLEDGED')`, tenant, deviceID); err != nil || first != v.SequenceNumber {
		return command.Record{}, ErrConflict
	}
	v.AttemptCount++
	v.Status = string(command.Delivered)
	now := time.Now().UTC()
	v.SentAt = &now
	_, err = tx.ExecContext(ctx, `UPDATE commands SET status='DELIVERED',attempt_count=$4,sent_at=$5,version=version+1,last_error=NULL WHERE tenant_id=$1 AND device_id=$2 AND command_id=$3 AND status='PENDING'`, tenant, deviceID, commandID, v.AttemptCount, now)
	if err != nil {
		return command.Record{}, err
	}
	_, err = tx.ExecContext(ctx, `INSERT INTO command_attempts(attempt_id,command_id,attempt_number,gateway_id,ownership_epoch) VALUES($1,$2,$3,$4,$5)`, auth.NewID(), commandID, v.AttemptCount, gatewayID, epoch)
	if err != nil {
		return command.Record{}, err
	}
	if err = insertAuditOutbox(ctx, tx, tenant, gatewayID, "COMMAND_DELIVERED", "command", commandID, "", "command.delivered.v1", map[string]interface{}{"command_id": commandID, "attempt": v.AttemptCount, "gateway_id": gatewayID, "ownership_epoch": epoch}); err != nil {
		return command.Record{}, err
	}
	return v, tx.Commit()
}

func (s *RegistryStore) ApplyCommandAck(ctx context.Context, principalTenant, principalDevice string, ack command.Ack) error {
	tx, err := s.DB.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	var v command.Record
	err = tx.GetContext(ctx, &v, `SELECT * FROM commands WHERE tenant_id=$1 AND device_id=$2 AND command_id=$3 FOR UPDATE`, principalTenant, principalDevice, ack.CommandID)
	if errors.Is(err, sql.ErrNoRows) {
		return ErrNotFound
	}
	if err != nil {
		return err
	}
	if ack.SequenceNumber != v.SequenceNumber {
		return ErrForbidden
	}
	if v.Status == string(command.Acknowledged) || command.IsTerminal(command.Status(v.Status)) {
		return tx.Commit()
	}
	if v.Status != string(command.Delivered) {
		return ErrInvalidTransition
	}
	now := time.Now().UTC()
	if ack.ReceivedAt.IsZero() {
		ack.ReceivedAt = now
	}
	accepted := ack.Status == "ACCEPTED" || ack.Status == "DUPLICATE"
	if accepted {
		_, err = tx.ExecContext(ctx, `UPDATE commands SET status='ACKNOWLEDGED',ack_status=$2,acknowledged_at=$3,version=version+1 WHERE command_id=$1 AND status='DELIVERED'`, v.CommandID, ack.Status, ack.ReceivedAt)
		if err == nil {
			_, err = tx.ExecContext(ctx, `UPDATE command_attempts SET completed_at=$2,result=$3 WHERE command_id=$1 AND attempt_number=$4`, v.CommandID, ack.ReceivedAt, ack.Status, v.AttemptCount)
		}
		if err == nil {
			_, err = tx.ExecContext(ctx, `UPDATE tasks SET status='IN_PROGRESS',started_at=COALESCE(started_at,$2),updated_at=NOW(),version=version+1 WHERE task_id=$1 AND status='ASSIGNED'`, v.TaskID, ack.ReceivedAt)
		}
	} else {
		next := string(command.Failed)
		taskNext := string(taskcore.Failed)
		if ack.Status == "EXPIRED" {
			next, taskNext = string(command.Expired), string(taskcore.Expired)
		}
		_, err = tx.ExecContext(ctx, `UPDATE commands SET status=$2,ack_status=$3,completed_at=$4,last_error=NULLIF($5,''),version=version+1 WHERE command_id=$1 AND status='DELIVERED'`, v.CommandID, next, ack.Status, ack.ReceivedAt, ack.Reason)
		if err == nil {
			_, err = tx.ExecContext(ctx, `UPDATE tasks SET status=$2,failed_at=$3,failure_reason=NULLIF($4,''),updated_at=NOW(),version=version+1 WHERE task_id=$1 AND status IN('ASSIGNED','IN_PROGRESS')`, v.TaskID, taskNext, ack.ReceivedAt, ack.Reason)
		}
		if err == nil {
			_, err = tx.ExecContext(ctx, `UPDATE device_assignments SET status='RELEASED',updated_at=NOW() WHERE task_id=$1 AND status='ACTIVE'`, v.TaskID)
		}
	}
	if err != nil {
		return err
	}
	action := "COMMAND_ACKNOWLEDGED"
	if !accepted {
		action = "COMMAND_REJECTED"
	}
	if err = insertAuditOutbox(ctx, tx, principalTenant, principalDevice, action, "command", v.CommandID, "", "command.acknowledged.v1", map[string]interface{}{"command_id": v.CommandID, "sequence_number": v.SequenceNumber, "status": ack.Status}); err != nil {
		return err
	}
	return tx.Commit()
}

func (s *RegistryStore) ApplyCommandResult(ctx context.Context, principalTenant, principalDevice string, result command.Result) error {
	tx, err := s.DB.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	var v command.Record
	err = tx.GetContext(ctx, &v, `SELECT * FROM commands WHERE tenant_id=$1 AND device_id=$2 AND command_id=$3 FOR UPDATE`, principalTenant, principalDevice, result.CommandID)
	if errors.Is(err, sql.ErrNoRows) {
		return ErrNotFound
	}
	if err != nil {
		return err
	}
	if result.SequenceNumber != v.SequenceNumber {
		return ErrForbidden
	}
	if command.IsTerminal(command.Status(v.Status)) {
		return tx.Commit()
	}
	if v.Status != string(command.Acknowledged) {
		return ErrInvalidTransition
	}
	if result.CompletedAt.IsZero() {
		result.CompletedAt = time.Now().UTC()
	}
	next, taskNext, action := string(command.Completed), string(taskcore.Completed), "COMMAND_COMPLETED"
	if result.Status != "SUCCEEDED" && result.Status != "COMPLETED" {
		next, taskNext, action = string(command.Failed), string(taskcore.Failed), "COMMAND_FAILED"
	}
	payload := result.Result
	if len(payload) == 0 {
		payload = []byte(`{}`)
	}
	_, err = tx.ExecContext(ctx, `UPDATE commands SET status=$2,result=$3,completed_at=$4,last_error=NULLIF($5,''),version=version+1 WHERE command_id=$1 AND status='ACKNOWLEDGED'`, v.CommandID, next, payload, result.CompletedAt, result.Reason)
	if err == nil {
		_, err = tx.ExecContext(ctx, `UPDATE tasks SET status=$2,completed_at=CASE WHEN $2='COMPLETED' THEN $3 ELSE completed_at END,failed_at=CASE WHEN $2='FAILED' THEN $3 ELSE failed_at END,failure_reason=NULLIF($4,''),updated_at=NOW(),version=version+1 WHERE task_id=$1 AND status='IN_PROGRESS'`, v.TaskID, taskNext, result.CompletedAt, result.Reason)
	}
	if err == nil {
		_, err = tx.ExecContext(ctx, `UPDATE device_assignments SET status='RELEASED',updated_at=NOW() WHERE task_id=$1 AND status='ACTIVE'`, v.TaskID)
	}
	if err != nil {
		return err
	}
	if err = insertAuditOutbox(ctx, tx, principalTenant, principalDevice, action, "command", v.CommandID, "", "command.result.v1", map[string]interface{}{"command_id": v.CommandID, "status": result.Status}); err != nil {
		return err
	}
	return tx.Commit()
}

func (s *RegistryStore) CancelTask(ctx context.Context, tenant, taskID, actor, requestID string) error {
	tx, err := s.DB.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	var status string
	if err = tx.GetContext(ctx, &status, `SELECT status FROM tasks WHERE tenant_id=$1 AND task_id=$2 FOR UPDATE`, tenant, taskID); errors.Is(err, sql.ErrNoRows) {
		return ErrNotFound
	}
	if err != nil {
		return err
	}
	if status == string(taskcore.Cancelled) {
		return tx.Commit()
	}
	var commandStatus sql.NullString
	_ = tx.GetContext(ctx, &commandStatus, `SELECT status FROM commands WHERE tenant_id=$1 AND task_id=$2 ORDER BY sequence_number DESC LIMIT 1`, tenant, taskID)
	if commandStatus.Valid && commandStatus.String != string(command.Pending) {
		return ErrInvalidTransition
	}
	result, err := tx.ExecContext(ctx, `UPDATE tasks SET status='CANCELLED',updated_at=NOW(),version=version+1 WHERE tenant_id=$1 AND task_id=$2 AND status IN('PENDING','ASSIGNING','ASSIGNED')`, tenant, taskID)
	if err != nil {
		return err
	}
	if n, _ := result.RowsAffected(); n != 1 {
		return ErrInvalidTransition
	}
	_, _ = tx.ExecContext(ctx, `UPDATE commands SET status='CANCELLED',completed_at=NOW(),version=version+1 WHERE tenant_id=$1 AND task_id=$2 AND status='PENDING'`, tenant, taskID)
	_, _ = tx.ExecContext(ctx, `UPDATE device_assignments SET status='RELEASED',updated_at=NOW() WHERE task_id=$1 AND status='ACTIVE'`, taskID)
	if err = insertAuditOutbox(ctx, tx, tenant, actor, "TASK_CANCELLED", "task", taskID, requestID, "task.cancelled.v1", map[string]string{"task_id": taskID}); err != nil {
		return err
	}
	return tx.Commit()
}

func (s *RegistryStore) RetryTask(ctx context.Context, tenant, taskID, actor, requestID string, expiresAt time.Time) error {
	tx, err := s.DB.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	result, err := tx.ExecContext(ctx, `UPDATE tasks SET status='PENDING',assigned_device_id=NULL,assigned_at=NULL,started_at=NULL,completed_at=NULL,failed_at=NULL,failure_reason=NULL,expires_at=$3,updated_at=NOW(),version=version+1 WHERE tenant_id=$1 AND task_id=$2 AND status IN('FAILED','EXPIRED')`, tenant, taskID, expiresAt)
	if err != nil {
		return err
	}
	if n, _ := result.RowsAffected(); n != 1 {
		return ErrInvalidTransition
	}
	_, _ = tx.ExecContext(ctx, `UPDATE device_assignments SET status='RELEASED',updated_at=NOW() WHERE task_id=$1 AND status='ACTIVE'`, taskID)
	if err = insertAuditOutbox(ctx, tx, tenant, actor, "TASK_RETRIED", "task", taskID, requestID, "task.retry.requested.v1", map[string]string{"task_id": taskID}); err != nil {
		return err
	}
	return tx.Commit()
}

func (s *RegistryStore) RetryCommand(ctx context.Context, tenant, commandID, actor, requestID string, force bool) error {
	tx, err := s.DB.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	var v command.Record
	if err = tx.GetContext(ctx, &v, `SELECT * FROM commands WHERE tenant_id=$1 AND command_id=$2 FOR UPDATE`, tenant, commandID); errors.Is(err, sql.ErrNoRows) {
		return ErrNotFound
	}
	if err != nil {
		return err
	}
	if v.Status != string(command.Delivered) && !(force && v.Status == string(command.Failed)) {
		return ErrInvalidTransition
	}
	if time.Now().After(v.ExpiresAt) || (!force && v.AttemptCount >= v.MaxAttempts) {
		return ErrInvalidTransition
	}
	if v.Status == string(command.Failed) {
		result, assignErr := tx.ExecContext(ctx, `UPDATE device_assignments a SET status='ACTIVE',lease_expires_at=$2,updated_at=NOW() WHERE task_id=$1 AND status IN('RELEASED','EXPIRED') AND NOT EXISTS(SELECT 1 FROM device_assignments x WHERE x.tenant_id=a.tenant_id AND x.device_id=a.device_id AND x.status='ACTIVE')`, v.TaskID, v.ExpiresAt)
		if assignErr != nil {
			return mapPQ(assignErr)
		}
		if n, _ := result.RowsAffected(); n != 1 {
			return ErrConflict
		}
		_, err = tx.ExecContext(ctx, `UPDATE tasks SET status='ASSIGNED',failure_reason=NULL,failed_at=NULL,updated_at=NOW(),version=version+1 WHERE task_id=$1 AND status='FAILED'`, v.TaskID)
		if err != nil {
			return err
		}
	}
	resetAttempts := v.AttemptCount
	if force && v.Status == string(command.Failed) {
		resetAttempts = 0
	}
	_, err = tx.ExecContext(ctx, `UPDATE commands SET status='PENDING',attempt_count=$2,available_at=NOW(),last_error=NULL,completed_at=NULL,version=version+1 WHERE command_id=$1`, commandID, resetAttempts)
	if err != nil {
		return err
	}
	v.Status = string(command.Pending)
	if err = insertAuditOutbox(ctx, tx, tenant, actor, "COMMAND_RETRIED", "command", commandID, requestID, "command.retry.requested.v1", v.Envelope()); err != nil {
		return err
	}
	return tx.Commit()
}

func (s *RegistryStore) ReconcileCommands(ctx context.Context, ackTimeout time.Duration) error {
	tx, err := s.DB.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	rows, err := tx.QueryxContext(ctx, `SELECT * FROM commands WHERE status='DELIVERED' AND sent_at < NOW()-$1::interval FOR UPDATE SKIP LOCKED`, ackTimeout.String())
	if err != nil {
		return err
	}
	due := []command.Record{}
	for rows.Next() {
		var v command.Record
		if err = rows.StructScan(&v); err != nil {
			rows.Close()
			return err
		}
		due = append(due, v)
	}
	rows.Close()
	for _, v := range due {
		if time.Now().After(v.ExpiresAt) || v.AttemptCount >= v.MaxAttempts {
			next := string(command.Failed)
			taskNext := string(taskcore.Failed)
			reason := "delivery attempts exhausted"
			if time.Now().After(v.ExpiresAt) {
				next, taskNext, reason = string(command.Expired), string(taskcore.Expired), "command expired"
			}
			_, err = tx.ExecContext(ctx, `UPDATE commands SET status=$2,completed_at=NOW(),last_error=$3,version=version+1 WHERE command_id=$1 AND status='DELIVERED'`, v.CommandID, next, reason)
			if err == nil {
				_, err = tx.ExecContext(ctx, `UPDATE tasks SET status=$2,failed_at=NOW(),failure_reason=$3,updated_at=NOW(),version=version+1 WHERE task_id=$1 AND status IN('ASSIGNED','IN_PROGRESS')`, v.TaskID, taskNext, reason)
			}
			if err == nil {
				_, err = tx.ExecContext(ctx, `UPDATE device_assignments SET status='RELEASED',updated_at=NOW() WHERE task_id=$1 AND status='ACTIVE'`, v.TaskID)
			}
			if err != nil {
				return err
			}
			continue
		}
		backoffSeconds := int(math.Min(math.Pow(2, float64(max(v.AttemptCount-1, 0))), 30))
		_, err = tx.ExecContext(ctx, `UPDATE commands SET status='PENDING',available_at=NOW()+$2::interval,last_error='ack timeout',version=version+1 WHERE command_id=$1 AND status='DELIVERED'`, v.CommandID, (time.Duration(backoffSeconds) * time.Second).String())
		if err != nil {
			return err
		}
		v.Status = string(command.Pending)
		if err = insertAuditOutbox(ctx, tx, v.TenantID, "reconciler", "COMMAND_RETRIED", "command", v.CommandID, "", "command.retry.requested.v1", v.Envelope()); err != nil {
			return err
		}
	}
	_, err = tx.ExecContext(ctx, `WITH expired AS (UPDATE commands SET status='EXPIRED',completed_at=NOW(),last_error='command expired',version=version+1 WHERE status='PENDING' AND expires_at<=NOW() RETURNING task_id) UPDATE tasks SET status='EXPIRED',failed_at=NOW(),failure_reason='command expired',updated_at=NOW(),version=version+1 WHERE task_id IN(SELECT task_id FROM expired) AND status IN('ASSIGNED','IN_PROGRESS')`)
	if err != nil {
		return err
	}
	_, err = tx.ExecContext(ctx, `UPDATE device_assignments a SET status='EXPIRED',updated_at=NOW() FROM tasks t WHERE a.task_id=t.task_id AND a.status='ACTIVE' AND t.status IN('EXPIRED','FAILED','CANCELLED','COMPLETED')`)
	return tx.Commit()
}

func (s *RegistryStore) FailExpiredPendingTasks(ctx context.Context) error {
	_, err := s.DB.ExecContext(ctx, `UPDATE tasks SET status='EXPIRED',failed_at=NOW(),failure_reason='task expired before assignment',updated_at=NOW(),version=version+1 WHERE status IN('PENDING','ASSIGNING') AND expires_at<=NOW()`)
	return err
}
```

---

## internal\adapter\repository\postgres_registry.go

```
package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/Akashpg-M/polaris/backend/internal/core/auth"
	"github.com/Akashpg-M/polaris/backend/internal/core/registry"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

var (
	ErrNotFound          = errors.New("registry resource not found")
	ErrConflict          = errors.New("registry resource already exists")
	ErrForbidden         = errors.New("registry operation forbidden")
	ErrInvalidTransition = errors.New("invalid lifecycle transition")
)

type RegistryStore struct{ DB *sqlx.DB }

func NewRegistryStore(postgresURL string) (*RegistryStore, error) {
	db, err := sqlx.Connect("postgres", postgresURL)
	if err != nil {
		return nil, err
	}
	// A process may host many concurrent device sessions, but PostgreSQL must
	// see a bounded pool rather than one backend per socket. The gateway and
	// engine each own one RegistryStore, so these limits also leave capacity
	// for archival, migrations, and operator diagnostics.
	db.SetMaxOpenConns(20)
	db.SetMaxIdleConns(10)
	db.SetConnMaxIdleTime(5 * time.Minute)
	db.SetConnMaxLifetime(30 * time.Minute)
	return &RegistryStore{DB: db}, nil
}
func (s *RegistryStore) Close() error { return s.DB.Close() }
func jsonOrEmpty(v json.RawMessage) []byte {
	if len(v) == 0 {
		return []byte(`{}`)
	}
	return v
}
func nullableJSON(v []byte) interface{} {
	if len(v) == 0 {
		return nil
	}
	return v
}
func mapPQ(err error) error {
	var p *pq.Error
	if errors.As(err, &p) && p.Code.Class() == "23" {
		if p.Code == "23505" {
			return ErrConflict
		}
		return fmt.Errorf("registry constraint: %w", err)
	}
	return err
}

func insertAuditOutbox(ctx context.Context, tx *sqlx.Tx, tenant, actorID, action, resourceType, resourceID, requestID, eventType string, payload interface{}) error {
	b, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	eventID := auth.NewID()
	if _, err = tx.ExecContext(ctx, `INSERT INTO audit_events(audit_id,tenant_id,actor_type,actor_id,action,resource_type,resource_id,request_id,outcome) VALUES($1,NULLIF($2,''),'OPERATOR',$3,$4,$5,$6,$7,'SUCCESS')`, auth.NewID(), tenant, actorID, action, resourceType, resourceID, requestID); err != nil {
		return err
	}
	_, err = tx.ExecContext(ctx, `INSERT INTO outbox_events(outbox_id,aggregate_type,aggregate_id,tenant_id,event_id,event_type,schema_version,payload,status) VALUES($1,$2,$3,$4,$5,$6,1,$7,'PENDING')`, auth.NewID(), resourceType, resourceID, tenant, eventID, eventType, b)
	return err
}

func (s *RegistryStore) BootstrapPlatformAdmin(ctx context.Context, raw string) error {
	if raw == "" {
		return nil
	}
	prefix, err := auth.TokenPrefix(raw)
	if err != nil {
		return err
	}
	_, err = s.DB.ExecContext(ctx, `INSERT INTO operator_api_keys(api_key_id,name,token_prefix,token_hash,role,status) VALUES($1,'development bootstrap',$2,$3,'PLATFORM_ADMIN','ACTIVE') ON CONFLICT(token_prefix) DO UPDATE SET token_hash=EXCLUDED.token_hash,status='ACTIVE',revoked_at=NULL`, auth.NewID(), prefix, auth.Hash(raw))
	return err
}

func (s *RegistryStore) ResolveOperator(ctx context.Context, raw string) (auth.OperatorPrincipal, error) {
	prefix, err := auth.TokenPrefix(raw)
	if err != nil {
		return auth.OperatorPrincipal{}, auth.ErrInvalidCredential
	}
	var row struct {
		ID      string         `db:"api_key_id"`
		Tenant  sql.NullString `db:"tenant_id"`
		Role    string         `db:"role"`
		Hash    []byte         `db:"token_hash"`
		Status  string         `db:"status"`
		Expires sql.NullTime   `db:"expires_at"`
	}
	err = s.DB.GetContext(ctx, &row, `SELECT api_key_id,tenant_id,role,token_hash,status,expires_at FROM operator_api_keys WHERE token_prefix=$1`, prefix)
	if err != nil || row.Status != "ACTIVE" || (row.Expires.Valid && time.Now().After(row.Expires.Time)) || !auth.Verify(raw, row.Hash) {
		return auth.OperatorPrincipal{}, auth.ErrInvalidCredential
	}
	_, _ = s.DB.ExecContext(ctx, `UPDATE operator_api_keys SET last_used_at=NOW() WHERE api_key_id=$1`, row.ID)
	return auth.OperatorPrincipal{APIKeyID: row.ID, TenantID: row.Tenant.String, Role: auth.Role(row.Role)}, nil
}

func (s *RegistryStore) ResolveDevice(ctx context.Context, raw string) (auth.DevicePrincipal, error) {
	prefix, err := auth.TokenPrefix(raw)
	if err != nil {
		return auth.DevicePrincipal{}, auth.ErrInvalidCredential
	}
	var row struct {
		CredentialID, TenantID, DeviceID, DeviceType, ProjectID, CredentialStatus, TenantStatus, DeviceStatus string
		Hash                                                                                                  []byte
		Expires                                                                                               sql.NullTime
	}
	err = s.DB.QueryRowxContext(ctx, `SELECT c.credential_id,c.tenant_id,c.device_id,d.device_type_id,COALESCE(d.project_id::text,''),c.status,c.token_hash,c.expires_at,t.status,d.lifecycle_status FROM device_credentials c JOIN devices d USING(tenant_id,device_id) JOIN tenants t USING(tenant_id) WHERE c.token_prefix=$1`, prefix).Scan(&row.CredentialID, &row.TenantID, &row.DeviceID, &row.DeviceType, &row.ProjectID, &row.CredentialStatus, &row.Hash, &row.Expires, &row.TenantStatus, &row.DeviceStatus)
	if err != nil || row.CredentialStatus != "ACTIVE" || row.TenantStatus != "ACTIVE" || row.DeviceStatus != "ACTIVE" || (row.Expires.Valid && time.Now().After(row.Expires.Time)) || !auth.Verify(raw, row.Hash) {
		return auth.DevicePrincipal{}, auth.ErrInvalidCredential
	}
	_, _ = s.DB.ExecContext(ctx, `UPDATE device_credentials SET last_used_at=NOW() WHERE credential_id=$1`, row.CredentialID)
	return auth.DevicePrincipal{TenantID: row.TenantID, DeviceID: row.DeviceID, CredentialID: row.CredentialID, DeviceType: row.DeviceType, ProjectID: row.ProjectID}, nil
}
func (s *RegistryStore) RevalidateDevice(ctx context.Context, p auth.DevicePrincipal) error {
	var ok bool
	err := s.DB.GetContext(ctx, &ok, `SELECT c.status='ACTIVE' AND (c.expires_at IS NULL OR c.expires_at>NOW()) AND d.lifecycle_status='ACTIVE' AND t.status='ACTIVE' FROM device_credentials c JOIN devices d USING(tenant_id,device_id) JOIN tenants t USING(tenant_id) WHERE c.credential_id=$1 AND c.tenant_id=$2 AND c.device_id=$3`, p.CredentialID, p.TenantID, p.DeviceID)
	if err != nil {
		return err
	}
	if !ok {
		return auth.ErrInvalidCredential
	}
	return nil
}

func (s *RegistryStore) CreateTenant(ctx context.Context, t registry.Tenant, actor, requestID string) error {
	tx, err := s.DB.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	_, err = tx.ExecContext(ctx, `INSERT INTO tenants(tenant_id,display_name,status,metadata) VALUES($1,$2,$3,$4)`, t.TenantID, t.DisplayName, t.Status, t.Metadata)
	if err != nil {
		return mapPQ(err)
	}
	if err = insertAuditOutbox(ctx, tx, t.TenantID, actor, "TENANT_CREATED", "tenant", t.TenantID, requestID, "tenant.registered.v1", t); err != nil {
		return err
	}
	return tx.Commit()
}
func (s *RegistryStore) GetTenant(ctx context.Context, id string) (registry.Tenant, error) {
	var v registry.Tenant
	err := s.DB.GetContext(ctx, &v, `SELECT * FROM tenants WHERE tenant_id=$1`, id)
	if errors.Is(err, sql.ErrNoRows) {
		err = ErrNotFound
	}
	return v, err
}
func (s *RegistryStore) SetTenantStatus(ctx context.Context, id, status, actor, requestID string) error {
	if status != "ACTIVE" && status != "SUSPENDED" && status != "DEACTIVATED" {
		return ErrInvalidTransition
	}
	tx, err := s.DB.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	result, err := tx.ExecContext(ctx, `UPDATE tenants SET status=$2,updated_at=NOW() WHERE tenant_id=$1`, id, status)
	if err != nil {
		return err
	}
	if n, _ := result.RowsAffected(); n != 1 {
		return ErrNotFound
	}
	if err = insertAuditOutbox(ctx, tx, id, actor, "TENANT_"+status, "tenant", id, requestID, "tenant.lifecycle.changed.v1", map[string]string{"tenant_id": id, "status": status}); err != nil {
		return err
	}
	return tx.Commit()
}

func (s *RegistryStore) CreateProject(ctx context.Context, p registry.Project, actor, requestID string) error {
	tx, err := s.DB.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	_, err = tx.ExecContext(ctx, `INSERT INTO projects(project_id,tenant_id,name,description,status,metadata) VALUES($1,$2,$3,$4,$5,$6)`, p.ProjectID, p.TenantID, p.Name, p.Description, p.Status, p.Metadata)
	if err != nil {
		return mapPQ(err)
	}
	if err = insertAuditOutbox(ctx, tx, p.TenantID, actor, "PROJECT_CREATED", "project", p.ProjectID, requestID, "project.registered.v1", p); err != nil {
		return err
	}
	return tx.Commit()
}
func (s *RegistryStore) ListProjects(ctx context.Context, tenant string) ([]registry.Project, error) {
	v := []registry.Project{}
	err := s.DB.SelectContext(ctx, &v, `SELECT * FROM projects WHERE tenant_id=$1 ORDER BY created_at DESC LIMIT 100`, tenant)
	return v, err
}
func (s *RegistryStore) GetProject(ctx context.Context, tenant, id string) (registry.Project, error) {
	var v registry.Project
	err := s.DB.GetContext(ctx, &v, `SELECT * FROM projects WHERE tenant_id=$1 AND project_id=$2`, tenant, id)
	if errors.Is(err, sql.ErrNoRows) {
		err = ErrNotFound
	}
	return v, err
}
func (s *RegistryStore) UpdateProject(ctx context.Context, tenant, id string, name, description, status *string, metadata []byte, actor, requestID string) error {
	if status != nil && *status != "ACTIVE" && *status != "ARCHIVED" {
		return ErrInvalidTransition
	}
	tx, err := s.DB.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	result, err := tx.ExecContext(ctx, `UPDATE projects SET name=COALESCE($3,name),description=COALESCE($4,description),status=COALESCE($5,status),metadata=CASE WHEN $6::jsonb IS NULL THEN metadata ELSE $6::jsonb END,updated_at=NOW() WHERE tenant_id=$1 AND project_id=$2`, tenant, id, name, description, status, nullableJSON(metadata))
	if err != nil {
		return mapPQ(err)
	}
	if n, _ := result.RowsAffected(); n != 1 {
		return ErrNotFound
	}
	if err = insertAuditOutbox(ctx, tx, tenant, actor, "PROJECT_UPDATED", "project", id, requestID, "project.updated.v1", map[string]string{"project_id": id}); err != nil {
		return err
	}
	return tx.Commit()
}

func (s *RegistryStore) CreateDevice(ctx context.Context, d registry.Device, actor, requestID string) error {
	tx, err := s.DB.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	_, err = tx.ExecContext(ctx, `INSERT INTO devices(tenant_id,device_id,project_id,device_type_id,display_name,lifecycle_status,firmware_version,software_version,model_version,metadata) VALUES($1,$2,$3,$4,$5,$6,$7,$8,$9,$10)`, d.TenantID, d.DeviceID, d.ProjectID, d.DeviceTypeID, d.DisplayName, d.LifecycleStatus, d.FirmwareVersion, d.SoftwareVersion, d.ModelVersion, d.Metadata)
	if err != nil {
		return mapPQ(err)
	}
	if err = insertAuditOutbox(ctx, tx, d.TenantID, actor, "DEVICE_REGISTERED", "device", d.DeviceID, requestID, "device.registered.v1", d); err != nil {
		return err
	}
	return tx.Commit()
}
func (s *RegistryStore) GetDevice(ctx context.Context, tenant, id string) (registry.Device, error) {
	var v registry.Device
	err := s.DB.GetContext(ctx, &v, `SELECT * FROM devices WHERE tenant_id=$1 AND device_id=$2`, tenant, id)
	if errors.Is(err, sql.ErrNoRows) {
		err = ErrNotFound
	}
	return v, err
}
func (s *RegistryStore) UpdateDevice(ctx context.Context, tenant, id string, displayName, firmware, software, model *string, metadata []byte, actor, requestID string) error {
	tx, err := s.DB.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	result, err := tx.ExecContext(ctx, `UPDATE devices SET display_name=COALESCE($3,display_name),firmware_version=COALESCE($4,firmware_version),software_version=COALESCE($5,software_version),model_version=COALESCE($6,model_version),metadata=CASE WHEN $7::jsonb IS NULL THEN metadata ELSE $7::jsonb END,updated_at=NOW() WHERE tenant_id=$1 AND device_id=$2 AND lifecycle_status<>'DECOMMISSIONED'`, tenant, id, displayName, firmware, software, model, nullableJSON(metadata))
	if err != nil {
		return err
	}
	if n, _ := result.RowsAffected(); n != 1 {
		return ErrNotFound
	}
	if err = insertAuditOutbox(ctx, tx, tenant, actor, "DEVICE_UPDATED", "device", id, requestID, "device.updated.v1", map[string]string{"device_id": id}); err != nil {
		return err
	}
	return tx.Commit()
}
func (s *RegistryStore) ListDevices(ctx context.Context, tenant string, limit int, cursor, projectID, deviceType, lifecycle, capability string) ([]registry.Device, error) {
	if limit < 1 || limit > 100 {
		limit = 50
	}
	args := []interface{}{tenant, limit}
	q := `SELECT d.* FROM devices d WHERE d.tenant_id=$1`
	add := func(clause string, v interface{}) { args = append(args, v); q += fmt.Sprintf(clause, len(args)) }
	if cursor != "" {
		add(` AND d.device_id>$%d`, cursor)
	}
	if projectID != "" {
		add(` AND d.project_id=$%d`, projectID)
	}
	if deviceType != "" {
		add(` AND d.device_type_id=$%d`, deviceType)
	}
	if lifecycle != "" {
		add(` AND d.lifecycle_status=$%d`, lifecycle)
	}
	if capability != "" {
		add(` AND EXISTS(SELECT 1 FROM device_capabilities dc WHERE dc.tenant_id=d.tenant_id AND dc.device_id=d.device_id AND dc.capability_id=$%d AND dc.enabled)`, capability)
	}
	q += ` ORDER BY d.device_id LIMIT $2`
	v := []registry.Device{}
	err := s.DB.SelectContext(ctx, &v, q, args...)
	return v, err
}

func (s *RegistryStore) SetLifecycle(ctx context.Context, tenant, id, next, actor, requestID string) error {
	tx, err := s.DB.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	var current string
	if err = tx.GetContext(ctx, &current, `SELECT lifecycle_status FROM devices WHERE tenant_id=$1 AND device_id=$2 FOR UPDATE`, tenant, id); errors.Is(err, sql.ErrNoRows) {
		return ErrNotFound
	}
	if err != nil {
		return err
	}
	if !registry.ValidTransition(current, next) {
		return ErrInvalidTransition
	}
	result, err := tx.ExecContext(ctx, `UPDATE devices SET lifecycle_status=$3,updated_at=NOW(),deactivated_at=CASE WHEN $3='DECOMMISSIONED' THEN NOW() ELSE deactivated_at END WHERE tenant_id=$1 AND device_id=$2`, tenant, id, next)
	if err != nil {
		return err
	}
	if n, _ := result.RowsAffected(); n != 1 {
		return ErrNotFound
	}
	payload := map[string]string{"device_id": id, "previous_status": current, "status": next}
	if err = insertAuditOutbox(ctx, tx, tenant, actor, "DEVICE_"+next, "device", id, requestID, "device.lifecycle.changed.v1", payload); err != nil {
		return err
	}
	return tx.Commit()
}

func (s *RegistryStore) ListCapabilities(ctx context.Context, tenant, id string) ([]registry.Capability, error) {
	v := []registry.Capability{}
	err := s.DB.SelectContext(ctx, &v, `SELECT c.capability_id,c.display_name,c.description,dc.configuration,dc.enabled FROM device_capabilities dc JOIN capabilities c USING(capability_id) WHERE dc.tenant_id=$1 AND dc.device_id=$2 ORDER BY c.capability_id`, tenant, id)
	return v, err
}
func (s *RegistryStore) AllCapabilities(ctx context.Context) ([]registry.Capability, error) {
	v := []registry.Capability{}
	err := s.DB.SelectContext(ctx, &v, `SELECT capability_id,display_name,description,'{}'::jsonb configuration,true enabled FROM capabilities ORDER BY capability_id`)
	return v, err
}
func (s *RegistryStore) PutCapability(ctx context.Context, tenant, id, capability string, configuration []byte, actor, requestID string) error {
	tx, err := s.DB.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	_, err = tx.ExecContext(ctx, `INSERT INTO device_capabilities(tenant_id,device_id,capability_id,configuration,enabled) VALUES($1,$2,$3,$4,true) ON CONFLICT(tenant_id,device_id,capability_id) DO UPDATE SET configuration=EXCLUDED.configuration,enabled=true`, tenant, id, capability, configuration)
	if err != nil {
		return mapPQ(err)
	}
	if err = insertAuditOutbox(ctx, tx, tenant, actor, "CAPABILITY_ADDED", "device", id, requestID, "device.capabilities.changed.v1", map[string]string{"device_id": id, "capability_id": capability}); err != nil {
		return err
	}
	return tx.Commit()
}
func (s *RegistryStore) RemoveCapability(ctx context.Context, tenant, id, capability, actor, requestID string) error {
	tx, err := s.DB.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	result, err := tx.ExecContext(ctx, `DELETE FROM device_capabilities WHERE tenant_id=$1 AND device_id=$2 AND capability_id=$3`, tenant, id, capability)
	if err != nil {
		return err
	}
	if n, _ := result.RowsAffected(); n != 1 {
		return ErrNotFound
	}
	if err = insertAuditOutbox(ctx, tx, tenant, actor, "CAPABILITY_REMOVED", "device", id, requestID, "device.capabilities.changed.v1", map[string]string{"device_id": id, "capability_id": capability, "enabled": "false"}); err != nil {
		return err
	}
	return tx.Commit()
}

func (s *RegistryStore) IssueCredential(ctx context.Context, tenant, id, actor, requestID string, expires *time.Time) (registry.CredentialMetadata, string, error) {
	raw, prefix, hash, err := auth.GenerateToken("dev")
	if err != nil {
		return registry.CredentialMetadata{}, "", err
	}
	credentialID := auth.NewID()
	tx, err := s.DB.BeginTxx(ctx, nil)
	if err != nil {
		return registry.CredentialMetadata{}, "", err
	}
	defer tx.Rollback()
	result, err := tx.ExecContext(ctx, `INSERT INTO device_credentials(credential_id,tenant_id,device_id,token_prefix,token_hash,status,expires_at) SELECT $1,$2,$3,$4,$5,'ACTIVE',$6 WHERE EXISTS(SELECT 1 FROM devices WHERE tenant_id=$2 AND device_id=$3 AND lifecycle_status<>'DECOMMISSIONED')`, credentialID, tenant, id, prefix, hash, expires)
	if err != nil {
		return registry.CredentialMetadata{}, "", mapPQ(err)
	}
	if n, _ := result.RowsAffected(); n != 1 {
		return registry.CredentialMetadata{}, "", ErrNotFound
	}
	if err = insertAuditOutbox(ctx, tx, tenant, actor, "CREDENTIAL_CREATED", "device", id, requestID, "device.credential.created.v1", map[string]string{"device_id": id, "credential_id": credentialID}); err != nil {
		return registry.CredentialMetadata{}, "", err
	}
	if err = tx.Commit(); err != nil {
		return registry.CredentialMetadata{}, "", err
	}
	return registry.CredentialMetadata{CredentialID: credentialID, TokenPrefix: prefix, Status: "ACTIVE", IssuedAt: time.Now().UTC(), ExpiresAt: expires}, raw, nil
}
func (s *RegistryStore) ListCredentials(ctx context.Context, tenant, id string) ([]registry.CredentialMetadata, error) {
	v := []registry.CredentialMetadata{}
	err := s.DB.SelectContext(ctx, &v, `SELECT credential_id,token_prefix,status,issued_at,expires_at,last_used_at,revoked_at FROM device_credentials WHERE tenant_id=$1 AND device_id=$2 ORDER BY issued_at DESC`, tenant, id)
	return v, err
}
func (s *RegistryStore) RevokeCredential(ctx context.Context, tenant, id, credentialID, actor, requestID string) error {
	tx, err := s.DB.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	result, err := tx.ExecContext(ctx, `UPDATE device_credentials SET status='REVOKED',revoked_at=NOW() WHERE credential_id=$1 AND tenant_id=$2 AND device_id=$3 AND status='ACTIVE'`, credentialID, tenant, id)
	if err != nil {
		return err
	}
	if n, _ := result.RowsAffected(); n != 1 {
		return ErrNotFound
	}
	if err = insertAuditOutbox(ctx, tx, tenant, actor, "CREDENTIAL_REVOKED", "device", id, requestID, "device.credential.revoked.v1", map[string]string{"device_id": id, "credential_id": credentialID}); err != nil {
		return err
	}
	return tx.Commit()
}
func (s *RegistryStore) RotateCredential(ctx context.Context, tenant, id, oldID, actor, requestID string) (registry.CredentialMetadata, string, error) {
	raw, prefix, hash, err := auth.GenerateToken("dev")
	if err != nil {
		return registry.CredentialMetadata{}, "", err
	}
	newID := auth.NewID()
	tx, err := s.DB.BeginTxx(ctx, nil)
	if err != nil {
		return registry.CredentialMetadata{}, "", err
	}
	defer tx.Rollback()
	result, err := tx.ExecContext(ctx, `UPDATE device_credentials SET status='REVOKED',revoked_at=NOW() WHERE credential_id=$1 AND tenant_id=$2 AND device_id=$3 AND status='ACTIVE'`, oldID, tenant, id)
	if err != nil {
		return registry.CredentialMetadata{}, "", err
	}
	if n, _ := result.RowsAffected(); n != 1 {
		return registry.CredentialMetadata{}, "", ErrNotFound
	}
	if _, err = tx.ExecContext(ctx, `INSERT INTO device_credentials(credential_id,tenant_id,device_id,token_prefix,token_hash,status) VALUES($1,$2,$3,$4,$5,'ACTIVE')`, newID, tenant, id, prefix, hash); err != nil {
		return registry.CredentialMetadata{}, "", err
	}
	if err = insertAuditOutbox(ctx, tx, tenant, actor, "CREDENTIAL_ROTATED", "device", id, requestID, "device.credential.rotated.v1", map[string]string{"device_id": id, "old_credential_id": oldID, "credential_id": newID}); err != nil {
		return registry.CredentialMetadata{}, "", err
	}
	if err = tx.Commit(); err != nil {
		return registry.CredentialMetadata{}, "", err
	}
	return registry.CredentialMetadata{CredentialID: newID, TokenPrefix: prefix, Status: "ACTIVE", IssuedAt: time.Now().UTC()}, raw, nil
}

func (s *RegistryStore) CreateTicket(ctx context.Context, tenant, id, credentialID string, ttl time.Duration) (string, error) {
	raw, prefix, hash, err := auth.GenerateToken("ticket")
	if err != nil {
		return "", err
	}
	_, err = s.DB.ExecContext(ctx, `INSERT INTO connection_tickets(ticket_prefix,ticket_hash,tenant_id,device_id,credential_id,expires_at) SELECT $1,$2,$3,$4,$5,NOW()+$6::interval WHERE EXISTS(SELECT 1 FROM devices WHERE tenant_id=$3 AND device_id=$4 AND lifecycle_status='ACTIVE')`, prefix, hash, tenant, id, credentialID, ttl.String())
	return raw, err
}
func (s *RegistryStore) ConsumeTicket(ctx context.Context, raw string) (auth.DevicePrincipal, error) {
	prefix, err := auth.TokenPrefix(raw)
	if err != nil {
		return auth.DevicePrincipal{}, auth.ErrInvalidCredential
	}
	tx, err := s.DB.BeginTxx(ctx, nil)
	if err != nil {
		return auth.DevicePrincipal{}, err
	}
	defer tx.Rollback()
	var p auth.DevicePrincipal
	var hash []byte
	var status string
	err = tx.QueryRowxContext(ctx, `SELECT x.tenant_id,x.device_id,x.credential_id,d.device_type_id,COALESCE(d.project_id::text,''),x.ticket_hash,d.lifecycle_status FROM connection_tickets x JOIN devices d USING(tenant_id,device_id) JOIN tenants t USING(tenant_id) WHERE x.ticket_prefix=$1 AND x.consumed_at IS NULL AND x.expires_at>NOW() AND t.status='ACTIVE' FOR UPDATE OF x`, prefix).Scan(&p.TenantID, &p.DeviceID, &p.CredentialID, &p.DeviceType, &p.ProjectID, &hash, &status)
	if err != nil || status != "ACTIVE" || !auth.Verify(raw, hash) {
		return auth.DevicePrincipal{}, auth.ErrInvalidCredential
	}
	if _, err = tx.ExecContext(ctx, `UPDATE connection_tickets SET consumed_at=NOW() WHERE ticket_prefix=$1`, prefix); err != nil {
		return auth.DevicePrincipal{}, err
	}
	if err = tx.Commit(); err != nil {
		return auth.DevicePrincipal{}, err
	}
	return p, nil
}
func (s *RegistryStore) CreateOperatorTicket(ctx context.Context, p auth.OperatorPrincipal, tenant string, ttl time.Duration) (string, error) {
	if p.Role != auth.PlatformAdmin {
		tenant = p.TenantID
	}
	if tenant == "" {
		return "", ErrForbidden
	}
	raw, prefix, hash, err := auth.GenerateToken("ticket")
	if err != nil {
		return "", err
	}
	_, err = s.DB.ExecContext(ctx, `INSERT INTO operator_connection_tickets(ticket_prefix,ticket_hash,api_key_id,tenant_id,role,expires_at) VALUES($1,$2,$3,$4,$5,NOW()+$6::interval)`, prefix, hash, p.APIKeyID, tenant, p.Role, ttl.String())
	return raw, err
}
func (s *RegistryStore) ConsumeOperatorTicket(ctx context.Context, raw string) (auth.OperatorPrincipal, error) {
	prefix, err := auth.TokenPrefix(raw)
	if err != nil {
		return auth.OperatorPrincipal{}, auth.ErrInvalidCredential
	}
	tx, err := s.DB.BeginTxx(ctx, nil)
	if err != nil {
		return auth.OperatorPrincipal{}, err
	}
	defer tx.Rollback()
	var p auth.OperatorPrincipal
	var hash []byte
	err = tx.QueryRowxContext(ctx, `SELECT api_key_id,tenant_id,role,ticket_hash FROM operator_connection_tickets WHERE ticket_prefix=$1 AND consumed_at IS NULL AND expires_at>NOW() FOR UPDATE`, prefix).Scan(&p.APIKeyID, &p.TenantID, &p.Role, &hash)
	if err != nil || !auth.Verify(raw, hash) {
		return auth.OperatorPrincipal{}, auth.ErrInvalidCredential
	}
	if _, err = tx.ExecContext(ctx, `UPDATE operator_connection_tickets SET consumed_at=NOW() WHERE ticket_prefix=$1`, prefix); err != nil {
		return auth.OperatorPrincipal{}, err
	}
	if err = tx.Commit(); err != nil {
		return auth.OperatorPrincipal{}, err
	}
	return p, nil
}

func (s *RegistryStore) Audit(ctx context.Context, tenant, actorID, action, resource, id, requestID, outcome string) error {
	_, err := s.DB.ExecContext(ctx, `INSERT INTO audit_events(audit_id,tenant_id,actor_type,actor_id,action,resource_type,resource_id,request_id,outcome) VALUES($1,NULLIF($2,''),'OPERATOR',$3,$4,$5,$6,$7,$8)`, auth.NewID(), tenant, actorID, action, resource, id, requestID, outcome)
	return err
}
func (s *RegistryStore) ListAudit(ctx context.Context, tenant string) ([]map[string]interface{}, error) {
	rows, err := s.DB.QueryxContext(ctx, `SELECT audit_id,tenant_id,actor_type,actor_id,action,resource_type,resource_id,request_id,outcome,created_at FROM audit_events WHERE tenant_id=$1 ORDER BY created_at DESC LIMIT 100`, tenant)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := []map[string]interface{}{}
	for rows.Next() {
		m := map[string]interface{}{}
		if err = rows.MapScan(m); err != nil {
			return nil, err
		}
		out = append(out, m)
	}
	return out, rows.Err()
}

type OutboxEvent struct {
	OutboxID, EventID, EventType, TenantID string
	Payload                                []byte
}

func (s *RegistryStore) ClaimOutbox(ctx context.Context, limit int) ([]OutboxEvent, error) {
	if limit < 1 {
		limit = 100
	}
	tx, err := s.DB.BeginTxx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()
	rows, err := tx.QueryxContext(ctx, `SELECT outbox_id,event_id,event_type,tenant_id,payload FROM outbox_events WHERE status IN('PENDING','RETRY_PENDING') AND next_attempt_at<=NOW() ORDER BY CASE WHEN event_type IN('command.created.v1','command.retry.requested.v1') THEN 0 ELSE 1 END,created_at LIMIT $1 FOR UPDATE SKIP LOCKED`, limit)
	if err != nil {
		return nil, err
	}
	events := []OutboxEvent{}
	for rows.Next() {
		var e OutboxEvent
		if err = rows.Scan(&e.OutboxID, &e.EventID, &e.EventType, &e.TenantID, &e.Payload); err != nil {
			return nil, err
		}
		events = append(events, e)
	}
	rows.Close()
	if len(events) > 0 {
		ids := make([]string, len(events))
		for i, e := range events {
			ids[i] = e.OutboxID
		}
		_, err = tx.ExecContext(ctx, `UPDATE outbox_events SET status='RETRY_PENDING',attempt_count=attempt_count+1,next_attempt_at=NOW()+INTERVAL '5 seconds' WHERE outbox_id=ANY($1)`, pq.Array(ids))
		if err != nil {
			return nil, err
		}
	}
	return events, tx.Commit()
}
func (s *RegistryStore) MarkOutboxPublishedBatch(ctx context.Context, ids []string) error {
	if len(ids) == 0 {
		return nil
	}
	_, err := s.DB.ExecContext(ctx, `UPDATE outbox_events SET status='PUBLISHED',published_at=NOW(),last_error=NULL WHERE outbox_id=ANY($1)`, pq.Array(ids))
	return err
}
func (s *RegistryStore) MarkOutboxFailed(ctx context.Context, id string, errText string) error {
	if len(errText) > 500 {
		errText = errText[:500]
	}
	_, err := s.DB.ExecContext(ctx, `UPDATE outbox_events SET status=CASE WHEN attempt_count>=10 THEN 'FAILED' ELSE 'RETRY_PENDING' END,last_error=$2,next_attempt_at=NOW()+INTERVAL '5 seconds' WHERE outbox_id=$1`, id, strings.ReplaceAll(errText, "\n", " "))
	return err
}
```

---

## internal\adapter\repository\redis_connection.go

```
package repository

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
)

type ConnectionOwnership struct {
	TenantID, DeviceID, GatewayID, ConnectionID, CredentialID string
	Epoch                                                     int64
	ConnectedAt, LeaseExpiresAt                               time.Time
}

type ConnectionOwnershipStore struct {
	client *redis.Client
	ttl    time.Duration
}

func NewConnectionOwnershipStore(client *redis.Client, ttl time.Duration) *ConnectionOwnershipStore {
	if ttl < 3*time.Second {
		ttl = 30 * time.Second
	}
	return &ConnectionOwnershipStore{client: client, ttl: ttl}
}

var claimOwnershipScript = redis.NewScript(`
local epoch = redis.call('INCR', KEYS[2])
redis.call('HSET', KEYS[1], 'gateway_id', ARGV[1], 'connection_id', ARGV[2], 'credential_id', ARGV[3], 'epoch', epoch, 'connected_at', ARGV[4], 'lease_expires_at', ARGV[5])
redis.call('PEXPIRE', KEYS[1], ARGV[6])
redis.call('SADD', KEYS[3], ARGV[7])
redis.call('ZADD', KEYS[4], ARGV[5], ARGV[7])
return epoch
`)

var refreshOwnershipScript = redis.NewScript(`
if redis.call('HGET', KEYS[1], 'gateway_id') ~= ARGV[1] or redis.call('HGET', KEYS[1], 'connection_id') ~= ARGV[2] or redis.call('HGET', KEYS[1], 'epoch') ~= ARGV[3] then return 0 end
redis.call('HSET', KEYS[1], 'lease_expires_at', ARGV[4])
redis.call('PEXPIRE', KEYS[1], ARGV[5])
redis.call('ZADD', KEYS[2], ARGV[4], ARGV[6])
return 1
`)

var releaseOwnershipScript = redis.NewScript(`
if redis.call('HGET', KEYS[1], 'gateway_id') ~= ARGV[1] or redis.call('HGET', KEYS[1], 'connection_id') ~= ARGV[2] or redis.call('HGET', KEYS[1], 'epoch') ~= ARGV[3] then return 0 end
redis.call('DEL', KEYS[1])
redis.call('SREM', KEYS[2], ARGV[4])
redis.call('ZREM', KEYS[3], ARGV[4])
return 1
`)

func ownershipKey(tenant, device string) string { return "polaris:connection:" + tenant + ":" + device }
func epochKey(tenant, device string) string {
	return "polaris:connection-epoch:" + tenant + ":" + device
}
func connectionMember(tenant, device string) string { return tenant + ":" + device }
func GatewayCommandChannel(gatewayID string) string {
	return "polaris:gateway:" + gatewayID + ":commands"
}

func (s *ConnectionOwnershipStore) Claim(ctx context.Context, tenant, device, gateway, connection, credential string) (ConnectionOwnership, error) {
	now := time.Now().UTC()
	expires := now.Add(s.ttl)
	member := connectionMember(tenant, device)
	epoch, err := claimOwnershipScript.Run(ctx, s.client, []string{ownershipKey(tenant, device), epochKey(tenant, device), "polaris:gateway:" + gateway + ":connections", "polaris:connections:lease-expiry"}, gateway, connection, credential, now.UnixMilli(), expires.UnixMilli(), s.ttl.Milliseconds(), member).Int64()
	return ConnectionOwnership{TenantID: tenant, DeviceID: device, GatewayID: gateway, ConnectionID: connection, CredentialID: credential, Epoch: epoch, ConnectedAt: now, LeaseExpiresAt: expires}, err
}

func (s *ConnectionOwnershipStore) Refresh(ctx context.Context, v ConnectionOwnership) (bool, error) {
	expires := time.Now().UTC().Add(s.ttl)
	result, err := refreshOwnershipScript.Run(ctx, s.client, []string{ownershipKey(v.TenantID, v.DeviceID), "polaris:connections:lease-expiry"}, v.GatewayID, v.ConnectionID, v.Epoch, expires.UnixMilli(), s.ttl.Milliseconds(), connectionMember(v.TenantID, v.DeviceID)).Int()
	return result == 1, err
}

func (s *ConnectionOwnershipStore) Release(ctx context.Context, v ConnectionOwnership) (bool, error) {
	result, err := releaseOwnershipScript.Run(ctx, s.client, []string{ownershipKey(v.TenantID, v.DeviceID), "polaris:gateway:" + v.GatewayID + ":connections", "polaris:connections:lease-expiry"}, v.GatewayID, v.ConnectionID, v.Epoch, connectionMember(v.TenantID, v.DeviceID)).Int()
	return result == 1, err
}

func (s *ConnectionOwnershipStore) Get(ctx context.Context, tenant, device string) (ConnectionOwnership, error) {
	values, err := s.client.HGetAll(ctx, ownershipKey(tenant, device)).Result()
	if err != nil {
		return ConnectionOwnership{}, err
	}
	if len(values) == 0 {
		return ConnectionOwnership{}, ErrNotFound
	}
	epoch, _ := strconv.ParseInt(values["epoch"], 10, 64)
	connected, _ := strconv.ParseInt(values["connected_at"], 10, 64)
	expires, _ := strconv.ParseInt(values["lease_expires_at"], 10, 64)
	return ConnectionOwnership{TenantID: tenant, DeviceID: device, GatewayID: values["gateway_id"], ConnectionID: values["connection_id"], CredentialID: values["credential_id"], Epoch: epoch, ConnectedAt: time.UnixMilli(connected), LeaseExpiresAt: time.UnixMilli(expires)}, nil
}

func (s *ConnectionOwnershipStore) Owns(ctx context.Context, v ConnectionOwnership) (bool, error) {
	current, err := s.Get(ctx, v.TenantID, v.DeviceID)
	if err != nil {
		return false, err
	}
	return current.GatewayID == v.GatewayID && current.ConnectionID == v.ConnectionID && current.Epoch == v.Epoch && current.LeaseExpiresAt.After(time.Now()), nil
}

func (s *ConnectionOwnershipStore) CleanExpired(ctx context.Context) error {
	now := time.Now().UnixMilli()
	members, err := s.client.ZRangeByScore(ctx, "polaris:connections:lease-expiry", &redis.ZRangeBy{Min: "-inf", Max: fmt.Sprint(now)}).Result()
	if err != nil {
		return err
	}
	for _, member := range members {
		parts := splitMember(member)
		if len(parts) != 2 {
			continue
		}
		current, getErr := s.Get(ctx, parts[0], parts[1])
		if getErr != nil || !current.LeaseExpiresAt.After(time.Now()) {
			_ = s.client.ZRem(ctx, "polaris:connections:lease-expiry", member).Err()
		}
	}
	return nil
}

func splitMember(member string) []string {
	for i := 0; i < len(member); i++ {
		if member[i] == ':' {
			return []string{member[:i], member[i+1:]}
		}
	}
	return nil
}
```

---

## internal\application\dispatch\dispatcher.go

```
package dispatch

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"sync/atomic"
	"time"

	"github.com/Akashpg-M/polaris/backend/internal/adapter/repository"
	"github.com/Akashpg-M/polaris/backend/internal/application/outbox"
	"github.com/Akashpg-M/polaris/backend/internal/core/command"
	"github.com/redis/go-redis/v9"
	"github.com/segmentio/kafka-go"
)

type Dispatcher struct {
	reader *kafka.Reader
	redis  *redis.Client
	owners *repository.ConnectionOwnershipStore
	health atomic.Int64
	done   chan struct{}
}

func New(broker string, redisClient *redis.Client, owners *repository.ConnectionOwnershipStore) *Dispatcher {
	reader := kafka.NewReader(kafka.ReaderConfig{Brokers: []string{broker}, Topic: outbox.CommandTopic, GroupID: "polaris-command-dispatcher", MinBytes: 1, MaxBytes: 10e6, CommitInterval: 0})
	d := &Dispatcher{reader: reader, redis: redisClient, owners: owners, done: make(chan struct{})}
	d.health.Store(time.Now().UnixMilli())
	return d
}

func (d *Dispatcher) Start(ctx context.Context) {
	defer close(d.done)
	defer d.reader.Close()
	for {
		message, err := d.reader.FetchMessage(ctx)
		if err != nil {
			if ctx.Err() != nil {
				return
			}
			slog.Error("command dispatcher fetch failed", "error", err)
			continue
		}
		var envelope command.Envelope
		if json.Unmarshal(message.Value, &envelope) != nil || envelope.FrameType != "COMMAND" || envelope.PartitionKey() != string(message.Key) {
			slog.Error("invalid durable command envelope", "partition", message.Partition, "offset", message.Offset)
			if err = d.reader.CommitMessages(ctx, message); err != nil {
				slog.Error("invalid command offset commit failed", "error", err)
			}
			continue
		}
		ownership, lookupErr := d.owners.Get(ctx, envelope.TenantID, envelope.DeviceID)
		if lookupErr == nil && ownership.LeaseExpiresAt.After(time.Now()) {
			if err = d.redis.Publish(ctx, repository.GatewayCommandChannel(ownership.GatewayID), message.Value).Err(); err != nil {
				slog.Error("command routing notification failed", "command_id", envelope.CommandID, "error", err)
				continue
			}
		} else if lookupErr != nil && !errors.Is(lookupErr, repository.ErrNotFound) {
			slog.Error("command ownership lookup failed", "command_id", envelope.CommandID, "error", lookupErr)
			continue
		}
		if err = d.reader.CommitMessages(ctx, message); err != nil {
			slog.Error("command dispatcher commit failed; notification may replay", "command_id", envelope.CommandID, "error", err)
			continue
		}
		d.health.Store(time.Now().UnixMilli())
	}
}

func (d *Dispatcher) Healthy() bool {
	return d.health.Load() > 0
}

func (d *Dispatcher) Wait(ctx context.Context) error {
	select {
	case <-d.done:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}
```

---

## internal\application\orchestration\metrics.go

```
package orchestration

import "sync/atomic"

type Metrics struct {
	TasksCreated, TasksCompleted, TasksFailed                                 atomic.Int64
	CommandsCreated, CommandsDelivered, CommandsAcknowledged, CommandsRetried atomic.Int64
	CommandsExpired, CommandsFailed, ActiveConnections, PendingCommands       atomic.Int64
}

func NewMetrics() *Metrics { return &Metrics{} }

func (m *Metrics) Snapshot() map[string]int64 {
	return map[string]int64{
		"tasks_created_total": m.TasksCreated.Load(), "tasks_completed_total": m.TasksCompleted.Load(), "tasks_failed_total": m.TasksFailed.Load(),
		"commands_created_total": m.CommandsCreated.Load(), "commands_delivered_total": m.CommandsDelivered.Load(), "commands_acknowledged_total": m.CommandsAcknowledged.Load(),
		"commands_retried_total": m.CommandsRetried.Load(), "commands_expired_total": m.CommandsExpired.Load(), "commands_failed_total": m.CommandsFailed.Load(),
		"active_device_connections": m.ActiveConnections.Load(), "pending_commands": m.PendingCommands.Load(),
	}
}
```

---

## internal\application\orchestration\service.go

```
package orchestration

import (
	"context"
	"encoding/json"
	"errors"
	"math"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/Akashpg-M/polaris/backend/internal/adapter/repository"
	"github.com/Akashpg-M/polaris/backend/internal/core/auth"
	"github.com/Akashpg-M/polaris/backend/internal/core/command"
	"github.com/Akashpg-M/polaris/backend/internal/core/extension"
	taskcore "github.com/Akashpg-M/polaris/backend/internal/core/task"
	twincore "github.com/Akashpg-M/polaris/backend/internal/core/twin"
	"github.com/redis/go-redis/v9"
)

type Service struct {
	store       *repository.RegistryStore
	redis       *redis.Client
	maxAttempts int
	metrics     *Metrics
	extensions  *extension.Registry
}

type CreateTaskInput struct {
	ProjectID     *string
	TaskType      string
	Priority      string
	Requirements  taskcore.Requirements
	Target        json.RawMessage
	ExpiresAt     time.Time
	CorrelationID string
}

type CreateResult struct {
	Task    taskcore.Task   `json:"task"`
	Command *command.Record `json:"command,omitempty"`
	Timing  *TaskPathTiming `json:"timing,omitempty"`
}

type TaskPathTiming struct {
	CandidateSelectionDurationUS int64 `json:"candidate_selection_duration_us"`
	RoutingDurationUS            int64 `json:"routing_duration_us"`
	PersistenceDurationUS        int64 `json:"persistence_duration_us"`
	TotalDurationUS              int64 `json:"total_duration_us"`
}

func (t *TaskPathTiming) add(other TaskPathTiming) {
	t.CandidateSelectionDurationUS += other.CandidateSelectionDurationUS
	t.RoutingDurationUS += other.RoutingDurationUS
	t.PersistenceDurationUS += other.PersistenceDurationUS
}

func NewService(store *repository.RegistryStore, redisClient *redis.Client, maxAttempts int, metrics *Metrics) *Service {
	registry := extension.NewRegistry()
	registry.RegisterTaskPlanner(extension.DefaultTaskPlanner{})
	return NewServiceWithRegistry(store, redisClient, maxAttempts, metrics, registry)
}

func NewServiceWithRegistry(store *repository.RegistryStore, redisClient *redis.Client, maxAttempts int, metrics *Metrics, registry *extension.Registry) *Service {
	if maxAttempts < 1 {
		maxAttempts = 5
	}
	if metrics == nil {
		metrics = NewMetrics()
	}
	if registry == nil {
		registry = extension.NewRegistry()
		registry.RegisterTaskPlanner(extension.DefaultTaskPlanner{})
	}
	return &Service{store: store, redis: redisClient, maxAttempts: maxAttempts, metrics: metrics, extensions: registry}
}

func (s *Service) CreateTask(ctx context.Context, tenant string, principal auth.OperatorPrincipal, requestID string, in CreateTaskInput) (result CreateResult, err error) {
	started := time.Now()
	timing := TaskPathTiming{}
	defer func() {
		timing.TotalDurationUS = time.Since(started).Microseconds()
		result.Timing = &timing
	}()
	if tenant == "" || in.TaskType == "" || !validPriority(in.Priority) || in.ExpiresAt.Before(time.Now()) {
		return CreateResult{}, ErrInvalidTask
	}
	if len(in.Target) == 0 || !json.Valid(in.Target) {
		return CreateResult{}, ErrInvalidTask
	}
	if capability := command.RequiredCapability(in.TaskType); capability == "" {
		return CreateResult{}, ErrUnsupportedCommand
	} else if !contains(in.Requirements.RequiredCapabilities, capability) {
		in.Requirements.RequiredCapabilities = append(in.Requirements.RequiredCapabilities, capability)
	}
	if in.Requirements.MinimumBattery < 0 || in.Requirements.MinimumBattery > 100 {
		return CreateResult{}, ErrInvalidTask
	}
	if in.TaskType == "NAVIGATE" || in.TaskType == "RELOCATE" {
		if in.Requirements.PlanningMode == "" {
			in.Requirements.PlanningMode = taskcore.PlanningDeviceLocal
		}
		if in.Requirements.PlanningMode != taskcore.PlanningDeviceLocal && in.Requirements.PlanningMode != taskcore.PlanningPolarisRequired {
			return CreateResult{}, ErrInvalidTask
		}
	} else if in.Requirements.PlanningMode != "" {
		return CreateResult{}, ErrInvalidTask
	}
	if in.ProjectID != nil && in.Requirements.ProjectID == "" {
		in.Requirements.ProjectID = *in.ProjectID
	}
	requirements, _ := json.Marshal(in.Requirements)
	correlation := in.CorrelationID
	if correlation == "" {
		correlation = auth.NewID()
	}
	v := taskcore.Task{TaskID: auth.NewID(), TenantID: tenant, ProjectID: in.ProjectID, TaskType: in.TaskType, Status: string(taskcore.Pending), Priority: in.Priority, Requirements: requirements, Target: in.Target, CorrelationID: correlation, CreatedBy: principal.APIKeyID, ExpiresAt: in.ExpiresAt}
	persistenceStarted := time.Now()
	if err := s.store.CreateTask(ctx, v, principal.APIKeyID, requestID); err != nil {
		timing.PersistenceDurationUS += time.Since(persistenceStarted).Microseconds()
		return CreateResult{}, err
	}
	timing.PersistenceDurationUS += time.Since(persistenceStarted).Microseconds()
	s.metrics.TasksCreated.Add(1)
	created, err := s.store.GetTask(ctx, tenant, v.TaskID)
	if err != nil {
		return CreateResult{}, err
	}
	cmd, assignTiming, assignErr := s.assignTimed(ctx, created, principal.APIKeyID, requestID)
	timing.add(assignTiming)
	if assignErr != nil && !errors.Is(assignErr, ErrNoEligibleDevice) && !errors.Is(assignErr, repository.ErrConflict) {
		return CreateResult{Task: created}, assignErr
	}
	created, _ = s.store.GetTask(ctx, tenant, v.TaskID)
	if assignErr == nil {
		s.metrics.CommandsCreated.Add(1)
		return CreateResult{Task: created, Command: &cmd}, nil
	}
	return CreateResult{Task: created}, nil
}

func (s *Service) Assign(ctx context.Context, v taskcore.Task, actor, requestID string) (command.Record, error) {
	record, _, err := s.assignTimed(ctx, v, actor, requestID)
	return record, err
}

func (s *Service) assignTimed(ctx context.Context, v taskcore.Task, actor, requestID string) (command.Record, TaskPathTiming, error) {
	timing := TaskPathTiming{}
	selectionStarted := time.Now()
	var requirements taskcore.Requirements
	if err := json.Unmarshal(v.Requirements, &requirements); err != nil {
		return command.Record{}, timing, ErrInvalidTask
	}
	candidates, err := s.store.EligibleDevices(ctx, v.TenantID, requirements)
	if err != nil {
		return command.Record{}, timing, err
	}
	eligible := make(map[string]repository.DeviceCandidate, len(candidates))
	eligibleIDs := make([]string, 0, len(candidates))
	for _, candidate := range candidates {
		eligible[candidate.DeviceID] = candidate
		eligibleIDs = append(eligibleIDs, candidate.DeviceID)
	}
	candidateTiming := &extension.CandidateTiming{}
	request := extension.CandidateRequest{TenantID: v.TenantID, EligibleDeviceIDs: eligibleIDs, RequiredCapabilities: requirements.RequiredCapabilities, DeviceTypeIDs: requirements.AllowedDeviceTypes, Limit: 50, Context: map[string]any{"maximum_distance_meters": requirements.MaximumDistanceM}, Timing: candidateTiming}
	if lat, lon, ok := targetCoordinates(v.Target); ok {
		request.Context["target_latitude"], request.Context["target_longitude"] = lat, lon
	}
	proposals := make([]extension.Candidate, 0, len(candidates))
	domainRanked := false
	if provider, providerErr := s.extensions.CandidateProvider(request); providerErr == nil {
		proposals, err = provider.Candidates(ctx, request)
		if err != nil {
			timing.RoutingDurationUS += candidateTiming.RoutingDuration.Microseconds()
			timing.CandidateSelectionDurationUS += (time.Since(selectionStarted) - candidateTiming.RoutingDuration).Microseconds()
			return command.Record{}, timing, err
		}
		// A domain provider ranks the candidates it understands; it must not
		// accidentally become an exclusive capability filter. Preserve every
		// core-eligible device as an unscored fallback (important for generic
		// non-spatial tasks and for spatial raw-result limits).
		proposals = includeUnrankedEligible(proposals, candidates)
		domainRanked = true
	} else {
		for _, candidate := range candidates {
			proposals = append(proposals, extension.Candidate{DeviceID: candidate.DeviceID})
		}
	}
	// Candidate evaluation is deliberately bounded. Providers may rank a
	// domain-specific subset and Core appends deterministic eligible fallbacks,
	// but a single task must not fan out Redis reads across an entire tenant.
	if request.Limit > 0 && len(proposals) > request.Limit {
		proposals = proposals[:request.Limit]
	}
	type candidateState struct {
		twin, connection *redis.MapStringStringCmd
	}
	states := make(map[string]candidateState, len(eligible))
	pipe := s.redis.Pipeline()
	for _, proposal := range proposals {
		if _, allowed := eligible[proposal.DeviceID]; !allowed {
			continue
		}
		if _, exists := states[proposal.DeviceID]; exists {
			continue
		}
		states[proposal.DeviceID] = candidateState{
			twin:       pipe.HGetAll(ctx, "polaris:twin:"+v.TenantID+":"+proposal.DeviceID),
			connection: pipe.HGetAll(ctx, "polaris:connection:"+v.TenantID+":"+proposal.DeviceID),
		}
	}
	if _, pipeErr := pipe.Exec(ctx); pipeErr != nil && !errors.Is(pipeErr, redis.Nil) {
		timing.RoutingDurationUS += candidateTiming.RoutingDuration.Microseconds()
		timing.CandidateSelectionDurationUS += (time.Since(selectionStarted) - candidateTiming.RoutingDuration).Microseconds()
		return command.Record{}, timing, pipeErr
	}
	ranked := make([]rankedCandidate, 0, len(proposals))
	for order, proposal := range proposals {
		candidate, allowed := eligible[proposal.DeviceID]
		if !allowed {
			continue
		}
		loaded := states[candidate.DeviceID]
		state, stateErr := loaded.twin.Result()
		connection, connectionErr := loaded.connection.Result()
		if stateErr != nil || connectionErr != nil || (state["connectivity_status"] != "ONLINE" && !activeConnectionState(connection)) {
			continue
		}
		var reported struct {
			Lat           float64 `json:"lat"`
			Lon           float64 `json:"lon"`
			EnergyPercent int32   `json:"energy_percent"`
		}
		hasReported := json.Unmarshal([]byte(state["reported_state"]), &reported) == nil
		_, _, hasSpatialTarget := targetCoordinates(v.Target)
		if (!hasReported && (requirements.MinimumBattery > 0 || hasSpatialTarget)) || (hasReported && reported.EnergyPercent < requirements.MinimumBattery) {
			continue
		}
		distance := 0.0
		if targetLat, targetLon, ok := targetCoordinates(v.Target); ok {
			distance = haversineMeters(reported.Lat, reported.Lon, targetLat, targetLon)
			if requirements.MaximumDistanceM > 0 && distance > requirements.MaximumDistanceM {
				continue
			}
		}
		ranked = append(ranked, rankedCandidate{id: candidate.DeviceID, distance: distance, battery: reported.EnergyPercent, domainScore: proposal.DomainScore, proposalOrder: order, twin: twinFromState(v.TenantID, candidate.DeviceID, state, connection)})
	}
	sort.Slice(ranked, func(i, j int) bool {
		if domainRanked {
			if ranked[i].domainScore != nil && ranked[j].domainScore != nil && *ranked[i].domainScore != *ranked[j].domainScore {
				return *ranked[i].domainScore < *ranked[j].domainScore
			}
			if (ranked[i].domainScore != nil) != (ranked[j].domainScore != nil) {
				return ranked[i].domainScore != nil
			}
			if ranked[i].proposalOrder != ranked[j].proposalOrder {
				return ranked[i].proposalOrder < ranked[j].proposalOrder
			}
		}
		if ranked[i].distance != ranked[j].distance {
			return ranked[i].distance < ranked[j].distance
		}
		if ranked[i].battery != ranked[j].battery {
			return ranked[i].battery > ranked[j].battery
		}
		return ranked[i].id < ranked[j].id
	})
	timing.RoutingDurationUS += candidateTiming.RoutingDuration.Microseconds()
	timing.CandidateSelectionDurationUS += (time.Since(selectionStarted) - candidateTiming.RoutingDuration).Microseconds()
	selectionStarted = time.Now()
	fresh, recheckErr := s.store.EligibleDevices(ctx, v.TenantID, requirements)
	if recheckErr != nil {
		timing.CandidateSelectionDurationUS += time.Since(selectionStarted).Microseconds()
		return command.Record{}, timing, recheckErr
	}
	stillEligible := make(map[string]struct{}, len(fresh))
	for _, candidate := range fresh {
		stillEligible[candidate.DeviceID] = struct{}{}
	}
	timing.CandidateSelectionDurationUS += time.Since(selectionStarted).Microseconds()
	for _, candidate := range ranked {
		if _, ok := stillEligible[candidate.id]; !ok {
			continue
		}
		twin := candidate.twin
		if twin.Connectivity != "ONLINE" {
			continue
		}
		plannedTask := v
		plannedTask.AssignedDeviceID = &candidate.id
		planners := s.extensions.TaskPlanners(plannedTask)
		if len(planners) == 0 {
			return command.Record{}, timing, ErrUnsupportedCommand
		}
		var plan extension.ExecutionPlan
		var planErr error
		planningStarted := time.Now()
		for _, planner := range planners {
			plan, planErr = planner.Plan(ctx, extension.PlanningRequest{Task: plannedTask, DeviceTwin: twin})
			if errors.Is(planErr, extension.ErrPlanningUnsupported) {
				continue
			}
			break
		}
		timing.RoutingDurationUS += time.Since(planningStarted).Microseconds()
		if planErr != nil {
			return command.Record{}, timing, planErr
		}
		persistenceStarted := time.Now()
		cmd, err := s.store.AssignTaskWithPlan(ctx, v, candidate.id, actor, requestID, s.maxAttempts, plan)
		timing.PersistenceDurationUS += time.Since(persistenceStarted).Microseconds()
		if err == nil {
			return cmd, timing, nil
		}
		if !errors.Is(err, repository.ErrConflict) {
			return command.Record{}, timing, err
		}
	}
	return command.Record{}, timing, ErrNoEligibleDevice
}

func includeUnrankedEligible(proposals []extension.Candidate, eligible []repository.DeviceCandidate) []extension.Candidate {
	seen := make(map[string]struct{}, len(proposals))
	for _, proposal := range proposals {
		seen[proposal.DeviceID] = struct{}{}
	}
	for _, candidate := range eligible {
		if _, ok := seen[candidate.DeviceID]; ok {
			continue
		}
		proposals = append(proposals, extension.Candidate{DeviceID: candidate.DeviceID})
	}
	return proposals
}

func (s *Service) RetryTask(ctx context.Context, tenant, id string, principal auth.OperatorPrincipal, requestID string, ttl time.Duration) (CreateResult, error) {
	if ttl <= 0 {
		ttl = 5 * time.Minute
	}
	if err := s.store.RetryTask(ctx, tenant, id, principal.APIKeyID, requestID, time.Now().Add(ttl)); err != nil {
		return CreateResult{}, err
	}
	v, err := s.store.GetTask(ctx, tenant, id)
	if err != nil {
		return CreateResult{}, err
	}
	cmd, assignErr := s.Assign(ctx, v, principal.APIKeyID, requestID)
	v, _ = s.store.GetTask(ctx, tenant, id)
	if assignErr == nil {
		s.metrics.CommandsCreated.Add(1)
		return CreateResult{Task: v, Command: &cmd}, nil
	}
	if errors.Is(assignErr, ErrNoEligibleDevice) {
		return CreateResult{Task: v}, nil
	}
	return CreateResult{Task: v}, assignErr
}

type rankedCandidate struct {
	id            string
	distance      float64
	battery       int32
	domainScore   *float64
	proposalOrder int
	twin          twincore.DeviceTwin
}

func (s *Service) loadTwin(ctx context.Context, tenant, device string) (twincore.DeviceTwin, error) {
	state, err := s.redis.HGetAll(ctx, "polaris:twin:"+tenant+":"+device).Result()
	if err != nil {
		return twincore.DeviceTwin{}, err
	}
	connection, connectionErr := s.redis.HGetAll(ctx, "polaris:connection:"+tenant+":"+device).Result()
	if connectionErr != nil && !errors.Is(connectionErr, redis.Nil) {
		return twincore.DeviceTwin{}, connectionErr
	}
	return twinFromState(tenant, device, state, connection), nil
}

func twinFromState(tenant, device string, state, connection map[string]string) twincore.DeviceTwin {
	twin := twincore.DeviceTwin{TenantID: tenant, DeviceID: device, Connectivity: state["connectivity_status"], Components: map[string]twincore.ComponentEnvelope{}}
	if twin.Connectivity != "ONLINE" && activeConnectionState(connection) {
		twin.Connectivity = "ONLINE"
	}
	for field, raw := range state {
		if strings.HasPrefix(field, "component:") {
			var c twincore.ComponentEnvelope
			if json.Unmarshal([]byte(raw), &c) == nil {
				twin.Components[strings.TrimPrefix(field, "component:")] = c
			}
		}
	}
	return twin
}

func activeConnectionState(state map[string]string) bool {
	if state["gateway_id"] == "" {
		return false
	}
	expiresAt, err := strconv.ParseInt(state["lease_expires_at"], 10, 64)
	return err == nil && expiresAt > time.Now().UnixMilli()
}

func targetCoordinates(raw json.RawMessage) (float64, float64, bool) {
	var value map[string]interface{}
	if json.Unmarshal(raw, &value) != nil {
		return 0, 0, false
	}
	lat, lok := number(value["lat"])
	lon, ook := number(value["lon"])
	return lat, lon, lok && ook && lat >= -90 && lat <= 90 && lon >= -180 && lon <= 180
}

func number(v interface{}) (float64, bool) {
	switch n := v.(type) {
	case float64:
		return n, true
	case string:
		f, err := strconv.ParseFloat(n, 64)
		return f, err == nil
	}
	return 0, false
}

func haversineMeters(lat1, lon1, lat2, lon2 float64) float64 {
	const earth = 6371000.0
	toRad := math.Pi / 180
	dLat, dLon := (lat2-lat1)*toRad, (lon2-lon1)*toRad
	a := math.Sin(dLat/2)*math.Sin(dLat/2) + math.Cos(lat1*toRad)*math.Cos(lat2*toRad)*math.Sin(dLon/2)*math.Sin(dLon/2)
	return 2 * earth * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
}

func contains(values []string, expected string) bool {
	for _, value := range values {
		if value == expected {
			return true
		}
	}
	return false
}

func validPriority(value string) bool {
	return value == "LOW" || value == "NORMAL" || value == "HIGH" || value == "CRITICAL"
}

var (
	ErrInvalidTask        = errors.New("invalid task")
	ErrUnsupportedCommand = errors.New("unsupported command type")
	ErrNoEligibleDevice   = errors.New("no eligible device")
)
```

---

## internal\application\orchestration\service_test.go

```
package orchestration

import (
	"encoding/json"
	"fmt"
	"math"
	"testing"
	"time"

	"github.com/Akashpg-M/polaris/backend/internal/adapter/repository"
	"github.com/Akashpg-M/polaris/backend/internal/core/extension"
)

func TestDeterministicHelpers(t *testing.T) {
	if d := haversineMeters(13.0067, 80.2206, 13.0067, 80.2206); math.Abs(d) > 0.001 {
		t.Fatalf("same point distance=%f", d)
	}
	if _, _, ok := targetCoordinates(json.RawMessage(`{"lat":91,"lon":80}`)); ok {
		t.Fatal("invalid latitude accepted")
	}
	if !validPriority("HIGH") || validPriority("URGENT") {
		t.Fatal("priority validation incorrect")
	}
}

func TestActiveConnectionStateRequiresLiveLease(t *testing.T) {
	if activeConnectionState(map[string]string{"gateway_id": "gateway-1", "lease_expires_at": "1"}) {
		t.Fatal("expired connection accepted")
	}
	future := time.Now().Add(time.Minute).UnixMilli()
	if !activeConnectionState(map[string]string{"gateway_id": "gateway-1", "lease_expires_at": fmt.Sprint(future)}) {
		t.Fatal("live connection rejected")
	}
}

func TestDomainProviderCannotDropCoreEligibleCandidates(t *testing.T) {
	score := 1.0
	got := includeUnrankedEligible(
		[]extension.Candidate{{DeviceID: "ranked-camera", DomainScore: &score}},
		[]repository.DeviceCandidate{{DeviceID: "vehicle"}, {DeviceID: "compute"}},
	)
	if len(got) != 3 || got[0].DeviceID != "ranked-camera" || got[1].DeviceID != "vehicle" || got[2].DeviceID != "compute" {
		t.Fatalf("eligible fallback candidates lost: %#v", got)
	}
	got = includeUnrankedEligible(got, []repository.DeviceCandidate{{DeviceID: "vehicle"}})
	if len(got) != 3 {
		t.Fatalf("duplicate fallback candidate added: %#v", got)
	}
}
```

---

## internal\application\orchestrator\predictive_strategy.go

```
package orchestrator

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	pb "github.com/Akashpg-M/polaris/backend/api/proto/v1"
	"github.com/jmoiron/sqlx"
)

// PredictiveZoneStrategy uses historical spatial clustering to predict demand
type PredictiveZoneStrategy struct {
	db *sqlx.DB
}

func NewPredictiveZoneStrategy(postgresURL string) (*PredictiveZoneStrategy, error) {
	db, err := sqlx.Connect("postgres", postgresURL)
	if err != nil {
		return nil, err
	}
	db.SetMaxOpenConns(4)
	db.SetMaxIdleConns(2)
	db.SetConnMaxIdleTime(5 * time.Minute)
	db.SetConnMaxLifetime(30 * time.Minute)
	return &PredictiveZoneStrategy{db: db}, nil
}

func (s *PredictiveZoneStrategy) Close() error { return s.db.Close() }

func (s *PredictiveZoneStrategy) GetTargetZones(ctx context.Context) []Zone {
	// ML Clustering via SQL:
	// We divide the map into a grid by rounding Lat/Lon to 2 decimal places (~1.1km accuracy).
	// We count how many pings happened in each grid over the last hour.
	// The top grids become our "Predicted Hotspots".

	query := `
		SELECT 
			ROUND(lat::numeric, 2) AS cluster_lat,
			ROUND(lon::numeric, 2) AS cluster_lon,
			COUNT(*) as ping_count
		FROM telemetry_history
		WHERE recorded_at >= NOW() - INTERVAL '1 hour'
		GROUP BY cluster_lat, cluster_lon
		ORDER BY ping_count DESC
		LIMIT 3; -- Pick the top 3 highest-density hotspots
	`

	rows, err := s.db.QueryxContext(ctx, query)
	if err != nil {
		slog.Error("Failed to run predictive clustering", "error", err)
		return []Zone{} // Fallback to empty if DB is busy
	}
	defer rows.Close()

	var zones []Zone
	hubIndex := 1

	for rows.Next() {
		var lat, lon float64
		var count int
		if err := rows.Scan(&lat, &lon, &count); err == nil {
			zones = append(zones, Zone{
				ID:             fmt.Sprintf("Predicted-Hotspot-%d", hubIndex),
				Lat:            lat,
				Lon:            lon,
				RadiusKm:       2.0, // Create a 2km radius catch-zone
				RequiredAssets: 5,   // Require 5 drones to pre-position here
				TargetClass:    pb.NodeType_NODE_TYPE_DRONE,
				TenantID:       "alpha_logistics",
			})
			hubIndex++
		}
	}

	if len(zones) > 0 {
		slog.Info("Predictive ML Engine updated hotspots", "clusters_found", len(zones))
	}
	return zones
}
```

---

## internal\application\orchestrator\strategies.go

```
package orchestrator

import (
	"context"
	pb "github.com/Akashpg-M/polaris/backend/api/proto/v1"
)

// StaticZoneStrategy simulates a database table of logistics hubs or smart-city sectors
type StaticZoneStrategy struct {}

func (s *StaticZoneStrategy) GetTargetZones(ctx context.Context) []Zone {
	return []Zone{
		{
			ID:             "Hub-Guindy",
			Lat:            13.0067,
			Lon:            80.2206,
			RadiusKm:       5.0,
			RequiredAssets: 3,
			TargetClass:    pb.NodeType_NODE_TYPE_DRONE,
      TenantID:       "alpha_logistics",
		},
		{
			ID:             "Hub-Adyar",
			Lat:            13.0012,
			Lon:            80.2565,
			RadiusKm:       3.0,
			RequiredAssets: 2,
			TargetClass:    pb.NodeType_NODE_TYPE_DRONE,
      TenantID:       "alpha_logistics",
		},
	}
}
```

---

## internal\application\orchestrator\zone.go

```
package orchestrator

import pb "github.com/Akashpg-M/polaris/backend/api/proto/v1"

// Zone is a read-only predicted demand view retained for the operator API.
// Predictions never issue commands directly; a future advisor may translate
// them into ordinary durable Polaris tasks.
type Zone struct {
	ID             string      `json:"id"`
	Lat            float64     `json:"lat"`
	Lon            float64     `json:"lon"`
	RadiusKm       float64     `json:"radius_km"`
	RequiredAssets int         `json:"required_assets"`
	TargetClass    pb.NodeType `json:"target_class"`
	TenantID       string      `json:"tenant_id"`
}
```

---

## internal\application\outbox\relay.go

```
package outbox

import (
	"context"
	"encoding/json"
	"log/slog"
	"time"

	"github.com/Akashpg-M/polaris/backend/internal/adapter/repository"
	"github.com/Akashpg-M/polaris/backend/internal/core/command"
	"github.com/segmentio/kafka-go"
)

const (
	LifecycleTopic     = "device.lifecycle.v1"
	TaskTopic          = "task.lifecycle.v1"
	CommandTopic       = "device.command.v1"
	CommandAckTopic    = "device.command.ack.v1"
	CommandResultTopic = "device.command.result.v1"
)

type Relay struct {
	store    *repository.RegistryStore
	writer   *kafka.Writer
	batch    int
	interval time.Duration
	done     chan struct{}
}

func New(store *repository.RegistryStore, broker string, batch int, interval time.Duration) *Relay {
	if batch < 1 {
		batch = 100
	}
	return &Relay{store: store, writer: &kafka.Writer{Addr: kafka.TCP(broker), Balancer: &kafka.Hash{}}, batch: batch, interval: interval, done: make(chan struct{})}
}
func (r *Relay) Start(ctx context.Context) {
	defer close(r.done)
	defer r.writer.Close()
	ticker := time.NewTicker(r.interval)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			r.flush(ctx)
		case <-ctx.Done():
			shutdown, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			r.flush(shutdown)
			cancel()
			return
		}
	}
}
func (r *Relay) flush(ctx context.Context) {
	events, err := r.store.ClaimOutbox(ctx, r.batch)
	if err != nil {
		slog.Error("outbox claim failed", "error", err)
		return
	}
	if len(events) == 0 {
		return
	}
	messages := make([]kafka.Message, 0, len(events))
	ids := make([]string, 0, len(events))
	for _, e := range events {
		topic, key, value := route(e)
		if topic == CommandTopic {
			var envelope command.Envelope
			if json.Unmarshal(value, &envelope) == nil {
				envelope.DeliveryObservation = &command.DeliveryObservation{RelayPublishedAt: time.Now().UTC()}
				if observed, marshalErr := json.Marshal(envelope); marshalErr == nil {
					value = observed
				}
			}
		}
		messages = append(messages, kafka.Message{Topic: topic, Key: []byte(key), Value: value})
		ids = append(ids, e.OutboxID)
	}
	if err = r.writer.WriteMessages(ctx, messages...); err != nil {
		// Kafka may have accepted a subset. Replay the complete batch: every
		// downstream consumer is idempotent and this preserves at-least-once.
		for _, e := range events {
			_ = r.store.MarkOutboxFailed(ctx, e.OutboxID, err.Error())
		}
		return
	}
	if err = r.store.MarkOutboxPublishedBatch(ctx, ids); err != nil {
		slog.Error("outbox batch publish marker failed; events will replay", "events", len(ids), "error", err)
	}
}

func route(e repository.OutboxEvent) (string, string, []byte) {
	key := e.TenantID + ":" + e.EventID
	value, _ := json.Marshal(map[string]interface{}{"event_id": e.EventID, "event_type": e.EventType, "schema_version": 1, "tenant_id": e.TenantID, "payload": json.RawMessage(e.Payload)})
	if e.EventType == "command.created.v1" || e.EventType == "command.retry.requested.v1" {
		var envelope struct {
			DeviceID string `json:"device_id"`
		}
		if json.Unmarshal(e.Payload, &envelope) == nil && envelope.DeviceID != "" {
			key = e.TenantID + ":" + envelope.DeviceID
		}
		return CommandTopic, key, e.Payload
	}
	if e.EventType == "command.acknowledged.v1" {
		return CommandAckTopic, key, value
	}
	if e.EventType == "command.result.v1" {
		return CommandResultTopic, key, value
	}
	if len(e.EventType) >= 5 && e.EventType[:5] == "task." {
		return TaskTopic, key, value
	}
	return LifecycleTopic, key, value
}
```

---

## internal\application\reconciliation\worker.go

```
package reconciliation

import (
	"context"
	"errors"
	"log/slog"
	"time"

	"github.com/Akashpg-M/polaris/backend/internal/adapter/repository"
	"github.com/Akashpg-M/polaris/backend/internal/application/orchestration"
)

type Worker struct {
	store      *repository.RegistryStore
	service    *orchestration.Service
	owners     *repository.ConnectionOwnershipStore
	interval   time.Duration
	ackTimeout time.Duration
}

func New(store *repository.RegistryStore, service *orchestration.Service, owners *repository.ConnectionOwnershipStore, interval, ackTimeout time.Duration) *Worker {
	if interval < 100*time.Millisecond {
		interval = time.Second
	}
	return &Worker{store: store, service: service, owners: owners, interval: interval, ackTimeout: ackTimeout}
}

func (w *Worker) Start(ctx context.Context) {
	ticker := time.NewTicker(w.interval)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			w.run(ctx)
		}
	}
}

func (w *Worker) run(ctx context.Context) {
	if err := w.store.ReconcileCommands(ctx, w.ackTimeout); err != nil {
		slog.Error("command reconciliation failed", "error", err)
	}
	if err := w.store.FailExpiredPendingTasks(ctx); err != nil {
		slog.Error("task expiry reconciliation failed", "error", err)
	}
	if err := w.owners.CleanExpired(ctx); err != nil {
		slog.Error("connection lease reconciliation failed", "error", err)
	}
	pending, err := w.store.PendingTasks(ctx, 50)
	if err != nil {
		slog.Error("pending task scan failed", "error", err)
		return
	}
	for _, v := range pending {
		_, err = w.service.Assign(ctx, v, "reconciler", "")
		if err != nil && !errors.Is(err, orchestration.ErrNoEligibleDevice) && !errors.Is(err, repository.ErrConflict) && !errors.Is(err, repository.ErrInvalidTransition) {
			slog.Error("pending task assignment failed", "task_id", v.TaskID, "error", err)
		}
	}
}
```

---

## internal\application\spatial\engine.go

```
package spatial

import (
	"hash/fnv"
	"sort"
	"strings"
	"sync"

	"github.com/Akashpg-M/polaris/backend/algo_/geo"
	pb "github.com/Akashpg-M/polaris/backend/api/proto/v1"
	"github.com/Akashpg-M/polaris/backend/internal/core/events"
	"google.golang.org/protobuf/proto"
)

// ShardCount dictates how many independent memory partitions exist.
// 32 is a standard default to minimize lock contention across highly concurrent goroutines.
const ShardCount = 32

type EngineShard struct {
	mu       sync.RWMutex
	nodes    map[string]*pb.SpatialObject
	versions map[string]stateVersion
}

type stateVersion struct {
	bootID        string
	bootStartedAt int64
	sequence      uint64
	retired       map[string]struct{}
}

type Classification string

const (
	Accepted     Classification = "ACCEPTED"
	Duplicate    Classification = "DUPLICATE"
	OutOfOrder   Classification = "OUT_OF_ORDER"
	NewBoot      Classification = "NEW_BOOT"
	RetiredBoot  Classification = "RETIRED_BOOT"
	BootConflict Classification = "BOOT_CONFLICT"
)

type Engine struct {
	shards []*EngineShard
}

// MatchResult is the DTO sent back to the dispatcher
type MatchResult struct {
	NodeID     string      `json:"node_id"`
	Type       pb.NodeType `json:"node_type"`
	Class      uint16      `json:"asset_class"`
	Lat        float64     `json:"lat"`
	Lon        float64     `json:"lon"`
	DistanceKm float64     `json:"distance_km"`
	ETASec     int         `json:"eta_seconds"`
	RouteType  string      `json:"route_type"`
}

func NewEngine() *Engine {
	shards := make([]*EngineShard, ShardCount)
	for i := 0; i < ShardCount; i++ {
		shards[i] = &EngineShard{nodes: make(map[string]*pb.SpatialObject), versions: make(map[string]stateVersion)}
	}
	return &Engine{shards: shards}
}

// getShard picks the correct memory partition using an FNV-1a hash of the NodeID
func (e *Engine) getShard(nodeID string) *EngineShard {
	h := fnv.New32a()
	h.Write([]byte(nodeID))
	return e.shards[h.Sum32()%ShardCount]
}

func (e *Engine) BatchUpdate(payloads []*pb.SpatialObject) {
	if len(payloads) == 0 {
		return
	}

	for _, p := range payloads {
		shard := e.getShard(p.Id)

		shard.mu.Lock()
		shard.nodes[p.Id] = p
		shard.mu.Unlock()
	}
}

// FindNearest is a compatibility projection for the Phase 0 endpoint. Redis is
// the freshness authority and Mobility owns indexed candidate discovery. This
// bounded linear scan deliberately avoids retaining the former non-subdividing
// "QuadTree" as a second production spatial authority.
func (e *Engine) FindNearest(tenantID string, lat, lon, radiusKm float64, reqType pb.NodeType) []MatchResult {
	var results []MatchResult
	for _, shard := range e.shards {
		shard.mu.RLock()
		for key, node := range shard.nodes {
			if node.TenantId != tenantID || node.Type != reqType {
				continue
			}
			dist := geo.Haversine(lat, lon, node.Lat, node.Lon)
			if dist <= radiusKm {
				results = append(results, MatchResult{NodeID: strings.TrimPrefix(key, tenantID+":"), Type: node.Type, Class: uint16(node.Type), Lat: node.Lat, Lon: node.Lon, DistanceKm: dist, ETASec: int((dist / 40.0) * 3600), RouteType: "COMPATIBILITY_ESTIMATE"})
			}
		}
		shard.mu.RUnlock()
	}

	sort.Slice(results, func(i, j int) bool { return results[i].ETASec < results[j].ETASec })
	if len(results) > 500 {
		return results[:500]
	}
	return results
}

func stateKey(tenantID, deviceID string) string { return tenantID + ":" + deviceID }

// ApplyEnvelope is a defensive projection guard. KafkaConsumer invokes it only
// after Redis—the canonical Phase 1 freshness authority—accepts the event.
func (e *Engine) ApplyEnvelope(envelope *events.TelemetryEnvelope) Classification {
	key := stateKey(envelope.TenantID, envelope.DeviceID)
	shard := e.getShard(key)
	shard.mu.Lock()
	current, exists := shard.versions[key]
	classification := Accepted
	apply := false
	if !exists {
		current = stateVersion{bootID: envelope.DeviceBootID, bootStartedAt: envelope.BootStartedAt, sequence: envelope.SequenceNumber, retired: make(map[string]struct{})}
		apply = true
	} else if envelope.DeviceBootID == current.bootID {
		switch {
		case envelope.SequenceNumber > current.sequence:
			current.sequence = envelope.SequenceNumber
			apply = true
		case envelope.SequenceNumber == current.sequence:
			classification = Duplicate
		default:
			classification = OutOfOrder
		}
	} else if _, retired := current.retired[envelope.DeviceBootID]; retired {
		classification = RetiredBoot
	} else if envelope.BootStartedAt > current.bootStartedAt {
		current.retired[current.bootID] = struct{}{}
		current.bootID = envelope.DeviceBootID
		current.bootStartedAt = envelope.BootStartedAt
		current.sequence = envelope.SequenceNumber
		classification = NewBoot
		apply = true
	} else if envelope.BootStartedAt == current.bootStartedAt {
		classification = BootConflict
	} else {
		classification = RetiredBoot
	}
	if apply {
		payload := proto.Clone(envelope.Payload).(*pb.SpatialObject)
		shard.nodes[key] = payload
		shard.versions[key] = current
	}
	shard.mu.Unlock()
	return classification
}
```

---

## internal\application\spatial\engine_reliability_test.go

```
package spatial

import (
	"testing"
	"time"

	pb "github.com/Akashpg-M/polaris/backend/api/proto/v1"
	"github.com/Akashpg-M/polaris/backend/internal/core/events"
)

func testEnvelope(boot string, bootStarted int64, sequence uint64, lat float64) *events.TelemetryEnvelope {
	now := time.Now().UTC()
	p := &pb.SpatialObject{Id: "device-1", TenantId: "tenant-1", DeviceBootId: boot, SequenceNumber: sequence,
		BootStartedAt: bootStarted, ObservedAt: now.UnixMilli(), SchemaVersion: 1, Type: pb.NodeType_NODE_TYPE_DRONE,
		Status: pb.NodeStatus_NODE_STATUS_ACTIVE, Lat: lat, Lon: 80, EnergyPercent: 80}
	return events.NewTelemetryEnvelope(p, now, "", "", "")
}

func TestLatestStateClassifications(t *testing.T) {
	e := NewEngine()
	started := time.Now().Add(-time.Minute).UnixMilli()
	cases := []struct {
		name     string
		envelope *events.TelemetryEnvelope
		want     Classification
	}{
		{"first", testEnvelope("boot-a", started, 1, 13.0), Accepted},
		{"duplicate", testEnvelope("boot-a", started, 1, 99.0), Duplicate},
		{"newer sequence", testEnvelope("boot-a", started, 3, 13.2), Accepted},
		{"out of order", testEnvelope("boot-a", started, 2, 99.0), OutOfOrder},
		{"boot conflict", testEnvelope("boot-conflict", started, 1, 99.0), BootConflict},
		{"new boot", testEnvelope("boot-b", started+1000, 1, 13.3), NewBoot},
		{"retired boot", testEnvelope("boot-a", started, 4, 99.0), RetiredBoot},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := e.ApplyEnvelope(tc.envelope); got != tc.want {
				t.Fatalf("got %s want %s", got, tc.want)
			}
		})
	}
	results := e.FindNearest("tenant-1", 13.3, 80, 1, pb.NodeType_NODE_TYPE_DRONE)
	if len(results) != 1 || results[0].NodeID != "device-1" {
		t.Fatalf("latest accepted state missing: %#v", results)
	}
}
```

---

## internal\application\stream\archiver.go

```
package stream

import (
	"context"
	"fmt"
	"log/slog"
	"sync/atomic"
	"time"

	"github.com/Akashpg-M/polaris/backend/internal/core/events"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/segmentio/kafka-go"
)

const KafkaTelemetryTopic = "telemetry.ingress"
const DeadLetterTopic = "telemetry.dead-letter.v1"

type KafkaPostgresArchiver struct {
	reader       *kafka.Reader
	writer       *kafka.Writer
	db           *sqlx.DB
	done         chan struct{}
	maxRetries   int
	lastProgress atomic.Int64
}

func NewKafkaPostgresArchiver(brokerURL, postgresURL string) (*KafkaPostgresArchiver, error) {
	db, err := sqlx.Connect("postgres", postgresURL)
	if err != nil {
		return nil, err
	}
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(5)
	db.SetConnMaxIdleTime(5 * time.Minute)
	db.SetConnMaxLifetime(30 * time.Minute)
	a := &KafkaPostgresArchiver{
		reader: kafka.NewReader(kafka.ReaderConfig{Brokers: []string{brokerURL}, Topic: KafkaTelemetryTopic, GroupID: "polaris_archive_group", CommitInterval: 0}),
		writer: &kafka.Writer{Addr: kafka.TCP(brokerURL), Topic: DeadLetterTopic, Balancer: &kafka.Hash{}},
		db:     db, done: make(chan struct{}), maxRetries: 5,
	}
	a.lastProgress.Store(time.Now().UnixMilli())
	return a, nil
}

func (a *KafkaPostgresArchiver) archive(ctx context.Context, e *events.TelemetryEnvelope) error {
	p := e.Payload
	_, err := a.db.ExecContext(ctx, `
		INSERT INTO telemetry_history
		(event_id, tenant_id, device_id, device_boot_id, sequence_number, asset_type,
		 lat, lon, geom, status, velocity_mps, heading_deg, battery, observed_at,
		 ingested_at, schema_version, correlation_id, recorded_at)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,ST_SetSRID(ST_MakePoint($8,$7),4326),$9,$10,$11,$12,$13::timestamptz,$14::timestamptz,$15,$16,($13::timestamptz AT TIME ZONE 'UTC'))
		ON CONFLICT DO NOTHING`,
		e.EventID, e.TenantID, e.DeviceID, e.DeviceBootID, e.SequenceNumber, int(p.Type),
		p.Lat, p.Lon, int(p.Status), p.VelocityMps, p.HeadingDeg, p.EnergyPercent,
		time.UnixMilli(e.ObservedAt), time.UnixMilli(e.IngestedAt), e.SchemaVersion, e.CorrelationID)
	return err
}

func (a *KafkaPostgresArchiver) sendToDLQ(ctx context.Context, msg kafka.Message, reason string) error {
	return a.writer.WriteMessages(ctx, kafka.Message{Key: msg.Key, Value: msg.Value, Headers: []kafka.Header{
		{Key: "error_reason", Value: []byte(reason)}, {Key: "source_topic", Value: []byte(msg.Topic)},
		{Key: "source_partition", Value: []byte(fmt.Sprint(msg.Partition))}, {Key: "source_offset", Value: []byte(fmt.Sprint(msg.Offset))},
		{Key: "failed_at", Value: []byte(time.Now().UTC().Format(time.RFC3339Nano))},
	}})
}

func (a *KafkaPostgresArchiver) process(ctx context.Context, msg kafka.Message) bool {
	e, err := events.Unmarshal(msg.Value)
	if err != nil {
		if dlqErr := a.sendToDLQ(ctx, msg, err.Error()); dlqErr != nil {
			slog.Error("archive poison event DLQ failed", "error", dlqErr)
			return false
		}
		return true
	}
	var lastErr error
	for attempt := 1; attempt <= a.maxRetries; attempt++ {
		if err := a.archive(ctx, e); err == nil {
			return true
		} else {
			lastErr = err
		}
		slog.Warn("transient PostgreSQL archive failure", "event_id", e.EventID, "attempt", attempt, "error", lastErr)
		if attempt < a.maxRetries {
			select {
			case <-time.After(time.Duration(attempt*50) * time.Millisecond):
			case <-ctx.Done():
				return false
			}
		}
	}
	if err := a.sendToDLQ(ctx, msg, "retry_exhausted: "+lastErr.Error()); err != nil {
		slog.Error("archive retry-exhausted DLQ failed", "error", err)
		return false
	}
	return true
}

func (a *KafkaPostgresArchiver) Start(ctx context.Context) {
	defer close(a.done)
	defer a.reader.Close()
	defer a.writer.Close()
	defer a.db.Close()
	slog.Info("Idempotent Kafka PostgreSQL archiver active")
	for {
		msg, err := a.reader.FetchMessage(ctx)
		if err != nil {
			if ctx.Err() != nil {
				slog.Info("archive consumer shutdown complete")
				return
			}
			slog.Error("archive fetch failed", "error", err)
			continue
		}
		for !a.process(ctx, msg) {
			select {
			case <-time.After(250 * time.Millisecond):
			case <-ctx.Done():
				slog.Error("archive shutdown with uncommitted message", "partition", msg.Partition, "offset", msg.Offset)
				return
			}
		}
		for {
			if err := a.reader.CommitMessages(ctx, msg); err == nil {
				break
			} else {
				slog.Error("archive Kafka commit failed; retrying same offset", "partition", msg.Partition, "offset", msg.Offset, "error", err)
			}
			select {
			case <-time.After(250 * time.Millisecond):
			case <-ctx.Done():
				return
			}
		}
		a.lastProgress.Store(time.Now().UnixMilli())
	}
}

func (a *KafkaPostgresArchiver) Ready(ctx context.Context) error { return a.db.PingContext(ctx) }
func (a *KafkaPostgresArchiver) Healthy() bool {
	select {
	case <-a.done:
		return false
	default:
		return true
	}
}
func (a *KafkaPostgresArchiver) Wait(ctx context.Context) error {
	select {
	case <-a.done:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}
```

---

## internal\application\stream\archiver_integration_test.go

```
//go:build integration

package stream

import (
	"context"
	"os"
	"testing"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func TestArchiveReplayIsIdempotent(t *testing.T) {
	postgresURL := os.Getenv("POSTGRES_URL")
	if postgresURL == "" {
		postgresURL = "postgres://polaris_user:polaris_password@localhost:5432/polaris_core?sslmode=disable"
	}
	db, err := sqlx.Connect("postgres", postgresURL)
	if err != nil {
		t.Skipf("PostgreSQL integration dependency unavailable: %v", err)
	}
	defer db.Close()
	e := streamEnvelope("archive-replay-integration", 1)
	a := &KafkaPostgresArchiver{db: db}
	ctx := context.Background()
	defer db.ExecContext(ctx, "DELETE FROM telemetry_history WHERE event_id=$1", e.EventID)
	if err := a.archive(ctx, e); err != nil {
		t.Fatal(err)
	}
	if err := a.archive(ctx, e); err != nil {
		t.Fatal(err)
	}
	var count int
	if err := db.GetContext(ctx, &count, "SELECT count(*) FROM telemetry_history WHERE event_id=$1", e.EventID); err != nil {
		t.Fatal(err)
	}
	if count != 1 {
		t.Fatalf("duplicate replay produced %d rows", count)
	}
}
```

---

## internal\application\stream\kafka_consumer.go

```
package stream

import (
	"context"
	"fmt"
	"log/slog"
	"sort"
	"sync/atomic"
	"time"

	"github.com/Akashpg-M/polaris/backend/internal/application/spatial"
	"github.com/Akashpg-M/polaris/backend/internal/core/events"
	"github.com/redis/go-redis/v9"
	"github.com/segmentio/kafka-go"
)

const DashboardUpdatesChannel = "spatial:updates"

type pendingTelemetry struct{ message kafka.Message }
type partitionBatch struct {
	items       []pendingTelemetry
	firstQueued time.Time
}

type telemetryReader interface {
	FetchMessage(context.Context) (kafka.Message, error)
	CommitMessages(context.Context, ...kafka.Message) error
	Close() error
}
type telemetryWriter interface {
	WriteMessages(context.Context, ...kafka.Message) error
	Close() error
}
type stateApplier interface {
	ApplyEnvelope(*events.TelemetryEnvelope) spatial.Classification
}
type latestProjector interface {
	Apply(context.Context, *events.TelemetryEnvelope) (spatial.Classification, error)
	Ready(context.Context) error
}

type KafkaConsumer struct {
	reader        telemetryReader
	dlq           telemetryWriter
	brokerURL     string
	engine        stateApplier
	projector     latestProjector
	batchSize     int
	flushInterval time.Duration
	maxRetries    int
	done          chan struct{}
	lastProgress  atomic.Int64
}

func NewKafkaConsumer(brokerURL string, engine stateApplier, redisClient *redis.Client) *KafkaConsumer {
	c := &KafkaConsumer{
		reader:    kafka.NewReader(kafka.ReaderConfig{Brokers: []string{brokerURL}, Topic: KafkaTelemetryTopic, GroupID: "polaris_engine_group", CommitInterval: 0}),
		dlq:       &kafka.Writer{Addr: kafka.TCP(brokerURL), Topic: DeadLetterTopic, Balancer: &kafka.Hash{}},
		brokerURL: brokerURL, engine: engine, projector: NewRedisProjector(redisClient),
		batchSize: 1000, flushInterval: 150 * time.Millisecond, maxRetries: 5, done: make(chan struct{}),
	}
	c.lastProgress.Store(time.Now().UnixMilli())
	return c
}

func partitionID(message kafka.Message) string {
	return fmt.Sprintf("%s:%d", message.Topic, message.Partition)
}

func (c *KafkaConsumer) sendToDLQ(ctx context.Context, msg kafka.Message, reason string) error {
	return c.dlq.WriteMessages(ctx, kafka.Message{Key: msg.Key, Value: msg.Value, Headers: []kafka.Header{
		{Key: "error_reason", Value: []byte(reason)}, {Key: "source_topic", Value: []byte(msg.Topic)},
		{Key: "source_partition", Value: []byte(fmt.Sprint(msg.Partition))}, {Key: "source_offset", Value: []byte(fmt.Sprint(msg.Offset))},
		{Key: "failed_at", Value: []byte(time.Now().UTC().Format(time.RFC3339Nano))},
	}})
}

func (c *KafkaConsumer) process(ctx context.Context, item pendingTelemetry) bool {
	envelope, err := events.Unmarshal(item.message.Value)
	if err != nil {
		if dlqErr := c.sendToDLQ(ctx, item.message, err.Error()); dlqErr != nil {
			slog.Error("permanent telemetry failure could not reach DLQ", "partition", item.message.Partition, "offset", item.message.Offset, "error", dlqErr)
			return false
		}
		slog.Warn("permanent telemetry failure sent to DLQ", "partition", item.message.Partition, "offset", item.message.Offset, "error", err)
		return true
	}
	var lastErr error
	for attempt := 1; attempt <= c.maxRetries; attempt++ {
		redisClass, projectionErr := c.projector.Apply(ctx, envelope)
		if projectionErr == nil {
			memoryClass := spatial.Classification("NOT_APPLIED")
			// A Redis DUPLICATE may be the replay that rebuilds an empty engine
			// after restart. Stale/retired/conflicting events must never enter it.
			if redisClass == spatial.Accepted || redisClass == spatial.NewBoot || redisClass == spatial.Duplicate {
				memoryClass = c.engine.ApplyEnvelope(envelope)
			}
			// Per-event classification is useful for diagnosis but is too noisy for
			// the steady-state INFO path under fleet load.
			slog.Debug("telemetry state classified", "event_id", envelope.EventID, "spatial", memoryClass, "redis", redisClass)
			return true
		}
		lastErr = projectionErr
		slog.Warn("transient Redis projection failure", "event_id", envelope.EventID, "attempt", attempt, "error", projectionErr)
		if attempt < c.maxRetries {
			select {
			case <-time.After(time.Duration(attempt*50) * time.Millisecond):
			case <-ctx.Done():
				return false
			}
		}
	}
	if err := c.sendToDLQ(ctx, item.message, "retry_exhausted: "+lastErr.Error()); err != nil {
		slog.Error("retry-exhausted telemetry could not reach DLQ", "offset", item.message.Offset, "error", err)
		return false
	}
	return true
}

// flushPartition commits only the highest contiguous successfully processed offset.
func (c *KafkaConsumer) flushPartition(ctx context.Context, key string, batch *partitionBatch) bool {
	if len(batch.items) == 0 {
		return true
	}
	sort.Slice(batch.items, func(i, j int) bool { return batch.items[i].message.Offset < batch.items[j].message.Offset })
	slog.Info("partition batch flush started", "partition", key, "messages", len(batch.items), "queue_age_ms", time.Since(batch.firstQueued).Milliseconds())
	succeeded := 0
	for _, item := range batch.items {
		if !c.process(ctx, item) {
			break
		}
		succeeded++
	}
	if succeeded == 0 {
		return false
	}
	highest := batch.items[succeeded-1].message
	if err := c.reader.CommitMessages(ctx, highest); err != nil {
		slog.Error("Kafka commit failed; successful effects will replay", "partition", key, "offset", highest.Offset, "error", err)
		return false
	}
	c.lastProgress.Store(time.Now().UnixMilli())
	batch.items = batch.items[succeeded:]
	if len(batch.items) == 0 {
		batch.firstQueued = time.Time{}
	} else {
		batch.firstQueued = time.Now()
	}
	slog.Info("partition batch committed", "partition", key, "messages", succeeded, "highest_offset", highest.Offset)
	return true
}

func (c *KafkaConsumer) Start(ctx context.Context, workerID string) {
	defer close(c.done)
	defer c.reader.Close()
	defer c.dlq.Close()
	slog.Info("Kafka partition-aware spatial consumer started", "worker_id", workerID, "batch_size", c.batchSize, "flush_interval", c.flushInterval)
	ticker := time.NewTicker(5 * time.Millisecond)
	defer ticker.Stop()
	batches := make(map[string]*partitionBatch)
	fetched := make(chan kafka.Message)
	fetchErrors := make(chan error, 1)
	go func() {
		defer close(fetched)
		for {
			message, err := c.reader.FetchMessage(ctx)
			if err != nil {
				if ctx.Err() == nil {
					fetchErrors <- err
				}
				return
			}
			select {
			case fetched <- message:
			case <-ctx.Done():
				return
			}
		}
	}()
	flushDue := func(flushCtx context.Context, force bool) {
		now := time.Now()
		for key, batch := range batches {
			if force || len(batch.items) >= c.batchSize || (!batch.firstQueued.IsZero() && now.Sub(batch.firstQueued) >= c.flushInterval) {
				c.flushPartition(flushCtx, key, batch)
				if len(batch.items) == 0 {
					delete(batches, key)
				}
			}
		}
	}
	shutdownFlush := func() {
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		flushDue(shutdownCtx, true)
		cancel()
		pending := 0
		for _, batch := range batches {
			pending += len(batch.items)
		}
		if pending > 0 {
			slog.Error("spatial consumer stopped with uncommitted messages", "pending", pending)
		} else {
			slog.Info("spatial consumer shutdown flush complete")
		}
	}
	for {
		select {
		case message, ok := <-fetched:
			if !ok {
				shutdownFlush()
				return
			}
			key := partitionID(message)
			batch := batches[key]
			if batch == nil {
				batch = &partitionBatch{firstQueued: time.Now()}
				batches[key] = batch
			}
			batch.items = append(batch.items, pendingTelemetry{message: message})
		case err := <-fetchErrors:
			slog.Error("Kafka fetch loop stopped", "error", err)
			shutdownFlush()
			return
		case <-ticker.C:
			c.lastProgress.Store(time.Now().UnixMilli())
			flushDue(ctx, false)
		case <-ctx.Done():
			shutdownFlush()
			return
		}
	}
}

func (c *KafkaConsumer) Ready(ctx context.Context) error {
	conn, err := kafka.DialContext(ctx, "tcp", c.brokerURL)
	if err != nil {
		return err
	}
	defer conn.Close()
	return c.projector.Ready(ctx)
}
func (c *KafkaConsumer) Healthy() bool {
	select {
	case <-c.done:
		return false
	default:
		return time.Since(time.UnixMilli(c.lastProgress.Load())) < 2*time.Minute
	}
}
func (c *KafkaConsumer) Wait(ctx context.Context) error {
	select {
	case <-c.done:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}
```

---

## internal\application\stream\kafka_consumer_reliability_test.go

```
package stream

import (
	"context"
	"errors"
	"sync"
	"testing"
	"time"

	pb "github.com/Akashpg-M/polaris/backend/api/proto/v1"
	"github.com/Akashpg-M/polaris/backend/internal/application/spatial"
	"github.com/Akashpg-M/polaris/backend/internal/core/events"
	"github.com/segmentio/kafka-go"
)

type fakeState struct {
	mu    sync.Mutex
	calls int
}

func (f *fakeState) ApplyEnvelope(*events.TelemetryEnvelope) spatial.Classification {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.calls++
	return spatial.Accepted
}

type fakeProjection struct {
	mu             sync.Mutex
	failures       int
	calls          int
	classification spatial.Classification
}

func (f *fakeProjection) Apply(context.Context, *events.TelemetryEnvelope) (spatial.Classification, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.calls++
	if f.calls <= f.failures {
		return "", errors.New("redis unavailable")
	}
	return f.classification, nil
}
func (*fakeProjection) Ready(context.Context) error { return nil }

type fakeWriter struct {
	mu       sync.Mutex
	messages []kafka.Message
}

func (f *fakeWriter) WriteMessages(_ context.Context, m ...kafka.Message) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.messages = append(f.messages, m...)
	return nil
}
func (*fakeWriter) Close() error { return nil }

type fakeReader struct {
	mu        sync.Mutex
	source    chan kafka.Message
	commits   []kafka.Message
	commitErr error
}

func (f *fakeReader) FetchMessage(ctx context.Context) (kafka.Message, error) {
	select {
	case m := <-f.source:
		return m, nil
	case <-ctx.Done():
		return kafka.Message{}, ctx.Err()
	}
}
func (f *fakeReader) CommitMessages(_ context.Context, m ...kafka.Message) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	if f.commitErr != nil {
		return f.commitErr
	}
	f.commits = append(f.commits, m...)
	return nil
}
func (*fakeReader) Close() error { return nil }

func streamEnvelope(device string, sequence uint64) *events.TelemetryEnvelope {
	now := time.Now().UTC()
	p := &pb.SpatialObject{Id: device, TenantId: "tenant-1", DeviceBootId: "boot-1", SequenceNumber: sequence,
		BootStartedAt: now.Add(-time.Minute).UnixMilli(), ObservedAt: now.UnixMilli(), SchemaVersion: 1, Type: pb.NodeType_NODE_TYPE_DRONE,
		Status: pb.NodeStatus_NODE_STATUS_ACTIVE, Lat: 13, Lon: 80, EnergyPercent: 50}
	return events.NewTelemetryEnvelope(p, now, "", "", "")
}
func kafkaEnvelope(t *testing.T, partition int, offset int64, device string) kafka.Message {
	t.Helper()
	data, err := streamEnvelope(device, uint64(offset+1)).Marshal()
	if err != nil {
		t.Fatal(err)
	}
	return kafka.Message{Topic: KafkaTelemetryTopic, Partition: partition, Offset: offset, Key: []byte("tenant-1:" + device), Value: data}
}

func TestRedisFailureDoesNotAdvanceSpatialState(t *testing.T) {
	state := &fakeState{}
	projection := &fakeProjection{failures: 2, classification: spatial.Accepted}
	writer := &fakeWriter{}
	c := &KafkaConsumer{engine: state, projector: projection, dlq: writer, maxRetries: 3}
	if !c.process(context.Background(), pendingTelemetry{message: kafkaEnvelope(t, 0, 0, "device-retry")}) {
		t.Fatal("event should succeed after retry")
	}
	if state.calls != 1 {
		t.Fatalf("spatial applied %d times before/after Redis recovery; want once", state.calls)
	}
}

func TestUnsupportedSchemaReachesDLQ(t *testing.T) {
	e := streamEnvelope("device-poison", 1)
	e.SchemaVersion = 99
	data, err := e.Marshal()
	if err != nil {
		t.Fatal(err)
	}
	writer := &fakeWriter{}
	c := &KafkaConsumer{engine: &fakeState{}, projector: &fakeProjection{}, dlq: writer, maxRetries: 1}
	if !c.process(context.Background(), pendingTelemetry{message: kafka.Message{Topic: KafkaTelemetryTopic, Partition: 2, Offset: 7, Value: data}}) {
		t.Fatal("successful DLQ is terminal")
	}
	if len(writer.messages) != 1 || string(writer.messages[0].Value) != string(data) {
		t.Fatal("original event was not preserved in DLQ")
	}
}

func TestCommitFailureLeavesSuccessfulBatchForReplay(t *testing.T) {
	reader := &fakeReader{commitErr: errors.New("broker unavailable")}
	batch := &partitionBatch{firstQueued: time.Now(), items: []pendingTelemetry{{message: kafkaEnvelope(t, 1, 10, "device-commit")}}}
	c := &KafkaConsumer{reader: reader, engine: &fakeState{}, projector: &fakeProjection{classification: spatial.Accepted}, dlq: &fakeWriter{}, maxRetries: 1}
	if c.flushPartition(context.Background(), "telemetry.ingress:1", batch) {
		t.Fatal("commit failure must not report success")
	}
	if len(batch.items) != 1 {
		t.Fatal("successful effects must remain queued for replay when commit fails")
	}
}

func TestGracefulShutdownFlushesPartitionsIndependently(t *testing.T) {
	reader := &fakeReader{source: make(chan kafka.Message, 2)}
	reader.source <- kafkaEnvelope(t, 0, 20, "device-p0")
	reader.source <- kafkaEnvelope(t, 2, 30, "device-p2")
	c := &KafkaConsumer{reader: reader, dlq: &fakeWriter{}, engine: &fakeState{}, projector: &fakeProjection{classification: spatial.Accepted}, batchSize: 1000, flushInterval: time.Hour, maxRetries: 1, done: make(chan struct{})}
	ctx, cancel := context.WithCancel(context.Background())
	go c.Start(ctx, "test")
	time.Sleep(30 * time.Millisecond)
	cancel()
	waitCtx, waitCancel := context.WithTimeout(context.Background(), time.Second)
	defer waitCancel()
	if err := c.Wait(waitCtx); err != nil {
		t.Fatal(err)
	}
	reader.mu.Lock()
	defer reader.mu.Unlock()
	if len(reader.commits) != 2 {
		t.Fatalf("committed %d partitions; want 2", len(reader.commits))
	}
	seen := map[int]bool{}
	for _, m := range reader.commits {
		seen[m.Partition] = true
	}
	if !seen[0] || !seen[2] {
		t.Fatalf("partition commits not independent: %#v", seen)
	}
}
```

---

## internal\application\stream\redis_projection.go

```
package stream

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/Akashpg-M/polaris/backend/internal/application/spatial"
	"github.com/Akashpg-M/polaris/backend/internal/core/events"
	twincore "github.com/Akashpg-M/polaris/backend/internal/core/twin"
	"github.com/redis/go-redis/v9"
)

var latestStateScript = redis.NewScript(`
local current_boot = redis.call('HGET', KEYS[1], 'device_boot_id')
local incoming_boot = ARGV[1]
local incoming_started = tonumber(ARGV[2])
local incoming_seq = ARGV[3]
local function compare_decimal(a, b)
  a = string.gsub(a, '^0+', ''); if a == '' then a = '0' end
  b = string.gsub(b, '^0+', ''); if b == '' then b = '0' end
  if string.len(a) < string.len(b) then return -1 end
  if string.len(a) > string.len(b) then return 1 end
  if a < b then return -1 end
  if a > b then return 1 end
  return 0
end
local classification = 'ACCEPTED'
if current_boot then
  local current_started = tonumber(redis.call('HGET', KEYS[1], 'boot_started_at'))
  local current_seq = redis.call('HGET', KEYS[1], 'sequence_number')
  if current_boot == incoming_boot then
    local sequence_order = compare_decimal(incoming_seq, current_seq)
    if sequence_order == 0 then return 'DUPLICATE' end
    if sequence_order < 0 then return 'OUT_OF_ORDER' end
  elseif redis.call('SISMEMBER', KEYS[2], incoming_boot) == 1 then
    return 'RETIRED_BOOT'
  elseif incoming_started > current_started then
    redis.call('SADD', KEYS[2], current_boot)
    classification = 'NEW_BOOT'
  elseif incoming_started == current_started then
    return 'BOOT_CONFLICT'
  else
    return 'RETIRED_BOOT'
  end
end
redis.call('HSET', KEYS[1],
  'device_boot_id', incoming_boot,
  'boot_started_at', ARGV[2],
  'sequence_number', ARGV[3],
  'event_id', ARGV[4],
  'reported_state', ARGV[5],
  'last_seen_at', ARGV[6],
  'connectivity_status', 'ONLINE',
  'component:spatial/v1', ARGV[9],
  'component:battery/v1', ARGV[10])
redis.call('ZADD', KEYS[3], ARGV[6], ARGV[8])
redis.call('PUBLISH', ARGV[7], ARGV[5])
return classification
`)

type RedisProjector struct{ client *redis.Client }

func NewRedisProjector(client *redis.Client) *RedisProjector { return &RedisProjector{client: client} }

func dashboardEnvelopeJSON(e *events.TelemetryEnvelope) ([]byte, error) {
	p := e.Payload
	return json.Marshal(map[string]interface{}{
		"event_id": e.EventID, "schema_version": e.SchemaVersion,
		"id": p.Id, "device_id": e.DeviceID, "tenant_id": e.TenantID,
		"device_boot_id": e.DeviceBootID, "sequence_number": e.SequenceNumber,
		"boot_started_at": e.BootStartedAt,
		"type":            p.Type, "status": p.Status, "lat": p.Lat, "lon": p.Lon,
		"velocity_mps": p.VelocityMps, "heading_deg": p.HeadingDeg,
		"energy_percent": p.EnergyPercent, "observed_at": e.ObservedAt,
		"ingested_at": e.IngestedAt, "timestamp": e.ObservedAt,
	})
}

func (p *RedisProjector) Apply(ctx context.Context, e *events.TelemetryEnvelope) (spatial.Classification, error) {
	data, err := dashboardEnvelopeJSON(e)
	if err != nil {
		return "", err
	}
	key := fmt.Sprintf("polaris:twin:%s:%s", e.TenantID, e.DeviceID)
	mobilityProfile := "ROAD_VEHICLE"
	if e.Payload.Type == 5 {
		mobilityProfile = "AERIAL_DRONE"
	} else if e.Payload.Type == 6 {
		mobilityProfile = "GROUND_ROBOT"
	} else if e.Payload.Type == 7 {
		mobilityProfile = "STATIC"
	}
	spatialPayload, _ := json.Marshal(map[string]interface{}{"latitude": e.Payload.Lat, "longitude": e.Payload.Lon, "heading_degrees": e.Payload.HeadingDeg, "speed_mps": e.Payload.VelocityMps, "mobility_profile": mobilityProfile})
	batteryPayload, _ := json.Marshal(map[string]interface{}{"percent": e.Payload.EnergyPercent})
	observedAt := time.UnixMilli(e.ObservedAt).UTC()
	spatialComponent, _ := json.Marshal(twincore.ComponentEnvelope{Type: "spatial/v1", SchemaVersion: 1, ObservedAt: observedAt, BootID: e.DeviceBootID, SequenceNumber: e.SequenceNumber, Payload: spatialPayload})
	batteryComponent, _ := json.Marshal(twincore.ComponentEnvelope{Type: "battery/v1", SchemaVersion: 1, ObservedAt: observedAt, BootID: e.DeviceBootID, SequenceNumber: e.SequenceNumber, Payload: batteryPayload})
	result, err := latestStateScript.Run(ctx, p.client, []string{key, key + ":retired_boots", "polaris:devices:last-seen"},
		e.DeviceBootID, e.BootStartedAt, e.SequenceNumber, e.EventID, string(data), time.Now().UTC().UnixMilli(), DashboardUpdatesChannel, e.TenantID+":"+e.DeviceID, string(spatialComponent), string(batteryComponent)).Text()
	return spatial.Classification(result), err
}

func (p *RedisProjector) Ready(ctx context.Context) error { return p.client.Ping(ctx).Err() }
```

---

## internal\application\stream\redis_projection_integration_test.go

```
//go:build integration

package stream

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/Akashpg-M/polaris/backend/internal/application/spatial"
	"github.com/redis/go-redis/v9"
)

func TestRedisProjectionAtomicClassifications(t *testing.T) {
	redisURL := os.Getenv("REDIS_URL")
	if redisURL == "" {
		redisURL = "redis://localhost:6379/0"
	}
	opts, err := redis.ParseURL(redisURL)
	if err != nil {
		t.Fatal(err)
	}
	client := redis.NewClient(opts)
	defer client.Close()
	ctx := context.Background()
	if err := client.Ping(ctx).Err(); err != nil {
		t.Skipf("Redis integration dependency unavailable: %v", err)
	}
	device := fmt.Sprintf("integration-%d", time.Now().UnixNano())
	key := "polaris:twin:tenant-1:" + device
	defer client.Del(ctx, key, key+":retired_boots")
	p := NewRedisProjector(client)
	started := time.Now().Add(-time.Minute).UnixMilli()
	first := streamEnvelope(device, 1)
	first.BootStartedAt = started
	first.Payload.BootStartedAt = started
	if got, err := p.Apply(ctx, first); err != nil || got != spatial.Accepted {
		t.Fatalf("first: %s %v", got, err)
	}
	if got, _ := p.Apply(ctx, first); got != spatial.Duplicate {
		t.Fatalf("duplicate: %s", got)
	}
	older := streamEnvelope(device, 1)
	older.SequenceNumber = 0
	older.Payload.SequenceNumber = 0
	if got, _ := p.Apply(ctx, older); got != spatial.OutOfOrder {
		t.Fatalf("out of order: %s", got)
	}
	newBoot := streamEnvelope(device, 1)
	newBoot.DeviceBootID = "boot-2"
	newBoot.Payload.DeviceBootId = "boot-2"
	newBoot.BootStartedAt = started + 1000
	newBoot.Payload.BootStartedAt = started + 1000
	if got, _ := p.Apply(ctx, newBoot); got != spatial.NewBoot {
		t.Fatalf("new boot: %s", got)
	}
	if got, _ := p.Apply(ctx, first); got != spatial.RetiredBoot {
		t.Fatalf("retired boot: %s", got)
	}
	conflict := streamEnvelope(device, 1)
	conflict.DeviceBootID = "boot-conflict"
	conflict.Payload.DeviceBootId = "boot-conflict"
	conflict.BootStartedAt = newBoot.BootStartedAt
	conflict.Payload.BootStartedAt = newBoot.BootStartedAt
	if got, _ := p.Apply(ctx, conflict); got != spatial.BootConflict {
		t.Fatalf("boot conflict: %s", got)
	}
	large := streamEnvelope(device, 9_007_199_254_740_993)
	large.DeviceBootID = "boot-2"
	large.Payload.DeviceBootId = "boot-2"
	large.BootStartedAt = newBoot.BootStartedAt
	large.Payload.BootStartedAt = newBoot.BootStartedAt
	if got, _ := p.Apply(ctx, large); got != spatial.Accepted {
		t.Fatalf("large sequence: %s", got)
	}
	largeOlder := streamEnvelope(device, 9_007_199_254_740_992)
	largeOlder.DeviceBootID = "boot-2"
	largeOlder.Payload.DeviceBootId = "boot-2"
	largeOlder.BootStartedAt = newBoot.BootStartedAt
	largeOlder.Payload.BootStartedAt = newBoot.BootStartedAt
	if got, _ := p.Apply(ctx, largeOlder); got != spatial.OutOfOrder {
		t.Fatalf("large out of order: %s", got)
	}
	values, err := client.HMGet(ctx, key, "device_boot_id", "sequence_number", "event_id", "reported_state", "last_seen_at").Result()
	if err != nil {
		t.Fatal(err)
	}
	if values[0] != "boot-2" || values[1] != "9007199254740993" || values[2] == nil || values[3] == nil || values[4] == nil {
		t.Fatalf("incomplete atomic twin hash: %#v", values)
	}
}
```

---

## internal\application\stream\state_fanout.go

```
package stream

import (
	"github.com/Akashpg-M/polaris/backend/internal/application/spatial"
	"github.com/Akashpg-M/polaris/backend/internal/core/events"
)

type StateFanout struct {
	Primary     stateApplier
	Projections []stateApplier
}

func (f *StateFanout) ApplyEnvelope(e *events.TelemetryEnvelope) spatial.Classification {
	result := f.Primary.ApplyEnvelope(e)
	for _, p := range f.Projections {
		p.ApplyEnvelope(e)
	}
	return result
}
```

---

## internal\application\twin\connectivity.go

```
package twin

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"strconv"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/segmentio/kafka-go"
)

var transitionScript = redis.NewScript(`
local seen=redis.call('HGET',KEYS[1],'last_seen_at')
if not seen or seen~=ARGV[1] then return 0 end
local current=redis.call('HGET',KEYS[1],'connectivity_status')
if current==ARGV[2] then return 0 end
if ARGV[2]=='STALE' and current~='ONLINE' then return 0 end
if ARGV[2]=='OFFLINE' and current~='ONLINE' and current~='STALE' then return 0 end
redis.call('HSET',KEYS[1],'connectivity_status',ARGV[2])
return 1`)

type Detector struct {
	redis                    *redis.Client
	writer                   *kafka.Writer
	stale, offline, interval time.Duration
	onTransition             func(tenant, device, status string)
}

func (d *Detector) SetTransitionHandler(fn func(tenant, device, status string)) { d.onTransition = fn }

func NewDetector(client *redis.Client, broker string, stale, offline, interval time.Duration) *Detector {
	return &Detector{redis: client, writer: &kafka.Writer{Addr: kafka.TCP(broker), Topic: "device.connectivity.v1", Balancer: &kafka.Hash{}}, stale: stale, offline: offline, interval: interval}
}
func (d *Detector) Start(ctx context.Context) {
	defer d.writer.Close()
	ticker := time.NewTicker(d.interval)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			d.scan(ctx, "OFFLINE", d.offline)
			d.scan(ctx, "STALE", d.stale)
		case <-ctx.Done():
			return
		}
	}
}
func (d *Detector) scan(ctx context.Context, next string, age time.Duration) {
	cutoff := time.Now().Add(-age).UnixMilli()
	entries, err := d.redis.ZRangeByScoreWithScores(ctx, "polaris:devices:last-seen", &redis.ZRangeBy{Min: "-inf", Max: strconv.FormatInt(cutoff, 10), Offset: 0, Count: 100}).Result()
	if err != nil {
		return
	}
	for _, z := range entries {
		member, ok := z.Member.(string)
		if !ok {
			continue
		}
		parts := strings.SplitN(member, ":", 2)
		if len(parts) != 2 {
			continue
		}
		last := strconv.FormatInt(int64(z.Score), 10)
		key := "polaris:twin:" + member
		changed, err := transitionScript.Run(ctx, d.redis, []string{key}, last, next).Int()
		if err != nil || changed != 1 {
			continue
		}
		if d.onTransition != nil {
			d.onTransition(parts[0], parts[1], next)
		}
		sum := sha256.Sum256([]byte(member + ":" + next + ":" + last))
		eventID := hex.EncodeToString(sum[:])
		payload, _ := json.Marshal(map[string]interface{}{"event_id": eventID, "event_type": "device.connectivity.changed.v1", "schema_version": 1, "tenant_id": parts[0], "device_id": parts[1], "connectivity_status": next, "last_seen_at": last})
		_ = d.writer.WriteMessages(ctx, kafka.Message{Key: []byte(member), Value: payload})
	}
}
```

---

## internal\core\auth\auth.go

```
package auth

import (
	"crypto/rand"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/hex"
	"errors"
	"strings"
)

type Role string

const (
	PlatformAdmin Role = "PLATFORM_ADMIN"
	TenantAdmin   Role = "TENANT_ADMIN"
	Operator      Role = "OPERATOR"
	Viewer        Role = "VIEWER"
)

type DevicePrincipal struct{ TenantID, DeviceID, CredentialID, DeviceType, ProjectID string }
type OperatorPrincipal struct {
	APIKeyID, TenantID string
	Role               Role
}

var ErrInvalidCredential = errors.New("invalid credential")

func NewID() string {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		panic(err)
	}
	b[6] = (b[6] & 0x0f) | 0x40
	b[8] = (b[8] & 0x3f) | 0x80
	h := hex.EncodeToString(b)
	return h[0:8] + "-" + h[8:12] + "-" + h[12:16] + "-" + h[16:20] + "-" + h[20:32]
}

func GenerateToken(kind string) (raw, prefix string, hash []byte, err error) {
	public := make([]byte, 8)
	secret := make([]byte, 32)
	if _, err = rand.Read(public); err != nil {
		return
	}
	if _, err = rand.Read(secret); err != nil {
		return
	}
	prefix = hex.EncodeToString(public)
	raw = "pol_" + kind + "_" + prefix + "." + hex.EncodeToString(secret)
	sum := sha256.Sum256([]byte(raw))
	hash = sum[:]
	return
}
func TokenPrefix(raw string) (string, error) {
	parts := strings.Split(raw, ".")
	if len(parts) != 2 {
		return "", ErrInvalidCredential
	}
	head := parts[0]
	i := strings.LastIndex(head, "_")
	if i < 0 || i == len(head)-1 {
		return "", ErrInvalidCredential
	}
	return head[i+1:], nil
}
func Verify(raw string, expected []byte) bool {
	sum := sha256.Sum256([]byte(raw))
	return len(expected) == len(sum) && subtle.ConstantTimeCompare(sum[:], expected) == 1
}
func Hash(raw string) []byte { sum := sha256.Sum256([]byte(raw)); return sum[:] }

func Can(role Role, permission string) bool {
	if role == PlatformAdmin {
		return true
	}
	switch permission {
	case "read":
		return role == TenantAdmin || role == Operator || role == Viewer
	case "mutate":
		return role == TenantAdmin
	case "orchestrate":
		return role == TenantAdmin || role == Operator
	case "admin_retry":
		return role == TenantAdmin
	case "audit":
		return role == TenantAdmin
	default:
		return false
	}
}
```

---

## internal\core\auth\auth_test.go

```
package auth

import "testing"

func TestGeneratedTokenIsHashedAndVerifiable(t *testing.T) {
	raw, prefix, hash, err := GenerateToken("dev")
	if err != nil {
		t.Fatal(err)
	}
	if len(hash) != 32 || len(raw) < 80 || prefix == "" {
		t.Fatalf("weak token shape: %d %d", len(raw), len(hash))
	}
	parsed, err := TokenPrefix(raw)
	if err != nil || parsed != prefix {
		t.Fatalf("prefix: %q %v", parsed, err)
	}
	if !Verify(raw, hash) {
		t.Fatal("valid token did not verify")
	}
	if Verify(raw+"x", hash) {
		t.Fatal("modified token verified")
	}
}
func TestRolePermissions(t *testing.T) {
	if !Can(PlatformAdmin, "mutate") || !Can(TenantAdmin, "mutate") || Can(Operator, "mutate") || Can(Viewer, "mutate") {
		t.Fatal("mutation permission matrix violated")
	}
	if !Can(Viewer, "read") || Can(Viewer, "audit") {
		t.Fatal("viewer permission matrix violated")
	}
	if !Can(Operator, "orchestrate") || Can(Viewer, "orchestrate") || Can(Operator, "admin_retry") || !Can(TenantAdmin, "admin_retry") {
		t.Fatal("orchestration permission matrix violated")
	}
}
```

---

## internal\core\command\model.go

```
package command

import (
	"encoding/json"
	"time"
)

const SchemaVersion = 1

type Status string

const (
	Pending      Status = "PENDING"
	Delivered    Status = "DELIVERED"
	Acknowledged Status = "ACKNOWLEDGED"
	Completed    Status = "COMPLETED"
	Failed       Status = "FAILED"
	Expired      Status = "EXPIRED"
	Cancelled    Status = "CANCELLED"
)

type Envelope struct {
	FrameType      string          `json:"frame_type"`
	CommandID      string          `json:"command_id"`
	CommandType    string          `json:"command_type"`
	SchemaVersion  int             `json:"schema_version"`
	TenantID       string          `json:"tenant_id"`
	DeviceID       string          `json:"device_id"`
	TaskID         string          `json:"task_id"`
	SequenceNumber int64           `json:"sequence_number"`
	CreatedAt      time.Time       `json:"created_at"`
	ExpiresAt      time.Time       `json:"expires_at"`
	CorrelationID  string          `json:"correlation_id"`
	CausationID    string          `json:"causation_id"`
	Payload        json.RawMessage `json:"payload"`
	// DeliveryObservation is volatile timing evidence added after the durable
	// command decision is read from the outbox. It is never persisted as part
	// of command identity and may differ across at-least-once delivery attempts.
	DeliveryObservation *DeliveryObservation `json:"delivery_observation,omitempty"`
}

type DeliveryObservation struct {
	RelayPublishedAt  time.Time `json:"relay_published_at"`
	GatewayReceivedAt time.Time `json:"gateway_received_at,omitempty"`
}

type Record struct {
	CommandID      string           `db:"command_id" json:"command_id"`
	TenantID       string           `db:"tenant_id" json:"tenant_id"`
	DeviceID       string           `db:"device_id" json:"device_id"`
	TaskID         string           `db:"task_id" json:"task_id"`
	CommandType    string           `db:"command_type" json:"command_type"`
	Payload        json.RawMessage  `db:"payload" json:"payload"`
	Status         string           `db:"status" json:"status"`
	SequenceNumber int64            `db:"sequence_number" json:"sequence_number"`
	CorrelationID  string           `db:"correlation_id" json:"correlation_id"`
	CausationID    string           `db:"causation_id" json:"causation_id"`
	AttemptCount   int              `db:"attempt_count" json:"attempt_count"`
	MaxAttempts    int              `db:"max_attempts" json:"max_attempts"`
	Version        int64            `db:"version" json:"version"`
	CreatedAt      time.Time        `db:"created_at" json:"created_at"`
	AvailableAt    time.Time        `db:"available_at" json:"available_at"`
	SentAt         *time.Time       `db:"sent_at" json:"sent_at,omitempty"`
	AcknowledgedAt *time.Time       `db:"acknowledged_at" json:"acknowledged_at,omitempty"`
	CompletedAt    *time.Time       `db:"completed_at" json:"completed_at,omitempty"`
	ExpiresAt      time.Time        `db:"expires_at" json:"expires_at"`
	AckStatus      *string          `db:"ack_status" json:"ack_status,omitempty"`
	Result         *json.RawMessage `db:"result" json:"result,omitempty"`
	LastError      *string          `db:"last_error" json:"last_error,omitempty"`
}

func (r Record) Envelope() Envelope {
	return Envelope{FrameType: "COMMAND", CommandID: r.CommandID, CommandType: r.CommandType, SchemaVersion: SchemaVersion, TenantID: r.TenantID, DeviceID: r.DeviceID, TaskID: r.TaskID, SequenceNumber: r.SequenceNumber, CreatedAt: r.CreatedAt, ExpiresAt: r.ExpiresAt, CorrelationID: r.CorrelationID, CausationID: r.CausationID, Payload: r.Payload}
}

func (e Envelope) PartitionKey() string { return e.TenantID + ":" + e.DeviceID }

type Ack struct {
	FrameType      string    `json:"frame_type"`
	CommandID      string    `json:"command_id"`
	SequenceNumber int64     `json:"sequence_number"`
	Status         string    `json:"status"`
	ReceivedAt     time.Time `json:"received_at"`
	Reason         string    `json:"reason,omitempty"`
}

type Result struct {
	FrameType      string          `json:"frame_type"`
	CommandID      string          `json:"command_id"`
	SequenceNumber int64           `json:"sequence_number"`
	Status         string          `json:"status"`
	CompletedAt    time.Time       `json:"completed_at"`
	Result         json.RawMessage `json:"result,omitempty"`
	Reason         string          `json:"reason,omitempty"`
}

func ValidTransition(from, to Status) bool {
	if from == to {
		return true
	}
	switch from {
	case Pending:
		return to == Delivered || to == Cancelled || to == Expired || to == Failed
	case Delivered:
		return to == Acknowledged || to == Pending || to == Expired || to == Failed
	case Acknowledged:
		return to == Completed || to == Failed
	}
	return false
}

func IsTerminal(status Status) bool {
	return status == Completed || status == Failed || status == Expired || status == Cancelled
}

func RequiredCapability(commandType string) string {
	switch commandType {
	case "RELOCATE":
		return "receive_relocation_command"
	case "NAVIGATE", "RETURN_HOME":
		return "navigate"
	case "CAPTURE_IMAGE":
		return "capture_image"
	case "RUN_MODEL":
		return "run_model"
	case "THERMAL_SCAN", "START_SCAN":
		return "thermal_scan"
	case "STOP":
		return "receive_relocation_command"
	default:
		return ""
	}
}
```

---

## internal\core\command\model_test.go

```
package command

import (
	"encoding/json"
	"testing"
	"time"
)

func TestCommandTransitionsAndCapabilities(t *testing.T) {
	if !ValidTransition(Pending, Delivered) || !ValidTransition(Delivered, Acknowledged) || !ValidTransition(Acknowledged, Completed) {
		t.Fatal("legal command lifecycle rejected")
	}
	if ValidTransition(Completed, Pending) || ValidTransition(Pending, Acknowledged) {
		t.Fatal("illegal command lifecycle accepted")
	}
	if RequiredCapability("RELOCATE") != "receive_relocation_command" || RequiredCapability("CAPTURE_IMAGE") != "capture_image" {
		t.Fatal("command capability mapping is incorrect")
	}
}

func TestDeliveryObservationIsVolatile(t *testing.T) {
	record := Record{CommandID: "command", TenantID: "tenant", DeviceID: "device", Payload: json.RawMessage(`{}`)}
	envelope := record.Envelope()
	if envelope.DeliveryObservation != nil {
		t.Fatal("durable record unexpectedly retained volatile delivery timing")
	}
	envelope.DeliveryObservation = &DeliveryObservation{RelayPublishedAt: time.Now().UTC()}
	if replay := record.Envelope(); replay.DeliveryObservation != nil || replay.CommandID != envelope.CommandID || string(replay.Payload) != string(envelope.Payload) {
		t.Fatal("delivery observation mutated durable command identity")
	}
}
```

---

## internal\core\domain\node.go

```
package domain

import "time"

// NodeStatus represents the universal FSM states for any IoT device
type NodeStatus string

const (
	StatusIdle        NodeStatus = "idle"         // Ready for orchestration
	StatusEnRoute     NodeStatus = "en_route"     // Actively moving to a target
	StatusActive      NodeStatus = "active"       // Performing its task (e.g., drone delivering)
	StatusMaintenance NodeStatus = "maintenance"  // Charging or broken
	StatusOffline     NodeStatus = "offline"      // Disconnected
)

// AssetClass uses Bitmasking to allow ultra-fast hardware-level filtering.
// By using uint16, we can support up to 16 distinct device categories natively.
type AssetClass uint16

const (
	ClassBike   AssetClass = 1 << 0 // 1  (Ground)
	ClassAuto   AssetClass = 1 << 1 // 2  (Ground)
	ClassSedan  AssetClass = 1 << 2 // 4  (Ground)
	ClassSUV    AssetClass = 1 << 3 // 8  (Ground)
	ClassDrone  AssetClass = 1 << 4 // 16 (Aerial - Bypasses street routing)
	ClassRobot  AssetClass = 1 << 5 // 32 (Ground - Warehouse/Sidewalk)
	ClassSensor AssetClass = 1 << 6 // 64 (Static - e.g., Smart Traffic Light)
)

// TelemetryPayload is the generic JSON structure expected from ANY device pinging the server.
type TelemetryPayload struct {
	TenantID  string                 `json:"tenant_id"` 
	NodeID    string                 `json:"node_id"`
	Class     AssetClass             `json:"asset_class"`
	Lat       float64                `json:"lat"`
	Lon       float64                `json:"lon"`
	Status    NodeStatus             `json:"status"`
	Battery   int                    `json:"battery,omitempty"`
	Metadata  map[string]interface{} `json:"metadata,omitempty"` // For custom device data
	Timestamp time.Time              `json:"timestamp"`
}
```

---

## internal\core\events\telemetry.go

```
package events

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"regexp"
	"strings"
	"time"

	pb "github.com/Akashpg-M/polaris/backend/api/proto/v1"
)

const (
	TelemetryEventType   = "polaris.telemetry.observed"
	CurrentSchemaVersion = uint32(1)
	GatewayProducer      = "polaris-gateway"
	MaxFrameBytes        = int64(64 * 1024)
)

// ':' is reserved as the unambiguous tenant/device partition-key separator.
var identityPattern = regexp.MustCompile(`^[A-Za-z0-9][A-Za-z0-9._-]{0,127}$`)

// TelemetryEnvelope is the canonical Kafka value. Device-owned facts stay in
// Payload; trusted platform metadata is added at the gateway boundary.
type TelemetryEnvelope struct {
	EventID        string            `json:"event_id"`
	EventType      string            `json:"event_type"`
	SchemaVersion  uint32            `json:"schema_version"`
	TenantID       string            `json:"tenant_id"`
	DeviceID       string            `json:"device_id"`
	DeviceBootID   string            `json:"device_boot_id"`
	SequenceNumber uint64            `json:"sequence_number"`
	BootStartedAt  int64             `json:"boot_started_at"`
	ObservedAt     int64             `json:"observed_at"`
	IngestedAt     int64             `json:"ingested_at"`
	CorrelationID  string            `json:"correlation_id"`
	CausationID    string            `json:"causation_id,omitempty"`
	Producer       string            `json:"producer"`
	Traceparent    string            `json:"traceparent,omitempty"`
	Payload        *pb.SpatialObject `json:"payload"`
}

func NewTelemetryEnvelope(p *pb.SpatialObject, ingestedAt time.Time, correlationID, causationID, traceparent string) *TelemetryEnvelope {
	observedAt := p.ObservedAt
	if observedAt == 0 {
		observedAt = p.Timestamp
	}
	identity := fmt.Sprintf("%s\x00%s\x00%s\x00%d", p.TenantId, p.Id, p.DeviceBootId, p.SequenceNumber)
	sum := sha256.Sum256([]byte(identity))
	eventID := hex.EncodeToString(sum[:])
	if correlationID == "" {
		correlationID = eventID
	}
	return &TelemetryEnvelope{
		EventID: eventID, EventType: TelemetryEventType, SchemaVersion: p.SchemaVersion,
		TenantID: p.TenantId, DeviceID: p.Id, DeviceBootID: p.DeviceBootId,
		SequenceNumber: p.SequenceNumber, BootStartedAt: p.BootStartedAt,
		ObservedAt: observedAt, IngestedAt: ingestedAt.UnixMilli(),
		CorrelationID: correlationID, CausationID: causationID,
		Producer: GatewayProducer, Traceparent: traceparent, Payload: p,
	}
}

func (e *TelemetryEnvelope) PartitionKey() string     { return e.TenantID + ":" + e.DeviceID }
func (e *TelemetryEnvelope) Marshal() ([]byte, error) { return json.Marshal(e) }
func Unmarshal(data []byte) (*TelemetryEnvelope, error) {
	var e TelemetryEnvelope
	if err := json.Unmarshal(data, &e); err != nil {
		return nil, fmt.Errorf("malformed envelope: %w", err)
	}
	if err := ValidateEnvelope(&e); err != nil {
		return nil, err
	}
	return &e, nil
}

func ValidateFrame(p *pb.SpatialObject, now time.Time) error {
	if p == nil {
		return errors.New("missing payload")
	}
	if !identityPattern.MatchString(p.TenantId) {
		return errors.New("invalid tenant_id")
	}
	if !identityPattern.MatchString(p.Id) {
		return errors.New("invalid device_id")
	}
	if !identityPattern.MatchString(p.DeviceBootId) {
		return errors.New("invalid device_boot_id")
	}
	if p.SequenceNumber == 0 {
		return errors.New("sequence_number must be positive")
	}
	if p.SequenceNumber > math.MaxInt64 {
		return errors.New("sequence_number exceeds supported range")
	}
	if p.SchemaVersion != CurrentSchemaVersion {
		return fmt.Errorf("unsupported schema_version %d", p.SchemaVersion)
	}
	if math.IsNaN(p.Lat) || math.IsInf(p.Lat, 0) || p.Lat < -90 || p.Lat > 90 {
		return errors.New("invalid latitude")
	}
	if math.IsNaN(p.Lon) || math.IsInf(p.Lon, 0) || p.Lon < -180 || p.Lon > 180 {
		return errors.New("invalid longitude")
	}
	if p.EnergyPercent < 0 || p.EnergyPercent > 100 {
		return errors.New("battery outside 0..100")
	}
	if math.IsNaN(p.VelocityMps) || math.IsInf(p.VelocityMps, 0) || p.VelocityMps < 0 || p.VelocityMps > 250 {
		return errors.New("invalid velocity")
	}
	if p.Type <= pb.NodeType_NODE_TYPE_UNKNOWN || p.Type > pb.NodeType_NODE_TYPE_STATIC_SENSOR {
		return errors.New("invalid device type")
	}
	observed := p.ObservedAt
	if observed == 0 {
		observed = p.Timestamp
	}
	if observed <= 0 {
		return errors.New("missing observed_at")
	}
	if p.BootStartedAt <= 0 || p.BootStartedAt > observed {
		return errors.New("invalid boot_started_at")
	}
	observedTime := time.UnixMilli(observed)
	if observedTime.Before(now.Add(-24*time.Hour)) || observedTime.After(now.Add(5*time.Minute)) {
		return errors.New("observation timestamp outside allowed window")
	}
	return nil
}

func ValidateEnvelope(e *TelemetryEnvelope) error {
	if e == nil || e.Payload == nil {
		return errors.New("missing envelope payload")
	}
	if e.EventID == "" || e.EventType != TelemetryEventType || e.Producer == "" || e.IngestedAt <= 0 || e.CorrelationID == "" {
		return errors.New("missing platform envelope metadata")
	}
	if e.SchemaVersion != CurrentSchemaVersion {
		return fmt.Errorf("unsupported schema_version %d", e.SchemaVersion)
	}
	if err := ValidateFrame(e.Payload, time.UnixMilli(e.IngestedAt)); err != nil {
		return err
	}
	if e.TenantID != e.Payload.TenantId || e.DeviceID != e.Payload.Id || e.DeviceBootID != e.Payload.DeviceBootId || e.SequenceNumber != e.Payload.SequenceNumber {
		return errors.New("envelope/payload identity mismatch")
	}
	if e.BootStartedAt != e.Payload.BootStartedAt {
		return errors.New("envelope/payload boot timestamp mismatch")
	}
	payloadObservedAt := e.Payload.ObservedAt
	if payloadObservedAt == 0 {
		payloadObservedAt = e.Payload.Timestamp
	}
	if e.ObservedAt != payloadObservedAt || e.SchemaVersion != e.Payload.SchemaVersion {
		return errors.New("envelope/payload schema or observation mismatch")
	}
	if strings.TrimSpace(e.Traceparent) != "" && !strings.HasPrefix(strings.ToLower(e.Traceparent), "00-") {
		return errors.New("invalid traceparent")
	}
	return nil
}
```

---

## internal\core\events\telemetry_test.go

```
package events

import (
	pb "github.com/Akashpg-M/polaris/backend/api/proto/v1"
	"testing"
	"time"
)

func validFrame(now time.Time) *pb.SpatialObject {
	return &pb.SpatialObject{Id: "device-1", TenantId: "tenant-1", DeviceBootId: "boot-1", SequenceNumber: 1,
		BootStartedAt: now.Add(-time.Minute).UnixMilli(), ObservedAt: now.UnixMilli(), SchemaVersion: 1,
		Type: pb.NodeType_NODE_TYPE_DRONE, Status: pb.NodeStatus_NODE_STATUS_ACTIVE, Lat: 13, Lon: 80, EnergyPercent: 50}
}

func TestEnvelopeRoundTripAndStableIdentity(t *testing.T) {
	now := time.Now().UTC()
	p := validFrame(now)
	a := NewTelemetryEnvelope(p, now, "", "", "")
	b := NewTelemetryEnvelope(p, now.Add(time.Second), "", "", "")
	if a.EventID != b.EventID {
		t.Fatal("device tuple must produce a stable replay identity")
	}
	data, err := a.Marshal()
	if err != nil {
		t.Fatal(err)
	}
	decoded, err := Unmarshal(data)
	if err != nil {
		t.Fatal(err)
	}
	if decoded.PartitionKey() != "tenant-1:device-1" {
		t.Fatalf("unexpected key %q", decoded.PartitionKey())
	}
}

func TestFrameValidationRejectsPermanentFailures(t *testing.T) {
	now := time.Now().UTC()
	tests := []struct {
		name   string
		mutate func(*pb.SpatialObject)
	}{
		{"unsupported schema", func(p *pb.SpatialObject) { p.SchemaVersion = 99 }},
		{"invalid coordinate", func(p *pb.SpatialObject) { p.Lat = 91 }},
		{"missing identity", func(p *pb.SpatialObject) { p.Id = "" }},
		{"invalid battery", func(p *pb.SpatialObject) { p.EnergyPercent = 101 }},
		{"invalid velocity", func(p *pb.SpatialObject) { p.VelocityMps = -1 }},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			p := validFrame(now)
			tc.mutate(p)
			if ValidateFrame(p, now) == nil {
				t.Fatal("expected rejection")
			}
		})
	}
}
```

---

## internal\core\extension\contracts.go

```
package extension

import (
	"context"
	"encoding/json"
	"time"

	taskcore "github.com/Akashpg-M/polaris/backend/internal/core/task"
	twincore "github.com/Akashpg-M/polaris/backend/internal/core/twin"
)

type ModuleState string

const (
	ModuleStarting ModuleState = "STARTING"
	ModuleReady    ModuleState = "READY"
	ModuleDegraded ModuleState = "DEGRADED"
	ModuleFailed   ModuleState = "FAILED"
	ModuleStopped  ModuleState = "STOPPED"
)

type ModuleStatus struct {
	State      ModuleState             `json:"state"`
	Message    string                  `json:"message,omitempty"`
	Components map[string]ModuleStatus `json:"components,omitempty"`
	Details    map[string]any          `json:"details,omitempty"`
}

type Module interface {
	Name() string
	Start(context.Context) error
	Ready(context.Context) ModuleStatus
	Close(context.Context) error
}

type CandidateRequest struct {
	TenantID             string
	EligibleDeviceIDs    []string
	RequiredCapabilities []string
	DeviceTypeIDs        []string
	ProjectIDs           []string
	Limit                int
	Context              map[string]any
	Timing               *CandidateTiming
}

// CandidateTiming is optional request-scoped instrumentation. Providers add
// only time spent in domain routing so Core can distinguish lookup/ranking
// from routing without coupling itself to a concrete capability module.
type CandidateTiming struct{ RoutingDuration time.Duration }

type Candidate struct {
	DeviceID    string         `json:"device_id"`
	DomainScore *float64       `json:"domain_score,omitempty"`
	Attributes  map[string]any `json:"attributes,omitempty"`
}

type CandidateProvider interface {
	Name() string
	Supports(CandidateRequest) bool
	Candidates(context.Context, CandidateRequest) ([]Candidate, error)
}

type PlanningRequest struct {
	Task       taskcore.Task
	DeviceTwin twincore.DeviceTwin
}

type ExecutionPlan struct {
	PlannerName   string          `json:"planner_name"`
	SchemaVersion uint32          `json:"schema_version"`
	CommandType   string          `json:"command_type"`
	Payload       json.RawMessage `json:"payload"`
	GeneratedAt   time.Time       `json:"generated_at"`
	ValidUntil    *time.Time      `json:"valid_until,omitempty"`
	Metadata      map[string]any  `json:"metadata,omitempty"`
}

type TaskPlanner interface {
	Name() string
	Supports(taskcore.Task) bool
	Plan(context.Context, PlanningRequest) (ExecutionPlan, error)
}
```

---

## internal\core\extension\default_planner.go

```
package extension

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	taskcore "github.com/Akashpg-M/polaris/backend/internal/core/task"
)

type DefaultTaskPlanner struct{}

func (DefaultTaskPlanner) Name() string                { return "core.default/v1" }
func (DefaultTaskPlanner) Supports(taskcore.Task) bool { return true }
func (DefaultTaskPlanner) Plan(_ context.Context, req PlanningRequest) (ExecutionPlan, error) {
	var requirements taskcore.Requirements
	if len(req.Task.Requirements) > 0 && json.Unmarshal(req.Task.Requirements, &requirements) != nil {
		return ExecutionPlan{}, errors.New("invalid task requirements")
	}
	if requirements.PlanningMode == taskcore.PlanningPolarisRequired {
		return ExecutionPlan{}, ErrPlanningRequired
	}
	if len(req.Task.Target) == 0 {
		return ExecutionPlan{}, errors.New("task target is empty")
	}
	now := time.Now().UTC()
	validUntil := req.Task.ExpiresAt
	return ExecutionPlan{PlannerName: "core.default/v1", SchemaVersion: 1, CommandType: req.Task.TaskType, Payload: req.Task.Target, GeneratedAt: now, ValidUntil: &validUntil}, nil
}
```

---

## internal\core\extension\registry.go

```
package extension

import (
	"context"
	"errors"
	"fmt"
	"sync"

	taskcore "github.com/Akashpg-M/polaris/backend/internal/core/task"
)

var (
	ErrNoCandidateProvider = errors.New("no candidate provider supports request")
	ErrNoTaskPlanner       = errors.New("no task planner supports task")
	ErrPlanningUnsupported = errors.New("planner does not support selected device")
	ErrPlanningRequired    = errors.New("compatible Polaris planner required")
)

// Registry is populated explicitly by application composition. It deliberately
// has no init hooks or dynamic plugin loading.
type Registry struct {
	mu        sync.RWMutex
	modules   []Module
	providers []CandidateProvider
	planners  []TaskPlanner
}

func NewRegistry() *Registry { return &Registry{} }

func (r *Registry) RegisterModule(m Module) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.modules = append(r.modules, m)
}
func (r *Registry) RegisterCandidateProvider(p CandidateProvider) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.providers = append(r.providers, p)
}
func (r *Registry) RegisterTaskPlanner(p TaskPlanner) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.planners = append(r.planners, p)
}

func (r *Registry) CandidateProvider(req CandidateRequest) (CandidateProvider, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	for _, p := range r.providers {
		if p.Supports(req) {
			return p, nil
		}
	}
	return nil, ErrNoCandidateProvider
}

func (r *Registry) TaskPlanner(v taskcore.Task) (TaskPlanner, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	for _, p := range r.planners {
		if p.Supports(v) {
			return p, nil
		}
	}
	return nil, ErrNoTaskPlanner
}

func (r *Registry) TaskPlanners(v taskcore.Task) []TaskPlanner {
	r.mu.RLock()
	defer r.mu.RUnlock()
	out := []TaskPlanner{}
	for _, p := range r.planners {
		if p.Supports(v) {
			out = append(out, p)
		}
	}
	return out
}

func (r *Registry) Start(ctx context.Context) error {
	r.mu.RLock()
	modules := append([]Module(nil), r.modules...)
	r.mu.RUnlock()
	started := make([]Module, 0, len(modules))
	for _, m := range modules {
		if err := m.Start(ctx); err != nil {
			for i := len(started) - 1; i >= 0; i-- {
				_ = started[i].Close(context.Background())
			}
			return fmt.Errorf("start module %s: %w", m.Name(), err)
		}
		started = append(started, m)
	}
	return nil
}

func (r *Registry) Close(ctx context.Context) error {
	r.mu.RLock()
	modules := append([]Module(nil), r.modules...)
	r.mu.RUnlock()
	var joined error
	for i := len(modules) - 1; i >= 0; i-- {
		joined = errors.Join(joined, modules[i].Close(ctx))
	}
	return joined
}

func (r *Registry) Status(ctx context.Context) map[string]ModuleStatus {
	r.mu.RLock()
	modules := append([]Module(nil), r.modules...)
	r.mu.RUnlock()
	out := make(map[string]ModuleStatus, len(modules))
	for _, m := range modules {
		out[m.Name()] = m.Ready(ctx)
	}
	return out
}
```

---

## internal\core\extension\registry_test.go

```
package extension

import (
	"context"
	"encoding/json"
	"errors"
	"testing"
	"time"

	taskcore "github.com/Akashpg-M/polaris/backend/internal/core/task"
)

type fixtureModule struct{ started, closed bool }

func (m *fixtureModule) Name() string                { return "fixture" }
func (m *fixtureModule) Start(context.Context) error { m.started = true; return nil }
func (m *fixtureModule) Ready(context.Context) ModuleStatus {
	if m.started && !m.closed {
		return ModuleStatus{State: ModuleReady}
	}
	return ModuleStatus{State: ModuleStopped}
}

func TestDefaultPlannerRejectsPolarisRequiredTask(t *testing.T) {
	requirements, _ := json.Marshal(taskcore.Requirements{PlanningMode: taskcore.PlanningPolarisRequired})
	_, err := (DefaultTaskPlanner{}).Plan(context.Background(), PlanningRequest{Task: taskcore.Task{TaskType: "NAVIGATE", Requirements: requirements, Target: []byte(`{"lat":13,"lon":80}`), ExpiresAt: time.Now().Add(time.Minute)}})
	if !errors.Is(err, ErrPlanningRequired) {
		t.Fatalf("expected planner requirement, got %v", err)
	}
}
func (m *fixtureModule) Close(context.Context) error { m.closed = true; return nil }
func TestExplicitModuleLifecycle(t *testing.T) {
	r := NewRegistry()
	m := &fixtureModule{}
	r.RegisterModule(m)
	if err := r.Start(context.Background()); err != nil || !m.started {
		t.Fatalf("module did not start: %v", err)
	}
	if r.Status(context.Background())["fixture"].State != ModuleReady {
		t.Fatal("readiness not exposed")
	}
	_ = r.Close(context.Background())
	if !m.closed {
		t.Fatal("module not closed")
	}
}
func TestMobilityDisabledLeavesGenericPlanning(t *testing.T) {
	r := NewRegistry()
	r.RegisterTaskPlanner(DefaultTaskPlanner{})
	for _, kind := range []string{"CAPTURE_IMAGE", "RUN_MODEL"} {
		task := taskcore.Task{TaskType: kind, Target: []byte(`{"fixture":true}`), ExpiresAt: time.Now().Add(time.Minute)}
		planner, err := r.TaskPlanner(task)
		if err != nil {
			t.Fatalf("%s unavailable with Mobility absent: %v", kind, err)
		}
		plan, err := planner.Plan(context.Background(), PlanningRequest{Task: task})
		if err != nil || plan.CommandType != kind {
			t.Fatalf("generic planning failed: %#v %v", plan, err)
		}
	}
	if _, err := r.TaskPlanner(taskcore.Task{TaskType: "NAVIGATE"}); err != nil {
		t.Fatalf("generic high-level NAVIGATE must remain available when Mobility is disabled: %v", err)
	}
}
```

---

## internal\core\ports\stream.go

```
package ports

import (
	"context"
	pb "github.com/Akashpg-M/polaris/backend/api/proto/v1"
)

type TelemetryPublisher interface {
	Publish(ctx context.Context, payload *pb.SpatialObject) error
}
```

---

## internal\core\registry\model.go

```
package registry

import (
	"encoding/json"
	"time"
)

type Tenant struct {
	TenantID    string          `db:"tenant_id" json:"tenant_id"`
	DisplayName string          `db:"display_name" json:"display_name"`
	Status      string          `db:"status" json:"status"`
	Metadata    json.RawMessage `db:"metadata" json:"metadata"`
	CreatedAt   time.Time       `db:"created_at" json:"created_at"`
	UpdatedAt   time.Time       `db:"updated_at" json:"updated_at"`
}
type Project struct {
	ProjectID   string          `db:"project_id" json:"project_id"`
	TenantID    string          `db:"tenant_id" json:"tenant_id"`
	Name        string          `db:"name" json:"name"`
	Description *string         `db:"description" json:"description,omitempty"`
	Status      string          `db:"status" json:"status"`
	Metadata    json.RawMessage `db:"metadata" json:"metadata"`
	CreatedAt   time.Time       `db:"created_at" json:"created_at"`
	UpdatedAt   time.Time       `db:"updated_at" json:"updated_at"`
}
type Device struct {
	TenantID        string          `db:"tenant_id" json:"tenant_id"`
	DeviceID        string          `db:"device_id" json:"device_id"`
	ProjectID       *string         `db:"project_id" json:"project_id,omitempty"`
	DeviceTypeID    string          `db:"device_type_id" json:"device_type_id"`
	DisplayName     string          `db:"display_name" json:"display_name"`
	LifecycleStatus string          `db:"lifecycle_status" json:"lifecycle_status"`
	FirmwareVersion *string         `db:"firmware_version" json:"firmware_version,omitempty"`
	SoftwareVersion *string         `db:"software_version" json:"software_version,omitempty"`
	ModelVersion    *string         `db:"model_version" json:"model_version,omitempty"`
	Metadata        json.RawMessage `db:"metadata" json:"metadata"`
	RegisteredAt    time.Time       `db:"registered_at" json:"registered_at"`
	UpdatedAt       time.Time       `db:"updated_at" json:"updated_at"`
	DeactivatedAt   *time.Time      `db:"deactivated_at" json:"deactivated_at,omitempty"`
}
type Capability struct {
	CapabilityID  string          `db:"capability_id" json:"capability_id"`
	DisplayName   string          `db:"display_name" json:"display_name"`
	Description   *string         `db:"description" json:"description,omitempty"`
	Configuration json.RawMessage `db:"configuration" json:"configuration"`
	Enabled       bool            `db:"enabled" json:"enabled"`
}
type CredentialMetadata struct {
	CredentialID string     `db:"credential_id" json:"credential_id"`
	TokenPrefix  string     `db:"token_prefix" json:"token_prefix"`
	Status       string     `db:"status" json:"status"`
	IssuedAt     time.Time  `db:"issued_at" json:"issued_at"`
	ExpiresAt    *time.Time `db:"expires_at" json:"expires_at,omitempty"`
	LastUsedAt   *time.Time `db:"last_used_at" json:"last_used_at,omitempty"`
	RevokedAt    *time.Time `db:"revoked_at" json:"revoked_at,omitempty"`
}

func ValidTransition(from, to string) bool {
	switch from + ">" + to {
	case "REGISTERED>ACTIVE", "REGISTERED>DECOMMISSIONED", "ACTIVE>SUSPENDED", "ACTIVE>DECOMMISSIONED", "SUSPENDED>ACTIVE", "SUSPENDED>DECOMMISSIONED":
		return true
	}
	return false
}
```

---

## internal\core\registry\model_test.go

```
package registry

import "testing"

func TestLifecycleTransitions(t *testing.T) {
	allowed := [][2]string{{"REGISTERED", "ACTIVE"}, {"ACTIVE", "SUSPENDED"}, {"SUSPENDED", "ACTIVE"}, {"ACTIVE", "DECOMMISSIONED"}}
	for _, v := range allowed {
		if !ValidTransition(v[0], v[1]) {
			t.Fatalf("expected %s -> %s", v[0], v[1])
		}
	}
	denied := [][2]string{{"DECOMMISSIONED", "ACTIVE"}, {"ACTIVE", "REGISTERED"}, {"REGISTERED", "SUSPENDED"}}
	for _, v := range denied {
		if ValidTransition(v[0], v[1]) {
			t.Fatalf("unexpected %s -> %s", v[0], v[1])
		}
	}
}
```

---

## internal\core\simulation\ca_runline.go

```
package simulation

import (
	"context"
	"math/rand"
	"time"
)

type CellularAutomataRunline struct {
	maxRiskThreshold float64
	computeDeadline  time.Duration
}

func NewCellularAutomataRunline(maxRisk float64, deadline time.Duration) *CellularAutomataRunline {
	return &CellularAutomataRunline{
		maxRiskThreshold: maxRisk,
		computeDeadline:  deadline,
	}
}

// ValidateAsync spins up an isolated forward-projection slice without blocking the routing path
func (ca *CellularAutomataRunline) ValidateAsync(ctx context.Context, req ValidationRequest) <-chan SimulationResult {
	out := make(chan SimulationResult, 1)

	go func() {
		defer close(out)
		startTime := time.Now()

		// Local timeout safety circuit bound to our real-time deadline parameter
		simCtx, cancel := context.WithTimeout(ctx, ca.computeDeadline)
		defer cancel()

		// Channel to catch the localized worker calculation results
		resultChan := make(chan float64, 1)

		go func() {
			// This simulates a discrete cellular automaton state matrix transition
			// In production, you map the target edge into a bitmapped cell sequence
			// representing vehicle slots (Nagel-Schreckenberg model matrix).
			
			// We execute a brief micro-step iteration simulating 20x real-world speed
			time.Sleep(3 * time.Millisecond) 
			
			// Generate a simulated density/risk value
			// In production, this calculates actual emergent bottlenecks or deadlocks
			r := rand.New(rand.NewSource(time.Now().UnixNano()))
			resultChan <- r.Float64()
		}()

		select {
		case <-simCtx.Done():
			// The simulation runline breached its real-time compute threshold constraint!
			// Fail-safe default strategy: Reject the risk proposal to preserve stability.
			out <- SimulationResult{
				RouteID:    req.RouteID,
				AllowRoute: false,
				RiskScore:  1.0,
				ComputeTime: time.Since(startTime),
			}
		case calculatedRisk := <-resultChan:
			isAllowed := calculatedRisk <= ca.maxRiskThreshold

			out <- SimulationResult{
				RouteID:    req.RouteID,
				AllowRoute: isAllowed,
				RiskScore:  calculatedRisk,
				ComputeTime: time.Since(startTime),
			}
		}
	}()

	return out
}
```

---

## internal\core\simulation\protocols.go

```
package simulation

import (
	"time"
)

// ValidationRequest specifies the isolated spatial tracks to simulate
type ValidationRequest struct {
	RouteID    string    `json:"route_id"`
	TargetEdge string    `json:"target_edge"`
	H3Sectors  []string  `json:"h3_sectors"` // Limits snapshot boundaries
	Timestamp  time.Time `json:"timestamp"`
}

// SimulationResult returns the deterministic risk evaluation
type SimulationResult struct {
	RouteID     string  `json:"route_id"`
	AllowRoute  bool    `json:"allow_route"`  // Strict binary gate
	RiskScore   float64 `json:"risk_score"`   // 0.0 (Safe) to 1.0 (Gridlock)
	ComputeTime time.Duration `json:"compute_time"`
}
```

---

## internal\core\task\model.go

```
package task

import (
	"encoding/json"
	"time"
)

type Status string
type PlanningMode string

const (
	Pending    Status = "PENDING"
	Assigning  Status = "ASSIGNING"
	Assigned   Status = "ASSIGNED"
	InProgress Status = "IN_PROGRESS"
	Completed  Status = "COMPLETED"
	Failed     Status = "FAILED"
	Cancelled  Status = "CANCELLED"
	Expired    Status = "EXPIRED"
)

const (
	PlanningDeviceLocal     PlanningMode = "DEVICE_LOCAL"
	PlanningPolarisRequired PlanningMode = "POLARIS_REQUIRED"
)

type Requirements struct {
	RequiredCapabilities []string     `json:"required_capabilities,omitempty"`
	MinimumBattery       int32        `json:"minimum_battery,omitempty"`
	AllowedDeviceTypes   []string     `json:"allowed_device_types,omitempty"`
	MaximumDistanceM     float64      `json:"max_distance_meters,omitempty"`
	ProjectID            string       `json:"project_id,omitempty"`
	PlanningMode         PlanningMode `json:"planning_mode,omitempty"`
	Custom               any          `json:"custom_constraints,omitempty"`
}

type Target struct {
	Latitude  *float64        `json:"lat,omitempty"`
	Longitude *float64        `json:"lon,omitempty"`
	H3Cell    string          `json:"h3_cell,omitempty"`
	Data      json.RawMessage `json:"data,omitempty"`
}

type Task struct {
	TaskID           string          `db:"task_id" json:"task_id"`
	TenantID         string          `db:"tenant_id" json:"tenant_id"`
	ProjectID        *string         `db:"project_id" json:"project_id,omitempty"`
	TaskType         string          `db:"task_type" json:"task_type"`
	Status           string          `db:"status" json:"status"`
	Priority         string          `db:"priority" json:"priority"`
	Requirements     json.RawMessage `db:"requirements" json:"requirements"`
	Target           json.RawMessage `db:"target" json:"target"`
	AssignedDeviceID *string         `db:"assigned_device_id" json:"assigned_device_id,omitempty"`
	CorrelationID    string          `db:"correlation_id" json:"correlation_id"`
	CreatedBy        string          `db:"created_by" json:"created_by"`
	Version          int64           `db:"version" json:"version"`
	CreatedAt        time.Time       `db:"created_at" json:"created_at"`
	UpdatedAt        time.Time       `db:"updated_at" json:"updated_at"`
	AssignedAt       *time.Time      `db:"assigned_at" json:"assigned_at,omitempty"`
	StartedAt        *time.Time      `db:"started_at" json:"started_at,omitempty"`
	CompletedAt      *time.Time      `db:"completed_at" json:"completed_at,omitempty"`
	FailedAt         *time.Time      `db:"failed_at" json:"failed_at,omitempty"`
	ExpiresAt        time.Time       `db:"expires_at" json:"expires_at"`
	FailureReason    *string         `db:"failure_reason" json:"failure_reason,omitempty"`
}

func ValidTransition(from, to Status) bool {
	if from == to {
		return true
	}
	switch from {
	case Pending:
		return to == Assigning || to == Cancelled || to == Expired || to == Failed
	case Assigning:
		return to == Assigned || to == Pending || to == Cancelled || to == Expired || to == Failed
	case Assigned:
		return to == InProgress || to == Cancelled || to == Expired || to == Failed
	case InProgress:
		return to == Completed || to == Failed || to == Expired
	}
	return false
}

func IsTerminal(status Status) bool {
	return status == Completed || status == Failed || status == Cancelled || status == Expired
}
```

---

## internal\core\task\model_test.go

```
package task

import "testing"

func TestTaskTransitions(t *testing.T) {
	allowed := [][2]Status{{Pending, Assigning}, {Assigning, Assigned}, {Assigned, InProgress}, {InProgress, Completed}, {Assigned, Failed}}
	for _, pair := range allowed {
		if !ValidTransition(pair[0], pair[1]) {
			t.Fatalf("expected %s -> %s", pair[0], pair[1])
		}
	}
	if ValidTransition(Completed, Pending) || ValidTransition(Cancelled, InProgress) {
		t.Fatal("terminal task state was reversible")
	}
}
```

---

## internal\core\twin\component.go

```
package twin

import (
	"encoding/json"
	"time"
)

// ComponentEnvelope is the storage contract for independently versioned twin
// components. Core code treats Payload as opaque; capability modules own it.
type ComponentEnvelope struct {
	Type           string          `json:"type"`
	SchemaVersion  uint32          `json:"schema_version"`
	ObservedAt     time.Time       `json:"observed_at"`
	BootID         string          `json:"boot_id"`
	SequenceNumber uint64          `json:"sequence_number"`
	Payload        json.RawMessage `json:"payload"`
}

type DeviceTwin struct {
	TenantID     string                       `json:"tenant_id"`
	DeviceID     string                       `json:"device_id"`
	Connectivity string                       `json:"connectivity"`
	Components   map[string]ComponentEnvelope `json:"components"`
}
```

---

## internal\infra\postgres\db.go

```
// package postgresinfra

// import (
// 	"log"

// 	"github.com/jmoiron/sqlx"
// 	_ "github.com/lib/pq"
// )

// func NewDB(url string) *sqlx.DB {
// 	db, err := sqlx.Connect("postgres", url)
// 	if err != nil {
// 		log.Fatalf("Postgres connection failed: %v", err)
// 	}

// 	return db
// }

package postgresinfra

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func NewDB(url string) (*sqlx.DB, error) {
	db, err := sqlx.Connect("postgres", url)
	if err != nil {
		return nil, fmt.Errorf("postgres connection failed: %w", err)
	}
	return db, nil
}
```

---

## internal\infra\redis\client.go

```
package redisinfra

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
)

func NewClient(url string) (*redis.Client, error) {
	opts, err := redis.ParseURL(url)
	if err != nil {
		return nil, fmt.Errorf("invalid redis url: %w", err)
	}

	client := redis.NewClient(opts)

	if err := client.Ping(context.Background()).Err(); err != nil {
		return nil, fmt.Errorf("redis connection failed: %w", err)
	}

	return client, nil
}
```

---

## internal\modules\mobility\config.go

```
package mobility

import (
	"errors"
	"os"
	"strconv"
	"time"
)

type Config struct {
	Enabled                   bool
	Required                  bool
	SpatialEnabled            bool
	RoutingEnabled            bool
	H3Resolution              int
	H3ShardResolution         int
	IndexMinMoveMeters        float64
	IndexMaxAge               time.Duration
	MaxH3Rings                int
	MaxSearchRadiusMeters     float64
	MaxRawCandidates          int
	MaxRoutedCandidates       int
	MaxActiveDevicesPerTenant int
	RoutingWorkers            int
	RoutingQueueCapacity      int
	RoutingTimeout            time.Duration
	MaxRouteExpansions        int
	MaxConcurrentRoutesTenant int
	MaxTrafficObservationAge  time.Duration
	TrafficRefreshInterval    time.Duration
	TrafficScope              string
	RoadGraphPath             string
	RoadGraphVersion          string
}

func DefaultConfig() Config {
	return Config{Enabled: true, SpatialEnabled: true, RoutingEnabled: true, H3Resolution: 8, H3ShardResolution: 6,
		IndexMinMoveMeters: 5, IndexMaxAge: 30 * time.Second, MaxH3Rings: 12, MaxSearchRadiusMeters: 10_000,
		MaxRawCandidates: 50, MaxRoutedCandidates: 8, MaxActiveDevicesPerTenant: 10_000,
		RoutingWorkers: 4, RoutingQueueCapacity: 64, RoutingTimeout: 2 * time.Second, MaxRouteExpansions: 250_000,
		MaxConcurrentRoutesTenant: 2, MaxTrafficObservationAge: 10 * time.Minute, TrafficRefreshInterval: 15 * time.Second, TrafficScope: "SHARED_TRUSTED",
		RoadGraphPath: "data/chennai-metro.osm.pbf", RoadGraphVersion: "chennai-v1"}
}

func LoadConfig() (Config, error) {
	c := DefaultConfig()
	c.Enabled = envBool("POLARIS_MODULE_MOBILITY_ENABLED", c.Enabled)
	c.Required = envBool("POLARIS_MODULE_MOBILITY_REQUIRED", c.Required)
	c.SpatialEnabled = envBool("MOBILITY_SPATIAL_ENABLED", c.SpatialEnabled)
	c.RoutingEnabled = envBool("MOBILITY_ROUTING_ENABLED", c.RoutingEnabled)
	c.H3Resolution = envInt("MOBILITY_H3_RESOLUTION", c.H3Resolution)
	c.H3ShardResolution = envInt("MOBILITY_H3_SHARD_RESOLUTION", c.H3ShardResolution)
	c.IndexMinMoveMeters = envFloat("MOBILITY_INDEX_MIN_MOVE_METERS", c.IndexMinMoveMeters)
	c.IndexMaxAge = envDuration("MOBILITY_INDEX_MAX_AGE", c.IndexMaxAge)
	c.MaxH3Rings = envInt("MOBILITY_MAX_H3_RINGS", c.MaxH3Rings)
	c.MaxSearchRadiusMeters = envFloat("MOBILITY_MAX_SEARCH_RADIUS_METERS", c.MaxSearchRadiusMeters)
	c.MaxRawCandidates = envInt("MOBILITY_MAX_RAW_CANDIDATES", c.MaxRawCandidates)
	c.MaxRoutedCandidates = envInt("MOBILITY_MAX_ROUTED_CANDIDATES", c.MaxRoutedCandidates)
	c.MaxActiveDevicesPerTenant = envInt("MOBILITY_MAX_ACTIVE_DEVICES_PER_TENANT", c.MaxActiveDevicesPerTenant)
	c.RoutingWorkers = envInt("MOBILITY_ROUTING_WORKERS", c.RoutingWorkers)
	c.RoutingQueueCapacity = envInt("MOBILITY_ROUTING_QUEUE_CAPACITY", c.RoutingQueueCapacity)
	c.RoutingTimeout = envDuration("MOBILITY_ROUTING_TIMEOUT", c.RoutingTimeout)
	c.MaxRouteExpansions = envInt("MOBILITY_MAX_ROUTE_EXPANSIONS", c.MaxRouteExpansions)
	c.MaxConcurrentRoutesTenant = envInt("MOBILITY_MAX_CONCURRENT_ROUTES_PER_TENANT", c.MaxConcurrentRoutesTenant)
	c.MaxTrafficObservationAge = envDuration("MOBILITY_MAX_TRAFFIC_OBSERVATION_AGE", c.MaxTrafficObservationAge)
	c.TrafficRefreshInterval = envDuration("MOBILITY_TRAFFIC_REFRESH_INTERVAL", c.TrafficRefreshInterval)
	if v := os.Getenv("MOBILITY_TRAFFIC_SCOPE"); v != "" {
		c.TrafficScope = v
	}
	if v := os.Getenv("MOBILITY_ROAD_GRAPH_PATH"); v != "" {
		c.RoadGraphPath = v
	}
	if v := os.Getenv("MOBILITY_ROAD_GRAPH_VERSION"); v != "" {
		c.RoadGraphVersion = v
	}
	return c, c.Validate()
}

func (c Config) Validate() error {
	if c.H3Resolution < 0 || c.H3Resolution > 15 || c.H3ShardResolution < 0 || c.H3ShardResolution > c.H3Resolution {
		return errors.New("invalid H3 resolutions")
	}
	if c.MaxH3Rings < 0 || c.MaxSearchRadiusMeters <= 0 || c.MaxRawCandidates < 1 || c.MaxRoutedCandidates < 1 || c.MaxRoutedCandidates > c.MaxRawCandidates {
		return errors.New("invalid candidate limits")
	}
	if c.RoutingWorkers < 1 || c.RoutingQueueCapacity < 1 || c.RoutingTimeout <= 0 || c.MaxRouteExpansions < 1 || c.MaxConcurrentRoutesTenant < 1 {
		return errors.New("invalid routing limits")
	}
	if c.MaxActiveDevicesPerTenant < 1 || c.IndexMinMoveMeters < 0 || c.IndexMaxAge <= 0 || c.MaxTrafficObservationAge <= 0 || c.TrafficRefreshInterval <= 0 {
		return errors.New("invalid spatial limits")
	}
	if c.TrafficScope != "SHARED_TRUSTED" {
		return errors.New("unsupported traffic scope; only SHARED_TRUSTED is currently implemented")
	}
	return nil
}

func envBool(k string, d bool) bool {
	if v, e := strconv.ParseBool(os.Getenv(k)); e == nil {
		return v
	}
	return d
}
func envInt(k string, d int) int {
	if v, e := strconv.Atoi(os.Getenv(k)); e == nil {
		return v
	}
	return d
}
func envFloat(k string, d float64) float64 {
	if v, e := strconv.ParseFloat(os.Getenv(k), 64); e == nil {
		return v
	}
	return d
}
func envDuration(k string, d time.Duration) time.Duration {
	if v, e := time.ParseDuration(os.Getenv(k)); e == nil {
		return v
	}
	return d
}
```

---

## internal\modules\mobility\config_test.go

```
package mobility

import "testing"

func TestTrafficScopeIsExplicitAndValidated(t *testing.T) {
	cfg := DefaultConfig()
	if cfg.TrafficScope != "SHARED_TRUSTED" || cfg.Validate() != nil {
		t.Fatalf("unexpected default traffic policy: %#v", cfg)
	}
	cfg.TrafficScope = "TENANT_PRIVATE"
	if cfg.Validate() == nil {
		t.Fatal("unimplemented tenant-private traffic policy was silently accepted")
	}
}
```

---

## internal\modules\mobility\module.go

```
package mobility

import (
	"context"
	"sync"

	"github.com/Akashpg-M/polaris/backend/internal/core/extension"
	"github.com/Akashpg-M/polaris/backend/internal/modules/mobility/model"
	"github.com/Akashpg-M/polaris/backend/internal/modules/mobility/routing"
	"github.com/Akashpg-M/polaris/backend/internal/modules/mobility/spatial"
)

type RebuildLoader func(context.Context) ([]model.SpatialState, error)
type Module struct {
	cfg     Config
	Spatial *spatial.Manager
	loader  RebuildLoader
	mu      sync.RWMutex
	state   extension.ModuleState
	message string
	router  *routing.Engine
	graph   *routing.RoadGraph
	traffic *routing.TrafficManager
}

func New(cfg Config, loader RebuildLoader) *Module {
	manager := spatial.NewManager(spatial.ManagerConfig{H3Resolution: cfg.H3Resolution, ShardResolution: cfg.H3ShardResolution, MinMoveMeters: cfg.IndexMinMoveMeters, MaxIndexAge: cfg.IndexMaxAge, MaxH3Rings: cfg.MaxH3Rings, MaxRadiusMeters: cfg.MaxSearchRadiusMeters, MaxDevicesPerTenant: cfg.MaxActiveDevicesPerTenant})
	return &Module{cfg: cfg, Spatial: manager, loader: loader, state: extension.ModuleStarting}
}
func (m *Module) Name() string { return "mobility" }
func (m *Module) Start(ctx context.Context) error {
	m.mu.Lock()
	m.state = extension.ModuleStarting
	m.mu.Unlock()
	if m.cfg.SpatialEnabled && m.loader != nil {
		states, err := m.loader(ctx)
		if err != nil {
			m.set(extension.ModuleDegraded, "spatial rebuild failed: "+err.Error())
		} else if err = m.Spatial.Rebuild(states); err != nil {
			m.set(extension.ModuleDegraded, "spatial rebuild failed: "+err.Error())
		}
	}
	if m.cfg.RoutingEnabled {
		graph, err := routing.LoadOSMPBF(ctx, m.cfg.RoadGraphPath, m.cfg.RoadGraphVersion)
		if err != nil {
			m.set(extension.ModuleDegraded, "road graph unavailable: "+err.Error())
			if m.cfg.Required {
				return err
			}
		} else {
			snapshots := routing.NewSnapshotStore(graph)
			engine := routing.NewEngine(graph, snapshots, routing.EngineConfig{Workers: m.cfg.RoutingWorkers, QueueCapacity: m.cfg.RoutingQueueCapacity, MaxExpansions: m.cfg.MaxRouteExpansions, MaxConcurrentPerTenant: m.cfg.MaxConcurrentRoutesTenant, Timeout: m.cfg.RoutingTimeout})
			engine.Start()
			m.mu.Lock()
			m.graph, m.router = graph, engine
			m.traffic = routing.NewTrafficManager(graph, snapshots, m.cfg.MaxTrafficObservationAge)
			m.mu.Unlock()
		}
	}
	m.mu.Lock()
	if m.state == extension.ModuleStarting {
		m.state = extension.ModuleReady
		m.message = "spatial and configured routing components started"
	}
	m.mu.Unlock()
	return nil
}
func (m *Module) set(s extension.ModuleState, msg string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.state, m.message = s, msg
}
func (m *Module) Ready(context.Context) extension.ModuleStatus {
	m.mu.RLock()
	defer m.mu.RUnlock()
	spatialState := extension.ModuleReady
	if !m.cfg.SpatialEnabled {
		spatialState = extension.ModuleStopped
	}
	routingState := extension.ModuleReady
	if !m.cfg.RoutingEnabled {
		routingState = extension.ModuleStopped
	} else if m.router == nil {
		routingState = extension.ModuleFailed
	}
	details := map[string]any{}
	if m.graph != nil {
		details["road_graph_version"] = m.graph.Version()
		details["road_nodes"] = m.graph.NodeCount()
		details["road_edges"] = m.graph.EdgeCount()
		details["routing_snapshot_version"] = m.router.SnapshotStore().Load().Version
		details["routing_runtime"] = m.router.Stats()
		details["traffic_scope"] = m.cfg.TrafficScope
		details["traffic_refresh_interval"] = m.cfg.TrafficRefreshInterval.String()
		if m.traffic != nil {
			details["traffic_edge_states"] = m.traffic.StateCount()
			details["traffic_overlay_bytes"] = m.traffic.OverlayBytes()
		}
	}
	return extension.ModuleStatus{State: m.state, Message: m.message, Components: map[string]extension.ModuleStatus{"spatial": {State: spatialState}, "routing": {State: routingState}}, Details: details}
}
func (m *Module) Route(ctx context.Context, req routing.RouteRequest) (routing.RouteResult, error) {
	m.mu.RLock()
	r := m.router
	m.mu.RUnlock()
	if r == nil {
		return routing.RouteResult{}, routing.ErrUnavailable
	}
	return r.Route(ctx, req)
}
func (m *Module) Close(context.Context) error {
	m.mu.Lock()
	if m.router != nil {
		m.router.Close()
		m.router = nil
	}
	m.state = extension.ModuleStopped
	m.message = "stopped"
	m.mu.Unlock()
	return nil
}
func (m *Module) Traffic() *routing.TrafficManager {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.traffic
}
```

---

## internal\modules\mobility\module_test.go

```
package mobility

import (
	"context"
	"testing"

	"github.com/Akashpg-M/polaris/backend/internal/core/extension"
)

func TestMissingRoadGraphDegradesOptionalModule(t *testing.T) {
	cfg := DefaultConfig()
	cfg.RoadGraphPath = "does-not-exist.osm.pbf"
	cfg.Required = false
	m := New(cfg, nil)
	if err := m.Start(context.Background()); err != nil {
		t.Fatal(err)
	}
	status := m.Ready(context.Background())
	if status.State != extension.ModuleDegraded || status.Components["spatial"].State != extension.ModuleReady || status.Components["routing"].State != extension.ModuleFailed {
		t.Fatalf("unexpected degradation: %#v", status)
	}
	_ = m.Close(context.Background())
}
func TestMissingRoadGraphFailsRequiredModule(t *testing.T) {
	cfg := DefaultConfig()
	cfg.RoadGraphPath = "does-not-exist.osm.pbf"
	cfg.Required = true
	m := New(cfg, nil)
	if err := m.Start(context.Background()); err == nil {
		t.Fatal("mandatory Mobility accepted a missing graph")
	}
}
func TestConfigurationRejectsUnsafeBounds(t *testing.T) {
	cfg := DefaultConfig()
	cfg.MaxRoutedCandidates = cfg.MaxRawCandidates + 1
	if err := cfg.Validate(); err == nil {
		t.Fatal("invalid candidate fanout accepted")
	}
	cfg = DefaultConfig()
	cfg.H3ShardResolution = cfg.H3Resolution + 1
	if err := cfg.Validate(); err == nil {
		t.Fatal("invalid H3 hierarchy accepted")
	}
}
```

---

## internal\modules\mobility\projector.go

```
package mobility

import (
	"time"

	pb "github.com/Akashpg-M/polaris/backend/api/proto/v1"
	legacy "github.com/Akashpg-M/polaris/backend/internal/application/spatial"
	"github.com/Akashpg-M/polaris/backend/internal/core/events"
	"github.com/Akashpg-M/polaris/backend/internal/modules/mobility/model"
	mobilityspatial "github.com/Akashpg-M/polaris/backend/internal/modules/mobility/spatial"
)

type TelemetryProjector struct{ Manager *mobilityspatial.Manager }

func profile(t pb.NodeType) model.MobilityProfile {
	switch t {
	case pb.NodeType_NODE_TYPE_DRONE:
		return model.MobilityAerialDrone
	case pb.NodeType_NODE_TYPE_ROBOT:
		return model.MobilityGroundRobot
	case pb.NodeType_NODE_TYPE_STATIC_SENSOR:
		return model.MobilityStatic
	default:
		return model.MobilityRoadVehicle
	}
}
func (p *TelemetryProjector) ApplyEnvelope(e *events.TelemetryEnvelope) legacy.Classification {
	speed := e.Payload.VelocityMps
	heading := e.Payload.HeadingDeg
	s := model.SpatialState{TenantID: e.TenantID, DeviceID: e.DeviceID, ReportedPosition: model.Position{Latitude: e.Payload.Lat, Longitude: e.Payload.Lon}, HeadingDegrees: &heading, SpeedMPS: &speed, MobilityProfile: profile(e.Payload.Type), ObservedAt: time.UnixMilli(e.ObservedAt).UTC(), BootID: e.DeviceBootID, BootStartedAt: time.UnixMilli(e.BootStartedAt).UTC(), SequenceNumber: e.SequenceNumber}
	err := p.Manager.Upsert(s)
	if err == nil {
		return legacy.Accepted
	}
	if err == mobilityspatial.ErrStaleVersion {
		return legacy.OutOfOrder
	}
	return legacy.OutOfOrder
}
```

---

## internal\modules\mobility\traffic_consumer.go

```
package mobility

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	pb "github.com/Akashpg-M/polaris/backend/api/proto/v1"
	"github.com/Akashpg-M/polaris/backend/internal/core/events"
	"github.com/Akashpg-M/polaris/backend/internal/modules/mobility/model"
	"github.com/Akashpg-M/polaris/backend/internal/modules/mobility/routing"
	"github.com/segmentio/kafka-go"
)

type TrafficConsumer struct {
	reader          *kafka.Reader
	dlq             *kafka.Writer
	traffic         *routing.TrafficManager
	refreshInterval time.Duration
	done            chan struct{}
}

func NewTrafficConsumer(broker string, traffic *routing.TrafficManager, refreshInterval time.Duration) *TrafficConsumer {
	return &TrafficConsumer{reader: kafka.NewReader(kafka.ReaderConfig{Brokers: []string{broker}, Topic: "telemetry.ingress", GroupID: "polaris_traffic_group", CommitInterval: 0}), dlq: &kafka.Writer{Addr: kafka.TCP(broker), Topic: "telemetry.dead-letter.v1", Balancer: &kafka.Hash{}}, traffic: traffic, refreshInterval: refreshInterval, done: make(chan struct{})}
}
func (c *TrafficConsumer) Start(ctx context.Context) {
	defer close(c.done)
	defer c.reader.Close()
	defer c.dlq.Close()
	refreshCtx, stopRefresh := context.WithCancel(ctx)
	defer stopRefresh()
	go func() {
		ticker := time.NewTicker(c.refreshInterval)
		defer ticker.Stop()
		for {
			select {
			case now := <-ticker.C:
				_ = c.traffic.Refresh(now.UTC())
			case <-refreshCtx.Done():
				return
			}
		}
	}()
	slog.Info("Mobility map-matched traffic consumer started", "refresh_interval", c.refreshInterval, "traffic_scope", "SHARED_TRUSTED")
	for {
		msg, err := c.reader.FetchMessage(ctx)
		if err != nil {
			if ctx.Err() != nil {
				return
			}
			continue
		}
		envelope, parseErr := events.Unmarshal(msg.Value)
		if parseErr != nil {
			for {
				err = c.dlq.WriteMessages(ctx, kafka.Message{Key: msg.Key, Value: msg.Value, Headers: []kafka.Header{{Key: "error_reason", Value: []byte(parseErr.Error())}, {Key: "consumer", Value: []byte("polaris_traffic_group")}, {Key: "source_partition", Value: []byte(fmt.Sprint(msg.Partition))}, {Key: "source_offset", Value: []byte(fmt.Sprint(msg.Offset))}}})
				if err == nil {
					break
				}
				select {
				case <-time.After(250 * time.Millisecond):
				case <-ctx.Done():
					return
				}
			}
		} else if roadTelemetry(envelope.Payload.Type) {
			heading := envelope.Payload.HeadingDeg
			_ = c.traffic.Observe(ctx, routing.TrafficObservation{Position: model.Position{Latitude: envelope.Payload.Lat, Longitude: envelope.Payload.Lon}, HeadingDegrees: &heading, SpeedMPS: envelope.Payload.VelocityMps, ObservedAt: time.UnixMilli(envelope.ObservedAt).UTC()})
		}
		for {
			if err = c.reader.CommitMessages(ctx, msg); err == nil {
				break
			}
			select {
			case <-time.After(250 * time.Millisecond):
			case <-ctx.Done():
				return
			}
		}
	}
}
func roadTelemetry(t pb.NodeType) bool {
	return t == pb.NodeType_NODE_TYPE_BIKE || t == pb.NodeType_NODE_TYPE_AUTO || t == pb.NodeType_NODE_TYPE_SEDAN || t == pb.NodeType_NODE_TYPE_SUV
}
func (c *TrafficConsumer) Wait(ctx context.Context) error {
	select {
	case <-c.done:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}
```

---

## internal\modules\mobility\matching\candidate_provider.go

```
package matching

import (
	"context"
	"errors"
	"sort"
	"time"

	"github.com/Akashpg-M/polaris/backend/internal/core/extension"
	"github.com/Akashpg-M/polaris/backend/internal/modules/mobility/model"
	"github.com/Akashpg-M/polaris/backend/internal/modules/mobility/routing"
	"github.com/Akashpg-M/polaris/backend/internal/modules/mobility/spatial"
)

const TargetPositionKey = "mobility.target_position"

type Provider struct {
	Spatial               *spatial.Manager
	Routing               routing.RoutingEngine
	RawLimit, RoutedLimit int
	MaxRadius             float64
}

func (p *Provider) Name() string { return "mobility.candidates/v1" }
func (p *Provider) Supports(req extension.CandidateRequest) bool {
	_, latOK := numeric(req.Context["target_latitude"])
	_, lonOK := numeric(req.Context["target_longitude"])
	return latOK && lonOK && p.Spatial != nil
}
func (p *Provider) Candidates(ctx context.Context, req extension.CandidateRequest) ([]extension.Candidate, error) {
	lat, latOK := numeric(req.Context["target_latitude"])
	lon, lonOK := numeric(req.Context["target_longitude"])
	if !latOK || !lonOK {
		return nil, errors.New("mobility target is missing")
	}
	target := model.Position{Latitude: lat, Longitude: lon}
	radius := p.MaxRadius
	if v, ok := req.Context["maximum_distance_meters"].(float64); ok && v > 0 && v < radius {
		radius = v
	}
	near, err := p.Spatial.Nearby(req.TenantID, target, radius, p.RawLimit)
	if err != nil {
		return nil, err
	}
	type ranked struct {
		candidate     extension.Candidate
		eta, distance float64
		hasRoute      bool
	}
	rankedItems := make([]ranked, 0, len(near))
	eligible := make(map[string]struct{}, len(req.EligibleDeviceIDs))
	for _, id := range req.EligibleDeviceIDs {
		eligible[id] = struct{}{}
	}
	routed := 0
	for _, c := range near {
		if len(eligible) > 0 {
			if _, ok := eligible[c.State.DeviceID]; !ok {
				continue
			}
		}
		item := ranked{candidate: extension.Candidate{DeviceID: c.State.DeviceID, Attributes: map[string]any{"distance_meters": c.DistanceMeters, "spatial_observed_at": c.State.ObservedAt, "mobility_profile": c.State.MobilityProfile}}, distance: c.DistanceMeters}
		if p.Routing != nil && routed < p.RoutedLimit && c.State.MobilityProfile == model.MobilityRoadVehicle {
			routed++
			routingStarted := time.Now()
			route, e := p.Routing.Route(ctx, routing.RouteRequest{TenantID: req.TenantID, MobilityProfile: c.State.MobilityProfile, Origin: c.State.ReportedPosition, Destination: target, Policy: routing.RouteFastest})
			if req.Timing != nil {
				req.Timing.RoutingDuration += time.Since(routingStarted)
			}
			if e == nil {
				item.hasRoute = true
				item.eta = route.EstimatedTime.Seconds()
				item.candidate.Attributes["route_eta_seconds"] = item.eta
				item.candidate.Attributes["routing_snapshot_version"] = route.SnapshotVersion
			}
		}
		score := item.distance
		if item.hasRoute {
			score = item.eta
		}
		item.candidate.DomainScore = &score
		rankedItems = append(rankedItems, item)
	}
	sort.Slice(rankedItems, func(i, j int) bool {
		if rankedItems[i].hasRoute != rankedItems[j].hasRoute {
			return rankedItems[i].hasRoute
		}
		if *rankedItems[i].candidate.DomainScore != *rankedItems[j].candidate.DomainScore {
			return *rankedItems[i].candidate.DomainScore < *rankedItems[j].candidate.DomainScore
		}
		return rankedItems[i].candidate.DeviceID < rankedItems[j].candidate.DeviceID
	})
	limit := req.Limit
	if limit <= 0 || limit > len(rankedItems) {
		limit = len(rankedItems)
	}
	out := make([]extension.Candidate, limit)
	for i := 0; i < limit; i++ {
		out[i] = rankedItems[i].candidate
	}
	return out, nil
}

func numeric(v any) (float64, bool) { n, ok := v.(float64); return n, ok }
```

---

## internal\modules\mobility\model\model.go

```
package model

import "time"

type MobilityProfile string

const (
	MobilityRoadVehicle MobilityProfile = "ROAD_VEHICLE"
	MobilityGroundRobot MobilityProfile = "GROUND_ROBOT"
	MobilityAerialDrone MobilityProfile = "AERIAL_DRONE"
	MobilityStatic      MobilityProfile = "STATIC"
)

type Position struct {
	Latitude  float64  `json:"latitude"`
	Longitude float64  `json:"longitude"`
	AltitudeM *float64 `json:"altitude_meters,omitempty"`
}

type SpatialState struct {
	TenantID         string          `json:"tenant_id"`
	DeviceID         string          `json:"device_id"`
	Position         Position        `json:"position"`
	ReportedPosition Position        `json:"reported_position"`
	HeadingDegrees   *float64        `json:"heading_degrees,omitempty"`
	SpeedMPS         *float64        `json:"speed_mps,omitempty"`
	MobilityProfile  MobilityProfile `json:"mobility_profile"`
	H3Cell           uint64          `json:"h3_cell"`
	ObservedAt       time.Time       `json:"observed_at"`
	IndexedAt        time.Time       `json:"indexed_at"`
	BootID           string          `json:"boot_id"`
	BootStartedAt    time.Time       `json:"boot_started_at"`
	SequenceNumber   uint64          `json:"sequence_number"`
	Quality          MobilityQuality `json:"quality"`
}

type MobilityQuality struct {
	Valid      bool     `json:"valid"`
	Confidence float64  `json:"confidence"`
	Anomalies  []string `json:"anomalies,omitempty"`
}

func (s SpatialState) NewerThan(current SpatialState) bool {
	if current.DeviceID == "" {
		return true
	}
	if s.BootID == current.BootID {
		return s.SequenceNumber > current.SequenceNumber
	}
	return s.BootStartedAt.After(current.BootStartedAt)
}
```

---

## internal\modules\mobility\planning\task_planner.go

```
package planning

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/Akashpg-M/polaris/backend/internal/core/extension"
	taskcore "github.com/Akashpg-M/polaris/backend/internal/core/task"
	"github.com/Akashpg-M/polaris/backend/internal/modules/mobility/model"
	"github.com/Akashpg-M/polaris/backend/internal/modules/mobility/routing"
)

type Planner struct {
	SpatialState func(tenant, device string) (model.SpatialState, bool)
	Routing      routing.RoutingEngine
	MaxPlanAge   time.Duration
}

func (p *Planner) Name() string { return "mobility.task-planner/v1" }
func (p *Planner) Supports(t taskcore.Task) bool {
	if t.TaskType != "NAVIGATE" && t.TaskType != "RELOCATE" {
		return false
	}
	var requirements taskcore.Requirements
	return json.Unmarshal(t.Requirements, &requirements) == nil && requirements.PlanningMode == taskcore.PlanningPolarisRequired
}

type target struct {
	Lat       *float64                `json:"lat"`
	Lon       *float64                `json:"lon"`
	Latitude  *float64                `json:"latitude"`
	Longitude *float64                `json:"longitude"`
	Policy    routing.RouteCostPolicy `json:"policy"`
}
type routePayload struct {
	RouteID                string                  `json:"route_id"`
	RouteSchemaVersion     uint32                  `json:"route_schema_version"`
	RoadGraphVersion       string                  `json:"road_graph_version"`
	RoutingSnapshotVersion uint64                  `json:"routing_snapshot_version"`
	GeneratedAt            time.Time               `json:"generated_at"`
	ValidUntil             time.Time               `json:"valid_until"`
	Origin                 model.Position          `json:"origin"`
	Destination            model.Position          `json:"destination"`
	Waypoints              []model.Position        `json:"waypoints"`
	DistanceMeters         float64                 `json:"distance_meters"`
	EstimatedDurationMS    int64                   `json:"estimated_duration_ms"`
	Policy                 routing.RouteCostPolicy `json:"policy"`
}

func (p *Planner) Plan(ctx context.Context, req extension.PlanningRequest) (extension.ExecutionPlan, error) {
	if p.Routing == nil {
		return extension.ExecutionPlan{}, routing.ErrUnavailable
	}
	if req.Task.AssignedDeviceID == nil {
		return extension.ExecutionPlan{}, errors.New("planning requires selected device")
	}
	state, ok := p.SpatialState(req.Task.TenantID, *req.Task.AssignedDeviceID)
	if !ok {
		return extension.ExecutionPlan{}, errors.New("selected device has no active spatial state")
	}
	var t target
	if json.Unmarshal(req.Task.Target, &t) != nil {
		return extension.ExecutionPlan{}, errors.New("invalid mobility target")
	}
	lat, lon := t.Lat, t.Lon
	if lat == nil {
		lat = t.Latitude
	}
	if lon == nil {
		lon = t.Longitude
	}
	if lat == nil || lon == nil {
		return extension.ExecutionPlan{}, errors.New("mobility target coordinates required")
	}
	policy := t.Policy
	if policy == "" {
		policy = routing.RouteFastest
	}
	destination := model.Position{Latitude: *lat, Longitude: *lon}
	route, err := p.Routing.Route(ctx, routing.RouteRequest{TenantID: req.Task.TenantID, MobilityProfile: state.MobilityProfile, Origin: state.ReportedPosition, Destination: destination, Policy: policy})
	if err != nil {
		if errors.Is(err, routing.ErrUnsupportedProfile) {
			return extension.ExecutionPlan{}, extension.ErrPlanningUnsupported
		}
		return extension.ExecutionPlan{}, err
	}
	now := time.Now().UTC()
	valid := now.Add(p.MaxPlanAge)
	if req.Task.ExpiresAt.Before(valid) {
		valid = req.Task.ExpiresAt
	}
	payload, err := json.Marshal(routePayload{route.RouteID, 1, route.GraphVersion, route.SnapshotVersion, now, valid, state.ReportedPosition, destination, route.Waypoints, route.DistanceMeters, route.EstimatedTime.Milliseconds(), route.Policy})
	if err != nil {
		return extension.ExecutionPlan{}, err
	}
	return extension.ExecutionPlan{PlannerName: p.Name(), SchemaVersion: 1, CommandType: req.Task.TaskType, Payload: payload, GeneratedAt: now, ValidUntil: &valid, Metadata: map[string]any{"route_id": route.RouteID, "road_graph_version": route.GraphVersion, "routing_snapshot_version": route.SnapshotVersion}}, nil
}
```

---

## internal\modules\mobility\planning\task_planner_test.go

```
package planning

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/Akashpg-M/polaris/backend/internal/core/extension"
	taskcore "github.com/Akashpg-M/polaris/backend/internal/core/task"
	"github.com/Akashpg-M/polaris/backend/internal/modules/mobility/model"
	"github.com/Akashpg-M/polaris/backend/internal/modules/mobility/routing"
)

type fakeRouter struct{}

func (fakeRouter) Route(context.Context, routing.RouteRequest) (routing.RouteResult, error) {
	return routing.RouteResult{RouteID: "route-1", GraphVersion: "chennai-v1", SnapshotVersion: 42, Policy: routing.RouteFastest, DistanceMeters: 1200, EstimatedTime: 90 * time.Second, Waypoints: []model.Position{{Latitude: 13, Longitude: 80}, {Latitude: 13.01, Longitude: 80.01}}, EdgeIDs: []routing.EdgeID{1}}, nil
}
func TestNavigatePlanContainsVersionAndValidity(t *testing.T) {
	device := "vehicle-1"
	state := model.SpatialState{TenantID: "tenant", DeviceID: device, ReportedPosition: model.Position{Latitude: 13, Longitude: 80}, MobilityProfile: model.MobilityRoadVehicle}
	planner := &Planner{SpatialState: func(string, string) (model.SpatialState, bool) { return state, true }, Routing: fakeRouter{}, MaxPlanAge: 2 * time.Minute}
	task := taskcore.Task{TaskID: "task-1", TenantID: "tenant", TaskType: "NAVIGATE", Target: []byte(`{"lat":13.01,"lon":80.01}`), AssignedDeviceID: &device, ExpiresAt: time.Now().Add(5 * time.Minute)}
	plan, err := planner.Plan(context.Background(), extension.PlanningRequest{Task: task})
	if err != nil {
		t.Fatal(err)
	}
	var payload routePayload
	if err = json.Unmarshal(plan.Payload, &payload); err != nil {
		t.Fatal(err)
	}
	if payload.RouteID != "route-1" || payload.RoadGraphVersion != "chennai-v1" || payload.RoutingSnapshotVersion != 42 || payload.ValidUntil.IsZero() {
		t.Fatalf("incomplete plan: %#v", payload)
	}
	saved := append([]byte(nil), plan.Payload...)
	if string(saved) != string(plan.Payload) {
		t.Fatal("durable replay payload changed")
	}
}

func TestPlannerOnlyClaimsPolarisRequiredMobilityTasks(t *testing.T) {
	planner := &Planner{}
	required, _ := json.Marshal(taskcore.Requirements{PlanningMode: taskcore.PlanningPolarisRequired})
	local, _ := json.Marshal(taskcore.Requirements{PlanningMode: taskcore.PlanningDeviceLocal})
	if !planner.Supports(taskcore.Task{TaskType: "NAVIGATE", Requirements: required}) {
		t.Fatal("Polaris-required navigation was not claimed")
	}
	if planner.Supports(taskcore.Task{TaskType: "NAVIGATE", Requirements: local}) {
		t.Fatal("device-local navigation must use the generic planner")
	}
}
```

---

## internal\modules\mobility\routing\engine.go

```
package routing

import (
	"context"
	"errors"
	"sync"
	"sync/atomic"
	"time"

	"github.com/Akashpg-M/polaris/backend/internal/core/auth"
	"github.com/Akashpg-M/polaris/backend/internal/modules/mobility/model"
	"github.com/Akashpg-M/polaris/backend/internal/modules/mobility/spatial"
)

type routeJob struct {
	ctx    context.Context
	req    RouteRequest
	result chan routeResponse
}
type routeResponse struct {
	route RouteResult
	err   error
}
type EngineConfig struct {
	Workers, QueueCapacity, MaxExpansions, MaxConcurrentPerTenant int
	Timeout                                                       time.Duration
}
type Engine struct {
	graph     *RoadGraph
	snapshots *SnapshotStore
	cfg       EngineConfig
	jobs      chan routeJob
	closed    chan struct{}
	wg        sync.WaitGroup
	started   atomic.Bool
	tenantMu  sync.Mutex
	tenant    map[string]*tenantLimiter
	requests  atomic.Uint64
	busy      atomic.Uint64
}

type tenantLimiter struct {
	slots chan struct{}
	refs  int
}

type EngineStats struct {
	Requests      uint64 `json:"requests"`
	Busy          uint64 `json:"routing_busy"`
	QueueDepth    int    `json:"queue_depth"`
	QueueCapacity int    `json:"queue_capacity"`
	ActiveTenants int    `json:"active_tenants"`
}

func NewEngine(g *RoadGraph, s *SnapshotStore, c EngineConfig) *Engine {
	return &Engine{graph: g, snapshots: s, cfg: c, jobs: make(chan routeJob, c.QueueCapacity), closed: make(chan struct{}), tenant: map[string]*tenantLimiter{}}
}
func (e *Engine) Start() {
	if !e.started.CompareAndSwap(false, true) {
		return
	}
	for i := 0; i < e.cfg.Workers; i++ {
		e.wg.Add(1)
		go func() {
			defer e.wg.Done()
			for {
				select {
				case <-e.closed:
					return
				case j := <-e.jobs:
					r, err := e.route(j.ctx, j.req)
					select {
					case j.result <- routeResponse{r, err}:
					case <-j.ctx.Done():
					}
				}
			}
		}()
	}
}
func (e *Engine) acquireTenant(id string) (*tenantLimiter, bool) {
	e.tenantMu.Lock()
	defer e.tenantMu.Unlock()
	limiter := e.tenant[id]
	if limiter == nil {
		limiter = &tenantLimiter{slots: make(chan struct{}, e.cfg.MaxConcurrentPerTenant)}
		e.tenant[id] = limiter
	}
	select {
	case limiter.slots <- struct{}{}:
		limiter.refs++
		return limiter, true
	default:
		return nil, false
	}
}
func (e *Engine) releaseTenant(id string, limiter *tenantLimiter) {
	e.tenantMu.Lock()
	defer e.tenantMu.Unlock()
	<-limiter.slots
	limiter.refs--
	if limiter.refs == 0 && e.tenant[id] == limiter {
		delete(e.tenant, id)
	}
}
func (e *Engine) Route(ctx context.Context, req RouteRequest) (RouteResult, error) {
	e.requests.Add(1)
	if !e.started.Load() {
		return RouteResult{}, ErrUnavailable
	}
	if req.MobilityProfile != model.MobilityRoadVehicle {
		return RouteResult{}, ErrUnsupportedProfile
	}
	limiter, acquired := e.acquireTenant(req.TenantID)
	if !acquired {
		e.busy.Add(1)
		return RouteResult{}, ErrBusy
	}
	defer e.releaseTenant(req.TenantID, limiter)
	timeout := e.cfg.Timeout
	if deadline, ok := ctx.Deadline(); ok && time.Until(deadline) < timeout {
		timeout = time.Until(deadline)
	}
	routeCtx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()
	result := make(chan routeResponse, 1)
	select {
	case e.jobs <- routeJob{routeCtx, req, result}:
	case <-routeCtx.Done():
		return RouteResult{}, classifyContext(routeCtx.Err())
	default:
		e.busy.Add(1)
		return RouteResult{}, ErrBusy
	}
	select {
	case r := <-result:
		return r.route, r.err
	case <-routeCtx.Done():
		return RouteResult{}, classifyContext(routeCtx.Err())
	}
}
func classifyContext(err error) error {
	if errors.Is(err, context.DeadlineExceeded) {
		return ErrTimeout
	}
	return err
}
func (e *Engine) route(ctx context.Context, req RouteRequest) (RouteResult, error) {
	if err := spatial.ValidatePosition(req.Origin); err != nil {
		return RouteResult{}, ErrOutsideRegion
	}
	if err := spatial.ValidatePosition(req.Destination); err != nil {
		return RouteResult{}, ErrOutsideRegion
	}
	from, err := e.graph.nodeIndex.Nearest(ctx, req.Origin)
	if err != nil {
		return RouteResult{}, ErrNoRoadNode
	}
	to, err := e.graph.nodeIndex.Nearest(ctx, req.Destination)
	if err != nil {
		return RouteResult{}, ErrNoRoadNode
	}
	snapshot := e.snapshots.Load()
	found, err := AStar(ctx, e.graph, snapshot, from.ID, to.ID, req.Policy, e.cfg.MaxExpansions)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			return RouteResult{}, ErrTimeout
		}
		return RouteResult{}, err
	}
	out := routeResult(e.graph, snapshot, req.Policy, found)
	out.RouteID = auth.NewID()
	if len(out.Waypoints) == 0 {
		out.Waypoints = []model.Position{req.Origin}
	}
	return out, nil
}
func (e *Engine) SnapshotStore() *SnapshotStore { return e.snapshots }
func (e *Engine) Stats() EngineStats {
	e.tenantMu.Lock()
	activeTenants := len(e.tenant)
	e.tenantMu.Unlock()
	return EngineStats{Requests: e.requests.Load(), Busy: e.busy.Load(), QueueDepth: len(e.jobs), QueueCapacity: cap(e.jobs), ActiveTenants: activeTenants}
}
func (e *Engine) Close() {
	if e.started.CompareAndSwap(true, false) {
		close(e.closed)
		e.wg.Wait()
	}
}
```

---

## internal\modules\mobility\routing\graph.go

```
package routing

import (
	"errors"
	"sort"
)

type RoadGraph struct {
	version     string
	nodes       map[int64]RoadNode
	edges       []RoadEdge
	adjacency   map[int64][]int
	incident    map[int64][]int
	nodeIndex   RoadNodeIndex
	maxSpeedMPS float64
}

func (g *RoadGraph) Version() string                { return g.version }
func (g *RoadGraph) NodeCount() int                 { return len(g.nodes) }
func (g *RoadGraph) EdgeCount() int                 { return len(g.edges) }
func (g *RoadGraph) Node(id int64) (RoadNode, bool) { v, ok := g.nodes[id]; return v, ok }
func (g *RoadGraph) Edge(index int) RoadEdge        { return g.edges[index] }
func (g *RoadGraph) Outgoing(id int64) []int        { return g.adjacency[id] }
func (g *RoadGraph) NodeIndex() RoadNodeIndex       { return g.nodeIndex }
func (g *RoadGraph) MaxSpeedMPS() float64           { return g.maxSpeedMPS }

type GraphBuilder struct {
	version string
	nodes   map[int64]RoadNode
	edges   []RoadEdge
	seen    map[[3]int64]struct{}
}

func NewGraphBuilder(version string) *GraphBuilder {
	return &GraphBuilder{version: version, nodes: map[int64]RoadNode{}, seen: map[[3]int64]struct{}{}}
}
func (b *GraphBuilder) AddNode(n RoadNode) error {
	if n.ID == 0 {
		return errors.New("invalid road node id")
	}
	if _, ok := b.nodes[n.ID]; ok {
		return nil
	}
	b.nodes[n.ID] = n
	return nil
}
func (b *GraphBuilder) AddEdge(e RoadEdge) error {
	if e.FromID == e.ToID || e.DistanceM <= 0 || e.BaseTravelTime <= 0 {
		return errors.New("invalid road edge")
	}
	if _, ok := b.nodes[e.FromID]; !ok {
		return ErrNoRoadNode
	}
	if _, ok := b.nodes[e.ToID]; !ok {
		return ErrNoRoadNode
	}
	k := [3]int64{e.FromID, e.ToID, int64(e.ID)}
	if _, ok := b.seen[k]; ok {
		return nil
	}
	b.seen[k] = struct{}{}
	b.edges = append(b.edges, e)
	return nil
}
func (b *GraphBuilder) Build() (*RoadGraph, error) {
	if len(b.nodes) == 0 {
		return nil, ErrNoRoadNode
	}
	sort.Slice(b.edges, func(i, j int) bool {
		if b.edges[i].FromID != b.edges[j].FromID {
			return b.edges[i].FromID < b.edges[j].FromID
		}
		if b.edges[i].ToID != b.edges[j].ToID {
			return b.edges[i].ToID < b.edges[j].ToID
		}
		return b.edges[i].ID < b.edges[j].ID
	})
	g := &RoadGraph{version: b.version, nodes: map[int64]RoadNode{}, edges: append([]RoadEdge(nil), b.edges...), adjacency: map[int64][]int{}, incident: map[int64][]int{}}
	nodes := make([]RoadNode, 0, len(b.nodes))
	for id, n := range b.nodes {
		g.nodes[id] = n
		nodes = append(nodes, n)
	}
	for i, e := range g.edges {
		g.adjacency[e.FromID] = append(g.adjacency[e.FromID], i)
		g.incident[e.FromID] = append(g.incident[e.FromID], i)
		g.incident[e.ToID] = append(g.incident[e.ToID], i)
		speed := e.DistanceM / e.BaseTravelTime.Seconds()
		if speed > g.maxSpeedMPS {
			g.maxSpeedMPS = speed
		}
	}
	for id, edges := range g.incident {
		if len(edges) > 2 {
			n := g.nodes[id]
			n.Type = NodeIntersection
			g.nodes[id] = n
		}
	}
	if g.maxSpeedMPS <= 0 {
		g.maxSpeedMPS = 40
	}
	g.nodeIndex = NewKDTree(nodes)
	return g, nil
}
```

---

## internal\modules\mobility\routing\kdtree.go

```
package routing

import (
	"context"
	"math"
	"sort"

	"github.com/Akashpg-M/polaris/backend/internal/modules/mobility/model"
)

type point3 struct {
	x, y, z float64
	n       RoadNode
}
type kdNode struct {
	p           point3
	axis        int
	left, right *kdNode
}
type KDTree struct{ root *kdNode }

func sphere(p model.Position) (x, y, z float64) {
	lat, lon := p.Latitude*math.Pi/180, p.Longitude*math.Pi/180
	return math.Cos(lat) * math.Cos(lon), math.Cos(lat) * math.Sin(lon), math.Sin(lat)
}
func NewKDTree(nodes []RoadNode) *KDTree {
	pts := make([]point3, len(nodes))
	for i, n := range nodes {
		pts[i].x, pts[i].y, pts[i].z = sphere(n.Position)
		pts[i].n = n
	}
	return &KDTree{root: buildKD(pts, 0)}
}
func buildKD(p []point3, depth int) *kdNode {
	if len(p) == 0 {
		return nil
	}
	axis := depth % 3
	value := func(v point3) float64 {
		if axis == 0 {
			return v.x
		}
		if axis == 1 {
			return v.y
		}
		return v.z
	}
	sort.Slice(p, func(i, j int) bool { return value(p[i]) < value(p[j]) })
	m := len(p) / 2
	return &kdNode{p: p[m], axis: axis, left: buildKD(append([]point3(nil), p[:m]...), depth+1), right: buildKD(append([]point3(nil), p[m+1:]...), depth+1)}
}
func (t *KDTree) Nearest(ctx context.Context, p model.Position) (RoadNode, error) {
	if t.root == nil {
		return RoadNode{}, ErrNoRoadNode
	}
	x, y, z := sphere(p)
	target := [3]float64{x, y, z}
	best := t.root.p
	bestD := math.Inf(1)
	var walk func(*kdNode)
	walk = func(n *kdNode) {
		if n == nil || ctx.Err() != nil {
			return
		}
		v := [3]float64{n.p.x, n.p.y, n.p.z}
		d := (v[0]-x)*(v[0]-x) + (v[1]-y)*(v[1]-y) + (v[2]-z)*(v[2]-z)
		if d < bestD {
			bestD, best = d, n.p
		}
		delta := target[n.axis] - v[n.axis]
		first, second := n.left, n.right
		if delta > 0 {
			first, second = second, first
		}
		walk(first)
		if delta*delta < bestD {
			walk(second)
		}
	}
	walk(t.root)
	if ctx.Err() != nil {
		return RoadNode{}, ctx.Err()
	}
	return best.n, nil
}
```

---

## internal\modules\mobility\routing\osm.go

```
package routing

import (
	"context"
	"io"
	"os"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/Akashpg-M/polaris/backend/internal/modules/mobility/model"
	"github.com/Akashpg-M/polaris/backend/internal/modules/mobility/spatial"
	"github.com/qedus/osmpbf"
)

type osmWay struct {
	id    int64
	nodes []int64
	tags  map[string]string
}

func LoadOSMPBF(ctx context.Context, path, version string) (*RoadGraph, error) {
	f, e := os.Open(path)
	if e != nil {
		return nil, e
	}
	defer f.Close()
	d := osmpbf.NewDecoder(f)
	if e = d.Start(runtime.GOMAXPROCS(0)); e != nil {
		return nil, e
	}
	nodes := map[int64]RoadNode{}
	ways := []osmWay{}
	for {
		if ctx.Err() != nil {
			return nil, ctx.Err()
		}
		v, e := d.Decode()
		if e == io.EOF {
			break
		}
		if e != nil {
			return nil, e
		}
		switch x := v.(type) {
		case *osmpbf.Node:
			nodes[x.ID] = RoadNode{ID: x.ID, Position: model.Position{Latitude: x.Lat, Longitude: x.Lon}, Type: NodeRoad}
		case *osmpbf.Way:
			if roadClass(x.Tags["highway"]) != "" {
				ways = append(ways, osmWay{x.ID, append([]int64(nil), x.NodeIDs...), x.Tags})
			}
		}
	}
	b := NewGraphBuilder(version)
	usedNodes := make(map[int64]struct{})
	for _, way := range ways {
		for _, id := range way.nodes {
			usedNodes[id] = struct{}{}
		}
	}
	for id := range usedNodes {
		if n, ok := nodes[id]; ok {
			_ = b.AddNode(n)
		}
	}
	var edgeID int64 = 1
	for _, w := range ways {
		class := roadClass(w.tags["highway"])
		speed := parseSpeed(w.tags["maxspeed"], defaultSpeed(class))
		reverseOnly := w.tags["oneway"] == "-1"
		oneway := reverseOnly || w.tags["oneway"] == "yes" || w.tags["oneway"] == "1" || w.tags["junction"] == "roundabout"
		for i := 0; i+1 < len(w.nodes); i++ {
			a, aok := nodes[w.nodes[i]]
			z, zok := nodes[w.nodes[i+1]]
			if !aok || !zok {
				continue
			}
			distance := spatial.DistanceMeters(a.Position, z.Position)
			from, to := a.ID, z.ID
			if reverseOnly {
				from, to = to, from
			}
			_ = b.AddEdge(RoadEdge{ID: EdgeID(edgeID), FromID: from, ToID: to, DistanceM: distance, BaseTravelTime: time.Duration(distance / speed * float64(time.Second)), RoadClass: class})
			edgeID++
			if !oneway {
				_ = b.AddEdge(RoadEdge{ID: EdgeID(edgeID), FromID: to, ToID: from, DistanceM: distance, BaseTravelTime: time.Duration(distance / speed * float64(time.Second)), RoadClass: class})
				edgeID++
			}
		}
	}
	return b.Build()
}
func roadClass(v string) RoadClass {
	switch v {
	case "motorway", "motorway_link", "trunk", "trunk_link":
		return RoadMotorway
	case "primary", "primary_link":
		return RoadPrimary
	case "secondary", "secondary_link":
		return RoadSecondary
	case "tertiary", "tertiary_link":
		return RoadTertiary
	case "residential", "living_street":
		return RoadResidential
	case "unclassified", "service":
		return RoadUnclassified
	}
	return ""
}
func defaultSpeed(c RoadClass) float64 {
	switch c {
	case RoadMotorway:
		return 27.8
	case RoadPrimary:
		return 19.4
	case RoadSecondary:
		return 16.7
	case RoadTertiary:
		return 13.9
	default:
		return 8.3
	}
}
func parseSpeed(raw string, fallback float64) float64 {
	if raw == "" {
		return fallback
	}
	value := strings.Fields(strings.Split(raw, ";")[0])
	if len(value) == 0 {
		return fallback
	}
	n, e := strconv.ParseFloat(value[0], 64)
	if e != nil || n <= 0 {
		return fallback
	}
	if strings.Contains(strings.ToLower(raw), "mph") {
		return n * .44704
	}
	return n / 3.6
}
```

---

## internal\modules\mobility\routing\routing_test.go

```
package routing

import (
	"context"
	"errors"
	"math"
	"sync"
	"testing"
	"time"

	"github.com/Akashpg-M/polaris/backend/internal/modules/mobility/model"
)

func testGraph(t *testing.T) *RoadGraph {
	b := NewGraphBuilder("test-v1")
	nodes := []RoadNode{{1, model.Position{Latitude: 0, Longitude: 0}, NodeRoad}, {2, model.Position{Latitude: 0, Longitude: .01}, NodeRoad}, {3, model.Position{Latitude: .01, Longitude: .01}, NodeRoad}, {4, model.Position{Latitude: .01, Longitude: 0}, NodeRoad}}
	for _, n := range nodes {
		if err := b.AddNode(n); err != nil {
			t.Fatal(err)
		}
	}
	edges := []RoadEdge{{1, 1, 2, 1000, time.Minute, RoadResidential}, {2, 2, 3, 1000, time.Minute, RoadResidential}, {3, 1, 4, 1200, time.Minute, RoadResidential}, {4, 4, 3, 1200, time.Minute, RoadResidential}}
	for _, e := range edges {
		if err := b.AddEdge(e); err != nil {
			t.Fatal(err)
		}
	}
	g, err := b.Build()
	if err != nil {
		t.Fatal(err)
	}
	return g
}
func TestAStarEqualsDijkstraAndHonorsDirection(t *testing.T) {
	g := testGraph(t)
	s := NewSnapshotStore(g).Load()
	for _, policy := range []RouteCostPolicy{RouteShortest, RouteFastest} {
		a, ea := AStar(context.Background(), g, s, 1, 3, policy, 100)
		d, ed := Dijkstra(context.Background(), g, s, 1, 3, policy, 100)
		if ea != nil || ed != nil || math.Abs(a.cost-d.cost) > 1e-9 {
			t.Fatalf("policy %s A*=%v/%v Dijkstra=%v/%v", policy, a.cost, ea, d.cost, ed)
		}
	}
	if _, err := AStar(context.Background(), g, s, 3, 1, RouteShortest, 100); !errors.Is(err, ErrNoRoute) {
		t.Fatalf("one-way route should fail, got %v", err)
	}
}
func TestSnapshotChangesCostWithoutMutatingTopology(t *testing.T) {
	g := testGraph(t)
	store := NewSnapshotStore(g)
	costs := append([]float64(nil), store.Load().EdgeCosts...)
	costs[0], costs[2] = 1000, 1000
	if err := store.Swap(RoutingCostSnapshot{Version: 2, EdgeCosts: costs}); err != nil {
		t.Fatal(err)
	}
	result, err := AStar(context.Background(), g, store.Load(), 1, 3, RouteFastest, 100)
	if err != nil {
		t.Fatal(err)
	}
	if result.edgeIndexes[0] != 1 {
		t.Fatalf("dynamic fastest route did not change: %v", result.edgeIndexes)
	}
	if g.EdgeCount() != 4 {
		t.Fatal("topology mutated")
	}
}
func TestSnapshotConcurrentReadersUseOneVersion(t *testing.T) {
	g := testGraph(t)
	store := NewSnapshotStore(g)
	var wg sync.WaitGroup
	versions := make(chan uint64, 200)
	for n := 0; n < 100; n++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			snapshot := store.Load()
			_, _ = AStar(context.Background(), g, snapshot, 1, 3, RouteFastest, 100)
			versions <- snapshot.Version
		}()
	}
	for version := uint64(2); version <= 5; version++ {
		old := store.Load()
		_ = store.Swap(RoutingCostSnapshot{Version: version, EdgeCosts: append([]float64(nil), old.EdgeCosts...)})
	}
	wg.Wait()
	close(versions)
	for version := range versions {
		if version < 1 || version > 5 {
			t.Fatalf("mixed/invalid snapshot version %d", version)
		}
	}
}
func TestKDTreeAndNodeTypeSemantics(t *testing.T) {
	g := testGraph(t)
	node, err := g.NodeIndex().Nearest(context.Background(), model.Position{Latitude: 0, Longitude: .0099})
	if err != nil || node.ID != 2 {
		t.Fatalf("nearest node=%v err=%v", node.ID, err)
	}
	if NodeRoad == NodeIntersection || NodeIntersection == NodeChargingStation {
		t.Fatal("NodeType enum values overlap")
	}
}
func TestBoundedEngineIsolationAndUnsupportedProfile(t *testing.T) {
	g := testGraph(t)
	e := NewEngine(g, NewSnapshotStore(g), EngineConfig{Workers: 1, QueueCapacity: 1, MaxExpansions: 100, MaxConcurrentPerTenant: 1, Timeout: time.Second})
	e.Start()
	defer e.Close()
	_, err := e.Route(context.Background(), RouteRequest{TenantID: "t", MobilityProfile: model.MobilityAerialDrone, Origin: g.nodes[1].Position, Destination: g.nodes[3].Position, Policy: RouteFastest})
	if !errors.Is(err, ErrUnsupportedProfile) {
		t.Fatalf("expected unsupported profile, got %v", err)
	}
}

func TestRoutingQueueSaturationIsExplicit(t *testing.T) {
	g := testGraph(t)
	engine := NewEngine(g, NewSnapshotStore(g), EngineConfig{Workers: 0, QueueCapacity: 1, MaxExpansions: 100, MaxConcurrentPerTenant: 1, Timeout: 50 * time.Millisecond})
	engine.Start()
	defer engine.Close()
	req := RouteRequest{TenantID: "tenant-a", MobilityProfile: model.MobilityRoadVehicle, Origin: g.nodes[1].Position, Destination: g.nodes[3].Position, Policy: RouteFastest}
	firstReq := req
	done := make(chan error, 1)
	go func() { _, err := engine.Route(context.Background(), firstReq); done <- err }()
	time.Sleep(5 * time.Millisecond)
	req.TenantID = "tenant-b"
	if _, err := engine.Route(context.Background(), req); !errors.Is(err, ErrBusy) {
		t.Fatalf("expected ROUTING_BUSY, got %v", err)
	}
	if err := <-done; !errors.Is(err, ErrTimeout) {
		t.Fatalf("queued request should time out, got %v", err)
	}
}
func TestInvalidSnapshotRejected(t *testing.T) {
	g := testGraph(t)
	store := NewSnapshotStore(g)
	if err := store.Swap(RoutingCostSnapshot{Version: 2, EdgeCosts: []float64{1}}); err == nil {
		t.Fatal("partial snapshot accepted")
	}
}

func TestTrafficObservationUpdatesExactlyOneDirectedEdge(t *testing.T) {
	g := testGraph(t)
	store := NewSnapshotStore(g)
	traffic := NewTrafficManager(g, store, time.Minute)
	heading := 90.0
	err := traffic.Observe(context.Background(), TrafficObservation{Position: model.Position{Latitude: 0, Longitude: .005}, HeadingDegrees: &heading, SpeedMPS: 1, ObservedAt: time.Now().UTC()})
	if err != nil {
		t.Fatal(err)
	}
	before := append([]float64(nil), store.Load().EdgeCosts...)
	if err = traffic.Refresh(time.Now().UTC()); err != nil {
		t.Fatal(err)
	}
	changed := 0
	for i, cost := range store.Load().EdgeCosts {
		if math.Abs(cost-before[i]) > 1e-9 {
			changed++
		}
	}
	if changed != 1 || traffic.StateCount() != 1 {
		t.Fatalf("one observation changed %d edges and retained %d states", changed, traffic.StateCount())
	}
}

func BenchmarkAStar(b *testing.B)    { benchmarkSearch(b, true) }
func BenchmarkDijkstra(b *testing.B) { benchmarkSearch(b, false) }
func benchmarkSearch(b *testing.B, astar bool) {
	const side = 30
	builder := NewGraphBuilder("benchmark-grid")
	for y := 0; y < side; y++ {
		for x := 0; x < side; x++ {
			id := int64(y*side + x + 1)
			_ = builder.AddNode(RoadNode{ID: id, Position: model.Position{Latitude: 13 + float64(y)*.001, Longitude: 80 + float64(x)*.001}, Type: NodeRoad})
		}
	}
	edgeID := EdgeID(1)
	for y := 0; y < side; y++ {
		for x := 0; x < side; x++ {
			from := int64(y*side + x + 1)
			if x+1 < side {
				to := from + 1
				_ = builder.AddEdge(RoadEdge{ID: edgeID, FromID: from, ToID: to, DistanceM: 110, BaseTravelTime: 10 * time.Second, RoadClass: RoadResidential})
				edgeID++
				_ = builder.AddEdge(RoadEdge{ID: edgeID, FromID: to, ToID: from, DistanceM: 110, BaseTravelTime: 10 * time.Second, RoadClass: RoadResidential})
				edgeID++
			}
			if y+1 < side {
				to := from + side
				_ = builder.AddEdge(RoadEdge{ID: edgeID, FromID: from, ToID: to, DistanceM: 110, BaseTravelTime: 10 * time.Second, RoadClass: RoadResidential})
				edgeID++
				_ = builder.AddEdge(RoadEdge{ID: edgeID, FromID: to, ToID: from, DistanceM: 110, BaseTravelTime: 10 * time.Second, RoadClass: RoadResidential})
				edgeID++
			}
		}
	}
	g, _ := builder.Build()
	snapshot := NewSnapshotStore(g).Load()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if astar {
			_, _ = AStar(context.Background(), g, snapshot, 1, int64(side*side), RouteFastest, side*side*2)
		} else {
			_, _ = Dijkstra(context.Background(), g, snapshot, 1, int64(side*side), RouteFastest, side*side*2)
		}
	}
}
```

---

## internal\modules\mobility\routing\search.go

```
package routing

import (
	"container/heap"
	"context"
	"math"
	"time"

	"github.com/Akashpg-M/polaris/backend/internal/modules/mobility/model"
	"github.com/Akashpg-M/polaris/backend/internal/modules/mobility/spatial"
)

type queueItem struct {
	node           int64
	priority, cost float64
	index          int
}
type priorityQueue []*queueItem

func (q priorityQueue) Len() int { return len(q) }
func (q priorityQueue) Less(i, j int) bool {
	if q[i].priority != q[j].priority {
		return q[i].priority < q[j].priority
	}
	return q[i].node < q[j].node
}
func (q priorityQueue) Swap(i, j int) { q[i], q[j] = q[j], q[i]; q[i].index = i; q[j].index = j }
func (q *priorityQueue) Push(v any)   { n := v.(*queueItem); n.index = len(*q); *q = append(*q, n) }
func (q *priorityQueue) Pop() any     { old := *q; n := old[len(old)-1]; *q = old[:len(old)-1]; return n }

type searchResult struct {
	edgeIndexes             []int
	cost, distance, seconds float64
	expanded                int
	heapOperations          int
}

type SearchAlgorithm string

const (
	AlgorithmAStar    SearchAlgorithm = "ASTAR"
	AlgorithmDijkstra SearchAlgorithm = "DIJKSTRA"
)

type SearchMetrics struct {
	Algorithm      SearchAlgorithm `json:"algorithm"`
	Cost           float64         `json:"cost"`
	DistanceMeters float64         `json:"distance_meters"`
	ExpandedNodes  int             `json:"expanded_nodes"`
	HeapOperations int             `json:"heap_operations"`
	EdgeCount      int             `json:"edge_count"`
}

func edgeCost(e RoadEdge, index int, policy RouteCostPolicy, s *RoutingCostSnapshot) float64 {
	if policy == RouteShortest {
		return e.DistanceM
	}
	if s != nil && index < len(s.EdgeCosts) && s.EdgeCosts[index] > 0 {
		return s.EdgeCosts[index]
	}
	return e.BaseTravelTime.Seconds()
}
func heuristic(g *RoadGraph, node, target int64, policy RouteCostPolicy) float64 {
	a := g.nodes[node].Position
	b := g.nodes[target].Position
	d := spatial.DistanceMeters(a, b)
	if policy == RouteFastest {
		return d / g.maxSpeedMPS
	}
	return d
}
func runSearch(ctx context.Context, g *RoadGraph, s *RoutingCostSnapshot, from, to int64, policy RouteCostPolicy, maxExpansions int, useHeuristic bool) (searchResult, error) {
	if _, ok := g.nodes[from]; !ok {
		return searchResult{}, ErrNoRoadNode
	}
	if _, ok := g.nodes[to]; !ok {
		return searchResult{}, ErrNoRoadNode
	}
	if from == to {
		return searchResult{}, nil
	}
	cost := map[int64]float64{from: 0}
	previous := map[int64]int{}
	visited := map[int64]bool{}
	q := &priorityQueue{}
	heap.Init(q)
	heap.Push(q, &queueItem{node: from})
	expanded, heapOperations := 0, 1
	for q.Len() > 0 {
		if ctx.Err() != nil {
			return searchResult{}, ctx.Err()
		}
		cur := heap.Pop(q).(*queueItem)
		heapOperations++
		if visited[cur.node] {
			continue
		}
		visited[cur.node] = true
		expanded++
		if expanded > maxExpansions {
			return searchResult{}, ErrTimeout
		}
		if cur.node == to {
			break
		}
		for _, ei := range g.adjacency[cur.node] {
			e := g.edges[ei]
			next := cost[cur.node] + edgeCost(e, ei, policy, s)
			old, ok := cost[e.ToID]
			if !ok || next < old {
				cost[e.ToID] = next
				previous[e.ToID] = ei
				h := 0.0
				if useHeuristic {
					h = heuristic(g, e.ToID, to, policy)
				}
				heap.Push(q, &queueItem{node: e.ToID, cost: next, priority: next + h})
				heapOperations++
			}
		}
	}
	if _, ok := cost[to]; !ok {
		return searchResult{}, ErrNoRoute
	}
	rev := []int{}
	at := to
	for at != from {
		ei, ok := previous[at]
		if !ok {
			return searchResult{}, ErrNoRoute
		}
		rev = append(rev, ei)
		at = g.edges[ei].FromID
	}
	edges := make([]int, len(rev))
	var distance, seconds float64
	for i := range rev {
		ei := rev[len(rev)-1-i]
		edges[i] = ei
		e := g.edges[ei]
		distance += e.DistanceM
		if s != nil && ei < len(s.EdgeCosts) {
			seconds += s.EdgeCosts[ei]
		} else {
			seconds += e.BaseTravelTime.Seconds()
		}
	}
	return searchResult{edgeIndexes: edges, cost: cost[to], distance: distance, seconds: seconds, expanded: expanded, heapOperations: heapOperations}, nil
}

func AStar(ctx context.Context, g *RoadGraph, s *RoutingCostSnapshot, from, to int64, p RouteCostPolicy, max int) (searchResult, error) {
	return runSearch(ctx, g, s, from, to, p, max, true)
}
func Dijkstra(ctx context.Context, g *RoadGraph, s *RoutingCostSnapshot, from, to int64, p RouteCostPolicy, max int) (searchResult, error) {
	return runSearch(ctx, g, s, from, to, p, max, false)
}

// MeasureSearch exposes algorithm work without exposing mutable graph details.
// It is used by the reproducible real-road benchmark and correctness oracle.
func MeasureSearch(ctx context.Context, g *RoadGraph, s *RoutingCostSnapshot, from, to int64, p RouteCostPolicy, max int, algorithm SearchAlgorithm) (SearchMetrics, error) {
	result, err := runSearch(ctx, g, s, from, to, p, max, algorithm == AlgorithmAStar)
	if err != nil {
		return SearchMetrics{}, err
	}
	return SearchMetrics{Algorithm: algorithm, Cost: result.cost, DistanceMeters: result.distance, ExpandedNodes: result.expanded, HeapOperations: result.heapOperations, EdgeCount: len(result.edgeIndexes)}, nil
}

func routeResult(g *RoadGraph, s *RoutingCostSnapshot, p RouteCostPolicy, r searchResult) RouteResult {
	waypoints := []model.Position{}
	ids := make([]EdgeID, len(r.edgeIndexes))
	if len(r.edgeIndexes) > 0 {
		waypoints = append(waypoints, g.nodes[g.edges[r.edgeIndexes[0]].FromID].Position)
	}
	for i, ei := range r.edgeIndexes {
		e := g.edges[ei]
		ids[i] = e.ID
		waypoints = append(waypoints, g.nodes[e.ToID].Position)
	}
	version := uint64(0)
	if s != nil {
		version = s.Version
	}
	return RouteResult{GraphVersion: g.Version(), SnapshotVersion: version, Policy: p, DistanceMeters: r.distance, EstimatedTime: timeDuration(r.seconds), Waypoints: waypoints, EdgeIDs: ids, ExpandedNodes: r.expanded}
}
func timeDuration(seconds float64) time.Duration {
	return time.Duration(math.Round(seconds * float64(time.Second)))
}
```

---

## internal\modules\mobility\routing\snapshot.go

```
package routing

import (
	"errors"
	"sync/atomic"
	"time"
)

type RoutingCostSnapshot struct {
	Version     uint64    `json:"version"`
	GeneratedAt time.Time `json:"generated_at"`
	EdgeCosts   []float64 `json:"edge_costs"`
}
type SnapshotStore struct {
	current atomic.Pointer[RoutingCostSnapshot]
}

func NewSnapshotStore(g *RoadGraph) *SnapshotStore {
	s := &SnapshotStore{}
	cost := make([]float64, len(g.edges))
	for i, e := range g.edges {
		cost[i] = e.BaseTravelTime.Seconds()
	}
	s.current.Store(&RoutingCostSnapshot{Version: 1, GeneratedAt: time.Now().UTC(), EdgeCosts: cost})
	return s
}
func (s *SnapshotStore) Load() *RoutingCostSnapshot { return s.current.Load() }
func (s *SnapshotStore) Swap(next RoutingCostSnapshot) error {
	old := s.current.Load()
	if old != nil && next.Version <= old.Version {
		return errors.New("snapshot version must increase")
	}
	if old != nil && len(next.EdgeCosts) != len(old.EdgeCosts) {
		return errors.New("snapshot edge count mismatch")
	}
	copyCosts := append([]float64(nil), next.EdgeCosts...)
	next.EdgeCosts = copyCosts
	if next.GeneratedAt.IsZero() {
		next.GeneratedAt = time.Now().UTC()
	}
	s.current.Store(&next)
	return nil
}
```

---

## internal\modules\mobility\routing\traffic.go

```
package routing

import (
	"context"
	"errors"
	"math"
	"sync"
	"time"

	"github.com/Akashpg-M/polaris/backend/internal/modules/mobility/model"
)

type EdgeTrafficState struct {
	SampleCount    uint64    `json:"sample_count"`
	EWMASpeedMPS   float64   `json:"ewma_speed_mps"`
	LastSpeedMPS   float64   `json:"last_speed_mps"`
	LastObservedAt time.Time `json:"last_observed_at"`
	Confidence     float64   `json:"confidence"`
}
type TrafficObservation struct {
	Position       model.Position
	HeadingDegrees *float64
	SpeedMPS       float64
	ObservedAt     time.Time
}
type TrafficManager struct {
	mu        sync.RWMutex
	graph     *RoadGraph
	snapshots *SnapshotStore
	states    map[int]EdgeTrafficState
	maxAge    time.Duration
}

func NewTrafficManager(g *RoadGraph, s *SnapshotStore, maxAge time.Duration) *TrafficManager {
	return &TrafficManager{graph: g, snapshots: s, states: map[int]EdgeTrafficState{}, maxAge: maxAge}
}
func (t *TrafficManager) Observe(ctx context.Context, o TrafficObservation) error {
	if o.SpeedMPS < 0 || time.Since(o.ObservedAt) > t.maxAge {
		return errors.New("traffic observation rejected")
	}
	ei, confidence, err := t.match(ctx, o)
	if err != nil {
		return err
	}
	t.mu.Lock()
	s := t.states[ei]
	alpha := .3
	if s.SampleCount == 0 {
		s.EWMASpeedMPS = o.SpeedMPS
	} else {
		s.EWMASpeedMPS = alpha*o.SpeedMPS + (1-alpha)*s.EWMASpeedMPS
	}
	s.SampleCount++
	s.LastSpeedMPS = o.SpeedMPS
	s.LastObservedAt = o.ObservedAt
	s.Confidence = confidence
	t.states[ei] = s
	t.mu.Unlock()
	return nil
}
func (t *TrafficManager) match(ctx context.Context, o TrafficObservation) (int, float64, error) {
	node, err := t.graph.nodeIndex.Nearest(ctx, o.Position)
	if err != nil {
		return 0, 0, err
	}
	best, bestScore := -1, math.Inf(1)
	for _, ei := range t.graph.incident[node.ID] {
		e := t.graph.edges[ei]
		a, b := t.graph.nodes[e.FromID].Position, t.graph.nodes[e.ToID].Position
		distance := segmentDistance(o.Position, a, b)
		score := distance
		if o.HeadingDegrees != nil {
			bearing := initialBearing(a, b)
			delta := math.Abs(*o.HeadingDegrees - bearing)
			if delta > 180 {
				delta = 360 - delta
			}
			score += delta * 2
		}
		if score < bestScore {
			best, bestScore = ei, score
		}
	}
	if best < 0 || bestScore > 100 {
		return 0, 0, ErrNoRoadNode
	}
	return best, math.Max(0, 1-bestScore/100), nil
}
func segmentDistance(p, a, b model.Position) float64 {
	latScale := 111320.0
	lonScale := latScale * math.Cos(p.Latitude*math.Pi/180)
	ax, ay := (a.Longitude-p.Longitude)*lonScale, (a.Latitude-p.Latitude)*latScale
	bx, by := (b.Longitude-p.Longitude)*lonScale, (b.Latitude-p.Latitude)*latScale
	dx, dy := bx-ax, by-ay
	if dx == 0 && dy == 0 {
		return math.Hypot(ax, ay)
	}
	u := -(ax*dx + ay*dy) / (dx*dx + dy*dy)
	u = math.Max(0, math.Min(1, u))
	return math.Hypot(ax+u*dx, ay+u*dy)
}
func initialBearing(a, b model.Position) float64 {
	lat1, lat2 := a.Latitude*math.Pi/180, b.Latitude*math.Pi/180
	dlon := (b.Longitude - a.Longitude) * math.Pi / 180
	y := math.Sin(dlon) * math.Cos(lat2)
	x := math.Cos(lat1)*math.Sin(lat2) - math.Sin(lat1)*math.Cos(lat2)*math.Cos(dlon)
	v := math.Atan2(y, x) * 180 / math.Pi
	if v < 0 {
		v += 360
	}
	return v
}
func (t *TrafficManager) Refresh(now time.Time) error {
	old := t.snapshots.Load()
	if old == nil {
		return ErrUnavailable
	}
	costs := make([]float64, len(t.graph.edges))
	t.mu.RLock()
	states := make(map[int]EdgeTrafficState, len(t.states))
	for edge, state := range t.states {
		states[edge] = state
	}
	t.mu.RUnlock()
	for i, e := range t.graph.edges {
		base := e.BaseTravelTime.Seconds()
		state, ok := states[i]
		if !ok {
			costs[i] = base
			continue
		}
		age := now.Sub(state.LastObservedAt)
		confidence := state.Confidence * math.Exp(-float64(age)/float64(t.maxAge))
		baseSpeed := e.DistanceM / base
		observed := math.Max(.5, state.EWMASpeedMPS)
		multiplier := math.Max(1, baseSpeed/observed)
		costs[i] = base * (1 + confidence*(multiplier-1))
	}
	// Expired observations no longer affect costs and must not accumulate for
	// the lifetime of the process. Compare timestamps before deletion so a
	// concurrent fresh observation for the same edge is never removed.
	t.mu.Lock()
	for edge, state := range states {
		if now.Sub(state.LastObservedAt) > 5*t.maxAge && t.states[edge].LastObservedAt.Equal(state.LastObservedAt) {
			delete(t.states, edge)
		}
	}
	t.mu.Unlock()
	return t.snapshots.Swap(RoutingCostSnapshot{Version: old.Version + 1, GeneratedAt: now.UTC(), EdgeCosts: costs})
}

func (t *TrafficManager) StateCount() int {
	t.mu.RLock()
	defer t.mu.RUnlock()
	return len(t.states)
}

func (t *TrafficManager) OverlayBytes() int64 {
	return int64(len(t.graph.edges)) * 8
}
```

---

## internal\modules\mobility\routing\types.go

```
package routing

import (
	"context"
	"errors"
	"time"

	"github.com/Akashpg-M/polaris/backend/internal/modules/mobility/model"
)

type EdgeID int64
type NodeType uint8

const (
	NodeUnknown NodeType = iota
	NodeRoad
	NodeIntersection
	NodeChargingStation
)

type RoadClass string

const (
	RoadMotorway     RoadClass = "MOTORWAY"
	RoadPrimary      RoadClass = "PRIMARY"
	RoadSecondary    RoadClass = "SECONDARY"
	RoadTertiary     RoadClass = "TERTIARY"
	RoadResidential  RoadClass = "RESIDENTIAL"
	RoadUnclassified RoadClass = "UNCLASSIFIED"
)

type RoadNode struct {
	ID       int64
	Position model.Position
	Type     NodeType
}
type RoadEdge struct {
	ID             EdgeID
	FromID, ToID   int64
	DistanceM      float64
	BaseTravelTime time.Duration
	RoadClass      RoadClass
}
type RouteCostPolicy string

const (
	RouteShortest RouteCostPolicy = "SHORTEST"
	RouteFastest  RouteCostPolicy = "FASTEST"
)

type RouteRequest struct {
	TenantID        string                `json:"tenant_id"`
	MobilityProfile model.MobilityProfile `json:"mobility_profile"`
	Origin          model.Position        `json:"origin"`
	Destination     model.Position        `json:"destination"`
	Policy          RouteCostPolicy       `json:"policy"`
}
type RouteResult struct {
	RouteID         string           `json:"route_id"`
	GraphVersion    string           `json:"road_graph_version"`
	SnapshotVersion uint64           `json:"snapshot_version"`
	Policy          RouteCostPolicy  `json:"policy"`
	DistanceMeters  float64          `json:"distance_meters"`
	EstimatedTime   time.Duration    `json:"estimated_time"`
	Waypoints       []model.Position `json:"waypoints"`
	EdgeIDs         []EdgeID         `json:"edge_ids"`
	ExpandedNodes   int              `json:"expanded_nodes"`
}
type RoutingEngine interface {
	Route(context.Context, RouteRequest) (RouteResult, error)
}
type RoadNodeIndex interface {
	Nearest(context.Context, model.Position) (RoadNode, error)
}

var (
	ErrNoRoute            = errors.New("NO_ROUTE")
	ErrNoRoadNode         = errors.New("NO_ROAD_NODE")
	ErrOutsideRegion      = errors.New("OUTSIDE_ROUTING_REGION")
	ErrUnsupportedProfile = errors.New("UNSUPPORTED_MOBILITY_PROFILE")
	ErrUnavailable        = errors.New("ROUTING_UNAVAILABLE")
	ErrBusy               = errors.New("ROUTING_BUSY")
	ErrTimeout            = errors.New("ROUTING_TIMEOUT")
)
```

---

## internal\modules\mobility\spatial\geo.go

```
package spatial

import (
	"errors"
	"math"

	"github.com/Akashpg-M/polaris/backend/internal/modules/mobility/model"
)

const EarthRadiusMeters = 6371008.8

func NormalizeLongitude(lon float64) float64 {
	if lon >= -180 && lon <= 180 {
		return lon
	}
	lon = math.Mod(lon+180, 360)
	if lon < 0 {
		lon += 360
	}
	return lon - 180
}

func ValidatePosition(p model.Position) error {
	if math.IsNaN(p.Latitude) || math.IsInf(p.Latitude, 0) || p.Latitude < -90 || p.Latitude > 90 {
		return errors.New("latitude outside [-90,90]")
	}
	if math.IsNaN(p.Longitude) || math.IsInf(p.Longitude, 0) || p.Longitude < -180 || p.Longitude > 180 {
		return errors.New("longitude outside [-180,180]")
	}
	return nil
}

func DistanceMeters(a, b model.Position) float64 {
	toRad := math.Pi / 180
	lat1, lat2 := a.Latitude*toRad, b.Latitude*toRad
	dLat := (b.Latitude - a.Latitude) * toRad
	dLonDeg := NormalizeLongitude(b.Longitude - a.Longitude)
	dLon := dLonDeg * toRad
	x := math.Sin(dLat/2)*math.Sin(dLat/2) + math.Cos(lat1)*math.Cos(lat2)*math.Sin(dLon/2)*math.Sin(dLon/2)
	return 2 * EarthRadiusMeters * math.Atan2(math.Sqrt(x), math.Sqrt(1-x))
}

func ValidateObservation(previous *model.SpatialState, next model.SpatialState) model.MobilityQuality {
	q := model.MobilityQuality{Valid: true, Confidence: 1}
	if err := ValidatePosition(next.ReportedPosition); err != nil {
		q.Valid = false
		q.Confidence = 0
		q.Anomalies = append(q.Anomalies, "INVALID_COORDINATES")
	}
	if next.SpeedMPS != nil && (*next.SpeedMPS < 0 || math.IsNaN(*next.SpeedMPS) || math.IsInf(*next.SpeedMPS, 0)) {
		q.Valid = false
		q.Confidence = 0
		q.Anomalies = append(q.Anomalies, "INVALID_SPEED")
	}
	if next.HeadingDegrees != nil {
		v := math.Mod(*next.HeadingDegrees, 360)
		if v < 0 {
			v += 360
		}
		next.HeadingDegrees = &v
	}
	if previous != nil && next.ObservedAt.After(previous.ObservedAt) {
		implied := DistanceMeters(previous.ReportedPosition, next.ReportedPosition) / next.ObservedAt.Sub(previous.ObservedAt).Seconds()
		limit := 90.0
		switch next.MobilityProfile {
		case model.MobilityGroundRobot:
			limit = 20
		case model.MobilityAerialDrone:
			limit = 120
		case model.MobilityStatic:
			limit = 2
		}
		if implied > limit {
			q.Confidence = .25
			q.Anomalies = append(q.Anomalies, "IMPLAUSIBLE_JUMP")
		}
		if implied < .15 {
			q.Anomalies = append(q.Anomalies, "STATIONARY_DEADBAND")
		}
	}
	return q
}
```

---

## internal\modules\mobility\spatial\index.go

```
package spatial

import (
	"sort"

	"github.com/Akashpg-M/polaris/backend/internal/modules/mobility/model"
)

type SpatialCandidate struct {
	State          model.SpatialState `json:"state"`
	DistanceMeters float64            `json:"distance_meters"`
}

type SpatialIndex interface {
	Upsert(model.SpatialState) error
	Remove(tenantID, deviceID string) error
	Nearest(target model.Position, limit int, radiusMeters float64) ([]SpatialCandidate, error)
	WithinRadius(target model.Position, radiusMeters float64) ([]SpatialCandidate, error)
}

type LinearSpatialIndex struct{ states map[string]model.SpatialState }

func NewLinearSpatialIndex() *LinearSpatialIndex {
	return &LinearSpatialIndex{states: map[string]model.SpatialState{}}
}
func indexKey(t, d string) string { return t + "\x00" + d }
func (i *LinearSpatialIndex) Upsert(s model.SpatialState) error {
	if err := ValidatePosition(s.Position); err != nil {
		return err
	}
	i.states[indexKey(s.TenantID, s.DeviceID)] = s
	return nil
}
func (i *LinearSpatialIndex) Remove(t, d string) error { delete(i.states, indexKey(t, d)); return nil }
func (i *LinearSpatialIndex) WithinRadius(p model.Position, radius float64) ([]SpatialCandidate, error) {
	if err := ValidatePosition(p); err != nil {
		return nil, err
	}
	out := []SpatialCandidate{}
	for _, s := range i.states {
		d := DistanceMeters(p, s.Position)
		if radius <= 0 || d <= radius {
			out = append(out, SpatialCandidate{s, d})
		}
	}
	sortCandidates(out)
	return out, nil
}
func (i *LinearSpatialIndex) Nearest(p model.Position, limit int, radius float64) ([]SpatialCandidate, error) {
	out, e := i.WithinRadius(p, radius)
	if limit > 0 && len(out) > limit {
		out = out[:limit]
	}
	return out, e
}
func sortCandidates(v []SpatialCandidate) {
	sort.Slice(v, func(a, b int) bool {
		if v[a].DistanceMeters != v[b].DistanceMeters {
			return v[a].DistanceMeters < v[b].DistanceMeters
		}
		return v[a].State.DeviceID < v[b].State.DeviceID
	})
}
```

---

## internal\modules\mobility\spatial\manager.go

```
package spatial

import (
	"errors"
	"math"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/Akashpg-M/polaris/backend/internal/modules/mobility/model"
	h3 "github.com/uber/h3-go/v4"
)

var (
	ErrStaleVersion   = errors.New("spatial source version is not newer")
	ErrTenantCapacity = errors.New("tenant spatial capacity reached")
)

type ManagerConfig struct {
	H3Resolution, ShardResolution int
	MinMoveMeters                 float64
	MaxIndexAge                   time.Duration
	MaxH3Rings                    int
	MaxRadiusMeters               float64
	MaxDevicesPerTenant           int
}
type location struct {
	region uint64
	state  model.SpatialState
}
type MobilityShard struct {
	mu       sync.RWMutex
	RegionID uint64
	Index    SpatialIndex
	devices  map[string]model.SpatialState
}

type Manager struct {
	cfg         ManagerConfig
	shardsMu    sync.RWMutex
	shards      map[string]map[uint64]*MobilityShard
	locationsMu sync.RWMutex
	locations   map[string]location
	deviceLocks [64]sync.Mutex
}

func NewManager(c ManagerConfig) *Manager {
	return &Manager{cfg: c, shards: map[string]map[uint64]*MobilityShard{}, locations: map[string]location{}}
}
func hashDevice(v string) uint {
	var h uint = 2166136261
	for i := 0; i < len(v); i++ {
		h = (h ^ uint(v[i])) * 16777619
	}
	return h
}
func (m *Manager) cell(p model.Position) (uint64, uint64, error) {
	c, e := h3.LatLngToCell(h3.NewLatLng(p.Latitude, p.Longitude), m.cfg.H3Resolution)
	if e != nil {
		return 0, 0, e
	}
	parent, e := c.Parent(m.cfg.ShardResolution)
	return uint64(c), uint64(parent), e
}
func (m *Manager) getShard(tenant string, region uint64, create bool) *MobilityShard {
	m.shardsMu.RLock()
	s := m.shards[tenant][region]
	m.shardsMu.RUnlock()
	if s != nil || !create {
		return s
	}
	m.shardsMu.Lock()
	defer m.shardsMu.Unlock()
	if m.shards[tenant] == nil {
		m.shards[tenant] = map[uint64]*MobilityShard{}
	}
	if m.shards[tenant][region] == nil {
		m.shards[tenant][region] = &MobilityShard{RegionID: region, Index: NewRTreeSpatialIndex(), devices: map[string]model.SpatialState{}}
	}
	return m.shards[tenant][region]
}

func (m *Manager) Upsert(in model.SpatialState) error {
	if e := ValidatePosition(in.ReportedPosition); e != nil {
		return e
	}
	cell, region, e := m.cell(in.ReportedPosition)
	if e != nil {
		return e
	}
	in.H3Cell = cell
	lock := &m.deviceLocks[hashDevice(in.TenantID+"\x00"+in.DeviceID)%uint(len(m.deviceLocks))]
	lock.Lock()
	defer lock.Unlock()
	k := indexKey(in.TenantID, in.DeviceID)
	m.locationsMu.RLock()
	old, exists := m.locations[k]
	m.locationsMu.RUnlock()
	if exists && !in.NewerThan(old.state) {
		return ErrStaleVersion
	}
	in.Quality = ValidateObservation(func() *model.SpatialState {
		if exists {
			return &old.state
		}
		return nil
	}(), in)
	if in.HeadingDegrees != nil {
		normalized := math.Mod(*in.HeadingDegrees, 360)
		if normalized < 0 {
			normalized += 360
		}
		in.HeadingDegrees = &normalized
	}
	if !in.Quality.Valid {
		return errors.New("invalid mobility observation")
	}
	now := time.Now().UTC()
	if in.IndexedAt.IsZero() {
		in.IndexedAt = now
	}
	in.Position = in.ReportedPosition
	if exists && old.region == region && old.state.H3Cell == cell && DistanceMeters(old.state.Position, in.ReportedPosition) < m.cfg.MinMoveMeters && now.Sub(old.state.IndexedAt) < m.cfg.MaxIndexAge {
		in.Position = old.state.Position
		in.IndexedAt = old.state.IndexedAt
	}
	if !exists {
		m.locationsMu.RLock()
		count := 0
		for key := range m.locations {
			if len(key) > len(in.TenantID) && key[:len(in.TenantID)] == in.TenantID && key[len(in.TenantID)] == '\x00' {
				count++
			}
		}
		m.locationsMu.RUnlock()
		if count >= m.cfg.MaxDevicesPerTenant {
			return ErrTenantCapacity
		}
	}
	newShard := m.getShard(in.TenantID, region, true)
	oldShard := newShard
	if exists && old.region != region {
		oldShard = m.getShard(in.TenantID, old.region, false)
	}
	locked := []*MobilityShard{}
	if oldShard != nil && oldShard != newShard {
		if oldShard.RegionID < newShard.RegionID {
			locked = []*MobilityShard{oldShard, newShard}
		} else {
			locked = []*MobilityShard{newShard, oldShard}
		}
	} else {
		locked = []*MobilityShard{newShard}
	}
	for _, s := range locked {
		s.mu.Lock()
	}
	defer func() {
		for i := len(locked) - 1; i >= 0; i-- {
			locked[i].mu.Unlock()
		}
	}()
	if exists && oldShard != nil && oldShard != newShard {
		_ = oldShard.Index.Remove(in.TenantID, in.DeviceID)
		delete(oldShard.devices, in.DeviceID)
	}
	if e = newShard.Index.Upsert(in); e != nil {
		return e
	}
	newShard.devices[in.DeviceID] = in
	m.locationsMu.Lock()
	m.locations[k] = location{region, in}
	m.locationsMu.Unlock()
	return nil
}

func (m *Manager) Remove(tenant, device string) error {
	k := indexKey(tenant, device)
	lock := &m.deviceLocks[hashDevice(k)%uint(len(m.deviceLocks))]
	lock.Lock()
	defer lock.Unlock()
	m.locationsMu.Lock()
	loc, ok := m.locations[k]
	if ok {
		delete(m.locations, k)
	}
	m.locationsMu.Unlock()
	if !ok {
		return nil
	}
	s := m.getShard(tenant, loc.region, false)
	if s == nil {
		return nil
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.devices, device)
	return s.Index.Remove(tenant, device)
}
func (m *Manager) RemoveTenant(tenant string) error {
	m.locationsMu.RLock()
	devices := []string{}
	for key := range m.locations {
		if strings.HasPrefix(key, tenant+"\x00") {
			devices = append(devices, strings.TrimPrefix(key, tenant+"\x00"))
		}
	}
	m.locationsMu.RUnlock()
	sort.Strings(devices)
	for _, device := range devices {
		if err := m.Remove(tenant, device); err != nil {
			return err
		}
	}
	return nil
}
func (m *Manager) Get(tenant, device string) (model.SpatialState, bool) {
	m.locationsMu.RLock()
	defer m.locationsMu.RUnlock()
	v, ok := m.locations[indexKey(tenant, device)]
	return v.state, ok
}

func (m *Manager) Nearby(tenant string, p model.Position, radius float64, limit int) ([]SpatialCandidate, error) {
	if e := ValidatePosition(p); e != nil {
		return nil, e
	}
	if radius <= 0 || radius > m.cfg.MaxRadiusMeters {
		return nil, errors.New("search radius outside configured bounds")
	}
	cell, _, e := m.cell(p)
	if e != nil {
		return nil, e
	}
	rings := int(math.Ceil(radius/700)) + 1
	if rings > m.cfg.MaxH3Rings {
		rings = m.cfg.MaxH3Rings
	}
	cells, e := h3.Cell(cell).GridDisk(rings)
	if e != nil {
		return nil, e
	}
	regions := map[uint64]struct{}{}
	for _, c := range cells {
		parent, pe := c.Parent(m.cfg.ShardResolution)
		if pe == nil {
			regions[uint64(parent)] = struct{}{}
		}
	}
	found := map[string]SpatialCandidate{}
	for region := range regions {
		s := m.getShard(tenant, region, false)
		if s == nil {
			continue
		}
		s.mu.RLock()
		v, qerr := s.Index.WithinRadius(p, radius)
		s.mu.RUnlock()
		if qerr != nil {
			return nil, qerr
		}
		for _, c := range v {
			if old, ok := found[c.State.DeviceID]; !ok || c.DistanceMeters < old.DistanceMeters {
				found[c.State.DeviceID] = c
			}
		}
	}
	out := make([]SpatialCandidate, 0, len(found))
	for _, v := range found {
		out = append(out, v)
	}
	sortCandidates(out)
	if limit > 0 && len(out) > limit {
		out = out[:limit]
	}
	return out, nil
}

func (m *Manager) EvictInactive(tenant, device, lifecycle, connectivity string) error {
	if lifecycle != "ACTIVE" || connectivity == "STALE" || connectivity == "OFFLINE" {
		return m.Remove(tenant, device)
	}
	return nil
}
func (m *Manager) Rebuild(states []model.SpatialState) error {
	sort.Slice(states, func(i, j int) bool {
		if states[i].TenantID != states[j].TenantID {
			return states[i].TenantID < states[j].TenantID
		}
		return states[i].DeviceID < states[j].DeviceID
	})
	for _, s := range states {
		if e := m.Upsert(s); e != nil && !errors.Is(e, ErrStaleVersion) {
			return e
		}
	}
	return nil
}
func (m *Manager) Snapshot() []model.SpatialState {
	m.locationsMu.RLock()
	defer m.locationsMu.RUnlock()
	out := make([]model.SpatialState, 0, len(m.locations))
	for _, v := range m.locations {
		out = append(out, v.state)
	}
	sort.Slice(out, func(i, j int) bool {
		if out[i].TenantID != out[j].TenantID {
			return out[i].TenantID < out[j].TenantID
		}
		return out[i].DeviceID < out[j].DeviceID
	})
	return out
}
```

---

## internal\modules\mobility\spatial\rtree.go

```
package spatial

import (
	"math"
	"sort"
	"sync"
	"sync/atomic"
	"time"

	"github.com/Akashpg-M/polaris/backend/internal/modules/mobility/model"
)

const rtreeFanout = 16

type box struct{ minX, minY, maxX, maxY float64 }

func pointBox(p model.Position) box { return box{p.Longitude, p.Latitude, p.Longitude, p.Latitude} }
func (a box) intersects(b box) bool {
	return a.minX <= b.maxX && a.maxX >= b.minX && a.minY <= b.maxY && a.maxY >= b.minY
}
func union(a, b box) box {
	return box{math.Min(a.minX, b.minX), math.Min(a.minY, b.minY), math.Max(a.maxX, b.maxX), math.Max(a.maxY, b.maxY)}
}

type rtreeEntry struct {
	state  model.SpatialState
	bounds box
}
type rtreeNode struct {
	bounds   box
	leaf     bool
	entries  []rtreeEntry
	children []*rtreeNode
}

// RTreeSpatialIndex is a packed STR R-tree. Mutations are O(1) and mark the
// hierarchy dirty; the next read atomically rebuilds it. This suits telemetry
// bursts and gives exact queries after geodesic post-filtering.
type RTreeSpatialIndex struct {
	mu               sync.RWMutex
	items            map[string]model.SpatialState
	root             *rtreeNode
	dirty            bool
	mutations        atomic.Uint64
	rebuilds         atomic.Uint64
	rebuildNanos     atomic.Uint64
	rebuildWaitNanos atomic.Uint64
}

type RTreeStats struct {
	Mutations        uint64        `json:"mutations"`
	Rebuilds         uint64        `json:"rebuilds"`
	TotalRebuildTime time.Duration `json:"total_rebuild_time"`
	TotalLockWait    time.Duration `json:"total_lock_wait"`
}

func NewRTreeSpatialIndex() *RTreeSpatialIndex {
	return &RTreeSpatialIndex{items: map[string]model.SpatialState{}, dirty: true}
}
func (i *RTreeSpatialIndex) Upsert(s model.SpatialState) error {
	if e := ValidatePosition(s.Position); e != nil {
		return e
	}
	i.mu.Lock()
	defer i.mu.Unlock()
	i.items[indexKey(s.TenantID, s.DeviceID)] = s
	i.dirty = true
	i.mutations.Add(1)
	return nil
}
func (i *RTreeSpatialIndex) Remove(t, d string) error {
	i.mu.Lock()
	defer i.mu.Unlock()
	delete(i.items, indexKey(t, d))
	i.dirty = true
	i.mutations.Add(1)
	return nil
}

func (i *RTreeSpatialIndex) ensure() {
	i.mu.RLock()
	dirty := i.dirty
	i.mu.RUnlock()
	if !dirty {
		return
	}
	waitStarted := time.Now()
	i.mu.Lock()
	i.rebuildWaitNanos.Add(uint64(time.Since(waitStarted)))
	defer i.mu.Unlock()
	if !i.dirty {
		return
	}
	rebuildStarted := time.Now()
	entries := make([]rtreeEntry, 0, len(i.items))
	for _, s := range i.items {
		entries = append(entries, rtreeEntry{s, pointBox(s.Position)})
	}
	i.root = buildSTR(entries)
	i.dirty = false
	i.rebuilds.Add(1)
	i.rebuildNanos.Add(uint64(time.Since(rebuildStarted)))
}

func (i *RTreeSpatialIndex) Stats() RTreeStats {
	return RTreeStats{Mutations: i.mutations.Load(), Rebuilds: i.rebuilds.Load(), TotalRebuildTime: time.Duration(i.rebuildNanos.Load()), TotalLockWait: time.Duration(i.rebuildWaitNanos.Load())}
}

func buildSTR(entries []rtreeEntry) *rtreeNode {
	if len(entries) == 0 {
		return nil
	}
	sort.Slice(entries, func(a, b int) bool { return entries[a].state.Position.Longitude < entries[b].state.Position.Longitude })
	leaves := make([]*rtreeNode, 0, (len(entries)+rtreeFanout-1)/rtreeFanout)
	for start := 0; start < len(entries); start += rtreeFanout {
		end := start + rtreeFanout
		if end > len(entries) {
			end = len(entries)
		}
		part := append([]rtreeEntry(nil), entries[start:end]...)
		sort.Slice(part, func(a, b int) bool { return part[a].state.Position.Latitude < part[b].state.Position.Latitude })
		n := &rtreeNode{leaf: true, entries: part, bounds: part[0].bounds}
		for _, e := range part[1:] {
			n.bounds = union(n.bounds, e.bounds)
		}
		leaves = append(leaves, n)
	}
	return buildParents(leaves)
}
func buildParents(nodes []*rtreeNode) *rtreeNode {
	if len(nodes) == 1 {
		return nodes[0]
	}
	sort.Slice(nodes, func(a, b int) bool { return nodes[a].bounds.minX < nodes[b].bounds.minX })
	next := make([]*rtreeNode, 0, (len(nodes)+rtreeFanout-1)/rtreeFanout)
	for s := 0; s < len(nodes); s += rtreeFanout {
		e := s + rtreeFanout
		if e > len(nodes) {
			e = len(nodes)
		}
		n := &rtreeNode{children: append([]*rtreeNode(nil), nodes[s:e]...), bounds: nodes[s].bounds}
		for _, c := range nodes[s+1 : e] {
			n.bounds = union(n.bounds, c.bounds)
		}
		next = append(next, n)
	}
	return buildParents(next)
}
func searchNode(n *rtreeNode, q box, out map[string]model.SpatialState) {
	if n == nil || !n.bounds.intersects(q) {
		return
	}
	if n.leaf {
		for _, e := range n.entries {
			if e.bounds.intersects(q) {
				out[indexKey(e.state.TenantID, e.state.DeviceID)] = e.state
			}
		}
		return
	}
	for _, c := range n.children {
		searchNode(c, q, out)
	}
}

func queryBoxes(p model.Position, radius float64) []box {
	if radius <= 0 || radius >= math.Pi*EarthRadiusMeters {
		return []box{{-180, -90, 180, 90}}
	}
	latDelta := radius / EarthRadiusMeters * 180 / math.Pi
	cos := math.Cos(p.Latitude * math.Pi / 180)
	lonDelta := 180.0
	if math.Abs(cos) > .000001 {
		lonDelta = math.Min(180, latDelta/math.Abs(cos))
	}
	minY, maxY := math.Max(-90, p.Latitude-latDelta), math.Min(90, p.Latitude+latDelta)
	minX, maxX := p.Longitude-lonDelta, p.Longitude+lonDelta
	if minX < -180 {
		return []box{{-180, minY, maxX, maxY}, {minX + 360, minY, 180, maxY}}
	}
	if maxX > 180 {
		return []box{{minX, minY, 180, maxY}, {-180, minY, maxX - 360, maxY}}
	}
	return []box{{minX, minY, maxX, maxY}}
}
func (i *RTreeSpatialIndex) WithinRadius(p model.Position, radius float64) ([]SpatialCandidate, error) {
	if e := ValidatePosition(p); e != nil {
		return nil, e
	}
	i.ensure()
	i.mu.RLock()
	defer i.mu.RUnlock()
	found := map[string]model.SpatialState{}
	for _, q := range queryBoxes(p, radius) {
		searchNode(i.root, q, found)
	}
	out := make([]SpatialCandidate, 0, len(found))
	for _, s := range found {
		d := DistanceMeters(p, s.Position)
		if radius <= 0 || d <= radius {
			out = append(out, SpatialCandidate{s, d})
		}
	}
	sortCandidates(out)
	return out, nil
}
func (i *RTreeSpatialIndex) Nearest(p model.Position, limit int, radius float64) ([]SpatialCandidate, error) {
	out, e := i.WithinRadius(p, radius)
	if limit > 0 && len(out) > limit {
		out = out[:limit]
	}
	return out, e
}
```

---

## internal\modules\mobility\spatial\rtree_contention_test.go

```
package spatial

import (
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/Akashpg-M/polaris/backend/internal/modules/mobility/model"
)

func populatedRTree(size int) (*RTreeSpatialIndex, []model.SpatialState) {
	rng := rand.New(rand.NewSource(41))
	index := NewRTreeSpatialIndex()
	states := make([]model.SpatialState, size)
	for n := range states {
		states[n] = model.SpatialState{TenantID: "tenant", DeviceID: fmt.Sprintf("device-%06d", n), Position: model.Position{Latitude: 12.9 + rng.Float64()*.25, Longitude: 80.1 + rng.Float64()*.25}}
		_ = index.Upsert(states[n])
	}
	return index, states
}

func TestRTreeDirtyReadRebuildIsBoundedAndExact(t *testing.T) {
	index, states := populatedRTree(5000)
	started := time.Now()
	result, err := index.WithinRadius(states[0].Position, 10)
	if err != nil || len(result) == 0 || result[0].State.DeviceID != states[0].DeviceID {
		t.Fatalf("dirty rebuild query was not exact: count=%d err=%v", len(result), err)
	}
	stats := index.Stats()
	if stats.Rebuilds != 1 || stats.Mutations != 5000 {
		t.Fatalf("unexpected rebuild accounting: %#v", stats)
	}
	if elapsed := time.Since(started); elapsed > 500*time.Millisecond {
		t.Fatalf("5k dirty-tree rebuild exceeded local safety budget: %v", elapsed)
	}
	_, _ = index.WithinRadius(states[1].Position, 10)
	if index.Stats().Rebuilds != 1 {
		t.Fatal("clean read rebuilt an unchanged tree")
	}
}

func BenchmarkRTreeMovementQueryMix(b *testing.B) {
	for _, size := range []int{1000, 5000, 10000} {
		for _, writes := range []int{80, 50, 20} {
			name := fmt.Sprintf("devices=%d/writes=%d", size, writes)
			b.Run(name, func(b *testing.B) {
				index, states := populatedRTree(size)
				_, _ = index.WithinRadius(states[0].Position, 5000)
				b.ResetTimer()
				for n := 0; n < b.N; n++ {
					state := states[n%len(states)]
					if n%100 < writes {
						state.Position.Latitude += float64(n%7) * .000001
						_ = index.Upsert(state)
					} else {
						_, _ = index.Nearest(state.Position, 10, 5000)
					}
				}
			})
		}
	}
}
```

---

## internal\modules\mobility\spatial\spatial_test.go

```
package spatial

import (
	"errors"
	"math/rand"
	"reflect"
	"sort"
	"testing"
	"time"

	"github.com/Akashpg-M/polaris/backend/internal/modules/mobility/model"
)

func state(tenant, id string, lat, lon float64, seq uint64) model.SpatialState {
	now := time.Unix(1_700_000_000+int64(seq), 0).UTC()
	return model.SpatialState{TenantID: tenant, DeviceID: id, Position: model.Position{Latitude: lat, Longitude: lon}, ReportedPosition: model.Position{Latitude: lat, Longitude: lon}, MobilityProfile: model.MobilityRoadVehicle, ObservedAt: now, BootID: "boot-1", BootStartedAt: time.Unix(1_699_999_000, 0).UTC(), SequenceNumber: seq}
}
func ids(v []SpatialCandidate) []string {
	out := make([]string, len(v))
	for i, c := range v {
		out[i] = c.State.DeviceID
	}
	return out
}
func TestGeodesicEdgeCases(t *testing.T) {
	cases := []struct {
		a, b     model.Position
		min, max float64
	}{{model.Position{Latitude: 0, Longitude: 179.9}, model.Position{Latitude: 0, Longitude: -179.9}, 20_000, 25_000}, {model.Position{Latitude: 89, Longitude: 0}, model.Position{Latitude: 89, Longitude: 90}, 150_000, 160_000}, {model.Position{Latitude: 13, Longitude: 80}, model.Position{Latitude: 13, Longitude: 80}, 0, .001}}
	for _, tc := range cases {
		got := DistanceMeters(tc.a, tc.b)
		if got < tc.min || got > tc.max {
			t.Fatalf("distance %f outside [%f,%f]", got, tc.min, tc.max)
		}
	}
}
func TestRTreeMatchesLinearOracleRandomized(t *testing.T) {
	rng := rand.New(rand.NewSource(42))
	tree := NewRTreeSpatialIndex()
	linear := NewLinearSpatialIndex()
	for n := 0; n < 1000; n++ {
		s := state("tenant-a", string(rune(1000+n)), rng.Float64()*170-85, rng.Float64()*360-180, 1)
		if err := tree.Upsert(s); err != nil {
			t.Fatal(err)
		}
		_ = linear.Upsert(s)
	}
	queries := []model.Position{{Latitude: 13, Longitude: 80}, {Latitude: 0, Longitude: 179.95}, {Latitude: 75, Longitude: -40}}
	for _, q := range queries {
		a, _ := tree.Nearest(q, 10, 0)
		b, _ := linear.Nearest(q, 10, 0)
		if !reflect.DeepEqual(ids(a), ids(b)) {
			t.Fatalf("nearest mismatch %v != %v", ids(a), ids(b))
		}
	}
}
func TestRTreeMutationAndAntimeridian(t *testing.T) {
	tree := NewRTreeSpatialIndex()
	_ = tree.Upsert(state("t", "west", 0, 179.95, 1))
	_ = tree.Upsert(state("t", "east", 0, -179.95, 1))
	got, err := tree.WithinRadius(model.Position{Latitude: 0, Longitude: 180}, 20_000)
	if err != nil || len(got) != 2 {
		t.Fatalf("antimeridian search got=%v err=%v", ids(got), err)
	}
	moved := state("t", "west", 20, 20, 2)
	_ = tree.Upsert(moved)
	_ = tree.Remove("t", "east")
	got, _ = tree.WithinRadius(model.Position{Latitude: 0, Longitude: 180}, 20_000)
	if len(got) != 0 {
		t.Fatalf("removed/moved device returned: %v", ids(got))
	}
}
func TestManagerVersionTenantAndReplayInvariants(t *testing.T) {
	cfg := ManagerConfig{H3Resolution: 8, ShardResolution: 6, MinMoveMeters: 5, MaxIndexAge: time.Minute, MaxH3Rings: 12, MaxRadiusMeters: 10_000, MaxDevicesPerTenant: 100}
	m := NewManager(cfg)
	s := state("a", "d1", 13.0067, 80.2206, 2)
	if err := m.Upsert(s); err != nil {
		t.Fatal(err)
	}
	old := state("a", "d1", 14, 81, 1)
	if err := m.Upsert(old); !errors.Is(err, ErrStaleVersion) {
		t.Fatalf("expected stale version, got %v", err)
	}
	_ = m.Upsert(state("b", "d2", 13.0067, 80.2206, 1))
	got, _ := m.Nearby("a", model.Position{Latitude: 13.0067, Longitude: 80.2206}, 1000, 10)
	if !reflect.DeepEqual(ids(got), []string{"d1"}) {
		t.Fatalf("tenant leak: %v", ids(got))
	}
	first := m.Snapshot()
	replay := NewManager(cfg)
	shuffled := append([]model.SpatialState(nil), first...)
	sort.Slice(shuffled, func(i, j int) bool { return shuffled[i].DeviceID > shuffled[j].DeviceID })
	if err := replay.Rebuild(shuffled); err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(first, replay.Snapshot()) {
		t.Fatal("replay was not deterministic")
	}
	_ = m.EvictInactive("a", "d1", "ACTIVE", "OFFLINE")
	got, _ = m.Nearby("a", model.Position{Latitude: 13.0067, Longitude: 80.2206}, 1000, 10)
	if len(got) != 0 {
		t.Fatal("offline device remained indexed")
	}
}
func TestMovementThresholdPreservesReportedState(t *testing.T) {
	m := NewManager(ManagerConfig{H3Resolution: 8, ShardResolution: 6, MinMoveMeters: 100, MaxIndexAge: time.Hour, MaxH3Rings: 12, MaxRadiusMeters: 10_000, MaxDevicesPerTenant: 10})
	a := state("t", "d", 13, 80, 1)
	_ = m.Upsert(a)
	b := state("t", "d", 13.00001, 80.00001, 2)
	_ = m.Upsert(b)
	got, _ := m.Get("t", "d")
	if got.ReportedPosition == got.Position {
		t.Fatal("reported and indexed positions should be independent below threshold")
	}
	if got.SequenceNumber != 2 {
		t.Fatal("latest source version not retained")
	}
}

func BenchmarkRTreeNearest(b *testing.B)  { benchmarkNearest(b, NewRTreeSpatialIndex()) }
func BenchmarkLinearNearest(b *testing.B) { benchmarkNearest(b, NewLinearSpatialIndex()) }
func benchmarkNearest(b *testing.B, index SpatialIndex) {
	for i := 0; i < 5000; i++ {
		_ = index.Upsert(state("bench", string(rune(i+1)), 13+float64(i%100)*.001, 80+float64(i/100)*.001, 1))
	}
	q := model.Position{Latitude: 13.04, Longitude: 80.02}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = index.Nearest(q, 10, 5000)
	}
}
```

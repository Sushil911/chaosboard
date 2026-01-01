# ChaosBoard: Mini Chaos Engineering Toolkit in Go

My very first attempt at building a chaos engineering tool in Go.  
Inspired by Litmus Chaos, but right now it’s super tiny and I’m still learning everything.

## What actually works today
- A basic Go HTTP server running on :8080
- Nothing is deployed to Kubernetes yet
- No CLI flags, no real experiments running in the cluster still figuring that out

## How I’m testing it right now
```bash
go run main.go
```

Then in another terminal:
```bash
curl -X POST http://localhost:8080/api/experiments \
     -H "Content-Type: application/json" \
     -d '{"type":"cpu-hog","duration":15}'
```

## Why I’m doing this?
- I want to master Go + Kubernetes + DevOps the hard way.
- I started with Next.js and tRPC and got completely lost in abstractions.
- Now I’m going back to basics so I never get confused again.This repo (and my other one https://github.com/Sushil911/go-devops-mastery) is me building in public; every mistake, every tiny step.

## What I plan to add (when I figure it out)
- Actually create CPU-hog pods in Kubernetes
- Add more experiments
- Maybe a proper CLI later
- Eventually apply for LFX Mentorship with this

PRs, advice, or even “this is wrong because of this this this reason” comments are welcome. 

Sushil - December 2025

## Update: 23 Days Progress (December 8, 2025 – December 31, 2025)

December 2025 has ended, and it's been 23 days since I started learning DevOps + Backend in Go seriously. I've studied on around 18 days (78% consistency), with an average of 4-5 hours per day (max 8 hours, min 2 hours). The last few days were especially intense, and I learned more in them than in the first 19 combined.

Total estimated hours: ~90–100. Progress was linear at first with some plateaus (due to skipped days), but it accelerated sharply in the last week as concepts started connecting.

## What actually works today
- A basic Go HTTP server running on :8080 with a JSON API
- Create chaos experiments via `POST /api/experiments` (types: `cpu-hog`, `memory-hog`, `disk-fill`)
- List experiments via `GET /api/experiments`
- Persistence with BoltDB (experiments survive restarts)
- Thread-safe in-memory cache with mutex
- Graceful shutdown with signal handling (SIGINT/SIGTERM)
- Prometheus metrics endpoint (`/metrics`) with:
  - Go runtime metrics (CPU, memory, goroutines, GC)
  - Custom HTTP metrics (requests total, request duration histogram)
  - Custom chaos metrics (experiments total by type, experiments active gauge)
- Multi-stage Dockerfile:
  - `deps` stage for dependency caching
  - `development` stage with hot-reload (air + `.air.toml`), golangci-lint, delve debugger
  - `builder` stage for static binary
  - `production` stage using distroless nonroot (secure & tiny)
- Docker Compose for dev (hot-reload + volumes) and prod (no volumes, restart policy)
- Prometheus + Grafana monitoring stack (via Docker Compose)

You can run this using:
```bash
docker compose -f docker-compose.dev.yml up
```

Then visit localhost:8080 for backend app, localhost:3000 for Grafana and localhost:9090 for Prometheus

- Sushil, January 2026





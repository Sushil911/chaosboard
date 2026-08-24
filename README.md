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

## Current Implementation(Dec 31,2025)
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


## Update: Some Ups & Downs (January 1, 2026 – July 31, 2026)

Somewhere along the way I got lost and jumped to other different domains which are way different from this domain that I was learning. So unfortunately I can't document those things in this github profile. But I'll be consistent from this day onwards and give my best to learn.

---

## Update: The Kubernetes Chapter (July 31, 2026 – August 24, 2026)
After getting the Docker and Docker Compose, I've setup CI/CD using Github Actions then I moved to Kubernetes. I set up a local Kind cluster, wrote deployment YAMLs, and got ChaosBoard running inside the cluster. The rollout strategy now uses the `:latest` tag with `imagePullPolicy: Always` so updates are easy. I also learned how the Control Plane and Worker Nodes work, and how traffic flows from a browser to a pod.

### Current Implementation (Latest)
* Working Go API with graceful shutdown
* BoltDB with in-memory cache
* Multi-stage Dockerfile (dev, builder, production)
* Docker Compose for local dev (with Prometheus + Grafana)
* GitHub Actions CI/CD (lint, test, build, push to Docker Hub)
* Kind cluster running locally
* ChaosBoard deployed as a Kubernetes Deployment
* Service exposing ChaosBoard internally
* Rollout restart with `:latest` image tag
* Prometheus + Grafana stack running locally (not yet inside K8s)

---

## How to Run It Locally

### With Docker Compose
```bash
docker compose -f docker-compose.dev.yml up --build
```

### With Kubernetes
```bash
kind create cluster --config kind-config.yaml --name chaosboard
kubectl apply -f chaosboard-deployment.yaml
kubectl apply -f chaosboard-service.yaml
kubectl port-forward service/chaosboard-service 8080:8080
```

---

## Future Plans
* Production-grade chaos experiments (context-aware, cgroup-aware)
* Deploy Prometheus and Grafana inside the Kubernetes cluster
* Replace Prometheus client with OpenTelemetry
* Write a Kubernetes controller and CRD
* Contribute to Kind or LitmusChaos to learn more

- Sushil, August 2026






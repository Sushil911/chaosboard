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




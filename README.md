# Kronos Scheduler

Kronos is a modular, queue-aware, and power-conscious Kubernetes scheduler extender designed to intelligently prioritize pod placement. It uses multi-level feedback scheduling logic, grounded in queuing theory, to reduce job wait times and optimize cost and energy consumption across the cluster.

This project models pod workloads using M/M/1 and M/G/1 queues, and applies different node scoring heuristics based on predicted job behavior. Kronos is designed to be extended with machine learning (for runtime prediction) and real-time energy metrics (via Kepler).

---

## Features

- Scheduler extender for Kubernetes (non-invasive, runs alongside default scheduler)
- Queue-type-aware scoring:
  - Short jobs (M/M/1): low delay penalty
  - Long jobs (M/G/1): higher delay sensitivity
- Custom pod scoring logic based on current node queue
- Label-based job classification (jobType: short / long)
- Designed for integration with:
  - ML-based job runtime predictors
  - Kepler for node-level energy usage
- Theoretical modeling using queuing theory (Wq, ρ, Var(S))

---

## Architecture

1. A pod is submitted to the cluster.
2. Kube-scheduler invokes Kronos via the /prioritize endpoint.
3. Kronos inspects the pod's jobType label and estimates expected wait cost.
4. Node scores are calculated using queue-aware logic.
5. The kube-scheduler schedules the pod on the highest scoring node.

```
Pod → kube-scheduler → /prioritize (Kronos) → score(node) → selected node → scheduled
```

Queue-based scoring strategy:

| Job Type | Queue Model | Penalty Weight | Reasoning |
|----------|-------------|----------------|-----------|
| short    | M/M/1       | Low (e.g. 5x)  | Short jobs are less sensitive to queueing delays |
| long     | M/G/1       | High (e.g. 15x) | High variance jobs suffer more from delay buildup |

---

## Tech Stack

| Component        | Tool                         |
|------------------|------------------------------|
| Scheduler Logic  | Go 1.21+ (Gin/Gorilla Mux)   |
| Kubernetes       | Minikube / kind              |
| Monitoring (optional) | Prometheus, Kepler           |
| ML Prediction (future) | ONNX / External API |

---

## Getting Started

1. Clone the repo:
   ```bash
   git clone https://github.com/priyxansh/kronos-scheduler.git
   cd kronos-scheduler
   ```

2. Install dependencies:
   ```bash
   go mod tidy
   ```

3. Build and run the extender:
   ```bash
   go build -o kronos-scheduler ./src
   ./kronos-scheduler
   ```

4. Patch your kube-scheduler to use the scheduler-config.yaml in this repo.

---

## Example: Pod YAML

```yaml
apiVersion: v1
kind: Pod
metadata:
  name: short-job
  labels:
    jobType: "short"
spec:
  containers:
    - name: worker
      image: busybox
      command: ["sh", "-c", "sleep 10"]
```

---

## License

MIT License

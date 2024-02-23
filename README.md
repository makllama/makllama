# kcd-shanghai-2024

## Background

With over 100,000 Docker Pulls, [Ollama](https://github.com/ollama/ollama) has become a popular choice for building and running language models on the local machine.

## Problem

However, Ollama is not designed to run on Kubernetes. This is a problem for many organizations that want to run Ollama at scale.

## Solution

We will use Virtual Kubelet + Containerd + shimv2 + runm to run Ollama on Kubernetes.

### Demo

Running Ollama on Kubernetes with macOS nodes:

```bash
```
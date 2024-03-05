# kcd-shanghai-2024

Alternative titles:

- Running LLMs on Kubernetes with Virtual Kubelet + Containerd + ShimV2 + runm
- Beyond container, orchestrate LLMs with Kubernetes on macOS

## Description

With the growing popularity of generative AI, there is an increasing demand for large language models (LLMs)
inference capabilities. Kubernetes, being the most popular orchestration platform, is a natural fit for these
inference needs. Although GPUs are expensive and often in short supply, Apple Silicon M-series chips
(with unified memory) have been proven to be an effective alternative for running LLMs
(see ggerganov/llama.cpp performance discussion). However, the prevalent Kubernetes ecosystem is predominantly
focused on Linux-based containers. In this presentation, we will showcase our efforts to facilitate LLMs inference
on Kubernetes using macOS nodes. We will demonstrate how to employ Virtual Kubelet, Containerd, ShimV2, and runm
(derived from llama.cpp: ggerganov/llama.cpp) for deploying open-source foundation models such as gemma, llama2,
and mistral on Kubernetes. Additionally, we will discuss our motivation and the challenges encountered during our
development journey. Our goal is to encourage the community to expand the Kubernetes ecosystem to inclusively
support the execution of LLMs on macOS platforms.

## Benefits to the Ecosystem

- Enable running LLMs on Kubernetes with macOS nodes
- Provide an alternative solution for running LLMs on Kubernetes
- Inspire the community to build a more inclusive Kubernetes ecosystem that supports running LLMs on macOS

## Background

With over 100,000 Docker Pulls, [Ollama](https://github.com/ollama/ollama) has become a popular choice for building and running language models on the local machine.

## Problem

However, Ollama is not designed to run on Kubernetes. This is a problem for many organizations that want to run Ollama at scale.

## Solution

We will use Virtual Kubelet + Containerd + shimv2 + runm to run Ollama on Kubernetes.

### Demo

Running Ollama on Kubernetes with macOS nodes:

```bash
git clone https://github.com/yeahdongcn/containerd.git
cd containerd
git checkout llm
make
make bin/containerd-shim-runc-v2

cd ..

git clone https://github.com/yeahdongcn/ollama.git
cd ollama
git checkout runm
go generate ./...
go build -o runm .
cp ./runm ../containerd/bin/runm

cd ..

git clone https://github.com/yeahdongcn/cri.git
cd cri
make clean && make build
cp ./bin/virtual-kubelet ../containerd/bin/virtual-kubelet

cd ..

cd containerd/bin
# Open a new terminal
sudo ./containerd

# Open a new terminal
# Make sure you have a running Kubernetes cluster
sudo ./virtual-kubelet --kubeconfig ~/.kube/config
# Remove taints of the virtual-kubelet node

# Open a new terminal
cd containerd/test/ollama
kubectl apply -f tinyllama-pod-11111.yaml
kubectl apply -f tinyllama-pod-11112.yaml

# 2 tinyllama pods are running on the virtual-kubelet node
kubectl get pods
NAME              READY   STATUS    RESTARTS   AGE
tinyllama-11111   1/1     Running   0          41m
tinyllama-11112   1/1     Running   0          41m
```

## References

llama 2 7B data:

| GPU                      | FP16  |    Q8 |    Q4 | price |
| ------------------------ | :---: | ----: | ----: | ----: |
| RTX 4060 Ti 16GB         | 19.10 | 33.86 | 57.87 |  4599 |
| RTX 4000                 | 23.93 | 42.43 | 73.38 |  7499 |
| RTX 3090                 |       |       | 87.34 |  7899 |
| M1 Ultra 48 core         | 33.92 | 55.69 | 74.93 |       |
| M1 Ultra 64 core         | 37.01 | 59.87 | 83.73 |       |
| M2 Ultra 60 core 128G 1T | 39.86 | 62.14 | 88.64 | 38999 |
| M2 Ultra 76 core 128G 1T | 41.02 | 66.64 | 94.27 | 46499 |

https://github.com/ggerganov/llama.cpp/discussions/4167

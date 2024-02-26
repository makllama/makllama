# kcd-shanghai-2024

Alternative titles:

- Running LLMs on Kubernetes with Virtual Kubelet + Containerd + ShimV2 + runm
- Beyond container, orchestrate LLMs with Kubernetes on macOS

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
<div align="center">
 <img alt="logo" height="200px" src="https://github.com/makllama/makllama/assets/2831050/1ff40147-f065-449d-b273-40a69995d980"></br>
 <i>Powered by DALL·E 3</i>
</div>

# MaKllama
[![Go Report Card](https://goreportcard.com/badge/github.com/makllama/makllama)](https://goreportcard.com/report/github.com/makllama/makllama)
<a href="https://github.com/makllama/makllama/graphs/contributors" alt="Contributors"><img src="https://img.shields.io/github/contributors/makllama/makllama" /></a>
<img alt="GitHub last commit (branch)" src="https://img.shields.io/github/last-commit/makllama/makllama/main" />
<img alt="GitHub" src="https://img.shields.io/github/license/makllama/makllama" />

Running and orchestrating large language models (LLMs) on Kubernetes with macOS nodes.

## Table of Contents

- [MaKllama](#makllama)
  - [Table of Contents](#table-of-contents)
  - [Main Components](#main-components)
  - [Quick Start (~ 1 minute)](#quick-start--1-minute)
    - [1. Prerequisites](#1-prerequisites)
    - [2. Start Containerd + Virtual Kubelet + BW](#2-start-containerd--virtual-kubelet--bw)
    - [3. Deploy TinyLlama with 2 Replicas](#3-deploy-tinyllama-with-2-replicas)
    - [4. Deploy Mods](#4-deploy-mods)
    - [5. Access OpenAI API Compatible Endpoint through Mods](#5-access-openai-api-compatible-endpoint-through-mods)
    - [6. Stop Containerd + Virtual Kubelet + BW](#6-stop-containerd--virtual-kubelet--bw)
  - [Community](#community)
  - [Session Submissions](#session-submissions)
    - [Title](#title)
    - [Description](#description)
    - [Benefits to the Ecosystem](#benefits-to-the-ecosystem)

## Main Components

To run and orchestrate LLMs on Kubernetes with macOS nodes, we need the following components:

- [Virtual Kubelet](https://github.com/makllama/cri): For running `pods` on macOS nodes (forked from [virtual-kubelet/cri](https://github.com/virtual-kubelet/cri)).
- [Containerd](https://github.com/makllama/containerd): For pulling and running Ollama LLM image (forked from [containerd/containerd](https://github.com/containerd/containerd)).
- Runm: A lightweight runtime derived from [llama.cpp](https://github.com/ggerganov/llama.cpp) for running LLMs on macOS nodes (source code will be available soon).
- Bronze Willow: CNI Plugin for macOS (source code will be available soon).

This project is inspired by [llama.cpp](https://github.com/ggerganov/llama.cpp), [Ollama](https://github.com/ollama/ollama) and [kind](https://kind.sigs.k8s.io/).

## Quick Start (~ 1 minute)

### 1. Prerequisites

* A Kubernetes cluster.
  * [kind](https://kind.sigs.k8s.io/) is not supported.
  * [Antrea](https://github.com/antrea-io/antrea) is preferred for CNI.
  * `kubeconfig` should locate at `~/.kube/config`.
* Mac with Apple Silicon chip.

### 2. Start Containerd + Virtual Kubelet + BW

```bash
$ make # optional
$ sudo ./bin/demo create
 ✓ Starting containerd 🚢
 ✓ Preparing virtual nodes 📦
 ✓ Creating network 🌐
$ kubectl get nodes
NAME            STATUS     ROLES           AGE    VERSION
bj-k8s01        Ready      control-plane   214d   v1.28.2
bj-k8s02        Ready      worker          214d   v1.28.2
bj-k8s03        Ready      worker          214d   v1.28.2
weiqiangt-mba   Ready      agent           23d    v1.15.2-vk-cri-fb9cc09-dev
xiaodong-m1     Ready      agent           23d    v1.15.2-vk-cri-fb9cc09-dev
```

After running the above commands, you should see the macOS nodes appear in the output of `kubectl get nodes`. In the example above, `weiqiangt-mba` and `xiaodong-m1` are the macOS nodes.

### 3. Deploy TinyLlama with 2 Replicas

```bash
$ kubectl apply -f k8s/tinyllama.yml
```

### 4. Deploy Mods

```bash
$ kubectl apply -f k8s/mods.yaml
```

### 5. Access OpenAI API Compatible Endpoint through Mods

```bash
# Retrieve the command for editing config file of mods.
$ echo sed -i \'s/localhost:11434/$(kubectl get svc -o json tinyllama-services | jq '.spec.clusterIP' -r)/g\' '~/.config/mods/mods.yml'
sed -i 's/localhost:11434/198.19.50.27/g' ~/.config/mods/mods.yml
# Copy the output.
$ kubectl exec -it $(kubectl get pods -l app=mods -o jsonpath='{.items[0].metadata.name}') -- bash
root@mods-deployment-77c464f4b8-zn6g5:/# echo "Execute the copied command."
root@mods-deployment-77c464f4b8-zn6g5:/# mods -f "What are some of the best ways to save money?"
```

### 6. Stop Containerd + Virtual Kubelet + BW

```bash
$ sudo ./bin/demo delete
 ✓ Deleting demo 🧹
```

## Community

* [Open an issue](https://github.com/makllama/makllama/issues/new)

## Session Submissions

- KCD Shanghai 2024 (Accepted)
- KubeCon + CloudNativeCon + Open Source Summit + AI_Dev China 2024 (In Evaluation)

### Title

Beyond Containers, Orchestrate LLMs with Kubernetes on macOS

### Description

With the growing popularity of generative AI, there is an increasing demand for large language models (LLMs)
inference capabilities. Kubernetes, being the most popular orchestration platform, is a natural fit for these
inference needs. Although GPUs are expensive and often in short supply, Apple Silicon M-series chips
(with Unified Memory Architecture) have been proven to be an effective alternative for running LLMs
(see ggerganov/llama.cpp performance discussion). However, the prevalent Kubernetes ecosystem is predominantly
focused on Linux-based containers. In this presentation, we will showcase our efforts to facilitate LLMs inference
on Kubernetes using macOS nodes. We will demonstrate how to employ Virtual Kubelet, Containerd, ShimV2, and runm
(derived from llama.cpp: ggerganov/llama.cpp) for deploying open-source foundation models such as gemma, llama2,
and mistral on Kubernetes. Additionally, we will discuss our motivation and the challenges encountered during our
development journey. Our goal is to encourage the community to expand the Kubernetes ecosystem to inclusively
support the execution of LLMs on macOS platforms.

### Benefits to the Ecosystem

- Enable running and orchestrating LLMs on Kubernetes with macOS nodes
- Provide an alternative solution for running LLMs on Kubernetes
- Inspire the community to build a more inclusive Kubernetes ecosystem that supports running LLMs on macOS

# Multicluster-Ingress-Controller

Multicluster-Ingress-Controller is a Go library for building Kubernetes ingress controllers that need to watch resources in multiple clusters. Multicluster-Ingress-Controller is accesible via REST API Interface. It is capable of delivering same events to multiple recipients. It is also capable of delivering specific kubernetes events to the registered servers.  


## Table of Contents

- [Introduction](#introduction)
- [Architecture](#architecture)
- [How it Works](#how-it-works)
- [Getting Started](#getting-started)
- [Configuration](#configuration)
- [API Reference](#api-reference)

## Introduction

Hybrid and multicloud architectures are becoming predominant. Customers are moving thier workloads to hybrid and multi cluster. Hybrid and Multicluster are very complex in nature. Kubernetes is standardizing the way clouds are operated: the same workflow can be used to manage resources in any cloud, whether public or private or even in on-prem. However managing workloads, exposing the application to external world (Route) across all these multiple clusters is still a challenge. Every cluster should have ingress controller to do this now. Managing life cycle of all such ingress controllers will be another challege. There are several platform orchastrators which allows to manage all clusters from one place. There tools to handle multiclsuer schedule, autoscale, upgrade etc from one place. Multicluster Ingress COntroller helps the admin to manage all the cluster from one Ingress controller which runs at one place (Inside k8s cluster/outside) and listen for all the ingress events and creates the routes.

## Architecture

![Multicluster Architecture](pkg/docs/images/Multicluster.png)
       <details>
       <summary>**Rest Interface**</summary>
	    Rest Interface has two components. First one is external API and another is openAPI documentation. OpenAPI documentation allows the user to understand the usage of API and an options to try out the Rest API. External API module accepts the user API request and process it. Rest interface invokes Config store and controller.  
       </details>
       <details>
       <summary>**Controller**</summary>
	    Controller has two module, first one is go client and the next  is dispatcher. Kubernetes Go client is being used for getting the events from kubernetes cluster. Dispatcer sends out the filtered events to the list of registered servers.
       </details>
       <details>
       <summary>**Config Store**</summary>
	    This is a persitant volume store in kubernetes. This helps to  keeps all the rest input. 
       </details>
       <details>
       <summary>**Scheduler**</summary>
	    Scheduler restart the controller for which it takes the input from the Config store. This module required for brining up the control module which is the core part of the MIC.
       </details>

## How it Works

Multicluster-Ingress-Controller is a rest based go library which exposes few Rest API. User can register and access the events from multiple cluster which can be send to multiple server.  



## Getting Started

Run locally
```
	1) Download/Clone the Multicluster-Ingress-Controller.
		
	2) Perform "make run" from build folder.
		This starts the Multicluster Ingress Controller. Go binary has to be installed for running MIC.

	3) Access the Rest Interface via http://localhost:8000/swagger/
```
Creating Binary and Docker Image
```
	1) Download/Clone the Multicluster-Ingress-Controller.
		
	2) Perform "make build" from build folder.
		This builds a go binary and docker image and stored under build folder (citrix-k8s-multicluster-ingress-controller).

	3) Copy binary and run in any machine (./citrix-k8s-multicluster-ingress-controller)
	4) Use Docker commands to load and run as docker images.
```

## Configuration

There is no specific configuratiosn required at initial phase. However right now user has to keep all the cluster kube config file where the controller is running. We are working on token based solution which will remove this dependancy.

## API reference

Following are Sample Rest API exposed by Multicluster Ingress controller.

![Rest API List](pkg/docs/images/RestApi.png)






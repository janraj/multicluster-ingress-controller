# Multicluster-Ingress-Controller

Multicluster-Ingress-Controller is a Go library for building Kubernetes ingress controllers that need to watch resources in multiple clusters. Multicluster-Ingress-Controller is accesible via REST API Interface. It is capable of delivering same events to multiple recipients. It is also capable of delivering specific kubernetes events to the registered servers.  


## Table of Contents

- [How it Works](#how-it-works)
- [Getting Started](#getting-started)
- [Configuration](#configuration)
- [API Reference](#api-reference)

## How it Works

Multicluster-Ingress-Controller is a rest based go library which exposes few Rest API. User can register access the events from multiple cluster which can be send to multiple server.  



## Getting Started

```
	1) Download/Clone the Multicluster-Ingress-Controller.
		
	2) Perform "make run" from build folder.
		This starts the Multicluster Ingress Controller

	3) Access the Rest Interface via http://localhost:8000/swagger/
```
	
## Configuration

There is no specific configuratiosn required at initial phase. However right now user has to keep all the cluster kube config file where the controller is running. We are working on token based solution which will remove this dependancy.

## API reference

Following are Sample Rest API exposed by Multicluster Ingress controller.

![Rest API List](pkg/docs/images/RestApi.png)






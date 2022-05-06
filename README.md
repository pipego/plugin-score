# plugin-score

[![Build Status](https://github.com/pipego/plugin-score/workflows/ci/badge.svg?branch=main&event=push)](https://github.com/pipego/plugin-score/actions?query=workflow%3Aci)
[![Go Report Card](https://goreportcard.com/badge/github.com/pipego/plugin-score)](https://goreportcard.com/report/github.com/pipego/plugin-score)
[![License](https://img.shields.io/github/license/pipego/plugin-score.svg)](https://github.com/pipego/plugin-score/blob/main/LICENSE)
[![Tag](https://img.shields.io/github/tag/pipego/plugin-score.svg)](https://github.com/pipego/plugin-score/tags)



## Introduction

*plugin-score* is the score plugin of [pipego](https://github.com/pipego) written in Go.



## Prerequisites

- Go >= 1.17.0



## Run

```bash
make lint
make build
./plugin-score-test
```



## Docker



## Usage

- `plugin/noderesourcesbalancedallocation.go`: Favors nodes that would obtain a more balanced resource usage if the Task is scheduled there.
- `plugin/noderesourcesfit.go`: Checks if the node has all the resources that the Task is requesting.



## Settings



## License

Project License can be found [here](LICENSE).



## Reference

- [go-plugin](https://github.com/hashicorp/go-plugin)
- [kube-scheduler-config](https://kubernetes.io/docs/reference/scheduling/config)
- [kube-scheduler-interface](https://github.com/kubernetes/kubernetes/blob/master/pkg/scheduler/framework/interface.go)
- [kube-scheduler-plugins](https://github.com/kubernetes/kubernetes/blob/master/pkg/scheduler/framework/plugins)

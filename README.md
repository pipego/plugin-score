# plugin-score

[![License](https://img.shields.io/github/license/pipego/plugin-score.svg)](https://github.com/pipego/plugin-score/blob/main/LICENSE)



## Introduction

*plugin-score* is the score plugin of [pipego](https://github.com/pipego) written in Go.



## Prerequisites

- Go >= 1.17.0



## Run

```bash
# Template
go env -w GOPROXY=https://goproxy.cn,direct

cd template
golangci-lint run
CGO_ENABLED=0 GOARCH=amd64 GOOS=linux go build -ldflags "-s -w" -o plugin/plugin-score-template main.go
upx plugin/plugin-score-template
```



```bash
# Test
go env -w GOPROXY=https://goproxy.cn,direct

cd test
golangci-lint run
CGO_ENABLED=0 GOARCH=amd64 GOOS=linux go build -ldflags "-s -w" -o bin/test main.go
upx bin/test

cp ../template/plugin/* bin/
./bin/test
```



## Docker



## Usage



## Settings



## License

Project License can be found [here](LICENSE).



## Reference

- [go-plugin](https://github.com/hashicorp/go-plugin)
- [kube-scheduler-config](https://kubernetes.io/docs/reference/scheduling/config)
- [kube-scheduler-interface](https://github.com/kubernetes/kubernetes/blob/master/pkg/scheduler/framework/interface.go)
- [kube-scheduler-plugins](https://github.com/kubernetes/kubernetes/blob/master/pkg/scheduler/framework/plugins)

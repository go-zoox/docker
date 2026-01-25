# Docker - Simply Docker SDK

[![PkgGoDev](https://pkg.go.dev/badge/github.com/go-zoox/docker)](https://pkg.go.dev/github.com/go-zoox/docker)
[![Build Status](https://github.com/go-zoox/docker/actions/workflows/lint.yml/badge.svg?branch=master)](https://github.com/go-zoox/docker/actions/workflows/lint.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/go-zoox/docker)](https://goreportcard.com/report/github.com/go-zoox/docker)
[![Coverage Status](https://coveralls.io/repos/github/go-zoox/docker/badge.svg?branch=master)](https://coveralls.io/github/go-zoox/docker?branch=master)
[![GitHub issues](https://img.shields.io/github/issues/go-zoox/docker.svg)](https://github.com/go-zoox/docker/issues)
[![Release](https://img.shields.io/github/tag/go-zoox/docker.svg?label=Release)](https://github.com/go-zoox/docker/tags)

一个简洁易用的 Docker SDK，为 Go 语言提供 Docker API 的封装。支持容器、镜像、网络、卷和系统信息的管理操作。

## 特性

- 🚀 **简洁的 API** - 使用函数式选项模式，API 设计简洁直观
- 📦 **完整功能** - 支持容器、镜像、网络、卷等 Docker 核心功能
- 🔧 **灵活配置** - 支持本地和远程 Docker 主机
- 🛠️ **CLI 工具** - 内置命令行工具，方便快速操作
- 📚 **类型安全** - 完整的类型定义，提供良好的 IDE 支持

## 安装

```bash
go get -u github.com/go-zoox/docker
```

## 快速开始

### 创建客户端

```go
package main

import (
    "context"
    "fmt"
    "github.com/go-zoox/docker"
)

func main() {
    // 使用默认配置（从环境变量读取）
    client, err := docker.New()
    if err != nil {
        panic(err)
    }

    // 或者指定 Docker 服务器地址
    client, err = docker.New(func(cfg *docker.Config) {
        cfg.Server = "tcp://localhost:2375"
    })
    if err != nil {
        panic(err)
    }

    // 使用客户端...
}
```

## API 文档

### 容器管理 (Container)

容器管理提供了完整的容器生命周期操作。

#### 列出容器

```go
containers, err := client.Container().List(ctx)
if err != nil {
    return err
}

for _, container := range containers {
    fmt.Printf("ID: %s, Name: %s, Status: %s\n", 
        container.ID, container.Names, container.Status)
}
```

#### 创建容器

```go
import (
    "github.com/go-zoox/docker/container"
    dc "github.com/docker/docker/api/types/container"
)

response, err := client.Container().Create(ctx, func(opt *container.CreateOptions) {
    opt.Name = "my-container"
    opt.Container.Image = "nginx:latest"
    opt.Container.Cmd = []string{"nginx", "-g", "daemon off;"}
    opt.Container.Env = []string{"ENV=production"}
})
```

#### 启动/停止/重启容器

```go
// 启动容器
err := client.Container().Start(ctx, "container-id")

// 停止容器
err := client.Container().Stop(ctx, "container-id", func(opt *container.StopOptions) {
    opt.Timeout = 30 * time.Second
})

// 重启容器
err := client.Container().Restart(ctx, "container-id")
```

#### 查看容器日志

```go
logs, err := client.Container().Logs(ctx, "container-id", func(opt *container.LogsConfig) {
    opt.ShowStdout = true
    opt.ShowStderr = true
    opt.Follow = true
    opt.Tail = "100"
})
if err != nil {
    return err
}
defer logs.Close()

io.Copy(os.Stdout, logs)
```

#### 执行容器命令

```go
term, err := client.Container().Exec(ctx, "container-id", func(opt *container.ExecOptions) {
    opt.Cmd = []string{"ls", "-la"}
    opt.Tty = true
})
if err != nil {
    return err
}
defer term.Close()

io.Copy(os.Stdout, term)
```

#### 查看容器统计信息

```go
stats, err := client.Container().Stats(ctx, "container-id")
if err != nil {
    return err
}
defer stats.Close()

// 读取统计信息流
io.Copy(os.Stdout, stats)
```

#### 运行容器（创建并启动）

```go
err := client.Container().Run(ctx, func(opt *container.RunOptions) {
    opt.Image = "nginx:latest"
    opt.Name = "my-nginx"
    opt.Remove = true // 退出时自动删除
})
```

### 镜像管理 (Image)

镜像管理提供了镜像的构建、拉取、推送等操作。

#### 列出镜像

```go
images, err := client.Image().List(ctx)
if err != nil {
    return err
}

for _, image := range images {
    fmt.Printf("ID: %s, Tags: %v, Size: %d\n", 
        image.ID, image.RepoTags, image.Size)
}
```

#### 拉取镜像

```go
err := client.Image().Pull(ctx, "nginx:latest", func(cfg *image.PullConfig) {
    cfg.All = false
})
```

#### 推送镜像

```go
err := client.Image().Push(ctx, "my-registry/nginx:latest", func(cfg *image.PushConfig) {
    cfg.All = false
})
```

#### 构建镜像

```go
err := client.Image().Build(ctx, "./dockerfile-dir", func(cfg *image.BuildConfig) {
    cfg.Dockerfile = "Dockerfile"
    cfg.Tags = []string{"my-image:latest"}
})
```

#### 删除镜像

```go
results, err := client.Image().Remove(ctx, "image-id", func(cfg *image.RemoveConfig) {
    cfg.Force = false
    cfg.PruneChildren = true
})
```

#### 清理未使用的镜像

```go
import "github.com/docker/docker/api/types/filters"

report, err := client.Image().Prune(ctx, func(cfg *image.PruneConfig) {
    *cfg = filters.NewArgs()
    cfg.Add("dangling", "true")
})
fmt.Printf("Deleted: %d, Space Reclaimed: %d\n", 
    len(report.ImagesDeleted), report.SpaceReclaimed)
```

### 网络管理 (Network)

网络管理提供了 Docker 网络的创建、删除和查询功能。

#### 列出网络

```go
networks, err := client.Network().List(ctx)
if err != nil {
    return err
}

for _, network := range networks {
    fmt.Printf("ID: %s, Name: %s, Driver: %s\n", 
        network.ID, network.Name, network.Driver)
}
```

#### 创建网络

```go
import (
    "github.com/go-zoox/docker/network"
    dnetwork "github.com/docker/docker/api/types/network"
)

response, err := client.Network().Create(ctx, "my-network", func(opt *network.CreateOption) {
    opt.Driver = "bridge"
    opt.IPAM = &dnetwork.IPAM{
        Config: []dnetwork.IPAMConfig{
            {Subnet: "172.20.0.0/16"},
        },
    }
})
```

#### 删除网络

```go
err := client.Network().Remove(ctx, "network-id")
```

#### 清理未使用的网络

```go
report, err := client.Network().Prune(ctx)
fmt.Printf("Deleted: %d\n", len(report.NetworksDeleted))
```

### 卷管理 (Volume)

卷管理提供了 Docker 数据卷的创建、删除和查询功能。

#### 列出卷

```go
volumes, err := client.Volume().List(ctx)
if err != nil {
    return err
}

for _, volume := range volumes {
    fmt.Printf("Name: %s, Driver: %s, Mountpoint: %s\n", 
        volume.Name, volume.Driver, volume.Mountpoint)
}
```

#### 创建卷

```go
import (
    "github.com/go-zoox/docker/volume"
    vo "github.com/docker/docker/api/types/volume"
)

vol, err := client.Volume().Create(ctx, func(opt *volume.CreateOption) {
    opt.Name = "my-volume"
    opt.Driver = "local"
    opt.Labels = map[string]string{
        "env": "production",
    }
})
```

#### 删除卷

```go
err := client.Volume().Remove(ctx, "volume-name", func(opt *volume.RemoveOption) {
    opt.Force = false
})
```

#### 清理未使用的卷

```go
import "github.com/docker/docker/api/types/filters"

report, err := client.Volume().Prune(ctx, func(opt *volume.PruneOption) {
    opt.Filters = filters.NewArgs()
    opt.Filters.Add("label", "env=test")
})
fmt.Printf("Deleted: %d, Space Reclaimed: %d\n", 
    len(report.VolumesDeleted), report.SpaceReclaimed)
```

### 系统信息 (Info)

系统信息提供了 Docker 系统相关的查询功能。

#### 获取系统信息

```go
info, err := client.Info().Get(ctx)
if err != nil {
    return err
}

fmt.Printf("Containers: %d, Images: %d\n", 
    info.Containers, info.Images)
fmt.Printf("Server Version: %s\n", info.ServerVersion)
```

#### 获取版本信息

```go
version, err := client.Info().Version(ctx)
if err != nil {
    return err
}

fmt.Printf("API Version: %s\n", version.APIVersion)
fmt.Printf("Go Version: %s\n", version.GoVersion)
```

#### 获取磁盘使用情况

```go
diskUsage, err := client.Info().Disk(ctx)
if err != nil {
    return err
}

fmt.Printf("Images Size: %d\n", diskUsage.ImagesSize)
fmt.Printf("Containers Size: %d\n", diskUsage.ContainersSize)
fmt.Printf("Volumes Size: %d\n", diskUsage.VolumesSize)
```

## CLI 工具

项目还提供了一个命令行工具，方便快速操作 Docker。

### 安装 CLI

```bash
go install github.com/go-zoox/docker/cmd/docker@latest
```

### 使用示例

```bash
# 列出容器
docker container list

# 查看容器日志
docker container logs <container-id>

# 执行容器命令
docker container exec <container-id> ls -la

# 列出镜像
docker image list

# 拉取镜像
docker image pull nginx:latest

# 列出网络
docker network list

# 创建网络
docker network create my-network

# 列出卷
docker volume list

# 创建卷
docker volume create my-volume
```

## 完整示例

```go
package main

import (
    "context"
    "fmt"
    "log"
    "time"

    "github.com/go-zoox/docker"
    "github.com/go-zoox/docker/container"
)

func main() {
    ctx := context.Background()

    // 创建客户端
    client, err := docker.New()
    if err != nil {
        log.Fatal(err)
    }

    // 创建并运行容器
    err = client.Container().Run(ctx, func(opt *container.RunOptions) {
        opt.Image = "nginx:latest"
        opt.Name = "test-nginx"
        opt.Remove = true
    })
    if err != nil {
        log.Fatal(err)
    }

    // 列出所有容器
    containers, err := client.Container().List(ctx)
    if err != nil {
        log.Fatal(err)
    }

    fmt.Printf("Found %d containers\n", len(containers))
    for _, c := range containers {
        fmt.Printf("- %s: %s\n", c.Names[0], c.Status)
    }

    // 获取系统信息
    info, err := client.Info().Get(ctx)
    if err != nil {
        log.Fatal(err)
    }

    fmt.Printf("\nDocker Info:\n")
    fmt.Printf("  Containers: %d\n", info.Containers)
    fmt.Printf("  Images: %d\n", info.Images)
    fmt.Printf("  Server Version: %s\n", info.ServerVersion)
}
```

## 许可证

GoZoox is released under the [MIT License](./LICENSE).

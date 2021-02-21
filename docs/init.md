# Initialize application

> Example application path `/home/demo`

### 1. 初始化模块

> 以`sketch`为例, 进入项目目录并执行初始化脚本.

```shell
cd /home/demo
go mod init sketch
```

### 2. 创建入口文件

> 在 `/path/demo` 目录下创建 `main.go` 文件, 内容如下

```go
package main

import (
	"fmt"

	"github.com/fuyibing/util/commands"
)

func main() {
	cmd := commands.Default()
	if err := cmd.Run(); err != nil {
		fmt.Printf("%c[%d;%d;%dm%s%c[0m\n", 0x1B, 0, 33, 41, err.Error(), 0x1B)
		return
	}
}
```

### 3. 创建项目目录

```shell
go run main.go make --type=path --name=sketch
```

### 4. 下载配置文件

> 从模板中下载项目配置文件, 选项`--name`指定模板名称, 多个模板使用逗号`,`分隔开。

```shell
go run main.go kv \
  --path=./config \
  --addr=udsdk.turboradio.cn \
  --name=go/app,go/db,go/log
```






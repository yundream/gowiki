# gowiki
## Plugin 시스템
```go
package plugin

import ()

type Plugin struct {
    version string
}

func (p Plugin) Version() string {
    return p.version
}
func (p *Plugin) SetVersion(ver string) {
    p.version = ver
}
```

```go
# plugin/helloworld
package helloworld

import (
    "bitbucket.org/dream_yun/testimport/plugin"
)

// plugin을 embeded 한다.
type Hello struct {
    plugin.Plugin
}
func New(p plugin.Plugin) *Hello {
    return &Hello{p}
}

// 플러그인 함수를 개발한다.
func (h Hello) Sum(a, b int) int {
    return a + b
}
```

테스트다. 지금은 하드코딩 했는데, 실제 구현에서는 템플릿으로 go 코드를 만들어서 빌드를 해야 할 것 같다.
# plugin 디렉토리에 있는 파일을 읽어서 import 목록을 구성한다.
# plugin 디렉토리에 있는 파일 이름으로 packagename.New()를 호출 플러그인 맵을 구성한다.
```go
package main

import (
    "bitbucket.org/dream_yun/testimport/plugin"
    "bitbucket.org/dream_yun/testimport/plugin/helloworld"
    "fmt"
    "reflect"
)

func main() {
    p := plugin.Plugin{}
    p.SetVersion("v3.0")

    // plugin을 읽어온다.
    // 실제 코드에서는 plugin/ 디렉토리에서 읽어서 처리하게 한다.
    plugin_helloworld := helloworld.New(p)
    method := reflect.TypeOf(plugin_helloworld)

    F := make(map[string]reflect.Value)
    for i := 0; i < method.NumMethod(); i++ {
        name := method.Method(i).Name
        fmt.Println("Reflect method ", name)
        reflectMethd := reflect.ValueOf(plugin_helloworld).MethodByName(name)
        F["helloworld/"+method.Method(i).Name] = reflectMethd

    }
    mm, ok := F["helloworld/Sum"]
    if ok {
        r := mm.Call([]reflect.Value{reflect.ValueOf(1), reflect.ValueOf(2)})
        fmt.Println(r[0].Int())
    }
}
```

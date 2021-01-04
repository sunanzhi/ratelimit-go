### 滑动窗口限流
- - -

> **算法思想** 将一个大的时间窗口分成多个小窗口，每次大窗口向后滑动一个小窗口，并保证大的窗口内流量不会超出最大值，这种实现比固定窗口的流量曲线更加平滑。

![avatar](https://static.sunanzhi.com/github/ratelimit-go/68747470733a2f2f63646e2e6c6561726e6b752e636f6d2f75706c6f6164732f696d616765732f3230323031322f30372f363936342f6d497a415578697042582e706e67.png)

----

#### 使用
> go get github.com/sunanzhi/ratelimit-go

#### 示例

```go
import(
    "fmt"
    "net/http"
    ratelimit "github.com/sunanzhi/ratelimit-go"
)

func main() {
    slideWindow, _ := ratelimit.Init(10, 5, 200)

    http.HandleFunc("/ratelimit", func(w http.ResponseWriter, r *http.Request) {
        err := slideWindow.Limiting()
        if err != nil {
            w.WriteHeader(http.StatusBadGateway)
            fmt.Println(err.Error())
        } else {
            w.Write([]byte("httpserver"))
        }
    })
    fmt.Println("Starting server ...")
    http.ListenAndServe(":9090", nil)
}
```

#### 原理实现

@todo

#### hey.exe 轻量压缩测试结果

> 用golang的boom来做压测,要提前安装boom工具
> go get github.com/rakyll/hey
> go install github.com/rakyll/hey

**执行以下命令**

> C:\Users\0\go\bin> .\hey.exe -c 6 -n 300 -q 6 -t 80 http://localhost:9090/ratelimit

<details>
  <summary>控制台输出</summary>

  ```json
    Summary:
    Total:        8.3399 secs
    Slowest:      0.0608 secs
    Fastest:      0.0001 secs
    Average:      0.0030 secs
    Requests/sec: 35.9718
    
    Total data:   2600 bytes
    Size/request: 8 bytes

    Response time histogram:
    0.000 [1]     |
    0.006 [269]   |■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■
    0.012 [14]    |■■
    0.018 [4]     |■
    0.024 [3]     |
    0.030 [4]     |■
    0.037 [1]     |
    0.043 [0]     |
    0.049 [2]     |
    0.055 [0]     |
    0.061 [2]     |


    Latency distribution:
    10% in 0.0002 secs
    25% in 0.0002 secs
    50% in 0.0002 secs
    75% in 0.0025 secs
    90% in 0.0062 secs
    95% in 0.0132 secs
    99% in 0.0461 secs

    Details (average, fastest, slowest):
    DNS+dialup:   0.0001 secs, 0.0001 secs, 0.0608 secs
    DNS-lookup:   0.0001 secs, 0.0000 secs, 0.0044 secs
    req write:    0.0000 secs, 0.0000 secs, 0.0001 secs
    resp wait:    0.0028 secs, 0.0001 secs, 0.0607 secs
    resp read:    0.0000 secs, 0.0000 secs, 0.0002 secs

    Status code distribution:
    [200] 200 responses
    [502] 100 responses
  ```
</details>
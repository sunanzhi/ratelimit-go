### 滑动窗口限流
- - -

> **算法思想** 将一个大的时间窗口分成多个小窗口，每次大窗口向后滑动一个小窗口，并保证大的窗口内流量不会超出最大值，这种实现比固定窗口的流量曲线更加平滑。

![slidewindow](https://static.sunanzhi.com/github/ratelimit-go/20190925171656739496.png)

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
    var (
        limitTime = 10
        bucketCount = 5
        limitCount = 200
    )
    slideWindow, _ := ratelimit.Init(limitTime, bucketCount, limitCount)

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


#### 原理实现

> 基于双向链表实现,减少滑动过程中过多的创建销毁

![doublyLinkedList](https://static.sunanzhi.com/github/ratelimit-go/20190925171656739497.png)

> `bukcetCount` 数值就是链表长度,每次初始化链表将会以当前时间末节点 + 节点的时间范畴开始 依次递减

**举例:** 

- 限速时间 `limitTime=4s` 
- 节点数量 `bucketCount=4`
- 且当前时间为第 `4` 秒

> 那么每个节点统计时间范畴为 `1s` 

**节点数据如下:(其他数据不具体展示)**

![nodeData](https://static.sunanzhi.com/github/ratelimit-go/20210105161233.png)

**窗口滑动:**

> 通过定时任务滑动,定时时间为每个节点统计的时间范畴值(单位毫秒)

**说明:**

- 从初始化到后续的滑动,链表有效时间的节点只有末节点
- 因此请求进来的都统一放入末节点 `count` 
- 因为滑动窗口会改变末节点,因此不需要遍历当前请求时间属于哪个节点


**注意:**

- 此包只有全部限流,不能映射任何 `key` 或者其他字段来做分组
- 如果有需求请自己创建线程安全 `map` 来映射限流数据
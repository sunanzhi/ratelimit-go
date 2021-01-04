package ratelimit

import (
	"container/ring"
	"errors"
	"fmt"
	"math"
	"sync"
	"time"
)

// 窗口节点
type bucketNode struct {
	startTime int64 // 节点开始时间
	endTime   int64 // 节点结束时间
	count     int   // 节点统计访问数量
}

// 滑动窗口
type SlideWindow struct {
	mu             sync.Mutex // 锁
	LimitTime      int        // 限流区间长度 单位(秒)
	BucketCount    int        // 滑动窗口个数
	LimitCount     int        // 限流数量
	Count          int        // 统计数量
	StartTime      int64      // 开始时间
	EndTime        int64      // 结束时间
	BucketInterval int        // 窗口区间时间 单位(毫秒)
	BucketChain    *ring.Ring // 窗口环形链
}

// 初始化
func Init(limitTime int, bucketCount int, limitCount int) (*SlideWindow, error) {
	if bucketCount <= 0 || limitTime <= 0 || limitCount <= 0 {
		return nil, errors.New("Parameter Error")
	}
	timeInterval := int(math.Floor(float64((limitTime * 1000) / bucketCount)))
	// 初始化滑动窗口
	bucketChain := ring.New(bucketCount)
	// 原始开始时间
	oriStartTime := time.Now().UnixNano() / 1e6
	bucketChain = bucketChain.Prev()
	for i := 0; i < bucketCount; i++ {
		bucketChain.Value = bucketNode{
			startTime: oriStartTime,
			endTime:   oriStartTime + int64(timeInterval),
			count:     0,
		}
		oriStartTime -= int64(timeInterval)
		bucketChain = bucketChain.Prev()
	}
	// 回到初始节点
	bucketChain = bucketChain.Next()
	slideWindow := &SlideWindow{
		LimitTime:      limitTime,
		BucketCount:    bucketCount,
		LimitCount:     limitCount,
		Count:          0,
		StartTime:      oriStartTime,
		EndTime:        oriStartTime + int64(bucketCount*timeInterval*limitTime),
		BucketInterval: timeInterval,
		BucketChain:    bucketChain,
	}
	go slideWindow.slide()

	return slideWindow, nil
}

// 定时滑动
func (slideWindow *SlideWindow) slide() {
	timer := time.NewTicker(time.Millisecond * time.Duration(slideWindow.BucketInterval))
	// 定时每隔 窗口区间时间 毫秒刷新一次滑动窗口数据
	for range timer.C {
		slideWindow.mu.Lock()
		headChain := slideWindow.BucketChain
		tailNodeValue := headChain.Prev().Value
		tailNode := tailNodeValue.(bucketNode)
		// 减去即将淘汰的节点数据统计
		headNode := headChain.Value.(bucketNode)
		slideWindow.Count -= headNode.count
		// 窗口移动到下一个
		headChain.Value = bucketNode{
			startTime: tailNode.endTime,
			// 避免有损耗 使用当前时间毫秒 + 区间
			endTime: (time.Now().UnixNano() / 1e6) + int64(slideWindow.BucketInterval),
			count:   0,
		}
		// 首位节点移动
		headChain = headChain.Next()
		slideWindow.BucketChain = headChain
		slideWindow.mu.Unlock()
		slideWindow.Print()
	}
}

// 限流
func (slideWindow *SlideWindow) Limiting() error {
	slideWindow.mu.Lock()
	defer slideWindow.mu.Unlock()
	tailChain := slideWindow.BucketChain.Prev()
	tailNodeValue := tailChain.Value
	tailNode := tailNodeValue.(bucketNode)
	// 单节点 || 全部节点 超过限制
	if tailNode.count+1 > slideWindow.LimitCount || slideWindow.Count+1 > slideWindow.LimitCount {
		return errors.New("Rate Limited")
	}
	tailNode.count++
	tailChain.Value = tailNode
	headChain := tailChain.Next()
	slideWindow.Count++
	slideWindow.BucketChain = headChain

	return nil
}

// 打印数据
func (slideWindow *SlideWindow) Print() {
	fmt.Println(slideWindow)
	// 打印节点
	chain := slideWindow.BucketChain
	for i := 0; i < slideWindow.BucketCount; i++ {
		fmt.Println(chain.Value)
		chain = chain.Next()
	}
}

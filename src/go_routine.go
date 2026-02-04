package main

import (
	"fmt"
	"image"
	"runtime"
	"sync"
	"time"
)

type TestGoRoutine struct{}

func hello() {
	fmt.Println("Hello from a goroutine!")
}

func test() {
	hello()
	fmt.Println("Main function is running.")
}

func test2() {
	go hello()
	fmt.Println("Main function is running.")
}

func test3() {
	go hello()
	fmt.Println("Main function is running.")
	time.Sleep(time.Second)
}

var wg sync.WaitGroup

func hello1(i int) {
	defer wg.Done()
	fmt.Printf("Hello from goroutine %d\n", i)
}

func test4() {
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go hello1(i)
	}
	wg.Wait()
}

// func test5() {
// 	go func() {
// 		i := 0
// 		for {
// 			i++
// 			fmt.Printf("Goroutine printing %d\n", i)
// 			time.Sleep(time.Second)
// 		}
// 	}()

// 	i := 0
// 	for {
// 		i++
// 		fmt.Printf("Main function printing %d\n", i)
// 		time.Sleep(time.Second)

// 		if i == 2 {
// 			break
// 		}
// 	}
// }

func test7() {
	go func(s string) {
		for i := 0; i < 2; i++ {
			fmt.Println(s)
		}
	}("world")

	for i := 0; i < 2; i++ {
		runtime.Gosched()
		fmt.Println("hello")
	}
}

func test8() {
	go func() {
		defer fmt.Println("A.defer")
		func() {
			defer fmt.Println("B.defer")
			// 结束协程
			runtime.Goexit()
			defer fmt.Println("C.defer")
			fmt.Println("B")
		}()
		fmt.Println("A")
	}()
	for {
	}
}

func test9() {
	// 创建带缓冲的 channel，容量为 1
	ch := make(chan int, 1)
	Log.Debug("Channel created.")

	// 发送数据不会阻塞（因为有缓冲空间）
	ch <- 10
	Log.Debug("Value sent to channel.")

	// 接收数据
	value := <-ch
	Log.Debug(fmt.Sprintf("Received value: %d", value))
}

func test10() {
	ch := make(chan int)
	go test11(ch)
	ch <- 10
	Log.Debug("send success")
}

func test11(c chan int) {
	value := <-c
	Log.Debug(fmt.Sprintf("Received value: %d", value))
}

func test12() {
	c := make(chan int)
	go func() {
		for i := 0; i < 5; i++ {
			c <- i
		}
		close(c)
	}()

	for {
		if data, ok := <-c; ok {
			Log.Debug(data)
		} else {
			break
		}
	}

	Log.Debug("main stop")
}

func test13(out chan<- int) {
	for i := 0; i < 100; i++ {
		out <- i
	}
	close(out)
}

func test14(out chan<- int, in <-chan int) {
	for i := range in {
		out <- i * i
	}
	close(out)
}

func test15(in <-chan int) {
	for i := range in {
		Log.Debug(i)
	}
}

func test16() {
	ch1 := make(chan int)
	ch2 := make(chan int)
	go test13(ch1)
	go test14(ch2, ch1)
	test15(ch2)
}

func test17() {
	ticker := time.NewTicker(2 * time.Second)
	i := 0

	go func() {
		for {
			i++
			Log.Debug(<-ticker.C)
			if i == 5 {
				ticker.Stop()
			}
		}
	}()
	for {

	}
}

func ts18(ch chan string) {
	time.Sleep(5 * time.Second)
	ch <- "data from ts18"
}

func ts19(ch chan string) {
	time.Sleep(2 * time.Second)
	ch <- "data from ts19"
}

func ts20() {
	ch1 := make(chan string)
	ch2 := make(chan string)

	go ts18(ch1)
	go ts19(ch2)

	select {
	case msg1 := <-ch1:
		Log.Debug("Received:", msg1)
	case msg2 := <-ch2:
		Log.Debug("Received:", msg2)
	}
}

func ts21() {
	int_chan := make(chan int, 1)
	string_chan := make(chan string, 1)

	go func() {
		int_chan <- 1
	}()

	go func() {
		string_chan <- "hello"
	}()

	select {
	case value := <-int_chan:
		Log.Debug("Received int:", value)
	case msg := <-string_chan:
		Log.Debug("Received string:", msg)
	}
	Log.Debug("stop")
}

func ts22() {
	output1 := make(chan string, 10)
	go ts23(output1)
	for s := range output1 {
		Log.Debug("read ", s)
		time.Sleep(time.Second)
	}
}

func ts23(ch chan<- string) {
	for {
		select {
		case ch <- "hello":
			Log.Debug("write hello")
		default:
			Log.Debug("channel full")
		}
		time.Sleep(500 * time.Millisecond)
	}
}

var x int64
var wg1 sync.WaitGroup
var lock sync.Mutex

func ts24() {
	for i := 0; i < 5000; i++ {
		lock.Lock()
		x = x + 1
		lock.Unlock()
	}
	wg1.Done()
}

func ts25() {
	wg1.Add(2)
	go ts24()
	go ts24()
	wg1.Wait()
	Log.Debug("x =", x)
}

var (
	x1     int64
	wg2    sync.WaitGroup
	lock2  sync.Mutex
	rwLock sync.RWMutex
)

func ts26() {
	rwLock.Lock()
	x1 = x1 + 1
	time.Sleep(10 * time.Millisecond)
	rwLock.Unlock()
	wg2.Done()
}

func ts27() {
	rwLock.RLock()
	time.Sleep(time.Millisecond)
	rwLock.RUnlock()
	wg2.Done()
}

func ts28() {
	start := time.Now()
	for i := 0; i < 10; i++ {
		wg2.Add(1)
		go ts26()
	}

	for i := 0; i < 1000; i++ {
		wg2.Add(1)
		go ts27()
	}

	wg2.Wait()
	end := time.Now()
	Log.Debug("x1 =", x1, "time taken:", end.Sub(start))
}

func ts29() {
	defer wg1.Done()
	Log.Debug("Hello, ts29")
}

func ts30() {
	wg1.Add(1)
	go ts29()
	Log.Debug("main ts29 done")
	wg1.Wait()
	Log.Debug("Done")
}

var icons map[string]image.Image
var loadIconOnce sync.Once

func loadIcon(filename string) image.Image {
	loadIconOnce.Do(func() {
		Log.Debug("Loading icon:", filename)
		// 模拟加载图标的耗时操作
		time.Sleep(2 * time.Second)
	})
	// 返回一个空的 image.Image 作为示例
	return image.NewRGBA(image.Rect(0, 0, 64, 64))
}

func ts31() {
	icons = map[string]image.Image{
		"left":  loadIcon("left.png"),
		"right": loadIcon("right.png"),
		"up":    loadIcon("up.png"),
		"down":  loadIcon("down.png"),
	}
}

func ts32(name string) image.Image {
	loadIconOnce.Do(ts31)
	return icons[name]
}

func (t TestGoRoutine) Test() {
	ts30()
	// ts21()
	// test17()
	// test8()
	// test()
	// test2()
	// test3()
	// test4()
	// test5()
	// test7()
}

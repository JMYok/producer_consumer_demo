package main

import (
	"pcdemo/model"
	"sync"
	"time"
)

// 缓冲池
var taskCh = make(chan model.Task, 10)

// 生产者需要生产的任务数量
const taskNum int64 = 10000

// 停止运行的信号
var done = make(chan struct{})

func producer(wo chan<- model.Task) {
	var i int64
	for i = 1; i <= taskNum; i++ {
		t := model.Task{
			ID: i,
		}
		wo <- t
	}
	// 单个生产者就可以直接关闭通道了，关闭后，消费者任然可以消费
	close(wo)
}

// 生产者在生产的时候，可能存在数据竞争问题
func producer2(wo chan<- model.Task, startNum int64, nums int64) {
	var i int64
	for i = startNum; i < startNum+nums; i++ {
		t := model.Task{
			ID: i,
		}
		wo <- t
	}
}

func producer3(wo chan<- model.Task, done chan struct{}) {
	var i int64
	for {
		if i >= taskNum { // 无限生产
			i = 0
		}
		i++
		t := model.Task{
			ID: i,
		}
		select {
		case wo <- t:
		case <-done:
			model.OutObj.Println("生产者退出")
			return
		}
	}
}

func consumer(ro <-chan model.Task) {
	for t := range ro {
		if t.ID != 0 {
			t.Run()
		}
	}
}

func consumer2(ro <-chan model.Task, done chan struct{}) {
	for {
		select {
		case t := <-ro:
			if t.ID != 0 {
				t.Run()
			}
		case <-done: // 这里如果直接退出的话，可能 channel 里面还有值没有被消费（有缓存区的情况）
			for t := range ro { // 生产者那边已经停止，消息不会再生产。消费者这里将所有消息消费后，就可以 退出了
				if t.ID != 0 {
					t.Run()
				}
			}
			model.OutObj.Println("消费者退出")
			return
		}
	}
}

func ExecOneOne() {
	wg := &sync.WaitGroup{}
	wg.Add(2)
	go func(wg *sync.WaitGroup) {
		defer wg.Done()
		producer(taskCh)
	}(wg)
	go func(wg *sync.WaitGroup) {
		defer wg.Done()
		consumer(taskCh)
	}(wg)

	wg.Wait()
	model.OutObj.OutPut()
}

func ExecOneMany() {
	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func(wg *sync.WaitGroup) {
		defer wg.Done()
		producer(taskCh)
	}(wg)

	var i int64
	for i = 0; i < taskNum; i++ {
		//每100个任务开一个消费者
		if i%100 == 0 {
			wg.Add(1)
			go func(wg *sync.WaitGroup) {
				defer wg.Done()
				consumer(taskCh)
			}(wg)
		}
	}
	wg.Wait()
	model.OutObj.Println("执行完成")
	model.OutObj.OutPut()
}

func ExecManyOne() {
	cwg := &sync.WaitGroup{}
	pwg := &sync.WaitGroup{}
	var i int64
	var nums int64 = 100
	for i = 0; i < taskNum; i += nums {
		if i >= taskNum {
			break
		}
		// 每个生产者生产 100 个任务
		pwg.Add(1)

		go func(i int64) {
			defer pwg.Done()
			producer2(taskCh, i, nums)
		}(i)
	}

	cwg.Add(1)
	go func() {
		defer cwg.Done()
		consumer(taskCh)
	}()

	pwg.Wait()
	go close(taskCh)

	cwg.Wait()
	model.OutObj.Println("执行成功")
	model.OutObj.OutPut()
}

func ExecManyMany() {
	// 多个生产者
	for i := 0; i < 8; i++ {
		go producer3(taskCh, done)
	}

	// 多个消费者
	for i := 0; i < 8; i++ {
		go consumer2(taskCh, done)
	}

	time.Sleep(time.Second * 5)
	// 一定要先关闭 done，再关闭通道。防止向已关闭的 channel 写入数据，报异常
	close(done)
	close(taskCh)

	model.OutObj.Println("执行成功")
	model.OutObj.OutPut()
}

func main() {
	//ExecOneOne()
	//ExecOneMany()
	//ExecManyOne()
	ExecManyMany()
}

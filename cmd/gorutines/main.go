package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"runtime"
	"sync"
	"syscall"
	"time"
)

func main() {
	//runtime.GOMAXPROCS(4)
	//runtime.Gosched()
	//println(runtime.NumCPU())
	//defer println("defer 1")
	//defer println("defer 2")
	//withWait()
	//writeWithoutConcurrent()
	//writeWitConcurrent()
	//nilChannel()
	//unbuffChan()
	//buffChan()
	//baseSelect()
	// gracefulShutdown()
	workerPool()
}

func workerPool() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	wg := &sync.WaitGroup{}
	numbersToProcess, processedNumbers := make(chan int, 5), make(chan int, 5)

	for i := 0; i < runtime.NumCPU(); i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			worker(ctx, numbersToProcess, processedNumbers)
		}()
	}

	go func() {
		for i := 0; i < 1000; i++ {
			numbersToProcess <- i
		}
		close(numbersToProcess)
	}()

	go func() {
		wg.Wait()
		close(processedNumbers)
	}()

}

func worker(ctx context.Context, toProcess <-chan int, processed chan<- int) {
	for {
		select {
		case <-ctx.Done():
			return
		case value, ok := <-toProcess:
			if !ok {
				return
			}
			time.Sleep(time.Millisecond)
			processed <- value * value
		}
	}
}

func gracefulShutdown() {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	timer := time.After(10 * time.Second)
	for {
		select {
		case <-timer:
			fmt.Println("times'up")
			return
		case sig := <-sigChan:
			fmt.Println("stopped by signal", sig)
			return
		}
	}
}

func nilChannel() {
	var nilChannel chan int
	fmt.Printf("%d, %d, %v \n", len(nilChannel), cap(nilChannel), nilChannel)
	//nilChannel <- 1
	//fmt.Println(<-nilChannel)
	close(nilChannel)
}

func unbuffChan() {
	unbuffChan := make(chan int)
	fmt.Printf("%d, %d, %v \n", len(unbuffChan), cap(unbuffChan), unbuffChan)
	unbuffChan <- 1
}

func withWait() {
	var wg sync.WaitGroup
	defer wg.Wait()

	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			println(i)

		}(i)
	}

	println("exit")
}

func writeWithoutConcurrent() {
	start := time.Now()
	var counter int
	for i := 0; i < 100; i++ {
		time.Sleep(time.Nanosecond)
		counter++
	}

	fmt.Printf("seconds %v\n", time.Since(start).Seconds())
}

func writeWitConcurrent() {
	start := time.Now()
	var wg sync.WaitGroup
	var counter uint32
	var rw sync.RWMutex
	wg.Add(1000)

	for i := 0; i < 1000; i++ {
		go func() {
			defer wg.Done()

			rw.RLock()
			time.Sleep(time.Nanosecond)
			_ = counter
			rw.RUnlock()

			rw.Lock()
			counter++
			rw.Unlock()
		}()

	}
	wg.Wait()
	fmt.Printf("seconds %v, counter %v\n", time.Since(start).Seconds(), counter)

}

func buffChan() {
	buffChan := make(chan int, 2)
	buffChan <- 1
	buffChan <- 2

	fmt.Println(<-buffChan)
	fmt.Println(<-buffChan)
}

func baseSelect() {
	buffChan := make(chan string, 2)
	//buffChan <- "first"

	select {
	case str := <-buffChan:
		fmt.Println(str)

	case buffChan <- "second":
		fmt.Println("write", <-buffChan, <-buffChan)

	}
}

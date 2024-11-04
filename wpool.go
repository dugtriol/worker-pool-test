package main

import (
	`fmt`
	`sync`
	`sync/atomic`
	`time`
)

type WPool struct {
	jobs    chan string
	results chan string
	wg      sync.WaitGroup
	workers chan struct{}
	done    chan struct{}
	count   int32
}

func NewWPool(maxWorkers int) *WPool {
	w := &WPool{
		jobs:    make(chan string),
		results: make(chan string),
		wg:      sync.WaitGroup{},
		workers: make(chan struct{}, maxWorkers),
		done:    make(chan struct{}),
		count:   int32(maxWorkers),
	}

	for worker := 1; worker <= maxWorkers; worker++ {
		w.wg.Add(1)
		go w.Work(worker)
	}

	return w
}

func (w *WPool) Work(id int) {
	defer w.wg.Done()
	for {
		select {
		case job, ok := <-w.jobs:
			if !ok {
				return
			}
			w.workers <- struct{}{}
			fmt.Printf("worker %d started task %s\n", id, job)
			time.Sleep(3 * time.Second)
			fmt.Printf("worker %d finished task %s\n", id, job)
			w.results <- fmt.Sprintf("task %s done", job)
			<-w.workers
		case <-w.done:
			return
		}
	}
}

func (w *WPool) Add() {
	count := atomic.AddInt32(&w.count, 1)
	w.wg.Add(1)
	go w.Work(int(count))
	fmt.Println("add new worker ", count)
}

func (w *WPool) Remove() {
	num := atomic.LoadInt32(&w.count)
	if num > 0 {
		w.done <- struct{}{}
		atomic.AddInt32(&w.count, -1)
		fmt.Println("remove worker ", num)
	}
}

func (w *WPool) Close() {
	w.wg.Wait()
	close(w.results)
	fmt.Println("all tasks are completed")
}

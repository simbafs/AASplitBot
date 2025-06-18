package workerpool

import (
	"slices"
	"sync"
)

type Pool struct {
	mutex  sync.Mutex
	wg     sync.WaitGroup
	tasks  []Task
	errs   []error
	closed bool
}

type Task func() error

// New returns a new worker pool.
func New() *Pool {
	return &Pool{}
}

// Do adds a task to the pool and runs it in a new goroutine.
func (p *Pool) Do(task Task) {
	p.mutex.Lock()
	if p.closed {
		p.mutex.Unlock()
		return // 如果 pool 已經關閉，就不接受新的任務
	}
	p.wg.Add(1)
	p.mutex.Unlock()

	go func() {
		defer p.wg.Done()
		err := task()
		if err != nil {
			p.mutex.Lock()
			p.errs = append(p.errs, err)
			p.mutex.Unlock()
		}
	}()
}

// Wait blocks until all tasks have completed.
func (p *Pool) Wait() {
	done := make(chan struct{})
	go func() {
		p.wg.Wait()
		p.mutex.Lock()
		p.closed = true
		p.mutex.Unlock()
		close(done)
	}()
	<-done
}

// Errors returns a slice of errors from tasks that failed.
func (p *Pool) Errors() []error {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	return slices.Clone(p.errs) // 回傳一個複製切片避免外部修改
}

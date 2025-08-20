package hw05parallelexecution

import (
	"context"
	"errors"
	"sync"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

func runTask(task Task, mu *sync.Mutex, countErr *int, m int, cancel context.CancelFunc) (stop bool) {
	mu.Lock()
	if *countErr >= m {
		mu.Unlock()
		cancel()
		return true
	}
	mu.Unlock()

	if err := task(); err != nil {
		mu.Lock()
		*countErr++
		if *countErr >= m {
			mu.Unlock()
			cancel()
			return true
		}
		mu.Unlock()
	}
	return false
}

func worker(
	ctx context.Context,
	ch <-chan Task,
	m int,
	mu *sync.Mutex,
	countErr *int,
	cancel context.CancelFunc,
	wg *sync.WaitGroup,
) {
	defer wg.Done()
	for {
		select {
		case <-ctx.Done():
			return
		case task, ok := <-ch:
			if !ok {
				return
			}
			if runTask(task, mu, countErr, m, cancel) {
				return
			}
		}
	}
}

func Run(tasks []Task, n, m int) error {
	if m <= 0 { // максимум 0 ошибок — сразу ошибка
		return ErrErrorsLimitExceeded
	}

	ch := make(chan Task, len(tasks))
	countErr := 0
	var mu sync.Mutex
	var wg sync.WaitGroup
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	for i := 0; i < n; i++ {
		wg.Add(1)
		go worker(ctx, ch, m, &mu, &countErr, cancel, &wg)
	}

	for _, val := range tasks {
		ch <- val
	}
	close(ch)
	wg.Wait()
	if countErr >= m {
		return ErrErrorsLimitExceeded
	}
	return nil
}

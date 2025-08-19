package hw05parallelexecution

import (
	"context"
	"errors"
	"sync"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

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
		go func(ctx context.Context) {
			defer wg.Done()
			for {
				select {
				case <-ctx.Done():
					return
				case task, ok := <-ch:
					if !ok {
						return
					}
					mu.Lock()
					if countErr >= m {
						mu.Unlock()
						cancel()
						return
					}
					mu.Unlock()

					// Выполняем задачу
					if err := task(); err != nil {
						mu.Lock()
						countErr++
						if countErr >= m {
							mu.Unlock()
							cancel()
							return
						}
						mu.Unlock()
					}
				}
			}
		}(ctx)

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

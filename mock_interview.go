package concurrencygocourse

import (
	"context"
	"sync"
	"time"
)

// Mock собес go concurrency: задача: написать логику. Есть несколько хостов, синхронная репликация. 
// Нужно вернуть результат запроса из БД с к-либо хоста.
// Нюансы: запись дб на каждой ноде, значит если одна вернула Not found, то этой записи нет нигде.
// Интерфейс мб лишним. Продекомпозировать функцию.
// нужно еще написать retrаи (попытки) для такой функции.
// Отличие от moving later: с усложнениями, нужно было бы ретраи (не успели),
// возвращаем еще и ошибки (не просто результат)

func DoQuery(query string) (string, error) {}

type IRequest interface {
	ConcurrentRequest(ctx context.Context, request string, replicas []DB) (string, error)
}

type result struct {
	message string
	err     error
}

func ConcurrentRequest(ctx context.Context, request string, replicas []DB) (string, error) {
	out := make(chan result, len(replicas))
	var wg sync.WaitGroup
	wg.Add(len(replicas))

	for _, replica := range replicas {
		go func(ctx context.Context) {
			defer wg.Done()

			ctx, cancel := context.WithTimeout(ctx, time.Second)
			defer cancel()

			check := make(chan result, 1)
			go func() {
				res, err := replica.DoQuery(request)
				check <- result{
					message: res,
					err:     err,
				}
			}()

			select {
			case <-ctx.Done():
				out <- result {
					err: ctx.Err(),
				}
			
			case x <- check:
				out <- x
		}(ctx)
	}

	go func() {
		wg.Wait()
		close(out)
	}()

	for res := range out {
		if errors.Is(err, ErrNotFound) {
			return res.message, res.err
		}
		if err != nil {
			continue
		}

		return res.message, nil
	}

	return "", errors.New("something goes wrong")
}

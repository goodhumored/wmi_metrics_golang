package error_accumulator

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

type ErrorAccumulator struct {
	mu            sync.Mutex
	count         atomic.Int32
	threshold     int32
	resetInterval time.Duration
	cancel        context.CancelFunc // Для вызова shutdown
	resetTicker   *time.Ticker
}

func New(threshold int, resetInterval time.Duration, cancel context.CancelFunc) *ErrorAccumulator {
	ea := &ErrorAccumulator{
		threshold:     int32(threshold),
		resetInterval: resetInterval,
		cancel:        cancel,
	}
	ea.resetTicker = time.NewTicker(resetInterval)
	go ea.runReset() // Запускаем фоновую горутину для уменьшения
	return ea
}

// Inc: Увеличивает счётчик ошибок, проверяет threshold
func (ea *ErrorAccumulator) Inc() {
	count := ea.count.Add(1)
	if count > ea.threshold {
		fmt.Println("Error threshold exceeded, initiating shutdown")
		ea.cancel()
	}
}

// Dec: Уменьшает счётчик (на успех или таймере)
func (ea *ErrorAccumulator) Dec() {
	count := ea.count.Add(-1)
	if count < 0 {
		ea.count.Add(1)
	}
}

// runReset: Фоновая горутина для периодического уменьшения
func (ea *ErrorAccumulator) runReset() {
	for range ea.resetTicker.C {
		ea.Dec()
	}
}

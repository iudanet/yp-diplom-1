package retry

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"strings"
	"time"

	"github.com/jackc/pgerrcode"
	"github.com/lib/pq"
)

var (
	ErrMaxRetriesReached = errors.New("maximum retries reached")
	// DefaultBackoffSchedule можно переопределить в тестах
	DefaultBackoffSchedule = []time.Duration{
		1 * time.Second,
		3 * time.Second,
		5 * time.Second,
		10 * time.Second,
	}
)

func WithRetry(op func() error) error {
	var lastErr error

	for _, backoff := range DefaultBackoffSchedule {
		err := op()
		if err == nil {
			return nil
		}

		if !isRetriable(err) {
			return err
		}

		lastErr = err
		log.Printf("retriable error occurred: %v, retrying in %v", err, backoff)
		time.Sleep(backoff)
	}

	return fmt.Errorf("%w: %v", ErrMaxRetriesReached, lastErr)
}

func isRetriable(err error) bool {
	// Проверяем ошибки PostgreSQL (pq)
	var pqErr *pq.Error
	if errors.As(err, &pqErr) {
		// Ошибки соединения
		switch pqErr.Code {
		case pgerrcode.ConnectionException,
			pgerrcode.ConnectionDoesNotExist,
			pgerrcode.ConnectionFailure,
			pgerrcode.SQLClientUnableToEstablishSQLConnection,
			pgerrcode.SQLServerRejectedEstablishmentOfSQLConnection,
			pgerrcode.TransactionResolutionUnknown:
			return true
		}
	}

	// Ошибки HTTP 5xx считаем retriable
	var httpErr interface {
		HTTPStatusCode() int
	}
	if errors.As(err, &httpErr) && httpErr.HTTPStatusCode() >= 500 {
		return true
	}

	// Проверяем сетевые ошибки
	var netErr net.Error
	if errors.As(err, &netErr) {
		return true
	}

	// Проверяем закрытые соединения
	if strings.Contains(err.Error(), "use of closed network connection") {
		return true
	}

	// Проверяем другие retriable ошибки
	if errors.Is(err, context.DeadlineExceeded) ||
		errors.Is(err, context.Canceled) {
		return true
	}

	return false
}

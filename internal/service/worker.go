package service

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/iudanet/yp-diplom-1/internal/models"
	"github.com/iudanet/yp-diplom-1/internal/repo"
)

type Worker struct {
	repo          repo.Repositories
	accrualClient AccrualClientInterface
	batchSize     int
	pollInterval  time.Duration
	logger        *log.Logger
}

func NewWorker(repo repo.Repositories, accrualClient AccrualClientInterface) *Worker {
	return &Worker{
		repo:          repo,
		accrualClient: accrualClient,
		batchSize:     10,
		pollInterval:  500 * time.Millisecond,
		logger:        log.New(log.Writer(), "worker: ", log.LstdFlags),
	}
}

func (w *Worker) Run(ctx context.Context) {
	w.logger.Println("Starting worker")
	defer w.logger.Println("Worker stopped")

	ticker := time.NewTicker(w.pollInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			if err := w.processBatch(ctx); err != nil {
				w.logger.Printf("Error processing batch: %v", err)
			}
		}
	}
}

func (w *Worker) processBatch(ctx context.Context) error {
	startTime := time.Now()
	// w.logger.Println("Processing new batch of orders")

	orders, err := w.repo.GetOrdersForProcessing(ctx, w.batchSize)
	if err != nil {
		return fmt.Errorf("failed to get orders for processing: %w", err)
	}

	if len(orders) == 0 {
		// w.logger.Println("No orders to process")
		return nil
	}

	w.logger.Printf("Processing %d orders", len(orders))

	for _, order := range orders {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			if err := w.processOrder(ctx, order); err != nil {
				w.logger.Printf("Error processing order %s: %v", order.Number, err)
			}
		}
	}

	w.logger.Printf("Batch processed in %v", time.Since(startTime))
	return nil
}

func (w *Worker) processOrder(ctx context.Context, order models.OrderUser) error {
	if order.Status == models.OrderUserStatusNew {
		w.logger.Printf("Updating order %s status to PROCESSING", order.Number)
		if err := w.repo.UpdateOrderStatus(ctx, order.Number, models.OrderUserStatusProcessing); err != nil {
			return fmt.Errorf("failed to update order status to PROCESSING: %w", err)
		}
	}

	w.logger.Printf("Getting accrual info for order %s", order.Number)
	accrualInfo, err := w.accrualClient.GetOrderAccrual(ctx, order.Number)
	if err != nil {
		return fmt.Errorf("failed to get accrual info: %w", err)
	}

	if accrualInfo == nil {
		w.logger.Printf("Order %s not found in accrual system, marking as INVALID", order.Number)
		if err := w.repo.UpdateOrderStatus(ctx, order.Number, models.OrderUserStatusInvalid); err != nil {
			return fmt.Errorf("failed to update order status to INVALID: %w", err)
		}
		return nil
	}

	switch accrualInfo.Status {
	case models.OrderAccrualStatusProcessed:
		if accrualInfo.Accrual == 0 {
			w.logger.Printf("No accrual for processed order %s", order.Number)
			return nil
		}
		// Конвертируем рубли в копейки перед сохранением
		accrualCents := int64(accrualInfo.Accrual * 100)
		w.logger.Printf(
			"Updating order %s status to PROCESSED with accrual %v (копеек: %d)",
			order.Number,
			accrualInfo.Accrual,
			accrualCents,
		)
		if err := w.repo.UpdateOrderAccrual(ctx, order.Number, models.OrderUserStatusProcessed, accrualCents); err != nil {
			return fmt.Errorf("failed to update order status to PROCESSED: %w", err)
		}

	case models.OrderAccrualStatusInvalid:
		w.logger.Printf("Updating order %s status to INVALID", order.Number)
		if err := w.repo.UpdateOrderStatus(ctx, order.Number, models.OrderUserStatusInvalid); err != nil {
			return fmt.Errorf("failed to update order status to INVALID: %w", err)
		}

	case models.OrderAccrualStatusProcessing, models.OrderAccrualStatusRegistered:
		w.logger.Printf("Order %s still processing in accrual system", order.Number)
	}

	return nil
}

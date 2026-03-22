package scraper

import (
	"context"
	"fmt"
	"time"
)

type RetryScraper struct {
	inner       Scraper
	maxRetries  int
	baseDelay   time.Duration
}

func NewRetryScraper(inner Scraper, maxRetries int, baseDelay time.Duration) *RetryScraper {
	return &RetryScraper{
		inner:      inner,
		maxRetries: maxRetries,
		baseDelay:  baseDelay,
	}
}

func (r *RetryScraper) Name() string {
	return fmt.Sprintf("Retry(%s)", r.inner.Name())
}

func (r *RetryScraper) CanSearchByHash() bool {
	return r.inner.CanSearchByHash()
}

func (r *RetryScraper) CanSearchByName() bool {
	return r.inner.CanSearchByName()
}

func (r *RetryScraper) SearchByHash(ctx context.Context, query SearchQuery) (*Metadata, error) {
	return r.doWithRetry(ctx, func() (*Metadata, error) {
		return r.inner.SearchByHash(ctx, query)
	})
}

func (r *RetryScraper) SearchByName(ctx context.Context, query SearchQuery) (*Metadata, error) {
	return r.doWithRetry(ctx, func() (*Metadata, error) {
		return r.inner.SearchByName(ctx, query)
	})
}

func (r *RetryScraper) doWithRetry(ctx context.Context, fn func() (*Metadata, error)) (*Metadata, error) {
	var lastErr error
	delay := r.baseDelay

	for i := 0; i <= r.maxRetries; i++ {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
		}

		meta, err := fn()
		if err == nil {
			return meta, nil
		}

		lastErr = err
		fmt.Printf("[%s] Intento %d fallido, reintentando en %v: %v\n", r.inner.Name(), i+1, delay, err)
		
		timer := time.NewTimer(delay)
		select {
		case <-ctx.Done():
			timer.Stop()
			return nil, ctx.Err()
		case <-timer.C:
			delay *= 2 // Backoff exponencial
		}
	}

	return nil, fmt.Errorf("máximos reintentos alcanzados: %w", lastErr)
}

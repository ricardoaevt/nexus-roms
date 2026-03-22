package scraper_test

import (
	"context"
	"errors"
	"romsRename/internal/scraper"
	"romsRename/internal/scraper/mocks"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestRetryScraperSearchByHash(t *testing.T) {
	mockInner := mocks.NewScraper(t)
	query := scraper.SearchQuery{HashMD5: "123"}
	ctx := context.Background()

	t.Run("success on first try", func(t *testing.T) {
		expectedMeta := &scraper.Metadata{Name: "Mario"}
		mockInner.On("SearchByHash", ctx, query).Return(expectedMeta, nil).Once()

		retryS := scraper.NewRetryScraper(mockInner, 3, 1*time.Millisecond)
		meta, err := retryS.SearchByHash(ctx, query)

		assert.NoError(t, err)
		assert.Equal(t, expectedMeta, meta)
		mockInner.AssertExpectations(t)
	})

	t.Run("success on second try with backoff", func(t *testing.T) {
		expectedMeta := &scraper.Metadata{Name: "Mario"}
		mockInner.On("SearchByHash", ctx, query).Return(nil, errors.New("temporary error")).Once()
		mockInner.On("SearchByHash", ctx, query).Return(expectedMeta, nil).Once()
		mockInner.On("Name").Return("Mock").Maybe()

		retryS := scraper.NewRetryScraper(mockInner, 3, 1*time.Millisecond)
		start := time.Now()
		meta, err := retryS.SearchByHash(ctx, query)
		duration := time.Since(start)

		assert.NoError(t, err)
		assert.Equal(t, expectedMeta, meta)
		assert.GreaterOrEqual(t, duration, 1*time.Millisecond)
		mockInner.AssertExpectations(t)
	})

	t.Run("fails after max retries", func(t *testing.T) {
		mockInner.On("SearchByHash", ctx, query).Return(nil, errors.New("persistent error")).Times(4)
		mockInner.On("Name").Return("Mock").Maybe()

		retryS := scraper.NewRetryScraper(mockInner, 3, 1*time.Millisecond)
		meta, err := retryS.SearchByHash(ctx, query)

		assert.Error(t, err)
		assert.Nil(t, meta)
		assert.Contains(t, err.Error(), "máximos reintentos alcanzados")
		mockInner.AssertExpectations(t)
	})

	t.Run("context cancellation", func(t *testing.T) {
		cancelCtx, cancel := context.WithCancel(ctx)
		mockInner.On("SearchByHash", cancelCtx, query).Return(nil, errors.New("error")).Once()
		mockInner.On("Name").Return("Mock").Maybe()

		retryS := scraper.NewRetryScraper(mockInner, 3, 100*time.Millisecond)

		go func() {
			time.Sleep(10 * time.Millisecond)
			cancel()
		}()

		meta, err := retryS.SearchByHash(cancelCtx, query)

		assert.Error(t, err)
		assert.True(t, errors.Is(err, context.Canceled))
		assert.Nil(t, meta)
	})
}

func TestRetryScraperMetadata(t *testing.T) {
	mockInner := mocks.NewScraper(t)
	retryS := scraper.NewRetryScraper(mockInner, 3, 1*time.Millisecond)

	mockInner.On("Name").Return("MockName")
	assert.Equal(t, "Retry(MockName)", retryS.Name())

	mockInner.On("CanSearchByHash").Return(true)
	assert.True(t, retryS.CanSearchByHash())

	mockInner.On("CanSearchByName").Return(false)
	assert.False(t, retryS.CanSearchByName())

	ctx := context.Background()
	query := scraper.SearchQuery{Filename: "Mario"}
	mockInner.On("SearchByName", ctx, query).Return(&scraper.Metadata{Name: "Mario"}, nil)
	meta, err := retryS.SearchByName(ctx, query)
	assert.NoError(t, err)
	assert.Equal(t, "Mario", meta.Name)
}

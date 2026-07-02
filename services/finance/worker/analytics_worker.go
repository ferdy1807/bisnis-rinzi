package worker

import (
	"bisnis-rinzi/packages/backend/logger"
	"bisnis-rinzi/services/finance/repository"
	"context"
	"time"
)

func StartAnalyticsAggregator(ctx context.Context, repo repository.FinanceRepository, interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	logger.Info("Background Analytics Aggregator berjalan dengan interval %s", interval.String())

	for {
		select {
		case <-ctx.Done():
			logger.Info("Background Analytics Aggregator dihentikan.")
			return
		case <-ticker.C:
			// Eksekusi Refresh
			errMonthly := repo.RefreshMonthlyAnalytics(ctx)
			if errMonthly != nil {
				logger.Error("Gagal refresh Monthly Analytics: %v", errMonthly)
			}

			errProduct := repo.RefreshProductAnalytics(ctx)
			if errProduct != nil {
				logger.Error("Gagal refresh Product Analytics: %v", errProduct)
			}

			errDaily := repo.RefreshDailyClosings(ctx)
			if errDaily != nil {
				logger.Error("Gagal sinkronisasi data Daily Closings: %v", errDaily)
			}

			if errMonthly == nil && errProduct == nil && errDaily == nil {
				// Matikan log sukses jika tidak ingin spamming console tiap menit.
				// logger.Info("Berhasil refresh tabel agregasi analytics di background.")
			}
		}
	}
}

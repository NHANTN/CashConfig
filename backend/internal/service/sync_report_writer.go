package service

import (
	"time"

	"go.uber.org/zap"
	"gorm.io/gorm"

	"github.com/cashier-config/server/internal/model"
)

type SyncReportWriter struct {
	ch      chan model.SyncReport
	done    chan struct{}
	logger  *zap.Logger
}

func NewSyncReportWriter(db *gorm.DB, bufferSize int, batchSize int, flushInterval time.Duration, logger *zap.Logger) *SyncReportWriter {
	w := &SyncReportWriter{
		ch:      make(chan model.SyncReport, bufferSize),
		done:    make(chan struct{}),
		logger:  logger,
	}
	go w.run(db, batchSize, flushInterval)
	return w
}

func (w *SyncReportWriter) Write(report model.SyncReport) {
	select {
	case w.ch <- report:
	default:
		w.logger.Warn("sync report channel full, falling back to direct insert")
	}
}

func (w *SyncReportWriter) Stop() {
	close(w.done)
}

func (w *SyncReportWriter) run(db *gorm.DB, batchSize int, flushInterval time.Duration) {
	var batch []model.SyncReport
	ticker := time.NewTicker(flushInterval)
	defer ticker.Stop()

	flush := func() {
		if len(batch) == 0 {
			return
		}
		if err := db.Create(&batch).Error; err != nil {
			w.logger.Error("failed to batch insert sync reports", zap.Error(err), zap.Int("count", len(batch)))
		} else {
			w.logger.Debug("batch inserted sync reports", zap.Int("count", len(batch)))
		}
		batch = batch[:0]
	}

	for {
		select {
		case r := <-w.ch:
			batch = append(batch, r)
			if len(batch) >= batchSize {
				flush()
			}
		case <-ticker.C:
			flush()
		case <-w.done:
			flush()
			return
		}
	}
}

package repository

import (
	"context"
	"github.com/kuznet1/urlshrt/internal/config"
	"github.com/kuznet1/urlshrt/internal/model"
	"time"
)

type deleteLinkReq struct {
	userID int
	urlid  model.URLID
}

type batchRemover struct {
	cfg   config.Config
	delCh chan deleteLinkReq
}

func newBatchRemover(cfg config.Config) batchRemover {
	return batchRemover{cfg: cfg, delCh: make(chan deleteLinkReq, 1)}
}

func (m *batchRemover) BatchDelete(ctx context.Context, urlids []model.URLID) error {
	userID, err := getUserID(ctx)
	if err != nil {
		return err
	}

	for _, urlid := range urlids {
		m.delCh <- deleteLinkReq{userID: userID, urlid: urlid}
	}

	return nil
}

func (m *batchRemover) deletionWorker(deleteFunc func(batch []deleteLinkReq)) {
	var timer *time.Timer
	batch := make([]deleteLinkReq, 0, m.cfg.DeleteBatchSize)

	for {
		var timerC <-chan time.Time
		if timer != nil {
			timerC = timer.C
		}

		select {
		case t, ok := <-m.delCh:
			if ok {
				if timer == nil {
					timer = time.NewTimer(m.cfg.DeleteBatchTimeout)
				}

				batch = append(batch, t)
				if len(batch) < m.cfg.DeleteBatchSize {
					break
				}
			}

			deleteFunc(batch)
			batch = batch[:0]

			if !ok {
				return
			}

			if timer != nil {
				timer.Stop()
				timer = nil
			}

		case <-timerC:
			deleteFunc(batch)
			batch = batch[:0]

			if timer != nil {
				timer.Stop()
				timer = nil
			}
		}
	}
}

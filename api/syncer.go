package api

import (
	"context"
	"github.com/Filecoin-Titan/titan/api/types"
	"github.com/gnasnik/titan-container-api/core/dao"
	"time"
)

type Syncer struct {
}

func NewSyncer() *Syncer {
	return &Syncer{}
}

func (s *Syncer) run(ctx context.Context) {
	interval := 10 * time.Minute
	ticker := time.NewTicker(interval)

	for {
		select {
		case <-ticker.C:
			s.updateProviders(ctx)
		case <-ctx.Done():
			return
		}
	}
}

func (s *Syncer) updateProviders(ctx context.Context) {
	for _, scheduler := range GlobalServer.GetSchedulers() {
		response, err := scheduler.Api.GetProviders(ctx, &types.GetProviderOption{Page: 1, Size: 100})
		if err != nil {
			log.Errorf("get provider: %v", err)
			continue
		}

		res, err := syncQueryResource(ctx, scheduler, response)
		if err != nil {
			log.Errorf("syncQueryResource: %v", err)
		}

		err = dao.AddProviderWithResource(ctx, res)
		if err != nil {
			log.Errorf("add provider: %v", err)
		}

	}

}

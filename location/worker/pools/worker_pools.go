package pools

import (
	"context"
	"sync"

	"github.com/Sharykhin/go-delivery-dymas/location/domain"
)

// LocationWorkerPools WorkerLocationPools Add count tasks in courierLocationQueue for handling these tasks and run count workers countWorkers
// LocationWorkerPools It needs when we have a lot of requests.
type LocationWorkerPool struct {
	courierLocationQueue chan *domain.CourierLocation
	courierService       domain.CourierLocationServiceInterface
	onceInit             sync.Once
	countTasks           int
	countWorkers         int
}

// Init inits workerPools define count task and count workers.
func (wl *LocationWorkerPools) Init(ctx context.Context) {
	wl.onceInit.Do(func() {
		wl.courierLocationQueue = make(chan domain.CourierLocation, wl.countTasks)

		for wl.countWorkers > 0 {
			go wl.handleTasks(ctx)
			wl.countWorkers--
		}
	})
}

func (wl *LocationWorkerPools) handleTasks(ctx context.Context) {
	for {
		courierLocation := <-wl.courierLocationQueue
		wl.courierService.SaveLatestCourierLocation(ctx, courierLocation)
	}
}

// AddTask adds task in LocationQueue.
func (wl *LocationWorkerPools) AddTask(courierLocation *domain.CourierLocation) {
	wl.courierLocationQueue <- courierLocation
}

// NewWorkerPools creates new worker pools.
func NewWorkerPools(
	courierLocationService domain.CourierLocationServiceInterface,
	countWorkers int,
	countTasks int,
) *LocationWorkerPools {
	return &LocationWorkerPools{
		courierService: courierLocationService,
		countWorkers:   countWorkers,
		countTasks:     countTasks,
	}
}

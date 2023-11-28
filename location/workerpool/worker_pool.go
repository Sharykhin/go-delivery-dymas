package workerpool

import (
	"context"
	"sync"

	"github.com/Sharykhin/go-delivery-dymas/location/domain"
)

// LocationPool add count tasks in courierLocationQueue for handling these tasks and run count workers countWorkers
// It needs when we have a lot of requests.
type LocationPool struct {
	courierLocationQueue chan domain.CourierLocation
	courierService       domain.CourierLocationServiceInterface
	onceInit             sync.Once
	countTasks           int
	countWorkers         int
}

// Init inits workerPools define count task and count workers.
func (wl *LocationPool) Init() {
	ctx := context.Background()
	wl.onceInit.Do(func() {
		wl.courierLocationQueue = make(chan domain.CourierLocation, wl.countTasks)
		i := 0
		for i < wl.countWorkers {
			go wl.handleTasks(ctx)
			i++
		}
	})
}

func (wl *LocationPool) handleTasks(ctx context.Context) {
	for {
		courierLocation := <-wl.courierLocationQueue
		wl.courierService.SaveLatestCourierLocation(ctx, &courierLocation)
	}
}

// AddTask adds task in LocationQueue.
func (wl *LocationPool) AddTask(courierLocation domain.CourierLocation) {
	wl.courierLocationQueue <- courierLocation
}

// NewWorkerPools creates new worker pools.
func NewWorkerPools(
	courierLocationService domain.CourierLocationServiceInterface,
	countWorkers int,
	countTasks int,
) *LocationPool {
	return &LocationPool{
		courierService: courierLocationService,
		countWorkers:   countWorkers,
		countTasks:     countTasks,
	}
}

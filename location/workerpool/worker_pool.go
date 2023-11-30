package workerpool

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/Sharykhin/go-delivery-dymas/location/domain"
)

// LocationPool add count tasks in courierLocationQueue for handling these tasks and run count workers countWorkers
// It needs when we have a lot of requests.
type LocationPool struct {
	courierLocationQueue chan domain.CourierLocation
	courierService       domain.CourierLocationServiceInterface
	countTasks           int
	countWorkers         int
}

// Init inits workerPools define count task and count workers.
func (wl *LocationPool) Init() {
	ctx := context.Background()
	wl.courierLocationQueue = make(chan domain.CourierLocation, wl.countTasks)
	i := 0
	for i < wl.countWorkers {
		go wl.handleTasks(ctx)
		i++
		if i == 1 {
			go wl.gracefulShutdown(ctx)
		}
	}
}

func (wl *LocationPool) handleTasks(ctx context.Context) {
	for courierLocation := range wl.courierLocationQueue {
		err := wl.courierService.SaveLatestCourierLocation(ctx, &courierLocation)
		if err != nil {
			log.Printf("failed to save latest position: %v\n", err)
		}
	}
}

func (wl *LocationPool) gracefulShutdown(ctx context.Context) {
	ctx, stop := signal.NotifyContext(ctx, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	<-ctx.Done()
	defer stop()
	close(wl.courierLocationQueue)

	fmt.Println("Stop Worker Pool")
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

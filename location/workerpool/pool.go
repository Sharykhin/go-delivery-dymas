package workerpool

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/Sharykhin/go-delivery-dymas/location/domain"
)

// LocationPool add count tasks in courierLocationQueue for handling these tasks and run count workers countWorkers
// It needs when we have a lot of requests.
type LocationPool struct {
	courierLocationQueue chan *domain.CourierLocation
	courierService       domain.CourierLocationServiceInterface
	countTasks           int
	countWorkers         int
}

// Init inits workerPools define count task and count workers.
func (wl *LocationPool) Init(ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()
	wl.courierLocationQueue = make(chan *domain.CourierLocation, wl.countTasks)
	chanTimeout := make(chan bool)
	var wgWorkerPool sync.WaitGroup
	for i := 0; i < wl.countWorkers; i++ {
		go wl.handleTasks(chanTimeout)
	}

	wgWorkerPool.Add(1)
	go wl.gracefulShutdown(ctx, &wgWorkerPool)
	wgWorkerPool.Wait()
	close(chanTimeout)
}

func (wl *LocationPool) handleTasks(timeoutSignal <-chan bool) {
	ctx := context.Background()
	for courierLocation := range wl.courierLocationQueue {
		select {
		case <-timeoutSignal:
			fmt.Println("Worker was stopped")
			break
		default:
			err := wl.courierService.SaveLatestCourierLocation(ctx, courierLocation)
			if err != nil {
				log.Printf("failed to save latest position: %v\n", err)
			}
		}
	}
}

func (wl *LocationPool) gracefulShutdown(ctx context.Context, wg *sync.WaitGroup) {
	<-ctx.Done()
	close(wl.courierLocationQueue)
	timeoutCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	exitGraceful := true
	for exitGraceful {
		select {
		case <-timeoutCtx.Done():
			exitGraceful = false
			fmt.Println("Exit with timeout")
		default:
			if len(wl.courierLocationQueue) == 0 {
				fmt.Println("Exit With empty task", len(wl.courierLocationQueue))
				exitGraceful = false
			}
		}
	}
	wg.Done()
	fmt.Println("Stop Worker Pool")
}

// AddTask adds task in LocationQueue.
func (wl *LocationPool) AddTask(courierLocation *domain.CourierLocation) {
	wl.courierLocationQueue <- courierLocation
}

// NewLocationPool creates new worker pools.
func NewLocationPool(
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

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
	courierLocationQueue    chan *domain.CourierLocation
	courierService          domain.CourierLocationServiceInterface
	countTasks              int
	countWorkers            int
	timeoutGracefulShutdown time.Duration
}

// Init inits workerPools define count task and count workers.
func (wl *LocationPool) Init(ctxSignalShutdown context.Context, wg *sync.WaitGroup) {
	defer wg.Done()
	wl.courierLocationQueue = make(chan *domain.CourierLocation, wl.countTasks)
	cancelCtx, cancel := context.WithCancel(context.Background())
	defer cancel()
	fmt.Println(wl.countTasks)
	fmt.Println("her")
	var wgWorkerPool sync.WaitGroup

	for i := 0; i < wl.countWorkers; i++ {
		go wl.handleTasks(cancelCtx)
	}

	<-ctxSignalShutdown.Done()

	close(wl.courierLocationQueue)
	wl.gracefulShutdown(&wgWorkerPool)
}

func (wl *LocationPool) handleTasks(ctxCancel context.Context) {
	ctx := context.Background()
	for courierLocation := range wl.courierLocationQueue {
		select {
		case <-ctxCancel.Done():
			return
		default:
			err := wl.courierService.SaveLatestCourierLocation(ctx, courierLocation)
			if err != nil {
				log.Printf("failed to save latest position: %v\n", err)
			}
		}
	}
}

func (wl *LocationPool) gracefulShutdown(wg *sync.WaitGroup) {
	timeoutCtx, cancel := context.WithTimeout(context.Background(), wl.timeoutGracefulShutdown*time.Second)
	defer cancel()
	for {
		select {
		case <-timeoutCtx.Done():
			return
		default:

			if len(wl.courierLocationQueue) == 0 {
				return
			}

		}
	}
}

// AddTask adds task in LocationQueue.
func (wl *LocationPool) AddTask(courierLocation *domain.CourierLocation) {
	fmt.Println("tsak was added")
	wl.courierLocationQueue <- courierLocation
	fmt.Println(len(wl.courierLocationQueue))
}

// NewLocationPool creates new worker pools.
func NewLocationPool(
	courierLocationService domain.CourierLocationServiceInterface,
	countWorkers int,
	countTasks int,
	timeoutGracefulShutdown time.Duration,

) *LocationPool {
	return &LocationPool{
		courierService:          courierLocationService,
		countWorkers:            countWorkers,
		countTasks:              countTasks,
		timeoutGracefulShutdown: timeoutGracefulShutdown,
	}
}

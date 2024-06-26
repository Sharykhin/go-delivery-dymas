// Code generated by http://github.com/gojuno/minimock (dev). DO NOT EDIT.

package mock

//go:generate minimock -i github.com/Sharykhin/go-delivery-dymas/location/domain.CourierLocationWorkerPool -o courier_location_worker_pool_mock_test.go -n CourierLocationWorkerPoolMock -p mock

import (
	"sync"
	mm_atomic "sync/atomic"
	mm_time "time"

	mm_domain "github.com/Sharykhin/go-delivery-dymas/location/domain"
	"github.com/gojuno/minimock/v3"
)

// CourierLocationWorkerPoolMock implements domain.CourierLocationWorkerPool
type CourierLocationWorkerPoolMock struct {
	t          minimock.Tester
	finishOnce sync.Once

	funcAddTask          func(courierLocation *mm_domain.CourierLocation)
	inspectFuncAddTask   func(courierLocation *mm_domain.CourierLocation)
	afterAddTaskCounter  uint64
	beforeAddTaskCounter uint64
	AddTaskMock          mCourierLocationWorkerPoolMockAddTask
}

// NewCourierLocationWorkerPoolMock returns a mock for domain.CourierLocationWorkerPool
func NewCourierLocationWorkerPoolMock(t minimock.Tester) *CourierLocationWorkerPoolMock {
	m := &CourierLocationWorkerPoolMock{t: t}

	if controller, ok := t.(minimock.MockController); ok {
		controller.RegisterMocker(m)
	}

	m.AddTaskMock = mCourierLocationWorkerPoolMockAddTask{mock: m}
	m.AddTaskMock.callArgs = []*CourierLocationWorkerPoolMockAddTaskParams{}

	t.Cleanup(m.MinimockFinish)

	return m
}

type mCourierLocationWorkerPoolMockAddTask struct {
	mock               *CourierLocationWorkerPoolMock
	defaultExpectation *CourierLocationWorkerPoolMockAddTaskExpectation
	expectations       []*CourierLocationWorkerPoolMockAddTaskExpectation

	callArgs []*CourierLocationWorkerPoolMockAddTaskParams
	mutex    sync.RWMutex
}

// CourierLocationWorkerPoolMockAddTaskExpectation specifies expectation struct of the CourierLocationWorkerPool.AddTask
type CourierLocationWorkerPoolMockAddTaskExpectation struct {
	mock      *CourierLocationWorkerPoolMock
	params    *CourierLocationWorkerPoolMockAddTaskParams
	paramPtrs *CourierLocationWorkerPoolMockAddTaskParamPtrs

	Counter uint64
}

// CourierLocationWorkerPoolMockAddTaskParams contains parameters of the CourierLocationWorkerPool.AddTask
type CourierLocationWorkerPoolMockAddTaskParams struct {
	courierLocation *mm_domain.CourierLocation
}

// CourierLocationWorkerPoolMockAddTaskParamPtrs contains pointers to parameters of the CourierLocationWorkerPool.AddTask
type CourierLocationWorkerPoolMockAddTaskParamPtrs struct {
	courierLocation **mm_domain.CourierLocation
}

// Expect sets up expected params for CourierLocationWorkerPool.AddTask
func (mmAddTask *mCourierLocationWorkerPoolMockAddTask) Expect(courierLocation *mm_domain.CourierLocation) *mCourierLocationWorkerPoolMockAddTask {
	if mmAddTask.mock.funcAddTask != nil {
		mmAddTask.mock.t.Fatalf("CourierLocationWorkerPoolMock.AddTask mock is already set by Set")
	}

	if mmAddTask.defaultExpectation == nil {
		mmAddTask.defaultExpectation = &CourierLocationWorkerPoolMockAddTaskExpectation{}
	}

	if mmAddTask.defaultExpectation.paramPtrs != nil {
		mmAddTask.mock.t.Fatalf("CourierLocationWorkerPoolMock.AddTask mock is already set by ExpectParams functions")
	}

	mmAddTask.defaultExpectation.params = &CourierLocationWorkerPoolMockAddTaskParams{courierLocation}
	for _, e := range mmAddTask.expectations {
		if minimock.Equal(e.params, mmAddTask.defaultExpectation.params) {
			mmAddTask.mock.t.Fatalf("Expectation set by When has same params: %#v", *mmAddTask.defaultExpectation.params)
		}
	}

	return mmAddTask
}

// ExpectCourierLocationParam1 sets up expected param courierLocation for CourierLocationWorkerPool.AddTask
func (mmAddTask *mCourierLocationWorkerPoolMockAddTask) ExpectCourierLocationParam1(courierLocation *mm_domain.CourierLocation) *mCourierLocationWorkerPoolMockAddTask {
	if mmAddTask.mock.funcAddTask != nil {
		mmAddTask.mock.t.Fatalf("CourierLocationWorkerPoolMock.AddTask mock is already set by Set")
	}

	if mmAddTask.defaultExpectation == nil {
		mmAddTask.defaultExpectation = &CourierLocationWorkerPoolMockAddTaskExpectation{}
	}

	if mmAddTask.defaultExpectation.params != nil {
		mmAddTask.mock.t.Fatalf("CourierLocationWorkerPoolMock.AddTask mock is already set by Expect")
	}

	if mmAddTask.defaultExpectation.paramPtrs == nil {
		mmAddTask.defaultExpectation.paramPtrs = &CourierLocationWorkerPoolMockAddTaskParamPtrs{}
	}
	mmAddTask.defaultExpectation.paramPtrs.courierLocation = &courierLocation

	return mmAddTask
}

// Inspect accepts an inspector function that has same arguments as the CourierLocationWorkerPool.AddTask
func (mmAddTask *mCourierLocationWorkerPoolMockAddTask) Inspect(f func(courierLocation *mm_domain.CourierLocation)) *mCourierLocationWorkerPoolMockAddTask {
	if mmAddTask.mock.inspectFuncAddTask != nil {
		mmAddTask.mock.t.Fatalf("Inspect function is already set for CourierLocationWorkerPoolMock.AddTask")
	}

	mmAddTask.mock.inspectFuncAddTask = f

	return mmAddTask
}

// Return sets up results that will be returned by CourierLocationWorkerPool.AddTask
func (mmAddTask *mCourierLocationWorkerPoolMockAddTask) Return() *CourierLocationWorkerPoolMock {
	if mmAddTask.mock.funcAddTask != nil {
		mmAddTask.mock.t.Fatalf("CourierLocationWorkerPoolMock.AddTask mock is already set by Set")
	}

	if mmAddTask.defaultExpectation == nil {
		mmAddTask.defaultExpectation = &CourierLocationWorkerPoolMockAddTaskExpectation{mock: mmAddTask.mock}
	}

	return mmAddTask.mock
}

// Set uses given function f to mock the CourierLocationWorkerPool.AddTask method
func (mmAddTask *mCourierLocationWorkerPoolMockAddTask) Set(f func(courierLocation *mm_domain.CourierLocation)) *CourierLocationWorkerPoolMock {
	if mmAddTask.defaultExpectation != nil {
		mmAddTask.mock.t.Fatalf("Default expectation is already set for the CourierLocationWorkerPool.AddTask method")
	}

	if len(mmAddTask.expectations) > 0 {
		mmAddTask.mock.t.Fatalf("Some expectations are already set for the CourierLocationWorkerPool.AddTask method")
	}

	mmAddTask.mock.funcAddTask = f
	return mmAddTask.mock
}

// AddTask implements domain.CourierLocationWorkerPool
func (mmAddTask *CourierLocationWorkerPoolMock) AddTask(courierLocation *mm_domain.CourierLocation) {
	mm_atomic.AddUint64(&mmAddTask.beforeAddTaskCounter, 1)
	defer mm_atomic.AddUint64(&mmAddTask.afterAddTaskCounter, 1)

	if mmAddTask.inspectFuncAddTask != nil {
		mmAddTask.inspectFuncAddTask(courierLocation)
	}

	mm_params := CourierLocationWorkerPoolMockAddTaskParams{courierLocation}

	// Record call args
	mmAddTask.AddTaskMock.mutex.Lock()
	mmAddTask.AddTaskMock.callArgs = append(mmAddTask.AddTaskMock.callArgs, &mm_params)
	mmAddTask.AddTaskMock.mutex.Unlock()

	for _, e := range mmAddTask.AddTaskMock.expectations {
		if minimock.Equal(*e.params, mm_params) {
			mm_atomic.AddUint64(&e.Counter, 1)
			return
		}
	}

	if mmAddTask.AddTaskMock.defaultExpectation != nil {
		mm_atomic.AddUint64(&mmAddTask.AddTaskMock.defaultExpectation.Counter, 1)
		mm_want := mmAddTask.AddTaskMock.defaultExpectation.params
		mm_want_ptrs := mmAddTask.AddTaskMock.defaultExpectation.paramPtrs

		mm_got := CourierLocationWorkerPoolMockAddTaskParams{courierLocation}

		if mm_want_ptrs != nil {

			if mm_want_ptrs.courierLocation != nil && !minimock.Equal(*mm_want_ptrs.courierLocation, mm_got.courierLocation) {
				mmAddTask.t.Errorf("CourierLocationWorkerPoolMock.AddTask got unexpected parameter courierLocation, want: %#v, got: %#v%s\n", *mm_want_ptrs.courierLocation, mm_got.courierLocation, minimock.Diff(*mm_want_ptrs.courierLocation, mm_got.courierLocation))
			}

		} else if mm_want != nil && !minimock.Equal(*mm_want, mm_got) {
			mmAddTask.t.Errorf("CourierLocationWorkerPoolMock.AddTask got unexpected parameters, want: %#v, got: %#v%s\n", *mm_want, mm_got, minimock.Diff(*mm_want, mm_got))
		}

		return

	}
	if mmAddTask.funcAddTask != nil {
		mmAddTask.funcAddTask(courierLocation)
		return
	}
	mmAddTask.t.Fatalf("Unexpected call to CourierLocationWorkerPoolMock.AddTask. %v", courierLocation)

}

// AddTaskAfterCounter returns a count of finished CourierLocationWorkerPoolMock.AddTask invocations
func (mmAddTask *CourierLocationWorkerPoolMock) AddTaskAfterCounter() uint64 {
	return mm_atomic.LoadUint64(&mmAddTask.afterAddTaskCounter)
}

// AddTaskBeforeCounter returns a count of CourierLocationWorkerPoolMock.AddTask invocations
func (mmAddTask *CourierLocationWorkerPoolMock) AddTaskBeforeCounter() uint64 {
	return mm_atomic.LoadUint64(&mmAddTask.beforeAddTaskCounter)
}

// Calls returns a list of arguments used in each call to CourierLocationWorkerPoolMock.AddTask.
// The list is in the same order as the calls were made (i.e. recent calls have a higher index)
func (mmAddTask *mCourierLocationWorkerPoolMockAddTask) Calls() []*CourierLocationWorkerPoolMockAddTaskParams {
	mmAddTask.mutex.RLock()

	argCopy := make([]*CourierLocationWorkerPoolMockAddTaskParams, len(mmAddTask.callArgs))
	copy(argCopy, mmAddTask.callArgs)

	mmAddTask.mutex.RUnlock()

	return argCopy
}

// MinimockAddTaskDone returns true if the count of the AddTask invocations corresponds
// the number of defined expectations
func (m *CourierLocationWorkerPoolMock) MinimockAddTaskDone() bool {
	for _, e := range m.AddTaskMock.expectations {
		if mm_atomic.LoadUint64(&e.Counter) < 1 {
			return false
		}
	}

	// if default expectation was set then invocations count should be greater than zero
	if m.AddTaskMock.defaultExpectation != nil && mm_atomic.LoadUint64(&m.afterAddTaskCounter) < 1 {
		return false
	}
	// if func was set then invocations count should be greater than zero
	if m.funcAddTask != nil && mm_atomic.LoadUint64(&m.afterAddTaskCounter) < 1 {
		return false
	}
	return true
}

// MinimockAddTaskInspect logs each unmet expectation
func (m *CourierLocationWorkerPoolMock) MinimockAddTaskInspect() {
	for _, e := range m.AddTaskMock.expectations {
		if mm_atomic.LoadUint64(&e.Counter) < 1 {
			m.t.Errorf("Expected call to CourierLocationWorkerPoolMock.AddTask with params: %#v", *e.params)
		}
	}

	// if default expectation was set then invocations count should be greater than zero
	if m.AddTaskMock.defaultExpectation != nil && mm_atomic.LoadUint64(&m.afterAddTaskCounter) < 1 {
		if m.AddTaskMock.defaultExpectation.params == nil {
			m.t.Error("Expected call to CourierLocationWorkerPoolMock.AddTask")
		} else {
			m.t.Errorf("Expected call to CourierLocationWorkerPoolMock.AddTask with params: %#v", *m.AddTaskMock.defaultExpectation.params)
		}
	}
	// if func was set then invocations count should be greater than zero
	if m.funcAddTask != nil && mm_atomic.LoadUint64(&m.afterAddTaskCounter) < 1 {
		m.t.Error("Expected call to CourierLocationWorkerPoolMock.AddTask")
	}
}

// MinimockFinish checks that all mocked methods have been called the expected number of times
func (m *CourierLocationWorkerPoolMock) MinimockFinish() {
	m.finishOnce.Do(func() {
		if !m.minimockDone() {
			m.MinimockAddTaskInspect()
			m.t.FailNow()
		}
	})
}

// MinimockWait waits for all mocked methods to be called the expected number of times
func (m *CourierLocationWorkerPoolMock) MinimockWait(timeout mm_time.Duration) {
	timeoutCh := mm_time.After(timeout)
	for {
		if m.minimockDone() {
			return
		}
		select {
		case <-timeoutCh:
			m.MinimockFinish()
			return
		case <-mm_time.After(10 * mm_time.Millisecond):
		}
	}
}

func (m *CourierLocationWorkerPoolMock) minimockDone() bool {
	done := true
	return done &&
		m.MinimockAddTaskDone()
}

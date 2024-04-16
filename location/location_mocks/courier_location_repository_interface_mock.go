// Code generated by http://github.com/gojuno/minimock (dev). DO NOT EDIT.

package location_mocks

//go:generate minimock -i github.com/Sharykhin/go-delivery-dymas/location/domain.CourierLocationRepositoryInterface -o courier_location_repository_interface_mock_test.go -n CourierLocationRepositoryInterfaceMock -p location_mocks

import (
	"context"
	"sync"
	mm_atomic "sync/atomic"
	mm_time "time"

	mm_domain "github.com/Sharykhin/go-delivery-dymas/location/domain"
	"github.com/gojuno/minimock/v3"
)

// CourierLocationRepositoryInterfaceMock implements domain.CourierLocationRepositoryInterface
type CourierLocationRepositoryInterfaceMock struct {
	t          minimock.Tester
	finishOnce sync.Once

	funcSaveLatestCourierGeoPosition          func(ctx context.Context, courierLocation *mm_domain.CourierLocation) (err error)
	inspectFuncSaveLatestCourierGeoPosition   func(ctx context.Context, courierLocation *mm_domain.CourierLocation)
	afterSaveLatestCourierGeoPositionCounter  uint64
	beforeSaveLatestCourierGeoPositionCounter uint64
	SaveLatestCourierGeoPositionMock          mCourierLocationRepositoryInterfaceMockSaveLatestCourierGeoPosition
}

// NewCourierLocationRepositoryInterfaceMock returns a mock for domain.CourierLocationRepositoryInterface
func NewCourierLocationRepositoryInterfaceMock(t minimock.Tester) *CourierLocationRepositoryInterfaceMock {
	m := &CourierLocationRepositoryInterfaceMock{t: t}

	if controller, ok := t.(minimock.MockController); ok {
		controller.RegisterMocker(m)
	}

	m.SaveLatestCourierGeoPositionMock = mCourierLocationRepositoryInterfaceMockSaveLatestCourierGeoPosition{mock: m}
	m.SaveLatestCourierGeoPositionMock.callArgs = []*CourierLocationRepositoryInterfaceMockSaveLatestCourierGeoPositionParams{}

	t.Cleanup(m.MinimockFinish)

	return m
}

type mCourierLocationRepositoryInterfaceMockSaveLatestCourierGeoPosition struct {
	mock               *CourierLocationRepositoryInterfaceMock
	defaultExpectation *CourierLocationRepositoryInterfaceMockSaveLatestCourierGeoPositionExpectation
	expectations       []*CourierLocationRepositoryInterfaceMockSaveLatestCourierGeoPositionExpectation

	callArgs []*CourierLocationRepositoryInterfaceMockSaveLatestCourierGeoPositionParams
	mutex    sync.RWMutex
}

// CourierLocationRepositoryInterfaceMockSaveLatestCourierGeoPositionExpectation specifies expectation struct of the CourierLocationRepositoryInterface.SaveLatestCourierGeoPosition
type CourierLocationRepositoryInterfaceMockSaveLatestCourierGeoPositionExpectation struct {
	mock    *CourierLocationRepositoryInterfaceMock
	params  *CourierLocationRepositoryInterfaceMockSaveLatestCourierGeoPositionParams
	results *CourierLocationRepositoryInterfaceMockSaveLatestCourierGeoPositionResults
	Counter uint64
}

// CourierLocationRepositoryInterfaceMockSaveLatestCourierGeoPositionParams contains parameters of the CourierLocationRepositoryInterface.SaveLatestCourierGeoPosition
type CourierLocationRepositoryInterfaceMockSaveLatestCourierGeoPositionParams struct {
	ctx             context.Context
	courierLocation *mm_domain.CourierLocation
}

// CourierLocationRepositoryInterfaceMockSaveLatestCourierGeoPositionResults contains results of the CourierLocationRepositoryInterface.SaveLatestCourierGeoPosition
type CourierLocationRepositoryInterfaceMockSaveLatestCourierGeoPositionResults struct {
	err error
}

// Expect sets up expected params for CourierLocationRepositoryInterface.SaveLatestCourierGeoPosition
func (mmSaveLatestCourierGeoPosition *mCourierLocationRepositoryInterfaceMockSaveLatestCourierGeoPosition) Expect(ctx context.Context, courierLocation *mm_domain.CourierLocation) *mCourierLocationRepositoryInterfaceMockSaveLatestCourierGeoPosition {
	if mmSaveLatestCourierGeoPosition.mock.funcSaveLatestCourierGeoPosition != nil {
		mmSaveLatestCourierGeoPosition.mock.t.Fatalf("CourierLocationRepositoryInterfaceMock.SaveLatestCourierGeoPosition mock is already set by Set")
	}

	if mmSaveLatestCourierGeoPosition.defaultExpectation == nil {
		mmSaveLatestCourierGeoPosition.defaultExpectation = &CourierLocationRepositoryInterfaceMockSaveLatestCourierGeoPositionExpectation{}
	}

	mmSaveLatestCourierGeoPosition.defaultExpectation.params = &CourierLocationRepositoryInterfaceMockSaveLatestCourierGeoPositionParams{ctx, courierLocation}
	for _, e := range mmSaveLatestCourierGeoPosition.expectations {
		if minimock.Equal(e.params, mmSaveLatestCourierGeoPosition.defaultExpectation.params) {
			mmSaveLatestCourierGeoPosition.mock.t.Fatalf("Expectation set by When has same params: %#v", *mmSaveLatestCourierGeoPosition.defaultExpectation.params)
		}
	}

	return mmSaveLatestCourierGeoPosition
}

// Inspect accepts an inspector function that has same arguments as the CourierLocationRepositoryInterface.SaveLatestCourierGeoPosition
func (mmSaveLatestCourierGeoPosition *mCourierLocationRepositoryInterfaceMockSaveLatestCourierGeoPosition) Inspect(f func(ctx context.Context, courierLocation *mm_domain.CourierLocation)) *mCourierLocationRepositoryInterfaceMockSaveLatestCourierGeoPosition {
	if mmSaveLatestCourierGeoPosition.mock.inspectFuncSaveLatestCourierGeoPosition != nil {
		mmSaveLatestCourierGeoPosition.mock.t.Fatalf("Inspect function is already set for CourierLocationRepositoryInterfaceMock.SaveLatestCourierGeoPosition")
	}

	mmSaveLatestCourierGeoPosition.mock.inspectFuncSaveLatestCourierGeoPosition = f

	return mmSaveLatestCourierGeoPosition
}

// Return sets up results that will be returned by CourierLocationRepositoryInterface.SaveLatestCourierGeoPosition
func (mmSaveLatestCourierGeoPosition *mCourierLocationRepositoryInterfaceMockSaveLatestCourierGeoPosition) Return(err error) *CourierLocationRepositoryInterfaceMock {
	if mmSaveLatestCourierGeoPosition.mock.funcSaveLatestCourierGeoPosition != nil {
		mmSaveLatestCourierGeoPosition.mock.t.Fatalf("CourierLocationRepositoryInterfaceMock.SaveLatestCourierGeoPosition mock is already set by Set")
	}

	if mmSaveLatestCourierGeoPosition.defaultExpectation == nil {
		mmSaveLatestCourierGeoPosition.defaultExpectation = &CourierLocationRepositoryInterfaceMockSaveLatestCourierGeoPositionExpectation{mock: mmSaveLatestCourierGeoPosition.mock}
	}
	mmSaveLatestCourierGeoPosition.defaultExpectation.results = &CourierLocationRepositoryInterfaceMockSaveLatestCourierGeoPositionResults{err}
	return mmSaveLatestCourierGeoPosition.mock
}

// Set uses given function f to mock the CourierLocationRepositoryInterface.SaveLatestCourierGeoPosition method
func (mmSaveLatestCourierGeoPosition *mCourierLocationRepositoryInterfaceMockSaveLatestCourierGeoPosition) Set(f func(ctx context.Context, courierLocation *mm_domain.CourierLocation) (err error)) *CourierLocationRepositoryInterfaceMock {
	if mmSaveLatestCourierGeoPosition.defaultExpectation != nil {
		mmSaveLatestCourierGeoPosition.mock.t.Fatalf("Default expectation is already set for the CourierLocationRepositoryInterface.SaveLatestCourierGeoPosition method")
	}

	if len(mmSaveLatestCourierGeoPosition.expectations) > 0 {
		mmSaveLatestCourierGeoPosition.mock.t.Fatalf("Some expectations are already set for the CourierLocationRepositoryInterface.SaveLatestCourierGeoPosition method")
	}

	mmSaveLatestCourierGeoPosition.mock.funcSaveLatestCourierGeoPosition = f
	return mmSaveLatestCourierGeoPosition.mock
}

// When sets expectation for the CourierLocationRepositoryInterface.SaveLatestCourierGeoPosition which will trigger the result defined by the following
// Then helper
func (mmSaveLatestCourierGeoPosition *mCourierLocationRepositoryInterfaceMockSaveLatestCourierGeoPosition) When(ctx context.Context, courierLocation *mm_domain.CourierLocation) *CourierLocationRepositoryInterfaceMockSaveLatestCourierGeoPositionExpectation {
	if mmSaveLatestCourierGeoPosition.mock.funcSaveLatestCourierGeoPosition != nil {
		mmSaveLatestCourierGeoPosition.mock.t.Fatalf("CourierLocationRepositoryInterfaceMock.SaveLatestCourierGeoPosition mock is already set by Set")
	}

	expectation := &CourierLocationRepositoryInterfaceMockSaveLatestCourierGeoPositionExpectation{
		mock:   mmSaveLatestCourierGeoPosition.mock,
		params: &CourierLocationRepositoryInterfaceMockSaveLatestCourierGeoPositionParams{ctx, courierLocation},
	}
	mmSaveLatestCourierGeoPosition.expectations = append(mmSaveLatestCourierGeoPosition.expectations, expectation)
	return expectation
}

// Then sets up CourierLocationRepositoryInterface.SaveLatestCourierGeoPosition return parameters for the expectation previously defined by the When method
func (e *CourierLocationRepositoryInterfaceMockSaveLatestCourierGeoPositionExpectation) Then(err error) *CourierLocationRepositoryInterfaceMock {
	e.results = &CourierLocationRepositoryInterfaceMockSaveLatestCourierGeoPositionResults{err}
	return e.mock
}

// SaveLatestCourierGeoPosition implements domain.CourierLocationRepositoryInterface
func (mmSaveLatestCourierGeoPosition *CourierLocationRepositoryInterfaceMock) SaveLatestCourierGeoPosition(ctx context.Context, courierLocation *mm_domain.CourierLocation) (err error) {
	mm_atomic.AddUint64(&mmSaveLatestCourierGeoPosition.beforeSaveLatestCourierGeoPositionCounter, 1)
	defer mm_atomic.AddUint64(&mmSaveLatestCourierGeoPosition.afterSaveLatestCourierGeoPositionCounter, 1)

	if mmSaveLatestCourierGeoPosition.inspectFuncSaveLatestCourierGeoPosition != nil {
		mmSaveLatestCourierGeoPosition.inspectFuncSaveLatestCourierGeoPosition(ctx, courierLocation)
	}

	mm_params := CourierLocationRepositoryInterfaceMockSaveLatestCourierGeoPositionParams{ctx, courierLocation}

	// Record call args
	mmSaveLatestCourierGeoPosition.SaveLatestCourierGeoPositionMock.mutex.Lock()
	mmSaveLatestCourierGeoPosition.SaveLatestCourierGeoPositionMock.callArgs = append(mmSaveLatestCourierGeoPosition.SaveLatestCourierGeoPositionMock.callArgs, &mm_params)
	mmSaveLatestCourierGeoPosition.SaveLatestCourierGeoPositionMock.mutex.Unlock()

	for _, e := range mmSaveLatestCourierGeoPosition.SaveLatestCourierGeoPositionMock.expectations {
		if minimock.Equal(*e.params, mm_params) {
			mm_atomic.AddUint64(&e.Counter, 1)
			return e.results.err
		}
	}

	if mmSaveLatestCourierGeoPosition.SaveLatestCourierGeoPositionMock.defaultExpectation != nil {
		mm_atomic.AddUint64(&mmSaveLatestCourierGeoPosition.SaveLatestCourierGeoPositionMock.defaultExpectation.Counter, 1)
		mm_want := mmSaveLatestCourierGeoPosition.SaveLatestCourierGeoPositionMock.defaultExpectation.params
		mm_got := CourierLocationRepositoryInterfaceMockSaveLatestCourierGeoPositionParams{ctx, courierLocation}
		if mm_want != nil && !minimock.Equal(*mm_want, mm_got) {
			mmSaveLatestCourierGeoPosition.t.Errorf("CourierLocationRepositoryInterfaceMock.SaveLatestCourierGeoPosition got unexpected parameters, want: %#v, got: %#v%s\n", *mm_want, mm_got, minimock.Diff(*mm_want, mm_got))
		}

		mm_results := mmSaveLatestCourierGeoPosition.SaveLatestCourierGeoPositionMock.defaultExpectation.results
		if mm_results == nil {
			mmSaveLatestCourierGeoPosition.t.Fatal("No results are set for the CourierLocationRepositoryInterfaceMock.SaveLatestCourierGeoPosition")
		}
		return (*mm_results).err
	}
	if mmSaveLatestCourierGeoPosition.funcSaveLatestCourierGeoPosition != nil {
		return mmSaveLatestCourierGeoPosition.funcSaveLatestCourierGeoPosition(ctx, courierLocation)
	}
	mmSaveLatestCourierGeoPosition.t.Fatalf("Unexpected call to CourierLocationRepositoryInterfaceMock.SaveLatestCourierGeoPosition. %v %v", ctx, courierLocation)
	return
}

// SaveLatestCourierGeoPositionAfterCounter returns a count of finished CourierLocationRepositoryInterfaceMock.SaveLatestCourierGeoPosition invocations
func (mmSaveLatestCourierGeoPosition *CourierLocationRepositoryInterfaceMock) SaveLatestCourierGeoPositionAfterCounter() uint64 {
	return mm_atomic.LoadUint64(&mmSaveLatestCourierGeoPosition.afterSaveLatestCourierGeoPositionCounter)
}

// SaveLatestCourierGeoPositionBeforeCounter returns a count of CourierLocationRepositoryInterfaceMock.SaveLatestCourierGeoPosition invocations
func (mmSaveLatestCourierGeoPosition *CourierLocationRepositoryInterfaceMock) SaveLatestCourierGeoPositionBeforeCounter() uint64 {
	return mm_atomic.LoadUint64(&mmSaveLatestCourierGeoPosition.beforeSaveLatestCourierGeoPositionCounter)
}

// Calls returns a list of arguments used in each call to CourierLocationRepositoryInterfaceMock.SaveLatestCourierGeoPosition.
// The list is in the same order as the calls were made (i.e. recent calls have a higher index)
func (mmSaveLatestCourierGeoPosition *mCourierLocationRepositoryInterfaceMockSaveLatestCourierGeoPosition) Calls() []*CourierLocationRepositoryInterfaceMockSaveLatestCourierGeoPositionParams {
	mmSaveLatestCourierGeoPosition.mutex.RLock()

	argCopy := make([]*CourierLocationRepositoryInterfaceMockSaveLatestCourierGeoPositionParams, len(mmSaveLatestCourierGeoPosition.callArgs))
	copy(argCopy, mmSaveLatestCourierGeoPosition.callArgs)

	mmSaveLatestCourierGeoPosition.mutex.RUnlock()

	return argCopy
}

// MinimockSaveLatestCourierGeoPositionDone returns true if the count of the SaveLatestCourierGeoPosition invocations corresponds
// the number of defined expectations
func (m *CourierLocationRepositoryInterfaceMock) MinimockSaveLatestCourierGeoPositionDone() bool {
	for _, e := range m.SaveLatestCourierGeoPositionMock.expectations {
		if mm_atomic.LoadUint64(&e.Counter) < 1 {
			return false
		}
	}

	// if default expectation was set then invocations count should be greater than zero
	if m.SaveLatestCourierGeoPositionMock.defaultExpectation != nil && mm_atomic.LoadUint64(&m.afterSaveLatestCourierGeoPositionCounter) < 1 {
		return false
	}
	// if func was set then invocations count should be greater than zero
	if m.funcSaveLatestCourierGeoPosition != nil && mm_atomic.LoadUint64(&m.afterSaveLatestCourierGeoPositionCounter) < 1 {
		return false
	}
	return true
}

// MinimockSaveLatestCourierGeoPositionInspect logs each unmet expectation
func (m *CourierLocationRepositoryInterfaceMock) MinimockSaveLatestCourierGeoPositionInspect() {
	for _, e := range m.SaveLatestCourierGeoPositionMock.expectations {
		if mm_atomic.LoadUint64(&e.Counter) < 1 {
			m.t.Errorf("Expected call to CourierLocationRepositoryInterfaceMock.SaveLatestCourierGeoPosition with params: %#v", *e.params)
		}
	}

	// if default expectation was set then invocations count should be greater than zero
	if m.SaveLatestCourierGeoPositionMock.defaultExpectation != nil && mm_atomic.LoadUint64(&m.afterSaveLatestCourierGeoPositionCounter) < 1 {
		if m.SaveLatestCourierGeoPositionMock.defaultExpectation.params == nil {
			m.t.Error("Expected call to CourierLocationRepositoryInterfaceMock.SaveLatestCourierGeoPosition")
		} else {
			m.t.Errorf("Expected call to CourierLocationRepositoryInterfaceMock.SaveLatestCourierGeoPosition with params: %#v", *m.SaveLatestCourierGeoPositionMock.defaultExpectation.params)
		}
	}
	// if func was set then invocations count should be greater than zero
	if m.funcSaveLatestCourierGeoPosition != nil && mm_atomic.LoadUint64(&m.afterSaveLatestCourierGeoPositionCounter) < 1 {
		m.t.Error("Expected call to CourierLocationRepositoryInterfaceMock.SaveLatestCourierGeoPosition")
	}
}

// MinimockFinish checks that all mocked methods have been called the expected number of times
func (m *CourierLocationRepositoryInterfaceMock) MinimockFinish() {
	m.finishOnce.Do(func() {
		if !m.minimockDone() {
			m.MinimockSaveLatestCourierGeoPositionInspect()
			m.t.FailNow()
		}
	})
}

// MinimockWait waits for all mocked methods to be called the expected number of times
func (m *CourierLocationRepositoryInterfaceMock) MinimockWait(timeout mm_time.Duration) {
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

func (m *CourierLocationRepositoryInterfaceMock) minimockDone() bool {
	done := true
	return done &&
		m.MinimockSaveLatestCourierGeoPositionDone()
}

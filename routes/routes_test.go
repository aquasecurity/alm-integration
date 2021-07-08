package routes

import (
	"testing"
)

func TestScheduling(t *testing.T) {
	stopCh := make(chan struct{})
	demoRoute1 := &InputRoute{}

	demoRoute1Stopped := false
	demoRoute1.StartScheduler()
	if !demoRoute1.IsSchedulerRun() {
		t.Errorf("Route 1 is not started")
	}
	go func() {
		<-demoRoute1.Scheduling
		demoRoute1Stopped = true
		stopCh <- struct{}{}
	}()
	demoRoute1.StopScheduler()
	<-stopCh

	if !demoRoute1Stopped {
		t.Errorf("Route 1 is not stopped")
	}

	demoRoute2 := &InputRoute{}
	if demoRoute2.IsSchedulerRun() {
		t.Errorf("Route 2 should not be started")
	}
	demoRoute2.StopScheduler()

}

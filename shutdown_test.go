package shutdown

import (
	"testing"
	"time"
)

func TestShutdown(t *testing.T) {
	handler := NewHandler(&blankLog{})

	service1 := serviceWithStop{}
	service2 := serviceWithStop{}
	service3 := serviceWithStop{}
	go service1.serve(100*time.Millisecond)
	go service2.serve(1000*time.Millisecond)
	go service3.serve(5000*time.Millisecond)

	handler.Add("service1", "", Init, &service1)
	handler.Add("service2", "service1", Before, &service2)
	handler.Add("service3", "service2", After, &service3)


	result := handler.debug()
	if result["service1"] < result["service2"] {
		t.Error("service1 should be bigger than service2")
	}
	if result["service3"] < result["service1"] {
		t.Error("service3 should be bigger than service1")
	}
	if result["service3"] < result["service2"] {
		t.Error("service3 should be bigger than service2")
	}
	if result["service1"] < mid {
		t.Errorf("service1 should be bigger than %d", mid)
	}
}
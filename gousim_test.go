package gousim

import (
	"testing"
	"fmt"
	// "time"
)

func TestRecordingSession(t *testing.T) {
	r := NewRecordingSession(500)
	fmt.Printf("Recording: %b\n", r.Recording)
	err := r.LogSample("AAMetric", 10 )
	if err == nil {
		t.Error("LogSample should have blown up.")
	}
	r.Start()
	// time.Sleep(5 * time.Second)
	fmt.Printf("Recording: %b\n", r.Recording)
	for i := 0; i < 400; i++ {
		fmt.Print("!")
		r.LogSample("AAMetric", 10 * i/2)
		r.LogSample("BBMetric", 10 * i/2)
		r.LogSample("CCMetric", 10 * i/2)
		r.LogSample("DDMetric", 10 * i/2)
	}
	r.Stop()
	fmt.Printf("Recording: %b\n", r.Recording)
	// time.Sleep(5 * time.Second)
	out := r.Export()
	fmt.Printf("%s \n", out)
	fmt.Println("Session Created")

}
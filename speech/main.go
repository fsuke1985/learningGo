package main

import (
	"context"
	"flag"
	"time"
	"k8s.io/klog/v2"
)

var (
    electionPort = flag.Int("electionPort", 4040,
		"Listen at this port for leader election updates. Set to zero to disable leader election")
    flushTimeout = flag.Duration("flushTimeout", 3000*time.Millisecond, "Emit any pending transcriptions after this time")
    lastIndex = 0
    pending  []string
)

func main() {
	flag.Parse()
		/*
			electionPort > 0

			then 
				creating HTTP/SERVER
			
			else
				method calling to sendAudio
		*/

    if *electionPort > 0 {

	} else {
		klog.Infof("ports is %i: call %v ", *electionPort, "sendAudio")
		sendAudio(context.Background())
	}
}

func sendAudio(ctx context.Context) {
	
	receiveChan := make(chan bool)
	for {
		select {
		case <-ctx.Done():
            klog.Infof("Context cancelled, exiting sender loop")	
			return
        case _, ok := <-receiveChan:
			if !ok { // ok is false
				klog.Info("receive channel closed, resetting stream")
				receiveChan = make(chan bool) // panic: close of closed channel
                resetIndex()
				continue
			}
		default:
			// todo
		}

		go receiveResponses(receiveChan, os.Getpid())
	}
}

func receiveResponses(receiveChan chan bool) {
	defer close(receiveChan)

    timer := time.NewTimer(*flushTimeout)
    go func() {
        <-timer.C
        flush()
    }()

    defer timer.Stop()

    for {
        if !timer.Stop() {
            return
        }
        timer.Reset(*flushTimeout)
    }
}

func resetIndex() {
    lastIndex = 0
    pending = nil
}
func flush() {}

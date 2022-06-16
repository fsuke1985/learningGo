package main

import (
	"context"
	"flag"


	"k8s.io/klog"
)

var (
	electionPort = flag.Int("Port", 4040, "Listen port default")
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
			klog.Infof("receive chanel cloesed")
			return
		case _, ok := <-receiveChan:
			if !ok { // ok is false
				klog.Info("receive chan closed")
				receiveChan = make(chan bool) // panic: close of closed channel
				//continue
			}
		default:
			// todo
		}

		go receiveResponses(receiveChan)
	}
}

func receiveResponses(receiveChan chan bool) {
	defer close(receiveChan)
}
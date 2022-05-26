package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"strings"

	"k8s.io/klog"
)

var (
	phrases = flag.String("phrases", "", "phrases")
	cancel  context.CancelFunc
	ctx     context.Context
)

func main() {

	contextPhrases := []string{}
	if *phrases != "" {
		fmt.Println("main")

		contextPhrases = strings.Split(*phrases, ",")
		klog.Infof("%d %+q", len(contextPhrases), contextPhrases)
	}

	server := &http.Server{Addr: "4000"}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel = context.WithCancel(context.Background())
		go sendAudio(ctx)
		w.WriteHeader(http.StatusOK)
	})

	klog.Infof("Starting leader election listener at port %s", "4000")

	server.ListenAndServe()
}

func sendAudio(ctx context.Context) {

	for {
		select {
		case <-ctx.Done():
			klog.Infof("exit ")
			return
		default:
		}
	}
}

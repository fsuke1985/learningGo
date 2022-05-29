package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	redis "github.com/go-redis/redis/v7"
	speechpb "google.golang.org/genproto/googleapis/cloud/speech/v1"
	"k8s.io/klog"
)

var (
	phrases = flag.String("phrases", "phrases", "phrases")
	//lint:ignore U1000 It's ok because this is just a example.

	redisHost     = flag.String("redisHost", "localhost", "")
	recoveryQueue = flag.String("recoveryQueue", "recoverq", "")
	electionPort  = flag.Int("electionPort", 4040, "")
	redisClient   *redis.Client
)

func main() {
	var ctx context.Context
	var cancel context.CancelFunc

	redisClient = redis.NewClient(&redis.Options{
		Addr:        *redisHost + ":6379",
		Password:    "", // no password set
		DB:          0,  // use default DB
		DialTimeout: 3 * time.Second,
		ReadTimeout: 4 * time.Second,
	})

	contextPhrases := []string{}
	if *phrases != "" {
		fmt.Println("main")

		contextPhrases = strings.Split(*phrases, ",")
		klog.Infof("%d %+q", len(contextPhrases), contextPhrases)
	}

	server := &http.Server{Addr: ":4000"}

	// Listen for interrupt. If so, cancel the context to stop transcriptions,
	// then shutdown the HTTP server, which will end the main thread.
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	go func() {
		<-ch
		if cancel != nil {
			cancel()
		}
		klog.Info("Received termination, stopping transciptions")
		if err := server.Shutdown(context.TODO()); err != nil {
			klog.Fatalf("Error shutting down HTTP server: %v", err)
		}
	}()

	webHandler := func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel = context.WithCancel(context.Background())

		if strings.Contains(r.URL.Path, "start") {

		}
		go sendAudio(ctx)
		w.WriteHeader(http.StatusOK)
	}

	http.HandleFunc("/", webHandler)
	klog.Infof("Starting leader election listener at port %s", "4000")

	server.ListenAndServe()
	//http.ListenAndServe(":4000", nil)
}

func sendAudio(ctx context.Context) {

	var stream speechpb.Speech_StreamingRecognizeClient
	reciveChan := make(chan bool)
	doRecovery := redisClient.Exists(*recoveryQueue).Val() > 0
	replaying := false

	for {
		select {
		case <-ctx.Done():
			klog.Infof("exit ")
			return
		case _, ok := <-reciveChan:
			if !ok && stream != nil {
				klog.Info("receive channel closed, resetting stream")
				stream = nil
				reciveChan = make(chan bool)
				continue
			}
		default:
			// proc
		}

		var result string
		var err error

		_, _ = result, err // never used workaround

		if doRecovery {
			// todo
			result, err = redisClient.RPop(*recoveryQueue).Result()

			if err == redis.Nil {
				doRecovery = false
				replaying = false
				continue

			}
			if err == nil && !replaying {
				replaying = true
			}
		} else {
			// todo
		}
	}
}

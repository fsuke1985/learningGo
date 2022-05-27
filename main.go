package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"strings"
	"time"

	redis "github.com/go-redis/redis/v7"
	speechpb "google.golang.org/genproto/googleapis/cloud/speech/v1"
	"k8s.io/klog"
)


var (
	phrases = flag.String("phrases", "", "phrases")
	//lint:ignore U1000 It's ok because this is just a example.
	cancel  context.CancelFunc
	ctx     context.Context

	redisHost = flag.String("redisHost", "localhost", "")
	recoveryQueue = flag.String("recoveryQueue", "recoverq", "")
	redisClient *redis.Client
)

func main() {

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

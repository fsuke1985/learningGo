package main

import (
	"fmt"
	"k8s.io/klog"
	"net/http"
)

var (
	err error
)

func smain() {

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println(r.Form)
		fmt.Println("path", r.URL.Scheme)
		w.WriteHeader(http.StatusOK)
	})

	err = http.ListenAndServe(":9000", nil)

	if err != nil {
		klog.Infof("ListenAndServer: ", err)
	}
}

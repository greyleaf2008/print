package main

import (
	"fmt"
	"html"
	"log"
	"net/http"
	"sync"
	"time"

	"os"
	"os/signal"
	"syscall"
)

func fooHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("Got connection:%s", r.Proto)
	w.Write([]byte("hello"))
}

const idleTimeout = 5 * time.Minute
const activeTimeout = 10 * time.Minute

func main() {
	var (
		wg sync.WaitGroup
	)

	srv := http.Server{
		Addr:           ":8888",
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	// http.Handle("/foo", fooHandler)
	http.HandleFunc("/bar", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "helloooo, %q", html.EscapeString(r.URL.Path))
	})

	wg.Add(1)
	go func() {
		wg.Done()
		if err := srv.ListenAndServe(); err != nil {
			fmt.Println("test1", err)
		} else {
			fmt.Println("test3")
		}
	}()
	wg.Wait()
	fmt.Println("test2")

	sigs := make(chan os.Signal, 1)
	done := make(chan bool, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		sig := <-sigs
		fmt.Println()
		fmt.Println(sig)
		done <- true
	}()
	fmt.Println("waiting signal")
	<-done
	fmt.Println("exiting")
}

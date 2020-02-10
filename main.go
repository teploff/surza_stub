package main

import (
	"bytes"
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/fatih/color"
	"github.com/gorilla/mux"
	"log"
	"math/rand"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const timeFormat = "2006.01.02 15:04:05.000000"

var (
	src  = flag.String("src", "127.0.0.1:8091", "server listen src")
	dest = flag.String("dest", "127.0.0.1:8092", "server listen src")
	freq = flag.Duration("freq", time.Second, "reactive power frequency")
)

type Payload struct {
	Q float64 `json:"q"`
}

type Task struct {
	cancel chan struct{}
	ticker *time.Ticker
}

func (t *Task) Run() {
	for {
		select {
		case <-t.cancel:
			return
		case <-t.ticker.C:
			work()
		default:
			time.Sleep(time.Millisecond * 10)
		}
	}
}

func work() {
	Q := Payload{Q: rand.Float64()}
	requestBody, err := json.Marshal(Q)
	if err != nil {
		log.Fatal(err)
	}

	destAddr := fmt.Sprintf("http://%s/surza", *dest)
	resp, err := http.Post(destAddr, "application/json", bytes.NewBuffer(requestBody))
	if err != nil {
		log.Println(err)
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		log.Println(resp.Body)
	}
	color.Cyan("%s Q is sent. It has value = %f", time.Now().Format(timeFormat), Q.Q)
}

func GetQEndpoint(_ http.ResponseWriter, r *http.Request) {
	var payload Payload
	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		color.Red("%s"+err.Error(), time.Now().Format(timeFormat))
	}
	color.Yellow("%s Q is received. It has value = %f", time.Now().Format(timeFormat), payload.Q)
}

func main() {
	flag.Parse()
	router := mux.NewRouter()
	router.HandleFunc("/surza", GetQEndpoint).Methods("POST")

	srv := &http.Server{
		Addr:    *src,
		Handler: router,
	}

	task := &Task{
		cancel: make(chan struct{}),
		ticker: time.NewTicker(*freq),
	}
	go task.Run()

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	color.Green("%s Server Started", time.Now().Format(timeFormat))

	<-done
	color.Green("%s Server Stopped", time.Now().Format(timeFormat))

	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer func() {
		// extra handling here
		cancel()
	}()

	if err := srv.Shutdown(ctx); err != nil {
		color.Red("%s Server Shutdown Failed:%+v", time.Now().Format(timeFormat), err)
	}
	task.cancel <- struct{}{}
	color.Green("%s Server Exited Properly", time.Now().Format(timeFormat))
}

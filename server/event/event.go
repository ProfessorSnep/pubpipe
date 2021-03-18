package event

import (
	"fmt"
	"io"
	"net/http"
	"sync"
)

type Subscription struct {
	ch   (chan []byte)
	done bool
}

var (
	subscriptions (map[string][]*Subscription)
	mutex         (*sync.RWMutex)
)

func HandleEvent(w http.ResponseWriter, r *http.Request) {
	key := r.URL.Path

	if r.Method == "POST" {
		mutex.RLock()
		defer mutex.RUnlock()

		subs, ok := subscriptions[key]
		if !ok {
			return
		}

		buf, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		query := r.URL.Query()
		event_id := query.Get("id")
		event_type := query.Get("type")

		body := ""
		if len(event_id) > 0 {
			body += fmt.Sprintf("id: %s\n", event_id)
		}
		if len(event_type) > 0 {
			body += fmt.Sprintf("event: %s\n", event_type)
		}
		body += fmt.Sprintf("data: %s\n\n", string(buf))

		body_buf := []byte(body)
		for _, sub := range subs {
			go func(sub *Subscription) {
				if !sub.done {
					sub.ch <- body_buf
				}
			}(sub)
		}
	} else if r.Method == "GET" {
		w.Header().Set("Connection", "keep-alive")
		w.Header().Set("Cache-Control", "no-cache")
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Access-Control-Allow-Origin", r.Header.Get("Origin"))
		w.Header().Set("Content-Type", "text/event-stream")
		w.Header().Set("X-Accel-Buffering", "no")
		w.Write([]byte(": connected\n"))
		w.WriteHeader(200)
		w.(http.Flusher).Flush()

		sub := Subscription{ch: make(chan []byte), done: false}

		mutex.Lock()
		subscriptions[key] = append(subscriptions[key], &sub)
		mutex.Unlock()

		for !sub.done {
			select {
			case bytes := <-sub.ch:
				w.Write(bytes)
				w.WriteHeader(200)
				w.(http.Flusher).Flush()
			case <-r.Context().Done():
				sub.done = true
				close(sub.ch)
			}
		}
	}
}

func Init() {
	subscriptions = make(map[string][]*Subscription)
	mutex = &sync.RWMutex{}
}

package queue

import (
	"io"
	"net/http"
	"sync"
)

type Stream struct {
	reader io.ReadCloser
	done   chan struct{}
}

var (
	channels (map[string]chan Stream)
	mutex    (*sync.Mutex)
)

func GetChannel(key string) chan Stream {
	mutex.Lock()
	defer mutex.Unlock()
	_, ok := channels[key]
	if !ok {
		channels[key] = make(chan Stream)
	}
	return channels[key]
}

func HandleQueue(w http.ResponseWriter, r *http.Request) {
	channel := GetChannel(r.URL.Path)

	if r.Method == "POST" {
		done := make(chan struct{})
		stream := Stream{reader: r.Body, done: done}
		select {
		case channel <- stream:
		case <-r.Context().Done():
		}
		<-done
	} else if r.Method == "GET" {
		select {
		case stream := <-channel:
			io.Copy(w, stream.reader)
			close(stream.done)
		case <-r.Context().Done():
		}
	}
}

func Init() {
	channels = make(map[string]chan Stream)
	mutex = &sync.Mutex{}
}

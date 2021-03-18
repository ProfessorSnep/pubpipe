package server

import (
	"log"
	"net/http"

	"github.com/ProfessorSnep/pubpipe/server/event"
	"github.com/ProfessorSnep/pubpipe/server/queue"
	"github.com/ProfessorSnep/pubpipe/server/sub"
)

func Start(port string) {
	log.Printf("Starting server on port %s\n", port)

	queue.Init()
	sub.Init()
	event.Init()

	http.HandleFunc("/queue/", queue.HandleQueue)
	http.HandleFunc("/sub/", sub.HandlePubsub)
	http.HandleFunc("/event/", event.HandleEvent)
	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		log.Fatal(err)
	}
}

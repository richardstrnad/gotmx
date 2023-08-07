package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/richardstrnad/gotmx/pkg/infra/inmemory"
	"github.com/richardstrnad/gotmx/pkg/server"
	"github.com/richardstrnad/gotmx/pkg/service"
)

var templateDir = "html/"

func main() {
	store := inmemory.NewInMemoryDataStore()
	service, err := service.New(
		service.WithInMemoryStore(store),
	)
	if err != nil {
		log.Fatal(err)
	}

	s := server.NewServer(service)

	tick := time.Tick(2 * time.Second)
	go func() {
		for range tick {
			s.Publish(
				[]byte(fmt.Sprintf(
					`
            <div id="notifications" hx-swap-oob="afterbegin">
              <div>
                %s
              </div>
            </div>
  `, time.Now().Format(time.RFC3339))))
		}
	}()

	log.Printf("Starting server on localhost:6060")
	http.ListenAndServe("localhost:6060", s)
}

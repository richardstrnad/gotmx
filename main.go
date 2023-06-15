package main

import (
	"fmt"
	"net/http"
	"time"
)

var templateDir = "html/"

func main() {
	store := NewInMemoryDataStore()
	server := NewServer(store)
	tick := time.Tick(2 * time.Second)
	go func() {
		for range tick {
			server.publish(
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

	http.ListenAndServe("localhost:6060", server)
}

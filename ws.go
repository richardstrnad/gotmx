package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"sync"
	"time"

	"golang.org/x/time/rate"
	"nhooyr.io/websocket"
)

type WebSocket struct {
	publishLimiter          *rate.Limiter
	subscriberMessageBuffer int
	subscribersMu           sync.Mutex
	subscribers             map[*subscriber]struct{}
}

type subscriber struct {
	msgs      chan []byte
	closeSlow func()
}

func (s *Server) subscribeHandler(w http.ResponseWriter, r *http.Request) {
	c, err := websocket.Accept(w, r,
		&websocket.AcceptOptions{OriginPatterns: []string{"*"}})
	if err != nil {
		log.Print(err)
		return
	}
	defer c.Close(websocket.StatusInternalError, "")

	err = s.subscribe(r.Context(), c)
	if errors.Is(err, context.Canceled) {
		return
	}
	if websocket.CloseStatus(err) == websocket.StatusNormalClosure ||
		websocket.CloseStatus(err) == websocket.StatusGoingAway {
		return
	}
	if err != nil {
		log.Print(err)
		return
	}
}

func (s *Server) subscribe(ctx context.Context, c *websocket.Conn) error {
	ctx = c.CloseRead(ctx)

	sub := &subscriber{
		msgs: make(chan []byte, s.websocket.subscriberMessageBuffer),
		closeSlow: func() {
			c.Close(websocket.StatusPolicyViolation, "connection too slow to keep up with messages")
		},
	}
	s.addSubscriber(sub)
	defer s.deleteSubscriber(sub)

	for {
		select {
		case msg := <-sub.msgs:
			err := writeTimeout(ctx, time.Second*5, c, msg)
			if err != nil {
				return err
			}
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}

// addSubscriber registers a subscriber.
func (s *Server) addSubscriber(sub *subscriber) {
	s.websocket.subscribersMu.Lock()
	s.websocket.subscribers[sub] = struct{}{}
	s.websocket.subscribersMu.Unlock()
}

// deleteSubscriber deletes the given subscriber.
func (s *Server) deleteSubscriber(sub *subscriber) {
	s.websocket.subscribersMu.Lock()
	delete(s.websocket.subscribers, sub)
	s.websocket.subscribersMu.Unlock()
}

func (s *Server) publish(msg []byte) {
	s.websocket.subscribersMu.Lock()
	defer s.websocket.subscribersMu.Unlock()

	s.websocket.publishLimiter.Wait(context.Background())

	for s := range s.websocket.subscribers {
		select {
		case s.msgs <- msg:
		default:
			go s.closeSlow()
		}
	}
}

func writeTimeout(ctx context.Context, timeout time.Duration, c *websocket.Conn, msg []byte) error {
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	return c.Write(ctx, websocket.MessageText, msg)
}

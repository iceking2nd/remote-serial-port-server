package services

import (
	"os"
	"sync"
	"testing"

	"github.com/iceking2nd/remote-serial-port-server/global"
	"github.com/sirupsen/logrus"
)

func TestMain(m *testing.M) {
	global.Log = logrus.New()
	global.Log.SetLevel(logrus.PanicLevel)
	os.Exit(m.Run())
}

func TestGetPortManagerSingleton(t *testing.T) {
	// Reset singleton for test
	once = sync.Once{}
	defaultPortManager = nil

	pm1 := GetPortManager()
	pm2 := GetPortManager()

	if pm1 != pm2 {
		t.Fatal("GetPortManager should return the same instance")
	}
}

func TestPortManagerRemoveSubscriberNoSession(t *testing.T) {
	once = sync.Once{}
	defaultPortManager = nil

	pm := GetPortManager()
	// Should not panic when removing a subscriber from a non-existent session
	pm.RemoveSubscriber("nonexistent_port", "sub1")
}

func TestPortSessionBroadcast(t *testing.T) {
	ps := &PortSession{
		PortName:    "TEST",
		subscribers: make(map[string]*Subscriber),
		done:        make(chan struct{}),
	}

	sub1 := &Subscriber{ID: "sub1", Send: make(chan []byte, 8)}
	sub2 := &Subscriber{ID: "sub2", Send: make(chan []byte, 8)}

	ps.mu.Lock()
	ps.subscribers[sub1.ID] = sub1
	ps.subscribers[sub2.ID] = sub2
	ps.mu.Unlock()

	data := []byte("hello")
	ps.broadcast(data)

	got1 := <-sub1.Send
	got2 := <-sub2.Send

	if string(got1) != "hello" {
		t.Fatalf("sub1 expected 'hello', got '%s'", string(got1))
	}
	if string(got2) != "hello" {
		t.Fatalf("sub2 expected 'hello', got '%s'", string(got2))
	}
}

func TestPortSessionBroadcastDropsOnFullChannel(t *testing.T) {
	ps := &PortSession{
		PortName:    "TEST",
		subscribers: make(map[string]*Subscriber),
		done:        make(chan struct{}),
	}

	// Channel with capacity 1, will be full after first write
	sub := &Subscriber{ID: "sub1", Send: make(chan []byte, 1)}
	ps.mu.Lock()
	ps.subscribers[sub.ID] = sub
	ps.mu.Unlock()

	// Fill the channel
	ps.broadcast([]byte("first"))
	// This should not block even though channel is full
	ps.broadcast([]byte("second"))

	got := <-sub.Send
	if string(got) != "first" {
		t.Fatalf("expected 'first', got '%s'", string(got))
	}
}

func TestPortSessionBroadcastClose(t *testing.T) {
	ps := &PortSession{
		PortName:    "TEST",
		subscribers: make(map[string]*Subscriber),
		done:        make(chan struct{}),
	}

	sub := &Subscriber{ID: "sub1", Send: make(chan []byte, 8)}
	ps.mu.Lock()
	ps.subscribers[sub.ID] = sub
	ps.mu.Unlock()

	ps.broadcastClose()

	_, ok := <-sub.Send
	if ok {
		t.Fatal("subscriber Send channel should be closed after broadcastClose")
	}
}

func TestPortManagerRemoveSubscriberLastOneClosesSession(t *testing.T) {
	once = sync.Once{}
	defaultPortManager = nil

	pm := GetPortManager()

	// Manually create a session to simulate an open port
	session := &PortSession{
		PortName:    "COM_TEST",
		subscribers: make(map[string]*Subscriber),
		done:        make(chan struct{}),
	}
	sub := &Subscriber{ID: "sub1", Send: make(chan []byte, 8)}
	session.subscribers[sub.ID] = sub

	pm.mu.Lock()
	pm.sessions["COM_TEST"] = session
	pm.mu.Unlock()

	// Removing the last subscriber should trigger closePort
	// Since we don't have a real serial.Port, closePort will error on Port.Close
	// but the session should still be removed from the map
	pm.RemoveSubscriber("COM_TEST", "sub1")

	pm.mu.RLock()
	_, exists := pm.sessions["COM_TEST"]
	pm.mu.RUnlock()

	if exists {
		t.Fatal("session should be removed from PortManager after last subscriber leaves")
	}
}

func TestPortManagerRemoveSubscriberKeepsSessionWhenOthersRemain(t *testing.T) {
	once = sync.Once{}
	defaultPortManager = nil

	pm := GetPortManager()

	session := &PortSession{
		PortName:    "COM_TEST",
		subscribers: make(map[string]*Subscriber),
		done:        make(chan struct{}),
	}
	sub1 := &Subscriber{ID: "sub1", Send: make(chan []byte, 8)}
	sub2 := &Subscriber{ID: "sub2", Send: make(chan []byte, 8)}
	session.subscribers[sub1.ID] = sub1
	session.subscribers[sub2.ID] = sub2

	pm.mu.Lock()
	pm.sessions["COM_TEST"] = session
	pm.mu.Unlock()

	pm.RemoveSubscriber("COM_TEST", "sub1")

	pm.mu.RLock()
	_, exists := pm.sessions["COM_TEST"]
	pm.mu.RUnlock()

	if !exists {
		t.Fatal("session should still exist when subscribers remain")
	}

	session.mu.RLock()
	count := len(session.subscribers)
	session.mu.RUnlock()

	if count != 1 {
		t.Fatalf("expected 1 subscriber, got %d", count)
	}
}

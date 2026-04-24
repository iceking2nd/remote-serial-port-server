package services

import (
	"fmt"
	"io"
	"sync"

	"github.com/iceking2nd/remote-serial-port-server/global"
	"github.com/sirupsen/logrus"
	"go.bug.st/serial"
)

// Subscriber represents a WebSocket client subscribed to a serial port.
type Subscriber struct {
	ID   string
	Send chan []byte
}

// PortSession manages a single opened serial port and its subscribers.
type PortSession struct {
	PortName    string
	Port        serial.Port
	Mode        *serial.Mode
	subscribers map[string]*Subscriber
	mu          sync.RWMutex
	done        chan struct{}
	once        sync.Once
}

// PortManager is a global singleton that manages all active serial port sessions.
type PortManager struct {
	sessions map[string]*PortSession
	mu       sync.RWMutex
}

var defaultPortManager *PortManager
var once sync.Once

// GetPortManager returns the singleton PortManager instance.
func GetPortManager() *PortManager {
	once.Do(func() {
		defaultPortManager = &PortManager{
			sessions: make(map[string]*PortSession),
		}
	})
	return defaultPortManager
}

// OpenPort opens a serial port or returns an existing session if the port is already open.
// It adds the subscriber to the session's subscriber list.
func (pm *PortManager) OpenPort(portName string, mode *serial.Mode, subscriber *Subscriber) (*PortSession, error) {
	log := global.Log.WithField("function", "services.PortManager.OpenPort")

	pm.mu.Lock()
	defer pm.mu.Unlock()

	if session, ok := pm.sessions[portName]; ok {
		log.WithField("port", portName).Debug("port already open, adding subscriber")
		session.mu.Lock()
		session.subscribers[subscriber.ID] = subscriber
		session.mu.Unlock()
		return session, nil
	}

	serialPort, err := serial.Open(portName, mode)
	if err != nil {
		log.WithError(err).WithField("port", portName).Error("failed to open serial port")
		return nil, fmt.Errorf("failed to open serial port %s: %w", portName, err)
	}

	session := &PortSession{
		PortName:    portName,
		Port:        serialPort,
		Mode:        mode,
		subscribers: map[string]*Subscriber{subscriber.ID: subscriber},
		done:        make(chan struct{}),
	}
	pm.sessions[portName] = session

	go session.readLoop()

	log.WithField("port", portName).Debug("serial port opened, read loop started")
	return session, nil
}

// RemoveSubscriber removes a subscriber from a port session.
// If this is the last subscriber, the port is closed and the session is removed.
func (pm *PortManager) RemoveSubscriber(portName string, subscriberID string) {
	log := global.Log.WithField("function", "services.PortManager.RemoveSubscriber")

	pm.mu.RLock()
	session, ok := pm.sessions[portName]
	pm.mu.RUnlock()

	if !ok {
		return
	}

	session.mu.Lock()
	delete(session.subscribers, subscriberID)
	remaining := len(session.subscribers)
	session.mu.Unlock()

	log.WithFields(logrus.Fields{
		"port":       portName,
		"subscriber": subscriberID,
		"remaining":  remaining,
	}).Debug("subscriber removed")

	if remaining == 0 {
		pm.closePort(portName, session)
	}
}

// closePort closes the serial port and removes the session.
func (pm *PortManager) closePort(portName string, session *PortSession) {
	log := global.Log.WithField("function", "services.PortManager.closePort")

	session.once.Do(func() {
		close(session.done)
	})

	pm.mu.Lock()
	delete(pm.sessions, portName)
	pm.mu.Unlock()

	if session.Port != nil {
		if err := session.Port.Close(); err != nil {
			log.WithError(err).WithField("port", portName).Error("failed to close serial port")
		} else {
			log.WithField("port", portName).Debug("serial port closed")
		}
	}
}

// readLoop reads data from the serial port and broadcasts to all subscribers.
func (ps *PortSession) readLoop() {
	log := global.Log.WithField("function", "PortSession.readLoop")
	buf := make([]byte, 128)

	for {
		select {
		case <-ps.done:
			return
		default:
		}

		n, err := ps.Port.Read(buf)
		if err != nil {
			if err == io.EOF {
				log.WithField("port", ps.PortName).Debug("serial port EOF")
			} else {
				log.WithError(err).WithField("port", ps.PortName).Error("error reading from serial port")
			}
			ps.broadcastClose()
			return
		}

		data := make([]byte, n)
		copy(data, buf[:n])
		ps.broadcast(data)
	}
}

// broadcast sends data to all subscribers.
func (ps *PortSession) broadcast(data []byte) {
	ps.mu.RLock()
	defer ps.mu.RUnlock()

	for _, sub := range ps.subscribers {
		select {
		case sub.Send <- data:
		default:
			log := global.Log.WithField("function", "PortSession.broadcast")
			log.WithField("subscriber", sub.ID).Warn("subscriber send channel full, dropping data")
		}
	}
}

// broadcastClose signals all subscribers that the serial port is closed.
func (ps *PortSession) broadcastClose() {
	ps.mu.RLock()
	defer ps.mu.RUnlock()

	for _, sub := range ps.subscribers {
		close(sub.Send)
	}
}

// Write writes data to the serial port.
func (ps *PortSession) Write(data []byte) (int, error) {
	return ps.Port.Write(data)
}

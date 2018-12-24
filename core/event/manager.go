package event

import (
	"sync"
)

// Manager Event manager
// Handles sending and receiving events for immediate handling
// Generally used for engine functionality, such as user input events, window
// management etc.
// Game entities should use their own event queue, and should not hook into this queue.
type Manager struct {
	listenerMap map[MessageType]map[EventHandle]IEventListenable
	mu          sync.Mutex
	eventQueue  []*queueItem
	runAsync    bool
}

// Listen Register a new listener to listen to an event
func (manager *Manager) Listen(eventName MessageType, listener IEventListenable) EventHandle {
	handle := newEventHandle()
	manager.mu.Lock()
	if _, ok := manager.listenerMap[eventName]; !ok {
		manager.listenerMap[eventName] = make(map[EventHandle]IEventListenable)
	}
	manager.listenerMap[eventName][handle] = listener
	manager.mu.Unlock()

	return handle
}

// RunConcurrent Runs the event queue in its own go routine
func (manager *Manager) RunConcurrent() {
	// Block double-running
	if manager.runAsync == true {
		return
	}
	manager.runAsync = true
	go func() {
		for manager.runAsync == true {
			manager.mu.Lock()
			queue := manager.eventQueue
			manager.mu.Unlock()

			if len(queue) > 0 {
				// FIFO - ensure dispatch order, and concurrency integrity
				item := queue[0]
				manager.mu.Lock()
				manager.eventQueue = manager.eventQueue[1:]

				// Fire event
				listeners := manager.listenerMap[item.EventName]
				manager.mu.Unlock()
				for _, listener := range listeners {
					listener.ReceiveMessage(item.Message)
				}
			}
		}
	}()
}

// Unlisten Remove a listener from listening for an event
func (manager *Manager) Unlisten(eventName MessageType, handle EventHandle) {
	manager.mu.Lock()
	if _, ok := manager.listenerMap[eventName][handle]; ok {
		delete(manager.listenerMap[eventName], handle)
	}
	manager.mu.Unlock()
}

// Dispatch Fires an event to all listening objects
func (manager *Manager) Dispatch(eventName MessageType, message IMessage) {
	message.SetType(eventName)
	queueItem := &queueItem{
		EventName: eventName,
		Message:   message,
	}
	manager.mu.Lock()
	manager.eventQueue = append(manager.eventQueue, queueItem)
	manager.mu.Unlock()
}

// Unregister Close the event manager
func (manager *Manager) Unregister() {
	// Ensure async event queue is halted
	manager.runAsync = false
}

var eventManager Manager

// GetEventManager return static eventmanager
func GetEventManager() *Manager {
	if eventManager.listenerMap == nil {
		eventManager.listenerMap = make(map[MessageType]map[EventHandle]IEventListenable, 512)
	}
	return &eventManager
}

// queueItem Event Queue item.
// Contains the event name, and a message
type queueItem struct {
	EventName MessageType
	Message   IMessage
}
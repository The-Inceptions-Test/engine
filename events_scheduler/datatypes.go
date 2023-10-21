// Structures required to create and process events

package events_scheduler

import (
	"sync"
	"time"

	"github.com/caffix/queue"
	"github.com/google/uuid"
)

// Event types (are used to query the Registry adn identify the action to be executed)
type EventType int

const (
	// EventTypeSay is used to print a message to the console
	EventTypeSay EventType = iota
	// EventTypeLog is used to log a message to the log file
	EventTypeLog
	// Add more event types here:
)

// Event states (are used to control the event flow)
type EventState int

const (
	StateDefault EventState = iota // Event is in default state
	// (normally used when the event is created)
	StateProcessable // Event is processable (all dependencies are met)
	StateWaiting     // Event is waiting (some dependencies are not met)
	StateDone        // Event is done (already processed)
	StateInProcess   // Event is in process (being processed)
	StateCancelled   // Event is cancelled (not processed)
	StateError       // Event is in error (not processed)
)

// Global variables
var (
	// zeroUUID is used to indicate that an event has no dependencies
	zeroUUID = uuid.UUID{}
)

// Event is the struct that represents an event
// This struct it's kind of the "currency of exchange" between the scheduler
// and the functions that create and process the events
type Event struct {
	UUID      uuid.UUID           /* Event UUID */
	Session   uuid.UUID           /* Session UUID */
	Name      string              /* Event name */
	Timestamp time.Time           /* Event timestamp */
	Type      EventType           /* Event type */
	State     EventState          /* Event state (processable, waiting, done, in process) */
	DependOn  []uuid.UUID         /* Events this event "depends on" */
	Action    func(e Event) error /* Event handler function (action) (normally populated by querying the
	-                                Registry)
	-                              */
	Priority    int /* Event priority (normally populated by querying the Registry) */
	RepeatEvery int /* Event repeat every X centiseconds (normally populated by querying
	-			       the Registry)
	-                */
	RepeatTimes int         /* Event repeat times (normally populated by querying the Registry) */
	Data        interface{} /* This field can hold any data type (normally populated by the function
	-                          that creates the event, and used by the function that processes the
	-                          event)
	-                        */
	timeout time.Time  /* Timeout timer (used to cancel the event if it's not processed in time) */
	s       *Scheduler /* Pointer to the scheduler that created the event */
}

// Scheduler is the struct that represents a scheduler
// We have 2 types of schedulers:
//   - Main scheduler, used to schedule and process events, it's the central scheduler and it's
//     allocated on the heap (it's a singleton) and it's initialized by calling
//     the MainSchedulerInit() function.
//   - Sub schedulers, used to schedule and process events, they are allocated on the stack and
//     they are initialized by calling the NewScheduler() function.
type Scheduler struct {
	q      queue.Queue          // Events Queue (Queue to store events)
	mutex  sync.Mutex           // Mutex to protect the queue when fetching the next event
	events map[uuid.UUID]*Event // Map to quickly look up events by UUID
}

// ProcessConfig is the struct that represents the configuration used to process the events
type ProcessConfig struct {
	ExitWhenEmpty bool
	CheckEvent    bool
	ExecuteAction bool
	ReturnIfFound bool
	DebugInfo     bool
	ActionTimeout int
}
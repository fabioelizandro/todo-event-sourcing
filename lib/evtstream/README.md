# Event Stream 

Small and simple lib for helping with event sourcing in golang. This libs aims to create a set of interface to accommodate almost any kind 
of event stream. The only built in implementation it has for now is a prevalent event stream where it'll be loaded in memory and persist 
in a storage like a database sql or s3 files.

## Usage

```go
package main

import (
	"fabioelizandro/todo-event-sourcing/lib/evtstream"
	"fmt"
)

func main () {
    // Loads eventstream from disk.
    store := evtstream.NewDiskPrevalentEventStore(
		"/tmp/todo-event-sourcing-stream",
		evtstream.NewInMemoryEventRegistry([]evtstream.Event{
			&EvtTaskCompleted{},
			&EvtTaskAdded{},
			&EvtTaskDescriptionChanged{},
		}),
	)

	envelopes, err := store.Load()
	if err != nil {
		panic(err)
	}
	
	stream := evtstream.NewPrevalentEventStream(store, envelopes, evtstream.NewUTCCLock())

    // Loads domain model from stream 
    task := &TaskDomainModel{}
    aggregateEvents, err := eventStream.ReadByCorrelationID("abc")
    if err != nil {
        panic(err)
    }

    for _, evt := range aggregateEvents {
        task.Apply(evt)
    }

    // Execute commands to domain model
    events := task.UpdateDescription("Do the dishes")
    
    err := stream.Write(events)
    if err != nil {
        panic(err)
    }
}
```
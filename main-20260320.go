package main

import (
	"fmt"
	"log"
	"net/http"

	// Using BoltDB which is present in your nostrlib folder
	"fiatjaf.com/nostr/eventstore/boltdb"
	"fiatjaf.com/nostr/khatru"
)

func main() {
	relay := khatru.NewRelay()

	// 1. Initialize BoltDB backend
	// This will create a local file named "relay.bolt"
	db := &boltdb.BoltBackend{Path: "relay.bolt"}
	if err := db.Init(); err != nil {
		log.Fatalf("Error: failed to initialize boltdb: %v", err)
	}

	// 2. Attach the storage backend
	// khatru will use BoltDB to store and query Nostr events
	relay.UseEventstore(db, 500)

	relay.Info.Name = "Zhihao's BoltDB Relay"
	relay.Info.Description = "Running with local nostrlib/boltdb persistence."

	port := ":3334"
	fmt.Printf("🚀 Relay is running with BoltDB at %s\n", port)
	
	if err := http.ListenAndServe(port, relay); err != nil {
		log.Fatal(err)
	}
}

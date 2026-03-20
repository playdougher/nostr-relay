package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strings"

	"fiatjaf.com/nostr"
	"fiatjaf.com/nostr/eventstore/boltdb"
	"fiatjaf.com/nostr/khatru"
	"fiatjaf.com/nostr/nip19"
)

func main() {
	relay := khatru.NewRelay()

	// 1. Initialize Storage
	db := &boltdb.BoltBackend{Path: "relay.bolt"}
	db.Init()
	relay.UseEventstore(db, 500) //

	// 2. Load Whitelist
	whitelist := make(map[string]bool)
	if data, err := os.ReadFile("whitelist.txt"); err == nil {
		for _, line := range strings.Split(string(data), "\n") {
			pk := strings.TrimSpace(line)
			if pk != "" && !strings.HasPrefix(pk, "#") {
				whitelist[pk] = true
			}
		}
	}
	fmt.Printf("Relay Start: %d keys in whitelist\n", len(whitelist))

	// 3. Simple Write Guard (OnEvent)
	// This ensures only whitelisted pubkeys can save events
        relay.OnEvent = func(ctx context.Context, ev nostr.Event) (bool, string) {
            // 1. 先拿到那个带前缀的原始字符串 (比如 "pk::06288d...")
            rawPubKey := ev.PubKey.String()
        
            // 2. 强力去掉开头的 "pk::"（如果有的话）
            // 这样 hexPK 就会变成干净的 "06288d..."
            hexPK := strings.TrimPrefix(rawPubKey, "pk::")
        
            // 3. 编码 npub 用于日志显示
            npub := nip19.EncodeNpub(ev.PubKey)
        
            // 4. 比对白名单
            if len(whitelist) > 0 && !whitelist[hexPK] {
                // 日志里我们看看脱皮后的 hexPK 对不对
                fmt.Printf("[DENY WRITE] %s (Checking Hex: %s)\n", npub, hexPK)
                return true, "auth-required: you are not on the whitelist"
            }
        
            fmt.Printf("[ACCEPT WRITE] %s\n", npub)
            return false, ""
        }
	// 4. Read Guard: Empty (Simplest)
	// We do NOT set relay.OnRequest. This allows anyone to view messages.
	// This solves your problem where "whitelisted users cannot view".

	port := ":3334"
	fmt.Printf("Simple Private Relay running on %s\n", port)
	http.ListenAndServe(port, relay)
}

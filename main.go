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
            // 1. 基础信息处理
            hexPK := strings.TrimPrefix(ev.PubKey.String(), "pk::")
            npub := nip19.EncodeNpub(ev.PubKey)
        
            // 2. 白名单检查 (DENY 逻辑)
            if len(whitelist) > 0 && !whitelist[hexPK] {
                // 使用与 HIT 一模一样的格式占位符
                // %-5d 让 Kind 左对齐，%x... 让 ID 看起来专业
                fmt.Printf("[DENY] Kind:%-5d | ID:%x... | From:%s\n", ev.Kind, ev.ID[:4], npub[:15])
                return true, "auth-required: you are not on the whitelist"
            }
        
            // 3. 过滤器 (HIT 逻辑)
            // isImportant := (ev.Kind == 1 || ev.Kind == 5 || ev.Kind == 7 || ev.Kind == 9735 || true)
            isImportant := (ev.Kind == 1)
        
            if isImportant {
                target := ""
                if ev.Kind == 5 {
                    for _, tag := range ev.Tags {
                        if len(tag) >= 2 && (tag[0] == "e" || tag[0] == "a") {
                            target = fmt.Sprintf(" -> Target:%s...", tag[1][:8])
                        }
                    }
                }
                // 保持格式高度一致
                fmt.Printf("[HIT ] Kind:%-5d | ID:%x...%s | From:%s\n", ev.Kind, ev.ID[:4], target, npub[:15])
            }
        
            return false, ""
        }
	// 4. Read Guard: Empty (Simplest)
	// We do NOT set relay.OnRequest. This allows anyone to view messages.
	// This solves your problem where "whitelisted users cannot view".

	port := ":3334"
	fmt.Printf("Simple Private Relay running on %s\n", port)
	http.ListenAndServe(port, relay)
}

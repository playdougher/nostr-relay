package main

import (
	"context"
	"fmt"
	"iter"
	"log"
	"net/http"

	"fiatjaf.com/nostr"
	"fiatjaf.com/nostr/khatru"
)

func main() {
	relay := khatru.NewRelay()

	relay.Info.Name = "我的中继器"
	relay.Info.Description = "修复了所有类型错误的版本"
	// 修复 1: MustPubKeyFromHex 返回的是数组值，而 Info.PubKey 现在需要一个指针
	pk := nostr.MustPubKeyFromHex("79be667ef9dcbbac55a06295ce870b07029bfcdb2dce28d959f2815b16f81798")
	relay.Info.PubKey = &pk

	// 修复 2: 存储使用 nostr.ID (数组) 作为键，nostr.Event (值) 作为值
	store := make(map[nostr.ID]nostr.Event)

	// 修复 3: StoreEvent 现在接收的是值 (nostr.Event)，不再是指针
	relay.StoreEvent = func(ctx context.Context, event nostr.Event) error {
		store[event.ID] = event
		return nil
	}

	// 修复 4: QueryStored 现在返回的是 iter.Seq[nostr.Event] (值序列)
	relay.QueryStored = func(ctx context.Context, filter nostr.Filter) iter.Seq[nostr.Event] {
		return func(yield func(nostr.Event) bool) {
			for _, evt := range store {
				// 修复 5: filter.Matches 现在直接接收值
				if filter.Matches(evt) {
					if !yield(evt) {
						break
					}
				}
			}
		}
	}

	// 修复 6: DeleteEvent 接收的是 nostr.ID 数组类型
	relay.DeleteEvent = func(ctx context.Context, id nostr.ID) error {
		delete(store, id)
		return nil
	}


	fmt.Println("Relay 正在运行在 :3334")
	if err := http.ListenAndServe(":3334", relay); err != nil {
		log.Fatal(err)
	}
}

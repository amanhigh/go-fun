package libraries

import (
	"fmt"

	"github.com/hashicorp/consul/api"
)

// HACK: #C Rewrite Consul as Ginkgo
func ConsulFun() {
	// Get a new client
	if client, err := api.NewClient(&api.Config{Address: "docker:8500"}); err == nil {
		// Get a handle to the KV API
		kv := client.KV()

		// PUT a new KV pair
		key := "aman/1"
		p := api.KVPair{Key: key, Value: []byte("2000")}
		if _, err = kv.Put(&p, nil); err == nil {
			if pair, _, err := kv.Get(key, nil); err == nil {
				fmt.Printf("KV: %v %s\n", pair.Key, pair.Value)
				//_, err = kv.Delete(key, nil)
			}

		}
		// Log Error
		if err != nil {
			panic(err)
		}
	}
}

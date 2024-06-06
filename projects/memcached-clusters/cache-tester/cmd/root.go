package cmd

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/bradfitz/gomemcache/memcache"
)

func Execute() {
	mcrouter := flag.String("mcrouter", "", "port to connect to mcrouter")
	memcacheds := flag.String("memcacheds", "", "comma separated list of memcached servers ports")
	flag.Parse()

	if *mcrouter == "" || *memcacheds == "" {
		flag.Usage()
		os.Exit(1)
	}

	// Connect to mcrouter
	mcr := memcache.New("localhost:" + *mcrouter)
	defer mcr.Close()

	// Connect to memcached servers
	ports := strings.Split(*memcacheds, ",")
	mcs := make([]*memcache.Client, 0)
	for i, port := range ports {
		mcs = append(mcs, memcache.New("localhost:"+port))
		defer mcs[i].Close()
	}

	// Set a key in mcrouter
	err := mcr.Set(&memcache.Item{Key: "hello", Value: []byte("world")})
	if err != nil {
		fmt.Printf("Error setting key in mcrouter: %v\n", err)
		os.Exit(1)
	}

	// Get the key from each memcached server
	items := make([]*memcache.Item, 0)
	item := new(memcache.Item)
	for i, mc := range mcs {
		item, err = mc.Get("hello")
		if err != nil {
			fmt.Printf("Error getting key from memcached server %v: %v\n", i, err)
		}
		if err == memcache.ErrCacheMiss {
			fmt.Println("Cache typology: sharded")
			return
		}
		items = append(items, item)
	}

	if allItemsEqual(items) {
		fmt.Println("Memcached cluster typology : replicated")
	}

}

func allItemsEqual(items []*memcache.Item) bool {
	for i := 1; i < len(items); i++ {
		if string(items[i].Value) != string(items[i-1].Value) {
			return false
		}
	}
	return true
}

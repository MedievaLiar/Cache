package main

import (
	"fmt"
	"time"
	
	"github.com/MedievaLiar/Cache/cache"
)

func main() {
    c := cache.NewCache(5*time.Minute, 10*time.Second, 1000)
    defer c.Stop()

    c.Set("user1", "Ann", 0)
    c.Set("user2", "Jim", 30*time.Second)

    if val, exists := c.Get("user1"); exists {
        fmt.Printf("User 1: %v\n", val)
    }

    keys := c.Keys()
    fmt.Printf("Ключи: %v\n", keys)
}

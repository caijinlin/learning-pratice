package main

import (
	"fmt"
	"github.com/gomodule/redigo/redis"
	"math/rand"
	_ "net"
	"net/http"
	"os"
	"time"
)

func main() {

	addrs := []string{"127.0.0.1:6379", "127.0.0.1:6380"}
	pool := newPool(addrs) // read pool

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {

		c := pool.Get()

		defer c.Close() // put or close

		_, er := c.Do("SET", "username", "nick")
		if er != nil {
			fmt.Fprintf(w, "redis set failed:"+er.Error())
		}
		username, err := redis.String(c.Do("GET", "username"))
		if err != nil {
			fmt.Fprintf(w, "redis get failed:"+err.Error())
		} else {
			fmt.Fprintf(w, "Got username:"+username)
		}
	})

	server := http.Server{
		Addr: ":3333",
	}

	err := server.ListenAndServe()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %s\n", err.Error())
	}
}

func newPool(addrs []string) *redis.Pool {
	var rand_gen = rand.New(rand.NewSource(time.Now().UnixNano()))
	return &redis.Pool{
		MaxIdle:     3,
		MaxActive:   6,
		IdleTimeout: 1 * time.Second,
		Wait:        false,
		// Dial or DialContext must be set. When both are set, DialContext takes precedence over Dial.
		Dial: func() (redis.Conn, error) {
			index := rand_gen.Intn(len(addrs))
			return redis.Dial("tcp", addrs[index])
		},
	}
}

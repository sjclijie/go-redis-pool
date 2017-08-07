package main

import (
	"github.com/garyburd/redigo/redis"
	"time"
	"net/http"
	"fmt"
)

var pool *redis.Pool

func newPool(server, passport string) *redis.Pool {
	return &redis.Pool{
		MaxIdle:     800,
		MaxActive:   1000,
		IdleTimeout: 200 * time.Second,
		Dial: func() (redis.Conn, error) {

			c, err := redis.Dial("tcp", server)

			if err != nil {
				return nil, err
			}

			if passport != "" {
				if _, err := c.Do("AUTH", passport); err != nil {
					c.Close()
					return nil, err
				}
			}

			return c, err
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			if time.Since(t) < time.Minute {
				return nil
			}

			_, err := c.Do("PING")

			return err
		},
	}
}

func init() {
	pool = newPool("172.18.1.7:6379", "")
}

func hello(w http.ResponseWriter, r *http.Request) {

	conn := pool.Get()

	test, _ := redis.String(conn.Do("GET", "name"))

	activeCount := pool.ActiveCount()
	idleCount := pool.IdleCount()

	fmt.Println("action: ", activeCount, idleCount)

	fmt.Fprint(w, test)

	return
}

func main() {

	http.HandleFunc("/hello", hello)
	http.ListenAndServe(":9999", nil)
}

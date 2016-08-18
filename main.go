package main

import (
	"encoding/base64"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/garyburd/redigo/redis"
	cache "github.com/patrickmn/go-cache"
)

var redisConnection redis.Conn
var server *http.Server
var localCache *cache.Cache

func establishRedis() {
	conn, err := redis.Dial("tcp", "127.0.0.1:6379")
	if err != nil {
		fmt.Printf("Could not connect: %s", err)
	} else {
		redisConnection = conn
		res, _ := redis.String(redisConnection.Do("PING"))
		fmt.Printf("Redis ping: %s\n", res)
	}
}

func establishCache() {
	c := cache.New(5*time.Second, 10*time.Second)
	localCache = c
}

func httpHandler(w http.ResponseWriter, r *http.Request) {
	// turn the URL into base64 so we dont pass special characters to Redis
	bytes := []byte(r.URL.Path)
	url64 := base64.StdEncoding.EncodeToString(bytes)
	key := fmt.Sprintf("%s::%s", "RedisWebContent", url64)

	fmt.Printf("Starting request for %s -- cache key: %s\n", r.URL.Path, key)

	haveResponse := false
	var response string

	data, found := localCache.Get(url64)
	if found {
		fmt.Println("  Found in local cache")
		haveResponse = true
		response = data.(string)
	} else {
		resp, err := redisConnection.Do("GET", key)
		if err == nil {
			content, err2 := redis.String(resp, nil)
			if err2 == nil {
				fmt.Println("  Found in Redis")
				localCache.Set(url64, content, 5*time.Second)
				haveResponse = true
				response = content
			} else {
				fmt.Println("  Not found in Redis or local cache")
			}
		} else {
			fmt.Println("  Error checking Redis")
		}
	}

	if haveResponse {
		io.WriteString(w, response)
		fmt.Println("Finished with success response")
	} else {
		fmt.Fprintf(w, "Could not find %s", r.URL.Path)
		fmt.Println("Finished with error response")
	}

}

func establishHTTPServer() {
	server = &http.Server{
		Addr:    ":8080",
		Handler: http.HandlerFunc(httpHandler),
	}
}

func main() {
	establishRedis()
	establishCache()
	establishHTTPServer()

	log.Fatal(server.ListenAndServe())
}

// Socket.getifaddrs.select{|i| i.name == 'en0'}.first.addr.getnameinfo[0].gsub(':', '').to_i(16).to_s(36)

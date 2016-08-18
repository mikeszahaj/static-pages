# Static Pages Server for Go -- Proof of Concept
### Serves content from Redis, with caching inside Go.  Nowhere near production-ready.

## Set up
- Install Redis
- Set up Go including your `$GOPATH`
- `go get github.com/garyburd/redigo/redis`
- `go get github.com/patrickmn/go-cache`


## Run
- `go build .`
- `./static-pages`

## Debug and Set Content
All of these run from within a Redis prompt (`redis-cli`)
- `monitor`: Watch incoming Redis commands
- `set RedisWebContent::Lw== "Hello"`: Set content for page `/` to `Hello` (`Lw==` is `/` in Base64)

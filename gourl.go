package main

import (
  "crypto/sha1"
  "fmt"
  "mustache"
  "redis"
  "regexp"
  "web"
)

var redisClient *redis.Client = new(redis.Client)

/*keyOf returns (part of) the SHA-1 hash of the data, as a hex string.*/
func keyOf(data []byte) string {
  sha := sha1.New()
  sha.Write(data)
  return fmt.Sprintf("%x", string(sha.Sum())[0:5])
}

/*returns the key as redis will see it*/
func getRedisKey(key string) string {
  return fmt.Sprintf(":urls:%s", key)
}

/*just prints a basic form to receive a URL*/
func index(ctx *web.Context) {
  ctx.WriteString (mustache.RenderFile("index.mustache"))
} 

/*shows the "shortened" URL to the user*/
func show(ctx *web.Context, key string) {
  ctx.WriteString (mustache.RenderFile("show.mustache", map[string]string {"key":key}))
} 

/*redirects the shortened URL to the real URL*/
func redirect(ctx *web.Context, key string) {
  /*fetch our URL*/
  url,_ := redisClient.Get(getRedisKey(key))
  if url == nil {
    printError(ctx, "I can't find that URL")
  } else {
    /*redirect*/
    ctx.Redirect(301, string(url)) 
  } 
} 

func printError(ctx *web.Context, error string){
  /*this is a 500 error condition*/
  ctx.StartResponse(500)
  /*print boilerplate error page with passed in message*/
  ctx.WriteString (mustache.RenderFile("error.mustache", map[string]string {"error":error}))
}

func shorten(ctx *web.Context){
  /*fetch URL and convert to string type*/
  url := fmt.Sprintf("%s", ctx.Request.Params["url"])
  /*crude REGEX to make sure URL is more or less a URL*/
  isURL, _ := regexp.MatchString("^http(s)?://.*", url)

  /*I think this is probably supposed to be a switch statement*/
  /*but it is my first Go app so I didn't get too crazy*/
  if url == "" {
    printError(ctx, "URL is missing. Return to Go. Do not collect $200.")
  } else if !isURL {
    printError(ctx, "That doesn't look like a URL.")
  } else {
    /*generate short key*/
    key := keyOf([]byte(url))
    /*set URL in Redis using a Redis-ized key*/
    err := redisClient.Set(getRedisKey(key), []byte(url))
    /*redirect to show page*/
    ctx.Redirect(301, fmt.Sprintf("/s/%s", key))
  } 
}

func main() {
  /*setup redis connection*/
  redisClient.Addr = "localhost:6379"
  redisClient.Db = 13

  /*setup web.go stuff*/
  web.Get("/", index)
  web.Get("/s/(.*)", show)
  web.Get("/u/(.*)", redirect)
  web.Post("/shorten", shorten)
  web.Run("0.0.0.0:8080")
}

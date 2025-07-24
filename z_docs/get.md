# general get request architecture 07/24/2025
## entrypoint: Get function
the Get function accepts a host, end[point], params, & headers 
(both slices of key-val pairs) and returns a response body, HTTP status code, &
error  
example that prints the response body: 
```go
body, _, err := Get("stats.nba.com", "/stats/commonplayerinfo", params, hdrs)
	if err != nil {
		log.Fatal(err)
	}
fmt.Println(string(body))
```
### example params & hdrs
```go

```
## Pair type
the Pair struct is used both for headers & URL params
```go
type Pair struct {
	Key string
	Val string
}
```
- both addHdrs and addParams accept a slice of Pair
    - addHdrs also acccepts a `http.Request`, loops through the slice & adds a 
    header to the request for each Pair
    - addParams also accepts the base url, loops through & adds the key value 
    parameter pairs to the url string

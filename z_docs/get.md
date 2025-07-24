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
var hdrs = []Pair{
	{"Accept", "application/json"},
	{"Connection", "keep-alive"},
	{"Referer", "https://www.nba.com"},
	{"Origin", "https://www.nba.com"},
	{"User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/87.0.4280.88 Safari/537.36"},
}

var params = []Pair{
	{"LeagueID", "10"},
	{"PlayerID", "2544"},
}
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

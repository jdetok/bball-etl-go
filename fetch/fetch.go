package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
)

type Pair struct {
	Key string
	Val string
}

func addParams(bUrl string, params []Pair) string {
	var url string = bUrl + "?"
	for i, p := range params {
		url = url + (p.Key + "=" + p.Val)
		if i < len(params)-1 {
			url += "&"
		}
	}
	return url
}

func baseUrl(host string, end string) string {
	return "https://" + host + end
}

func addHdrs(r *http.Request, hdrs []Pair) {
	for _, h := range hdrs {
		r.Header.Add(h.Key, h.Val)
	}
}

func ClientDo(req *http.Request) ([]byte, int, error) {
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Printf("Error occured: %e\n", err)
		if res != nil {
			fmt.Printf("Error Status Code: %d", res.StatusCode)
			return nil, res.StatusCode, err
		}
		return nil, 0, err
	}
	defer res.Body.Close()
	fmt.Printf("Status Code: %d\n", res.StatusCode)
	body, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Printf("Error occured: %e\n", err)
		return nil, res.StatusCode, err
	}
	return body, res.StatusCode, nil
}

func Get(host string, end string, params []Pair, hdrs []Pair) ([]byte, int, error) {
	bUrl := baseUrl(host, end)
	url := addParams(bUrl, params)

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		fmt.Printf("Error occured: %e\n", err)
		return nil, 0, err
	}
	addHdrs(req, hdrs)
	body, status, err := ClientDo(req)
	if err != nil {
		return nil, status, fmt.Errorf("%d: HTTP Request Error: %e", status, err)
	}
	return body, status, nil
}

func main() {
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

	body, _, err := Get("stats.nba.com", "/stats/commonplayerinfo", params, hdrs)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(body))
}

package main

import (
	"fmt"
	"io"
	"net/http"
)

type GetReq struct {
	Host     string
	Endpoint string
	Params   []Pair
	Headers  []Pair
}

type Pair struct {
	Key string
	Val string
}

func (gr *GetReq) GetRespBody() ([]byte, int, error) {
	bUrl := baseUrl(gr.Host, gr.Endpoint)
	url := addParams(bUrl, gr.Params)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		fmt.Printf("Error occured: %e\n", err)
		return nil, 0, err
	}
	addHdrs(req, gr.Headers)
	body, status, err := ClientDo(req)
	if err != nil {
		return nil, status, fmt.Errorf("%d: HTTP Request Error: %e", status, err)
	}
	return body, status, nil
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

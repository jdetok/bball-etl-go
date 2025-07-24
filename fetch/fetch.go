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
	bUrl := gr.baseUrl()
	url := gr.addParams(bUrl)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		fmt.Printf("Error occured: %e\n", err)
		return nil, 0, err
	}
	gr.addHdrs(req)
	body, status, err := ClientDo(req)
	if err != nil {
		return nil, status, fmt.Errorf("%d: HTTP Request Error: %e", status, err)
	}
	return body, status, nil
}

func (gr *GetReq) addParams(bUrl string) string {
	var url string = bUrl + "?"
	for i, p := range gr.Params {
		url = url + (p.Key + "=" + p.Val)
		if i < len(gr.Params)-1 {
			url += "&"
		}
	}
	return url
}

func (gr *GetReq) baseUrl() string {
	return "https://" + gr.Host + gr.Endpoint
}

func (gr *GetReq) addHdrs(r *http.Request) {
	for _, h := range gr.Headers {
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

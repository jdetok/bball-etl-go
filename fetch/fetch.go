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

func (gr *GetReq) GetRespBody() ([]byte, error) {
	bUrl := gr.baseUrl()
	url := gr.addParams(bUrl)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	gr.addHdrs(req)
	body, err := ClientDo(req)
	if err != nil {
		return nil, err
	}
	return body, nil
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

func ClientDo(req *http.Request) ([]byte, error) {
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		if res != nil {
			fmt.Printf("Error Status Code: %d", res.StatusCode)
			return nil, fmt.Errorf(
				"*Response status %d - HTTP client error occured: %e\n",
				res.StatusCode, err)
		}
		return nil, fmt.Errorf(
			"*HTTP client error occured, no response received: %e\n", err)
	}
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Printf("Error occured: %e\n", err)
		return nil, fmt.Errorf(
			"*Response status %d error occured reading response body: %e\n",
			res.StatusCode, err)
	}
	return body, nil
}

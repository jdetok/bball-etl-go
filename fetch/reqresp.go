package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// json response body from stats.nba.com unmarshal into Resp
type Resp struct {
	Resource   string      `json:"resource"`
	Parameters any         `json:"parameters"`
	ResultSets []ResultSet `json:"resultSets"`
}

// main json object in response body after endpoint/params
type ResultSet struct {
	Name    string   `json:"name"`
	Headers []string `json:"headers"`
	RowSet  [][]any  `json:"rowSet"`
}

// build GetReq types to request data from new endpoints
type GetReq struct {
	Host     string
	Endpoint string
	Params   []Pair
	Headers  []Pair
}

// basic key value type
type Pair struct {
	Key string
	Val string
}

func MakeGameLogReq(league string, season string, plTm string,
	dateFrom string, dateTo string) GetReq {
	var gr = GetReq{
		Host:     HOST,
		Headers:  HDRS,
		Endpoint: "/stats/leaguegamelog",
		Params: []Pair{
			{"LeagueID", league},
			{"Season", season},
			{"SeasonType", "Regular+Season"},
			{"Counter", "0"},
			{"Sorter", "DATE"},
			{"Direction", "DESC"},
			{"PlayerOrTeam", plTm},
			{"DateFrom", dateFrom},
			{"DateTo", dateTo},
		},
	}
	return gr
}

// pass a defined GetReq struct, unmarshals body & returns as Resp struct
func RequestResp(gr GetReq) (Resp, error) {
	var resp Resp
	body, err := gr.BodyFromReq()
	if err != nil {
		return resp, fmt.Errorf("error getting response: %e", err)
	}
	resp, err = UnmarshalInto(body)
	if err != nil {
		return resp, fmt.Errorf("error unmarshaling: %e", err)
	}
	return resp, nil
}

/*
pass resp returned from RequestResp
placeholder print `header - val` to console
*/
func ProcessResp(resp Resp) {
	// fmt.Println(resp.ResultSets[0].RowSet[0]...)
	for _, r := range resp.ResultSets[0].RowSet {
		for i, x := range r {
			fmt.Printf("%v: %v\n", resp.ResultSets[0].Headers[i], x)
		}
		fmt.Println("*******")
	}
}

// unmarshal []byte body into Resp struct
func UnmarshalInto(body []byte) (Resp, error) {
	var resp Resp
	if err := json.Unmarshal(body, &resp); err != nil {
		return resp, fmt.Errorf("error unmarshaling: %e", err)
	}
	return resp, nil
}

/*
endptURL to concat endpoint to base url
makeQryStr to loop through gr.Params & make query string
*/
func (gr *GetReq) MakeFulLURL() string {
	bUrl := gr.endptURL()
	return gr.makeQryStr(bUrl)
}

/*
make new request with url returned from MakeFullURL
add gr.Headers to req with addHdrs
use RespFromClient to do the http req, return the resp body []byte
*/
func (gr *GetReq) BodyFromReq() ([]byte, error) {
	req, err := http.NewRequest(http.MethodGet, gr.MakeFulLURL(), nil)
	if err != nil {
		return nil, err
	}
	gr.addHdrs(req)
	body, err := RespFromClient(req)
	if err != nil {
		return nil, err
	}
	return body, nil
}

// concat endpoint to host
func (gr *GetReq) endptURL() string {
	return "https://" + gr.Host + gr.Endpoint
}

// makes the query string from gr.Params
func (gr *GetReq) makeQryStr(bUrl string) string {
	var url string = bUrl + "?"
	for i, p := range gr.Params {
		url = url + (p.Key + "=" + p.Val)
		if i < len(gr.Params)-1 {
			url += "&"
		}
	}
	return url
}

// loop through gr.Headers & add each as a header to the request
func (gr *GetReq) addHdrs(r *http.Request) {
	for _, h := range gr.Headers {
		r.Header.Add(h.Key, h.Val)
	}
}

/*
use http client to perform http request
get & return body as []byte
*/
func RespFromClient(req *http.Request) ([]byte, error) {
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		if res != nil {
			fmt.Printf("Error Status Code: %d", res.StatusCode)
			return nil, fmt.Errorf(
				"*Response status %d - HTTP client error occured: %e",
				res.StatusCode, err)
		}
		return nil, fmt.Errorf(
			"*HTTP client error occured, no response received: %e", err)
	}
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Printf("Error occured: %e\n", err)
		return nil, fmt.Errorf(
			"*Response status %d error occured reading response body: %e",
			res.StatusCode, err)
	}
	return body, nil
}

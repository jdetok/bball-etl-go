package main

import (
	"encoding/json"
	"fmt"
)

type Resp struct {
	Resource   string      `json:"resource"`
	Parameters any         `json:"parameters"`
	ResultSets []ResultSet `json:"resultSets"`
}

type ResultSet struct {
	Name    string   `json:"name"`
	Headers []string `json:"headers"`
	RowSet  [][]any  `json:"rowSet"`
}

func RequestResp(gr GetReq) (Resp, error) {
	var resp Resp
	body, err := gr.GetRespBody()
	if err != nil {
		return resp, fmt.Errorf("error getting response: %e", err)
	}

	resp, err = UnmarshalInto(body)
	if err != nil {
		return resp, fmt.Errorf("error unmarshaling: %e", err)
	}
	return resp, nil
}

func UnmarshalInto(body []byte) (Resp, error) {
	var resp Resp
	if err := json.Unmarshal(body, &resp); err != nil {
		return resp, fmt.Errorf("error unmarshaling: %e", err)
	}
	return resp, nil
}

func ProcessResp(resp Resp) {
	fmt.Println(resp.ResultSets[0].RowSet[0]...)
	for _, r := range resp.ResultSets[0].RowSet {
		for i, x := range r {
			fmt.Printf("%v: %v\n", resp.ResultSets[0].Headers[i], x)
		}
		fmt.Println("*******")
	}
}

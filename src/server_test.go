package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strconv"
	"testing"
)

// WrongTypedPair is the pair with wrong input type
type WrongTypedPair struct {
	A string
	B string
}

func newRequest(method, url string, body interface{}) (*http.Request, error) {
	var r *http.Request
	var buf bytes.Buffer
	r, err := http.NewRequest(method, url, nil)
	if err != nil {
		return nil, err
	}
	if err := json.NewEncoder(&buf).Encode(body); err != nil {
		return nil, err
	}
	r.Body = ioutil.NopCloser(&buf)
	r.Header.Set("Content-Type", "application/json;charset=UTF-8")
	return r, nil
}

func decodeSumResponse(resp *http.Response) (interface{}, error) {
	var response AddResult
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, err
	}
	return &response, nil
}

func TestAdd(t *testing.T) {
	router := new(addHandler)
	ts := httptest.NewServer(router)
	defer ts.Close()

	type want struct {
		ResponseCode int
		Response     interface{}
	}
	var testCases = []struct {
		Label   string
		Request interface{}
		Want    want
	}{
		{
			Label: "Post two integers to obtain their sum",
			Request: &Pair{
				A: 1,
				B: 2,
			},
			Want: want{
				ResponseCode: http.StatusOK,
				Response:     &AddResult{Sum: 3},
			},
		},
		{
			Label: "Post non-integer inputs",
			Request: &WrongTypedPair{
				A: "hello",
				B: "world",
			},
			Want: want{
				ResponseCode: http.StatusBadRequest,
			},
		},
	}

	for i, tc := range testCases {
		t.Run(strconv.Itoa(i)+"/"+tc.Label, func(t *testing.T) {
			var req *http.Request
			var resp *http.Response
			var err error
			req, err = newRequest("POST", ts.URL+"/", tc.Request)
			if err != nil {
				t.Error(err)
			}
			resp, err = http.DefaultClient.Do(req)
			if err != nil {
				t.Error(err)
			}
			defer resp.Body.Close()
			if tc.Want.Response != nil {
				res, err := decodeSumResponse(resp)
				if err != nil {
					t.Error(err)
				}
				if !reflect.DeepEqual(res, tc.Want.Response) {
					t.Errorf("Response %d %v; want %v", i, resp, tc.Want.Response)
				} else {
					t.Logf("A: %v, B: %v, Sum: %v\n", tc.Request.(*Pair).A, tc.Request.(*Pair).B, res.(*AddResult).Sum)
				}
			}
			if resp.StatusCode != tc.Want.ResponseCode {
				t.Errorf("Response code %v; want %v", resp.StatusCode, tc.Want.ResponseCode)
			}
		})
	}
}

package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"os"
	"reflect"
	"testing"
)

const (
	SortJson1 = `
{
	"array": [1, 3, 2, 3, 4],
	"uniq": true
}`
	SortJson2 = `
{
	"array": [1, 2, 0, 4, 5, 5, 5],
	"uniq": false
}`
	SortJson3 = `
{
	"array": [],
	"uniq": true
}`
)

func TestMain(m *testing.M) {
	go InitServer()
	code := m.Run()
	os.Exit(code)
}

func TestSortRouteWithUniq(t *testing.T) {
	shouldBeArray := []int {
		1, 2, 3, 4,
	}
	body := bytes.NewReader([]byte(SortJson1))
	resp, err := http.DefaultClient.Post("http://localhost:8080/api/sort", "application/json", body)
	if err != nil {
		t.Fatal(err)
		t.Fail()
	}
	var sortResponse SortResponse
	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&sortResponse)
	if !reflect.DeepEqual(sortResponse.Array, shouldBeArray) {
		t.Fatal("Not equal")
		t.Fail()
	}
}

func TestSortRouteWithoutUniq(t *testing.T) {
	shouldBeArray := []int {
		0, 1, 2, 4, 5, 5, 5,
	}
	body := bytes.NewReader([]byte(SortJson2))
	resp, err := http.DefaultClient.Post("http://localhost:8080/api/sort", "application/json", body)
	if err != nil {
		t.Fatal(err)
		t.Fail()
	}
	var sortResponse SortResponse
	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&sortResponse)
	if !reflect.DeepEqual(sortResponse.Array, shouldBeArray) {
		t.Fatal("Not equal")
		t.Fail()
	}
}

func TestSortRouteEmpty(t *testing.T) {
	body := bytes.NewReader([]byte(SortJson3))
	resp, err := http.DefaultClient.Post("http://localhost:8080/api/sort", "application/json", body)
	if err != nil {
		t.Fatal(err)
		t.Fail()
	}
	if resp.StatusCode != http.StatusBadRequest {
		t.Fatal("Status is not 400")
		t.Fail()
	}
}

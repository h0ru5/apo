package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	req, err := http.NewRequest("GET", "http://localhost:3000/token", nil)
	check(err)
	req.SetBasicAuth("foo", "bar")
	res, err := http.DefaultClient.Do(req)
	check(err)

	defer res.Body.Close()
	token, err := ioutil.ReadAll(res.Body)
	check(err)

	fmt.Println(string(token))

	req, err = http.NewRequest("GET", "http://localhost:3001/open", nil)
	check(err)
	req.Header.Add("Authorization", "Bearer "+string(token))
	res, err = http.DefaultClient.Do(req)
	check(err)

	defer res.Body.Close()
	result, err := ioutil.ReadAll(res.Body)
	check(err)

	fmt.Println(string(result))
}

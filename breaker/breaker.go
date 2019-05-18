package main

import (
	"fmt"
	"log"
	"github.com/sony/gobreaker"
	"io/ioutil"
	"net/http"
)

var cb *gobreaker.CircuitBreaker

func init()  {
	var st gobreaker.Settings
	st.Name = "HTTP GET"
	st.ReadyToTrip = func(counts gobreaker.Counts) bool {
		failureRatio := float64(counts.TotalFailures) / float64(counts.Requests)
		return counts.Requests >= 3 && failureRatio >= 0.6
	}
	cb = gobreaker.NewCircuitBreaker(st)
}

func Get(url string) ([]byte, error) {
	body, err := cb.Execute(func() (i interface{}, e error) {
		resp, err := http.Get(url)
		if err != nil{
			return nil, err
		}
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil{
			return nil, err
		}
		return body, nil
	})
	if err != nil{
		return nil, err
	}
	return body.([]byte), nil
}

func main()  {
	body, err := Get("http://www.google.com/robots.txt")
	if err != nil{
		log.Fatal(err)
	}
	fmt.Println(string(body))

}
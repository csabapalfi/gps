package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"time"
)

var myClient = &http.Client{Timeout: 10 * time.Second}

func debug(data []byte, err error) {
	if err == nil {
		fmt.Printf("%s\n\n", data)
	} else {
		log.Fatalf("%s\n\n", err)
	}
}

func getJson(url string, target interface{}) error {
	response, err := myClient.Get(url)
	if err != nil {
		return err
	}
	defer response.Body.Close()
	debug((httputil.DumpResponse(response, true)))

	return json.NewDecoder(response.Body).Decode(target)
}

type PageSpeedResult struct {
	RuleGroups struct {
		Speed struct {
			Score int
		}
	}
}

func main() {
	result := &PageSpeedResult{}
	getJson("https://www.googleapis.com/pagespeedonline/v2/runPagespeed?url=http://www.example.com&strategy=mobile", result)
	println(result.RuleGroups.Speed.Score)
}

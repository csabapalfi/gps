package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
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

func getJson(pageSpeedUrl string, target interface{}) error {
	response, err := myClient.Get(pageSpeedUrl)
	if err != nil {
		return err
	}
	defer response.Body.Close()
	debug((httputil.DumpResponse(response, true)))

	return json.NewDecoder(response.Body).Decode(target)
}

func buildPageSpeedUrl(targetUrl string, strategy string) string {
	pageSpeedUrl := &url.URL{
		Host:   "www.googleapis.com",
		Scheme: "https",
		Path: "pagespeedonline/v2/runPagespeed",
	}
	q := pageSpeedUrl.Query()
	q.Set("url", targetUrl)
	q.Set("strategy", strategy)
	pageSpeedUrl.RawQuery = q.Encode()
	return pageSpeedUrl.String();
}

type PageSpeedResult struct {
	RuleGroups struct {
		Speed struct {
			Score int
		}
	}
}

func main() {
	flag.Parse()
	var args = flag.Args()
	targetUrl, slackChannel := args[0], args[1]

	result := &PageSpeedResult{}
	getJson(buildPageSpeedUrl(targetUrl, "mobile"), result)
	println(result.RuleGroups.Speed.Score)
	println(slackChannel)
}

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

func buildPageSpeedUrl(targetUrl string, strategy string) string {
	pageSpeedUrl := &url.URL{
		Host:   "www.googleapis.com",
		Scheme: "https",
		Path:   "pagespeedonline/v2/runPagespeed",
	}
	q := pageSpeedUrl.Query()
	q.Set("url", targetUrl)
	q.Set("strategy", strategy)
	pageSpeedUrl.RawQuery = q.Encode()
	return pageSpeedUrl.String()
}

func getScore(targetUrl string, strategy string, verbose bool) (error, int) {
	pageSpeedUrl := buildPageSpeedUrl(targetUrl, strategy)
	response, err := myClient.Get(pageSpeedUrl)
	if err != nil {
		return err, -1
	}
	defer response.Body.Close()
	if verbose {
		debug((httputil.DumpResponse(response, true)))
	}

	result := &PageSpeedResult{}
	json.NewDecoder(response.Body).Decode(result)
	return nil, result.RuleGroups.Speed.Score
}

func getIcon(score, threshold int) string {
	if score >= threshold {
		return "✅"
	}
	return "❌"
}

type PageSpeedResult struct {
	RuleGroups struct {
		Speed struct {
			Score int
		}
	}
}

func main() {
	verbose := flag.Bool("v", false, "verbose")
	mobileThreshold := flag.Int("mobile", 0, "mobile threshold")
	desktopThreshold := flag.Int("desktop", 0, "desktop threshold")
	flag.Parse()
	var args = flag.Args()
	targetUrl, slackChannel := args[0], args[1]

	_, mobileScore := getScore(targetUrl, "mobile", *verbose)
	_, desktopScore := getScore(targetUrl, "desktop", *verbose)

	mobileIcon := getIcon(mobileScore, *mobileThreshold)
	desktopIcon := getIcon(desktopScore, *desktopThreshold)

	message := fmt.Sprintf("📱 %d %s  🖥 %d %s", mobileScore, mobileIcon, desktopScore, desktopIcon)

	println(slackChannel, message)
}

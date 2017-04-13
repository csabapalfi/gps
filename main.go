package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
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

func getJson(url string, target interface{}, verbose bool) error {
	response, err := myClient.Get(url)
	if err != nil {
		return err
	}
	defer response.Body.Close()
	if verbose {
		debug((httputil.DumpResponse(response, true)))
	}
	return json.NewDecoder(response.Body).Decode(target)
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

func getPageSpeedScore(targetUrl string, strategy string, verbose bool) int {
	type PageSpeedResult struct {
		RuleGroups struct {
			Speed struct {
				Score int
			}
		}
	}
	pageSpeedUrl := buildPageSpeedUrl(targetUrl, strategy)
	result := &PageSpeedResult{}
	err := getJson(pageSpeedUrl, result, verbose)
	if err != nil {
		handleError(err)
	}
	return result.RuleGroups.Speed.Score
}

func buildSlackUrl(channel string, message string) string {
	slackUrl := &url.URL{
		Host:   "slack.com",
		Scheme: "https",
		Path:   "api/chat.postMessage",
	}
	token := os.Getenv("SLACK_TOKEN")
	q := slackUrl.Query()
	q.Set("token", token)
	q.Set("channel", channel)
	q.Set("text", message)
	q.Set("as_user", "true")
	slackUrl.RawQuery = q.Encode()
	return slackUrl.String()
}

func postToSlack(channel string, message string, verbose bool) {
	type SlackResult struct {
		Ok    bool
		Error string
	}
	url := buildSlackUrl(channel, message)
	result := &SlackResult{}
	err := getJson(url, result, verbose)
	if err != nil {
		handleError(err)
	}
	if err == nil && !result.Ok {
		handleError(errors.New(result.Error))
	}
}

func getIcon(score, threshold int) string {
	if score >= threshold {
		return "✅"
	}
	return "❌"
}

func handleError(err error) {
	log.Fatalf("%s\n\n", err)
	os.Exit(1)
}

func main() {
	verbose := flag.Bool("v", false, "verbose")
	mobileThreshold := flag.Int("mobile", 0, "mobile threshold")
	desktopThreshold := flag.Int("desktop", 0, "desktop threshold")
	flag.Parse()
	var args = flag.Args()
	targetUrl, slackChannel := args[0], args[1]

	mobileScore := getPageSpeedScore(targetUrl, "mobile", *verbose)
	desktopScore := getPageSpeedScore(targetUrl, "desktop", *verbose)

	mobileIcon := getIcon(mobileScore, *mobileThreshold)
	desktopIcon := getIcon(desktopScore, *desktopThreshold)

	slackMessage := fmt.Sprintf("📱 %d %s  🖥 %d %s", mobileScore, mobileIcon, desktopScore, desktopIcon)
	postToSlack(slackChannel, slackMessage, *verbose)
	fmt.Printf("Posted to %s: %s", slackChannel, slackMessage)
}

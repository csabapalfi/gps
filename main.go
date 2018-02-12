package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"time"
)

var myClient = &http.Client{Timeout: 10 * time.Second}

func debug(data []byte, err error) {
	if err == nil {
		fmt.Fprintf(os.Stderr, "%s\n\n", data)
	} else {
		fmt.Fprintf(os.Stderr, "%s\n\n", err)
	}
}

func getJSON(url string, target interface{}, verbose bool) error {
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

func buildPageSpeedURL(targetURL string, strategy string) string {
	pageSpeedURL := &url.URL{
		Host:   "www.googleapis.com",
		Scheme: "https",
		Path:   "pagespeedonline/v4/runPagespeed",
	}
	q := pageSpeedURL.Query()
	q.Set("url", targetURL)
	q.Set("strategy", strategy)
	pageSpeedURL.RawQuery = q.Encode()
	return pageSpeedURL.String()
}

func getPageSpeedScore(targetURL string, strategy string, verbose bool) int {
	type PageSpeedResult struct {
		RuleGroups struct {
			Speed struct {
				Score int
			}
		}
	}
	pageSpeedURL := buildPageSpeedURL(targetURL, strategy)
	result := &PageSpeedResult{}
	err := getJSON(pageSpeedURL, result, verbose)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n\n", err)
		os.Exit(1)
	}
	return result.RuleGroups.Speed.Score
}

func getIcon(score, threshold int) string {
	if score >= threshold {
		return "✅"
	}
	return "❌"
}

func main() {
	verbose := flag.Bool("v", false, "verbose output")
	mobileThreshold := flag.Int("mobile", 0, "mobile threshold")
	desktopThreshold := flag.Int("desktop", 0, "desktop threshold")
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage of gps:\n")
		fmt.Fprintf(os.Stderr, "gps [flags...] <url>\n")
		flag.PrintDefaults()
	}

	flag.Parse()
	targetURL := flag.Arg(0)

	mobileScore := getPageSpeedScore(targetURL, "mobile", *verbose)
	desktopScore := getPageSpeedScore(targetURL, "desktop", *verbose)

	mobileIcon := getIcon(mobileScore, *mobileThreshold)
	desktopIcon := getIcon(desktopScore, *desktopThreshold)

	fmt.Printf("📱 %d %s  🖥 %d %s  %s", mobileScore, mobileIcon, desktopScore, desktopIcon, targetURL)
}

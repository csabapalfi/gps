# Google PageSpeed -> Slack (WIP)

Post Google PageSpeed scores to Slack.

The idea is to run this from your CI/CD periodically and/or after deployments.

## usage

```sh
GPS_SLACK_TOKEN=1212 gps "http://example.com" "#my-channel"
```

## todo

* [ ] parse args
* [ ] run both mobile/desktop strategy
* [ ] add threshold flags
* [ ] build result message
* [ ] post to slack

# Google PageSpeed -> Slack (WIP)

Post *G*oogle *P*ageSpeed scores to *S*lack.

The idea is to run this from your CI/CD periodically and/or after deployments.

I built a node version of this for one of my clients but want to have a lighter single-binary version.

## usage

```sh
GPS_SLACK_TOKEN=1212 gps "http://example.com" "#my-channel"
```

## todo

* [x] basic PageSpeed request
* [x] parse score from PageSpeed response
* [ ] parse args
* [ ] run both mobile/desktop strategy
* [ ] add threshold flags
* [ ] build result message
* [ ] post to slack

# Google PageSpeed Score CLI

Check **G**oogle **P**ageSpeed **S**cores from the command line.

The idea is to run this from your CI/CD periodically and/or after deployments (and post the output message to Slack by piping it to slackcat)

I built a node version of this for one of my clients but want to have a lighter single-binary version.

## installation

Grab the [latest release from Github](https://github.com/csabapalfi/gps/releases/latest).

## basic usage

```sh
gps "http://example.com"
```
This should output something like this:
```
📱 99 ✅  🖥 99 ✅
```
If you want to post to Slack then check out the brilliant [slackcat](https://github.com/crewjam/slackcat) tool:
```sh
gps "http://example.com" | slackcat -tee -token=$SLACK_TOKEN -channel=$YOUR_CHANNEL
```


## defining thresholds

You can define thresholds to check your score using the `mobile` and `desktop` flags:
```sh
gps -mobile 100 -desktop 90 "http://example.com"
```
Scores above the threshold get a green tick, scores below get a red X:
```
📱 99 ❌  🖥 99 ✅
```

## verbose output

To log PageSpeed API responses just pass the `-v` flag:

```sh
gps -v "http://example.com"
```
```
HTTP/2.0 200 OK
...
{
 "kind": "pagespeedonline#result",
 "id": "http://example.com",
 "responseCode": 200,
...
HTTP/2.0 200 OK
...
{
 "kind": "pagespeedonline#result",
 "id": "http://example.com",
 "responseCode": 200,
...
📱 99 ✅  🖥 99 ✅
```

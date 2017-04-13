# Google PageSpeed -> Slack

Post *G*oogle *P*ageSpeed scores to *S*lack.

The idea is to run this from your CI/CD periodically and/or after deployments.

I built a node version of this for one of my clients but want to have a lighter single-binary version.

## basic usage

A Slack token is expected to be set as the `SLACK_TOKEN` environment variable.
```sh
export SLACK_TOKEN=1212
gps "http://example.com" "my-channel"
```
This should output something like this:
```
Posted to my-channel: 📱 99 ✅  🖥 99 ✅
```
...and of course the message should pop up on Slack, too.


## defining thresholds

You can define thresholds to check your score using the `mobile` and `desktop` flags:
```sh
gps -mobile 100 -desktop 90 "http://example.com" "my-channel"
```
Scores above the threshold get a green tick, scores below get a red X:
```
Posted to my-channel: 📱 99 ❌  🖥 99 ✅
```

## verbose output

To log PageSpeed and Slack API responses just pass the `-v` flag:

```sh
gps -v "http://example.com" "my-channel"
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
HTTP/1.1 200 OK
...
{"ok":true,
...
Posted to my-channel: 📱 99 ✅  🖥 99 ✅
```

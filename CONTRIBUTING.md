# Contributing to unaspoller

## PRs

I'm not looking for PRs (code or docs) at the moment as I'm still working on this myself.

I am, however, happy to receive ideas, suggestions, feedback (good or bad) via issues in this repo.

## Probe data

You can perform a single pass of the known URLs against a target, it takes about 20 seconds:
```
$ ./unaspoller -target=192.168.1.17 -verifyssl=false -probe
21:43:14.471 main ▶ INFO 001 unaspoller version v0.0.1
21:43:14.596 LoginUNAS ▶ INFO 002 login successful
21:43:14.596 doProbeKnownURLs ▶ INFO 003 performing probe of 10 known URLs and then exiting
21:43:14.596 doProbeKnownURLs ▶ INFO 004 Output will be saved in [unaspoller-probe.txt]
21:43:34.227 Probe ▶ INFO 005 Probe mode finished. Everything looks good!
21:43:34.227 Probe ▶ INFO 006 Your data may still be useful so please consider reading https://github.com/alexgreenbank/unaspoller/CONTRIBUTING.md
21:43:34.227 main ▶ INFO 007 Exiting after probe mode 
```

The `unaspoller-probe.txt` file will then contain the JSON from the various endpoints, one line per endpoint.
```
$ cat unaspoller-probe.txt | cut -b 1-120
PROBE:version=[v0.0.1]
PROBE:[/proxy/drive/api/v2/systems/device-info]: resp=[{"networkInterfaces":[{"interface":"ethernet","interfaceName":"en
PROBE:[/proxy/drive/api/v1/systems/performance/file-operations]: resp=[{"err":null,"type":"single","data":{"busy":false}
PROBE:[/proxy/drive/api/v1/systems/storage?type=detail]: resp=[{"err":null,"type":"single","data":{"diskInfo":{"needMore
PROBE:[/proxy/drive/api/v2/systems/disk-stats?start=1777840100&end=1777927400&interval=900]: resp=[{"series":{"disks":[{
PROBE:[/proxy/drive/api/v2/systems/network-io]: resp=[{"receiveKBPS":1.3940776826437187,"transmitKBPS":0.641392883401206
PROBE:[/proxy/users/drive/api/v1/systems/identity]: resp=[{"err":null,"type":"single","data":{"autoSendInvitations":fals
PROBE:[/proxy/users/drive/api/v1/systems/info]: resp=[{"err":null,"type":"single","data":{"console":{"deviceId":"6C63F85
PROBE:[/proxy/users/drive/api/v2/drives]: resp=[{"drives":[{"id":"41c20687-6c0b-4bc5-af30-abce54982b84","type":"shared",
PROBE:[/proxy/users/drive/api/v2/groups]: resp=[{"groups":[]}]
PROBE:[/proxy/users/drive/api/v2/storage]: resp=[{"pools":[{"number":1,"id":"aa460908-1e83-4acb-ab65-436913517d61","pref
PROBE:DONE
```

If you are going to forward this data on then please make sure you check it for IP addresses, PII, folder names, or any other information you do not want sharing.

Once you're happy with it you can [create an issue](https://github.com/alexgreenbank/unaspoller/issues) with the subject `New Probe Data`, include the contents of the `unaspolller-probe.txt` file, and I will take a look at it, make the necessary adjustments and close the issue once resolved.

This kind of data will help enormously as I have no idea what other field values/data are out there as I only have a single UNAS device with a limited set of data.

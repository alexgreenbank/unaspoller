# unaspoller

UNAS Prometheus Exporter

Following on from https://github.com/unpoller/unpoller/issues/785 I thought I'd have a go at creating a Prometheus exporter/poller for the UNAS Pro units.

This connects directly to the Drive application on the device.

As noted in the issue linked above I can't see a simple way to get this info via the Remote API (given that the API Key scope only covers "Network" and "Protect" and not "Drive").

The console on ui.com is able to access the drive data remotely, but I can't work out from a 10 minute play with Chrome Developer Tools how this is done.

For now I'll just make a poller that can poll one UNAS device (which should easily be extendible to n devices) and expose the resulting metrics for Prometheus to scrape.

If you just want to dive in then go to the [Getting Started](#usage) section.

## Drive API and Metrics

The metrics produced are documented at [METRICS.md](METRICS.md) along with details of what I've found out about the Drive API.

## Current Status

Right now (v0.0.1) the poller is very limited. There are two major areas it needs work on right now:

### Limited data for API responses

The Drive API is not public so I'm having to work everything out from responses that I see, but I obviously don't get to see all of the possible responses.

I only have access to one UNAS device (a UNAS Pro). It doesn't have any cache slots. It only has one type of pool. All of the drives are healthy.

I don't have any example API responses for whole sections of the API surface. Joy.

### Parsing JSON is often not as simple as you may think

Anyone who has written JSON parsing code will know the fun of an endpoint that returns `{"isValid": true}` and a different endpoint returns `{"isValid": "true"}`. Having to deal with this in golang is a chore, but it's not impossible.

Again, with no documented API I've no idea what is guaranteed to be returned, what is optional, etc.

The current code is not very tolerant of problems.

### How you can help

The thing I need the most is access to as much response data as I can get from different devices. Devices with unhealthy pools. Devices with cache slots populated. Devices in the middle of some kind of operation (repairing, etc).

I'm not looking for code/documentation PRs right now. (If only to make it easier in terms of handing over the code to other projects.)

Whilst the poller may work well for you (it should get more and more tolerant/capable as I am able to see and fix problems) it would be good if you can report any JSON dumps that represent unknown/unhandled/unseen API responses. It does not send them to me automatically.

More info at [CONTRIBUTING.md](CONTRIBUTING.md).

## TODO

* Command line flags
    * Metrics host (for listening, currently defaults to 0.0.0.0)
    * Login credentials in a file?
* Add remaining API endpoints to poll (drives, etc)
* VALIDATE: find all values so far from own data...
* PROBE MODE: Validate JSON
    * PROBE MODE: Validate values
* Scrub data before logging
* Also log similar dumps when running if data not as expected
* Option to send log to file?
* Sort out logging of unexpected/unparsed data outputs
* Document things (specifically API and example responses in METRICS.md)
* errors output with %w?

* Clean up the `_test.go` files and add them to repo
* Get feedback with real data (and bugs)
* Find more possible endpoints (use developer console in Chrome and browse around UI lots)
* Different poll frequency for different endpoints?
* Handle multiple targets?
    * Via a config file simplifies chance of different usernames/passwords
    * label per target
* Dockerfile

## Usage

Either build using `go build .` and then run as `./unaspoller <options>` or just run with `go run . <options>`.

### Creating a user in the Unifi Drive app

The poller needs to login in the UNAS device to gather the data. I've used `unifipoller` as a username but you can choose whatever you want. Same with the `First Name` and `Last Name` fields.

1. In the Drive UI go to `Admin & Users` (`https://<IP>/drive/admins/users`) and click on the `Create New` button.
1. I used `First Name` = `Unifi`, `Last Name` = `Poller`, no email
1. Select the `Admin` checkbox
1. Select the `Restrict to Local Access Only` checkbox (which opens up the ability to set a username and password for this new user)
1. `Username` = `unifipoller` and give it a password (you'll need to remeber this password to provide it to the `unaspoller` when you run it.
1. The `Use a Predefined Role` checkbox should be selected, you want the user to be a `Super Admin` (TODO - investigate if lower privs will suffice)
1. Everything else can stay as default.
1. Click the `Create` button

That should be everything setup on the UNAS side.

### Probe mode

To avoid putting the password on the command line you can put it in the `UNAS_PW` environment variable before running `unaspoller`:-

```
 export UNAS_PW=SuperSecret123
```

If you just want to run a single pass of the known URLs against a target:
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

### Running the exporter

Assuming you can get data in `-probe` mode then switching to the exporter mode is easy:
```
$ ./unaspoller -target=192.168.1.17 -verifyssl=false
21:48:31.650 main ▶ INFO 001 unaspoller version v0.0.1
21:48:31.781 LoginUNAS ▶ INFO 002 login successful
...
```

It's not very chatty at all if not in debug mode, adding `-debug` makes it more verbose:
```
$ ./unaspoller -target=192.168.1.17 -verifyssl=false -debug
21:49:52.126 main ▶ INFO 001 unaspoller version v0.0.1
21:49:52.126 registerAPIDef ▶ DEBU 002 registering apidef [/proxy/drive/api/v2/storage]
21:49:52.126 registerAPIDef ▶ DEBU 003 registering apidef [/proxy/drive/api/v2/systems/device-info]
21:49:52.126 LoginUNAS ▶ DEBU 004 Attempt 1/5 to login
21:49:52.249 LoginUNAS ▶ INFO 005 login successful
21:49:52.250 main ▶ DEBU 006 Logged in as 'unifipoller'
21:49:52.250 mainPollLoop ▶ DEBU 007 Starting poll
21:49:52.250 mainPollLoop ▶ DEBU 008 Polling /proxy/drive/api/v2/storage
21:49:52.250 doRequest ▶ DEBU 009 Attempt 1/5 to perform GET https://192.168.1.17/proxy/drive/api/v2/storage
21:49:52.302 doDriveAPIDef ▶ DEBU 00a unmarshalling /proxy/drive/api/v2/storage
21:49:52.302 doDriveAPIDef ▶ DEBU 00b unmarshalled /proxy/drive/api/v2/storage ok
21:49:52.302 mainPollLoop ▶ DEBU 00c Polling /proxy/drive/api/v2/systems/device-info
21:49:52.302 doRequest ▶ DEBU 00d Attempt 1/5 to perform GET https://192.168.1.17/proxy/drive/api/v2/systems/device-info
21:49:52.309 doDriveAPIDef ▶ DEBU 00e unmarshalling /proxy/drive/api/v2/systems/device-info
21:49:52.309 doDriveAPIDef ▶ DEBU 00f unmarshalled /proxy/drive/api/v2/systems/device-info ok
21:49:52.309 mainPollLoop ▶ DEBU 010 Sleeping for 15s between polls
```

Once it is running you can see what metrics are being produced using curl of the `/metrics` endpoint of the host that is running the exporter:
```
$ curl -s http://192.168.1.189:8090/metrics
# HELP unas_cpu_load CPU load of UNAS device
# TYPE unas_cpu_load gauge
unas_cpu_load 0.09699999999999999
# HELP unas_cpu_temperature Temperature of CPU in deg C
# TYPE unas_cpu_temperature gauge
unas_cpu_temperature 63
# HELP unas_disk_bad_sector_count The count of bad sectors reported by the disk
# TYPE unas_disk_bad_sector_count gauge
unas_disk_bad_sector_count{serial="N8GUTKVY",slotId="2"} 0
...
```

### Configuring Prometheus

TODO

### All options

```
Usage of ./unaspoller:
  -debug
    	enable debug mode
  -lograwjson
    	log the raw JSON received even if it is good
  -max429retries int
    	Maximum number of retries if 429 received (default 5)
  -metricport int
    	Port to listen on for /metrics endpoint (default 8090)
  -metrixprefix string
    	Prefix for metrics (default "unas")
  -password string
    	password to connect to UNAS device
  -pollinterval string
    	Sleep duration between polls (default "15s")
  -probe
    	probe known URLs once dumping info and stopping
  -probefile string
    	Filename for -probe mode output (default "unaspoller-probe.txt")
  -probeinterval string
    	Sleep duration between polls in probe mode (default "2s")
  -retrydelay string
    	Time to wait until retry if 429 received (default "2s")
  -target string
    	IP address of UNAS device (default "192.168.1.1")
  -username string
    	username to connect to UNAS device (default "unifipoller")
  -verifyssl
    	verify SSL certificates (default true)
  -version
    	print version information
```

## FAQs

Q: I've got more than one UNAS device, how can I distinguish between the metrics for each of them?

A: If you scrape the exporter from Prometheus then you'll already be getting an `instance` label filled in with the `target` value used for the poll.

If you need something else (like a hostname, internal name, etc) then you can add a label in the scrape target in `prometheus.yml`, for example:
```
  - job_name: unas01
    # Scrape metrics for UNAS (192.168.1.17)
    static_configs:
      - targets: ['192.168.1.11:8090']
        labels:
          deviceip: '192.168.1.17'
```

Notes:
* `192.168.1.17` is the IP of my UNAS Pro. `192.168.1.11` is the IP of my homelab server I run exporters/Prometheus/etc on
* (FUTURE) When/if I get round to adding an option to poll multiple UNAS devices from a single exporter instance I will include the ability to set a unique label/value for each target

## Acknowledgements

The https://unpoller.com/ project (also https://github.com/unpoller/) for the inspiration along with https://github.com/unpoller/unpoller/issues/785 that piqued my interest.

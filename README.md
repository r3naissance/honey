# honey
Incredibly lightweight honeypot for anomalous traffic
### 
![Run](https://github.com/r3naissance/honey/blob/main/img/run.gif)
## TL;DR
### Benefits
- Less than 120 lines of code
- Analyze traffic behavior
- Alert to network mapping (IDS supplement)
- JSON output for easy consumption/push
- Connection is dropped immediately after data is logged
- OS agnostic
### TODO
- UDP support
## Installation
### Build from source
```
git clone https://github.com/r3naissance/honey
cd honey
go build .
```
### Install it
```
go install github.com/r3naissance/honey@latest
```
## Usage
### No options
```
└─$ honey
ERRO[0000] You must provide at least 1 port
honey - Usage:
  -output string
        Where to save the logs
         -output /tmp/honey.log (default "stdout")
  -ports string
        Ports to listen to (comma separated and/or ranges)
         > Pro Tip: nmap -oG - -v --top-ports 25 2>/dev/null | grep Ports
         -ports 21,22,80,443
         -ports 20-25,80,81,8000-8100
```
### Running to stdout
```
└─$ honey -ports 8000,8001,8080-8100
{"level":"info","msg":"Honey Options","output":"stdout","ports":"8000,8001,8080-8100","time":"2023-04-28T16:53:03-06:00"}
```
### Running to log
```
└─$ ./honey -ports 8000,8001,8080-8100 -output /tmp/honey.log &
[1] 3005205

└─$ tail -f /tmp/honey.log
{"level":"info","msg":"Honey Options","output":"/tmp/honey.log","ports":"8000,8001,8080-8100","time":"2023-04-28T17:02:28-06:00"}
```

# Airbrake Hook for Logrus <img src="http://i.imgur.com/hTeVwmJ.png" width="40" height="40" alt=":walrus:" class="emoji" title=":walrus:" />&nbsp;[![Build Status](https://travis-ci.org/Invoiced/logrus-airbrake-hook.svg?branch=master)](https://travis-ci.org/Invoiced/logrus-airbrake-hook)&nbsp;[![godoc reference](https://godoc.org/github.com/gemnasium/logrus-airbrake-hook?status.png)](https://godoc.org/gopkg.in/gemnasium/logrus-airbrake-hook.v2)

Use this hook to send your errors to [Airbrake](https://airbrake.io/).
This hook is using the [official airbrake go package](https://github.com/airbrake/gobrake), and will hit the api V3.
The hook is synchronous and will send the error for `log.Error`, `log.Fatal` and `log.Panic` levels.

All logrus fields will be sent as context fields on Airbrake.

## Usage

The hook must be configured with:

* A project ID (found in your your Airbrake project settings)
* An API key ID (found in your your Airbrake project settings)
* The name of the current environment ("development", "staging", "production", ...)
* The Stack Level Trace You Want to Begin With
* If you want Airbrake to fire synchronous or not

```go
import (
    "log/syslog"
    "github.com/Invoiced/logrus"
    "github.com/Invoiced/logrus-airbrake-hook" // the package is named 
    )

func main() {
    log.AddHook(airbrake.NewHook(123, "xyz", "production",3,true))
    log.Error("some logging message") // The error is sent to airbrake in background
}
```



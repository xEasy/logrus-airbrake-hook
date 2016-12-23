package airbrake

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/Sirupsen/logrus"
	"gopkg.in/airbrake/gobrake.v2"
)

// AirbrakeHook to send exceptions to an exception-tracking service compatible
// with the Airbrake API.
type airbrakeHook struct {
	Airbrake        *gobrake.Notifier
	StackTraceLevel int
	Synchronous     bool
}

type Error struct {
	msg string
}

func (e Error) Error() string {
	return e.msg
}

func NewHook(projectID int64, apiKey, env string, stackTraceLevel int, synchronous bool) *airbrakeHook {
	airbrake := gobrake.NewNotifier(projectID, apiKey)
	airbrake.AddFilter(func(notice *gobrake.Notice) *gobrake.Notice {
		if env == "development" {
			return nil
		}
		notice.Context["environment"] = env
		return notice
	})
	hook := &airbrakeHook{
		airbrake,
		stackTraceLevel,
		synchronous,
	}
	return hook
}

func (hook *airbrakeHook) Fire(entry *logrus.Entry) error {
	var notifyErr error
	err, ok := entry.Data["error"].(error)
	if ok {
		notifyErr = err
	} else {
		notifyErr = errors.New(entry.Message)
	}
	var req *http.Request
	for k, v := range entry.Data {
		if r, ok := v.(*http.Request); ok {
			req = r
			delete(entry.Data, k)
			break
		}
	}
	notice := hook.Airbrake.Notice(notifyErr, req, hook.StackTraceLevel)
	for k, v := range entry.Data {
		notice.Env[k] = fmt.Sprintf("%s", v)
	}

	hook.Airbrake.SendNoticeAsync(notice)

	if hook.Synchronous {
		hook.Airbrake.Flush()
	}

	return nil
}

func (hook *airbrakeHook) Levels() []logrus.Level {
	return []logrus.Level{
		logrus.ErrorLevel,
		logrus.FatalLevel,
		logrus.PanicLevel,
	}
}

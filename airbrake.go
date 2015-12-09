package airbrake

import (
	"github.com/Invoiced/logrus"
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

	entryString, err := entry.String()

	if err != nil {
		return err
	}

	//so go brake displays type as error
	cerr := Error{entryString}

	notice := hook.Airbrake.Notice(cerr, nil, hook.StackTraceLevel)

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

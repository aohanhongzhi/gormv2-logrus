package gormv2logrus

import (
	"time"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm/logger"
)

// withLogrus obviously
// option to truncate loooong requests
// option do hide requests containing passwords etc...
// option log request latency
// option to change write format

// options represents all available starting options for gormv2_logrus.
// options are optional parameters that can be passed to init function of gormv2_logrus
// to configure log policy.
type options struct {
	// if a query contains one of bannedKeywords, it will not be logged, it's useful for preventing passwords and secrets
	// for being logged.
	bannedKeywords []BannedKeyword

	// pointer to your logrusEntry instance
	logrusEntry *logrus.Entry

	lr *logrus.Logger

	SlowThreshold time.Duration

	LogLevel logger.LogLevel

	// if set tO 0, nothing wil be truncated, else you can set it to the value you want to avoid
	// logging too big SQL queries
	truncateLen uint

	// if set to true, it will add latency informations for your queries
	logLatency bool
}

// BannedKeyword represents a rule for scanning for Keyword in log output.
type BannedKeyword struct {
	// Keyword represent the string watched, for example : "password"
	Keyword string
	// CaseMatters if set to false, the Keyword matching will occur depending the case.
	// if set to true, Keyword will stricly match input messages
	CaseMatters bool
}

type GormOptions struct {
	SlowThreshold time.Duration
	LogLevel      logger.LogLevel
	TruncateLen   uint
	LogLatency    bool
}

func defaultOptions() options {
	return options{
		logrusEntry:    nil,
		truncateLen:    0,
		bannedKeywords: nil,
		logLatency:     false,
	}
}

// Option interface is used to configure gormv2_logrus options.
type Option interface {
	apply(*options)
}

// funcOption wraps a function that modifies options into an
// implementation of the Option interface.
type funcOption struct {
	f func(*options)
}

// apply is used in mstack init function to read given parameters.
func (fo *funcOption) apply(do *options) {
	fo.f(do)
}

// newGormLogOption is implemented by function that save parameters.
func newGormLogOption(f func(*options)) *funcOption {
	return &funcOption{
		f: f,
	}
}

// WithLogrusEntry Option (not compatible with WithLogrus) used to specify your logrusEntry instance.
// If you don't set a logrusEntry isntance or if your logrusInstance is nil, Gormlog will consider
// that you want log to be printed on stdout.
// It's useful on developpement purposes when you want to see your logs directly in terminal.
func WithLogrusEntry(logrusEntry *logrus.Entry) Option {
	return newGormLogOption(func(o *options) {
		o.logrusEntry = logrusEntry
	})
}

// WithLogrus Option is (not compatible with WithLogrusEntry) is used to set your logrus isntance.
func WithLogrus(lr *logrus.Logger) Option {
	return newGormLogOption(func(o *options) {
		o.lr = lr
	})
}

func WithBannedKeyword(bannedKeywords []BannedKeyword) Option {
	return newGormLogOption(func(o *options) {
		o.bannedKeywords = bannedKeywords
	})
}

func WithGormOptions(gormOpt GormOptions) Option {
	return newGormLogOption(func(o *options) {
		o.logLatency = gormOpt.LogLatency
		o.LogLevel = gormOpt.LogLevel
		o.SlowThreshold = gormOpt.SlowThreshold
	})
}

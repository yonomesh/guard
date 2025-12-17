package guard

import (
	"encoding/json"
	"fmt"
	"io"
	"time"
)

type Logging struct {
	// 非结构化的日志，如何记录
	// Sink is the destination for all unstructured logs emitted
	// from Go's standard library logger. These logs are common
	// in dependencies that are not designed specifically for use
	// in Caddy. Because it is global and unstructured, the sink
	// lacks most advanced features and customizations.
	Sink *SinkLog `json:"sink,omitempty"`

	// 对日志的自定义
	// Logs are your logs, keyed by an arbitrary name of your
	// choosing. The default log can be customized by defining
	// a log called "default". You can further define other logs
	// and filter what kinds of entries they accept.
	Logs map[string]*CustomLog `json:"logs,omitempty"`

	// writerKeys 则保证了当日志配置不再使用时，可以正确地关闭已打开的日志流，从而避免资源泄露。
	// a list of all keys for open writers; all writers
	// that are opened to provision this logging config
	// must have their keys added to this list so they
	// can be closed when cleaning up
	writerKeys []string
}

// SinkLog configures the default Go standard library
// global logger in the log package. This is necessary because
// module dependencies which are not built specifically for
// Caddy will use the standard logger. This is also known as
// the "sink" logger.
type SinkLog struct {
	BaseLog
}

// CustomLog represents a custom logger configuration.
//
// By default, a log will emit all log entries. Some entries
// will be skipped if sampling is enabled. Further, the Include
// and Exclude parameters define which loggers (by name) are
// allowed or rejected from emitting in this log. If both Include
// and Exclude are populated, their values must be mutually
// exclusive, and longer namespaces have priority. If neither
// are populated, all logs are emitted.
type CustomLog struct {
	BaseLog

	// Include defines the names of loggers to emit in this
	// log. For example, to include only logs emitted by the
	// admin API, you would include "admin.api".
	Include []string `json:"include,omitempty"`

	// Exclude defines the names of loggers that should be
	// skipped by this log. For example, to exclude only
	// HTTP access logs, you would exclude "http.log.access".
	Exclude []string `json:"exclude,omitempty"`
}

// BaseLog contains the common logging parameters for logging.
type BaseLog struct {
	// The module that writes out log entries for the sink.
	WriterRaw json.RawMessage `json:"writer,omitempty" caddy:"namespace=caddy.logging.writers inline_key=output"`

	// The encoder is how the log entries are formatted or encoded.
	EncoderRaw json.RawMessage `json:"encoder,omitempty" caddy:"namespace=caddy.logging.encoders inline_key=format"`

	// Tees entries through a zap.Core module which can extract
	// log entry metadata and fields for further processing.
	CoreRaw json.RawMessage `json:"core,omitempty" caddy:"namespace=caddy.logging.cores inline_key=module"`

	// 大于此等级的Log才被记录
	// Level is the minimum level to emit, and is inclusive.
	// Possible levels: DEBUG, INFO, WARN, ERROR, PANIC, and FATAL
	Level string `json:"level,omitempty"`

	// 日志采样器，记录部分组件
	// Sampling configures log entry sampling. If enabled,
	// only some log entries will be emitted. This is useful
	// for improving performance on extremely high-pressure
	// servers.
	Sampling *LogSampling `json:"sampling,omitempty"`

	// 包含调用者文件名和行数
	// If true, the log entry will include the caller's
	// file name and line number. Default off.
	WithCaller bool `json:"with_caller,omitempty"`

	// 且此值非零，则会跳过指定的堆栈帧来确定调用者。默认值为 0。
	// If non-zero, and `with_caller` is true, this many
	// stack frames will be skipped when determining the
	// caller. Default 0.
	WithCallerSkip int `json:"with_caller_skip,omitempty"`

	// 如果非空，并且日志级别高于或等于指定值时，将包括堆栈跟踪。默认关闭。
	// If not empty, the log entry will include a stack trace
	// for all logs at the given level or higher. See `level`
	// for possible values. Default off.
	WithStacktrace string `json:"with_stacktrace,omitempty"`

	writerOpener WriterOpener
	writer       io.WriteCloser
	encoder      string //zapcore.Encoder
	levelEnabler string //zapcore.LevelEnabler
	core         string // zapcore.Core
}

// WriterOpener is a module that can open a log writer.
// It can return a human-readable string representation
// of itself so that operators can understand where
// the logs are going.
type WriterOpener interface {
	fmt.Stringer

	// WriterKey is a string that uniquely identifies this
	// writer configuration. It is not shown to humans.
	WriterKey() string

	// OpenWriter opens a log for writing. The writer
	// should be safe for concurrent use but need not
	// be synchronous.
	OpenWriter() (io.WriteCloser, error)
}

// LogSampling configures log entry sampling.
type LogSampling struct {
	// The window over which to conduct sampling.
	Interval time.Duration `json:"interval,omitempty"`

	// Log this many entries within a given level and
	// message for each interval.
	First int `json:"first,omitempty"`

	// If more entries with the same level and message
	// are seen during the same interval, keep one in
	// this many entries until the end of the interval.
	Thereafter int `json:"thereafter,omitempty"`
}

// Log represents the log data format.
type LogEntry struct {
	Time     string   `json:"ts"`       // Timestamp of the log entry
	Level    string   `json:"level"`    // Log level (e.g., Trace, Debug, Info, Warning, Error, Fataland Panic)
	Category string   `json:"category"` // Category or type of the log (e.g., user-action)
	Tags     []string `json:"tags"`     // Tags related to the log
	Msg      Msger    `json:"msg"`      // Msg content, implemented via the interface for customization
	Extra    Extra    `json:"extra"`    // Ectra content, implemented via the interface for customization
}

type Msger interface {
	MsgToString() (string, error)
}

type Extra interface {
	ExtraToString() (string, error)
}

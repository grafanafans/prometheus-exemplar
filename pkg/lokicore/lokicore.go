package lokicore

import (
	"fmt"
	"strings"
	"time"

	"github.com/afiskon/promtail-client/promtail"
	"go.uber.org/zap/zapcore"
)

var emptyZapEntCaller = zapcore.EntryCaller{}

var promtailLevel = map[zapcore.Level]promtail.LogLevel{
	zapcore.DebugLevel:  promtail.DEBUG,
	zapcore.InfoLevel:   promtail.INFO,
	zapcore.WarnLevel:   promtail.WARN,
	zapcore.ErrorLevel:  promtail.ERROR,
	zapcore.DPanicLevel: promtail.ERROR,
	zapcore.PanicLevel:  promtail.ERROR,
	zapcore.FatalLevel:  promtail.ERROR,
}

type LokiClientConfig struct {
	URL                string
	LevelName          string
	SendLevel          zapcore.Level
	Labels             map[string]string
	BatchWait          time.Duration
	BatchEntriesNumber int
}

func (c *LokiClientConfig) setDefault() {
	if c.URL == "" {
		c.URL = "http://localhost:3100/api/prom/push"
	}
	if c.LevelName == "" {
		c.LevelName = "severity"
	}
	if len(c.Labels) == 0 {
		c.Labels = map[string]string{
			"source": "test",
			"job":    "job",
		}
	}

	if c.BatchWait == 0 {
		c.BatchWait = 5 * time.Second
	}

	if c.BatchEntriesNumber == 0 {
		c.BatchEntriesNumber = 10000
	}
}

func (c *LokiClientConfig) genLabelsWithLogLevel(level string) string {
	c.Labels[c.LevelName] = level
	labelsList := []string{}
	for k, v := range c.Labels {
		labelsList = append(labelsList, fmt.Sprintf(`%s="%s"`, k, v))
	}
	labelString := fmt.Sprintf(`{%s}`, strings.Join(labelsList, ", "))
	return labelString
}

// LokiCore the zapcore of loki
type LokiCore struct {
	cfg                  *LokiClientConfig
	clients              map[zapcore.Level]promtail.Client
	zapcore.LevelEnabler                        // LevelEnabler interface
	fields               map[string]interface{} // save Fields
}

func NewLokiCore(c *LokiClientConfig) (*LokiCore, error) {
	var err error
	if c == nil {
		c = &LokiClientConfig{}
	}
	c.setDefault()
	conf := promtail.ClientConfig{
		PushURL:            c.URL,
		BatchWait:          c.BatchWait,
		BatchEntriesNumber: c.BatchEntriesNumber,
		SendLevel:          promtailLevel[c.SendLevel],
		PrintLevel:         promtail.DISABLE,
	}

	clients := make(map[zapcore.Level]promtail.Client)
	for k := range promtailLevel {
		conf.Labels = c.genLabelsWithLogLevel(k.String())
		clients[k], err = promtail.NewClientJson(conf)
		if err != nil {
			return nil, fmt.Errorf("unable to init promtail client: %v", err)
		}
	}
	return &LokiCore{
		cfg:          c,
		clients:      clients,
		fields:       make(map[string]interface{}),
		LevelEnabler: c.SendLevel,
	}, nil
}

func (c *LokiCore) with(fs []zapcore.Field) *LokiCore {
	m := make(map[string]interface{}, len(c.fields))
	for k, v := range c.fields {
		m[k] = v
	}

	enc := zapcore.NewMapObjectEncoder()
	for _, f := range fs {
		f.AddTo(enc)
	}

	for k, v := range enc.Fields {
		m[k] = v
	}

	return &LokiCore{
		cfg:          c.cfg,
		clients:      c.clients,
		fields:       m,
		LevelEnabler: c.LevelEnabler,
	}
}

func (c *LokiCore) With(fs []zapcore.Field) zapcore.Core {
	return c.with(fs)
}

func (c *LokiCore) Check(ent zapcore.Entry, ce *zapcore.CheckedEntry) *zapcore.CheckedEntry {
	if c.cfg.SendLevel.Enabled(ent.Level) {
		return ce.AddCore(ent, c)
	}
	return ce
}

func (c *LokiCore) Write(ent zapcore.Entry, fs []zapcore.Field) error {
	clone := c.with(fs)

	message := fmt.Sprintf(`msg=%q`, ent.Message)
	for k, v := range clone.fields {
		message += fmt.Sprintf(` %s=%v`, k, v)
	}

	if ent.Caller != emptyZapEntCaller {
		message += fmt.Sprintf(` caller="%s:%d"`, ent.Caller.File, ent.Caller.Line)
	}
	if ent.Stack != "" {
		message += fmt.Sprintf(` stacktrace=%q`, ent.Stack)
	}

	lvl := promtailLevel[ent.Level]
	switch lvl {
	case promtail.DEBUG:
		c.clients[ent.Level].Debugf(message)
	case promtail.INFO:
		c.clients[ent.Level].Infof(message)
	case promtail.WARN:
		c.clients[ent.Level].Warnf(message)
	case promtail.ERROR:
		c.clients[ent.Level].Errorf(message)
	default:
		return fmt.Errorf("unknown log level")
	}
	return nil
}

func (c *LokiCore) Sync() error {
	return nil
}

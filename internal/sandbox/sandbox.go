// internal/sandbox/sandbox.go
package sandbox

// Options controls the sandbox server's runtime behaviour.
type Options struct {
	Host     string // listen address, default "127.0.0.1"
	Port     int    // 0 = OS-assigned random port
	LogLevel string // "info" | "warn" | "error" | "silent"
	LogFile  string // path to append JSON log; "" = no file
	Format   string // "auto" | "schema" | "faker"
}

func (o Options) hostOrDefault() string {
	if o.Host == "" {
		return "127.0.0.1"
	}
	return o.Host
}

func (o Options) logLevelOrDefault() string {
	if o.LogLevel == "" {
		return "info"
	}
	return o.LogLevel
}

func (o Options) formatOrDefault() string {
	if o.Format == "" {
		return "auto"
	}
	return o.Format
}

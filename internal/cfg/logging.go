package cfg

type Logging struct {
	// ONE of trace, debug, info, warn, error, fatal, panic
	Level         string `json:"level"`
	// File to write logs to
	LogFile       string `json:"log_file"`
	// Whether to prettify the output. Set to false to log JSON.
	PrettyConsole bool   `json:"pretty_console"`
}

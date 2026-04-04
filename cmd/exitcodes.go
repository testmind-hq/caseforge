// cmd/exitcodes.go
package cmd

// Exit codes returned by caseforge commands.
// Each exit code has a distinct semantic meaning so that callers (CI systems,
// AI agents, shell scripts) can handle different failure modes precisely.
const (
	ExitOK             = 0 // command completed successfully
	ExitGeneralError   = 1 // unclassified error (also used by cobra for flag errors)
	ExitSpecParseError = 2 // OpenAPI spec could not be parsed
	ExitLintError      = 3 // lint found issues and --fail-on threshold was reached
	ExitNoOutput       = 4 // LLM required but unavailable (and --no-ai not set)
	ExitWriteError     = 5 // output could not be written or rendered
	ExitPartialSuccess = 6 // some operations succeeded but others failed
)

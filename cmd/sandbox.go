// cmd/sandbox.go
package cmd

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/spf13/cobra"
	"github.com/testmind-hq/caseforge/internal/sandbox"
	"github.com/testmind-hq/caseforge/internal/spec"
)

var (
	sandboxSpec     string
	sandboxPort     int
	sandboxHost     string
	sandboxLogLevel string
	sandboxLogFile  string
	sandboxFormat   string
)

var sandboxCmd = &cobra.Command{
	Use:   "sandbox",
	Short: "Start a local mock API server from an OpenAPI spec",
	Long: `Sandbox starts a local HTTP server that responds to requests described
by an OpenAPI spec. It generates realistic mock responses using a
strategy chain (example → schema → faker) and tracks stateful CRUD flows.

Prints the listening address to stdout on startup. Stop with Ctrl-C.

Examples:
  caseforge sandbox --spec openapi.yaml
  caseforge sandbox --spec openapi.yaml --port 8080 --log-level info --log-file sandbox.log`,
	RunE: runSandbox,
}

func init() {
	rootCmd.AddCommand(sandboxCmd)
	sandboxCmd.Flags().StringVar(&sandboxSpec, "spec", "", "OpenAPI spec file or URL (required)")
	_ = sandboxCmd.MarkFlagRequired("spec")
	sandboxCmd.Flags().IntVar(&sandboxPort, "port", 0, "Listen port (0 = random)")
	sandboxCmd.Flags().StringVar(&sandboxHost, "host", "127.0.0.1", "Listen address")
	sandboxCmd.Flags().StringVar(&sandboxLogLevel, "log-level", "info", "Log level: info|warn|error|silent")
	sandboxCmd.Flags().StringVar(&sandboxLogFile, "log-file", "", "Append JSON logs to file (optional)")
	sandboxCmd.Flags().StringVar(&sandboxFormat, "format", "auto", "Response strategy: auto|schema|faker")
}

func runSandbox(cmd *cobra.Command, _ []string) error {
	ps, err := spec.NewLoader().Load(sandboxSpec)
	if err != nil {
		return fmt.Errorf("loading spec: %w", err)
	}

	opts := sandbox.Options{
		Host:     sandboxHost,
		Port:     sandboxPort,
		LogLevel: sandboxLogLevel,
		LogFile:  sandboxLogFile,
		Format:   sandboxFormat,
	}

	srv := sandbox.NewSandboxServer(ps, opts)
	if err := srv.Start(opts.Host, opts.Port); err != nil {
		return err
	}

	fmt.Fprintf(cmd.OutOrStdout(), "caseforge sandbox listening on http://%s\n", srv.Addr())

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	<-ctx.Done()

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return srv.Shutdown(shutdownCtx)
}

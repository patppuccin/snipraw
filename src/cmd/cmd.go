package cmd

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/patppuccin/snipraw/src/config"
	"github.com/patppuccin/snipraw/src/console"
	"github.com/patppuccin/snipraw/src/consts"
	"github.com/patppuccin/snipraw/src/helpers"
	"github.com/patppuccin/snipraw/src/server"
	"github.com/spf13/cobra"
)

// setup flags
var flagConfig string
var runtime = &server.Runtime{}

// command entry point
var SRCmd = &cobra.Command{
	Use:   consts.AppName,
	Short: consts.AppDesc,
	Long:  console.Banner(consts.AppDesc),
	Version: fmt.Sprintf("%s (commit: %s, built: %s)",
		consts.AppVersion,
		consts.BuildCommit,
		helpers.FormatDateString(consts.BuildDate),
	),
	SilenceUsage:  true,
	SilenceErrors: true,
	PreRun: func(cmd *cobra.Command, args []string) {
		// ensure no arguments are passed
		if len(args) > 0 {
			console.Error("No arguments are allowed for " + consts.AppName)
			os.Exit(1)
		}

		// ensure the snippets directory is specified explicitly
		if runtime.Dir == "" {
			console.Error("Requires a directory to serve snippets from. Use --dir to specify.")
			os.Exit(1)
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
		// resolve snippets directory and ensure it exists
		dir, err := helpers.ResolvePath(runtime.Dir)
		if err != nil {
			console.Error("Failed to resolve snippets directory: " + err.Error())
			os.Exit(1)
		}

		if _, err := os.Stat(dir); os.IsNotExist(err) {
			console.Error("Snippets directory does not exist: " + dir)
			os.Exit(1)
		}

		runtime.Dir = dir

		// resolve configuration and load it
		cfg := config.Default(runtime.Host, runtime.Port)

		if flagConfig != "" {
			cfgPath, err := helpers.ResolvePath(flagConfig)
			if err != nil {
				console.Error("Failed to resolve config path: " + err.Error())
				os.Exit(1)
			}

			if _, err := os.Stat(cfgPath); os.IsNotExist(err) {
				if err := config.Write(cfgPath, cfg); err != nil {
					console.Error("Failed to create config: " + err.Error())
					os.Exit(1)
				}
			} else {
				loaded, err := config.Load(cfgPath)
				if err != nil {
					console.Error("Failed to load config: " + err.Error())
					os.Exit(1)
				}
				cfg = loaded
			}
		}

		// setup context and run server
		ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
		defer cancel()

		if err := server.Run(ctx, runtime, cfg); err != nil {
			console.Error("server encountered an error: " + err.Error())
			os.Exit(1)
		}
	},
}

func init() {
	// tweak cobra behaviour
	cobra.EnableCommandSorting = false
	SRCmd.CompletionOptions.DisableDefaultCmd = true
	SRCmd.SetHelpCommand(&cobra.Command{Hidden: true})
	SRCmd.Flags().SortFlags = false

	// setup and parse flags
	SRCmd.Flags().StringVar(&runtime.Host, "host", "127.0.0.1", "host to bind to")
	SRCmd.Flags().IntVar(&runtime.Port, "port", 8245, "port to bind to")
	SRCmd.Flags().StringVar(&runtime.Dir, "dir", "", "directory to serve snippets from")
	SRCmd.Flags().StringVar(&runtime.LogLevel, "log-level", "info", "log level")
	SRCmd.Flags().StringVar(&flagConfig, "config", "", "path to config file")

	// append notes to help template
	SRCmd.SetHelpTemplate(SRCmd.HelpTemplate() +
		"\nSet env var NO_COLOR to disable colored output" +
		"\nDocumentation at https://snipraw.patppuccin.com\n",
	)

	// setup custom error handler
	SRCmd.SetFlagErrorFunc(func(cmd *cobra.Command, err error) error {
		cmd.Help()
		os.Stdout.WriteString("\n")
		console.Error(err.Error())
		os.Stdout.WriteString("\n")
		os.Exit(1)
		return nil
	})
}

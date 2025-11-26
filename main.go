package main

//go:generate go run ./internal/generator/generate.go code providers.json internal/libdnsregistry/registry.go
import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"runtime/debug"
	"slices"
	"strings"
	"time"

	"github.com/project0/external-dns-libdns-webhook/internal/externaldns"
	"github.com/project0/external-dns-libdns-webhook/internal/libdnsregistry"
	"github.com/urfave/cli/v3"
	webhookApi "sigs.k8s.io/external-dns/provider/webhook/api"
)

const (
	envPrefix               = "LIBDNS_"
	flagLogLevel            = "log.level"
	flagLogFormat           = "log.format"
	flagProviderName        = "provider.name"
	flagProviderConfig      = "provider.config"
	flagProviderZones       = "provider.zones"
	flagWebhookListen       = "webhook.listen"
	flagWebhookReadTimeout  = "webhook.read-timeout"
	flagWebhookWriteTimeout = "webhook.write-timeout"
)

func flagEnv(name string) cli.ValueSourceChain {
	return cli.EnvVars(
		strings.ToUpper(envPrefix +
			strings.ReplaceAll(
				strings.ReplaceAll(name, "-", "_"),
				".", "_",
			),
		))
}

func flagAlternatives(name string) cli.ValueSourceChain {
	// IDEA: Add more sources here
	return cli.NewValueSourceChain(
		flagEnv(name).Chain...,
	)
}

const description = `A generic external-dns webhook provider supporting several libdns based DNS providers.

Supported Providers:
%s
`

func main() {
	var version string
	if info, ok := debug.ReadBuildInfo(); ok {
		version = info.Main.Version
	}
	startedChan := make(chan struct{}, 1)
	cmd := &cli.Command{
		Name:        "external-dns-webhook-libdns",
		Version:     version,
		Usage:       "Webhook for external-dns using libdns providers.",
		Description: fmt.Sprintf(description, strings.Join(libdnsregistry.List(), ", ")),
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    flagLogLevel,
				Usage:   "The log level (debug, info, warn, error)",
				Value:   "info",
				Sources: flagAlternatives(flagLogLevel),
				Validator: func(s string) error {
					if _, err := parseLogLevel(s); err != nil {
						return fmt.Errorf("cannot set log level: %w", err)
					}
					return nil
				},
			},
			&cli.StringFlag{
				Name:    flagLogFormat,
				Usage:   "The log format (text, json)",
				Value:   "text",
				Sources: flagAlternatives(flagLogFormat),
				Validator: func(s string) error {
					if s != "text" && s != "json" {
						return errors.New("log format must be text or json")
					}

					return nil
				},
			},
			&cli.StringFlag{
				Name:     flagProviderName,
				Usage:    "The name of the libdns provider",
				Sources:  flagAlternatives(flagProviderName),
				Required: true,
				Validator: func(s string) error {
					if !slices.Contains(libdnsregistry.List(), s) {
						return fmt.Errorf("libdns provider %s not found", s)
					}

					return nil
				},
			},
			&cli.StringFlag{
				Name:    flagProviderConfig,
				Usage:   "The json config for the libdns provider as a string",
				Sources: flagAlternatives(flagProviderConfig),
				Value:   "{}",
			},
			&cli.StringSliceFlag{
				Name:     flagProviderZones,
				Usage:    "The name of the dns zones that the provider will manage",
				Required: true,
				Sources:  flagAlternatives(flagProviderZones),
			},
			&cli.StringFlag{
				Name:    flagWebhookListen,
				Usage:   "The webhook server address to listen on (host:port)",
				Value:   "localhost:8888",
				Sources: flagAlternatives(flagWebhookListen),
			},
			&cli.DurationFlag{
				Name:    flagWebhookReadTimeout,
				Usage:   "The webhook server read timeout",
				Value:   time.Minute,
				Sources: flagAlternatives(flagWebhookReadTimeout),
			},
			&cli.DurationFlag{
				Name:    flagWebhookWriteTimeout,
				Usage:   "The webhook server write timeout",
				Value:   time.Minute,
				Sources: flagAlternatives(flagWebhookWriteTimeout),
			},
		},

		Before: func(ctx context.Context, cmd *cli.Command) (context.Context, error) {
			// setup logging
			if cmd.String(flagLogFormat) == "json" {
				slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, nil)))
			}

			level, err := parseLogLevel(cmd.String(flagLogLevel))
			if err != nil {
				return ctx, fmt.Errorf("failed to parse log level: %w", err)
			}

			slog.Info("Setting log level", "level", level.String())
			slog.SetLogLoggerLevel(level)

			return ctx, nil
		},

		Action: func(_ context.Context, cmd *cli.Command) error {
			confs := [][]byte{}
			if conf := cmd.String(flagProviderConfig); conf != "" {
				confs = append(confs, []byte(conf))
			}

			if len(confs) == 0 {
				return errors.New("no provider config is provided")
			}

			providerName := cmd.String(flagProviderName)

			libdnsProvider, err := libdnsregistry.New(providerName, confs)
			if err != nil {
				return fmt.Errorf("failed to create provider %s: %w", providerName, err)
			}

			webhookApi.StartHTTPApi(
				*externaldns.NewWebhookProvider(
					cmd.StringSlice(flagProviderZones),
					libdnsProvider,
				),
				startedChan,
				cmd.Duration(flagWebhookReadTimeout),
				cmd.Duration(flagWebhookWriteTimeout),
				cmd.String(flagWebhookListen),
			)
			<-startedChan
			slog.Info("Startup completed")
			return nil
		},
	}

	if err := cmd.Run(context.Background(), os.Args); err != nil {
		slog.Error("Failed to run", "err", err)
	}
}

func parseLogLevel(s string) (slog.Level, error) {
	var level slog.Level
	var err = level.UnmarshalText([]byte(s))
	return level, err
}

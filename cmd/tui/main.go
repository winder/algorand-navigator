package main

import (
	"context"
	"fmt"
	"github.com/urfave/cli/v3"
	"os"

	"github.com/algorand/node-ui/tui"
	"github.com/algorand/node-ui/tui/args"
	"github.com/algorand/node-ui/version"
)

func main() {
	err := makeCommand().Run(context.Background(), os.Args)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Problem running command: %s\n", err.Error())
	}
}

func run(args args.Arguments) {
	if args.VersionFlag {
		fmt.Println(version.LongVersion())
		os.Exit(0)
	}
	tui.Start(args)
}

func makeCommand() *cli.Command {
	var args args.Arguments
	return &cli.Command{
		Name:  "node-ui",
		Usage: "Launch the Algorand Node UI.",
		Flags: []cli.Flag{
			&cli.Uint64Flag{
				Name:        "tui-port",
				Aliases:     []string{"p"},
				Usage:       "Port address to host TUI from, set to 0 to run directly.",
				Value:       0,
				Sources:     cli.EnvVars("TUI_PORT"),
				Destination: &args.TuiPort,
			},
			&cli.StringFlag{
				Name:        "algod-url",
				Aliases:     []string{"u"},
				Usage:       "Algod URL and port to monitor, formatted like localhost:1234.",
				Value:       "",
				Sources:     cli.EnvVars("ALGOD_URL"),
				Destination: &args.AlgodURL,
			},
			&cli.StringFlag{
				Name:        "algod-token",
				Aliases:     []string{"t"},
				Usage:       "Algod REST API token.",
				Value:       "",
				Sources:     cli.EnvVars("ALGOD_TOKEN"),
				Destination: &args.AlgodToken,
			},
			&cli.StringFlag{
				Name:        "algod-admin-token",
				Aliases:     []string{"a"},
				Usage:       "Algod REST API Admin token.",
				Value:       "",
				Sources:     cli.EnvVars("ALGOD_ADMIN_TOKEN"),
				Destination: &args.AlgodAdminToken,
			},
			&cli.StringFlag{
				Name:        "algod-data-dir",
				Aliases:     []string{"d"},
				Usage:       "Path to Algorand data directory.",
				Value:       "",
				Sources:     cli.EnvVars("ALGORAND_DATA"),
				Destination: &args.AlgodDataDir,
			},
			&cli.StringSliceFlag{
				Name:        "watch-list",
				Aliases:     []string{"w"},
				Usage:       "Account addresses to watch in the accounts tab, may provide more than once to watch multiple accounts. Use comma separated values if providing more than one account with an environment variable.",
				Value:       nil,
				Sources:     cli.EnvVars("WATCH_LIST"),
				Destination: &args.AddressWatchList,
			},
			&cli.BoolFlag{
				Name:        "version",
				Aliases:     []string{"v"},
				Usage:       "Print version information and exit.",
				Value:       false,
				Destination: &args.VersionFlag,
			},
		},
		Action: func(c *cli.Context) error {
			run(args)
			return nil
		},
	}
}

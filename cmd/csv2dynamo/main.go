package main

import (
    "github.com/maito1201/csv2dynamo"
	"github.com/urfave/cli/v2"
	"log"
	"os"
)

func main() {
	app := &cli.App{
		Name:  "csv2dynamo",
		Usage: "reimport csv from export function of aws DynamoDB",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "file",
				Aliases:  []string{"f", "csv"},
				Usage:    "file to import e.g ./tablename.csv (required)",
				Required: true,
			},
			&cli.StringFlag{
				Name:     "table-name",
				Aliases:  []string{"t"},
				Usage:    "target dynamo db tabe name (required)",
				Required: true,
			},
			&cli.StringFlag{
				Name:    "endpoint",
				Aliases: []string{"e"},
				Usage:   "endpoint of DynamoDB",
			},
			&cli.StringFlag{
				Name:    "profile",
				Aliases: []string{"p"},
				Usage:   "profile of aws cli",
			},
			&cli.StringFlag{
				Name:    "output",
				Aliases: []string{"o", "out"},
				Usage:   "target output (default: stdout), no file will be created if execute option is enabled",
			},
			&cli.BoolFlag{
				Name:  "execute",
				Usage: "is directly execute import command",
			},
		},
		Action: csv2dynamo.Execute,
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

package main

import (
	"fmt"
	"github.com/maito1201/clearout"
	"github.com/maito1201/csv2dynamo/csvreader"
	"github.com/urfave/cli/v2"
	"log"
	"os"
	"os/exec"
	"strings"
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
		Action: execute,
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func execute(c *cli.Context) error {
	cmd := []string{"aws", "dynamodb", "put-item"}
	if c.String("profile") != "" {
		cmd = append(cmd, []string{"--profile", c.String("profile")}...)
	}
	if c.String("endpoint") != "" {
		cmd = append(cmd, []string{"--endpoint", c.String("endpoint")}...)
	}
	in, err := csvreader.ReadCSV(c.String("file"))
	if err != nil {
		return err
	}

	var outs [][]string
	for _, v := range in {
		o := append(cmd, []string{"--table-name", c.String("table-name"), "--item", v.ToJsonString(c.Bool("execute"))}...)
		outs = append(outs, o)
	}

	if c.Bool("execute") {
		for _, v := range outs {
			cmd := exec.Command(v[0], v[1:]...)
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			err := cmd.Run()
			if err != nil {
				log.Fatal(err)
			}
		}
		return nil
	}

	if c.String("output") != "" {
		writeFile(outs, c.String("output"))
	} else {
		for _, v := range outs {
			fmt.Println(strings.Join(v, " "))
		}
	}
	return nil
}

func writeFile(outs [][]string, path string) error {
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return err
	}
	defer f.Close()
	cout := clearout.Output{Prefix: fmt.Sprintf("writing commands to %v\n", path)}
	for i := 0; i < len(outs); i++ {
		cout.Printf("progress: %d/%d\n", i+1, len(outs))
		if i == len(outs)-1 {
			cout.Println("complete!")
		}
		cout.Render()
		fmt.Fprintln(f, strings.Join(outs[i], " "))
	}
	return nil
}

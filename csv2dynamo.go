package csv2dynamo

import (
	"fmt"
	"github.com/maito1201/clearout"
	"github.com/urfave/cli/v2"
	"log"
	"os"
	"os/exec"
	"strings"
)

func Execute(c *cli.Context) error {
	cmd := []string{"aws", "dynamodb", "put-item"}
	if c.String("profile") != "" {
		cmd = append(cmd, []string{"--profile", c.String("profile")}...)
	}
	if c.String("endpoint") != "" {
		cmd = append(cmd, []string{"--endpoint", c.String("endpoint")}...)
	}
	in, err := readCSV(c.String("file"))
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

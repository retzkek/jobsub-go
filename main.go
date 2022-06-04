package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"

	"github.com/urfave/cli/v2"
)

func condor_wrapper(command string) func(ctx *cli.Context) error {
	return func(ctx *cli.Context) error {
		if ctx.NArg() < 1 {
			return fmt.Errorf("must specify at least one job")
		}
		jobs := make([]Job, 0)
		args := make([]string, ctx.Args().Len())
		for _, jid := range ctx.Args().Slice() {
			j := Job{ID: jid}
			if err := j.DecomposeID(); err == nil {
				jobs = append(jobs, j)
			} else {
				// not a job id, must be an argument to pass through?
				args = append(args, jid)
			}
		}

		// get creds
		group := ctx.String("group")
		if group == "" {
			var err error
			group, err = getExp()
			if err != nil {
				log.Fatalf("error determining experiment: %s", err)
			}
		}
		log.Print(group)

		role, err := getRole()
		if err != nil {
			log.Fatalf("error determining role: %s", err)
		}
		log.Print(role)

		if err := getToken(group, role); err != nil {
			log.Fatalf("%s", err)
		}

		// run the command for each job
		for _, j := range jobs {
			jargs := append(args, "-name", j.Schedd)
			if j.Cluster {
				jargs = append(jargs, j.Seq)
			} else {
				jargs = append(jargs, j.Seq+"."+j.Proc)
			}
			condor_command := "condor_" + command
			log.Printf("running %s with args %v", condor_command, jargs)
			cmd := exec.Command(condor_command, jargs...)
			cmd.Stderr = os.Stderr
			cmd.Stdout = os.Stdout
			if err := cmd.Run(); err != nil {
				log.Fatalf("error running condor command: %s", err)
			}
		}
		return nil
	}
}

func main() {
	app := &cli.App{
		Usage: "FIFE job management client",
		Flags: []cli.Flag{
			&cli.IntFlag{
				Name:    "debug",
				Aliases: []string{"v"},
				Value:   0,
				Usage:   "debug level",
			},
			&cli.StringFlag{
				Name:    "group",
				Aliases: []string{"G"},
				Usage:   "experiment/vo override",
			},
		},
		Commands: []*cli.Command{
			{
				Name:      "fetchlog",
				Aliases:   []string{},
				Usage:     "fetch job logs (stdout etc)",
				ArgsUsage: "JOBID",
			},
			{
				Name:      "hold",
				Aliases:   []string{},
				Usage:     "hold job(s)",
				ArgsUsage: "JOBID",
				Action:    condor_wrapper("hold"),
			},
			{
				Name:      "queue",
				Aliases:   []string{"q"},
				Usage:     "list current job status",
				ArgsUsage: "[JOBID]",
			},
			{
				Name:      "release",
				Aliases:   []string{},
				Usage:     "release held job(s)",
				ArgsUsage: "JOBID",
				Action:    condor_wrapper("release"),
			},
			{
				Name:      "rm",
				Aliases:   []string{},
				Usage:     "remove job(s) from queue",
				ArgsUsage: "JOBID",
				Action:    condor_wrapper("rm"),
			},
			{
				Name:    "submit",
				Aliases: []string{},
				Usage:   "submit job(s) to batch system queue",
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

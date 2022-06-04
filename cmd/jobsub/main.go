package main

import (
	"log"
	"os"

	"github.com/urfave/cli/v2"

	"github.com/hepcloud/jobsub-go"
)

func main() {
	app := &cli.App{
		Usage: "FIFE job management client",
		Flags: jobsub.GlobalFlags,
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
				Action:    jobsub.CondorWrapper("hold"),
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
				Action:    jobsub.CondorWrapper("release"),
			},
			{
				Name:      "rm",
				Aliases:   []string{},
				Usage:     "remove job(s) from queue",
				ArgsUsage: "JOBID",
				Action:    jobsub.CondorWrapper("rm"),
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
package jobsub

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/clok/kemba"
	"github.com/urfave/cli/v2"
)

var FetchlogFlags = []cli.Flag{
	&cli.StringFlag{
		Name:     "jobid",
		Aliases:  []string{"J", "job"},
		Usage:    "job/submission ID",
		Required: true,
	},
	&cli.StringFlag{
		Name:    "destdir",
		Aliases: []string{"dest-dir", "unzipdir"},
		Usage:   "Directory to automatically unarchive logs into",
	},
	&cli.StringFlag{
		Name:  "archive-format",
		Usage: "format for downloaded archive:\"tar\" (default,compressed) or \"zip\"",
		Value: "tar",
	},
}

func Fetchlog(ctx *cli.Context) error {
	k := kemba.New("jobsub:fetchlog")

	// get creds
	if err := CheckCreds(ctx); err != nil {
		return err
	}

	// decompose job ID so we can build the condor command
	j := Job{ID: ctx.String("jobid")}
	if err := j.DecomposeID(); err != nil {
		return fmt.Errorf("error decomposing job id: %w", err)
	}

	// run condor_transfer_data to get output from schedd
	condor_command := "condor_transfer_data"
	args := []string{"-name", j.Schedd}
	if j.Cluster {
		args = append(args, j.Seq)
	} else {
		args = append(args, j.Seq+"."+j.Proc)
	}

	k.Printf("running %s with args %v", condor_command, args)
	cmd := exec.Command(condor_command, args...)
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("error running condor command: %w", err)
	}

	// TODO build tarball

	return nil
}

package jobsub

import (
	"fmt"
	"os"

	"github.com/clok/kemba"
	"github.com/retzkek/htcondor-go"
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

	// determine where condor_transfer_data will put output
	iwd, err := j.GetAttribute("SUBMIT_Iwd")
	if err != nil {
		return fmt.Errorf("error determining log location: %w", err)
	}
	k.Printf("will look for output in %s", iwd)

	// run condor_transfer_data to get output from schedd
	ccmd := htcondor.NewCommand("condor_transfer_data").WithName(j.Schedd)
	if j.Cluster {
		ccmd = ccmd.WithArg(j.Seq)
	} else {
		ccmd = ccmd.WithArg(j.Seq + "." + j.Proc)
	}

	k.Printf("running %s with args %v", ccmd.Command, ccmd.MakeArgs())
	cmd := ccmd.Cmd()
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("error running condor command: %w", err)
	}

	// TODO build tarball

	return nil
}

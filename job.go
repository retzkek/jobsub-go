package main

import (
	"fmt"
	"regexp"
)

var jobIDRegexp *regexp.Regexp

func init() {
	jobIDRegexp = regexp.MustCompile("(\\d+)(?:\\.(\\d+))?@([\\w\\.]+)")
}

// Job is a single HTCondor batch job or cluster
type Job struct {
	ID      string
	Seq     string
	Proc    string
	Schedd  string
	Cluster bool
}

// DecomposeID breaks out the job's ID into components.
func (j *Job) DecomposeID() error {
	if matches := jobIDRegexp.FindStringSubmatch(j.ID); len(matches) == 4 {
		j.Seq = matches[1]
		j.Proc = matches[2]
		j.Schedd = matches[3]
	} else {
		return fmt.Errorf("error parsing job ID %s", j.ID)
	}
	if j.Proc == "" {
		j.Proc = "0"
		j.Cluster = true
	}
	return nil
}

// ComposeID builds the canonical job ID out of components.
func (j *Job) ComposeID() string {
	if j.Seq == "" {
		// well it looks like the ID was never decomposed (or components set). Maybe ID is right.
		return j.ID
	}
	if j.Cluster {
		return fmt.Sprintf("%s@%s", j.Seq, j.Schedd)
	}
	return fmt.Sprintf("%s.%s@%s", j.Seq, j.Proc, j.Schedd)
}

// String returns a string representation of the job (job ID)
func (j *Job) String() string {
	if j.ID != "" {
		return j.ID
	}
	return j.ComposeID()
}

package main

import (
	"fmt"
	"os"
	"os/exec"
	"os/user"
	"strings"
	"time"

	"github.com/lestrrat-go/jwx/jwt"
)

const (
	DEFAULT_VAULT_HOST = "fermicloud543.fnal.gov"
	DEFAULT_ROLE       = "Analysis"
)

// getExp tries to determine the user's experiment/vo
func getExp() (string, error) {
	// check if a recognized env var is set
	for _, ev := range []string{"GROUP", "EXPERIMENT", "SAM_EXPERIMENT"} {
		if g, found := os.LookupEnv(ev); found {
			return g, nil
		}
	}
	// use the primary group name
	u, err := user.Current()
	if err != nil {
		return "", fmt.Errorf("unable to determine experiment name from group: %w", err)
	}
	g, err := user.LookupGroupId(u.Gid)
	if err != nil {
		return "", fmt.Errorf("unable to determine experiment name from group: %w", err)
	}

	return g.Name, nil
}

// getRole determines the user's role from the environment
func getRole() (string, error) {
	u, err := user.Current()
	if err != nil {
		return "", fmt.Errorf("unable to determine current user: %w", err)
	}
	if strings.HasSuffix(u.Name, "pro") {
		return "Production", nil
	}
	return DEFAULT_ROLE, nil
}

// getToken checks for a valid token, otherwise obtains one and sets
// BEARER_TOKEN_FILE
func getToken(exp, role string) error {
	pid := os.Getpid()
	tmp := os.TempDir()
	role = strings.ToLower(role)
	issuer := exp
	if exp == "samdev" {
		issuer = "fermilab"
	}
	tokenfile := fmt.Sprintf("%s/bt_token_%s_%s_%d", tmp, issuer, role, pid)
	if err := os.Setenv("BEARER_TOKEN_FILE", tokenfile); err != nil {
		return fmt.Errorf("error setting BEARER_TOKEN_FILE: %w", err)
	}
	token, err := jwt.ReadFile(tokenfile)
	if err == nil {
		// there's already a token, is it still good?
		if time.Now().Before(token.Expiration()) {
			return nil
		}
	}
	// get new token
	cmd := exec.Command("htgettoken", "-a", DEFAULT_VAULT_HOST, "-i", issuer, "-r", role)
	cmd.Stdout = os.Stderr
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("error running htgettoken: %w", err)
	}
	return nil
}

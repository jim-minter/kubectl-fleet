package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
)

func main() {
	if err := cmdFleet.ExecuteContext(context.Background()); err != nil {
		var exitErr *exec.ExitError
		if errors.As(err, &exitErr) {
			os.Exit(exitErr.ExitCode())
		}

		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

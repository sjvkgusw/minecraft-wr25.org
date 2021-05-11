package main

import (
	"fmt"
	"os/exec"
)

func checkoutBranch(url string, commitish string, target string) error {
	fmt.Printf("%s %s %s %s\n", "git", "clone", url, target)
	cmd := exec.Command("git", "clone", url, target)
	output, err := cmd.Output()
	if err != nil {
		fmt.Printf("Couldn't run git clone command\n")
		fmt.Printf("output: %s\n", output)
		return err
	}

	cmd = exec.Command("git", "checkout", commitish)
	cmd.Dir = target
	output, err = cmd.Output()
	if err != nil {
		fmt.Printf("Couldnt checkout commitish\n")
		fmt.Printf("output: %s", output)
		return err
	}
	return nil
}

func pullBranch(url string, commitish string, target string) error {
	cmd := exec.Command("git", "pull", url, commitish)
	cmd.Dir = target
	output, err := cmd.Output()
	if err != nil {
		fmt.Printf("Couldnt pull commitish\n")
		fmt.Printf("output: %s", output)
		return err
	}
	return nil
}

func pushBranch(target string) error {
	cmd := exec.Command("git", "push")
	cmd.Dir = target
	output, err := cmd.Output()
	if err != nil {
		fmt.Printf("Couldnt push commit\n")
		fmt.Printf("output: %s", output)
		return err
	}
	return nil
}

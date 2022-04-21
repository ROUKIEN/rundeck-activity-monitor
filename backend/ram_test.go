package main

import "testing"

func TestCmdServe(t *testing.T) {
	args := []string{
		"serve",
	}
	run(args)
}

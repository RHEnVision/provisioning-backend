package main

import (
	"fmt"
	"os"
	"os/exec"
	"runtime/debug"
)

func update() {
	goCmd, _ := exec.LookPath("go")
	if len(os.Args) >= 2 {
		goCmd = os.Args[2]
	}
	fmt.Printf("updating project dependencies one by one using %s\n", goCmd)

	if bi, ok := debug.ReadBuildInfo(); ok {
		for _, dep := range bi.Deps {
			fmt.Printf("updating %s\n", dep.Path)
			cmd := exec.Command(goCmd, "get", "-u", dep.Path)
			output, _ := cmd.CombinedOutput()
			if len(output) > 0 {
				fmt.Println(string(output))
			}
			if err := cmd.Run(); err != nil {
				fmt.Println("error updating", dep.Path)
			}
		}
	} else {
		panic("could not read build info")
	}
}

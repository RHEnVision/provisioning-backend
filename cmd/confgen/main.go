package main

import (
	"bufio"
	"fmt"
	"strings"

	_ "github.com/RHEnVision/provisioning-backend/internal/clients/http/azure"
	_ "github.com/RHEnVision/provisioning-backend/internal/clients/http/ec2"
	_ "github.com/RHEnVision/provisioning-backend/internal/clients/http/gcp"
	"github.com/RHEnVision/provisioning-backend/internal/config"
)

func main() {
	config.Initialize("")

	fmt.Printf("#\n# Copy this file as configs/api.env, the syntax is KEY=value.\n")
	fmt.Printf("#\n# Values from configs/{worker,migrate,typesctl,test} take precedence.\n")
	fmt.Printf("# This file was generated by 'make generate-example-config'.\n")

	help, err := config.HelpText()
	if err != nil {
		panic(err)
	}

	s := bufio.NewScanner(strings.NewReader(help))
	for s.Scan() {
		fmt.Printf("# %s\n", s.Text())
	}

	fmt.Printf("#\n\n")
}

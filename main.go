/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package main

import (
	"github.com/caibo86/logger"
	"github.com/skea3344/try/cmd"
)

func main() {
	logger.Init(
		logger.SetIsRedirectErr(false),
	)
	cmd.Execute()
}

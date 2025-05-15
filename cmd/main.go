/*
Copyright © 2025 Cristian Oliveira me@cristianoliveira.dev
*/
package cmd

import "os"

func Run() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

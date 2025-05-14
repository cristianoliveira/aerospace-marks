/*
Copyright © 2025 Cristian Oliveira me@cristianoliveira.dev
*/
package main

import "os"

func main() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

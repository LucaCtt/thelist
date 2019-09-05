//go:generate mockery -all -output $PWD/mocks

package main

import "github.com/LucaCtt/thelist/cmd"

func main() {
	// Start cobra
	cmd.Execute()
}

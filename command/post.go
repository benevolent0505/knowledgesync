package command

import (
	"strings"
)

type PostCommand struct {
	Meta
}

func (c *PostCommand) Run(args []string) int {
	// Write your code here

	return 0
}

func (c *PostCommand) Synopsis() string {
	return ""
}

func (c *PostCommand) Help() string {
	helpText := `

`
	return strings.TrimSpace(helpText)
}

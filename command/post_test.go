package command

import (
	"testing"

	"github.com/mitchellh/cli"
)

func TestPostCommand_implement(t *testing.T) {
	var _ cli.Command = &PostCommand{}
}

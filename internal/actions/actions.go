package actions

import (
	"fmt"
	"github.com/maaslalani/slides/internal/actions/hooks"
	"strings"
)

// Actions holds the current action-context
type Actions struct {
	Prefix    string
	Buffer    string
	StatusBar string
}

func (a *Actions) IsCapturing() bool {
	return a.Prefix != ""
}

func (a *Actions) Begin(prefix string) {
	a.Reset()
	a.Prefix = prefix
}

func (a *Actions) Reset() {
	a.Prefix = ""
	a.Buffer = ""
	a.StatusBar = ""
}

///

func (a *Actions) CreateCtx(m hooks.Model) *hooks.Ctx {
	var command = strings.TrimSpace(a.Buffer)
	var args []string
	if strings.Contains(command, " ") {
		args = strings.Split(command, " ")[1:]
		command = command[:strings.Index(command, " ")]
	}
	return &hooks.Ctx{
		Prefix:  a.Prefix,
		Command: command,
		Args:    args,
		Model:   m,
	}
}

func (a *Actions) GetStatus() string {
	if a.StatusBar != "" {
		return fmt.Sprintf("> %s", a.StatusBar)
	}
	if !a.IsCapturing() {
		return ""
	}
	// Special case: forward- and backward-search
	if a.Prefix == "/" || a.Prefix == "?" {
		var dir = "->"
		if a.Prefix == "?" {
			dir = "<-"
		}
		return fmt.Sprintf("| Search %s [%s]%s", dir, a.Prefix, a.Buffer)
	}
	return fmt.Sprintf("[ %s(%s) ]", a.Prefix, a.Buffer)
}

func (a *Actions) Execute(m hooks.Model) {
	ctx := a.CreateCtx(m)
	a.Reset()
	for _, hook := range hooks.Hooks {
		if msg, ok := hook(ctx); ok {
			a.StatusBar = msg
			return
		}
	}
	a.StatusBar = "err: command (" + ctx.Command + ") not found"
}

package segments

import (
	"regexp"

	"github.com/fatih/color"
)

var (
	colorDim     = color.New(color.Faint)
	colorBlue    = color.New(color.FgCyan)
	colorGreen   = color.New(color.FgGreen)
	colorYellow  = color.New(color.FgYellow)
	colorRed     = color.New(color.FgRed)
	colorMagenta = color.New(color.FgMagenta)
	colorGray    = color.New(color.FgHiBlack)
	colorCyan    = color.New(color.FgCyan)

	branchRegex = regexp.MustCompile(`^(?:No commits yet on )?([^\.\n]+)`)
)

// OSC8 codes
const (
	OSC8Start = "\033]8;;"
	OSC8Sep   = "\033\\"
	OSC8End   = "\033]8;;\033\\"
)

// Icons
const (
	IconModel     = "󰚩"
	IconFolder    = "󰉋"
	IconGitBranch = "󰊢"
	IconWorktree  = "󰙅"
	IconContext   = "󰧑"
	IconCost      = "󰥏"
	IconDuration  = "󱎫"
	IconLines     = "󰆓"
	IconVimNormal = "󰌌"
	IconVimInsert = "󰏫"
	IconAutorun   = "󰄙"
	IconMax       = "󰓅"
	IconSession   = "󰋩"
	IconFastMode  = "⚡"
	IconEffort    = "󰗡"
	IconThinking  = "󰜗"
	IconPR        = "󰐃"
	IconAgent     = "󰀈"
	IconRateLimit = "󰌉"
	IconCache     = "󰋊"
	IconTokenIn   = "↓"
	IconTokenOut  = "↑"
	IconExceeds   = "󰄘"
	IconBarFilled = "#"
	IconBarEmpty  = "-"
)

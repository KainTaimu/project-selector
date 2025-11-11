package cli

const (
	SaveCursor    = "\033[s"
	RestoreCursor = "\033[u"
	HideCursor    = "\033[?25l"
	ShowCursor    = "\033[?25h"
	HomeCursor    = "\033[H"
	ClearScreen   = "\033[J"
	ClearToEnd    = "\033[0J"
)

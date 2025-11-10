package cli

const (
	Bold = "\033[1m"
	End  = "\033[0m"
)

func ColorEntry(s string) string {
	return "\033[38;2;170;196;116m" + s + End
}

func ColorPop(s string) string {
	return Bold + "\033[38;2;170;66;66m" + s + End
}

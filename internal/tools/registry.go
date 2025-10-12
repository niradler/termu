package tools

var KnownCommands = map[string]bool{
	"sd":    true,
	"fd":    true,
	"rg":    true,
	"bat":   true,
	"xsv":   true,
	"jaq":   true,
	"yq":    true,
	"dua":   true,
	"eza":   true,
	"ls":    true,
	"cat":   true,
	"grep":  true,
	"find":  true,
	"echo":  true,
	"pwd":   true,
	"cd":    true,
	"git":   true,
	"curl":  true,
	"wget":  true,
	"rm":    true,
	"mv":    true,
	"cp":    true,
	"mkdir": true,
	"touch": true,
	"chmod": true,
	"chown": true,
	"sudo":  true,
}

func IsKnownCommand(cmd string) bool {
	return KnownCommands[cmd]
}

func GetDefaultAllowedCommands() []string {
	return []string{
		"sd", "fd", "rg", "bat", "xsv", "jaq", "yq", "dua", "eza",
		"ls", "cat", "grep", "find", "echo", "pwd", "git",
	}
}

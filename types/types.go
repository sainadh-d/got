package types

type Args struct {
	Command string   `arg:"positional"`
	Options []string `arg:"positional"`
}

type Repository struct {
	Worktree string
	Gitdir   string
	Conf     string
}

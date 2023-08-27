package types

type CheckoutCmd struct {
	Branch string `arg:"positional"`
	Track  bool   `arg:"-t"`
}

type InitCmd struct {
	WorkTree string `arg:"positional"`
}

type CatFileCmd struct {
	Hash        string `arg:"positional"`
	PrettyPrint bool   `arg:"-p"`
	Type        bool   `arg:"-t"`
	Size        bool   `arg:"-s"`
}

type HashObjectCmd struct {
	FileName string `arg:"positional"`
	Write    bool   `arg:"-w"`
}

type Test struct {
}

type Args struct {
	Checkout   *CheckoutCmd   `arg:"subcommand:checkout"`
	CatFile    *CatFileCmd    `arg:"subcommand:cat-file"`
	HashObject *HashObjectCmd `arg:"subcommand:hash-object"`
	Init       *InitCmd       `arg:"subcommand:init"`

	Test *Test `arg:"subcommand:test"`
}

package main

import (
	"fmt"
	"got/types"
	"os"

	"github.com/alexflint/go-arg"
	"github.com/bigkevmcd/go-configparser"
)

func main() {
	args := types.Args{}
	arg.MustParse(&args)
	switch cmd := args.Command; cmd {
	case "add":
		add(args)
		/*
			case "cat-file":
				cat_file(args)
			case "check-ignore":
				check_ignore(args)
			case "checkout":
				checkout(args)
			case "commit":
				commit(args)
			case "hash-object":
				hash_object(args)
		*/
	case "init":
		initialize(args)
		/*
			case "log":
				log(args)
			case "ls-files":
				ls_files(args)
			case "ls-tree":
				ls_tree(args)
			case "merge":
				merge(args)
			case "rebase":
				rebase(args)
			case "rev-parse":
				rev_parse(args)
			case "rm":
				rm(args)
			case "show-ref":
				show_ref(args)
			case "status":
				status(args)
			case "tag":
				tag(args)
		*/
	default:
		print("Bad command.")
	}
}

func createRepository(worktree string) types.Repository {

	return types.Repository{
		Worktree: worktree,
		Gitdir:   fmt.Sprintf("%s/.git", worktree),
	}
}

func add(args types.Args) {
}

func initialize(args types.Args) {
	worktree := args.Options[0]

	// Check if worktree exists and its a directory
	d, err := os.Stat(worktree)

	if os.IsNotExist(err) {
		fmt.Println(worktree, "doesn't exist")
		return
	}

	if !d.IsDir() {
		fmt.Println(worktree, "isn't a directory")
		return
	}

	repo := createRepository(worktree)

	// Create branches, objects, refs/tags, refs/heads directories
	if err = os.MkdirAll(fmt.Sprintf("%s/branches", repo.Gitdir), 0755); err != nil {
		fmt.Println("Failed to create branches", err)
	}

	if err = os.MkdirAll(fmt.Sprintf("%s/objects", repo.Gitdir), 0755); err != nil {
		fmt.Println("Failed to create objects", err)
	}

	if err = os.MkdirAll(fmt.Sprintf("%s/refs/tags", repo.Gitdir), 0755); err != nil {
		fmt.Println("Failed to create refs/tags", err)
	}

	if err = os.MkdirAll(fmt.Sprintf("%s/refs/heads", repo.Gitdir), 0755); err != nil {
		fmt.Println("Failed to create refs/heads", err)
	}

	// Create .git/description
	if f, err := os.Create(fmt.Sprintf("%s/description", repo.Gitdir)); err != nil {
		fmt.Println("Failed to create .git/description", err)
	} else {
		f.WriteString("Unnamed repository; edit this file 'description' to name the repository.\n")
		defer f.Close()
	}

	// Create .git/HEAD
	if f, err := os.Create(fmt.Sprintf("%s/HEAD", repo.Gitdir)); err != nil {
		fmt.Println("Failed to create .git/HEAD", err)
	} else {
		f.WriteString("ref: refs/heads/master\n")
		defer f.Close()
	}

	// Create .git/config
	gitconfig := fmt.Sprintf("%s/config", repo.Gitdir)
	if f, err := os.Create(gitconfig); err != nil {
		fmt.Println("Failed to create .git/config", err)
	} else {
		// Write default config
		p, err := configparser.NewConfigParserFromFile(gitconfig)
		if err != nil {
			fmt.Println("Failed to parse .git/config file", err)
		}

		// Set some basic config
		_ = p.AddSection("core")
		_ = p.Set("core", "repositoryformatversion", "0")
		_ = p.Set("core", "filemode", "false")
		_ = p.Set("core", "bare", "false")

		p.SaveWithDelimiter(gitconfig, "=")
		defer f.Close()
	}

}

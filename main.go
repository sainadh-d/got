package main

import (
	"fmt"
	"got/types"
	"got/utils"
	"os"

	"github.com/alexflint/go-arg"
	"github.com/bigkevmcd/go-configparser"
)

func main() {
	args := types.Args{}
	arg.MustParse(&args)
	switch {
	case args.CatFile != nil:
		cat_file(args)
	case args.Test != nil:
		test(args)
	case args.Init != nil:
		initialize(args)
		/*
		   case "add":
		       add(args)
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

	worktree := "."
	// if len(args.Input) == 1 {
	// 	worktree = args.Input[0]
	// }

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

func test(args types.Args) {
	repo := createRepository(".")

	content :=
		`tree 28758b08be51d44589fbbc1f1e2626a687e1c7f0
author Bubbles <sainadh976@gmail.com> 1693037671 -0700
committer GitHub <noreply@github.com> 1693037671 -0700
gpgsig -----BEGIN PGP SIGNATURE-----

 wsBcBAABCAAQBQJk6bRnCRBK7hj4Ov3rIwAAHiQIAE8v3y4RH0HvpN7OEy8C5Xaz
 QoO0Y9HkGgGW2IokAizG+ZD/05oiLo/QM6jOJ+4ZdWpMySyYs91wV6pOhRv5Z/P1
 05HUJVrjUiWlqVj8yARrLOB1+8bH2beqWl2L+5QQqLINyMPBzDfWhq71bu/Z87x0
 QoIB/iKpKSD2uBQIFPahNel0Zgu1+hpL4ixRuBNxfR7xK6lBZv5ZLp9hrYmNeikr
 0njxboOz1ptbk1bg+QtsGZYmyzCXETh7lfywxKevZ0hPquCnxaM0M2ABwBERpyDE
 617a7D16Z66wml99aN/6CNTC5oik7CIQ0zdGGpiOJTfBsR995UnX+LfOd1JRr3Y=
 =6xDz
 -----END PGP SIGNATURE-----


Initial commit`

	fmt.Println("Successfully wrote the content to file")
	gc := types.GitObject{
		Type: "commit",
		Size: len(content),
		Data: []byte(content),
	}

	// Write Something and Read it
	hash, err := repo.WriteObject(gc)
	if err != nil {
		return
	}
	fmt.Println(hash)

	c, e := repo.ReadObject(hash)
	if e != nil {
		fmt.Println(e)
		return
	}
	fmt.Println(string(c.Data))
	fmt.Println("Successfully read the content to file")
}

func cat_file(args types.Args) {
	repoRoot, err := utils.FindRepoRoot(".")
	if err != nil {
		fmt.Println(err)
		return
	}

	repo := createRepository(repoRoot)

	obj, err := repo.ReadObject(args.CatFile.Hash)
	if err != nil {
		fmt.Println(err)
		return
	}

	if args.CatFile.Size {
		fmt.Println(obj.Size)
	}

	if args.CatFile.Type {
		fmt.Println(obj.Type)
	}

	if args.CatFile.PrettyPrint {
		fmt.Println(string(obj.Data))
	}

}

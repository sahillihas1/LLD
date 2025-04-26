package main

import (
	"fmt"
)

// ----------------------------
// File Model
// ----------------------------
type File struct {
	Name    string
	Content string
}

func NewFile(name, content string) File {
	return File{Name: name, Content: content}
}

// ----------------------------
// Commit Model
// ----------------------------
type Commit struct {
	Message string
	Files   []File
}

func NewCommit(message string, files []File) Commit {
	copied := make([]File, len(files))
	copy(copied, files)
	return Commit{
		Message: message,
		Files:   copied,
	}
}

// ----------------------------
// Branch Model
// ----------------------------
type Branch struct {
	Name    string
	Commits []Commit
	Staged  []File
}

func NewBranch(name string) *Branch {
	return &Branch{
		Name:    name,
		Commits: []Commit{},
		Staged:  []File{},
	}
}

func (b *Branch) StageFile(f File) {
	b.Staged = append(b.Staged, f)
}

func (b *Branch) Commit(msg string) {
	if len(b.Staged) == 0 {
		fmt.Println("Nothing to commit.")
		return
	}
	newCommit := NewCommit(msg, b.Staged)
	b.Commits = append(b.Commits, newCommit)
	b.Staged = []File{}
	fmt.Printf("Committed to branch '%s': %s\n", b.Name, msg)
}

func (b *Branch) Clone(newName string) *Branch {
	newBranch := NewBranch(newName)
	newBranch.Commits = append(newBranch.Commits, b.Commits...)
	return newBranch
}

func (b *Branch) Rollback() {
	if len(b.Commits) == 0 {
		fmt.Println("Nothing to rollback.")
		return
	}
	b.Commits = b.Commits[:len(b.Commits)-1]
	fmt.Println("Rollback successful.")
}

func (b *Branch) ShowCommits() {
	fmt.Printf("\nCommits on branch '%s':\n", b.Name)
	for i, c := range b.Commits {
		fmt.Printf("  %d. %s\n", i+1, c.Message)
	}
	fmt.Println()
}

// ----------------------------
// Version Control Interface
// ----------------------------
type VersionControl interface {
	AddFile(name, content string)
	Commit(message string)
	CreateBranch(name string)
	SwitchBranch(name string)
	Rollback()
	ShowHistory()
}

// ----------------------------
// Git System Implementation
// ----------------------------
type GitSystem struct {
	Branches   map[string]*Branch
	CurrBranch *Branch
}

func NewGitSystem() VersionControl {
	main := NewBranch("main")
	return &GitSystem{
		Branches:   map[string]*Branch{"main": main},
		CurrBranch: main,
	}
}

func (g *GitSystem) AddFile(name, content string) {
	g.CurrBranch.StageFile(NewFile(name, content))
}

func (g *GitSystem) Commit(message string) {
	g.CurrBranch.Commit(message)
}

func (g *GitSystem) CreateBranch(name string) {
	if _, exists := g.Branches[name]; exists {
		fmt.Printf("Branch '%s' already exists.\n", name)
		return
	}
	g.Branches[name] = g.CurrBranch.Clone(name)
	fmt.Printf("Branch '%s' created.\n", name)
}

func (g *GitSystem) SwitchBranch(name string) {
	if branch, ok := g.Branches[name]; ok {
		g.CurrBranch = branch
		fmt.Printf("Switched to branch '%s'\n", name)
	} else {
		fmt.Printf("Branch '%s' does not exist.\n", name)
	}
}

func (g *GitSystem) Rollback() {
	g.CurrBranch.Rollback()
}

func (g *GitSystem) ShowHistory() {
	g.CurrBranch.ShowCommits()
}

// ----------------------------
// Main Function (Testing)
// ----------------------------
func main() {
	vcs := NewGitSystem()

	vcs.AddFile("file1.txt", "Hello World")
	vcs.Commit("Initial commit")

	vcs.CreateBranch("feature-x")
	vcs.SwitchBranch("feature-x")

	vcs.AddFile("file2.txt", "Feature X started")
	vcs.Commit("Feature X - First Commit")

	vcs.ShowHistory()

	vcs.Rollback()
	vcs.ShowHistory()

	vcs.SwitchBranch("main")
	vcs.ShowHistory()
}

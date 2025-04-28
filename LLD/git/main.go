package main

import (
	"fmt"
	"time"
)

// File represents a basic file structure.
type File struct {
	Name    string
	Content string
}

// ICommit defines the contract for a commit.
type ICommit interface {
	GetID() int
	GetFiles() map[string]string
	GetMessage() string
	GetTimestamp() time.Time
	GetParent() ICommit
}

// Commit is the concrete implementation of ICommit.
type Commit struct {
	id        int
	files     map[string]string
	message   string
	timestamp time.Time
	parent    ICommit
}

func (c *Commit) GetID() int                  { return c.id }
func (c *Commit) GetFiles() map[string]string { return c.files }
func (c *Commit) GetMessage() string          { return c.message }
func (c *Commit) GetTimestamp() time.Time     { return c.timestamp }
func (c *Commit) GetParent() ICommit          { return c.parent }

// IBranch defines the contract for a branch.
type IBranch interface {
	GetName() string
	GetHead() ICommit
	SetHead(commit ICommit)
	AddCommit(commit ICommit)
	GetCommits() []ICommit
	SetCommits(commits []ICommit)
}

// Branch is the concrete implementation of IBranch.
type Branch struct {
	name       string
	head       ICommit
	commitList []ICommit
}

func (b *Branch) GetName() string              { return b.name }
func (b *Branch) GetHead() ICommit             { return b.head }
func (b *Branch) SetHead(commit ICommit)       { b.head = commit }
func (b *Branch) AddCommit(commit ICommit)     { b.commitList = append(b.commitList, commit) }
func (b *Branch) GetCommits() []ICommit        { return b.commitList }
func (b *Branch) SetCommits(commits []ICommit) { b.commitList = commits }

// ICommand defines the command interface for actions like commit, add, rollback.
type ICommand interface {
	Execute()
}

// AddFileCommand adds a file to the staging area.
type AddFileCommand struct {
	vc   *VersionControl
	file File
}

func (cmd *AddFileCommand) Execute() {
	cmd.vc.stagingArea[cmd.file.Name] = cmd.file.Content
}

// CommitCommand handles committing staged files.
type CommitCommand struct {
	vc      *VersionControl
	message string
}

func (cmd *CommitCommand) Execute() {
	files := make(map[string]string)
	head := cmd.vc.current.GetHead()
	if head != nil {
		for k, v := range head.GetFiles() {
			files[k] = v
		}
	}
	for k, v := range cmd.vc.stagingArea {
		files[k] = v
	}
	commit := &Commit{
		id:        cmd.vc.commitID,
		files:     files,
		message:   cmd.message,
		timestamp: time.Now(),
		parent:    head,
	}
	cmd.vc.current.SetHead(commit)
	cmd.vc.current.AddCommit(commit)
	cmd.vc.commitID++
	cmd.vc.stagingArea = make(map[string]string)
}

// RollbackCommand handles rollback functionality.
type RollbackCommand struct {
	vc       *VersionControl
	commitID int
}

func (cmd *RollbackCommand) Execute() {
	commits := cmd.vc.current.GetCommits()
	for i := len(commits) - 1; i >= 0; i-- {
		if commits[i].GetID() == cmd.commitID {
			cmd.vc.current.SetHead(commits[i])
			cmd.vc.current.SetCommits(commits[:i+1])
			fmt.Println("Rolled back to commit", cmd.commitID)
			return
		}
	}
	fmt.Println("Commit ID not found")
}

// RevertCommand reverts to a specific commit, preserving history.
type RevertCommand struct {
	vc       *VersionControl
	commitID int
}

func (cmd *RevertCommand) Execute() {
	var target ICommit
	for _, c := range cmd.vc.current.GetCommits() {
		if c.GetID() == cmd.commitID {
			target = c
			break
		}
	}
	if target == nil {
		fmt.Println("Commit ID not found")
		return
	}

	head := cmd.vc.current.GetHead()
	revertCommit := &Commit{
		id:        cmd.vc.commitID,
		files:     target.GetFiles(),
		message:   fmt.Sprintf("Revert to commit %d", cmd.commitID),
		timestamp: time.Now(),
		parent:    head,
	}
	cmd.vc.current.SetHead(revertCommit)
	cmd.vc.current.AddCommit(revertCommit)
	cmd.vc.commitID++
	fmt.Println("Reverted to commit", cmd.commitID)
}

// VersionControl orchestrates version control features.
type VersionControl struct {
	branches    map[string]IBranch
	current     IBranch
	stagingArea map[string]string
	commitID    int
}

func NewVersionControl() *VersionControl {
	master := &Branch{name: "master"}
	vc := &VersionControl{
		branches:    map[string]IBranch{"master": master},
		current:     master,
		stagingArea: make(map[string]string),
		commitID:    0,
	}
	return vc
}

func (vc *VersionControl) RunCommand(cmd ICommand) {
	cmd.Execute()
}

func (vc *VersionControl) CreateBranch(name string) {
	newBranch := &Branch{
		name:       name,
		head:       vc.current.GetHead(),
		commitList: append([]ICommit{}, vc.current.GetCommits()...),
	}
	vc.branches[name] = newBranch
}

func (vc *VersionControl) CheckoutBranch(name string) {
	if branch, ok := vc.branches[name]; ok {
		vc.current = branch
	} else {
		fmt.Println("Branch does not exist")
	}
}

func main() {
	vc := NewVersionControl()

	vc.RunCommand(&AddFileCommand{vc, File{"file1.txt", "Hello World"}})
	vc.RunCommand(&CommitCommand{vc, "Initial commit"})

	vc.RunCommand(&AddFileCommand{vc, File{"file2.txt", "Another file"}})
	vc.RunCommand(&CommitCommand{vc, "Added file2"})

	vc.CreateBranch("feature")
	vc.CheckoutBranch("feature")

	vc.RunCommand(&AddFileCommand{vc, File{"file1.txt", "Updated in feature"}})
	vc.RunCommand(&CommitCommand{vc, "Updated file1 in feature branch"})

	vc.RunCommand(&RevertCommand{vc, 0}) // Revert to first commit while preserving history

	fmt.Println("Current HEAD ID:", vc.current.GetHead().GetID())
	fmt.Println("Current HEAD files:", vc.current.GetHead().GetFiles())
}

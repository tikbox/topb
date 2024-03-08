package main

// gen:topb
//
//go:generate go run main.go -in user.go -pb your-git.com/your-pb-path/pb
type User struct {
	Uid  uint32
	Name string
}

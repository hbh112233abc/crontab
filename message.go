package main

import "fmt"

type Message struct {
}

func Msg() *Message {
	return &Message{}
}

func (m *Message) Info(msg string) {
	fmt.Printf("\x1b[%dm "+msg+" \x1b[0m\n", 34)
}
func (m *Message) Error(msg string) {
	fmt.Printf("\x1b[%dm "+msg+" \x1b[0m\n", 31)
}
func (m *Message) Warning(msg string) {
	fmt.Printf("\x1b[%dm "+msg+" \x1b[0m\n", 33)
}
func (m *Message) Success(msg string) {
	fmt.Printf("\x1b[%dm "+msg+" \x1b[0m\n", 32)
}
func (m *Message) Default(msg string) {
	fmt.Printf("\x1b[%dm "+msg+" \x1b[0m\n", 37)
}

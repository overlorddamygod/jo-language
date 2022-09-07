package io

import (
	"fmt"
	"sync"

	"github.com/bwmarrin/discordgo"
)

var WaitingForInput = Waiting{
	Val: false,
}

type DiscordIO struct {
	dg      *discordgo.Session
	channel string
}

func NewDiscordIO(dg *discordgo.Session, channel string) DiscordIO {
	return DiscordIO{
		dg:      dg,
		channel: channel,
		// waitingForInput: false,
	}
}

func (d DiscordIO) Print(a ...interface{}) (n int, err error) {
	d.dg.ChannelMessageSend(d.channel, fmt.Sprint(a...))
	return 1, nil
}

func (d DiscordIO) Println(a ...interface{}) (n int, err error) {
	d.dg.ChannelMessageSend(d.channel, fmt.Sprint(a...))
	return 1, nil
}

func (d DiscordIO) Error(a string) (n int, err error) {
	return d.Println(a)
}

func (d DiscordIO) Printf(format string, a ...interface{}) (n int, err error) {
	d.dg.ChannelMessageSend(d.channel, fmt.Sprintf(format, a...))
	return 1, nil
}

type Waiting struct {
	mu  sync.Mutex
	Val bool
}

func (d DiscordIO) Input() string {
	WaitingForInput.mu.Lock()

	WaitingForInput.Val = true
	textChan := make(chan string)
	defer func() {
		close(textChan)
	}()
	// d.dg.
	d.dg.AddHandlerOnce(func(s *discordgo.Session, m *discordgo.MessageCreate) {

		// Ignore all messages created by the bot itself
		// This isn't required in this specific example but it's a good practice.
		if m.Author.ID == s.State.User.ID {
			return
		}

		fmt.Println("INPUT", m.Content)

		textChan <- m.Content
	})
	val := <-textChan

	WaitingForInput.Val = false
	WaitingForInput.mu.Unlock()
	return val
}

package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
	"github.com/overlorddamygod/jo/cmd/discord-interpreter/io"
	"github.com/overlorddamygod/jo/internal/eval"
	"github.com/overlorddamygod/jo/pkg/lexer"
	"github.com/overlorddamygod/jo/pkg/parser"
	"github.com/overlorddamygod/jo/pkg/stdio"
)

type DiscordInterpreter struct {
	lexer       lexer.Lexer
	io          io.DiscordIO
	Interpreter *eval.Evaluator
}

func NewDiscordInterpreter(io io.DiscordIO) *DiscordInterpreter {
	lexer := *lexer.NewLexer("")

	di := &DiscordInterpreter{
		lexer:       lexer,
		io:          io,
		Interpreter: eval.NewEvaluator(&lexer, []parser.Node{}),
	}
	return di
}

func (d *DiscordInterpreter) Interpret(src string) {
	if io.WaitingForInput.Val {
		return
	}
	d.lexer = *lexer.NewLexer(src)

	_, _, err := d.lexer.Lex()
	if err != nil {
		// stdio.Io.Print(tokens)
		stdio.Io.Print("[Lexer]\n\n", err)
		return
	}

	parser := parser.NewParser(&d.lexer)

	node, err := parser.Parse()

	if err != nil {
		stdio.Io.Print("[Parser]\n\n", err)
		return
	}

	d.Interpreter.SetLexerNode(&d.lexer, node)

	_, err = d.Interpreter.Eval()

	if err != nil {
		stdio.Io.Printf("[Evaluator]\n\n%s", err)
	}
}

func (d *DiscordInterpreter) messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {

	// Ignore all messages created by the bot itself
	if m.Author.ID == s.State.User.ID {
		return
	}

	fmt.Println(m.Author.Username, m.Content, m.ChannelID)

	d.Interpret(m.Content)
}

// Variables used for command line parameters

func main() {
	token := os.Getenv("discord_bot_token")

	if token == "" {
		panic("Discord bot token not provided.")
	}
	// Create a new Discord session using the provided bot token.
	dg, err := discordgo.New("Bot " + token)
	if err != nil {
		fmt.Println("error creating Discord session,", err)
		return
	}

	// Register the messageCreate func as a callback for MessageCreate events.
	discordIO := io.NewDiscordIO(dg, "1015964107763093654")
	stdio.SetIO(discordIO)

	di := NewDiscordInterpreter(discordIO)
	dg.AddHandler(di.messageCreate)
	// Open a websocket connection to Discord and begin listening.``
	err = dg.Open()
	if err != nil {
		fmt.Println("error opening connection,", err)
		return
	}

	// Wait here until CTRL-C or other term signal is received.
	fmt.Println("Bot is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	// Cleanly close down the Discord session.
	dg.Close()
}

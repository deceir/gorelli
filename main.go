package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/bwmarrin/discordgo"
)

func main() {
	// Read the bot token from an environment variable.
	//
	// Do not paste your token directly into the source code.
	// Accidentally uploading a token to GitHub would allow other people
	// to control your bot.
	token := os.Getenv("DISCORD_BOT_TOKEN")

	if token == "" {
		log.Fatal("DISCORD_BOT_TOKEN environment variable is not set")
	}

	// Create a new Discord session.
	//
	// Discord requires bot authentication strings to begin with "Bot ".
	session, err := discordgo.New("Bot " + token)
	if err != nil {
		log.Fatalf("could not create Discord session: %v", err)
	}

	// Tell Discord which types of events this bot wants to receive.
	//
	// IntentsGuilds:
	// Gives us basic server/guild-related events.
	//
	// IntentsGuildMessages:
	// Lets us receive events when messages are posted in servers.
	//
	// IntentsMessageContent:
	// Lets us read the actual text inside those messages.
	session.Identify.Intents =
		discordgo.IntentsGuilds |
			discordgo.IntentsGuildMessages |
			discordgo.IntentsMessageContent

	// Register messageCreate as the function that should run
	// whenever Discord sends us a new-message event.
	session.AddHandler(messageCreate)

	// Open the WebSocket connection to Discord.
	err = session.Open()
	if err != nil {
		log.Fatalf("could not connect to Discord: %v", err)
	}

	// Close the Discord connection when main() exits.
	defer func() {
		if err := session.Close(); err != nil {
			log.Printf("error closing Discord session: %v", err)
		}
	}()

	fmt.Println("Bot is online.")
	fmt.Println("Type !ping in your Discord server.")
	fmt.Println("Press Ctrl+C to stop the bot.")

	// Create a channel that receives operating-system shutdown signals.
	//
	// This prevents main() from immediately finishing and disconnecting
	// the bot. It waits here until you press Ctrl+C or terminate the app.
	stop := make(chan os.Signal, 1)

	signal.Notify(
		stop,
		os.Interrupt,
		syscall.SIGTERM,
	)

	// Block the program until a shutdown signal arrives.
	<-stop

	fmt.Println("\nShutting down...")
}

// messageCreate is called every time the bot receives a message event.
func messageCreate(session *discordgo.Session, message *discordgo.MessageCreate) {
	// Ignore messages sent by the bot itself.
	//
	// Without this check, a bot could respond to its own messages
	// and accidentally create an infinite response loop.
	if message.Author.ID == session.State.User.ID {
		return
	}

	// Remove leading/trailing whitespace and make the command lowercase.
	//
	// That means "!PING", "!ping", and "  !ping  " all work.
	command := strings.ToLower(strings.TrimSpace(message.Content))

	switch command {
	case "!ping":
		// Send a message in the same channel where !ping was received.
		_, err := session.ChannelMessageSend(message.ChannelID, "Pong! 🏓")
		if err != nil {
			log.Printf("could not send message: %v", err)
		}

	case "!hello":
		// message.Author.Mention() creates a Discord mention,
		// such as <@123456789>.
		reply := fmt.Sprintf(
			"Hello, %s! The Go bot is working.",
			message.Author.Mention(),
		)

		_, err := session.ChannelMessageSend(message.ChannelID, reply)
		if err != nil {
			log.Printf("could not send message: %v", err)
		}
	}
}

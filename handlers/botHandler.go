package handlers

import (
	"context"
	"log"
	"os"
	"os/signal"
	"sync"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

type Message struct {
	UserID int64
	Msg    string
}

var MessageChannel = make(chan Message, 100)
var botInstance *bot.Bot
var botContext context.Context
var botCancel context.CancelFunc
var botWG sync.WaitGroup

func SendMessageToUser(ctx context.Context, userID int64, message string) {
	log.Printf("Sending message to user %d ...", userID)
	if botInstance == nil {
		log.Println("Bot instance is not initialized.")
		return
	}
	_, err := botInstance.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: userID,
		Text:   message,
	})
	if err != nil {
		log.Printf("Error sending message to user %d: %v", userID, err)
	}
}

func NewBotHandler(token string) (*bot.Bot, error) {
    log.Println("Starting new bot handler ...")

    ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
    botContext = ctx
    botCancel = cancel

    opts := []bot.Option{
        bot.WithNotAsyncHandlers(),
        bot.WithDefaultHandler(handler),
    }

    b, err := bot.New(token, opts...)
    if err != nil {
        return nil, err
    }
    botInstance = b

    botWG.Add(1)
    go func() {
        defer botWG.Done()
        b.Start(botContext)
        log.Println("Bot stopped")
    }()

    botWG.Add(1)
    go func() {
        defer botWG.Done()
        for msg := range MessageChannel {
            SendMessageToUser(botContext, msg.UserID, msg.Msg)
        }
        log.Println("Message channel closed")
        botCancel()
    }()

    return b, nil
}

func handler(ctx context.Context, b *bot.Bot, update *models.Update) {
	log.Println("Handling update ...")
	if update.Message == nil {
		return
	}
}

func SendMessagesAndShutdown(messages []Message) {
    for _, msg := range messages {
        MessageChannel <- msg
    }
    close(MessageChannel)
    botWG.Wait()
}
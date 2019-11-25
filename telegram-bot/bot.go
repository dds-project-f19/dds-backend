package telegram_bot

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"log"
	"os"
	"time"
)

type a struct {
	Username string `gorm:"unique_index;not null"`
	ChatID   string `gorm:"unique"`
}

var commandKeyboard = tgbotapi.NewReplyKeyboard(
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("/schedule"),
	),
)

func sendDelayedMessage(bot *tgbotapi.BotAPI, message string, chatID int64) {
	time.Sleep(2 * time.Minute)
	msg := tgbotapi.NewMessage(chatID, message)
	_, err := bot.Send(msg)
	if err != nil {
		log.Panic(err.Error())
	}
}

func LaunchBot() {
	bot, err := tgbotapi.NewBotAPI(os.Getenv("DDS_TELEGRAM_BOT_APIKEY"))
	if err != nil {
		log.Panic(err)
	}
	bot.Debug = false
	log.Printf("Authorized on account %s", bot.Self.UserName)

	var updateConf tgbotapi.UpdateConfig = tgbotapi.NewUpdate(0)
	updateConf.Timeout = 60
	updatesChan, err := bot.GetUpdatesChan(updateConf)

	// TODO authorize user, possibly via t.me/dds_project_f19_bot?start=some123key456
	// TODO handle /schedule command
	// TODO send notifications
	for {
		select {
		case update := <-updatesChan:
			var msg tgbotapi.MessageConfig
			if update.Message == nil {
				log.Println("Message is nil")
				return
			}
			if update.Message.IsCommand() && update.Message.Command() == "start" {
				var cmd string
				var key string
				readNum, _ := fmt.Sscanf(update.Message.Text, "%s %s", &cmd, &key)
				msg.ChatID = update.Message.Chat.ID
				if readNum != 2 {
					msg.Text = "Sorry, you don't have access to this bot."
				} else {
					// validate key, register chatid for this user
					msg.Text = fmt.Sprintf("You are registering with key: %s", key)
				}
			} else if update.Message.IsCommand() {
				// check if not registered and refuse further communication in that case
				msg = tgbotapi.NewMessage(update.Message.Chat.ID, "")
				switch update.Message.Command() {
				case "help":
					msg.Text = "type /schedule to know your schedule"
				case "schedule":
					msg.Text = "This is your schedule: \n1 2 3"
				default:
					msg.Text = "I don't know that command"
				}
			} else {
				chatID := update.Message.Chat.ID
				msg = tgbotapi.NewMessage(chatID, "Bazinga")
				go sendDelayedMessage(bot, "Delayed Message", chatID)
			}

			msg.ReplyMarkup = commandKeyboard
			_, err := bot.Send(msg)
			if err != nil {
				log.Println(err.Error())
			}
		}
	}
}

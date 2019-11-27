package services

import (
	"dds-backend/common"
	"dds-backend/database"
	"dds-backend/models"
	"errors"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"log"
	"os"
	"time"
)

const (
	RegistrationTokenExpirationDuration = time.Minute * 10
)

func GetChatRegistrationLink(username string) (string, error) {
	chat := models.TelegramChat{
		Username: username,
	}
	res := database.DB.Model(&models.TelegramChat{}).Where(&chat).First(&chat)
	if res.RecordNotFound() {
		// need to be registered
		chat.RegistrationToken = common.GenerateNewToken()
		chat.TokenExpiration = time.Now().Add(RegistrationTokenExpirationDuration)
		res = database.DB.Model(&models.TelegramChat{}).Create(&chat)
		if res.Error != nil {
			return "", res.Error
		}
	} else if res.Error != nil {
		return "", res.Error
	} else {
		chat.RegistrationToken = common.GenerateNewToken()
		chat.TokenExpiration = time.Now().Add(RegistrationTokenExpirationDuration)
		res = database.DB.Model(&models.TelegramChat{}).Save(&chat)
		if res.Error != nil {
			return "", res.Error
		}
	}
	return fmt.Sprintf("t.me/dds_project_f19_bot?start=%s", chat.RegistrationToken), nil
}

func GetUsernameByChat(chatID int64) (string, error) {
	chat := models.TelegramChat{ChatID: chatID}
	res := database.DB.Model(&models.TelegramChat{}).Where(&chat).First(&chat)
	if res.Error != nil {
		return "", errors.New("could not get this chat")
	}
	return chat.Username, nil
}

func GetChatIDByUsername(username string) (int64, error) {
	chat := models.TelegramChat{Username: username}
	res := database.DB.Model(&models.TelegramChat{}).Where(&chat).First(&chat)
	if res.Error != nil {
		return 0, errors.New("could not get this chat")
	}
	return chat.ChatID, nil
}

// Requested when user logs in via Telegram
func ValidateChat(registrationToken string, chatID int64) error {
	chat := models.TelegramChat{
		RegistrationToken: registrationToken,
	}
	res := database.DB.Model(&models.TelegramChat{}).Where(&chat).First(&chat)
	if res.RecordNotFound() {
		// registration token does not exist
		return errors.New("can't validate this token")
	} else if res.Error != nil {
		// unexpected error
		return errors.New("something went wrong")
	} else {
		// ok
		if chat.TokenExpiration.Before(time.Now()) {
			return errors.New("token has expired")
		}
		chat.ChatID = chatID
		chat.RegistrationToken = common.GenerateNewToken(chat.Username) // to erase previous token
		res = database.DB.Model(&models.TelegramChat{}).Save(&chat)
		if res.Error != nil {
			return errors.New("registration failed")
		}
	}
	return nil
}

var commandKeyboard = tgbotapi.NewReplyKeyboard(
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("/schedule"),
	),
)

func sendDelayedMessage(bot *tgbotapi.BotAPI, message string, chatID int64) {
}

var BotInstance *tgbotapi.BotAPI

func SendNotification(username string, text string) error {
	chatID, err := GetChatIDByUsername(username)
	if err != nil {
		return err
	}
	if BotInstance != nil {
		msg := tgbotapi.NewMessage(chatID, text)
		_, err := BotInstance.Send(msg)
		if err != nil {
			log.Panic(err.Error())
		}
	} else {
		return errors.New("bot instance is nil")
	}
	return nil
}

func LaunchBot() {
	BotInstance, err := tgbotapi.NewBotAPI(os.Getenv("DDS_TELEGRAM_BOT_APIKEY"))
	if err != nil {
		log.Panic(err)
	}
	BotInstance.Debug = false
	log.Printf("Authorized on account %s", BotInstance.Self.UserName)

	var updateConf tgbotapi.UpdateConfig = tgbotapi.NewUpdate(0)
	updateConf.Timeout = 60
	updatesChan, err := BotInstance.GetUpdatesChan(updateConf)

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
			// handle login of new user
			if update.Message.IsCommand() && update.Message.Command() == "start" {
				var cmd string
				var key string
				readNum, _ := fmt.Sscanf(update.Message.Text, "%s %s", &cmd, &key)
				msg.ChatID = update.Message.Chat.ID
				if readNum != 2 {
					msg.Text = "Sorry, you don't have access to this BotInstance."
				} else {
					// validate key, register chatid for this user
					msg.Text = fmt.Sprintf("You are registering with key: %s", key)
					err := ValidateChat(key, msg.ChatID)
					if err != nil {
						msg.Text = err.Error()
					} else {
						msg.Text = "Welcome to DDS Schedule Bot!\n You are successfully registered. See /help for available commands."
					}
				}
				// handle commands of existing user
			} else if update.Message.IsCommand() {
				// check if not registered and refuse further communication in that case
				username, err := GetUsernameByChat(update.Message.Chat.ID)
				if err != nil {
					msg.Text = ""
				} else {
					msg = tgbotapi.NewMessage(update.Message.Chat.ID, "")
					switch update.Message.Command() {
					case "help":
						msg.Text = "type /schedule to know your schedule"
					case "schedule":
						msg.Text = fmt.Sprintf("Schedule for %s", username)
					default:
						msg.Text = "I don't know that command"
					}
				}
			} else {
				chatID := update.Message.Chat.ID
				msg = tgbotapi.NewMessage(chatID, "Bazinga")
				go sendDelayedMessage(BotInstance, "Delayed Message", chatID)
			}

			msg.ReplyMarkup = commandKeyboard
			_, err := BotInstance.Send(msg)
			if err != nil {
				log.Println(err.Error())
			}
		}
	}
}

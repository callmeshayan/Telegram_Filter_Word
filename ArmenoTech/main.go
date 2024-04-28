package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	_ "github.com/lib/pq"
)

const (
	host   = "localhost"
	port   = 5432
	user   = "postgres"
	dbname = "Telegram_Fliter_bot"
)

var filteredWord string

func main() {
	// Get password from the terminal
	password := os.Getenv("POSTGRES_PASSWORD")

	// Initialize the bot
	token := os.Getenv("BOT_TOKEN")
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Fatal(err)
	}

	psqlInfo := fmt.Sprintf("postgresql://postgres:%s@192.168.1.16/Telegram_Fliter_bot?sslmode=disable", password)

	bot.Debug = true

	// Set up long polling to receive updates
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)
	if err != nil {
		log.Fatal(err)
	}

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	// Handle incoming messages
	for update := range updates {
		if update.Message != nil {
			if update.Message.IsCommand() {
				handleCommand(bot, db, update)
			} else {
				handleMessage(bot, db, update)
			}
		}
	}
}

func handleCommand(bot *tgbotapi.BotAPI, db *sql.DB, update tgbotapi.Update) {
	switch update.Message.Command() {
	case "start":
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Welcome to the Telegram Filter Bot!")
		bot.Send(msg)
	case "filter":
		handleFilterCommand(bot, db, update)
	case "help":
		handleHelpCommand(bot, update)
	case "stop":
		handleStopCommand(bot, update)
	default:
		// Ignore unsupported commands
	}
}

func handleFilterCommand(bot *tgbotapi.BotAPI, db *sql.DB, update tgbotapi.Update) {
	filteredWord = update.Message.CommandArguments()
	log.Println("Received command arguments:", filteredWord)
	if filteredWord == "" {
		// If no word is provided after "/filter", notify the user
		reply := tgbotapi.NewMessage(update.Message.Chat.ID, "Please provide a word after '/filter' command.")
		bot.Send(reply)
		return
	}
	sendingDate := time.Unix(int64(update.Message.Date), 0)
	senderID := update.Message.From.ID
	messageID := update.Message.MessageID
	messageText := update.Message.Text
	err := insertfilteredWord(db, filteredWord, sendingDate, senderID, messageID, messageText)
	if err != nil {
		log.Println("Error inserting filter word:", err)
		reply := tgbotapi.NewMessage(update.Message.Chat.ID, "Failed to add filter word.")
		bot.Send(reply)
		return
	}

	replyText := fmt.Sprintf("Filter word '%s' added successfully!", filteredWord)
	reply := tgbotapi.NewMessage(update.Message.Chat.ID, replyText)
	bot.Send(reply)

	// Print whether the word was added or not in the chat box
	//reply = tgbotapi.NewMessage(update.Message.Chat.ID, replyText)
	//bot.Send(reply)
}

func handleHelpCommand(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	// Construct the help message
	helpMessage := "Welcome to the Telegram Filter Bot!\n\n"
	helpMessage += "Available commands:\n"
	helpMessage += "/filter 'word': Filter messages containing <word>\n"
	helpMessage += "/help: Display this help message\n"
	helpMessage += "/stop: Stop the bot\n\n"
	helpMessage += "I.E:\n"
	helpMessage += "To filter messages containing the word 'example', use /filter example\n"

	// Send the help message to the user
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, helpMessage)
	_, err := bot.Send(msg)
	if err != nil {
		log.Println("Error sending help message:", err)
	}
}

func handleStopCommand(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Bot stopped.")
	bot.Send(msg)
	os.Exit(0)
}

func handleMessage(bot *tgbotapi.BotAPI, db *sql.DB, update tgbotapi.Update) {

	// Extract necessary details from the Telegram message
	sendingDate := time.Unix(int64(update.Message.Date), 0)
	senderID := update.Message.From.ID
	messageID := update.Message.MessageID
	messageText := update.Message.Text

	// Check if the message contains the filter word
	messageContainsfilteredWord := containsfilteredWord(messageText, filteredWord)

	// Send a message indicating whether the message contains the filter word or not
	var resultMessage string
	if messageContainsfilteredWord {
		resultMessage = fmt.Sprintln("Message ID contains the filter word: ", filteredWord)
	} else {
		resultMessage = fmt.Sprintln("Message ID does not contain the filter word:", filteredWord)
	}

	reply := tgbotapi.NewMessage(update.Message.Chat.ID, resultMessage)
	_, err := bot.Send(reply)
	if err != nil {
		log.Println("Error sending result message:", err)
	}

	// Insert message details into the database
	if messageContainsfilteredWord {
		// Message contains the filter word, insert into filtered_messages table
		err = insertfilteredWord(db, filteredWord, sendingDate, senderID, messageID, messageText)
	} else {
		// Message does not contain the filter word, insert into non_filtered_messages table
		err = insertNonFilteredMessage(db, sendingDate, senderID, messageID, messageText)
	}
	if err != nil {
		log.Println("Error inserting message details into database:", err)
	}
}

// Check if a message contains the filter word
func containsfilteredWord(message, filteredWord string) bool {
	// Convert both message and filteredWord to lowercase for case-insensitive matching
	message = strings.ToLower(message)
	filteredWord = strings.ToLower(filteredWord)

	// Check if the message contains the filter word
	return strings.Contains(message, filteredWord)
}

// Insert the filter word into the database
func insertfilteredWord(db *sql.DB, word string, sendingDate time.Time, senderID int, messageID int, messageText string) error {
	_, err := db.Exec("INSERT INTO filtered_messages (word, sending_date, sender_id, message_id, message_text) VALUES ($1, $2, $3, $4, $5)", word, sendingDate, senderID, messageID, messageText)
	return err
}
func insertNonFilteredMessage(db *sql.DB, sendingDate time.Time, senderID int, messageID int, messageText string) error {
	_, err := db.Exec("INSERT INTO non_filtered_messages (sending_date, sender_id, message_id, message_text) VALUES ($1, $2, $3, $4)", sendingDate, senderID, messageID, messageText)
	return err
}

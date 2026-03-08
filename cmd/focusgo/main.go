package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"focusgo/internal/database"
)

var bot *tgbotapi.BotAPI

func main() {
	// Инициализация БД
	if err := database.InitDB("focusgo.db"); err != nil {
		log.Fatalf("❌ Ошибка инициализации БД: %v", err)
	}
	defer database.CloseDB()

	// Автоматические бэкапы каждые 24 часа
	database.AutoBackup(24 * time.Hour)

	// Обработка сигналов
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		log.Println("🛑 Получен сигнал завершения...")
		database.CloseDB()
		os.Exit(0)
	}()

	// Токен бота
	token := os.Getenv("TELEGRAM_BOT_TOKEN")
	if token == "" {
		token = "YOUR_BOT_TOKEN_HERE"
		fmt.Println("⚠️  Установите TELEGRAM_BOT_TOKEN")
	}

	var err error
	bot, err = tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true
	log.Printf("✅ Бот авторизован: %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message != nil {
			handleMessage(update.Message)
		} else if update.CallbackQuery != nil {
			handleCallback(update.CallbackQuery)
		}
	}
}

func handleMessage(message *tgbotapi.Message) {
	if message.IsCommand() {
		handleCommand(message)
	}
}

func handleCommand(message *tgbotapi.Message) {
	chatID := message.Chat.ID
	command := message.Command()

	switch command {
	case "start":
		sendStart(chatID)
	case "backup":
		handleBackup(chatID)
	case "help":
		sendHelp(chatID)
	default:
		sendUnknown(chatID)
	}
}

func sendStart(chatID int64) {
	text := `🎮 <b>FOCUSGO</b>

Симулятор борьбы с искушениями!

Команды:
/backup — Бэкап БД
/help — Справка`

	msg := tgbotapi.NewMessage(chatID, text)
	msg.ParseMode = "HTML"
	bot.Send(msg)
}

func handleBackup(chatID int64) {
	backupPath, err := database.CreateBackup()
	if err != nil {
		text := fmt.Sprintf("❌ Ошибка: %v", err)
		bot.Send(tgbotapi.NewMessage(chatID, text))
		return
	}

	text := fmt.Sprintf("💾 Бэкап создан: %s", backupPath)
	bot.Send(tgbotapi.NewMessage(chatID, text))
}

func sendHelp(chatID int64) {
	text := `📖 <b>СПРАВКА</b>

Команды:
/start — Старт
/backup — Бэкап БД
/help — Справка

Уведомления:
• 9:00 — Квесты
• 20:00 — Битва
• 22:00 — Незавершённые квесты

Автосохранение: каждые 24 часа`

	msg := tgbotapi.NewMessage(chatID, text)
	msg.ParseMode = "HTML"
	bot.Send(msg)
}

func sendUnknown(chatID int64) {
	text := "❌ Неизвестная команда. Используйте /help"
	bot.Send(tgbotapi.NewMessage(chatID, text))
}

func handleCallback(callback *tgbotapi.CallbackQuery) {
	bot.Request(tgbotapi.NewCallback(callback.ID, ""))
}

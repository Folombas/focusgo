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

	// Настраиваем Menu кнопку
	setupMenuButton()

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

// setupMenuButton настраивает Menu кнопку
func setupMenuButton() {
	// Устанавливаем кнопку Menu с командами
	cmd := tgbotapi.NewSetMyCommands(
		tgbotapi.BotCommand{Command: "start", Description: "Начать игру"},
		tgbotapi.BotCommand{Command: "menu", Description: "Главное меню"},
		tgbotapi.BotCommand{Command: "profile", Description: "Профиль"},
		tgbotapi.BotCommand{Command: "skills", Description: "Навыки"},
		tgbotapi.BotCommand{Command: "quests", Description: "Квесты"},
		tgbotapi.BotCommand{Command: "stats", Description: "Статистика"},
		tgbotapi.BotCommand{Command: "backup", Description: "Бэкап БД"},
		tgbotapi.BotCommand{Command: "help", Description: "Справка"},
	)

	_, err := bot.Request(cmd)
	if err != nil {
		log.Printf("⚠️  Ошибка установки команд: %v", err)
	}

	log.Println("✅ Menu кнопка настроена")
}

func handleMessage(message *tgbotapi.Message) {
	chatID := message.Chat.ID

	// Обработка текста (не команд)
	if message.Text != "" && !message.IsCommand() {
		// Показываем главное меню
		showMainMenu(chatID)
		return
	}

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
	case "menu":
		showMainMenu(chatID)
	case "game":
		showGameMenu(chatID)
	case "profile":
		showProfile(chatID)
	case "skills":
		showSkills(chatID)
	case "quests":
		showQuests(chatID)
	case "stats":
		showStats(chatID)
	default:
		sendUnknown(chatID)
	}
}

func sendStart(chatID int64) {
	text := `🎮 <b>FOCUSGO — Temptation Simulator</b>

Добро пожаловать в симулятор борьбы с искушениями!

🎯 <b>Твоя цель:</b>
Сопротивляться искушениям, изучать Go и достичь уровня Go-Мастера!

📋 <b>Команды:</b>
/menu — Главное меню
/profile — Твой профиль
/skills — Дерево навыков
/quests — Ежедневные квесты
/stats — Статистика
/backup — Бэкап БД
/help — Справка

💪 <b>Помни:</b>
Каждая строка кода на Go — кирпичик в фундаменте твоей карьеры!`

	msg := tgbotapi.NewMessage(chatID, text)
	msg.ParseMode = "HTML"

	// Добавляем клавиатуру с кнопками
	keyboard := tgbotapi.ReplyKeyboardMarkup{
		ResizeKeyboard: true,
		Keyboard: [][]tgbotapi.KeyboardButton{
			{
				tgbotapi.NewKeyboardButton("🎮 Меню"),
			},
			{
				tgbotapi.NewKeyboardButton("👤 Профиль"),
				tgbotapi.NewKeyboardButton("🌳 Навыки"),
			},
			{
				tgbotapi.NewKeyboardButton("📋 Квесты"),
				tgbotapi.NewKeyboardButton("📊 Статистика"),
			},
		},
	}
	msg.ReplyMarkup = keyboard

	bot.Send(msg)
}

// showMainMenu показывает главное меню с inline-клавиатурой
func showMainMenu(chatID int64) {
	text := `🎮 <b>ГЛАВНОЕ МЕНЮ</b>

Выберите действие:`

	msg := tgbotapi.NewMessage(chatID, text)
	msg.ParseMode = "HTML"
	msg.ReplyMarkup = createMainMenuKeyboard()

	bot.Send(msg)
}

// createMainMenuKeyboard создаёт главную inline-клавиатуру
func createMainMenuKeyboard() tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🎮 Начать игру", "cb_start_game"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("📚 Учить Go (30 мин)", "cb_study_30"),
			tgbotapi.NewInlineKeyboardButtonData("📚 Учить Go (60 мин)", "cb_study_60"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("💤 Отдохнуть (15 мин)", "cb_rest_15"),
			tgbotapi.NewInlineKeyboardButtonData("💤 Отдохнуть (30 мин)", "cb_rest_30"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("📋 Квесты", "cb_quests"),
			tgbotapi.NewInlineKeyboardButtonData("🌳 Навыки", "cb_skills"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("📊 Статистика", "cb_stats"),
			tgbotapi.NewInlineKeyboardButtonData("👤 Профиль", "cb_profile"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("💾 Сохранить", "cb_save"),
			tgbotapi.NewInlineKeyboardButtonData("🌙 Завершить день", "cb_end_day"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🔙 Назад к командам", "cb_back"),
		),
	)
}

// showGameMenu показывает игровое меню
func showGameMenu(chatID int64) {
	text := `🎮 <b>ИГРОВОЕ МЕНЮ</b>

Выберите действие:`

	msg := tgbotapi.NewMessage(chatID, text)
	msg.ParseMode = "HTML"
	msg.ReplyMarkup = createGameMenuKeyboard()

	bot.Send(msg)
}

// createGameMenuKeyboard создаёт игровую inline-клавиатуру
func createGameMenuKeyboard() tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("⚔️  Финальная битва", "cb_battle"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("📚 Учить Go", "cb_study_30"),
			tgbotapi.NewInlineKeyboardButtonData("💤 Отдых", "cb_rest_15"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🔙 Главное меню", "cb_main_menu"),
		),
	)
}

// showProfile показывает профиль
func showProfile(chatID int64) {
	text := `👤 <b>ТВОЙ ПРОФИЛЬ</b>

Здесь будет твоя статистика...

⚠️  Функция в разработке!`

	msg := tgbotapi.NewMessage(chatID, text)
	msg.ParseMode = "HTML"
	msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🔙 Главное меню", "cb_main_menu"),
		),
	)

	bot.Send(msg)
}

// showSkills показывает навыки
func showSkills(chatID int64) {
	text := `🌳 <b>ДЕРЕВО НАВЫКОВ</b>

Здесь будет дерево навыков...

⚠️  Функция в разработке!`

	msg := tgbotapi.NewMessage(chatID, text)
	msg.ParseMode = "HTML"
	msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🔙 Главное меню", "cb_main_menu"),
		),
	)

	bot.Send(msg)
}

// showQuests показывает квесты
func showQuests(chatID int64) {
	text := `📋 <b>ЕЖЕДНЕВНЫЕ КВЕСТЫ</b>

Здесь будут твои квесты...

⚠️  Функция в разработке!`

	msg := tgbotapi.NewMessage(chatID, text)
	msg.ParseMode = "HTML"
	msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🔙 Главное меню", "cb_main_menu"),
		),
	)

	bot.Send(msg)
}

// showStats показывает статистику
func showStats(chatID int64) {
	text := `📊 <b>СТАТИСТИКА</b>

Здесь будет твоя статистика...

⚠️  Функция в разработке!`

	msg := tgbotapi.NewMessage(chatID, text)
	msg.ParseMode = "HTML"
	msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🔙 Главное меню", "cb_main_menu"),
		),
	)

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

<b>Команды:</b>
/start — Старт
/menu — Главное меню
/game — Игровое меню
/profile — Твой профиль
/skills — Дерево навыков
/quests — Ежедневные квесты
/stats — Статистика
/backup — Бэкап БД
/help — Справка

<b>Уведомления:</b>
• 9:00 — Напоминание о квестах
• 20:00 — Напоминание о финальной битве
• 22:00 — Напоминание о незавершённых квестах

<b>Menu кнопка:</b>
Нажми на кнопку "Menu" слева от поля ввода, чтобы выбрать команду!

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
	chatID := callback.Message.Chat.ID
	data := callback.Data

	// Отправляем подтверждение
	bot.Request(tgbotapi.NewCallback(callback.ID, ""))

	// Обработчики inline-кнопок
	switch data {
	case "cb_start_game":
		sendCallbackMessage(chatID, "🎮 Игра начинается! (в разработке)")
	case "cb_study_30":
		sendCallbackMessage(chatID, "📚 Изучение Go: 30 минут (в разработке)")
	case "cb_study_60":
		sendCallbackMessage(chatID, "📚 Изучение Go: 60 минут (в разработке)")
	case "cb_rest_15":
		sendCallbackMessage(chatID, "💤 Отдых: 15 минут (в разработке)")
	case "cb_rest_30":
		sendCallbackMessage(chatID, "💤 Отдых: 30 минут (в разработке)")
	case "cb_quests":
		showQuests(chatID)
	case "cb_skills":
		showSkills(chatID)
	case "cb_stats":
		showStats(chatID)
	case "cb_profile":
		showProfile(chatID)
	case "cb_save":
		handleBackup(chatID)
	case "cb_end_day":
		sendCallbackMessage(chatID, "🌙 День завершён! (в разработке)")
	case "cb_battle":
		sendCallbackMessage(chatID, "⚔️  Финальная битва! (в разработке)")
	case "cb_main_menu":
		showMainMenu(chatID)
	case "cb_back":
		sendCallbackMessage(chatID, "🔙 Возврат к командам")
	}
}

// sendCallbackMessage отправляет сообщение в ответ на callback
func sendCallbackMessage(chatID int64, text string) {
	msg := tgbotapi.NewMessage(chatID, text)
	msg.ParseMode = "HTML"
	bot.Send(msg)
}

package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var (
	bot     *tgbotapi.BotAPI
	players = make(map[int64]*Player) // Кэш игроков в памяти
)

func main() {
	// Инициализация базы данных
	if err := InitDB("focusgo.db"); err != nil {
		log.Fatalf("❌ Ошибка инициализации БД: %v", err)
	}
	defer CloseDB()

	// Обработка сигналов для корректного завершения
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		log.Println("🛑 Получен сигнал завершения, закрываем БД...")
		CloseDB()
		os.Exit(0)
	}()

	// Получаем токен из переменной окружения или используем заглушку
	token := os.Getenv("TELEGRAM_BOT_TOKEN")
	if token == "" {
		token = "YOUR_BOT_TOKEN_HERE"
		fmt.Println("⚠️  Установите переменную окружения TELEGRAM_BOT_TOKEN")
	}

	var err error
	bot, err = tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true
	log.Printf("✅ Бот авторизован: %s", bot.Self.UserName)

	// Создаем конфиг обновлений
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	// Получаем канал обновлений
	updates := bot.GetUpdatesChan(u)

	// Обрабатываем обновления
	for update := range updates {
		if update.Message != nil {
			handleMessage(update.Message)
		} else if update.CallbackQuery != nil {
			handleCallback(update.CallbackQuery)
		}
	}
}

// handleMessage обрабатывает входящие сообщения
func handleMessage(message *tgbotapi.Message) {
	chatID := message.Chat.ID

	// Команды
	if message.IsCommand() {
		handleCommand(message)
		return
	}

	// Обработка текста для активной игры
	if players[chatID] != nil && players[chatID].GameActive {
		handleGameInput(message)
	}
}

// handleCommand обрабатывает команды
func handleCommand(message *tgbotapi.Message) {
	chatID := message.Chat.ID
	command := message.Command()

	switch command {
	case "start":
		sendStartMessage(chatID)
	case "play":
		startGame(chatID, message.From.FirstName)
	case "stats":
		sendStats(chatID)
	case "save":
		saveGame(chatID)
	case "load":
		loadGame(chatID)
	case "help":
		sendHelp(chatID)
	case "profile":
		sendProfile(chatID)
	case "skills":
		sendSkills(chatID)
	case "quests":
		sendQuests(chatID)
	case "leaderboard":
		sendLeaderboard(chatID)
	default:
		sendUnknownCommand(chatID)
	}
}

// sendStartMessage отправляет приветственное сообщение
func sendStartMessage(chatID int64) {
	text := `🎮 <b>FOCUSGO — Temptation Simulator</b>

Добро пожаловать в симулятор борьбы с искушениями!

🎯 <b>Твоя цель:</b>
Сопротивляться искушениям, изучать Go и достичь уровня Go-Мастера!

📋 <b>Команды:</b>
/play — Начать игру
/profile — Твой профиль
/skills — Дерево навыков
/quests — Ежедневные квесты
/stats — Статистика
/leaderboard — Таблица лидеров
/save — Сохранить прогресс
/help — Помощь

💪 <b>Помни:</b>
Каждая строка кода на Go — кирпичик в фундаменте твоей карьеры!`

	msg := tgbotapi.NewMessage(chatID, text)
	msg.ParseMode = "HTML"

	// Добавляем кнопки
	msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🎮 Начать игру", "start_game"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("📊 Профиль", "show_profile"),
			tgbotapi.NewInlineKeyboardButtonData("🌳 Навыки", "show_skills"),
		),
	)

	bot.Send(msg)
}

// sendHelp отправляет справку
func sendHelp(chatID int64) {
	text := `📖 <b>СПРАВКА</b>

<b>Основные команды:</b>
/start — Приветственное сообщение
/play — Начать новую игру
/profile — Показать твой профиль
/skills — Дерево навыков
/quests — Ежедневные квесты
/stats — Расширенная статистика
/leaderboard — Таблица лидеров
/save — Сохранить прогресс
/load — Загрузить сохранение
/help — Эта справка

<b>Как играть:</b>
1. Начни игру командой /play
2. Выбирай действия из кнопок
3. Изучай Go и сопротивляйся искушениям
4. Выполняй ежедневные квесты
5. Улучшай навыки
6. Сохраняй прогресс

<b>Советы:</b>
• Изучение Go даёт опыт и очки навыков
• Сопротивление искушениям укрепляет силу воли
• Выполняй квесты для бонусов
• Не поддавайся прокрастинации!

🚀 Удачи в изучении Go!`

	msg := tgbotapi.NewMessage(chatID, text)
	msg.ParseMode = "HTML"
	bot.Send(msg)
}

// sendUnknownCommand отправляет сообщение о неизвестной команде
func sendUnknownCommand(chatID int64) {
	text := "❌ Неизвестная команда. Используйте /help для справки."
	msg := tgbotapi.NewMessage(chatID, text)
	bot.Send(msg)
}

// handleCallback обрабатывает нажатия на inline-кнопки
func handleCallback(callback *tgbotapi.CallbackQuery) {
	chatID := callback.Message.Chat.ID
	data := callback.Data

	// Отправляем подтверждение
	bot.Request(tgbotapi.NewCallback(callback.ID, ""))

	switch data {
	case "start_game":
		startGame(chatID, callback.From.FirstName)
	case "show_profile":
		sendProfile(chatID)
	case "show_skills":
		sendSkills(chatID)
	case "show_quests":
		sendQuests(chatID)
	case "show_stats":
		sendStats(chatID)
	case "study_go_30":
		handleStudyGo30(chatID)
	case "study_go_60":
		handleStudyGo60(chatID)
	case "rest_15":
		handleRest15(chatID)
	case "rest_30":
		handleRest30(chatID)
	case "check_quests":
		sendQuests(chatID)
	case "upgrade_skill":
		handleUpgradeSkill(chatID)
	case "save_game":
		saveGame(chatID)
	case "end_day":
		handleEndDay(chatID)
	case "main_menu":
		showMainMenu(chatID)
	case "final_battle_fight":
		handleFinalBattle(chatID)
	}

	// Обработка улучшения навыков
	if len(data) > 8 && data[:8] == "upgrade_" {
		skillID := data[8:]
		handleSkillUpgrade(chatID, skillID)
	}
}

// handleGameInput обрабатывает ввод во время игры
func handleGameInput(message *tgbotapi.Message) {
	chatID := message.Chat.ID
	player := players[chatID]

	if player == nil {
		return
	}

	// Обработка числового ввода (например, минуты изучения)
	// Здесь можно добавить парсинг ввода пользователя
}

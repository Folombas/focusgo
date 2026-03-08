package bot

import (
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"focusgo/internal/database"
	"focusgo/internal/game"
	"focusgo/internal/models"
	"focusgo/internal/notifications"
)

var (
	players  = make(map[int64]*models.Player)
	telegramBot *tgbotapi.BotAPI
)

// SetBot устанавливает бота для использования в пакете
func SetBot(bot *tgbotapi.BotAPI) {
	telegramBot = bot
	notifications.SetBot(bot)
}

// HandleMessage обрабатывает входящие сообщения
func HandleMessage(bot *tgbotapi.BotAPI, message *tgbotapi.Message) {
	chatID := message.Chat.ID

	// Команды
	if message.IsCommand() {
		handleCommand(bot, message)
		return
	}

	// Обработка текста для активной игры
	if players[chatID] != nil && players[chatID].GameActive {
		handleGameInput(bot, message)
	}
}

// HandleCallback обрабатывает нажатия на inline-кнопки
func HandleCallback(bot *tgbotapi.BotAPI, callback *tgbotapi.CallbackQuery) {
	chatID := callback.Message.Chat.ID
	data := callback.Data

	// Отправляем подтверждение
	bot.Request(tgbotapi.NewCallback(callback.ID, ""))

	switch data {
	case "start_game":
		game.StartGame(bot, chatID, callback.From.FirstName, players)
	case "show_profile":
		game.SendProfile(bot, chatID, players)
	case "show_skills":
		game.SendSkills(bot, chatID, players)
	case "show_quests":
		game.SendQuests(bot, chatID, players)
	case "show_stats":
		game.SendStats(bot, chatID, players)
	case "study_go_30":
		game.HandleStudyGo30(bot, chatID, players)
	case "study_go_60":
		game.HandleStudyGo60(bot, chatID, players)
	case "rest_15":
		game.HandleRest15(bot, chatID, players)
	case "rest_30":
		game.HandleRest30(bot, chatID, players)
	case "check_quests":
		game.SendQuests(bot, chatID, players)
	case "upgrade_skill":
		game.HandleUpgradeSkill(bot, chatID, players)
	case "save_game":
		game.SaveGame(bot, chatID, players)
	case "end_day":
		game.HandleEndDay(bot, chatID, players)
	case "main_menu":
		game.ShowMainMenu(bot, chatID)
	case "final_battle_fight":
		game.HandleFinalBattle(bot, chatID, players)
	case "notif_toggle_all":
		handleToggleAllNotifications(bot, chatID)
	case "notif_toggle_quests":
		handleToggleQuestsNotifications(bot, chatID)
	case "notif_toggle_battle":
		handleToggleBattleNotifications(bot, chatID)
	case "notif_toggle_unfinished":
		handleToggleUnfinishedNotifications(bot, chatID)
	}

	// Обработка улучшения навыков
	if len(data) > 8 && data[:8] == "upgrade_" {
		skillID := data[8:]
		game.HandleSkillUpgrade(bot, chatID, skillID, players)
	}
}

// handleCommand обрабатывает команды
func handleCommand(bot *tgbotapi.BotAPI, message *tgbotapi.Message) {
	chatID := message.Chat.ID
	command := message.Command()

	switch command {
	case "start":
		sendStartMessage(bot, chatID)
	case "play":
		game.StartGame(bot, chatID, message.From.FirstName, players)
	case "stats":
		game.SendStats(bot, chatID, players)
	case "save":
		game.SaveGame(bot, chatID, players)
	case "load":
		game.LoadGame(bot, chatID, players)
	case "help":
		sendHelp(bot, chatID)
	case "profile":
		game.SendProfile(bot, chatID, players)
	case "skills":
		game.SendSkills(bot, chatID, players)
	case "quests":
		game.SendQuests(bot, chatID, players)
	case "leaderboard":
		game.SendLeaderboard(bot, chatID)
	case "remind":
		game.SendRemindSettings(bot, chatID, players)
	case "backup":
		handleBackup(bot, chatID)
	default:
		sendUnknownCommand(bot, chatID)
	}
}

// sendStartMessage отправляет приветственное сообщение
func sendStartMessage(bot *tgbotapi.BotAPI, chatID int64) {
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
/remind — Настройки уведомлений
/backup — Бэкап прогресса
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
func sendHelp(bot *tgbotapi.BotAPI, chatID int64) {
	text := `📖 <b>СПРАВКА</b>

<b>Основные команды:</b>
/start — Приветственное сообщение
/play — Начать новую игру
/profile — Показать твой профиль
/skills — Дерево навыков
/quests — Ежедневные квесты
/stats — Расширенная статистика
/leaderboard — Таблица лидеров
/remind — Настройки уведомлений
/backup — Бэкап прогресса
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

<b>Уведомления:</b>
• 9:00 — Напоминание о квестах
• 20:00 — Напоминание о финальной битве
• 22:00 — Напоминание о незавершённых квестах

Настрой уведомления командой /remind

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
func sendUnknownCommand(bot *tgbotapi.BotAPI, chatID int64) {
	text := "❌ Неизвестная команда. Используйте /help для справки."
	msg := tgbotapi.NewMessage(chatID, text)
	bot.Send(msg)
}

// handleGameInput обрабатывает ввод во время игры
func handleGameInput(bot *tgbotapi.BotAPI, message *tgbotapi.Message) {
	// Обработка числового ввода (например, минуты изучения)
	// Здесь можно добавить парсинг ввода пользователя
}

// handleToggleAllNotifications переключает все уведомления
func handleToggleAllNotifications(bot *tgbotapi.BotAPI, chatID int64) {
	settings := notifications.GetSettings(chatID)
	settings.Enabled = !settings.Enabled
	settings.DailyQuestsEnabled = settings.Enabled
	settings.FinalBattleEnabled = settings.Enabled
	settings.UnfinishedEnabled = settings.Enabled

	notifications.SaveSettings(settings)

	text := fmt.Sprintf("🔔 Уведомления %s", map[bool]string{true: "включены", false: "выключены"}[settings.Enabled])
	msg := tgbotapi.NewMessage(chatID, text)
	msg.ParseMode = "HTML"
	bot.Send(msg)
	game.SendRemindSettings(bot, chatID, players)
}

// handleToggleQuestsNotifications переключает уведомления о квестах
func handleToggleQuestsNotifications(bot *tgbotapi.BotAPI, chatID int64) {
	settings := notifications.GetSettings(chatID)
	settings.DailyQuestsEnabled = !settings.DailyQuestsEnabled

	notifications.SaveSettings(settings)

	text := fmt.Sprintf("📋 Уведомления о квестах %s", map[bool]string{true: "включены", false: "выключены"}[settings.DailyQuestsEnabled])
	msg := tgbotapi.NewMessage(chatID, text)
	msg.ParseMode = "HTML"
	bot.Send(msg)
	game.SendRemindSettings(bot, chatID, players)
}

// handleToggleBattleNotifications переключает уведомления о битве
func handleToggleBattleNotifications(bot *tgbotapi.BotAPI, chatID int64) {
	settings := notifications.GetSettings(chatID)
	settings.FinalBattleEnabled = !settings.FinalBattleEnabled

	notifications.SaveSettings(settings)

	text := fmt.Sprintf("⚔️  Уведомления о битве %s", map[bool]string{true: "включены", false: "выключены"}[settings.FinalBattleEnabled])
	msg := tgbotapi.NewMessage(chatID, text)
	msg.ParseMode = "HTML"
	bot.Send(msg)
	game.SendRemindSettings(bot, chatID, players)
}

// handleToggleUnfinishedNotifications переключает уведомления о незавершённых квестах
func handleToggleUnfinishedNotifications(bot *tgbotapi.BotAPI, chatID int64) {
	settings := notifications.GetSettings(chatID)
	settings.UnfinishedEnabled = !settings.UnfinishedEnabled

	notifications.SaveSettings(settings)

	text := fmt.Sprintf("⏰ Уведомления о незавершённых квестах %s", map[bool]string{true: "включены", false: "выключены"}[settings.UnfinishedEnabled])
	msg := tgbotapi.NewMessage(chatID, text)
	msg.ParseMode = "HTML"
	bot.Send(msg)
	game.SendRemindSettings(bot, chatID, players)
}

// handleBackup обрабатывает команду бэкапа
func handleBackup(bot *tgbotapi.BotAPI, chatID int64) {
	backupPath, err := database.CreateBackup()
	if err != nil {
		text := fmt.Sprintf("❌ <b>ОШИБКА БЭКАПА!</b>\n\n%s", err.Error())
		msg := tgbotapi.NewMessage(chatID, text)
		msg.ParseMode = "HTML"
		bot.Send(msg)
		return
	}

	text := fmt.Sprintf("💾 <b>БЭКАП СОЗДАН!</b>\n\nФайл: %s", backupPath)
	msg := tgbotapi.NewMessage(chatID, text)
	msg.ParseMode = "HTML"
	bot.Send(msg)
}

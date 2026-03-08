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
	"focusgo/internal/game"
)

var (
	bot          *tgbotapi.BotAPI
	gameStates   = make(map[int64]*game.GameState)
	questSystems = make(map[int64]*game.QuestSystem)
	skillTrees   = make(map[int64]*game.SkillTree)
	achievementSystems = make(map[int64]*game.AchievementSystem)
)

func main() {
	// Инициализация БД
	if err := database.InitDB("focusgo.db"); err != nil {
		log.Fatalf("❌ Ошибка инициализации БД: %v", err)
	}
	defer database.CloseDB()

	// Автоматические бэкапы
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

	// Настройка Menu кнопки
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

	if message.Text != "" && !message.IsCommand() {
		if message.Text == "🎮 Меню" {
			showMainMenu(chatID)
			return
		}
		if message.Text == "👤 Профиль" {
			showProfile(chatID)
			return
		}
		if message.Text == "🌳 Навыки" {
			showSkills(chatID)
			return
		}
		if message.Text == "📋 Квесты" {
			showQuests(chatID)
			return
		}
		if message.Text == "📊 Статистика" {
			showStats(chatID)
			return
		}
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
		sendStart(chatID, message.From.FirstName)
	case "menu":
		showMainMenu(chatID)
	case "profile":
		showProfile(chatID)
	case "skills":
		showSkills(chatID)
	case "quests":
		showQuests(chatID)
	case "stats":
		showStats(chatID)
	case "leaderboard":
		showLeaderboard(chatID)
	case "achievements":
		showAchievements(chatID)
	case "backup":
		handleBackup(chatID)
	case "help":
		sendHelp(chatID)
	default:
		sendUnknown(chatID)
	}
}

func sendStart(chatID int64, name string) {
	// Загружаем или создаём игру
	state, _ := game.LoadGameState(chatID)
	if state == nil {
		state = game.NewGameState(chatID, name)
		state.SaveGameState()
		
		// Создаём квесты
		qs := game.NewQuestSystem(chatID)
		qs.GenerateDailyQuests()
		questSystems[chatID] = qs
		
		// Создаём дерево навыков
		tree := game.NewSkillTree(chatID)
		tree.SaveSkillTree()
		skillTrees[chatID] = tree
		
		// Применяем бонусы
		state.ApplySkillBonuses(tree)
	} else {
		gameStates[chatID] = state
		// Загружаем квесты
		qs := game.NewQuestSystem(chatID)
		qs.GenerateDailyQuests()
		questSystems[chatID] = qs
		
		// Загружаем дерево навыков и применяем бонусы
		tree, _ := game.LoadSkillTree(chatID)
		if tree != nil {
			skillTrees[chatID] = tree
			state.ApplySkillBonuses(tree)
		}
	}

	text := fmt.Sprintf(`🎮 <b>FOCUSGO — Temptation Simulator</b>

👋 Привет, %s!

🎯 <b>Твоя цель:</b>
Сопротивляться искушениям, изучать Go и достичь уровня Go-Мастера!

%s

💪 <b>Помни:</b>
Каждая строка кода на Go — кирпичик в фундаменте твоей карьеры!

Используй /menu или кнопку "🎮 Меню" для игры!`,
		name, state.GetStatus())

	msg := tgbotapi.NewMessage(chatID, text)
	msg.ParseMode = "HTML"

	keyboard := tgbotapi.ReplyKeyboardMarkup{
		ResizeKeyboard: true,
		Keyboard: [][]tgbotapi.KeyboardButton{
			{tgbotapi.NewKeyboardButton("🎮 Меню")},
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

func showMainMenu(chatID int64) {
	// Загружаем состояние
	state, _ := game.LoadGameState(chatID)
	if state == nil {
		text := "⚠️  <b>ИГРА НЕ НАЧАТА</b>\n\nИспользуйте /start для начала игры!"
		bot.Send(tgbotapi.NewMessage(chatID, text))
		return
	}

	gameStates[chatID] = state

	text := `🎮 <b>ГЛАВНОЕ МЕНЮ</b>

Выберите действие:`

	msg := tgbotapi.NewMessage(chatID, text)
	msg.ParseMode = "HTML"
	msg.ReplyMarkup = createMainMenuKeyboard()

	bot.Send(msg)
}

func createMainMenuKeyboard() tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
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
	)
}

func handleCallback(callback *tgbotapi.CallbackQuery) {
	chatID := callback.Message.Chat.ID
	data := callback.Data

	bot.Request(tgbotapi.NewCallback(callback.ID, ""))

	// Загружаем состояние
	state, _ := game.LoadGameState(chatID)
	if state == nil && data != "cb_start_game" {
		text := "⚠️  <b>ИГРА НЕ НАЧАТА</b>\n\nИспользуйте /start!"
		bot.Send(tgbotapi.NewMessage(chatID, text))
		return
	}

	switch data {
	case "cb_study_30":
		msg, xp, knowledge := state.StudyGo(30)
		state.SaveGameState()

		// Проверяем квесты
		qs := questSystems[chatID]
		if qs != nil {
			completed, reward := qs.UpdateProgress("study_go_30min", 30)
			if completed {
				bot.Send(tgbotapi.NewMessage(chatID, fmt.Sprintf("✅ Квест выполнен! +%d очков навыков", reward)))
				// Добавляем очки навыков в дерево
				tree, _ := game.LoadSkillTree(chatID)
				if tree != nil {
					tree.EarnSkillPoints(reward)
					tree.SaveSkillTree()
					skillTrees[chatID] = tree
				}
			}
		}

		// Проверяем искушение
		if state.CheckTemptation() {
			temptMsg, _ := state.ResistTemptation(70)
			bot.Send(tgbotapi.NewMessage(chatID, temptMsg))
		}
		
		response := fmt.Sprintf("%s\n\n✨ +%d опыта\n🧠 +%d знаний", msg, xp, knowledge)
		bot.Send(tgbotapi.NewMessage(chatID, response))
		
	case "cb_study_60":
		msg, xp, knowledge := state.StudyGo(60)
		state.SaveGameState()
		
		qs := questSystems[chatID]
		if qs != nil {
			completed, reward := qs.UpdateProgress("study_go_30min", 60)
			if completed {
				bot.Send(tgbotapi.NewMessage(chatID, fmt.Sprintf("✅ Квест выполнен! +%d очков навыков", reward)))
			}
		}
		
		response := fmt.Sprintf("%s\n\n✨ +%d опыта\n🧠 +%d знаний", msg, xp, knowledge)
		bot.Send(tgbotapi.NewMessage(chatID, response))
		
	case "cb_rest_15":
		msg := state.Rest(15)
		state.SaveGameState()
		bot.Send(tgbotapi.NewMessage(chatID, msg))
		
	case "cb_rest_30":
		msg := state.Rest(30)
		state.SaveGameState()
		bot.Send(tgbotapi.NewMessage(chatID, msg))
		
	case "cb_quests":
		showQuests(chatID)
		
	case "cb_skills":
		showSkills(chatID)

	case "cb_stats":
		showStats(chatID)

	case "cb_profile":
		showProfile(chatID)

	case "cb_upgrade_go_basics", "cb_upgrade_concurrency", "cb_upgrade_interfaces",
	     "cb_upgrade_web_frameworks", "cb_upgrade_database", "cb_upgrade_microservices",
	     "cb_upgrade_focus_master", "cb_upgrade_meditation", "cb_upgrade_anti_procrastination",
	     "cb_upgrade_willpower", "cb_upgrade_discipline", "cb_upgrade_money_management":
		// Извлекаем ID навыка из callback_data
		skillID := data[11:] // убираем "cb_upgrade_"
		handleUpgradeSkill(chatID, skillID)
		
	case "cb_save":
		if state != nil {
			state.SaveGameState()
			bot.Send(tgbotapi.NewMessage(chatID, "💾 Прогресс сохранён!"))
		}
		
	case "cb_end_day":
		if state != nil {
			// Финальная битва
			bosses := []struct{ name string; power int }{
				{"👹 CAPCUT МОНТЁР", 95},
				{"👹 ИГРОВОЙ ЗАВИСИМОН", 92},
				{"👹 СОЦСЕТЕЙ ДЕМОНИУС", 88},
				{"👹 АЛКОГОЛЬНЫЙ ПРИЗРАК", 90},
				{"👹 ДЕПРЕССИЯ МАКСИМА", 98},
			}
			boss := bosses[time.Now().Unix()%int64(len(bosses))]

			won, battleMsg := state.FinalBattle(boss.name, boss.power)
			_ = won // используем переменную
			state.SaveGameState()

			// Проверяем квесты
			qs := questSystems[chatID]
			if qs != nil {
				completed := qs.GetCompletedCount() == 5
				qs.CheckDayStreak(completed)
				reward := qs.ClaimRewards()
				if reward > 0 {
					bot.Send(tgbotapi.NewMessage(chatID, fmt.Sprintf("💰 Получено %d очков навыков за квесты!", reward)))
				}
				
				// Проверяем достижения
				as, _ := game.LoadAchievementSystem(chatID)
				if as != nil {
					tree := skillTrees[chatID]
					unlocked := as.CheckAchievements(state, tree, qs)
					if len(unlocked) > 0 {
						as.SaveAchievementSystem()
						achievementSystems[chatID] = as
						for _, achievement := range unlocked {
							bot.Send(tgbotapi.NewMessage(chatID, fmt.Sprintf("🏆 <b>НОВОЕ ДОСТИЖЕНИЕ!</b>\n\n%s", achievement)))
						}
					}
				}
			}

			bot.Send(tgbotapi.NewMessage(chatID, battleMsg))
		}
		
	case "cb_main_menu":
		showMainMenu(chatID)
		
	case "cb_back":
		bot.Send(tgbotapi.NewMessage(chatID, "🔙 Возврат к командам. Используйте /menu"))
	}
}

func showProfile(chatID int64) {
	state, err := game.LoadGameState(chatID)
	if err != nil || state == nil {
		text := "⚠️  <b>ИГРА НЕ НАЧАТА</b>\n\nИспользуйте /start!"
		bot.Send(tgbotapi.NewMessage(chatID, text))
		return
	}

	// Загружаем дерево навыков для применения бонусов
	tree, _ := game.LoadSkillTree(chatID)
	if tree != nil {
		state.ApplySkillBonuses(tree)
		skillTrees[chatID] = tree
	}

	// Формируем строку с бонусами
	bonusText := ""
	if state.SkillBonuses["focus"] > 0 {
		bonusText += fmt.Sprintf("🎯 Фокус: +%d\n", state.SkillBonuses["focus"])
	}
	if state.SkillBonuses["willpower"] > 0 {
		bonusText += fmt.Sprintf("💪 Сила воли: +%d\n", state.SkillBonuses["willpower"])
	}
	if state.SkillBonuses["knowledge"] > 0 {
		bonusText += fmt.Sprintf("📚 Знание Go: +%d\n", state.SkillBonuses["knowledge"])
	}
	if state.SkillBonuses["money"] > 0 {
		bonusText += fmt.Sprintf("💰 Деньги: +%d\n", state.SkillBonuses["money"])
	}
	if state.SkillBonuses["dopamine"] > 0 {
		bonusText += fmt.Sprintf("✨ Дофамин: +%d\n", state.SkillBonuses["dopamine"])
	}

	if bonusText != "" {
		bonusText = "\n🌳 <b>БОНУСЫ ОТ НАВЫКОВ:</b>\n" + bonusText
	}

	text := fmt.Sprintf(`👤 <b>ПРОФИЛЬ ИГРОКА</b>
━━━━━━━━━━━━━━━━━━━━

%s

🏅 Рейтинг: %s%s`, state.GetStatus(), state.GetRating(), bonusText)

	msg := tgbotapi.NewMessage(chatID, text)
	msg.ParseMode = "HTML"
	msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🔙 Главное меню", "cb_main_menu"),
		),
	)

	bot.Send(msg)
}

func showSkills(chatID int64) {
	// Загружаем дерево навыков
	tree, err := game.LoadSkillTree(chatID)
	if err != nil {
		log.Printf("Ошибка загрузки дерева навыков: %v", err)
	}
	
	if tree == nil {
		tree = game.NewSkillTree(chatID)
	}
	
	skillTrees[chatID] = tree

	text := tree.Display()

	msg := tgbotapi.NewMessage(chatID, text)
	msg.ParseMode = "HTML"
	
	// Создаём клавиатуру с доступными улучшениями
	keyboard := createSkillsKeyboard(tree)
	if len(keyboard) > 0 {
		msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(keyboard...)
	} else {
		msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("🔙 Главное меню", "cb_main_menu"),
			),
		)
	}

	bot.Send(msg)
}

// createSkillsKeyboard создаёт клавиатуру для улучшения навыков
func createSkillsKeyboard(tree *game.SkillTree) [][]tgbotapi.InlineKeyboardButton {
	var keyboard [][]tgbotapi.InlineKeyboardButton

	if tree.SkillPoints > 0 {
		for id, skill := range tree.Skills {
			if skill.Unlocked && skill.Level < skill.MaxLevel {
				cost := skill.CostPerLevel
				if tree.SkillPoints >= cost {
					keyboard = append(keyboard, tgbotapi.NewInlineKeyboardRow(
						tgbotapi.NewInlineKeyboardButtonData(
							fmt.Sprintf("⬆️ %s (%d очк.)", skill.Name, cost),
							fmt.Sprintf("cb_upgrade_%s", id),
						),
					))
				}
			}
		}
	}

	keyboard = append(keyboard, tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("🔙 Главное меню", "cb_main_menu"),
	))

	return keyboard
}

func showQuests(chatID int64) {
	qs, exists := questSystems[chatID]
	if !exists {
		qs = game.NewQuestSystem(chatID)
		qs.GenerateDailyQuests()
		questSystems[chatID] = qs
	}

	text := qs.Display()

	msg := tgbotapi.NewMessage(chatID, text)
	msg.ParseMode = "HTML"
	msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🔙 Главное меню", "cb_main_menu"),
		),
	)

	bot.Send(msg)
}

// showLeaderboard показывает таблицу лидеров
func showLeaderboard(chatID int64) {
	leaderboard, err := database.GetLeaderboard(10)
	if err != nil {
		bot.Send(tgbotapi.NewMessage(chatID, "❌ Ошибка загрузки таблицы лидеров"))
		return
	}

	text := "🏆 <b>ТАБЛИЦА ЛИДЕРОВ</b>\n━━━━━━━━━━━━━━━━━━━━\n\nТоп-10 игроков:\n\n"

	for i, entry := range leaderboard {
		rank := ""
		switch i {
		case 0:
			rank = "🥇"
		case 1:
			rank = "🥈"
		case 2:
			rank = "🥉"
		default:
			rank = fmt.Sprintf("%d.", i+1)
		}

		text += fmt.Sprintf("%s <b>%s</b> — Ур.%d | Рейтинг: %d\n",
			rank, entry["name"], entry["level"], entry["rating"])
	}

	totalPlayers, _ := database.GetTotalPlayers()
	text += fmt.Sprintf("\n📊 Всего игроков: %d", totalPlayers)

	msg := tgbotapi.NewMessage(chatID, text)
	msg.ParseMode = "HTML"
	msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🔙 Главное меню", "cb_main_menu"),
		),
	)

	bot.Send(msg)
}

// showAchievements показывает достижения игрока
func showAchievements(chatID int64) {
	// Загружаем систему достижений
	as, err := game.LoadAchievementSystem(chatID)
	if err != nil {
		bot.Send(tgbotapi.NewMessage(chatID, "❌ Ошибка загрузки достижений"))
		return
	}

	achievementSystems[chatID] = as

	text := as.Display()

	msg := tgbotapi.NewMessage(chatID, text)
	msg.ParseMode = "HTML"
	msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🔙 Главное меню", "cb_main_menu"),
		),
	)

	bot.Send(msg)
}

func showStats(chatID int64) {
	state, err := game.LoadGameState(chatID)
	if err != nil || state == nil {
		text := "⚠️  <b>ИГРА НЕ НАЧАТА</b>\n\nИспользуйте /start!"
		bot.Send(tgbotapi.NewMessage(chatID, text))
		return
	}

	text := fmt.Sprintf(`📊 <b>СТАТИСТИКА</b>
━━━━━━━━━━━━━━━━━━━━

🏆 Уровень: %d
⭐ Опыт: %d/%d
🕐 Время в игре: %d минут
📅 Дней сыграно: %d
🏅 Рейтинг: %s

🚀 Продолжай учиться и достигнешь цели!`,
		state.Level, state.Experience, state.NextLevelXP,
		state.PlayTime, state.DaysPlayed, state.GetRating())

	msg := tgbotapi.NewMessage(chatID, text)
	msg.ParseMode = "HTML"
	msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🔙 Главное меню", "cb_main_menu"),
		),
	)

	bot.Send(msg)
}

// handleUpgradeSkill обрабатывает улучшение навыка
func handleUpgradeSkill(chatID int64, skillID string) {
	tree, exists := skillTrees[chatID]
	if !exists {
		var err error
		tree, err = game.LoadSkillTree(chatID)
		if err != nil {
			bot.Send(tgbotapi.NewMessage(chatID, "❌ Ошибка загрузки дерева навыков"))
			return
		}
		skillTrees[chatID] = tree
	}

	success, msg := tree.UpgradeSkill(skillID)
	
	if success {
		tree.SaveSkillTree()
		
		// Применяем бонусы к игроку
		state, _ := game.LoadGameState(chatID)
		if state != nil {
			state.ApplySkillBonuses(tree)
			state.SaveGameState()
			gameStates[chatID] = state
		}
		
		bot.Send(tgbotapi.NewMessage(chatID, msg))
		// Показываем обновлённое дерево
		showSkills(chatID)
	} else {
		bot.Send(tgbotapi.NewMessage(chatID, msg))
	}
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
/start — Начать игру
/menu — Главное меню
/profile — Твой профиль
/skills — Дерево навыков
/quests — Ежедневные квесты
/stats — Статистика
/leaderboard — Таблица лидеров
/achievements — Достижения
/backup — Бэкап БД
/help — Справка

<b>Как играть:</b>
1. Начни игру командой /start
2. Изучай Go и получай опыт
3. Выполняй ежедневные квесты
4. Сопротивляйся искушениям
5. Заверши день победой над боссом
6. Улучшай навыки и получай бонусы
7. Соревнуйся с другими в таблице лидеров
8. Разблокируй достижения!

<b>Уведомления:</b>
• 9:00 — Напоминание о квестах
• 20:00 — Напоминание о финальной битве
• 22:00 — Напоминание о незавершённых квестах

Автосохранение: каждые 24 часа`

	msg := tgbotapi.NewMessage(chatID, text)
	msg.ParseMode = "HTML"
	bot.Send(msg)
}

func sendUnknown(chatID int64) {
	text := "❌ Неизвестная команда. Используйте /help"
	bot.Send(tgbotapi.NewMessage(chatID, text))
}

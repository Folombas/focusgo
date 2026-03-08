package main

import (
	"fmt"
	"log"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// startGame начинает новую игру
func startGame(chatID int64, name string) {
	// Сначала пробуем загрузить из БД
	existingPlayer, err := LoadPlayer(chatID)
	if err != nil {
		log.Printf("Ошибка загрузки игрока: %v", err)
	}

	if existingPlayer != nil {
		// Игрок найден, загружаем его
		players[chatID] = existingPlayer
		text := fmt.Sprintf(`🎮 <b>ИГРА ЗАГРУЖЕНА!</b>

👋 С возвращением, %s!

Твой прогресс загружен из базы данных.
Продолжаем путь к становлению Go-разработчика!

%s`,
			name,
			existingPlayer.DisplayStatus(),
		)

		msg := tgbotapi.NewMessage(chatID, text)
		msg.ParseMode = "HTML"
		msg.ReplyMarkup = createGameKeyboard(existingPlayer)

		bot.Send(msg)
		return
	}

	// Создаем нового игрока
	player := NewPlayer(chatID, name)
	players[chatID] = player

	// Сохраняем в БД
	if err := SavePlayer(player); err != nil {
		log.Printf("Ошибка сохранения игрока: %v", err)
	}

	text := fmt.Sprintf(`🎮 <b>НОВАЯ ИГРА НАЧАТА!</b>

👋 Привет, %s!

Ты начинаешь свой путь к становлению Go-разработчика.
Твоя цель — сопротивляться искушениям и изучить Go!

%s

🎯 <b>Первый шаг:</b>
Изучи Go в течение 30 минут, чтобы выполнить первый квест!`,
		name,
		player.DisplayStatus(),
	)

	msg := tgbotapi.NewMessage(chatID, text)
	msg.ParseMode = "HTML"
	msg.ReplyMarkup = createGameKeyboard(player)

	bot.Send(msg)

	// Показываем ежедневные квесты
	player.Quests.GenerateDailyQuests()
	sendQuests(chatID)
}

// createGameKeyboard создает игровую клавиатуру
func createGameKeyboard(player *Player) tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("📚 Учить Go (30 мин)", "study_go_30"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("📚 Учить Go (60 мин)", "study_go_60"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("💤 Отдохнуть (15 мин)", "rest_15"),
			tgbotapi.NewInlineKeyboardButtonData("💤 Отдохнуть (30 мин)", "rest_30"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("📋 Квесты", "check_quests"),
			tgbotapi.NewInlineKeyboardButtonData("🌳 Навыки", "upgrade_skill"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("📊 Статистика", "show_stats"),
			tgbotapi.NewInlineKeyboardButtonData("💾 Сохранить", "save_game"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🌙 Завершить день", "end_day"),
		),
	)
}

// handleStudyGo30 обрабатывает изучение Go 30 минут
func handleStudyGo30(chatID int64) {
	player := players[chatID]
	if player == nil {
		sendNoGameMessage(chatID)
		return
	}

	result := player.StudyGo(30)
	player.PlayTime += 30
	player.Hour += 1

	if player.Hour >= 24 {
		player.Hour = 8
		player.CurrentDay++
	}

	// Сохраняем в БД после каждого действия
	if err := SavePlayer(player); err != nil {
		log.Printf("Ошибка сохранения: %v", err)
	}

	msg := tgbotapi.NewMessage(chatID, result)
	msg.ParseMode = "HTML"
	msg.ReplyMarkup = createGameKeyboard(player)

	bot.Send(msg)
	checkRandomEvents(chatID, player)
	updatePlayerStatus(chatID, player)
}

// handleStudyGo60 обрабатывает изучение Go 60 минут
func handleStudyGo60(chatID int64) {
	player := players[chatID]
	if player == nil {
		sendNoGameMessage(chatID)
		return
	}

	result := player.StudyGo(60)
	player.PlayTime += 60
	player.Hour += 2

	if player.Hour >= 24 {
		player.Hour = 8
		player.CurrentDay++
	}

	// Сохраняем в БД
	if err := SavePlayer(player); err != nil {
		log.Printf("Ошибка сохранения: %v", err)
	}

	msg := tgbotapi.NewMessage(chatID, result)
	msg.ParseMode = "HTML"
	msg.ReplyMarkup = createGameKeyboard(player)

	bot.Send(msg)
	checkRandomEvents(chatID, player)
	checkRandomEvents(chatID, player)
	updatePlayerStatus(chatID, player)
}

// handleRest15 обрабатывает отдых 15 минут
func handleRest15(chatID int64) {
	player := players[chatID]
	if player == nil {
		sendNoGameMessage(chatID)
		return
	}

	result := player.Rest(15)
	player.PlayTime += 15
	player.Hour += 1

	if player.Hour >= 24 {
		player.Hour = 8
		player.CurrentDay++
	}

	// Сохраняем в БД
	if err := SavePlayer(player); err != nil {
		log.Printf("Ошибка сохранения: %v", err)
	}

	msg := tgbotapi.NewMessage(chatID, result)
	msg.ParseMode = "HTML"
	msg.ReplyMarkup = createGameKeyboard(player)

	bot.Send(msg)
}

// handleRest30 обрабатывает отдых 30 минут
func handleRest30(chatID int64) {
	player := players[chatID]
	if player == nil {
		sendNoGameMessage(chatID)
		return
	}

	result := player.Rest(30)
	player.PlayTime += 30
	player.Hour += 2

	if player.Hour >= 24 {
		player.Hour = 8
		player.CurrentDay++
	}

	// Сохраняем в БД
	if err := SavePlayer(player); err != nil {
		log.Printf("Ошибка сохранения: %v", err)
	}

	msg := tgbotapi.NewMessage(chatID, result)
	msg.ParseMode = "HTML"
	msg.ReplyMarkup = createGameKeyboard(player)

	bot.Send(msg)
}

// handleSkillUpgrade обрабатывает улучшение навыка
func handleSkillUpgrade(chatID int64, skillID string) {
	player := players[chatID]
	if player == nil {
		sendNoGameMessage(chatID)
		return
	}

	skill := player.SkillTree.Skills[skillID]
	if skill == nil {
		return
	}

	if player.SkillTree.UpgradeSkill(skillID) {
		player.ApplySkillBonuses()

		text := fmt.Sprintf(`🎉 <b>НАВЫК УЛУЧШЕН!</b>

%s повышен до уровня %d!

+%d к %s`,
			skill.Name,
			skill.Level,
			skill.BonusValue,
			translateBonusType(skill.BonusType))

		// Сохраняем в БД
		if err := SavePlayer(player); err != nil {
			log.Printf("Ошибка сохранения: %v", err)
		}

		msg := tgbotapi.NewMessage(chatID, text)
		msg.ParseMode = "HTML"
		msg.ReplyMarkup = createGameKeyboard(player)

		bot.Send(msg)
		sendSkills(chatID)
	} else {
		text := `❌ <b>НЕ УДАЛОСЬ УЛУЧШИТЬ</b>

Недостаточно очков навыков или навык заблокирован.`

		msg := tgbotapi.NewMessage(chatID, text)
		msg.ParseMode = "HTML"
		bot.Send(msg)
	}
}

// translateBonusType переводит тип бонуса
func translateBonusType(bonusType string) string {
	switch bonusType {
	case "focus":
		return "Фокус"
	case "willpower":
		return "Сила воли"
	case "knowledge":
		return "Знание Go"
	case "money":
		return "Деньги"
	case "dopamine":
		return "Дофамин"
	default:
		return bonusType
	}
}

// handleUpgradeSkill обрабатывает улучшение навыков
func handleUpgradeSkill(chatID int64) {
	player := players[chatID]
	if player == nil {
		sendNoGameMessage(chatID)
		return
	}

	sendSkills(chatID)
}

// handleEndDay обрабатывает завершение дня
func handleEndDay(chatID int64) {
	player := players[chatID]
	if player == nil {
		sendNoGameMessage(chatID)
		return
	}

	// Финальная битва
	boss := GenerateBossTemptation()

	text := fmt.Sprintf(`🌙 <b>ВЕЧЕР. ФИНАЛЬНАЯ БИТВА!</b>

⚠️  Появляется ОСОБО ОПАСНОЕ ИСКУШЕНИЕ!

👹 <b>%s</b>
Сила: %d%%
Описание: %s

━━━━━━━━━━━━━━━━━━━━

Ваш шанс на победу: %d%%

Нажмите кнопку, чтобы сразиться!`,
		boss.Name,
		boss.Power,
		boss.Description,
		(player.Willpower*2+player.Focus)/3,
	)

	msg := tgbotapi.NewMessage(chatID, text)
	msg.ParseMode = "HTML"
	msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("⚔️ В БОЙ!", "final_battle_fight"),
		),
	)

	bot.Send(msg)
}

// handleFinalBattle обрабатывает финальную битву
func handleFinalBattle(chatID int64) {
	player := players[chatID]
	if player == nil {
		sendNoGameMessage(chatID)
		return
	}

	boss := GenerateBossTemptation()
	won := player.FinalBattle(boss)

	var text string
	if won {
		text = fmt.Sprintf(`🎉 <b>ПОБЕДА!</b>

Вы победили %s и успешно завершили день!

✨ +200 опыта
🎯 Фокус восстановлен: 100%%
💪 Сила воли восстановлена: 100%%
✨ Дофамин +500

🏆 Получено достижение: "Победил финальное искушение"`,
			boss.Name)

		player.Achievements = append(player.Achievements, "Победитель искушений")
	} else {
		text = fmt.Sprintf(`💔 <b>ПОРАЖЕНИЕ...</b>

%s оказался сильнее...

🎯 Фокус: 30%%
💪 Сила воли: 40%%
✨ Дофамин -300

Не сдавайся! Завтра будет новый день!`,
			boss.Name)
	}

	msg := tgbotapi.NewMessage(chatID, text)
	msg.ParseMode = "HTML"

	// Завершаем день
	questsCompleted := player.Quests.ClaimRewards()
	player.CurrentDay++
	player.Hour = 8

	// Сохраняем сессию
	score := player.CalculateScore()
	saveGameSession(chatID, player.CurrentDay-1, score, player.PlayTime, won, questsCompleted)

	// Сохраняем серию дней
	player.Quests.CheckDayStreak(true)
	saveDayStreak(player)

	// Генерируем новые квесты
	player.Quests.GenerateDailyQuests()

	// Сохраняем игрока
	if err := SavePlayer(player); err != nil {
		log.Printf("Ошибка сохранения: %v", err)
	}

	msg.ReplyMarkup = createGameKeyboard(player)
	bot.Send(msg)

	// Показываем статистику
	sendStats(chatID)
}

// checkRandomEvents проверяет случайные события
func checkRandomEvents(chatID int64, player *Player) {
	// 35% шанс искушения
	if randIntn(100) < 35 {
		temptation := GenerateTemptation()
		result := player.HandleTemptation(temptation)

		msg := tgbotapi.NewMessage(chatID, result)
		msg.ParseMode = "HTML"
		bot.Send(msg)

		// Сохраняем после искушения
		if err := SavePlayer(player); err != nil {
			log.Printf("Ошибка сохранения: %v", err)
		}
	}

	// 30% шанс мотивации
	if randIntn(100) < 30 {
		motivation := GetRandomMotivation()
		player.AddExperience(motivation.XPBonus)

		text := fmt.Sprintf(`💪 <b>МОТИВАЦИЯ!</b>

"%s"

✨ +%d опыта
📈 Эффект: %s`,
			motivation.Text,
			motivation.XPBonus,
			motivation.Effect)

		msg := tgbotapi.NewMessage(chatID, text)
		msg.ParseMode = "HTML"
		bot.Send(msg)

		// Сохраняем после мотивации
		if err := SavePlayer(player); err != nil {
			log.Printf("Ошибка сохранения: %v", err)
		}
	}
}

// updatePlayerStatus обновляет статус игрока
func updatePlayerStatus(chatID int64, player *Player) {
	text := player.DisplayStatus()

	msg := tgbotapi.NewMessage(chatID, text)
	msg.ParseMode = "HTML"
	bot.Send(msg)
}

// sendProfile отправляет профиль игрока
func sendProfile(chatID int64) {
	player := players[chatID]
	if player == nil {
		sendNoGameMessage(chatID)
		return
	}

	text := player.DisplayProfile()

	msg := tgbotapi.NewMessage(chatID, text)
	msg.ParseMode = "HTML"
	msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🔙 Главное меню", "main_menu"),
		),
	)

	bot.Send(msg)
}

// sendSkills отправляет дерево навыков
func sendSkills(chatID int64) {
	player := players[chatID]
	if player == nil {
		sendNoGameMessage(chatID)
		return
	}

	text := player.SkillTree.Display()

	msg := tgbotapi.NewMessage(chatID, text)
	msg.ParseMode = "HTML"

	// Добавляем кнопки для улучшения навыков
	keyboard := createSkillsKeyboard(player)
	if len(keyboard) > 0 {
		msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(keyboard...)
	}

	bot.Send(msg)
}

// createSkillsKeyboard создает клавиатуру навыков
func createSkillsKeyboard(player *Player) [][]tgbotapi.InlineKeyboardButton {
	var keyboard [][]tgbotapi.InlineKeyboardButton

	if player.SkillTree.SkillPoints > 0 {
		// Добавляем кнопки для каждого улучшаемого навыка
		for id, skill := range player.SkillTree.Skills {
			if skill.Unlocked && skill.Level < skill.MaxLevel {
				cost := skill.CostPerLevel
				if player.SkillTree.SkillPoints >= cost {
					keyboard = append(keyboard, tgbotapi.NewInlineKeyboardRow(
						tgbotapi.NewInlineKeyboardButtonData(
							fmt.Sprintf("⬆️ %s (%d очк.)", skill.Name, cost),
							fmt.Sprintf("upgrade_%s", id),
						),
					))
				}
			}
		}
	}

	// Добавляем кнопку назад
	keyboard = append(keyboard, tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("🔙 Назад", "main_menu"),
	))

	return keyboard
}

// sendQuests отправляет квесты
func sendQuests(chatID int64) {
	player := players[chatID]
	if player == nil {
		sendNoGameMessage(chatID)
		return
	}

	text := player.Quests.DisplayQuests()

	msg := tgbotapi.NewMessage(chatID, text)
	msg.ParseMode = "HTML"
	msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🔙 Главное меню", "main_menu"),
		),
	)

	bot.Send(msg)
}

// sendStats отправляет статистику
func sendStats(chatID int64) {
	player := players[chatID]
	if player == nil {
		sendNoGameMessage(chatID)
		return
	}

	text := player.DisplayStatistics()

	msg := tgbotapi.NewMessage(chatID, text)
	msg.ParseMode = "HTML"
	msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🔙 Главное меню", "main_menu"),
		),
	)

	bot.Send(msg)
}

// showMainMenu показывает главное меню
func showMainMenu(chatID int64) {
	text := `🎮 <b>ГЛАВНОЕ МЕНЮ</b>

Выберите действие:`

	msg := tgbotapi.NewMessage(chatID, text)
	msg.ParseMode = "HTML"
	msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("📚 Учить Go", "study_go"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("📋 Квесты", "check_quests"),
			tgbotapi.NewInlineKeyboardButtonData("🌳 Навыки", "upgrade_skill"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("📊 Статистика", "show_stats"),
			tgbotapi.NewInlineKeyboardButtonData("💾 Сохранить", "save_game"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🌙 Завершить день", "end_day"),
		),
	)

	bot.Send(msg)
}

// saveGame сохраняет игру
func saveGame(chatID int64) {
	player := players[chatID]
	if player == nil {
		sendNoGameMessage(chatID)
		return
	}

	// Сохраняем в БД
	if err := SavePlayer(player); err != nil {
		log.Printf("Ошибка сохранения: %v", err)
		text := `❌ <b>ОШИБКА СОХРАНЕНИЯ!</b>

Произошла ошибка при сохранении прогресса.`

		msg := tgbotapi.NewMessage(chatID, text)
		msg.ParseMode = "HTML"
		bot.Send(msg)
		return
	}

	text := `💾 <b>ПРОГРЕСС СОХРАНЁН!</b>

Ваш прогресс сохранён в базе данных.
Вы можете продолжить игру в любой момент!

📊 Текущий статус:
` + player.DisplayStatus()

	msg := tgbotapi.NewMessage(chatID, text)
	msg.ParseMode = "HTML"
	bot.Send(msg)
}

// loadGame загружает игру
func loadGame(chatID int64) {
	player, err := LoadPlayer(chatID)
	if err != nil {
		text := `❌ <b>ОШИБКА ЗАГРУЗКИ!</b>

Произошла ошибка при загрузке сохранения.`

		msg := tgbotapi.NewMessage(chatID, text)
		msg.ParseMode = "HTML"
		bot.Send(msg)
		return
	}

	if player == nil {
		text := `⚠️  <b>СОХРАНЕНИЕ НЕ НАЙДЕНО</b>

У вас нет сохранённого прогресса.
Начните новую игру командой /play`

		msg := tgbotapi.NewMessage(chatID, text)
		msg.ParseMode = "HTML"
		msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("🎮 Начать игру", "start_game"),
			),
		)
		bot.Send(msg)
		return
	}

	// Загружаем игрока в кэш
	players[chatID] = player

	text := `💾 <b>ИГРА ЗАГРУЖЕНА!</b>

Ваш прогресс успешно загружён.

` + player.DisplayStatus()

	msg := tgbotapi.NewMessage(chatID, text)
	msg.ParseMode = "HTML"
	msg.ReplyMarkup = createGameKeyboard(player)
	bot.Send(msg)
}

// sendNoGameMessage отправляет сообщение, если игра не начата
func sendNoGameMessage(chatID int64) {
	text := `⚠️  <b>ИГРА НЕ НАЧАТА</b>

Сначала начните игру командой /play

Или нажмите кнопку "Начать игру".`

	msg := tgbotapi.NewMessage(chatID, text)
	msg.ParseMode = "HTML"
	msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🎮 Начать игру", "start_game"),
		),
	)

	bot.Send(msg)
}

// sendLeaderboard отправляет таблицу лидеров
func sendLeaderboard(chatID int64) {
	leaderboard, err := GetLeaderboard(10)
	if err != nil {
		text := `❌ <b>ОШИБКА!</b>

Не удалось загрузить таблицу лидеров.`

		msg := tgbotapi.NewMessage(chatID, text)
		msg.ParseMode = "HTML"
		bot.Send(msg)
		return
	}

	text := `🏆 <b>ТАБЛИЦА ЛИДЕРОВ</b>
━━━━━━━━━━━━━━━━━━━━

Топ-10 игроков FocusGo:

`

	for i, entry := range leaderboard {
		text += fmt.Sprintf("%d. <b>%s</b> — Ур.%d | Рейтинг: %d\n",
			i+1, entry["name"], entry["level"], entry["rating"])
	}

	// Добавляем общую статистику
	totalPlayers, _ := GetTotalPlayers()
	text += fmt.Sprintf("\n📊 Всего игроков: %d", totalPlayers)

	msg := tgbotapi.NewMessage(chatID, text)
	msg.ParseMode = "HTML"
	msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🔙 Главное меню", "main_menu"),
		),
	)

	bot.Send(msg)
}

// Вспомогательная функция для случайных чисел
func randIntn(n int) int {
	return int(time.Now().UnixNano()%int64(n)) % n
}

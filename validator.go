package main

import (
	"fmt"
	"log"
)

// ============================================================================
// КОНСТАНТЫ ВАЛИДАЦИИ
// ============================================================================

const (
	// Минимальные и максимальные значения характеристик
	MinStatValue     = 0
	MaxStatValue     = 100
	MinMoneyValue    = 0
	MaxMoneyValue    = 999999
	MinDopamineValue = 0
	MaxDopamineValue = 999

	// Минимальные и максимальные значения опыта и уровня
	MinLevel       = 1
	MaxLevel       = 100
	MinExperience  = 0
	MaxExperience  = 999999

	// Минимальные и максимальные значения времени
	MinPlayTime   = 0
	MaxPlayTime   = 999999
	MinDaysPlayed = 1
	MaxDaysPlayed = 9999

	// Минимальные и максимальные значения игрового времени
	MinHour = 0
	MaxHour = 23

	// Максимальная длина строк
	MaxNameLength        = 50
	MaxAchievementLength = 200
	MaxTemptationLength  = 100
)

// ============================================================================
// ФУНКЦИИ CLAMP (ОГРАНИЧЕНИЕ ДИАПАЗОНА)
// ============================================================================

// ClampInt ограничивает целое число диапазоном [min, max]
func ClampInt(value, min, max int) int {
	if value < min {
		return min
	}
	if value > max {
		return max
	}
	return value
}

// ClampStat ограничивает характеристику диапазоном [0, 100]
func ClampStat(value int) int {
	return ClampInt(value, MinStatValue, MaxStatValue)
}

// ClampMoney ограничивает деньги диапазоном [0, 999999]
func ClampMoney(value int) int {
	return ClampInt(value, MinMoneyValue, MaxMoneyValue)
}

// ClampDopamine ограничивает дофамин диапазоном [0, 999]
func ClampDopamine(value int) int {
	return ClampInt(value, MinDopamineValue, MaxDopamineValue)
}

// ClampExperience ограничивает опыт диапазоном [0, 999999]
func ClampExperience(value int) int {
	return ClampInt(value, MinExperience, MaxExperience)
}

// ClampLevel ограничивает уровень диапазоном [1, 100]
func ClampLevel(value int) int {
	return ClampInt(value, MinLevel, MaxLevel)
}

// ClampHour ограничивает час диапазоном [0, 23]
func ClampHour(value int) int {
	return ClampInt(value, MinHour, MaxHour)
}

// ClampStringLength ограничивает длину строки
func ClampStringLength(s string, maxLen int) string {
	if len(s) > maxLen {
		return s[:maxLen]
	}
	return s
}

// ============================================================================
// ВАЛИДАЦИЯ ИГРОКА
// ============================================================================

// ValidatePlayer проверяет валидность всех характеристик игрока
func ValidatePlayer(player *Player) []string {
	var errors []string

	// Проверка имени
	if player.Name == "" {
		errors = append(errors, "Имя игрока не может быть пустым")
	} else if len(player.Name) > MaxNameLength {
		errors = append(errors, fmt.Sprintf("Имя игрока слишком длинное (макс. %d символов)", MaxNameLength))
	}

	// Проверка chat_id
	if player.ChatID == 0 {
		errors = append(errors, "ChatID не может быть нулевым")
	}

	// Проверка уровня
	if player.Level < MinLevel || player.Level > MaxLevel {
		errors = append(errors, fmt.Sprintf("Уровень вне диапазона [%d-%d]: %d", MinLevel, MaxLevel, player.Level))
	}

	// Проверка опыта
	if player.Experience < MinExperience || player.Experience > MaxExperience {
		errors = append(errors, fmt.Sprintf("Опыт вне диапазона [%d-%d]: %d", MinExperience, MaxExperience, player.Experience))
	}

	// Проверка характеристик (0-100)
	if player.Focus < MinStatValue || player.Focus > MaxStatValue {
		errors = append(errors, fmt.Sprintf("Фокус вне диапазона [%d-%d]: %d", MinStatValue, MaxStatValue, player.Focus))
	}

	if player.Willpower < MinStatValue || player.Willpower > MaxStatValue {
		errors = append(errors, fmt.Sprintf("Сила воли вне диапазона [%d-%d]: %d", MinStatValue, MaxStatValue, player.Willpower))
	}

	if player.GoKnowledge < MinStatValue || player.GoKnowledge > MaxStatValue {
		errors = append(errors, fmt.Sprintf("Знание Go вне диапазона [%d-%d]: %d", MinStatValue, MaxStatValue, player.GoKnowledge))
	}

	// Проверка денег
	if player.Money < MinMoneyValue || player.Money > MaxMoneyValue {
		errors = append(errors, fmt.Sprintf("Деньги вне диапазона [%d-%d]: %d", MinMoneyValue, MaxMoneyValue, player.Money))
	}

	// Проверка дофамина
	if player.Dopamine < MinDopamineValue || player.Dopamine > MaxDopamineValue {
		errors = append(errors, fmt.Sprintf("Дофамин вне диапазона [%d-%d]: %d", MinDopamineValue, MaxDopamineValue, player.Dopamine))
	}

	// Проверка времени игры
	if player.PlayTime < MinPlayTime || player.PlayTime > MaxPlayTime {
		errors = append(errors, fmt.Sprintf("Время игры вне диапазона [%d-%d]: %d", MinPlayTime, MaxPlayTime, player.PlayTime))
	}

	// Проверка дней сыграно
	if player.DaysPlayed < MinDaysPlayed || player.DaysPlayed > MaxDaysPlayed {
		errors = append(errors, fmt.Sprintf("Дней сыграно вне диапазона [%d-%d]: %d", MinDaysPlayed, MaxDaysPlayed, player.DaysPlayed))
	}

	// Проверка текущего дня
	if player.CurrentDay < 1 || player.CurrentDay > MaxDaysPlayed {
		errors = append(errors, fmt.Sprintf("Текущий день вне диапазона [1-%d]: %d", MaxDaysPlayed, player.CurrentDay))
	}

	// Проверка игрового часа
	if player.Hour < MinHour || player.Hour > MaxHour {
		errors = append(errors, fmt.Sprintf("Игровой час вне диапазона [%d-%d]: %d", MinHour, MaxHour, player.Hour))
	}

	return errors
}

// SanitizePlayer исправляет некорректные значения характеристик игрока
func SanitizePlayer(player *Player) {
	// Ограничиваем характеристики диапазоном [0, 100]
	player.Focus = ClampStat(player.Focus)
	player.Willpower = ClampStat(player.Willpower)
	player.GoKnowledge = ClampStat(player.GoKnowledge)

	// Ограничиваем деньги и дофамин
	player.Money = ClampMoney(player.Money)
	player.Dopamine = ClampDopamine(player.Dopamine)

	// Ограничиваем опыт и уровень
	player.Experience = ClampExperience(player.Experience)
	player.Level = ClampLevel(player.Level)

	// Ограничиваем время
	player.PlayTime = ClampInt(player.PlayTime, MinPlayTime, MaxPlayTime)
	player.DaysPlayed = ClampInt(player.DaysPlayed, MinDaysPlayed, MaxDaysPlayed)
	player.CurrentDay = ClampInt(player.CurrentDay, 1, MaxDaysPlayed)
	player.Hour = ClampHour(player.Hour)

	// Ограничиваем длину имени
	player.Name = ClampStringLength(player.Name, MaxNameLength)
	if player.Name == "" {
		player.Name = "Игрок"
	}

	// Проверяем chat_id
	if player.ChatID == 0 {
		log.Printf("⚠️  WARNING: ChatID игрока %s равен 0, устанавливаем 1", player.Name)
		player.ChatID = 1
	}
}

// ValidateAndSanitize проверяет и исправляет игрока, возвращает ошибки
func ValidateAndSanitize(player *Player) []string {
	// Сначала sanitaze, потом validate
	SanitizePlayer(player)
	return ValidatePlayer(player)
}

// ============================================================================
// ВАЛИДАЦИЯ НАВЫКОВ
// ============================================================================

// ValidateSkill проверяет валидность навыка
func ValidateSkill(skill *Skill) []string {
	var errors []string

	if skill.ID == "" {
		errors = append(errors, "ID навыка не может быть пустым")
	}

	if skill.Name == "" {
		errors = append(errors, "Название навыка не может быть пустым")
	}

	if skill.Level < 0 || skill.Level > skill.MaxLevel {
		errors = append(errors, fmt.Sprintf("Уровень навыка %s вне диапазона [0-%d]: %d",
			skill.Name, skill.MaxLevel, skill.Level))
	}

	if skill.MaxLevel < 1 || skill.MaxLevel > 10 {
		errors = append(errors, fmt.Sprintf("Максимальный уровень навыка %s вне диапазона [1-10]: %d",
			skill.Name, skill.MaxLevel))
	}

	if skill.CostPerLevel < 1 || skill.CostPerLevel > 100 {
		errors = append(errors, fmt.Sprintf("Стоимость улучшения навыка %s вне диапазона [1-100]: %d",
			skill.Name, skill.CostPerLevel))
	}

	if skill.BonusValue < 0 || skill.BonusValue > 1000 {
		errors = append(errors, fmt.Sprintf("Бонус навыка %s вне диапазона [0-1000]: %d",
			skill.Name, skill.BonusValue))
	}

	return errors
}

// SanitizeSkill исправляет некорректные значения навыка
func SanitizeSkill(skill *Skill) {
	if skill.Level < 0 {
		skill.Level = 0
	}
	if skill.Level > skill.MaxLevel {
		skill.Level = skill.MaxLevel
	}

	if skill.MaxLevel < 1 {
		skill.MaxLevel = 1
	}
	if skill.MaxLevel > 10 {
		skill.MaxLevel = 10
	}

	if skill.CostPerLevel < 1 {
		skill.CostPerLevel = 1
	}
	if skill.CostPerLevel > 100 {
		skill.CostPerLevel = 100
	}

	if skill.BonusValue < 0 {
		skill.BonusValue = 0
	}
	if skill.BonusValue > 1000 {
		skill.BonusValue = 1000
	}
}

// ValidateSkillTree проверяет всё дерево навыков
func ValidateSkillTree(tree *SkillTree) []string {
	var errors []string

	if tree.SkillPoints < 0 {
		errors = append(errors, fmt.Sprintf("Очки навыков отрицательные: %d", tree.SkillPoints))
		tree.SkillPoints = 0
	}

	if tree.TotalPoints < 0 {
		errors = append(errors, fmt.Sprintf("Всего очков навыков отрицательное: %d", tree.TotalPoints))
		tree.TotalPoints = 0
	}

	// Проверяем каждый навык
	for id, skill := range tree.Skills {
		skillErrors := ValidateSkill(skill)
		for _, err := range skillErrors {
			errors = append(errors, fmt.Sprintf("Навык %s: %s", id, err))
		}
	}

	return errors
}

// ============================================================================
// ВАЛИДАЦИЯ КВЕСТОВ
// ============================================================================

// ValidateQuest проверяет валидность квеста
func ValidateQuest(quest *DailyQuest) []string {
	var errors []string

	if quest.ID == "" {
		errors = append(errors, "ID квеста не может быть пустым")
	}

	if quest.Title == "" {
		errors = append(errors, "Название квеста не может быть пустым")
	} else if len(quest.Title) > 100 {
		errors = append(errors, fmt.Sprintf("Название квеста слишком длинное (макс. 100 символов): %d", len(quest.Title)))
	}

	if quest.Goal < 1 || quest.Goal > 10000 {
		errors = append(errors, fmt.Sprintf("Цель квеста вне диапазона [1-10000]: %d", quest.Goal))
	}

	if quest.Progress < 0 {
		errors = append(errors, fmt.Sprintf("Прогресс квеста отрицательный: %d", quest.Progress))
	}

	if quest.Reward < 0 || quest.Reward > 1000 {
		errors = append(errors, fmt.Sprintf("Награда за квест вне диапазона [0-1000]: %d", quest.Reward))
	}

	return errors
}

// SanitizeQuest исправляет некорректные значения квеста
func SanitizeQuest(quest *DailyQuest) {
	if quest.Progress < 0 {
		quest.Progress = 0
	}

	if quest.Goal < 1 {
		quest.Goal = 1
	}
	if quest.Goal > 10000 {
		quest.Goal = 10000
	}

	if quest.Reward < 0 {
		quest.Reward = 0
	}
	if quest.Reward > 1000 {
		quest.Reward = 1000
	}

	quest.Title = ClampStringLength(quest.Title, 100)
}

// ValidateQuestSystem проверяет систему квестов
func ValidateQuestSystem(qs *QuestSystem) []string {
	var errors []string

	if qs.DayStreak < 0 {
		errors = append(errors, fmt.Sprintf("Серия дней отрицательная: %d", qs.DayStreak))
		qs.DayStreak = 0
	}

	if qs.TotalCompleted < 0 {
		errors = append(errors, fmt.Sprintf("Всего выполнено квестов отрицательное: %d", qs.TotalCompleted))
		qs.TotalCompleted = 0
	}

	// Проверяем каждый квест
	for i, quest := range qs.Quests {
		questErrors := ValidateQuest(quest)
		for _, err := range questErrors {
			errors = append(errors, fmt.Sprintf("Квест #%d: %s", i+1, err))
		}
	}

	return errors
}

// ============================================================================
// ВАЛИДАЦИЯ ИСКУШЕНИЙ
// ============================================================================

// ValidateTemptation проверяет валидность искушения
func ValidateTemptation(t *Temptation) []string {
	var errors []string

	if t.Name == "" {
		errors = append(errors, "Название искушения не может быть пустым")
	} else if len(t.Name) > MaxTemptationLength {
		errors = append(errors, fmt.Sprintf("Название искушения слишком длинное (макс. %d символов)", MaxTemptationLength))
	}

	if t.Power < 0 || t.Power > 100 {
		errors = append(errors, fmt.Sprintf("Сила искушения вне диапазона [0-100]: %d", t.Power))
	}

	if t.XPLoss < 0 || t.XPLoss > 1000 {
		errors = append(errors, fmt.Sprintf("Потеря опыта от искушения вне диапазона [0-1000]: %d", t.XPLoss))
	}

	return errors
}

// ============================================================================
// ВАЛИДАЦИЯ МОТИВАЦИИ
// ============================================================================

// ValidateMotivation проверяет валидность мотивации
func ValidateMotivation(m *Motivation) []string {
	var errors []string

	if m.Text == "" {
		errors = append(errors, "Текст мотивации не может быть пустым")
	} else if len(m.Text) > 500 {
		errors = append(errors, fmt.Sprintf("Текст мотивации слишком длинный (макс. 500 символов): %d", len(m.Text)))
	}

	if m.XPBonus < 0 || m.XPBonus > 10000 {
		errors = append(errors, fmt.Sprintf("Бонус опыта мотивации вне диапазона [0-10000]: %d", m.XPBonus))
	}

	return errors
}

// ============================================================================
// ВАЛИДАЦИЯ ПРИ СОХРАНЕНИИ В БД
// ============================================================================

// ValidateBeforeSave проверяет игрока перед сохранением в БД
func ValidateBeforeSave(player *Player) error {
	errors := ValidateAndSanitize(player)

	if len(errors) > 0 {
		log.Printf("⚠️  WARNING: Валидация игрока %s перед сохранением:", player.Name)
		for _, err := range errors {
			log.Printf("  - %s", err)
		}

		// Логируем ошибки, но не блокируем сохранение
		// SanitizePlayer уже исправил критичные значения
	}

	return nil
}

// ============================================================================
// ВАЛИДАЦИЯ ПРИ ЗАГРУЗКЕ ИЗ БД
// ============================================================================

// ValidateAfterLoad проверяет игрока после загрузки из БД
func ValidateAfterLoad(player *Player) error {
	errors := ValidateAndSanitize(player)

	if len(errors) > 0 {
		log.Printf("⚠️  WARNING: Валидация игрока %s после загрузки:", player.Name)
		for _, err := range errors {
			log.Printf("  - %s", err)
		}
	}

	return nil
}

// ============================================================================
// ВСПОМОГАТЕЛЬНЫЕ ФУНКЦИИ
// ============================================================================

// IsValidName проверяет корректность имени
func IsValidName(name string) bool {
	if name == "" || len(name) > MaxNameLength {
		return false
	}
	// Можно добавить проверку на допустимые символы
	return true
}

// IsValidChatID проверяет корректность chat_id
func IsValidChatID(chatID int64) bool {
	return chatID > 0
}

// FormatValidationErrors форматирует ошибки валидации для вывода
func FormatValidationErrors(errors []string) string {
	if len(errors) == 0 {
		return ""
	}

	result := "❌ <b>ОШИБКИ ВАЛИДАЦИИ:</b>\n"
	for i, err := range errors {
		result += fmt.Sprintf("%d. %s\n", i+1, err)
	}
	return result
}

// LogValidationErrors логирует ошибки валидации
func LogValidationErrors(context string, errors []string) {
	if len(errors) == 0 {
		return
	}

	log.Printf("⚠️  Валидация: %s", context)
	for _, err := range errors {
		log.Printf("  - %s", err)
	}
}

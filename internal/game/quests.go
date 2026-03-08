package game

import (
	"fmt"
	"time"
)

// Quest представляет квест
type Quest struct {
	ID          string
	Title       string
	Description string
	Goal        int
	Progress    int
	Reward      int // очки навыков
	Completed   bool
	Deadline    string
}

// QuestSystem система квестов
type QuestSystem struct {
	ChatID       int64
	Quests       []*Quest
	DayStreak    int
	TotalCompleted int
}

// NewQuestSystem создаёт систему квестов
func NewQuestSystem(chatID int64) *QuestSystem {
	return &QuestSystem{
		ChatID: chatID,
		Quests: make([]*Quest, 0),
	}
}

// GenerateDailyQuests генерирует ежедневные квесты
func (qs *QuestSystem) GenerateDailyQuests() {
	today := time.Now().Format("2006-01-02")

	qs.Quests = []*Quest{
		{
			ID:          "study_go_30min",
			Title:       "30 минут Go",
			Description: "Изучи Go в течение 30 минут",
			Goal:        30,
			Progress:    0,
			Reward:      2,
			Completed:   false,
			Deadline:    today,
		},
		{
			ID:          "resist_temptation",
			Title:       "Борец с искушениями",
			Description: "Сопротивляйся 3 искушениям",
			Goal:        3,
			Progress:    0,
			Reward:      3,
			Completed:   false,
			Deadline:    today,
		},
		{
			ID:          "code_practice",
			Title:       "Практика кода",
			Description: "Напиши 50+ строк кода на Go",
			Goal:        50,
			Progress:    0,
			Reward:      3,
			Completed:   false,
			Deadline:    today,
		},
		{
			ID:          "morning_routine",
			Title:       "Утренний ритуал",
			Description: "Выполни утренний ритуал",
			Goal:        1,
			Progress:    0,
			Reward:      1,
			Completed:   false,
			Deadline:    today,
		},
		{
			ID:          "no_social_media",
			Title:       "Цифровой детокс",
			Description: "Не заходи в соцсети 4 часа",
			Goal:        4,
			Progress:    0,
			Reward:      2,
			Completed:   false,
			Deadline:    today,
		},
	}
}

// UpdateProgress обновляет прогресс квеста
func (qs *QuestSystem) UpdateProgress(questID string, progress int) (bool, int) {
	for _, quest := range qs.Quests {
		if quest.ID == questID && !quest.Completed {
			quest.Progress += progress
			if quest.Progress >= quest.Goal {
				quest.Completed = true
				qs.TotalCompleted++
				return true, quest.Reward
			}
			return false, 0
		}
	}
	return false, 0
}

// GetCompletedCount возвращает количество выполненных квестов
func (qs *QuestSystem) GetCompletedCount() int {
	count := 0
	for _, quest := range qs.Quests {
		if quest.Completed {
			count++
		}
	}
	return count
}

// Display отображает квесты
func (qs *QuestSystem) Display() string {
	if len(qs.Quests) == 0 {
		return "📋 <b>КВЕСТЫ</b>\n\nНет активных квестов.\n\nИспользуйте /start для начала игры."
	}

	text := "📋 <b>ЕЖЕДНЕВНЫЕ КВЕСТЫ</b>\n━━━━━━━━━━━━━━━━━━━━\n\n"

	for _, quest := range qs.Quests {
		status := "⏳"
		if quest.Completed {
			status = "✅"
		}

		progress := quest.Progress
		if progress > quest.Goal {
			progress = quest.Goal
		}

		text += fmt.Sprintf("%s <b>%s</b>\n   %s\n   Прогресс: %d/%d | Награда: %d очк.\n\n",
			status, quest.Title, quest.Description,
			progress, quest.Goal, quest.Reward)
	}

	text += fmt.Sprintf("🔥 Серия дней: %d\n", qs.DayStreak)
	text += fmt.Sprintf("🏆 Выполнено квестов: %d\n", qs.TotalCompleted)

	return text
}

// CheckDayStreak проверяет серию дней
func (qs *QuestSystem) CheckDayStreak(allCompleted bool) {
	if allCompleted {
		qs.DayStreak++
	} else {
		if qs.DayStreak > 0 {
			qs.DayStreak = 0
		}
	}
}

// ClaimRewards забирает награды за выполненные квесты
func (qs *QuestSystem) ClaimRewards() int {
	totalReward := 0
	for _, quest := range qs.Quests {
		if quest.Completed && quest.Reward > 0 {
			totalReward += quest.Reward
			quest.Reward = 0 // Сбрасываем награду
		}
	}
	return totalReward
}

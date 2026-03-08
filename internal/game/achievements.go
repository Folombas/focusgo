package game

import (
	"fmt"
	"strings"
	"time"

	"focusgo/internal/database"
)

// Achievement представляет достижение игрока
type Achievement struct {
	ID          string
	Name        string
	Description string
	Icon        string
	Unlocked    bool
	UnlockedAt  time.Time
}

// AchievementSystem система достижений игрока
type AchievementSystem struct {
	ChatID       int64
	Achievements map[string]*Achievement
	TotalUnlocked int
}

// Список всех достижений (20+)
var allAchievements = map[string]*Achievement{
	// Достижения за уровни
	"first_steps": {
		ID:          "first_steps",
		Name:        "Первые шаги",
		Description: "Достичь уровня 2",
		Icon:        "👶",
	},
	"go_newbie": {
		ID:          "go_newbie",
		Name:        "Go-Новичок",
		Description: "Достичь уровня 5",
		Icon:        "🌱",
	},
	"go_developer": {
		ID:          "go_developer",
		Name:        "Go-Разработчик",
		Description: "Достичь уровня 10",
		Icon:        "💻",
	},
	"go_master": {
		ID:          "go_master",
		Name:        "Go-Мастер",
		Description: "Достичь уровня 20",
		Icon:        "🎓",
	},
	"go_legend": {
		ID:          "go_legend",
		Name:        "Go-Легенда",
		Description: "Достичь уровня 30",
		Icon:        "🏆",
	},

	// Достижения за квесты
	"quest_beginner": {
		ID:          "quest_beginner",
		Name:        "Начинающий квестер",
		Description: "Выполнить 10 квестов",
		Icon:        "📋",
	},
	"quest_hunter": {
		ID:          "quest_hunter",
		Name:        "Охотник за квестами",
		Description: "Выполнить 50 квестов",
		Icon:        "🎯",
	},
	"quest_master": {
		ID:          "quest_master",
		Name:        "Мастер квестов",
		Description: "Выполнить 100 квестов",
		Icon:        "👑",
	},

	// Достижения за серию дней
	"week_streak": {
		ID:          "week_streak",
		Name:        "Недельная серия",
		Description: "7 дней подряд в игре",
		Icon:        "🔥",
	},
	"month_streak": {
		ID:          "month_streak",
		Name:        "Месячная серия",
		Description: "30 дней подряд в игре",
		Icon:        "🔥🔥",
	},

	// Достижения за искушения
	"temptation_resister": {
		ID:          "temptation_resister",
		Name:        "Борец с искушениями",
		Description: "Преодолеть 10 искушений",
		Icon:        "🛡️",
	},
	"temptation_warrior": {
		ID:          "temptation_warrior",
		Name:        "Воин искушений",
		Description: "Преодолеть 50 искушений",
		Icon:        "⚔️",
	},
	"temptation_legend": {
		ID:          "temptation_legend",
		Name:        "Легенда сопротивления",
		Description: "Преодолеть 100 искушений",
		Icon:        "🏅",
	},

	// Достижения за боссов
	"boss_slayer": {
		ID:          "boss_slayer",
		Name:        "Убийца боссов",
		Description: "Победить 10 боссов",
		Icon:        "🗡️",
	},
	"boss_legend": {
		ID:          "boss_legend",
		Name:        "Легенда боссов",
		Description: "Победить 50 боссов",
		Icon:        "🐉",
	},

	// Достижения за изучение Go
	"go_student": {
		ID:          "go_student",
		Name:        "Студент Go",
		Description: "Изучить Go на 50%",
		Icon:        "📖",
	},
	"go_scholar": {
		ID:          "go_scholar",
		Name:        "Учёный Go",
		Description: "Изучить Go на 100%",
		Icon:        "🎓",
	},

	// Достижения за навыки
	"skill_master": {
		ID:          "skill_master",
		Name:        "Мастер навыков",
		Description: "Улучшить любой навык до максимума",
		Icon:        "⭐",
	},
	"skill_collector": {
		ID:          "skill_collector",
		Name:        "Коллекционер навыков",
		Description: "Разблокировать все навыки",
		Icon:        "🌟",
	},

	// Специальные достижения
	"perfectionist": {
		ID:          "perfectionist",
		Name:        "Перфекционист",
		Description: "Выполнить все квесты за день",
		Icon:        "💎",
	},
	"marathon_runner": {
		ID:          "marathon_runner",
		Name:        "Марафонец",
		Description: "Играть 1000 минут",
		Icon:        "🏃",
	},
	"early_bird": {
		ID:          "early_bird",
		Name:        "Ранняя пташка",
		Description: "Начать день до 9 утра",
		Icon:        "🌅",
	},
	"night_owl": {
		ID:          "night_owl",
		Name:        "Ночная сова",
		Description: "Завершить день после 23 часов",
		Icon:        "🦉",
	},
}

// NewAchievementSystem создаёт новую систему достижений
func NewAchievementSystem(chatID int64) *AchievementSystem {
	as := &AchievementSystem{
		ChatID:       chatID,
		Achievements: make(map[string]*Achievement),
		TotalUnlocked: 0,
	}

	// Копируем все достижения
	for id, achievement := range allAchievements {
		as.Achievements[id] = &Achievement{
			ID:          achievement.ID,
			Name:        achievement.Name,
			Description: achievement.Description,
			Icon:        achievement.Icon,
			Unlocked:    false,
		}
	}

	return as
}

// LoadAchievementSystem загружает систему достижений из БД
func LoadAchievementSystem(chatID int64) (*AchievementSystem, error) {
	as := NewAchievementSystem(chatID)

	query := `SELECT achievement_id, unlocked_at FROM achievements WHERE chat_id = ?`
	rows, err := database.DB.Query(query, chatID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var achievementID string
		var unlockedAt time.Time

		if err := rows.Scan(&achievementID, &unlockedAt); err != nil {
			continue
		}

		if achievement, exists := as.Achievements[achievementID]; exists {
			achievement.Unlocked = true
			achievement.UnlockedAt = unlockedAt
			as.TotalUnlocked++
		}
	}

	return as, rows.Err()
}

// SaveAchievementSystem сохраняет систему достижений в БД
func (as *AchievementSystem) SaveAchievementSystem() error {
	for _, achievement := range as.Achievements {
		if achievement.Unlocked {
			query := `INSERT OR REPLACE INTO achievements (chat_id, achievement_id, unlocked_at) VALUES (?, ?, ?)`
			_, err := database.DB.Exec(query, as.ChatID, achievement.ID, achievement.UnlockedAt)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// CheckAndUnlock проверяет и разблокирует достижение
func (as *AchievementSystem) CheckAndUnlock(achievementID string, condition bool) bool {
	if condition {
		if achievement, exists := as.Achievements[achievementID]; exists {
			if !achievement.Unlocked {
				achievement.Unlocked = true
				achievement.UnlockedAt = time.Now()
				as.TotalUnlocked++
				return true
			}
		}
	}
	return false
}

// CheckAchievements проверяет все достижения игрока
func (as *AchievementSystem) CheckAchievements(state *GameState, tree *SkillTree, questSystem *QuestSystem) []string {
	unlocked := []string{}

	// Достижения за уровни
	if as.CheckAndUnlock("first_steps", state.Level >= 2) {
		unlocked = append(unlocked, as.Achievements["first_steps"].Name)
	}
	if as.CheckAndUnlock("go_newbie", state.Level >= 5) {
		unlocked = append(unlocked, as.Achievements["go_newbie"].Name)
	}
	if as.CheckAndUnlock("go_developer", state.Level >= 10) {
		unlocked = append(unlocked, as.Achievements["go_developer"].Name)
	}
	if as.CheckAndUnlock("go_master", state.Level >= 20) {
		unlocked = append(unlocked, as.Achievements["go_master"].Name)
	}
	if as.CheckAndUnlock("go_legend", state.Level >= 30) {
		unlocked = append(unlocked, as.Achievements["go_legend"].Name)
	}

	// Достижения за квесты
	if questSystem != nil {
		if as.CheckAndUnlock("quest_beginner", questSystem.TotalCompleted >= 10) {
			unlocked = append(unlocked, as.Achievements["quest_beginner"].Name)
		}
		if as.CheckAndUnlock("quest_hunter", questSystem.TotalCompleted >= 50) {
			unlocked = append(unlocked, as.Achievements["quest_hunter"].Name)
		}
		if as.CheckAndUnlock("quest_master", questSystem.TotalCompleted >= 100) {
			unlocked = append(unlocked, as.Achievements["quest_master"].Name)
		}

		// Достижения за серию дней
		if as.CheckAndUnlock("week_streak", questSystem.DayStreak >= 7) {
			unlocked = append(unlocked, as.Achievements["week_streak"].Name)
		}
		if as.CheckAndUnlock("month_streak", questSystem.DayStreak >= 30) {
			unlocked = append(unlocked, as.Achievements["month_streak"].Name)
		}
	}

	// Достижения за изучение Go
	if as.CheckAndUnlock("go_student", state.GoKnowledge >= 50) {
		unlocked = append(unlocked, as.Achievements["go_student"].Name)
	}
	if as.CheckAndUnlock("go_scholar", state.GoKnowledge >= 100) {
		unlocked = append(unlocked, as.Achievements["go_scholar"].Name)
	}

	// Достижения за время игры
	if as.CheckAndUnlock("marathon_runner", state.PlayTime >= 1000) {
		unlocked = append(unlocked, as.Achievements["marathon_runner"].Name)
	}

	// Достижения за навыки
	if tree != nil {
		allUnlocked := true
		maxLevel := false
		for _, skill := range tree.Skills {
			if !skill.Unlocked {
				allUnlocked = false
			}
			if skill.Level >= skill.MaxLevel {
				maxLevel = true
			}
		}
		if as.CheckAndUnlock("skill_collector", allUnlocked) {
			unlocked = append(unlocked, as.Achievements["skill_collector"].Name)
		}
		if as.CheckAndUnlock("skill_master", maxLevel) {
			unlocked = append(unlocked, as.Achievements["skill_master"].Name)
		}
	}

	return unlocked
}

// Display отображает достижения
func (as *AchievementSystem) Display() string {
	var sb strings.Builder

	sb.WriteString("🏆 <b>ДОСТИЖЕНИЯ</b>\n")
	sb.WriteString("━━━━━━━━━━━━━━━━━━━━\n\n")
	sb.WriteString(fmt.Sprintf("Всего разблокировано: %d/%d\n\n", as.TotalUnlocked, len(allAchievements)))

	// Разблокированные
	sb.WriteString("<b>🔓 Разблокированные:</b>\n")
	hasUnlocked := false
	for _, achievement := range as.Achievements {
		if achievement.Unlocked {
			sb.WriteString(fmt.Sprintf("%s %s — %s\n", achievement.Icon, achievement.Name, achievement.Description))
			hasUnlocked = true
		}
	}
	if !hasUnlocked {
		sb.WriteString("Пока нет разблокированных достижений\n")
	}

	sb.WriteString("\n")

	// Заблокированные
	sb.WriteString("<b>🔒 Заблокированные:</b>\n")
	hasLocked := false
	for _, achievement := range as.Achievements {
		if !achievement.Unlocked {
			sb.WriteString(fmt.Sprintf("🔒 %s — %s\n", achievement.Name, achievement.Description))
			hasLocked = true
		}
	}
	if !hasLocked {
		sb.WriteString("Все достижения разблокированы! 🎉\n")
	}

	return sb.String()
}

// GetUnlockedCount возвращает количество разблокированных достижений
func (as *AchievementSystem) GetUnlockedCount() int {
	return as.TotalUnlocked
}

// GetTotalCount возвращает общее количество достижений
func (as *AchievementSystem) GetTotalCount() int {
	return len(allAchievements)
}

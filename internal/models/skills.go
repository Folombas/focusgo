package models

import (
	"fmt"
	"log"
	"strings"
	"time"
)

// Skill представляет навык в дереве развития
type Skill struct {
	ID            string
	Name          string
	Description   string
	Level         int      // Уровень навыка (0-5)
	MaxLevel      int      // Максимальный уровень
	CostPerLevel  int      // Стоимость улучшения за уровень
	BonusType     string   // Тип бонуса: focus, willpower, knowledge, money, dopamine
	BonusValue    int      // Значение бонуса за уровень
	Unlocked      bool     // Разблокирован ли навык
	Prerequisites []string // Требуемые навыки для разблокировки
}

// SkillTree представляет дерево навыков игрока
type SkillTree struct {
	Skills      map[string]*Skill
	SkillPoints int // Доступные очки навыков
	TotalPoints int // Всего заработано очков
}

// DailyQuest представляет ежедневный квест
type DailyQuest struct {
	ID        string
	Title     string
	Desc      string
	Goal      int  // Цель
	Progress  int  // Текущий прогресс
	Reward    int  // Награда (очки навыков)
	Completed bool // Выполнен ли квест
	Deadline  string
}

// QuestSystem управляет ежедневными квестами
type QuestSystem struct {
	Quests       []*DailyQuest
	DayStreak    int // Серия успешных дней
	TotalCompleted int
}

// NewSkillTree создает новое дерево навыков
func NewSkillTree() *SkillTree {
	tree := &SkillTree{
		Skills:      make(map[string]*Skill),
		SkillPoints: 0,
		TotalPoints: 0,
	}

	tree.initSkills()
	return tree
}

// initSkills инициализирует навыки
func (st *SkillTree) initSkills() {
	// Базовые навыки Go
	st.Skills["go_basics"] = &Skill{
		ID: "go_basics", Name: "Основы Go",
		Description: "Синтаксис, типы данных, функции",
		Level: 1, MaxLevel: 5, CostPerLevel: 1,
		BonusType: "knowledge", BonusValue: 5, Unlocked: true,
	}

	st.Skills["concurrency"] = &Skill{
		ID: "concurrency", Name: "Конкурентность",
		Description: "Горутины, каналы, sync package",
		Level: 0, MaxLevel: 5, CostPerLevel: 2,
		BonusType: "knowledge", BonusValue: 8, Unlocked: false,
		Prerequisites: []string{"go_basics"},
	}

	st.Skills["interfaces"] = &Skill{
		ID: "interfaces", Name: "Интерфейсы",
		Description: "Интерфейсы, полиморфизм",
		Level: 0, MaxLevel: 5, CostPerLevel: 2,
		BonusType: "knowledge", BonusValue: 7, Unlocked: false,
		Prerequisites: []string{"go_basics"},
	}

	// Навыки фокуса
	st.Skills["focus_master"] = &Skill{
		ID: "focus_master", Name: "Мастер Фокуса",
		Description: "Концентрация на задачах",
		Level: 0, MaxLevel: 5, CostPerLevel: 1,
		BonusType: "focus", BonusValue: 5, Unlocked: true,
	}

	st.Skills["meditation"] = &Skill{
		ID: "meditation", Name: "Медитация",
		Description: "Восстановление фокуса",
		Level: 0, MaxLevel: 5, CostPerLevel: 2,
		BonusType: "dopamine", BonusValue: 10, Unlocked: false,
		Prerequisites: []string{"focus_master"},
	}

	// Навыки силы воли
	st.Skills["willpower"] = &Skill{
		ID: "willpower", Name: "Сила Воли",
		Description: "Сопротивление искушениям",
		Level: 0, MaxLevel: 5, CostPerLevel: 2,
		BonusType: "willpower", BonusValue: 8, Unlocked: true,
	}

	st.Skills["discipline"] = &Skill{
		ID: "discipline", Name: "Дисциплина",
		Description: "Следование плану",
		Level: 0, MaxLevel: 5, CostPerLevel: 3,
		BonusType: "willpower", BonusValue: 10, Unlocked: false,
		Prerequisites: []string{"willpower"},
	}

	st.Skills["money_management"] = &Skill{
		ID: "money_management", Name: "Управление деньгами",
		Description: "Экономия и инвестиции",
		Level: 0, MaxLevel: 5, CostPerLevel: 2,
		BonusType: "money", BonusValue: 50, Unlocked: false,
		Prerequisites: []string{"willpower"},
	}

	// Продвинутые навыки
	st.Skills["web_frameworks"] = &Skill{
		ID: "web_frameworks", Name: "Web Фреймворки",
		Description: "Gin, Echo, Fiber",
		Level: 0, MaxLevel: 5, CostPerLevel: 3,
		BonusType: "knowledge", BonusValue: 10, Unlocked: false,
		Prerequisites: []string{"concurrency", "interfaces"},
	}

	st.Skills["database"] = &Skill{
		ID: "database", Name: "Базы данных",
		Description: "PostgreSQL, MongoDB, Redis",
		Level: 0, MaxLevel: 5, CostPerLevel: 3,
		BonusType: "knowledge", BonusValue: 10, Unlocked: false,
		Prerequisites: []string{"concurrency"},
	}

	st.Skills["microservices"] = &Skill{
		ID: "microservices", Name: "Микросервисы",
		Description: "gRPC, Docker, Kubernetes",
		Level: 0, MaxLevel: 5, CostPerLevel: 4,
		BonusType: "knowledge", BonusValue: 12, Unlocked: false,
		Prerequisites: []string{"web_frameworks", "database"},
	}

	st.Skills["clean_architecture"] = &Skill{
		ID: "clean_architecture", Name: "Чистая архитектура",
		Description: "SOLID, DDD, паттерны",
		Level: 0, MaxLevel: 5, CostPerLevel: 4,
		BonusType: "knowledge", BonusValue: 12, Unlocked: false,
		Prerequisites: []string{"interfaces", "web_frameworks"},
	}

	st.Skills["anti_procrastination"] = &Skill{
		ID: "anti_procrastination", Name: "Борьба с прокрастинацией",
		Description: "Pomodoro, тайм-менеджмент",
		Level: 0, MaxLevel: 5, CostPerLevel: 2,
		BonusType: "focus", BonusValue: 8, Unlocked: false,
		Prerequisites: []string{"focus_master"},
	}

	st.Skills["cold_shower"] = &Skill{
		ID: "cold_shower", Name: "Холодный душ",
		Description: "Закаливание и ритуалы",
		Level: 0, MaxLevel: 3, CostPerLevel: 1,
		BonusType: "willpower", BonusValue: 5, Unlocked: true,
	}
}

// Display отображает дерево навыков
func (st *SkillTree) Display() string {
	var sb strings.Builder

	sb.WriteString("🌳 <b>ДЕРЕВО НАВЫКОВ</b>\n")
	sb.WriteString("━━━━━━━━━━━━━━━━━━━━\n")
	sb.WriteString(fmt.Sprintf("✨ Очки навыков: %d (всего: %d)\n\n", st.SkillPoints, st.TotalPoints))

	// Группируем по категориям
	categories := map[string][]string{
		"📚 GO-НАВЫКИ": {"go_basics", "concurrency", "interfaces", "web_frameworks", "database", "microservices", "clean_architecture"},
		"🎯 ФОКУС":     {"focus_master", "meditation", "anti_procrastination"},
		"💪 СИЛА ВОЛИ": {"willpower", "discipline", "cold_shower"},
		"💰 ФИНАНСЫ":   {"money_management"},
	}

	for categoryName, skillIDs := range categories {
		sb.WriteString(fmt.Sprintf("<b>%s</b>\n", categoryName))
		sb.WriteString("────────────────────\n")

		for _, skillID := range skillIDs {
			skill := st.Skills[skillID]
			sb.WriteString(st.displaySkill(skill))
		}
		sb.WriteString("\n")
	}

	return sb.String()
}

// displaySkill отображает один навык
func (st *SkillTree) displaySkill(skill *Skill) string {
	status := "🔒"
	if skill.Unlocked {
		status = "✅"
	}

	levelBar := ""
	for i := 0; i < skill.MaxLevel; i++ {
		if i < skill.Level {
			levelBar += "█"
		} else {
			levelBar += "░"
		}
	}

	cost := ""
	if skill.Unlocked && skill.Level < skill.MaxLevel {
		cost = fmt.Sprintf(" [%d очк.]", skill.CostPerLevel)
	}

	return fmt.Sprintf("  %s %-18s [%s] Ур.%d/%d%s\n     %s\n",
		status, skill.Name, levelBar, skill.Level, skill.MaxLevel, cost, skill.Description)
}

// UpgradeSkill улучшает навык
func (st *SkillTree) UpgradeSkill(skillID string) bool {
	skill, exists := st.Skills[skillID]
	if !exists {
		return false
	}

	if !skill.Unlocked {
		return false
	}

	if skill.Level >= skill.MaxLevel {
		return false
	}

	cost := skill.CostPerLevel
	if st.SkillPoints < cost {
		return false
	}

	st.SkillPoints -= cost
	skill.Level++
	st.checkUnlocks(skillID)

	return true
}

// checkUnlocks проверяет разблокировку навыков
func (st *SkillTree) checkUnlocks(skillID string) {
	for _, skill := range st.Skills {
		if skill.Unlocked {
			continue
		}

		allMet := true
		for _, prereq := range skill.Prerequisites {
			prereqSkill := st.Skills[prereq]
			if prereqSkill == nil || prereqSkill.Level == 0 {
				allMet = false
				break
			}
		}

		if allMet {
			skill.Unlocked = true
		}
	}
}

// EarnSkillPoints начисляет очки навыков
func (st *SkillTree) EarnSkillPoints(points int) {
	// Валидация
	if points < 0 {
		log.Printf("⚠️  WARNING: Отрицательные очки навыков: %d, установлено 0", points)
		points = 0
	}

	st.SkillPoints = clampInt(st.SkillPoints + points, 0, 10000)
	st.TotalPoints = clampInt(st.TotalPoints + points, 0, 100000)
}

// GetTotalBonus возвращает суммарный бонус по типу
func (st *SkillTree) GetTotalBonus(bonusType string) int {
	total := 0
	for _, skill := range st.Skills {
		if skill.BonusType == bonusType {
			total += skill.Level * skill.BonusValue
		}
	}
	return total
}

// GetTotalBonuses возвращает все бонусы
func (st *SkillTree) GetTotalBonuses() map[string]int {
	bonuses := make(map[string]int)
	bonusTypes := []string{"focus", "willpower", "knowledge", "money", "dopamine"}

	for _, bt := range bonusTypes {
		bonuses[bt] = st.GetTotalBonus(bt)
	}

	return bonuses
}

// NewQuestSystem создает систему квестов
func NewQuestSystem() *QuestSystem {
	return &QuestSystem{
		Quests:       make([]*DailyQuest, 0),
		DayStreak:    0,
		TotalCompleted: 0,
	}
}

// GenerateDailyQuests генерирует ежедневные квесты
func (qs *QuestSystem) GenerateDailyQuests() {
	today := time.Now().Format("2006-01-02")

	qs.Quests = []*DailyQuest{
		{
			ID: "study_go_30min", Title: "30 минут Go",
			Desc: "Изучи новую тему Go в течение 30 минут",
			Goal: 30, Progress: 0, Reward: 2, Completed: false, Deadline: today,
		},
		{
			ID: "resist_temptation", Title: "Борец с искушениями",
			Desc: "Сопротивляйся 3 искушениям сегодня",
			Goal: 3, Progress: 0, Reward: 3, Completed: false, Deadline: today,
		},
		{
			ID: "code_practice", Title: "Практика кода",
			Desc: "Напиши 50+ строк кода на Go",
			Goal: 50, Progress: 0, Reward: 3, Completed: false, Deadline: today,
		},
		{
			ID: "morning_routine", Title: "Утренний ритуал",
			Desc: "Выполни утренний ритуал",
			Goal: 1, Progress: 0, Reward: 1, Completed: false, Deadline: today,
		},
		{
			ID: "no_social_media", Title: "Цифровой детокс",
			Desc: "Не заходи в соцсети 4 часа",
			Goal: 4, Progress: 0, Reward: 2, Completed: false, Deadline: today,
		},
	}
}

// UpdateQuestProgress обновляет прогресс квеста
func (qs *QuestSystem) UpdateQuestProgress(questID string, progress int) {
	// Валидация прогресса
	if progress < 0 {
		log.Printf("⚠️  WARNING: Отрицательный прогресс квеста %s: %d, установлено 0", questID, progress)
		progress = 0
	}

	for _, quest := range qs.Quests {
		if quest.ID == questID && !quest.Completed {
			quest.Progress = clampInt(quest.Progress + progress, 0, quest.Goal*10)
			if quest.Progress >= quest.Goal {
				quest.Completed = true
				qs.TotalCompleted++
			}
			return
		}
	}
}

// DisplayQuests отображает квесты
func (qs *QuestSystem) DisplayQuests() string {
	var sb strings.Builder

	sb.WriteString("📋 <b>ЕЖЕДНЕВНЫЕ КВЕСТЫ</b>\n")
	sb.WriteString("━━━━━━━━━━━━━━━━━━━━\n\n")

	if len(qs.Quests) == 0 {
		sb.WriteString("Нет активных квестов\n")
		return sb.String()
	}

	for _, quest := range qs.Quests {
		status := "⏳"
		if quest.Completed {
			status = "✅"
		}

		sb.WriteString(fmt.Sprintf("%s <b>%s</b>\n", status, quest.Title))
		sb.WriteString(fmt.Sprintf("   %s\n", quest.Desc))
		sb.WriteString(fmt.Sprintf("   Прогресс: %d/%d | Награда: %d очк.\n\n",
			quest.Progress, quest.Goal, quest.Reward))
	}

	sb.WriteString(fmt.Sprintf("🔥 Серия дней: %d\n", qs.DayStreak))
	sb.WriteString(fmt.Sprintf("🏆 Выполнено квестов: %d\n", qs.TotalCompleted))

	return sb.String()
}

// ClaimRewards забирает награды
func (qs *QuestSystem) ClaimRewards() int {
	totalReward := 0
	for _, quest := range qs.Quests {
		if quest.Completed && quest.Reward > 0 {
			totalReward += quest.Reward
			quest.Reward = 0
		}
	}
	return totalReward
}

// CheckDayStreak проверяет серию дней
func (qs *QuestSystem) CheckDayStreak(allCompleted bool) {
	if allCompleted {
		qs.DayStreak++
	} else {
		qs.DayStreak = 0
	}
}

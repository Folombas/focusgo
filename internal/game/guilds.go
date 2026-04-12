package game

import (
	"fmt"
	"time"
)

// Guild представляет гильдию игроков
type Guild struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	MemberCount int       `json:"member_count"`
	Level       int       `json:"level"`
	TotalXP     int       `json:"total_xp"`
	GuildPoints int       `json:"guild_points"`
}

// GuildMember представляет участника гильдии
type GuildMember struct {
	PlayerID   int64     `json:"player_id"`
	PlayerName string    `json:"player_name"`
	Role       string    `json:"role"` // leader, officer, member
	JoinedAt   time.Time `json:"joined_at"`
	ContributionXP int   `json:"contribution_xp"`
	ContributionDays int `json:"contribution_days"`
}

// GuildQuest представляет квест гильдии
type GuildQuest struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Target      int       `json:"target"`
	Progress    int       `json:"progress"`
	RewardXP    int       `json:"reward_xp"`
	RewardPoints int      `json:"reward_points"`
	Difficulty  int       `json:"difficulty"`
	ExpiresAt   time.Time `json:"expires_at"`
	CompletedBy []int64   `json:"completed_by"`
}

// GuildSystem управляет системой гильдий
type GuildSystem struct {
	Guild         *Guild        `json:"guild"`
	Members       []GuildMember `json:"members"`
	Quests        []GuildQuest  `json:"quests"`
	PlayerID      int64         `json:"player_id"`
	WeeklyProgress int          `json:"weekly_progress"`
}

// NewGuildSystem создаёт новую систему гильдий
func NewGuildSystem(playerID int64, playerName string) *GuildSystem {
	gs := &GuildSystem{
		PlayerID: playerID,
		Members:  make([]GuildMember, 0),
		Quests:   make([]GuildQuest, 0),
	}

	// Создаём гильдию для игрока
	gs.Guild = &Guild{
		ID:          fmt.Sprintf("guild_%d", playerID),
		Name:        fmt.Sprintf("%s's Guild", playerName),
		Description: "Гильдия целеустремлённых разработчиков",
		CreatedAt:   time.Now(),
		MemberCount: 1,
		Level:       1,
		TotalXP:     0,
		GuildPoints: 0,
	}

	// Добавляем игрока как лидера
	gs.Members = append(gs.Members, GuildMember{
		PlayerID:         playerID,
		PlayerName:       playerName,
		Role:             "leader",
		JoinedAt:         time.Now(),
		ContributionXP:   0,
		ContributionDays: 0,
	})

	// Генерируем квесты гильдии
	gs.GenerateGuildQuests()

	return gs
}

// GenerateGuildQuests генерирует квесты гильдии
func (gs *GuildSystem) GenerateGuildQuests() {
	gs.Quests = make([]GuildQuest, 0)
	now := time.Now()
	expiresAt := now.Add(7 * 24 * time.Hour) // Неделю

	gs.Quests = append(gs.Quests, GuildQuest{
		ID:          "guild_study_500min",
		Name:        "📚 Командное обучение",
		Description: "Суммарно изучите Go 500 минут",
		Target:      500,
		Progress:    0,
		RewardXP:    300,
		RewardPoints: 150,
		Difficulty:  2,
		ExpiresAt:   expiresAt,
	})

	gs.Quests = append(gs.Quests, GuildQuest{
		ID:          "guild_resist_50",
		Name:        "💪 Стальная гильдия",
		Description: "Сопротивитесь 50 искушениям суммарно",
		Target:      50,
		Progress:    0,
		RewardXP:    400,
		RewardPoints: 200,
		Difficulty:  3,
		ExpiresAt:   expiresAt,
	})

	gs.Quests = append(gs.Quests, GuildQuest{
		ID:          "guild_quiz_100",
		Name:        "🧩 Команда знатоков",
		Description: "Ответьте правильно на 100 вопросов",
		Target:      100,
		Progress:    0,
		RewardXP:    500,
		RewardPoints: 250,
		Difficulty:  3,
		ExpiresAt:   expiresAt,
	})
}

// ContributeXP добавляет вклад игрока в гильдию
func (gs *GuildSystem) ContributeXP(playerID int64, xp int) string {
	for i := range gs.Members {
		if gs.Members[i].PlayerID == playerID {
			gs.Members[i].ContributionXP += xp
			gs.Guild.TotalXP += xp
			gs.Guild.GuildPoints += xp / 10

			// Проверяем уровень гильдии
			gs.checkGuildLevel()

			return fmt.Sprintf("✨ Внесён вклад: %d XP\n🏆 Гильдия: %d очков", xp, gs.Guild.GuildPoints)
		}
	}
	return "❌ Участник не найден"
}

// checkGuildLevel проверяет и обновляет уровень гильдии
func (gs *GuildSystem) checkGuildLevel() {
	requiredPoints := gs.Guild.Level * 500
	if gs.Guild.GuildPoints >= requiredPoints {
		gs.Guild.Level++
		gs.Guild.MemberCount++
	}
}

// UpdateQuestProgress обновляет прогресс квеста гильдии
func (gs *GuildSystem) UpdateQuestProgress(questID string, amount int) (completed bool, rewardMsg string) {
	for i := range gs.Quests {
		if gs.Quests[i].ID == questID {
			gs.Quests[i].Progress += amount

			if gs.Quests[i].Progress >= gs.Quests[i].Target && gs.Quests[i].Progress-amount < gs.Quests[i].Target {
				gs.Quests[i].CompletedBy = append(gs.Quests[i].CompletedBy, gs.PlayerID)
				completed = true
				rewardMsg = fmt.Sprintf(
					"🎉 Квест гильдии \"%s\" выполнен!\n✨ +%d XP\n🏆 +%d очков гильдии",
					gs.Quests[i].Name,
					gs.Quests[i].RewardXP,
					gs.Quests[i].RewardPoints,
				)
			}
			return
		}
	}
	return false, ""
}

// Display возвращает текстовое представление гильдии
func (gs *GuildSystem) Display() string {
	guild := gs.Guild

	text := fmt.Sprintf("🏆 <b>ГИЛЬДИЯ: %s</b>\n━━━━━━━━━━━━━━━━━━━━\n\n"+
		"📊 Уровень гильдии: %d\n"+
		"👥 Участников: %d\n"+
		"✨ Суммарно XP: %d\n"+
		"🏆 Очков гильдии: %d\n\n",
		guild.Name,
		guild.Level,
		guild.MemberCount,
		guild.TotalXP,
		guild.GuildPoints)

	// Квесты гильдии
	text += "<b>📋 КВЕСТЫ ГИЛЬДИИ</b>\n\n"
	for i, quest := range gs.Quests {
		status := "⬜"
		if quest.Progress >= quest.Target {
			status = "✅"
		}

		difficulty := ""
		switch quest.Difficulty {
		case 1:
			difficulty = "🟢"
		case 2:
			difficulty = "🟡"
		case 3:
			difficulty = "🔴"
		}

		progress := quest.Progress
		if progress > quest.Target {
			progress = quest.Target
		}

		daysLeft := int(time.Until(quest.ExpiresAt).Hours() / 24)

		text += fmt.Sprintf("%d. %s %s %s\n", i+1, status, difficulty, quest.Name)
		text += fmt.Sprintf("   %s\n", quest.Description)
		text += fmt.Sprintf("   Прогресс: %d/%d\n", progress, quest.Target)
		text += fmt.Sprintf("   Награда: ✨%d XP | 🏆%d очков\n", quest.RewardXP, quest.RewardPoints)
		text += fmt.Sprintf("   ⏰ Осталось дней: %d\n\n", daysLeft)
	}

	// Топ участников
	text += "<b>🏅 ТОП УЧАСТНИКОВ</b>\n\n"
	for i, member := range gs.Members {
		if i >= 5 {
			break
		}

		role := ""
		switch member.Role {
		case "leader":
			role = "👑"
		case "officer":
			role = "⭐"
		default:
			role = "👤"
		}

		text += fmt.Sprintf("%s %s: %d XP вклад\n", role, member.PlayerName, member.ContributionXP)
	}

	text += "\n━━━━━━━━━━━━━━━━━━━━\n💡 Вноси вклад в гильдию, чтобы получать бонусы!"

	return text
}

// GetWeeklyReport возвращает отчёт за неделю
func (gs *GuildSystem) GetWeeklyReport() string {
	totalContrib := 0
	for _, member := range gs.Members {
		totalContrib += member.ContributionXP
	}

	return fmt.Sprintf(
		"📊 <b>НЕДЕЛЬНЫЙ ОТЧЁТ ГИЛЬДИИ</b>\n\n"+
			"✨ Общий вклад: %d XP\n"+
			"🏆 Уровень гильдии: %d\n"+
			"📋 Выполнено квестов: %d/%d\n\n"+
			"💪 Отличная работа, команда!",
		totalContrib,
		gs.Guild.Level,
		gs.getCompletedQuestsCount(),
		len(gs.Quests),
	)
}

func (gs *GuildSystem) getCompletedQuestsCount() int {
	count := 0
	for _, quest := range gs.Quests {
		if quest.Progress >= quest.Target {
			count++
		}
	}
	return count
}

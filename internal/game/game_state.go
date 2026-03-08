package game

import (
	"database/sql"
	"fmt"
	"math/rand"
	"time"

	"focusgo/internal/database"
)

// GameState представляет состояние игры игрока
type GameState struct {
	ChatID       int64
	Name         string
	Level        int
	Experience   int
	NextLevelXP  int
	GoKnowledge  int
	Focus        int
	Willpower    int
	Money        int
	Dopamine     int
	PlayTime     int // минуты
	DaysPlayed   int
	CurrentDay   int
	CurrentHour  int // 8-23
	IsPlaying    bool
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

// NewGameState создаёт новое состояние игры
func NewGameState(chatID int64, name string) *GameState {
	return &GameState{
		ChatID:      chatID,
		Name:        name,
		Level:       1,
		Experience:  0,
		NextLevelXP: 100,
		GoKnowledge: 40,
		Focus:       70,
		Willpower:   65,
		Money:       500,
		Dopamine:    200,
		PlayTime:    0,
		DaysPlayed:  1,
		CurrentDay:  1,
		CurrentHour: 8,
		IsPlaying:   true,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
}

// LoadGameState загружает состояние игры из БД
func LoadGameState(chatID int64) (*GameState, error) {
	query := `
		SELECT chat_id, name, level, experience, go_knowledge, focus, willpower,
		       money, dopamine, play_time, days_played, current_day, hour, game_active,
		       created_at, updated_at
		FROM players
		WHERE chat_id = ?
	`

	row := database.DB.QueryRow(query, chatID)

	state := &GameState{}
	var gameActive int
	var createdAt, updatedAt time.Time

	err := row.Scan(
		&state.ChatID, &state.Name, &state.Level, &state.Experience,
		&state.GoKnowledge, &state.Focus, &state.Willpower, &state.Money,
		&state.Dopamine, &state.PlayTime, &state.DaysPlayed, &state.CurrentDay,
		&state.CurrentHour, &gameActive, &createdAt, &updatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	state.IsPlaying = gameActive == 1
	state.CreatedAt = createdAt
	state.UpdatedAt = updatedAt
	state.NextLevelXP = state.Level * 100

	return state, nil
}

// SaveGameState сохраняет состояние игры в БД
func (s *GameState) SaveGameState() error {
	query := `
		INSERT OR REPLACE INTO players 
		(chat_id, name, level, experience, go_knowledge, focus, willpower,
		 money, dopamine, play_time, days_played, current_day, hour, game_active,
		 created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	s.UpdatedAt = time.Now()

	_, err := database.DB.Exec(query,
		s.ChatID, s.Name, s.Level, s.Experience, s.GoKnowledge,
		s.Focus, s.Willpower, s.Money, s.Dopamine, s.PlayTime,
		s.DaysPlayed, s.CurrentDay, s.CurrentHour, s.IsPlaying,
		s.CreatedAt, s.UpdatedAt,
	)

	return err
}

// AddExperience добавляет опыт и проверяет повышение уровня
func (s *GameState) AddExperience(xp int) int {
	if xp < 0 {
		xp = 0
	}

	s.Experience += xp
	levelsGained := 0

	for s.Experience >= s.NextLevelXP {
		s.Experience -= s.NextLevelXP
		s.Level++
		s.NextLevelXP = s.Level * 100
		levelsGained++

		// Бонус за уровень
		s.Focus = 100
		s.Willpower = 100
	}

	return levelsGained
}

// StudyGo изучает Go
func (s *GameState) StudyGo(minutes int) (string, int, int) {
	if minutes <= 0 {
		minutes = 30
	}

	// Базовый опыт
	baseXP := minutes / 2
	knowledgeGained := minutes / 5

	s.AddExperience(baseXP)
	s.GoKnowledge += knowledgeGained
	if s.GoKnowledge > 100 {
		s.GoKnowledge = 100
	}

	// Прогресс времени
	s.PlayTime += minutes
	s.CurrentHour += minutes / 60
	if s.CurrentHour >= 24 {
		s.CurrentHour = 8
		s.CurrentDay++
	}

	s.Dopamine += minutes / 3
	if s.Dopamine > 500 {
		s.Dopamine = 500
	}

	msg := fmt.Sprintf("📚 <b>ИЗУЧЕНИЕ GO: %d минут</b>\n\n✨ +%d опыта\n🧠 +%d к знанию Go",
		minutes, baseXP, knowledgeGained)

	return msg, baseXP, knowledgeGained
}

// Rest отдыхает
func (s *GameState) Rest(minutes int) string {
	if minutes <= 0 {
		minutes = 15
	}

	focusRecovered := minutes / 2
	dopamineRecovered := minutes / 3

	s.Focus += focusRecovered
	if s.Focus > 100 {
		s.Focus = 100
	}

	s.Dopamine += dopamineRecovered
	if s.Dopamine > 500 {
		s.Dopamine = 500
	}

	s.PlayTime += minutes
	s.CurrentHour += minutes / 60
	if s.CurrentHour >= 24 {
		s.CurrentHour = 8
		s.CurrentDay++
	}

	msg := fmt.Sprintf("💤 <b>ОТДЫХ: %d минут</b>\n\n😌 Фокус + %d → %d%%\n✨ Дофамин + %d → %d",
		minutes, focusRecovered, s.Focus, dopamineRecovered, s.Dopamine)

	return msg
}

// CheckTemptation проверяет искушение
func (s *GameState) CheckTemptation() bool {
	return rand.Intn(100) < 35 // 35% шанс
}

// ResistTemptation сопротивляется искушению
func (s *GameState) ResistTemptation(power int) (string, int) {
	resistChance := s.Willpower - power + 50
	if resistChance > 90 {
		resistChance = 90
	}
	if resistChance < 20 {
		resistChance = 20
	}

	roll := rand.Intn(100)

	if roll < resistChance {
		// Победа
		xpReward := power / 2
		s.AddExperience(xpReward)
		s.Focus += 10
		if s.Focus > 100 {
			s.Focus = 100
		}
		s.Willpower += 5
		if s.Willpower > 100 {
			s.Willpower = 100
		}
		s.Dopamine += 50

		msg := fmt.Sprintf("✅ <b>СОПРОТИВЛЕНИЕ!</b>\n\nВы успешно сопротивлялись искушению!\n\n✨ +%d опыта\n🎯 Фокус +10 → %d%%\n💪 Сила воли +5 → %d%%",
			xpReward, s.Focus, s.Willpower)

		return msg, xpReward
	} else {
		// Поражение
		if s.Experience >= 20 {
			s.Experience -= 20
		}
		s.Focus -= 20
		if s.Focus < 0 {
			s.Focus = 0
		}
		s.Willpower -= 10
		if s.Willpower < 0 {
			s.Willpower = 0
		}
		s.Dopamine -= 100
		if s.Dopamine < 0 {
			s.Dopamine = 0
		}

		msg := fmt.Sprintf("❌ <b>ПОРАЖЕНИЕ...</b>\n\nВы поддались искушению.\n\n💀 -20 опыта\n🎯 Фокус -20 → %d%%\n💪 Сила воли -10 → %d%%",
			s.Focus, s.Willpower)

		return msg, 0
	}
}

// FinalBattle финальная битва с боссом
func (s *GameState) FinalBattle(bossName string, bossPower int) (bool, string) {
	successChance := (s.Willpower*2 + s.Focus) / 3
	if successChance > 95 {
		successChance = 95
	}
	if successChance < 20 {
		successChance = 20
	}

	roll := rand.Intn(100)

	if roll < successChance {
		// Победа
		s.Focus = 100
		s.Willpower = 100
		s.Dopamine += 500
		if s.Dopamine > 1000 {
			s.Dopamine = 1000
		}
		s.AddExperience(200)

		// Завершение дня
		s.CurrentDay++
		s.CurrentHour = 8

		msg := fmt.Sprintf(`🎉 <b>ПОБЕДА!</b>

Вы победили %s и успешно завершили день!

✨ +200 опыта
🎯 Фокус восстановлен: 100%%
💪 Сила воли восстановлена: 100%%
✨ Дофамин +500

🌅 Новый день начался!`, bossName)

		return true, msg
	} else {
		// Поражение
		s.Focus = 30
		s.Willpower = 40
		s.Dopamine -= 300
		if s.Dopamine < 0 {
			s.Dopamine = 0
		}

		// Завершение дня
		s.CurrentDay++
		s.CurrentHour = 8

		msg := fmt.Sprintf(`💔 <b>ПОРАЖЕНИЕ...</b>

%s оказался сильнее...

🎯 Фокус: 30%%
💪 Сила воли: 40%%
✨ Дофамин -300

🌅 Новый день начался!
Не сдавайся!`, bossName)

		return false, msg
	}
}

// GetStatus возвращает строку статуса
func (s *GameState) GetStatus() string {
	return fmt.Sprintf(`👤 <b>%s</b>
━━━━━━━━━━━━━━━━━━━━
🏆 Уровень: %d (Опыт: %d/%d)
📚 Знание Go: %d/100
🎯 Фокус: %d%%
💪 Сила воли: %d%%
💰 Деньги: %d₽
✨ Дофамин: %d

📅 День: %d | ⏰ %02d:00`,
		s.Name, s.Level, s.Experience, s.NextLevelXP,
		s.GoKnowledge, s.Focus, s.Willpower, s.Money, s.Dopamine,
		s.CurrentDay, s.CurrentHour)
}

// GetRating возвращает рейтинг игрока
func (s *GameState) GetRating() string {
	rating := s.GoKnowledge*10 + s.Focus*5 + s.Willpower*3 + s.Dopamine/10 + s.Level*100

	if rating < 500 {
		return "🌱 Начинающий гофер"
	} else if rating < 1500 {
		return "🌿 Ученик разработчика"
	} else if rating < 3000 {
		return "🌳 Junior Go Developer"
	} else if rating < 5000 {
		return "🏢 Middle Go Developer"
	} else {
		return "🚀 Senior Go Master"
	}
}

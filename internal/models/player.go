package models

import (
	"fmt"
	"log"
	"math/rand"
	"time"
)

// Player представляет игрока в Telegram-версии
type Player struct {
	// Основные характеристики
	Name        string
	ChatID      int64
	Level       int
	Experience  int
	GoKnowledge int
	Focus       int
	Willpower   int
	Money       int
	Dopamine    int

	// Прогресс
	PlayTime   int // в минутах
	DaysPlayed int

	// Списки
	Temptations  []string
	Achievements []string

	// Система навыков и квестов
	SkillTree *SkillTree
	Quests    *QuestSystem

	// Состояние игры
	GameActive bool
	CurrentDay int
	Hour       int // Текущий игровой час (8-24)
}

// NewPlayer создает нового игрока
func NewPlayer(chatID int64, name string) *Player {
	rand.Seed(time.Now().UnixNano())

	player := &Player{
		Name:         name,
		ChatID:       chatID,
		Level:        1,
		Experience:   0,
		GoKnowledge:  40,
		Focus:        70,
		Willpower:    65,
		Money:        500,
		Dopamine:     200,
		PlayTime:     0,
		DaysPlayed:   1,
		CurrentDay:   1,
		Hour:         8, // Начинаем в 8 утра
		Temptations:  make([]string, 0),
		Achievements: make([]string, 0),
		GameActive:   true,
	}

	// Инициализируем системы
	player.SkillTree = NewSkillTree()
	player.Quests = NewQuestSystem()

	// Валидация после создания
	if err := validatePlayer(player); err != nil {
		log.Printf("⚠️  WARNING: Ошибки валидации нового игрока: %v", err)
	}

	// Применяем начальные бонусы
	player.ApplySkillBonuses()

	return player
}

// DisplayStatus возвращает строку со статусом игрока
func (p *Player) DisplayStatus() string {
	return fmt.Sprintf(`👤 <b>%s</b>
━━━━━━━━━━━━━━━━━━━━
🏆 Уровень: %d (Опыт: %d/%d)
📚 Знание Go: %d/100
🎯 Фокус: %d%%
💪 Сила воли: %d%%
💰 Деньги: %d₽
✨ Дофамин: %d

📅 День: %d | ⏰ %02d:00`,
		p.Name,
		p.Level,
		p.Experience,
		p.Level*100,
		p.GoKnowledge,
		p.Focus,
		p.Willpower,
		p.Money,
		p.Dopamine,
		p.CurrentDay,
		p.Hour,
	)
}

// DisplayProfile возвращает подробный профиль
func (p *Player) DisplayProfile() string {
	bonuses := p.SkillTree.GetTotalBonuses()

	text := fmt.Sprintf(`👤 <b>ПРОФИЛЬ ИГРОКА</b>
━━━━━━━━━━━━━━━━━━━━

📛 Имя: %s
🏆 Уровень: %d
⭐ Опыт: %d/%d

📊 <b>ХАРАКТЕРИСТИКИ:</b>
📚 Знание Go: %d/100 (бонус: +%d)
🎯 Фокус: %d%% (бонус: +%d)
💪 Сила воли: %d%% (бонус: +%d)
💰 Деньги: %d₽ (бонус: +%d)
✨ Дофамин: %d (бонус: +%d)

📈 <b>ПРОГРЕСС:</b>
🕐 Время в игре: %d мин
📅 Дней сыграно: %d
🔥 Серия дней: %d

🎯 Выполнено квестов: %d
💪 Преодолено искушений: %d
🏆 Достижений: %d`,
		p.Name,
		p.Level,
		p.Experience,
		p.Level*100,
		p.GoKnowledge, bonuses["knowledge"],
		p.Focus, bonuses["focus"],
		p.Willpower, bonuses["willpower"],
		p.Money, bonuses["money"],
		p.Dopamine, bonuses["dopamine"],
		p.PlayTime,
		p.DaysPlayed,
		p.Quests.DayStreak,
		p.Quests.TotalCompleted,
		len(p.Temptations),
		len(p.Achievements),
	)

	return text
}

// AddExperience добавляет опыт и проверяет повышение уровня
func (p *Player) AddExperience(xp int) int {
	// Валидация XP
	if xp < 0 {
		log.Printf("⚠️  WARNING: Отрицательный опыт: %d, установлено 0", xp)
		xp = 0
	}

	p.Experience = clampExperience(p.Experience + xp)

	// Проверяем повышение уровня
	xpNeeded := p.Level * 100
	levelsGained := 0

	for p.Experience >= xpNeeded {
		p.Experience -= xpNeeded
		p.Level = clampLevel(p.Level + 1)
		levelsGained++

		// Восстанавливаем характеристики
		p.Focus = clampStat(100)
		p.Willpower = clampStat(100)

		// Начисляем очки навыков
		skillPoints := 2 + (p.Level / 5)
		p.SkillTree.EarnSkillPoints(skillPoints)

		// Добавляем достижение
		p.Achievements = append(p.Achievements, fmt.Sprintf("Достигнут уровень %d", p.Level))

		// Особые достижения
		if p.Level == 5 {
			p.Achievements = append(p.Achievements, "🏆 Go-Новичок: Уровень 5")
		}
		if p.Level == 10 {
			p.Achievements = append(p.Achievements, "🏆 Go-Разработчик: Уровень 10")
		}
		if p.Level == 20 {
			p.Achievements = append(p.Achievements, "🏆 Go-Мастер: Уровень 20")
		}
		if p.Level == 30 {
			p.Achievements = append(p.Achievements, "🏆 Go-Легенда: Уровень 30")
		}
	}

	return levelsGained
}

// StudyGo изучает Go
func (p *Player) StudyGo(minutes int) string {
	// Валидация минут
	if minutes < 0 {
		log.Printf("⚠️  WARNING: Отрицательные минуты изучения: %d, установлено 0", minutes)
		minutes = 0
	}

	baseXP := minutes / 2
	knowledgeBonus := p.SkillTree.GetTotalBonus("knowledge")
	totalXP := baseXP + (knowledgeBonus / 5)

	p.AddExperience(totalXP)
	p.GoKnowledge = clampStat(p.GoKnowledge + minutes/5)

	// Обновляем квесты
	p.Quests.UpdateQuestProgress("study_go_30min", minutes)
	p.Quests.UpdateQuestProgress("code_practice", minutes/2)

	// Восстанавливаем дофамин
	p.Dopamine = clampDopamine(p.Dopamine + minutes/3)

	return fmt.Sprintf(`📚 <b>ИЗУЧЕНИЕ GO: %d минут</b>

✨ +%d опыта
🧠 +%d к знанию Go (всего: %d/100)
✨ +%d дофамина`,
		minutes, totalXP, minutes/5, p.GoKnowledge, minutes/3)
}

// Rest отдыхает
func (p *Player) Rest(minutes int) string {
	// Валидация минут
	if minutes < 0 {
		log.Printf("⚠️  WARNING: Отрицательные минуты отдыха: %d, установлено 0", minutes)
		minutes = 0
	}

	p.Focus = clampStat(p.Focus + minutes/2)
	p.Dopamine = clampDopamine(p.Dopamine + minutes/3)

	return fmt.Sprintf(`💤 <b>ОТДЫХ: %d минут</b>

😌 Фокус восстановлен: %d%%
✨ Дофамин: %d`,
		minutes, p.Focus, p.Dopamine)
}

// HandleTemptation обрабатывает искушение
func (p *Player) HandleTemptation(t Temptation) string {
	resistChance := p.Willpower - t.Power + 50
	resistChance = clampInt(resistChance, 10, 100)

	roll := rand.Intn(100)

	if roll < resistChance {
		// Успешное сопротивление
		xpReward := t.Power / 2
		p.AddExperience(xpReward)
		p.Focus = clampStat(p.Focus + 10)
		p.Willpower = clampStat(p.Willpower + 5)
		p.Dopamine = clampDopamine(p.Dopamine + 50)
		p.Temptations = append(p.Temptations, t.Name)

		// Обновляем квест
		p.Quests.UpdateQuestProgress("resist_temptation", 1)

		return fmt.Sprintf(`✅ <b>СОПРОТИВЛЕНИЕ!</b>

Вы успешно сопротивлялись искушению "%s"!

✨ +%d опыта
🎯 Фокус +10 → %d%%
💪 Сила воли +5 → %d%%
✨ Дофамин +50`,
			t.Name, xpReward, p.Focus, p.Willpower)
	} else {
		// Поражение
		xpLoss := t.XPLoss
		if p.Experience >= xpLoss {
			p.Experience = clampExperience(p.Experience - xpLoss)
		}

		p.Focus = clampStat(p.Focus - 20)
		p.Willpower = clampStat(p.Willpower - 10)
		p.Dopamine = clampDopamine(p.Dopamine - 100)

		return fmt.Sprintf(`❌ <b>ПОРАЖЕНИЕ...</b>

Вы поддались искушению "%s".

💀 -%d опыта
🎯 Фокус -20 → %d%%
💪 Сила воли -10 → %d%%
✨ Дофамин -100`,
			t.Name, xpLoss, p.Focus, p.Willpower)
	}
}

// FinalBattle финальная битва с боссом
func (p *Player) FinalBattle(boss Temptation) bool {
	successChance := (p.Willpower*2 + p.Focus) / 3
	successChance = clampInt(successChance, 10, 95)

	roll := rand.Intn(100)

	if roll < successChance {
		// Победа
		p.Focus = clampStat(100)
		p.Willpower = clampStat(100)
		p.Dopamine = clampDopamine(p.Dopamine + 500)
		p.Achievements = append(p.Achievements, "Победил финальное искушение")
		p.AddExperience(200)
		return true
	} else {
		// Поражение
		p.Focus = clampStat(30)
		p.Willpower = clampStat(40)
		p.Dopamine = clampDopamine(p.Dopamine - 300)
		return false
	}
}

// CalculateScore вычисляет счет
func (p *Player) CalculateScore() int {
	return p.GoKnowledge*10 + p.Focus*5 + p.Willpower*3 + p.Dopamine/10
}

// ApplySkillBonuses применяет бонусы от навыков
func (p *Player) ApplySkillBonuses() {
	bonuses := p.SkillTree.GetTotalBonuses()

	p.Focus = clampStat(p.Focus + bonuses["focus"])
	p.Willpower = clampStat(p.Willpower + bonuses["willpower"])
	p.GoKnowledge = clampStat(p.GoKnowledge + bonuses["knowledge"])
	p.Money = clampMoney(p.Money + bonuses["money"])
	p.Dopamine = clampDopamine(p.Dopamine + bonuses["dopamine"])
}

// GetRating возвращает рейтинг игрока
func (p *Player) GetRating() string {
	rating := p.CalculateScore() + (p.Level * 100) + (len(p.Achievements) * 50)

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

// DisplayStatistics возвращает расширенную статистику
func (p *Player) DisplayStatistics() string {
	rating := p.CalculateScore() + (p.Level*100) + (len(p.Achievements) * 50)

	return fmt.Sprintf(`📊 <b>СТАТИСТИКА ИГРЫ</b>
━━━━━━━━━━━━━━━━━━━━

🏆 Уровень: %d
⭐ Опыт: %d/%d
🕐 Время в игре: %d минут
📅 Дней сыграно: %d
💪 Преодолено искушений: %d
🏆 Достижений: %d
🎯 Выполнено квестов: %d
🔥 Текущая серия: %d дней
✨ Очков навыков: %d (всего: %d)

🏅 ОБЩИЙ РЕЙТИНГ: %d
Ранг: %s`,
		p.Level,
		p.Experience,
		p.Level*100,
		p.PlayTime,
		p.DaysPlayed,
		len(p.Temptations),
		len(p.Achievements),
		p.Quests.TotalCompleted,
		p.Quests.DayStreak,
		p.SkillTree.SkillPoints,
		p.SkillTree.TotalPoints,
		rating,
		p.GetRating(),
	)
}

package game

import (
	"fmt"
	"math/rand"
	"time"
)

// DailyChallenge представляет ежедневный челлендж
type DailyChallenge struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Difficulty  int       `json:"difficulty"`  // 1-3 (easy, medium, hard)
	Target      int       `json:"target"`      // Целевое значение
	Progress    int       `json:"progress"`    // Текущий прогресс
	RewardXP    int       `json:"reward_xp"`   // Награда опыта
	RewardCoins int       `json:"reward_coins"` // Награда монет
	Completed   bool      `json:"completed"`
	Date        time.Time `json:"date"`
}

// ChallengeType определяет тип челленджа
type ChallengeType string

const (
	ChallengeStudyTime     ChallengeType = "study_time"     // Изучать Go X минут
	ChallengeResistTempts  ChallengeType = "resist_tempts"  // Сопротивиться X искушениям
	ChallengeCodeLines     ChallengeType = "code_lines"     // Написать X строк кода
	ChallengeQuizScore     ChallengeType = "quiz_score"      // Набрать X очков в викторине
	ChallengeStudyStreak   ChallengeType = "study_streak"   // Играть X дней подряд
	ChallengeSkillUpgrade  ChallengeType = "skill_upgrade"  // Улучшить навык X раз
	ChallengeEarlyBird     ChallengeType = "early_bird"     // Начать утром
	ChallengeNightOwl      ChallengeType = "night_owl"      // Играть ночью
	ChallengePerfectDay    ChallengeType = "perfect_day"    // Не проиграть боссу
	ChallengeComboMaster   ChallengeType = "combo_master"   // Достичь X комбо
)

// DailyChallengeSystem управляет ежедневными челленджами
type DailyChallengeSystem struct {
	Challenges    []DailyChallenge `json:"challenges"`
	CurrentDate   time.Time        `json:"current_date"`
	PlayerID      int64            `json:"player_id"`
	StreakDays    int              `json:"streak_days"`
	LastChallenge time.Time        `json:"last_challenge"`
}

// NewDailyChallengeSystem создаёт новую систему челленджей
func NewDailyChallengeSystem(playerID int64) *DailyChallengeSystem {
	return &DailyChallengeSystem{
		Challenges:  make([]DailyChallenge, 0),
		CurrentDate: time.Now(),
		PlayerID:    playerID,
	}
}

// GenerateDailyChallenges генерирует 3 ежедневных челленджа
func (dcs *DailyChallengeSystem) GenerateDailyChallenges() []DailyChallenge {
	// Проверяем, нужно ли генерировать новые
	if len(dcs.Challenges) > 0 && dcs.CurrentDate.Day() == time.Now().Day() {
		return dcs.Challenges
	}

	dcs.Challenges = make([]DailyChallenge, 0)
	dcs.CurrentDate = time.Now()

	// Пул челленджей
	challengePool := dcs.getChallengePool()

	// Выбираем 3 случайных челленджа
	rand.Seed(time.Now().UnixNano())
	perm := rand.Perm(len(challengePool))

	for i := 0; i < 3 && i < len(perm); i++ {
		challenge := challengePool[perm[i]]
		challenge.Date = dcs.CurrentDate
		challenge.Progress = 0
		challenge.Completed = false
		dcs.Challenges = append(dcs.Challenges, challenge)
	}

	return dcs.Challenges
}

// getChallengePool возвращает пул доступных челленджей
func (dcs *DailyChallengeSystem) getChallengePool() []DailyChallenge {
	return []DailyChallenge{
		{
			ID:          "study_30min",
			Name:        "📚 30 минут Go",
			Description: "Изучай Go минимум 30 минут сегодня",
			Difficulty:  1,
			Target:      30,
			RewardXP:    100,
			RewardCoins: 50,
		},
		{
			ID:          "study_60min",
			Name:        "📚 Час кодирования",
			Description: "Потрать 60 минут на изучение Go",
			Difficulty:  2,
			Target:      60,
			RewardXP:    250,
			RewardCoins: 100,
		},
		{
			ID:          "study_120min",
			Name:        "📚 Марафон Go",
			Description: "120 минут изучения Go!",
			Difficulty:  3,
			Target:      120,
			RewardXP:    500,
			RewardCoins: 200,
		},
		{
			ID:          "resist_5",
			Name:        "💪 Борец с искушениями",
			Description: "Сопротивись 5 искушениям",
			Difficulty:  1,
			Target:      5,
			RewardXP:    150,
			RewardCoins: 75,
		},
		{
			ID:          "resist_10",
			Name:        "💪 Стальная воля",
			Description: "Сопротивись 10 искушениям",
			Difficulty:  2,
			Target:      10,
			RewardXP:    300,
			RewardCoins: 150,
		},
		{
			ID:          "quiz_5_correct",
			Name:        "🧩 Знаток Go",
			Description: "Ответь правильно на 5 вопросов в викторине",
			Difficulty:  1,
			Target:      5,
			RewardXP:    100,
			RewardCoins: 50,
		},
		{
			ID:          "quiz_10_correct",
			Name:        "🧩 Гуру викторин",
			Description: "Ответь правильно на 10 вопросов",
			Difficulty:  2,
			Target:      10,
			RewardXP:    250,
			RewardCoins: 100,
		},
		{
			ID:          "upgrade_2_skills",
			Name:        "🌳 Улучшатель",
			Description: "Улучши 2 навыка",
			Difficulty:  1,
			Target:      2,
			RewardXP:    150,
			RewardCoins: 75,
		},
		{
			ID:          "early_bird",
			Name:        "🌅 Ранняя пташка",
			Description: "Начни игру до 9 утра",
			Difficulty:  1,
			Target:      1,
			RewardXP:    100,
			RewardCoins: 50,
		},
		{
			ID:          "boss_victory",
			Name:        "👹 Победитель боссов",
			Description: "Победи босса сегодня",
			Difficulty:  2,
			Target:      1,
			RewardXP:    200,
			RewardCoins: 100,
		},
		{
			ID:          "no_temptations",
			Name:        "✨ Чистый день",
			Description: "Пройди день без проигрыша искушениям",
			Difficulty:  3,
			Target:      1,
			RewardXP:    400,
			RewardCoins: 200,
		},
		{
			ID:          "code_100_lines",
			Name:        "💻 Сотня строк",
			Description: "Напиши 100 строк кода на Go",
			Difficulty:  2,
			Target:      100,
			RewardXP:    300,
			RewardCoins: 150,
		},
	}
}

// UpdateProgress обновляет прогресс челленджа
func (dcs *DailyChallengeSystem) UpdateProgress(challengeID string, amount int) (completed bool, rewards string) {
	for i := range dcs.Challenges {
		if dcs.Challenges[i].ID == challengeID && !dcs.Challenges[i].Completed {
			dcs.Challenges[i].Progress += amount

			if dcs.Challenges[i].Progress >= dcs.Challenges[i].Target {
				dcs.Challenges[i].Completed = true
				completed = true
				rewards = fmt.Sprintf(
					"✅ Челлендж \"%s\" выполнен!\n✨ +%d XP\n💰 +%d монет",
					dcs.Challenges[i].Name,
					dcs.Challenges[i].RewardXP,
					dcs.Challenges[i].RewardCoins,
				)
			}
			return
		}
	}
	return false, ""
}

// GetTotalRewards возвращает суммарные награды за день
func (dcs *DailyChallengeSystem) GetTotalRewards() (totalXP, totalCoins int) {
	for _, c := range dcs.Challenges {
		if c.Completed {
			totalXP += c.RewardXP
			totalCoins += c.RewardCoins
		}
	}
	return
}

// GetCompletionCount возвращает количество выполненных челленджей
func (dcs *DailyChallengeSystem) GetCompletionCount() int {
	count := 0
	for _, c := range dcs.Challenges {
		if c.Completed {
			count++
		}
	}
	return count
}

// Display возвращает текстовое представление челленджей
func (dcs *DailyChallengeSystem) Display() string {
	if len(dcs.Challenges) == 0 {
		dcs.GenerateDailyChallenges()
	}

	text := fmt.Sprintf("🎯 <b>ЕЖЕДНЕВНЫЕ ЧЕЛЛЕНДЖИ</b>\n━━━━━━━━━━━━━━━━━━━━\n\n📅 Дата: %s\n\n",
		dcs.CurrentDate.Format("02.01.2006"))

	for i, c := range dcs.Challenges {
		status := "⬜"
		if c.Completed {
			status = "✅"
		}

		difficulty := ""
		switch c.Difficulty {
		case 1:
			difficulty = "🟢"
		case 2:
			difficulty = "🟡"
		case 3:
			difficulty = "🔴"
		}

		progress := c.Progress
		if progress > c.Target {
			progress = c.Target
		}

		text += fmt.Sprintf("%d. %s %s %s\n", i+1, status, difficulty, c.Name)
		text += fmt.Sprintf("   %s\n", c.Description)
		text += fmt.Sprintf("   Прогресс: %d/%d\n", progress, c.Target)
		text += fmt.Sprintf("   Награда: ✨%d XP | 💰%d монет\n\n", c.RewardXP, c.RewardCoins)
	}

	completed := dcs.GetCompletionCount()
	totalXP, totalCoins := dcs.GetTotalRewards()

	text += fmt.Sprintf("━━━━━━━━━━━━━━━━━━━━\n📊 Выполнено: %d/3\n✨ Заработано: %d XP\n💰 Заработано: %d монет",
		completed, totalXP, totalCoins)

	if completed == 3 {
		text += "\n\n🎉 <b>ВСЕ ЧЕЛЛЕНДЖИ ВЫПОЛНЕНЫ!</b> 🎉"
	}

	return text
}

// CheckStreak проверяет серию дней
func (dcs *DailyChallengeSystem) CheckStreak() int {
	now := time.Now()
	if dcs.LastChallenge.IsZero() {
		dcs.StreakDays = 1
		dcs.LastChallenge = now
		return dcs.StreakDays
	}

	daysDiff := int(now.Sub(dcs.LastChallenge).Hours() / 24)

	if daysDiff == 1 {
		dcs.StreakDays++
	} else if daysDiff > 1 {
		dcs.StreakDays = 1
	}

	dcs.LastChallenge = now
	return dcs.StreakDays
}

// GetStreakBonus возвращает бонус за серию дней
func (dcs *DailyChallengeSystem) GetStreakBonus() int {
	switch {
	case dcs.StreakDays >= 30:
		return 500
	case dcs.StreakDays >= 14:
		return 300
	case dcs.StreakDays >= 7:
		return 200
	case dcs.StreakDays >= 3:
		return 100
	default:
		return 0
	}
}

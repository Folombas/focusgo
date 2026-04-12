package game

import (
	"fmt"
	"time"
)

// PlayerStatsDay представляет статистику за один день
type PlayerStatsDay struct {
	Date           time.Time `json:"date"`
	StudyMinutes   int       `json:"study_minutes"`
	RestMinutes    int       `json:"rest_minutes"`
	TemptationsWon int       `json:"temptations_won"`
	TemptationsLost int      `json:"temptations_lost"`
	QuizCorrect    int       `json:"quiz_correct"`
	QuizTotal      int       `json:"quiz_total"`
	XPEarned       int       `json:"xp_earned"`
	CoinsEarned    int       `json:"coins_earned"`
	SkillsUpgraded int       `json:"skills_upgraded"`
	BossDefeated   bool      `json:"boss_defeated"`
	ChallengesDone int       `json:"challenges_done"`
}

// PlayerStatistics хранит полную статистику игрока
type PlayerStatistics struct {
	PlayerID    int64                    `json:"player_id"`
	DailyStats  map[string]PlayerStatsDay `json:"daily_stats"` // key: "2006-01-02"
	TotalDays   int                       `json:"total_days"`
	BestStreak  int                       `json:"best_streak"`
	CurrentStreak int                     `json:"current_streak"`
}

// NewPlayerStatistics создаёт новую систему статистики
func NewPlayerStatistics(playerID int64) *PlayerStatistics {
	return &PlayerStatistics{
		PlayerID:   playerID,
		DailyStats: make(map[string]PlayerStatsDay),
	}
}

// RecordDailyStats записывает статистику за день
func (ps *PlayerStatistics) RecordDailyStats(stats PlayerStatsDay) {
	dateKey := stats.Date.Format("2006-01-02")
	ps.DailyStats[dateKey] = stats
	ps.TotalDays++
}

// GetWeeklyReport возвращает отчёт за последние 7 дней
func (ps *PlayerStatistics) GetWeeklyReport() string {
	now := time.Now()
	weekAgo := now.Add(-7 * 24 * time.Hour)

	totalStudy := 0
	totalXP := 0
	totalCoins := 0
	totalTemptations := 0
	totalQuizCorrect := 0
	totalQuizTotal := 0
	totalChallenges := 0
	bossVictories := 0

	for dateStr, stats := range ps.DailyStats {
		date, _ := time.Parse("2006-01-02", dateStr)
		if date.After(weekAgo) {
			totalStudy += stats.StudyMinutes
			totalXP += stats.XPEarned
			totalCoins += stats.CoinsEarned
			totalTemptations += stats.TemptationsWon
			totalQuizCorrect += stats.QuizCorrect
			totalQuizTotal += stats.QuizTotal
			totalChallenges += stats.ChallengesDone
			if stats.BossDefeated {
				bossVictories++
			}
		}
	}

	quizPercent := 0
	if totalQuizTotal > 0 {
		quizPercent = (totalQuizCorrect * 100) / totalQuizTotal
	}

	return fmt.Sprintf(
		"📊 <b>НЕДЕЛЬНЫЙ ОТЧЁТ</b>\n━━━━━━━━━━━━━━━━━━━━\n\n"+
			"📅 Дней в игре: %d\n"+
			"🔥 Текущая серия: %d дней\n"+
			"🏆 Лучшая серия: %d дней\n\n"+
			"<b>Последние 7 дней:</b>\n\n"+
			"📚 Изучение Go: %d минут (%.1f часов)\n"+
			"✨ Получено XP: %d\n"+
			"💰 Заработано монет: %d\n"+
			"💪 Побеждено искушений: %d\n"+
			"🧩 Викторина: %d/%d (%d%%)\n"+
			"🎯 Выполнено челленджей: %d\n"+
			"👹 Побеждено боссов: %d\n\n"+
			"💪 Отличная работа! Продолжай учиться!",
		ps.TotalDays,
		ps.CurrentStreak,
		ps.BestStreak,
		totalStudy,
		float64(totalStudy)/60.0,
		totalXP,
		totalCoins,
		totalTemptations,
		totalQuizCorrect,
		totalQuizTotal,
		quizPercent,
		totalChallenges,
		bossVictories,
	)
}

// GetMonthlyReport возвращает отчёт за последние 30 дней
func (ps *PlayerStatistics) GetMonthlyReport() string {
	now := time.Now()
	monthAgo := now.Add(-30 * 24 * time.Hour)

	totalStudy := 0
	totalXP := 0
	totalCoins := 0
	totalTemptations := 0
	daysPlayed := 0

	for dateStr, stats := range ps.DailyStats {
		date, _ := time.Parse("2006-01-02", dateStr)
		if date.After(monthAgo) {
			totalStudy += stats.StudyMinutes
			totalXP += stats.XPEarned
			totalCoins += stats.CoinsEarned
			totalTemptations += stats.TemptationsWon
			daysPlayed++
		}
	}

	avgDailyStudy := 0
	if daysPlayed > 0 {
		avgDailyStudy = totalStudy / daysPlayed
	}

	return fmt.Sprintf(
		"📊 <b>МЕСЯЧНЫЙ ОТЧЁТ</b>\n━━━━━━━━━━━━━━━━━━━━\n\n"+
			"📅 Дней в игре: %d\n"+
			"🔥 Текущая серия: %d дней\n\n"+
			"<b>Последние 30 дней:</b>\n\n"+
			"📚 Общее время изучения: %d минут (%.1f часов)\n"+
			"📊 Среднее в день: %d минут\n"+
			"✨ Получено XP: %d\n"+
			"💰 Заработано монет: %d\n"+
			"💪 Побеждено искушений: %d\n\n"+
			"💡 <b>Совет:</b> Старайся учиться каждый день хотя бы по 30 минут!",
		ps.TotalDays,
		ps.CurrentStreak,
		totalStudy,
		float64(totalStudy)/60.0,
		avgDailyStudy,
		totalXP,
		totalCoins,
		totalTemptations,
	)
}

// GetStudyHeatMap возвращает heat map изучения за последние 30 дней
func (ps *PlayerStatistics) GetStudyHeatMap() string {
	now := time.Now()
	text := "🗓️ <b>КАРТА АКТИВНОСТИ (30 дней)</b>\n\n"

	for i := 29; i >= 0; i-- {
		date := now.Add(-time.Duration(i) * 24 * time.Hour)
		dateKey := date.Format("2006-01-02")
		stats, exists := ps.DailyStats[dateKey]

		dayStr := date.Format("02.01")
		icon := "⬜"

		if exists && stats.StudyMinutes > 0 {
			switch {
			case stats.StudyMinutes >= 120:
				icon = "🟩"
			case stats.StudyMinutes >= 60:
				icon = "🟨"
			case stats.StudyMinutes >= 30:
				icon = "🟧"
			case stats.StudyMinutes > 0:
				icon = "🟥"
			}
		}

		text += fmt.Sprintf("%s %s", dayStr, icon)

		if i%7 == 0 {
			text += "\n"
		} else {
			text += " "
		}
	}

	text += "\n🟥 <10мин | 🟧 30мин | 🟨 60мин | 🟩 120мин+"

	return text
}

// UpdateStreak обновляет серию дней
func (ps *PlayerStatistics) UpdateStreak() {
	now := time.Now()
	dateKey := now.Format("2006-01-02")

	// Проверяем последний день
	lastDate := now.Add(-24 * time.Hour)
	lastKey := lastDate.Format("2006-01-02")

	_, lastExists := ps.DailyStats[lastKey]
	_, todayExists := ps.DailyStats[dateKey]

	if todayExists {
		return // Уже отмечено сегодня
	}

	if lastExists {
		ps.CurrentStreak++
	} else {
		ps.CurrentStreak = 1
	}

	if ps.CurrentStreak > ps.BestStreak {
		ps.BestStreak = ps.CurrentStreak
	}
}

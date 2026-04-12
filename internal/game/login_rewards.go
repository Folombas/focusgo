package game

import (
	"fmt"
	"time"
)

// DailyLoginReward представляет награду за ежедневный вход
type DailyLoginReward struct {
	DayNumber int    `json:"day_number"`
	RewardXP  int    `json:"reward_xp"`
	RewardCoins int  `json:"reward_coins"`
	BonusItem string `json:"bonus_item,omitempty"`
	Claimed   bool   `json:"claimed"`
}

// LoginRewardSystem управляет системой ежедневных наград
type LoginRewardSystem struct {
	PlayerID      int64              `json:"player_id"`
	CurrentDay    int                `json:"current_day"`
	LastLogin     time.Time          `json:"last_login"`
	StreakDays    int                `json:"streak_days"`
	Rewards       []DailyLoginReward `json:"rewards"`
	TotalClaimed  int                `json:"total_claimed"`
}

// NewLoginRewardSystem создаёт новую систему наград за вход
func NewLoginRewardSystem(playerID int64) *LoginRewardSystem {
	lrs := &LoginRewardSystem{
		PlayerID:   playerID,
		CurrentDay: 1,
		Rewards:    make([]DailyLoginReward, 30),
	}

	// Инициализируем награды на 30 дней
	lrs.initializeRewards()
	return lrs
}

// initializeRewards инициализирует награды на 30 дней
func (lrs *LoginRewardSystem) initializeRewards() {
	for i := 0; i < 30; i++ {
		dayNum := i + 1
		reward := DailyLoginReward{
			DayNumber: dayNum,
			RewardXP:  50 + (i * 25),
			RewardCoins: 25 + (i * 15),
		}

		// Бонусные предметы на определённых днях
		switch dayNum {
		case 7:
			reward.BonusItem = "🎁 Недельный бонус: +100 XP к следующему квесту"
		case 14:
			reward.BonusItem = "🎁 Двухнедельный бонус: Удвоение XP на 1 час"
		case 21:
			reward.BonusItem = "🎁 Трёхнедельный бонус: +50 монет"
		case 28:
			reward.BonusItem = "🎁 Почти месяц! +200 XP"
		case 30:
			reward.BonusItem = "🏆 МЕСЯЦ В ИГРЕ! Легендарная награда: +500 XP, +250 монет"
		}

		lrs.Rewards[i] = reward
	}
}

// ClaimDailyReward получает награду за текущий день
func (lrs *LoginRewardSystem) ClaimDailyReward() (rewardText string, xp int, coins int) {
	now := time.Now()

	// Проверяем, не получал ли уже сегодня
	if !lrs.LastLogin.IsZero() && lrs.LastLogin.Day() == now.Day() {
		if lrs.Rewards[lrs.CurrentDay-1].Claimed {
			return "⚠️  Ты уже получил награду за сегодня. Возвращайся завтра!", 0, 0
		}
	}

	// Проверяем серию
	daysDiff := 0
	if !lrs.LastLogin.IsZero() {
		daysDiff = int(now.Sub(lrs.LastLogin).Hours() / 24)
	}

	if daysDiff == 1 {
		// Продолжаем серию
		lrs.StreakDays++
		lrs.CurrentDay = min(lrs.CurrentDay+1, 30)
	} else if daysDiff > 1 {
		// Серия прервана, начинаем сначала
		if lrs.StreakDays > 0 {
			rewardText += fmt.Sprintf("⚠️  Серия прервана! Было %d дней, начинаем сначала.\n\n", lrs.StreakDays)
		}
		lrs.StreakDays = 1
		lrs.CurrentDay = 1
	} else if daysDiff == 0 {
		// Тот же день - уже получал
		return "⚠️  Ты уже получил награду за сегодня. Возвращайся завтра!", 0, 0
	} else {
		// Первый вход
		lrs.StreakDays = 1
	}

	// Получаем награду
	reward := lrs.Rewards[lrs.CurrentDay-1]
	reward.Claimed = true
	lrs.Rewards[lrs.CurrentDay-1] = reward
	lrs.LastLogin = now
	lrs.TotalClaimed++

	xp = reward.RewardXP
	coins = reward.RewardCoins

	// Формируем текст награды
	rewardText = fmt.Sprintf(
		"🎁 <b>НАГРАДА ЗА ДЕНЬ %d</b>\n━━━━━━━━━━━━━━━━━━━━\n\n"+
		"✨ Получено: %d XP\n"+
		"💰 Получено: %d монет\n"+
		"🔥 Серия входа: %d дней\n\n",
		lrs.CurrentDay,
		xp,
		coins,
		lrs.StreakDays,
	)

	if reward.BonusItem != "" {
		rewardText += fmt.Sprintf("%s\n\n", reward.BonusItem)
	}

	// Показываем следующую награду
	if lrs.CurrentDay < 30 {
		nextReward := lrs.Rewards[lrs.CurrentDay]
		rewardText += fmt.Sprintf("🔜 <b>Завтра:</b> День %d\n✨ %d XP | 💰 %d монет",
			nextReward.DayNumber,
			nextReward.RewardXP,
			nextReward.RewardCoins)
	} else {
		rewardText += "🏆 <b>Цикл завершён! Начинаем заново с бонусом!</b>"
		lrs.CurrentDay = 1
		lrs.initializeRewards()
	}

	return rewardText, xp, coins
}

// Display возвращает текстовое представление системы наград
func (lrs *LoginRewardSystem) Display() string {
	text := fmt.Sprintf("🎁 <b>ЕЖЕДНЕВНЫЕ НАГРАДЫ</b>\n━━━━━━━━━━━━━━━━━━━━\n\n"+
		"🔥 Серия: %d дней\n"+
		"📅 Текущий день: %d/30\n"+
		"✅ Всего получено: %d наград\n\n",
		lrs.StreakDays,
		lrs.CurrentDay,
		lrs.TotalClaimed)

	// Показываем ближайшие 7 дней
	text += "<b>Ближайшие награды:</b>\n\n"
	for i := 0; i < 7 && (lrs.CurrentDay-1+i) < 30; i++ {
		idx := lrs.CurrentDay - 1 + i
		reward := lrs.Rewards[idx]

		status := "⬜"
		if reward.Claimed {
			status = "✅"
		}

		if i == 0 {
			status = "👉"
		}

		text += fmt.Sprintf("%s День %d: ✨%d XP | 💰%d монет\n", status, reward.DayNumber, reward.RewardXP, reward.RewardCoins)

		if reward.BonusItem != "" {
			text += fmt.Sprintf("   %s\n", reward.BonusItem)
		}
	}

	text += "\n━━━━━━━━━━━━━━━━━━━━\n💡 Заходи каждый день, чтобы получать больше!"

	return text
}

// GetStreakBonus возвращает бонус за серию дней
func (lrs *LoginRewardSystem) GetStreakBonus() int {
	return lrs.StreakDays * 10
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

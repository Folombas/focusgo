package models

// clampInt ограничивает целое число диапазоном [min, max]
func clampInt(value, min, max int) int {
	if value < min {
		return min
	}
	if value > max {
		return max
	}
	return value
}

// clampStat ограничивает характеристику диапазоном [0, 100]
func clampStat(value int) int {
	return clampInt(value, 0, 100)
}

// clampMoney ограничивает деньги диапазоном [0, 999999]
func clampMoney(value int) int {
	return clampInt(value, 0, 999999)
}

// clampDopamine ограничивает дофамин диапазоном [0, 999]
func clampDopamine(value int) int {
	return clampInt(value, 0, 999)
}

// clampExperience ограничивает опыт диапазоном [0, 999999]
func clampExperience(value int) int {
	return clampInt(value, 0, 999999)
}

// clampLevel ограничивает уровень диапазоном [1, 100]
func clampLevel(value int) int {
	return clampInt(value, 1, 100)
}

// validatePlayer проверяет валидность игрока
func validatePlayer(player *Player) error {
	// Базовая валидация
	if player.Name == "" {
		player.Name = "Игрок"
	}
	if player.ChatID == 0 {
		player.ChatID = 1
	}
	if player.Level < 1 {
		player.Level = 1
	}
	if player.Experience < 0 {
		player.Experience = 0
	}
	// Остальные поля ограничиваются clamp функциями при изменении
	return nil
}

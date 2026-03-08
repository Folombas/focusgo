package game

import (
	"database/sql"
	"fmt"
	"strings"

	"focusgo/internal/database"
)

// Skill представляет навык в дереве развития
type Skill struct {
	ID            string
	Name          string
	Description   string
	Level         int      // 0-5
	MaxLevel      int      // максимальный уровень
	CostPerLevel  int      // стоимость улучшения за уровень
	BonusType     string   // focus, willpower, knowledge, money, dopamine
	BonusValue    int      // значение бонуса за уровень
	Unlocked      bool     // разблокирован ли
	Prerequisites []string // требуемые навыки
}

// SkillTree представляет дерево навыков игрока
type SkillTree struct {
	ChatID      int64
	Skills      map[string]*Skill
	SkillPoints int // доступные очки
	TotalPoints int // всего заработано
}

// NewSkillTree создаёт новое дерево навыков
func NewSkillTree(chatID int64) *SkillTree {
	tree := &SkillTree{
		ChatID:      chatID,
		Skills:      make(map[string]*Skill),
		SkillPoints: 0,
		TotalPoints: 0,
	}

	tree.initSkills()
	return tree
}

// initSkills инициализирует 12 навыков в 4 категориях
func (st *SkillTree) initSkills() {
	// 📚 GO-НАВЫКИ (6 навыков)
	st.Skills["go_basics"] = &Skill{
		ID:            "go_basics",
		Name:          "📘 Основы Go",
		Description:   "Синтаксис, типы данных, функции",
		Level:         1, // начальный навык
		MaxLevel:      5,
		CostPerLevel:  1,
		BonusType:     "knowledge",
		BonusValue:    5,
		Unlocked:      true,
		Prerequisites: []string{},
	}

	st.Skills["concurrency"] = &Skill{
		ID:            "concurrency",
		Name:          "⚡ Конкурентность",
		Description:   "Горутины, каналы, sync package",
		Level:         0,
		MaxLevel:      5,
		CostPerLevel:  2,
		BonusType:     "knowledge",
		BonusValue:    8,
		Unlocked:      false,
		Prerequisites: []string{"go_basics"},
	}

	st.Skills["interfaces"] = &Skill{
		ID:            "interfaces",
		Name:          "🔌 Интерфейсы",
		Description:   "Интерфейсы, полиморфизм",
		Level:         0,
		MaxLevel:      5,
		CostPerLevel:  2,
		BonusType:     "knowledge",
		BonusValue:    7,
		Unlocked:      false,
		Prerequisites: []string{"go_basics"},
	}

	st.Skills["web_frameworks"] = &Skill{
		ID:            "web_frameworks",
		Name:          "🌐 Web Фреймворки",
		Description:   "Gin, Echo, Fiber — создание API",
		Level:         0,
		MaxLevel:      5,
		CostPerLevel:  3,
		BonusType:     "knowledge",
		BonusValue:    10,
		Unlocked:      false,
		Prerequisites: []string{"concurrency", "interfaces"},
	}

	st.Skills["database"] = &Skill{
		ID:            "database",
		Name:          "🗄️ Базы данных",
		Description:   "PostgreSQL, MongoDB, Redis",
		Level:         0,
		MaxLevel:      5,
		CostPerLevel:  3,
		BonusType:     "knowledge",
		BonusValue:    10,
		Unlocked:      false,
		Prerequisites: []string{"concurrency"},
	}

	st.Skills["microservices"] = &Skill{
		ID:            "microservices",
		Name:          "🔧 Микросервисы",
		Description:   "gRPC, Docker, Kubernetes",
		Level:         0,
		MaxLevel:      5,
		CostPerLevel:  4,
		BonusType:     "knowledge",
		BonusValue:    12,
		Unlocked:      false,
		Prerequisites: []string{"web_frameworks", "database"},
	}

	// 🎯 ФОКУС (3 навыка)
	st.Skills["focus_master"] = &Skill{
		ID:            "focus_master",
		Name:          "🎯 Мастер Фокуса",
		Description:   "Умение концентрироваться на задачах",
		Level:         0,
		MaxLevel:      5,
		CostPerLevel:  1,
		BonusType:     "focus",
		BonusValue:    5,
		Unlocked:      true,
		Prerequisites: []string{},
	}

	st.Skills["meditation"] = &Skill{
		ID:            "meditation",
		Name:          "🧘 Медитация",
		Description:   "Восстановление фокуса и энергии",
		Level:         0,
		MaxLevel:      5,
		CostPerLevel:  2,
		BonusType:     "dopamine",
		BonusValue:    10,
		Unlocked:      false,
		Prerequisites: []string{"focus_master"},
	}

	st.Skills["anti_procrastination"] = &Skill{
		ID:            "anti_procrastination",
		Name:          "⏰ Борьба с прокрастинацией",
		Description:   "Техники Pomodoro, тайм-менеджмент",
		Level:         0,
		MaxLevel:      5,
		CostPerLevel:  2,
		BonusType:     "focus",
		BonusValue:    8,
		Unlocked:      false,
		Prerequisites: []string{"focus_master"},
	}

	// 💪 СИЛА ВОЛИ (2 навыка)
	st.Skills["willpower"] = &Skill{
		ID:            "willpower",
		Name:          "💪 Сила Воли",
		Description:   "Сопротивление искушениям",
		Level:         0,
		MaxLevel:      5,
		CostPerLevel:  2,
		BonusType:     "willpower",
		BonusValue:    8,
		Unlocked:      true,
		Prerequisites: []string{},
	}

	st.Skills["discipline"] = &Skill{
		ID:            "discipline",
		Name:          "📅 Дисциплина",
		Description:   "Ежедневное следование плану",
		Level:         0,
		MaxLevel:      5,
		CostPerLevel:  3,
		BonusType:     "willpower",
		BonusValue:    10,
		Unlocked:      false,
		Prerequisites: []string{"willpower"},
	}

	// 💰 ФИНАНСЫ (1 навык)
	st.Skills["money_management"] = &Skill{
		ID:            "money_management",
		Name:          "💰 Управление деньгами",
		Description:   "Экономия и инвестиции в обучение",
		Level:         0,
		MaxLevel:      5,
		CostPerLevel:  2,
		BonusType:     "money",
		BonusValue:    50,
		Unlocked:      false,
		Prerequisites: []string{"willpower"},
	}
}

// LoadSkillTree загружает дерево навыков из БД
func LoadSkillTree(chatID int64) (*SkillTree, error) {
	tree := NewSkillTree(chatID)

	// Загружаем очки навыков
	query := `SELECT skill_points, total_points FROM skill_trees WHERE chat_id = ?`
	row := database.DB.QueryRow(query, chatID)

	err := row.Scan(&tree.SkillPoints, &tree.TotalPoints)
	if err == sql.ErrNoRows {
		// Создаём новую запись
		tree.SaveSkillTree()
		return tree, nil
	}
	if err != nil {
		return nil, err
	}

	// Загружаем уровни навыков
	query = `SELECT skill_id, level, unlocked FROM skills WHERE chat_id = ?`
	rows, err := database.DB.Query(query, chatID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var skillID string
		var level int
		var unlocked int

		if err := rows.Scan(&skillID, &level, &unlocked); err != nil {
			continue
		}

		if skill, exists := tree.Skills[skillID]; exists {
			skill.Level = level
			skill.Unlocked = unlocked == 1
		}
	}

	// Проверяем разблокировки
	tree.checkUnlocks()

	return tree, nil
}

// SaveSkillTree сохраняет дерево навыков в БД
func (st *SkillTree) SaveSkillTree() error {
	// Сохраняем дерево навыков
	query := `INSERT OR REPLACE INTO skill_trees (chat_id, skill_points, total_points) VALUES (?, ?, ?)`
	_, err := database.DB.Exec(query, st.ChatID, st.SkillPoints, st.TotalPoints)
	if err != nil {
		return err
	}

	// Сохраняем каждый навык
	for id, skill := range st.Skills {
		query = `INSERT OR REPLACE INTO skills (chat_id, skill_id, level, unlocked) VALUES (?, ?, ?, ?)`
		_, err := database.DB.Exec(query, st.ChatID, id, skill.Level, boolToInt(skill.Unlocked))
		if err != nil {
			return err
		}
	}

	return nil
}

// EarnSkillPoints начисляет очки навыков
func (st *SkillTree) EarnSkillPoints(points int) {
	if points < 0 {
		points = 0
	}
	st.SkillPoints += points
	st.TotalPoints += points
}

// UpgradeSkill улучшает навык
func (st *SkillTree) UpgradeSkill(skillID string) (bool, string) {
	skill, exists := st.Skills[skillID]
	if !exists {
		return false, "❌ Навык не найден"
	}

	if !skill.Unlocked {
		return false, "❌ Навык заблокирован. Изучите требуемые навыки сначала"
	}

	if skill.Level >= skill.MaxLevel {
		return false, "⚠️  Навык уже максимального уровня"
	}

	cost := skill.CostPerLevel
	if st.SkillPoints < cost {
		return false, fmt.Sprintf("❌ Недостаточно очков навыков (нужно %d, есть %d)", cost, st.SkillPoints)
	}

	st.SkillPoints -= cost
	skill.Level++

	// Проверяем разблокировку следующих навыков
	st.checkUnlocks()

	return true, fmt.Sprintf("🎉 Навык \"%s\" улучшен до уровня %d!\n+%d к %s",
		skill.Name, skill.Level, skill.BonusValue, translateBonusType(skill.BonusType))
}

// checkUnlocks проверяет и разблокирует новые навыки
func (st *SkillTree) checkUnlocks() {
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

// Display отображает дерево навыков
func (st *SkillTree) Display() string {
	var sb strings.Builder

	sb.WriteString("🌳 <b>ДЕРЕВО НАВЫКОВ</b>\n")
	sb.WriteString("━━━━━━━━━━━━━━━━━━━━\n\n")
	sb.WriteString(fmt.Sprintf("✨ Очки навыков: %d (всего: %d)\n\n", st.SkillPoints, st.TotalPoints))

	// Группируем по категориям
	categories := map[string][]string{
		"📚 GO-НАВЫКИ": {"go_basics", "concurrency", "interfaces", "web_frameworks", "database", "microservices"},
		"🎯 ФОКУС":     {"focus_master", "meditation", "anti_procrastination"},
		"💪 СИЛА ВОЛИ": {"willpower", "discipline"},
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

	return fmt.Sprintf("  %s %-20s [%s] Ур.%d/%d%s\n     %s\n",
		status, skill.Name, levelBar, skill.Level, skill.MaxLevel, cost, skill.Description)
}

// GetUpgradeKeyboard создаёт клавиатуру для улучшения навыков
func (st *SkillTree) GetUpgradeKeyboard() [][]InlineKeyboardButton {
	var keyboard [][]InlineKeyboardButton

	if st.SkillPoints > 0 {
		for id, skill := range st.Skills {
			if skill.Unlocked && skill.Level < skill.MaxLevel {
				cost := skill.CostPerLevel
				if st.SkillPoints >= cost {
					keyboard = append(keyboard, []InlineKeyboardButton{
						{Text: fmt.Sprintf("⬆️ %s (%d очк.)", skill.Name, cost), CallbackData: fmt.Sprintf("cb_upgrade_%s", id)},
					})
				}
			}
		}
	}

	keyboard = append(keyboard, []InlineKeyboardButton{
		{Text: "🔙 Назад", CallbackData: "cb_main_menu"},
	})

	return keyboard
}

// translateBonusType переводит тип бонуса
func translateBonusType(bonusType string) string {
	switch bonusType {
	case "focus":
		return "Фокус"
	case "willpower":
		return "Сила воли"
	case "knowledge":
		return "Знание Go"
	case "money":
		return "Деньги"
	case "dopamine":
		return "Дофамин"
	default:
		return bonusType
	}
}

// boolToInt конвертирует bool в int
func boolToInt(b bool) int {
	if b {
		return 1
	}
	return 0
}

// InlineKeyboardButton для клавиатуры
type InlineKeyboardButton struct {
	Text         string
	CallbackData string
}

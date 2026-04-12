package game

import (
	"fmt"
	"math/rand"
	"time"
)

// MissionNode представляет узел миссии (сценарий)
type MissionNode struct {
	ID          string   `json:"id"`
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Type        string   `json:"type"` // study, temptation, quiz, boss, choice
	Choices     []MissionChoice `json:"choices,omitempty"`
	RewardXP    int      `json:"reward_xp"`
	RewardCoins int      `json:"reward_coins"`
	RequiredLevel int    `json:"required_level"`
	Completed   bool     `json:"completed"`
}

// MissionChoice представляет выбор в миссии
type MissionChoice struct {
	Text        string `json:"text"`
	NextNodeID  string `json:"next_node_id"`
	RewardXP    int    `json:"reward_xp"`
	Consequence string `json:"consequence"` // Описание последствия выбора
}

// MissionChain представляет цепочку миссий
type MissionChain struct {
	ID          string       `json:"id"`
	Name        string       `json:"name"`
	Description string       `json:"description"`
	Nodes       []MissionNode `json:"nodes"`
	CurrentNode string       `json:"current_node"`
	Started     time.Time    `json:"started"`
	Completed   bool         `json:"completed"`
	Difficulty  int          `json:"difficulty"` // 1-3
}

// MissionSystem управляет системой миссий
type MissionSystem struct {
	PlayerID      int64                     `json:"player_id"`
	ActiveMissions []MissionChain            `json:"active_missions"`
	CompletedMissions []string               `json:"completed_missions"`
	MissionHistory []string                  `json:"mission_history"`
}

// NewMissionSystem создаёт новую систему миссий
func NewMissionSystem(playerID int64) *MissionSystem {
	return &MissionSystem{
		PlayerID:      playerID,
		ActiveMissions: make([]MissionChain, 0),
		CompletedMissions: make([]string, 0),
	}
}

// GenerateMission генерирует новую миссию с ветвящимся сценарием
func (ms *MissionSystem) GenerateMission(playerLevel int) *MissionChain {
	missionTemplates := ms.getMissionTemplates()

	// Выбираем миссию по уровню
	var suitableMissions []MissionChain
	for _, template := range missionTemplates {
		if template.Difficulty <= (playerLevel/10)+1 {
			suitableMissions = append(suitableMissions, template)
		}
	}

	if len(suitableMissions) == 0 {
		suitableMissions = missionTemplates[:1]
	}

	// Выбираем случайную
	rand.Seed(time.Now().UnixNano())
	mission := suitableMissions[rand.Intn(len(suitableMissions))]

	// Сбрасываем состояние
	mission.CurrentNode = mission.Nodes[0].ID
	mission.Started = time.Now()
	mission.Completed = false

	ms.ActiveMissions = append(ms.ActiveMissions, mission)
	return &mission
}

// getMissionTemplates возвращает шаблоны миссий
func (ms *MissionSystem) getMissionTemplates() []MissionChain {
	return []MissionChain{
		{
			ID:          "mission_first_day",
			Name:        "🌟 Первый день разработчика",
			Description: "Твой первый день в мире Go! Что выберешь?",
			Difficulty:  1,
			Nodes: []MissionNode{
				{
					ID:   "start",
					Title: "Начало пути",
					Description: "Ты открыл учебник по Go. Первая глава: переменные.\nЧто будешь делать?",
					Type: "choice",
					Choices: []MissionChoice{
						{
							Text: "📚 Читать теорию внимательно",
							NextNodeID: "study_careful",
							RewardXP: 50,
							Consequence: "Ты понял основы лучше, но потратил больше времени",
						},
						{
							Text: "💻 Сразу писать код",
							NextNodeID: "practice_fast",
							RewardXP: 40,
							Consequence: "Быстро начал, но некоторые темы упустил",
						},
					},
					RewardXP: 0,
				},
				{
					ID:   "study_careful",
					Title: "Внимательное изучение",
					Description: "Ты внимательно прочитал про переменные, типы и объявление.\nТеперь попробуй написать код!",
					Type: "study",
					RewardXP: 80,
					RewardCoins: 30,
				},
				{
					ID:   "practice_fast",
					Title: "Быстрая практика",
					Description: "Ты начал писать код, но столкнулся с ошибками типов.\nПридётся вернуться к теории!",
					Type: "study",
					RewardXP: 60,
					RewardCoins: 20,
				},
			},
		},
		{
			ID:          "mission_temptation",
			Name:        "🎮 Искушение игр",
			Description: "Друг зовёт играть в Steam. Как поступишь?",
			Difficulty:  2,
			Nodes: []MissionNode{
				{
					ID:   "start",
					Title: "Зов друга",
					Description: "Привет! Погнали в CS, уже 3 человека собрались!\nТы только начал учить Go 30 минут назад...",
					Type: "choice",
					Choices: []MissionChoice{
						{
							Text: "💪 Отказать и учиться",
							NextNodeID: "resist",
							RewardXP: 100,
							Consequence: "Сила воли +10! Продолжил изучение Go",
						},
						{
							Text: "🎮 Пойти играть",
							NextNodeID: "give_in",
							RewardXP: 10,
							Consequence: "Играл 3 часа. Чувствуешь вину...",
						},
						{
							Text: "⏰ Сначала 1 час Go, потом игра",
							NextNodeID: "compromise",
							RewardXP: 70,
							Consequence: "Отличный баланс! И поучился, и отдохнул",
						},
					},
					RewardXP: 0,
				},
				{
					ID:   "resist",
					Title: "Победа над искушением",
					Description: "Ты отказался и продолжил учить конкурентность в Go.\nЧерез час ты уже понимаешь горутины!",
					Type: "study",
					RewardXP: 120,
					RewardCoins: 50,
				},
				{
					ID:   "give_in",
					Title: "Поддался искушению",
					Description: "Ты играл до полуночи. Завтра обещай себе начать раньше!",
					Type: "temption",
					RewardXP: 10,
					RewardCoins: 5,
				},
				{
					ID:   "compromise",
					Title: "Мудрое решение",
					Description: "Ты позанимался час и пошёл играть с чистой совестью.\nОтличный баланс между учёбой и отдыхом!",
					Type: "study",
					RewardXP: 90,
					RewardCoins: 40,
				},
			},
		},
		{
			ID:          "mission_quiz_challenge",
			Name:        "🧩 Испытание викториной",
			Description: "Викторина с выбором сложности",
			Difficulty:  2,
			Nodes: []MissionNode{
				{
					ID:   "start",
					Title: "Выбор сложности",
					Description: "Ты вошёл в Go-викторину! Какой уровень выберешь?",
					Type: "choice",
					Choices: []MissionChoice{
						{
							Text: "🟢 Лёгкий (безопасно)",
							NextNodeID: "easy_quiz",
							RewardXP: 50,
							Consequence: "Лёгкие вопросы, но мало XP",
						},
						{
							Text: "🟡 Средний (баланс)",
							NextNodeID: "medium_quiz",
							RewardXP: 100,
							Consequence: "Интересные вопросы, хороший XP",
						},
						{
							Text: "🔴 Сложный (хардкор!)",
							NextNodeID: "hard_quiz",
							RewardXP: 200,
							Consequence: "Очень сложно, но много XP",
						},
					},
					RewardXP: 0,
				},
				{
					ID:   "easy_quiz",
					Title: "Лёгкий уровень",
					Description: "Ты ответил на 8/10 вопросов правильно!\nХороший результат для начала!",
					Type: "quiz",
					RewardXP: 80,
					RewardCoins: 30,
				},
				{
					ID:   "medium_quiz",
					Title: "Средний уровень",
					Description: "Ты ответил на 7/10 вопросов!\nОтлично! Конкурентность далась тяжело.",
					Type: "quiz",
					RewardXP: 150,
					RewardCoins: 60,
				},
				{
					ID:   "hard_quiz",
					Title: "Сложный уровень",
					Description: "Ты ответил на 5/10 вопросов...\nНо столько узнал нового! Это победа!",
					Type: "quiz",
					RewardXP: 200,
					RewardCoins: 80,
				},
			},
		},
		{
			ID:          "mission_boss_fight",
			Name:        "👹 Битва с боссом: Прокрастинация",
			Description: "Финальная битва дня! Готов ли ты?",
			Difficulty:  3,
			Nodes: []MissionNode{
				{
					ID:   "start",
					Title: "Появление босса",
					Description: "👹 ПРОКРАСТИНАЦИЯ МАКСИМА появилась!\n\nСила: 90\nСпособности: 'Ещё 5 минут', 'Завтра начну', 'Я устал'",
					Type: "choice",
					Choices: []MissionChoice{
						{
							Text: "⚔️ Атаковать знанием Go",
							NextNodeID: "attack_knowledge",
							RewardXP: 150,
							Consequence: "Используешь силу знаний!",
						},
						{
							Text: "🛡️ Защититься дисциплиной",
							NextNodeID: "defend_discipline",
							RewardXP: 130,
							Consequence: "Дисциплина - щит от прокрастинации!",
						},
						{
							Text: "🧘 Медитировать",
							NextNodeID: "meditate",
							RewardXP: 120,
							Consequence: "Спокойный ум сильнее искушений",
						},
					},
					RewardXP: 0,
				},
				{
					ID:   "attack_knowledge",
					Title: "Атака знанием!",
					Description: "Ты использовал знание горутин и каналов!\nБосс ослаблен на 50%!\n\nНо он использует 'Завтра начну'!",
					Type: "boss",
					RewardXP: 200,
					RewardCoins: 100,
				},
				{
					ID:   "defend_discipline",
					Title: "Защита дисциплиной!",
					Description: "Твой щит дисциплины выдержал!\nБосс отступает с позором!\n\nПобеда!",
					Type: "boss",
					RewardXP: 180,
					RewardCoins: 90,
				},
				{
					ID:   "meditate",
					Title: "Сила медитации",
					Description: "Ты закрыл глаза и медитировал 5 минут.\nБосс рассеялся как дым...\n\nМудрая победа!",
					Type: "boss",
					RewardXP: 160,
					RewardCoins: 80,
				},
			},
		},
	}
}

// GetCurrentNode возвращает текущий узел миссии
func (ms *MissionSystem) GetCurrentNode(missionID string) *MissionNode {
	for i := range ms.ActiveMissions {
		if ms.ActiveMissions[i].ID == missionID {
			for j := range ms.ActiveMissions[i].Nodes {
				if ms.ActiveMissions[i].Nodes[j].ID == ms.ActiveMissions[i].CurrentNode {
					return &ms.ActiveMissions[i].Nodes[j]
				}
			}
		}
	}
	return nil
}

// MakeChoice делает выбор в миссии
func (ms *MissionSystem) MakeChoice(missionID string, choiceIndex int) (result string, xp int, coins int) {
	mission := ms.getActiveMission(missionID)
	if mission == nil {
		return "❌ Миссия не найдена", 0, 0
	}

	currentNode := ms.GetCurrentNode(missionID)
	if currentNode == nil || len(currentNode.Choices) == 0 {
		return "❌ Нет доступных выборов", 0, 0
	}

	if choiceIndex < 0 || choiceIndex >= len(currentNode.Choices) {
		return "❌ Неверный выбор", 0, 0
	}

	choice := currentNode.Choices[choiceIndex]

	// Переходим к следующему узлу
	mission.CurrentNode = choice.NextNodeID
	nextNode := ms.GetCurrentNode(missionID)

	if nextNode != nil {
		result = fmt.Sprintf(
			"📖 %s\n\n%s\n\n✨ +%d XP\n💰 +%d монет\n\n🔜 %s",
			nextNode.Title,
			nextNode.Description,
			choice.RewardXP,
			nextNode.RewardCoins,
			choice.Consequence,
		)

		xp = choice.RewardXP + nextNode.RewardXP
		coins = nextNode.RewardCoins
	} else {
		// Миссия завершена
		mission.Completed = true
		ms.CompletedMissions = append(ms.CompletedMissions, missionID)

		result = fmt.Sprintf(
			"🎉 <b>МИССИЯ ЗАВЕРШЕНА!</b>\n\n%s\n\n✨ Итого: +%d XP\n💰 Итого: +%d монет\n\n🏆 Миссия добавлена в историю!",
			choice.Consequence,
			choice.RewardXP,
			mission.getTotalRewardCoins(),
		)

		xp = choice.RewardXP
		coins = mission.getTotalRewardCoins()
	}

	return result, xp, coins
}

// getActiveMission возвращает активную миссию
func (ms *MissionSystem) getActiveMission(missionID string) *MissionChain {
	for i := range ms.ActiveMissions {
		if ms.ActiveMissions[i].ID == missionID {
			return &ms.ActiveMissions[i]
		}
	}
	return nil
}

// Display возвращает текстовое представление системы миссий
func (ms *MissionSystem) Display() string {
	if len(ms.ActiveMissions) == 0 {
		return "📖 <b>МИССИИ</b>\n━━━━━━━━━━━━━━━━━━━━\n\nНет активных миссий.\nИспользуй /mission_start для начала!"
	}

	text := "📖 <b>АКТИВНЫЕ МИССИИ</b>\n━━━━━━━━━━━━━━━━━━━━\n\n"

	for _, mission := range ms.ActiveMissions {
		status := "🟡"
		if mission.Completed {
			status = "✅"
		}

		difficulty := ""
		switch mission.Difficulty {
		case 1:
			difficulty = "🟢"
		case 2:
			difficulty = "🟡"
		case 3:
			difficulty = "🔴"
		}

		text += fmt.Sprintf("%s %s %s\n", status, difficulty, mission.Name)
		text += fmt.Sprintf("   %s\n\n", mission.Description)

		currentNode := ms.GetCurrentNode(mission.ID)
		if currentNode != nil && !mission.Completed {
			text += fmt.Sprintf("<b>Текущий этап:</b> %s\n", currentNode.Title)
			text += fmt.Sprintf("%s\n\n", currentNode.Description)

			if len(currentNode.Choices) > 0 {
				text += "<b>Выбор:</b>\n\n"
				for i, choice := range currentNode.Choices {
					text += fmt.Sprintf("%d. %s\n", i+1, choice.Text)
				}
				text += "\nИспользуй /mission_choice <номер>"
			}
		}

		text += "\n━━━━━━━━━━━━━━━━━━━━\n"
	}

	return text
}

// getTotalReward XP возвращает суммарную награду миссии
func (mc *MissionChain) getTotalRewardXP() int {
	total := 0
	for _, node := range mc.Nodes {
		total += node.RewardXP
	}
	return total
}

// getTotalRewardCoins возвращает суммарные монеты миссии
func (mc *MissionChain) getTotalRewardCoins() int {
	total := 0
	for _, node := range mc.Nodes {
		total += node.RewardCoins
	}
	return total
}

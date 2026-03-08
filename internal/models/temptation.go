package models

import (
	"math/rand"
	"time"
)

// Temptation представляет искушение
type Temptation struct {
	Name        string
	Power       int    // Сила искушения (0-100)
	Description string
	XPLoss      int    // Потеря опыта при поддаче
	Category    string // Категория
}

// Motivation представляет мотивацию
type Motivation struct {
	Text    string
	Effect  string // "focus+", "willpower+", "knowledge+"
	XPBonus int
}

var temptations = []Temptation{
	// Цифровые искушения
	{"CapCut Видеомонтаж", 85, "Установить программу и монтировать видео", 50, "digital"},
	{"Видеоигры", 70, "Поиграть в новую RPG-игру", 40, "digital"},
	{"Соцсети", 60, "Прокрутить ленту 2 часа", 30, "digital"},
	{"YouTube", 65, "Посмотреть 'ещё одно' видео", 35, "digital"},
	{"Telegram каналы", 55, "Читать всё подряд вместо учёбы", 25, "digital"},
	{"Мобильные игры", 60, "Проверить ежедневные награды", 30, "digital"},

	// Социальные искушения
	{"Бары/Клубы", 75, "Сходить в бар с друзьями", 45, "social"},
	{"Встречи без цели", 50, "Пойти на встречу 'просто поболтать'", 25, "social"},
	{"Звонки друзьям", 45, "Долгий разговор вместо учёбы", 20, "social"},

	// Потребительские искушения
	{"Покупки онлайн", 55, "Купить ненужные вещи", 30, "shopping"},
	{"Новое железо", 70, "Заказать новую клавиатуру", 40, "shopping"},
	{"Курсы 'всё в одном'", 60, "Купить ещё один курс", 35, "shopping"},

	// Еда и здоровье
	{"Фастфуд", 50, "Съесть пиццу вместо здорового ужина", 20, "health"},
	{"Сладкое", 40, "Съесть шоколадку", 15, "health"},
	{"Алкоголь", 80, "Выпить пива 'для расслабления'", 50, "health"},

	// Прокрастинация
	{"Прокрастинация", 80, "Отложить изучение Go на завтра", 45, "procrastination"},
	{"Перфекционизм", 75, "Переписывать код вместо движения вперёд", 40, "procrastination"},
	{"Синдром самозванца", 70, "Читать форумы вместо написания кода", 35, "procrastination"},
	{"Поиск 'идеального' решения", 65, "Искать лучший фреймворк вместо работы", 35, "procrastination"},
}

// Босс-искушения
var bossTemptations = []Temptation{
	{"👹 CAPCUT МОНТЁР", 95, "Установи CapCut и монтируй 8 часов!", 80, "boss"},
	{"👹 ИГРОВОЙ ЗАВИСИМОН", 92, "Новая AAA-игра, просто 'один квест'!", 75, "boss"},
	{"👹 СОЦСЕТЕЙ ДЕМОНИУС", 88, "Листай ленту до 3 часов ночи!", 65, "boss"},
	{"👹 АЛКОГОЛЬНЫЙ ПРИЗРАК", 90, "Выпей 'для расслабления'!", 70, "boss"},
	{"👹 ДЕПРЕССИЯ МАКСИМА", 98, "Лежи и жалей себя целыми днями!", 85, "boss"},
}

var motivations = []Motivation{
	{"Каждая строка кода на Go — кирпичик в фундаменте карьеры", "knowledge+", 50},
	{"Сегодняшний дискомфорт — завтрашний комфорт зарплаты", "focus+", 40},
	{"Распыляться — значит стоять на месте. Фокус — значит расти", "willpower+", 45},
	{"Хобби подождут, когда будет стабильный доход", "focus+", 40},
	{"Гофер — символ эффективности и простоты", "knowledge+", 45},
	{"Каждый коммит приближает к офису с видом на город", "knowledge+", 50},
	{"Ошибки — не провалы, а инструкции к улучшению", "focus+", 35},
	{"Go учит не только программировать, но и мыслить системно", "knowledge+", 55},
	{"Сила воли в программировании важнее, чем в спортзале", "willpower+", 50},
	{"Экосистема Go — твой новый город возможностей", "knowledge+", 45},
	{"150К+ зарплаты ждут тех, кто не сдался", "willpower+", 60},
	{"Горутины масштабируются, как твоя карьера", "knowledge+", 55},
	{"Интерфейсы в Go проще, чем переговоры с заказчиком", "knowledge+", 40},
	{"Каждый день без Go — шаг назад к доставке", "focus+", 50},
	{"Твой GitHub — это твоё резюме. Пиши код!", "knowledge+", 65},
	{"Контекст не отменяется, и твоя мечта тоже", "willpower+", 45},
	{"Мьютексы блокируют гонки, а ты блокируй прокрастинацию", "focus+", 50},
	{"Defer, Panic, Recover — твой план на случай ошибок", "knowledge+", 40},
	{"Channel без буфера — как фокус без отвлечений", "focus+", 45},
	{"Go modules решают зависимости, а ты решай задачи", "knowledge+", 50},
}

func init() {
	rand.Seed(time.Now().UnixNano())
}

// GenerateTemptation генерирует случайное искушение
func GenerateTemptation() Temptation {
	index := rand.Intn(len(temptations))
	return temptations[index]
}

// GetRandomMotivation возвращает случайную мотивацию
func GetRandomMotivation() Motivation {
	index := rand.Intn(len(motivations))
	return motivations[index]
}

// GenerateBossTemptation генерирует босс-искушение
func GenerateBossTemptation() Temptation {
	index := rand.Intn(len(bossTemptations))
	return bossTemptations[index]
}

// GetTemptationsByCategory возвращает искушения по категории
func GetTemptationsByCategory(category string) []Temptation {
	var result []Temptation
	for _, t := range temptations {
		if t.Category == category {
			result = append(result, t)
		}
	}
	return result
}

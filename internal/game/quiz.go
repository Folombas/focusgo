package game

import (
	"fmt"
	"math/rand"
	"time"
)

// Question представляет вопрос викторины по Go
type Question struct {
	ID            int
	Question      string   // Текст вопроса
	Options       []string // Варианты ответов (4 варианта)
	CorrectAnswer int      // Индекс правильного ответа (0-3)
	Explanation   string   // Объяснение правильного ответа
	Category      string   // Категория: basics, concurrency, interfaces, errors, types
	Difficulty    int      // Сложность: 1 (лёгкий), 2 (средний), 3 (сложный)
	XPReward      int      // Награда за правильный ответ
}

// QuizSession представляет сессию викторины
type QuizSession struct {
	ChatID        int64
	Questions     []*Question
	CurrentIndex  int
	CorrectCount  int
	TotalCount    int
	TotalXP       int
	IsActive      bool
	StartedAt     time.Time
	CompletedAt   time.Time
	Category      string // Категория квиза или "mixed" для смешанного
	Difficulty    int    // Средняя сложность
}

// База вопросов по Go (50+ вопросов)
var goQuestionsBank = []*Question{
	// === ОСНОВЫ GO (basics) ===
	{
		ID:            1,
		Question:      "Как объявить переменную в Go?",
		Options:       []string{"var x int", "int x", "x : int", "dimension x as Integer"},
		CorrectAnswer: 0,
		Explanation:   "В Go переменные объявляются через ключевое слово var: var x int. Также можно использовать краткую форму x := 0",
		Category:      "basics",
		Difficulty:    1,
		XPReward:      10,
	},
	{
		ID:            2,
		Question:      "Что вернёт len(\"привет\") в Go?",
		Options:       []string{"6", "12", "4", "Ошибка компиляции"},
		CorrectAnswer: 1,
		Explanation:   "Функция len() возвращает количество байт в строке. Строка \"привет\" в UTF-8 занимает 12 байт (каждая кириллическая буква — 2 байта)",
		Category:      "basics",
		Difficulty:    2,
		XPReward:      15,
	},
	{
		ID:            3,
		Question:      "Какое ключевое слово используется для создания слайса?",
		Options:       []string{"slice", "array", "make", "new"},
		CorrectAnswer: 2,
		Explanation:   "Слайсы создаются через make([]T, length, capacity) или литерал []T{values}. Ключевое слово slice не существует в Go",
		Category:      "basics",
		Difficulty:    1,
		XPReward:      10,
	},
	{
		ID:            4,
		Question:      "Что выведет fmt.Printf(\"%T\", 42)?",
		Options:       []string{"int", "int64", "int32", "number"},
		CorrectAnswer: 0,
		Explanation:   "По умолчанию целочисленные литералы в Go имеют тип int. %T выводит тип значения",
		Category:      "basics",
		Difficulty:    1,
		XPReward:      10,
	},
	{
		ID:            5,
		Question:      "Как правильно создать map в Go?",
		Options:       []string{
			"map[string]int{}",
			"new map[string]int",
			"create map[string]int",
			"map[string]int.new()",
		},
		CorrectAnswer: 0,
		Explanation:   "Map создаётся через map[string]int{} или make(map[string]int). Ключевое слово new не используется для map",
		Category:      "basics",
		Difficulty:    1,
		XPReward:      10,
	},
	{
		ID:            6,
		Question:      "Что такое zero value в Go?",
		Options:       []string{
			"Значение по умолчанию для типа",
			"Нулевой указатель",
			"Ошибка компиляции",
			"Пустая строка",
		},
		CorrectAnswer: 0,
		Explanation:   "Zero value — это значение по умолчанию, которое получает переменная при объявлении без инициализации (0 для int, \"\" для string, false для bool)",
		Category:      "basics",
		Difficulty:    1,
		XPReward:      10,
	},
	{
		ID:            7,
		Question:      "Как объявить константу в Go?",
		Options:       []string{"const", "final", "define", "static"},
		CorrectAnswer: 0,
		Explanation:   "В Go константы объявляются через ключевое слово const: const Pi = 3.14",
		Category:      "basics",
		Difficulty:    1,
		XPReward:      10,
	},
	{
		ID:            8,
		Question:      "Что вернёт cap(make([]int, 5, 10))?",
		Options:       []string{"5", "10", "15", "Ошибка"},
		CorrectAnswer: 1,
		Explanation:   "Функция cap() возвращает ёмкость слайса. В данном случае ёмкость явно указана как 10",
		Category:      "basics",
		Difficulty:    2,
		XPReward:      15,
	},
	{
		ID:            9,
		Question:      "Как правильно импортировать пакет в Go?",
		Options:       []string{
			"import \"fmt\"",
			"using fmt",
			"include fmt",
			"require fmt",
		},
		CorrectAnswer: 0,
		Explanation:   "В Go импорты объявляются через ключевое слово import с путём в кавычках",
		Category:      "basics",
		Difficulty:    1,
		XPReward:      10,
	},
	{
		ID:            10,
		Question:      "Что такое blank identifier (_) в Go?",
		Options:       []string{
			"Игнорирование значения",
			"Приватная переменная",
			"Указатель nil",
			"Пустой интерфейс",
		},
		CorrectAnswer: 0,
		Explanation:   "Blank identifier (_) используется для игнорирования значений, которые не нужны: _, err := fn()",
		Category:      "basics",
		Difficulty:    1,
		XPReward:      10,
	},

	// === КОНКУРЕНТНОСТЬ (concurrency) ===
	{
		ID:            11,
		Question:      "Как создать горутину в Go?",
		Options:       []string{
			"go func()",
			"async func()",
			"thread func()",
			"spawn func()",
		},
		CorrectAnswer: 0,
		Explanation:   "Горутина создаётся добавлением ключевого слова go перед вызовом функции: go func()",
		Category:      "concurrency",
		Difficulty:    1,
		XPReward:      10,
	},
	{
		ID:            12,
		Question:      "Что такое channel в Go?",
		Options:       []string{
			"Типизированный канал для связи между горутинами",
			"Функция для ожидания горутин",
			"Мьютекс для синхронизации",
			"Таймер для отложенного выполнения",
		},
		CorrectAnswer: 0,
		Explanation:   "Channel — это типизированный канал для безопасной передачи данных между горутинами",
		Category:      "concurrency",
		Difficulty:    1,
		XPReward:      10,
	},
	{
		ID:            13,
		Question:      "Что делает select в Go?",
		Options:       []string{
			"Ждёт операции на нескольких каналах",
			"Выбирает случайное число",
			"Создаёт новую горутину",
			"Закрывает канал",
		},
		CorrectAnswer: 0,
		Explanation:   "Select позволяет ожидать операции на нескольких каналах одновременно и выполняет первую доступную",
		Category:      "concurrency",
		Difficulty:    2,
		XPReward:      15,
	},
	{
		ID:            14,
		Question:      "Как правильно закрыть channel?",
		Options:       []string{
			"close(ch)",
			"ch.close()",
			"delete ch",
			"ch = nil",
		},
		CorrectAnswer: 0,
		Explanation:   "Канал закрывается функцией close(ch). После закрытия в канал нельзя отправлять значения",
		Category:      "concurrency",
		Difficulty:    1,
		XPReward:      10,
	},
	{
		ID:            15,
		Question:      "Что такое sync.WaitGroup?",
		Options:       []string{
			"Счётчик для ожидания завершения горутин",
			"Канал для синхронизации",
			"Мьютекс для защиты данных",
			"Таймер для горутин",
		},
		CorrectAnswer: 0,
		Explanation:   "WaitGroup — это счётчик, который позволяет дождаться завершения всех зарегистрированных горутин",
		Category:      "concurrency",
		Difficulty:    2,
		XPReward:      15,
	},
	{
		ID:            16,
		Question:      "Что такое race condition?",
		Options:       []string{
			"Ситуация, когда результат зависит от порядка выполнения горутин",
			"Ошибка компиляции кода",
			"Бесконечный цикл",
			"Утечка памяти",
		},
		CorrectAnswer: 0,
		Explanation:   "Race condition возникает, когда несколько горутин обращаются к общим данным и хотя бы одна из них записывает",
		Category:      "concurrency",
		Difficulty:    2,
		XPReward:      15,
	},
	{
		ID:            17,
		Question:      "Как защитить общие данные от race condition?",
		Options:       []string{
			"Использовать мьютексы (sync.Mutex)",
			"Объявить переменные глобальными",
			"Использовать package unsafe",
			"Никак, Go сам решает",
		},
		CorrectAnswer: 0,
		Explanation:   "Для защиты общих данных используются мьютексы (sync.Mutex) или каналы для передачи владения",
		Category:      "concurrency",
		Difficulty:    2,
		XPReward:      15,
	},
	{
		ID:            18,
		Question:      "Что делает sync.Mutex.Lock()?",
		Options:       []string{
			"Блокирует мьютекс для эксклюзивного доступа",
			"Создаёт новую горутину",
			"Закрывает канал",
			"Ожидает таймер",
		},
		CorrectAnswer: 0,
		Explanation:   "Lock() захватывает мьютекс. Если мьютекс уже захвачен, вызов блокируется до его освобождения",
		Category:      "concurrency",
		Difficulty:    2,
		XPReward:      15,
	},
	{
		ID:            19,
		Question:      "Что такое buffered channel?",
		Options:       []string{
			"Канал с ёмкостью для хранения значений",
			"Канал без ёмкости",
			"Закрытый канал",
			"Канал только для чтения",
		},
		CorrectAnswer: 0,
		Explanation:   "Buffered channel имеет ёмкость: ch := make(chan int, 10). Отправка не блокируется, пока буфер не заполнен",
		Category:      "concurrency",
		Difficulty:    2,
		XPReward:      15,
	},
	{
		ID:            20,
		Question:      "Как отменить выполнение горутины?",
		Options:       []string{
			"Использовать context.Context",
			"Вызвать kill()",
			"Установить переменную в nil",
			"Никак, горутины работают до завершения",
		},
		CorrectAnswer: 0,
		Explanation:   "Для отмены горутин используется context.Context с отменой (context.WithCancel)",
		Category:      "concurrency",
		Difficulty:    3,
		XPReward:      20,
	},

	// === ИНТЕРФЕЙСЫ (interfaces) ===
	{
		ID:            21,
		Question:      "Как объявить интерфейс в Go?",
		Options:       []string{
			"type Interface interface { Method() }",
			"interface Interface { Method() }",
			"class Interface implements { Method() }",
			"abstract Interface { Method() }",
		},
		CorrectAnswer: 0,
		Explanation:   "Интерфейс объявляется через type Interface interface { ... } с перечислением методов",
		Category:      "interfaces",
		Difficulty:    1,
		XPReward:      10,
	},
	{
		ID:            22,
		Question:      "Что такое empty interface в Go?",
		Options:       []string{
			"interface{} — может хранить любое значение",
			"Интерфейс без методов",
			"Пустой интерфейс nil",
			"Интерфейс с одним методом",
		},
		CorrectAnswer: 0,
		Explanation:   "Empty interface interface{} не имеет методов и может хранить значения любого типа",
		Category:      "interfaces",
		Difficulty:    1,
		XPReward:      10,
	},
	{
		ID:            23,
		Question:      "Как Go реализует интерфейсы?",
		Options:       []string{
			"Неявно через структурную типизацию",
			"Через ключевое слово implements",
			"Через наследование",
			"Через аннотации",
		},
		CorrectAnswer: 0,
		Explanation:   "Go использует структурную типизацию: тип реализует интерфейс неявно, если реализует все его методы",
		Category:      "interfaces",
		Difficulty:    2,
		XPReward:      15,
	},
	{
		ID:            24,
		Question:      "Что такое type assertion?",
		Options:       []string{
			"Преобразование интерфейса к конкретному типу: v.(T)",
			"Проверка типа переменной",
			"Объявление нового типа",
			"Импорт типа из пакета",
		},
		CorrectAnswer: 0,
		Explanation:   "Type assertion извлекает конкретное значение из интерфейса: value, ok := iface.(Type)",
		Category:      "interfaces",
		Difficulty:    2,
		XPReward:      15,
	},
	{
		ID:            25,
		Question:      "Что такое type switch?",
		Options:       []string{
			"switch v.(type) { case int: ... }",
			"switch typeof(v) { ... }",
			"switch v.type { ... }",
			"switch cast(v) { ... }",
		},
		CorrectAnswer: 0,
		Explanation:   "Type switch позволяет выполнить разные ветки кода в зависимости от динамического типа интерфейса",
		Category:      "interfaces",
		Difficulty:    2,
		XPReward:      15,
	},
	{
		ID:            26,
		Question:      "Что такое io.Reader?",
		Options:       []string{
			"Интерфейс с методом Read(p []byte) (int, error)",
			"Функция для чтения файлов",
			"Структура для работы с файлами",
			"Пакет для ввода-вывода",
		},
		CorrectAnswer: 0,
		Explanation:   "io.Reader — базовый интерфейс для чтения данных: Read(p []byte) (n int, err error)",
		Category:      "interfaces",
		Difficulty:    2,
		XPReward:      15,
	},
	{
		ID:            27,
		Question:      "Что такое io.Writer?",
		Options:       []string{
			"Интерфейс с методом Write(p []byte) (int, error)",
			"Функция для записи файлов",
			"Структура для работы с файлами",
			"Пакет для вывода данных",
		},
		CorrectAnswer: 0,
		Explanation:   "io.Writer — базовый интерфейс для записи данных: Write(p []byte) (n int, err error)",
		Category:      "interfaces",
		Difficulty:    2,
		XPReward:      15,
	},
	{
		ID:            28,
		Question:      "Что такое io.Closer?",
		Options:       []string{
			"Интерфейс с методом Close() error",
			"Функция для закрытия файлов",
			"Пакет для управления ресурсами",
			"Структура для освобождения памяти",
		},
		CorrectAnswer: 0,
		Explanation:   "io.Closer — простой интерфейс с одним методом Close() error для освобождения ресурсов",
		Category:      "interfaces",
		Difficulty:    1,
		XPReward:      10,
	},
	{
		ID:            29,
		Question:      "Что такое io.ReadWriter?",
		Options:       []string{
			"Комбинация io.Reader и io.Writer",
			"Интерфейс только для чтения",
			"Интерфейс только для записи",
			"Пакет для работы с файлами",
		},
		CorrectAnswer: 0,
		Explanation:   "io.ReadWriter объединяет интерфейсы Reader и Writer в одном типе",
		Category:      "interfaces",
		Difficulty:    2,
		XPReward:      15,
	},
	{
		ID:            30,
		Question:      "Что такое fmt.Stringer?",
		Options:       []string{
			"Интерфейс с методом String() string",
			"Функция для преобразования в строку",
			"Пакет для форматирования строк",
			"Тип для работы с текстом",
		},
		CorrectAnswer: 0,
		Explanation:   "fmt.Stringer — интерфейс с методом String() string, используется для строкового представления",
		Category:      "interfaces",
		Difficulty:    2,
		XPReward:      15,
	},

	// === ОБРАБОТКА ОШИБОК (errors) ===
	{
		ID:            31,
		Question:      "Как правильно обработать ошибку в Go?",
		Options:       []string{
			"if err != nil { return err }",
			"try { ... } catch (err) { ... }",
			"throw err",
			"on error resume next",
		},
		CorrectAnswer: 0,
		Explanation:   "В Go ошибки обрабатываются явной проверкой: if err != nil { ... }",
		Category:      "errors",
		Difficulty:    1,
		XPReward:      10,
	},
	{
		ID:            32,
		Question:      "Что такое error в Go?",
		Options:       []string{
			"Встроенный интерфейс с методом Error() string",
			"Класс исключений",
			"Тип данных для ошибок",
			"Макрос для обработки ошибок",
		},
		CorrectAnswer: 0,
		Explanation:   "error — встроенный интерфейс: type error interface { Error() string }",
		Category:      "errors",
		Difficulty:    1,
		XPReward:      10,
	},
	{
		ID:            33,
		Question:      "Как создать пользовательскую ошибку?",
		Options:       []string{
			"errors.New(\"message\") или fmt.Errorf(\"format\", args)",
			"new Error(\"message\")",
			"throw new Error(\"message\")",
			"raise \"message\"",
		},
		CorrectAnswer: 0,
		Explanation:   "Ошибки создаются через errors.New() или fmt.Errorf() для форматирования",
		Category:      "errors",
		Difficulty:    1,
		XPReward:      10,
	},
	{
		ID:            34,
		Question:      "Что делает errors.Wrap() из pkg/errors?",
		Options:       []string{
			"Оборачивает ошибку с добавлением контекста",
			"Создаёт новую ошибку",
			"Печатает стек вызовов",
			"Игнорирует ошибку",
		},
		CorrectAnswer: 0,
		Explanation:   "errors.Wrap() оборачивает ошибку, добавляя контекст и стек вызовов (в Go 1.13+ есть errors.As/Is)",
		Category:      "errors",
		Difficulty:    3,
		XPReward:      20,
	},
	{
		ID:            35,
		Question:      "Что такое sentinel error?",
		Options:       []string{
			"Глобальная переменная ошибки для сравнения",
			"Ошибка с кодом состояния",
			"Паника вместо ошибки",
			"Ошибка времени выполнения",
		},
		CorrectAnswer: 0,
		Explanation:   "Sentinel error — это глобальная переменная ошибки: var ErrNotFound = errors.New(\"not found\")",
		Category:      "errors",
		Difficulty:    3,
		XPReward:      20,
	},
	{
		ID:            36,
		Question:      "Как проверить тип ошибки в Go 1.13+?",
		Options:       []string{
			"errors.As(err, &targetType)",
			"err instanceof Type",
			"type(err) == Type",
			"err is Type",
		},
		CorrectAnswer: 0,
		Explanation:   "errors.As() проверяет, можно ли привести ошибку к целевому типу",
		Category:      "errors",
		Difficulty:    3,
		XPReward:      20,
	},
	{
		ID:            37,
		Question:      "Что делает errors.Is()?",
		Options:       []string{
			"Проверяет, равна ли ошибка целевой (с учётом обёрток)",
			"Сравнивает два error",
			"Создаёт новую ошибку",
			"Печатает ошибку",
		},
		CorrectAnswer: 0,
		Explanation:   "errors.Is() рекурсивно проверяет цепочку обёрнутых ошибок на равенство целевой",
		Category:      "errors",
		Difficulty:    3,
		XPReward:      20,
	},
	{
		ID:            38,
		Question:      "Что такое panic в Go?",
		Options:       []string{
			"Невосстанавливаемая ошибка времени выполнения",
			"Обычная ошибка",
			"Предупреждение компилятора",
			"Тип исключения",
		},
		CorrectAnswer: 0,
		Explanation:   "panic — это невосстанавливаемая ошибка, которая прерывает выполнение программы (если не recover)",
		Category:      "errors",
		Difficulty:    2,
		XPReward:      15,
	},
	{
		ID:            39,
		Question:      "Что делает defer?",
		Options:       []string{
			"Откладывает выполнение функции до выхода из текущей",
			"Запускает горутину",
			"Ожидает завершения функции",
			"Повторяет вызов функции",
		},
		CorrectAnswer: 0,
		Explanation:   "defer откладывает вызов функции до момента выхода из окружающей функции",
		Category:      "errors",
		Difficulty:    2,
		XPReward:      15,
	},
	{
		ID:            40,
		Question:      "Как восстановить программу после panic?",
		Options:       []string{
			"Использовать recover() в defer",
			"try-catch блок",
			"Никак, программа завершится",
			"Перезапустить горутину",
		},
		CorrectAnswer: 0,
		Explanation:   "recover() вызывается внутри defer-функции для восстановления после panic",
		Category:      "errors",
		Difficulty:    3,
		XPReward:      20,
	},

	// === ТИПЫ ДАННЫХ (types) ===
	{
		ID:            41,
		Question:      "Какой размер int в Go?",
		Options:       []string{
			"Зависит от архитектуры (32 или 64 бита)",
			"Всегда 32 бита",
			"Всегда 64 бита",
			"16 бит",
		},
		CorrectAnswer: 0,
		Explanation:   "int имеет размер, зависящий от архитектуры: 32 бита на 32-битных системах, 64 бита на 64-битных",
		Category:      "types",
		Difficulty:    2,
		XPReward:      15,
	},
	{
		ID:            42,
		Question:      "Что такое rune в Go?",
		Options:       []string{
			"Псевдоним для int32, представляет Unicode code point",
			"Символ ASCII",
			"Строка из одного символа",
			"Байт текста",
		},
		CorrectAnswer: 0,
		Explanation:   "rune — это псевдоним для int32, используется для представления Unicode code points",
		Category:      "types",
		Difficulty:    2,
		XPReward:      15,
	},
	{
		ID:            43,
		Question:      "Что такое byte в Go?",
		Options:       []string{
			"Псевдоним для uint8",
			"Псевдоним для int8",
			"Отдельный тип данных",
			"Строка из одного символа",
		},
		CorrectAnswer: 0,
		Explanation:   "byte — это псевдоним для uint8, представляет один байт данных",
		Category:      "types",
		Difficulty:    1,
		XPReward:      10,
	},
	{
		ID:            44,
		Question:      "Как преобразовать int в string?",
		Options:       []string{
			"strconv.Itoa(int) или fmt.Sprintf(\"%d\", int)",
			"string(int)",
			"int.toString()",
			"text(int)",
		},
		CorrectAnswer: 0,
		Explanation:   "Для преобразования int в string используется strconv.Itoa() или fmt.Sprintf()",
		Category:      "types",
		Difficulty:    2,
		XPReward:      15,
	},
	{
		ID:            45,
		Question:      "Что такое struct в Go?",
		Options:       []string{
			"Составной тип данных с полями",
			"Класс с методами",
			"Интерфейс с методами",
			"Массив фиксированного размера",
		},
		CorrectAnswer: 0,
		Explanation:   "struct — это составной тип, объединяющий поля разных типов в одну структуру",
		Category:      "types",
		Difficulty:    1,
		XPReward:      10,
	},
	{
		ID:            46,
		Question:      "Что такое pointer в Go?",
		Options:       []string{
			"Переменная, хранящая адрес другой переменной",
			"Указатель на функцию",
			"Ссылка на объект",
			"Адрес в памяти",
		},
		CorrectAnswer: 0,
		Explanation:   "Pointer (*) хранит адрес переменной. & получает адрес, * разыменовывает указатель",
		Category:      "types",
		Difficulty:    2,
		XPReward:      15,
	},
	{
		ID:            47,
		Question:      "Что такое array в Go?",
		Options:       []string{
			"Массив фиксированного размера [N]T",
			"Динамический массив",
			"Список значений",
			"Коллекция элементов",
		},
		CorrectAnswer: 0,
		Explanation:   "Array — это массив фиксированного размера [N]T. Размер является частью типа",
		Category:      "types",
		Difficulty:    2,
		XPReward:      15,
	},
	{
		ID:            48,
		Question:      "Чем slice отличается от array?",
		Options:       []string{
			"Slice — динамический, array — фиксированный размер",
			"Ничем, это одно и то же",
			"Slice хранит указатели, array — значения",
			"Array быстрее slice",
		},
		CorrectAnswer: 0,
		Explanation:   "Slice ([]T) — это динамическое представление array, может расти через append",
		Category:      "types",
		Difficulty:    2,
		XPReward:      15,
	},
	{
		ID:            49,
		Question:      "Что такое variadic function?",
		Options:       []string{
			"Функция с переменным числом аргументов: fn(args ...T)",
			"Функция без аргументов",
			"Функция с одним аргументом",
			"Функция с аргументами по умолчанию",
		},
		CorrectAnswer: 0,
		Explanation:   "Variadic function принимает переменное число аргументов одного типа: func fn(args ...int)",
		Category:      "types",
		Difficulty:    2,
		XPReward:      15,
	},
	{
		ID:            50,
		Question:      "Что такое type alias в Go?",
		Options:       []string{
			"type Alias = OriginalType",
			"alias Alias for OriginalType",
			"using Alias = OriginalType",
			"define Alias OriginalType",
		},
		CorrectAnswer: 0,
		Explanation:   "Type alias создаёт псевдоним типа: type Alias = OriginalType (Go 1.9+)",
		Category:      "types",
		Difficulty:    3,
		XPReward:      20,
	},
}

// NewQuizSession создаёт новую сессию викторины
func NewQuizSession(chatID int64, category string, count int) *QuizSession {
	rand.Seed(time.Now().UnixNano())

	session := &QuizSession{
		ChatID:       chatID,
		CurrentIndex: 0,
		CorrectCount: 0,
		TotalCount:   count,
		TotalXP:      0,
		IsActive:     true,
		StartedAt:    time.Now(),
		Category:     category,
	}

	// Выбираем вопросы
	session.Questions = selectQuestions(category, count)
	session.TotalCount = len(session.Questions)

	return session
}

// selectQuestions выбирает вопросы по категории
func selectQuestions(category string, count int) []*Question {
	var filtered []*Question

	for _, q := range goQuestionsBank {
		if category == "mixed" || q.Category == category {
			filtered = append(filtered, q)
		}
	}

	// Перемешиваем и выбираем count вопросов
	rand.Shuffle(len(filtered), func(i, j int) {
		filtered[i], filtered[j] = filtered[j], filtered[i]
	})

	if count > len(filtered) {
		count = len(filtered)
	}

	return filtered[:count]
}

// GetCurrentQuestion возвращает текущий вопрос
func (qs *QuizSession) GetCurrentQuestion() *Question {
	if qs.CurrentIndex >= len(qs.Questions) {
		return nil
	}
	return qs.Questions[qs.CurrentIndex]
}

// AnswerQuestion обрабатывает ответ на вопрос
func (qs *QuizSession) AnswerQuestion(answerIndex int) (bool, string, int) {
	question := qs.GetCurrentQuestion()
	if question == nil {
		return false, "Викторина завершена", 0
	}

	isCorrect := answerIndex == question.CorrectAnswer

	if isCorrect {
		qs.CorrectCount++
		qs.TotalXP += question.XPReward
	}

	qs.CurrentIndex++

	var result string
	if isCorrect {
		result = fmt.Sprintf("✅ <b>ПРАВИЛЬНО!</b>\n\n%s\n\n📘 <b>Объяснение:</b>\n%s",
			question.Question, question.Explanation)
	} else {
		result = fmt.Sprintf("❌ <b>НЕПРАВИЛЬНО!</b>\n\nПравильный ответ: %s\n\n%s\n\n📘 <b>Объяснение:</b>\n%s",
			question.Options[question.CorrectAnswer], question.Question, question.Explanation)
	}

	return isCorrect, result, question.XPReward
}

// IsFinished проверяет завершение викторины
func (qs *QuizSession) IsFinished() bool {
	return qs.CurrentIndex >= qs.TotalCount
}

// GetResults возвращает результаты викторины
func (qs *QuizSession) GetResults() string {
	percentage := 0
	if qs.TotalCount > 0 {
		percentage = (qs.CorrectCount * 100) / qs.TotalCount
	}

	rating := "🔰 Новичок"
	if percentage >= 90 {
		rating = "🏆 Go-Эксперт"
	} else if percentage >= 70 {
		rating = "🎓 Go-Продвинутый"
	} else if percentage >= 50 {
		rating = "📚 Go-Студент"
	}

	return fmt.Sprintf(`🧩 <b>РЕЗУЛЬТАТЫ ВИКТОРИНЫ</b>
━━━━━━━━━━━━━━━━━━━━

✅ Правильных ответов: %d из %d
📊 Процент правильности: %d%%
✨ Заработано опыта: %d
🏅 Рейтинг: %s

⏱️ Время прохождения: %v`,
		qs.CorrectCount, qs.TotalCount, percentage, qs.TotalXP, rating,
		time.Since(qs.StartedAt).Round(time.Second))
}

// GetCategoryQuestions возвращает количество вопросов в категории
func GetCategoryQuestions(category string) int {
	count := 0
	for _, q := range goQuestionsBank {
		if category == "mixed" || q.Category == category {
			count++
		}
	}
	return count
}

// GetAllCategories возвращает все доступные категории
func GetAllCategories() []string {
	categories := []string{"mixed", "basics", "concurrency", "interfaces", "errors", "types"}
	return categories
}

// GetCategoryName возвращает название категории на русском
func GetCategoryName(category string) string {
	names := map[string]string{
		"mixed":       "🔀 Все темы",
		"basics":      "📘 Основы Go",
		"concurrency": "⚡ Конкурентность",
		"interfaces":  "🔌 Интерфейсы",
		"errors":      "⚠️ Обработка ошибок",
		"types":       "📊 Типы данных",
	}

	if name, exists := names[category]; exists {
		return name
	}
	return category
}

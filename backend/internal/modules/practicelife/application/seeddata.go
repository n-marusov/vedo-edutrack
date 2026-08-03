// Package practicelife provides the practicelife bounded context application layer.
//
// This file holds the M2 launch content catalog: stories (why knowledge
// matters in the real world) and cross-subject project ideas. The catalog is
// ontology-sourced in production; launch content is bundled as seed data so
// the product is demonstrable before the VEDO Hub starter ontology is ready
// (M2 exit criteria: ≥ 50 stories, ≥ 30 project ideas).
package practicelife

import domain "vedo-edutrack/backend/internal/modules/practicelife/domain"

// LaunchStories returns the M2 launch story catalog (≥ 50 stories across
// grades 5–11 subjects). Every story has 1–3 linked modules and a mandatory
// real-world application section.
func LaunchStories() []domain.Story {
	return []domain.Story{
		// --- Mathematics (8) ---
		{ID: "story-math-percent", Title: "Проценты в жизни", Text: "Скидки, банковские вклады, кредиты — всё это проценты.", LinkedModules: []string{"percent", "math-5-11"}, Subjects: []string{"Математика"}, RealWorld: "Аналитики считают процентные изменения курсов валют и инфляции.", ReadingMinutes: 3},
		{ID: "story-math-fractions", Title: "Дроби на кухне", Text: "Рецепты требуют точных пропорций — это дробные отношения.", LinkedModules: []string{"fractions"}, Subjects: []string{"Математика"}, RealWorld: "Повара и фармацевты ежедневно работают с дробными пропорциями.", ReadingMinutes: 3},
		{ID: "story-math-equations", Title: "Уравнения в инженерии", Text: "Баланс нагрузок в мостах описывается линейными уравнениями.", LinkedModules: []string{"linear-equations", "math-7-1"}, Subjects: []string{"Математика", "Физика"}, RealWorld: "Инженеры решают системы уравнений при проектировании конструкций.", ReadingMinutes: 4},
		{ID: "story-math-functions", Title: "Функции вокруг нас", Text: "Графики функций описывают рост цен, движение и прогресс.", LinkedModules: []string{"functions", "math-7-3"}, Subjects: []string{"Математика"}, RealWorld: "Экономисты строят функции спроса и предложения для прогнозов.", ReadingMinutes: 4},
		{ID: "story-math-geometry", Title: "Геометрия в архитектуре", Text: "Пифагор, углы и подобие — основа чертежей зданий.", LinkedModules: []string{"geometry", "pythagorean-theorem"}, Subjects: []string{"Математика", "ИЗО"}, RealWorld: "Архитекторы используют теорему Пифагора для расчёта скатов крыш.", ReadingMinutes: 4},
		{ID: "story-math-probability", Title: "Вероятность в играх и жизни", Text: "Шансы выигрыша в лотереях и прогнозы погоды — это вероятность.", LinkedModules: []string{"probability", "math-9-1"}, Subjects: []string{"Математика"}, RealWorld: "Страховые компании оценивают риски через теорию вероятностей.", ReadingMinutes: 4},
		{ID: "story-math-logarithms", Title: "Логарифмы: от музыки до сейсмологии", Text: "Равномерно темперированный строй и шкала Рихтера — логарифмические.", LinkedModules: []string{"logarithms", "math-10-2"}, Subjects: []string{"Математика", "Физика"}, RealWorld: "Сейсмологи измеряют силу землетрясений по логарифмической шкале.", ReadingMinutes: 5},
		{ID: "story-math-vectors", Title: "Векторы в навигации", Text: "Скорость и направление ветра — векторные величины.", LinkedModules: []string{"vectors", "math-9-4"}, Subjects: []string{"Математика", "Физика"}, RealWorld: "Пилоты и моряки складывают векторы скоростей для курса.", ReadingMinutes: 4},

		// --- Biology (7) ---
		{ID: "story-bio-photosynthesis", Title: "Фотосинтез — химия в каждом листе", Text: "Растения превращают свет и углекислый газ в глюкозу.", LinkedModules: []string{"photosynthesis", "chemistry"}, Subjects: []string{"Биология", "Химия"}, RealWorld: "Понимание фотосинтеза ведёт к созданию искусственного фотосинтеза для топлива.", ReadingMinutes: 5},
		{ID: "story-bio-cells", Title: "Клетка — мини-фабрика", Text: "Митохондрии, рибосомы, мембраны — каждый органоид на своём месте.", LinkedModules: []string{"cells", "bio-8-1"}, Subjects: []string{"Биология"}, RealWorld: "Клеточные технологии используются для выращивания тканей и лекарств.", ReadingMinutes: 4},
		{ID: "story-bio-genetics", Title: "Генетика: почему мы похожи", Text: "ДНК кодирует наследственные признаки от цвета глаз до болезней.", LinkedModules: []string{"genetics", "bio-9-2"}, Subjects: []string{"Биология"}, RealWorld: "Генетические тесты помогают подбирать лечение и оценивать риски.", ReadingMinutes: 4},
		{ID: "story-bio-ecology", Title: "Экология двора и города", Text: "Пищевые цепи и круговорот веществ работают вокруг нас.", LinkedModules: []string{"ecology", "bio-7-3"}, Subjects: []string{"Биология", "География"}, RealWorld: "Городские экологи рассчитывают, сколько деревьев нужно для чистого воздуха.", ReadingMinutes: 4},
		{ID: "story-bio-anatomy", Title: "Как работает тело", Text: "Кровь, сердце, лёгкие — единая транспортная система.", LinkedModules: []string{"human-anatomy", "bio-8-3"}, Subjects: []string{"Биология", "Физика"}, RealWorld: "Кардиологи используют физику давления для диагностики сердца.", ReadingMinutes: 4},
		{ID: "story-bio-population", Title: "Популяции и статистика", Text: "Численность популяций считают с помощью процентов и выборок.", LinkedModules: []string{"population-genetics", "percent"}, Subjects: []string{"Биология", "Математика"}, RealWorld: "Биологи-статистики оценивают численность видов по выборочным данным.", ReadingMinutes: 4},
		{ID: "story-bio-evolution", Title: "Эволюция: отбор в действии", Text: "Естественный отбор объясняет разнообразие видов.", LinkedModules: []string{"evolution", "bio-9-5"}, Subjects: []string{"Биология"}, RealWorld: "Эволюционные алгоритмы используются в программировании для оптимизации.", ReadingMinutes: 4},

		// --- Physics (7) ---
		{ID: "story-phys-mechanics", Title: "Механика: почему тела движутся", Text: "Силы, масса и ускорение — три кита классической механики.", LinkedModules: []string{"mechanics", "vectors"}, Subjects: []string{"Физика", "Математика"}, RealWorld: "Проектировщики автомобилей рассчитывают силу торможения.", ReadingMinutes: 4},
		{ID: "story-phys-optics", Title: "Оптика: от линз до лазеров", Text: "Преломление света лежит в основе очков, камер и микроскопов.", LinkedModules: []string{"optics", "phys-8-2"}, Subjects: []string{"Физика"}, RealWorld: "Лазерная коррекция зрения — практическое применение оптики.", ReadingMinutes: 4},
		{ID: "story-phys-electricity", Title: "Электричество в каждом доме", Text: "Ток, напряжение и сопротивление — язык электросетей.", LinkedModules: []string{"electricity", "phys-8-4"}, Subjects: []string{"Физика"}, RealWorld: "Энергетики балансируют сети по закону Ома.", ReadingMinutes: 4},
		{ID: "story-phys-waves", Title: "Волны: звук и свет", Text: "Всё — от голоса до Wi-Fi — распространяется волнами.", LinkedModules: []string{"waves", "phys-9-3"}, Subjects: []string{"Физика"}, RealWorld: "Инженеры связи кодируют данные в электромагнитные волны.", ReadingMinutes: 4},
		{ID: "story-phys-thermo", Title: "Термодинамика: тепло и энергия", Text: "Тепловые двигатели превращают тепло в работу.", LinkedModules: []string{"thermodynamics", "phys-10-1"}, Subjects: []string{"Физика"}, RealWorld: "Теплоэлектростанции и холодильники работают по законам термодинамики.", ReadingMinutes: 4},
		{ID: "story-phys-pressure", Title: "Давление: от атмосферы до воды", Text: "Атмосферное давление объясняет, почему летают самолёты.", LinkedModules: []string{"pressure", "phys-7-2"}, Subjects: []string{"Физика", "География"}, RealWorld: "Метеорологи прогнозируют погоду по изменению давления.", ReadingMinutes: 3},
		{ID: "story-phys-quantum", Title: "Квантовая физика и технологии", Text: "Квантовые эффекты управляют транзисторами и лазерами.", LinkedModules: []string{"quantum", "phys-11-2"}, Subjects: []string{"Физика"}, RealWorld: "Современные процессоры работают на квантовых принципах полупроводников.", ReadingMinutes: 5},

		// --- Chemistry (6) ---
		{ID: "story-chem-solutions", Title: "Растворы вокруг нас", Text: "Концентрация растворов в медицине и кулинарии.", LinkedModules: []string{"solutions", "chemistry"}, Subjects: []string{"Химия"}, RealWorld: "Физраствор, уксус и морская вода — растворы разной концентрации.", ReadingMinutes: 4},
		{ID: "story-chem-reactions", Title: "Химические реакции в природе", Text: "Окисление, горение, фотосинтез — реакции вокруг нас.", LinkedModules: []string{"chemical-reactions", "chemistry"}, Subjects: []string{"Химия", "Биология"}, RealWorld: "Коррозия металлов — химическая реакция, стоящая миллиарды.", ReadingMinutes: 4},
		{ID: "story-chem-atoms", Title: "Атомы: из чего всё состоит", Text: "Периодическая таблица — карта всех известных элементов.", LinkedModules: []string{"atoms", "periodic-table"}, Subjects: []string{"Химия"}, RealWorld: "Материаловеды создают сплавы, комбинируя элементы.", ReadingMinutes: 4},
		{ID: "story-chem-acids", Title: "Кислоты и основания в быту", Text: "Лимонный сок, сода, мыло — кислоты и основания вокруг нас.", LinkedModules: []string{"acids-bases", "chem-8-3"}, Subjects: []string{"Химия"}, RealWorld: "Химики регулируют pH почвы и воды для сельского хозяйства.", ReadingMinutes: 3},
		{ID: "story-chem-organic", Title: "Органическая химия: химия жизни", Text: "Углеводороды, белки, полимеры — молекулы жизни.", LinkedModules: []string{"organic-chemistry", "chem-10-2"}, Subjects: []string{"Химия", "Биология"}, RealWorld: "Пластик, лекарства и топливо — продукты органической химии.", ReadingMinutes: 4},
		{ID: "story-chem-electrolysis", Title: "Электролиз: ток творит химию", Text: "Электрический ток разлагает вещества на составляющие.", LinkedModules: []string{"electrolysis", "electricity"}, Subjects: []string{"Химия", "Физика"}, RealWorld: "Электролиз используется для получения алюминия и очистки металлов.", ReadingMinutes: 4},

		// --- History (6) ---
		{ID: "story-hist-ancient", Title: "Древний мир: колыбель цивилизаций", Text: "Египет, Греция, Рим — основы права и науки.", LinkedModules: []string{"ancient-history", "hist-5-1"}, Subjects: []string{"История"}, RealWorld: "Римское право лежит в основе современных кодексов.", ReadingMinutes: 4},
		{ID: "story-hist-middle", Title: "Средневековье: замки и университеты", Text: "Феодализм, цеха и первые университеты.", LinkedModules: []string{"middle-ages", "hist-6-2"}, Subjects: []string{"История"}, RealWorld: "Университетская система — наследие средневековых школ.", ReadingMinutes: 4},
		{ID: "story-hist-modern", Title: "Новое время: революции и наука", Text: "Научная революция изменила представление о мире.", LinkedModules: []string{"modern-history", "hist-7-3"}, Subjects: []string{"История"}, RealWorld: "Метод эксперимента, рождённый в Новое время, — основа науки.", ReadingMinutes: 4},
		{ID: "story-hist-wars", Title: "Мировые войны и XX век", Text: "Технологии и общество изменились после двух мировых войн.", LinkedModules: []string{"world-wars", "hist-9-2"}, Subjects: []string{"История"}, RealWorld: "Организация Объединённых Наций — ответ на войны XX века.", ReadingMinutes: 5},
		{ID: "story-hist-culture", Title: "Культура и искусство в истории", Text: "Искусство отражает эпоху: от наскальных рисунков до авангарда.", LinkedModules: []string{"cultural-history", "hist-8-1"}, Subjects: []string{"История", "Литература"}, RealWorld: "Реставраторы используют исторические знания для сохранения памятников.", ReadingMinutes: 4},
		{ID: "story-hist-archeology", Title: "Археология: чтение по слоям", Text: "Раскопки открывают быт древних людей.", LinkedModules: []string{"archeology", "hist-5-4"}, Subjects: []string{"История"}, RealWorld: "Археологи датируют находки радиоуглеродным методом — физикой и химией.", ReadingMinutes: 4},

		// --- Literature (5) ---
		{ID: "story-lit-classics", Title: "Русская классика: зеркало души", Text: "Толстой, Достоевский, Чехов — о человеке и обществе.", LinkedModules: []string{"russian-classics", "lit-9-1"}, Subjects: []string{"Литература"}, RealWorld: "Сценаристы и психологи обращаются к классике за пониманием человека.", ReadingMinutes: 4},
		{ID: "story-lit-poetry", Title: "Поэзия: ритм и смысл", Text: "Размер, рифма и образ — инструменты поэта.", LinkedModules: []string{"poetry", "lit-6-2"}, Subjects: []string{"Литература"}, RealWorld: "Авторы песен используют поэтические приёмы для создания хитов.", ReadingMinutes: 3},
		{ID: "story-lit-fairytale", Title: "Сказки: мудрость в метафорах", Text: "Народные сказки кодируют жизненные уроки.", LinkedModules: []string{"fairy-tales", "lit-5-1"}, Subjects: []string{"Литература"}, RealWorld: "Психологи используют сказки в терапии и маркетинге.", ReadingMinutes: 3},
		{ID: "story-lit-drama", Title: "Драматургия: конфликт на сцене", Text: "Пьесы строятся на конфликтах и диалогах.", LinkedModules: []string{"drama", "lit-8-3"}, Subjects: []string{"Литература"}, RealWorld: "Сценарии фильмов и сериалов — современная драматургия.", ReadingMinutes: 4},
		{ID: "story-lit-science-fiction", Title: "Фантастика: наука вперёд времени", Text: "Жюль Верн и Брэдбери предсказали технологии.", LinkedModules: []string{"science-fiction", "lit-10-1"}, Subjects: []string{"Литература", "Физика"}, RealWorld: "Идеи фантастов вдохновляют инженеров: от подводных лодок до роботов.", ReadingMinutes: 4},

		// --- Geography (6) ---
		{ID: "story-geo-continents", Title: "Континенты: планы и границы", Text: "Дрейф материков объясняет расположение континентов.", LinkedModules: []string{"continents", "geo-5-1"}, Subjects: []string{"География"}, RealWorld: "Логистические компании планируют маршруты через континенты.", ReadingMinutes: 3},
		{ID: "story-geo-climate", Title: "Климат и климатические зоны", Text: "Почему на экваторе жарко, а у полюсов холодно.", LinkedModules: []string{"climate", "geo-6-2"}, Subjects: []string{"География"}, RealWorld: "Климатологи моделируют изменения климата и их последствия.", ReadingMinutes: 4},
		{ID: "story-geo-demography", Title: "Демография: население мира", Text: "Рождаемость, смертность, миграция — демографические процессы.", LinkedModules: []string{"demography", "percent"}, Subjects: []string{"География", "Математика"}, RealWorld: "Правительства планируют школы и больницы по демографическим данным.", ReadingMinutes: 4},
		{ID: "story-geo-rivers", Title: "Реки: артерии цивилизаций", Text: "Великие цивилизации выросли на берегах рек.", LinkedModules: []string{"rivers", "geo-6-4"}, Subjects: []string{"География", "История"}, RealWorld: "Гидроэлектростанции используют энергию рек.", ReadingMinutes: 4},
		{ID: "story-geo-map", Title: "Карты: как мы ориентируемся", Text: "Проекции и масштаб — язык картографии.", LinkedModules: []string{"maps", "geo-5-3"}, Subjects: []string{"География", "Математика"}, RealWorld: "GPS-навигация основана на картографических проекциях.", ReadingMinutes: 3},
		{ID: "story-geo-resources", Title: "Природные ресурсы и экономика", Text: "Нефть, газ, металлы — ресурсы определяют экономики стран.", LinkedModules: []string{"natural-resources", "geo-9-2"}, Subjects: []string{"География", "Химия"}, RealWorld: "Геологи и химики оценивают запасы и качество руд.", ReadingMinutes: 4},

		// --- Computer Science (5) ---
		{ID: "story-cs-algorithms", Title: "Алгоритмы: инструкции для машин", Text: "Поиск, сортировка, сжатие — базовые алгоритмы.", LinkedModules: []string{"algorithms", "cs-8-1"}, Subjects: []string{"Информатика"}, RealWorld: "Поисковые системы ранжируют страницы алгоритмами.", ReadingMinutes: 4},
		{ID: "story-cs-programming", Title: "Программирование: разговор с машиной", Text: "Языки программирования — способ управлять компьютерами.", LinkedModules: []string{"programming-basics", "cs-7-1"}, Subjects: []string{"Информатика"}, RealWorld: "Программисты создают приложения, которыми пользуются миллиарды.", ReadingMinutes: 4},
		{ID: "story-cs-python", Title: "Python: язык для всех", Text: "От школьных задач до искусственного интеллекта.", LinkedModules: []string{"python", "cs-9-3"}, Subjects: []string{"Информатика"}, RealWorld: "Python — основной язык анализа данных и машинного обучения.", ReadingMinutes: 4},
		{ID: "story-cs-data", Title: "Данные: новая нефть", Text: "Статистика и визуализация помогают понимать большие данные.", LinkedModules: []string{"data-analysis", "statistics"}, Subjects: []string{"Информатика", "Математика"}, RealWorld: "Аналитики данных находят закономерности в поведении пользователей.", ReadingMinutes: 4},
		{ID: "story-cs-security", Title: "Кибербезопасность: защита информации", Text: "Шифрование, пароли, фишинг — основы цифровой безопасности.", LinkedModules: []string{"cybersecurity", "cs-10-2"}, Subjects: []string{"Информатика", "Обществознание"}, RealWorld: "Специалисты по безопасности защищают банки и государство от атак.", ReadingMinutes: 4},

		// --- Social Studies (5) ---
		{ID: "story-soc-economics", Title: "Экономика: как работает рынок", Text: "Спрос, предложение и цена — базовые силы рынка.", LinkedModules: []string{"economics", "social-8-1"}, Subjects: []string{"Обществознание", "Математика"}, RealWorld: "Финансовые аналитики прогнозируют рынки по экономическим моделям.", ReadingMinutes: 4},
		{ID: "story-soc-taxes", Title: "Налоги: откуда деньги у государства", Text: "НДФЛ, НДС и налоги — финансирование общественных благ.", LinkedModules: []string{"taxes", "percent"}, Subjects: []string{"Обществознание", "Математика"}, RealWorld: "Бухгалтеры рассчитывают налоги с помощью процентных ставок.", ReadingMinutes: 4},
		{ID: "story-soc-law", Title: "Право: правила общества", Text: "Конституция, законы и права человека.", LinkedModules: []string{"law", "social-9-1"}, Subjects: []string{"Обществознание"}, RealWorld: "Юристы защищают права граждан в судах.", ReadingMinutes: 4},
		{ID: "story-soc-society", Title: "Общество и социальные группы", Text: "Стратификация, роли и социальная мобильность.", LinkedModules: []string{"society", "social-10-1"}, Subjects: []string{"Обществознание"}, RealWorld: "Социологи исследуют поведение групп для бизнеса и власти.", ReadingMinutes: 4},
		{ID: "story-soc-inflation", Title: "Инфляция: почему дорожают товары", Text: "Инфляция — рост общего уровня цен.", LinkedModules: []string{"inflation", "percent"}, Subjects: []string{"Обществознание", "Математика"}, RealWorld: "Центральные банки регулируют инфляцию процентными ставками.", ReadingMinutes: 4},
	}
}

// LaunchProjects returns the M2 launch project idea catalog (≥ 30 ideas).
// Every project requires modules from ≥ 2 subjects and has a graded
// difficulty level (basic / medium / advanced).
func LaunchProjects() []domain.ProjectIdea {
	return []domain.ProjectIdea{
		// Cross math + science
		{ID: "proj-math-bio-lab", Title: "Биохимическая лаборатория дома", Modules: []string{"solutions", "chemistry", "cells"}, DifficultyLevel: domain.DifficultyMedium, ExpectedOutcome: "Провести серию опытов по концентрации растворов и описать результаты."},
		{ID: "proj-math-eco", Title: "Экология двора", Modules: []string{"chemistry", "percent", "ecology"}, DifficultyLevel: domain.DifficultyBasic, ExpectedOutcome: "Рассчитать процент загрязнения и предложить меры."},
		{ID: "proj-math-physics", Title: "Самодельный барометр", Modules: []string{"pressure", "mechanics", "functions"}, DifficultyLevel: domain.DifficultyMedium, ExpectedOutcome: "Собрать барометр и построить график изменения давления за неделю."},
		{ID: "proj-math-data", Title: "Статистика школьной столовой", Modules: []string{"statistics", "data-analysis", "percent"}, DifficultyLevel: domain.DifficultyBasic, ExpectedOutcome: "Собрать данные о питании и представить в виде графиков."},

		// Physics + engineering
		{ID: "proj-phys-bridge", Title: "Мост из спагетти", Modules: []string{"mechanics", "geometry", "linear-equations"}, DifficultyLevel: domain.DifficultyAdvanced, ExpectedOutcome: "Построить мост, выдерживающий нагрузку, и рассчитать силы."},
		{ID: "proj-phys-solar", Title: "Солнечная батарея из подручных средств", Modules: []string{"electricity", "optics", "energy"}, DifficultyLevel: domain.DifficultyAdvanced, ExpectedOutcome: "Собрать фотоэлемент и измерить его мощность."},
		{ID: "proj-phys-audio", Title: "Музыкальный синтезатор на частотах", Modules: []string{"waves", "logarithms", "programming-basics"}, DifficultyLevel: domain.DifficultyAdvanced, ExpectedOutcome: "Создать программу, генерирующую звуки разных частот."},
		{ID: "proj-phys-flight", Title: "Планер из бумаги", Modules: []string{"pressure", "mechanics", "vectors"}, DifficultyLevel: domain.DifficultyBasic, ExpectedOutcome: "Сконструировать планер и объяснить его полёт силами."},

		// Chemistry + biology
		{ID: "proj-chem-water", Title: "Очистка воды в миниатюре", Modules: []string{"solutions", "chemical-reactions", "ecology"}, DifficultyLevel: domain.DifficultyMedium, ExpectedOutcome: "Собрать фильтр и проверить качество воды химическими тестами."},
		{ID: "proj-chem-food", Title: "Химия на кухне: наука о выпечке", Modules: []string{"chemical-reactions", "acids-bases", "fractions"}, DifficultyLevel: domain.DifficultyBasic, ExpectedOutcome: "Изучить, как разрыхлители и кислота влияют на тесто."},
		{ID: "proj-chem-photosynthesis", Title: "Фотосинтез в пробирке", Modules: []string{"photosynthesis", "chemical-reactions", "optics"}, DifficultyLevel: domain.DifficultyAdvanced, ExpectedOutcome: "Провести эксперимент по выделению кислорода водорослями на свету."},
		{ID: "proj-chem-ph", Title: "Карта кислотности района", Modules: []string{"acids-bases", "demography", "data-analysis"}, DifficultyLevel: domain.DifficultyMedium, ExpectedOutcome: "Измерить pH почвы в разных местах и построить карту."},

		// History + literature + art
		{ID: "proj-hist-museum", Title: "Мини-музей эпохи", Modules: []string{"ancient-history", "cultural-history", "russian-classics"}, DifficultyLevel: domain.DifficultyMedium, ExpectedOutcome: "Создать экспозицию о выбранной эпохе с артефактами и текстами."},
		{ID: "proj-hist-family", Title: "Генеалогическое древо семьи", Modules: []string{"middle-ages", "world-wars", "modern-history"}, DifficultyLevel: domain.DifficultyBasic, ExpectedOutcome: "Исследовать историю семьи и связать её с событиями страны."},
		{ID: "proj-lit-play", Title: "Постановка сцены из классики", Modules: []string{"drama", "russian-classics", "poetry"}, DifficultyLevel: domain.DifficultyMedium, ExpectedOutcome: "Подготовить и сыграть сцену с декорациями."},
		{ID: "proj-lit-almanac", Title: "Альманах собственных стихов", Modules: []string{"poetry", "literature"}, DifficultyLevel: domain.DifficultyBasic, ExpectedOutcome: "Написать сборник стихов в разных жанрах."},

		// Geography + social studies
		{ID: "proj-geo-trip", Title: "Виртуальное путешествие", Modules: []string{"continents", "climate", "demography"}, DifficultyLevel: domain.DifficultyBasic, ExpectedOutcome: "Разработать маршрут путешествия с расчётом бюджета."},
		{ID: "proj-geo-climate-report", Title: "Климатический отчёт школы", Modules: []string{"climate", "data-analysis", "statistics"}, DifficultyLevel: domain.DifficultyMedium, ExpectedOutcome: "Собрать метеоданные и построить климатический график."},
		{ID: "proj-geo-economics", Title: "Экономика моего города", Modules: []string{"economics", "taxes", "natural-resources"}, DifficultyLevel: domain.DifficultyMedium, ExpectedOutcome: "Проанализировать бюджет города и источники доходов."},
		{ID: "proj-geo-river", Title: "Исследование реки", Modules: []string{"rivers", "maps", "chemistry"}, DifficultyLevel: domain.DifficultyMedium, ExpectedOutcome: "Изучить экосистему реки и составить карту загрязнений."},

		// Computer science + everything
		{ID: "proj-cs-game", Title: "Обучающая игра в Scratch или Python", Modules: []string{"programming-basics", "python", "algorithms"}, DifficultyLevel: domain.DifficultyMedium, ExpectedOutcome: "Создать игру, обучающую школьной теме."},
		{ID: "proj-cs-bot", Title: "Чат-бот по школьному расписанию", Modules: []string{"python", "data-analysis", "statistics"}, DifficultyLevel: domain.DifficultyAdvanced, ExpectedOutcome: "Разработать бота, отвечающего на вопросы о расписании."},
		{ID: "proj-cs-infographic", Title: "Инфографика научного факта", Modules: []string{"data-analysis", "statistics", "science-fiction"}, DifficultyLevel: domain.DifficultyBasic, ExpectedOutcome: "Визуализировать научный факт в виде плаката."},
		{ID: "proj-cs-security", Title: "Шифр Цезаря своими руками", Modules: []string{"cybersecurity", "algorithms", "poetry"}, DifficultyLevel: domain.DifficultyBasic, ExpectedOutcome: "Написать программу шифрования и зашифровать послание."},

		// Interdisciplinary flagship projects
		{ID: "proj-flagship-space", Title: "Космическая миссия: от физики до истории", Modules: []string{"mechanics", "optics", "modern-history", "science-fiction"}, DifficultyLevel: domain.DifficultyAdvanced, ExpectedOutcome: "Спроектировать миссию с расчётом орбиты и историческим контекстом."},
		{ID: "proj-flagship-city", Title: "Умный город: проект будущего", Modules: []string{"electricity", "data-analysis", "economics", "maps"}, DifficultyLevel: domain.DifficultyAdvanced, ExpectedOutcome: "Разработать концепцию умного города с энергосбережением."},
		{ID: "proj-flagship-ocean", Title: "Океан: экосистема и человек", Modules: []string{"ecology", "rivers", "natural-resources", "inflation"}, DifficultyLevel: domain.DifficultyMedium, ExpectedOutcome: "Исследовать влияние человека на океан и предложить меры."},
		{ID: "proj-flagship-health", Title: "Здоровье: биология, химия и данные", Modules: []string{"human-anatomy", "chemical-reactions", "statistics", "data-analysis"}, DifficultyLevel: domain.DifficultyAdvanced, ExpectedOutcome: "Провести исследование здоровья школьников и представить анализ."},
		{ID: "proj-flagship-history", Title: "Хронология технологий", Modules: []string{"world-wars", "modern-history", "algorithms", "electricity"}, DifficultyLevel: domain.DifficultyMedium, ExpectedOutcome: "Составить таймлайн изобретений и их влияния на общество."},
		{ID: "proj-flagship-budget", Title: "Семейный бюджет: математика и обществознание", Modules: []string{"percent", "taxes", "inflation", "economics"}, DifficultyLevel: domain.DifficultyBasic, ExpectedOutcome: "Составить и проанализировать семейный бюджет с учётом инфляции."},
		{ID: "proj-flagship-poetry-science", Title: "Поэзия и наука: метафоры явлений", Modules: []string{"poetry", "optics", "waves", "russian-classics"}, DifficultyLevel: domain.DifficultyMedium, ExpectedOutcome: "Написать стихи о научных явлениях и объяснить их термины."},
	}
}

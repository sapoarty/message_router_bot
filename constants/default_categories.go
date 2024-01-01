package constants

func GetDefaultCategories() map[string][]string {
    return map[string][]string{
        "Почитать": {
            "книга", "роман", "беллетристика", "чтиво", "рекомендую почитать", 
            "обложка", "писатель", "литература", "издание", "поэзия", "стихи", 
            "научная работа", "журнал", "газета", "публицистика", "нерецензируемая литература",
            "автобиография", "энциклопедия", "новинка в мире литературы", "читальный зал",
        },
        "Посмотреть кино": {
            "боевик", "драма", "комедия", "триллер", "блокбастер", "кинофильм", "режиссер", 
            "сценарист", "актер", "актриса", "голливуд", "кинопремьера", "фильмография",
            "анимация", "документальный фильм", "фестиваль кино", "критика кино", "субтитры",
            "кинотеатр", "ночь кино", "оскар", "золотая пальмовая ветвь",
        },
        "Музыка": {
            "альбом", "сингл", "композиция", "исполнитель", "группа", "джаз", 
            "рок", "поп-музыка", "классика", "хип-хоп", "электронная музыка", "саундтрек",
            "музыкальный клип", "концерт", "живое выступление", "гитара", "пианино",
            "музыкальный фестиваль", "ноты", "аранжировка", "вокал", "хор",
        },
        "Наука": {
            "исследование", "открытие", "ученый", "лаборатория", "эксперимент", 
            "публикация", "научная работа", "конференция", "семинар", "лекция", "доклад",
            "теория", "гипотеза", "научный метод", "статистика", "опыт",
            "анализ данных", "пробирка", "микроскоп", "SCI-журнал", "академия наук",
        },
        "Посмотреть youtube": {
            "youtube.", "vlog", "подкаст", "обучающий ролик", "реакция", "ютуб-канал",
            "стрим", "прямая трансляция", "блогер", "вайрал", "распаковка",
            "интервью", "мастер-класс", "скетчи", "образовательное видео", "дайджест",
            "геймплей на ютуб", "рукоделие", "кулинария", "юмористическое шоу", "научно-популярное",
        },
        "Кулинария": {
            "рецепт", "блюдо", "приготовление", "кухня", "ингредиент",
            "повар", "ресторан", "кулинарное шоу", "выпекать", "жарить",
            "варить", "резать", "запекать", "гарнир", "десерт",
            "салат", "гастроном", "еда на вынос", "дегустация", "пищевой блогер",
        },
        "Программирование": {
            "программирование", "кодинг", "разработка", "code", "dev", "go", "python", "javascript",
            "frontend", "backend", "full-stack", "бэкенд", "фронтенд", "веб-разработка",
            "алгоритмы", "структуры данных", "интерфейсы", "функции", "объектно-ориентированное программирование",
            "программист", "разработчик", "код", "переменные", "циклы", "условия", "массивы", "списки", "карты", "горутины",
            "каналы", "слайсы", "методы", "пакеты", "инкапсуляция", "наследование", "полиморфизм", "абстракция",
            "инкапсуляция", "рекурсия", "стек", "очередь", "сортировка", "список", "бинарное дерево", "хэш-таблица",
            "граф", "Дейкстры", "алгоритм", "сортировки", "булевый тип", "целочисленный тип", "строковый тип", "байтовый тип",
            "указатели", "массив структур", "управляющие конструкции", "циклическая структура", "условный оператор",
            "многопоточное программирование", "параллельное программирование", "асинхронное программирование",
            "функциональное программирование", "модули", "библиотеки", "области видимости", "контекст", "замыкания",
            "аргументы функций", "рекурсивные функции", "стек вызовов", "область памяти", "управление памятью",
            "типы данных", "регистры процессора", "отладка", "тестирование", "рефакторинг", "habr.", "документация", "Python", "Java", "JavaScript", "C#", "C++", "Ruby", "Go", "PHP",
            "swift", "kotlin", "typescript", "scala", "perl", "rust", "objective-c", "visual basic", "dart", "r", "sql", "shell",
            "groovy", "powershell", "haskell", "lua", "f#", "elixir",
            "фреймворк", "angular", "react", "vue.js", "django", "flask", "ruby on rails", "spring", "asp.net",
            "laravel", "express", "node.js", "bootstrap", "jquery",
            "среда разработки", "ide", "visual studio code", "intellij idea", "eclipse", "pycharm", "sublime text", "atom",
            "netbeans", "xcode", "android studio", "github", "git", "svn", "mercurial",
            "система сборки", "maven", "gradle", "ant", "webpack", "gulp", "grunt",
            "базы данных", "mysql", "postgresql", "sqlite", "mongodb", "cassandra", "redis", "oracle database", "microsoft sql server",
            "devops", "docker", "kubernetes", "ansible", "terraform", "jenkins", "ci/cd",
            "тестирование", "junit", "selenium", "testng", "cypress", "mocha", "jest",
            "облачные платформы", "aws", "azure", "google cloud platform", "heroku", "машинное обучение", "tensorflow", "keras", "pytorch", 
            "scikit-learn", "opencv", "библиотеки", "numpy", "pandas", "matplotlib", "d3.js", "дизайн", "figma", "sketch", "adobe xd", 
            "invision", "zeplin", "аппаратное обеспечение", "arduino", "raspberry pi", "api", "rest", "graphql", "soap", 
            "версионирование", "git", "subversion", "mercurial",
        },
        "Финансы": {
            "финансы", "инвестиции", "биржа", "акции", "облигации", "фондовый рынок", "трейдинг", "дивиденды", 
            "портфель инвестиций", "доходность", "риски", "доход", "капитал", "долг", "пассивный доход", "арбитраж", "ликвидность", 
            "рыночная стоимость", "рыночная капитализация", "денежный поток", "страхование", "пенсионное планирование", "налоги", 
            "кредитование", "стоимость капитала",
        },
        "Купить": {
            "продукты", "одежда", "электроника", "игрушки", "скидка", "скидки", "маркет", "ali.", "aliexpress.", "отзыв", 
            "акция", "распродажа", "аутлет", "бренд", "интернет-магазин", "доставка", "гарантия", "новинка", "эксклюзив", 
            "ритейлер", "оптовая закупка", "торговый центр", "распаковка", "бытовая техника", "косметика", "аксессуары", 
            "украшения", "мебель", "автомобиль", "заказ онлайн", "курьерская служба", "обзор продукции", "вернуть товар", 
            "потребительский кредит", "характеристики товара", "сертификат качества", "сравнение цен", "рейтинг продавца",
        },
        "Путешествия": {
            "ехать", "поход", "туризм", "билет", "самолет", "поезд",
            "гостиница", "отель", "апартаменты", "чартерный рейс", "бронирование", "страховка для путешественников", 
            "экскурсия", "гид", "путеводитель", "туристическая группа", "направление", "курорт", "виза", "пакетный тур", 
            "экзотические страны", "местный колорит", "памятники культуры", "исторические места", "природные заповедники", 
            "туристический сезон", "недорогой отдых", "высококлассный отдых", "круиз", "культурное наследие", 
            "медицинский туризм", "палаточный лагерь", "снаряжение для путешествий", "карта путешественника", 
            "авиакомпания", "aэропортовый сбор", "трансфер", "aviasales",
        },

    }
}

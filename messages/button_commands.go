package messages

var CommandStart = map[string]string{
    "en": "/start",
    "ru": "/start",
}

var CommandHelp = map[string]string{
    "en": "/help",
    "ru": "/help",
}

var CommandAddByCategory = map[string]string{
    "en": "Set the keywords using the category",
    "ru": "Установить ключевые слова для группы, используя категорию",
}

var CommandAddKeywords = map[string]string{
    "en": "Add keywords/keyphrases",
    "ru": "Добавить ключевые слова/фразы",
}
var CommandDeleteKeywords = map[string]string{
    "en": "Delete keywords/keyphrases",
    "ru": "Удалить ключевые слова/фразы",
}
var CommandSetDefaultGroup = map[string]string{
    "en": "Set group as Default",
    "ru": "Установить как группу для несортированных",
}

var CommandPrintAllKeywords = map[string]string{
    "en": "Print all set keywords and chats",
    "ru": "Вывести все ключевые слова по группам",
}

var KeywordsAreNotSetYet = map[string]string{
    "en": "Keywords are not set yet, use /help for information",
    "ru": "Ключевые слова не заданы, /help для информации",
}

var CommandChangeLang = map[string]string{
    "en": "Eng/Rus",
    "ru": "Eng/Rus",
}
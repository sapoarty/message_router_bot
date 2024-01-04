package messages

var GreetingsMessage = map[string]string{
	"en": "Welcome! This bot is here to help you sort through the many messages you have in your 'Saved Messages' group, " +
	"or forward those solitary messages directly to the bot, thus saving you time. " + 
	"For optimal performance, this bot needs information about the keywords or keyphrases for each group to which you wish to forward messages. " +
	"To allow the Bot use all the functions should set the Bot as admin of needed groups." +
	"If the keyword or keyphrase is not found, then the message will be forwarded to the group you've set as the 'Default' - this is a mandatory step.",
	"ru": "Добро пожаловать! Этот бот поможет вам сортировать множество сообщений в вашей группе 'Сохраненные сообщения', " +
	"или пересылать отдельные сообщения прямо боту, что позволит сэкономить время.\n" +
	"Для оптимальной работы боту необходима информация о ключевых словах или фразах для каждой группы, в которую вы хотите переслать сообщения.\n" +
	"Для полной функциональности бота следует установить администратором групп.\n" +
	"Если ключевое слово или фраза не найдены, то сообщение будет переслано в группу, которую вы установили как 'По-умолчанию' - это обязательный шаг.\n\n ",
}

var DefaulGroupAlias = map[string]string{
	"en": "Default group (for unsorted messages)",
	"ru": "Группа по-умолчанию (для несортированных сообщений)",
}

var DefaulGroupSetSuccessfully = map[string]string{
	"en": "Group '%s' is set as default group (for unsorted messages)",
	"ru": "Группа '%s' установлена, как группа по-умолчанию (для несортированных сообщений)",
}

var DefaulGroupResetSuccessfully = map[string]string{
	"en": "Old default group '%s' is reset",
	"ru": "Старая группа '%s' по-умолчанию сброшена",
}

var InputKeywordsListRequest = map[string]string{
	"en": "Please, send keyword and phrases list, separated by comma",
	"ru": "Пожалуйста, отправьте список ключевых слов и фраз, разделенных запятой",
}

var KeywordsListEmpty = map[string]string{
	"en": "Keywords list is empty, please add one or the list using ',' separator",
	"ru": "Список ключевых слов пуст, пожалуйста, добавьте одно или несколько слов с использованием разделителя ','",
}

var KeywordsListForGroup = map[string]string{
	"en": "Keyword and keyphrases list for group '%s': %s" ,
	"ru": "Ключевые слова и фразы для группы '%s': %s",
}

var AwaitingCommand = map[string]string{
	"en": "Please select a command or send messages to the [Bot](%s) for sorting",
	"ru": "Пожалуйста, выберите команду или отправьте сообщения [Боту](%s) для сортировки",
}

var DefaulGroupIsNotSet = map[string]string{
	"en": "Default group (for unsorted messages) is not set" ,
	"ru": "Группа по-умолчанию (для несортированных сообщений) не установлена",
}

var UseMenuCommands = map[string]string{
	"en": "Please, use menu",
	"ru": "Пожалуйста, используйте команды из меню",
}

var LangIsChanged = map[string]string{
	"en": "Bot language has been changed to Eng",
	"ru": "Язык меню был изменен на Рус",
}

var KeywordsNotFound = map[string]string{
	"en": "Keywords not found: %s",
	"ru": "Ключевые слова не найдены: %s",
}

var AllKeywordsFoundAndDeleted = map[string]string{
	"en": "All requested keywords found and deleted",
	"ru": "Все запрошенные ключевые слова найдены и удалены",
}

var KeywordsListIsChanged = map[string]string{
	"en": "List has been updated. Keywords list [%s] can be found in group description",
	"ru": "Список обновлен. Актуальный список ключевых слов и фраз [%s] вы всегда можете найти в описании группы",
}

var GroupDesc = map[string]string{
	"en": "Keywords list for this group: %s",
	"ru": "Список ключевых слов/фраз для этой группы: %s",
}

var DefaultGroupDesc = map[string]string{
	"en": "Default group (for unsorted messages)",
	"ru": "Группа по-умолчанию (для несортированных сообщений)",
}

var GroupHasKeywordsAndCanNotBeDefault = map[string]string{
	"en": "This group can not be used as unsorted messages group as it's already has some set keywords [%s]",
	"ru": "Эта группа не может быть использована, как группа для несортированных сообщений, так как уже имеет установленные ключевые слова [%s]",
}

var GroupIsAlreadyDefaultGroup = map[string]string{
	"en": "This group is already set as default group (for unsorted messages).",
	"ru": "Эта группа уже установлена, как группа по-умолчанию (для несортированных сообщений).",
}

var SetAnotherGroupAsDefault = map[string]string{
	"en": "Please, set another group as default to use commands",
	"ru": "Пожалуйста, установите другую группу, как группу по-умолчанию для использования команд",
}

var ChooseCategory = map[string]string{
	"en": "Choose the category below:",
	"ru": "Выберите категорию из списка:",
}

var ChatDeleted = map[string]string{
	"en": "The group chat '%s' has been deleted. If no other chats with the keyword '%s' are found, the message will be forwarded to the default group.",
	"ru": "Групповой чат '%s' был удален. Если другие группы с ключевым словом '%s' не будут найдены, сообщение будет перенаправлено в группу по умолчанию.",
}

var  MessageHasBeenForwardedToChatsUsingKeywords = map[string]string{
	"en": "Message with text \n\n [%s] \n\n has been forwarded to the chat/s using (keywords): %s",
	"ru": "Сообщение с текстом \n\n [%s] \n\n было переотправлено в чат/ы, используя (ключевые слова): %s",
}

var  MessageHasBeenForwardedToDefaultGroup = map[string]string{
	"en": "Message with text \n\n [%s] \n\n has been forwarded to the Default group: %s",
	"ru": "Сообщение с текстом \n\n [%s] \n\n было переотправлено в группу по-умолчанию: %s",
}


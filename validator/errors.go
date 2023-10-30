package validator

var errorMessages = map[string]string{
	"required":      "Поле %s обязательно",
	"email_custom":  "Email %s не действителен",
	"str_gt":        "Поле %s должно иметь больше %s символов",
	"str_lt":        "Поле %s должно иметь меньше %s символов",
	"has_lowercase": "Поле %s должно иметь хотя бы один маленькую букву",
	"has_uppercase": "Поле %s должно быть хотя бы один большую букву",
	"has_special":   "Поле %s должно содержать хотя бы один специальный символ",
}

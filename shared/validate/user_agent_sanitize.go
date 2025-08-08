package shared_validate

import "strings"

func UserAgentSanitize(ua string) string {
	// user_agent будет использоваться для антифлуда, то есть нас не интересует абсолютное значение, достаточно относительное. Поэтому если вдруг user_agent > len, то просто обрезаем его вместо того, чтоб вернуть ошибку
	// также мы подсчитываем ascii символы, считать unicode это медленнее, и не имеет смысла, user_agent в ascii

	// Убираем не-ASCII или управляющие символы
	clean := strings.Map(func(r rune) rune {
		if r < 32 || r > 126 {
			return -1
		}
		return r
	}, ua)

	// Обрезаем до 512 символов
	if len(clean) > 512 {
		clean = clean[:512]
	}

	return clean
}

package utils

import (
	"strings"
	"unicode"
)

// GenerateSlug gera um slug a partir de um nome
// Remove acentos básicos, converte para minúsculas e substitui espaços por hífens
func GenerateSlug(name string) string {
	// Mapa de caracteres acentuados para não acentuados
	accentMap := map[rune]rune{
		'á': 'a', 'à': 'a', 'ã': 'a', 'â': 'a', 'ä': 'a',
		'é': 'e', 'è': 'e', 'ê': 'e', 'ë': 'e',
		'í': 'i', 'ì': 'i', 'î': 'i', 'ï': 'i',
		'ó': 'o', 'ò': 'o', 'õ': 'o', 'ô': 'o', 'ö': 'o',
		'ú': 'u', 'ù': 'u', 'û': 'u', 'ü': 'u',
		'ç': 'c', 'ñ': 'n',
		'Á': 'A', 'À': 'A', 'Ã': 'A', 'Â': 'A', 'Ä': 'A',
		'É': 'E', 'È': 'E', 'Ê': 'E', 'Ë': 'E',
		'Í': 'I', 'Ì': 'I', 'Î': 'I', 'Ï': 'I',
		'Ó': 'O', 'Ò': 'O', 'Õ': 'O', 'Ô': 'O', 'Ö': 'O',
		'Ú': 'U', 'Ù': 'U', 'Û': 'U', 'Ü': 'U',
		'Ç': 'C', 'Ñ': 'N',
	}

	// Converter e normalizar
	var builder strings.Builder
	for _, r := range name {
		// Substituir acentos
		if replacement, ok := accentMap[r]; ok {
			r = replacement
		}

		// Converter para minúsculas
		r = unicode.ToLower(r)

		// Manter apenas letras, dígitos e hífens
		if unicode.IsLetter(r) || unicode.IsDigit(r) {
			builder.WriteRune(r)
		} else if r == ' ' || r == '-' || r == '_' {
			// Adicionar hífen apenas se o último caractere não for hífen
			if builder.Len() > 0 {
				lastChar := builder.String()[builder.Len()-1]
				if lastChar != '-' {
					builder.WriteRune('-')
				}
			}
		}
	}

	slug := builder.String()

	// Remover hífens do início e fim
	slug = strings.Trim(slug, "-")

	// Remover hífens duplicados
	for strings.Contains(slug, "--") {
		slug = strings.ReplaceAll(slug, "--", "-")
	}

	return slug
}


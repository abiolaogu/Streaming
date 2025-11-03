package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"
)

// SupportedLocales maps locale codes to language names
var SupportedLocales = map[string]string{
	"en-US": "English (US)",
	"es-ES": "Spanish (Spain)",
	"fr-FR": "French (France)",
	"de-DE": "German (Germany)",
	"pt-BR": "Portuguese (Brazil)",
	"ja-JP": "Japanese (Japan)",
	"zh-CN": "Chinese (Simplified)",
}

const LocaleContextKey = "locale"
const DefaultLocale = "en-US"

// I18nMiddleware extracts locale from Accept-Language header
// Issue #30: i18n (Internationalization) Support
func I18nMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		locale := DefaultLocale

		// Parse Accept-Language header
		acceptLang := c.GetHeader("Accept-Language")
		if acceptLang != "" {
			// Parse: "en-US,en;q=0.9,es;q=0.8"
			parts := strings.Split(acceptLang, ",")
			if len(parts) > 0 {
				firstLang := strings.TrimSpace(strings.Split(parts[0], ";")[0])
				
				// Check if exact match exists
				if _, exists := SupportedLocales[firstLang]; exists {
					locale = firstLang
				} else {
					// Try language code only (e.g., "en" -> "en-US")
					langCode := strings.Split(firstLang, "-")[0]
					for key := range SupportedLocales {
						if strings.HasPrefix(key, langCode+"-") {
							locale = key
							break
						}
					}
				}
			}
		}

		// Set locale in context
		c.Set(LocaleContextKey, locale)
		c.Next()
	}
}

// GetLocale gets the current locale from context
func GetLocale(c *gin.Context) string {
	if locale, exists := c.Get(LocaleContextKey); exists {
		if loc, ok := locale.(string); ok {
			return loc
		}
	}
	return DefaultLocale
}


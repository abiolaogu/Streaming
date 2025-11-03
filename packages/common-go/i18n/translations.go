package i18n

import (
	"fmt"
	"time"
)

// Translations map for different languages
var translations = map[string]map[string]string{
	"en-US": {
		"subscription.active":     "Your %s subscription is active",
		"subscription.expired":    "Your subscription has expired",
		"content.not_found":       "Content not found",
		"auth.invalid_credentials": "Invalid email or password",
		"payment.success":         "Payment processed successfully",
		"payment.failed":          "Payment failed. Please try again",
	},
	"es-ES": {
		"subscription.active":     "Tu suscripción %s está activa",
		"subscription.expired":    "Tu suscripción ha expirado",
		"content.not_found":       "Contenido no encontrado",
		"auth.invalid_credentials": "Correo electrónico o contraseña inválidos",
		"payment.success":         "Pago procesado exitosamente",
		"payment.failed":          "Pago fallido. Por favor, inténtalo de nuevo",
	},
	"fr-FR": {
		"subscription.active":     "Votre abonnement %s est actif",
		"subscription.expired":    "Votre abonnement a expiré",
		"content.not_found":       "Contenu introuvable",
		"auth.invalid_credentials": "Email ou mot de passe invalide",
		"payment.success":         "Paiement traité avec succès",
		"payment.failed":          "Échec du paiement. Veuillez réessayer",
	},
	"de-DE": {
		"subscription.active":     "Ihr %s Abonnement ist aktiv",
		"subscription.expired":    "Ihr Abonnement ist abgelaufen",
		"content.not_found":       "Inhalt nicht gefunden",
		"auth.invalid_credentials": "Ungültige E-Mail oder Passwort",
		"payment.success":         "Zahlung erfolgreich verarbeitet",
		"payment.failed":          "Zahlung fehlgeschlagen. Bitte versuchen Sie es erneut",
	},
	"pt-BR": {
		"subscription.active":     "Sua assinatura %s está ativa",
		"subscription.expired":    "Sua assinatura expirou",
		"content.not_found":       "Conteúdo não encontrado",
		"auth.invalid_credentials": "Email ou senha inválidos",
		"payment.success":         "Pagamento processado com sucesso",
		"payment.failed":          "Pagamento falhou. Por favor, tente novamente",
	},
	"ja-JP": {
		"subscription.active":     "あなたの%sサブスクリプションは有効です",
		"subscription.expired":    "サブスクリプションが期限切れです",
		"content.not_found":       "コンテンツが見つかりません",
		"auth.invalid_credentials": "メールアドレスまたはパスワードが無効です",
		"payment.success":         "支払いが正常に処理されました",
		"payment.failed":          "支払いに失敗しました。もう一度お試しください",
	},
	"zh-CN": {
		"subscription.active":     "您的%s订阅已激活",
		"subscription.expired":    "您的订阅已过期",
		"content.not_found":       "未找到内容",
		"auth.invalid_credentials": "电子邮件或密码无效",
		"payment.success":         "支付处理成功",
		"payment.failed":          "支付失败。请重试",
	},
}

// LocalizerFunc localizes a message key with arguments
type LocalizerFunc func(key string, args ...interface{}) string

// GetLocalizer returns a localizer function for the given locale
func GetLocalizer(locale string) LocalizerFunc {
	lang, exists := translations[locale]
	if !exists {
		lang = translations["en-US"] // Fallback to English
	}

	return func(key string, args ...interface{}) string {
		message, exists := lang[key]
		if !exists {
			return key // Return key if translation not found
		}

		if len(args) > 0 {
			return fmt.Sprintf(message, args...)
		}
		return message
	}
}

// FormatCurrency formats currency based on locale
func FormatCurrency(locale string, amount float64, currency string) string {
	switch locale {
	case "en-US", "pt-BR":
		return fmt.Sprintf("%s%.2f", currency, amount)
	case "es-ES", "fr-FR", "de-DE":
		return fmt.Sprintf("%.2f %s", amount, currency)
	case "ja-JP", "zh-CN":
		return fmt.Sprintf("%s%.0f", currency, amount)
	default:
		return fmt.Sprintf("%s%.2f", currency, amount)
	}
}

// FormatDate formats date based on locale
func FormatDate(locale string, t time.Time) string {
	switch locale {
	case "en-US":
		return t.Format("01/02/2006")
	case "es-ES", "pt-BR":
		return t.Format("02/01/2006")
	case "fr-FR", "de-DE":
		return t.Format("02.01.2006")
	case "ja-JP", "zh-CN":
		return t.Format("2006年01月02日")
	default:
		return t.Format("2006-01-02")
	}
}


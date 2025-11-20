package middleware

import (
	"fmt"
	"net/http"
	"runtime"

	"github.com/WahyuSiddarta/be_saham_go/helper"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog"
)

var defaultRecoverConfig = RecoverConfig{
	StackTraceSize:                 4 << 10, // 4 KB
	PrintStackTraceOfAllGoroutines: false,
}

// RecoverConfig defines the config for Recover middleware.
type RecoverConfig struct {
	// Size allocated on memory for stack trace.
	StackTraceSize int

	// If stack trace is enabled, this is to print stack traces of all goroutines.
	PrintStackTraceOfAllGoroutines bool

	// The panic was happened, and it was handled and logged gracefully.
	// What's next?
	// This function is called to handle the error of panic.
	ErrorHandler func(c echo.Context, err error)
}

func RecoverWithConfig(config RecoverConfig, logger *zerolog.Logger) echo.MiddlewareFunc {
	if config.StackTraceSize == 0 {
		config.StackTraceSize = defaultRecoverConfig.StackTraceSize
	}

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			defer func() {
				if err := recover(); err != nil {
					e := func() error {
						if e, ok := err.(error); ok {
							return e
						} else {
							return fmt.Errorf("panic: %v", err)
						}
					}()

					// Set error on context
					c.Error(e)

					// Get stack trace
					stack := make([]byte, config.StackTraceSize)
					stackLen := runtime.Stack(stack, config.PrintStackTraceOfAllGoroutines)
					errorMessage := string(stack[:stackLen])

					// Log the error immediately
					logger.Error().Msgf("PANIC RECOVERED: %s\nStack trace:\n%s", e.Error(), errorMessage)

					// Capture panic in Sentry
					CaptureRecovery(c, e, errorMessage)

					// TODO: Telegram notification (commented out for now)
					// Uncomment and configure the following block to enable telegram notifications:
					/*
						go func() {
							defer func() {
								if teleErr := recover(); teleErr != nil {
									logger.Error().Msgf("Error in telegram notification: %v", teleErr)
								}
							}()

							client := http.Client{
								Timeout: 15 * time.Second,
							}

							// Configure these constants in your environment:
							// TELEGRAM_BOT_ID, TELEGRAM_CHAT_ID, TELEGRAM_URL
							telegramMsg := "be_saham_go panic: " + sanitizeTelegramMessage(errorMessage)
							URI := "https://api.telegram.org" + "/bot" + "YOUR_BOT_TOKEN" + "/sendMessage?chat_id=" + "YOUR_CHAT_ID" + "&parse_mode=Markdown&text=" + url.QueryEscape(telegramMsg)

							resp, err := client.Get(URI)
							if err != nil {
								logger.Error().Msgf("RECOVERY telegram error %s", err.Error())
							} else {
								logger.Info().Msgf("RECOVERY telegram success %s", resp.Status)
								if resp.Body != nil {
									resp.Body.Close()
								}
							}
						}()
					*/

					// Ensure response is not already written
					if !c.Response().Committed {
						// Try to return a proper HTTP error response
						if err := helper.ErrorResponse(c, http.StatusInternalServerError, "Terjadi kesalahan internal server", nil); err != nil {
							logger.Error().Msgf("Failed to send error response: %v", err)
						}
					}

					if config.ErrorHandler != nil {
						config.ErrorHandler(c, e)
					}
				}
			}()
			return next(c)
		}
	}
}

// Recover returns a middleware that recovers from panics anywhere in the chain
func Recover() echo.MiddlewareFunc {
	return RecoverWithConfig(defaultRecoverConfig, Logger)
}

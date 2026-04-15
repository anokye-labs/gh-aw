//go:build !js && !wasm

package console

import (
	"errors"

	"charm.land/huh/v2"
	"github.com/github/gh-aw/pkg/logger"
	"github.com/github/gh-aw/pkg/styles"
	"github.com/github/gh-aw/pkg/tty"
)

var inputLog = logger.New("console:input")

// PromptSecretInput shows an interactive password input prompt with masking
// The input is masked for security and includes validation
// Returns the entered secret value or an error
func PromptSecretInput(title, description string) (string, error) {
	inputLog.Printf("PromptSecretInput: prompting for %q", title)

	// Check if stdin is a TTY - if not, we can't show interactive forms
	if !tty.IsStderrTerminal() {
		inputLog.Print("PromptSecretInput: TTY not available, cannot show interactive form")
		return "", errors.New("interactive input not available (not a TTY)")
	}

	var value string

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title(title).
				Description(description).
				EchoMode(huh.EchoModePassword). // Masks input for security
				Validate(func(s string) error {
					if len(s) == 0 {
						return errors.New("value cannot be empty")
					}
					return nil
				}).
				Value(&value),
		),
	).WithTheme(styles.HuhTheme).WithAccessible(IsAccessibleMode())

	inputLog.Print("PromptSecretInput: running interactive form")
	if err := form.Run(); err != nil {
		inputLog.Printf("PromptSecretInput: form failed: %v", err)
		return "", err
	}

	inputLog.Print("PromptSecretInput: secret value received")
	return value, nil
}

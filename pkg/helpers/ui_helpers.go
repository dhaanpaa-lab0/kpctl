package helpers

import "github.com/pterm/pterm"

func GetInputAsString(prompt string, defaultValue string) string {
	inputBox := pterm.DefaultInteractiveTextInput.WithDefaultValue(defaultValue)
	result, _ := inputBox.Show(prompt)
	return result
}

func SelectOptionAsString(prompt string, options []string) string {
	selectBox := pterm.DefaultInteractiveSelect.WithOptions(options)
	result, _ := selectBox.Show(prompt)
	return result
}

func RenderTable(data pterm.TableData) error {
	return pterm.DefaultTable.WithHasHeader().WithData(data).Render()
}

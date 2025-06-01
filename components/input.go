package components

// Input is a component that handles text input.

type Input struct {
	// Text is the current text in the input field.
	Text string
	// Cursor is the position of the cursor in the input field.
	Cursor int
	// Width is the width of the input field.
	Width int
	// Height is the height of the input field.
	Height int
	// Focused indicates whether the input field is focused.
	Focused bool
	// Placeholder is the placeholder text for the input field.
	Placeholder string
	// OnSubmit is a function that is called when the input is submitted.
	OnSubmit func(text string)
	// OnChange is a function that is called when the input text changes.
	OnChange func(text string)
}

func NewInput(width, height int) *Input {
	return &Input{
		Text:        "",
		Cursor:      0,
		Width:       width,
		Height:      height,
		Focused:     false,
		Placeholder: "",
		OnSubmit:    nil,
		OnChange:    nil,
	}
}

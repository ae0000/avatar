package avatar

import (
	"image/color"
	"testing"
)

func TestParseHexColorFast(t *testing.T) {
	tests := []struct {
		input    string
		expected color.RGBA
		wantErr  bool
	}{
		{"#FF0000", color.RGBA{255, 0, 0, 255}, false},
		{"#00FF00", color.RGBA{0, 255, 0, 255}, false},
		{"#0000FF", color.RGBA{0, 0, 255, 255}, false},
		{"#FFFFFF", color.RGBA{255, 255, 255, 255}, false},
		{"#000000", color.RGBA{0, 0, 0, 255}, false},
		{"#ABC", color.RGBA{170, 187, 204, 255}, false}, // Short form
		{"#F0F", color.RGBA{255, 0, 255, 255}, false},   // Short form
		{"FF0000", color.RGBA{}, true},                   // Missing #
		{"#GG0000", color.RGBA{}, true},                  // Invalid hex
		{"#FF00", color.RGBA{}, true},                    // Wrong length
	}

	for _, tt := range tests {
		result, err := parseHexColorFast(tt.input)
		if tt.wantErr {
			if err == nil {
				t.Errorf("parseHexColorFast(%s) expected error but got none", tt.input)
			}
		} else {
			if err != nil {
				t.Errorf("parseHexColorFast(%s) unexpected error: %v", tt.input, err)
				continue
			}
			if result != tt.expected {
				t.Errorf("parseHexColorFast(%s) = %v, want %v", tt.input, result, tt.expected)
			}
		}
	}
}

func TestDefaultColor(t *testing.T) {
	tests := []struct {
		initial string
		color   string
	}{
		{"A", "Red"},
		{"B", "Green"},
		{"C", "Yellow"},
		{"D", "Blue"},
		{"E", "Orange"},
		{"F", "Purple"},
		{"G", "Lime"},
		{"H", "Magenta"},
		{"I", "Pink"},
		{"J", "Cyan"},
		{"K", "Teal"},
		{"L", "Lavender"},
		{"M", "Brown"},
		{"N", "Beige"},
		{"O", "Maroon"},
		{"P", "Mint"},
		{"Q", "Olive"},
		{"R", "Coral"},
		{"S", "Navy"},
		{"T", "Gold"},
		{"U", "Red"},
		{"V", "Green"},
		{"W", "Blue"},
		{"X", "Orange"},
		{"Y", "Purple"},
		{"Z", "Cyan"},
		{"0", "Red"},
		{"1", "Green"},
		{"9", "Cyan"},
		{"#", "Grey"}, // Unknown character
	}

	for _, tt := range tests {
		result := defaultColor(tt.initial)
		// Just check it doesn't panic and returns a non-nil color
		r, g, b, a := result.RGBA()
		if r == 0 && g == 0 && b == 0 && a == 0 {
			t.Errorf("defaultColor(%s) returned zero value color", tt.initial)
		}
	}
}

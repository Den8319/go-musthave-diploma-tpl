package luhn

import "testing"

func TestValid(t *testing.T) {

	tests := []struct {
		name   string
		number string
		want   bool
	}{
		{name: "Valid 4 digits", number: "1230", want: true},
		{name: "Invalid 4 digits", number: "1235", want: false},
		{name: "Valid 16 digits", number: "49927398716", want: true},

		{name: "Valid 11 digits", number: "79927398713", want: true},
		{name: "Invalid 11 digits", number: "79927398714", want: false},

		{name: "Zero", number: "0", want: false},
		{name: "Negative number", number: "-1234", want: false},
		{name: "Single digit valid", number: "0", want: false}, // по стандарту одиночные цифры обычно невалидны
	}

	// Запуск тестов в цикле
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Valid(tt.number)
			if got != tt.want {
				t.Errorf("Valid(%s) = %v; want %v", tt.number, got, tt.want)
			}
		})
	}
}
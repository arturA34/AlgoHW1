package main

import (
	"testing"
)

func TestCountScore(t *testing.T) {
	tests := []struct {
		name     string
		k        int
		chars    map[int]int
		expected int
	}{
		{
			name: "Пример 1 из условия",
			k:    3,
			chars: map[int]int{
				1: 2,
				2: 8,
				3: 1,
			},
			expected: 2, // цифра 1 (2 <= 6) и 3 (1 <= 6), цифра 2 (8 > 6)
		},
		{
			name: "Пример 2 из условия",
			k:    4,
			chars: map[int]int{
				1: 10,
				9: 6,
			},
			expected: 1, // цифра 9 (6 <= 8), цифра 1 (10 > 8)
		},
		{
			name:     "Граничный: Пустое поле (только точки)",
			k:        1,
			chars:    map[int]int{},
			expected: 0,
		},
		{
			name: "Граничный: k = 1 (минимум), точная граница 2k = 2",
			k:    1,
			chars: map[int]int{
				1: 2, // ровно 2 <= 2*1 -> +1 балл
				2: 3, // 3 > 2*1 -> не успевают
			},
			expected: 1,
		},
		{
			name: "Граничный: k = 5 (максимум), все цифры от 1 до 9 встречаются по 1 разу",
			k:    5,
			chars: map[int]int{
				1: 1, 2: 1, 3: 1, 4: 1, 5: 1, 6: 1, 7: 1, 8: 1, 9: 1,
			},
			expected: 9,
		},
		{
			name: "Граничный: Все 16 клавиш заняты одной цифрой (k=5, 2k=10 < 16)",
			k:    5,
			chars: map[int]int{
				7: 16,
			},
			expected: 0,
		},
		{
			name: "Граничный: Ровно 2k клавиш для 4 разных цифр (k=2, 2k=4)",
			k:    2,
			chars: map[int]int{
				1: 4, 2: 4, 3: 4, 4: 4,
			},
			expected: 4,
		},
		{
			name: "Граничный: Сравнение 2k и 2k+1 (k=3, 2k=6)",
			k:    3,
			chars: map[int]int{
				1: 6, // 6 <= 6 -> зачтено
				2: 7, // 7 > 6 -> не зачтено
			},
			expected: 1,
		},
		{
			name: "Присутствуют только старшие цифры (8 и 9)",
			k:    2,
			chars: map[int]int{
				8: 4,
				9: 3,
			},
			expected: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := countScore(tt.k, tt.chars)
			if got != tt.expected {
				t.Errorf("countScore(%d, %v) = %d; ожидалось %d", tt.k, tt.chars, got, tt.expected)
			}
		})
	}
}

package service

import "math"

// Вспомогательная функция для округления до 2 знаков
func roundToTwoDecimals(value float64) float64 {
	return math.Round(value*100) / 100
}

package mathutils

// Factorial вычисляет факториал числа n
// n int - число, факториал которого нужно вычислить
//  возвращает результат вычисления факториала
func Factorial(n int) int {
	if n < 0 {
		return 0 // Факториал не определен для отрицательных чисел
	}
	if n == 0 || n == 1 {
		return 1
	}

	result := 1
	for i := 2; i <= n; i++ {
		result *= i
	}
	return result
}

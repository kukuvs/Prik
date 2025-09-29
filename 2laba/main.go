package main

import (
	"fmt"
)

// Проверка четности числа
func isEven(n int) bool {
	return n%2 == 0
}

// Определение знака числа
func numberSign(n int) string {
	if n > 0 {
		return "Positive"
	} else if n < 0 {
		return "Negative"
	} else {
		return "Zero"
	}
}

// Длина строки
func strLen(s string) int {
	return len(s)
}

// Структура Rectangle с методом для площади
type Rectangle struct {
	Width  float64
	Height float64
}

func (r Rectangle) Area() float64 {
	return r.Width * r.Height
}

// Среднее двух целых чисел
func avg(a, b int) float64 {
	return float64(a+b) / 2.0
}

func main() {
	// 1. Четное/нечетное
	var n int
	fmt.Print("Введите число: ")
	fmt.Scan(&n)
	if isEven(n) {
		fmt.Println("Число четное")
	} else {
		fmt.Println("Число нечетное")
	}

	// 2. Positive/Negative/Zero
	fmt.Println("Знак числа:", numberSign(n))

	// 3. Числа от 1 до 10
	fmt.Println("Числа от 1 до 10:")
	for i := 1; i <= 10; i++ {
		fmt.Print(i, " ")
	}
	fmt.Println()

	// 4. Длина строки
	var s string
	fmt.Print("Введите строку: ")
	fmt.Scanln(&s)
	fmt.Println("Длина строки:", strLen(s))

	// 5. Rectangle и площадь
	r := Rectangle{Width: 3, Height: 4}
	fmt.Println("Площадь прямоугольника 3x4:", r.Area())

	// 6. Среднее двух чисел
	a, b := 4, 8
	fmt.Printf("Среднее для %d и %d: %.2f\n", a, b, avg(a, b))
}

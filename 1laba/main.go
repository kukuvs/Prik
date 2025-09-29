package main

import (
	"fmt"
	"time"
)

// Функция для вычисления суммы и разности двух чисел с плавающей запятой
func sumAndDiff(a, b float64) (float64, float64) {
	return a + b, a - b
}

// Функция для вычисления среднего значения трех чисел
func average(a, b, c float64) float64 {
	return (a + b + c) / 3.0
}

func main() {
	fmt.Println("=== Лабораторная работа 1 ===")

	// 1. Вывод текущего времени и даты
	fmt.Println("\n1. Текущее время и дата:")
	now := time.Now()
	fmt.Printf("Дата и время: %s\n", now.Format("2006-01-02 15:04:05"))

	// 2. Создание переменных различных типов
	fmt.Println("\n2. Переменные различных типов:")
	var intVar int = 42
	var floatVar float64 = 3.14159
	var stringVar string = "Привет, Go!"
	var boolVar bool = true

	fmt.Printf("int: %d\n", intVar)
	fmt.Printf("float64: %.5f\n", floatVar)
	fmt.Printf("string: %s\n", stringVar)
	fmt.Printf("bool: %t\n", boolVar)

	// 3. Краткая форма объявления переменных
	fmt.Println("\n3. Краткая форма объявления переменных:")
	shortInt := 100
	shortFloat := 2.718
	shortString := "Короткая форма"
	shortBool := false

	fmt.Printf("shortInt := %d\n", shortInt)
	fmt.Printf("shortFloat := %.3f\n", shortFloat)
	fmt.Printf("shortString := %s\n", shortString)
	fmt.Printf("shortBool := %t\n", shortBool)

	// 4. Арифметические операции с двумя целыми числами
	fmt.Println("\n4. Арифметические операции:")
	num1, num2 := 15, 7
	fmt.Printf("Числа: %d и %d\n", num1, num2)
	fmt.Printf("Сложение: %d + %d = %d\n", num1, num2, num1+num2)
	fmt.Printf("Вычитание: %d - %d = %d\n", num1, num2, num1-num2)
	fmt.Printf("Умножение: %d * %d = %d\n", num1, num2, num1*num2)
	fmt.Printf("Деление: %d / %d = %d\n", num1, num2, num1/num2)
	fmt.Printf("Остаток: %d %% %d = %d\n", num1, num2, num1%num2)

	// 5. Функция для суммы и разности чисел с плавающей запятой
	fmt.Println("\n5. Сумма и разность чисел с плавающей запятой:")
	float1, float2 := 10.5, 3.2
	sum, diff := sumAndDiff(float1, float2)
	fmt.Printf("Числа: %.1f и %.1f\n", float1, float2)
	fmt.Printf("Сумма: %.5f\n", sum)
	fmt.Printf("Разность: %.5f\n", diff)

	// 6. Среднее значение трех чисел
	fmt.Println("\n6. Среднее значение трех чисел:")
	a, b, c := 12.0, 8.5, 15.3
	avg := average(a, b, c)
	fmt.Printf("Числа: %.1f, %.1f, %.1f\n", a, b, c)
	fmt.Printf("Среднее значение: %.5f\n", avg)
}

package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func main() {
	fmt.Println("=== Лабораторная работа 4 ===")

	// Задание 1-3: Работа с картой людей и их возрастов
	fmt.Println("\n1. Работа с картой людей:")
	people := make(map[string]int)

	// Добавляем изначальных людей
	people["Алексей"] = 25
	people["Мария"] = 30
	people["Иван"] = 22

	// Добавляем нового человека
	people["Елена"] = 28

	// Выводим все записи
	fmt.Println("Все записи:")
	for name, age := range people {
		fmt.Printf("%s: %d лет\n", name, age)
	}

	// Вычисляем средний возраст
	avgAge := averageAge(people)
	fmt.Printf("Средний возраст: %.2f лет\n", avgAge)

	// Удаляем запись
	deleteName := "Иван"
	delete(people, deleteName)
	fmt.Printf("Удалена запись: %s\n", deleteName)

	fmt.Println("Записи после удаления:")
	for name, age := range people {
		fmt.Printf("%s: %d лет\n", name, age)
	}

	// Задание 4: Строка в верхнем регистре
	fmt.Println("\n4. Ввод строки и вывод в верхнем регистре:")
	fmt.Print("Введите строку: ")
	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)
	fmt.Printf("Строка в верхнем регистре: %s\n", strings.ToUpper(input))

	// Задание 5: Сумма чисел
	fmt.Println("\n5. Сумма введенных чисел:")
	fmt.Print("Введите числа через пробел: ")
	numbersInput, _ := reader.ReadString('\n')
	numbersInput = strings.TrimSpace(numbersInput)

	sum := 0
	numbers := strings.Fields(numbersInput)
	for _, numStr := range numbers {
		if num, err := strconv.Atoi(numStr); err == nil {
			sum += num
		}
	}
	fmt.Printf("Сумма чисел: %d\n", sum)

	// Задание 6: Массив в обратном порядке
	fmt.Println("\n6. Массив в обратном порядке:")
	fmt.Print("Введите целые числа через пробел: ")
	arrayInput, _ := reader.ReadString('\n')
	arrayInput = strings.TrimSpace(arrayInput)

	var array []int
	arrayNumbers := strings.Fields(arrayInput)
	for _, numStr := range arrayNumbers {
		if num, err := strconv.Atoi(numStr); err == nil {
			array = append(array, num)
		}
	}

	fmt.Print("Массив в обратном порядке: ")
	for i := len(array) - 1; i >= 0; i-- {
		fmt.Printf("%d ", array[i])
	}
	fmt.Println()
}

// Функция для вычисления среднего возраста
func averageAge(people map[string]int) float64 {
	if len(people) == 0 {
		return 0
	}

	sum := 0
	for _, age := range people {
		sum += age
	}

	return float64(sum) / float64(len(people))
}

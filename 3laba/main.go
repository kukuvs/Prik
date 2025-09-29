package main

import (
	"3lab/mathutils"
	"3lab/stringutils"
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func main() {
	reader := bufio.NewReader(os.Stdin)

	// Задание 2: Использование пакета mathutils для вычисления факториала
	fmt.Print("Введите число для вычисления факториала: ")
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)
	n, err := strconv.Atoi(input)
	if err != nil {
		fmt.Println("Ошибка ввода числа")
		return
	}

	factorial := mathutils.Factorial(n)
	fmt.Printf("Факториал %d = %d\n", n, factorial)

	// Задание 3: Использование пакета stringutils для переворота строки
	fmt.Print("Введите строку для переворота: ")
	str, _ := reader.ReadString('\n')
	str = strings.TrimSpace(str)

	reversed := stringutils.Reverse(str)
	fmt.Printf("Перевернутая строка: %s\n", reversed)

	// Задание 4: Создание массива из 5 целых чисел
	var arr [5]int
	fmt.Println("Введите 5 целых чисел:")
	for i := 0; i < 5; i++ {
		fmt.Printf("Элемент %d: ", i+1)
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)
		num, err := strconv.Atoi(input)
		if err != nil {
			fmt.Println("Ошибка ввода числа, используется 0")
			arr[i] = 0
		} else {
			arr[i] = num
		}
	}

	fmt.Println("Массив:")
	for i, v := range arr {
		fmt.Printf("arr[%d] = %d\n", i, v)
	}

	// Задание 5: Создание среза из массива и операции с ним
	slice := arr[1:4] // Создаем срез из элементов с индексами 1-3
	fmt.Printf("Исходный срез: %v\n", slice)

	// Добавление элементов
	slice = append(slice, 100, 200)
	fmt.Printf("После добавления элементов: %v\n", slice)

	// Удаление элемента (удаляем элемент с индексом 1)
	if len(slice) > 1 {
		slice = append(slice[:1], slice[2:]...)
		fmt.Printf("После удаления элемента с индексом 1: %v\n", slice)
	}

	// Задание 6: Поиск самой длинной строки в срезе
	strings := []string{"привет", "мир"}
	fmt.Printf("Срез строк: %v\n", strings)

	longest := findLongestString(strings)
	fmt.Printf("Самая длинная строка: %s (длина: %d)\n", longest, len([]rune(longest)))
}

// Функция для поиска самой длинной строки
func findLongestString(strings []string) string {
	if len(strings) == 0 {
		return ""
	}

	longest := strings[0]
	for _, str := range strings {
		if len([]rune(str)) > len([]rune(longest)) {
			longest = str
		}
	}
	return longest
}

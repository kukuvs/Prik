package main

import (
	"fmt"
	"math/rand"
	"strings"
	"sync"
	"time"
)

// ============================================
// ЗАДАНИЕ 1: Создание и запуск горутин
// ============================================

func factorial(n int) {
	fmt.Printf("Начало вычисления факториала %d\n", n)
	time.Sleep(1 * time.Second)

	result := 1
	for i := 1; i <= n; i++ {
		result *= i
	}

	fmt.Printf("Факториал %d = %d\n", n, result)
}

func generateRandom(count int) {
	fmt.Println("Начало генерации случайных чисел")
	time.Sleep(500 * time.Millisecond)

	rand.Seed(time.Now().UnixNano())
	for i := 0; i < count; i++ {
		fmt.Printf("Случайное число %d: %d\n", i+1, rand.Intn(100))
		time.Sleep(200 * time.Millisecond)
	}
}

func sumSeries(n int) {
	fmt.Printf("Начало вычисления суммы ряда 1..%d\n", n)
	time.Sleep(800 * time.Millisecond)

	sum := 0
	for i := 1; i <= n; i++ {
		sum += i
	}

	fmt.Printf("Сумма ряда 1..%d = %d\n", n, sum)
}

func task1() {
	fmt.Println("\n=== ЗАДАНИЕ 1: Создание и запуск горутин ===")
	fmt.Println("Запуск параллельных горутин:")

	go factorial(5)
	go generateRandom(3)
	go sumSeries(10)

	time.Sleep(3 * time.Second)
	fmt.Println("\nВсе горутины завершены")
}

// ============================================
// ЗАДАНИЕ 2: Использование каналов
// ============================================

func fibonacci(n int, ch chan int) {
	a, b := 0, 1

	for i := 0; i < n; i++ {
		ch <- a
		a, b = b, a+b
		time.Sleep(300 * time.Millisecond)
	}

	close(ch)
}

func reader(ch chan int) {
	for num := range ch {
		fmt.Printf("Получено число Фибоначчи: %d\n", num)
	}
	fmt.Println("Канал закрыт, чтение завершено")
}

func task2() {
	fmt.Println("\n=== ЗАДАНИЕ 2: Использование каналов для передачи данных ===")

	ch := make(chan int)

	fmt.Println("Запуск генератора Фибоначчи:")

	go fibonacci(10, ch)
	reader(ch)

	fmt.Println("\nЗадание 2 завершено")
}

// ============================================
// ЗАДАНИЕ 3: Применение select
// ============================================

func randomGenerator(ch chan int, done chan bool) {
	rand.Seed(time.Now().UnixNano())

	for i := 0; i < 10; i++ {
		num := rand.Intn(100)
		ch <- num
		time.Sleep(500 * time.Millisecond)
	}

	done <- true
}

func parityChecker(inputCh chan int, outputCh chan string, done chan bool) {
	for num := range inputCh {
		var parity string
		if num%2 == 0 {
			parity = fmt.Sprintf("Число %d - ЧЁТНОЕ", num)
		} else {
			parity = fmt.Sprintf("Число %d - НЕЧЁТНОЕ", num)
		}
		outputCh <- parity
		time.Sleep(200 * time.Millisecond)
	}

	done <- true
}

func task3() {
	fmt.Println("\n=== ЗАДАНИЕ 3: Применение select для управления каналами ===")

	numberCh := make(chan int)
	parityCh := make(chan string)
	done1 := make(chan bool)
	done2 := make(chan bool)

	fmt.Println("Демонстрация работы select с несколькими каналами:")

	go randomGenerator(numberCh, done1)
	go parityChecker(numberCh, parityCh, done2)

	activeChannels := 2

	for activeChannels > 0 {
		select {
		case result := <-parityCh:
			fmt.Println("✓", result)

		case <-done1:
			fmt.Println("\n[Генератор завершён]")
			close(numberCh)
			activeChannels--

		case <-done2:
			fmt.Println("[Проверка чётности завершена]")
			activeChannels--

		case <-time.After(2 * time.Second):
			fmt.Println("\n⏱ Таймаут ожидания данных")
			activeChannels = 0
		}
	}

	fmt.Println("\nЗадание 3 завершено")
}

// ============================================
// ЗАДАНИЕ 4: Синхронизация с мьютексами
// ============================================

var (
	counter int
	mutex   sync.Mutex
)

func incrementWithoutMutex(wg *sync.WaitGroup, id int) {
	defer wg.Done()

	for i := 0; i < 1000; i++ {
		counter++
	}
	fmt.Printf("Горутина %d завершена (без мьютекса)\n", id)
}

func incrementWithMutex(wg *sync.WaitGroup, id int) {
	defer wg.Done()

	for i := 0; i < 1000; i++ {
		mutex.Lock()
		counter++
		mutex.Unlock()
	}
	fmt.Printf("Горутина %d завершена (с мьютексом)\n", id)
}

func task4() {
	fmt.Println("\n=== ЗАДАНИЕ 4: Синхронизация с помощью мьютексов ===")
	fmt.Println("\n--- Демонстрация работы БЕЗ мьютекса ---")
	counter = 0
	var wg1 sync.WaitGroup

	for i := 1; i <= 5; i++ {
		wg1.Add(1)
		go incrementWithoutMutex(&wg1, i)
	}

	wg1.Wait()
	fmt.Printf("Итоговое значение счётчика БЕЗ мьютекса: %d (ожидалось: 5000)\n", counter)
	fmt.Println("❌ Результат неверный из-за гонки данных!")

	time.Sleep(1 * time.Second)

	fmt.Println("--- Демонстрация работы С мьютексом ---")
	counter = 0
	var wg2 sync.WaitGroup

	for i := 1; i <= 5; i++ {
		wg2.Add(1)
		go incrementWithMutex(&wg2, i)
	}

	wg2.Wait()
	fmt.Printf("Итоговое значение счётчика С мьютексом: %d (ожидалось: 5000)\n", counter)
	fmt.Println("✓ Результат корректный благодаря мьютексу!")
}

// ============================================
// ЗАДАНИЕ 5: Многопоточный калькулятор
// ============================================

type Request struct {
	operation string
	a, b      float64
	result    chan float64
}

func calculator(requests chan Request, wg *sync.WaitGroup) {
	defer wg.Done()

	for req := range requests {
		var result float64

		switch req.operation {
		case "+":
			result = req.a + req.b
		case "-":
			result = req.a - req.b
		case "*":
			result = req.a * req.b
		case "/":
			if req.b != 0 {
				result = req.a / req.b
			} else {
				fmt.Println("Ошибка: деление на ноль!")
				req.result <- 0
				continue
			}
		default:
			fmt.Println("Неизвестная операция:", req.operation)
			req.result <- 0
			continue
		}

		req.result <- result
	}
}

func clientRequest(id int, operation string, a, b float64, requests chan Request, wg *sync.WaitGroup) {
	defer wg.Done()

	resultCh := make(chan float64)

	requests <- Request{
		operation: operation,
		a:         a,
		b:         b,
		result:    resultCh,
	}

	result := <-resultCh
	fmt.Printf("Клиент %d: %.2f %s %.2f = %.2f\n", id, a, operation, b, result)
}

func task5() {
	fmt.Println("\n=== ЗАДАНИЕ 5: Многопоточный калькулятор ===")

	requests := make(chan Request, 10)
	var serverWg sync.WaitGroup
	var clientWg sync.WaitGroup

	for i := 0; i < 3; i++ {
		serverWg.Add(1)
		go calculator(requests, &serverWg)
	}

	operations := []struct {
		op   string
		a, b float64
	}{
		{"+", 10, 5},
		{"-", 20, 8},
		{"*", 7, 6},
		{"/", 100, 4},
		{"+", 15.5, 4.5},
		{"*", 3.14, 2},
		{"/", 50, 2},
		{"-", 100, 25},
	}

	for i, op := range operations {
		clientWg.Add(1)
		go clientRequest(i+1, op.op, op.a, op.b, requests, &clientWg)
	}

	clientWg.Wait()
	close(requests)
	serverWg.Wait()

	fmt.Println("\nКалькулятор завершил работу")
}

// ============================================
// ЗАДАНИЕ 6: Создание пула воркеров
// ============================================

type Task struct {
	id   int
	text string
}

type Result struct {
	taskID   int
	reversed string
}

func reverseString(s string) string {
	runes := []rune(s)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}

func worker(id int, tasks <-chan Task, results chan<- Result, wg *sync.WaitGroup) {
	defer wg.Done()

	for task := range tasks {
		fmt.Printf("Воркер %d обрабатывает задачу %d: %s\n", id, task.id, task.text)

		reversed := reverseString(task.text)

		results <- Result{
			taskID:   task.id,
			reversed: reversed,
		}
	}

	fmt.Printf("Воркер %d завершил работу\n", id)
}

func resultCollector(results <-chan Result, done chan bool, numTasks int) {
	collected := make([]Result, 0, numTasks)

	for i := 0; i < numTasks; i++ {
		result := <-results
		collected = append(collected, result)
	}

	fmt.Println("\n=== РЕЗУЛЬТАТЫ ОБРАБОТКИ ===")
	for _, res := range collected {
		fmt.Printf("Задача %d: %s\n", res.taskID, res.reversed)
	}

	done <- true
}

func task6() {
	fmt.Println("\n=== ЗАДАНИЕ 6: Создание пула воркеров ===")

	inputStrings := []string{
		"Hello World",
		"Golang",
		"Concurrency",
		"Worker Pool",
		"Goroutines",
		"Channels",
		"Synchronization",
		"Mutex",
	}

	numWorkers := 3
	numTasks := len(inputStrings)

	fmt.Printf("Количество воркеров: %d\n", numWorkers)
	fmt.Printf("Количество задач: %d\n\n", numTasks)

	tasks := make(chan Task, numTasks)
	results := make(chan Result, numTasks)
	done := make(chan bool)

	var wg sync.WaitGroup

	for i := 1; i <= numWorkers; i++ {
		wg.Add(1)
		go worker(i, tasks, results, &wg)
	}

	go resultCollector(results, done, numTasks)

	for i, str := range inputStrings {
		tasks <- Task{
			id:   i + 1,
			text: str,
		}
	}

	close(tasks)
	wg.Wait()
	close(results)
	<-done

	fmt.Println("\n✓ Все задачи обработаны успешно!")
}

// ============================================
// ГЛАВНОЕ МЕНЮ
// ============================================

func printMenu() {

	fmt.Println("\nВыберите задание:")
	fmt.Println("  1 - Создание и запуск горутин")
	fmt.Println("  2 - Использование каналов для передачи данных")
	fmt.Println("  3 - Применение select для управления каналами")
	fmt.Println("  4 - Синхронизация с помощью мьютексов")
	fmt.Println("  5 - Многопоточный калькулятор")
	fmt.Println("  6 - Создание пула воркеров")
	fmt.Println("  7 - Выполнить все задания подряд")
	fmt.Println("  0 - Выход")
	fmt.Print("\nВведите номер задания: ")
}

func runAllTasks() {

	task1()
	time.Sleep(1 * time.Second)

	task2()
	time.Sleep(1 * time.Second)

	task3()
	time.Sleep(1 * time.Second)

	task4()
	time.Sleep(1 * time.Second)

	task5()
	time.Sleep(1 * time.Second)

	task6()

	fmt.Println("\n✓ Все задания выполнены успешно!")
}

func main() {
	var choice int

	for {
		printMenu()

		_, err := fmt.Scanf("%d", &choice)
		if err != nil {
			fmt.Println("\n❌ Ошибка ввода! Введите число от 0 до 7.")
			var discard string
			fmt.Scanln(&discard)
			continue
		}

		switch choice {
		case 1:
			task1()
		case 2:
			task2()
		case 3:
			task3()
		case 4:
			task4()
		case 5:
			task5()
		case 6:
			task6()
		case 7:
			runAllTasks()
		case 0:
			fmt.Println("А всё тут пусто ")
			return
		default:
			fmt.Println("\n❌ Неверный выбор! Введите число от 0 до 7.")
		}

		fmt.Println("\n" + strings.Repeat("─", 50))
		fmt.Print("Нажмите Enter для продолжения...")
		fmt.Scanln()
		fmt.Scanln()
	}
}

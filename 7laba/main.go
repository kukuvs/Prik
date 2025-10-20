package main

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/gorilla/websocket"
)

// ============================================
// ЗАДАНИЕ 1: Простой TCP-сервер
// ============================================

func task1() {
	fmt.Println("\n=== ЗАДАНИЕ 1: TCP-сервер ===")
	PORT := ":8080"

	fmt.Printf("Запуск TCP сервера на порту %s\n", PORT)

	listener, err := net.Listen("tcp", PORT)
	if err != nil {
		fmt.Println("Ошибка запуска сервера:", err)
		return
	}
	defer listener.Close()

	fmt.Println("✓ Сервер запущен. Ожидание подключения...")
	fmt.Println("Запустите клиент (Задание 2) в другом терминале")
	fmt.Println("Нажмите Ctrl+C для остановки")

	conn, err := listener.Accept()
	if err != nil {
		fmt.Println("Ошибка принятия соединения:", err)
		return
	}
	defer conn.Close()

	fmt.Printf("✓ Клиент подключен: %s\n", conn.RemoteAddr().String())

	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		message := scanner.Text()
		fmt.Printf("Получено: %s\n", message)

		if strings.TrimSpace(message) == "STOP" {
			fmt.Println("Получена команда остановки")
			break
		}

		response := "Сервер: сообщение получено\n"
		conn.Write([]byte(response))
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Ошибка чтения:", err)
	}

	fmt.Println("Соединение закрыто")
}

// ============================================
// ЗАДАНИЕ 2: TCP-клиент
// ============================================

func task2() {
	fmt.Println("\n=== ЗАДАНИЕ 2: TCP-клиент ===")
	SERVER := "localhost:8080"

	fmt.Printf("Подключение к серверу %s\n", SERVER)

	conn, err := net.Dial("tcp", SERVER)
	if err != nil {
		fmt.Println("❌ Ошибка подключения:", err)
		fmt.Println("Убедитесь, что сервер запущен (Задание 1 или 3)")
		return
	}
	defer conn.Close()

	fmt.Println("✓ Подключено к серверу")
	fmt.Println("Введите сообщение (STOP для выхода):")

	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Print(">> ")
		if !scanner.Scan() {
			break
		}

		message := scanner.Text()

		_, err = fmt.Fprintf(conn, "%s", message+"\n")
		if err != nil {
			fmt.Println("Ошибка отправки:", err)
			return
		}

		if strings.TrimSpace(message) == "STOP" {
			fmt.Println("Завершение работы...")
			return
		}

		response, err := bufio.NewReader(conn).ReadString('\n')
		if err != nil {
			fmt.Println("Ошибка чтения ответа:", err)
			return
		}

		fmt.Printf("Ответ: %s", response)
	}
}

// ============================================
// ЗАДАНИЕ 3: Асинхронная обработка
// ============================================

var (
	clientCount int
	mu          sync.Mutex
	wg          sync.WaitGroup
)

func handleClient(ctx context.Context, conn net.Conn, clientID int) {
	defer wg.Done()
	defer conn.Close()

	fmt.Printf("[Клиент %d] Подключен: %s\n", clientID, conn.RemoteAddr().String())

	scanner := bufio.NewScanner(conn)

	for {
		select {
		case <-ctx.Done():
			fmt.Printf("[Клиент %d] Graceful shutdown\n", clientID)
			conn.Write([]byte("Сервер завершает работу\n"))
			return

		default:
			conn.SetReadDeadline(time.Now().Add(1 * time.Second))

			if scanner.Scan() {
				message := scanner.Text()
				fmt.Printf("[Клиент %d] Сообщение: %s\n", clientID, message)

				response := fmt.Sprintf("Сервер [%d]: получено\n", clientID)
				conn.Write([]byte(response))
			} else {
				if err := scanner.Err(); err != nil {
					if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
						continue
					}
					fmt.Printf("[Клиент %d] Ошибка: %v\n", clientID, err)
					return
				}
				fmt.Printf("[Клиент %d] Отключен\n", clientID)
				return
			}
		}
	}
}

func task3() {
	fmt.Println("\n=== ЗАДАНИЕ 3: Асинхронная обработка с Graceful Shutdown ===")
	PORT := ":8080"

	fmt.Printf("Запуск сервера на порту %s\n", PORT)

	listener, err := net.Listen("tcp", PORT)
	if err != nil {
		fmt.Println("Ошибка запуска:", err)
		return
	}
	defer listener.Close()

	fmt.Println("✓ Сервер запущен (поддержка множества клиентов)")
	fmt.Println("Нажмите Ctrl+C для graceful shutdown")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigChan
		fmt.Println("\n\n⚠ Получен сигнал остановки")
		fmt.Println("Ожидание завершения активных соединений...")
		cancel()
		listener.Close()
	}()

	for {
		conn, err := listener.Accept()
		if err != nil {
			select {
			case <-ctx.Done():
				fmt.Println("Сервер остановлен")
				wg.Wait()
				fmt.Println("✓ Все соединения завершены")
				return
			default:
				fmt.Println("Ошибка принятия соединения:", err)
				continue
			}
		}

		mu.Lock()
		clientCount++
		currentClientID := clientCount
		mu.Unlock()

		wg.Add(1)
		go handleClient(ctx, conn, currentClientID)
	}
}

// ============================================
// ЗАДАНИЕ 4: HTTP-сервер
// ============================================

func helloHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
		return
	}

	fmt.Fprintf(w, "Привет! Добро пожаловать на HTTP-сервер Go!\n")
	fmt.Fprintf(w, "Время: %s\n", time.Now().Format("15:04:05"))
}

func dataHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Ожидается POST запрос", http.StatusMethodNotAllowed)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Ошибка чтения тела запроса", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	var data map[string]interface{}
	if err := json.Unmarshal(body, &data); err != nil {
		http.Error(w, "Неверный формат JSON", http.StatusBadRequest)
		return
	}

	fmt.Println("\n--- Получены данные POST /data ---")
	jsonData, _ := json.MarshalIndent(data, "", "  ")
	fmt.Println(string(jsonData))
	fmt.Println("----------------------------------")

	w.Header().Set("Content-Type", "application/json")
	response := map[string]interface{}{
		"status":  "success",
		"message": "Данные получены",
		"data":    data,
	}
	json.NewEncoder(w).Encode(response)
}

func task4() {
	fmt.Println("\n=== ЗАДАНИЕ 4: HTTP-сервер ===")

	http.HandleFunc("/hello", helloHandler)
	http.HandleFunc("/data", dataHandler)

	PORT := ":8081"
	fmt.Printf("HTTP сервер запущен на http://localhost%s\n", PORT)
	fmt.Println("\nДоступные эндпоинты:")
	fmt.Println("  GET  http://localhost:8081/hello")
	fmt.Println("  POST http://localhost:8081/data")
	fmt.Println("\nПример запроса:")
	fmt.Println("  curl http://localhost:8081/hello")
	fmt.Println("  curl -X POST http://localhost:8081/data -d '{\"name\":\"Test\",\"value\":123}'")
	fmt.Println("\nНажмите Ctrl+C для остановки")

	if err := http.ListenAndServe(PORT, nil); err != nil {
		fmt.Println("Ошибка запуска сервера:", err)
	}
}

// ============================================
// ЗАДАНИЕ 5: Маршрутизация и Middleware
// ============================================

func loggingMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		fmt.Printf("[%s] %s %s - Начало обработки\n",
			start.Format("15:04:05"),
			r.Method,
			r.URL.Path,
		)

		next(w, r)

		duration := time.Since(start)
		fmt.Printf("[%s] %s %s - Завершено за %v\n",
			time.Now().Format("15:04:05"),
			r.Method,
			r.URL.Path,
			duration,
		)
	}
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Домашняя страница\n")
	time.Sleep(100 * time.Millisecond) // Имитация работы
}

func aboutHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Страница О нас\n")
	fmt.Fprintf(w, "Лабораторная работа №7\n")
	time.Sleep(50 * time.Millisecond)
}

func apiHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	response := map[string]string{
		"status":  "ok",
		"message": "API работает",
		"version": "1.0",
	}
	json.NewEncoder(w).Encode(response)
	time.Sleep(80 * time.Millisecond)
}

func task5() {
	fmt.Println("\n=== ЗАДАНИЕ 5: Маршрутизация и Middleware ===")

	http.HandleFunc("/", loggingMiddleware(homeHandler))
	http.HandleFunc("/about", loggingMiddleware(aboutHandler))
	http.HandleFunc("/api", loggingMiddleware(apiHandler))

	PORT := ":8082"
	fmt.Printf("HTTP сервер с middleware запущен на http://localhost%s\n", PORT)
	fmt.Println("\nДоступные маршруты:")
	fmt.Println("  GET http://localhost:8082/")
	fmt.Println("  GET http://localhost:8082/about")
	fmt.Println("  GET http://localhost:8082/api")
	fmt.Println("\nMiddleware логирует каждый запрос с временем выполнения")
	fmt.Println("Нажмите Ctrl+C для остановки")

	if err := http.ListenAndServe(PORT, nil); err != nil {
		fmt.Println("Ошибка запуска сервера:", err)
	}
}

// ============================================
// ЗАДАНИЕ 6: WebSocket чат
// ============================================

var (
	upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true // В продакшене нужна проверка origin
		},
	}
	clients   = make(map[*websocket.Conn]string)
	clientsMu sync.Mutex
	broadcast = make(chan Message)
)

type Message struct {
	Username string `json:"username"`
	Text     string `json:"text"`
	Time     string `json:"time"`
}

func handleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Ошибка WebSocket upgrade:", err)
		return
	}
	defer conn.Close()

	username := r.URL.Query().Get("username")
	if username == "" {
		username = "Аноним"
	}

	clientsMu.Lock()
	clients[conn] = username
	clientsMu.Unlock()

	fmt.Printf("[WebSocket] Подключен: %s\n", username)

	joinMsg := Message{
		Username: "Система",
		Text:     fmt.Sprintf("%s присоединился к чату", username),
		Time:     time.Now().Format("15:04:05"),
	}
	broadcast <- joinMsg

	defer func() {
		clientsMu.Lock()
		delete(clients, conn)
		clientsMu.Unlock()

		leaveMsg := Message{
			Username: "Система",
			Text:     fmt.Sprintf("%s покинул чат", username),
			Time:     time.Now().Format("15:04:05"),
		}
		broadcast <- leaveMsg
	}()

	for {
		var msg Message
		err := conn.ReadJSON(&msg)
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("Ошибка: %v", err)
			}
			break
		}

		msg.Username = username
		msg.Time = time.Now().Format("15:04:05")
		fmt.Printf("[Чат] %s: %s\n", username, msg.Text)

		broadcast <- msg
	}
}

func handleBroadcast() {
	for {
		msg := <-broadcast

		clientsMu.Lock()
		for client := range clients {
			err := client.WriteJSON(msg)
			if err != nil {
				log.Printf("Ошибка отправки: %v", err)
				client.Close()
				delete(clients, client)
			}
		}
		clientsMu.Unlock()
	}
}

func chatPageHandler(w http.ResponseWriter, r *http.Request) {
	html := `<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <title>WebSocket Чат</title>
    <style>
        body { font-family: Arial, sans-serif; max-width: 600px; margin: 50px auto; padding: 20px; }
        #messages { border: 1px solid #ccc; height: 400px; overflow-y: scroll; padding: 10px; margin-bottom: 10px; }
        .message { margin: 5px 0; }
        .system { color: #666; font-style: italic; }
        .username { font-weight: bold; color: #0066cc; }
        .time { color: #999; font-size: 0.9em; }
        #input { width: 80%; padding: 10px; }
        #send { width: 18%; padding: 10px; }
    </style>
</head>
<body>
    <h1>WebSocket Чат</h1>
    <div id="messages"></div>
    <input id="input" type="text" placeholder="Введите сообщение..." />
    <button id="send">Отправить</button>

    <script>
        const username = prompt("Введите ваше имя:") || "Аноним";
        const ws = new WebSocket("ws://localhost:8083/ws?username=" + encodeURIComponent(username));
        const messages = document.getElementById("messages");
        const input = document.getElementById("input");
        const send = document.getElementById("send");

        ws.onmessage = function(event) {
            const msg = JSON.parse(event.data);
            const div = document.createElement("div");
            div.className = "message";
            
            if (msg.username === "Система") {
                div.className += " system";
                div.innerHTML = msg.text + " <span class='time'>[" + msg.time + "]</span>";
            } else {
                div.innerHTML = "<span class='username'>" + msg.username + "</span>: " + 
                                msg.text + " <span class='time'>[" + msg.time + "]</span>";
            }
            
            messages.appendChild(div);
            messages.scrollTop = messages.scrollHeight;
        };

        send.onclick = function() {
            if (input.value) {
                ws.send(JSON.stringify({text: input.value}));
                input.value = "";
            }
        };

        input.addEventListener("keypress", function(event) {
            if (event.key === "Enter") {
                send.click();
            }
        });
    </script>
</body>
</html>`
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write([]byte(html))
}

func task6() {
	fmt.Println("\n=== ЗАДАНИЕ 6: WebSocket Чат ===")

	go handleBroadcast()

	http.HandleFunc("/", chatPageHandler)
	http.HandleFunc("/ws", handleWebSocket)

	PORT := ":8083"
	fmt.Printf("WebSocket чат-сервер запущен на http://localhost%s\n", PORT)
	fmt.Println("\nОткройте в браузере: http://localhost:8083")
	fmt.Println("Можно открыть несколько вкладок для тестирования")
	fmt.Println("Нажмите Ctrl+C для остановки")

	if err := http.ListenAndServe(PORT, nil); err != nil {
		fmt.Println("Ошибка запуска сервера:", err)
	}
}

// ============================================
// ГЛАВНОЕ МЕНЮ
// ============================================

func printMenu() {
	fmt.Println("\n╔════════════════════════════════════════════════╗")
	fmt.Println("║   ЛАБОРАТОРНАЯ РАБОТА 7 - СЕТЕВОЕ ПРОГРАММИРОВАНИЕ  ║")
	fmt.Println("╚════════════════════════════════════════════════╝")
	fmt.Println("\nВыберите задание:")
	fmt.Println("  1 - TCP-сервер (простой)")
	fmt.Println("  2 - TCP-клиент")
	fmt.Println("  3 - TCP-сервер с асинхронной обработкой")
	fmt.Println("  4 - HTTP-сервер (GET, POST)")
	fmt.Println("  5 - HTTP с маршрутизацией и middleware")
	fmt.Println("  6 - WebSocket чат-сервер")
	fmt.Println("  0 - Выход")
	fmt.Print("\nВведите номер задания: ")
}

func main() {
	var choice int

	for {
		printMenu()

		_, err := fmt.Scanf("%d", &choice)
		if err != nil {
			fmt.Println("\n❌ Ошибка ввода! Введите число от 0 до 6.")
			var discard string
			fmt.Scanln(&discard)
			continue
		}

		// Очищаем буфер
		fmt.Scanln()

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
		case 0:
			fmt.Println("\n╔════════════════════════════════════════════════╗")
			fmt.Println("║         Программа завершена. До встречи!      ║")
			fmt.Println("╚════════════════════════════════════════════════╝")
			return
		default:
			fmt.Println("\n❌ Неверный выбор! Введите число от 0 до 6.")
		}

		fmt.Println("\n" + strings.Repeat("─", 50))
		fmt.Print("Нажмите Enter для возврата в меню...")
		bufio.NewReader(os.Stdin).ReadBytes('\n')
	}
}

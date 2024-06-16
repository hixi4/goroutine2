package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

// Генеруємо випадкові числа та надсилаємо їх у канал
func generateNumbers(numCh chan<- int, wg *sync.WaitGroup, limit int) {
	defer wg.Done()
	for i := 0; i < limit; i++ { // Працюємо, поки лічильник не досягне ліміту
		num := rand.Intn(10) // Генеруємо випадкове число від 0 до 9
		numCh <- num
		time.Sleep(time.Second) // Затримка для наочності
	}
	close(numCh) // Закриваємо канал після завершення генерації
}

// Отримуємо випадкові числа та знаходимо найбільше й найменше числа
func findMinMax(numCh <-chan int, minMaxCh chan<- [2]int, wg *sync.WaitGroup) {
	defer wg.Done()
	var min, max int
	min, max = int(^uint(0)>>1), -int(^uint(0)>>1)-1 // Ініціалізуємо min та max

	for num := range numCh { // Читаємо числа з каналу
		if num < min {
			min = num
		}
		if num > max {
			max = num
		}
		minMaxCh <- [2]int{min, max}
	}
	close(minMaxCh) // Закриваємо канал після завершення обробки
}

func main() {
	numCh := make(chan int)
	minMaxCh := make(chan [2]int)
	limit := 10 // Ліміт на кількість згенерованих елементів

	var wg sync.WaitGroup
	wg.Add(2)

	go generateNumbers(numCh, &wg, limit)
	go findMinMax(numCh, minMaxCh, &wg)

	// Виводимо результати з minMaxCh
	go func() {
		for minMax := range minMaxCh {
			fmt.Printf("Найменше число: %d, Найбільше число: %d\n", minMax[0], minMax[1])
		}
	}()

	// Запобігаємо завершенню програми
	wg.Wait()
}

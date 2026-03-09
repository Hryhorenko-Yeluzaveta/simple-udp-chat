package main

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
)

func sendMessage(connection *net.UDPConn) {
	friendAddr, err := net.ResolveUDPAddr("udp", "192.168.0.100:9000")
	if err != nil {
		log.Printf("Не вдалося визначити адресу іншого хоста: %v\n", err)
		return
	}

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		message := scanner.Text()
		if len(message) == 0 {
			continue
		}
		_, err = connection.WriteToUDP([]byte(message), friendAddr)
		if err != nil {
			fmt.Print("Не вдалося відправити повідомлення: ", err)
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Println("Помилка зчитування вводу: ", err)
	}
}

func main() {
	localAddr, err := net.ResolveUDPAddr("udp", "192.168.0.103:9000")
	if err != nil {
		log.Fatal("Помилка підготовки адреси: ", err)
	}
	connection, err := net.ListenUDP("udp", localAddr)
	if err != nil {
		log.Fatal("Помилка з'єднання: ", err)
	}
	defer connection.Close()

	closeSignal := make(chan os.Signal, 1)
	signal.Notify(closeSignal, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-closeSignal
		fmt.Println("\nРозриваю з'єднання")
		connection.Close()
	}()

	fmt.Println("Слухаємо на порту 9000. Напишіть повідомлення для відправки:")
	fmt.Println("---------------------------------------------------------")

	go sendMessage(connection)

	// Читаємо повідомлення відправлені з іншого хоста
	buffer := make([]byte, 1024)
	for {
		msg, addrFrom, err := connection.ReadFromUDP(buffer)
		if err != nil {
			if errors.Is(err, net.ErrClosed) {
				break
			}
			fmt.Printf("Помилка мережі: %v\n", err)
			break
		}
		fmt.Println("Отримано від: ", addrFrom, "Повідомлення: ", string(buffer[:msg]))
	}
}

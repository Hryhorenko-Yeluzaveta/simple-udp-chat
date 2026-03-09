package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	localAddr, err := net.ResolveUDPAddr("udp", ":9000")
	if err != nil {
		log.Fatal("Помилка прив'язки порту (bind): ", err)
	}
	connection, err := net.ListenUDP("udp", localAddr)
	if err != nil {
		log.Fatal("Помилка з'єднання: ", err)
	}

	closeSignal := make(chan os.Signal, 1)
	signal.Notify(closeSignal, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-closeSignal
		connection.Close()
		fmt.Println("З'єднання розірвано")
		os.Exit(0)
	}()

	fmt.Println("Слухаємо повідомлення на 9000")

	buffer := make([]byte, 1024)
	for {
		msg, addrFrom, err := connection.ReadFromUDP(buffer)
		if err != nil {
			fmt.Println("помилка читання повідомлення: ", err)
			break
		}
		fmt.Println("Отримано від: ", addrFrom, "Повідомлення: ", string(buffer[:msg]))
	}
}

package main

import (
	"fmt"
	"log"
	"net"
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
	defer connection.Close()

	fmt.Println("Слухаємо повідомлення на 9000")

	buffer := make([]byte, 1024)
	for {
		msg, addrFrom, err := connection.ReadFromUDP(buffer)
		if err != nil {
			fmt.Println("помилка читання повідомлення: ", err)
			continue
		}
		fmt.Println("Отримано від: ", addrFrom, "Повідомлення: ", string(buffer[:msg]))
	}
}

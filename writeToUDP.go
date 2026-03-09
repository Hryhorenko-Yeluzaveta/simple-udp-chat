package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
)

func main() {
	localAddr, err := net.ResolveUDPAddr("udp", ":9001")
	if err != nil {
		log.Fatal("Помилка прив'язки порту (bind): ", err)
	}
	connection, err := net.ListenUDP("udp", localAddr)
	if err != nil {
		log.Fatal("Помилка з'єднання: ", err)
	}
	defer connection.Close()

	fmt.Println("Відправляємо повідомлення на 9000")

	friendAddr, err := net.ResolveUDPAddr("udp", "192.168.0.100:9000")
	if err != nil {
		log.Fatal("Не вдалося визначити адресу іншого хоста: ", err)
	}

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		var message string
		message = scanner.Text()
		if err := scanner.Err(); err != nil {
			fmt.Println("помилка читання повідомлення: ", err)
			continue
		}

		_, err = connection.WriteToUDP([]byte(message), friendAddr)
		if err != nil {
			fmt.Print("Не вдалося відправити повідомлення", err)
			continue
		}
	}
}

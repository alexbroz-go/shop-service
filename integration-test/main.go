package main

import (
	"fmt"
	"os"
	"os/exec"
)

func main() {
	fmt.Println("Запуск интеграционных тестов...")
	cmd := exec.Command("go", "test", "-v", "./...") // запускает все тесты в папке
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	if err != nil {
		fmt.Println("Ошибка при выполнении тестов:", err)
		os.Exit(1)
	}
}

package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/kirillat6/go-basis/internal/task"
)


func main() {
	tasks := task.TaskManager{}
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Println("Выберите действие:")
		fmt.Println("1. Добавить Задачу")
		fmt.Println("2. Удалить Задачу")
		fmt.Println("3. Выполнить Задачу")
		fmt.Println("4. Просмотреть все задачи")
		fmt.Print("Введите ваш вариант: ")
		num, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Такого варианта нет!")
			continue
		}

		num = strings.TrimSpace(num)

		n, err := strconv.Atoi(num)
		if err != nil {
			fmt.Println("Введите число")
			continue
		}

		if n == 1 {
			fmt.Print("Введите название задачи: ")
			title, err := reader.ReadString('\n')
			if err != nil {
				fmt.Println("Считывание ответа прошло не успешно")
				continue
			}
			title = strings.TrimSpace(title)
			tasks.CreateTask(title)
		}
		if n == 2 {
			fmt.Print("Введите номер задачи: ")
			strId, err := reader.ReadString('\n')
			if err != nil {
				fmt.Println("Считывание ответа прошло не успешно")
				continue
			}
			strId = strings.TrimSpace(strId)
			id, err := strconv.Atoi(strId)
			if err != nil {
				fmt.Println("Преобразование id прошло не успешно")
				continue
			}
			err = tasks.DeleteTask(id)
			if err != nil {
				fmt.Println("Задача не найдена!")
				continue
			}
		}
		if n == 3 {
			fmt.Print("Введите номер задачи: ")
			strId, err := reader.ReadString('\n')

			if err != nil {
				fmt.Println("Считывание ответа прошло не успешно")
				continue
			}
			strId = strings.TrimSpace(strId)

			id, err := strconv.Atoi(strId)
			if err != nil {
				fmt.Println("Преобразование id прошло не успешно")
				continue
			}
			err = tasks.CompleteTask(id)
			if err != nil {
				fmt.Println("Задача не найдена!")
				continue
			}
		}
		if n == 4 {
			allTasks := tasks.GetTasks()
			for _,task := range allTasks {
				fmt.Println(task)
			}
		}
	}
}

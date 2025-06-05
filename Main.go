package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func ChooseClass() string {
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Выберите класс: Воин, Маг, Разбойник")
	classInput, _ := reader.ReadString('\n')
	classInput = strings.TrimSpace(classInput)
	return classInput
}

func main() {
	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Введите имя героя: ")
	name, _ := reader.ReadString('\n')
	name = strings.TrimSpace(name)

	class := ChooseClass()
	hero := ChooseHero(name, class)
	hero.ActiveQuest = append(hero.ActiveQuest, Quest{
		ID:          1,
		Description: "Убей 3 гоблинов",
		GoalType:    "kill",
		Target:      "Гоблин",
		Required:    3,
		RewardXP:    100,
	})

	enemies := []Enemy{
		{Name: "Гоблин", Health: 60},
		{Name: "Орк", Health: 80},
		{Name: "Скелет", Health: 50},
	}

	for {
		fmt.Println("\nВыберите действие:")
		fmt.Println("1 - Атаковать врага")
		fmt.Println("2 - Использовать спец. способность")
		fmt.Println("3 - Вылечиться")
		fmt.Println("4 - Показать статус")
		fmt.Println("5 - Показать инвентарь")
		fmt.Println("6 - Добавить предмет")
		fmt.Println("7 - Переместиться")
		fmt.Println("8 - Использовать зелье")
		fmt.Println("9 - Показать активные квесты")
		fmt.Println("10 - Выйти из игры")
		fmt.Print("Введите номер действия: ")
		choiceStr, _ := reader.ReadString('\n')
		choiceStr = strings.TrimSpace(choiceStr)
		choiceInt, err := strconv.Atoi(choiceStr)
		if err != nil || choiceInt < 1 || choiceInt > 10 {
			fmt.Println("Ошибка: введите число от 1 до 9.")
			continue
		}

		switch choiceInt {
		case 1:
			if len(enemies) == 0 {
				fmt.Println("Врагов больше нет!")
				continue
			}
			fmt.Println("Список врагов:")
			for i, enemy := range enemies {
				fmt.Printf("%d: %s (HP: %d)\n", i+1, enemy.Name, enemy.Health)
			}
			fmt.Print("Выберите врага по номеру: ")
			indexStr, _ := reader.ReadString('\n')
			indexStr = strings.TrimSpace(indexStr)
			index, err := strconv.Atoi(indexStr)
			if err != nil || index < 1 || index > len(enemies) {
				fmt.Println("Некорректный выбор врага.Введите числоот 1 до", len(enemies))
				continue
			}
			selected := &enemies[index-1]
			hero.Attack(selected, hero.AttackPower)
			if selected.Health <= 0 {
				fmt.Println(selected.Name, "побежден!")
				enemies = append(enemies[:index-1], enemies[index:]...)
				hero.gainXP(50)
			}
		case 2:
			if len(enemies) == 0 {
				fmt.Println("Нет врагов для применения способности.")
				continue
			}
			hero.UseSpecial(&enemies[0]) // применяет к первому врагу
		case 3:
			hero.heal(30)
		case 4:
			hero.status()
		case 5:
			hero.ShowInventory()
		case 6:
			fmt.Print("Введите имя предмета: ")
			itemName, _ := reader.ReadString('\n')
			itemName = strings.TrimSpace(itemName)

			fmt.Print("Введите тип предмета (Оружие / Броня / Зелье): ")
			itemType, _ := reader.ReadString('\n')
			itemType = strings.TrimSpace(itemType)
			itemType = strings.ToLower(itemType)
			switch itemType {
			case "оружие", "weapon":
				itemType = "Оружие"
			case "броня", "armor":
				itemType = "Броня"
			case "зелье", "potion":
				itemType = "Зелье"
			default:
				fmt.Println("Неизвестный тип предмета.")
				continue
			}

			fmt.Print("Введите значение предмета (число): ")
			valueStr, _ := reader.ReadString('\n')
			valueStr = strings.TrimSpace(valueStr)
			itemValue, err := strconv.Atoi(valueStr)
			if err != nil {
				fmt.Println("Ошибка ввода. Введите корректное число.")
				continue
			}

			newItem := Item{
				Name:  itemName,
				Type:  itemType,
				Value: itemValue,
			}

			hero.AddItem(newItem)
			SaveGame(&GameState{Hero: hero, Enemies: enemies})
		case 7:
			fmt.Print("Куда переместиться? ")
			location, _ := reader.ReadString('\n')
			location = strings.TrimSpace(location)
			hero.moveTo(location)
			SaveGame(&GameState{Hero: hero, Enemies: enemies})
		case 8:
			hero.UsePotion()
		case 9:
			if len(hero.ActiveQuest) == 0 {
				fmt.Println("Нет активных квестов.")
			} else {
				fmt.Println("🎯 Активные квесты:")
				for _, q := range hero.ActiveQuest {
					fmt.Printf("- %s (%s): %d/%d\n", q.Description, q.GoalType, q.Progress, q.Required)
				}
			}
		case 10:
			SaveGame(&GameState{Hero: hero, Enemies: enemies})
			fmt.Println("Выход из игры.")
			return
		default:
			fmt.Println("Неверный выбор. Повторите ввод.")
		}
	}
}

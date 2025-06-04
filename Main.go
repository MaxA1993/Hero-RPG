package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"strings"
)

type Hero struct {
	Name           string
	Health         int
	Location       string
	Inventory      []Item
	XP             int
	Level          int
	Class          string
	AttackPower    int
	EquippedWeapon *Item
	EquippedArmor  *Item
}

type Enemy struct {
	Name        string
	Health      int
	AttackPower int
}

type GameState struct {
	Hero    *Hero
	Enemies []Enemy
}

type Item struct {
	Name  string
	Type  string
	Value int
}

func ChooseClass() string {
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Выберите класс: Воин, Маг, Разбойник")
	classInput, _ := reader.ReadString('\n')
	classInput = strings.TrimSpace(classInput)
	return classInput
}

func ChooseHero(name string, class string) *Hero {
	hero := Hero{}
	hero.Name = name
	hero.Class = class
	if class == "Воин" {
		hero.Health = 150
		hero.AttackPower = 20
	} else if class == "Маг" {
		hero.Health = 80
		hero.AttackPower = 30
	} else if class == "Разбойник" {
		hero.Health = 100
		hero.AttackPower = 25
	}
	hero.Level = 1
	hero.XP = 0
	hero.Location = "Начальная локация"
	return &hero
}

func (h *Hero) UseSpecial(enemy *Enemy) {
	if h.Class == "Воин" {
		h.Health += 20
		fmt.Println("Воин активировал щит", h.Health)
	} else if h.Class == "Маг" {
		enemy.Health -= 40
		fmt.Println("Маг использует Фаерболл и наносит 40 урона", enemy.Name)
	} else if h.Class == "Разбойник" {
		fmt.Println("Разбойник уклонился от атаки ")
	}

}

func (h *Hero) gainXP(amount int) {
	h.XP += amount
	fmt.Printf("%s получил %d опыта\n", h.Name, amount)

	for h.XP >= 100 {
		h.Level++
		h.XP -= 100
		fmt.Printf("%s повысил уровень! Теперь уровень %d, здоровье: %d\n", h.Name, h.Level, h.Health)
	}
}

func (h *Hero) moveTo(location string) {
	h.Location = location
	fmt.Println(h.Name, "переместился в", location)
}

func (h *Hero) takeDamage(amount int) {
	if h.EquippedArmor != nil {
		amount -= h.EquippedArmor.Value
		if amount < 0 {
			amount = 0
		}
	}
	h.Health -= amount
	fmt.Println(h.Name, "получил повреждения на", amount, "HP")
}

func (h *Hero) heal(amount int) {
	h.Health += amount
	if h.Health > 100 {
		h.Health = 100
	}
	fmt.Println(h.Name, "восстановил здоровья на", amount, "HP")
}

func (h *Hero) Attack(enemy *Enemy, damage int) {
	if h.EquippedWeapon != nil {
		damage += h.EquippedWeapon.Value
		fmt.Println(h.Name, "атакует с помощью", h.EquippedWeapon.Name, "(+", h.EquippedWeapon.Value, "урона)")
	}
	enemy.Health -= damage
	fmt.Println(h.Name, "атакует", enemy.Name, "и наносит урона:", damage)
	if enemy.Health <= 0 {
		enemy.Health = 0
		fmt.Println(enemy.Name, "побежден")
	} else if enemy.Health != 0 {
		damage := rand.Intn(20) + 10
		fmt.Println(enemy.Name, "атакует", h.Name)
		h.takeDamage(damage)
		if h.Health == 0 {
			fmt.Println("Game over")
		}
	}

}

func (h *Hero) AddItem(item Item) {
	h.Inventory = append(h.Inventory, item)
	fmt.Println(h.Name, "получил предмет:", item)
}

func (h *Hero) ShowInventory() {
	if len(h.Inventory) == 0 {
		fmt.Println("Инвентарь пуст")
		return
	}
	fmt.Println("Инвентарь у", h.Name, ":")
	for _, item := range h.Inventory {
		fmt.Printf("- %s (%s, +%d)\n", item.Name, item.Type, item.Value)
	}
}

func (h *Hero) UsePotion() {
	for i, potion := range h.Inventory {
		if potion.Type == "Зелье" {
			h.Health += 30
			if h.Health > 100 {
				h.Health = 100
			}
			fmt.Println(h.Name, "использовал зелье. Здоровье восстановлено до", h.Health)
			h.Inventory = append(h.Inventory[:i], h.Inventory[i+1:]...)
			return
		}

	}
	fmt.Println("В инвентаре нет зелья.")
}
func (h *Hero) UseWeapon() {
	if h.EquippedWeapon != nil {
		fmt.Println("Оружие уже экипировано:", h.EquippedWeapon.Name)
		return
	}
	for i, item := range h.Inventory {
		if item.Type == "Оружие " {
			h.EquippedWeapon = &h.Inventory[i]
			h.AttackPower += item.Value
			fmt.Println(h.Name, "экипировал оружие:", item.Name, "Теперь сила атаки:", h.AttackPower)
			return
		}
	}
	fmt.Println("Нет оружия в инвентаре.")
}

func (h *Hero) UseArmor() {
	if h.EquippedArmor != nil {
		fmt.Println("Броня уже экипирована", h.EquippedArmor.Name)
		return
	}
	for i, item := range h.Inventory {
		if item.Type == "Броня" {
			h.EquippedArmor = &h.Inventory[i]
			fmt.Println(h.Name, "экипировал броню:", item.Name)
			return
		}
	}
	fmt.Println("Нет брони в инвентаре.")
}

func (h *Hero) status() {
	fmt.Printf("Имя:%s\nЗдоровье:%d\nЛокация:%s\n", h.Name, h.Health, h.Location)
}

func main() {
	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Введите имя героя: ")
	name, _ := reader.ReadString('\n')
	name = strings.TrimSpace(name)

	class := ChooseClass()
	hero := ChooseHero(name, class)

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
		fmt.Println("9 - Выйти из игры")
		fmt.Print("Введите номер действия: ")
		choiceStr, _ := reader.ReadString('\n')
		choiceStr = strings.TrimSpace(choiceStr)
		choiceInt, err := strconv.Atoi(choiceStr)
		if err != nil || choiceInt < 1 || choiceInt > 9 {
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
		case 7:
			fmt.Print("Куда переместиться? ")
			location, _ := reader.ReadString('\n')
			location = strings.TrimSpace(location)
			hero.moveTo(location)
		case 8:
			hero.UsePotion()
		case 9:
			fmt.Println("Выход из игры.")
			return
		default:
			fmt.Println("Неверный выбор. Повторите ввод.")
		}
	}
}

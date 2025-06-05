package main

import (
	"fmt"
	"math/rand"
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
	ActiveQuest    []Quest
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

func (h *Hero) CheckQuestProgress(eventType, target string) {
	for i := range h.ActiveQuest {
		quest := &h.ActiveQuest[i]
		if quest.Completed {
			continue
		}
		if quest.GoalType == eventType && quest.Target == target {
			quest.Progress++
			fmt.Printf("📜 Прогресс по квесту: %s (%d/%d)\n", quest.Description, quest.Progress, quest.Required)
			if quest.Progress >= quest.Required {
				quest.Completed = true
				fmt.Println("✅ Квест завершён:", quest.Description)
				h.gainXP(quest.RewardXP)
			}
		}
	}
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

	h.CheckQuestProgress("travel", location)
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

		h.CheckQuestProgress("kill", enemy.Name)
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

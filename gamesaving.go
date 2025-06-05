package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
)

func SaveGame(state *GameState) {
	data, err := json.MarshalIndent(state, "", "  ")
	if err != nil {
		log.Println("Ошибка при сохранении:", err)
		return
	}
	err = os.WriteFile("save.json", data, 0644)
	if err != nil {
		log.Println("Ошибка записи в файл:", err)
	} else {
		fmt.Println("💾 Игра сохранена.")
	}
}

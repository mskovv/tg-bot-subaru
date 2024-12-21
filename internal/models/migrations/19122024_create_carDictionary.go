package migrations

import (
	"github.com/mskovv/tg-bot-subaru96/internal/models"
	"gorm.io/gorm"
	"log"
)

func CreateCarDictionary(db *gorm.DB) {
	cars := []models.CarDictionary{
		{Mark: "Subaru", CarModel: "Impreza GF/GC"},
		{Mark: "Subaru", CarModel: "Impreza GG/GD"},
		{Mark: "Subaru", CarModel: "Impreza GH/GЕ"},
		{Mark: "Subaru", CarModel: "Impreza GJ"},
		{Mark: "Subaru", CarModel: "Impreza GR/GV"},
		{Mark: "Subaru", CarModel: "Impreza GP/GJ"},
		{Mark: "Subaru", CarModel: "Impreza GP(XV)"},

		{Mark: "Subaru", CarModel: "Forest SF"},
		{Mark: "Subaru", CarModel: "Forest SG"},
		{Mark: "Subaru", CarModel: "Forest SH"},
		{Mark: "Subaru", CarModel: "Forest SJ"},
		{Mark: "Subaru", CarModel: "Forest SK"},

		{Mark: "Subaru", CarModel: "Outback BG"},
		{Mark: "Subaru", CarModel: "Outback BH-BHE"},
		{Mark: "Subaru", CarModel: "Outback BP9-BPE"},
		{Mark: "Subaru", CarModel: "Outback BM9-BR9"},
		{Mark: "Subaru", CarModel: "Outback BS"},

		{Mark: "Toyota", CarModel: "Mark II GX80"},
		{Mark: "Toyota", CarModel: "Mark II JZX81"},
		{Mark: "Toyota", CarModel: "Mark II JZX90"},
		{Mark: "Toyota", CarModel: "Mark II JZX100"},
		{Mark: "Toyota", CarModel: "Mark II JZX110"},

		{Mark: "Toyota", CarModel: "Crown S130"},
		{Mark: "Toyota", CarModel: "Crown JZS141"},
		{Mark: "Toyota", CarModel: "Crown JZS151"},
		{Mark: "Toyota", CarModel: "Crown JZS171"},
		{Mark: "Toyota", CarModel: "Crown GRS182"},

		{Mark: "Toyota", CarModel: "Chaser JZX81"},
		{Mark: "Toyota", CarModel: "Chaser JZX90"},
		{Mark: "Toyota", CarModel: "Chaser JZX100"},

		{Mark: "Toyota", CarModel: "Altezza"},
	}

	var err error
	for _, car := range cars {
		err = db.FirstOrCreate(&car, models.CarDictionary{CarModel: car.CarModel}).Error
	}

	if err != nil {
		log.Fatalf("Не удалось добавить Словарь: %v", err)
	} else {
		log.Println("Словарь успешно загружен")
	}
}

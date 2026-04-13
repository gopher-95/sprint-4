package daysteps

import (
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/Yandex-Practicum/tracker/internal/spentcalories"
)

const (
	// Длина одного шага в метрах
	stepLength = 0.65
	// Количество метров в одном километре
	mInKm = 1000
)

func parsePackage(data string) (int, time.Duration, error) {

	//поделили строку на слайс строк
	stepsAndTime := strings.Split(data, ",")

	//проверили, что длина слайса равна 2
	if len(stepsAndTime) != 2 {
		return 0, 0, errors.New("the length of slice != 2")
	}

	//преобразовали первый элемент строки в тип int, обработали возможную ошибку и проверили, что число шагов неотрицательное
	steps, err := strconv.Atoi(stepsAndTime[0])

	if err != nil {
		return 0, 0, fmt.Errorf("произошла ошибка преобразования %w", err)
	}

	if steps <= 0 {
		return 0, 0, errors.New("steps are incorrect")
	}

	//преобразовали второй элемент строки в интервал времени, обработали возможную ошибку
	duration, err := time.ParseDuration(stepsAndTime[1])
	if err != nil {
		return 0, 0, fmt.Errorf("произошла ошибка преобразования %w", err)
	}

	if duration <= 0 {
		return 0, 0, errors.New("неверная продолжительность")
	}

	return steps, duration, nil
}

func DayActionInfo(data string, weight, height float64) string {

	//получаем количество шагов с помощью parsePackage
	steps, duration, err := parsePackage(data)
	if err != nil {
		log.Println(err)
		return ""
	}

	distance := float64(steps) * stepLength

	distanceInKilometers := distance / mInKm

	calories, err := spentcalories.WalkingSpentCalories(steps, weight, height, duration)
	if err != nil {
		log.Println(err)
		return ""
	}
	return fmt.Sprintf("Количество шагов: %d.\nДистанция составила %.2f км.\nВы сожгли %.2f ккал.\n", steps, distanceInKilometers, calories)

}

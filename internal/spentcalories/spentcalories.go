package spentcalories

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

// Основные константы, необходимые для расчетов.
const (
	lenStep                    = 0.65 // средняя длина шага.
	mInKm                      = 1000 // количество метров в километре.
	minInH                     = 60   // количество минут в часе.
	stepLengthCoefficient      = 0.45 // коэффициент для расчета длины шага на основе роста.
	walkingCaloriesCoefficient = 0.5  // коэффициент для расчета калорий при ходьбе
)

func parseTraining(data string) (int, string, time.Duration, error) {

	// сплитуем исходную строку и проверяем длину, если длина не равна 3 - создаем ошибку и ретёрним
	sliceInfo := strings.Split(data, ",")
	if len(sliceInfo) != 3 {
		return 0, "", 0, errors.New("длина слайса != 3")
	}

	//преобразуем строку в тип int
	steps, err := strconv.Atoi(sliceInfo[0])
	if err != nil {
		return 0, "", 0, fmt.Errorf("произошла ошибка преобразования: %w", err)
	}

	if steps <= 0 {
		return 0, "", 0, errors.New("неверное количество шагов")
	}

	//преобразуем строку в интервал времени
	duration, err := time.ParseDuration(sliceInfo[2])
	if err != nil {
		return 0, "", 0, fmt.Errorf("произошла ошибка преобразования: %w", err)
	}

	if duration <= 0 {
		return 0, "", 0, errors.New("неверная продолжительность")
	}

	activity := sliceInfo[1]

	return steps, activity, duration, nil
}

func distance(steps int, height float64) float64 {

	//длина шага
	stepLength := height * stepLengthCoefficient

	//дистанция
	distance := stepLength * float64(steps)

	//дистанция в километрах
	return distance / mInKm

}

func meanSpeed(steps int, height float64, duration time.Duration) float64 {

	//проверка, что дистанция строго больше нуля
	if duration <= 0 {
		return 0
	}

	//для подсчета дистанции используем предыдущую функцию
	distance := distance(steps, height)

	//средняя скорость
	averageSpeed := distance / duration.Hours()

	return averageSpeed
}

func TrainingInfo(data string, weight, height float64) (string, error) {
	steps, activity, duration, err := parseTraining(data)
	if err != nil {
		return "", err
	}

	switch activity {
	case "Бег":
		distanceRun := distance(steps, height)
		averageSpeedRun := meanSpeed(steps, height, duration)
		caloriesRun, err := RunningSpentCalories(steps, weight, height, duration)
		if err != nil {
			return "", nil
		}
		return fmt.Sprintf("Тип тренировки: %s\nДлительность: %.2f ч.\nДистанция: %.2f км.\nСкорость: %.2f км/ч\nСожгли калорий: %.2f\n", activity, duration.Hours(), distanceRun, averageSpeedRun, caloriesRun), nil
	case "Ходьба":
		distanceWalking := distance(steps, height)
		averageSpeedWalking := meanSpeed(steps, height, duration)
		caloriesWalking, err := WalkingSpentCalories(steps, weight, height, duration)
		if err != nil {
			return "", nil
		}
		return fmt.Sprintf("Тип тренировки: %s\nДлительность: %.2f ч.\nДистанция: %.2f км.\nСкорость: %.2f км/ч\nСожгли калорий: %.2f\n", activity, duration.Hours(), distanceWalking, averageSpeedWalking, caloriesWalking), nil
	}
	return "", errors.New("неизвестный тип тренировки")
}

func RunningSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {

	//проверка значений на отрицательно, либо равенство нулю
	if steps <= 0 || weight <= 0 || height <= 0 || duration <= 0 {
		return 0, errors.New("эти значения не могут быть меньше, либо равны нулю")
	}

	//расчет средней скорости, возпользовались предыдущей функцией
	averageSpeed := meanSpeed(steps, height, duration)

	//перевод интервала в минуты
	durationMin := duration.Minutes()

	//количество калорий
	calories := weight * averageSpeed * durationMin

	//количество калорий в час
	caloriesPerHour := calories / minInH

	return caloriesPerHour, nil
}

func WalkingSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	if steps <= 0 || weight <= 0 || height <= 0 || duration <= 0 {
		return 0, errors.New("эти значения не могут быть меньше, либо равны нулю")
	}

	//расчет средней скорости, возпользовались предыдущей функцией
	averageSpeed := meanSpeed(steps, height, duration)

	//перевод интервала в минуты
	durationMin := duration.Minutes()

	//количество калорий
	calories := weight * averageSpeed * durationMin

	//количество калорий в час
	caloriesPerHour := calories / minInH

	walkingCalories := caloriesPerHour * walkingCaloriesCoefficient
	return walkingCalories, nil
}

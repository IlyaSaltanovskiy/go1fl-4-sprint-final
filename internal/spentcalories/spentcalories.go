package spentcalories

import (
	"errors"
	"fmt"
	"log"
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
	parts := strings.Split(data, ",")

	if len(parts) != 3 {
		return 0, "", 0, fmt.Errorf("неверный формат данных: %s", data)
	}

	steps, err := strconv.Atoi(parts[0])
	if err != nil {
		return 0, "", 0, err
	}

	if steps <= 0 {
		return 0, "", 0, fmt.Errorf("количество шагов должно быть больше 0")
	}

	activity := parts[1]

	duration, err := time.ParseDuration(parts[2])
	if err != nil {
		return 0, "", 0, err
	}

	if duration <= 0 {
		return 0, "", 0, fmt.Errorf("продолжительность должна быть больше 0")
	}

	return steps, activity, duration, nil
}

func distance(steps int, height float64) float64 {
	stepLength := height * stepLengthCoefficient
	distanceM := float64(steps) * stepLength
	return distanceM / mInKm

}

func meanSpeed(steps int, height float64, duration time.Duration) float64 {
	if duration <= 0 {
		return 0
	}

	dist := distance(steps, height)
	return dist / duration.Hours()
}
func TrainingInfo(data string, weight, height float64) (string, error) {
	steps, activity, duration, err := parseTraining(data)
	if err != nil {
		log.Println(err)
		return "", err
	}

	switch activity {
	case "Бег":
		calories, err := RunningSpentCalories(steps, weight, height, duration)
		if err != nil {
			log.Println(err)
			return "", err
		}
		return fmt.Sprintf(
			"Тип тренировки: %s\nДлительность: %.2f ч.\nДистанция: %.2f км.\nСкорость: %.2f км/ч\nСожгли калорий: %.2f\n",
			activity, duration.Hours(), distance(steps, height), meanSpeed(steps, height, duration), calories,
		), nil

	case "Ходьба":
		calories, err := WalkingSpentCalories(steps, weight, height, duration)
		if err != nil {
			log.Println(err)
			return "", err
		}
		return fmt.Sprintf(
			"Тип тренировки: %s\nДлительность: %.2f ч.\nДистанция: %.2f км.\nСкорость: %.2f км/ч\nСожгли калорий: %.2f\n",
			activity, duration.Hours(), distance(steps, height), meanSpeed(steps, height, duration), calories,
		), nil

	default:
		return "", errors.New("неизвестный тип тренировки")
	}
}

func RunningSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	if steps <= 0 {
		return 0, fmt.Errorf("количество шагов должно быть больше 0")
	}
	if weight <= 0 {
		return 0, fmt.Errorf("вес должен быть больше 0")
	}
	if height <= 0 {
		return 0, fmt.Errorf("рост должен быть больше 0")
	}
	if duration <= 0 {
		return 0, fmt.Errorf("продолжительность должна быть больше 0")
	}

	speed := meanSpeed(steps, height, duration)
	durationInMinutes := duration.Minutes()

	calories := (weight * speed * durationInMinutes) / minInH
	return calories, nil
}

func WalkingSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	if steps <= 0 {
		return 0, fmt.Errorf("количество шагов должно быть больше 0")
	}
	if weight <= 0 {
		return 0, fmt.Errorf("вес должен быть больше 0")
	}
	if height <= 0 {
		return 0, fmt.Errorf("рост должен быть больше 0")
	}
	if duration <= 0 {
		return 0, fmt.Errorf("продолжительность должна быть больше 0")
	}

	speed := meanSpeed(steps, height, duration)
	durationInMinutes := duration.Minutes()

	calories := (weight * speed * durationInMinutes) / minInH
	calories *= walkingCaloriesCoefficient

	return calories, nil
}

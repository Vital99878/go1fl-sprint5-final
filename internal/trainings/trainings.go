package trainings

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/Yandex-Practicum/tracker/internal/personaldata"
	"github.com/Yandex-Practicum/tracker/internal/spentenergy"
)

type Training struct {
	Steps        int
	TrainingType string
	Duration     time.Duration
	personaldata.Personal
}

func (t *Training) Parse(datastring string) (err error) {
	parts := strings.Split(datastring, ",")
	if len(parts) != 3 {
		return errors.New("ожидаются три значения: шаги, тип тренировки, длительность")
	}

	stepCount, err := strconv.Atoi(parts[0])
	if err != nil {
		return err
	}
	if stepCount <= 0 {
		return errors.New("количество шагов должно быть больше нуля")
	}

	trainingType := parts[1]
	if trainingType == "" {
		return errors.New("должен быть указан тип тренировки")
	}

	trainingDuration, err := time.ParseDuration(parts[2])
	if err != nil {
		return err
	}
	if trainingDuration <= 0 {
		return errors.New("длительность тренировки должна быть больше нуля")
	}

	t.Steps = stepCount
	t.TrainingType = trainingType
	t.Duration = trainingDuration

	return nil
}

func (t Training) ActionInfo() (string, error) {
	distance := spentenergy.Distance(t.Steps, t.Height)
	speed := spentenergy.MeanSpeed(t.Steps, t.Height, t.Duration)

	var calories float64
	var err error

	switch t.TrainingType {
	case "Бег":
		calories, err = spentenergy.RunningSpentCalories(t.Steps, t.Weight, t.Height, t.Duration)
	case "Ходьба":
		calories, err = spentenergy.WalkingSpentCalories(t.Steps, t.Weight, t.Height, t.Duration)
	default:
		return "", errors.New("неизвестный тип тренировки")
	}

	if err != nil {
		return "", err
	}

	message := fmt.Sprintf(
		"Тип тренировки: %s\n"+
			"Длительность: %.2f ч.\n"+
			"Дистанция: %.2f км.\n"+
			"Скорость: %.2f км/ч\n"+
			"Сожгли калорий: %.2f\n",
		t.TrainingType,
		t.Duration.Hours(),
		distance,
		speed,
		calories,
	)
	return message, nil
}

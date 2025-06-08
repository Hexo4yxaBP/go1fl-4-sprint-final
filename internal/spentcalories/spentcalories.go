// Package spentcalories implements functions for calculating calories based on training and user parameters.
// Also contains functions collecting training report.
package spentcalories

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

// Basic constants for calculations.
const (
	lenStep                    = 0.65 // average step length.
	mInKm                      = 1000 // number of meters in a kilometer.
	minInH                     = 60   // number of minutes in an hour.
	stepLengthCoefficient      = 0.45 // coefficient for calculating step length based on height.
	walkingCaloriesCoefficient = 0.5  // coefficient for calculating calories from walking
)

// Error vars definition.

var (
	ErrUnknownTrainType = errors.New("неизвестный тип тренировки") // Unknown Training Type error
	ErrIncorrectParams  = errors.New("incorrect input parameters") // Incorrect input parameters error
	ErrWrongData        = errors.New("wrong data format")          // Wrong data format error
)

// parseTraining parces training data string and returns steps count, activity type and training duration.
func parseTraining(data string) (int, string, time.Duration, error) {

	splittedData := strings.Split(data, ",")

	if len(splittedData) != 3 {
		return 0, "", 0, ErrWrongData
	}

	steps, err := strconv.Atoi(splittedData[0])

	if err != nil || steps <= 0 {
		return 0, "", 0, errors.Join(ErrWrongData, err)
	}

	activity := splittedData[1]

	if len(activity) == 0 {
		return 0, "", 0, ErrUnknownTrainType
	}

	duration, err := time.ParseDuration(splittedData[2])

	if err != nil || duration <= 0 {
		return 0, "", 0, errors.Join(ErrWrongData, err)
	}

	return steps, activity, duration, nil

}

// distance calculates the distance based on step count and user height.
func distance(steps int, height float64) float64 {

	return height * stepLengthCoefficient * float64(steps) / mInKm
}

// meanSpeed calculates the meanSpeed based on step count, user height and activity duration.
func meanSpeed(steps int, height float64, duration time.Duration) float64 {

	if duration <= 0 {
		return 0
	}

	return distance(steps, height) / duration.Hours()
}

// TrainingInfo function make a training report and returns it in string format.
func TrainingInfo(data string, weight, height float64) (string, error) {

	steps, activityType, trainingDuration, err := parseTraining(data)

	if err != nil {
		return "", err
	}

	var spentCalories float64

	switch activityType {
	case "Бег":
		spentCalories, err = RunningSpentCalories(steps, weight, height, trainingDuration)
	case "Ходьба":
		spentCalories, err = WalkingSpentCalories(steps, weight, height, trainingDuration)
	default:
		return "", ErrUnknownTrainType
	}

	if err != nil {
		return "", err
	}

	trainingResult := fmt.Sprintf("Тип тренировки: %s\nДлительность: %.2f ч.\n", activityType, trainingDuration.Hours())
	trainingResult += fmt.Sprintf("Дистанция: %.2f км.\n", distance(steps, height))
	trainingResult += fmt.Sprintf("Скорость: %.2f км/ч\n", meanSpeed(steps, height, trainingDuration))
	trainingResult += fmt.Sprintf("Сожгли калорий: %.2f\n", spentCalories)

	return trainingResult, nil

}

// RunningSpentCalories calculates spent calories for running activity based on steps count, user weight and height and activity duration.
// Returns 0 and "Incorrect input parameters" error if any parameter == 0.
func RunningSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {

	if steps <= 0 || weight <= 0 || height <= 0 || duration <= 0 {
		return 0, ErrIncorrectParams
	}

	return weight * meanSpeed(steps, height, duration) * duration.Minutes() / minInH, nil

}

// WalkingSpentCalories calculates spent calories for walking activity based on steps count, user weight and height and activity duration.
// Returns 0 and "Incorrect input parameters" error if any parameter == 0.
func WalkingSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {

	if steps <= 0 || weight <= 0 || height <= 0 || duration <= 0 {
		return 0, ErrIncorrectParams
	}

	return weight * meanSpeed(steps, height, duration) * duration.Minutes() / minInH * walkingCaloriesCoefficient, nil

}

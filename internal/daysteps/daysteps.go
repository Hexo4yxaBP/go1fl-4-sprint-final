// Package daysteps implements calculating information about user day activity.
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
	stepLength = 0.65 // Length of one step in meters

	mInKm = 1000 // Meters in one kilometer
)

// Error vars definition.
var (
	ErrWrongData = errors.New("wrong data") // Wrong data format error
)

// parcePackage converts string containing steps and time to int and time.Duration formats.
// Returns int steps, time.Duration Duration or error.
func parsePackage(data string) (int, time.Duration, error) {

	// split string

	splittedData := strings.Split(data, ",")
	if len(splittedData) != 2 {
		return 0, 0, ErrWrongData
	}

	// convert data

	steps, err := strconv.Atoi(splittedData[0])

	if err != nil || steps <= 0 {
		return 0, 0, errors.Join(ErrWrongData, err)
	}

	walkDuration, err := time.ParseDuration(splittedData[1])

	if err != nil || walkDuration <= 0 {
		return 0, 0, errors.Join(ErrWrongData, err)
	}

	return steps, walkDuration, nil
}

// DayActionInfo calculate and return information in text format about day activity of user.
// Output string countains information about steps count, passed distance in km and burned kCals.
func DayActionInfo(data string, weight, height float64) string {

	steps, walkDuration, err := parsePackage(data)

	if err != nil {
		log.Print(err.Error())
		return ""
	}

	if steps <= 0 {
		return ""
	}

	var distanseKm float64 = float64(steps) * stepLength / mInKm

	kCal, err := spentcalories.WalkingSpentCalories(steps, weight, height, walkDuration)

	if err != nil {
		fmt.Println(err)
		return ""
	}

	return fmt.Sprintf("Количество шагов: %d.\nДистанция составила %.2f км.\nВы сожгли %.2f ккал.\n", steps, distanseKm, kCal)

}

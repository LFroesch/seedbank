package generator

import "time"

var fixedReferenceDate = time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)

func referenceDate() time.Time {
	return fixedReferenceDate
}

func ageAt(dob, ref time.Time) int {
	age := ref.Year() - dob.Year()
	if ref.YearDay() < dob.YearDay() {
		age--
	}
	return age
}

package main

import (
	"fmt"
)

type SessionInterface interface {
	validateMvs() *SessionEntity
	validateStatus() *SessionEntity
	getTotalCost() float64
}

type SessionEntity struct {
	meterValueArray  []float64
	statusArray      []string
	totalCost        float64
	finished         bool
	needManualReview bool
	errorMessage     string
}

func (s *SessionEntity) validateMvs() *SessionEntity {
	if s.errorMessage != "" {
		panic(s.errorMessage)
	}

	if len(s.meterValueArray) == 0 {
		s.setErrorMessage("Invalide meter values")
		return s
	}

	currentValue := s.meterValueArray[0]

	for _, v := range s.meterValueArray {
		if v < currentValue {
			s.setErrorMessage("Invalide meter values")
			return s
		}
		currentValue = v
	}

	return s
}

func (s *SessionEntity) validateStatus() *SessionEntity {
	if s.errorMessage != "" {
		panic(s.errorMessage)
	}

	for _, v := range s.statusArray {
		if v == "FAULTED" {
			s.setErrorMessage("Invalide status array")
			return s
		}
	}
	return s
}

func (s *SessionEntity) getTotalCost() (res float64) {
	if !s.finished {
		res = -1
		return
	}

	if s.finished && len(s.meterValueArray) == 0 {
		res = -1
		return
	}

	res = s.meterValueArray[len(s.meterValueArray)-1] - s.meterValueArray[0]

	return
}

func (s *SessionEntity) setErrorMessage(err string) {
	s.errorMessage = err
}

func main() {
	mvs := make([]float64, 10)
	mvs = append(mvs, 0, 1, 2, 3, 4, 5, 6, 7, 8)

	status := make([]string, 10)
	status = append(status, "charing", "finishing")

	session := SessionEntity{totalCost: -1, finished: true, needManualReview: false, meterValueArray: mvs, statusArray: status}

	var sessionInterface SessionInterface = &session

	sessionInterface.validateMvs().validateStatus()

	totalCost := sessionInterface.getTotalCost()

	if totalCost > 0 {
		fmt.Printf("We need to bill this customer %f\n", totalCost)
	}

}

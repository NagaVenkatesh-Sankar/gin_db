package main

import (
	"errors"
	"fmt"
)

//AndroidDevice section
type AndroidDevice struct {
}

//DeviceInfo details
func (device *AndroidDevice) DeviceInfo() string {
	return fmt.Sprintf("Android device")
}

// Lock function
func (device *AndroidDevice) Lock() string {
	return fmt.Sprintf("Android device locked %T", device)
}

// Unlock function
func (device *AndroidDevice) Unlock() string {
	return fmt.Sprintf("Android device Unlock")
}

var n int = 0

// SendSms function
func (device *AndroidDevice) SendSms(context smsSchema) (string, error) {
	n++
	// mock, if the SIM is not available, send error
	if n%2 == 0 {
		return "Android SMS failed", errors.New(fmt.Sprint("SIM Card not available"))
	}
	fmt.Println("Android SMS", context)

	//on successful sms sent
	return "Android SMS sent", nil
}

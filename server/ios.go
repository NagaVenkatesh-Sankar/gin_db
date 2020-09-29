package main

import (
	"errors"
	"fmt"
)

//IosDevice section
type IosDevice struct {
}

//DeviceInfo details
func (device *IosDevice) DeviceInfo() string {
	return fmt.Sprintf("iOS device")
}

// Lock function
func (device *IosDevice) Lock() string {
	return fmt.Sprintf("iOS device locked")
}

// Unlock function
func (device *IosDevice) Unlock() string {
	return fmt.Sprintf("iOS device Unlock")
}

// SendSms function
func (device *IosDevice) SendSms(context smsSchema) (string, error) {
	n++
	// mock, if the SIM is not available, send error
	if n%2 == 0 {
		return "iOS SMS failed", errors.New(fmt.Sprint("SIM Card not available"))
	}
	fmt.Println("iOS SMS", context)

	//on successful sms sent
	return "iOS SMS sent", nil
}

package main

import (
	"errors"
	"fmt"
)

// device types supported
const (
	Android = "android"
	Ios     = "ios"
)

//Device definitions
type Device interface {
	DeviceInfo() string
	Lock() string
	Unlock() string
	SendSms(smsSchema) (string, error)
}

//CreateDevice based on the device type definied
func CreateDevice(deviceType string) (Device, error) {
	switch deviceType {
	case Android:
		return new(AndroidDevice), nil
	case Ios:
		return new(IosDevice), nil
	default:
		return nil, errors.New(fmt.Sprint("Invalid Device type specified"))
	}
}

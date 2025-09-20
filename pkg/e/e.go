package e

import "fmt"

var (
	ErrStationNotFound = fmt.Errorf("station not found")
	ErrStationExists   = fmt.Errorf("station is exists")

	ErrStoreNotFound = fmt.Errorf("store not found")
	ErrStoreExists   = fmt.Errorf("store exists")

	ErrLocationNotFound = fmt.Errorf("location not found")
	ErrLocationExists   = fmt.Errorf("location is exists")
)

func Wrap(msg string, err error) error {
	return fmt.Errorf("%s: %w", msg, err)
}

func WrapWithFunc(funcName, msg string, err error) error {
	return Wrap(fmt.Sprintf("[%s] %s", funcName, msg), err)
}

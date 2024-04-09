package utils

import "fmt"

func HandleError(err error) error {
	fmt.Println(err)
	return err
}

package utils

import "encoding/json"

func Unmarshal[T any](data []byte) (T, error) {

	res := new(T)

	err := json.Unmarshal(data, res)
	if err != nil {
		return *res, err
	}

	return *res, nil
}

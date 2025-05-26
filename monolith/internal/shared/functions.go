package shared

import (
	"encoding/json"
)

func ParseParams(params map[string]interface{}, target interface{}) error {
	data, err := json.Marshal(params)
	if err != nil {
		return err
	}

	err = json.Unmarshal(data, target)
	if err != nil {
		return err
	}
	return nil
}

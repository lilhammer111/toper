package util

import (
	"errors"
	"go.uber.org/zap"
	"strconv"
)

func StrConvertUint(id string) (uint64, error) {
	parsedID, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		zap.S().Error("error parsing string to uint: ", err)
		return 0, err
	}
	// If running on a 32-bit system, make sure that parsedID does not exceed the maximum value of uint32
	if parsedID > uint64(^uint(0)) {
		zap.S().Error("parsed ID exceeds the maximum value for uint")
		return 0, errors.New("parsed ID exceeds the maximum value for uint")
	}

	return parsedID, nil
}

package scheduler

import (
	"context"
	"go.uber.org/zap"
	"to-persist/server/global"
)

const (
	SetUndone = "SetToperIDUndone"
)

var (
	TaskFunctionsMap = map[string]func(toperID string) func(){
		SetUndone: SetToperIDUndone,
	}
)

func SetToperIDUndone(toperIDStr string) func() {
	return func() {
		value, err := global.RedisClient.Get(context.Background(), toperIDStr).Result()
		if err != nil {
			zap.S().Errorf("failed to get %s : %s", toperIDStr, err)
			//todo  retry
			return
		}

		if value == "done" {
			_, err := global.RedisClient.Set(context.Background(), toperIDStr, "undone", 0).Result()
			if err != nil {
				zap.S().Errorf("failed to set %s: %s", toperIDStr, err)
				// todo retry
			}
		}
	}
}

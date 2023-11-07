package scheduler

import (
	"context"
	"go.uber.org/zap"
	"to-persist/server/global"
	"to-persist/server/model"
	"to-persist/server/util"
)

const (
	CheckDoneStatus = "ResetAndSubmitToperDoneStatus"
)

var (
	TaskFunctionsMap = map[string]func(toperID string) func(){
		CheckDoneStatus: ResetAndSubmitToperDoneStatus,
	}
)

func ResetAndSubmitToperDoneStatus(toperID string) func() {
	return func() {
		value, err := global.RedisClient.Get(context.Background(), toperID).Result()
		if err != nil {
			zap.S().Errorf("failed to get %s : %s", toperID, err)
			//todo  retry
			return
		}

		if value == global.ToperStatusUndone {
			var doneHistory model.DoneHistory
			doneHistory.Done = global.ToperStatusUndone

			id, err := util.StrConvertUint(toperID)
			if err != nil {
				return
			}
			doneHistory.ToperID = uint(id)

			if res := global.MysqlDB.Create(&doneHistory); res.RowsAffected == 0 {
				zap.S().Errorf("failed to create toper history for toper id %s : %v", toperID, res.Error)
				return
			}

		}
		_, err = global.RedisClient.Set(context.Background(), toperID, "undone", 0).Result()
		if err != nil {
			zap.S().Errorf("failed to set %s: %s", toperID, err)
			// todo retry
		}

	}
}

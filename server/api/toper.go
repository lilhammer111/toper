package api

import (
	"context"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
	"strconv"
	"strings"
	"to-persist/server/global"
	"to-persist/server/model"
	"to-persist/server/util/scheduler"
)

type ListResp struct {
	ID      int    `json:"id,omitempty"`
	Acronym string `json:"acronym,omitempty"`
	Desc    string `json:"desc,omitempty"`
	DueDate string `json:"due-date,omitempty"`
	Period  string `json:"period,omitempty"`
	Done    string `json:"done,omitempty"`
}

func Create(c *gin.Context) {
	userID, exists := c.Get("user-id")
	if !exists {
		c.Status(http.StatusUnauthorized)
		zap.S().Error("failed to get user id in gin's context ")
		return
	}
	toper := model.Toper{}
	err := c.ShouldBindJSON(&toper)
	if err != nil {
		c.Status(http.StatusBadRequest)
		zap.S().Error("bad request while create toper: ", err)
		return
	}

	if id, ok := userID.(string); ok {
		parsedID, err := strconv.ParseUint(id, 10, 64)
		if err != nil {
			zap.S().Error("error parsing string to uint: ", err)
			c.Status(http.StatusBadRequest)
			return
		}
		// If running on a 32-bit system, make sure that parsedID does not exceed the maximum value of uint32
		if parsedID > uint64(^uint(0)) {
			zap.S().Error("parsed ID exceeds the maximum value for uint")
			c.Status(http.StatusBadRequest)
			return
		}

		toper.UserID = uint(parsedID)
	}

	tx := global.MysqlDB.Begin()
	if res := tx.Create(&toper); res.Error != nil {
		zap.S().Errorf("failed to create toper in mysql: %s", err)
		c.Status(http.StatusInternalServerError)
		return
	}

	expr, err := parseForExpr(toper.Period, toper.DueDate)
	if err != nil {
		tx.Rollback()
		zap.S().Errorw("failed to parse flags while creating toper", toper.Period)
		c.Status(http.StatusInternalServerError)
		return
	}

	toperIDStr := strconv.FormatUint(uint64(toper.ID), 10)
	_, err = global.RedisClient.Set(context.Background(), toperIDStr, "undone", 0).Result()
	if err != nil {
		tx.Rollback()
		zap.S().Errorf("failed to set the key of toper ID: %s", err)
		c.Status(http.StatusInternalServerError)
		return
	}

	taskScheduler := scheduler.NewTaskScheduler()
	err = taskScheduler.AddTask(toperIDStr, expr, scheduler.SetUndone)
	if err != nil {
		tx.Rollback()
		zap.S().Errorf("failed to add task for %d : %s", toper.ID, err)
		c.Status(http.StatusInternalServerError)
		return
	}

	tx.Commit()

	c.Status(http.StatusOK)

}

func List(c *gin.Context) {
	userID, exists := c.Get("user-id")
	if !exists {
		c.Status(http.StatusUnauthorized)
		zap.S().Error("failed to get user id in gin's context ")
		return
	}

	topers := make([]model.Toper, 0)
	res := global.MysqlDB.Where("user_id = ?", userID).Find(&topers)
	if res.Error != nil {
		c.Status(http.StatusInternalServerError)
		zap.S().Errorf("failed to retrieve topers for user %v: %v", userID, res.Error)
		return
	}
	if res.RowsAffected == 0 {
		c.Status(http.StatusNotFound)
		zap.S().Errorf("there is no toper for user %v ", userID)
		return
	}

	listResponses := make([]ListResp, 0)
	for i, toper := range topers {
		listResp := ListResp{}

		toperID := strconv.FormatUint(uint64(toper.ID), 10)
		done, err := global.RedisClient.Get(context.Background(), toperID).Result()
		if err != nil {
			zap.S().Errorf("failed to get toper id while response list api: %s", err)
			c.Status(http.StatusInternalServerError)
			return
		}

		if done == "done" {
			listResp.Done = "√"
		} else if done == "undone" {
			listResp.Done = "×"
		} else {
			zap.S().Error("the key of 'done' is neither 'done' nor 'undone'")
			c.Status(http.StatusInternalServerError)
			return
		}
		listResp.ID = i + 1
		listResp.Acronym = toper.Acronym
		listResp.Desc = toper.Description
		listResp.DueDate = toper.DueDate
		listResp.Period = toper.Period

		listResponses = append(listResponses, listResp)
	}

	c.JSON(http.StatusOK, listResponses)
}

func Done(c *gin.Context) {

}

func History(c *gin.Context) {

}

func Alter(c *gin.Context) {

}

func parseForExpr(period string, dueDate string) (string, error) {

	zap.S().Infof("period's vlaue: %s", period)

	hourAndMinute := strings.Split(dueDate, ":")
	if len(hourAndMinute) != 2 {
		return "", errors.New("wrong due date format")
	}

	var expr string
	if strings.ToLower(period) == "everyday" {
		expr = fmt.Sprintf("%s %s * * *", hourAndMinute[1], hourAndMinute[0])
	} else {
		expr = fmt.Sprintf("%s %s * * %s", hourAndMinute[1], hourAndMinute[0], period)
	}

	return expr, nil
}

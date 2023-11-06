package scheduler

import (
	"errors"
	"github.com/robfig/cron/v3"
	"go.uber.org/zap"
	"sync"
	"to-persist/server/global"
	"to-persist/server/model"
)

type TaskScheduler struct {
	cr    *cron.Cron
	tasks map[string]cron.EntryID
	mu    sync.Mutex
}

var (
	instance *TaskScheduler
	once     sync.Once
)

func NewTaskScheduler() *TaskScheduler {
	once.Do(func() {
		instance = &TaskScheduler{
			cr:    cron.New(),
			tasks: make(map[string]cron.EntryID),
		}
	})
	return instance
}

func (ts *TaskScheduler) Start() {
	ts.cr.Start()
}

func (ts *TaskScheduler) Stop() {
	ts.cr.Stop()
}

func (ts *TaskScheduler) AddTask(toperID, expr string, taskType string) error {
	if _, ok := TaskFunctionsMap[taskType]; !ok {
		return errors.New("failed to match the predefined task function")
	}

	task := model.Task{
		ToperID:      toperID,
		Expression:   expr,
		TaskFuncType: taskType,
	}

	tx := global.MysqlDB.Begin()
	if res := global.MysqlDB.Create(&task); res.RowsAffected == 0 {
		zap.S().Errorf("failed to create the record of task: %v", res.Error)
		return errors.New("failed to create the record of task")
	}

	entryID, err := ts.cr.AddFunc(expr, TaskFunctionsMap[taskType](toperID))
	if err != nil {
		zap.S().Errorf("failed to add func: %v", err)
		tx.Rollback()
		return err
	}

	ts.mu.Lock()
	defer ts.mu.Unlock()
	ts.tasks[toperID] = entryID

	tx.Commit()
	return nil
}

func (ts *TaskScheduler) RemoveTask(toperID string) error {
	task := model.Task{ToperID: toperID}
	res := global.MysqlDB.Where("toper_id = ?", toperID).First(&task)
	if res.RowsAffected == 0 {
		zap.S().Errorf("no record with toper_id %s: %v\n", toperID, res.Error)
		return res.Error
	}

	tx := global.MysqlDB.Begin()
	res = global.MysqlDB.Delete(&task)
	if res.RowsAffected == 0 {
		zap.S().Errorf("failed to delete the record of task: %v", res.Error)
		return res.Error
	}

	ts.mu.Lock()
	defer ts.mu.Unlock()
	if entryID, ok := ts.tasks[toperID]; ok {
		ts.cr.Remove(entryID)
		delete(ts.tasks, toperID)
	}
	tx.Commit()
	return nil
}

func (ts *TaskScheduler) ReInitTasksFromDB() error {
	var total int64
	err := global.MysqlDB.Model(&model.Task{}).Count(&total).Error
	if err != nil {
		zap.S().Errorf("failed to count the number of tasks: %v", err)
		return err
	}

	var wg sync.WaitGroup
	wg.Add(1)
	ch := make(chan []model.Task)
	go func() {
		defer wg.Done()
		for tasks := range ch {
			for _, task := range tasks {
				entryID, err := ts.cr.AddFunc(task.Expression, TaskFunctionsMap[task.TaskFuncType](task.ToperID))
				if err != nil {
					// todo : record the failure operation and retry later instead of breaking for loop
					zap.S().Errorf("failed to init task %s: %v", task.ToperID, err)
					continue
				}

				ts.mu.Lock()
				ts.tasks[task.ToperID] = entryID
				ts.mu.Unlock()
			}
		}
	}()

	batchSize := 100
	for offset := 0; offset < int(total); offset += batchSize {
		var tasks []model.Task
		err = global.MysqlDB.Limit(batchSize).Offset(offset).Find(&tasks).Error
		if err != nil {
			zap.S().Errorf("failed to init tasks: %v", err)
			return err
		}

		ch <- tasks
	}
	close(ch)
	wg.Wait()
	return nil
}

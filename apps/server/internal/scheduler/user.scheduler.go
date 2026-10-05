package scheduler

import (
	"encoding/json"
	"log"
	"time"

	"github.com/giridhar-m-a/custom_form_app/constants"
	"github.com/giridhar-m-a/custom_form_app/internal/dto/scheduler_dto"
	"github.com/hibiken/asynq"
)

func ScheduleDeleteUser(userID string, scheduleTime time.Time) (*asynq.TaskInfo, error) {
	client := NewClient()
	defer client.Close()
	deleteUserPayload := scheduler_dto.UserSchedulerPayload{
		UserID: userID,
	}
	payload, err := json.Marshal(deleteUserPayload)
	if err != nil {
		log.Printf("[Delete User Scheduler] failed to parse json, %s", err.Error())
	}
	task := asynq.NewTask(constants.TaskTypeDeleteUser, payload)

	info, err := client.Enqueue(task, asynq.ProcessAt(scheduleTime), asynq.Queue(constants.QueueDeleteUser), asynq.MaxRetry(0))
	if err != nil {
		log.Printf("Error scheduling delete user for %s: %v", userID, err)
		return nil, err
	}
	log.Printf("Delete user scheduled for %s at %v, Task ID: %s", userID, scheduleTime, info.ID)
	return info, nil

}


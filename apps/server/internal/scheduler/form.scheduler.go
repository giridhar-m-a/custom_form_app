package scheduler

import (
	"encoding/json"
	"log"
	"time"

	"github.com/giridhar-m-a/custom_form_app/constants"
	"github.com/giridhar-m-a/custom_form_app/internal/dto/scheduler_dto"
	"github.com/hibiken/asynq"
)

func FormStatusUpdateSchedule(formID string, scheduleTime time.Time) (*asynq.TaskInfo, error) {
	client := NewClient()
	defer client.Close()
	formPayload := scheduler_dto.InvitationSchedulerPayload{
		FormID: formID,
	}
	payload, _ := json.Marshal(formPayload)
	task := asynq.NewTask(constants.TaskTypeFormStatusUpdate, payload)

	info, err := client.Enqueue(task, asynq.ProcessAt(scheduleTime))
	if err != nil {
		log.Printf("Error scheduling form for form %s: %v", formID, err)
		return nil, err
	}
	log.Printf("Form scheduled for form %s at %v, Task ID: %s", formID, scheduleTime, info.ID)
	return info, nil

}

func CancelFormStatusUpdateSchedule(scheduleId string) error {
	inspector := asynq.NewInspector(NewRedisClientOpt())
	err := inspector.DeleteTask(constants.TaskTypeFormStatusUpdate, scheduleId)
	if err != nil {
		return err
	}
	log.Printf("Form cancelled for schedule ID: %s", scheduleId)
	return nil
}

func UpdateFormStatusUpdateSchedule(scheduleID string, scheduleTime time.Time, formID string) (*asynq.TaskInfo, error) {
	err := CancelFormStatusUpdateSchedule(scheduleID)
	if err != nil {
		return nil, err
	}
	info, err := FormStatusUpdateSchedule(formID, scheduleTime)
	if err != nil {
		return nil, err
	}
	log.Printf("Form updated for schedule ID: %s, New Task ID: %s", scheduleID, info.ID)
	return info, nil
}

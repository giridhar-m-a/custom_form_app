package scheduler

import (
	"encoding/json"
	"log"
	"time"

	"github.com/giridhar-m-a/custom_form_app/constants"
	"github.com/giridhar-m-a/custom_form_app/internal/dto/scheduler_dto"
	"github.com/hibiken/asynq"
)

func ScheduleInvitation(formID string, scheduleTime time.Time) (*asynq.TaskInfo, error) {
	client := NewClient()
	defer client.Close()
	invitationPayload := scheduler_dto.InvitationSchedulerPayload{
		FormID: formID,
	}
	payload, _ := json.Marshal(invitationPayload)
	task := asynq.NewTask(constants.TaskTypeInvitationSchedule, payload)

	info, err := client.Enqueue(task, asynq.ProcessAt(scheduleTime))
	if err != nil {
		log.Printf("Error scheduling invitation for form %s: %v", formID, err)
		return nil, err
	}
	log.Printf("Invitation scheduled for form %s at %v, Task ID: %s", formID, scheduleTime, info.ID)
	return info, nil

}

func CancelInvitationSchedule(scheduleId string) error {
	inspector := asynq.NewInspector(NewRedisClientOpt())
	err := inspector.DeleteTask(constants.TaskTypeInvitationSchedule, scheduleId)
	if err != nil {
		return err
	}
	log.Printf("Invitation cancelled for schedule ID: %s", scheduleId)
	return nil
}

func UpdateInvitationSchedule(scheduleID string, scheduleTime time.Time, formID string) (*asynq.TaskInfo, error) {
	err := CancelInvitationSchedule(scheduleID)
	if err != nil {
		return nil, err
	}
	info, err := ScheduleInvitation(formID, scheduleTime)
	if err != nil {
		return nil, err
	}
	log.Printf("Invitation updated for schedule ID: %s, New Task ID: %s", scheduleID, info.ID)
	return info, nil
}

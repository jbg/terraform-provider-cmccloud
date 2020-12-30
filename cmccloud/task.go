package cmccloud

import (
	"errors"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

// func waitForTaskFinished(taskID string, meta interface{}) (interface{}, error) {
// 	log.Printf("[INFO] Waiting for server with task id (%s) to be created", taskID)
// 	stateConf := &resource.StateChangeConf{
// 		Pending:    []string{"WAIT", "PROCESSING"},
// 		Target:     []string{"DONE"},
// 		Refresh:    taskStateRefreshfunc(taskID, meta),
// 		Timeout:    60 * time.Second,
// 		Delay:      20 * time.Second,
// 		MinTimeout: 3 * time.Second,
// 	}
// 	return stateConf.WaitForState()
// }

func taskStateRefreshfunc(taskID string, meta interface{}) resource.StateRefreshFunc {
	client := meta.(*CombinedConfig).goCMCClient()
	return func() (interface{}, string, error) {
		// Get task result from cloud server API
		resp, err := client.Task.Get(taskID)
		if err != nil {
			return nil, "", err
		}
		// if the task is not ready, we need to wait for a moment
		if resp.Status == "ERROR" {
			log.Println("[DEBUG] Task is failed")
			return nil, "", errors.New(fmt.Sprint(resp))
		}

		if resp.Status == "DONE" {
			return resp, "DONE", nil
		}

		log.Println("[DEBUG] Task is not done")
		return nil, "", nil
	}
}

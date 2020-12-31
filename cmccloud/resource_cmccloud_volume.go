package cmccloud

import (
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceCMCCloudVolume() *schema.Resource {
	return &schema.Resource{
		Create: resourceCMCCloudVolumeCreate,
		Read:   resourceCMCCloudVolumeRead,
		Update: resourceCMCCloudVolumeUpdate,
		Delete: resourceCMCCloudVolumeDelete,
		Importer: &schema.ResourceImporter{
			State: resourceCMCCloudVolumeImport,
		},
		SchemaVersion: 1,
		Schema:        volumeSchema(),

		// CustomizeDiff: func(ctx context.Context, diff *schema.ResourceDiff, v interface{}) error {
		// 	// if the new size of the volume is smaller than the old one return an error since
		// 	// only expanding the volume is allowed
		// 	oldSize, newSize := diff.GetChange("size")
		// 	if newSize.(int) < oldSize.(int) {
		// 		return fmt.Errorf("volumes `size` can only be expanded and not shrunk")
		// 	}

		// 	return nil
		// },
		CustomizeDiff: customdiff.All(
			// customdiff.ValidateChange("size", func (old, new, meta interface{}) error {
			//     // If we are increasing "size" then the new value must be
			//     // a multiple of the old value.
			//     if new.(int) <= old.(int) {
			//         return fmt.Errorf("volumes `size` can only be expanded and not shrunk")
			//     }
			//     return nil
			// }),
			customdiff.ForceNewIfChange("size", func(old, new, meta interface{}) bool {
				// "size" can only increase in-place, so we must create a new resource
				// if it is decreased.
				return new.(int) < old.(int)
			}),
		),
	}
}

func resourceCMCCloudVolumeCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).goCMCClient()
	_, taskStatus, err := client.Volume.Create(map[string]interface{}{
		"name":      d.Get("name").(string),
		"size":      d.Get("size").(int),
		"type":      d.Get("type").(string),
		"region":    strings.ToLower(d.Get("region").(string)),
		"server_id": d.Get("server_id").(string),
	})
	if err != nil {
		return fmt.Errorf("Error creating Volume: %s", err)
	}
	d.SetId(taskStatus.ResultID)

	return resourceCMCCloudVolumeRead(d, meta)
}

func resourceCMCCloudVolumeRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).goCMCClient()
	volume, err := client.Volume.Get(d.Id())
	if err != nil {
		return fmt.Errorf("Error retrieving Volume %s: %v", d.Id(), err)
	}

	_ = d.Set("name", volume.Name)
	_ = d.Set("size", volume.Size)
	_ = d.Set("type", volume.Type)
	_ = d.Set("region", strings.ToLower(volume.Region))
	_ = d.Set("state", volume.State)
	_ = d.Set("created", volume.Created)
	_ = d.Set("server_id", volume.ServerID)
	return nil
}

func resourceCMCCloudVolumeUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).goCMCClient()
	id := d.Id()
	if d.HasChange("name") {
		// Resize Volume to new flavor
		err := client.Volume.Rename(id, d.Get("name").(string))
		if err != nil {
			return fmt.Errorf("Error when rename Volume [%s]: %v", id, err)
		}
	}

	if d.HasChange("size") {
		_, _, err := client.Volume.Resize(id, d.Get("size").(int))
		if err != nil {
			return fmt.Errorf("Error when resize volume [%s]: %v", id, err)
		}
	}

	if d.HasChange("server_id") {
		serverID := d.Get("server_id").(string)
		if serverID == "" {
			_, err := client.Volume.Detach(id)
			if err != nil {
				return fmt.Errorf("Error when detach volume [%s]: %v", id, err)
			}
		} else {
			// get current server that volume is attached
			volume, err := client.Volume.Get(id)
			if err != nil {
				return fmt.Errorf("Not found volume [%s]: %v", id, err)
			}
			currentServerID := volume.ServerID
			if currentServerID != "" && currentServerID != serverID {
				_, err := client.Volume.Detach(id)
				if err != nil {
					return fmt.Errorf("Error when detach old server [%s]: %v", currentServerID, err)
				}
			}
			_, err = client.Volume.Attach(id, serverID)
			if err != nil {
				return fmt.Errorf("Error when attach volume to server [%s]: %v", serverID, err)
			}
		}
	}

	// if d.HasChange("enable_backup") {
	// 	if d.Get("enable_backup").(bool) {
	// 		// Enable auto backup
	// 		intervalType, scheduleTime := parseBackupSchedule(d)
	// 		_, err := client.Volume.EnableBackup(id, intervalType, scheduleTime)
	// 		if err != nil {
	// 			return fmt.Errorf("Error enabling auto backup on volume (%s): %s", d.Id(), err)
	// 		}
	// 	} else {
	// 		// Disable auto backup
	// 		_, err := client.Volume.DisableBackup(id)
	// 		if err != nil {
	// 			return fmt.Errorf("Error disabling auto backup on volume (%s): %s", d.Id(), err)
	// 		}
	// 	}
	// }
	// if d.HasChange("backup_schedule") {
	// 	if d.Get("enable_backup").(bool) {
	// 		intervalType, scheduleTime := parseBackupSchedule(d)
	// 		_, err := client.Volume.UpdateScheduleTime(id, intervalType, scheduleTime)
	// 		if err != nil {
	// 			return fmt.Errorf("Update schedule time on volume (%s) error: %s", d.Id(), err)
	// 		}
	// 	}
	// }

	return resourceCMCCloudVolumeRead(d, meta)
}

func resourceCMCCloudVolumeDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).goCMCClient()
	_, err := client.Volume.Delete(d.Id())

	if err != nil {
		return fmt.Errorf("Error delete cloud volume: %v", err)
	}
	return nil
}

func resourceCMCCloudVolumeImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	err := resourceCMCCloudVolumeRead(d, meta)
	return []*schema.ResourceData{d}, err
}

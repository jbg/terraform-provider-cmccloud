package cmccloud

import (
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/cmc-cloud/gocmcapi"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceCMCCloudServer() *schema.Resource {
	return &schema.Resource{
		Create: resourceCMCCloudServerCreate,
		Read:   resourceCMCCloudServerRead,
		Update: resourceCMCCloudServerUpdate,
		Delete: resourceCMCCloudServerDelete,
		Importer: &schema.ResourceImporter{
			State: resourceCMCCloudServerImport,
		},
		SchemaVersion: 1,
		Schema:        serverSchema(),
	}
}

func resourceCMCCloudServerCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).goCMCClient()

	_, taskStatus, err := client.Server.Create(map[string]interface{}{
		"name":                   d.Get("name").(string),
		"cpu":                    d.Get("cpu").(int),
		"ram":                    d.Get("ram").(int),
		"root":                   d.Get("root").(int),
		"gpu":                    d.Get("gpu").(int),
		"image_type":             d.Get("image_type").(string),
		"image_id":               d.Get("image_id").(string),
		"region":                 strings.ToUpper(d.Get("region").(string)),
		"enable_private_network": d.Get("enable_private_network").(bool),
		"enable_backup":          d.Get("enable_backup").(bool),
		// "num_ip_public":          d.Get("num_ip_public").(int),
	})
	if err != nil {
		return fmt.Errorf("Error creating server: %s", err)
	}
	id := taskStatus.ResultID
	d.SetId(id)

	if _, ok := d.GetOk("backup_schedule"); ok {
		if d.Get("auto_backup").(bool) {
			intervalType, scheduleTime := parseBackupSchedule(d)
			_, err := client.Server.UpdateScheduleTime(id, intervalType, scheduleTime)
			if err != nil {
				return fmt.Errorf("Update schedule time on server (%s) error: %s", id, err)
			}
		}
	}

	return resourceCMCCloudServerRead(d, meta)
}

func resourceCMCCloudServerRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).goCMCClient()
	server, err := client.Server.Get(d.Id())
	if err != nil {
		if errors.Is(err, gocmcapi.ErrNotFound) {
			log.Printf("[WARN] CMC Cloud Server with id = (%s) is not found", d.Id())
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error retrieving server %s: %v", d.Id(), err)
	}
	d.Set("cpu", server.CPU)
	d.Set("ram", server.RAM)
	d.Set("root", server.Root)
	d.Set("gpu", server.GPU)
	d.Set("image_id", server.ImageID)
	d.Set("image_type", IfThenElse(server.ImageType == "", "image", server.ImageType))
	d.Set("region", strings.ToLower(server.RegionName))

	d.Set("name", server.Name)
	d.Set("display_name", server.DisplayName)
	d.Set("created", server.Created)
	d.Set("image_name", server.ImageName)
	d.Set("bits", server.Bits)
	d.Set("state", server.State)
	d.Set("auto_backup", server.AutoBackup)
	d.Set("main_ip_address", server.MainIPAddress)
	d.Set("backup_schedule", server.BackupSchedule)

	isPrivate := false
	for _, nic := range server.Nics {
		if nic.IsPrivate {
			isPrivate = true
		}
	}
	d.Set("enable_private_network", isPrivate)
	return nil
}

func resourceCMCCloudServerUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).goCMCClient()
	id := d.Id()
	if d.HasChange("cpu") || d.HasChange("ram") || d.HasChange("root") || d.HasChange("gpu") {
		// Resize server to new flavor
		_, _, err := client.Server.Resize(id, d.Get("cpu").(int), d.Get("ram").(int), d.Get("root").(int), d.Get("gpu").(int))
		if err != nil {
			return fmt.Errorf("Error when resize server [%s]: %v", id, err)
		}
	}

	if d.HasChange("name") {
		old, new := d.GetChange("name")
		_, err := client.Server.Rename(id, d.Get("name").(string))
		if err != nil {
			return fmt.Errorf("Error when rename server [%s] from `%s` to `%s`: %v", id, old, new, err)
		}
	}

	if d.HasChange("enable_private_network") {
		if d.Get("enable_private_network").(bool) {
			// Enable private network
			_, err := client.Server.EnablePrivateNetwork(id)
			if err != nil {
				return fmt.Errorf("Error enabling private network on server (%s): %s", d.Id(), err)
			}
		} else {
			// Disable private network
			_, err := client.Server.DisablePrivateNetwork(id)
			if err != nil {
				return fmt.Errorf("Error disabling private network on server (%s): %s", d.Id(), err)
			}
		}
	}

	if d.HasChange("auto_backup") {
		if d.Get("auto_backup").(bool) {
			// Enable auto backup
			intervalType, scheduleTime := parseBackupSchedule(d)
			_, _, err := client.Server.EnableBackup(id, intervalType, scheduleTime)
			if err != nil {
				return fmt.Errorf("Error enabling auto backup on server (%s): %s", d.Id(), err)
			}
		} else {
			// Disable auto backup
			_, err := client.Server.DisableBackup(id)
			if err != nil {
				return fmt.Errorf("Error disabling auto backup on server (%s): %s", d.Id(), err)
			}
		}
	}
	if d.HasChange("backup_schedule") {
		if d.Get("auto_backup").(bool) {
			intervalType, scheduleTime := parseBackupSchedule(d)
			_, err := client.Server.UpdateScheduleTime(id, intervalType, scheduleTime)
			if err != nil {
				return fmt.Errorf("Update schedule time on server (%s) error: %s", d.Id(), err)
			}
		}
	}

	return resourceCMCCloudServerRead(d, meta)
}
func parseBackupSchedule(d *schema.ResourceData) (string, string) {
	backupSchedule := d.Get("backup_schedule").(string)
	if backupSchedule == "" {
		backupSchedule = "daily-02:30"
	}
	intervalType := strings.Split(backupSchedule, "-")[0]
	scheduleTime := strings.Split(backupSchedule, "-")[1]
	return intervalType, scheduleTime
}

func resourceCMCCloudServerDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).goCMCClient()
	// destroy the cloud server
	_, err := client.Server.Delete(d.Id())

	if err != nil {
		return fmt.Errorf("Error delete cloud server: %v", err)
	}
	/*var rootDiskID string
	for _, v := range server.AttachedVolumes {
		if v.AttachedType == attachTypeRootDisk {
			rootDiskID = v.ID
		}
	}
	err = client.Server.Delete(context.Background(), d.Id())
	if err != nil {
		return fmt.Errorf("Error delete cloud server %v", err)
	}
	err = client.Volume.Delete(context.Background(), rootDiskID)
	if err != nil {
		return fmt.Errorf("Error deleting rootdisk of cloud server %s: %v", d.Id(), err)
	} */
	return nil
}

func resourceCMCCloudServerImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	resourceCMCCloudServerRead(d, meta)
	return []*schema.ResourceData{d}, nil
}

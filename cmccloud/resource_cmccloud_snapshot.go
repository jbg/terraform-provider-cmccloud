package cmccloud

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceCMCCloudSnapshot() *schema.Resource {
	return &schema.Resource{
		Create: resourceCMCCloudSnapshotCreate,
		Read:   resourceCMCCloudSnapshotRead,
		Update: resourceCMCCloudSnapshotUpdate,
		Delete: resourceCMCCloudSnapshotDelete,
		Importer: &schema.ResourceImporter{
			State: resourceCMCCloudSnapshotImport,
		},
		SchemaVersion: 1,
		Schema:        snapshotSchema(),
	}
}

func resourceCMCCloudSnapshotCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).goCMCClient()
	name := ""
	if _, ok := d.GetOk("name"); ok {
		name = d.Get("name").(string)
	}
	if _, ok := d.GetOk("volume_id"); ok {
		_, taskStatus, err := client.Snapshot.Create(d.Get("volume_id").(string), name)
		if err != nil {
			return fmt.Errorf("Error creating Snapshot: %s", err)
		}
		d.SetId(taskStatus.ResultID)
		return resourceCMCCloudSnapshotRead(d, meta)
	}

	if _, ok := d.GetOk("server_id"); ok {
		_, taskStatus, err := client.Server.TakeSnapshot(d.Get("server_id").(string), name)
		if err != nil {
			return fmt.Errorf("Error creating Snapshot: %s", err)
		}
		d.SetId(taskStatus.ResultID)
		return resourceCMCCloudSnapshotRead(d, meta)
	}
	return fmt.Errorf("volume_id or server_id is required")
}

func resourceCMCCloudSnapshotRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).goCMCClient()
	snapshot, err := client.Snapshot.Get(d.Id())
	if err != nil {
		return fmt.Errorf("Error retrieving Snapshot %s: %v", d.Id(), err)
	}

	d.Set("name", snapshot.Name)
	d.Set("size", snapshot.Size)
	d.Set("state", snapshot.State)
	d.Set("created", snapshot.Created)
	d.Set("volume_id", snapshot.VolumeID)
	return nil
}

func resourceCMCCloudSnapshotUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).goCMCClient()
	id := d.Id()
	if d.HasChange("name") {
		// Resize Snapshot to new flavor
		err := client.Snapshot.Rename(id, d.Get("name").(string))
		if err != nil {
			return fmt.Errorf("Error when rename Snapshot [%s]: %v", id, err)
		}
	}
	return resourceCMCCloudSnapshotRead(d, meta)
}

func resourceCMCCloudSnapshotDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).goCMCClient()
	_, err := client.Snapshot.Delete(d.Id())

	if err != nil {
		return fmt.Errorf("Error delete cloud volume: %v", err)
	}
	return nil
}

func resourceCMCCloudSnapshotImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	resourceCMCCloudSnapshotRead(d, meta)
	return []*schema.ResourceData{d}, nil
}

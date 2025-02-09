package bizflycloud

import (
	"context"
	"fmt"
	"github.com/bizflycloud/gobizfly"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceBizflyCloudSimpleStorageBucket() *schema.Resource {
	return &schema.Resource{
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Create: resourceBizflyCloudSimpleStorageBucketCreate,
		Read:   resourceBizflyCloudSimpleStorageBucketRead,
		Delete: resourceBizflyCloudSimpleStorageBucketDelete,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"location": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"default_storage_class": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(20 * time.Minute),
		},
	}
}

func resourceBizflyCloudSimpleStorageBucketCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).gobizflyClient()
	bcr := gobizfly.BucketCreateRequest{
		Name:                d.Get("name").(string),
		Location:            d.Get("location").(string),
		Acl:                 "private",
		DefaultStorageClass: d.Get("default_storage_class").(string),
	}
	ss, err := client.CloudSimpleStorage.Create(context.Background(), &bcr)
	if err != nil {
		return fmt.Errorf("Error when creating simple store bucket: %v", err)
	}
	d.SetId(ss.Name)
	return resourceBizflyCloudSimpleStorageBucketRead(d, meta)
}

func resourceBizflyCloudSimpleStorageBucketRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).gobizflyClient()
	simpleStores, err := client.CloudSimpleStorage.List(context.Background(), &gobizfly.ListOptions{})
	if err != nil {
		return fmt.Errorf("Error retrieving simple store bucket: %v", err)
	}
	for _, simpleStore := range simpleStores {
		_ = d.Set("name", simpleStore.Name)
		_ = d.Set("location", simpleStore.Location)
		_ = d.Set("created_at", simpleStore.CreatedAt)
		_ = d.Set("num_objects", simpleStore.NumObjects)
		_ = d.Set("size_kb", simpleStore.SizeKb)
	}
	return nil
}

func resourceBizflyCloudSimpleStorageBucketDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).gobizflyClient()
	err := client.CloudSimpleStorage.Delete(context.Background(), d.Id())
	if err != nil {
		return fmt.Errorf("Error deleting simple store bucket: %v", err)
	}
	return nil
}

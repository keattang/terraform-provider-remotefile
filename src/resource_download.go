package main

import (
	"os"

	"github.com/hashicorp/terraform/helper/schema"
)

func resourceDownloadCreate(d *schema.ResourceData, m interface{}) error {
	filePath := d.Get("destination").(string)
	url := d.Get("source_url").(string)
	base64sha256, md5, etag, err := downloadFile(filePath, url)

	if err != nil {
		return err
	}

	d.Set("file_base64sha256", base64sha256)
	d.Set("file_md5", md5)
	d.Set("etag", etag)
	d.SetId(base64sha256)

	return nil
}

func resourceDownloadRead(d *schema.ResourceData, m interface{}) error {
	destFilePath := d.Get("destination").(string)
	url := d.Get("source_url").(string)
	oldBase64sha256 := d.Get("file_base64sha256").(string)
	etag := d.Get("etag").(string)

	// If the destination file doesn't exist we need to run create again
	if _, err := os.Stat(destFilePath); os.IsNotExist(err) {
		d.SetId("")
		return nil
	}

	// If the destination file isn't the same as when we last ran we need to run create again
	currentBase64Sha256, err := hashFileBase64Sha256(destFilePath)
	if err != nil {
		return err
	} else if currentBase64Sha256 != oldBase64sha256 {
		d.SetId("")
		return nil
	}

	// Check if the remote file has changed, if it has we need to run create again
	remoteChanged, err := checkIfRemoteFileChanged(url, oldBase64sha256, etag)
	if err != nil {
		return err
	} else if remoteChanged {
		d.SetId("")
	}

	return nil
}

func resourceDownloadDelete(d *schema.ResourceData, m interface{}) error {
	filePath := d.Get("destination").(string)
	os.Remove(filePath)
	return nil
}

func resourceDownload() *schema.Resource {
	return &schema.Resource{
		Create: resourceDownloadCreate,
		Read:   resourceDownloadRead,
		Delete: resourceDownloadDelete,

		Schema: map[string]*schema.Schema{
			"source_url": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"destination": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"file_base64sha256": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"file_md5": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"etag": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

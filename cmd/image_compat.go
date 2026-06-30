package cmd

import (
	"io"
	"strings"

	"github.com/spf13/cobra"

	sdk "github.com/ucloud/ucloud-sdk-go/ucloud"

	"github.com/ucloud/ucloud-cli/base"
	"github.com/ucloud/ucloud-cli/internal/common"
	"github.com/ucloud/ucloud-cli/model/cli"
	"github.com/ucloud/ucloud-cli/pkg/command"
)

// Transition shim — image migrated to products/image (Part 5), but cmd/ still
// has unmigrated image consumers that keep the cmd-local (base.BizClient)
// definitions, with the SAME names/signatures, so cmd/ keeps compiling:
//   - describeImageByID: cmd/uhost.go ~:372 (uhost create image lookup), ~:1510 (create-image poll)
//   - getImageList:      cmd/uhost.go ~:554 (uhost create --image-id completion)
//   - ImageRow / NewCmdUImageList: cmd/uhost_test.go (TestUhost live integration,
//     fetches an image ID for the uhost create flow; moves with uhost in Part 6)
// REMOVE this file when uhost migrates (Part 6) and these become product-local.
// NewCmdUImageList is intentionally NOT registered in cmd/root.go — `image` is
// now served by the products/image package, so this constructor exists solely
// for the test and never appears in the command tree.

func getImageList(states []string, imageType, project, region, zone string) []string {
	req := base.BizClient.NewDescribeImageRequest()
	req.ProjectId = &project
	req.Region = &region
	req.Zone = &zone
	req.Limit = sdk.Int(1000)
	if imageType != cli.IMAGE_ALL {
		req.ImageType = sdk.String(imageType)
	}
	resp, err := base.BizClient.DescribeImage(req)
	if err != nil {
		return nil
	}
	list := []string{}
	for _, image := range resp.ImageSet {
		for _, s := range states {
			if image.State == s {
				list = append(list, image.ImageId+"/"+image.ImageName)
			}
		}
	}
	return list
}

// ImageRow 表格行 — cmd-local copy kept for cmd/uhost_test.go. Mirrors
// products/image/internal/image.ImageRow.
type ImageRow struct {
	ImageName         string
	ImageID           string
	ImageType         string
	BasicImage        string
	ExtensibleFeature string
	CreationTime      string
	State             string
}

// NewCmdUImageList builds the `image list` command (cmd-local, base.BizClient).
// Kept for cmd/uhost_test.go only; not registered in the command tree.
func NewCmdUImageList(out io.Writer) *cobra.Command {
	req := base.BizClient.NewDescribeImageRequest()
	cmd := &cobra.Command{
		Use:     "list",
		Short:   "List image",
		Long:    "List image",
		Example: "ucloud image list --image-type Base",
		Run: func(cmd *cobra.Command, args []string) {
			resp, err := base.BizClient.DescribeImage(req)
			if err != nil {
				base.HandleError(err)
				return
			}
			list := make([]ImageRow, 0)
			for _, image := range resp.ImageSet {
				row := ImageRow{}
				row.ImageName = image.ImageName
				row.ImageID = image.ImageId
				row.ImageType = image.ImageType
				row.BasicImage = image.OsName
				row.ExtensibleFeature = strings.Join(image.Features, ",")
				row.CreationTime = common.FormatDate(image.CreateTime)
				row.State = image.State
				if row.State == "Available" {
					list = append(list, row)
				}
			}
			base.PrintList(list, out)
		},
	}
	req.ProjectId = cmd.Flags().String("project-id", base.ConfigIns.ProjectID, "Optional. Assign project-id")
	req.Region = cmd.Flags().String("region", base.ConfigIns.Region, "Optional. Assign region")
	req.Zone = cmd.Flags().String("zone", "", "Optional. Assign availability zone")
	req.ImageType = cmd.Flags().String("image-type", "Base", "Optional. 'Base',Standard image; 'Business',image market; 'Custom',custom image")
	req.OsType = cmd.Flags().String("os-type", "", "Optional. Linux or Windows. Return all types by default")
	req.ImageId = cmd.Flags().String("image-id", "", "Optional. Resource ID of image")
	req.Offset = cmd.Flags().Int("offset", 0, "Optional. Offset default 0")
	req.Limit = cmd.Flags().Int("limit", 500, "Optional. Max count")
	command.SetFlagValues(cmd, "image-type", "Base", "Business", "Custom")
	return cmd
}

func describeImageByID(imageID, project, region, zone string) (interface{}, error) {
	req := base.BizClient.NewDescribeImageRequest()
	req.ImageId = sdk.String(imageID)
	req.ProjectId = sdk.String(project)
	req.Region = sdk.String(region)
	req.Zone = sdk.String(zone)
	req.Limit = sdk.Int(50)
	resp, err := base.BizClient.DescribeImage(req)
	if err != nil {
		return nil, err
	}
	if len(resp.ImageSet) < 1 {
		return nil, nil
	}
	return &resp.ImageSet[0], nil
}

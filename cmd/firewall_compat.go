package cmd

import (
	"github.com/ucloud/ucloud-sdk-go/services/unet"
	sdk "github.com/ucloud/ucloud-sdk-go/ucloud"

	"github.com/ucloud/ucloud-cli/base"
)

// Cross-product transition shim.
//
// The firewall product has moved to products/firewall (Part 2 of the batch-1
// migration). But cmd/uhost.go still registers its --firewall-id completion via
// getFirewallIDNames, and uhost has not been migrated yet. These cmd-local
// copies (identical logic to the originals in the now-deleted cmd/firewall.go)
// keep package cmd compiling until uhost is migrated, at which point this shim
// can be removed. getFirewall is intentionally NOT copied here — uhost does not
// use it.
func getFirewallIDNames(project, region string) (idNames []string) {
	list, err := getAllFirewallIns(project, region)
	if err != nil {
		return
	}
	for _, f := range list {
		idNames = append(idNames, f.FWId+"/"+f.Name)
	}
	return
}

func getAllFirewallIns(project, region string) ([]unet.FirewallDataSet, error) {
	req := base.BizClient.NewDescribeFirewallRequest()
	req.ProjectId = sdk.String(project)
	req.Region = sdk.String(region)
	list := []unet.FirewallDataSet{}
	for offset, limit := 0, 100; ; offset += limit {
		req.Offset = sdk.Int(offset)
		req.Limit = sdk.Int(limit)
		resp, err := base.BizClient.DescribeFirewall(req)
		if err != nil {
			return nil, err
		}
		for _, fw := range resp.DataSet {
			list = append(list, fw)
		}
		if resp.TotalCount < offset+limit {
			break
		}
	}
	return list, nil
}

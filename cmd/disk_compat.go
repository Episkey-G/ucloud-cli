// Copyright © 2018 NAME HERE tony.li@ucloud.cn
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"fmt"
	"io"

	"github.com/ucloud/ucloud-sdk-go/services/udisk"
	sdk "github.com/ucloud/ucloud-sdk-go/ucloud"
	"github.com/ucloud/ucloud-sdk-go/ucloud/request"

	"github.com/ucloud/ucloud-cli/base"
	"github.com/ucloud/ucloud-cli/model/status"
)

// disk_compat.go is a cross-product transition shim. The udisk product moved to
// products/udisk (its own copy lives there as the exported DetachUdisk), but
// cmd/uhost.go still calls detachUdisk in its delete flow (reinstall/detach
// before removing the host). Until uhost itself migrates, this keeps the
// cmd-local detachUdisk (and its describeUdiskByID dependency) with the exact
// same name+signature as the original cmd/disk.go, so cmd/uhost.go compiles
// unchanged. Only what uhost needs is kept here.

func detachUdisk(async bool, udiskID string, out io.Writer) error {
	any, err := describeUdiskByID(udiskID, nil)
	if err != nil {
		return err
	}
	if any == nil {
		return fmt.Errorf("udisk[%v] is not exist", any)
	}
	ins, ok := any.(*udisk.UDiskDataSet)
	if !ok {
		return fmt.Errorf("%#v convert to udisk failed", any)
	}
	req := base.BizClient.NewDetachUDiskRequest()
	req.UHostId = sdk.String(ins.UHostId)
	req.UDiskId = sdk.String(udiskID)
	resp, err := base.BizClient.DetachUDisk(req)
	if err != nil {
		return err
	}
	text := fmt.Sprintf("udisk[%s] is detaching from uhost[%s]", resp.UDiskId, resp.UHostId)
	if async {
		fmt.Fprintln(out, text)
	} else {
		poller := base.NewSpoller(describeUdiskByID, out)
		poller.Spoll(udiskID, text, []string{status.DISK_AVAILABLE, status.DISK_FAILED})
	}
	return nil
}

func describeUdiskByID(udiskID string, commonBase *request.CommonBase) (interface{}, error) {
	req := base.BizClient.NewDescribeUDiskRequest()
	if commonBase != nil {
		req.CommonBase = *commonBase
	}
	req.UDiskId = sdk.String(udiskID)
	req.Limit = sdk.Int(50)
	resp, err := base.BizClient.DescribeUDisk(req)
	if err != nil {
		return nil, err
	}
	if len(resp.DataSet) < 1 {
		return nil, nil
	}
	return &resp.DataSet[0], nil
}

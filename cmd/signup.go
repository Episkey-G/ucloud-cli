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
	"os/exec"
	"runtime"

	"github.com/spf13/cobra"
	"github.com/ucloud/ucloud-cli/base"
)

// NewCmdSignup ucloud signup
func NewCmdSignup() *cobra.Command {
	var cmd = &cobra.Command{
		Use:     "signup",
		Short:   fmt.Sprintf("Launch %s sign up page in browser", base.BrandName),
		Long:    fmt.Sprintf(`Launch %s sign up page in browser`, base.BrandName),
		Args:    cobra.NoArgs,
		Example: fmt.Sprintf("%s signup", base.BrandNameLower),
		Run: func(cmd *cobra.Command, args []string) {
			url := fmt.Sprintf("https://passport.%s/#register", base.BrandURL)
			openbrowser(url)
		},
	}
	return cmd
}

func openbrowser(url string) {
	var err error

	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", url).Start()
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin":
		err = exec.Command("open", url).Start()
	default:
		err = fmt.Errorf("Open url: %s in your browser", url)
	}
	if err != nil {
		fmt.Printf("Open url: %s in your browser\n", url)
	}
}

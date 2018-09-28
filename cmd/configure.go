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
	"strings"

	"github.com/spf13/cobra"

	. "github.com/ucloud/ucloud-cli/util"
)

var config = ConfigInstance

//NewCmdConfig ucloud config
func NewCmdConfig() *cobra.Command {
	var configDesc = `Public-key and private-key could be acquired from https://console.ucloud.cn/uapi/apikey.`
	var helloUcloud = `
  _   _      _ _         _   _ _____ _                 _ 
  | | | |    | | |       | | | /  __ \ |               | |
  | |_| | ___| | | ___   | | | | /  \/ | ___  _   _  __| |
  |  _  |/ _ \ | |/ _ \  | | | | |   | |/ _ \| | | |/ _\ |
  | | | |  __/ | | (_) | | |_| | \__/\ | (_) | |_| | (_| |
  \_| |_/\___|_|_|\___/   \___/ \____/_|\___/ \__,_|\__,_|
  `

	var configCmd = &cobra.Command{
		Use:     "config",
		Short:   "Config UCloud CLI options",
		Long:    `Config UCloud CLI options such as private-key,public-key,default region and default project-id.`,
		Example: "ucloud config; ucloud config set region cn-bj2; ucloud config set project org-xxx",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println(configDesc)
			if len(config.PrivateKey) != 0 && len(config.PublicKey) != 0 {
				fmt.Printf("Your have already configured public-key and private-key. Do you want to overwrite it? (y/n):")
				var overwrite string
				_, err := fmt.Scanf("%s\n", &overwrite)
				if err != nil {
					fmt.Println(err)
					return
				}
				overwrite = strings.Trim(overwrite, " ")
				overwrite = strings.ToLower(overwrite)
				if overwrite != "yes" && overwrite != "y" {
					return
				}
			}
			config.ClearConfig()
			ClientConfig.Region = ""
			ClientConfig.ProjectId = ""
			config.ConfigPublicKey()
			config.ConfigPrivateKey()

			region, zone, err := getDefaultRegion()
			if err != nil {
				Cxt.Println(err)
				return
			}
			config.Region = region
			config.Zone = zone
			Cxt.Printf("Configured default region:%s zone:%s\n", region, zone)

			projectId, projectName, err := getDefaultProject()
			if err != nil {
				Cxt.Println(err)
				return
			}
			config.ProjectID = projectId
			Cxt.Printf("Configured default project:%s %s\n", projectId, projectName)

			config.SaveConfig()

			userInfo, err := getUserInfo()

			Cxt.Printf("You are logged in as: [%s]\n", userInfo.UserEmail)

			certified := isUserCertified(userInfo)
			if err != nil {
				fmt.Println(err)
			} else if certified == false {
				fmt.Println("\nWarning: Please authenticate the account with your valid documentation at 'https://accountv2.ucloud.cn/authentication'.")
			}
			fmt.Println(helloUcloud)
		},
	}

	configCmd.AddCommand(NewCmdConfigList())
	configCmd.AddCommand(NewCmdConfigClear())
	configCmd.AddCommand(NewCmdConfigSet())

	originHelpFunc := configCmd.HelpFunc()

	configCmd.SetHelpFunc(func(cmd *cobra.Command, args []string) {
		rootCmd := cmd.Parent()
		rootCmd.Flags().MarkHidden("region")
		rootCmd.Flags().MarkHidden("project-id")
		originHelpFunc(cmd, args)
	})
	return configCmd
}

//NewCmdConfigList ucloud config ls
func NewCmdConfigList() *cobra.Command {
	var configListCmd = &cobra.Command{
		Use:   "list",
		Short: "list all settings",
		Long:  `list all settings`,
		Run: func(cmd *cobra.Command, args []string) {
			config.ListConfig(global.json)
		},
	}
	return configListCmd
}

//NewCmdConfigClear ucloud config clear
func NewCmdConfigClear() *cobra.Command {
	var configClearCmd = &cobra.Command{
		Use:   "clear",
		Short: "clear all settings",
		Long:  "clear all settings",
		Run: func(cmd *cobra.Command, args []string) {
			config.ClearConfig()
		},
	}
	return configClearCmd
}

//NewCmdConfigSet ucloud config set
func NewCmdConfigSet() *cobra.Command {

	var configSetCmd = &cobra.Command{
		Use:     "set",
		Short:   "Set a config value",
		Long:    "Set a config value, including private-key public-key region and project-id.",
		Example: "ucloud config set region cn-bj2",
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) != 2 {
				fmt.Printf("Error: accepts 2 arg(s), received %d\n", len(args))
				return
			}
			switch args[0] {
			case "region":
				config.Region = args[1]
			case "zone":
				config.Zone = args[1]
			case "project-id":
				config.ProjectID = args[1]
			case "public-key":
				config.PublicKey = args[1]
			case "private-key":
				config.PrivateKey = args[1]
			default:
				Cxt.Println("Only public-key, private-key, region, zone and project-id supported")
			}
			config.SaveConfig()
		},
	}
	return configSetCmd
}

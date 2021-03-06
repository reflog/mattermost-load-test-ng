// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

package main

import (
	"fmt"
	"os"

	"github.com/mattermost/mattermost-load-test-ng/deployment"
	"github.com/mattermost/mattermost-load-test-ng/deployment/terraform"
	"github.com/mattermost/mattermost-load-test-ng/logger"
	"github.com/mattermost/mattermost-server/v5/mlog"

	"github.com/spf13/cobra"
)

func RunCreateCmdF(cmd *cobra.Command, args []string) error {
	config, err := getConfig(cmd)
	if err != nil {
		mlog.Error(err.Error())
		return nil
	}

	t := terraform.New(config)
	err = t.Create()
	if err != nil {
		mlog.Error(err.Error())
	}
	return nil
}

func RunDestroyCmdF(cmd *cobra.Command, args []string) error {
	config, err := getConfig(cmd)
	if err != nil {
		mlog.Error(err.Error())
		return nil
	}

	t := terraform.New(config)
	err = t.Destroy()
	if err != nil {
		mlog.Error(err.Error())
	}
	return nil
}

func getConfig(cmd *cobra.Command) (*deployment.Config, error) {
	configFilePath, _ := cmd.Flags().GetString("config")
	cfg, err := deployment.ReadConfig(configFilePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config: %w", err)
	}

	if err := cfg.IsValid(); err != nil {
		return nil, fmt.Errorf("failed to validate config: %w", err)
	}

	logger.Init(&cfg.LogSettings)
	return cfg, nil
}

func main() {
	rootCmd := &cobra.Command{
		Use:   "deployer",
		Short: "Create and destroy load test environments",
	}
	rootCmd.PersistentFlags().StringP("config", "c", "", "path to the configuration file to use")

	commands := []*cobra.Command{
		{
			Use:   "create",
			Short: "Deploy a load test environment",
			RunE:  RunCreateCmdF,
		},
		{
			Use:   "destroy",
			Short: "Destroy a load test environment",
			RunE:  RunDestroyCmdF,
		},
	}

	rootCmd.AddCommand(commands...)
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

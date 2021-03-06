// Copyright 2017 The OpenSDS Authors.
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
/*
This module implements a entry into the OpenSDS service.

*/

package cli

import (
	"log"
	"os"
	"strconv"

	"github.com/opensds/opensds/pkg/model"
	"github.com/spf13/cobra"
)

var volumeCommand = &cobra.Command{
	Use:   "volume",
	Short: "manage volumes in the cluster",
	Run:   volumeAction,
}

var volumeCreateCommand = &cobra.Command{
	Use:     "create <size>",
	Short:   "create a volume in the cluster",
	Example: "osdsctl volume create 1 --name vol-name",
	Run:     volumeCreateAction,
}

var volumeShowCommand = &cobra.Command{
	Use:   "show <id>",
	Short: "show a volume in the cluster",
	Run:   volumeShowAction,
}

var volumeListCommand = &cobra.Command{
	Use:   "list",
	Short: "list all volumes in the cluster",
	Run:   volumeListAction,
}

var volumeDeleteCommand = &cobra.Command{
	Use:   "delete <id>",
	Short: "delete a volume in the cluster",
	Run:   volumeDeleteAction,
}

var volumeUpdateCommand = &cobra.Command{
	Use:   "update <id>",
	Short: "update a volume in the cluster",
	Run:   volumeUpdateAction,
}

var volumeExtendCommand = &cobra.Command{
	Use:   "extend <id> <new size>",
	Short: "extend a volume in the cluster",
	Run:   volumeExtendAction,
}

var (
	profileId string
	volName   string
	volDesp   string
	volAz     string
	volSnap   string
)

var (
	volLimit          string
	volOffset         string
	volSortDir        string
	volSortKey        string
	volId             string
	volTenantId       string
	volUserId         string
	volStatus         string
	volPoolId         string
	volProfileId      string
	volGroupId        string
	snapshotFromCloud bool
)

func init() {
	volumeListCommand.Flags().StringVarP(&volLimit, "limit", "", "50", "the number of entries displayed per page")
	volumeListCommand.Flags().StringVarP(&volOffset, "offset", "", "0", "all requested data offsets")
	volumeListCommand.Flags().StringVarP(&volSortDir, "sortDir", "", "desc", "the sort direction of all requested data. supports asc or desc(default)")
	volumeListCommand.Flags().StringVarP(&volSortKey, "sortKey", "", "id",
		"the sort key of all requested data. supports id(default), name, status, availabilityzone, profileid, tenantid, size, poolid, description")
	volumeListCommand.Flags().StringVarP(&volId, "id", "", "", "list volume by id")
	volumeListCommand.Flags().StringVarP(&volName, "name", "", "", "list volume by name")
	volumeListCommand.Flags().StringVarP(&volDesp, "description", "", "", "list volume by description")
	volumeListCommand.Flags().StringVarP(&volTenantId, "tenantId", "", "", "list volume by tenantId")
	volumeListCommand.Flags().StringVarP(&volUserId, "userId", "", "", "list volume by storage userId")
	volumeListCommand.Flags().StringVarP(&volStatus, "status", "", "", "list volume by status")
	volumeListCommand.Flags().StringVarP(&volPoolId, "poolId", "", "", "list volume by poolId")
	volumeListCommand.Flags().StringVarP(&volAz, "availabilityZone", "", "", "list volume by availability zone")
	volumeListCommand.Flags().StringVarP(&volProfileId, "profileId", "", "", "list volume by profile id")
	volumeListCommand.Flags().StringVarP(&volGroupId, "groupId", "", "", "list volume by volume group id")

	volumeCommand.PersistentFlags().StringVarP(&profileId, "profile", "p", "", "the id of profile configured by admin")

	volumeCommand.AddCommand(volumeCreateCommand)
	volumeCreateCommand.Flags().StringVarP(&volName, "name", "n", "", "the name of created volume")
	volumeCreateCommand.Flags().StringVarP(&volDesp, "description", "d", "", "the description of created volume")
	volumeCreateCommand.Flags().StringVarP(&volAz, "az", "a", "", "the availability zone of created volume")
	volumeCreateCommand.Flags().StringVarP(&volSnap, "snapshot", "s", "", "the snapshot to create volume")
	volumeCreateCommand.Flags().StringVarP(&poolId, "pool", "l", "", "the pool to create volume")
	volumeCreateCommand.Flags().BoolVarP(&snapshotFromCloud, "snapshotFromCloud", "c", false, "download snapshot from cloud")
	volumeCommand.AddCommand(volumeShowCommand)
	volumeCommand.AddCommand(volumeListCommand)
	volumeCommand.AddCommand(volumeDeleteCommand)
	volumeCommand.AddCommand(volumeUpdateCommand)
	volumeUpdateCommand.Flags().StringVarP(&volName, "name", "n", "", "the name of updated volume")
	volumeUpdateCommand.Flags().StringVarP(&volDesp, "description", "d", "", "the description of updated volume")
	volumeCommand.AddCommand(volumeExtendCommand)

	volumeCommand.AddCommand(volumeSnapshotCommand)
	volumeCommand.AddCommand(volumeAttachmentCommand)
	volumeCommand.AddCommand(volumeGroupCommand)
	volumeCommand.AddCommand(replicationCommand)
}

func volumeAction(cmd *cobra.Command, args []string) {
	cmd.Usage()
	os.Exit(1)
}

var volFormatters = FormatterList{"Metadata": JsonFormatter}

func volumeCreateAction(cmd *cobra.Command, args []string) {
	ArgsNumCheck(cmd, args, 1)
	size, err := strconv.Atoi(args[0])
	if err != nil {
		Fatalln("input size is not valid. It only support integer.")
		log.Fatalf("error parsing size %s: %+v", args[0], err)
	}

	vol := &model.VolumeSpec{
		Name:              volName,
		Description:       volDesp,
		AvailabilityZone:  volAz,
		Size:              int64(size),
		ProfileId:         profileId,
		PoolId:            poolId,
		SnapshotId:        volSnap,
		SnapshotFromCloud: snapshotFromCloud,
	}

	resp, err := client.CreateVolume(vol)
	if err != nil {
		Fatalln(HttpErrStrip(err))
	}

	keys := KeyList{"Id", "CreatedAt", "Name", "Description", "Size", "AvailabilityZone",
		"Status", "PoolId", "ProfileId", "Metadata", "GroupId", "MultiAttach"}
	PrintDict(resp, keys, volFormatters)
}

func volumeShowAction(cmd *cobra.Command, args []string) {
	ArgsNumCheck(cmd, args, 1)
	resp, err := client.GetVolume(args[0])
	if err != nil {
		Fatalln(HttpErrStrip(err))
	}
	keys := KeyList{"Id", "CreatedAt", "UpdatedAt", "Name", "Description", "Size",
		"AvailabilityZone", "Status", "PoolId", "ProfileId", "Metadata", "GroupId", "SnapshotId",
		"MultiAttach"}
	PrintDict(resp, keys, volFormatters)
}

func volumeListAction(cmd *cobra.Command, args []string) {
	ArgsNumCheck(cmd, args, 0)

	var opts = map[string]string{"limit": volLimit, "offset": volOffset, "sortDir": volSortDir,
		"sortKey": volSortKey, "Id": volId,
		"Name": volName, "Description": volDesp, "UserId": volUserId, "AvailabilityZone": volAz,
		"Status": volStatus, "PoolId": volPoolId, "ProfileId": volProfileId, "GroupId": volGroupId}

	resp, err := client.ListVolumes(opts)
	if err != nil {
		Fatalln(HttpErrStrip(err))
	}
	keys := KeyList{"Id", "Name", "Description", "Size", "Status", "ProfileId", "AvailabilityZone"}
	PrintList(resp, keys, volFormatters)
}

func volumeDeleteAction(cmd *cobra.Command, args []string) {
	ArgsNumCheck(cmd, args, 1)
	vol := &model.VolumeSpec{
		ProfileId: profileId,
	}
	err := client.DeleteVolume(args[0], vol)
	if err != nil {
		Fatalln(HttpErrStrip(err))
	}
}

func volumeUpdateAction(cmd *cobra.Command, args []string) {
	ArgsNumCheck(cmd, args, 1)
	vol := &model.VolumeSpec{
		Name:        volName,
		Description: volDesp,
	}

	resp, err := client.UpdateVolume(args[0], vol)
	if err != nil {
		Fatalln(HttpErrStrip(err))
	}
	keys := KeyList{"Id", "UpdatedAt", "Name", "Description", "Size", "AvailabilityZone",
		"Status", "PoolId", "ProfileId", "Metadata", "GroupId", "MultiAttach"}
	PrintDict(resp, keys, volFormatters)
}

func volumeExtendAction(cmd *cobra.Command, args []string) {
	ArgsNumCheck(cmd, args, 2)
	newSize, err := strconv.Atoi(args[1])
	if err != nil {
		log.Fatalf("error parsing new size %s: %+v", args[1], err)
	}

	body := &model.ExtendVolumeSpec{
		NewSize: int64(newSize),
	}

	resp, err := client.ExtendVolume(args[0], body)
	if err != nil {
		Fatalln(HttpErrStrip(err))
	}
	keys := KeyList{"Id", "CreatedAt", "UpdatedAt", "Name", "Description", "Size",
		"AvailabilityZone", "Status", "PoolId", "ProfileId", "Metadata", "GroupId", "MultiAttach"}
	PrintDict(resp, keys, volFormatters)
}

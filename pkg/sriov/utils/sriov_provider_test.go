// Copyright (c) 2020 Doc.ai and/or its affiliates.
//
// SPDX-License-Identifier: Apache-2.0
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at:
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package utils_test

import (
	"bytes"
	"context"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/networkservicemesh/sdk-sriov/pkg/sriov/utils"
)

const (
	totalVfFile          = "sriov_totalvfs"
	configuredVfFile     = "sriov_numvfs"
	physicalFunctionPath = "physfn"
	netInterfacesPath    = "net"
	sriovTestDir         = "nsm/sriov/test"
	pciAddr              = "0000:01:00:0"
	pciAddrShortForm     = "01:00:0"
	pciAddr2             = "0000:02:00:0"
	iommuGroupNumber     = 42
	driverName           = "i40e"
)

var (
	pciDevicesPath = filepath.Join(os.TempDir(), sriovTestDir, "devices")
	pciDriversPath = filepath.Join(os.TempDir(), sriovTestDir, "drivers")
	iommuPath      = filepath.Join(os.TempDir(), sriovTestDir, "iommu_groups")

	devicePath            = filepath.Join(pciDevicesPath, pciAddr)
	iommuGroupPath        = filepath.Join(iommuPath, strconv.Itoa(iommuGroupNumber))
	iommuGroupDevicesPath = filepath.Join(iommuGroupPath, "devices")
	driverPath            = filepath.Join(pciDriversPath, driverName)

	configuredVfPath = filepath.Join(devicePath, configuredVfFile)
)

func Test_IsDeviceSriovCapable(t *testing.T) {
	u := utils.NewSriovProvider(pciDevicesPath, pciDriversPath, iommuPath)

	err := os.RemoveAll(devicePath)
	assert.Nil(t, err)

	capable := u.IsDeviceSriovCapable(context.Background(), pciAddr)
	assert.False(t, capable)

	err = os.MkdirAll(devicePath, 0750)
	assert.Nil(t, err)

	_, err = os.Create(filepath.Join(devicePath, totalVfFile))
	assert.Nil(t, err)

	capable = u.IsDeviceSriovCapable(context.Background(), pciAddr)
	assert.True(t, capable)

	err = os.RemoveAll(devicePath)
	assert.Nil(t, err)
}

func Test_IsSriovVirtualFunction(t *testing.T) {
	u := utils.NewSriovProvider(pciDevicesPath, pciDriversPath, iommuPath)

	err := os.RemoveAll(devicePath)
	assert.Nil(t, err)

	capable := u.IsSriovVirtualFunction(context.Background(), pciAddr)
	assert.False(t, capable)

	err = os.MkdirAll(devicePath, 0750)
	assert.Nil(t, err)

	_, err = os.Create(filepath.Join(devicePath, physicalFunctionPath))
	assert.Nil(t, err)

	capable = u.IsSriovVirtualFunction(context.Background(), pciAddr)
	assert.True(t, capable)

	err = os.RemoveAll(devicePath)
	assert.Nil(t, err)
}

func Test_GetConfiguredVirtualFunctionsNumber(t *testing.T) {
	u := utils.NewSriovProvider(pciDevicesPath, pciDriversPath, iommuPath)

	err := os.RemoveAll(devicePath)
	assert.Nil(t, err)

	_, err = u.GetConfiguredVirtualFunctionsNumber(context.Background(), pciAddr)
	assert.NotNil(t, err)

	err = os.MkdirAll(devicePath, 0750)
	assert.Nil(t, err)

	_, err = os.Create(configuredVfPath)
	assert.Nil(t, err)

	err = ioutil.WriteFile(configuredVfPath, []byte("invalid number"), 0600)
	assert.Nil(t, err)

	_, err = u.GetConfiguredVirtualFunctionsNumber(context.Background(), pciAddr)
	assert.NotNil(t, err)

	numVfs := 7
	err = ioutil.WriteFile(configuredVfPath, []byte(strconv.FormatInt(int64(numVfs), 10)), 0600)
	assert.Nil(t, err)

	gotNumVfs, err := u.GetConfiguredVirtualFunctionsNumber(context.Background(), pciAddr)
	assert.Nil(t, err)
	assert.Equal(t, numVfs, gotNumVfs)

	err = os.RemoveAll(devicePath)
	assert.Nil(t, err)
}

func Test_IsSriovConfigured(t *testing.T) {
	u := utils.NewSriovProvider(pciDevicesPath, pciDriversPath, iommuPath)

	err := os.RemoveAll(devicePath)
	assert.Nil(t, err)

	conf := u.IsSriovConfigured(context.Background(), pciAddr)
	assert.False(t, conf)

	err = os.MkdirAll(devicePath, 0750)
	assert.Nil(t, err)

	_, err = os.Create(configuredVfPath)
	assert.Nil(t, err)

	err = ioutil.WriteFile(configuredVfPath, []byte("invalid number"), 0600)
	assert.Nil(t, err)

	conf = u.IsSriovConfigured(context.Background(), pciAddr)
	assert.False(t, conf)

	numVfs := 7
	err = ioutil.WriteFile(configuredVfPath, []byte(strconv.FormatInt(int64(numVfs), 10)), 0600)
	assert.Nil(t, err)

	conf = u.IsSriovConfigured(context.Background(), pciAddr)
	assert.True(t, conf)

	err = os.RemoveAll(devicePath)
	assert.Nil(t, err)
}

func Test_GetSriovVirtualFunctionsCapacity(t *testing.T) {
	u := utils.NewSriovProvider(pciDevicesPath, pciDriversPath, iommuPath)

	err := os.RemoveAll(devicePath)
	assert.Nil(t, err)

	_, err = u.GetSriovVirtualFunctionsCapacity(context.Background(), pciAddr)
	assert.NotNil(t, err)

	err = os.MkdirAll(devicePath, 0750)
	assert.Nil(t, err)

	totalVfPath := filepath.Join(devicePath, totalVfFile)
	_, err = os.Create(totalVfPath)
	assert.Nil(t, err)

	err = ioutil.WriteFile(totalVfPath, []byte("invalid number"), os.ModePerm)
	assert.Nil(t, err)

	_, err = u.GetSriovVirtualFunctionsCapacity(context.Background(), pciAddr)
	assert.NotNil(t, err)

	numVfs := 7
	err = ioutil.WriteFile(totalVfPath, []byte(strconv.FormatInt(int64(numVfs), 10)), 0600)
	assert.Nil(t, err)

	gotNumVfs, err := u.GetSriovVirtualFunctionsCapacity(context.Background(), pciAddr)
	assert.Nil(t, err)
	assert.Equal(t, numVfs, gotNumVfs)

	err = os.RemoveAll(devicePath)
	assert.Nil(t, err)
}

func Test_IsDeviceExists(t *testing.T) {
	u := utils.NewSriovProvider(pciDevicesPath, pciDriversPath, iommuPath)

	err := os.RemoveAll(devicePath)
	assert.Nil(t, err)

	_, err = u.IsDeviceExists(context.Background(), "invalid PCI address")
	assert.NotNil(t, err)

	exists, err := u.IsDeviceExists(context.Background(), pciAddr)
	assert.Nil(t, err)
	assert.False(t, exists)

	exists, err = u.IsDeviceExists(context.Background(), pciAddrShortForm)
	assert.Nil(t, err)
	assert.False(t, exists)

	err = os.MkdirAll(devicePath, 0750)
	assert.Nil(t, err)

	exists, err = u.IsDeviceExists(context.Background(), pciAddr)
	assert.Nil(t, err)
	assert.True(t, exists)

	exists, err = u.IsDeviceExists(context.Background(), pciAddrShortForm)
	assert.Nil(t, err)
	assert.True(t, exists)

	err = os.RemoveAll(devicePath)
	assert.Nil(t, err)
}

func Test_GetNetInterfacesNames(t *testing.T) {
	u := utils.NewSriovProvider(pciDevicesPath, pciDriversPath, iommuPath)

	err := os.RemoveAll(devicePath)
	assert.Nil(t, err)

	_, err = u.GetNetInterfacesNames(context.Background(), pciAddr)
	assert.NotNil(t, err)

	netIfacesPath := filepath.Join(devicePath, netInterfacesPath)
	err = os.MkdirAll(netIfacesPath, 0750)
	assert.Nil(t, err)

	netIfaces, err := u.GetNetInterfacesNames(context.Background(), pciAddr)
	assert.Nil(t, err)
	assert.Empty(t, netIfaces)

	iface1 := "enp1s0"
	_, err = os.Create(filepath.Join(netIfacesPath, iface1))
	assert.Nil(t, err)

	netIfaces, err = u.GetNetInterfacesNames(context.Background(), pciAddr)
	assert.Nil(t, err)
	assert.Equal(t, []string{iface1}, netIfaces)

	iface2 := "wlp2s0"
	_, err = os.Create(filepath.Join(netIfacesPath, iface2))
	assert.Nil(t, err)

	netIfaces, err = u.GetNetInterfacesNames(context.Background(), pciAddr)
	assert.Nil(t, err)
	// is this array network interfaces are sorted alphabetically by their names
	assert.Equal(t, []string{iface1, iface2}, netIfaces)

	err = os.RemoveAll(devicePath)
	assert.Nil(t, err)
}

func Test_CreateVirtualFunctions(t *testing.T) {
	u := utils.NewSriovProvider(pciDevicesPath, pciDriversPath, iommuPath)

	err := os.RemoveAll(devicePath)
	assert.Nil(t, err)

	err = u.CreateVirtualFunctions(context.Background(), pciAddr, -123)
	assert.NotNil(t, err)

	err = u.CreateVirtualFunctions(context.Background(), pciAddr, 0)
	assert.NotNil(t, err)

	numVfs := 7
	err = u.CreateVirtualFunctions(context.Background(), pciAddr, numVfs)
	assert.NotNil(t, err)

	err = os.MkdirAll(devicePath, 0750)
	assert.Nil(t, err)

	err = u.CreateVirtualFunctions(context.Background(), pciAddr, numVfs)
	assert.Nil(t, err)

	gotVfs, err := ioutil.ReadFile(filepath.Clean(configuredVfPath))
	assert.Nil(t, err)
	gotVfs = bytes.TrimSpace(gotVfs)
	gotNumVfs, err := strconv.Atoi(string(gotVfs))
	assert.Nil(t, err)
	assert.Equal(t, numVfs, gotNumVfs)

	err = u.CreateVirtualFunctions(context.Background(), pciAddr, 15)
	assert.NotNil(t, err)

	err = os.RemoveAll(devicePath)
	assert.Nil(t, err)
}

func Test_GetVirtualFunctionsList(t *testing.T) {
	u := utils.NewSriovProvider(pciDevicesPath, pciDriversPath, iommuPath)

	err := os.RemoveAll(devicePath)
	assert.Nil(t, err)

	_, err = u.GetVirtualFunctionsList(context.Background(), pciAddr)
	assert.NotNil(t, err)

	err = os.MkdirAll(devicePath, 0750)
	assert.Nil(t, err)

	vfs, err := u.GetVirtualFunctionsList(context.Background(), pciAddr)
	assert.Nil(t, err)
	assert.Empty(t, vfs)

	vf1PciAddr := "0000:01:00:1"
	vf1Path := filepath.Join(pciDevicesPath, vf1PciAddr)
	err = os.Mkdir(vf1Path, 0750)
	assert.Nil(t, err)
	err = os.Symlink(vf1Path, filepath.Join(devicePath, "virtfn1"))
	assert.Nil(t, err)

	vfs, err = u.GetVirtualFunctionsList(context.Background(), pciAddr)
	assert.Nil(t, err)
	assert.Equal(t, []string{vf1PciAddr}, vfs)

	vf2PciAddr := "0000:01:00:2"
	vf2Path := filepath.Join(pciDevicesPath, vf2PciAddr)
	err = os.Mkdir(vf2Path, 0750)
	assert.Nil(t, err)
	err = os.Symlink(vf2Path, filepath.Join(devicePath, "virtfn2"))
	assert.Nil(t, err)

	vfs, err = u.GetVirtualFunctionsList(context.Background(), pciAddr)
	assert.Nil(t, err)
	assert.Equal(t, []string{vf1PciAddr, vf2PciAddr}, vfs)

	err = os.RemoveAll(devicePath)
	assert.Nil(t, err)
	err = os.RemoveAll(vf1Path)
	assert.Nil(t, err)
	err = os.RemoveAll(vf2Path)
	assert.Nil(t, err)
}

func Test_GetIommuGroupNumber(t *testing.T) {
	u := utils.NewSriovProvider(pciDevicesPath, pciDriversPath, iommuPath)

	err := os.RemoveAll(devicePath)
	assert.Nil(t, err)
	err = os.RemoveAll(iommuGroupPath)
	assert.Nil(t, err)

	_, err = u.GetIommuGroupNumber(context.Background(), pciAddr)
	assert.NotNil(t, err)

	err = os.MkdirAll(devicePath, 0750)
	assert.Nil(t, err)

	_, err = u.GetIommuGroupNumber(context.Background(), pciAddr)
	assert.NotNil(t, err)

	err = os.MkdirAll(iommuGroupPath, 0750)
	assert.Nil(t, err)
	err = os.Symlink(iommuGroupPath, filepath.Join(devicePath, "iommu_group"))
	assert.Nil(t, err)

	groupNumber, err := u.GetIommuGroupNumber(context.Background(), pciAddr)
	assert.Nil(t, err)
	assert.Equal(t, iommuGroupNumber, groupNumber)

	err = os.RemoveAll(devicePath)
	assert.Nil(t, err)
	err = os.RemoveAll(iommuGroupPath)
	assert.Nil(t, err)
}

func Test_GetIommuGroupDevices(t *testing.T) {
	u := utils.NewSriovProvider(pciDevicesPath, pciDriversPath, iommuPath)

	err := os.RemoveAll(iommuGroupPath)
	assert.Nil(t, err)

	_, err = u.GetIommuGroupDevices(context.Background(), iommuGroupNumber)
	assert.NotNil(t, err)

	err = os.MkdirAll(iommuGroupDevicesPath, 0750)
	assert.Nil(t, err)

	devices, err := u.GetIommuGroupDevices(context.Background(), iommuGroupNumber)
	assert.Nil(t, err)
	assert.Empty(t, devices)

	// can use any existing location instead of sriovTestDir, link won't be evaluated
	err = os.Symlink(sriovTestDir, filepath.Join(iommuGroupDevicesPath, pciAddr))
	assert.Nil(t, err)

	devices, err = u.GetIommuGroupDevices(context.Background(), iommuGroupNumber)
	assert.Nil(t, err)
	assert.Equal(t, []string{pciAddr}, devices)

	err = os.Symlink(sriovTestDir, filepath.Join(iommuGroupDevicesPath, pciAddr2))
	assert.Nil(t, err)

	devices, err = u.GetIommuGroupDevices(context.Background(), iommuGroupNumber)
	assert.Nil(t, err)
	assert.Equal(t, []string{pciAddr, pciAddr2}, devices)

	err = os.RemoveAll(iommuGroupPath)
	assert.Nil(t, err)
}

func Test_GetBoundDriver(t *testing.T) {
	u := utils.NewSriovProvider(pciDevicesPath, pciDriversPath, iommuPath)

	err := os.RemoveAll(devicePath)
	assert.Nil(t, err)
	err = os.RemoveAll(driverPath)
	assert.Nil(t, err)

	_, err = u.GetBoundDriver(context.Background(), pciAddr)
	assert.NotNil(t, err)

	err = os.MkdirAll(devicePath, 0750)
	assert.Nil(t, err)

	driver, err := u.GetBoundDriver(context.Background(), pciAddr)
	assert.Nil(t, err)
	assert.Equal(t, "", driver)

	err = os.MkdirAll(driverPath, 0750)
	assert.Nil(t, err)
	err = os.Symlink(driverPath, filepath.Join(devicePath, "driver"))
	assert.Nil(t, err)

	driver, err = u.GetBoundDriver(context.Background(), pciAddr)
	assert.Nil(t, err)
	assert.Equal(t, driverName, driver)

	err = os.RemoveAll(devicePath)
	assert.Nil(t, err)
	err = os.RemoveAll(driverPath)
	assert.Nil(t, err)
}

func Test_UnbindDriver(t *testing.T) {
	u := utils.NewSriovProvider(pciDevicesPath, pciDriversPath, iommuPath)

	err := os.RemoveAll(devicePath)
	assert.Nil(t, err)
	err = os.RemoveAll(driverPath)
	assert.Nil(t, err)

	err = u.UnbindDriver(context.Background(), pciAddr)
	assert.NotNil(t, err)

	err = os.MkdirAll(devicePath, 0750)
	assert.Nil(t, err)

	err = u.UnbindDriver(context.Background(), pciAddr)
	assert.Nil(t, err)

	err = os.MkdirAll(driverPath, 0750)
	assert.Nil(t, err)
	err = os.Symlink(driverPath, filepath.Join(devicePath, "driver"))
	assert.Nil(t, err)

	err = u.UnbindDriver(context.Background(), pciAddr)
	assert.Nil(t, err)

	addrBytes, err := ioutil.ReadFile(filepath.Clean(filepath.Join(driverPath, "unbind")))
	assert.Nil(t, err)
	addr := string(bytes.TrimSpace(addrBytes))
	assert.Equal(t, pciAddr, addr)

	err = os.RemoveAll(devicePath)
	assert.Nil(t, err)
	err = os.RemoveAll(driverPath)
	assert.Nil(t, err)
}

func Test_BindDriver(t *testing.T) {
	u := utils.NewSriovProvider(pciDevicesPath, pciDriversPath, iommuPath)

	err := os.RemoveAll(driverPath)
	assert.Nil(t, err)
	err = os.RemoveAll(devicePath)
	assert.Nil(t, err)

	err = u.BindDriver(context.Background(), pciAddr, driverName)
	assert.NotNil(t, err)

	err = os.MkdirAll(devicePath, 0750)
	assert.Nil(t, err)

	err = u.BindDriver(context.Background(), pciAddr, driverName)
	assert.NotNil(t, err)

	err = os.MkdirAll(driverPath, 0750)
	assert.Nil(t, err)

	err = u.BindDriver(context.Background(), pciAddr, driverName)
	assert.Nil(t, err)

	addrBytes, err := ioutil.ReadFile(filepath.Clean(filepath.Join(driverPath, "bind")))
	assert.Nil(t, err)
	addr := string(bytes.TrimSpace(addrBytes))
	assert.Equal(t, pciAddr, addr)

	err = os.RemoveAll(driverPath)
	assert.Nil(t, err)
	err = os.RemoveAll(devicePath)
	assert.Nil(t, err)
}

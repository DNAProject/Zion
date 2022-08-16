/*
 * Copyright (C) 2021 The Zion Authors
 * This file is part of The Zion library.
 *
 * The Zion is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Lesser General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * The Zion is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU Lesser General Public License for more details.
 *
 * You should have received a copy of the GNU Lesser General Public License
 * along with The Zion.  If not, see <http://www.gnu.org/licenses/>.
 */

package maas_config

import (
	"fmt"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/contracts/native/utils"
)

const contractName = "maas config"

const (

	// abi
	MaasConfigABI = "[{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"bool\",\"name\":\"doBlock\",\"type\":\"bool\"}],\"name\":\"BlockAccount\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"oldOwner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"ChangeOwner\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bool\",\"name\":\"doEnable\",\"type\":\"bool\"}],\"name\":\"EnableGasManage\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address[]\",\"name\":\"addrs\",\"type\":\"address[]\"},{\"indexed\":false,\"internalType\":\"bool\",\"name\":\"addOrRemove\",\"type\":\"bool\"}],\"name\":\"SetAdmins\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"bool\",\"name\":\"isManager\",\"type\":\"bool\"}],\"name\":\"SetGasManager\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address[]\",\"name\":\"addrs\",\"type\":\"address[]\"},{\"indexed\":false,\"internalType\":\"bool\",\"name\":\"addOrRemove\",\"type\":\"bool\"}],\"name\":\"SetGasUsers\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"},{\"internalType\":\"bool\",\"name\":\"doBlock\",\"type\":\"bool\"}],\"name\":\"blockAccount\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"}],\"name\":\"changeOwner\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bool\",\"name\":\"doEnable\",\"type\":\"bool\"}],\"name\":\"enableGasManage\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getAdminList\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getBlacklist\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getGasManagerList\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getGasUserList\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getOwner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"}],\"name\":\"isAdmin\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"}],\"name\":\"isBlocked\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"isGasManageEnabled\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"}],\"name\":\"isGasManager\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"}],\"name\":\"isGasUser\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"name\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address[]\",\"name\":\"addrs\",\"type\":\"address[]\"},{\"internalType\":\"bool\",\"name\":\"addOrRemove\",\"type\":\"bool\"}],\"name\":\"setAdmins\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"},{\"internalType\":\"bool\",\"name\":\"isManager\",\"type\":\"bool\"}],\"name\":\"setGasManager\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address[]\",\"name\":\"addrs\",\"type\":\"address[]\"},{\"internalType\":\"bool\",\"name\":\"addOrRemove\",\"type\":\"bool\"}],\"name\":\"setGasUsers\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]"

	// method name
	MethodName         = "name"
	MethodChangeOwner  = "changeOwner"
	MethodGetOwner     = "getOwner"
	MethodBlockAccount = "blockAccount"
	MethodIsBlocked    = "isBlocked"
	MethodGetBlacklist = "getBlacklist"

	MethodEnableGasManage    = "enableGasManage"
	MethodIsGasManageEnabled = "isGasManageEnabled"
	MethodSetGasManager      = "setGasManager"
	MethodIsGasManager       = "isGasManager"
	MethodGetGasManagerList  = "getGasManagerList"

	MethodSetGasUsers    = "setGasUsers"
	MethodIsGasUser      = "isGasUser"
	MethodGetGasUserList = "getGasUserList"

	MethodSetAdmins    = "setAdmins"
	MethodIsAdmin      = "isAdmin"
	MethodGetAdminList = "getAdminList"

	EventChangeOwner     = "ChangeOwner"
	EventBlockAccount    = "BlockAccount"
	EventEnableGasManage = "EnableGasManage"
	EventSetGasManager   = "SetGasManager"
	EventSetGasUsers     = "SetGasUsers"
	EventSetAdmins       = "SetAdmins"
)

func InitABI() {
	ab, err := abi.JSON(strings.NewReader(MaasConfigABI))
	if err != nil {
		panic(fmt.Sprintf("failed to load abi json string: [%v]", err))
	}
	ABI = &ab
}

var (
	ABI  *abi.ABI
	this = utils.MaasConfigContractAddress
)

type MethodContractNameOutput struct {
	Name string
}

func (m *MethodContractNameOutput) Encode() ([]byte, error) {
	m.Name = contractName
	return utils.PackOutputs(ABI, MethodName, m.Name)
}
func (m *MethodContractNameOutput) Decode(payload []byte) error {
	return utils.UnpackOutputs(ABI, MethodName, m, payload)
}

type MethodBoolOutput struct {
	Success bool
}

func (m *MethodBoolOutput) Encode(methodName string) ([]byte, error) {
	return utils.PackOutputs(ABI, methodName, m.Success)
}

func (m *MethodBoolOutput) Decode(payload []byte, methodName string) error {
	return utils.UnpackOutputs(ABI, methodName, m, payload)
}

type MethodAddressOutput struct {
	Addr common.Address
}

func (m *MethodAddressOutput) Encode(methodName string) ([]byte, error) {
	return utils.PackOutputs(ABI, methodName, m.Addr)
}

func (m *MethodAddressOutput) Decode(payload []byte, methodName string) error {
	return utils.UnpackOutputs(ABI, methodName, m, payload)
}

type MethodChangeOwnerInput struct {
	Addr common.Address
}

func (m *MethodChangeOwnerInput) Encode() ([]byte, error) {
	return utils.PackMethod(ABI, MethodChangeOwner, m.Addr)
}

func (m *MethodChangeOwnerInput) Decode(payload []byte) error {
	var data struct {
		Addr common.Address
	}
	if err := utils.UnpackMethod(ABI, MethodChangeOwner, &data, payload); err != nil {
		return err
	}
	m.Addr = data.Addr
	return nil
}

type MethodBlockAccountInput struct {
	Addr    common.Address
	DoBlock bool
}

func (m *MethodBlockAccountInput) Encode() ([]byte, error) {
	return utils.PackMethod(ABI, MethodBlockAccount, m.Addr, m.DoBlock)
}

func (m *MethodBlockAccountInput) Decode(payload []byte) error {
	var data struct {
		Addr    common.Address
		DoBlock bool
	}
	if err := utils.UnpackMethod(ABI, MethodBlockAccount, &data, payload); err != nil {
		return err
	}
	m.Addr = data.Addr
	m.DoBlock = data.DoBlock
	return nil
}

type MethodIsBlockedInput struct {
	Addr common.Address
}

func (m *MethodIsBlockedInput) Encode() ([]byte, error) {
	return utils.PackMethod(ABI, MethodIsBlocked, m.Addr)
}

func (m *MethodIsBlockedInput) Decode(payload []byte) error {
	var data struct {
		Addr common.Address
	}
	if err := utils.UnpackMethod(ABI, MethodIsBlocked, &data, payload); err != nil {
		return err
	}
	m.Addr = data.Addr
	return nil
}

type MethodStringOutput struct {
	Result string
}

func (m *MethodStringOutput) Encode(methodName string) ([]byte, error) {
	return utils.PackOutputs(ABI, methodName, m.Result)
}

func (m *MethodStringOutput) Decode(payload []byte, methodName string) error {
	return utils.UnpackOutputs(ABI, methodName, m, payload)
}

type MethodEnableGasManageInput struct {
	DoEnable bool
}

func (m *MethodEnableGasManageInput) Encode() ([]byte, error) {
	return utils.PackMethod(ABI, MethodEnableGasManage, m.DoEnable)
}

func (m *MethodEnableGasManageInput) Decode(payload []byte) error {
	return utils.UnpackMethod(ABI, MethodEnableGasManage, m, payload)
}

type MethodSetGasManagerInput struct {
	Addr      common.Address
	IsManager bool
}

func (m *MethodSetGasManagerInput) Encode() ([]byte, error) {
	return utils.PackMethod(ABI, MethodSetGasManager, m.Addr, m.IsManager)
}

func (m *MethodSetGasManagerInput) Decode(payload []byte) error {
	return utils.UnpackMethod(ABI, MethodSetGasManager, m, payload)
}

type MethodIsGasManagerInput struct {
	Addr common.Address
}

func (m *MethodIsGasManagerInput) Encode() ([]byte, error) {
	return utils.PackMethod(ABI, MethodIsGasManager, m.Addr)
}

func (m *MethodIsGasManagerInput) Decode(payload []byte) error {
	return utils.UnpackMethod(ABI, MethodIsGasManager, m, payload)
}

type MethodSetGasUsersInput struct {
	Addrs       []common.Address
	AddOrRemove bool
}

func (m *MethodSetGasUsersInput) Encode() ([]byte, error) {
	return utils.PackMethod(ABI, MethodSetGasUsers, m.Addrs, m.AddOrRemove)
}

func (m *MethodSetGasUsersInput) Decode(payload []byte) error {
	return utils.UnpackMethod(ABI, MethodSetGasUsers, m, payload)
}

type MethodIsGasUserInput struct {
	Addr common.Address
}

func (m *MethodIsGasUserInput) Encode() ([]byte, error) {
	return utils.PackMethod(ABI, MethodIsGasUser, m.Addr)
}

func (m *MethodIsGasUserInput) Decode(payload []byte) error {
	return utils.UnpackMethod(ABI, MethodIsGasUser, m, payload)
}

type MethodSetAdminsInput struct {
	Addrs       []common.Address
	AddOrRemove bool
}

func (m *MethodSetAdminsInput) Encode() ([]byte, error) {
	return utils.PackMethod(ABI, MethodSetAdmins, m.Addrs, m.AddOrRemove)
}

func (m *MethodSetAdminsInput) Decode(payload []byte) error {
	return utils.UnpackMethod(ABI, MethodSetAdmins, m, payload)
}

type MethodIsAdminInput struct {
	Addr common.Address
}

func (m *MethodIsAdminInput) Encode() ([]byte, error) {
	return utils.PackMethod(ABI, MethodIsAdmin, m.Addr)
}

func (m *MethodIsAdminInput) Decode(payload []byte) error {
	return utils.UnpackMethod(ABI, MethodIsAdmin, m, payload)
}

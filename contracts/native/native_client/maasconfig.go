// Copyright 2016 The go-ethereum Authors
// This file is part of the go-ethereum library.
//
// The go-ethereum library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-ethereum library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.

package native_client

import (
	"errors"
	"fmt"
	"math/big"
	"strconv"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/contracts/native"
	"github.com/ethereum/go-ethereum/contracts/native/utils"
	"github.com/ethereum/go-ethereum/core/state"
	"github.com/ethereum/go-ethereum/log"
)

// MaasConfigABI is the ABI used to call native contract maasConfig
const MaasConfigABI = "[{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"bool\",\"name\":\"doBlock\",\"type\":\"bool\"}],\"name\":\"BlockAccount\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"},{\"internalType\":\"bool\",\"name\":\"doBlock\",\"type\":\"bool\"}],\"name\":\"blockAccount\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"}],\"name\":\"isBlocked\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"name\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]"
const MethodIsBlocked = "isBlocked"

var ErrAccountBlocked = errors.New("account is in blacklist")

func IsBlocked(state *state.StateDB, address *common.Address) bool {
	log.Debug("### isBlocked called")
	if address == nil {
		return false
	}
	caller := common.EmptyAddress
	ref := native.NewContractRef(state, caller, caller, big.NewInt(-1), common.EmptyHash, 0, nil)

	ab, err := abi.JSON(strings.NewReader(MaasConfigABI))
	if err != nil {
		panic(fmt.Sprintf("failed to load abi json string: [%v]", err))
	}
	ABI := &ab

	payload, err := utils.PackMethod(ABI, MethodIsBlocked, address) // (&maas_config.MethodIsBlockedInput{Addr: *address}).Encode()
	if err != nil {
		log.Error("[PackMethod]", "pack `isBlocked` input failed", err)
		return false
	}
	enc, _, err := ref.NativeCall(caller, utils.MaasConfigContractAddress, payload)
	if err != nil {
		return false
	}
	var data struct {
		Result bool
	}
	utils.UnpackOutputs(ABI, MethodIsBlocked, &data, enc)

	log.Debug("IsBlocked: " + address.String() + ", " + strconv.FormatBool(data.Result))
	return data.Result
}

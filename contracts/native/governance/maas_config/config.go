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
	"errors"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/contracts/native"
	"github.com/ethereum/go-ethereum/contracts/native/contract"
	"github.com/ethereum/go-ethereum/contracts/native/utils"
	"github.com/ethereum/go-ethereum/log"
)

var (
	gasTable = map[string]uint64{
		MethodName:         0,
		MethodChangeOwner:  30000,
		MethodGetOwner:     0,
		MethodBlockAccount: 30000,
		MethodIsBlocked:    0,
	}
	// initial owner
	owner = common.HexToAddress("0x2D3913c12ACa0E4A2278f829Fb78A682123c0125")
)

func InitMaasConfig() {
	InitABI()
	native.Contracts[this] = RegisterMaasConfigContract
}

func RegisterMaasConfigContract(s *native.NativeContract) {
	s.Prepare(ABI, gasTable)

	s.Register(MethodName, Name)
	s.Register(MethodChangeOwner, ChangeOwner)
	s.Register(MethodGetOwner, GetOwner)
	s.Register(MethodBlockAccount, BlockAccount)
	s.Register(MethodIsBlocked, IsBlocked)
}

func Name(s *native.NativeContract) ([]byte, error) {
	return new(MethodContractNameOutput).Encode()
}

// change owner
func ChangeOwner(s *native.NativeContract) ([]byte, error) {
	ctx := s.ContractRef().CurrentContext()
	caller := ctx.Caller

	// check authority
	if err := contract.ValidateOwner(s, caller); err != nil {
		log.Trace("blockAccount", "ValidateOwner caller failed", err)
		return utils.ByteFailed, errors.New("invalid authority for caller")
	}

	currentOwner := getOwner(s)
	if err := contract.ValidateOwner(s, currentOwner); err != nil {
		log.Trace("blockAccount", "ValidateOwner owner failed", err)
		return utils.ByteFailed, errors.New("invalid authority for owner")
	}

	// decode input
	input := new(MethodChangeOwnerInput)
	if err := input.Decode(ctx.Payload); err != nil {
		log.Trace("ChangeOwner", "decode input failed", err)
		return utils.ByteFailed, errors.New("invalid input")
	}

	// store blacklist
	key := getOwnerKey()
	set(s, key, input.Addr.Bytes())

	// emit event log
	if err := s.AddNotify(ABI, []string{EventChangeOwner}, currentOwner, input.Addr); err != nil {
		log.Trace("propose", "emit event log failed", err)
		return utils.ByteFailed, errors.New("emit EventChangeOwner error")
	}

	log.Debug("ChangeOwner: " + input.Addr.String())
	return utils.ByteSuccess, nil
}

// get owner
func GetOwner(s *native.NativeContract) ([]byte, error) {
	output := &MethodGetOwnerOutput{Addr: getOwner(s)}
	return output.Encode()
}

func getOwner(s *native.NativeContract) common.Address {
	// get value
	key := getOwnerKey()
	value, _ := get(s, key)
	if len(value) == 0 {
		return owner
	}
	return common.BytesToAddress(value)
}

// block account(add account to blacklist map) or unblock account
func BlockAccount(s *native.NativeContract) ([]byte, error) {
	ctx := s.ContractRef().CurrentContext()
	caller := ctx.Caller

	// check authority
	if err := contract.ValidateOwner(s, caller); err != nil {
		log.Trace("blockAccount", "ValidateOwner caller failed", err)
		return utils.ByteFailed, errors.New("invalid authority for caller")
	}

	currentOwner := getOwner(s)
	if err := contract.ValidateOwner(s, currentOwner); err != nil {
		log.Trace("blockAccount", "ValidateOwner owner failed", err)
		return utils.ByteFailed, errors.New("invalid authority for owner")
	}

	// decode input
	input := new(MethodBlockAccountInput)
	if err := input.Decode(ctx.Payload); err != nil {
		log.Trace("blockAccount", "decode input failed", err)
		return utils.ByteFailed, errors.New("invalid input")
	}

	// store blacklist
	key := getBlacklistKey(input.Addr)
	if input.DoBlock {
		set(s, key, utils.ByteSuccess)
	} else {
		del(s, key)
	}

	// emit event log
	if err := s.AddNotify(ABI, []string{EventBlockAccount}, input.Addr, input.DoBlock); err != nil {
		log.Trace("propose", "emit event log failed", err)
		return utils.ByteFailed, errors.New("emitBlockAccount error")
	}

	log.Debug("BlockAccount: "+input.Addr.String(), input.DoBlock)
	return utils.ByteSuccess, nil
}

func getBlacklistKey(addr common.Address) []byte {
	return utils.ConcatKey(this, []byte(BLACKLIST), addr.Bytes())
}

func getOwnerKey() []byte {
	return utils.ConcatKey(this, []byte(OWNER))
}

// check if account is blocked
func IsBlocked(s *native.NativeContract) ([]byte, error) {
	ctx := s.ContractRef().CurrentContext()

	// decode input
	input := new(MethodIsBlockedInput)
	if err := input.Decode(ctx.Payload); err != nil {
		log.Trace("IsBlocked", "decode input failed", err)
		return utils.ByteFailed, errors.New("invalid input")
	}

	// get value
	key := getBlacklistKey(input.Addr)
	value, _ := get(s, key)
	output := &MethodIsBlockedOutput{Success: len(value) > 0}
	return output.Encode()
}

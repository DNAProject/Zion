package maas_config

import (
	"crypto/rand"
	"errors"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/contracts/native"
	"github.com/ethereum/go-ethereum/core/rawdb"
	"github.com/ethereum/go-ethereum/core/state"
	"github.com/stretchr/testify/assert"
	"math/big"
	"os"
	"strconv"
	"testing"
)

func TestMain(m *testing.M) {
	InitMaasConfig()
	os.Exit(m.Run())
}

var (
	testStateDB  *state.StateDB
	testEmptyCtx *native.NativeContract

	testSupplyGas    uint64 = 100000000000000000
	testGenesisNum   int    = 4
	testCaller       common.Address
	//testGenesisEpoch *EpochInfo
)


func generateNativeContractRef(origin common.Address, blockNum int) *native.ContractRef {
	token := make([]byte, common.HashLength)
	rand.Read(token)
	hash := common.BytesToHash(token)
	return native.NewContractRef(testStateDB, origin, origin, big.NewInt(int64(blockNum)), hash, testSupplyGas, nil)
}

func generateNativeContract(origin common.Address, blockNum int) *native.NativeContract {
	ref := generateNativeContractRef(origin, blockNum)
	return native.NewNativeContract(testStateDB, ref)
}

func resetTestContext() {
	db := rawdb.NewMemoryDatabase()
	testStateDB, _ = state.New(common.Hash{}, state.NewDatabase(db), nil)
	testEmptyCtx = native.NewNativeContract(testStateDB, nil)
	//testGenesisPeers := generateTestPeers(testGenesisNum)
	//testGenesisEpoch, _ = storeGenesisEpoch(testStateDB, testGenesisPeers)
	testCaller = testAddresses[0]
}

func TestChangeOwner(t *testing.T) {
	type TestCase struct {
		BlockNum      int
		StartHeight   uint64
		Payload       []byte
		BeforeHandler func(c *TestCase, ctx *native.NativeContract)
		AfterHandler  func(c *TestCase, ctx *native.NativeContract)
		Index         int
		Expect        error
		ReturnData	  bool
	}

	cases := []*TestCase{
		{
			BlockNum:    3,
			StartHeight: 2,
			Index:       1,
			BeforeHandler: func(c *TestCase, ctx *native.NativeContract) {
				input := &MethodChangeOwnerInput{Addr: common.HexToAddress("0x1231231")}
				c.Payload, _ = input.Encode()
			},
			Expect: nil,
			ReturnData: true,
		},
	}

	for _, v := range cases {
		resetTestContext()
		ctx := generateNativeContract(testCaller, v.BlockNum)
		if v.BeforeHandler != nil {
			v.BeforeHandler(v, ctx)
		}
		result, _, err := ctx.ContractRef().NativeCall(testCaller, this, v.Payload)
		assert.Equal(t, v.Expect, err)
		res, err := strconv.ParseBool(string(result))
		assert.NoError(t, err)
		t.Log("change owner result", res)
		assert.Equal(t, v.ReturnData, res)
		if v.AfterHandler != nil {
			v.AfterHandler(v, ctx)
		}
	}
}

func TestGetOwner(t *testing.T) {
	type TestCase struct {
		BlockNum      int
		StartHeight   uint64
		Payload       []byte
		BeforeHandler func(c *TestCase, ctx *native.NativeContract)
		AfterHandler  func(c *TestCase, ctx *native.NativeContract)
		Index         int
		Expect        error
	}

	cases := []*TestCase{
		{
			BlockNum:    3,
			StartHeight: 2,
			Index:       1,
			BeforeHandler: func(c *TestCase, ctx *native.NativeContract) {
				input := &MethodGetOwnerInput{}
				payload, err := input.Encode()
				assert.NoError(t, err)
				c.Payload = payload
			},
			Expect: nil,
		},
	}

	for _, v := range cases {
		resetTestContext()
		ctx := generateNativeContract(testCaller, v.BlockNum)
		if v.BeforeHandler != nil {
			v.BeforeHandler(v, ctx)
		}
		result, _, err := ctx.ContractRef().NativeCall(testCaller, this, v.Payload)
		assert.Equal(t, v.Expect, err)
		t.Log(hexutil.Encode(result))
		if v.AfterHandler != nil {
			v.AfterHandler(v, ctx)
		}
	}

}

func TestMethodBlockAccount(t *testing.T) {
	type TestCase struct {
		BlockNum      int
		StartHeight   uint64
		Payload       []byte
		BeforeHandler func(c *TestCase, ctx *native.NativeContract)
		AfterHandler  func(c *TestCase, ctx *native.NativeContract)
		Index         int
		Expect        error
	}

	cases := []*TestCase{
		{
			BlockNum:    3,
			StartHeight: 2,
			Index:       1,
			BeforeHandler: func(c *TestCase, ctx *native.NativeContract) {
				input := &MethodBlockAccountInput{Addr: common.HexToAddress("0x123"), DoBlock: true}
				c.Payload, _ = input.Encode()
			},
			Expect: errors.New("invalid authority for owner"),
		},
	}

	for _, v := range cases {
		resetTestContext()
		ctx := generateNativeContract(testCaller, v.BlockNum)
		if v.BeforeHandler != nil {
			v.BeforeHandler(v, ctx)
		}
		result, _, err := ctx.ContractRef().NativeCall(testCaller, this, v.Payload)
		assert.Equal(t, v.Expect, err)
		t.Log(hexutil.Encode(result))
		if v.AfterHandler != nil {
			v.AfterHandler(v, ctx)
		}
	}
}

func TestMethodGetBlacklist(t *testing.T) {
	type TestCase struct {
		BlockNum      int
		StartHeight   uint64
		Payload       []byte
		BeforeHandler func(c *TestCase, ctx *native.NativeContract)
		AfterHandler  func(c *TestCase, ctx *native.NativeContract)
		Index         int
		Expect        error
	}

	cases := []*TestCase{
		{
			BlockNum:    3,
			StartHeight: 2,
			Index:       1,
			BeforeHandler: func(c *TestCase, ctx *native.NativeContract) {
				input := &MethodGetBlacklistInput{}
				c.Payload, _ = input.Encode()
			},
			Expect: nil,
		},
	}

	for _, v := range cases {
		resetTestContext()
		ctx := generateNativeContract(testCaller, v.BlockNum)
		if v.BeforeHandler != nil {
			v.BeforeHandler(v, ctx)
		}
		result, _, err := ctx.ContractRef().NativeCall(testCaller, this, v.Payload)
		assert.Equal(t, v.Expect, err)
		t.Log(hexutil.Encode(result))
		if v.AfterHandler != nil {
			v.AfterHandler(v, ctx)
		}
	}
}

func TestMethodIsBlocked(t *testing.T) {
	type TestCase struct {
		BlockNum      int
		StartHeight   uint64
		Payload       []byte
		BeforeHandler func(c *TestCase, ctx *native.NativeContract)
		AfterHandler  func(c *TestCase, ctx *native.NativeContract)
		Index         int
		Expect        error
	}

	cases := []*TestCase{
		{
			BlockNum:    3,
			StartHeight: 2,
			Index:       1,
			BeforeHandler: func(c *TestCase, ctx *native.NativeContract) {
				input := &MethodIsBlockedInput{Addr: common.HexToAddress("0x0")}
				c.Payload, _ = input.Encode()
			},
			Expect: nil,
		},
	}

	for _, v := range cases {
		resetTestContext()
		ctx := generateNativeContract(testCaller, v.BlockNum)
		if v.BeforeHandler != nil {
			v.BeforeHandler(v, ctx)
		}
		result, _, err := ctx.ContractRef().NativeCall(testCaller, this, v.Payload)
		assert.Equal(t, v.Expect, err)
		t.Log(hexutil.Encode(result))
		if v.AfterHandler != nil {
			v.AfterHandler(v, ctx)
		}
	}
}

func TestMethodName(t *testing.T){
	type TestCase struct {
		BlockNum      int
		StartHeight   uint64
		Payload       []byte
		BeforeHandler func(c *TestCase, ctx *native.NativeContract)
		AfterHandler  func(c *TestCase, ctx *native.NativeContract)
		Index         int
		Expect        error
	}

	cases := []*TestCase{
		{
			BlockNum:    3,
			StartHeight: 2,
			Index:       1,
			BeforeHandler: func(c *TestCase, ctx *native.NativeContract) {
				input := &MethodContractNameInput{}
				c.Payload, _ = input.Encode()
			},
			Expect: nil,
		},
	}

	for _, v := range cases {
		resetTestContext()
		ctx := generateNativeContract(testCaller, v.BlockNum)
		if v.BeforeHandler != nil {
			v.BeforeHandler(v, ctx)
		}
		result, _, err := ctx.ContractRef().NativeCall(testCaller, this, v.Payload)
		assert.Equal(t, v.Expect, err)
		t.Log(hexutil.Encode(result))
		if v.AfterHandler != nil {
			v.AfterHandler(v, ctx)
		}
	}
}
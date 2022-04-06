package maas_config

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/contracts/native/utils"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestABIShowJonString(t *testing.T) {
	t.Log(MaasConfigABI)
	for name, v := range ABI.Methods {
		t.Logf("method %s, id %s", name, hexutil.Encode(v.ID))
	}
}

func TestABIMethodContractName(t *testing.T) {

	enc, err := utils.PackOutputs(ABI, MethodName, contractName)
	assert.NoError(t, err)
	params := new(MethodContractNameOutput)
	assert.NoError(t, utils.UnpackOutputs(ABI, MethodName, params, enc))
	assert.Equal(t, contractName, params.Name)
}

func TestABIMethodChangeOwnerInput(t *testing.T) {
	expect := &MethodChangeOwnerInput{Addr: common.HexToAddress("0x2D3913c12ACa0E4A2278f829Fb78A682123c0125")}
	enc, err := expect.Encode()
	assert.NoError(t, err)
	methodId := hexutil.Encode(crypto.Keccak256([]byte("changeOwner(address)"))[:4])
	t.Log(methodId)
	t.Log(hexutil.Encode(enc)[:10])
	assert.Equal(t, methodId, hexutil.Encode(enc)[:10])
	got := new(MethodChangeOwnerInput)
	assert.NoError(t, got.Decode(enc))
	assert.Equal(t, expect, got)
}

//func TestABIMethodChangeOwnerOutput(t *testing.T) {
//	var cases = []struct {
//		Result bool
//	}{
//
//	}
//}
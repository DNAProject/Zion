package params

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/log"
)

var instance *bsnConfig
var mu sync.RWMutex
var once sync.Once

// BsnConfig is the config which determines the blockchain settings about ipwhitelist, address blacklist and manage accounts.
type bsnConfigFile struct {
	IpWhitelist       []string `json:"ipWhitelist"`       // only the ip in IpWhitelist can connect to this node with p2p network
	P2pManageAccounts []string `json:"p2pManageAccounts"` // transfer manage accounts
	P2pTransferEnable bool     `json:"p2pTransferEnable"` // only manage accounts can be from or to address when the P2pTransferEnable is false
	Blacklist         []string `json:"blacklist"`         // the address in blacklist can't be from or to address
}

type bsnConfig struct {
	ipWhitelist       map[string]struct{}
	p2pManageAccounts map[string]struct{}
	p2pTransferEnable bool
	blacklist         map[string]struct{}
}

func StartBsnTask(filepath string) {
	go func() {
		initialStat, err := os.Stat(filepath)
		if err != nil {
			log.Error("### initialStat error. " + err.Error())
		}

		for {
			stat, err := os.Stat(filepath)
			if err != nil {
				log.Warn("### os.Stat failed. " + err.Error())
			}

			if stat.Size() != initialStat.Size() || stat.ModTime() != initialStat.ModTime() {
				err = loadBsnConfig(filepath)
				if err != nil {
					log.Warn("### LoadBsnConfig file failed. " + err.Error())
				}
			}

			time.Sleep(time.Second)
		}
	}()
}

func loadBsnConfig(path string) error {
	config, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}
	var res *bsnConfigFile
	if err := json.Unmarshal(config, &res); err != nil {
		return err
	}

	instance = GetBsnConfig()
	mu.Lock()
	defer mu.Unlock()
	instance.ipWhitelist = convertToMap(res.IpWhitelist)
	instance.p2pManageAccounts = convertToMap(res.P2pManageAccounts)
	instance.p2pTransferEnable = res.P2pTransferEnable
	instance.blacklist = convertToMap(res.Blacklist)
	return nil
}

func convertToMap(list []string) map[string]struct{} {
	var result = make(map[string]struct{})
	for _, s := range list {
		// todo: add address format check
		// store address as lower case
		result[strings.ToLower(s)] = struct{}{}
	}
	return result
}

func GetBsnConfig() *bsnConfig {
	once.Do(func() {
		instance = &bsnConfig{p2pTransferEnable: true}
	})

	return instance
}

func (config *bsnConfig) IsInIpWhitelist(address *common.Address) bool {
	mu.RLock()
	defer mu.RUnlock()

	if address == nil || config.ipWhitelist == nil {
		return false
	}
	_, ok := config.ipWhitelist[strings.ToLower(address.String())]
	return ok
}

func (config *bsnConfig) IsInP2pManageAccounts(address *common.Address) bool {
	mu.RLock()
	defer mu.RUnlock()

	if address == nil || config.p2pManageAccounts == nil {
		return false
	}
	_, ok := config.p2pManageAccounts[strings.ToLower(address.String())]
	return ok
}

func (config *bsnConfig) IsInBlacklist(address *common.Address) bool {
	mu.RLock()
	defer mu.RUnlock()

	if address == nil || config.blacklist == nil {
		return false
	}
	_, ok := config.blacklist[strings.ToLower(address.String())]
	return ok
}

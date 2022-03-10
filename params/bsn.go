package params

import (
	"encoding/json"
	"errors"
	"io/fs"
	"io/ioutil"
	"net"
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

var (
	// ErrIpWhitelist is returned if the node ip is not in whitelist
	ErrIpWhitelist = errors.New("### BsnConfig ip address is not in whitelist")

	// ErrBlacklist is returned if a transaction's address is in blacklist
	ErrBlacklist = errors.New("### BsnConfig transaction's address is in blacklist")

	// ErrGasManage is returned if a transaction's address is not in GasManage Accounts
	ErrGasManage = errors.New("### BsnConfig transaction's address is not in GasManage Accounts")
)

// BsnConfig is the config which determines the blockchain settings about ipwhitelist, address blacklist and manage accounts.
type bsnConfigFile struct {
	IpWhitelistEnable bool     `json:"ipWhitelistEnable"` // only the ip in IpWhitelist can connect to this node with p2p network when the IpWhitelistEnable is true, default false
	IpWhitelist       []string `json:"ipWhitelist"`       // ip whitelist
	GasManageEnable   bool     `json:"gasManageEnable"`   // only manage accounts can be from or to address when the GasManageEnable is true, default false
	GasManageAccounts []string `json:"gasManageAccounts"` // gas transfer manage accounts
	BlacklistEnable   bool     `json:"blacklistEnable"`   // blacklist function is enabled when the blacklistEnable is true, default false
	Blacklist         []string `json:"blacklist"`         // the address in blacklist can't be from or to address
}

type bsnConfig struct {
	ipWhitelistEnable bool
	ipWhitelist       map[string]struct{}
	gasManageEnable   bool
	gasManageAccounts map[string]struct{}
	blacklistEnable   bool
	blacklist         map[string]struct{}
}

func StartBsnTask(filepath string) {
	go func() {
		var initialStat fs.FileInfo

		for {
			stat, err := os.Stat(filepath)
			if err != nil {
				log.Warn("### StartBsnTask os.Stat error. " + err.Error())
				time.Sleep(time.Second)
				continue
			}

			if initialStat == nil || stat.Size() != initialStat.Size() || stat.ModTime() != initialStat.ModTime() {
				err = loadBsnConfig(filepath)
				initialStat = stat
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
	instance.ipWhitelistEnable = res.IpWhitelistEnable
	instance.ipWhitelist = convertToIPMap(res.IpWhitelist)
	instance.gasManageEnable = res.GasManageEnable
	instance.gasManageAccounts = convertToAddressMap(res.GasManageAccounts)
	instance.blacklistEnable = res.BlacklistEnable
	instance.blacklist = convertToAddressMap(res.Blacklist)
	return nil
}

func convertToIPMap(list []string) map[string]struct{} {
	var result = make(map[string]struct{})
	for _, s := range list {
		// ip format check
		if net.ParseIP(s) != nil {
			result[s] = struct{}{}
		}
	}
	return result
}

func convertToAddressMap(list []string) map[string]struct{} {
	var result = make(map[string]struct{})
	for _, s := range list {
		// address format check
		if common.IsHexAddress(s) {
			// store address as lower case
			result[strings.ToLower(s)] = struct{}{}
		}
	}
	return result
}

func GetBsnConfig() *bsnConfig {
	once.Do(func() {
		instance = &bsnConfig{}
	})

	return instance
}

func (config *bsnConfig) CheckIpWhitelist(ip net.IP) bool {
	mu.RLock()
	defer mu.RUnlock()

	if !config.ipWhitelistEnable {
		return true
	}

	if ip == nil {
		return false
	}

	return isInMap(ip.String(), config.ipWhitelist)
}

func (config *bsnConfig) CheckGasManage(from *common.Address, to *common.Address) bool {
	mu.RLock()
	defer mu.RUnlock()

	if !config.gasManageEnable {
		return true
	}
	return isAddressInMap(from, config.gasManageAccounts) || isAddressInMap(to, config.gasManageAccounts)
}

func (config *bsnConfig) CheckBlacklist(from *common.Address, to *common.Address) bool {
	mu.RLock()
	defer mu.RUnlock()

	if !config.blacklistEnable {
		return true
	}
	return !isAddressInMap(from, config.blacklist) && !isAddressInMap(to, config.blacklist)
}

func isAddressInMap(address *common.Address, m map[string]struct{}) bool {
	if address == nil {
		return false
	}

	return isInMap(strings.ToLower(address.String()), m)
}

func isInMap(s string, m map[string]struct{}) bool {
	if len(m) == 0 {
		return false
	}

	_, ok := m[s]
	return ok
}

# 1. 私链环境搭建

## 1.1. 安装环境

Ubuntu 16或18
CentOS 7
MacOS
golang版本1.14及以上，同时需要C编译器

## 1.2. 获取源代码

```shell
$ git clone https://github.com/DNAProject/Zion
```
## 1.3. 编译源代码

* 编译geth

```shell
$ make geth
```

* 全部编译

```shell
$ make all
```

编译后可在build/bin目录中生成二进制可执行文件

## 1.4. 环境搭建

环境搭建以单机搭建4节点私链网络为例。

### 1.4.1. 生成节点初始配置信息

利用 genesisTool 生成节点初始化配置文件
```
./geth genesisTool generate --basePath <outputPath> --nodeCount <nodeCount> --nodePass <nodePass>
```

以下是生成4节点配置文件示例：

```
./geth genesisTool generate --basePath ./genesisOutput --nodeCount 4 --nodePass 1234
```

生成三份文件：
* genesis.json

```json
{
	"config": {
		"chainId": 10898,
		"homesteadBlock": 0,
		"eip150Block": 0,
		"eip155Block": 0,
		"eip158Block": 0,
		"byzantiumBlock": 0,
		"constantinopleBlock": 0,
		"petersburgBlock": 0,
		"istanbulBlock": 0,
		"hotstuff": {
			"protocol": "basic"
		}
	},
	"alloc": {
		"0x49525E980345C81498fE0e30a9ACC7f4dC9E237B": {
			"publicKey": "0x0213db218e3638d64ae0cb440482c5cfda460ad02759c51a0b53a42f4954e40137",
			"balance": "100000000000000000000000000000"
		},
		"0xA29cfe2827fFf2d38e300be374c8a89214fa5C95": {
			"publicKey": "0x033b2d6b8db288cffe1b10de45e3c920942b069dc6db2a4110a63194fa147352f9",
			"balance": "100000000000000000000000000000"
		},
		"0xAD048c8a4Fc1002B8414F23e4a0105799e9A232D": {
			"publicKey": "0x037595cdec137c1c81fffc70a9bdd77af3add53e91e28cce4675e856de79128cf6",
			"balance": "100000000000000000000000000000"
		},
		"0xeffc0210C58fFE4c523309F0e0918b89911C0985": {
			"publicKey": "0x023bfcaab2a46272bbc5fa5a1fe8c8e19af32120e6299ad7006b6038dd892510ed",
			"balance": "100000000000000000000000000000"
		}
	},
	"coinbase": "0x0000000000000000000000000000000000000000",
	"difficulty": "0x1",
	"extraData": "0x0000000000000000000000000000000000000000000000000000000000000000f89bf8549449525e980345c81498fe0e30a9acc7f4dc9e237b94a29cfe2827fff2d38e300be374c8a89214fa5c9594ad048c8a4fc1002b8414f23e4a0105799e9a232d94effc0210c58ffe4c523309f0e0918b89911c0985b8410000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000c080",
	"gasLimit": "0xffffffff",
	"nonce": "0x4510809143055965",
	"mixhash": "0x0000000000000000000000000000000000000000000000000000000000000000",
	"parentHash": "0x0000000000000000000000000000000000000000000000000000000000000000",
	"timestamp": "0x00"
}
```
* nodes.json 节点的私钥、公钥、地址、静态通讯地址

```json
[
	{
		"address": "0x1D3a781db87D57a2091f968734186c8C72353116",
		"nodeKey": "0x1ca9dcf44053a4241b2b8350b08dcade1b52ecba340237be69004caf76d6ab43",
		"keystore": {
			"address": "1d3a781db87d57a2091f968734186c8c72353116",
			"crypto": {
				"cipher": "aes-128-ctr",
				"ciphertext": "2d5c5cc5ff6c7e2f2e49cb53859df5457fbada13b179f41e160bb6d80897d87f",
				"cipherparams": {
					"iv": "978375653c56f3255763bd509336782b"
				},
				"kdf": "scrypt",
				"kdfparams": {
					"dklen": 32,
					"n": 262144,
					"p": 1,
					"r": 8,
					"salt": "de8296b9cdb58b1f220841c1c84bbc2452a87fbce2d4d7c9cfbaf16e5501dee2"
				},
				"mac": "40cecc8a51db200f9afde4b85d7715951be0295723d148dfaa46ca7da20b0643"
			},
			"id": "93415801-0d57-4d62-8329-f82ebae37939",
			"version": 3
		},
		"pubKey": "0x02b68faf75dd7ec8f2e9a3a23354795069c77e7955ee6b2bdb42e8858d06968a2a",
		"static": "enode://b68faf75dd7ec8f2e9a3a23354795069c77e7955ee6b2bdb42e8858d06968a2a4bfe1073f3a74e3b61b3efa86bb25d655b664a2c106c5fbc07b4455039f84742@127.0.0.1:30300?discport=0"
	},
    ........
```
* static-nodes.json 静态节点通讯配置

```json
[
	"enode://f460929ebaf0ec94f872246c653a5e47c137ca9c8bddb2c872c6cd96d209311ec2380242061e0f5cc6015d137a51430877167c4442099d86a608d2af8b857004@127.0.0.1:30300?discport=0",
	"enode://93a5568672fa325a1c36b0a73c974489c36d6658f71bf56da7b3aa6cc46a3397688a6ae289f7cde8c585c6368aa5a92e0eeedc62888f9cb16c274681e1f040c2@127.0.0.1:30300?discport=0",
	"enode://fdacbff85c9544af0c4dd072d5c570e4854fd9ee7d1677384a1bd6e2d13b245491109e1c2a50b3625fa5ea59dd1682ad7a67f6a340fce3d896f270d92bd1778a@127.0.0.1:30300?discport=0",
	"enode://999ae3f263795e025fb89f96a177287fe620e0509c0c511f2c0c144bbd77b5c52c43bd681d888f1a39a295b7d655e142eac4456c3c1bfcb72b6a602a200047e6@127.0.0.1:30300?discport=0"
]
```

### 1.4.2. 拷贝安装辅助文件目录到安装文件夹

```shell
$ cd zion/docs/install_guide/install_file
$ ls
init.sh setup start.sh stop.sh
$ cp -r setup /your/install/folder/.
$ cp *.sh /your/install/folder/.

$ cd zion/build/bin
$ cp geth /your/install/folder/.

$ cd /path/to/genesisOutput
$ ls genesis.json nodes.json static-nodes.json
$ cp *.json /your/instal/folder/setup/.
```

### 1.4.3. 将setup/node i的目录中的nodekey和pubkey改成setup/nodes.json对应的key

```shell
cd setup
```

顺序修改node0-node3中的nodekey和pubkey，nodekey为keystore, pubkey为pubKey

### 1.4.4. 修改setup/static-nodes.json的ip和端口

将各节点的和端口改为不同的数字

### 1.4.5. 修改start.sh

顺序修改start.sh中的coinbase为setup/nodes.json中的各节点address

### 1.4.6. 执行init.sh，初始化各个节点

按顺序执行4遍init.sh, 在console的交互中输入0-3的节点号

### 1.4.7. 启动节点

执行4遍start.sh, 在console的交互中顺序输入0-3的节点号

### 1.4.8. 停止节点

执行4遍stop.sh，在console的交互中顺序输入0-3的节点号



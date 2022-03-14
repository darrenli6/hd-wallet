package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"math/big"
	"os"

	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/tyler-smith/go-bip39"
)

//生成助记词
func test_mnemonic() {
	// Entropy 生成

	b, err := bip39.NewEntropy(128)
	if err != nil {
		log.Panic("failed to NewEntropy", err, b)
	}
	fmt.Println(b)

	// 生成助记词
	nm, err := bip39.NewMnemonic(b)
	if err != nil {
		log.Panic("failed to NewMnemonic ", err)
	}
	fmt.Println(nm)
}

// 测试助记词有效
func test_ganache() {
	// 助记词转化为种子
	f, err := os.OpenFile(".secret", os.O_RDONLY, 0600)
	defer f.Close()
	if err != nil {
		fmt.Println(err.Error())
	}
	contentByte, err := ioutil.ReadAll(f)

	nm := string(contentByte)
	// 种子转化为账户地址
	//推导路径  再获取钱包
	path := MustParseDerivationPath("m/44'/60'/0'/0/0")

	wallet, err := NewFromMnemonic(nm, "")
	if err != nil {
		log.Panic("faile to NewFromMnemonic", err)
	}
	account, err := wallet.Derive(path, false)
	if err != nil {
		log.Panic("faile to Derive", err)
	}

	fmt.Println(account.Address.Hex())

	path = MustParseDerivationPath("m/44'/60'/0'/0/1")

	account, err = wallet.Derive(path, false)
	if err != nil {
		log.Panic("faile to Derive", err)
	}

	fmt.Println(account.Address.Hex())

}

func test_keystore() {
	// 助记词转化为种子
	f, err := os.OpenFile(".secret", os.O_RDONLY, 0600)
	defer f.Close()
	if err != nil {
		fmt.Println(err.Error())
	}
	contentByte, err := ioutil.ReadAll(f)

	nm := string(contentByte)
	// 种子转化为账户地址
	//推导路径  再获取钱包
	path := MustParseDerivationPath("m/44'/60'/0'/0/0")

	wallet, err := NewFromMnemonic(nm, "")
	if err != nil {
		log.Panic("faile to NewFromMnemonic", err)
	}
	account, err := wallet.Derive(path, false)
	if err != nil {
		log.Panic("faile to Derive", err)
	}

	// fmt.Println(account.Address.Hex())

	pkey, err := wallet.derivePrivateKey(path)
	if err != nil {
		log.Panic("failed to derivePrivateKey", err)
	}

	// 得到私钥
	fmt.Println(*pkey)

	key := NewKeyFromECDSA(pkey)

	hkds := NewHDKeyStore("./data")

	//password
	err = hkds.StoreKey(hkds.JoinPath(account.Address.Hex()), key, "")

	if err != nil {
		log.Panic("failed to StoreKey", err)
	}

}

func test_keystore2() {
	// 助记词转化为种子
	f, err := os.OpenFile(".secret", os.O_RDONLY, 0600)
	defer f.Close()
	if err != nil {
		fmt.Println(err.Error())
	}
	contentByte, err := ioutil.ReadAll(f)

	nm := string(contentByte)
	// 种子转化为账户地址
	//推导路径  再获取钱包
	path := MustParseDerivationPath("m/44'/60'/0'/0/0")

	wallet, err := NewFromMnemonic(nm, "")
	if err != nil {
		log.Panic("faile to NewFromMnemonic", err)
	}
	account, err := wallet.Derive(path, false)
	if err != nil {
		log.Panic("faile to Derive", err)
	}

	account, err = wallet.Derive(path, false)
	if err != nil {
		log.Panic("failed to Derive", err)
	}

	fmt.Println(account.Address.Hex())

	pkey, err := wallet.derivePrivateKey(path)
	if err != nil {
		log.Panic("failed to derivePrivateKey", err)
	}

	// 得到私钥
	fmt.Println(*pkey)

	key := NewKeyFromECDSA(pkey)

	hkds := NewHDKeyStore("./data")

	//password
	err = hkds.StoreKey(hkds.JoinPath(account.Address.Hex()), key, "123")

	if err != nil {
		log.Panic("failed to StoreKey", err)
	}

}

// 解析keystore
func test_des_keystore() {
	hdks := NewHDKeyStore("./data")

	//解析keystore秘钥
	_, err := hdks.GetKey(common.HexToAddress("0x3E6ba59EfFa8d6f3f550287F33b3E6C250C83af0"), hdks.JoinPath("0x3E6ba59EfFa8d6f3f550287F33b3E6C250C83af0"), "123")
	if err != nil {
		log.Panic("failed to getkey", err)
	}

	fmt.Println("get Key ok")
	//拿到key 就能交易
	// 创建交易
	nonce := uint64(0)
	amount := big.NewInt(300000000000000000)
	toAccount := common.HexToAddress("0xD3Bbf16bdCb574803782f75d67e53513Ec2271Bc")
	gasLimit := uint64(21000)
	gasPrice := big.NewInt(2100000000)
	var data []byte
	tx := types.NewTransaction(nonce, toAccount, amount, gasLimit, gasPrice, data)
	// 签名
	account := accounts.Account{common.HexToAddress("0x3E6ba59EfFa8d6f3f550287F33b3E6C250C83af0"), accounts.URL{"", ""}}

	stx, err := hdks.SignTx(account, tx, nil)
	if err != nil {
		log.Panic("failed to SignTx", err)
	}

	// 发送网络节点
	cli, err := ethclient.Dial("http://127.0.0.1:7545")

	if err != nil {
		log.Panic("failed to Dial", err)
	}

	defer cli.Close()

	err = cli.SendTransaction(context.Background(), stx)
	if err != nil {
		log.Panic("failed to SendTransaction", err)
	}

}

func main() {

	// test_mnemonic()
	// test_ganache()

	// test_keystore()
	// test_keystore2()
	test_des_keystore()
}

package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"

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

func main() {

	// test_mnemonic()
	test_ganache()
}

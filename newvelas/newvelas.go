/*
 * Copyright 2018 The OpenWallet Authors
 * This file is part of the OpenWallet library.
 *
 * The OpenWallet library is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Lesser General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * The OpenWallet library is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
 * GNU Lesser General Public License for more details.
 */

package newvelas

import (
	"encoding/hex"
	"encoding/json"
	"github.com/assetsadapterstore/newvelas-adapter/newvelas_addrdec"
	"github.com/astaxie/beego/config"
	"github.com/blocktree/openwallet/v2/log"
	"github.com/blocktree/quorum-adapter/quorum"
	"github.com/blocktree/quorum-adapter/quorum_rpc"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"strings"
)

const (
	Symbol = "NVLX"
)

type WalletManager struct {
	*quorum.WalletManager

	BackupClient *quorum_rpc.Client // 节点客户端
}

func NewWalletManager() *WalletManager {
	wm := WalletManager{}
	wm.WalletManager = quorum.NewWalletManager()
	wm.Blockscanner = NewBlockScanner(&wm)
	wm.Config = quorum.NewConfig(Symbol)
	wm.Log = log.NewOWLogger(wm.Symbol())
	wm.Decoder = newvelas_addrdec.NewAddressDecoder()
	wm.CustomAddressEncodeFunc = wm.CustomAddressEncode
	wm.CustomAddressDecodeFunc = wm.CustomAddressDecode
	return &wm
}

//FullName 币种全名
func (wm *WalletManager) FullName() string {
	return "new velas"
}

func (wm *WalletManager) CustomAddressEncode(address string) string {
	hashHex := strings.TrimPrefix(address, "0x")
	hash, err := hex.DecodeString(hashHex)
	if err != nil {
		return address
	}
	a, err := wm.Decoder.AddressEncode(hash)
	if err != nil {
		return address
	}
	return a
}
func (wm *WalletManager) CustomAddressDecode(address string) string {
	hash, err := wm.Decoder.AddressDecode(address)
	if err != nil {
		return address
	}
	return "0x" + hex.EncodeToString(hash)
}

func (wm *WalletManager) GetBlockByNum(blockNum uint64, showTransactionSpec bool) (*quorum.EthBlock, error) {
	params := []interface{}{
		hexutil.EncodeUint64(blockNum),
		showTransactionSpec,
	}
	var ethBlock quorum.EthBlock

	result, err := wm.WalletClient.Call("eth_getBlockByNumber", params)
	if err != nil {
		return nil, err
	}
	if !result.IsObject() {
		result, err = wm.BackupClient.Call("eth_getBlockByNumber", params)
		if err != nil {
			return nil, err
		}
	}

	if showTransactionSpec {
		err = json.Unmarshal([]byte(result.Raw), &ethBlock)
	} else {
		err = json.Unmarshal([]byte(result.Raw), &ethBlock.BlockHeader)
	}
	if err != nil {
		return nil, err
	}
	ethBlock.BlockHeight, err = hexutil.DecodeUint64(ethBlock.BlockNumber)
	if err != nil {
		return nil, err
	}
	return &ethBlock, nil
}

func (wm *WalletManager) GetTransactionReceipt(transactionId string) (*quorum.TransactionReceipt, error) {
	params := []interface{}{
		transactionId,
	}

	var ethReceipt types.Receipt
	result, err := wm.WalletClient.Call("eth_getTransactionReceipt", params)
	if err != nil {
		return nil, err
	}
	if !result.IsObject() {
		result, err = wm.BackupClient.Call("eth_getTransactionReceipt", params)
		if err != nil {
			return nil, err
		}
	}

	err = ethReceipt.UnmarshalJSON([]byte(result.Raw))
	if err != nil {
		return nil, err
	}

	txReceipt := &quorum.TransactionReceipt{ETHReceipt: &ethReceipt, Raw: result.Raw}

	return txReceipt, nil

}

//LoadAssetsConfig 加载外部配置
func (wm *WalletManager) LoadAssetsConfig(c config.Configer) error {
	wm.WalletManager.LoadAssetsConfig(c)

	backupAPI := c.String("backupAPI")
	client := &quorum_rpc.Client{BaseURL: backupAPI, BroadcastURL: backupAPI, Debug: false}
	wm.BackupClient = client

	return nil

}

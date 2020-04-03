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
	"github.com/assetsadapterstore/newvelas-adapter/newvelas_addrdec"
	"github.com/blocktree/openwallet/v2/log"
	"github.com/blocktree/quorum-adapter/quorum"
	"strings"
)

const (
	Symbol    = "NVLX"
)

type WalletManager struct {
	*quorum.WalletManager
}

func NewWalletManager() *WalletManager {
	wm := WalletManager{}
	wm.WalletManager = quorum.NewWalletManager()
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
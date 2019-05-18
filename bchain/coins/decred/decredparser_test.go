// +build unittest

package decred

import (
	"blockbook/bchain"
	"blockbook/tests/dbtestdata"
	"encoding/hex"
	"fmt"
	"math/big"
	"reflect"
	"testing"
)

func TestDecredParser_GetAddrDescFromAddress(t *testing.T) {
	type args struct {
		address string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "with 0x prefix",
			args: args{address: "0x81b7e08f65bdf5648606c89998a9cc8164397647"},
			want: "81b7e08f65bdf5648606c89998a9cc8164397647",
		},
		{
			name: "without 0x prefix",
			args: args{address: "47526228d673e9f079630d6cdaff5a2ed13e0e60"},
			want: "47526228d673e9f079630d6cdaff5a2ed13e0e60",
		},
		{
			name:    "address of wrong length",
			args:    args{address: "7526228d673e9f079630d6cdaff5a2ed13e0e60"},
			want:    "",
			wantErr: true,
		},
		{
			name:    "ErrAddressMissing",
			args:    args{address: ""},
			want:    "",
			wantErr: true,
		},
		{
			name:    "error - not eth address",
			args:    args{address: "1JKgN43B9SyLuZH19H5ECvr4KcfrbVHzZ6"},
			want:    "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := NewDecredParser(1)
			got, err := p.GetAddrDescFromAddress(tt.args.address)
			if (err != nil) != tt.wantErr {
				t.Errorf("DecredParser.GetAddrDescFromAddress() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			h := hex.EncodeToString(got)
			if !reflect.DeepEqual(h, tt.want) {
				t.Errorf("DecredParser.GetAddrDescFromAddress() = %v, want %v", h, tt.want)
			}
		})
	}
}

var testTx1, testTx2 bchain.Tx

func init() {

	testTx1 = bchain.Tx{
		Blocktime: 1534858022,
		Time:      1534858022,
		Txid:      "0xcd647151552b5132b2aef7c9be00dc6f73afc5901dde157aab131335baaa853b",
		Vin: []bchain.Vin{
			{
				Addresses: []string{"0x3e3a3d69dc66ba10737f531ed088954a9ec89d97"},
			},
		},
		Vout: []bchain.Vout{
			{
				ValueSat: *big.NewInt(1999622000000000000),
				ScriptPubKey: bchain.ScriptPubKey{
					Addresses: []string{"0x555ee11fbddc0e49a9bab358a8941ad95ffdb48f"},
				},
			},
		},
		CoinSpecificData: completeTransaction{
			Tx: &rpcTransaction{
				AccountNonce:     "0xb26c",
				GasPrice:         "0x430e23400",
				GasLimit:         "0x5208",
				To:               "0x555ee11fbddc0e49a9bab358a8941ad95ffdb48f",
				Value:            "0x1bc0159d530e6000",
				Payload:          "0x",
				Hash:             "0xcd647151552b5132b2aef7c9be00dc6f73afc5901dde157aab131335baaa853b",
				BlockNumber:      "0x41eee8",
				From:             "0x3e3a3d69dc66ba10737f531ed088954a9ec89d97",
				TransactionIndex: "0xa",
			},
			Receipt: &rpcReceipt{
				GasUsed: "0x5208",
				Status:  "0x1",
				Logs:    []*rpcLog{},
			},
		},
	}

	testTx2 = bchain.Tx{
		Blocktime: 1534858022,
		Time:      1534858022,
		Txid:      "0xa9cd088aba2131000da6f38a33c20169baee476218deea6b78720700b895b101",
		Vin: []bchain.Vin{
			{
				Addresses: []string{"0x20cd153de35d469ba46127a0c8f18626b59a256a"},
			},
		},
		Vout: []bchain.Vout{
			{
				ValueSat: *big.NewInt(0),
				ScriptPubKey: bchain.ScriptPubKey{
					Addresses: []string{"0x4af4114f73d1c1c903ac9e0361b379d1291808a2"},
				},
			},
		},
		CoinSpecificData: completeTransaction{
			Tx: &rpcTransaction{
				AccountNonce:     "0xd0",
				GasPrice:         "0x9502f9000",
				GasLimit:         "0x130d5",
				To:               "0x4af4114f73d1c1c903ac9e0361b379d1291808a2",
				Value:            "0x0",
				Payload:          "0xa9059cbb000000000000000000000000555ee11fbddc0e49a9bab358a8941ad95ffdb48f00000000000000000000000000000000000000000000021e19e0c9bab2400000",
				Hash:             "0xa9cd088aba2131000da6f38a33c20169baee476218deea6b78720700b895b101",
				BlockNumber:      "0x41eee8",
				From:             "0x20cd153de35d469ba46127a0c8f18626b59a256a",
				TransactionIndex: "0x0"},
			Receipt: &rpcReceipt{
				GasUsed: "0xcb39",
				Status:  "0x1",
				Logs: []*rpcLog{
					{
						Address: "0x4af4114f73d1c1c903ac9e0361b379d1291808a2",
						Data:    "0x00000000000000000000000000000000000000000000021e19e0c9bab2400000",
						Topics: []string{
							"0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef",
							"0x00000000000000000000000020cd153de35d469ba46127a0c8f18626b59a256a",
							"0x000000000000000000000000555ee11fbddc0e49a9bab358a8941ad95ffdb48f",
						},
					},
				},
			},
		},
	}
}

func TestDecredParser_PackTx(t *testing.T) {
	type args struct {
		tx        *bchain.Tx
		height    uint32
		blockTime int64
	}
	tests := []struct {
		name    string
		p       *DecredParser
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "1",
			args: args{
				tx:        &testTx1,
				height:    4321000,
				blockTime: 1534858022,
			},
			want: dbtestdata.EthTx1Packed,
		},
		{
			name: "2",
			args: args{
				tx:        &testTx2,
				height:    4321000,
				blockTime: 1534858022,
			},
			want: dbtestdata.EthTx2Packed,
		},
	}
	p := NewEthereumParser(1)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := p.PackTx(tt.args.tx, tt.args.height, tt.args.blockTime)
			if (err != nil) != tt.wantErr {
				t.Errorf("DecredParser.PackTx() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			h := hex.EncodeToString(got)
			if !reflect.DeepEqual(h, tt.want) {
				t.Errorf("DecredParser.PackTx() = %v, want %v", h, tt.want)
			}
		})
	}
}

func TestDecredParser_UnpackTx(t *testing.T) {
	type args struct {
		hex string
	}
	tests := []struct {
		name    string
		p       *DecredParser
		args    args
		want    *bchain.Tx
		want1   uint32
		wantErr bool
	}{
		{
			name:  "1",
			args:  args{hex: dbtestdata.EthTx1Packed},
			want:  &testTx1,
			want1: 4321000,
		},
		{
			name:  "2",
			args:  args{hex: dbtestdata.EthTx2Packed},
			want:  &testTx2,
			want1: 4321000,
		},
	}
	p := NewDecredParser(1)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b, err := hex.DecodeString(tt.args.hex)
			if err != nil {
				panic(err)
			}
			got, got1, err := p.UnpackTx(b)
			if (err != nil) != tt.wantErr {
				t.Errorf("DecredParser.UnpackTx() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			// DeepEqual has problems with pointers in completeTransaction
			gs := got.CoinSpecificData.(completeTransaction)
			ws := tt.want.CoinSpecificData.(completeTransaction)
			gc := *got
			wc := *tt.want
			gc.CoinSpecificData = nil
			wc.CoinSpecificData = nil
			if fmt.Sprint(gc) != fmt.Sprint(wc) {
				// if !reflect.DeepEqual(gc, wc) {
				t.Errorf("DecredParser.UnpackTx() gc got = %+v, want %+v", gc, wc)
			}
			if !reflect.DeepEqual(gs.Tx, ws.Tx) {
				t.Errorf("DecredParser.UnpackTx() gs.Tx got = %+v, want %+v", gs.Tx, ws.Tx)
			}
			if !reflect.DeepEqual(gs.Receipt, ws.Receipt) {
				t.Errorf("DecredParser.UnpackTx() gs.Receipt got = %+v, want %+v", gs.Receipt, ws.Receipt)
			}
			if got1 != tt.want1 {
				t.Errorf("DecredParser.UnpackTx() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

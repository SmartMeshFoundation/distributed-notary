package models

import (
	"github.com/asdine/storm"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
	"github.com/btcsuite/btcutil"
	"github.com/jinzhu/gorm"
	"github.com/kataras/go-errors"
)

// BTCOutpointStatus :
type BTCOutpointStatus int

const (
	// BTCOutpointStatusUsable 初始状态,即可用状态
	BTCOutpointStatusUsable = iota
	// BTCOutpointStatusUsed 已使用状态
	BTCOutpointStatusUsed
)
const (
	//UTXOPBFTStatusNotUsed 该UTXO刚刚分配,还未使用
	UTXOPBFTStatusNotUsed = iota
	//UTXOPBFTStatusPreSelect 该UTXO已经被主节点标记选中,不要再重复使用了. 如果是不同的主节点选中的,可以继续使用
	UTXOPBFTStatusPreSelect
	//UTXOPBFTStatusCommit 该主节点已经达成共识
	UTXOPBFTStatusCommit
)

// BTCOutpoint 保存公证人分布式私钥对应地址上可用的普通utxo
type BTCOutpoint struct {
	PublicKeyHashStr string            `json:"public_key_hash_str"`
	TxHashStr        string            `json:"tx_hash" gorm:"primary_key"` // utxo所在的txHash
	Index            int               `json:"index"`                      // utxo在tx中的index
	Amount           btcutil.Amount    `json:"amount"`                     // 金额
	Status           BTCOutpointStatus `json:"status"`                     // 0-可用 1-已使用
	CreateTime       int64             `json:"create_time"`                // 创建时间
	UseTime          int64             `json:"use_time"`                   // 使用的时间

	//PBFT use
	PBFTStatus int
	View       int    //如果是UTXOPBFTStatusPreSelect和UTXOPBFTStatusCommit都是分配给谁了
	PBFTReason string //为什么消费此UTXO
}

// GetOutpoint :
func (o *BTCOutpoint) GetOutpoint() *wire.OutPoint {
	utxoTxHash, err := chainhash.NewHashFromStr(o.TxHashStr)
	if err != nil {
		panic(err)
	}
	return wire.NewOutPoint(utxoTxHash, uint32(o.Index))
}

// GetPKScript :
func (o *BTCOutpoint) GetPKScript(net *chaincfg.Params) []byte {
	addr, err := btcutil.DecodeAddress(o.PublicKeyHashStr, net)
	if err != nil {
		panic(err)
	}
	pkScript, err := txscript.PayToAddrScript(addr)
	if err != nil {
		panic(err)
	}
	return pkScript
}

//NewBTCOutpoint :
func (db *DB) NewBTCOutpoint(outpoint *BTCOutpoint) error {
	return db.Create(outpoint).Error
}

//GetBTCOutpoint :
func (db *DB) GetBTCOutpoint(txHashStr string) (outpoint *BTCOutpoint, err error) {
	err = db.Where(&BTCOutpoint{
		TxHashStr: txHashStr,
	}).First(outpoint).Error
	return
}

// GetBTCOutpointList 条件查询
func (db *DB) GetBTCOutpointList(status BTCOutpointStatus) (list []*BTCOutpoint) {
	if status == -1 {
		err := db.Find(&list).Error
		if err == gorm.ErrRecordNotFound {
			err = nil
		}
		return
	}
	err := db.Where(&BTCOutpoint{
		Status: status,
	}).Find(&list).Error
	if err == storm.ErrNotFound {
		err = nil
	}
	return
}

/*
GetAvailableUTXO 用户在侧链prepareLockout以后,公证人需要在主链进行prepareLockout
这时候选取合适的UTXO进行消费,选取原则简单粗暴:
1. 选金额最大的,如果不行
2. 选取top10,然后逐步相加知道超过这个金额
3. 如果仍然没有找到,报错

其他: amount应该包含手续费,
*/
func (db *DB) GetAvailableUTXO(view int, amount btcutil.Amount) (outpoints []*BTCOutpoint, err error) {
	//找最接近amount的
	err = db.Where("amount>=? and (pbft_status=? or (pbft_status=? and view<>?))",
		amount, UTXOPBFTStatusNotUsed, UTXOPBFTStatusPreSelect,
		view).Limit(10).Order("amount asc").Find(&outpoints).Error
	//找不到也会返回nil err
	if err == nil && len(outpoints) > 0 {
		outpoints = []*BTCOutpoint{outpoints[0]} //第一个就够了
		return
	}
	//没有找到,那就选top10,然后相加
	err = db.Where("pbft_status=? or (pbft_status=? and view<>?)",
		UTXOPBFTStatusNotUsed, UTXOPBFTStatusPreSelect,
		view).Limit(10).Order("amount desc").Find(&outpoints).Error
	if err != nil {
		return
	}
	var res []*BTCOutpoint
	var sum btcutil.Amount
	for _, o := range outpoints {
		sum += o.Amount
		res = append(res, o)
		if sum >= amount {
			return res, nil
		}
	}
	err = errUTXONoAvailabe
	return
}

var (
	errUTXODuplciateUase = errors.New("utxo duplicate use")
	errUTXOAlreadyUsed   = errors.New("utxo already used")
	errUTXONoAvailabe    = errors.New("cannout find availalbe BTC UTXO")
)

//PreSelectUTXO PBFT主节点决定选用此UTXO或者收到来自主节点的UTXO选用通知
//针对同一个UTXO,不能被同一个主节点反复使用,
func (db *DB) PreSelectUTXO(TxHashStr string, view int, reason string) error {
	o := &BTCOutpoint{
		TxHashStr: TxHashStr,
	}
	err := db.Where(o).First(o).Error
	if err != nil {
		return err
	}
	if o.PBFTStatus == UTXOPBFTStatusCommit {
		return errUTXOAlreadyUsed
	}
	if o.PBFTStatus == UTXOPBFTStatusPreSelect && o.View == view {
		return errUTXODuplciateUase
	}
	o.PBFTStatus = UTXOPBFTStatusPreSelect
	o.View = view
	return db.Model(o).Updates(
		BTCOutpoint{
			View:       view,
			PBFTStatus: o.PBFTStatus,
		}).Error
}

//CommitUTXO 标记UTXO最终被选定,
func (db *DB) CommitUTXO(TxHashStr string, view int, reason string) error {
	o := &BTCOutpoint{
		TxHashStr: TxHashStr,
	}
	return db.Model(o).Update(BTCOutpoint{
		View:       view,
		PBFTStatus: UTXOPBFTStatusCommit,
		PBFTReason: reason,
	}).Error
}

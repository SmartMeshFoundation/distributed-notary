package service

import (
	"bytes"
	"fmt"
	"math/big"
	"strings"

	"github.com/SmartMeshFoundation/distributed-notary/pbft/pbft"
	utils "github.com/nkbai/goutils"
	"github.com/nkbai/log"

	"github.com/btcsuite/btcutil"
	"github.com/kataras/go-errors"
)

type btcPBFTService struct {
	*PBFTService
}

func (ps *btcPBFTService) UpdateSeq(seq int, op, auxiliary string) {
	nonce, err := ps.db.GetNonce(ps.key)
	if err != nil {
		panic(err)
	}
	if seq > nonce {
		err = ps.db.UpdateNonce(ps.key, seq)
		if err != nil {
			panic(err)
		}
	}
}

/*
	GetOpAuxiliary 根据来自用户的op构造相应的辅助信息,
	对于以太坊来说,就很简单,就是op的hash值
	对于比特币来说就是,分配出去的UTXO列表
	主节点预分配UTXO后应该立即标记选中,否则会造成重复选中
*/
func (ps *btcPBFTService) GetOpAuxiliary(op string, view int) (string, error) {
	ss := strings.Split(op, "-")
	if len(ss) != 3 {
		panic(fmt.Sprintf("op in btc must be chainname-secrethash-amount,op=%s", op))
	}
	if ps.chain != ss[0] {
		panic(fmt.Sprintf("op=%s,expect chain=%s", op, ps.chain))
	}
	amount, b := new(big.Int).SetString(ss[2], 0)
	if !b {
		panic(fmt.Sprintf("op format error,op=%s", op))
	}
	os, err := ps.db.GetAvailableUTXO(view, btcutil.Amount(amount.Int64()))
	if err != nil {
		return "", err
	}
	//auxiliary string格式为digest-outpoint1-outpoint2
	bf := new(bytes.Buffer)
	digest := pbft.Digest(op)
	bf.WriteString(digest)
	bf.WriteByte(byte('-'))
	for _, o := range os {
		_, err = bf.WriteString(o.TxHashStr)
		err = bf.WriteByte(byte('-'))
		err = ps.db.PreSelectUTXO(o.TxHashStr, view, op)
		if err != nil {
			return "", err
		}
	}
	bs := bf.Bytes()
	bs = bs[:len(bs)-1] //去掉最后一个-
	return string(bs), nil
}

//PrepareSeq 实现PBFTAuxiliary
//对于以太坊来说,只要主节点不恶意,都是不会重复的.
//如果恶意,在CommitSeq中会被检测出来.
func (ps *btcPBFTService) PrepareSeq(view, seq int, op string, auxiliary string) error {
	//对于来自我自己的选择,已经选过一次了
	if ps.dispatchService.getSelfNotaryInfo().ID == view {
		return nil
	}
	if len(auxiliary) == 0 {
		return errors.New("auxiliary is empty")
	}
	ss := strings.Split(auxiliary, "-")
	if len(ss) <= 1 {
		return errors.New("auxiliary format is digest-outpoint1-outpoint2")
	}
	for i, tx := range ss {
		if i == 0 {
			continue
		}
		err := ps.db.PreSelectUTXO(tx, view, op)
		if err != nil {
			return err
		}
	}
	return nil
}

/*
	CommitSeq 在集齐验证prepare消息后,验证op对应的auxiliary是否有效.
	对于以太坊来说,总是有效的
	对于比特币来说,可能因为分配出去的utxo已经使用,金额不够等原因造成失败
*/
func (ps *btcPBFTService) CommitSeq(view, seq int, op string, auxiliary string) error {
	if len(auxiliary) == 0 {
		return errors.New("auxiliary is empty")
	}
	ss := strings.Split(auxiliary, "-")
	for i, tx := range ss {
		if i == 0 {
			continue
		}
		err := ps.db.CommitUTXO(tx, view, op)
		if err != nil {
			return err
		}
	}
	return nil
}

func (ps *btcPBFTService) newUTXO(op string) (utxos string, err error) {
	log.Trace(fmt.Sprintf("ps[%s] new nonce for %s", ps.key, op))
	ps.lock.Lock()
	c, ok := ps.nonces[op]
	if ok {
		ps.lock.Unlock()
		err = fmt.Errorf("already exist req %s", op)
		return
	}
	c = make(chan pbft.OpResult, 1)
	ps.nonces[op] = c
	ps.lock.Unlock()
	ps.client.Start(op)
	r := <-c
	log.Trace(fmt.Sprintf("ps[%s] newUTXO return %s", ps.key, utils.StringInterface(r, 3)))
	if r.Error != nil {
		err = r.Error
		return
	}
	ss := strings.Split(r.Auxiliary, "-")
	if len(ss) <= 1 {
		err = errors.New("no valid utxo on reply,it's impossible")
		return
	}
	utxos = strings.Join(ss[1:], "-")
	return
}

package mecdsa

import (
	"math/big"
	"testing"

	"github.com/SmartMeshFoundation/distributed-notary/curv/share"
	"github.com/SmartMeshFoundation/distributed-notary/models"
	"github.com/SmartMeshFoundation/distributed-notary/params"
	"github.com/SmartMeshFoundation/distributed-notary/utils"
	"github.com/stretchr/testify/assert"
)

func init() {
	params.ThresholdCount = 2
	params.ShareCount = 5
}

func newFiveNotaryService() (ns0, ns1, ns2, ns3, ns4 *NotaryService) {
	ns0 = &NotaryService{
		NotaryShareArg: &params.NotaryShareArg{Index: 0},
	}
	ns1 = &NotaryService{
		NotaryShareArg: &params.NotaryShareArg{Index: 1},
	}
	ns2 = &NotaryService{
		NotaryShareArg: &params.NotaryShareArg{Index: 2},
	}
	ns3 = &NotaryService{
		NotaryShareArg: &params.NotaryShareArg{Index: 3},
	}
	ns4 = &NotaryService{
		NotaryShareArg: &params.NotaryShareArg{Index: 4},
	}
	return
}
func newFiveNotaryLockedIn() (l0, l1, l2, l3, l4 *ThresholdPrivKeyGenerator) {
	ns0, ns1, ns2, ns3, ns4 := newFiveNotaryService()
	key := utils.NewRandomHash()
	l0 = &ThresholdPrivKeyGenerator{
		db:           models.SetupTestDB2("l0"),
		srv:          ns0,
		PrivateKeyID: key,
	}
	l1 = &ThresholdPrivKeyGenerator{
		db:           models.SetupTestDB2("l1"),
		srv:          ns1,
		PrivateKeyID: key,
	}
	l2 = &ThresholdPrivKeyGenerator{
		db:           models.SetupTestDB2("l2"),
		srv:          ns2,
		PrivateKeyID: key,
	}
	l3 = &ThresholdPrivKeyGenerator{
		db:           models.SetupTestDB2("l3"),
		srv:          ns3,
		PrivateKeyID: key,
	}
	l4 = &ThresholdPrivKeyGenerator{
		db:           models.SetupTestDB2("l4"),
		srv:          ns4,
		PrivateKeyID: key,
	}
	return
}
func TestLockedIn(t *testing.T) {
	var finish bool
	var err error

	//步骤1 -----------------------------------------------
	l0, l1, l2, l3, l4 := newFiveNotaryLockedIn()
	msg10, err := l0.GeneratePhase1PubKeyProof()
	assert.EqualValues(t, err, nil)
	msg11, err := l1.GeneratePhase1PubKeyProof()
	assert.EqualValues(t, err, nil)
	msg12, err := l2.GeneratePhase1PubKeyProof()
	assert.EqualValues(t, err, nil)
	msg13, err := l3.GeneratePhase1PubKeyProof()
	assert.EqualValues(t, err, nil)
	msg14, err := l4.GeneratePhase1PubKeyProof()
	assert.EqualValues(t, err, nil)

	//步骤1 :0 添加其他人
	finish, err = l0.ReceivePhase1PubKeyProof(msg11, 1)
	assert.EqualValues(t, err, nil)
	assert.EqualValues(t, finish, false)
	finish, err = l0.ReceivePhase1PubKeyProof(msg11, 1)
	assert.EqualValues(t, err == nil, false)
	assert.EqualValues(t, finish, false)
	finish, err = l0.ReceivePhase1PubKeyProof(msg12, 2)
	assert.EqualValues(t, err, nil)
	assert.EqualValues(t, finish, false)
	finish, err = l0.ReceivePhase1PubKeyProof(msg13, 3)
	assert.EqualValues(t, err, nil)
	assert.EqualValues(t, finish, false)
	finish, err = l0.ReceivePhase1PubKeyProof(msg14, 4)
	assert.EqualValues(t, err, nil)
	assert.EqualValues(t, finish, true)

	//1 添加其他人
	finish, err = l1.ReceivePhase1PubKeyProof(msg10, 0)
	assert.EqualValues(t, err, nil)
	assert.EqualValues(t, finish, false)
	finish, err = l1.ReceivePhase1PubKeyProof(msg12, 2)
	assert.EqualValues(t, err, nil)
	assert.EqualValues(t, finish, false)
	finish, err = l1.ReceivePhase1PubKeyProof(msg13, 3)
	assert.EqualValues(t, err, nil)
	assert.EqualValues(t, finish, false)
	finish, err = l1.ReceivePhase1PubKeyProof(msg14, 4)
	assert.EqualValues(t, err, nil)
	assert.EqualValues(t, finish, true)

	//2 添加其他人
	finish, err = l2.ReceivePhase1PubKeyProof(msg10, 0)
	assert.EqualValues(t, err, nil)
	assert.EqualValues(t, finish, false)
	finish, err = l2.ReceivePhase1PubKeyProof(msg11, 1)
	assert.EqualValues(t, err, nil)
	assert.EqualValues(t, finish, false)
	finish, err = l2.ReceivePhase1PubKeyProof(msg13, 3)
	assert.EqualValues(t, err, nil)
	assert.EqualValues(t, finish, false)
	finish, err = l2.ReceivePhase1PubKeyProof(msg14, 4)
	assert.EqualValues(t, err, nil)
	assert.EqualValues(t, finish, true)

	//3 添加其他人
	finish, err = l3.ReceivePhase1PubKeyProof(msg10, 0)
	assert.EqualValues(t, err, nil)
	assert.EqualValues(t, finish, false)
	finish, err = l3.ReceivePhase1PubKeyProof(msg12, 2)
	assert.EqualValues(t, err, nil)
	assert.EqualValues(t, finish, false)
	finish, err = l3.ReceivePhase1PubKeyProof(msg11, 1)
	assert.EqualValues(t, err, nil)
	assert.EqualValues(t, finish, false)
	finish, err = l3.ReceivePhase1PubKeyProof(msg14, 4)
	assert.EqualValues(t, err, nil)
	assert.EqualValues(t, finish, true)

	//4 添加其他人
	finish, err = l4.ReceivePhase1PubKeyProof(msg10, 0)
	assert.EqualValues(t, err, nil)
	assert.EqualValues(t, finish, false)
	finish, err = l4.ReceivePhase1PubKeyProof(msg12, 2)
	assert.EqualValues(t, err, nil)
	assert.EqualValues(t, finish, false)
	finish, err = l4.ReceivePhase1PubKeyProof(msg13, 3)
	assert.EqualValues(t, err, nil)
	assert.EqualValues(t, finish, false)
	finish, err = l4.ReceivePhase1PubKeyProof(msg11, 1)
	assert.EqualValues(t, err, nil)
	assert.EqualValues(t, finish, true)

	//步骤二------------------------------------------
	msg20, err := l0.GeneratePhase2PaillierKeyProof()
	msg21, err := l1.GeneratePhase2PaillierKeyProof()
	msg22, err := l2.GeneratePhase2PaillierKeyProof()
	msg23, err := l3.GeneratePhase2PaillierKeyProof()
	msg24, err := l4.GeneratePhase2PaillierKeyProof()

	//0 添加其他人
	finish, err = l0.ReceivePhase2PaillierPubKeyProof(msg21, 1)
	assert.EqualValues(t, err, nil)
	assert.EqualValues(t, finish, false)
	finish, err = l0.ReceivePhase2PaillierPubKeyProof(msg22, 2)
	assert.EqualValues(t, err, nil)
	assert.EqualValues(t, finish, false)
	finish, err = l0.ReceivePhase2PaillierPubKeyProof(msg23, 3)
	assert.EqualValues(t, err, nil)
	assert.EqualValues(t, finish, false)
	finish, err = l0.ReceivePhase2PaillierPubKeyProof(msg24, 4)
	assert.EqualValues(t, err, nil)
	assert.EqualValues(t, finish, true)

	//1 添加其他人
	finish, err = l1.ReceivePhase2PaillierPubKeyProof(msg20, 0)
	assert.EqualValues(t, err, nil)
	assert.EqualValues(t, finish, false)
	finish, err = l1.ReceivePhase2PaillierPubKeyProof(msg22, 2)
	assert.EqualValues(t, err, nil)
	assert.EqualValues(t, finish, false)
	finish, err = l1.ReceivePhase2PaillierPubKeyProof(msg23, 3)
	assert.EqualValues(t, err, nil)
	assert.EqualValues(t, finish, false)
	finish, err = l1.ReceivePhase2PaillierPubKeyProof(msg24, 4)
	assert.EqualValues(t, err, nil)
	assert.EqualValues(t, finish, true)

	//2 添加其他人
	finish, err = l2.ReceivePhase2PaillierPubKeyProof(msg20, 0)
	assert.EqualValues(t, err, nil)
	assert.EqualValues(t, finish, false)
	finish, err = l2.ReceivePhase2PaillierPubKeyProof(msg21, 1)
	assert.EqualValues(t, err, nil)
	assert.EqualValues(t, finish, false)
	finish, err = l2.ReceivePhase2PaillierPubKeyProof(msg23, 3)
	assert.EqualValues(t, err, nil)
	assert.EqualValues(t, finish, false)
	finish, err = l2.ReceivePhase2PaillierPubKeyProof(msg24, 4)
	assert.EqualValues(t, err, nil)
	assert.EqualValues(t, finish, true)

	//3 添加其他人
	finish, err = l3.ReceivePhase2PaillierPubKeyProof(msg20, 0)
	assert.EqualValues(t, err, nil)
	assert.EqualValues(t, finish, false)
	finish, err = l3.ReceivePhase2PaillierPubKeyProof(msg22, 2)
	assert.EqualValues(t, err, nil)
	assert.EqualValues(t, finish, false)
	finish, err = l3.ReceivePhase2PaillierPubKeyProof(msg21, 1)
	assert.EqualValues(t, err, nil)
	assert.EqualValues(t, finish, false)
	finish, err = l3.ReceivePhase2PaillierPubKeyProof(msg24, 4)
	assert.EqualValues(t, err, nil)
	assert.EqualValues(t, finish, true)

	//4 添加其他人
	finish, err = l4.ReceivePhase2PaillierPubKeyProof(msg20, 0)
	assert.EqualValues(t, err, nil)
	assert.EqualValues(t, finish, false)
	finish, err = l4.ReceivePhase2PaillierPubKeyProof(msg22, 2)
	assert.EqualValues(t, err, nil)
	assert.EqualValues(t, finish, false)
	finish, err = l4.ReceivePhase2PaillierPubKeyProof(msg23, 3)
	assert.EqualValues(t, err, nil)
	assert.EqualValues(t, finish, false)
	finish, err = l4.ReceivePhase2PaillierPubKeyProof(msg21, 1)
	assert.EqualValues(t, err, nil)
	assert.EqualValues(t, finish, true)

	//步骤三 定向广播secret share--------------------------------------------------

	msg30, err := l0.GeneratePhase3SecretShare()
	//如果生成私钥的过程中某个公证人不诚实呢? 如果他用的不是一开始声明的UI,将会造成广播出去的pk不一致,也会被拒绝.
	//if true {
	//	p, err := l1.db.LoadPrivatedKeyInfo(l1.PrivateKeyID)
	//	if err != nil {
	//		return
	//	}
	//	p.UI.D = new(big.Int).Add(p.UI.D, big.NewInt(1))
	//	err = l1.db.TestSave(p)
	//	if err != nil {
	//		panic(err.Error())
	//	}
	//}
	msg31, err := l1.GeneratePhase3SecretShare()
	if false {
		//不能传递给别人一个错的secret share,会被检测出来.
		msg31[0].SecretShare.D = new(big.Int).Add(msg31[0].SecretShare.D, big.NewInt(1))
	}
	msg32, err := l2.GeneratePhase3SecretShare()
	msg33, err := l3.GeneratePhase3SecretShare()
	msg34, err := l4.GeneratePhase3SecretShare()

	//0 添加其他人的shares
	finish, err = l0.ReceivePhase3SecretShare(msg31[0], 1)
	assert.EqualValues(t, err, nil)
	assert.EqualValues(t, finish, false)
	finish, err = l0.ReceivePhase3SecretShare(msg32[0], 2)
	assert.EqualValues(t, err, nil)
	assert.EqualValues(t, finish, false)
	finish, err = l0.ReceivePhase3SecretShare(msg33[0], 3)
	assert.EqualValues(t, err, nil)
	assert.EqualValues(t, finish, false)
	finish, err = l0.ReceivePhase3SecretShare(msg34[0], 4)
	assert.EqualValues(t, err, nil)
	assert.EqualValues(t, finish, true)

	//1 添加其他人的shares
	finish, err = l1.ReceivePhase3SecretShare(msg30[1], 0)
	assert.EqualValues(t, err, nil)
	assert.EqualValues(t, finish, false)
	finish, err = l1.ReceivePhase3SecretShare(msg32[1], 2)
	assert.EqualValues(t, err, nil)
	assert.EqualValues(t, finish, false)
	finish, err = l1.ReceivePhase3SecretShare(msg33[1], 3)
	assert.EqualValues(t, err, nil)
	assert.EqualValues(t, finish, false)
	finish, err = l1.ReceivePhase3SecretShare(msg34[1], 4)
	assert.EqualValues(t, err, nil)
	assert.EqualValues(t, finish, true)

	//2 添加其他人的shares
	finish, err = l2.ReceivePhase3SecretShare(msg30[2], 0)
	assert.EqualValues(t, err, nil)
	assert.EqualValues(t, finish, false)
	finish, err = l2.ReceivePhase3SecretShare(msg31[2], 1)
	assert.EqualValues(t, err, nil)
	assert.EqualValues(t, finish, false)
	finish, err = l2.ReceivePhase3SecretShare(msg33[2], 3)
	assert.EqualValues(t, err, nil)
	assert.EqualValues(t, finish, false)
	finish, err = l2.ReceivePhase3SecretShare(msg34[2], 4)
	assert.EqualValues(t, err, nil)
	assert.EqualValues(t, finish, true)

	//3 添加其他人的shares
	finish, err = l3.ReceivePhase3SecretShare(msg30[3], 0)
	assert.EqualValues(t, err, nil)
	assert.EqualValues(t, finish, false)
	finish, err = l3.ReceivePhase3SecretShare(msg32[3], 2)
	assert.EqualValues(t, err, nil)
	assert.EqualValues(t, finish, false)
	finish, err = l3.ReceivePhase3SecretShare(msg31[3], 1)
	assert.EqualValues(t, err, nil)
	assert.EqualValues(t, finish, false)
	finish, err = l3.ReceivePhase3SecretShare(msg34[3], 4)
	assert.EqualValues(t, err, nil)
	assert.EqualValues(t, finish, true)

	//4 添加其他人的shares
	finish, err = l4.ReceivePhase3SecretShare(msg30[4], 0)
	assert.EqualValues(t, err, nil)
	assert.EqualValues(t, finish, false)
	finish, err = l4.ReceivePhase3SecretShare(msg32[4], 2)
	assert.EqualValues(t, err, nil)
	assert.EqualValues(t, finish, false)
	finish, err = l4.ReceivePhase3SecretShare(msg33[4], 3)
	assert.EqualValues(t, err, nil)
	assert.EqualValues(t, finish, false)
	finish, err = l4.ReceivePhase3SecretShare(msg31[4], 1)
	assert.EqualValues(t, err, nil)
	assert.EqualValues(t, finish, true)

	//步骤四 ,验证最后key---------------------------------------
	msg40, err := l0.GeneratePhase4PubKeyProof()
	msg41, err := l1.GeneratePhase4PubKeyProof()
	msg42, err := l2.GeneratePhase4PubKeyProof()
	msg43, err := l3.GeneratePhase4PubKeyProof()
	msg44, err := l4.GeneratePhase4PubKeyProof()

	//0 添加其他人
	finish, err = l0.ReceivePhase4VerifyTotalPubKey(msg41, 1)
	assert.EqualValues(t, err, nil)
	assert.EqualValues(t, finish, false)
	finish, err = l0.ReceivePhase4VerifyTotalPubKey(msg42, 2)
	assert.EqualValues(t, err, nil)
	assert.EqualValues(t, finish, false)
	finish, err = l0.ReceivePhase4VerifyTotalPubKey(msg43, 3)
	assert.EqualValues(t, err, nil)
	assert.EqualValues(t, finish, false)
	finish, err = l0.ReceivePhase4VerifyTotalPubKey(msg44, 4)
	assert.EqualValues(t, err, nil)
	assert.EqualValues(t, finish, true)

	//1 添加其他人
	finish, err = l1.ReceivePhase4VerifyTotalPubKey(msg40, 0)
	assert.EqualValues(t, err, nil)
	assert.EqualValues(t, finish, false)
	finish, err = l1.ReceivePhase4VerifyTotalPubKey(msg42, 2)
	assert.EqualValues(t, err, nil)
	assert.EqualValues(t, finish, false)
	finish, err = l1.ReceivePhase4VerifyTotalPubKey(msg43, 3)
	assert.EqualValues(t, err, nil)
	assert.EqualValues(t, finish, false)
	finish, err = l1.ReceivePhase4VerifyTotalPubKey(msg44, 4)
	assert.EqualValues(t, err, nil)
	assert.EqualValues(t, finish, true)

	//2 添加其他人
	finish, err = l2.ReceivePhase4VerifyTotalPubKey(msg41, 1)
	assert.EqualValues(t, err, nil)
	assert.EqualValues(t, finish, false)
	finish, err = l2.ReceivePhase4VerifyTotalPubKey(msg40, 0)
	assert.EqualValues(t, err, nil)
	assert.EqualValues(t, finish, false)
	finish, err = l2.ReceivePhase4VerifyTotalPubKey(msg43, 3)
	assert.EqualValues(t, err, nil)
	assert.EqualValues(t, finish, false)
	finish, err = l2.ReceivePhase4VerifyTotalPubKey(msg44, 4)
	assert.EqualValues(t, err, nil)
	assert.EqualValues(t, finish, true)

	//3 添加其他人
	finish, err = l3.ReceivePhase4VerifyTotalPubKey(msg41, 1)
	assert.EqualValues(t, err, nil)
	assert.EqualValues(t, finish, false)
	finish, err = l3.ReceivePhase4VerifyTotalPubKey(msg42, 2)
	assert.EqualValues(t, err, nil)
	assert.EqualValues(t, finish, false)
	finish, err = l3.ReceivePhase4VerifyTotalPubKey(msg40, 0)
	assert.EqualValues(t, err, nil)
	assert.EqualValues(t, finish, false)
	finish, err = l3.ReceivePhase4VerifyTotalPubKey(msg44, 4)
	assert.EqualValues(t, err, nil)
	assert.EqualValues(t, finish, true)

	//4 添加其他人
	finish, err = l4.ReceivePhase4VerifyTotalPubKey(msg41, 1)
	assert.EqualValues(t, err, nil)
	assert.EqualValues(t, finish, false)
	finish, err = l4.ReceivePhase4VerifyTotalPubKey(msg42, 2)
	assert.EqualValues(t, err, nil)
	assert.EqualValues(t, finish, false)
	finish, err = l4.ReceivePhase4VerifyTotalPubKey(msg43, 3)
	assert.EqualValues(t, err, nil)
	assert.EqualValues(t, finish, false)
	finish, err = l4.ReceivePhase4VerifyTotalPubKey(msg40, 0)
	assert.EqualValues(t, err, nil)
	assert.EqualValues(t, finish, true)

	p0, err := l0.db.LoadPrivatedKeyInfo(l0.PrivateKeyID)
	p1, err := l1.db.LoadPrivatedKeyInfo(l0.PrivateKeyID)
	p2, err := l2.db.LoadPrivatedKeyInfo(l0.PrivateKeyID)
	p3, err := l3.db.LoadPrivatedKeyInfo(l0.PrivateKeyID)
	p4, err := l4.db.LoadPrivatedKeyInfo(l0.PrivateKeyID)

	//校验私钥是否正确分配
	var xi []share.SPrivKey
	xi = append(xi, p0.XI)
	xi = append(xi, p1.XI)
	xi = append(xi, p2.XI)
	xi = append(xi, p3.XI)
	xi = append(xi, p4.XI)

	totalPrivKey := msg30[1].Vss.Reconstruct([]int{0, 1, 2}, xi[0:3])
	sum := share.PrivKeyZero.Clone()
	share.ModAdd(sum, p0.UI)
	share.ModAdd(sum, p1.UI)
	share.ModAdd(sum, p2.UI)
	share.ModAdd(sum, p3.UI)
	share.ModAdd(sum, p4.UI)

	if sum.D.Cmp(totalPrivKey.D) != 0 {
		t.Error("not equal")
	}
	pubx, puby := share.S.ScalarBaseMult(sum.Bytes())
	if pubx.Cmp(p0.PublicKeyX) != 0 || puby.Cmp(p0.PublicKeyY) != 0 {
		t.Error("pub key error")
	}
	//t.Logf("p0=%s", utils.StringInterface(p0, 2))
	//t.Logf("p1=%s", utils.StringInterface(p1, 2))
	//t.Logf("p2=%s", utils.StringInterface(p2, 2))
	//t.Logf("p3=%s", utils.StringInterface(p3, 2))
	//t.Logf("p4=%s", utils.StringInterface(p4, 2))
}

func newTestLockin(t *testing.T) (l0, l1, l2, l3, l4 *ThresholdPrivKeyGenerator) {
	var finish bool
	var err error

	//步骤1 -----------------------------------------------
	l0, l1, l2, l3, l4 = newFiveNotaryLockedIn()
	msg10, err := l0.GeneratePhase1PubKeyProof()
	assert.EqualValues(t, err, nil)
	msg11, err := l1.GeneratePhase1PubKeyProof()
	assert.EqualValues(t, err, nil)
	msg12, err := l2.GeneratePhase1PubKeyProof()
	assert.EqualValues(t, err, nil)
	msg13, err := l3.GeneratePhase1PubKeyProof()
	assert.EqualValues(t, err, nil)
	msg14, err := l4.GeneratePhase1PubKeyProof()
	assert.EqualValues(t, err, nil)

	//步骤1 :0 添加其他人
	finish, err = l0.ReceivePhase1PubKeyProof(msg11, 1)
	assert.EqualValues(t, err, nil)
	assert.EqualValues(t, finish, false)
	finish, err = l0.ReceivePhase1PubKeyProof(msg11, 1)
	assert.EqualValues(t, err == nil, false)
	assert.EqualValues(t, finish, false)
	finish, err = l0.ReceivePhase1PubKeyProof(msg12, 2)
	assert.EqualValues(t, err, nil)
	assert.EqualValues(t, finish, false)
	finish, err = l0.ReceivePhase1PubKeyProof(msg13, 3)
	assert.EqualValues(t, err, nil)
	assert.EqualValues(t, finish, false)
	finish, err = l0.ReceivePhase1PubKeyProof(msg14, 4)
	assert.EqualValues(t, err, nil)
	assert.EqualValues(t, finish, true)

	//1 添加其他人
	finish, err = l1.ReceivePhase1PubKeyProof(msg10, 0)
	assert.EqualValues(t, err, nil)
	assert.EqualValues(t, finish, false)
	finish, err = l1.ReceivePhase1PubKeyProof(msg12, 2)
	assert.EqualValues(t, err, nil)
	assert.EqualValues(t, finish, false)
	finish, err = l1.ReceivePhase1PubKeyProof(msg13, 3)
	assert.EqualValues(t, err, nil)
	assert.EqualValues(t, finish, false)
	finish, err = l1.ReceivePhase1PubKeyProof(msg14, 4)
	assert.EqualValues(t, err, nil)
	assert.EqualValues(t, finish, true)

	//2 添加其他人
	finish, err = l2.ReceivePhase1PubKeyProof(msg10, 0)
	assert.EqualValues(t, err, nil)
	assert.EqualValues(t, finish, false)
	finish, err = l2.ReceivePhase1PubKeyProof(msg11, 1)
	assert.EqualValues(t, err, nil)
	assert.EqualValues(t, finish, false)
	finish, err = l2.ReceivePhase1PubKeyProof(msg13, 3)
	assert.EqualValues(t, err, nil)
	assert.EqualValues(t, finish, false)
	finish, err = l2.ReceivePhase1PubKeyProof(msg14, 4)
	assert.EqualValues(t, err, nil)
	assert.EqualValues(t, finish, true)

	//3 添加其他人
	finish, err = l3.ReceivePhase1PubKeyProof(msg10, 0)
	assert.EqualValues(t, err, nil)
	assert.EqualValues(t, finish, false)
	finish, err = l3.ReceivePhase1PubKeyProof(msg12, 2)
	assert.EqualValues(t, err, nil)
	assert.EqualValues(t, finish, false)
	finish, err = l3.ReceivePhase1PubKeyProof(msg11, 1)
	assert.EqualValues(t, err, nil)
	assert.EqualValues(t, finish, false)
	finish, err = l3.ReceivePhase1PubKeyProof(msg14, 4)
	assert.EqualValues(t, err, nil)
	assert.EqualValues(t, finish, true)

	//4 添加其他人
	finish, err = l4.ReceivePhase1PubKeyProof(msg10, 0)
	assert.EqualValues(t, err, nil)
	assert.EqualValues(t, finish, false)
	finish, err = l4.ReceivePhase1PubKeyProof(msg12, 2)
	assert.EqualValues(t, err, nil)
	assert.EqualValues(t, finish, false)
	finish, err = l4.ReceivePhase1PubKeyProof(msg13, 3)
	assert.EqualValues(t, err, nil)
	assert.EqualValues(t, finish, false)
	finish, err = l4.ReceivePhase1PubKeyProof(msg11, 1)
	assert.EqualValues(t, err, nil)
	assert.EqualValues(t, finish, true)

	//步骤二------------------------------------------
	msg20, err := l0.GeneratePhase2PaillierKeyProof()
	msg21, err := l1.GeneratePhase2PaillierKeyProof()
	msg22, err := l2.GeneratePhase2PaillierKeyProof()
	msg23, err := l3.GeneratePhase2PaillierKeyProof()
	msg24, err := l4.GeneratePhase2PaillierKeyProof()

	//0 添加其他人
	finish, err = l0.ReceivePhase2PaillierPubKeyProof(msg21, 1)
	assert.EqualValues(t, err, nil)
	assert.EqualValues(t, finish, false)
	finish, err = l0.ReceivePhase2PaillierPubKeyProof(msg22, 2)
	assert.EqualValues(t, err, nil)
	assert.EqualValues(t, finish, false)
	finish, err = l0.ReceivePhase2PaillierPubKeyProof(msg23, 3)
	assert.EqualValues(t, err, nil)
	assert.EqualValues(t, finish, false)
	finish, err = l0.ReceivePhase2PaillierPubKeyProof(msg24, 4)
	assert.EqualValues(t, err, nil)
	assert.EqualValues(t, finish, true)

	//1 添加其他人
	finish, err = l1.ReceivePhase2PaillierPubKeyProof(msg20, 0)
	assert.EqualValues(t, err, nil)
	assert.EqualValues(t, finish, false)
	finish, err = l1.ReceivePhase2PaillierPubKeyProof(msg22, 2)
	assert.EqualValues(t, err, nil)
	assert.EqualValues(t, finish, false)
	finish, err = l1.ReceivePhase2PaillierPubKeyProof(msg23, 3)
	assert.EqualValues(t, err, nil)
	assert.EqualValues(t, finish, false)
	finish, err = l1.ReceivePhase2PaillierPubKeyProof(msg24, 4)
	assert.EqualValues(t, err, nil)
	assert.EqualValues(t, finish, true)

	//2 添加其他人
	finish, err = l2.ReceivePhase2PaillierPubKeyProof(msg20, 0)
	assert.EqualValues(t, err, nil)
	assert.EqualValues(t, finish, false)
	finish, err = l2.ReceivePhase2PaillierPubKeyProof(msg21, 1)
	assert.EqualValues(t, err, nil)
	assert.EqualValues(t, finish, false)
	finish, err = l2.ReceivePhase2PaillierPubKeyProof(msg23, 3)
	assert.EqualValues(t, err, nil)
	assert.EqualValues(t, finish, false)
	finish, err = l2.ReceivePhase2PaillierPubKeyProof(msg24, 4)
	assert.EqualValues(t, err, nil)
	assert.EqualValues(t, finish, true)

	//3 添加其他人
	finish, err = l3.ReceivePhase2PaillierPubKeyProof(msg20, 0)
	assert.EqualValues(t, err, nil)
	assert.EqualValues(t, finish, false)
	finish, err = l3.ReceivePhase2PaillierPubKeyProof(msg22, 2)
	assert.EqualValues(t, err, nil)
	assert.EqualValues(t, finish, false)
	finish, err = l3.ReceivePhase2PaillierPubKeyProof(msg21, 1)
	assert.EqualValues(t, err, nil)
	assert.EqualValues(t, finish, false)
	finish, err = l3.ReceivePhase2PaillierPubKeyProof(msg24, 4)
	assert.EqualValues(t, err, nil)
	assert.EqualValues(t, finish, true)

	//4 添加其他人
	finish, err = l4.ReceivePhase2PaillierPubKeyProof(msg20, 0)
	assert.EqualValues(t, err, nil)
	assert.EqualValues(t, finish, false)
	finish, err = l4.ReceivePhase2PaillierPubKeyProof(msg22, 2)
	assert.EqualValues(t, err, nil)
	assert.EqualValues(t, finish, false)
	finish, err = l4.ReceivePhase2PaillierPubKeyProof(msg23, 3)
	assert.EqualValues(t, err, nil)
	assert.EqualValues(t, finish, false)
	finish, err = l4.ReceivePhase2PaillierPubKeyProof(msg21, 1)
	assert.EqualValues(t, err, nil)
	assert.EqualValues(t, finish, true)

	//步骤三 定向广播secret share--------------------------------------------------

	msg30, err := l0.GeneratePhase3SecretShare()
	msg31, err := l1.GeneratePhase3SecretShare()
	msg32, err := l2.GeneratePhase3SecretShare()
	msg33, err := l3.GeneratePhase3SecretShare()
	msg34, err := l4.GeneratePhase3SecretShare()

	//0 添加其他人的shares
	finish, err = l0.ReceivePhase3SecretShare(msg31[0], 1)
	assert.EqualValues(t, err, nil)
	assert.EqualValues(t, finish, false)
	finish, err = l0.ReceivePhase3SecretShare(msg32[0], 2)
	assert.EqualValues(t, err, nil)
	assert.EqualValues(t, finish, false)
	finish, err = l0.ReceivePhase3SecretShare(msg33[0], 3)
	assert.EqualValues(t, err, nil)
	assert.EqualValues(t, finish, false)
	finish, err = l0.ReceivePhase3SecretShare(msg34[0], 4)
	assert.EqualValues(t, err, nil)
	assert.EqualValues(t, finish, true)

	//1 添加其他人的shares
	finish, err = l1.ReceivePhase3SecretShare(msg30[1], 0)
	assert.EqualValues(t, err, nil)
	assert.EqualValues(t, finish, false)
	finish, err = l1.ReceivePhase3SecretShare(msg32[1], 2)
	assert.EqualValues(t, err, nil)
	assert.EqualValues(t, finish, false)
	finish, err = l1.ReceivePhase3SecretShare(msg33[1], 3)
	assert.EqualValues(t, err, nil)
	assert.EqualValues(t, finish, false)
	finish, err = l1.ReceivePhase3SecretShare(msg34[1], 4)
	assert.EqualValues(t, err, nil)
	assert.EqualValues(t, finish, true)

	//2 添加其他人的shares
	finish, err = l2.ReceivePhase3SecretShare(msg30[2], 0)
	assert.EqualValues(t, err, nil)
	assert.EqualValues(t, finish, false)
	finish, err = l2.ReceivePhase3SecretShare(msg31[2], 1)
	assert.EqualValues(t, err, nil)
	assert.EqualValues(t, finish, false)
	finish, err = l2.ReceivePhase3SecretShare(msg33[2], 3)
	assert.EqualValues(t, err, nil)
	assert.EqualValues(t, finish, false)
	finish, err = l2.ReceivePhase3SecretShare(msg34[2], 4)
	assert.EqualValues(t, err, nil)
	assert.EqualValues(t, finish, true)

	//3 添加其他人的shares
	finish, err = l3.ReceivePhase3SecretShare(msg30[3], 0)
	assert.EqualValues(t, err, nil)
	assert.EqualValues(t, finish, false)
	finish, err = l3.ReceivePhase3SecretShare(msg32[3], 2)
	assert.EqualValues(t, err, nil)
	assert.EqualValues(t, finish, false)
	finish, err = l3.ReceivePhase3SecretShare(msg31[3], 1)
	assert.EqualValues(t, err, nil)
	assert.EqualValues(t, finish, false)
	finish, err = l3.ReceivePhase3SecretShare(msg34[3], 4)
	assert.EqualValues(t, err, nil)
	assert.EqualValues(t, finish, true)

	//4 添加其他人的shares
	finish, err = l4.ReceivePhase3SecretShare(msg30[4], 0)
	assert.EqualValues(t, err, nil)
	assert.EqualValues(t, finish, false)
	finish, err = l4.ReceivePhase3SecretShare(msg32[4], 2)
	assert.EqualValues(t, err, nil)
	assert.EqualValues(t, finish, false)
	finish, err = l4.ReceivePhase3SecretShare(msg33[4], 3)
	assert.EqualValues(t, err, nil)
	assert.EqualValues(t, finish, false)
	finish, err = l4.ReceivePhase3SecretShare(msg31[4], 1)
	assert.EqualValues(t, err, nil)
	assert.EqualValues(t, finish, true)

	//步骤四 ,验证最后key---------------------------------------
	msg40, err := l0.GeneratePhase4PubKeyProof()
	msg41, err := l1.GeneratePhase4PubKeyProof()
	msg42, err := l2.GeneratePhase4PubKeyProof()
	msg43, err := l3.GeneratePhase4PubKeyProof()
	msg44, err := l4.GeneratePhase4PubKeyProof()

	//0 添加其他人
	finish, err = l0.ReceivePhase4VerifyTotalPubKey(msg41, 1)
	assert.EqualValues(t, err, nil)
	assert.EqualValues(t, finish, false)
	finish, err = l0.ReceivePhase4VerifyTotalPubKey(msg42, 2)
	assert.EqualValues(t, err, nil)
	assert.EqualValues(t, finish, false)
	finish, err = l0.ReceivePhase4VerifyTotalPubKey(msg43, 3)
	assert.EqualValues(t, err, nil)
	assert.EqualValues(t, finish, false)
	finish, err = l0.ReceivePhase4VerifyTotalPubKey(msg44, 4)
	assert.EqualValues(t, err, nil)
	assert.EqualValues(t, finish, true)

	//1 添加其他人
	finish, err = l1.ReceivePhase4VerifyTotalPubKey(msg40, 0)
	assert.EqualValues(t, err, nil)
	assert.EqualValues(t, finish, false)
	finish, err = l1.ReceivePhase4VerifyTotalPubKey(msg42, 2)
	assert.EqualValues(t, err, nil)
	assert.EqualValues(t, finish, false)
	finish, err = l1.ReceivePhase4VerifyTotalPubKey(msg43, 3)
	assert.EqualValues(t, err, nil)
	assert.EqualValues(t, finish, false)
	finish, err = l1.ReceivePhase4VerifyTotalPubKey(msg44, 4)
	assert.EqualValues(t, err, nil)
	assert.EqualValues(t, finish, true)

	//2 添加其他人
	finish, err = l2.ReceivePhase4VerifyTotalPubKey(msg41, 1)
	assert.EqualValues(t, err, nil)
	assert.EqualValues(t, finish, false)
	finish, err = l2.ReceivePhase4VerifyTotalPubKey(msg40, 0)
	assert.EqualValues(t, err, nil)
	assert.EqualValues(t, finish, false)
	finish, err = l2.ReceivePhase4VerifyTotalPubKey(msg43, 3)
	assert.EqualValues(t, err, nil)
	assert.EqualValues(t, finish, false)
	finish, err = l2.ReceivePhase4VerifyTotalPubKey(msg44, 4)
	assert.EqualValues(t, err, nil)
	assert.EqualValues(t, finish, true)

	//3 添加其他人
	finish, err = l3.ReceivePhase4VerifyTotalPubKey(msg41, 1)
	assert.EqualValues(t, err, nil)
	assert.EqualValues(t, finish, false)
	finish, err = l3.ReceivePhase4VerifyTotalPubKey(msg42, 2)
	assert.EqualValues(t, err, nil)
	assert.EqualValues(t, finish, false)
	finish, err = l3.ReceivePhase4VerifyTotalPubKey(msg40, 0)
	assert.EqualValues(t, err, nil)
	assert.EqualValues(t, finish, false)
	finish, err = l3.ReceivePhase4VerifyTotalPubKey(msg44, 4)
	assert.EqualValues(t, err, nil)
	assert.EqualValues(t, finish, true)

	//4 添加其他人
	finish, err = l4.ReceivePhase4VerifyTotalPubKey(msg41, 1)
	assert.EqualValues(t, err, nil)
	assert.EqualValues(t, finish, false)
	finish, err = l4.ReceivePhase4VerifyTotalPubKey(msg42, 2)
	assert.EqualValues(t, err, nil)
	assert.EqualValues(t, finish, false)
	finish, err = l4.ReceivePhase4VerifyTotalPubKey(msg43, 3)
	assert.EqualValues(t, err, nil)
	assert.EqualValues(t, finish, false)
	finish, err = l4.ReceivePhase4VerifyTotalPubKey(msg40, 0)
	assert.EqualValues(t, err, nil)
	assert.EqualValues(t, finish, true)

	p0, err := l0.db.LoadPrivatedKeyInfo(l0.PrivateKeyID)
	p1, err := l1.db.LoadPrivatedKeyInfo(l0.PrivateKeyID)
	p2, err := l2.db.LoadPrivatedKeyInfo(l0.PrivateKeyID)
	p3, err := l3.db.LoadPrivatedKeyInfo(l0.PrivateKeyID)
	p4, err := l4.db.LoadPrivatedKeyInfo(l0.PrivateKeyID)

	//校验私钥是否正确分配
	var xi []share.SPrivKey
	xi = append(xi, p0.XI)
	xi = append(xi, p1.XI)
	xi = append(xi, p2.XI)
	xi = append(xi, p3.XI)
	xi = append(xi, p4.XI)

	totalPrivKey := msg30[1].Vss.Reconstruct([]int{0, 1, 2}, xi[0:3])
	sum := share.PrivKeyZero.Clone()
	share.ModAdd(sum, p0.UI)
	share.ModAdd(sum, p1.UI)
	share.ModAdd(sum, p2.UI)
	share.ModAdd(sum, p3.UI)
	share.ModAdd(sum, p4.UI)

	if sum.D.Cmp(totalPrivKey.D) != 0 {
		t.Error("not equal")
	}
	pubx, puby := share.S.ScalarBaseMult(sum.Bytes())
	if pubx.Cmp(p0.PublicKeyX) != 0 || puby.Cmp(p0.PublicKeyY) != 0 {
		t.Error("pub key error")
	}
	return
}

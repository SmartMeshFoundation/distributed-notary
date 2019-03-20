package pbft

type pbftHistory interface {
	Seq() int
	IsCheckPoint() bool
	Digest() string
}
type HistoryState struct {
	SeqVal    int
	DigestVal string
}

func (hs *HistoryState) Seq() int {
	return hs.SeqVal
}
func (hs *HistoryState) Digest() string {
	return hs.DigestVal
}

type historyOpState struct {
	HistoryState
}

func newHistoryOpState(seq int, digest string) *historyOpState {
	return &historyOpState{
		HistoryState: HistoryState{
			SeqVal:    seq,
			DigestVal: digest,
		},
	}
}
func (hs *historyOpState) IsCheckPoint() bool {
	return false
}

/*
checkpoint state如何计算
借鉴区块链的思路, 假设checkpointdiv是2
s0:=hash<empty>
h1:=hash<op0,op1>
s1:=hash<h1,s0> 第1个checkpoint
h2:=hash<op2,op3>
s2:=hash<h2,s1> 第二个checkpoint
*/
type historyCheckPointState struct {
	HistoryState
}

func newHistoryCheckPointState(seq int, digest string) *historyCheckPointState {
	return &historyCheckPointState{
		HistoryState: HistoryState{
			SeqVal:    seq,
			DigestVal: digest,
		},
	}
}
func (hs *historyCheckPointState) IsCheckPoint() bool {
	return true
}

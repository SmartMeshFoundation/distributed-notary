package pbft

import (
	"encoding/gob"
	"fmt"
	"time"

	"github.com/nkbai/log"

	utils "github.com/nkbai/goutils"
)

var (
	changeViewTimeout = 5000 * time.Millisecond
	checkpointDiv     = 100
)

type ApplySaver interface {
	UpdateSeq(seq int)
}

/*
以客户端发来的Op作为唯一的区分,不论哪个客户端发来的op,只要是一样,就认为是相同的
客户端需要确保op处理完了以后要保存下来,不要很久以前的op
*/
type Server struct {
	id      int //公证人编号
	msgChan chan interface{}
	sender  MessageSender
	tSeq    int // Total sequence number of next request
	//SeqVal               []int              // Sequence number of next client request
	seqmap            map[string]int     // Use to map op to global sequence number for all prepared message 处理过的op应该保存在数据库中
	view              int                //当前view
	apply             int                // Sequence number of last executed request
	log               map[entryID]*entry //客户端的一次请求对应的entry
	cps               map[int]*CheckPoint
	h                 int  //下沿 收到的seq 不能比这个低
	H                 int  //上沿 不能比这个高
	f                 int  //恶意节点数
	monitor           bool //是否监控主节点是否正常
	change            *time.Timer
	changing          bool // Indicate if this node is changing view
	myViewChangeCount int  //我自己第几次连续尝试view change,在view change的过程中,可能有连续的主节点崩溃,从而导致系统无法正常使用
	replicas          []int
	// Deterministic state machine's state 将stable checkpoint hash,这样保证state中的数量是有限的
	// 参考链的设计思路state保存的是<SeqVal,ischeckpoint,DigestVal> 和<SeqVal,notcheckpoint,DigestVal>
	state      []pbftHistory
	applyQueue *PriorityQueue
	vcs        map[int][]*ViewChangeArgs
	lastcp     int //上一个check point
	as         ApplySaver
}

//客户端请求的唯一编号
type entryID struct {
	v int
	n int
}

type entry struct {
	pp         *PrePrepareArgs
	p          []*PrepareArgs //收到的prepare消息,包括自己
	sendCommit bool           //是否已经发送过commit了
	c          []*CommitArgs  //收到的commit消息
	sendReply  bool           //应答是否已经发送过了
	r          *ResponseArgs  //发送给客户端的请求
}

// pendingTask is use to represent the pending task
type pendingTask struct {
	seq  int
	dgt  string
	ent  *entry
	args CommitArgs
}

// Request is use to process incoming clients's request, if this peer is leader then process this request,
// otherwise it will relay the request to leader and start view change timer
func (s *Server) Request(args RequestArgs) error {
	leader := s.view % len(s.replicas)
	if !s.changing {
		if s.id == leader {
			return s.handleRequestLeader(&args)
		}
		return s.handleRequestFollower(&args)
	}
	return nil
}

func (s *Server) handleRequestLeader(args *RequestArgs) error {
	gseq, ok := s.seqmap[args.Op]
	if !ok {
		if s.tSeq >= s.H {
			log.Warn(fmt.Sprintf("%s  cannot receive more request,args=%+v", s, args))
			return nil
		}
		log.Trace("%s[Record/NewRequest]:RequestArgs:%+v\n", s, args)
		// Extract continuous requests from message queue of client and update sequence
		//一次性处理这一个用户的所有请求
		ppArgs := PrePrepareArgs{
			View:    s.view,
			Seq:     s.tSeq,          //从0开始
			Digest:  Digest(args.Op), //op作为唯一的标识,无论来自哪个客户端的重复op,都认为是相同的
			Message: *args,
		}
		s.seqmap[args.Op] = s.tSeq
		s.tSeq++
		log.Trace("%s[B/PrePrepare]:PrePrepareArgs:%+v\n", s, ppArgs)
		s.broadcast(true, newPrePrepareMessage(&ppArgs))

	} else {
		// Old request, if it already finish execute, try to find result and reply again
		// if it haven't done yet, just start again by broadcast
		if gseq <= s.lastcp {
			// The log already been removed
			log.Trace("%s[Request/Fail] log already been removed @ %v, CP:%v", s, gseq, s.lastcp)
		} else {
			// Base on assert : leader will have all preprepare message for request after last checkpoint
			ent := s.getEntry(entryID{s.view, gseq})
			if ent.pp == nil {
				st := fmt.Sprintf("%s old request didn't have preprepare @ %v", s, args.Op)
				log.Trace("st=%s", st)
				return nil
			}
			if ent.pp.Digest == "" {
				log.Trace("%s[Request/Fail]:RequestArgs:%+v", s, args)
			} else {
				if !ent.sendReply { //没有发送过应答,重来协商一遍?
					log.Trace("%s[Re/B/PrePrepare]:PrePrepareArgs:%+v", s, *(ent.pp))
					s.broadcast(true, newPrePrepareMessage(ent.pp))
				} else {
					s.reply(ent)
				}
			}
		}

	}

	return nil
}

func (s *Server) handleRequestFollower(args *RequestArgs) error {
	log.Trace(fmt.Sprintf("%s follower receive request %v", s, args))
	if !s.monitor {
		if s.change != nil { //有可能二次view change正在运行
			s.change.Stop()
			log.Trace(fmt.Sprintf("%s handleRequestFollower stop change=%v", s, s.change))
		}
		//监控主节点是否拒绝服务,不关心主节点是否有错
		s.monitor = true
		//收到主节点的PrePrepare,prePare都停止,只是处理主节点不响应请求这种情况
		log.Trace(fmt.Sprintf("%s startviewchange timer", s))
		s.change = time.AfterFunc(changeViewTimeout, s.startViewChange)
		log.Trace(fmt.Sprintf("%s handleRequestFollower change=%v", s, s.change))
		//转发请求到主节点
		go s.sender.SendMessage(newRequestMessage(args), s.view%len(s.replicas))

	}

	return nil
}

// PrePrepareArgs is the argument for RPC handler Server.PrePrepare
type PrePrepareArgs struct {
	View   int
	Seq    int
	Digest string
	// Signature
	Message RequestArgs
}

// PrePrepare is the handler of RPC Server.PrePrepare
func (s *Server) PrePrepare(args PrePrepareArgs) error {
	if s.changing {
		log.Info(fmt.Sprintf("%s ignore PrePrepare because changing args=%+v", s, args))
		return nil
	}
	s.stopTimer()
	if s.view == args.View && s.h <= args.Seq && args.Seq < s.H {
		log.Trace("%s[R/PrePrepare]:PrePrepareArgs:%+v", s, args)

		ent := s.getEntry(entryID{args.View, args.Seq})
		//没收到过这个prePrepare或者收到重复的PrePrepare
		if ent.pp == nil || (ent.pp.Digest == args.Digest && ent.pp.Seq == args.Seq) {
			pArgs := PrepareArgs{
				View:   args.View,
				Seq:    args.Seq,
				Digest: args.Digest,
				Rid:    s.id,
			}
			ent.pp = &args

			log.Trace("%s[B/Prepare]:PrepareArgs:%+v", s, pArgs)

			s.broadcast(true, newPrepareMessage(&pArgs))
		}
	}
	return nil
}

// PrepareArgs is the argument for RPC handler Server.Prepare
type PrepareArgs struct {
	View   int
	Seq    int
	Digest string
	Rid    int
	// Signature
}

// Prepare is the handler of RPC Server.Prepare
func (s *Server) Prepare(args PrepareArgs) error {
	if s.changing {
		log.Info(fmt.Sprintf("%s ignore Prepare because changing args=%+v", s, args))
		return nil
	}
	s.stopTimer()
	if s.view == args.View && s.h <= args.Seq && args.Seq < s.H {
		ent := s.getEntry(entryID{args.View, args.Seq})
		if ent.pp != nil && ent.pp.Digest != args.Digest {
			return fmt.Errorf("receive Prepare,but DigestVal not match args=%+v,ent=%s", args, utils.StringInterface(ent, 3))
		}
		log.Trace("%s[R/Prepare]:PrepareArgs:%+v", s, args)
		ent.p = append(ent.p, &args)
		//todo 发送commit条件,prepare消息数量够了,或者已经发送过commit, 这会导致commit消息会重复发送几次 如果不重复,在view change的过程中会失败
		if ent.pp != nil && !ent.sendCommit && s.prepared(ent) {
			//避免重复发送commit
			s.seqmap[ent.pp.Message.Op] = args.Seq
			cArgs := CommitArgs{
				View:   args.View,
				Seq:    ent.pp.Seq,
				Digest: ent.pp.Digest,
				Rid:    s.id,
			}
			//每个节点都会发两次
			log.Trace("%s[B/Commit]:CommitArgs:%+v", s, cArgs)
			s.broadcast(true, newCommitMessage(&cArgs))
			ent.sendCommit = true

		}
	}
	return nil
}

//收到了足够的prepared消息
func (s *Server) prepared(ent *entry) bool {
	if len(ent.p) > 2*s.f {
		// Key is the id of sender replica
		validSet := make(map[int]bool)
		for i, sz := 0, len(ent.p); i < sz; i++ {
			if ent.p[i].View == ent.pp.View && ent.p[i].Seq == ent.pp.Seq && ent.p[i].Digest == ent.pp.Digest {
				validSet[ent.p[i].Rid] = true
			}
		}
		return len(validSet) > 2*s.f
	}
	return false
}

// CommitArgs is the argument for  handler Commit
type CommitArgs struct {
	View   int
	Seq    int
	Digest string
	Rid    int
}

// Commit is the handler of Commit
func (s *Server) Commit(args CommitArgs) error {
	if s.changing {
		log.Info(fmt.Sprintf("%s ignore Commit because changing args=%+v", s, args))
		return nil
	}
	s.stopTimer()

	if s.view == args.View && s.h <= args.Seq && args.Seq < s.H {
		ent := s.getEntry(entryID{args.View, args.Seq})
		if ent.pp != nil && ent.pp.Digest != args.Digest {
			return fmt.Errorf("receive commit,but DigestVal not match args=%+v,ent=%s", args, utils.StringInterface(ent, 3))
		}
		ent.c = append(ent.c, &args)
		log.Trace(fmt.Sprintf("%s[R/Commit]:CommitArgs:%+v", s, args))
		if !ent.sendReply && ent.sendCommit && s.committed(ent) {
			log.Trace("%s start execute %v @ %v", s, ent.pp.Message.Op, args.Seq)
			// Execute will make sure there only one execution of one request
			s.execute(args.Seq, ent, args)
		} else {
			//s.reply(ent)
		} //todo 如果确保只reply一次,那会造成view change的时候无法成功

	}
	return nil
}

func (s *Server) committed(ent *entry) bool {
	if len(ent.c) > 2*s.f {
		// Key is replica id
		validSet := make(map[int]bool)
		for i, sz := 0, len(ent.c); i < sz; i++ {
			if ent.c[i].View == ent.pp.View && ent.c[i].Seq == ent.pp.Seq && ent.c[i].Digest == ent.pp.Digest {
				validSet[ent.c[i].Rid] = true
			}
		}
		return len(validSet) > 2*s.f
	}
	return false
}

func (s *Server) reply(ent *entry) {
	if ent.r != nil {
		go func() {
			log.Trace("%s[S/Reply]:ResponseArgs:%+v", s, ent.r)
			s.sender.SendMessage(newResponseMessage(ent.r), ent.pp.Message.ID)
		}()
	} else {
		log.Trace("%s[Trying reply] Fail", s)
	}
}
func (s *Server) replyByPendingTask(pt pendingTask) {
	if pt.ent.r == nil {
		rArgs := ResponseArgs{
			View: pt.args.View,
			Seq:  pt.seq, //集体序号
			Cid:  pt.ent.pp.Message.ID,
			Rid:  s.id,
			Res:  pt.ent.pp.Message.Op,
		}
		pt.ent.r = &rArgs
	}
	if s.as != nil {
		s.as.UpdateSeq(pt.seq)
	}
	s.reply(pt.ent)
}

func (s *Server) execute(seq int, ent *entry, args CommitArgs) {

	if seq <= s.apply {
		// This is an old request, try to find it's result, if can't, return false
		return
	}
	elem := PQElem{
		Pri: seq,
		C: pendingTask{
			seq:  seq,
			dgt:  args.Digest,
			ent:  ent,
			args: args,
		},
	}

	log.Trace("%s[AddQueue] %v @ %v", s, ent.pp.Message.Op, seq)

	inserted := s.applyQueue.Insert(elem)
	if !inserted { //todo 在切换view的过程中会出问题
		panic(fmt.Sprintf("Already insert some request with same sequence SeqVal=%d", seq))
	}
	ent.sendReply = true // 保证不会重复执行 进入队列,后续其他序号补上以后保证会被执行.

	for i, sz := 0, s.applyQueue.Length(); i < sz; i++ {
		m, err := s.applyQueue.GetMin()
		if err != nil {
			break
		}
		//按顺序执行,如果没有执行到,等待上一个序列完成
		if s.apply+1 == m.Pri {
			s.apply++
			pt := m.C.(pendingTask)
			if pt.dgt != "" {
				//state中的要是严格有序的
				s.state = append(s.state, newHistoryOpState(pt.seq, pt.ent.pp.Message.Op))
			}
			//m.C.(pendingTask).done <- m.C.(pendingTask).op
			_, err = s.applyQueue.ExtractMin()
			if err != nil {
				panic(err)
			}
			//发送结果给相应的client
			log.Trace("%s[Execute] Op:%v @ %v", s, pt.ent.pp.Message.Op, pt.seq)
			s.replyByPendingTask(pt)
			if s.apply%checkpointDiv == 0 {
				s.startCheckPoint(s.view, s.apply)
			}
		} else if s.apply+1 > m.Pri {
			panic("This should already done")
		} else {
			//s.Apply< m.Pri 表示有乱序
			break
		}
	}
}

// CheckPoint is the reply of FetchCheckPoint, signature is only set when it transmit by RPC
type CheckPoint struct {
	Seq    int
	Stable bool
	Digest string //所有记录的digest
	View   int
	Proof  []*CheckPointArgs
	// Signature
}

func (s *Server) startCheckPoint(v int, n int) {
	cp := s.getCheckPoint(n)
	// There is no newer stable checkpoint
	if cp != nil {
		cp.Seq = n
		cp.Stable = false
		cpArgs := CheckPointArgs{
			Seq:    n,
			Digest: cp.Digest,
			Rid:    s.id,
			View:   s.view,
		}

		log.Trace("%s[B/CheckPoint]: Seq:%d,digest=%s", s, n, cp.Digest)

		s.broadcast(true, newCheckPointMessage(&cpArgs))
	}
}

// CheckPointArgs is the argument for  handler CheckPoint
type CheckPointArgs struct {
	Seq    int
	Digest string
	Rid    int
	View   int //当前view是多少,为了让主节点能够参与同步
}

/*
以大多数人的checkpoint为准,来更新我自己的checkpoint,我的checkpoint可能是错的
*/
func (s *Server) CheckPoint(args CheckPointArgs) error {
	if args.Seq%checkpointDiv != 0 {
		return fmt.Errorf("receive CheckPointArgs=%s,but SeqVal is error", utils.StringInterface(args, 3))
	}
	cp := s.getCheckPoint(args.Seq)
	if cp != nil {

		if !cp.Stable {
			cp.Proof = append(cp.Proof, &args)

			log.Trace("%s[R/CheckPoint]:Args:%+v,Total:%d", s, args, len(cp.Proof))

			if len(cp.Proof) > 2*s.f {
				//因为不超过f个恶意节点,所有可以确定有超过2f个节点会有相同的digest以及view,以他们的为准就行了.
				countDigest := make(map[string]int)
				countView := make(map[int]int)
				for i, sz := 0, len(cp.Proof); i < sz; i++ {
					countDigest[cp.Proof[i].Digest]++
					countView[cp.Proof[i].View]++
					// Stablize checkpoint,我自己的也在里面,如果我自己的state不在里面,后面将不会在同步.
					/*
						同步checkpoint的过程比较麻烦
						虽然有超过2f+1个达成一致了,但是我本地的并没有和大多数人保持一致怎么办?
						目前这个做法是假设我在大多数人里面,但是如果state和别人不一致,我的digest也和别人不一致,
						简单同步别人的信息,更新了stable checkpoint也没啥意义
					*/
					if countDigest[cp.Proof[i].Digest] > 2*s.f && countView[cp.Proof[i].View] > 2*s.f {
						log.Trace("%s[Stablize]:Seq:%d,Digest:%s", s, args.Seq, cp.Proof[i].Digest)
						s.stablizeCP(s.cps[args.Seq], cp.Proof[i])
						break
					}
				}
			}
		}
	}

	return nil
}

// Make sure this is the newest checkpoint before call this function
func (s *Server) stablizeCP(cp *CheckPoint, proof *CheckPointArgs) {
	log.Trace("%s[Update]:Checkpoint:{%v,%v}", s, cp.Seq, cp.Stable)
	if proof.View > s.view {
		//这种方式未必合理,在checkpoint中带着view,主要是用于节点重启以后和其他人view不同步的问题
		s.view = proof.View
		s.changing = false
		s.stopTimer()
	}
	cp.Stable = true
	cp.Digest = proof.Digest //以proof的为准,有可能我自己的checkpoint是错的
	if s.apply < proof.Seq {
		s.apply = cp.Seq
	}
	s.removeOldLog(cp.Seq)
	s.removeOldCheckPoint(cp.Seq)
	s.removeOldState(cp.Seq)
	s.removeOldSeqMap(cp.Seq - 2*checkpointDiv) //不能清理seqmap太及时了,否则在checkpoint临界值会导致为统一op分配不同的seq,
	s.insertStateAt0(newHistoryCheckPointState(proof.Seq, proof.Digest))
	s.h = cp.Seq
	s.H = s.h + 2*checkpointDiv
	for i, sz := 0, s.applyQueue.Length(); i < sz; i++ {
		m, err := s.applyQueue.GetMin()
		if err != nil {
			break
		}
		//对于由于某些原因我这边没有同步的request,直接忽略,因为其他超过2/3的节点已经认可了.
		if m.Pri <= cp.Seq {
			m, err = s.applyQueue.ExtractMin()
			if err != nil {
				panic(err)
			}
			s.replyByPendingTask(m.C.(pendingTask))
		} else {
			break
		}
	}
	s.lastcp = cp.Seq
}

/*如果主节点作恶，它可能会给不同的请求编上相同的序号，或者不去分配序号，或者让相邻的序号不连续。
备份节点应当有职责来主动检查这些序号的合法性。如果主节点掉线或者作恶不广播客户端的请求，客户端设置超时机制，
超时的话，向所有副本节点广播请求消息。副本节点检测出主节点作恶或者下线，发起View Change协议。

只有检测到主节点不作为,没有检测主节点乱作为问题
*/
func (s *Server) startViewChange() {
	log.Trace("%s[StartViewChange]", s)
	s.msgChan <- &internalTimerMessage{}
}

func (s *Server) handleTimer() {
	log.Trace(fmt.Sprintf("%s[handleTimer]", s))
	vcArgs := s.generateViewChange()
	s.myViewChangeCount++ //下次就进入新的view了

	log.Trace("%s[B/ViewChange]:Args:%+v,cp=%s", s, *vcArgs, utils.StringInterface(vcArgs.CP, 2))
	s.broadcast(true, newViewChangeMessage(vcArgs))
	//再次超时启动start view change,针对连续主节点崩溃问题
	s.change = time.AfterFunc(changeViewTimeout, s.startViewChange)
	log.Trace(fmt.Sprintf("%s handleTimer change=%v", s, s.change))
	s.monitor = true
}

func (s *Server) generateViewChange() *ViewChangeArgs {
	s.changing = true
	cp := s.getStableCheckPoint()

	vcArgs := ViewChangeArgs{
		View: s.view + s.myViewChangeCount + 1,
		//View: s.view + 1,
		Rid: s.id,
		CP: &CheckPointArgs{
			Seq:    cp.Seq,
			View:   cp.View,
			Rid:    s.id,
			Digest: cp.Digest,
		},
	}

	/*
		   副本节点向其他节点广播<VIEW-CHANGE, v+1, n, C, P, i>消息。n是最新的stable checkpoint的编号，C是2f+1验证过的CheckPoint消息集合，
		按照论文中的意思:
		要广播哪些还没有stable的checkpoint以后的所有已经达成工共识的交易
		也就是说大部分人已经同意了这个序号,无论有没有发送reply给client.
	*/
	for k, v := range s.log {
		if k.n > cp.Seq && v.pp != nil && (v.sendCommit || s.prepared(v)) {
			pm := Pm{
				PP: v.pp,
				P:  v.p,
			}
			vcArgs.P = append(vcArgs.P, &pm)
		}
	}

	return &vcArgs
}

// Pm is use to hold a preprepare message and at least 2f corresponding prepare message
type Pm struct {
	PP *PrePrepareArgs
	P  []*PrepareArgs
}

// ViewChangeArgs is the argument for  handler ViewChange
type ViewChangeArgs struct {
	View int
	CP   *CheckPointArgs
	P    []*Pm
	Rid  int
}

// ViewChange is the handler of  ViewChange
func (s *Server) ViewChange(args ViewChangeArgs) error {

	log.Trace("%s[R/ViewChange]:Args:%+v", s, args)

	// Ignore old viewchange message
	if args.View <= s.view {
		return nil
	}

	// Insert this view change message to its log ,todo 对于来自同一rid的重复的viewchange,是否应该替换,还有我如何验证viewchange的有效性呢? checkpoint,pm这些信息是否正确?
	// 看了sawtooth-pbft也是一样没变
	s.vcs[args.View] = append(s.vcs[args.View], &args)
	//不验证重复消息么?
	// Leader entering new view
	/*
				当主节点p = v + 1 mod |R|收到2f个有效的VIEW-CHANGE消息后，向其他节点广播<NEW-VIEW, v+1, V, O>消息。
			V是有效的VIEW-CHANGE消息集合。O是主节点重新发起的未经完成的PRE-PREPARE消息集合。PRE-PREPARE消息集合的选取规则：

			1. 选取V中最小的stable checkpoint编号min-s，选取V中prepare消息的最大编号max-s。

			2. 在min-s和max-s之间，如果存在P消息集合，则创建<<PRE-PREPARE, v+1, n, d>, m>消息。否则创建一个空的PRE-PREPARE消息，
		即：<<PRE-PREPARE, v+1, n, d(null)>, m(null)>, m(null)空消息，d(null)空消息摘要。

			副本节点收到主节点的NEW-VIEW消息，验证有效性，有效的话，进入v+1状态，并且开始O中的PRE-PREPARE消息处理流程。
	*/
	if (args.View%len(s.replicas) == s.id) && len(s.vcs[args.View]) >= 2*s.f+1 { //todo 都不用做去重? 这样明显感觉有问题
		nvArgs := NewViewArgs{
			View: args.View,
			V:    s.vcs[args.View],
		}
		//nvArgs.V = append(nvArgs.V, s.generateViewChange())
		//mins 是checkpoint的seq
		//maxs 是preprepared的seq
		mins, maxs, pprepared := s.calcMinMaxspp(&nvArgs)
		pps := s.calcPPS(args.View, mins, maxs, pprepared)

		nvArgs.O = pps

		log.Trace("%s[B/NewView]:Args:%+v", s, nvArgs)

		s.broadcast(false, newNewViewMessage(&nvArgs))

		s.enteringNewView(&nvArgs, mins, maxs, pps)
	}

	return nil
}

func (s *Server) enteringNewView(nvArgs *NewViewArgs, mins int, maxs int, pps []*PrePrepareArgs) []*PrepareArgs {
	log.Trace("%s[EnterNextView]:%v,mins=%d,maxs=%d", s, nvArgs.View, mins, maxs)

	scp := s.getStableCheckPoint()
	//如果我自己的checkpoint比最低的还低,按照最新的来就好了.
	if mins > scp.Seq {
		for i, sz := 0, len(nvArgs.V); i < sz; i++ {
			if nvArgs.V[i].CP.Seq == mins {
				cp2 := nvArgs.V[i].CP
				s.cps[mins] = &CheckPoint{
					Seq:    cp2.Seq,
					Stable: true,
					Digest: cp2.Digest,
					View:   cp2.View,
				}
				break
			}
		}
		cp := s.cps[mins]
		s.stablizeCP(cp, &CheckPointArgs{Seq: cp.Seq, Rid: s.id, Digest: cp.Digest, View: s.view}) //这里checkpoint的state岂不是越来越大了?每次消息发送那么复杂的数据,像测试中这种实现,还可能互相影响彼此的state,因为是符号拷贝
	}

	s.tSeq = maxs + 1
	ps := make([]*PrepareArgs, len(pps))
	for i, sz := 0, len(pps); i < sz; i++ {
		s.seqmap[pps[i].Message.Op] = pps[i].Seq
		ent := s.getEntry(entryID{nvArgs.View, pps[i].Seq})
		ent.pp = pps[i]

		pArgs := PrepareArgs{
			View:   nvArgs.View,
			Seq:    pps[i].Seq,
			Digest: pps[i].Digest,
			Rid:    s.id,
		}
		ps[i] = &pArgs
	}
	s.view = nvArgs.View
	s.removeOldLog(mins)
	s.stopTimer()
	s.changing = false
	s.removeOldViewChange(nvArgs.View)

	/*
		副本节点收到主节点的NEW-VIEW消息，验证有效性，有效的话，进入v+1状态，并且开始O(pps)中的PRE-PREPARE消息处理流程。
	*/
	go func() {
		for i, sz := 0, len(ps); i < sz; i++ {
			log.Trace("%s[B/Prepare]:Args:%+v", s, ps[i])
			s.broadcast(true, newPrepareMessage(ps[i]))
			time.Sleep(5 * time.Millisecond)
		}
	}()

	return ps
}

/*
	1. 选取V中最小的stable checkpoint编号min-s，选取V中prepare消息的最大编号max-s。

			2. 在min-s和max-s之间，如果存在P消息集合，则创建<<PRE-PREPARE, v+1, n, d>, m>消息。否则创建一个空的PRE-PREPARE消息，
		即：<<PRE-PREPARE, v+1, n, d(null)>, m(null)>, m(null)空消息，d(null)空消息摘要。
*/
func (s *Server) calcMinMaxspp(nvArgs *NewViewArgs) (int, int, map[int]*PrePrepareArgs) {
	mins, maxs := -1, -1
	pprepared := make(map[int]*PrePrepareArgs)
	for i, sz := 0, len(nvArgs.V); i < sz; i++ {
		if nvArgs.V[i].CP.Seq > mins {
			mins = nvArgs.V[i].CP.Seq
		}
		for j, psz := 0, len(nvArgs.V[i].P); j < psz; j++ { //这里的问题就在于我如何从众多view change中找出prepared集合. 这里只是加定信任.
			if nvArgs.V[i].P[j].PP.Seq > maxs {
				maxs = nvArgs.V[i].P[j].PP.Seq
			}
			//以哪个的preprepared为准呢?还是说到了这一步,确保prepared都是有效的?
			pprepared[nvArgs.V[i].P[j].PP.Seq] = nvArgs.V[i].P[j].PP
		}
	}
	if maxs < mins {
		maxs = mins //checkpoint都确认了,那么可能造成max是负数,从而导致后续错误
	}
	return mins, maxs, pprepared
}

/*

		2. 在min-s和max-s之间，如果存在P消息集合，则创建<<PRE-PREPARE, v+1, n, d>, m>消息。否则创建一个空的PRE-PREPARE消息，
	即：<<PRE-PREPARE, v+1, n, d(null)>, m(null)>, m(null)空消息，d(null)空消息摘要。

*/
func (s *Server) calcPPS(view int, mins int, maxs int, pprepared map[int]*PrePrepareArgs) []*PrePrepareArgs {
	pps := make([]*PrePrepareArgs, 0)
	for i := mins + 1; i <= maxs; i++ {
		v, ok := pprepared[i]
		if ok {
			pps = append(pps, &PrePrepareArgs{
				View:    view,
				Seq:     i,
				Digest:  v.Digest,
				Message: v.Message,
			})
		} else {
			pps = append(pps, &PrePrepareArgs{
				View:   view,
				Seq:    i,
				Digest: "",
			})
		}
	}
	return pps
}

// NewViewArgs is the argument for  handler NewView
type NewViewArgs struct {
	View int
	V    []*ViewChangeArgs
	O    []*PrePrepareArgs
}

// NewView is the handler of NewView
func (s *Server) NewView(args NewViewArgs) error {
	// Verify signature
	if args.View <= s.view {
		return nil
	}
	log.Trace("%s[R/NewView]:Args:%+v", s, args)

	// Verify V sest
	vcs := make(map[int]bool)
	for i, sz := 0, len(args.V); i < sz; i++ {
		if args.V[i].View == args.View {
			vcs[args.V[i].Rid] = true
		}
	}
	if len(vcs) <= 2*s.f {
		log.Trace("%s[V/NewView/Fail] view change message is not enough", s)
		return nil
	}

	// Verify O set
	mins, maxs, pprepared := s.calcMinMaxspp(&args)
	pps := s.calcPPS(args.View, mins, maxs, pprepared)

	for i, sz := 0, len(pps); i < sz; i++ {
		if pps[i].View != args.O[i].View || pps[i].Seq != args.O[i].Seq || pps[i].Digest != args.O[i].Digest {
			log.Trace("%s[V/NewView/Fail] PrePrepare message missmatch : %+v", s, pps[i])
			return nil
		}
	}

	s.enteringNewView(&args, mins, maxs, pps)

	log.Trace("%s[NowInNewView]:%v", s, args.View)

	return nil
}

func (s *Server) getEntry(id entryID) *entry {
	_, ok := s.log[id]
	if !ok {
		s.log[id] = &entry{
			pp:         nil,
			p:          make([]*PrepareArgs, 0),
			sendCommit: false,
			c:          make([]*CommitArgs, 0),
			sendReply:  false,
			r:          nil,
		}
	}
	return s.log[id]
}

var invalidState = newHistoryOpState(-1, Digest(""))

func (s *Server) getState(seq int) pbftHistory {
	for _, state := range s.state {
		if state.Seq() == seq {
			return state
		}
	}
	log.Warn(fmt.Sprintf("%s getState %d,but cannot found", s, seq))
	return invalidState
}
func (s *Server) insertStateAt0(state pbftHistory) {
	s.state = append([]pbftHistory{state}, s.state...)
}

/*
包含from,不包含to
有可能出现如下情况:
1. 收到来自别人的checkpoint,但是我自己to还没有走到checkpoint这一步, 这种情况后续计算会保证更新
2. 因为自己意外重启,导致丢了很多消息,从而from确实,但是从from+某个值到to是完整的
3. 软件升级,到时所有的节点from到to不完整,但是从from+某个值到to是完整的
*/
func (s *Server) getStateRange(from, to int) []pbftHistory {
	//s.state中的是严格单增且顺序的.
	var i = -1
	var state pbftHistory
	for i, state = range s.state {
		/*
			就算是没有完整的from-to之间所有的state,有部分state也算数
		*/
		if state.Seq() >= from && state.Seq() < to {
			break
		}
	}
	if i < 0 || i == len(s.state)-1 {
		log.Trace(fmt.Sprintf("%s get StateRange %d-%d,but cannot found", s, from, to))
		return nil
	}
	end := i + to - from
	if end > len(s.state) {
		log.Trace(fmt.Sprintf("%s get StateRange end is too large ,end=%d,from=%d,to=%d", s, end, from, to))
		end = len(s.state)
	}
	return s.state[i:end]
}

/*
收到来自其他节点的checkpoint或者我想发送checkpoint的时候调用.
*/
func (s *Server) getCheckPoint(seq int) *CheckPoint {
	_, ok := s.cps[seq]
	if !ok {
		//收到了checkpoint,但是我已经stable,就不更新了
		for k, v := range s.cps {
			if k > seq && v.Stable {
				return nil
			}
		}
		s.cps[seq] = &CheckPoint{
			Seq:    seq,
			Stable: false,
			Digest: "",
			Proof:  make([]*CheckPointArgs, 0),
		}
	}
	cp := s.cps[seq]
	//多算几次digest,因为有可能我先收到别人的checkpoint,这时候我算出来自己的digest是不完整的.
	switch {
	case seq < 0:
	//没有历史交易,直接忽略.
	case seq == 0:
		cp.Digest = Digest2(s.getState(seq-1).Digest(), s.getState(seq).Digest())
	//第一笔交易
	case seq >= checkpointDiv && seq%checkpointDiv == 0:
		//计算hash<lastCheckpointDigest,SeqVal-checkpointDiv+1,...SeqVal>
		opDigest := Digest(s.getStateRange(seq-checkpointDiv+1, seq+1))
		cp.Digest = Digest2(s.getState(seq-checkpointDiv).Digest(), opDigest)
	default:
		panic("unkown checkpoint ")

	}
	return cp
}

func (s *Server) getStableCheckPoint() *CheckPoint {
	for _, v := range s.cps {
		if v.Stable {
			return v
		}
	}
	panic("No stable checkpoint")
}

func (s *Server) removeOldCheckPoint(seq int) {
	for k, v := range s.cps {
		if v.Seq < seq {
			delete(s.cps, k)
		}
	}
}

func (s *Server) removeOldLog(seq int) {
	for k := range s.log {
		if k.n < seq {
			delete(s.log, k)
		}
	}
}

func (s *Server) removeOldViewChange(seq int) {
	for k := range s.vcs {
		if k < seq {
			delete(s.vcs, k)
		}
	}
}

func (s *Server) removeEntry(id entryID) {
	delete(s.log, id)
}
func (s *Server) removeOldState(seq int) {
	var findex = -1
	for i, state := range s.state {
		if state.Seq() == seq {
			findex = i
			break
		} else if state.Seq() > seq {
			break
		}
	}
	if findex >= 0 {
		l := len(s.state)
		s.state = s.state[findex+1 : l]
	}
}
func (s *Server) removeOldSeqMap(seq int) {
	for k, v := range s.seqmap {
		if v <= seq {
			delete(s.seqmap, k)
		}
	}
}
func (s *Server) stopTimer() {
	if s.change != nil {
		log.Trace(fmt.Sprintf("%s stopTimer change=%v", s, s.change))
		s.change.Stop()
		s.change = nil
	}
	s.monitor = false
	s.myViewChangeCount = 0 //重置为0,不再继续增长
}

func (s *Server) String() string {
	return fmt.Sprintf("{S:ID:%d,tSeq:%d,View:%d,CP:%v,Apply:%v,h:%v,H:%v}", s.id, s.tSeq, s.view, s.lastcp, s.apply, s.h, s.H)
}

//sc: 是否同步,等请求处理完了再返回
//toself: 广播是否给自己
func (s *Server) broadcast(toself bool, req interface{}) {

	for i, sz := 0, len(s.replicas); i < sz; i++ {
		if toself && s.id == i {
			s.handleReq(req)
		} else {
			go func(rid int) {
				s.sender.SendMessage(req, rid)
			}(i)
		}
	}

}
func (s *Server) isPrimary() bool {
	return s.view%len(s.replicas) == s.id
}
func (s *Server) handleReq(req interface{}) {
	var err error
	id := utils.RandomString(10)
	primary := "primary"
	if !s.isPrimary() {
		primary = "follower"
	}
	log.Trace("%s %s receive call %+v  id=%s", s, primary, req, id)
	switch r := req.(type) {
	case *RequestMessage:
		log.Trace("args=%+v", r.Arg.Op)
		err = s.Request(*r.Arg)
	case *PrePrepareMessage:
		err = s.PrePrepare(*r.Arg)
	case *PrepareMessage:
		err = s.Prepare(*r.Arg)

	case *CommitMessage:
		//如果seq执行不是严格按照顺序,则有可能阻塞,所以必须并发
		err = s.Commit(*r.Arg)
	case *CheckPointMessage:

		err = s.CheckPoint(*r.Arg)

	case *ViewChangeMessage:
		err = s.ViewChange(*r.Arg)
	case *NewViewMessage:
		err = s.NewView(*r.Arg)
	case *internalTimerMessage:
		s.handleTimer()
	default:
		panic(fmt.Sprintf("unkown request %+v", r))

	}
	if err != nil {
		log.Error(fmt.Sprintf("request %s complete err=%s", id, err))
	} else {
		log.Trace("request %s complete,err=%v", id, err)
	}
}
func (s *Server) loop() {

	for {
		r := <-s.msgChan
		s.handleReq(r)
	}
}

// NewPBFTServer use given information create a new pbft server
func NewPBFTServer(rid, f, initSeq int, msgChan chan interface{}, sender MessageSender, nodes []int, as ApplySaver) *Server {
	s := Server{
		id:         rid,
		msgChan:    msgChan,
		tSeq:       initSeq, //可以支持tSeq指定初始化
		seqmap:     make(map[string]int),
		view:       0,
		apply:      initSeq - 1,
		log:        make(map[entryID]*entry),
		cps:        make(map[int]*CheckPoint),
		h:          initSeq,
		H:          initSeq + 2*checkpointDiv,
		f:          f,
		monitor:    false,
		change:     nil,
		changing:   false,
		replicas:   nodes,
		sender:     sender,
		state:      make([]pbftHistory, 1),
		applyQueue: NewPriorityQueue(),
		vcs:        make(map[int][]*ViewChangeArgs),
		lastcp:     -1,
		as:         as,
	}
	s.state[0] = invalidState
	// Put an initial stable checkpoint
	cp := s.getCheckPoint(-1)
	cp.Stable = true
	cp.Digest = s.state[0].Digest()
	go s.loop()
	return &s
}

func init() {
	gob.Register(historyOpState{})
	gob.Register(historyCheckPointState{})
	gob.Register(HistoryState{})
}

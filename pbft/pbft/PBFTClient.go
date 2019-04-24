package pbft

import (
	"fmt"
	"time"

	"errors"

	"github.com/nkbai/log"
)

type ResponseFunc func(op string, res interface{})

var (
	retransmitTimeout = 500 * time.Millisecond //注意retransmitTimeout 和view change的时间关系,必须远小于view change time out
	outstanding       = 10                     //client请求并发数量
)

type RequestArgs struct {
	Op        string //也是cmd
	ID        int    //client 唯一编号
	Timestamp time.Time
	// Signature
}

type RequestReply struct {
	Err string
}
type OpResult struct {
	Cmd       string
	Auxiliary string
	Seq       int
	Error     error
}
type cEntry struct {
	apply bool //是否已经收到足够的reply
	req   *RequestMessage
	res   []*ResponseArgs //收到的响应,可能有不一致的时候
	t     *time.Timer
	ft    bool //单点请求还是广播请求公证人 第一次请求是单点发送,如果失败,后续尝试会广播
}

/*
Client 控制发送速度,一次不要发送太多的请求,否则会导致反复重试,从而堵塞业务处理.
*/
type Client struct {
	id          int //当前client编号
	msgChan     chan interface{}
	view        int   //主要是指明下次请求发给谁,没有其他更多意义
	replicas    []int //所有的公证人节点列表
	f           int   // Max number of fail node
	log         map[string]*cEntry
	reqQueue    *StringQueue  // Use to send request to sendRequest routine ,处理完毕一个,再压入下一个
	backupQueue *StringQueue  // Use to receive request but not in reqQueue
	Apply       chan OpResult //通知某个cmd完成
	retry       chan *cEntry
	count       int
	sender      MessageSender
}

func (c *Client) loop() {
	for {
		//不能阻塞
		select {
		case r := <-c.msgChan:
			c.call(r)
		case ent := <-c.retry: //重试的不放在backup中,直接进入队列
			c.sendRequest(ent)
		}
	}

}
func (c *Client) call(r interface{}) {
	var err error
	switch r2 := r.(type) {
	case *StartMessage:
		err = c.startInternal(r2)
		if err != nil {
			c.notifyResult(r2.Arg, "", -1, err)
		}
	case *ResponseMessage:
		err = c.Response(*r2.Arg)
		if err != nil {
			log.Error(fmt.Sprintf("%s client receive response err=%s", c, err))
		}

	default:
		panic(fmt.Sprintf("unkown req %v", r))
	}
}
func (c *Client) startInternal(msg *StartMessage) error {
	log.Trace("%s[Start]:Cmd:%s\n", c, msg.Arg)

	args := RequestArgs{
		Op:        msg.Arg,
		ID:        c.id,
		Timestamp: time.Now(),
	}
	ent, exist := c.newCEntry(args.Op)
	if exist {
		return errors.New("cmd duplicate")
	}
	ent.req = newRequestMessage(&args)

	if c.reqQueue.Length() <= outstanding {
		log.Trace("%s[Add2Req]%+v", c, args)
		c.reqQueue.Insert(msg.Arg, ent)
		go func() {
			// Notify sendRequest to send request
			c.retry <- ent
		}()
	} else {
		log.Trace("%s[Add2Backup]%+v", c, args)
		c.backupQueue.Insert(msg.Arg, ent)
	}

	return nil
}

func (c *Client) Start(cmd string) {
	c.msgChan <- newStartMessage(cmd)
	return

}

// ResponseArgs is the argument for  handler  Response
type ResponseArgs struct {
	View      int    //处理此次的cmdview是多少
	Seq       int    //server协商出来,在所有client请求中的排序结果
	Cid       int    // Client id
	Rid       int    // Replica id
	Res       string //op request中的op
	Auxiliary string //对于以太坊系列来说可有可无,对于btc系列来说包含的是可用outpoint
	// Signature
}

func (c *Client) notifyResult(cmd, auxiliary string, seq int, err error) {
	log.Trace(fmt.Sprintf("%s cmd=%s,SeqVal=%d complete", c, cmd, seq))
	select {
	case c.Apply <- OpResult{
		Cmd:       cmd,
		Auxiliary: auxiliary,
		Seq:       seq,
		Error:     err,
	}:
	default:
		log.Trace("notify result cmd:%s,SeqVal:%d,err=%v fail", cmd, seq, err)
	}
	delete(c.log, cmd) //清理log
}

// Response is the handler of  Response
func (c *Client) Response(args ResponseArgs) (err error) {
	ent, exist := c.log[args.Res]
	if !exist {
		log.Info(fmt.Sprintf("receive response args=%+v,but not found entry,may duplicate response", args))
		return nil
	}
	//更新下次请求的view,错了也没关系,大不了重试一次就好了
	if args.View > c.view {
		c.view = args.View
	}

	log.Trace("%s[R/Response]:Args:%+v", c, args)

	if ent.t != nil {
		ent.t.Stop()
	}
	//收集够了足够的reply,后续的会被忽略
	if !ent.apply {
		ent.res = append(ent.res, &args)
		//有f+1个节点达成了共识(Res,Seq相同,并且是不同的Rid)
		if len(ent.res) > c.f {
			// Map result to replicas list who send that result
			count := make(map[string][]int) //op(cmd)--> 已reply的公证人集合
			for i, sz := 0, len(ent.res); i < sz; i++ {
				key := fmt.Sprintf("%s-%d", ent.res[i].Res, ent.res[i].Seq)
				count[key] = append(count[key], ent.res[i].Rid)
				if len(count[key]) > c.f && DifferentElemInSlice(count[key]) > c.f {
					ent.apply = true
					if !c.reqQueue.Remove(ent.req.Arg.Op) {
						panic("must exist")
					}
					c.updateQueue(args.Seq)
					c.notifyResult(args.Res, args.Auxiliary, args.Seq, nil)
					break
				}
			}
		}
	}
	return
}

// Update response queue and move some request from backupQueue to reqQueue
func (c *Client) updateQueue(seq int) {
	// Move request from backup queue to request queue
	for {
		key, elem, err := c.backupQueue.GetMin()
		if err != nil {
			// Nothing in backup queue
			break
		}
		if c.reqQueue.Length() <= outstanding {
			log.Trace("%s[MoveBackup2Req]@%v", c, key)
			_, _, err := c.backupQueue.ExtractMin()
			if err != nil {
				panic(err)
			}
			c.reqQueue.Insert(key, elem)
			go func() { // reqQueue队列中有新成员,通知即可,保证reqQueue中的任务处理完了再进新的任务
				c.retry <- elem.(*cEntry)
			}()
		} else {
			break
		}
	}
}

func (c *Client) sendRequest(ent *cEntry) {

	if ent.t != nil {
		ent.t.Stop()
	}
	if !ent.apply { //考虑到重试,有可能重试进入队列的请求,取出来的时候已经成功了.
		ent.t = time.AfterFunc(retransmitTimeout, func() {
			log.Trace("%s[Retransmit]%+v", c, ent.req)
			c.retry <- ent
		})

		if ent.ft {
			log.Trace("%s[S/Request]:Args:%+v", c, *(ent.req))

			ent.ft = false //在一段时间内没有得到响应,则下次广播该请求
			leader := c.view % len(c.replicas)
			go c.sender.SendMessage(ent.req, leader)
		} else {
			log.Trace("%s[B/Request]:Args:%+v", c, *(ent.req))

			for i, sz := 0, len(c.replicas); i < sz; i++ {
				go func(rid int) {
					c.sender.SendMessage(ent.req, rid)
				}(i)
			}
		}
	}
}

func (c *Client) String() string {
	s := fmt.Sprintf("{C:ID:%d,View:%d}", c.id, c.view)
	return s
}

func (c *Client) newCEntry(op string) (ent *cEntry, alreadyExist bool) {
	_, ok := c.log[op]
	if !ok {
		c.log[op] = &cEntry{
			apply: false,
			req:   nil,
			res:   make([]*ResponseArgs, 0),
			t:     nil,
			ft:    true,
		}
	}
	return c.log[op], ok
}

// NewPBFTClient use given information to create a new pbft client
func NewPBFTClient(cid int, msgChan chan interface{}, sender MessageSender, f int, nodes []int) *Client {
	c := Client{
		id:          cid,
		view:        0,
		f:           f,
		log:         make(map[string]*cEntry),
		reqQueue:    NewStringQueue(),
		backupQueue: NewStringQueue(),
		Apply:       make(chan OpResult, outstanding*10),
		retry:       make(chan *cEntry, outstanding),
		msgChan:     msgChan,
		sender:      sender,
		replicas:    nodes,
	}
	go c.loop()
	return &c
}

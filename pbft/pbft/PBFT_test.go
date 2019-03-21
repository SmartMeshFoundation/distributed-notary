package pbft

import (
	"fmt"
	"net/http"
	"os"

	"github.com/stretchr/testify/require"

	"sort"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"

	_ "net/http/pprof"

	"github.com/nkbai/log"
)

func init() {
	log.Root().SetHandler(log.LvlFilterHandler(log.LvlTrace, log.DefaultStreamHandler(os.Stderr)))
	retransmitTimeout = 500 * time.Millisecond //注意retransmitTimeout 和view change的时间关系,必须远小于view change time out
	outstanding = 1                            //client请求并发数量,要小于checkpointdiv,否则会反复重试,导致卡死
	changeViewTimeout = 1000 * time.Millisecond
	checkpointDiv = 2 //checkpoint不能是1,否则会出现主节点连续broadcast checkpoint,还没来得及处理其他人的checkpoint,从而导致错误
}

type config struct {
	useMockServer   bool
	cn              int //client 总数
	sn              int //公证人节点数量
	n               int //公证人加client总和
	f               int //恶意节点数量
	servers         []interface{}
	clients         []*Client
	serverAddresses []chan interface{}
	clientAddresses []chan interface{}
	hub             *MessageHub
	initSeq         int //协商的初始seq,可以不从0开始,考虑到集体重启,以及版本升级问题.
	ni              *Network
}

func testBasic(t *testing.T, cnum int, reqnum int, initSeq int, pause time.Duration, useMockServer bool) {
	go func() {
		http.ListenAndServe("localhost:6060", nil)
	}()
	seqm := make(map[int]string)
	sl := sync.Mutex{}
	assertNotDuplicate := func(seq int, cmd string) {
		sl.Lock()
		ocmd, ok := seqm[seq]
		if ok {
			panic(fmt.Sprintf("SeqVal %d already assigned to %s,new =%s", seq, ocmd, cmd))
		}
		seqm[seq] = cmd
		sl.Unlock()
	}
	defer func() {
		if err := recover(); err != nil {
			log.Trace("ccc")
		}
	}()
	c := config{
		useMockServer: useMockServer,
	}
	log.Trace("aa")
	c.init(cnum, 4, initSeq)
	log.Trace("testbasic init...")
	wg := sync.WaitGroup{}

	for i := 0; i < c.cn; i++ {
		wg.Add(1)
		go func(cid int) {
			for j := 0; j < reqnum; j++ {
				c.clients[cid].Start(strconv.Itoa(j) + "-" + strconv.Itoa(cid))
				time.Sleep(pause)
			}
			wg.Done()
		}(i)
		wg.Add(1)
		go func(cid int) {
			for j := 0; j < reqnum; j++ {
				res := <-c.clients[cid].Apply
				if res.Error != nil {
					panic(fmt.Sprintf("res error  %+v", res))
				}
				assertNotDuplicate(res.Seq, res.Cmd)
				log.Trace(fmt.Sprintf("%s:[Apply]:%s-%d\n", c.clients[cid], res.Cmd, res.Seq))
			}
			wg.Done()
		}(i)
	}

	wg.Wait()
	time.Sleep(time.Second)
}

func TestOne(t *testing.T) {
	testBasic(t, 1, 1, 0, 2*time.Millisecond, false)
}
func TestTwo(t *testing.T) {
	testBasic(t, 2, 3*checkpointDiv, 0, 2*time.Millisecond, false)
}

//测试跨越checkpoint边界,但是起始seq不为0的情况
func TestInitSeq(t *testing.T) {
	testBasic(t, 2, checkpointDiv, 3, 2*time.Millisecond, false)
}
func TestMultiple(t *testing.T) {
	testBasic(t, 10, checkpointDiv*10, 0, 2*time.Millisecond, false)
}

func TestMultiClientLowLoad(t *testing.T) {
	changeViewTimeout = 10000 * time.Millisecond
	testBasic(t, 2, 1000*checkpointDiv, 0, 2*time.Millisecond, false)
}

func TestDisconnectFollower(t *testing.T) {
	testDisconnect(t, 1, checkpointDiv*2, 2*time.Millisecond)
}

func splitres(res interface{}) (cid int, seq int) {
	var err error
	ss := strings.Split(res.(string), "-")
	if len(ss) != 2 {
		panic("res error")
	}
	cid, err = strconv.Atoi(ss[1])
	if err != nil {
		panic(res)
	}
	seq, err = strconv.Atoi(ss[0])
	if err != nil {
		panic(res)
	}
	return
}
func TestDisconnectLeader(t *testing.T) {
	seqm := make(map[int]string)
	sl := sync.Mutex{}
	assertNotDuplicate := func(seq int, cmd string) {
		sl.Lock()
		ocmd, ok := seqm[seq]
		if ok {
			panic(fmt.Sprintf("SeqVal %d already assigned to %s,new =%s", seq, ocmd, cmd))
		}
		seqm[seq] = cmd
		sl.Unlock()
	}
	l := sync.Mutex{}
	m := make(map[int][]int)
	c := config{useMockServer: false}
	cn := 1
	c.init(cn, 4, 0)
	totalNum := checkpointDiv * 3
	wg := sync.WaitGroup{}

	wg.Add(totalNum * cn)
	go func() {
		for i := 0; i < cn; i++ {
			go func(index int) {
				for j := 0; j < totalNum; j++ {
					res := <-c.clients[index].Apply
					if res.Error != nil {
						panic(fmt.Sprintf("res error  %+v", res))
					}
					assertNotDuplicate(res.Seq, res.Cmd)
					log.Trace(fmt.Sprintf("%s:[Apply]:%+v\n", c.clients[index], res))
					l.Lock()
					cid, seq := splitres(res.Cmd)
					t.Logf("add-res %s", res.Cmd)
					m[cid] = append(m[cid], seq)
					l.Unlock()
					wg.Done()
				}

			}(i)
		}

	}()

	c.ni.Enable(0, false)
	go func() {
		for i := 0; i < cn; i++ {
			go func(index int) {
				for j := 0; j < totalNum; j++ {
					c.clients[index].Start(strconv.Itoa(j) + "-" + strconv.Itoa(index))
					time.Sleep(2 * time.Millisecond)
				}
			}(i)
		}
	}()

	wg.Wait()
	time.Sleep(5 * time.Second)
	t.Logf("Apply....")
	for k, v := range m {
		t.Logf("cid %d------", k)
		sort.Sort(sort.IntSlice(v))
		t.Logf("%d Apply:%v", k, v)
	}
}

//测试连续两个主节点切换问题
func TestDisconnectLeader2(t *testing.T) {
	l := sync.Mutex{}
	m := make(map[int][]int)
	c := config{useMockServer: false}
	cn := 1
	c.init(cn, 7, 0)
	totalNum := checkpointDiv * 3
	wg := sync.WaitGroup{}

	wg.Add(cn)
	go func() {
		for i := 0; i < cn; i++ {
			go func(index int) {
				for j := 0; j < totalNum; j++ {
					res := <-c.clients[index].Apply
					if res.Error != nil {
						panic(fmt.Sprintf("res error  %+v", res))
					}
					log.Trace(fmt.Sprintf("%s:[Apply]:%+v\n", c.clients[index], res))
					l.Lock()
					cid, seq := splitres(res.Cmd)
					t.Logf("add-res %s", res.Cmd)
					m[cid] = append(m[cid], seq)
					l.Unlock()
				}
				wg.Done()
			}(i)
		}

	}()

	wg.Add(cn * 2)
	go func() {

		for i := 0; i < cn; i++ {
			go func(index int) {
				for j := 0; j < totalNum/2; j++ {
					c.clients[index].Start(strconv.Itoa(j) + "-" + strconv.Itoa(index))
					time.Sleep(2 * time.Millisecond)
				}
				wg.Done()
			}(i)
		}
		c.ni.Enable(0, false)
		c.ni.Enable(1, false)
		//c.ni.Enable(5+c.cn, false)
		//c.ni.Enable(7+c.cn, false)

		log.Trace("Disconnect replica 0")

		for i := 0; i < cn; i++ {
			go func(index int) {
				for j := totalNum / 2; j < totalNum; j++ {
					c.clients[index].Start(strconv.Itoa(j) + "-" + strconv.Itoa(index))
					time.Sleep(2 * time.Millisecond)
				}
				wg.Done()
			}(i)
		}
	}()

	wg.Wait()
	time.Sleep(2 * time.Second)
	t.Logf("Apply....")
	for k, v := range m {
		t.Logf("cid %d------", k)
		sort.Sort(sort.IntSlice(v))
		t.Logf("%d Apply:%v", k, v)
	}
}

func testDisconnect(t *testing.T, rid int, reqnum int, pause time.Duration) {
	c := config{useMockServer: false}
	c.init(1, 4, 0)

	c.ni.Enable(rid, false)

	wg := sync.WaitGroup{}

	for i := 0; i < c.cn; i++ {
		wg.Add(1)
		go func(cid int) {
			for j := 0; j < reqnum; j++ {
				c.clients[cid].Start(strconv.Itoa(j) + "-" + strconv.Itoa(cid))
				time.Sleep(pause)
			}
			wg.Done()
		}(i)
		wg.Add(1)
		go func(cid int) {
			for j := 0; j < reqnum; j++ {
				res := <-c.clients[cid].Apply
				if res.Error != nil {
					panic(fmt.Sprintf("res error  %+v", res))
				}
				log.Trace(fmt.Sprintf("%s:[Apply]:%+v\n", c.clients[cid], res))
			}
			wg.Done()
		}(i)
	}

	wg.Wait()
	time.Sleep(5 * time.Second)
}

/*
测试跟随节点一开始就掉线,工作一会儿上线,再一会儿其他跟随节点再下线,保证可以工作
*/
func TestDisconnectFollowerAndReConnect(t *testing.T) {
	testNodeDisconnectAndReconnect(t, 1, 2*time.Millisecond)
}

/*
测试主节点一开始就掉线,工作一会儿主节点上线,再一会儿当前主节点再下线,保证可以工作
*/
func TestDisconnectLeaderAndReConnectLeader(t *testing.T) {
	testNodeDisconnectAndReconnect(t, 0, 2*time.Millisecond)
}

func testNodeDisconnectAndReconnect(t *testing.T, rid int, pause time.Duration) {
	/*
		必须是3倍的checkpint,
		先断开rid,然后处理第一个checkpointdiv
		再连上rid,处理第二个checkpointdiv,这时候rid就会和其他节点进行同步完成
		然后再断开rid+1,这样保证还是仅有2/3的节点有效,确保仍然可以正常工作.
	*/
	reqnum := checkpointDiv*3 + 1
	/*
		请求的第一部分让rid断网,待处理完毕以后再让rid连网,这样看看rid能否正常处理业务
		注意有没有rid的参与,pbft都能正常工作.
	*/
	rq := require.New(t)
	rq.True(reqnum/3 > 0, "reqnum must not too small")
	firstPartNumber := reqnum / 3
	secondPartNumber := reqnum / 3 * 2
	firstChan := make(chan bool, 1)
	c := config{useMockServer: false}
	c.init(1, 4, 0)

	c.ni.Enable(rid, false)

	wg := sync.WaitGroup{}

	for i := 0; i < c.cn; i++ {
		wg.Add(1)
		go func(cid int) {
			for j := 0; j < reqnum; j++ {
				c.clients[cid].Start(strconv.Itoa(j) + "-" + strconv.Itoa(cid))
				if j == firstPartNumber {
					<-firstChan //等待第一部分完成
					//rid联网成功
					c.ni.Enable(rid, true)
				}
				if j == secondPartNumber {
					<-firstChan //等待第二部分完成,这时候断开rid+1
					c.ni.Enable(rid+1, false)
				}
				time.Sleep(pause)
			}
			wg.Done()
		}(i)
		wg.Add(1)
		go func(cid int) {
			for j := 0; j < reqnum; j++ {
				res := <-c.clients[cid].Apply
				if res.Error != nil {
					panic(fmt.Sprintf("res error  %+v", res))
				}
				log.Trace(fmt.Sprintf("%s:[Apply]:%+v\n", c.clients[cid], res))
				if j == firstPartNumber {
					time.Sleep(time.Second)
					log.Trace(fmt.Sprintf("first part finished"))
					firstChan <- true //第一部分完成了,通知第二部分可以启动了
				}
				if j == secondPartNumber {
					time.Sleep(time.Second)
					log.Trace(fmt.Sprintf("second part finished"))
					firstChan <- true //第一部分完成了,通知第二部分可以启动了
				}
			}
			wg.Done()
		}(i)
	}

	wg.Wait()
	time.Sleep(5 * time.Second)
}
func (c *config) init(cn int, sn int, initSeq int) {
	c.cn = cn
	c.sn = sn
	c.n = c.cn + c.sn
	c.f = 1
	c.initSeq = initSeq
	c.servers = make([]interface{}, c.sn)
	c.clients = make([]*Client, c.cn)
	c.serverAddresses = make([]chan interface{}, c.sn)
	c.clientAddresses = make([]chan interface{}, c.cn)
	c.ni = NewNetwork(c.sn)
	c.hub = NewMessageHub(c.serverAddresses, c.clientAddresses, c.ni)
	for i := 0; i < c.sn; i++ {
		c.serverAddresses[i] = make(chan interface{}, 10) //必须有缓冲区,否则会失败
	}
	for i := 0; i < c.cn; i++ {
		c.clientAddresses[i] = make(chan interface{}, 10)
	}

	nodes := make([]int, c.sn)
	for i := 0; i < c.sn; i++ {
		nodes[i] = i
	}
	for i := 0; i < c.cn; i++ {
		c.clients[i] = NewPBFTClient(i, c.clientAddresses[i], c.hub, c.f, nodes)
	}

	for i := 0; i < c.sn; i++ {
		if c.useMockServer {
			c.servers[i] = NewMockServer(i, c.f, c.initSeq, c.serverAddresses[i], c.hub, nodes)
		} else {
			c.servers[i] = NewPBFTServer(i, c.f, c.initSeq, c.serverAddresses[i], c.hub, nodes, nil)
		}
	}
}

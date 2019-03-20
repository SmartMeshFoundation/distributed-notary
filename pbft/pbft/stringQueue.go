package pbft

import "errors"

type StringQueue struct {
	Q []string
	M map[string]interface{}
}

func (sq *StringQueue) ExtractMin() (key string, c interface{}, err error) {
	if len(sq.Q) != 0 {
		key = sq.Q[0]
		c = sq.M[key]
		sq.Q = sq.Q[1:]
		delete(sq.M, key)
		return
	}
	err = errors.New("no element")
	return
}
func (sq *StringQueue) Remove(key string) bool {
	delete(sq.M, key)
	for i := 0; i < len(sq.Q); i++ {
		if sq.Q[i] == key {
			sq.Q = append(sq.Q[0:i], sq.Q[i+1:]...)
			return true
		}
	}
	return false
}
func (sq *StringQueue) GetMin() (key string, c interface{}, err error) {
	if len(sq.Q) != 0 {
		key = sq.Q[0]
		c = sq.M[key]
		return
	}
	err = errors.New("no element")
	return
}

//不重复即可
func (sq *StringQueue) Insert(key string, c interface{}) bool {
	_, ok := sq.M[key]
	if ok {
		return false
	}
	sq.Q = append(sq.Q, key)
	sq.M[key] = c
	return true
}

func (sq *StringQueue) Length() int {
	return len(sq.Q)
}

func NewStringQueue() *StringQueue {
	return &StringQueue{
		M: make(map[string]interface{}),
	}
}

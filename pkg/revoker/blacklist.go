package revoker

import (
	"sync"
	"time"
)

type BlackList struct {
	elements []Element
	mx       sync.RWMutex
}

type Element struct {
	Time       time.Time
	Conditions []Condition
}

type Condition struct {
	Key   string
	Value string
}

func NewBlackList() *BlackList {
	bl := &BlackList{
		elements: make([]Element, 0),
	}
	go bl.cleaner()

	return bl
}

func (bl *BlackList) Add(data map[string]string) {
	bl.mx.Lock()
	defer bl.mx.Unlock()

	var conditions []Condition
	for k, v := range data {
		conditions = append(conditions, Condition{Key: k, Value: v})
	}

	bl.elements = append(bl.elements, Element{
		Time:       time.Now(),
		Conditions: conditions,
	})
}

func (bl *BlackList) Test(data map[string]any) bool {
	bl.mx.RLock()
	defer bl.mx.RUnlock()

	for _, element := range bl.elements {
		iat, ok := data["iat"].(float64)
		if !ok {
			continue
		}
		iatTime := time.Unix(int64(iat), 0)

		if !element.Time.After(iatTime) {
			continue
		}

		flag := true
		for _, v := range element.Conditions {
			condition, ok := data[v.Key]
			if !ok || condition != v.Value {
				flag = false
				break
			}
		}

		if !flag {
			continue
		}

		return true
	}

	return false
}

func (bl *BlackList) cleaner() {
	ticker := time.NewTicker(time.Minute * 30)
	defer ticker.Stop()

	for range ticker.C {
		bl.clean()
	}
}

func (bl *BlackList) clean() {
	bl.mx.Lock()
	defer bl.mx.Unlock()

	for k := len(bl.elements) - 1; k >= 0; k-- {
		t := bl.elements[k].Time.Add(time.Hour * 6)
		if t.After(time.Now()) {
			bl.elements = append(bl.elements[:k], bl.elements[k+1:]...)
		}
	}
}

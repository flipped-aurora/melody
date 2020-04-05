package refresh

import "sync"

var RefreshList = NewRefreshList()

type List struct {
	Head *Refresh
	Tail *Refresh
	Size int
	lock sync.Mutex
}

type Refresh struct {
	Value *chan int
	Next  *Refresh
	Pre   *Refresh
}

func NewRefreshList() *List {
	head := &Refresh{}
	tail := &Refresh{
		Pre: head,
	}
	head.Next = tail
	return &List{
		Head: head,
		Tail: tail,
	}
}

func (this *List) Add(node *Refresh) {
	this.lock.Lock()
	node.Pre, node.Next = this.Tail.Pre, this.Tail
	node.Pre.Next, node.Next.Pre = node, node
	this.Size++
	this.lock.Unlock()
}

func (this *List) Remove(node *Refresh) {
	this.lock.Lock()
	node.Pre.Next, node.Next.Pre = node.Next, node.Pre
	this.Size--
	this.lock.Unlock()
}

func (this *List) Front() *Refresh {
	return this.Head.Next
}

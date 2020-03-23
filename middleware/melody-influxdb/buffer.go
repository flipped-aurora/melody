package influxdb

import "github.com/influxdata/influxdb/client/v2"

func NewBuffer(size int) *Buffer {
	return &Buffer{
		data: []client.BatchPoints{},
		size: size,
	}
}

type Buffer struct {
	data []client.BatchPoints
	size int
}

func (b *Buffer) Add(ps ...client.BatchPoints) {
	b.data = append(b.data, ps...)
	if len(b.data) > b.size {
		b.data = b.data[len(b.data)-b.size:]
	}
}

func (b *Buffer) Elements() []client.BatchPoints {
	var res []client.BatchPoints
	res, b.data = b.data, []client.BatchPoints{}
	return res
}

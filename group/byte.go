package mycache

type Byte struct { //缓存的值  byte 可以转很多
	data []byte
}

func (b Byte) Len() int {
	return len(b.data)
}

func (b Byte) ToString() string {
	return string(b.data)
}

func (b Byte) Copy() []byte {
	temp := make([]byte, b.Len())
	copy(temp, b.data)
	return temp
}

//func (b *Byte) ToJson() interface{}

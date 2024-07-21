# buf

go buf, write var int


```
buf := new(bytes.Buffer)
pbW := NewWriter(buf)
pbW.WriteVarBytes([]byte{'a','b','c'}) 
pbW.WriteVarString("abc") 
pbW.WriteVarInt(uint64(0xfffffff))

pbR := NewReader(bytes.NewBuffer(buf.Bytes()))
dataBytes, _, err := pbR.ReadVarBytes()  
dataStr, _, err := pbR.ReadVarString() 
dataInt, _, err := pbR.ReadVarInt()
```
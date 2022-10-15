package pack

import (
	"bytes"
	bigendian "common.bojiu.com/utils/bigendian"
	"encoding/binary"
	"errors"
	"fmt"
)

// Pkg 收到的客户端包
// Pkg golang 编码网络字节序为小端

type CPackage struct {
	ProtoNum  uint16
	Secret    [2]byte
	RandNum   [4]byte
	ProtoData []byte

	ContronMsg int   //文本  二进制  ping pong close  用户收到包的判断
	CErr       error // 发送过来的错误
}

func NewCPackage() CPackage {
	return CPackage{}
}

func (pkg *CPackage) PkgBgData(protoNum uint16, secret [2]byte, randNum [4]byte, protoData []byte) []byte {
	var data []byte
	if protoData != nil {
		dataLen := binary.Size(protoData)
		data = make([]byte, dataLen+8)
	} else {
		data = make([]byte, 8)
	}

	data[0] = byte(protoNum)      // int8 == byte
	data[1] = byte(protoNum >> 8) //

	// 这样的写法 业务逻辑不明显
	//data = append(data,secret[:]...)
	data[2] = secret[0]
	data[3] = secret[1]

	data[4] = randNum[0]
	data[5] = randNum[1]
	data[6] = randNum[2]
	data[7] = randNum[3]

	if protoData != nil {
		for k, _ := range protoData {
			data[k+8] = protoData[k]
		}
	}
	//copy(data[8:], protoData[:])

	return data
}

func (pkg *CPackage) UnPkgBgData(wmsg WsMessage) error {
	data := wmsg.Data // 多少byte
	r := bytes.NewReader(data)

	if r.Size() < 8 {
		return errors.New(fmt.Sprintf("客户端包长度不够 %+v", data))
	}
	t := readNByte(2, r)
	pkg.ProtoNum = bigendian.FromUint16([2]byte{t[0], t[1]})

	//protoNumStr := fmt.Sprintf("%d", pkg.ProtoNum)
	//log.ZapLog.With(zap.String("protoNumStr", protoNumStr)).Info("UnPkgBgData解包")
	//if strings.Contains(msg.CMDS, protoNumStr) == false {
	//	return errors.New("客户端包协议号错误")
	//}

	t = readNByte(2, r)
	pkg.Secret = [2]byte{t[0], t[1]}
	t = readNByte(4, r)
	pkg.RandNum = [4]byte{t[0], t[1], t[2], t[3]}
	pkg.ProtoData = readNByte(r.Len(), r)

	pkg.ContronMsg = wmsg.ContronMsg
	pkg.CErr = wmsg.Err
	return nil
}

func readNByte(n int, r *bytes.Reader) (s []byte) {
	for i := 0; i < n; i++ {
		t, _ := r.ReadByte()
		s = append(s, t)
	}
	return s
}

// PkgSmBgData 小端发送
func (pkg *CPackage) PkgSmBgData(protoNum uint16, secret [2]byte, randNum [4]byte, protoData []byte) []byte {
	var data []byte
	if protoData != nil {
		dataLen := binary.Size(protoData)
		data = make([]byte, dataLen+8)
	} else {
		data = make([]byte, 8)
	}

	data[0] = byte(protoNum >> 8) // int8 == byte
	data[1] = byte(protoNum)      //
	data[2] = secret[0]
	data[3] = secret[1]
	data[4] = randNum[0]
	data[5] = randNum[1]
	data[6] = randNum[2]
	data[7] = randNum[3]

	if protoData != nil {
		for k, _ := range protoData {
			data[k+8] = protoData[k]
		}
	}
	return data
}

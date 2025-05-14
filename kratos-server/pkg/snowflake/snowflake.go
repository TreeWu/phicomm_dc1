// Package snowflake provides a very simple Twitter snowflake generator and parser.
package snowflake

import (
	"encoding/base64"
	"encoding/binary"
	"errors"
	"fmt"
	"strconv"
	"sync"
	"time"
)

var (
	// Epoch 设置为 Twitter Snowflake 的纪元时间，即 2010 年 11 月 4 日 01:42:54 UTC 的毫秒表示
	// 您可以根据应用程序的需求自定义此值以设置不同的纪元时间。
	Epoch int64 = 1288834974657

	// NodeBits 保存用于节点（Node）的比特位数
	// 请记住，您总共有 22 位可在节点（Node）和步长（Step）之间共享。
	NodeBits uint8 = 10

	// StepBits 保存用于步长（Step）的比特位数
	// 请记住，您总共有 22 位可在节点（Node）和步长（Step）之间共享。
	StepBits uint8 = 12

	// 已弃用：以下四个变量将在将来的版本中移除。
	mu        sync.Mutex
	nodeMax   int64 = -1 ^ (-1 << NodeBits)
	nodeMask        = nodeMax << StepBits
	stepMask  int64 = -1 ^ (-1 << StepBits)
	timeShift       = NodeBits + StepBits
	nodeShift       = StepBits
)

const encodeBase32Map = "ybndrfg8ejkmcpqxot1uwisza345h769"

var decodeBase32Map [256]byte

const encodeBase58Map = "123456789abcdefghijkmnopqrstuvwxyzABCDEFGHJKLMNPQRSTUVWXYZ"

var decodeBase58Map [256]byte

// 如果提供了无效的 ID，则在 UnmarshalJSON 中返回 JSONSyntaxError。
type JSONSyntaxError struct{ original []byte }

func (j JSONSyntaxError) Error() string {
	return fmt.Sprintf("invalid  snowflake ID: %q", string(j.original))
}

// 在给定无效 []byte 时，由 ParseBase58 返回 ErrInvalidBase58。
var ErrInvalidBase58 = errors.New("invalid  base58")

// 在给定无效 []byte 时，由 ParseBase32 返回 ErrInvalidBase32。
var ErrInvalidBase32 = errors.New("invalid  base32")

// 为解码 Base58/Base32 创建映射。
// 这极大地加快了解码过程。
func init() {

	for i := 0; i < len(decodeBase58Map); i++ {
		decodeBase58Map[i] = 0xFF
	}

	for i := 0; i < len(encodeBase58Map); i++ {
		decodeBase58Map[encodeBase58Map[i]] = byte(i)
	}

	for i := 0; i < len(decodeBase32Map); i++ {
		decodeBase32Map[i] = 0xFF
	}

	for i := 0; i < len(encodeBase32Map); i++ {
		decodeBase32Map[encodeBase32Map[i]] = byte(i)
	}
}

// Node 结构体保存用于 snowflake 生成器的基本信息
// 节点
type Node struct {
	mu    sync.Mutex
	epoch time.Time
	time  int64
	node  int64
	step  int64

	nodeMax   int64
	nodeMask  int64
	stepMask  int64
	timeShift uint8
	nodeShift uint8
}

// ID 是用于 snowflake ID 的自定义类型。这样可以将方法附加到 ID 上。
type ID int64

// NewNode 返回一个新的 snowflake 节点，可用于生成 snowflake ID
func NewNode(node int64) (*Node, error) {

	if NodeBits+StepBits > 22 {
		return nil, errors.New("请记住，您总共有 22 位可在节点（Node）和步长（Step）之间共享")
	}
	// 在自定义 NodeBits 或 StepBits 设置的情况下重新计算
	// 已弃用：以下块将在将来的版本中移除。
	mu.Lock()
	nodeMax = -1 ^ (-1 << NodeBits)
	nodeMask = nodeMax << StepBits
	stepMask = -1 ^ (-1 << StepBits)
	timeShift = NodeBits + StepBits
	nodeShift = StepBits
	mu.Unlock()

	n := Node{}
	n.node = node
	n.nodeMax = -1 ^ (-1 << NodeBits)
	n.nodeMask = n.nodeMax << StepBits
	n.stepMask = -1 ^ (-1 << StepBits)
	n.timeShift = NodeBits + StepBits
	n.nodeShift = StepBits

	if n.node < 0 || n.node > n.nodeMax {
		return nil, errors.New("节点编号必须介于 0 和 " + strconv.FormatInt(n.nodeMax, 10) + " 之间")
	}

	var curTime = time.Now()
	// 添加 time.Duration 到 curTime 以确保使用单调时钟（如果可用）
	n.epoch = curTime.Add(time.Unix(Epoch/1000, (Epoch%1000)*1000000).Sub(curTime))

	return &n, nil
}

// Generate 创建并返回一个唯一的 snowflake ID
// 为了确保唯一性
// - 确保您的系统保持准确的系统时间
// - 确保您从未运行具有相同节点 ID 的多个节点
func (n *Node) Generate() ID {

	n.mu.Lock()
	defer n.mu.Unlock()

	now := time.Since(n.epoch).Milliseconds()

	if now == n.time {
		n.step = (n.step + 1) & n.stepMask

		if n.step == 0 {
			for now <= n.time {
				now = time.Since(n.epoch).Milliseconds()
			}
		}
	} else {
		n.step = 0
	}

	n.time = now

	r := ID((now)<<n.timeShift |
		(n.node << n.nodeShift) |
		(n.step),
	)

	return r
}

// Int64 返回 snowflake ID 的 int64 表示
func (f ID) Int64() int64 {
	return int64(f)
}

// ParseInt64 将 int64 转换为 snowflake ID
func ParseInt64(id int64) ID {
	return ID(id)
}

// String 返回 snowflake ID 的字符串表示
func (f ID) String() string {
	return strconv.FormatInt(int64(f), 10)
}

// ParseString 将字符串转换为 snowflake ID
func ParseString(id string) (ID, error) {
	i, err := strconv.ParseInt(id, 10, 64)
	return ID(i), err
}

// Base2 返回 snowflake ID 的二进制字符串表示
func (f ID) Base2() string {
	return strconv.FormatInt(int64(f), 2)
}

// ParseBase2 将二进制字符串转换为 snowflake ID
func ParseBase2(id string) (ID, error) {
	i, err := strconv.ParseInt(id, 2, 64)
	return ID(i), err
}

// Base32 使用 z-base-32 字符集，但编码和解码类似于 base58，允许创建更小的结果字符串。
// 注意：有许多不同的 base32 实现，因此在进行任何交互操作时要小心。
func (f ID) Base32() string {

	if f < 32 {
		return string(encodeBase32Map[f])
	}

	b := make([]byte, 0, 12)
	for f >= 32 {
		b = append(b, encodeBase32Map[f%32])
		f /= 32
	}
	b = append(b, encodeBase32Map[f])

	for x, y := 0, len(b)-1; x < y; x, y = x+1, y-1 {
		b[x], b[y] = b[y], b[x]
	}

	return string(b)
}

// ParseBase32 将 base32 []byte 解析为 snowflake ID
// 注意：有许多不同的 base32 实现，因此在进行任何交互操作时要小心。
func ParseBase32(b []byte) (ID, error) {

	var id int64

	for i := range b {
		if decodeBase32Map[b[i]] == 0xFF {
			return -1, ErrInvalidBase32
		}
		id = id*32 + int64(decodeBase32Map[b[i]])
	}

	return ID(id), nil
}

// Base36 返回 snowflake ID 的 base36 字符串表示
func (f ID) Base36() string {
	return strconv.FormatInt(int64(f), 36)
}

// ParseBase36 将 Base36 字符串转换为 snowflake ID
func ParseBase36(id string) (ID, error) {
	i, err := strconv.ParseInt(id, 36, 64)
	return ID(i), err
}

// Base58 返回 snowflake ID 的 base58 字符串表示
func (f ID) Base58() string {

	if f < 58 {
		return string(encodeBase58Map[f])
	}

	b := make([]byte, 0, 11)
	for f >= 58 {
		b = append(b, encodeBase58Map[f%58])
		f /= 58
	}
	b = append(b, encodeBase58Map[f])

	for x, y := 0, len(b)-1; x < y; x, y = x+1, y-1 {
		b[x], b[y] = b[y], b[x]
	}

	return string(b)
}

// ParseBase58 将 base58 []byte 解析为 snowflake ID
func ParseBase58(b []byte) (ID, error) {

	var id int64

	for i := range b {
		if decodeBase58Map[b[i]] == 0xFF {
			return -1, ErrInvalidBase58
		}
		id = id*58 + int64(decodeBase58Map[b[i]])
	}

	return ID(id), nil
}

// Base64 返回 snowflake ID 的 base64 字符串表示
func (f ID) Base64() string {
	return base64.StdEncoding.EncodeToString(f.Bytes())
}

// ParseBase64 将 base64 字符串转换为 snowflake ID
func ParseBase64(id string) (ID, error) {
	b, err := base64.StdEncoding.DecodeString(id)
	if err != nil {
		return -1, err
	}
	return ParseBytes(b)

}

// Bytes 返回 snowflake ID 的字节切片
func (f ID) Bytes() []byte {
	return []byte(f.String())
}

// ParseBytes 将字节切片转换为 snowflake ID
func ParseBytes(id []byte) (ID, error) {
	i, err := strconv.ParseInt(string(id), 10, 64)
	return ID(i), err
}

// IntBytes 返回 snowflake ID 的字节数组，以大端整数表示。
func (f ID) IntBytes() [8]byte {
	var b [8]byte
	binary.BigEndian.PutUint64(b[:], uint64(f))
	return b
}

// ParseIntBytes 将字节数组以大端整数表示转换为 snowflake ID
func ParseIntBytes(id [8]byte) ID {
	return ID(int64(binary.BigEndian.Uint64(id[:])))
}

// Time 返回 snowflake ID 时间的 int64 UNIX 时间戳（以毫秒为单位）
// 已弃用：以下函数将在将来的版本中移除。
func (f ID) Time() int64 {
	return (int64(f) >> timeShift) + Epoch
}

// Node 返回 snowflake ID 节点号的 int64 表示
// 已弃用：以下函数将在将来的版本中移除。
func (f ID) Node() int64 {
	return int64(f) & nodeMask >> nodeShift
}

// Step 返回 snowflake 步骤（或序列）号的 int64 表示
// 已弃用：以下函数将在将来的版本中移除。
func (f ID) Step() int64 {
	return int64(f) & stepMask
}

// MarshalJSON 返回 snowflake ID 的 JSON 字节数组字符串表示。
func (f ID) MarshalJSON() ([]byte, error) {
	buff := make([]byte, 0, 22)
	buff = append(buff, '"')
	buff = strconv.AppendInt(buff, int64(f), 10)
	buff = append(buff, '"')
	return buff, nil
}

// UnmarshalJSON 将 snowflake ID 的 JSON 字节数组转换为 ID 类型。
func (f *ID) UnmarshalJSON(b []byte) error {
	if len(b) < 3 || b[0] != '"' || b[len(b)-1] != '"' {
		return JSONSyntaxError{b}
	}

	i, err := strconv.ParseInt(string(b[1:len(b)-1]), 10, 64)
	if err != nil {
		return err
	}

	*f = ID(i)
	return nil
}

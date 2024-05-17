package logger

import (
	"time"

	"github.com/gopi-frame/contract/logger"
	"go.uber.org/zap/buffer"
	"go.uber.org/zap/zapcore"
)

func newEncoder(formatter logger.Formatter) zapcore.Encoder {
	return &encoder{
		pool:             buffer.NewPool(),
		formatter:        formatter,
		MapObjectEncoder: zapcore.NewMapObjectEncoder(),
	}
}

type encoder struct {
	pool      buffer.Pool
	formatter logger.Formatter
	*zapcore.MapObjectEncoder
}

func (enc *encoder) Clone() zapcore.Encoder {
	return &encoder{
		pool:             enc.pool,
		formatter:        enc.formatter,
		MapObjectEncoder: zapcore.NewMapObjectEncoder(),
	}
}

func (enc *encoder) EncodeEntry(e zapcore.Entry, fields []zapcore.Field) (*buffer.Buffer, error) {
	enc = enc.Clone().(*encoder)
	for _, field := range fields {
		field.AddTo(enc)
	}
	ent := &entry{Entry: e, fields: enc.Fields}
	buf := enc.pool.Get()
	buf.AppendString(enc.formatter.Format(ent))
	return buf, nil
}

// func (enc *encoder) AddArray(key string, marshaler zapcore.ArrayMarshaler) error {
// 	s := new(sliceEncoder)
// 	err := marshaler.MarshalLogArray(s)
// 	enc.fields[key] = s.elements
// 	return err
// }

// func (enc *encoder) AddObject(key string, marshaler zapcore.ObjectMarshaler) error {
// 	m := zapcore.NewMapObjectEncoder()
// 	err := marshaler.MarshalLogObject(m)
// 	enc.fields[key] = m.Fields
// 	return err
// }

// func (enc *encoder) AddBinary(key string, value []byte)          { enc.fields[key] = value }
// func (enc *encoder) AddByteString(key string, value []byte)      { enc.fields[key] = string(value) }
// func (enc *encoder) AddBool(key string, value bool)              { enc.fields[key] = value }
// func (enc *encoder) AddComplex128(key string, value complex128)  { enc.fields[key] = value }
// func (enc *encoder) AddComplex64(key string, value complex64)    { enc.fields[key] = value }
// func (enc *encoder) AddDuration(key string, value time.Duration) { enc.fields[key] = value }
// func (enc *encoder) AddFloat64(key string, value float64)        { enc.fields[key] = value }
// func (enc *encoder) AddFloat32(key string, value float32)        { enc.fields[key] = value }
// func (enc *encoder) AddInt(key string, value int)                { enc.fields[key] = value }
// func (enc *encoder) AddInt64(key string, value int64)            { enc.fields[key] = value }
// func (enc *encoder) AddInt32(key string, value int32)            { enc.fields[key] = value }
// func (enc *encoder) AddInt16(key string, value int16)            { enc.fields[key] = value }
// func (enc *encoder) AddInt8(key string, value int8)              { enc.fields[key] = value }
// func (enc *encoder) AddString(key, value string)                 { enc.fields[key] = value }
// func (enc *encoder) AddTime(key string, value time.Time)         { enc.fields[key] = value }
// func (enc *encoder) AddUint(key string, value uint)              { enc.fields[key] = value }
// func (enc *encoder) AddUint64(key string, value uint64)          { enc.fields[key] = value }
// func (enc *encoder) AddUint32(key string, value uint32)          { enc.fields[key] = value }
// func (enc *encoder) AddUint16(key string, value uint16)          { enc.fields[key] = value }
// func (enc *encoder) AddUint8(key string, value uint8)            { enc.fields[key] = value }
// func (enc *encoder) AddUintptr(key string, value uintptr)        { enc.fields[key] = value }
// func (enc *encoder) AddReflected(key string, value interface{}) error {
// 	enc.fields[key] = value
// 	return nil
// }
// func (enc *encoder) OpenNamespace(key string) {
// 	m := make(map[string]any)
// 	enc.fields[key] = m
// 	enc.final = &enc.fields
// 	enc.fields = m
// }

type sliceEncoder struct{ elements []any }

func (s *sliceEncoder) AppendBool(value bool)              { s.elements = append(s.elements, value) }
func (s *sliceEncoder) AppendByteString(value []byte)      { s.elements = append(s.elements, string(value)) }
func (s *sliceEncoder) AppendComplex128(value complex128)  { s.elements = append(s.elements, value) }
func (s *sliceEncoder) AppendComplex64(value complex64)    { s.elements = append(s.elements, value) }
func (s *sliceEncoder) AppendFloat64(value float64)        { s.elements = append(s.elements, value) }
func (s *sliceEncoder) AppendFloat32(value float32)        { s.elements = append(s.elements, value) }
func (s *sliceEncoder) AppendInt(value int)                { s.elements = append(s.elements, value) }
func (s *sliceEncoder) AppendInt64(value int64)            { s.elements = append(s.elements, value) }
func (s *sliceEncoder) AppendInt32(value int32)            { s.elements = append(s.elements, value) }
func (s *sliceEncoder) AppendInt16(value int16)            { s.elements = append(s.elements, value) }
func (s *sliceEncoder) AppendInt8(value int8)              { s.elements = append(s.elements, value) }
func (s *sliceEncoder) AppendString(value string)          { s.elements = append(s.elements, value) }
func (s *sliceEncoder) AppendUint(value uint)              { s.elements = append(s.elements, value) }
func (s *sliceEncoder) AppendUint64(value uint64)          { s.elements = append(s.elements, value) }
func (s *sliceEncoder) AppendUint32(value uint32)          { s.elements = append(s.elements, value) }
func (s *sliceEncoder) AppendUint16(value uint16)          { s.elements = append(s.elements, value) }
func (s *sliceEncoder) AppendUint8(value uint8)            { s.elements = append(s.elements, value) }
func (s *sliceEncoder) AppendUintptr(value uintptr)        { s.elements = append(s.elements, value) }
func (s *sliceEncoder) AppendDuration(value time.Duration) { s.elements = append(s.elements, value) }
func (s *sliceEncoder) AppendTime(value time.Time)         { s.elements = append(s.elements, value) }
func (s *sliceEncoder) AppendArray(value zapcore.ArrayMarshaler) error {
	enc := new(sliceEncoder)
	err := value.MarshalLogArray(enc)
	s.elements = append(s.elements, enc.elements)
	return err
}
func (s *sliceEncoder) AppendObject(value zapcore.ObjectMarshaler) error {
	enc := zapcore.NewMapObjectEncoder()
	err := value.MarshalLogObject(enc)
	s.elements = append(s.elements, enc.Fields)
	return err
}
func (s *sliceEncoder) AppendReflected(value interface{}) error {
	s.elements = append(s.elements, value)
	return nil
}

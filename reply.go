package redis

import (
	"github.com/redis/go-redis/v9"
)

type Reply struct {
	cmd *redis.Cmd

	Err error
}

func (r *Reply) Nil() bool {
	if r.Err == redis.Nil {
		return true
	}

	return false
}

func (r *Reply) Bool() (bool, error) {
	return r.cmd.Bool()
}

func (r *Reply) BoolSlice() ([]bool, error) {
	return r.cmd.BoolSlice()
}

func (r *Reply) Float32() (float32, error) {
	return r.cmd.Float32()
}

func (r *Reply) Float32Slice() ([]float32, error) {
	return r.cmd.Float32Slice()
}

func (r *Reply) Float64() (float64, error) {
	return r.cmd.Float64()
}

func (r *Reply) Float64Slice() ([]float64, error) {
	return r.cmd.Float64Slice()
}

func (r *Reply) Int64() (int64, error) {
	return r.cmd.Int64()
}

func (r *Reply) Int64Slice() ([]int64, error) {
	return r.cmd.Int64Slice()
}

func (r *Reply) String() (string, error) {
	return r.cmd.Text()
}

func (r *Reply) StringSlice() ([]string, error) {
	return r.cmd.StringSlice()
}

func (r *Reply) Uint64() (uint64, error) {
	return r.cmd.Uint64()
}

func (r *Reply) Uint64Slice() ([]uint64, error) {
	return r.cmd.Uint64Slice()
}

func (r *Reply) Int() (int, error) {
	return r.cmd.Int()
}

func (r *Reply) Slice() ([]interface{}, error) {
	return r.cmd.Slice()
}

func (r *Reply) Value() interface{} {
	return r.cmd.Val()
}

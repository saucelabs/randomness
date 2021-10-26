// Copyright 2021 The randomness Authors. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package randomness

import (
	"crypto/rand"
	"log"
	"math"
	"math/big"

	"github.com/saucelabs/customerror"
	"github.com/saucelabs/sypl"
	"github.com/saucelabs/sypl/fields"
	"github.com/saucelabs/sypl/level"
	"github.com/saucelabs/sypl/options"
)

var l = sypl.NewDefault("randomgenerator", level.Info)

var (
	ErrFailedToGenerateRandomness      = customerror.NewFailedToError("to generate randomness", "", nil)
	ErrFailedToGenerateRangeSaturated  = customerror.NewFailedToError("to generate, range saturated", "", nil)
	ErrFailedToGenerateReachedMaxRetry = customerror.NewFailedToError("to generate, reached max retry", "", nil)
	ErrInvalidMax                      = customerror.NewInvalidError("params. `max` is less than 0", "", nil)
	ErrInvalidMin                      = customerror.NewInvalidError("params. `min` is less than 1", "", nil)
	ErrInvalidMinBiggerThanMax         = customerror.NewInvalidError("param. Min can't be bigger than max", "", nil)
)

//////
// Helpers.
//////

type Randomness struct {
	CollisionFree bool

	Max int
	Min int

	maxRetry *int
	memory   []int64
}

// Generate returns a random number.
func (r *Randomness) Generate() (int64, error) {
	// calculate the max we will be using
	bg := big.NewInt(int64(r.Max - r.Min + 1))

	n, err := rand.Int(rand.Reader, bg)
	if err != nil {
		return 0, customerror.Wrap(ErrFailedToGenerateRandomness, err)
	}

	port := n.Int64() + int64(r.Min)

	if r.memory == nil {
		return port, nil
	}

	if len(r.memory) > (r.Max - r.Min) {
		return 0, ErrFailedToGenerateRangeSaturated
	}

	for _, p := range r.memory {
		if p == port {
			if r.maxRetry != nil {
				if *r.maxRetry == 0 {
					return 0, ErrFailedToGenerateReachedMaxRetry
				}

				l.PrintlnWithOptions(&options.Options{
					Fields: fields.Fields{"retry": *r.maxRetry},
				}, level.Debug, "Retrying...")
				*r.maxRetry--
			}

			return r.Generate()
		}
	}

	r.memory = append(r.memory, port)

	return port, nil
}

// MustGenerate is like `RandomPortGenerator`, but will panic in case of any
// error.
func (r *Randomness) MustGenerate() int64 {
	n, err := r.Generate()
	if err != nil {
		log.Panicln(err)
	}

	return n
}

// GenerateMany returns an slice of `n` numbers.
func (r *Randomness) GenerateMany(n int) ([]int64, error) {
	numbers := []int64{}

	for i := 0; i < n; i++ {
		n, err := r.Generate()
		if err != nil {
			return nil, err
		}

		numbers = append(numbers, n)
	}

	return numbers, nil
}

// MustGenerateMany is like `GenerateMany`, but will panic in case of any error.
func (r *Randomness) MustGenerateMany(n int) []int64 {
	numbers, err := r.GenerateMany(n)
	if err != nil {
		log.Panicln(err)
	}

	return numbers
}

//////
// Factory
//////

// New is the Randomness factory.
func New(min, max, maxRetry int, collisionFree bool) (*Randomness, error) {
	if min < 1 {
		return nil, ErrInvalidMin
	}

	if max < 0 {
		return nil, ErrInvalidMax
	}

	if max == 0 {
		// TODO: Switch to `math.MaxInt`, only available in Go 1.17+.
		max = math.MaxInt32
	}

	if max < min {
		return nil, ErrInvalidMinBiggerThanMax
	}

	var mR *int

	if maxRetry > 0 {
		mR = &maxRetry
	}

	var memory []int64

	if collisionFree {
		memory = []int64{}
	}

	return &Randomness{
		CollisionFree: collisionFree,
		Max:           max,
		Min:           min,

		maxRetry: mR,
		memory:   memory,
	}, nil
}

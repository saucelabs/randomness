// Copyright 2021 The randomness Authors. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package randomness

import (
	"crypto/rand"
	"fmt"
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
	ErrFailedToGenerateRangeSaturated  = customerror.NewFailedToError("to generate, range saturated", "", nil)
	ErrFailedToGenerateReachedMaxRetry = customerror.NewFailedToError("to generate, reached max retry", "", nil)
	ErrInvalidMinBiggerThanMax         = customerror.NewInvalidError("param. Min can't be bigger than max", "", nil)
	ErrInvalidMinOrMaxLessThanZero     = customerror.NewInvalidError("params. `min`/`max` need to be bigger than zero", "", nil)
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

func (r *Randomness) Generate() (int64, error) {
	nBigInt := (r.Max - r.Min + 1)

	bigInt, err := rand.Int(rand.Reader, big.NewInt(int64(nBigInt)))
	if err != nil {
		log.Panicln(fmt.Errorf("failed to generated random int. %w", err))
	}

	port := bigInt.Int64() + 1

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

func New(min, max, maxRetry int, collisionFree bool) (*Randomness, error) {
	if min < 0 || max < 0 {
		return nil, ErrInvalidMinOrMaxLessThanZero
	}

	if max < min {
		return nil, ErrInvalidMinBiggerThanMax
	}

	if max == 0 {
		max = math.MaxInt
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

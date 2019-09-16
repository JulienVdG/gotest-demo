// Copyright 2019 Splitted-Desktop Systems. All rights reserved
// Copyright 2019 Julien Viard de Galbert
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"regexp"
	"testing"
	"time"

	"github.com/JulienVdG/tastevin/pkg/testsuite"
	exp "github.com/google/goexpect"
)

var (
	args = []string{"go", "run", "."}
)

func TestDemo(t *testing.T) {
	opts, warn := testsuite.ExpectOptions("")
	if warn != nil {
		t.Log(warn)
	}

	e, errchan, err := exp.SpawnWithArgs(args, 1*time.Second, opts...)
	if err != nil {
		t.Fatalf("error spawning demo: %v", err)
	}

	// match output
	out, _, err := e.Expect(regexp.MustCompile("type something"), 20*time.Second)
	if err != nil {
		t.Fatalf("error opening demo: %v (got %v)", err, out)
	}

	t.Run("LongTest", func(t *testing.T) {
		for i := 1; i < 5; i++ {
			msg := fmt.Sprintf("Step%d", i)
			testBatcher := []exp.Batcher{
				&exp.BExpT{R: "-> ", T: 1},
				&exp.BSnd{S: msg + "\n"},
				&testsuite.BExpTLog{
					L: "Matched " + msg,
					R: msg,
					T: 3,
				}}
			res, err := e.ExpectBatch(testBatcher, 0)
			if err != nil {
				t.Errorf("%s: %v", msg, testsuite.DescribeBatcherErr(testBatcher, res, err))
			}
		}
	})

	t.Run("FailedTest", func(t *testing.T) {
		testBatcher := []exp.Batcher{
			&exp.BExpT{R: "-> ", T: 1},
			&exp.BSnd{S: "Different\n"},
			&testsuite.BExpTLog{
				L: "Matched same",
				R: "same",
				T: 3,
			}}
		res, err := e.ExpectBatch(testBatcher, 0)
		if err != nil {
			fmt.Printf("Expected failure\n")
			t.Errorf("Expected failure: %v", testsuite.DescribeBatcherErr(testBatcher, res, err))
			e.Send("\n")
			time.Sleep(3 * time.Second)
		}
	})

	t.Run("MoreTest", func(t *testing.T) {
		for i := 5; i < 7; i++ {
			msg := fmt.Sprintf("Step%d", i)
			testBatcher := []exp.Batcher{
				&exp.BExpT{R: "-> ", T: 1},
				&exp.BSnd{S: msg + "\n"},
				&testsuite.BExpTLog{
					L: "Matched " + msg,
					R: msg,
					T: 3,
				}}
			res, err := e.ExpectBatch(testBatcher, 0)
			if err != nil {
				t.Errorf("%s: %v", msg, testsuite.DescribeBatcherErr(testBatcher, res, err))
			}
		}
	})

	quitBatcher := []exp.Batcher{
		&exp.BExpT{R: "-> ", T: 1},
		&exp.BSnd{S: "bye\n"},
		&testsuite.BExpTLog{
			L: "Done",
			R: "Done.",
			T: 1,
		}}

	res, err := e.ExpectBatch(quitBatcher, 0)
	if err != nil {
		t.Errorf("Quit: %v", testsuite.DescribeBatcherErr(quitBatcher, res, err))
	}

	err = e.Close()
	if err != nil {
		t.Errorf("error closing Spawn: %v", err)
	}

	// make sure the expect session is done and screenlog are closed
	<-errchan
}

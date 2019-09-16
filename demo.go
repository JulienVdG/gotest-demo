// Copyright 2019 Splitted-Desktop Systems. All rights reserved
// Copyright 2019 Julien Viard de Galbert
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("type something")
	for {
		fmt.Print("\r-> ")
		text, _ := reader.ReadString('\n')
		fmt.Println("\r")
		// convert CRLF to LF
		text = strings.Replace(text, "\n", "", -1)

		if strings.Compare("bye", text) == 0 {
			break
		}
		time.Sleep(2 * time.Second)

		fmt.Printf("Got: %s\n", text)
	}
	fmt.Println("Done.")
}

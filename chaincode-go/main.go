/*
SPDX-License-Identifier: Apache-2.0
*/

package main

import (
	"bufio"
	"fmt"
	"os"
	"sync"
	"time"
)

var buffer map[string][]string

var buffer_read map[string][]string

var mutex_map map[string]*mutex

type mutex struct {
	mutex1 sync.Mutex
	mutex2 sync.Mutex
}

var devices []string

func UpdateBuffer(idDevice string, c chan string) error {

	noErr := true
	for noErr {

		//open file of encrypted data of the devices
		f, err := os.Open("Device[" + idDevice + "]_ript_1.csv")

		if err != nil {
			fmt.Println("eeeeee")
			continue
		}

		// remember to close the file at the end of the program
		defer f.Close()

		//parse data and insert in the leadger
		scanner := bufio.NewScanner(f)
		for scanner.Scan() {
			time.Sleep(100 * time.Millisecond)

			c <- scanner.Text()

		}

		close(c)

		noErr = false

	}

	return nil
}

func funz(c chan string) {
	var n int
	n = 0
	for i := range c {
		fmt.Println(i)
		n++
		if n == 4 {
			break
		}
	}

	fmt.Println("separate")
	n = 0
	for i := range c {

		fmt.Println(i)
		n++

		if n == 4 {
			break
		}
	}

}

func main() {

	var c chan string
	m := make(map[string]chan string)
	c = make(chan string, 4)
	m["10:89"] = c
	if m["1"] == nil {
		fmt.Println("yes")

	}
	go UpdateBuffer("10:89", m["10:89"])
	go funz(m["10:89"])

	//fmt.Println(EncryptDecrypt("\"init\",\"resp\",\"dist\",\"conf\"", "P"))
	//a, _ := int64(strconv.Atoi("10974441875307669990"))
	//fmt.Println(a)
	//fmt.Println(strconv.FormatUint(getHash("\"init\",\"resp\",\"dist\",\"conf\"\n"), 10))

}

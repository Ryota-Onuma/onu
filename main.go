package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/Ryota-Onuma/onu/core"
)

func main() {
	args := os.Args
	if len(args) == 2 {
		filePath := args[1]
		execWithFile(filePath)
	} else {
		callREPL()
	}
}

func execWithFile(filePath string) {
	env := core.Env{}
	f, err := os.Open(filePath)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	// fileの中身を読み込んで、stringに変換
	scanner := bufio.NewScanner(f)
	var inputs []string
	for scanner.Scan() {
		inputs = append(inputs, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	for _, input := range inputs {
		core.Execute(input, env)
	}
}

func callREPL() {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	// 対話型のRPELループを開始
	reader := bufio.NewReader(os.Stdin)
	done := make(chan struct{})
	env := core.Env{}
	go func() {
		for {
			fmt.Print(">> ")
			input, err := reader.ReadString('\n')
			if err != nil {
				done <- struct{}{}
				return
			}

			// 改行文字を削除
			input = strings.TrimSpace(input)

			// 入力が "exit" だったらループを終了
			if input == "exit" {
				fmt.Println("RPELを終了します。")
				done <- struct{}{}
				return
			}

			// 入力を評価し、結果を表示
			core.Execute(input, env)
		}
	}()

	select {
	case <-sigChan:
		os.Exit(1)
	case <-done:
	}
}

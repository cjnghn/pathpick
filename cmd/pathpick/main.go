package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/cjnghn/pathpick/internal/tree"
	"github.com/cjnghn/pathpick/internal/ui"
)

func main() {
	// CLI 플래그 정의
	path := flag.String("path", ".", "시작 디렉토리 경로")
	pattern := flag.String("pattern", "", "파일 패턴 (예: *.go)")
	flag.Parse()

	// 파일시스템 워커 초기화
	walker := &tree.Walker{
		Pattern: *pattern,
	}

	// UI 초기화 및 시작
	display := ui.NewDisplay(walker)
	if err := display.Start(*path); err != nil {
		fmt.Fprintf(os.Stderr, "오류: %v\n", err)
		os.Exit(1)
	}
}

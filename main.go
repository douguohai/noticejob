package main

import (
	"fmt"
	"github.com/go-git/go-git/v5"
	"strings"
)

func main() {
	r, _ := git.PlainOpen("./")
	rHead, _ := r.Head()
	// 将获得的
	rHeadStr := rHead.String()
	rHeadIdx := strings.Index(rHeadStr, " ")
	lastCommitHash := rHeadStr[:rHeadIdx]
	fmt.Println(lastCommitHash)

}

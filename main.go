package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

var startStr = `package %s

import (
	"context"
	"github.com/go-kit/kit/endpoint"
)
`
var endpointStr = `func make%sEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(*%sRequest)
		res, err := s.%s(ctx, req)
		return res, err
	}
}`

func main() {
	if len(os.Args) != 3 {
		return
	}

	pkNames := strings.Split(os.Args[1], "/")
	if len(pkNames) <= 0 {
		return
	}
	pkName := pkNames[len(pkNames)-1]

	rf, err := os.Open(os.Args[1] + "/" + os.Args[2])
	if err != nil {
		log.Println(err)
		return
	}
	defer rf.Close()

	wf, err := os.Open(os.Args[1] + "/endpoint.go")
	if err != nil {
		log.Println(err)
		return
	}
	defer wf.Close()

	w := bufio.NewWriter(wf)
	defer w.Flush()

	var begin bool
	br := bufio.NewReader(rf)
	for {
		a, _, c := br.ReadLine()
		if c == io.EOF {
			break
		}
		if !begin {
			if strings.Contains(string(a), "interface") {
				begin = true
			}
			w.WriteString(fmt.Sprintf(startStr, pkName))
		} else {
			// 结束
			if strings.Contains(string(a), "}") {
				return
			}
			ts := strings.Split(strings.TrimLeft(string(a), " "), "(")
			if len(ts) <= 0 {
				continue
			}
			w.WriteString(fmt.Sprintf(endpointStr, ts[0]))
		}

	}
}

package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
)

var startStr = `package %s

import (
	"context"
	"net/http"

	kit "github.com/go-kit/kit/endpoint"

	"github.com/viile/server/components/endpoint"
	"github.com/viile/server/components/errors"
)
`
var endpointStr = `func make%sEndpoint(s Service) kit.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req,ok := request.(%s)
		if !ok || req == nil{
			return nil,errors.ErrParser
		}
		res, err := s.%s(ctx, req)
		return res, err
	}
}
`

var decodeStr = `func decode%sRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	var err error
	req := &%s{}
	if err = endpoint.Decodefunc(ctx, r, req); err != nil {
		return nil, errors.ErrInput.With(err)
	}

	if err = req.Parser(ctx); err != nil {
		return nil, errors.ErrInput.With(err)
	}

	return req, nil
}
`

func getEndpoint(name, req string) string {
	name = strings.Trim(name, "\t")
	name = strings.Trim(name, " ")
	req = strings.Trim(req, "\t")
	req = strings.Trim(req, " ")
	rreq := strings.Trim(req, "*")
	return fmt.Sprintf(endpointStr, name, req, name) + fmt.Sprintf(decodeStr, name, rreq)
}

func main() {
	fmt.Println(os.Args)
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
		fmt.Println(err)
		return
	}
	defer rf.Close()

	wf, err := os.Create(os.Args[1] + "/endpoint.go")
	if err != nil {
		fmt.Println(err)
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
				w.WriteString(fmt.Sprintf(startStr, pkName))
			}
		} else {
			// 结束
			if strings.Contains(string(a), "}") {
				return
			}
			ts := strings.Split(string(a), "(")
			if len(ts) <= 0 {
				continue
			}
			t1 := strings.Split(string(a), ")")
			if len(t1) <= 0 {
				continue
			}
			t2 := strings.Split(t1[0], " ")
			if len(t2) <= 0 {
				continue
			}
			fmt.Println(string(a))
			w.WriteString(getEndpoint(ts[0], t2[len(t2)-1]))
		}

	}
}

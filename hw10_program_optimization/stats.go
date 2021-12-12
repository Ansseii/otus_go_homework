package hw10programoptimization

import (
	"bufio"
	"io"
	"strings"

	"github.com/valyala/fastjson"
)

type DomainStat map[string]int

func GetDomainStat(r io.Reader, domain string) (DomainStat, error) {
	result := make(DomainStat)

	s := bufio.NewScanner(r)
	for s.Scan() {
		email := fastjson.GetString(s.Bytes(), "Email")

		if strings.HasSuffix(email, "."+domain) {
			subdomain := strings.SplitN(email, "@", 2)[1]
			result[strings.ToLower(subdomain)]++
		}
	}
	return result, nil
}

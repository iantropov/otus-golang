package hw10programoptimization

import (
	"bufio"
	"fmt"
	"io"
	"regexp"
	"strings"

	"github.com/valyala/fastjson"
)

type User struct {
	ID       int
	Name     string
	Username string
	Email    string
	Phone    string
	Password string
	Address  string
}

type DomainStat map[string]int

func GetDomainStat(r io.Reader, domain string) (DomainStat, error) {
	domainStat, err := getDomainStat(r, domain)
	if err != nil {
		return nil, fmt.Errorf("get users error: %w", err)
	}
	return domainStat, nil
}

func getDomainStat(r io.Reader, domain string) (DomainStat, error) {
	domainStat := make(DomainStat)
	scanner := bufio.NewScanner(r)
	p := fastjson.Parser{}
	domainRx, err := regexp.Compile("\\." + domain)
	if err != nil {
		return nil, err
	}

	for scanner.Scan() {
		parsedJson, err := p.Parse(scanner.Text())
		if err != nil {
			return nil, err
		}
		email := string(parsedJson.GetStringBytes("Email"))
		if domainRx.MatchString(email) {
			domainPart := strings.SplitN(email, "@", 2)[1]
			domainStat[strings.ToLower(domainPart)]++
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return domainStat, nil
}

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
	u, err := getUsers(r)
	if err != nil {
		return nil, fmt.Errorf("get users error: %w", err)
	}
	return countDomains(u, domain)
}

type userEmails []string

func getUsers(r io.Reader) (result userEmails, err error) {
	scanner := bufio.NewScanner(r)
	var value *fastjson.Value
	p := fastjson.Parser{}
	for scanner.Scan() {
		value, err = p.Parse(scanner.Text())
		if err != nil {
			return
		}
		s := string(value.GetStringBytes("Email"))
		result = append(result, s)
	}
	if err = scanner.Err(); err != nil {
		return
	}

	return
}

func countDomains(u userEmails, domain string) (DomainStat, error) {
	result := make(DomainStat)

	domainRx, err := regexp.Compile("\\." + domain)
	if err != nil {
		return nil, err
	}
	for _, userEmail := range u {
		if domainRx.MatchString(userEmail) {
			domainPart := strings.SplitN(userEmail, "@", 2)[1]
			result[strings.ToLower(domainPart)]++
		}
	}
	return result, nil
}

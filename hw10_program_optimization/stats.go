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

type userEmails [100_000]*string

func getUsers(r io.Reader) (result userEmails, err error) {
	// jsonDecoder := json.NewDecoder(bufio.NewReaderSize(r, 409_600))

	// for i := 0; ; i++ {
	// 	var user User
	// 	if err = jsonDecoder.Decode(&user); err == io.EOF {
	// 		err = nil
	// 		break
	// 	} else if err != nil {
	// 		return
	// 	}
	// 	result[i] = user
	// }

	i := 0
	scanner := bufio.NewScanner(r)
	var value *fastjson.Value
	p := fastjson.Parser{}
	for scanner.Scan() {
		value, err = p.Parse(scanner.Text())
		if err != nil {
			return
		}
		s := string(value.GetStringBytes("Email"))
		result[i] = &s
		i++
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
		if userEmail == nil {
			return result, nil
		}

		if domainRx.MatchString(*userEmail) {
			domainPart := strings.SplitN(*userEmail, "@", 2)[1]
			result[strings.ToLower(domainPart)]++
		}
	}
	return result, nil
}

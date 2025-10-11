package hw10programoptimization

import (
	"bufio"
	"io"
	"regexp"
	"strings"

	//nolint:depguard
	"github.com/buger/jsonparser"
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
	result := make(DomainStat)
	re := regexp.MustCompile("\\." + domain + "$")

	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := scanner.Bytes()

		email, err := jsonparser.GetString(line, "Email")
		if err != nil {
			continue
		}

		email = strings.ToLower(email)
		if re.MatchString(email) {
			parts := strings.SplitN(email, "@", 2)
			if len(parts) == 2 {
				domainPart := parts[1]
				result[domainPart]++
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return result, nil
}

package hw10programoptimization

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"regexp"
	"strings"
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
	re := regexp.MustCompile(fmt.Sprintf(`(?i)@[^@]+\.%s$`, regexp.QuoteMeta(domain)))

	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		var data struct {
			Email string `json:"Email"`
		}
		if err := json.Unmarshal(scanner.Bytes(), &data); err != nil {
			continue
		}

		email := strings.ToLower(data.Email)
		if email == "" {
			continue
		}

		if re.MatchString(email) {
			parts := strings.SplitN(email, "@", 2)
			if len(parts) == 2 {
				domainPart := parts[1]
				result[domainPart]++
			}
		}
	}
	return result, nil
}

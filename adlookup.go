package main

import (
	"crypto/tls"
	"fmt"
	"log"

	"gopkg.in/ldap.v2"
)

var l *ldap.Conn

const (
	ldapServer     = ""
	ldapServerPort = 389
	ldapUser       = "" // in "PS1\\<user>" format
	ldapPassword   = ""
	baseDN         = ""
)

func connectToADServer() {
	var err error
	l, err = ldap.Dial("tcp", fmt.Sprintf("%s:%d", ldapServer, ldapServerPort))
	if err != nil {
		log.Fatal(err)
	}

	// Reconnect with TLS
	err = l.StartTLS(&tls.Config{InsecureSkipVerify: true})
	if err != nil {
		log.Fatal(err)
	}

	// First bind with a read only user
	err = l.Bind(ldapUser, ldapPassword)
	if err != nil {
		log.Fatal(err)
	}
}

func getRFIDTagsFor(username string) []string {
	searchRequest := ldap.NewSearchRequest(
		baseDN, // The base dn to search
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
		fmt.Sprintf("(&(objectClass=Person)(distinguishedName=%s))", username),
		[]string{"otherPager", "userAccountControl"},
		nil,
	)

	sr, err := l.Search(searchRequest)
	if err != nil {
		log.Fatal(err)
	}

	var pagers []string

	for _, entry := range sr.Entries {
		// Here's how we shut out inactive members: If the userAccountControl value is
		// not 512 (active), then we simply return, so the tags are not added to the
		// array and thus we're not going to let this person use the tool (assuming
		// that person had the tag we're querying)
		var userActive = entry.GetAttributeValue("userAccountControl")
		if userActive != "512" {
			log.Println("Ignoring", username, "because the member is not active")
			return pagers
		}

		pagers = make([]string, len(entry.GetAttributeValues("otherPager")))
		copy(pagers, entry.GetAttributeValues("otherPager"))
	}

	return pagers
}

func GetUsersInGroup(group string) ([]string, error) {
	base := fmt.Sprintf("cn=%s,ou=Authorization Groups,ou=Domain Groups,dc=pumpingstationone,dc=org", group)
	searchRequest := ldap.NewSearchRequest(base, ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
		"(objectClass=group)",
		[]string{"member"},
		nil,
	)

	sr, err := l.Search(searchRequest)
	if err != nil {
		return nil, err
	}

	groups := []string{}
	entries := sr.Entries[0].GetAttributeValues("member")
	for _, entry := range entries {
		groups = append(groups, entry)
	}

	return groups, nil
}

/*
func main() {
	connectToADServer()

	users, _ := GetUsersInGroup("Tormach Authorized Users")
	for _, user := range users {
		tags := getRFIDTagsFor(user)
		for _, tag := range tags {
			fmt.Println(user, " has tag ", tag)
		}
	}

	l.Close()
}
*/

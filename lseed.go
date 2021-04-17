package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/cheggaaa/pb/v3"
	"github.com/go-ldap/ldap"
)

var ouDN, bindUser, bindPassword, bindHost, photoPath, namePrefix string
var bindPort, numGroups, numMembersPerGroup int
var help bool

func main() {
	flag.StringVar(&ouDN, "ou", "ou=loadtest,dc=mm,dc=test,dc=com", "the organizational unit that will contain the seeded data")
	flag.StringVar(&bindUser, "user", "cn=admin,dc=mm,dc=test,dc=com", "the bind user")
	flag.StringVar(&bindPassword, "password", "mostest", "the bind password")
	flag.StringVar(&bindHost, "host", "0.0.0.0", "the bind host")
	flag.IntVar(&bindPort, "port", 389, "the bind port")
	flag.IntVar(&numGroups, "groups", 1, "the number of groups")
	flag.IntVar(&numMembersPerGroup, "members", 10, "the number of members per group")
	flag.StringVar(&photoPath, "photo", "", "the path to the profile photo")
	flag.StringVar(&namePrefix, "prefix", "test.", "the user uid and group cn prefix")
	flag.BoolVar(&help, "help", false, "show help")
	flag.Parse()

	if help {
		flag.PrintDefaults()
		os.Exit(0)
	}

	l, err := ldap.Dial("tcp", fmt.Sprintf("%s:%d", bindHost, bindPort))
	if err != nil {
		log.Fatal(err)
	}
	defer l.Close()
	err = l.Bind(bindUser, bindPassword)
	if err != nil {
		log.Fatal(err)
	}

	// create org. unit
	err = l.Add(&ldap.AddRequest{
		DN: ouDN,
		Attributes: []ldap.Attribute{
			{Type: "objectclass", Vals: []string{"organizationalunit"}},
		},
	})
	if err != nil && !ldap.IsErrorWithCode(err, ldap.LDAPResultEntryAlreadyExists) {
		log.Fatal(err)
	}

	// get profile photo data
	var strData string
	if len(photoPath) > 0 {
		imageData, err := ioutil.ReadFile(photoPath)
		if err != nil {
			log.Fatal(err)
			os.Exit(1)
		}
		strData = string(imageData)
	}

	userEntriesCount := numGroups * numMembersPerGroup

	bar := pb.StartNew(userEntriesCount + (numGroups * numMembersPerGroup))

	// create users
	for i := 0; i < userEntriesCount; i++ {
		username := fmt.Sprintf("%s%d", namePrefix, i)
		attributes := []ldap.Attribute{
			{Type: "objectclass", Vals: []string{"iNetOrgPerson"}},
			{Type: "cn", Vals: []string{fmt.Sprintf("Test%d", i)}},
			{Type: "sn", Vals: []string{"User"}},
			{Type: "mail", Vals: []string{fmt.Sprintf("%s@test.com", username)}},
			{Type: "userPassword", Vals: []string{"Password1"}},
		}
		if len(strData) > 0 {
			attributes = append(attributes, ldap.Attribute{Type: "jpegPhoto", Vals: []string{strData}})
		}
		err = l.Add(&ldap.AddRequest{
			DN:         fmt.Sprintf("uid=%s,%s", username, ouDN),
			Attributes: attributes,
		})
		if err != nil {
			log.Fatal(err)
		}

		bar.Increment()
	}

	for i := 0; i < numGroups; i++ {
		groupDN := fmt.Sprintf("cn=%s%d,%s", namePrefix, i, ouDN)

		var uniqueMembers []string
		for j := 0; j < numMembersPerGroup; j++ {
			username := fmt.Sprintf("%s%d", namePrefix, j+(numMembersPerGroup*i))
			uniqueMembers = append(uniqueMembers, fmt.Sprintf("uid=%s,%s", username, ouDN))
		}

		err = l.Add(&ldap.AddRequest{
			DN: groupDN,
			Attributes: []ldap.Attribute{
				{Type: "objectclass", Vals: []string{"groupOfUniqueNames"}},
				{Type: "uniqueMember", Vals: []string{uniqueMembers[0]}},
			},
		})
		if err != nil {
			log.Fatal(err)
		}

		for _, member := range uniqueMembers[1:] {
			err = l.Modify(&ldap.ModifyRequest{
				DN: groupDN,
				Changes: []ldap.Change{
					{Operation: ldap.AddAttribute, Modification: ldap.PartialAttribute{Type: "uniqueMember", Vals: []string{member}}},
				},
			})
			if err != nil {
				log.Fatal(err)
			}
			bar.Increment()
		}
		bar.Increment()
	}

	bar.Finish()

	fmt.Println("\nSuccessfully completed.")
}

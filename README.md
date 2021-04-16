# lseed

Utility to seed an LDAP instance with data.

## Help

```bash
$ go run lseed.go --help
  -groups int
    	the number of groups (default 2)
  -help
    	show help
  -host string
    	the bind host (default "0.0.0.0")
  -members int
    	the number of members per group (default 10)
  -ou string
    	the organizational unit that will contain the seeded data (default "ou=loadtest,dc=mm,dc=test,dc=com")
  -password string
    	the bind password (default "mostest")
  -photo string
    	the path to the profile photo
  -port int
    	the bind port (default 389)
  -prefix string
    	the user uid and group cn prefix (default "test.")
  -user string
    	the bind user (default "cn=admin,dc=mm,dc=test,dc=com")
```

## Examples

Seed with all of the defaults:

```bash
$ go run lseed.go
```

Seed with each user having a profile photo:

```bash
$ go run lseed.go -photo ~/Pictures/test.jpeg
```

Seed 1 group with 100,000 users and then add 30,000 more groups each with 10 users:

```bash
$ go run lseed.go -groups 1 -members 100000 -prefix "seed1."
$ go run lseed.go -groups 30000 -members 10 -prefix "seed2."
```

## LDAP count queries

Group:

```bash
$ ldapsearch -LLL -x -D "cn=admin,dc=mm,dc=test,dc=com" -w "mostest" -b "dc=mm,dc=test,dc=com" \
-h "0.0.0.0"  "(objectClass=groupOfUniqueNames)" dn | grep "dn:" | wc -l
```

Users:

```bash
$ ldapsearch -LLL -x -D "cn=admin,dc=mm,dc=test,dc=com" -w "mostest" -b "dc=mm,dc=test,dc=com" \
-h "0.0.0.0"  "(objectClass=inetOrgPerson)" dn | grep "dn:" | wc -l
```

Group members:

```bash
$ ldapsearch -LLL -x -D "cn=admin,dc=mm,dc=test,dc=com" -w "mostest" -b "dc=mm,dc=test,dc=com" \
-h "0.0.0.0"  "(objectClass=groupOfUniqueNames)" uniqueMember | grep "uniqueMember:" | wc -l
```
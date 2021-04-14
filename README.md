# lseed

Utility to seed an LDAP instance with data.

## Help

```
go run lseed.go --help
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
    	the username prefix (default "test.")
  -user string
    	the bind user (default "cn=admin,dc=mm,dc=test,dc=com")
```

## Examples

Seed with all of the defaults:
```
$ go run lseed.go
```

Seed with each user having a profile photo:
```
$ go run lseed.go -photo ~/Pictures/test.jpeg
```

Seed 1 group with 100,000 users and then add 30,000 more groups each with 10 users:
```
$ go run lseed.go -groups 1 -members 100000
$ go run lseed.go -groups 30000 -members 10 -ou "ou=loadtest2,dc=mm,dc=test,dc=com" -prefix "test2."
```

## Mattermost-specific

Login as each user:

```
$ ldapsearch -x -D "cn=admin,dc=mm,dc=test,dc=com" -w "mostest" -b "dc=mm,dc=test,dc=com" \
-h "0.0.0.0"  "(objectClass=inetOrgPerson)" | grep "uid:" | cut -d':' -f 2 | \
xargs -I {} curl 'http://localhost:8065/api/v4/users/login' -X POST \
-H 'X-Requested-With: XMLHttpRequest' \
-d '{"device_id":"","login_id":"{}","password":"Password1","token":""}'
```
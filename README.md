# sts-helper


```
function aws_sudo() {
     read -p "Helper profile:" prof
     read -p "MFA token: " token
     eval $(go run main.go --helper-profile $prof --token $token)
}
```


## Config file

~/.sts-helper.yaml

```
---
s3-stuff:
  duration: 3600
  role-arn: arn:aws:iam::1234566778899:role/foo
  mfa-arn: arn:aws:iam::1234566778899:mfa/some-arn
  session-name: foobar123
  clear-env: true
rds-readonly:
  duration: 3600
  role-arn: arn:aws:iam::1234566778899:role/foo
  mfa-arn: arn:aws:iam::1234566778899:mfa/some-arn
  session-name: foobar123
  clear-env: true

```
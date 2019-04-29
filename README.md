# sts-helper

`sts-helper` is a CLI tool used to easily use AWS STS from bash. 




## Setup 

### Create a `~/.sts-helper.yaml` file:

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

### Paste this function into your .bashrc

```
function sts_assume() {
     echo "Available sts-help profiles:"
     sts-helper list-helper-profiles
     read -p "Helper profile: " prof
     read -p "MFA token: " token

     eval $(sts-helper assume-role --helper-profile $prof --token $token)
}
```

```
source ~/.bashrc
sts_assume
```


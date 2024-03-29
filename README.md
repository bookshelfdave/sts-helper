# sts-helper

`sts-helper` is a CLI tool used to make assuming AWS roles more convenient.

## Prebuilt binaries

Only Darwin and Linux 64 bit binaries are available on the [releases](https://github.com/metadave/sts-helper/releases) page.

## Build and install from source

```bash
go get github.com/metadave/sts-helper
cd ${GOPATH}/src/github.com/metadave/sts-helper
dep ensure
go install
```

## Setup 

### Create an `~/.sts-helper.yaml` file:

```yaml
---
s3-stuff:
  duration-seconds: 3600
  role-arn: arn:aws:iam::1234566778899:role/foo
  mfa-arn: arn:aws:iam::1234566778899:mfa/some-arn
  session-name: foobar123
  clear-env: true
rds-readonly:
  duration-seconds: 3600
  role-arn: arn:aws:iam::1234566778899:role/foo
  mfa-arn: arn:aws:iam::1234566778899:mfa/some-arn
  session-name: foobar123
  clear-env: true

```

### Paste this function into your .bashrc

```bash
function sts_assume() {
     read -p "Helper profile: " prof
     read -p "MFA token: " token
     eval $(sts-helper assume-role --helper-profile $prof --token $token)
}
```

or just use a snippet like:

```bash
eval $(sts-helper assume-role --helper-profile my_profile --token 123456)
```

If you want to see a list of sts-helper profiles, use a function like this:

```bash
function sts_assume_with_list() {
     echo "Available sts-help profiles:"
     sts-helper list-helper-profiles
     read -p "Helper profile: " prof
     read -p "MFA token: " token
     eval $(sts-helper assume-role --helper-profile $prof --token $token)
}
```


## Example

```bash
source ~/.bashrc
export AWS_PROFILE=some_profile

$ sts_assume_with_list 
Available sts-help profiles:
s3-stuff arn:aws:iam:: 1234566778899:role/foobar123
rds-readonly arn:aws:iam::1234566778899:role/barbaz123
Helper profile: s3-stuff
MFA token: 123456

$ env | grep AWS
AWS_SESSION_TOKEN=...
AWS_SECRET_ACCESS_KEY=...
AWS_ACCESS_KEY_ID=...
```

> By default, `AWS_PROFILE` will be unset as part of the output from `sts-helper`

---

# License

[Apache Software License 2.0](https://github.com/metadave/ksv/blob/master/LICENSE)


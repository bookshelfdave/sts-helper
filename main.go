package main

import (
	"bufio"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"strings"
	"text/template"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sts"
	"github.com/urfave/cli"
	"gopkg.in/yaml.v2"
)

type profile struct {
	Duration    int64  `yaml:"duration"`
	RoleArn     string `yaml:"role-arn"`
	MFAArn      string `yaml:"mfa-arn"`
	SessionName string `yaml:"session-name"`
	ClearEnv    bool   `yaml:"clear-env"`
}

func readProfile(profileName string, path string) (*sts.AssumeRoleInput, error) {
	yamlFile, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	profiles := make(map[string]*profile)
	err = yaml.Unmarshal(yamlFile, &profiles)
	if err != nil {
		return nil, err
	}
	if p, ok := profiles[profileName]; ok {
		return &sts.AssumeRoleInput{
			DurationSeconds: &p.Duration,
			RoleArn:         &p.RoleArn,
			RoleSessionName: &p.SessionName,
			SerialNumber:    &p.MFAArn,
		}, nil
	}
	return nil, errors.New("profile not found")
}

func getInput(c *cli.Context) (*sts.AssumeRoleInput, error) {

	var assumeRoleInput *sts.AssumeRoleInput
	var err error
	if c.String("helper-profile") != "" {
		// read from profile
		assumeRoleInput, err = readProfile(c.String("helper-profile"), c.String("helper-profile-path"))
		if err != nil {
			return nil, err
		}
	} else {
		assumeRoleInput = &sts.AssumeRoleInput{
			DurationSeconds: aws.Int64(c.Int64("duration")),
			RoleArn:         aws.String(c.String("role")),
			RoleSessionName: aws.String(c.String("session-name")),
			SerialNumber:    aws.String(c.String("mfa")),
		}
	}

	var mfaTokenValue string
	if c.String("token") != "" {
		mfaTokenValue = c.String("token")
	} else {
		fmt.Printf("MFA token:")
		reader := bufio.NewReader(os.Stdin)
		mfaTokenValue, _ = reader.ReadString('\n')
	}
	mfaTokenValue = strings.TrimSpace(mfaTokenValue)
	assumeRoleInput.TokenCode = &mfaTokenValue
	return assumeRoleInput, nil
}

func doit(input *sts.AssumeRoleInput) (*sts.AssumeRoleOutput, error) {
	svc := sts.New(session.New())

	result, err := svc.AssumeRole(input)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func showIt(clearEnv bool, result *sts.AssumeRoleOutput) {
	if clearEnv {
		clearTemplate := template.Must(template.New("clear-template").Parse(`
unset AWS_SECRET_ACCESS_KEY;
unset AWS_SESSION_TOKEN;
unset AWS_PROFILE;`))
		clearTemplate.Execute(os.Stdout, nil)
	}

	setTemplate := template.Must(template.New("setenv-template").Parse(`
export AWS_ACCESS_KEY_ID='{{.AccessKeyId}}'
export AWS_SECRET_ACCESS_KEY='{{.SecretAccessKey}}'
export AWS_SESSION_TOKEN='{{.SessionToken}}'
`))
	setTemplate.Execute(os.Stdout, *result.Credentials)
}

func main() {

	app := cli.NewApp()
	app.Name = "sts-helper"
	app.HelpName = "sts-helper"
	app.Usage = "Assume an AWS role and display shell code to eval"
	app.Version = "0.1.0"
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   "role, r",
			Usage:  "Role to assume",
			EnvVar: "STS_ROLE_ARN_TO_ASSUME",
		},
		cli.StringFlag{
			Name:   "mfa, m",
			Usage:  "MFA arn",
			EnvVar: "STS_MFA_ARN",
		},
		cli.StringFlag{
			Name:  "token, t",
			Usage: "MFA token value",
		},
		cli.StringFlag{
			Name:  "session-name, n",
			Usage: "STS session name",
		},
		cli.Int64Flag{
			Name:  "duration, d",
			Usage: "Duration in seconds",
			Value: 3600,
		},
		cli.StringFlag{
			Name:  "helper-profile, p",
			Usage: "sts-helper profile",
		},
		cli.StringFlag{
			Name:  "helper-profile-path",
			Usage: "Path to sts-helper config file",
			Value: path.Join(os.Getenv("HOME"), ".sts-helper.yaml"),
		},

		cli.BoolTFlag{
			Name:  "clear-env, c",
			Usage: "Clear current aws profile values",
		},
	}
	app.Action = func(c *cli.Context) error {
		input, err := getInput(c)
		if err != nil {
			panic(err)
		}
		result, err := doit(input)
		if err != nil {
			panic(err)
		}
		showIt(c.BoolT("clear-env"), result)
		return nil
	}
	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

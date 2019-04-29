package main

import (
	"fmt"
	"log"
	"os"
	"text/template"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sts"
	"github.com/urfave/cli"
)

func doit(c *cli.Context) {
	svc := sts.New(session.New())
	input := &sts.AssumeRoleInput{
		DurationSeconds: aws.Int64(c.Int64("duration")),
		RoleArn:         aws.String(c.String("role")),
		RoleSessionName: aws.String("Foo"),
		TokenCode:       aws.String(c.String("token")),
		SerialNumber:    aws.String(c.String("mfa")),
	}

	result, err := svc.AssumeRole(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case sts.ErrCodeMalformedPolicyDocumentException:
				fmt.Println(sts.ErrCodeMalformedPolicyDocumentException, aerr.Error())
			case sts.ErrCodePackedPolicyTooLargeException:
				fmt.Println(sts.ErrCodePackedPolicyTooLargeException, aerr.Error())
			case sts.ErrCodeRegionDisabledException:
				fmt.Println(sts.ErrCodeRegionDisabledException, aerr.Error())
			default:
				fmt.Println(aerr.Error())
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			fmt.Println(err.Error())
		}
		return
	}

	t, err := template.New("foo").Parse(`
export AWS_ACCESS_KEY_ID='{{.AccessKeyId}}'
export AWS_SECRET_ACCESS_KEY='{{.SecretAccessKey}}'
export AWS_SESSION_TOKEN='{{.SessionToken}}'`)
	if err != nil {
		panic(err)
	}
	t.Execute(os.Stdout, *result.Credentials)
}

func main() {
	app := cli.NewApp()
	app.Name = "sts-helper"
	app.Usage = "Assume an AWS role and set environment variables"
	app.Version = "0.1.0"
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   "role, r",
			Usage:  "language for the greeting",
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
		cli.Int64Flag{
			Name:  "duration, d",
			Usage: "Duration",
			Value: 3600,
		},
	}
	app.Action = func(c *cli.Context) error {
		doit(c)
		return nil
	}
	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

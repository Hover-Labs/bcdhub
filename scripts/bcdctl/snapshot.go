package main

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/baking-bad/bcdhub/internal/search"
	"github.com/pkg/errors"
)

type snapshotCommand struct{}

var snapshotCmd snapshotCommand

// Execute
func (x *snapshotCommand) Execute(_ []string) error {
	if err := uploadMappings(ctx.Searcher, creds); err != nil {
		return err
	}
	if err := listRepositories(ctx.Searcher); err != nil {
		return err
	}
	name, err := askQuestion("Please, enter target repository name:")
	if err != nil {
		return err
	}
	snapshotName := fmt.Sprintf("snapshot_%s", strings.ToLower(time.Now().UTC().Format(time.RFC3339)))
	return ctx.Searcher.CreateSnapshots(name, snapshotName, search.Indices)
}

type restoreCommand struct{}

var restoreCmd restoreCommand

// Execute
func (x *restoreCommand) Execute(_ []string) error {
	if err := listRepositories(ctx.Searcher); err != nil {
		return err
	}
	name, err := askQuestion("Please, enter target repository name:")
	if err != nil {
		return err
	}

	if err := listSnapshots(ctx.Searcher, name); err != nil {
		return err
	}
	snapshotName, err := askQuestion("Please, enter target snapshot name:")
	if err != nil {
		return err
	}
	return ctx.Searcher.RestoreSnapshots(name, snapshotName, search.Indices)
}

type setPolicyCommand struct{}

var setPolicyCmd setPolicyCommand

// Execute
func (x *setPolicyCommand) Execute(_ []string) error {
	if err := listPolicies(ctx.Searcher); err != nil {
		return err
	}
	policyID, err := askQuestion("Please, enter target new or existing policy ID:")
	if err != nil {
		return err
	}
	repository, err := askQuestion("Please, enter target repository name:")
	if err != nil {
		return err
	}
	schedule, err := askQuestion("Please, enter schedule in cron format (https://www.elastic.co/guide/en/elasticsearch/reference/current/trigger-schedule.html#schedule-cron):")
	if err != nil {
		return err
	}
	expiredAfter, err := askQuestion("Please, enter expiration in days:")
	if err != nil {
		return err
	}
	iExpiredAfter, err := strconv.ParseInt(expiredAfter, 10, 64)
	if err != nil {
		return err
	}
	return ctx.Searcher.SetSnapshotPolicy(policyID, schedule, policyID, repository, iExpiredAfter)
}

type reloadSecureSettingsCommand struct{}

var reloadSecureSettingsCmd reloadSecureSettingsCommand

// Execute
func (x *reloadSecureSettingsCommand) Execute(_ []string) error {
	return ctx.Searcher.ReloadSecureSettings()
}

func listPolicies(storage search.Searcher) error {
	policies, err := storage.GetAllPolicies()
	if err != nil {
		return err
	}

	fmt.Println("")
	fmt.Println("Available snapshot policies")
	fmt.Println("=======================================")
	for i := range policies {
		fmt.Println(policies[i])
	}
	fmt.Println("")
	return nil
}

func listRepositories(storage search.Searcher) error {
	listRepos, err := storage.ListRepositories()
	if err != nil {
		return err
	}

	fmt.Println("")
	fmt.Println("Available repositories")
	fmt.Println("=======================================")
	for i := range listRepos {
		fmt.Print(listRepos[i].String())
	}
	fmt.Println("")
	return nil
}

func listSnapshots(storage search.Searcher, repository string) error {
	listSnaps, err := storage.ListSnapshots(repository)
	if err != nil {
		return err
	}
	fmt.Println("")
	fmt.Println(listSnaps)
	fmt.Println("")
	return nil
}

func uploadMappings(storage search.Searcher, creds awsData) error {
	mappings, err := storage.GetMappings(search.Indices)
	if err != nil {
		return err
	}

	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String(creds.Region),
		Credentials: credentials.NewEnvCredentials(),
	})
	if err != nil {
		return err
	}
	uploader := s3manager.NewUploader(sess)

	for key, value := range mappings {
		fileName := fmt.Sprintf("mappings/%s.json", key)
		body := strings.NewReader(value)

		if _, err := uploader.Upload(&s3manager.UploadInput{
			Bucket:      aws.String(creds.BucketName),
			Key:         aws.String(fileName),
			Body:        body,
			ContentType: aws.String("application/json"),
		}); err != nil {
			return errors.Errorf("failed to upload file, %v", err)
		}
	}
	return nil
}

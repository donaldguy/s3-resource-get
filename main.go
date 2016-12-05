package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/concourse/fly/rc"
	"github.com/concourse/go-concourse/concourse"
	"github.com/concourse/s3-resource"

	flags "github.com/jessevdk/go-flags"
)

type Options struct {
	Target   rc.TargetName `short:"t" long:"target" description:"Fly target to monitor" env:"FLIGHT_TRACKER_SERVER"`
	FileName string        `short:"o" long:"output" description:"Where to output file. Default is $PWD/<basename of remote file>"`
}

type S3Source struct {
	Bucket          string
	AccessKeyID     string
	SecretAccessKey string
	RegionName      string
	Endpoint        string
}

func dieIf(err error) {
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func main() {
	var opts Options

	parser := flags.NewParser(&opts, flags.HelpFlag|flags.PassDoubleDash)
	parser.NamespaceDelimiter = "-"

	args, err := parser.Parse()
	dieIf(err)

	target, err := rc.LoadTarget(opts.Target)
	dieIf(err)

	if len(args) > 1 {
		fmt.Fprintf(os.Stderr, "Extra argument, ignoring %v", args[1:])
	}
	var pipelineName, resourceName string
	if len(args) > 0 {
		resourceComponents := strings.Split(args[0], "/")
		if len(resourceComponents) == 2 {
			pipelineName = resourceComponents[0]
			resourceName = resourceComponents[1]
		}
	}

	if len(args) < 1 || pipelineName == "" {
		dieIf(fmt.Errorf("Must provide a resource in pipeline/name format"))
	}

	team := target.Team()
	config, _, _, _, err := team.PipelineConfig(pipelineName)

	var source s3resource.Source
	for _, res := range config.Resources {
		if res.Name == resourceName {
			if res.Type != "s3" {
				dieIf(fmt.Errorf("Resource %s is type %s not type s3", resourceName, res.Type))
			}

			j, e := json.Marshal(res.Source)
			dieIf(e)
			e = json.Unmarshal(j, &source)
			dieIf(e)
			break
		}
	}

	if source.Bucket == "" {
		dieIf(fmt.Errorf("Could not find a resource named %s in %s", resourceName, pipelineName))
	}

	awsConfig := s3resource.NewAwsConfig(
		source.AccessKeyID,
		source.SecretAccessKey,
		source.RegionName,
		source.Endpoint,
		source.DisableSSL,
	)

	client := s3resource.NewS3Client(
		os.Stdout,
		awsConfig,
		source.UseV2Signing,
	)

	// TODO: allow grabbing specific version
	resVers, _, _, err := team.ResourceVersions(pipelineName, resourceName, concourse.Page{Limit: 10})
	dieIf(err)

	// TODO: support versionIds
	resPath := resVers[0].Version["path"]

	var localPath string
	if opts.FileName == "" {
		wd, err := os.Getwd()
		dieIf(err)
		name := path.Base(resPath)
		localPath = path.Join(wd, name)
		fmt.Printf("Downloading %s\n:", name)
	} else {
		localPath = opts.FileName
	}

	client.DownloadFile(
		source.Bucket,
		resPath,
		"",
		localPath,
	)
}
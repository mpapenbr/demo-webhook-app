package releaseupdater

import (
	"context"
	b64 "encoding/base64"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/google/go-github/v44/github"
	"github.com/ktrysmt/go-bitbucket"
)

func ProcessNewRelease(ctx Context, release *github.ReleaseEvent) {
	fmt.Printf("Incoming release event from %s\n", *release.Repo.FullName)
	for _, action := range ctx.Config.Actions {
		if action.From == *release.Repo.Name {
			commitComponent := release.Repo.Name
			if action.Component != "" {
				commitComponent = &action.Component
			}
			for _, update := range action.Update {
				log.Printf("%v\n", update)
				repoOwner := *release.Repo.Owner.Login
				replacer := func(content string) (string, string) {
					return ReplaceVersionString(content, update.Regex, *release.Release.TagName),
						fmt.Sprintf("pkg: Bump %s to %s", *commitComponent, *release.Release.TagName)
				}
				switch strings.ToLower(update.RepoType) {
				case "bitbucket":
					handleBitbucket(ctx, update, replacer)
				default:
					handleGithub(ctx, repoOwner, update, replacer)
				}

			}
		}
	}
}

func handleGithub(ctx Context, repoOwner string, update Update, processContent func(content string) (string, string)) {
	for _, toUpdateFile := range update.Files {
		content, _, resp, err := ctx.ProbotCtx.GitHub.Repositories.GetContents(
			context.Background(),
			repoOwner,
			update.Repo,
			toUpdateFile,
			&github.RepositoryContentGetOptions{})
		if err != nil {
			log.Printf("error reading source: %+v\n", err)
			continue
		}
		fileContent, _ := b64.StdEncoding.DecodeString(*content.Content)
		log.Printf("Content: <%s>\n", fileContent)
		log.Printf("RegEx: <%s>\n", update.Regex)
		log.Printf("Resp: %+v\n", *resp)
		newVersion, message := processContent(string(fileContent))
		if string(newVersion) != string(fileContent) {
			fmt.Printf("Updating file %s\n", toUpdateFile)
			fmt.Printf("NewContent: <%s>\n", string(newVersion))

			_, resp, err = ctx.ProbotCtx.GitHub.Repositories.UpdateFile(
				context.Background(),
				repoOwner,
				update.Repo,
				toUpdateFile,
				&github.RepositoryContentFileOptions{
					Content: []byte(newVersion),
					Message: github.String(message),
					SHA:     github.String(*content.SHA),
				})

			if err != nil {
				log.Printf("error updating %s: %+v\n", toUpdateFile, err)
				continue
			}
		} else {
			log.Println("No changes detected")
		}
	}
}

func handleBitbucket(ctx Context, update Update, processContent func(content string) (string, string)) {
	for _, toUpdateFile := range update.Files {
		fileContent, err := ctx.BitbucketClient.Repositories.Repository.GetFileBlob(
			&bitbucket.RepositoryBlobOptions{Owner: "mpapenbr",
				RepoSlug: update.Repo, Path: toUpdateFile, Ref: "master"})
		if err != nil {
			log.Printf("error reading source: %+v\n", err)
			continue
		}
		log.Println(fileContent.String())
		newVersion, message := processContent(string(fileContent.String()))
		if string(newVersion) != string(fileContent.String()) {
			fmt.Printf("Updating file %s\n", toUpdateFile)
			fmt.Printf("NewContent: <%s>\n", string(newVersion))

			f, _ := os.CreateTemp("", "bbupload")
			f.WriteString(newVersion)
			f.Close()

			err := ctx.BitbucketClient.Repositories.Repository.WriteFileBlob(&bitbucket.RepositoryBlobWriteOptions{
				Owner:    "mpapenbr",
				RepoSlug: update.Repo,
				FilePath: f.Name(),
				FileName: toUpdateFile,
				Branch:   "master",
				Message:  message,
			})
			log.Printf("Deleting temp file %s\n", f.Name())
			err = os.Remove(f.Name())
			if err != nil {
				log.Printf("Error deleting temp file %s: %v\n", f.Name(), err)
			}
		} else {
			log.Println("No changes detected")
		}
	}
}

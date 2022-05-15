package releaseupdater

import (
	"context"
	b64 "encoding/base64"
	"fmt"
	"log"

	"github.com/google/go-github/github"
	"github.com/swinton/go-probot/probot"
)

func ProcessNewRelease(config *Config, ctx *probot.Context, release *github.ReleaseEvent) {
	fmt.Printf("Incoming release event from %s\n", *release.Repo.FullName)
	for _, action := range config.Actions {
		if action.From == *release.Repo.Name {
			commitComponent := release.Repo.Name
			if action.Component != "" {
				commitComponent = &action.Component
			}
			for _, update := range action.Update {
				repoOwner := *release.Repo.Owner.Login
				log.Printf("Fetching %s from %s/%s\n", update.File, repoOwner, update.Repo)
				content, _, resp, err := ctx.GitHub.Repositories.GetContents(context.Background(), repoOwner, update.Repo, update.File, &github.RepositoryContentGetOptions{})
				if err != nil {
					log.Printf("error reading source: %+v\n", err)
					continue
				}
				fileContent, _ := b64.StdEncoding.DecodeString(*content.Content)
				log.Printf("Content: <%s>\n", fileContent)
				log.Printf("RegEx: <%s>\n", update.Regex)
				log.Printf("Resp: %+v\n", *resp)
				newVersion := ReplaceVersion(fileContent, update.Regex, *release.Release.TagName)
				if string(newVersion) != string(fileContent) {
					fmt.Printf("Updating file %s\n", update.File)
					fmt.Printf("NewContent: <%s>\n", string(newVersion))

					message := fmt.Sprintf("pkg: Bump %s to %s", *commitComponent, *release.Release.TagName)
					_, resp, err = ctx.GitHub.Repositories.UpdateFile(context.Background(), repoOwner, update.Repo, update.File, &github.RepositoryContentFileOptions{
						Content: []byte(newVersion),
						Message: github.String(message),
						SHA:     github.String(*content.SHA),
					})
					if err != nil {
						log.Printf("error updating %s: %+v\n", update.File, err)
						continue
					}

				} else {
					log.Println("No changes detected")
				}

			}
		}
	}
}

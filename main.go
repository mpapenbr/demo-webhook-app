package main

import (
	"context"
	"fmt"
	"log"

	"demo-webhook-app/releaseupdater"

	"github.com/google/go-github/v44"
	"github.com/mpapenbr/go-probot/probot"
)

func main() {
	config, err := releaseupdater.GetConfig("config.yml")
	if err != nil {
		log.Fatal("Could not read config")
	}

	probot.HandleEvent("release", func(ctx *probot.Context) error {
		// Because we're listening for "release" we know the payload is a *github.ReleaseEvent
		event := ctx.Payload.(*github.ReleaseEvent)
		if *event.Action == "published" {

			// log.Printf("got release published %+v\n", event)
			log.Printf("got release published from %s\n", *event.GetRepo().Name)
			releaseupdater.ProcessNewRelease(config, ctx, event)
			repoOwner := "mpapenbr"
			repo := "demo_deploy1"
			fileRef := "versions.properties"
			content, _, resp, err := ctx.GitHub.Repositories.GetContents(context.Background(), repoOwner, repo, fileRef, &github.RepositoryContentGetOptions{})
			_, resp, err = ctx.GitHub.Repositories.UpdateFile(context.Background(), repoOwner, repo, fileRef, &github.RepositoryContentFileOptions{
				Content: []byte("demo_app_version1: " + *event.Release.TagName),
				Message: github.String("changed by bot for " + *event.Release.TagName),
				SHA:     github.String(*content.SHA),
			})
			if err != nil {

				log.Fatalf("UpdateFileRest %v", err)
			}

			if false {
				fmt.Printf("%v", resp)
			}
		} else {
			log.Printf("I'm only interested in published releases.\n")
		}

		return nil
	})

	probot.Start()
}

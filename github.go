package main

import (
	"context"
	"encoding/json"
	"html/template"
	"log"
	"strings"

	"github.com/google/go-github/v56/github"
)

const GithubUsername = "faangbait"

func getClient() *github.Client {
	client := github.NewClient(nil)
	return client
}

func getBasicInfo(c *github.Client) (ghUser, *github.Response) {
	u, resp, err := c.Users.Get(context.TODO(), GithubUsername)
	if err != nil {
		log.Println(err.Error())
		return ghUser{}, resp
	}

	return ghUser{
		Name:      u.GetLogin(),
		URL:       u.GetHTMLURL(),
		AvatarURL: u.GetAvatarURL(),
	}, resp
}

func getRepositories(c *github.Client) ([]ghRepository, *github.Response) {
	var repoList []ghRepository
	opts := github.RepositoryListOptions{
		Visibility:  "public",
		Affiliation: "owner,collaborator",
		ListOptions: github.ListOptions{
			Page:    1,
			PerPage: 100,
		},
	}

	repos, resp, err := c.Repositories.List(context.TODO(), GithubUsername, &opts)
	if err != nil {
		log.Println(err.Error())
		return repoList, resp
	}

	for _, r := range repos {
		repoList = append(repoList, ghRepository{
			Owner: ghUser{
				Name:      r.Owner.GetLogin(),
				URL:       r.Owner.GetHTMLURL(),
				AvatarURL: r.Owner.GetAvatarURL(),
			},
			Name:        r.GetName(),
			Language:    r.GetLanguage(),
			Description: r.GetDescription(),
			URL:         r.GetHTMLURL(),
			Fork:        r.GetFork(),
			Archived:    r.GetArchived(),
		})
	}

	return repoList, resp
}

func getEvents(c *github.Client) ([]ghEvent, *github.Response) {
	var eventList []ghEvent
	opts := github.ListOptions{
		Page:    1,
		PerPage: 100,
	}

	events, resp, err := c.Activity.ListEventsPerformedByUser(context.TODO(), GithubUsername, false, &opts)
	if err != nil {
		log.Println(err.Error())
		return eventList, resp
	}

	for _, e := range events {
		var pushEvent github.PushEvent
		var commits []ghCommit
		pushPayload := ghPush{}

		err := json.Unmarshal(e.GetRawPayload(), &pushEvent)

		if err == nil {
			for _, hc := range pushEvent.GetCommits() {
				commits = append(commits, ghCommit{
					Message: template.HTML(strings.Replace(hc.GetMessage(), "\n", "<br />", -1)),
					SHA:     hc.GetSHA(),
				})
			}

			pushPayload = ghPush{
				Head:    pushEvent.GetHead(),
				Commits: commits,
			}
		}

		repoName, _ := strings.CutPrefix(e.Repo.GetName(), "faangbait/")

		eventList = append(eventList, ghEvent{
			Actor: ghUser{
				Name:      e.Actor.GetLogin(),
				URL:       e.Actor.GetHTMLURL(),
				AvatarURL: e.Actor.GetAvatarURL(),
			},
			Repo: ghRepository{
				Owner: ghUser{
					Name:      e.Repo.Owner.GetLogin(),
					URL:       e.Repo.Owner.GetHTMLURL(),
					AvatarURL: e.Repo.Owner.GetAvatarURL(),
				},
				Name:        repoName,
				Language:    e.Repo.GetLanguage(),
				Description: e.Repo.GetDescription(),
				URL:         e.Repo.GetHTMLURL(),
				Fork:        e.Repo.GetFork(),
				Archived:    e.Repo.GetArchived(),
			},
			Type:    e.GetType(),
			Payload: pushPayload,
		})
	}

	return eventList, resp
}

type ghData struct {
	UserInfo ghUser
	Repos    []ghRepository
	Events   []ghEvent
}

type ghUser struct {
	Name      string
	URL       string
	AvatarURL string
}

type ghRepository struct {
	Owner       ghUser
	Name        string
	Language    string
	Description string
	URL         string
	Fork        bool
	Archived    bool
}

type ghEvent struct {
	Actor   ghUser
	Repo    ghRepository
	Type    string
	Payload ghPush
}

type ghCommit struct {
	Message template.HTML
	SHA     string
}

type ghPush struct {
	Head    string
	Commits []ghCommit
}

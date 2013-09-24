package main

import (
	"sort"

	"code.google.com/p/goauth2/oauth"
	"github.com/google/go-github/github"
)

type Repos []github.Repository

type ReposByUpdated struct {
	Repos Repos
}

func (s ReposByUpdated) Less(i, j int) bool {
	return s.Repos[i].UpdatedAt.UnixNano() > s.Repos[j].UpdatedAt.UnixNano()
}

func (s ReposByUpdated) Len() int      { return len(s.Repos) }
func (s ReposByUpdated) Swap(i, j int) { s.Repos[i], s.Repos[j] = s.Repos[j], s.Repos[i] }

func GithubUpdate(token string, limit int) (*github.User, *[]github.Repository) {
	t := &oauth.Transport{
		Token: &oauth.Token{AccessToken: token},
	}

	client := github.NewClient(t.Client())
	user, _, _ := client.Users.Get("")
	repos, _, _ := client.Repositories.List("", &github.RepositoryListOptions{Type: "public"})

	sort.Sort(ReposByUpdated{Repos: repos})

	if len(repos) > limit {
		repos = repos[:limit]
	}

	return user, &repos
}

package clients

import (
	"fmt"
	"time"

	"code.google.com/p/goauth2/oauth"
	"github.com/google/go-github/github"
)

var (
	GithubUser  *github.User
	GithubRepos *[]github.Repository
)

func GithubUpdate(token string) {
	t := &oauth.Transport{
		Token: &oauth.Token{AccessToken: token},
	}

	client := github.NewClient(t.Client())
	user, _, _ := client.Users.Get("")
	repos, _, _ := client.Repositories.List("", &github.RepositoryListOptions{Type: "public"})
	repos = repos[:5]

	GithubUser = user
	GithubRepos = &repos

	fmt.Println("Github updated", time.Now())
}

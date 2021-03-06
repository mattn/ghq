package main

import (
	"errors"
	"fmt"
	"net/url"
	"os/exec"
	"strings"

	"github.com/motemen/ghq/utils"
	. "github.com/onsi/gomega"
)
import "testing"

func parseURL(urlString string) *url.URL {
	u, err := url.Parse(urlString)
	if err != nil {
		panic(err)
	}
	return u
}

func TestNewRemoteRepositoryGitHub(t *testing.T) {
	RegisterTestingT(t)

	var (
		repo RemoteRepository
		err  error
	)

	repo, err = NewRemoteRepository(parseURL("https://github.com/motemen/pusheen-explorer"))
	Expect(err).To(BeNil())
	Expect(repo.IsValid()).To(Equal(true))
	Expect(repo.VCS()).To(Equal(GitBackend))

	repo, err = NewRemoteRepository(parseURL("https://github.com/motemen/pusheen-explorer/"))
	Expect(err).To(BeNil())
	Expect(repo.IsValid()).To(Equal(true))
	Expect(repo.VCS()).To(Equal(GitBackend))

	repo, err = NewRemoteRepository(parseURL("https://github.com/motemen/pusheen-explorer/blob/master/README.md"))
	Expect(err).To(BeNil())
	Expect(repo.IsValid()).To(Equal(false))

	repo, err = NewRemoteRepository(parseURL("https://example.com/motemen/pusheen-explorer"))
	Expect(err).NotTo(BeNil())
}

func TestNewRemoteRepositoryGoogleCode(t *testing.T) {
	RegisterTestingT(t)

	var (
		repo RemoteRepository
		err  error
	)

	repo, err = NewRemoteRepository(parseURL("https://code.google.com/p/vim/"))
	Expect(err).To(BeNil())
	Expect(repo.IsValid()).To(Equal(true))
	utils.CommandRunner = NewFakeRunner(map[string]error{
		"hg identify":   nil,
		"git ls-remote": errors.New(""),
	})
	Expect(repo.VCS()).To(Equal(MercurialBackend))

	repo, err = NewRemoteRepository(parseURL("https://code.google.com/p/git-core"))
	Expect(err).To(BeNil())
	Expect(repo.IsValid()).To(Equal(true))
	utils.CommandRunner = NewFakeRunner(map[string]error{
		"hg identify":   errors.New(""),
		"git ls-remote": nil,
	})
	Expect(repo.VCS()).To(Equal(GitBackend))
}

func NewFakeRunner(dispatch map[string]error) utils.RunFunc {
	return func(cmd *exec.Cmd) error {
		cmdString := strings.Join(cmd.Args, " ")
		for cmdPrefix, err := range dispatch {
			if strings.Index(cmdString, cmdPrefix) == 0 {
				return err
			}
		}
		panic(fmt.Sprintf("No fake dispatch found for: %s", cmdString))
	}
}

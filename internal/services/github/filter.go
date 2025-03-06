package github

import "slices"

var ReposToExclude = []string{
	"homebrew-formulae",
	"excalith-start-page",
}

func ShouldIncludeRepo(repoName string) bool {
	return !slices.Contains(ReposToExclude, repoName)
}

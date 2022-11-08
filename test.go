package gotest

import (
	"fmt"
	"strings"
)

// A Test represents an actual test.
type Test struct {
	id          string
	Name        string
	Description string
	Author      string
	Tags        []string
	Setup       SetupFunction
	Function    TestFunction
	SuiteName   string
}

// matchTags matches all tags if --match-all is set otherwise matches any tag
func (t Test) matchTags(tags []string, matchAll bool) bool {

	m := make([]string, len(tags))
	copy(m, tags)

	if !matchAll {
		return t.matchAnyTags(m)
	}

	return t.matchAllTags(m)
}

// matchAllTags returns true if all incoming tags are matching minus exclusions
func (t Test) matchAllTags(tags []string) bool {

	if len(tags) == 0 {
		return true
	}

	for _, incoming := range tags {
		if strings.HasPrefix(incoming, "~") {
			if t.hasTag(strings.TrimPrefix(incoming, "~")) {
				return false
			}

			continue
		}

		if !t.hasTag(incoming) {
			return false
		}
	}

	return true
}

// matchAnyTags returns true if any incoming tags are matching
func (t Test) matchAnyTags(tags []string) bool {

	if len(tags) == 0 {
		return true
	}

	for _, incoming := range tags {
		if t.hasTag(incoming) {
			return true
		}
	}

	return false
}

// hasTag returns true if the slice has the tag
func (t Test) hasTag(tag string) bool {

	tags := t.Tags
	if t.SuiteName != "" {
		tags = append(tags, fmt.Sprintf("suite:%s", strings.ToLower(t.SuiteName)))
	}

	for _, testTag := range tags {
		if tag == testTag {
			return true
		}
	}

	return false
}

func (t Test) String() string {
	return fmt.Sprintf(`  id         : %s
  name       : %s
  desc       : %s
  author     : %s
  categories : %s
`, t.id, t.Name, t.Description, t.Author, strings.Join(t.Tags, ", "))
}

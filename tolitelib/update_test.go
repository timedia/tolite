package tolitelib_test

import (
	"fmt"
	"log"
	"strings"
	"testing"

	. "."
)

const yamlData = `
users:
  k_shigata:
    email: user@timedia.co.jp
    keys: key1
  kansei_s:
    keys:
      - key2
      - key3: raw+public_key+string
groups:
  group_1:
    - k_shigata
    - kansei_s
  group_2: [kansei_s]
  group_3: [k_shigata, kansei_s]
  group_4: kansei_s
  internal-staff: k_shigata
repos:
  "@all":
    R: internal-staff
  some_repo:
    R: group_1
    RW: [group_2]
    RW+:
      - group_3
      - group_4
    "-": all
adminRepos:
  "@all":
    R: internal-staff
  some_repo:
    R: group_1
    RW: [group_2]
    RW+:
      - group_3
      - group_4
    "-": all
`

func TestGenerateConfHeader(t *testing.T) {
	const expect = `#
# Users
#

`
	arrived := GenerateConfHeader("Users")
	if arrived != expect {
		log.Fatalf("error: expect:'%s' arrived:'%s'", expect, arrived)
	}
}

func TestParseYaml(t *testing.T) {
	data, err := ParseYaml([]byte(yamlData))
	fmt.Printf("%v", data)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
}

func TestGenerateUser(t *testing.T) {
	id := "k_shigata"
	user := User{Keys: []string{"k_shigata", "kansei_s"}}
	expect := "@k_shigata           = k_shigata kansei_s"
	arrived := GenerateUser(id, user)
	if arrived != expect {
		log.Fatalf("error: expected: %s, arrived: %s", expect, arrived)
	}
}

var userData = map[string]User{
	"k_shigata": User{
		Keys: "key1",
	},
	"kansei_s": User{
		Keys: []string{"key2", "key3"},
	},
}

func TestGenerateUsers(t *testing.T) {
	data := userData
	expect := `@k_shigata           = key1
@kansei_s            = key2 key3`
	arrived := GenerateUsers(data)
	for _, e := range strings.Split(expect, "\n") {
		find := false
		for _, a := range strings.Split(arrived, "\n") {
			if e == a {
				find = true
				break
			}
		}
		if !find {
			log.Fatalf("error: expected: \n%s\n, arrived: \n%s\n", expect, arrived)
		}
	}
}

func TestGenerateRepoGroups(t *testing.T) {
	perm := "R"
	groups := []string{"group_a", "group_b"}
	expect := "\tR\t= @group_a @group_b\n"
	arrived := GenerateRepoGroups(perm, groups)
	if expect != arrived {
		log.Fatalf("error: expected: \n%s\n, arrived: \n%s\n", expect, arrived)
	}
}

func TestGenerateRepo(t *testing.T) {
	id := "@all"
	repo := Repo{
		Read: []string{"internal-staff"},
	}
	expect := `repo	@all
	R	= @internal-staff
`
	arrived := GenerateRepo(id, repo)
	if expect != arrived {
		log.Fatalf("error: expected: \n%s\n, arrived: \n%s\n", expect, arrived)
	}
	id = "some_repo"
	repo = Repo{
		Read:          "group_1",
		ReadWrite:     []string{"group_2", "group_3"},
		ReadWritePlus: []string{"group_4"},
		Deny:          "all",
	}
	expect = `repo	some_repo
	R	= @group_1
	RW	= @group_2 @group_3
	RW+	= @group_4
	-	= @all
`
	arrived = GenerateRepo(id, repo)
	if expect != arrived {
		log.Fatalf("error: expected: \n%s\n, arrived: \n%s\n", expect, arrived)
	}
}

var reposData = map[string]Repo{
	"@all": Repo{Read: []string{"internal-staff"}},
	"some_repo": Repo{
		Read:          "group_1",
		ReadWrite:     []string{"group_2"},
		ReadWritePlus: []string{"group_3", "group_4"},
		Deny:          "all",
	},
}

func TestGenerateRepos(t *testing.T) {
	repos := reposData
	expect := `repo	@all
	R	= @internal-staff

repo	some_repo
	R	= @group_1
	RW	= @group_2
	RW+	= @group_3 @group_4
	-	= @all

`
	arrived := GenerateRepos(repos)
	if expect != arrived {
		log.Fatalf("error: expected: \n%s\n, arrived: \n%s\n", expect, arrived)
	}
}

func TestGenerateGroup(t *testing.T) {
	id := "a_group"
	group := "k_shigata"
	expect := `@a_group             = @k_shigata`
	arrived := GenerateGroup(id, group)
	if expect != arrived {
		log.Fatalf("error: expected: \n%s\n, arrived: \n%s\n", expect, arrived)
	}
}

var groupsData = map[string]interface{}{
	"group_1":        []string{"k_shigata", "kansei_s"},
	"group_2":        "kansei_s",
	"group_3":        []string{"k_shigata", "kansei_s"},
	"group_4":        "kansei_s",
	"internal-staff": "k_shigata",
}

func TestGenerateGroups(t *testing.T) {
	groups := groupsData
	expect := `@group_1             = @k_shigata @kansei_s
@group_2             = @kansei_s
@group_3             = @k_shigata @kansei_s
@group_4             = @kansei_s
@internal-staff      = @k_shigata

`
	arrived := GenerateGroups(groups)
	if expect != arrived {
		log.Fatalf("error: expected: \n%s\n, arrived: \n%s\n", expect, arrived)
	}
}

const confData = `#
# Users
#

@k_shigata           = key1
@kansei_s            = key2 key3

#
# Groups
#

@group_1             = @k_shigata @kansei_s
@group_2             = @kansei_s
@group_3             = @k_shigata @kansei_s
@group_4             = @kansei_s
@internal-staff      = @k_shigata

#
# Repos
#

repo	@all
	R	= @internal-staff

repo	some_repo
	R	= @group_1
	RW	= @group_2
	RW+	= @group_3 @group_4
	-	= @all

#
# Admin Repos
#

repo	@all
	R	= @internal-staff

repo	some_repo
	R	= @group_1
	RW	= @group_2
	RW+	= @group_3 @group_4
	-	= @all
`

func TestGenerateConf(t *testing.T) {
	data := Data{
		Users:      userData,
		Groups:     groupsData,
		Repos:      reposData,
		AdminRepos: reposData,
	}

	expected := confData
	arrived := GenerateConf(data)
	if expected != arrived {
		log.Fatalf("error: expected: \n%s\n, arrived: \n%s\n", expected, arrived)
	}
}

func TestYamlToConf(t *testing.T) {
	data := yamlData
	expected := confData
	out, err := ParseYaml([]byte(data))
	if err != nil {
		log.Fatalf("error")
	}
	fmt.Printf("%v", out)
	arrived := GenerateConf(out)
	if expected != arrived {
		log.Fatalf("error: expected: \n%s\n, arrived: \n%s\n", expected, arrived)
	}
}

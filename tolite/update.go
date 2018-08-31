package tolite

import (
	"fmt"
	"log"
	"sort"
	"strings"

	"gopkg.in/yaml.v2"
)

type Repo struct {
	ReadWritePlus interface{}
	ReadWrite     interface{}
	Read          interface{}
	Deny          interface{}
}
type User struct {
	Email string
	Keys  interface{}
}

type Data struct {
	Users      map[string]User        `yaml:"users"`
	Groups     map[string]interface{} `yaml:"groups"`
	Repos      map[string]Repo        `yaml:"repos"`
	AdminRepos map[string]Repo        `yaml:"adminRepos"`
}

func ParseYaml(in []byte) (Data, error) {
	var o = make(map[string]interface{})
	err := yaml.Unmarshal(in, o)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	users := map[string]User{}
	for key, value := range o["users"].(map[interface{}]interface{}) {
		user := value.(map[interface{}]interface{})
		u := User{}
		if e, ok := user["email"].(string); ok {
			u.Email = e
			fmt.Println("email!", e)
		}
		if k, ok := user["keys"].([]interface{}); ok {
			u.Keys = k
			fmt.Println("keys!", k)
		}
		if k, ok := user["keys"].(string); ok {
			u.Keys = k
			fmt.Println("keys!", k)
		}
		users[key.(string)] = u
	}
	groups := map[string]interface{}{}
	for key, value := range o["groups"].(map[interface{}]interface{}) {
		groups[key.(string)] = value
	}
	repos := map[string]Repo{}
	for key, value := range o["repos"].(map[interface{}]interface{}) {
		repo := value.(map[interface{}]interface{})
		re := Repo{}
		if r, ok := repo["R"].(interface{}); ok {
			re.Read = r
		}
		if rw, ok := repo["RW"].(interface{}); ok {
			re.ReadWrite = rw
		}
		if rwp, ok := repo["RW+"].(interface{}); ok {
			re.ReadWritePlus = rwp
		}
		if d, ok := repo["-"].(interface{}); ok {
			re.Deny = d
		}
		repos[key.(string)] = re
	}
	admin_repos := map[string]Repo{}
	for key, value := range o["adminRepos"].(map[interface{}]interface{}) {
		repo := value.(map[interface{}]interface{})
		re := Repo{}
		if r, ok := repo["R"].(interface{}); ok {
			re.Read = r
		}
		if rw, ok := repo["RW"].(interface{}); ok {
			re.ReadWrite = rw
		}
		if rwp, ok := repo["RW+"].(interface{}); ok {
			re.ReadWritePlus = rwp
		}
		if d, ok := repo["-"].(interface{}); ok {
			re.Deny = d
		}
		admin_repos[key.(string)] = re
	}
	out := Data{
		Users:      users,
		Groups:     groups,
		Repos:      repos,
		AdminRepos: admin_repos,
	}
	return out, err
}

func generateConfHeader(title string) string {
	return fmt.Sprintf(`#
# %s
#

`, title)
}

func convArrayableValue(v interface{}, c func(string) string) (out string) {
	if value, ok := v.([]string); ok {
		for _, x := range value {
			out += c(x)
		}
	} else if value, ok := v.([]interface{}); ok {
		for _, x := range value {
			if v, ok := x.(string); ok {
				out += c(v)
			} else if a, ok := x.([]string); ok {
				for _, v := range a {
					out += c(v)
				}
			} else if m, ok := x.(map[interface{}]interface{}); ok {
				for k, _ := range m {
					out += c(k.(string))
				}
			}
		}
	} else if value, ok := v.(string); ok {
		out += c(value)
	}
	return out
}

func generateUser(id string, user User) string {
	out := fmt.Sprintf("@%s", id)
	if len(out) < 21 {
		out += strings.Repeat(" ", 21-len(out))
	}
	out += "= "
	out += convArrayableValue(user.Keys, func(v string) string {
		return v + " "
	})
	out = strings.TrimRight(out, " ")
	return out
}

func generateUsers(in map[string]User) (out string) {
	keys := []string{}
	for k, _ := range in {
		keys = append(keys, k)
	}
	sort.Sort(sort.StringSlice(keys))

	for i := 0; i < len(keys); i += 1 {
		out += generateUser(keys[i], in[keys[i]]) + "\n"
	}
	return out + "\n"
}

func generateGroup(id string, group interface{}) (out string) {
	out = "@" + id
	for i := len(id); i < 20; i += 1 {
		out += " "
	}
	out += "= "
	out += convArrayableValue(group, func(u string) string {
		return "@" + u + " "
	})
	out = strings.TrimRight(out, " ")
	return out
}

func generateGroups(in map[string]interface{}) (out string) {
	keys := []string{}
	for k, _ := range in {
		keys = append(keys, k)
	}
	sort.Sort(sort.StringSlice(keys))

	for i := 0; i < len(keys); i += 1 {
		out += generateGroup(keys[i], in[keys[i]]) + "\n"
	}
	return out + "\n"
}

func generateRepos(in map[string]Repo) (out string) {
	keys := []string{}
	for k, _ := range in {
		keys = append(keys, k)
	}
	sort.Sort(sort.StringSlice(keys))

	for i := 0; i < len(keys); i += 1 {
		out += generateRepo(keys[i], in[keys[i]]) + "\n"
	}
	return out
}

func generateRepoGroups(perm string, groups interface{}) (out string) {
	out = "\t" + perm + "\t= "
	out += convArrayableValue(groups, func(g string) string {
		return "@" + g + " "
	})
	out = strings.TrimRight(out, " ")
	return out + "\n"
}
func generateRepo(path string, repo Repo) (out string) {
	out = "repo\t" + path + "\n"
	if repo.Read != nil {
		out += generateRepoGroups("R", repo.Read)
	}
	if repo.ReadWrite != nil {
		out += generateRepoGroups("RW", repo.ReadWrite)
	}
	if repo.ReadWritePlus != nil {
		out += generateRepoGroups("RW+", repo.ReadWritePlus)
	}
	if repo.Deny != nil {
		out += generateRepoGroups("-", repo.Deny)
	}
	return out
}

func GenerateConf(in Data) (out string) {
	out = generateConfHeader("Users")
	out += generateUsers(in.Users)
	out += generateConfHeader("Groups")
	out += generateGroups(in.Groups)
	out += generateConfHeader("Repos")
	out += generateRepos(in.Repos)
	out += generateConfHeader("Admin Repos")
	out += generateRepos(in.AdminRepos)
	return strings.TrimRight(out, "\n") + "\n"
}

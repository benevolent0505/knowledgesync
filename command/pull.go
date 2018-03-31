package command

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

type Knowledge struct {
	ID             int      `json:"knowledgeId"`
	Title          string   `json:"title"`
	Content        string   `json:"content"`
	CommentCount   int      `json:"commentCount"`
	InsertUser     int      `json:"insertUser"`
	InsertDatetime string   `json:"insertDatetime"`
	LikeCount      int      `json:"likeCount"`
	PublicFlag     int      `json:"publicFlag"`
	Tags           []string `json:"tags"`
	Template       string   `json:"template"`
	UpdateDatetime string   `json:"updateDatetime"`
	UpdateUser     int      `json:"updateUser"`
	Viewers        struct {
		Groups []group `json:"groups"`
		Users  []user  `json:"users"`
	} `json:"viewers"`
}

type group struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

type user struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

type UserDetail struct {
	UserId         int    `json:"userId"`
	Username       string `json:"userName"`
	InsertUser     int    `json:"insertUser"`
	InsertDatetime string `json:"insertDatetime"`
	UpdateUser     int    `json:"updateUser"`
}

type PullCommand struct {
	Meta
	Config
}

func (c *PullCommand) Run(args []string) int {
	// user一覧を取得
	userIdAndNameMap := make(map[int]string)
	for offset := 0; len(userIdAndNameMap) >= offset; offset += 10 {
		reqUrl := fmt.Sprintf("http://%s/api/users?offset=%d", c.Config.Host, offset)
		req, err := http.NewRequest("GET", reqUrl, nil)
		if err != nil {
			fmt.Println("new request error", err)
			return 1
		}
		req.Header.Set("PRIVATE-TOKEN", c.Config.PrivateToken)
		client := &http.Client{}
		res, err := client.Do(req)
		if err != nil {
			fmt.Println("response error: ", err)
			return 1
		}
		defer res.Body.Close()

		var userDetails []UserDetail
		err = json.NewDecoder(res.Body).Decode(&userDetails)
		if err != nil {
			fmt.Println("decode error: ", err)
			return 1
		}
		for _, u := range userDetails {
			userIdAndNameMap[u.UserId] = u.Username
		}
	}

	type knowledge struct {
		Id      int
		User    string
		Content string
	}
	var knowledgeList []knowledge = make([]knowledge, 0)
	for offset := 0; len(knowledgeList) >= offset; offset += 10 {
		reqUrl := fmt.Sprintf("http://%s/api/knowledges?offset=%d", c.Config.Host, offset)
		req, err := http.NewRequest("GET", reqUrl, nil)
		if err != nil {
			fmt.Println("new request error", err)
			return 1
		}
		req.Header.Set("PRIVATE-TOKEN", c.Config.PrivateToken)
		client := &http.Client{}
		res, err := client.Do(req)
		if err != nil {
			fmt.Println("response error: ", err)
			return 1
		}
		defer res.Body.Close()

		var knowkedges []Knowledge
		err = json.NewDecoder(res.Body).Decode(&knowkedges)
		if err != nil {
			fmt.Println("decode error: ", err)
			return 1
		}

		for _, k := range knowkedges {
			knowledgeList = append(knowledgeList, knowledge{
				Id:      k.ID,
				User:    userIdAndNameMap[k.InsertUser],
				Content: k.Content,
			})
		}
	}

	for _, name := range userIdAndNameMap {
		dir := filepath.Join(c.Config.LocalRoot, c.Config.Host, name)
		err := os.MkdirAll(dir, 0755)
		if err != nil {
			fmt.Println("mkdir error: ", err)
			return 1
		}
	}

	baseDir := filepath.Join(c.Config.LocalRoot, c.Config.Host)
	for _, k := range knowledgeList {
		f, err := os.Create(filepath.Join(baseDir, k.User, fmt.Sprint(k.Id, ".md")))
		if err != nil {
			fmt.Println("mkfile error: ", err)
			return 1
		}
		_, err = f.WriteString(k.Content)
		if err != nil {
			fmt.Println("write content error: ", err)
			return 1
		}
		err = f.Close()
		if err != nil {
			fmt.Println("file close error: ", err)
			return 1
		}
	}

	return 0
}

func (c *PullCommand) Synopsis() string {
	return c.Config.Host
}

func (c *PullCommand) Help() string {
	helpText := `

`
	return strings.TrimSpace(helpText)
}

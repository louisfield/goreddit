package main

import (
	"github.com/vartanbeno/go-reddit/v2/reddit"
	"encoding/json"
	"fmt"
	"os"
	"context"
	"time"
	"io/ioutil"
)

type Config struct {
    Username  string `json:"username"`
    Password string `json:"password"`
	Secret string `json:"secret"`
	Id string `json:"id"`
}

var ctx = context.Background()

func LoadConfiguration(file string) Config {
    var config Config
    configFile, err := os.Open(file)
	byteValue, _ := ioutil.ReadAll(configFile)
	var result map[string]interface{}
    json.Unmarshal([]byte(byteValue), &result)

    fmt.Println(result)
	json.Unmarshal(byteValue, &config)
    if err != nil {
        fmt.Println(err.Error())
    }
	defer configFile.Close()
    return config
}

func getAllComments(client *reddit.Client , comments []*reddit.Comment, comments_list []string) []string {
	if len(comments) == 0 {
		return comments_list
	}
	for _, comment := range comments {
		comments_list = append(comments_list, comment.Body)
		client.Comment.LoadMoreReplies(ctx, comment)
		time.Sleep(15 * time.Millisecond)
		comments_list = getAllComments(client, comment.Replies.Comments, comments_list)
	}
	return comments_list
}

func main() {
	configuration := LoadConfiguration("./config/conf.json")
	credentials := reddit.Credentials{ID: configuration.Id, Secret: configuration.Secret, Username: configuration.Username, Password: configuration.Password}
    client, _ := reddit.NewClient(credentials)
	posts, _, err := client.Subreddit.TopPosts(ctx, "leagueoflegends", &reddit.ListPostOptions{
		ListOptions: reddit.ListOptions{
			Limit: 10,
		},
		Time: "today",
	})
	if err != nil {
		panic(err)
	}
	comments := []string{}
	for _, post := range posts {
		fmt.Println(post.ID)
		fmt.Println(post.Title)
		thread, _, err := client.Post.Get(ctx, post.ID)
		if err != nil {
			panic(err)
		}
		comments = getAllComments(client, thread.Comments, comments)

	}
	fmt.Println(len(comments))
}

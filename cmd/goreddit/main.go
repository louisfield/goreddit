package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/vartanbeno/go-reddit/v2/reddit"
)

type Config struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Secret   string `json:"secret"`
	Id       string `json:"id"`
}

var leagueSlang = []string{"tyler1", "faker", "guma", "lec", "lcs", "lck", "lpl", "msi", "perkz", "rekkles", "caps", "flakked", "brokenblade", "bb", "jankos", "taragmas", "oner", "zeus", "doublelift", "ls", "keria", "gumaushyi", "caster", "baron", "drake", "hextech", "mountain", "infernal", "aatrox", "atrox", "ahri", "akali", "alistar", "ali", "amumu", "anivia", "annie", "ashe", "azir", "blitzcrank", "blitz", "brand", "braum", "caitlyn", "cait", "cassiopeia", "cassio", "chogath", "cho", "corki", "darius", "diana", "dr. mundo", "dr mundo", "dr mundo", "mundo", "draven", "elise", "evelynn", "ezreal", "fiddlesticks", "fiddle", "fiora", "fizz", "galio", "gangplank", "garen", "gnar", "gragas", "graves", "hecarim", "heca", "heimerdinger", "heim", "irelia", "janna", "jarvan iv", "j4", "jax", "jayce", "jinx", "kalista", "karma", "karthus", "kassadin", "kassa", "katarina", "kata", "kayle", "kennen", "kha", "khazix", "kogmaw", "leblanc", "lee sin", "lee", "leona", "lissandra", "liss", "lucian", "lulu", "lux", "malphite", "malph", "malzahar", "malz", "maokai", "master yi", "miss fortune", "mf", "mordekaiser", "morde", "morgana", "morg", "nami", "nasus", "nautilus", "naut", "nidalee", "nid", "nida", "nocturne", "noc", "nunu", "olaf", "orianna", "ori", "pantheon", "poppy", "quinn", "rammus", "reksai", "renekton", "rengar", "riven", "rumble", "ryze", "sejuani", "sej", "shaco", "shen", "shyvana", "shyv", "singed", "sion", "sivir", "skarner", "sona", "soraka", "swain", "syndra", "talon", "taric", "teemo", "thresh", "tristana", "trist", "trundle", "tryndamere", "tryn", "twisted fate", "twitch", "udyr", "urgot", "varus", "vayne", "veigar", "velkoz", "vi", "viktor", "vladimir", "vlad", "volibear", "voli", "warwick", "wukong", "xerath", "xin zhao", "xin", "yasuo", "yas", "yorick", "zac", "zed", "ziggs", "zilean", "zyra", "bard", "ekko", "tahm kench", "tahm", "unbench the kench", "inting", "int", "kindred", "illaoi", "jhin", "aurelion sol", "sol", "taliyah", "kled", "ivern", "camille", "rakan", "xayah", "kayn", "ornn", "zoe", "kaisa", "pyke", "neeko", "sylas", "yuumi", "qiyana", "senna", "aphelios", "sett", "lillia", "yone", "samira", "seraphine", "rell", "viego", "gwen", "akshan", "vex", "t1", "skt", "tl", "c9", "fnc", "g2", "tsm", "tower", "aa", "ace", "ad", "carry", "ad", "afk", "aggro", "babysit", "base", "bf", "biscuit", "bench", "bork", "bot", "cass", "cs", "feeding", "gank"}

type CommentsList struct {
	Comments []string
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

func getAllComments(client *reddit.Client, comments []*reddit.Comment, comments_list []string) []string {
	if len(comments) == 0 {
		return comments_list
	}
	for _, comment := range comments {
		reg, err := regexp.Compile(`[^\w\s]`)
		if err != nil {
			log.Fatal(err)
		}
		processedString := strings.ToLower(comment.Body)
		for _, word := range leagueSlang {

			if strings.Contains(processedString, word) {
				processedString = strings.Replace(processedString, word, "", -1)
			}
		}

		processedString = reg.ReplaceAllString(processedString, "")
		var b strings.Builder
		for _, word := range strings.Split(processedString, " ") {
			if len(word) > 3 {
				b.WriteString(word)
				b.WriteString(" ")
			}
		}
		comments_list = append(comments_list, b.String())
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

	f, err := os.Create("data.txt")

	if err != nil {
		log.Fatal(err)
	}

	defer f.Close()
	commentsToParse := CommentsList{comments}

	file, _ := json.MarshalIndent(commentsToParse, "", " ")
	_ = ioutil.WriteFile("output.json", file, 0644)

}

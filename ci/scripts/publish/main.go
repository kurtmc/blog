package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

func main() {
	client := NewClient("https://dev.to", os.Getenv("DEV_TO_API_KEY"))

	articles, err := client.GetArticles()
	if err != nil {
		panic(err)
	}

	articleDirectories := []string{
		"2022-07/lightweight-artifact-repository-with-python-and-github",
		"2024-06/unifi-autobackup-data-recovery-and-restore",
	}

	for _, articleDirectory := range articleDirectories {
		title, err := getTitle(articleDirectory)
		if err != nil {
			panic(err)
		}
		if title == "" {
			panic(fmt.Errorf("title cannot be blank"))
		}
		markdownBody, err := getMarkdownBody(articleDirectory)
		if err != nil {
			panic(err)
		}

		articleID := 0
		for _, article := range articles {
			if article.Title == title {
				articleID = article.ID
			}
		}

		if articleID == 0 {
			a, err := client.CreateNewArticle(&ArticleRequest{
				Article: struct {
					Title        string   `json:"title"`
					BodyMarkdown string   `json:"body_markdown"`
					Published    bool     `json:"published,omitempty"`
					Tags         []string `json:"tags,omitempty"`
				}{
					Title:        title,
					BodyMarkdown: markdownBody,
				},
			})
			if err != nil {
				panic(err)
			}
			fmt.Println(a)
		} else {
			a, err := client.UpdateArticle(articleID, &ArticleRequest{
				Article: struct {
					Title        string   `json:"title"`
					BodyMarkdown string   `json:"body_markdown"`
					Published    bool     `json:"published,omitempty"`
					Tags         []string `json:"tags,omitempty"`
				}{
					Title:        title,
					BodyMarkdown: markdownBody,
				},
			})
			if err != nil {
				panic(err)
			}
			fmt.Println(a)
		}
	}
}

func getTitle(path string) (string, error) {
	file, err := os.Open("../../../" + path + "/post.md")
	if err != nil {
		return "", err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		t := scanner.Text()
		if strings.HasPrefix(t, "title: ") {
			return strings.TrimPrefix(t, "title: "), nil
		}
	}

	if err := scanner.Err(); err != nil {
		return "", err
	}

	return "", fmt.Errorf("could not find title")
}
func getMarkdownBody(path string) (string, error) {
	b, err := ioutil.ReadFile("../../../" + path + "/post.md")
	return string(b), err
}

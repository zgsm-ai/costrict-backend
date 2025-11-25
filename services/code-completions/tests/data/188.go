package main

import (
	"fmt"
	"time"
)

type User struct {
	ID       int       `json:"id"`
	Username string    `json:"username"`
	Email    string    `json:"email"`
	JoinDate time.Time `json:"join_date"`
}

type Post struct {
	ID        int       `json:"id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	Author    *User     `json:"author"`
	CreatedAt time.Time `json:"created_at"`
	Tags      []string  `json:"tags"`
	Published bool      `json:"published"`
}

type Comment struct {
	ID        int       `json:"id"`
	Content   string    `json:"content"`
	Author    *User     `json:"author"`
	PostID    int       `json:"post_id"`
	CreatedAt time.Time `json:"created_at"`
}

type Blog struct {
	Users    map[int]*User    `json:"users"`
	Posts    map[int]*Post    `json:"posts"`
	Comments map[int]*Comment `json:"comments"`
	nextID   int
}

func NewBlog() *Blog {
	return &Blog{
		Users:    make(map[int]*User),
		Posts:    make(map[int]*Post),
		Comments: make(map[int]*Comment),
		nextID:   1,
	}
}

func (b *Blog) getNextID() int {
	id := b.nextID
	b.nextID++
	return id
}

func (b *Blog) AddUser(username, email string) *User {
	user := &User{
		ID:       b.getNextID(),
		Username: username,
		Email:    email,
		JoinDate: time.Now(),
	}
	b.Users[user.ID] = user
	return user
}

func (b *Blog) AddPost(title, content string, author *User, tags []string) *Post {
	post := &Post{
		ID:        b.getNextID(),
		Title:     title,
		Content:   content,
		Author:    author,
		CreatedAt: time.Now(),
		Tags:      tags,
		Published: false,
	}
	b.Posts[post.ID] = post
	return post
}

func (b *Blog) PublishPost(id int) bool {
	if post, ok := b.Posts[id]; ok {
		post.Published = true
		return true
	}
	return false
}

func (b *Blog) AddComment(content string, author *User, postID int) *Comment {
	comment := &Comment{
		ID:        b.getNextID(),
		Content:   content,
		Author:    author,
		PostID:    postID,
		CreatedAt: time.Now(),
	}
	b.Comments[comment.ID] = comment
	return comment
}

func (b *Blog) GetCommentsByPost(postID int) []*Comment {
	var comments []*Comment
	for _, comment := range b.Comments {
		if comment.PostID == postID {
			comments = append(comments, comment)
		}
	}
	return comments
}

func (b *Blog) GetPublishedPosts() []*Post {
	var posts []*Post
	for _, post := range b.Posts {
		if post.Published {
			posts = append(posts, post)
		}
	}
	return posts
}

func (b *Blog) GetPostsByAuthor(authorID int) []*Post {
	var posts []*Post
	for _, post := range b.Posts {
		if post.Author.ID == authorID {
			posts = append(posts, post)
		}
	}
	return posts
}

func DisplayUser(user *User) {
	fmt.Printf("ID: %d\n", user.ID)
	fmt.Printf("Username: %s\n", user.Username)
	fmt.Printf("Email: %s\n", user.Email)
}

func DisplayPost(post *Post, showComments bool, blog *Blog) {
	fmt.Printf("ID: %d\n", post.ID)
	fmt.Printf("Title: %s\n", post.Title)
	fmt.Printf("Author: %s\n", post.Author.Username)
	fmt.Printf("Created: %s\n", post.CreatedAt.Format("2006-01-02 15:04:05"))
	fmt.Printf("Published: %t\n", post.Published)

	if len(post.Tags) > 0 {
		fmt.Printf("Tags: ")
		for i, tag := range post.Tags {
			if i > 0 {
				fmt.Printf(", ")
			}
			fmt.Printf("#%s", tag)
		}
		fmt.Printf("\n")
	}

	fmt.Printf("Content:\n%s\n", post.Content)

	if showComments {
		comments := blog.GetCommentsByPost(post.ID)
		if len(comments) > 0 {
			fmt.Printf("\nComments (%d):\n", len(comments))
			for _, comment := range comments {
				fmt.Printf("%s (%d - %s):\n", comment.Author.Username,
					comment.ID, comment.CreatedAt.Format("2006-01-02 15:04:05"))
				fmt.Printf("%s\n", comment.Content)
			}
		}
	}
}

func GenerateBlog() *Blog {
	blog := NewBlog()

	john := blog.AddUser("johnsmith", "john@example.com")
	jane := blog.AddUser("janedoe", "jane@example.com")

	post1 := blog.AddPost("Getting Started with Go",
		"Go is a powerful programming language developed by Google. It's known for its simplicity and efficiency.",
		john, []string{"programming", "golang"})
	blog.PublishPost(post1.ID)

	post2 := blog.AddPost("Design Principles for Modern UI",
		"Creating a modern UI requires understanding of fundamental design principles.",
		jane, []string{"design", "ui"})
	blog.PublishPost(post2.ID)

	blog.AddComment("Great introduction to Go!", jane, post1.ID)
	blog.AddComment("This really helped me understand Go's concurrency model.", john, post1.ID)

	return blog
}

func main() {
	blog := GenerateBlog()

	fmt.Printf("=== Blog Statistics ===\n")
	fmt.Printf("Total Users: %d\n", len(blog.Users))
	fmt.Printf("Total Posts: %d\n", len(blog.Posts))

	fmt.Printf("\n=== All Published Posts ===\n")
	publishedPosts := blog.GetPublishedPosts()
	for _, post := range publishedPosts {
		fmt.Printf("- %s by %s\n", post.Title, post.Author.Username)
	}

	fmt.Printf("\n=== Detailed Post ===\n")
	post := blog.GetPublishedPosts()[0]
	DisplayPost(post, true, blog)
}

package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"sort"
	"strings"

	"github.com/google/go-github/v60/github"
	"golang.org/x/oauth2"
)

const (
	votingCommentMarker = "<!-- hackathon-voting-comment -->"
	votingEmoji         = "+1"
)

// issueVote represents an issue with its vote count
type issueVote struct {
	issue *github.Issue
	votes int
}

func main() {
	var (
		owner          string
		repo           string
		discussionID   string
		command        string
		hackathonLabel string
	)

	flag.StringVar(&owner, "owner", "", "GitHub repository owner")
	flag.StringVar(&repo, "repo", "", "GitHub repository name")
	flag.StringVar(&discussionID, "discussion-id", "", "Discussion ID for voting results (only for close-voting)")
	flag.StringVar(&command, "command", "", "Command to run: start-voting or close-voting")
	flag.StringVar(&hackathonLabel, "label", "", "Label to filter issues (e.g., 'hackathon-2025-autumn')")
	flag.Parse()

	if command == "" {
		log.Fatal("command is required: start-voting or close-voting")
	}
	if owner == "" || repo == "" {
		log.Fatal("owner and repo are required")
	}

	token := os.Getenv("GITHUB_TOKEN")
	if token == "" {
		log.Fatal("GITHUB_TOKEN environment variable is required")
	}

	ctx := context.Background()
	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token})
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)

	switch command {
	case "start-voting":
		if err := startVoting(ctx, client, owner, repo, hackathonLabel); err != nil {
			log.Fatalf("Error starting voting: %v", err)
		}
	case "close-voting":
		if discussionID == "" {
			log.Fatal("discussion-id is required for close-voting")
		}
		if err := closeVoting(ctx, client, owner, repo, hackathonLabel, discussionID); err != nil {
			log.Fatalf("Error closing voting: %v", err)
		}
	default:
		log.Fatalf("Unknown command: %s", command)
	}
}

// startVoting adds voting comments to all open issues with the specified label
// and creates a discussion listing all topics
func startVoting(ctx context.Context, client *github.Client, owner, repo, label string) error {
	issues, err := getOpenIssues(ctx, client, owner, repo, label)
	if err != nil {
		return fmt.Errorf("failed to get open issues: %w", err)
	}

	if len(issues) == 0 {
		log.Println("No open issues found with the specified label")
		return nil
	}

	log.Printf("Found %d open issues", len(issues))

	// Add voting comments to each issue
	for _, issue := range issues {
		if err := addVotingComment(ctx, client, owner, repo, issue); err != nil {
			log.Printf("Warning: failed to add voting comment to issue #%d: %v", *issue.Number, err)
			continue
		}
		log.Printf("Added voting comment to issue #%d: %s", *issue.Number, *issue.Title)
	}

	// Create discussion body with list of all topics
	discussionBody := createTopicsListMarkdown(owner, repo, issues)

	fmt.Println("\n=== DISCUSSION CONTENT ===")
	fmt.Println("Create a new discussion with the following content:")
	fmt.Println("Title: Hackathon Topic Voting")
	fmt.Println("\nBody:")
	fmt.Println(discussionBody)
	fmt.Println("=== END DISCUSSION CONTENT ===")

	return nil
}

// closeVoting counts votes on all issues and updates the discussion with ranked results
// Note: discussionID is validated but the actual discussion update happens via GitHub Actions
func closeVoting(ctx context.Context, client *github.Client, owner, repo, label, _ string) error {
	issues, err := getOpenIssues(ctx, client, owner, repo, label)
	if err != nil {
		return fmt.Errorf("failed to get open issues: %w", err)
	}

	if len(issues) == 0 {
		log.Println("No open issues found with the specified label")
		return nil
	}

	log.Printf("Found %d open issues", len(issues))

	// Count votes for each issue
	var issueVotes []issueVote
	for _, issue := range issues {
		votes, err := countVotes(ctx, client, owner, repo, issue)
		if err != nil {
			log.Printf("Warning: failed to count votes for issue #%d: %v", *issue.Number, err)
			votes = 0
		}
		issueVotes = append(issueVotes, issueVote{issue: issue, votes: votes})
		log.Printf("Issue #%d: %s - %d votes", *issue.Number, *issue.Title, votes)
	}

	// Sort by votes (descending)
	sort.Slice(issueVotes, func(i, j int) bool {
		return issueVotes[i].votes > issueVotes[j].votes
	})

	// Create ranked list markdown
	rankedList := createRankedListMarkdown(owner, repo, issueVotes)

	fmt.Println("\n=== VOTING RESULTS ===")
	fmt.Println("Update the discussion with the following content:")
	fmt.Println("\nRanked Topics:")
	fmt.Println(rankedList)
	fmt.Println("=== END VOTING RESULTS ===")

	return nil
}

// getOpenIssues retrieves all open issues with the specified label
func getOpenIssues(ctx context.Context, client *github.Client, owner, repo, label string) ([]*github.Issue, error) {
	opts := &github.IssueListByRepoOptions{
		State: "open",
		ListOptions: github.ListOptions{
			PerPage: 100,
		},
	}

	if label != "" {
		opts.Labels = []string{label}
	}

	var allIssues []*github.Issue
	for {
		issues, resp, err := client.Issues.ListByRepo(ctx, owner, repo, opts)
		if err != nil {
			return nil, err
		}

		// Filter out pull requests (GitHub API returns PRs as issues too)
		for _, issue := range issues {
			if issue.PullRequestLinks == nil {
				allIssues = append(allIssues, issue)
			}
		}

		if resp.NextPage == 0 {
			break
		}
		opts.Page = resp.NextPage
	}

	return allIssues, nil
}

// addVotingComment adds a voting comment to an issue if one doesn't already exist
func addVotingComment(ctx context.Context, client *github.Client, owner, repo string, issue *github.Issue) error {
	// Check if voting comment already exists
	comments, _, err := client.Issues.ListComments(ctx, owner, repo, *issue.Number, &github.IssueListCommentsOptions{
		ListOptions: github.ListOptions{PerPage: 100},
	})
	if err != nil {
		return fmt.Errorf("failed to list comments: %w", err)
	}

	for _, comment := range comments {
		if comment.Body != nil && strings.Contains(*comment.Body, votingCommentMarker) {
			log.Printf("Voting comment already exists on issue #%d", *issue.Number)
			return nil
		}
	}

	// Create voting comment
	commentBody := fmt.Sprintf(`%s

## 🗳️ Vote for this Topic!

If you're interested in working on this topic during the hackathon, please react to this comment with a 👍 (thumbs up) emoji.

**How to vote:**
1. Click the smiley face icon below this comment
2. Select the 👍 reaction

The topics with the most votes will help us prioritize and form teams!

---
*This is an automated voting comment. Votes are counted based on 👍 reactions to this comment.*
`, votingCommentMarker)

	comment, _, err := client.Issues.CreateComment(ctx, owner, repo, *issue.Number, &github.IssueComment{
		Body: github.String(commentBody),
	})
	if err != nil {
		return fmt.Errorf("failed to create comment: %w", err)
	}

	// Note: GitHub API doesn't support pinning comments directly
	// The comment will need to be pinned manually or via GraphQL API
	log.Printf("Created voting comment (ID: %d) - Note: Pin this comment manually", *comment.ID)

	return nil
}

// countVotes counts the thumbs up reactions on the voting comment of an issue
func countVotes(ctx context.Context, client *github.Client, owner, repo string, issue *github.Issue) (int, error) {
	comments, _, err := client.Issues.ListComments(ctx, owner, repo, *issue.Number, &github.IssueListCommentsOptions{
		ListOptions: github.ListOptions{PerPage: 100},
	})
	if err != nil {
		return 0, fmt.Errorf("failed to list comments: %w", err)
	}

	for _, comment := range comments {
		if comment.Body != nil && strings.Contains(*comment.Body, votingCommentMarker) {
			// Get reactions on this comment
			reactions, _, err := client.Reactions.ListIssueCommentReactions(ctx, owner, repo, *comment.ID, &github.ListOptions{
				PerPage: 100,
			})
			if err != nil {
				return 0, fmt.Errorf("failed to list reactions: %w", err)
			}

			// Count thumbs up reactions
			count := 0
			for _, reaction := range reactions {
				if reaction.Content != nil && *reaction.Content == votingEmoji {
					count++
				}
			}
			return count, nil
		}
	}

	return 0, nil
}

// createTopicsListMarkdown creates a markdown list of all topics linking to issues
func createTopicsListMarkdown(owner, repo string, issues []*github.Issue) string {
	var sb strings.Builder
	sb.WriteString("# Hackathon Topic Voting\n\n")
	sb.WriteString("Vote for the topics you're interested in! Go to each issue and react with 👍 on the voting comment.\n\n")
	sb.WriteString("## Topics\n\n")

	for _, issue := range issues {
		sb.WriteString(fmt.Sprintf("- [#%d - %s](https://github.com/%s/%s/issues/%d)\n",
			*issue.Number, *issue.Title, owner, repo, *issue.Number))
	}

	sb.WriteString("\n---\n")
	sb.WriteString("*Voting results will be posted here once voting closes.*\n")

	return sb.String()
}

// createRankedListMarkdown creates a ranked markdown list of topics with vote counts
func createRankedListMarkdown(owner, repo string, issueVotes []issueVote) string {
	var sb strings.Builder
	sb.WriteString("# 🏆 Hackathon Topic Voting Results\n\n")
	sb.WriteString("## Ranked Topics\n\n")
	sb.WriteString("| Rank | Topic | Votes |\n")
	sb.WriteString("|:----:|:------|:-----:|\n")

	for i, iv := range issueVotes {
		medal := ""
		switch i {
		case 0:
			medal = "🥇 "
		case 1:
			medal = "🥈 "
		case 2:
			medal = "🥉 "
		}
		sb.WriteString(fmt.Sprintf("| %s%d | [#%d - %s](https://github.com/%s/%s/issues/%d) | %d |\n",
			medal, i+1, *iv.issue.Number, *iv.issue.Title, owner, repo, *iv.issue.Number, iv.votes))
	}

	sb.WriteString("\n---\n")
	sb.WriteString("*Thank you for voting! Teams will be formed based on these results and participant preferences.*\n")

	return sb.String()
}

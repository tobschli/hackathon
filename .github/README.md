# Hackathon Voting System

GitHub Actions workflows for managing hackathon topic voting.

## Workflows

### Start Voting (`start-voting.yml`)

Initiates the voting phase:
1. Creates a GitHub Discussion for voting
2. Adds a comment for each open issue in the discussion
3. Each comment includes the issue title, description preview, and voting instructions

### Close Voting (`close-voting.yml`)

Closes voting and tallies results:
1. Automatically finds the active voting discussion
2. Counts upvotes on each topic comment
3. Collects 🚀 reactions to determine team assignments
4. Updates the discussion body with ranked results
5. Assigns users who reacted with 🚀 to their respective issues

## How to Vote

1. Go to the voting discussion
2. Find the comment for the topic you're interested in
3. **Upvote** the comment to vote for the topic
4. **React with 🚀** if you want to work on it (you'll be assigned when voting closes)

## Setup

### Prerequisites

- GitHub Discussions must be enabled for the repository
- A Personal Access Token (PAT) with discussion permissions

### Setting Up the Discussion Token

1. **Create a Fine-grained PAT:**
   - Go to https://github.com/settings/tokens?type=beta
   - Click **"Generate new token"**
   - **Token name:** `hackathon-voting`
   - **Expiration:** Set as needed
   - **Repository access:** Select your hackathon repo
   - **Permissions → Repository permissions:**
     - **Discussions:** Read and write
     - **Issues:** Read and write
   - Click **"Generate token"** and copy it

2. **Add the token as a repository secret:**
   - Go to your repository → **Settings** → **Secrets and variables** → **Actions**
   - Click **"New repository secret"**
   - **Name:** `DISCUSSION_TOKEN`
   - **Secret:** Paste the token

## Usage

### Starting Voting

1. Go to **Actions** → **"Start Hackathon Voting"**
2. Click **"Run workflow"**
3. (Optional) Change the discussion category
4. Done! A voting discussion is created with all open issues as topics.

### Closing Voting

1. Go to **Actions** → **"Close Hackathon Voting"**
2. Click **"Run workflow"**
3. Done! Results are posted and team members are assigned.

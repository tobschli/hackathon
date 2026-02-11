# Hackathon Voting System

This repository includes GitHub Actions workflows for managing hackathon topic voting.

## Workflows

### Start Voting (`start-voting.yml`)

Initiates the voting phase:
1. Creates a GitHub Discussion for voting
2. Adds a comment for each topic issue in the discussion
3. Each comment includes the issue title, description preview, and voting instructions

### Close Voting (`close-voting.yml`)

Closes voting and tallies results:
1. Counts 👍 reactions on each topic comment (votes)
2. Collects 🚀 reactions to determine team assignments
3. Updates the discussion body with ranked results
4. Assigns users who reacted with 🚀 to their respective issues

## How to Vote

1. Go to the voting discussion
2. Find the comment for the topic you're interested in
3. React with:
   - 👍 = I want this topic to happen (vote)
   - 🚀 = I want to work on this topic (you'll be assigned when voting closes)

## Usage

### Prerequisites

- GitHub Discussions must be enabled for the repository
- A Personal Access Token (PAT) with discussion permissions

### Setting Up the Discussion Token (Required)

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

### Starting Voting

1. Go to **Actions** → **"Start Hackathon Voting"**
2. Click **"Run workflow"**
3. Enter the hackathon label (e.g., `hackathon-2025-autumn`)
4. Optionally specify a discussion category (default: `General`)

### Closing Voting

1. Go to **Actions** → **"Close Hackathon Voting"**
2. Click **"Run workflow"**
3. Enter the hackathon label
4. Enter the discussion number (from the URL)

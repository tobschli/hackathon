# Hackathon Voting System

This repository includes GitHub Actions workflows for managing hackathon topic voting.

## Workflows

### Start Voting (`start-voting.yml`)

Initiates the voting phase:
1. Adds a voting comment to each open issue with the specified label
2. Pins the voting comment to each issue
3. Creates a GitHub Discussion listing all topics

### Close Voting (`close-voting.yml`)

Closes voting and tallies results:
1. Counts 👍 reactions on each voting comment
2. Updates the discussion body with ranked results
3. Removes voting comments from all issues

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

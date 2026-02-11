# Hackathon Voting System

This directory contains tooling for managing hackathon topic voting via GitHub Actions.

## Overview

The voting system consists of two GitHub Actions workflows:

1. **Start Voting** (`start-voting.yml`) - Initiates the voting phase
2. **Close Voting** (`close-voting.yml`) - Closes voting and tallies results

## How It Works

### Topic Collection Phase

1. Participants create issues using the "Hackathon Topic Proposal" template
2. Issues should be labeled with the hackathon identifier (e.g., `hackathon-2025-autumn`)

### Voting Phase

When topic collection ends, an organizer triggers the "Start Hackathon Voting" workflow:

1. The workflow adds a voting comment to each open issue with the specified label
2. The voting comment is automatically pinned to the issue
3. A GitHub Discussion is created listing all topics with links

Participants vote by:
- Going to each topic issue
- Adding a 👍 reaction to the pinned voting comment

### Results Phase

When voting ends, an organizer triggers the "Close Hackathon Voting" workflow:

1. The workflow counts 👍 reactions on each voting comment
2. Topics are ranked by vote count
3. The discussion body is updated with the ranked results

## Usage

### Prerequisites

- GitHub Discussions must be enabled for the repository
- Ensure the repository has an appropriate discussion category (default: "General")
- A Personal Access Token (PAT) with discussion permissions (see setup below)

### Setting Up the Discussion Token (Required)

The default `GITHUB_TOKEN` cannot create/update discussions via GraphQL. You must create a PAT:

1. **Create a Fine-grained PAT:**
   - Go to https://github.com/settings/tokens?type=beta
   - Click **"Generate new token"**
   - **Token name:** `hackathon-voting`
   - **Expiration:** Set as needed (e.g., 90 days)
   - **Repository access:** Select **"Only select repositories"** → choose your hackathon repo
   - **Permissions → Repository permissions:**
     - **Discussions:** Read and write
     - **Issues:** Read and write
   - Click **"Generate token"** and copy it

2. **Add the token as a repository secret:**
   - Go to your repository → **Settings** → **Secrets and variables** → **Actions**
   - Click **"New repository secret"**
   - **Name:** `DISCUSSION_TOKEN`
   - **Secret:** Paste the token
   - Click **"Add secret"**

### Starting Voting

1. Go to Actions → "Start Hackathon Voting"
2. Click "Run workflow"
3. Enter:
   - `hackathon_label`: Label to filter issues (e.g., `hackathon-2025-autumn`)
   - `discussion_category`: Discussion category name (default: `General`)
4. Click "Run workflow"

### Closing Voting

1. Go to Actions → "Close Hackathon Voting"
2. Click "Run workflow"
3. Enter:
   - `hackathon_label`: Same label used when starting voting
   - `discussion_number`: The number of the voting discussion (visible in the URL)
4. Click "Run workflow"

## Technical Details

### Voting Comment Marker

Voting comments contain a hidden HTML marker to identify them:

```html
<!-- hackathon-voting-comment -->
```

This allows the workflow to locate the correct comment when counting votes.

### Permissions Required

The GitHub Actions workflows require a PAT (`DISCUSSION_TOKEN`) with:
- **Discussions:** Read and write
- **Issues:** Read and write

### Troubleshooting

**"Resource not accessible by integration"**
- Ensure you created the `DISCUSSION_TOKEN` secret (see setup above)
- Verify the PAT has `Discussions: Read and write` permission
- Check that the PAT hasn't expired

**"Discussion category not found"**
- Ensure GitHub Discussions is enabled for the repository
- Verify the category name matches exactly (case-insensitive)

## Local Development (Optional)

The Go tool in this directory is provided for local testing. The GitHub Actions workflows handle everything automatically.

### Building

```bash
cd tools/voting
go mod download
go build -o voting-tool .
```

### Running Locally

```bash
export GITHUB_TOKEN="your-github-token"

# Start voting
./voting-tool \
  -command=start-voting \
  -owner=gardener-community \
  -repo=hackathon \
  -label=hackathon-2025-autumn

# Close voting
./voting-tool \
  -command=close-voting \
  -owner=gardener-community \
  -repo=hackathon \
  -label=hackathon-2025-autumn \
  -discussion-id=D_xxxxx
```

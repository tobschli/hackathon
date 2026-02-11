# Gardener Hackathons

The content of this repository has been moved to https://gardener.cloud/community/hackathons/ ([source](https://github.com/gardener/documentation/blob/master/website/community/hackathons/)).

I want to use this repo to organize hackathon votings, and assignments via github actions.
You should implement the following features:
- Pipeline to run when topic collection is over
  - Each open issue should get a "voting comment". People vote for this topic by reacting with a thumbs up emoji.
  - This comment should also be pinned to the issue.
  - Create a discussion with a list of all topics, linking to the issues.
- Pipeline to run when voting is over
  - When the voting is over, count the votes for each topic.
  - on the discussion, create a ranked list of the topic issues.

Instead of implementing this through the github api directly, you should use a popular framework, preferably in go.
# CI/CD in the wild and best practices

Through my work as a consultant I have been exposed to some varying opinions
and implementations of CI/CD practices which I believe has put me in a good
position to talk about some of the benefits and disadvantages of the approaches
I have seen. I would also like to describe an approach that takes into account
the needs of many organisations and I consider best practice.

## Trunk based development

Most organisations strive to run multiple environments for their application
stack. At a minimum organisations will have a "dev" and a "prod" environment.
As a quick shortcut these organisations will decide to map git branches to
these environments and setup their CI system to watch changes to these branches
and deploy those changes directly to the environment they are associated with.
They have it setup like the diagram below, where each environment could be
isolated by AWS account for example.

![](https://github.com/kurtmc/blog/raw/master/drafts/ci-cd-in-the-wild-and-best-practices/images/branch_based_deployements.png)

This is incredibly easy to setup and get going and I think that's why most
organisations are doing this. The pitfall in this approach is that it pushes
all the responsibility of promoting changes through the environments to
operations on the git branches. The way most organisations deal with this is to
set the "dev" branch as the default branch and have changes first merge into
that branch, which automatically triggers the CI system to deploy the "dev"
branch to the "dev" environment. Some manual or automated testing might happen
in this environment and when the developer is satisfied they will create a pull
request from the "dev" branch into the "prod" branch. This works well enough
for low volume changes and if your team has the discipline to perform the
changes to the branches in this prescribed way. What inevitably happens though,
is that a change is merged only into the "prod" branch, and subsequently
someone comes along and tries to merge to "dev" then do the "dev" to "prod"
merge and there are conflicts. Again these small issues come up once in a while
and they are dealt with, with some careful investigation.

A bigger issue arises when the organisation would like to manage more than just
these two environments. There is a need for a new environment, let's call it
"staging" that needs to be closer to the "prod" environment to do some
performance testing for example. In this system, easy, we just create a new
branch! Now when a developer would like to get a change to production they need
to do the following steps:

1. Create a feature branch off "dev"
2. Make the changes and create a pull request to the "dev" branch
3. Create a pull request from the "dev" branch to the "staging" branch
4. Create a pull request from the "staging" branch to the "prod" branch

That's 3 pull requests, imagine trying to get a small bug fix merged into
"prod" in this scenario, depending on your pull request review process this
could take a very long time.


There is a better way, it's called trunk based development and it simplifies
this process but requires you to build enough automated testing so that you are
confident changes that pass the tests are safe to deploy to production.



# draft notes:

```
Structure:
1. Describe real approaches that I have seen, criticise
2. Describe high level ideal world
3. Describe my approach in terms of AWS, Concourse CI, Terraform

Ideas:
1. Using tools in the CI/CD space and calling it CI/CD
2. Duplicating complex logic into many projects, diverging over time, fixing bugs in multiple projects
3. Environments mapped to AWS (Cloud) accounts
4. Automatic promotion if changes
5. Testing before merging
6. Mutating the low impact environments from developer machines
7. Feature flags
8. Advanced features, post deployment integration tests, alarm checking
9. Reducing sources of truth / single source of truth
10. Pipeline per branch, merging branches to promote changes
11. Automate the shit out of everything
12. Separate the infrastructure and code deployment methodology
13. Downloading dependencies on every step run
14. Caching
15. Speed is important
16. Updating many project configurations en masse
17. Execing into the environment running the scripts that build/deploy
18. Hardcoded parameters
19. Avoid creating cycles
20. Deploying multiple environments to the same accounts is a bad idea
21. You should avoid creating differences between the IAM relationships, each environment should contain the entire IAM relationship so that it can be developed before deploying to production
```

# CI/CD in the wild and best practices

Through my work as a consultant I have been exposed to some varying opinions and implementations of CI/CD practices which I believe has put me in a good position to talk about some of the benefits and disadvantages of the approaches I have seen. I would also like to describe an approach that takes into account the needs of many organisations and I consider best practice.

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

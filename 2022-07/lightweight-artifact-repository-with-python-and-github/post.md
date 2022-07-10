---
title: Lightweight artifact repository with Python and GitHub
published: true
tags: python, git, cicd
---
# Lightweight artifact repository with Python and GitHub

Code reuse is fundamental to reducing the cost of software development, reusing a function implementation rather than developing it again is faster. Fixing a bug in a library rather than in multiple re-implementations is easier. A common way to reuse code is to package up related functionality and publish it as a library. Usually you could publish to [Arifactory](https://jfrog.com/artifactory/) or [Nexus](https://www.sonatype.com/products/nexus-repository) but occasionally you may have a business constraint that makes it painful and slow to onboard a new tool, often for valid reasons, they maybe be expensive to onboard and support.

Coming from a background of Node.js and Go programming, it had been quite a shock to me when I saw the state of Python dependency management. Node.js and Go have canonical dependency management practices, with node you have `npm` and `yarn` and Go it's built into the toolchain. When you are ready to abstract some logic into it’s own library, it’s as simple as creating a new git repository, writing the appropriate metadata files (`package.json`, `go.mod`) push and you have a dependency you can import into your project!

Here are some examples of importing directly from git in the tools I am familiar with:

- Node.js: `yarn add https://github.com/octokit/rest.js.git`
- Go: `go get github.com/google/go-github/v45`

This is so incredibly easy and I want to have the same experience with Python. Is this possible? Almost!

I discovered that I could achieve similar ergonomics using existing and widely used tools, and I would like to demonstrate that here.

## Python Library Project setup

The first step is to setup a Python library project, the best way to go about doing this is to follow the official documentation which can be found here: https://packaging.python.org/en/latest/tutorials/packaging-projects/

I will summarise what needs to be done to demonstrate a working example.

You will need to create the following structure in your git repository:

```
.
├── pyproject.toml
├── README.md
├── src
│   ├── example_package
│   │   ├── example.py
│   │   └── __init__.py
└── tests
    └── test_example.py
```

The contents of these files are listed here, you should update the fields in
`pyproject.toml` to match your organisation or project.

- `pyproject.toml` Update this to match your organisation.

```toml
[build-system]
  build-backend = "hatchling.build"
  requires = ["hatchling"]

[project]
  classifiers = ["Programming Language :: Python :: 3", "License :: OSI Approved :: MIT License", "Operating System :: OS Independent"]
  dependencies = ["boto3==1.23.6"]
  description = "A small example package"
  name = "python-library-test"
  readme = "README.md"
  requires-python = ">=3.8"
  version = "0.0.1"

  [[project.authors]]
    email = "kurt.mcalpine@sourcedgroup.com"
    name = "Kurt McAlpine"

  [project.urls]
    "Bug Tracker" = "https://github.com/kurtmc/python-library-test/issues"
    Homepage = "https://github.com/kurtmc/python-library-test"
```

**Note:** you may add additional dependencies this project may have to the dependencies field under [project] . In this example I have added boto3 as a dependency.

- `src/example_package/example.py` example code

```python
def add_one(number):
    return number + 1
```

- `tests/test_example.py` example test

```python
import unittest

import sys
import os
sys.path.append(os.path.dirname(os.path.realpath(__file__)) + "/../src")
from example_package.example import add_one


class TestExamplePackage(unittest.TestCase):

    def test_add_one(self):
        expected = 1
        actual = add_one(0)
        self.assertEqual(expected, actual)


if __name__ == '__main__':
    unittest.main()
```

Once you have this structure setup, you can commit it and push it to your git repository. The next incredibly useful feature to add will be automatic versioning and tagging. We can use GitHub actions to automatically increment a version number and apply git tags. Later we will use the git tag to specify exactly which version of the library we want to include as a dependency to a new project.

Create the following files:

- `.github/workflows/update-version.yml` You may want to change the branch name if main is not your default branch name.

```yaml
name: Updates version and tags
on:
  push:
    branches:
      - main
permissions:
  contents: write
jobs:
  update_version_and_tag:
    runs-on: ubuntu-latest
    steps:
    - uses: actions/checkout@v2
    - name: Install Python 3
      uses: actions/setup-python@v2
      with:
        python-version: 3.8
    - name: Update version
      uses: kurtmc/github-action-python-versioner@v1
```

Now any changes to `main` will be tagged and the version in `pyproject.toml` will be updated by GitHub actions:

![](https://github.com/kurtmc/blog/raw/master/2022-07/lightweight-artifact-repository-with-python-and-github/images/1.png)

Whilst we are here, we should add a GitHub action that runs on pull requests to enforce code style consistency and validate that the unit tests pass.

Create `.github/workflows/pull-request.yml`:

```yaml
name: Run tests against pull requests
on: pull_request
jobs:
  lint:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v2
        with:
          ref: ${{ github.head_ref }}
      - name: Install Python 3
        uses: actions/setup-python@v2
        with:
          python-version: 3.8
      - name: Lint
        run: |
          pip install flake8==4.0.1
          flake8 ./src --ignore E501
  unit_tests:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v2
        with:
          ref: ${{ github.head_ref }}
      - name: Install Python 3
        uses: actions/setup-python@v2
        with:
          python-version: 3.8
      - name: Install dependencies
        run: |
          pip install -e .
      - name: Run Python unittest
        run: |
          python -m unittest tests/*.py
```

Now we can ensure that all new code added to the library follows consistent code style and the unit tests pass.

![](https://github.com/kurtmc/blog/raw/master/2022-07/lightweight-artifact-repository-with-python-and-github/images/2.png)

We now have a Python library project in GitHub, following code style best practices and automatic version incrementing. How do we import it into a Python project?

Using the git URL in `requirements.txt`:

```
example-python-library @ git+https://github.com/YourOrg/example-python-library.git@0.0.1
```

Now lets try install it:

```shell
$ pip install -r requirements.txt
Collecting example-python-library@ git+https://github.com/YourOrg/example-python-library.git@0.0.1
  Cloning https://github.com/YourOrg/example-python-library.git (to revision 0.0.1) to /tmp/pip-install-1h4qrmmg/example-python-library_22bdfa1c6ab242c18e0e17b700c1be60
  Running command git clone --filter=blob:none --quiet https://github.com/YourOrg/example-python-library.git /tmp/pip-install-1h4qrmmg/example-python-library_22bdfa1c6ab242c18e0e17b700c1be60
Username for 'https://github.com':
Password for 'https://github.com':
  remote: Repository not found.
  fatal: Authentication failed for 'https://github.com/YourOrg/example-python-library.git/'
  error: subprocess-exited-with-error

  × git clone --filter=blob:none --quiet https://github.com/YourOrg/example-python-library.git /tmp/pip-install-1h4qrmmg/example-python-library_22bdfa1c6ab242c18e0e17b700c1be60 did not run successfully.
  │ exit code: 128
  ╰─> See above for output.

  note: This error originates from a subprocess, and is likely not a problem with pip.
error: subprocess-exited-with-error

× git clone --filter=blob:none --quiet https://github.com/YourOrg/example-python-library.git /tmp/pip-install-1h4qrmmg/example-python-library_22bdfa1c6ab242c18e0e17b700c1be60 did not run successfully.
│ exit code: 128
╰─> See above for output.

note: This error originates from a subprocess, and is likely not a problem with pip.
```

This fails to install because in the case of a private repository. We need to tell git to use our SSH credentials when cloning this private repository, which we can do with `git config`:

```shell
git config --global url."git@github.com:".insteadOf "https://github.com/"
```

Attempting the install again:

```shell
$ pip install -r requirements.txt
Collecting example-python-library@ git+https://github.com/YourOrg/example-python-library.git@0.0.1
  Cloning https://github.com/YourOrg/example-python-library.git (to revision 0.0.1) to /tmp/pip-install-z0q1jh7e/example-python-library_0e007d22fbd1439d9481e28d224387bf
  Running command git clone --filter=blob:none --quiet https://github.com/YourOrg/example-python-library.git /tmp/pip-install-z0q1jh7e/example-python-library_0e007d22fbd1439d9481e28d224387bf
  Running command git checkout -q 32600d1874df73fc209736eef6bbd09553cf2dc0
  Resolved https://github.com/YourOrg/example-python-library.git to commit 32600d1874df73fc209736eef6bbd09553cf2dc0
  Installing build dependencies ... done
  Getting requirements to build wheel ... done
  Installing backend dependencies ... done
  Preparing metadata (pyproject.toml) ... done
Requirement already satisfied: boto3==1.23.6 in /usr/local/lib/python3.10/site-packages (from example-python-library@ git+https://github.com/YourOrg/example-python-library.git@0.0.1->-r requirements.txt (line 1)) (1.23.6)
```

:partying_face: This is working locally! Now we should configure our CI/CD platform in the same way if we use SSH authentication, however if you are using personal access tokens registered against a service user you can configure git like this (assuming the personal access token is available under the `GITHUB_TOKEN` environment variable):

```shell
git config --global url."https://${GITHUB_TOKEN}@github.com/".insteadOf "https://github.com/"
```

## Conclusion

Above, is a demonstration on how to build your own private dependency management system for Python using git and GitHub actions. Is this the best solution for private dependency management? Probably not, if you are in the position to pick technologies and services or are starting a greenfield project, you will be able to pick something that works out of the box (examples include: [Artifactory](https://jfrog.com/artifactory/), [Nexus](https://www.sonatype.com/products/nexus-repository), [AWS CodeArtifact](https://aws.amazon.com/codeartifact/)) and establish best practices from the beginning. Not everyone is so lucky, and you may not be able to onboard a new tool so you need to stick with what you already have, and you almost certainly already have GitHub, this may be a solution for you.

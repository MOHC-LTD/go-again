# How to contribute

## Contents

- [⚓️ Git hooks](#-git-hooks)
- [💬 Committing and merging](#-committing-and-merging)

## ⚓ Git hooks

This project includes three git hooks

- `pre-commit`
- `pre-push`
- `commit-msg`

To enable these, run the following command

```sh
git config core.hooksPath .githooks
```

> Note - on a unix machine (mac and linux) you may need to make both files executable

```sh
chmod +x .githooks/pre-commit
chmod +x .githooks/pre-push
chmod +x .githooks/commit-msg
```

## 💬 Committing and merging

We use [commitlint](https://commitlint.js.org/) to ensure that all commit messages meet a certain style and standard. You can find out more about conventional commits [here](https://www.conventionalcommits.org/).

When merging pull requests into the `development` branch on GitHub, use the `Squash and merge` option and make sure the squash commit message follows the style from commitlint and contains the pull request number. Use the normal `Merge pull request` for everything else including the `main` branch.

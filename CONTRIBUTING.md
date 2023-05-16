# Contributing

[fork]: https://github.com/robolaunch/template/fork
[pr]: https://github.com/robolaunch/template/compare
[style]: STYLEGUIDE
[code-of-conduct]: CODE_OF_CONDUCT.md

This project is open for contributions only from the [robolaunch](https://github.com/robolaunch) users for now. To request a new image or feature, please consider opening an issue.

Please note that this project is released with a [Contributor Code of Conduct][code-of-conduct]. By participating in this project you agree to abide by its terms.

## Adding a Pipeline

- First, create an issue for your request in this format: eg. **Add Pipeline for Freecad**
- Create a branch using the issue number: eg. **23-freecad**
- Push the empty branch, and then link it to your issue: eg. `git push origin 23-freecad`
- Do your changes, commit & push them to your branch.
  - Add your Dockerfiles under `dockerfiles/`.
  - Add your pipeline configuration under `pipelines/`.
  - Avoid rebuilding existing images.
  - Version your components and images, use build arguments for seperating them.
- Open a pull request for merging your changes. eg. from `23-freecad` to `main`

## Commit Styling

Use [semantic commit messages](https://gist.github.com/joshbuchea/6f47e86d2510bce28f8e7f42ae84c716):

```bash
# examples
# ........
# use feat for adding a new component
git commit -m "feat(freecad): add dockerfile and pipeline for freecad images"
# use refactor for small changes
git commit -m "refactor(freecad): change linux user in dockerfile"
# use fix for small changes
git commit -m "fix(freecad): fix dockerfile syntax error"
```

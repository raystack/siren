# Release

Siren release tags follow [SEMVER](https://semver.org/) convention.

Github workflow is used to build and push the built docker image to Docker hub.

A release is triggered when a github tag of format `vM.m.p` is pushed. After release job is succeeded, a docker image of
format `M.m.p` is pushed to docker hub.

#!/usr/bin/env bash
set -e


DOCKER_ORG="${DOCKER_ORG-puzzle}"
DOCKER_REPO="${DOCKER_REPO-cryptopus-k8s-secretcontroller}"
DOCKER_REGISTRY="${DOCKER_REGISTRY-docker.io}"


GO_DOCKER_IMAGE=golang:1.11
PACKAGENAME=cryptopus-k8s-secretcontroller
FULL_PACKAGENAME=github.com/puzzle/$PACKAGENAME

if [ -z "$BUILD_VERSION" ]; then
    TAG=1.0-SNAPSHOT
else
    TAG=$BUILD_VERSION
fi

DOCKER_TAG_LATEST="${DOCKER_TAG-${DOCKER_REGISTRY}/${DOCKER_ORG}/${DOCKER_REPO}:latest}"
DOCKER_TAG="${DOCKER_TAG-${DOCKER_REGISTRY}/${DOCKER_ORG}/${DOCKER_REPO}:${TAG}}"

echo "INFO: Building Version: $TAG"
echo "INFO: Prepare file system"
mkdir -p $PWD/build/
echo "INFO: running install"
docker run --rm -v $PWD://go/src/$FULL_PACKAGENAME -v $PWD/build/://go/bin -w //go/src/$FULL_PACKAGENAME -e "CGO_ENABLED=0" -e "GOOS=linux" $GO_DOCKER_IMAGE go install -a -installsuffix nocgo .
echo "INFO: coping Dockerfile"
cp $PWD/Dockerfile $PWD/build/
echo "INFO: changing into build dir"
cd $PWD/build/
echo "🐳 INFO: building image latest"
docker build --pull --tag $DOCKER_TAG_LATEST .
echo "INFO: tagging image version"
docker tag $DOCKER_TAG_LATEST $DOCKER_TAG

if [ -z "$BUILD_VERSION" ]; then
    echo "⚠️  INFO: developer build. Image only locally available"
else
    echo "⏫  INFO: release build. Pushing the image to the registry"
    docker push $DOCKER_TAG_LATEST
    docker push $DOCKER_TAG
fi

echo "✅ INFO: ... build successful !"

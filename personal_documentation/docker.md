Docker
======

Build
-----

```
docker build path/to/Dockerfile

# You can also tag while building
docker build -t quay.io/cbone/scalp:http path/to/Dockerfile
```

Run
---

```
docker run <img_identifier>
```

Tagging
-------

See `docker help tag`

Example:

```
docker tag <image_id> quay.io/cbone/scalp:http
```

Pushing
-------

See `docker help push`

Example:

```
# This will push just this tagged image
docker push quay.io/cbone/scalp:http
# This will push ALL tagged images of this one
docker push quay.io/cbone/scalp
```

Dangling Volumes
----------------

You can remove all volumes by doing something like:

```
docker volume ls -qf dangling=true | xargs docker volume rm
```

About Images and Containers
---------------------------

http://merrigrove.blogspot.co.uk/2015/10/visualizing-docker-containers-and-images.html

Echo UDP Server
---------------

Useful to test that things like Datadog or statsd are actually sending data:
https://hub.docker.com/r/eexit/dumudp-server/

Volumes
-------

More information on volumes:
http://container-solutions.com/understanding-volumes-docker/

docker-compose
--------------

version 2 apparently by default creates a network through which all containers
can talk to one another so you don't really need the "links" option anymore
(this would create the network between 2 containers and do startup order). All
you need is "depends_on":
https://medium.com/@giorgioto/docker-compose-yml-from-v1-to-v2-3c0f8bb7a48e#.1u9f16kj1

version 2 you can configure aspects about the network over which docker
containers communicate with. I'm really not familiar with it:
https://docs.docker.com/compose/networking/.

Create And Run Example
---------------------

docker container create --name lucas2 -it busybox
docker start -i lucas2

## Inject Environment Variables Into Docker Image
Whenever you build a docker image the ONLY way to pass in arbitrary
arguments to it is with the `ARG` directive. For example, if you have
this dockerfile:

```
FROM alpine
ARG SOME_ARG
RUN echo "$SOME_ARG"
```

Then if you do `docker build --build-arg SOME_ARG=hello .` `hello`
will be echo'd when the image is being built. I think `ARG`
essentially lets you create environment variables which are ONLY
accessible during the build process.

Another directive is the `ENV` directive. What this does is create an
environment variable in the image that will be built. Like `ARG` this
environment variable also exists during the build process (but
persisting an environment variable into the resulting image is the
most important thing). For example if you have this dockerfile:

```
FROM alpine
ENV SOME_ENV_VAR="wow there"
RUN echo "$SOME_ENV_VAR"
```

Then when you build it with `docker build .` it will echo out `wow
there` AND (more importantly) any containers you start from this image
will have the environment variable `SOME_ENV_VAR` with the value `wow
there`.

If you combine these two concepts you get the ability to inject
environment variables into the image:

```
FROM alpine
ARG MY_FIRST_ARG
ARG MY_SECOND_ARG
ENV MY_FIRST_ENV="$MY_FIRST_ARG" \
    MY_SECOND_ENV="$MY_SECOND_ARG" \
    MY_THIRD_ENV="$PWD" \
    MY_FOURTH_ENV="HELLO"
```

If you build this image `docker build -t lucas-test --build-arg
MY_FIRST_ARG=ooweee --build-arg MY_SECOND_ARG=goshhh .` then the
environment variables will be:
- MY_FIRST_ENV=oowee
- MY_SECOND_ENV=goshhh
- MY_THIRD_ENV=
- MY_FOURTH_ENV=HELLO

`MY_THIRD_ENV` is empty because when building the docker image it does
NOT have access to your local environment variables. You can have the
`ARG` and environment variable have the same name if you want:

```
FROM alpine
ARG MY_ENV_VAR
ENV MY_ENV_VAR="$MY_ENV_VAR"
```

This is a very useful technique! One application is to inject
environment variables indicating which github hash or version is
deployed which could hypothetically help when debugging.

When passing `--build-arg` through docker-compose you'll use this sort
of code:

```
build:
  context: .
  args:
    - buildno=1
    - password=secret
	- value_from_env
```

`buildno` and `password` will have the values you see above but since
you left out a value for `value_from_env` it will get its value from
the environment variable with the same name. So that snippet of code
is equivalent to: `docker build --build-arg buildno=1 --build-arg
password=secret --build-arg value_from_env=$value_from_env`.

## Caching
You want to leverage caching as much as possible so your builds will
be as fast as possilbe. So we need to be aware of docker's algorithm
for when a cached image can be used instead of created:
https://docs.docker.com/develop/develop-images/dockerfile_best-practices/#build-cache

The basic algorithm seems to be:
1. We start at the FROM image
2. We look at the next instruction:
   - If the instruction is the same as it was before then we use the
     cache.
   - Otherwise we run that instruction and ALL subsequent steps will
     be redone.
3. While there are more instructions go to step 2

With that in mind
1. Put instructions that do not change much at the top of the
   Dockerfile because if something at the top of the Dockerfile
   changes then ALL subsequent steps will be re-ran. Note that when
   doing something like "COPY . ." if ANY file you are copying has
   changed then there will be no cache hit and all subsequent steps
   will be re-run.
2. Make use of multi-stage builds whenever useful to get more out of
   the caching. For example I was trying to Dockerize a legacy
   application where essentially the WORKDIR path will have the
   version being deployed in it. So that means WORKDIR will always be
   changing which will cause all subsequent commands to not use the
   cache. So having a separate build stage where we downloaded
   dependencies and then copied those over to the new location helped.

## Scratch
The scratch image is very limited and will probably need extra things
like CA certs. For instance one time I dealt with an application which
set the timezone to something other than the default of UTC so the
timezone data /usr/share/zoneinfo/America/Chicago had to be copied
into the scratch image. For CA certs, the file containing all of them
lives in /etc/ssl/certs/ca-certificates.crt. This tutorial sums that
up nicely:
https://sebest.github.io/post/create-a-small-docker-image-for-a-golang-binary/

## Docker vs Virtual Machine (VM)
This one has confused me for a while and I think still confuses me a
bit. I'll document here things that I've learned.

In one sentence: a container (i.e. docker) is process and a VM is a
server

Containers are (I *think*) a combination of two linux kernel features:
1. namespace - controls which resources are visible to a process
2. control group - controls how much of certain resources can be used

Containers have been around before docker came along but docker made
interacting with them much easier. In addition to dealing with
namespaces and c groups docker also uses a union file system which
promotes reusability:
https://medium.freecodecamp.org/a-beginner-friendly-introduction-to-containers-vms-and-docker-79a9e3e119b

Here is what a container looks like:
```
                +------------+
                | Container  |
+------------+  +------------+  +------------+
|   App A    |  |   App B    |  |   App C    |
+------------+  +------------+  +------------+
|  Bins/Libs |  |  Bins/Libs |  |  Bins/Libs |
+------------+--+------------+--+------------+
|                   Docker                   |
+--------------------------------------------+
|                  Host OS                   |
+--------------------------------------------+
|               Infrastructure               |
+--------------------------------------------+
```

Here is what a VM looks like:
```
                +------------+
                |     VM     |
+------------+  +------------+  +------------+
|   App A    |  |   App B    |  |   App C    |
+------------+  +------------+  +------------+
|  Bins/Libs |  |  Bins/Libs |  |  Bins/Libs |
+------------+  +------------+  +------------+
|  Guest OS  |  |  Guest OS  |  |  Guest OS  |
+------------+--+------------+--+------------+
|                 Hypervisor                 |
+--------------------------------------------+
|                  Host OS                   |
+--------------------------------------------+
|               Infrastructure               |
+--------------------------------------------+
```

So a container is a way to separate processes so it *looks* like they
operate by themselves but each container shares the same kernel (i.e.
part of the operating system that deals with devices, process
management, memory management, and system calls). Note that ALL linux
distributions share the same kernel and the only difference is
userland stuff:
https://serverfault.com/questions/755607/why-do-we-use-a-os-base-image-with-docker-if-containers-have-no-guest-os
All docker containers run on linux and when I do docker stuff on my
mac those containers do indeed run on a linux kernel. I think there's
some linux VM'ing going on but for just the kernel so I think that
would mean that there's a hypervisor? A hypervisor by the way is a
process that sits between an OS and the hardware. A hypervisor runs
your VMs. I believe that hypervisors contain kernel logic within them
(which an OS needs in order to run).

TODO: Is there any reason why ALL processes don't run inside a
continer of some sort? Seems like a nice idea?

## Copy from container to local filesystem
`docker cp container:container/path local/path`

https://stackoverflow.com/questions/22049212/copying-files-from-docker-container-to-host

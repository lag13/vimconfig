REST
====

Roy Fieldings's Dissertation IS REST
--------------------------------

I have learned today (01/25/2017) that Roy Fieldings dissertation originally
brought REST into being:
https://www.ics.uci.edu/~fielding/pubs/dissertation/top.htm. In other words
this is the true definition of REST and I should probably read it along with
his blog (http://roy.gbiv.com/untangled/tag/rest). Judging from  my previous
difficulties in getting a solid definition it would seem that REST's
definition is different to a lot of people (probably a product of
people misinterpreting his dissertation, not realizing it, and spreading that
knowledge. 

Links
-----

A tutorial:

http://www.restapitutorial.com/

Other book recommendations:

https://killring.org/2014/11/18/api-design-reading-list/

http://stackoverflow.com/questions/630453/put-vs-post-in-rest/2590281#2590281

https://apihandyman.io/do-you-really-know-why-you-prefer-rest-over-rpc/

https://www.troyhunt.com/your-api-versioning-is-wrong-which-is/ and this one
linked from that one
(http://codebetter.com/howarddierking/2012/11/09/versioning-restful-services/
<<--- I really like this guy. Check out his other stuff too. His current blog
is here: http://blog.howarddierking.com/)

https://github.com/interagent/http-api-design

https://leanpub.com/build-apis-you-wont-hate

http://swagger.io/

Not REST per se, just about API's in general and why they are good:
http://fr.slideshare.net/faberNovel/6-reasons-why-apis-are-reshaping-your-business

URL vs URI: https://danielmiessler.com/study/url-uri/

https://www.quora.com/What-are-some-guidelines-to-build-a-good-RESTful-API#ans1036119

http://barelyenough.org/blog/2008/05/versioning-rest-web-services/

http://timelessrepo.com/haters-gonna-hateoas

http://blog.steveklabnik.com/posts/2011-07-03-nobody-understands-rest-or-http

http://ivanzuzak.info/2010/04/03/why-understanding-rest-is-hard-and-what-we-should-do-about-it-systematization-models-and-terminology-for-rest.html

The URI Is The Thing: https://www.flickr.com/photos/psd/2918889380/

Another nice discussion of REST: http://homepages.tig.com.au/~ijoyner/Ian_Joyner/REST.html

REST vs SOAP: http://spf13.com/post/soap-vs-rest/

Snarky Comment
--------------

http://stackoverflow.com/questions/15056878/rest-vs-json-rpc

The first comment on this response probably perfectly sums up my frustration
with REST: http://stackoverflow.com/a/20643993. Basically REST specifications
seem sort of arbitrary sometimes and then some guy comes along and says this:

```
This answer shows the all-too usual misconception of what REST actually is.
REST is definitely not just a mapping of CRUD to HTTP methods. The idea that
it is a problem to "add another method" clearly indicates that REST is
misunderstood as RPC over HTTP, which it is not at all. Try reading Roy
Fieldings blog or his dissertation - Google will help you find it - you are
not describing REST at all in your answer.
```

Now maybe this guy is right (I haven't read Roy Fielding's dissertation and
maybe I should: https://www.ics.uci.edu/~fielding/pubs/dissertation/top.htm)
but the fact is based on most of the tutorials and such that I've looked at,
almost everyone thinks of REST as being a consistent approach to doing CRUD
operations and instead of making some snarky comment about "oh, you don't
really understand how rest works" it should be more constructive or something.

Upon rereading that comment (2018-10-16) it doesn't really seem THAT
snarky. I was probably frustrated at the time about snarky comments
similar to that one so I projected more snarkiness into the tone than
was actually there.

This Person Puts It Nicely
--------------------------

From:
http://codebetter.com/glennblock/2012/03/11/you-cant-achieve-rest-without-client-and-server-participation/.

```
I think that the mere fact that a REST expert like yourself has to be
constantly changing his mind  just proves that the term REST almost serves no
purpose other than confuse people... let's just stick to HTTP web apis with a
set of recommended practices for both client/server that developers can choose
to follow or not.
```

The author then had a nice response and linked to a nice site:
https://theamiableapi.com/2012/03/18/rest-design-philosophies/

```
I agree that "REST" has been abused too much to make it meaningful any longer.
I found this post useful for defining terms: https://theamiableapi.com/2012/03/18/rest-design-philosophies/
```

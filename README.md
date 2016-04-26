# randomImage
Used to grab random image from configured sources according to tag

# Config
Each source is responsible for its own config.
## Tumblr
A list of json objects
```
[
        {
                "url"          : "http://emergencykittens.tumblr.com/",
                "img_path"     : "//*[@class=\"photo-wrapper-inner\"]//img",
                "size_path"    : "//*[@id=\"pagination\"]/a",
                "size_pattern" : "data-total-pages=\"(\\d+)",
                "tags"         : [ "cute", "kittens" ]
        }
]
```

# TODO
- [ ] Write documentation
- [ ] Refactor code to a less hacky structure

The only reason the code is published in this state is because someone asked to see it right away

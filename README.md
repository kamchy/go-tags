# Goal
This program prints to stdout all tags found in  metadata (so called _frontmatter_) of `.md` files under given directory 
(think: Hugo blog's `content/` directory.)

**Caution !**

This is my second (!) program written in Go so far so it is not production-ready. It is just language-learning exercise. You've been warned.
(C'mon, nobody is ever going to run this program except myself so this disclaimer is just not needed.)

## Building:

```bash
go build .
```

##Executing:
The program takes single flag `-p` specifying the directory that will be searched for .md files. If not given, current dirrectory is assumed.

```bash
./tags -p ~/blog/content/
```

It outputs a "histogram" of tags used in frontmatter of all found .md files.
Example output:

```bash
Searching path: /home/user/blog/content/
  18:  java
   9:  grafika
   6:  deno
   6:  javascript
   5:  aoc
   5:  blog
   5:  programowanie
   5:  python
   3:  go
   3:  javafx
   2:  bash
   2:  howto
   2:  typescript
   1:  jshell
   1:  matematyka
   1:  maven
   1:  nikola
   1:  programowaie
   1:  ripgrep
   1:  sed
   1:  unicode
```

## Note
This program assumes simple TOML metadata with tags specified like:
```toml
tags=["bash", "maven"]
```


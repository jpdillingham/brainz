```
       ,. brainz
 (¬º-°)¬ 
```

# brainz
A lightweight Musicbrainz CLI in Go.

### Example (interactive mode)

Allow the application to prompt for input:

```
$ go run brainz.go

       ,. brainz
 (¬º-°)¬

Enter artist: nine inch nails
Enter album: the fragile

Searching for artists matching 'nine inch nails'...
```

Or, supply it as command line arguments:

```
$ go run brainz.go --artist "nine inch nails" --album "the fragile"

       ,. brainz
 (¬º-°)¬

Searching for artists matching 'nine inch nails'...

--> 100%        Nine Inch Nails
     47%        Tapeworm (Nine Inch Nails side project)
     44%        Static X (Nine Inch Nails trance remix artist)
     41%        Kobold (Nine Inch Nails remixer, first name "Marie")
     34%        Nine Inch Elvis

Best artist: Nine Inch Nails (b7ffd2af-418f-4be2-bdd1-22f8b48613da) (Score: 100%)

Searching for release group matching 'the fragile'...

--> 100%        The Fragile             Album
     73%        The Fragile Remixes     Album, Remix
     61%        Fragiliglade            Album, Live
     55%        Hellraisers             Album, Live
     53%        The Slip                Album

Best release group: The Fragile (fc4d7589-a6b9-35ca-a6b6-8c8e70c3baaa) (Score: 100%)

Fetching releases...

17 releases fetched. Computing canonical track count(s)...

-->   9x        12 + 11
      3x        8 + 8 + 9
      2x        28
      1x        9 + 10 + 9 + 9 + 37
      1x        25
      1x        12 + 14

Canonical track count(s): 12 + 11

Probable canonical releases:

1999-09-20      The Fragile                     2xCD
1999-09-20      The Fragile                     2xCD
1999-09-20      The Fragile                     2xCD
1999-09-21      The Fragile (BMG club edition)  2xCD
1999-09-21      The Fragile                     2xCD
1999-09-21      The Fragile                     2xCassette
1999-09-22      The Fragile                     2xCD
1999-12-31      The Fragile                     2xCassette
2006-05-17      The Fragile                     2xCD

Best probable canonical release: The Fragile (1999-09-20, 2xCD, 7b781387-5a2f-3984-aee4-a14f09ea7376) (Score: 99%):

  101   Somewhat Damaged                                100%
  102   The Day the World Went Away                     100%
  103   The Frail                                       100%
  104   The Wretched                                    100%
  105   We’re in This Together                          100%
  106   The Fragile                                     100%
  107   Just Like You Imagined                          100%
  108   Even Deeper                                     100%
  109   Pilgrimage                                      100%
  110   No, You Don’t                                   100%
  111   La Mer                                          100%
  112   The Great Below                                 100%
  201   The Way Out Is Through                          100%
  202   Into the Void                                   100%
  203   Where Is Everybody?                             100%
  204   The Mark Has Been Made                          100%
  205   Please                                           89%    Please (+Appendage)
  206   Starfuckers, Inc.                               100%
  207   Complication                                    100%
  208   I’m Looking Forward to Joining You, Finally     100%
  209   The Big Come Down                               100%
  210   Underneath It All                               100%
  211   Ripe (With Decay)                                89%    Ripe (with Decay)
```

### Example (automated mode, with pretty printed json)

```
$ go run brainz.go --artist "nine inch nails" --album "the fragile" --output json

{
    "artist": "Nine Inch Nails",
    "album": "The Fragile",
    "mbid": "7b781387-5a2f-3984-aee4-a14f09ea7376",
    "score": 0.9903381642512078,
    "tracks": [{
        "disc": 1,
        "position": 1,
        "number": "101",
        "title": "Somewhat Damaged",
        "length": 271693,
        "score": 1
    }, {
        "disc": 1,
        "position": 2,
        "number": "102",
        "title": "The Day the World Went Away",
        "length": 273666,
        "score": 1
    }, {
        "disc": 1,
        "position": 3,
        "number": "103",
        "title": "The Frail",
        "length": 114040,
        "score": 1
    }, {
        "disc": 1,
        "position": 4,
        "number": "104",
        "title": "The Wretched",
        "length": 325893,
        "score": 1
    }, {
        "disc": 1,
        "position": 5,
        "number": "105",
        "title": "We’re in This Together",
        "length": 436000,
        "score": 1
    }, {
        "disc": 1,
        "position": 6,
        "number": "106",
        "title": "The Fragile",
        "length": 275466,
        "score": 1
    }, {
        "disc": 1,
        "position": 7,
        "number": "107",
        "title": "Just Like You Imagined",
        "length": 229866,
        "score": 1
    }, {
        "disc": 1,
        "position": 8,
        "number": "108",
        "title": "Even Deeper",
        "length": 347933,
        "score": 1
    }, {
        "disc": 1,
        "position": 9,
        "number": "109",
        "title": "Pilgrimage",
        "length": 211933,
        "score": 1
    }, {
        "disc": 1,
        "position": 10,
        "number": "110",
        "title": "No, You Don’t",
        "length": 215933,
        "score": 1
    }, {
        "disc": 1,
        "position": 11,
        "number": "111",
        "title": "La Mer",
        "length": 277866,
        "score": 1
    }, {
        "disc": 1,
        "position": 12,
        "number": "112",
        "title": "The Great Below",
        "length": 317173,
        "score": 1
    }, {
        "disc": 2,
        "position": 1,
        "number": "201",
        "title": "The Way Out Is Through",
        "length": 259413,
        "score": 1
    }, {
        "disc": 2,
        "position": 2,
        "number": "202",
        "title": "Into the Void",
        "length": 291293,
        "score": 1
    }, {
        "disc": 2,
        "position": 3,
        "number": "203",
        "title": "Where Is Everybody?",
        "length": 342186,
        "score": 1
    }, {
        "disc": 2,
        "position": 4,
        "number": "204",
        "title": "The Mark Has Been Made",
        "length": 317146,
        "score": 1
    }, {
        "disc": 2,
        "position": 5,
        "number": "205",
        "title": "Please",
        "length": 212693,
        "score": 0.8888888888888888,
        "alternate-titles": ["Please (+Appendage)"]
    }, {
        "disc": 2,
        "position": 6,
        "number": "206",
        "title": "Starfuckers, Inc.",
        "length": 302213,
        "score": 1
    }, {
        "disc": 2,
        "position": 7,
        "number": "207",
        "title": "Complication",
        "length": 152760,
        "score": 1
    }, {
        "disc": 2,
        "position": 8,
        "number": "208",
        "title": "I’m Looking Forward to Joining You, Finally",
        "length": 255346,
        "score": 1
    }, {
        "disc": 2,
        "position": 9,
        "number": "209",
        "title": "The Big Come Down",
        "length": 254853,
        "score": 1
    }, {
        "disc": 2,
        "position": 10,
        "number": "210",
        "title": "Underneath It All",
        "length": 168120,
        "score": 1
    }, {
        "disc": 2,
        "position": 11,
        "number": "211",
        "title": "Ripe (With Decay)",
        "length": 394320,
        "score": 0.8888888888888888,
        "alternate-titles": ["Ripe (with Decay)"]
    }]
}
```
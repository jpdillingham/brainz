```
       ,. brainz
 (¬º-°)¬ 
```

# brainz
A lightweight Musicbrainz CLI in Go.

### Example

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

Determining canonical track listing...

Compiling formats..

-->   7x        2xCD                            12 + 11
      3x        3x12" Vinyl                     8 + 8 + 9
      2x        2xCassette                      12 + 11
      2x        Digital Media                   28
      1x        4x12" Vinyl + 1xDigital Media   9 + 10 + 9 + 9 + 37
      1x        2xDigital Media                 12 + 14
      1x        Digital Media                   25

Canonical format: 2xCD, canonical track count: 12 + 11

Probable canonical releases:

   1999-09-20   The Fragile                     2xCD
   1999-09-20   The Fragile                     2xCD
   1999-09-20   The Fragile                     2xCD
   1999-09-21   The Fragile                     2xCD
   1999-09-21   The Fragile (BMG club edition)  2xCD
   1999-09-22   The Fragile                     2xCD
   2006-05-17   The Fragile                     2xCD

Probable canonical track listing:

   100% 101     Somewhat Damaged
   100% 102     The Day the World Went Away
   100% 103     The Frail
   100% 104     The Wretched
   100% 105     We’re in This Together
   100% 106     The Fragile
   100% 107     Just Like You Imagined
   100% 108     Even Deeper
   100% 109     Pilgrimage
   100% 110     No, You Don’t
   100% 111     La Mer
   100% 112     The Great Below
   100% 201     The Way Out Is Through
   100% 202     Into the Void
   100% 203     Where Is Everybody?
   100% 204     The Mark Has Been Made
   100% 205     Please
   100% 206     Starfuckers, Inc.
   100% 207     Complication
   100% 208     I’m Looking Forward to Joining You, Finally
   100% 209     The Big Come Down
   100% 210     Underneath It All
   100% 211     Ripe (With Decay)
```
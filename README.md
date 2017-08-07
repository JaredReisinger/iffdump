# iffdump

Tools to dump IFF-formatted files.


## design

There are some different approaches to understanding an IFF file:

* Run through start-to-finish, dealing with chunks as they are encountered
* unmarshal into strongly-typed structure.

While the former seems like a reasonable first approach, it makes it much harder
to *do* anything with the data... for example, in attempting to understand a
Quetzal save file, we really do need to understand the IFZS, IFhd, CMem, UMem,
Stks, and IntD chunks.

So... let's look/compare the JSON decode (and other unmarshal) logic.  Here, a
very generic IFF-chunk reader *is* useful, as a utility.  Really, we just want
something that, given an io.Reader (or iff.ReadAtSeeker), can give us a series
of chunks.  Then, each chunk has its own unmarshalling, which may use the
"series of chunks" logic again.

So, we could create something that acts like a BinaryUnmarshaler, but that
interface reads from a `[]byte`... which we may not want to assume!  (That
said, `json.NewDecoder()` takes an `io.Reader`... so go figure!)

In order to successfully unmarshal/decode, the decoder needs hooks for every
chunk type, so it can call them to handle the data...  Unlike `json.Decode()`,
we can't provide the typed struct first... we have to read a little bit in order
to determine the struct needed.  Thus, the instead of decoding *into* a struct,
the decoders must *create* the struct/data.

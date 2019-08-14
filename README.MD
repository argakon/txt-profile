# This implementation is rewrited from Apache Nutch.
An implementation of a page signature. It calculates an MD5 hash of a plain text "profile" of a page.
The algorithm to calculate a page "profile" takes the plain text version of a page and performs the following steps:

* remove all characters except letters and digits, and bring all characters to lower case,
* split the text into tokens (all consecutive non-whitespace characters),
* discard tokens equal or shorter than MIN_TOKEN_LEN (default 2 characters),
* sort the list of tokens by decreasing frequency,
* round down the counts of tokens to the nearest multiple of QUANT (`QUANT = QUANT_RATE * maxFreq`, where `QUANT_RATE` is 0.01f by default, and `maxFreq` is the maximum token frequency). If `maxFreq` is higher than 1, then QUANT is always higher than 2 (which means that tokens with frequency 1 are always discarded).
* tokens, which frequency after quantization falls below QUANT, are discarded.
* create a list of tokens and their quantized frequency, separated by spaces, in the order of decreasing frequency.

This list is then submitted to an MD5 hash calculation.

#### Example
```
echo "The Rocky Horror Picture Show premiered in London, the first in its record-breaking run in cinemas, which continues in limited release." | go run textprofile.go
```

#### Result
```
592325199d7e02121eaf3304d10cc681
```
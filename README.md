# `pwr`: transfer setting from a java .properties file into another one.

use it as follows:

	pwr -from src.properties -to dst.properties

Assume, this is `src.properties` file:

	maintainer: Rabbit <rabbit@europe.eu>
	telephon: +49 4433 012345

`dst.properties` has the following form:

	package: pwr
	description: <to be inserted>
	maintainer: <to be inserted>
	version: 0.0.1

After executing this program `dst.properties` has changed to:


	package: pwr
	description: <to be inserted>
	maintainer: Rabbit <rabbit@europe.eu>
	version: 0.0.1

Herein, setting for `maintainer was transferred from `src.properties` to `dst.properties`.


## Presuppositions:

* both files must obey the [standard](https://en.wikipedia.org/wiki/.properties) for a .properties file.
* both files should be encoded in ISO 8859.1
* both files should use the same eol style.  either both files use Window's crlf style or both files' lines end with Unix' lf style.
* the last line in `src.properties` should be complete, i.e. ending with eol, i.e. crlf or lf accoring to eol style.

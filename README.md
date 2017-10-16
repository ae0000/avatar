A V A T A R
===========

Create avatars based on names.

Example
-------

```
package main

import (
	"github.com/ae0000/avatar"
)

func main() {
	avatar.ToDisk("AE", "../ae.png")

    // Which is the same as
    avatar.ToDisk("Andrew Edwards", "../ae.png")
}


```
[![Example](https://raw.githubusercontent.com/ae0000/avatar/master/ae.png)](https://raw.githubusercontent.com/ae0000/avatar/master/ae.png)

You can pass in a single character as well

```
// Single initial as well..
avatar.ToDisk("Jet", "../j.png")
```
[![Example](https://raw.githubusercontent.com/ae0000/avatar/master/j.png)](https://raw.githubusercontent.com/ae0000/avatar/master/j.png)
# sedgf

A wrapper around sed to avoid typing common patterns based on gf. All credits for almost of this code go to @tomnomnom.

## Changes from gf

I use sed command combined with gf to fuzz the params in the gf patterns. So I made some changes to use the same patterns with sed. If you wanna use it, you have to add the line 'flags-sed' in the config .json file to add options to the sed command.

For example:

```
{
  "flags": "-iE",
  "flags-sed" : "-E",
  "patterns": [
	"image_url=",
	"open=",

  ...
```

## How it works?

For example, if you want to find some open redirects, and your tools works behind the word FUZZ and obviously you don't want to fuzz every param. You can do something like this (It was intended to work with qsreplace too) :

```
echo "https://example.com?redirect=https://evil.com&param1=test1&param2=test2" | qsreplace newval | sedgf -replace newval open-redirect
```

Given the following result:

```
https://example.com?redirect=FUZZ&param1=newval&param2=newval
```

Now you can connect the output with any tool or fuzzer.

## Example

An example more complex searching open redirect:

```
echo "testphp.vulnweb.com" | gau | gf open-redirect | qsreplace newval | sedgf -replace newval open-redirect > urls.txt
python3.7 openredirex.py -l urls.txt -p payloads.txt --keyword FUZZ 
```

Another example if you wanna replace directly the payload, you can use the option -new-value (don't forget escape characters like '/','[',']'):

```
echo "testphp.vulnweb.com" | gau | gf open-redirect | qsreplace newval | sedgf -replace newval -new-value https:\/\/myinteract.sh open-redirect | rush 'curl -L -k -s -v {}' > /dev/null  

```



## Options

```

Usage of sedgf:

  -dump
        prints the sed command rather than executing it
  -list
        list available patterns
  -new-value string
        new value to replace. (default "FUZZ")
  -replace string
        value to replace for FUZZ. By default search empty params

```


## Install 

```
go get -u github.com/c4rrilat0r/sedgf

```

## TODO

- Not depend of qsreplace
- Look for another more efficient way.

## References

- https://github.com/lc/gau
- https://github.com/tomnomnom/gf
- https://github.com/shenwei356/rush
- https://github.com/devanshbatham/OpenRedireX
- https://github.com/tomnomnom/qsreplace

## Credits

- Idea from https://github.com/tomnomnom/gf - @tomnomnom.

# gograbber
A horizontal and vertical web content enumerator by swarley7 (@swarley777)

# Introduction

`gograbber` is a shitty program I made to solve a number of problems I have with content enumerators:
 - They only typically work on a single URL
 - They produce poor output (xml!? or non-greppable garbage)
 - Tools that produce screenshots of webpages typically only test `/`
 - aquatone rules, but only if you have a domain to throw at it... on most pentests we have a list of net ranges or ip addresses :(
 - They can be very slow (I'm very impatient)

`gograbber` attempts to solve these problems!
- Supply a list of urls, hosts, ip addresses, CIDR ranges, ports, web paths, whatever... and `gograbber` will attempt to discover stuff there.
- Screenshot discovered content! (can be tuned to prevent excessive output)
- output is greppable

# TODO

- [x] write gograbber
- [ ] make it not shit

```
Examples for ./gograbber:

>> Scan and dirbust the hosts from hosts.txt.
./gograbber -i hosts.txt -w wordlist.txt -t 2000 -scan -dirbust

>> Scan and dirbust the hosts from hosts.txt, and screenshot discovered web resources.
./gograbber -i hosts.txt -w wordlist.txt -t 2000 -scan  -dirbust -screenshot

>> Scan, dirbust, and screenshot the hosts from hosts.txt on common web application ports. Additionally, set the number of phantomjs processes to 3.
./gograbber -i hosts.txt -w wordlist.txt -t 2000 -p_procs=3 -p top -scan -dirbust -screenshot

>> Screenshot the URLs from urls.txt. Additionally, use a custom phantomjs path.
./gograbber -U urls.txt -t 200 -j 400 -phantomjs /my/path/to/phantomjs -screenshot

>> Screenshot the supplied URL. Additionally, use a custom phantomjs path.
./gograbber -u http://example.com/test -t 200 -j 400 -phantomjs /my/path/to/phantomjs -screenshot

>> EASY MODE/I DON'T WANT TO READ STUFF LEMME HACK OK?.
./gograbber -i hosts.txt -w wordlist.txt -easy
```

# Acks

- OJ Reeves: this project borrows heavily from `github.com/OJ/gobuster` (which is awesome, btw.) so thanks!
- Aquatone (name?): this project was heavily inspired by aquatone (which is also awesome)
- C_Sto: thx for forcing me to learn Golang and laughing at my extreme incompetence

# Donate?

If you like this *thing* and want to shout a beer maybe do that?
- ETH: `0x486b0faea72a17425ed7241e44dc9ed627f9e492`
- BTC: `1Jdz37JDyZYnK7tRDkF9ZW8QJ2bk2DNHzh`

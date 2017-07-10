# ipv4

[![GoDoc](https://godoc.org/github.com/signalsciences/ipv4?status.svg)](https://godoc.org/github.com/signalsciences/ipv4) [![Build Status](https://travis-ci.org/signalsciences/ipv4.svg?branch=master)](https://travis-ci.org/signalsciences/ipv4)

Package for conveniently working with IPv4 and CIDR ranges.

## Examples

* Get start and end IPs in a CIDR range:

```
left, right, err := ipv4.CIDR2Range("199.27.72.0/21")
if err != nil {
  log.Fatal(err)
}
fmt.Println(left, right)
```

Output:

```
199.27.72.0 199.27.79.255
```

* Check if IP is IPv4 (works for CIDR too):

```
fmt.Println(ipv4.IsIPv4("10.0.0.0"))
fmt.Println(ipv4.IsIPv4("10.0.0.0/8"))
```

Output:

```
true
true
```

* Check if IP is private:

```
fmt.Println(ipv4.IsPrivate("10.0.0.0"))
```

Output:
```
true
```

See [GoDoc](http://godoc.org/github.com/signalsciences/ipv4) for more.

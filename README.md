# A leaky way to view `dhcpd` DHCP leases

Allows one to turn this:

```
# The format of this file is documented in the dhcpd.leases(5) manual page.
# This lease file was written by isc-dhcp-4.3.6-P1

# authoring-byte-order entry is generated, DO NOT DELETE
authoring-byte-order little-endian;

lease 192.168.1.5 {
  starts 6 2019/04/27 03:24:45;
  ends 6 2019/04/27 03:34:45;
  tstp 6 2019/04/27 03:34:45;
  cltt 6 2019/04/27 03:24:45;
  binding state free;
  hardware ethernet ba:dc0:0f:fe:ca:dd:47;
  uid "";
  client-hostname "derranged";
}
```

to show [localhost:3000/](localhost:3000/) as


![This](https://github.com/npotts/lease-leaker/raw/master/.github/scr.png)

and [localhost:3000/json](localhost:3000/json) as

```json
{
	"error": null,
	"leases": [{
		"ip": "192.168.1.5",
		"starts": "2019-04-27T03:24:45Z",
		"ends": "2019-04-27T03:34:45Z",
		"tstp": "2019-04-27T03:34:45Z",
		"tsfp": "0001-01-01T00:00:00Z",
		"atsfp": "0001-01-01T00:00:00Z",
		"cllt": "2019-04-27T03:24:45Z",
		"binding-state": "free",
		"next-binding-state": "",
		"hardware": {
			"hardware": "ethernet",
			"mac": "ba:dc0:0f:fe:ca:dd:47"
		},
		"uid": "",
		"client-hostname": "derranged"
	}]
}
```


Im not sure anyone outside of me would want this - I plan on using it in a pfsense router with an internal port exposed.  It allows me to lazily get a list of what hosts are currently connected to the router for troubleshooting.



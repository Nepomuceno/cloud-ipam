## cloud-ipam env add

Add or update environment

```
cloud-ipam env add [flags]
```

### Options

```
  -h, --help             help for add
  -i, --id string        ID of the environment (be default the same as name lowercase and dashes for spaces)
      --ip-range ipNet   IP range to be used by the environment (default 10.0.0.0/8 (default 10.0.0.0/8)
  -n, --name string      Name of the environment
```

### Options inherited from parent commands

```
  -c, --config string                 config file
      --log-file string               log file (default ".cloud-ipam.log")
      --log-level string              log level (default "INFO")
  -s, --storage-account-name string   storage account name
      --table-name string             ipam storage table name (default "cloudIpam")
```

### SEE ALSO

* [cloud-ipam env](cloud-ipam_env.md)	 - Manage cloud ipam environments

###### Auto generated by spf13/cobra on 22-Aug-2022

[![Build Status](https://travis-ci.com/ahmet2mir/terraform-provider-freeipa.svg?branch=master)](https://travis-ci.com/ahmet2mir/terraform-provider-freeipa)

/!\/!\/!\/!\/!\/!\/!\

This a fork of https://github.com/fiveai/terraform-provider-freeipa
Provider use my fork of GoIPA package from which is a fork of https://github.com/fiveai/goapi
which is a fork from https://github.com/ubccr/goipa

See https://github.com/ubccr/goipa/issues/4

/!\/!\/!\/!\/!\/!\/!\

# terraform-provider-freeipa

## Build

Makefile has been made basically and will attempt to crosscompile for all distros. Due to kerberos and
using a non-native (no golang) library this requires gcc compiler which breaks very easy across the
distributions so your probably end up with failures.

TODO: separate make tasks to specify distro on build for local

```bash
make dist
```

## Usage Example

### Provider

For a test example you might not of configured your own certificates. If so
please download the self-signed certs from `<your-domain>/ipa/config/ca.crt` and
put save to `/etc/ipa/ca.crt`, you could also trust IPA CA on system wide.

```tf
provider "freeipa" {
  host     = var.freeipa_host
  username = var.freeipa_username
  password = var.freeipa_password
  base_dn  = var.freeipa_base_dn
}
```

Arguments:
* `host` - host name of IPA server, example *ipa.example.com* (string, required)
* `username` - username for authentication (string, required)
* `password` - password for authentication (string, required)
* `base_dn` - domain components (dc), example: *dc=ipa,dc=example,dc=com* (string, required)

### Resource User

```tf
resource "freeipa_user" "user" {
  uid        = var.uid
  first_name = var.first_name
  last_name  = var.last_name
  email      = var.email
  groups     = var.groups
  uid_number = var.uid_number
  gid_number = var.gid_number
}
```

Arguments:
* `uid` - User name (string, required)
* `first_name` - First name (string, required)
* `last_name` - Last name (string, required)
* `email` - Email address (string, optional, computed)
* `groups` - Groups user will be added to (set/list, optional, computed)
* `uid_number` - User ID (string, optional, computed)
* `gid_number` - Group ID (string, optional, computed)

### Datasouce User

```tf
data "freeipa_user" "user" {
  uid = var.uid
}
```

Arguments:
* `uid` - User name (string, required)

### Resource Group

```tf
resource "freeipa_user" "user" {
  gid           = var.gid
  gid_number    = var.gid_number
  description   = var.description
  groups        = var.groups
  group_members = var.group_members
}
```

Arguments:
* `gid` - Group Name (string, required)
* `gid_number` - Group ID (string, optional, computed)
* `description` - Description for group (string, optional, default: "")
* `groups` - Groups group will be added to (set/list, optional, computed)(attribute memberof)
* `group_members` - Group groups will be added to (set/list, optional, computed)(attribute member)

### Datasouce Group

```tf
data "freeipa_group" "group" {
  gid = var.gid
}
```

Arguments:
* `gid` - Group Name (string, required)

### Import

```bash
# terraform import <module path> <ipauniqueid for group>

terraform import freeipa_group.group <ipauniqueid for group>
terraform import freeipa_user.user <ipauniqueid for user>
```

If your using FreeIPA then your likely to find that the interface does not expose
the `ipauniqueid`. To get this you will need to log onto the IPA Server and run the
following command:

```bash
ipa user-show --all <username>
ipa group-show --all <groupname>
```


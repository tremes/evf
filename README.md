# evf
evf stands for Errata verion finder and this simple tool is for finding all the Bugzillas for the given product, component, version and status of bugs. It scans all the comments in every bug and tries to find related errata. For every errata, it tries to find the corresponding release version of the product. This was originally implemented to find z-stream release versions of the OpenShift Container platform (for particular Bugzilla bugs).

## Configuration

Configuration is defined in the `config.yaml` file in the root of this repository.**This file is required** and an example looks like the following (you can also check the `example-config.yaml`):

```yaml
bugzilla:
  url: https://bugzilla.redhat.com/rest
  key: "<your Bugzilla Key>"
  params:
    product: "OpenShift Container Platform"
    component: "Insights Operator"
    key: "CLOSED,VERIFIED"
    version: "4.9"
errata:
  url: "https://errata.devel.redhat.com/api/v1/erratum/"
  kerberos-conf: "/etc/krb5.conf"
  username: "<your Kerberos username>"
  realm: "<Kerberos Realm>"
  password: "<your Kerberos password>"

```

 You have to provide your Bugzilla key (to be able to communicate with the Bugzilla API) and your Kerberos settings (to be able to comunicate with the errata API).

## How to run

You can run the evf with:

```bash
make run
```

or with:

```bash
go run cmd/evf/main.go
```

The output is printed to the standard output.
# evf
evf stands for Errata verion finder and this simple tool is for finding all the Jiras for the given product, component, version and status of bugs. It scans all the comments in every bug and tries to find related errata. For every errata, it tries to find the corresponding release version of the product. This was originally implemented to find z-stream release versions of the OpenShift Container platform (for particular Bugzilla bugs).

## Configuration

Configuration is defined in the `config.yaml` file in the root of this repository.**This file is required** and an example looks like the following (you can also check the `example-config.yaml`):

```yaml
errata:
  url: "https://errata.devel.redhat.com/api/v1/erratum/"
  kerberos-conf: "/etc/krb5.conf"
  username: "<your Kerberos username>"
  realm: "<Kerberos Realm>"
jira:
  url: "https://issues.redhat.com"
  token: "<your Jira Access Token>"
  params:
    jql: "<jql issues search query>"
    # Example jql for all IO OCPBUGS closed in 4.12: 
    # project = "OCPBUGS" AND component="Insights Operator" AND status = Closed AND affectedversion = 4.12
```

 You have to provide your Jira token (to be able to communicate with the Jira API) and your Kerberos settings (to be able to comunicate with the errata API).

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
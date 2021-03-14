# iasst
iasst is a tool to help you take inventory your aws resources.

# How to use

```
NAME:
   iasst sg - Shows list of resources to which specified security group are attached

USAGE:
   iasst sg [command options] [arguments...]

OPTIONS:
   --security-group-id value, --id value  Sets security group id
   --check-security-group, -s             Shows security groups to which the security group is being used as rule. (default: false)
   --check-eni, -e                        Shows ENIs to which the security group is being attached. (default: false)
   --help, -h                             show help (default: false)
```
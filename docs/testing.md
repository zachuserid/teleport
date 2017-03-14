```
Feature                                            | 2.0 | 1.0 -> 2.0 |
---------------------------------------------------|-----|------------|-----
Adding Nodes via Valid Static Token                |  ✔  |            |
Adding Nodes via Valid Short-lived Tokens          |     |            |
Adding Nodes via Invalid Static Token Fails        |  ✔  |            |
Adding Nodes via Invalid Short-lived Tokens Fails  |     |            |
Revoking Node Invitation                           |  ✔  |            |
                                                   |     |            |
Static Labels                                      |  ✔  |            |
Dynamic Labels                                     |  ✔  |            |
                                                   |     |            |
Adding Trusted Cluster Valid Static Token          |  ✔  |            |
Adding Trusted Cluster Valid Short-lived Token     |  ✔  |            |
Adding Trusted Cluster Invalid Static Token        |  ✔  |            |
Adding Trusted Cluster Invalid Short-lived Token   |  ✔  |            |
Removing Trusted Cluster                           |  ✔  |            |
                                                   |     |            |
Namespaces                                         |     |            |
RBAC                                               |     |            |
                                                   |     |            |
Adding Users TOTP                                  |  ✔  |            |
Adding Users U2F                                   |     |            |
Deleting Users                                     |     |            |
Login TOTP                                         |  ✔  |            |
Login OIDC                                         |  ✔  |            |
Login OIDC (Google)                                |     |            |
Login U2F                                          |     |            |
                                                   |     |            |
Backend: etcd                                      |     |            |
Backend: dynamodb                                  |     |            |
Backend: boltdb                                    |     |            |
Backend: dir                                       |     |            |
                                                   |     |            |
tsh ssh <regular-node>                             |  ✔  |            |
tsh ssh <trusted-node>                             |  ✔  |            |
tsh join <regular-node>                            |  ✔  |            |
tsh join <trusted-node>                            |     |            |
tsh play <regular-node>                            |     |            |
tsh play <trusted-node>                            |     |            |
tsh scp <regular-node>                             |  ✔  |            |
tsh scp <trusted-node>                             |  ✔  |            |
tsh ssh -L <regular-node>                          |  ✔  |            |
tsh ssh -L <trusted-node>                          |  ✔  |            |
tsh ls                                             |  ✔  |            |
tsh clusters                                       |  ✔  |            |
```

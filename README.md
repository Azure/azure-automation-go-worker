The azure automation worker is mainly used to run script on Azure virtual machine. These script can be part of the update management solution or standalone to run automation tasks.

# Requirement 
Go 1.10

# Build
On Windows :
```cmd
build.cmd
```
On Linux :
```sh
make
```

After building, the `/bin` folder will contain 2 executable (one for the main worker and an other one for sandboxes).

# Worker configuration
A configuration which contains the following required key is required to run the hybrid worker.

```json
{
  "jrds_cert_path" : "",
  "jrds_key_path" : "",
  "jrds_base_uri" : "",

  "account_id" : "",
  "machine_id" : "",
  "hybrid_worker_group_name" : "",
  "worker_version" : "",
  "working_directory_path" : ""
}
```

# Executing the worker
To start the hybrid worker run :
```sh
./worker <path_to_your_configuration>
```

# Missing features
- Proxy support
- Python automation assets
- Signature validation
- Http client retry logic

# Contributing

This project welcomes contributions and suggestions.  Most contributions require you to agree to a
Contributor License Agreement (CLA) declaring that you have the right to, and actually do, grant us
the rights to use your contribution. For details, visit https://cla.microsoft.com.

When you submit a pull request, a CLA-bot will automatically determine whether you need to provide
a CLA and decorate the PR appropriately (e.g., label, comment). Simply follow the instructions
provided by the bot. You will only need to do this once across all repos using our CLA.

This project has adopted the [Microsoft Open Source Code of Conduct](https://opensource.microsoft.com/codeofconduct/).
For more information see the [Code of Conduct FAQ](https://opensource.microsoft.com/codeofconduct/faq/) or
contact [opencode@microsoft.com](mailto:opencode@microsoft.com) with any additional questions or comments.
